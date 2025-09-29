package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup <file>",
	Short: "Backup your configuration to a file",
	Long:  `Creates a backup copy of your config.json`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupPath := args[0]

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Load current config
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Save to backup file
		if err := cfg.Save(backupPath); err != nil {
			fmt.Printf("Error creating backup: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Configuration backed up to: %s\n", backupPath)

		total := len(cfg.Taps) + len(cfg.Brews) + len(cfg.Casks)
		fmt.Printf("  📊 Backed up %d packages\n", total)
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore <file>",
	Short: "Restore configuration from a backup file",
	Long:  `Restores your config.json from a backup file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupPath := args[0]

		// Check if backup file exists
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			fmt.Printf("Backup file not found: %s\n", backupPath)
			os.Exit(1)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Create backup of current config if it exists
		if !cmd.Flags().Changed("no-backup") {
			if _, err := os.Stat(configPath); err == nil {
				backupCurrent := configPath + ".bak"
				if err := copyFile(configPath, backupCurrent); err != nil {
					fmt.Printf("Warning: Could not backup current config: %v\n", err)
				} else {
					fmt.Printf("📋 Current config backed up to: %s\n", backupCurrent)
				}
			}
		}

		// Load backup config
		backupCfg, err := config.Load(backupPath)
		if err != nil {
			fmt.Printf("Error loading backup file: %v\n", err)
			os.Exit(1)
		}

		// Save as current config
		if err := backupCfg.Save(configPath); err != nil {
			fmt.Printf("Error restoring configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Configuration restored from: %s\n", backupPath)

		total := len(backupCfg.Taps) + len(backupCfg.Brews) + len(backupCfg.Casks)
		fmt.Printf("  📊 Restored %d packages\n", total)
		if len(backupCfg.Taps) > 0 {
			fmt.Printf("    📋 %d taps\n", len(backupCfg.Taps))
		}
		if len(backupCfg.Brews) > 0 {
			fmt.Printf("    🍺 %d brews\n", len(backupCfg.Brews))
		}
		if len(backupCfg.Casks) > 0 {
			fmt.Printf("    📦 %d casks\n", len(backupCfg.Casks))
		}
	},
}

func copyFile(src, dst string) error {
	sourceData, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, sourceData, 0644)
}

func init() {
	restoreCmd.Flags().Bool("no-backup", false, "Don't backup current config before restoring")

	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
}