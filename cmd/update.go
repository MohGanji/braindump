package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/MohGanji/braindump/pkg/models"
	"github.com/spf13/cobra"
)

var (
	updateTitle   string
	updateContent string
	updateTags    string
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a note",
	Example: `  braindump update a1b2c3d4 --content "new content"
  braindump update a1b2c3d4 --title "New Title" --tags "tag1,tag2"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

var appendCmd = &cobra.Command{
	Use:   "append <id> <content>",
	Short: "Append content to a note",
	Example: `  braindump append a1b2c3d4 "Additional information"`,
	Args: cobra.ExactArgs(2),
	RunE: runAppend,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(appendCmd)

	updateCmd.Flags().StringVar(&updateTitle, "title", "", "new title")
	updateCmd.Flags().StringVar(&updateContent, "content", "", "new content")
	updateCmd.Flags().StringVar(&updateTags, "tags", "", "comma-separated tags")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	idOrTitle := args[0]

	if updateTitle == "" && updateContent == "" && updateTags == "" {
		return fmt.Errorf("at least one of --title, --content, or --tags must be provided")
	}

	note, err := findNote(idOrTitle)
	if err != nil {
		return err
	}

	if updateTitle != "" {
		note.Title = updateTitle
	}

	if updateContent != "" {
		note.Content = updateContent
	}

	if updateTags != "" {
		tags := strings.Split(updateTags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		note.Tags = tags
	}

	note.Updated = time.Now()

	if err := store.Update(note); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	fmt.Printf("✓ Updated note: \"%s\" (id: %s)\n", note.Title, note.ID[:8])
	return nil
}

func runAppend(cmd *cobra.Command, args []string) error {
	idOrTitle := args[0]
	appendContent := args[1]

	note, err := findNote(idOrTitle)
	if err != nil {
		return err
	}

	note.Content = note.Content + "\n" + appendContent
	note.Updated = time.Now()

	if err := store.Update(note); err != nil {
		return fmt.Errorf("failed to append to note: %w", err)
	}

	fmt.Printf("✓ Appended to note: \"%s\" (id: %s)\n", note.Title, note.ID[:8])
	return nil
}

func findNote(idOrTitle string) (*models.Note, error) {
	note, err := store.Get(idOrTitle)
	if err == nil {
		return note, nil
	}

	notes, err := store.List("")
	if err != nil {
		return nil, fmt.Errorf("failed to search for note: %w", err)
	}

	var idMatches []*models.Note
	var titleMatches []*models.Note

	for _, n := range notes {
		if strings.HasPrefix(n.ID, idOrTitle) {
			idMatches = append(idMatches, n)
		}
		if n.Title == idOrTitle {
			titleMatches = append(titleMatches, n)
		}
	}

	if len(idMatches) == 1 {
		return idMatches[0], nil
	}

	if len(idMatches) > 1 {
		fmt.Printf("Multiple notes found with ID prefix \"%s\":\n", idOrTitle)
		for _, n := range idMatches {
			fmt.Printf("  [%s] %s (id: %s)\n", n.Category, n.Title, n.ID[:8])
		}
		return nil, fmt.Errorf("please specify a longer ID prefix")
	}

	if len(titleMatches) == 0 {
		return nil, fmt.Errorf("note not found: %s", idOrTitle)
	}

	if len(titleMatches) == 1 {
		return titleMatches[0], nil
	}

	fmt.Printf("Multiple notes found with title \"%s\":\n", idOrTitle)
	for _, n := range titleMatches {
		fmt.Printf("  [%s] %s (id: %s)\n", n.Category, n.Title, n.ID[:8])
	}
	return nil, fmt.Errorf("please specify by ID")
}
