package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List all categories",
	RunE:  runCategories,
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags",
	RunE:  runTags,
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(tagsCmd)
}

func runCategories(cmd *cobra.Command, args []string) error {
	categories, err := store.GetCategories()
	if err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}

	if len(categories) == 0 {
		fmt.Println("No categories found")
		return nil
	}

	if formatFlag == "json" {
		return outputJSON(categories)
	}

	sort.Strings(categories)

	fmt.Println("Categories:")
	for _, cat := range categories {
		notes, _ := store.List(cat)
		fmt.Printf("  %s (%d note(s))\n", cat, len(notes))
	}

	return nil
}

func runTags(cmd *cobra.Command, args []string) error {
	tags, err := store.GetTags()
	if err != nil {
		return fmt.Errorf("failed to get tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found")
		return nil
	}

	if formatFlag == "json" {
		return outputJSON(tags)
	}

	sort.Strings(tags)

	fmt.Println("Tags:")
	for _, tag := range tags {
		fmt.Printf("  %s\n", tag)
	}

	return nil
}
