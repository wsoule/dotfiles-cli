package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:     "backup [file]",
	GroupID: "advanced",
	Short:   "💾 Backup your configuration to a file",
	Long: `💾 Create Configuration Backup

Creates a timestamped backup of your config.json.
If no file path is provided, uses automatic timestamped naming.

Examples:
  dotfiles backup                        # Auto-timestamped: config-backup-20230615-143022.json
  dotfiles backup mybackup.json          # Custom name
  dotfiles backup --dir ~/backups        # Custom directory`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupDir, _ := cmd.Flags().GetString("dir")

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := GetConfigPath(cmd)

		// Load current config
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Determine backup path
		var backupPath string
		if len(args) > 0 {
			backupPath = args[0]
		} else {
			// Auto-generate timestamped backup name
			timestamp := time.Now().Format("20060102-150405")
			filename := fmt.Sprintf("config-backup-%s.json", timestamp)

			if backupDir != "" {
				backupPath = filepath.Join(backupDir, filename)
			} else {
				backupPath = filepath.Join(home, ".dotfiles", "backups", filename)
			}
		}

		// Ensure directory exists
		backupDirPath := filepath.Dir(backupPath)
		if err := os.MkdirAll(backupDirPath, 0755); err != nil {
			fmt.Printf("❌ Error creating backup directory: %v\n", err)
			os.Exit(1)
		}

		// Save to backup file
		if err := cfg.Save(backupPath); err != nil {
			fmt.Printf("❌ Error creating backup: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Configuration backed up successfully!")
		fmt.Println()
		fmt.Printf("   Location: %s\n", backupPath)
		fmt.Println()

		total := len(cfg.Taps) + len(cfg.Brews) + len(cfg.Casks) + len(cfg.Stow)
		fmt.Printf("   📊 Backed up %d items:\n", total)
		if len(cfg.Taps) > 0 {
			fmt.Printf("      • %d taps\n", len(cfg.Taps))
		}
		if len(cfg.Brews) > 0 {
			fmt.Printf("      • %d brews\n", len(cfg.Brews))
		}
		if len(cfg.Casks) > 0 {
			fmt.Printf("      • %d casks\n", len(cfg.Casks))
		}
		if len(cfg.Stow) > 0 {
			fmt.Printf("      • %d stow packages\n", len(cfg.Stow))
		}
		fmt.Println()
		fmt.Println("💡 To restore this backup:")
		fmt.Printf("   dotfiles restore %s\n", backupPath)
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

		// Create timestamped backup of current config if it exists
		if !cmd.Flags().Changed("no-backup") {
			if _, err := os.Stat(configPath); err == nil {
				timestamp := time.Now().Format("20060102-150405")
				backupsDir := filepath.Join(home, ".dotfiles", "backups")
				os.MkdirAll(backupsDir, 0755)
				backupCurrent := filepath.Join(backupsDir, fmt.Sprintf("config-pre-restore-%s.json", timestamp))

				if err := copyFile(configPath, backupCurrent); err != nil {
					fmt.Printf("⚠️  Warning: Could not backup current config: %v\n", err)
				} else {
					fmt.Printf("💾 Current config backed up to:\n   %s\n", backupCurrent)
					fmt.Println()
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
	backupCmd.Flags().String("dir", "", "Custom backup directory")
	restoreCmd.Flags().Bool("no-backup", false, "Don't backup current config before restoring")

	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
}