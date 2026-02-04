package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/moganji/jot/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	store      storage.Store
	storePath  string
	formatFlag string
)

var rootCmd = &cobra.Command{
	Use:   "jot",
	Short: "Agent-friendly local memory",
	Long:  `Store and search notes across conversations. Fast, local, and persistent.`,
	Example: `  jot add api-creds --title "Stripe Key" --content "sk_test_..."
  jot search "stripe"
  jot list api-creds
  jot get api-creds "stripe"`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initStore)
	rootCmd.PersistentFlags().StringVar(&storePath, "store", getDefaultStorePath(), "path to notes directory")
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "text", "output format (text|json)")
}

func initStore() {
	var err error
	store, err = storage.NewFileStore(storePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize store: %v\n", err)
		os.Exit(1)
	}
}

func getDefaultStorePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".jot"
	}
	return filepath.Join(home, ".jot")
}
