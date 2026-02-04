package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/moganji/jot/pkg/models"
	"github.com/spf13/cobra"
)

var (
	addTitle   string
	addContent string
	addTags    string
)

var addCmd = &cobra.Command{
	Use:   "add <category> [title] [content]",
	Short: "Add a new note",
	Example: `  jot add api-creds --title "Stripe Key" --content "sk_test_..."
  jot add api-creds "Stripe Key" "sk_test_..."
  echo "sk_test_..." | jot add api-creds --title "Stripe Key"
  jot add api-creds --title "Stripe" --content "..." --tags "stripe,payment"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addTitle, "title", "", "note title")
	addCmd.Flags().StringVar(&addContent, "content", "", "note content")
	addCmd.Flags().StringVar(&addTags, "tags", "", "comma-separated tags")
}

func runAdd(cmd *cobra.Command, args []string) error {
	category := args[0]

	title := addTitle
	content := addContent

	if title == "" && len(args) > 1 {
		title = args[1]
	}

	if content == "" && len(args) > 2 {
		content = args[2]
	}

	if content == "" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = strings.TrimSpace(string(bytes))
		}
	}

	if title == "" {
		return fmt.Errorf("title is required (use --title flag or provide as argument)")
	}

	if content == "" {
		return fmt.Errorf("content is required (use --content flag, provide as argument, or pipe via stdin)")
	}

	var tags []string
	if addTags != "" {
		tags = strings.Split(addTags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	note := models.NewNote(category, title, content, tags)

	if err := store.Add(note); err != nil {
		return fmt.Errorf("failed to add note: %w", err)
	}

	if formatFlag == "json" {
		return outputJSON(note)
	}

	fmt.Printf("✓ Added note to %s: \"%s\" (id: %s)\n", category, title, note.ID[:8])
	return nil
}

func outputJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
