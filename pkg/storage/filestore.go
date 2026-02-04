package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/moganji/jot/pkg/models"
	_ "modernc.org/sqlite"
	"gopkg.in/yaml.v3"
)

// FileStore implements Store using markdown files + SQLite FTS5
type FileStore struct {
	basePath string
	searchDB *sql.DB
}

// Note metadata for YAML frontmatter
type NoteMeta struct {
	ID       string    `yaml:"id"`
	Title    string    `yaml:"title"`
	Created  time.Time `yaml:"created"`
	Updated  time.Time `yaml:"updated"`
	Tags     []string  `yaml:"tags,omitempty"`
	Category string    `yaml:"category"`
}

func NewFileStore(basePath string) (*FileStore, error) {
	// Create base directory
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// Create index directory
	indexDir := filepath.Join(basePath, ".index")
	if err := os.MkdirAll(indexDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create index directory: %w", err)
	}

	// Open SQLite database for FTS5
	dbPath := filepath.Join(indexDir, "search.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open search database: %w", err)
	}

	store := &FileStore{
		basePath: basePath,
		searchDB: db,
	}

	// Initialize FTS5 index
	if err := store.initSearchIndex(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *FileStore) initSearchIndex() error {
	schema := `
	CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
		id UNINDEXED,
		title,
		content,
		tags,
		category UNINDEXED,
		filepath UNINDEXED
	);
	`
	_, err := s.searchDB.Exec(schema)
	return err
}

