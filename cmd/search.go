package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/moganji/jot/pkg/models"
	"github.com/spf13/cobra"
)

var (
	searchCategory string
	searchTags     string
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search notes",
	Example: `  jot search "stripe"
  jot search "oauth" --in api-quirks
  jot search "api" --tag payment,sandbox`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchCategory, "in", "", "search only in this category")
	searchCmd.Flags().StringVar(&searchTags, "tag", "", "filter by tags (comma-separated)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	var tags []string
	if searchTags != "" {
		tags = strings.Split(searchTags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	results, err := store.Search(query, searchCategory, tags)
	if err != nil {
		return fmt.Errorf("failed to search: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No notes found")
		return nil
	}

	scoredResults := rankResults(results, query)
	sort.Slice(scoredResults, func(i, j int) bool {
		return scoredResults[i].score > scoredResults[j].score
	})

	if formatFlag == "json" {
		notes := make([]*models.Note, len(scoredResults))
		for i, r := range scoredResults {
			notes[i] = r.note
		}
		return outputJSON(notes)
	}

	fmt.Printf("Found %d note(s):\n\n", len(scoredResults))

	for _, result := range scoredResults {
		note := result.note
		fmt.Printf("  [%s] %s (%s)\n", note.Category, note.Title, note.ID[:8])

		preview := getMatchPreview(note.Content, query)
		if preview != "" {
			fmt.Printf("  > %s\n", preview)
		}
		fmt.Println()
	}

	return nil
}

type scoredResult struct {
	note  *models.Note
	score int
}

func rankResults(notes []*models.Note, query string) []scoredResult {
	queryLower := strings.ToLower(query)
	results := make([]scoredResult, len(notes))

	for i, note := range notes {
		score := 0
		titleLower := strings.ToLower(note.Title)
		contentLower := strings.ToLower(note.Content)

		if titleLower == queryLower {
			score += 100
		} else if strings.Contains(titleLower, queryLower) {
			score += 50
		}

		if strings.Contains(contentLower, queryLower) {
			if strings.HasPrefix(contentLower, queryLower) {
				score += 30
			} else {
				score += 10
			}
		}

		results[i] = scoredResult{note: note, score: score}
	}

	return results
}

func getMatchPreview(content, query string) string {
	queryLower := strings.ToLower(query)
	contentLower := strings.ToLower(content)

	idx := strings.Index(contentLower, queryLower)
	if idx == -1 {
		if len(content) > 80 {
			return content[:80] + "..."
		}
		return content
	}

	start := idx - 20
	if start < 0 {
		start = 0
	}

	end := idx + len(query) + 40
	if end > len(content) {
		end = len(content)
	}

	preview := content[start:end]
	preview = strings.ReplaceAll(preview, "\n", " ")

	if start > 0 {
		preview = "..." + preview
	}
	if end < len(content) {
		preview = preview + "..."
	}

	return preview
}
