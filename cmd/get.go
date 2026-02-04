package cmd

import (
	"fmt"
	"strings"

	"github.com/moganji/jot/pkg/models"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <category> [pattern]",
	Short: "Get note(s) from a category",
	Example: `  jot get api-creds
  jot get api-creds "stripe"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	category := args[0]
	var titlePattern string
	if len(args) > 1 {
		titlePattern = args[1]
	}

	notes, err := store.List(category)
	if err != nil {
		return fmt.Errorf("failed to get notes: %w", err)
	}

	if titlePattern != "" {
		notes = filterByTitle(notes, titlePattern)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found")
		return nil
	}

	if formatFlag == "json" {
		return outputJSON(notes)
	}

	for i, note := range notes {
		if i > 0 {
			fmt.Println()
		}
		printNote(note)
	}

	return nil
}

func filterByTitle(notes []*models.Note, pattern string) []*models.Note {
	patternLower := strings.ToLower(pattern)
	var filtered []*models.Note
	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Title), patternLower) {
			filtered = append(filtered, note)
		}
	}
	return filtered
}

func printNote(note *models.Note) {
	fmt.Printf("%s (%s)\n", note.Title, note.ID[:8])
	fmt.Println(strings.Repeat("-", len(note.Title)+11))
	fmt.Println(note.Content)
	fmt.Println()
	fmt.Printf("Created: %s\n", note.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("Category: %s\n", note.Category)
	if len(note.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(note.Tags, ", "))
	}
}
