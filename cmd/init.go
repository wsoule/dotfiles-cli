package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	GroupID: "dotfiles",
	Short:   "Initialize a new dotfiles configuration",
	Long:    `Creates a new config.json file in ~/.dotfiles/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}

		configDir := filepath.Join(home, ".dotfiles")
		configPath := filepath.Join(configDir, "config.json")

		// Check if config already exists
		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("Configuration already exists at:", configPath)
			return nil
		}

		// Create initial empty config
		cfg := &config.Config{
			Brews: []string{"git"},
			Casks: []string{},
			Taps:  []string{},
		}

		if err := cfg.Save(configPath); err != nil {
			return fmt.Errorf("error creating configuration: %w", err)
		}

		fmt.Println("✓ Created configuration at:", configPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
