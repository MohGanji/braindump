package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [category]",
	Short: "List notes",
	Example: `  braindump list
  braindump list api-creds`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	var category string
	if len(args) > 0 {
		category = args[0]
	}

	notes, err := store.List(category)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found")
		return nil
	}

	if formatFlag == "json" {
		return outputJSON(notes)
	}

	sort.Slice(notes, func(i, j int) bool {
		if notes[i].Category == notes[j].Category {
			return notes[i].Created.Before(notes[j].Created)
		}
		return notes[i].Category < notes[j].Category
	})

	currentCategory := ""
	for _, note := range notes {
		if note.Category != currentCategory {
			if currentCategory != "" {
				fmt.Println()
			}
			currentCategory = note.Category
			fmt.Printf("[%s]\n", currentCategory)
		}

		preview := note.Content
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		preview = strings.ReplaceAll(preview, "\n", " ")

		fmt.Printf("  %s - %s\n", note.Title, preview)
		fmt.Printf("    ID: %s | Created: %s\n", note.ID[:8], note.Created.Format("2006-01-02 15:04"))
	}

	fmt.Printf("\nTotal: %d note(s)\n", len(notes))
	return nil
}
