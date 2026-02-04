package cmd

import (
	"fmt"
	"strings"

	"github.com/moganji/jot/pkg/models"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a note",
	Example: `  jot delete a1b2c3d4`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	idOrTitle := args[0]

	note, err := store.Get(idOrTitle)
	if err == nil {
		if err := store.Delete(note.ID); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}
		fmt.Printf("✓ Deleted note: \"%s\" (id: %s)\n", note.Title, note.ID[:8])
		return nil
	}

	notes, err := store.List("")
	if err != nil {
		return fmt.Errorf("failed to search for note: %w", err)
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
		if err := store.Delete(idMatches[0].ID); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}
		fmt.Printf("✓ Deleted note: \"%s\" (id: %s)\n", idMatches[0].Title, idMatches[0].ID[:8])
		return nil
	}

	if len(idMatches) > 1 {
		fmt.Printf("Multiple notes found with ID prefix \"%s\":\n", idOrTitle)
		for _, n := range idMatches {
			fmt.Printf("  [%s] %s (id: %s)\n", n.Category, n.Title, n.ID[:8])
		}
		return fmt.Errorf("please specify a longer ID prefix")
	}

	if len(titleMatches) == 0 {
		return fmt.Errorf("note not found: %s", idOrTitle)
	}

	if len(titleMatches) == 1 {
		if err := store.Delete(titleMatches[0].ID); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}
		fmt.Printf("✓ Deleted note: \"%s\" (id: %s)\n", titleMatches[0].Title, titleMatches[0].ID[:8])
		return nil
	}

	fmt.Printf("Multiple notes found with title \"%s\":\n", idOrTitle)
	for _, n := range titleMatches {
		fmt.Printf("  [%s] %s (id: %s)\n", n.Category, n.Title, n.ID[:8])
	}
	fmt.Println("\nPlease delete by specific ID")

	return nil
}