func (s *FileStore) Add(note *models.Note) error {
	// Create category directory
	categoryPath := filepath.Join(s.basePath, note.Category)
	if err := os.MkdirAll(categoryPath, 0755); err != nil {
		return fmt.Errorf("failed to create category directory: %w", err)
	}

	// Generate filename from title (slugify)
	filename := slugify(note.Title) + ".md"
	filePath := filepath.Join(categoryPath, filename)

	// Format as markdown with YAML frontmatter
	content, err := s.formatMarkdown(note)
	if err != nil {
		return fmt.Errorf("failed to format markdown: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Update search index
	relPath, _ := filepath.Rel(s.basePath, filePath)
	_, err = s.searchDB.Exec(`
		INSERT INTO notes_fts (id, title, content, tags, category, filepath)
		VALUES (?, ?, ?, ?, ?, ?)
	`, note.ID, note.Title, note.Content, strings.Join(note.Tags, " "), note.Category, relPath)

	return err
}

func (s *FileStore) Get(id string) (*models.Note, error) {
	// Search index for file path
	var filePath string
	err := s.searchDB.QueryRow(`
		SELECT filepath FROM notes_fts WHERE id = ?
	`, id).Scan(&filePath)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("note not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	// Read and parse markdown file
	fullPath := filepath.Join(s.basePath, filePath)
	return s.parseMarkdownFile(fullPath)
}

func (s *FileStore) GetByTitle(category, title string) (*models.Note, error) {
	// Try exact match first
	filename := slugify(title) + ".md"
	exactPath := filepath.Join(s.basePath, category, filename)

	if _, err := os.Stat(exactPath); err == nil {
		return s.parseMarkdownFile(exactPath)
	}

	// Fall back to search
	var filePath string
	err := s.searchDB.QueryRow(`
		SELECT filepath FROM notes_fts
		WHERE category = ? AND title = ?
		LIMIT 1
	`, category, title).Scan(&filePath)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("note not found: %s/%s", category, title)
	}
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(s.basePath, filePath)
	return s.parseMarkdownFile(fullPath)
}

func (s *FileStore) List(category string) ([]*models.Note, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT filepath FROM notes_fts WHERE category = ? ORDER BY filepath`
		args = []interface{}{category}
	} else {
		query = `SELECT filepath FROM notes_fts ORDER BY filepath`
	}

	rows, err := s.searchDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			continue
		}

		fullPath := filepath.Join(s.basePath, filePath)
		note, err := s.parseMarkdownFile(fullPath)
		if err != nil {
			continue
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *FileStore) Update(note *models.Note) error {
	// Delete old entry
	var oldPath string
	err := s.searchDB.QueryRow(`SELECT filepath FROM notes_fts WHERE id = ?`, note.ID).Scan(&oldPath)
	if err != nil {
		return fmt.Errorf("note not found: %s", note.ID)
	}

	fullOldPath := filepath.Join(s.basePath, oldPath)

	// Delete from index
	_, err = s.searchDB.Exec(`DELETE FROM notes_fts WHERE id = ?`, note.ID)
	if err != nil {
		return err
	}

	// Delete old file
	os.Remove(fullOldPath)

	// Add as new (handles category changes)
	return s.Add(note)
}

func (s *FileStore) Delete(id string) error {
	// Get file path
	var filePath string
	err := s.searchDB.QueryRow(`SELECT filepath FROM notes_fts WHERE id = ?`, id).Scan(&filePath)
	if err == sql.ErrNoRows {
		return fmt.Errorf("note not found: %s", id)
	}
	if err != nil {
		return err
	}

	// Delete from index
	_, err = s.searchDB.Exec(`DELETE FROM notes_fts WHERE id = ?`, id)
	if err != nil {
		return err
	}

	// Delete file
	fullPath := filepath.Join(s.basePath, filePath)
	return os.Remove(fullPath)
}

func (s *FileStore) Search(query string, category string, tags []string) ([]*models.Note, error) {
	// Build FTS5 query
	sqlQuery := `SELECT filepath, rank FROM notes_fts WHERE notes_fts MATCH ?`
	args := []interface{}{query}

	if category != "" {
		sqlQuery += ` AND category = ?`
		args = append(args, category)
	}

	sqlQuery += ` ORDER BY rank LIMIT 100`

	rows, err := s.searchDB.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var filePath string
		var rank float64
		if err := rows.Scan(&filePath, &rank); err != nil {
			continue
		}

		fullPath := filepath.Join(s.basePath, filePath)
		note, err := s.parseMarkdownFile(fullPath)
		if err != nil {
			continue
		}

		// Filter by tags if specified
		if len(tags) > 0 && !hasAnyTag(note.Tags, tags) {
			continue
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (s *FileStore) GetCategories() ([]string, error) {
	rows, err := s.searchDB.Query(`SELECT DISTINCT category FROM notes_fts ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (s *FileStore) GetTags() ([]string, error) {
	// Get all unique tags from index
	rows, err := s.searchDB.Query(`SELECT DISTINCT tags FROM notes_fts WHERE tags != ''`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagSet := make(map[string]bool)
	for rows.Next() {
		var tagsStr string
		if err := rows.Scan(&tagsStr); err != nil {
			continue
		}
		for _, tag := range strings.Fields(tagsStr) {
			tagSet[tag] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *FileStore) Close() error {
	return s.searchDB.Close()
}

// Helper functions

func (s *FileStore) formatMarkdown(note *models.Note) (string, error) {
	meta := NoteMeta{
		ID:       note.ID,
		Title:    note.Title,
		Created:  note.Created,
		Updated:  note.Updated,
		Tags:     note.Tags,
		Category: note.Category,
	}

	yamlBytes, err := yaml.Marshal(meta)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("---\n%s---\n\n%s\n", string(yamlBytes), note.Content), nil
}

func (s *FileStore) parseMarkdownFile(path string) (*models.Note, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)

	// Parse YAML frontmatter
	if !strings.HasPrefix(content, "---\n") {
		return nil, fmt.Errorf("invalid markdown format: missing frontmatter")
	}

	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid markdown format: malformed frontmatter")
	}

	var meta NoteMeta
	if err := yaml.Unmarshal([]byte(parts[0]), &meta); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	note := &models.Note{
		ID:       meta.ID,
		Title:    meta.Title,
		Content:  strings.TrimSpace(parts[1]),
		Tags:     meta.Tags,
		Created:  meta.Created,
		Updated:  meta.Updated,
		Category: meta.Category,
		Metadata: make(map[string]string),
	}

	return note, nil
}

func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and special chars with hyphens
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '_' {
			return '-'
		}
		return -1
	}, s)

	// Remove consecutive hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	// Trim hyphens from ends
	s = strings.Trim(s, "-")

	// Limit length
	if len(s) > 100 {
		s = s[:100]
	}

	return s
}

func hasAnyTag(noteTags []string, searchTags []string) bool {
	for _, searchTag := range searchTags {
		for _, noteTag := range noteTags {
			if strings.EqualFold(noteTag, searchTag) {
				return true
			}
		}
	}
	return false
}
