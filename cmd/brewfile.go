package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var brewfileCmd = &cobra.Command{
	Use:   "brewfile",
	Short: "Generate a Brewfile from your configuration",
	Long:  `Creates a Brewfile based on your config.json packages`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		brewfileContent := cfg.GenerateBrewfile()

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "./Brewfile"
		}

		if err := os.WriteFile(output, []byte(brewfileContent), 0644); err != nil {
			fmt.Printf("Error writing Brewfile: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Generated Brewfile at: %s\n", output)
	},
}

func init() {
	brewfileCmd.Flags().StringP("output", "o", "./Brewfile", "Output path for the Brewfile")
	rootCmd.AddCommand(brewfileCmd)
}