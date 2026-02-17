package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var brewfileCmd = &cobra.Command{
	Use: "brewfile",
	GroupID: "advanced",
	Short: "Generate a Brewfile from your configuration",
	Long:    `Creates a Brewfile based on your config.json packages`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("error loading configuration: %w", err)
		}

		brewfileContent := cfg.GenerateBrewfile()

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "./Brewfile"
		}

		if err := os.WriteFile(output, []byte(brewfileContent), 0644); err != nil {
			return fmt.Errorf("error writing brewfile: %w", err)
		}

		fmt.Printf(" Generated Brewfile at: %s\n", output)
		return nil
	},
}

func init() {
	brewfileCmd.Flags().StringP("output", "o", "./Brewfile", "Output path for the Brewfile")
	rootCmd.AddCommand(brewfileCmd)
}
