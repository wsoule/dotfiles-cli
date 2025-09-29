package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new dotfiles configuration",
	Long:  `Creates a new config.json file in ~/.dotfiles/`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configDir := filepath.Join(home, ".dotfiles")
		configPath := filepath.Join(configDir, "config.json")

		// Check if config already exists
		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("Configuration already exists at:", configPath)
			return
		}

		// Create initial empty config
		cfg := &config.Config{
			Brews: []string{"git"},
			Casks: []string{},
			Taps:  []string{},
		}

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("Error creating configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Created configuration at:", configPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}