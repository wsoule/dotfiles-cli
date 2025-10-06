package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

type Snapshot struct {
	Timestamp   string         `json:"timestamp"`
	Description string         `json:"description"`
	Config      *config.Config `json:"config"`
	Metadata    map[string]string `json:"metadata"`
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "📸 Manage configuration snapshots",
	Long: `📸 Snapshot Management

Create timestamped backups of your configuration for easy rollback.
Snapshots are stored in ~/.dotfiles/snapshots/ with timestamps.

Examples:
  dotfiles snapshot create                       # Create snapshot with auto description
  dotfiles snapshot create -m "Before update"    # Create with custom message
  dotfiles snapshot list                         # List all snapshots
  dotfiles snapshot restore <timestamp>          # Restore from snapshot
  dotfiles snapshot delete <timestamp>           # Delete a snapshot
  dotfiles snapshot auto                         # Create before major operations`,
}

var snapshotCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snapshot",
	Long: `Create a timestamped snapshot of your current configuration.
Useful before making major changes to enable easy rollback.`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Create snapshots directory
		snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")
		if err := os.MkdirAll(snapshotsDir, 0755); err != nil {
			fmt.Printf("❌ Error creating snapshots directory: %v\n", err)
			os.Exit(1)
		}

		// Create snapshot
		timestamp := time.Now().Format("20060102-150405")
		if message == "" {
			message = "Manual snapshot"
		}

		snapshot := Snapshot{
			Timestamp:   timestamp,
			Description: message,
			Config:      cfg,
			Metadata: map[string]string{
				"created_by": "dotfiles snapshot",
				"platform":   "darwin",
			},
		}

		// Save snapshot
		snapshotPath := filepath.Join(snapshotsDir, timestamp+".json")
		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			fmt.Printf("❌ Error marshaling snapshot: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(snapshotPath, data, 0644); err != nil {
			fmt.Printf("❌ Error saving snapshot: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📸 Snapshot created successfully!")
		fmt.Println()
		fmt.Printf("   Timestamp: %s\n", timestamp)
		fmt.Printf("   Description: %s\n", message)
		fmt.Printf("   Location: snapshots/%s.json\n", timestamp)
		fmt.Println()
		fmt.Printf("   Packages: %d brews, %d casks, %d taps, %d stow\n",
			len(cfg.Brews), len(cfg.Casks), len(cfg.Taps), len(cfg.Stow))
		fmt.Println()
		fmt.Println("💡 To restore this snapshot:")
		fmt.Printf("   dotfiles snapshot restore %s\n", timestamp)
	},
}

var snapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")
		if _, err := os.Stat(snapshotsDir); os.IsNotExist(err) {
			fmt.Println("📸 No snapshots found")
			fmt.Println()
			fmt.Println("💡 Create your first snapshot:")
			fmt.Println("   dotfiles snapshot create -m 'Initial snapshot'")
			return
		}

		entries, err := os.ReadDir(snapshotsDir)
		if err != nil {
			fmt.Printf("❌ Error reading snapshots: %v\n", err)
			os.Exit(1)
		}

		snapshots := []Snapshot{}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				snapshotPath := filepath.Join(snapshotsDir, entry.Name())
				data, err := os.ReadFile(snapshotPath)
				if err != nil {
					continue
				}

				var snapshot Snapshot
				if err := json.Unmarshal(data, &snapshot); err != nil {
					continue
				}

				snapshots = append(snapshots, snapshot)
			}
		}

		if len(snapshots) == 0 {
			fmt.Println("📸 No snapshots found")
			return
		}

		// Sort by timestamp (newest first)
		sort.Slice(snapshots, func(i, j int) bool {
			return snapshots[i].Timestamp > snapshots[j].Timestamp
		})

		fmt.Printf("📸 Found %d snapshot(s):\n", len(snapshots))
		fmt.Println("=" + strings.Repeat("=", 25))
		fmt.Println()

		for i, snapshot := range snapshots {
			// Parse timestamp for display
			t, err := time.Parse("20060102-150405", snapshot.Timestamp)
			displayTime := snapshot.Timestamp
			if err == nil {
				displayTime = t.Format("Jan 02, 2006 at 3:04 PM")
			}

			fmt.Printf("%d. %s\n", i+1, snapshot.Timestamp)
			fmt.Printf("   Created: %s\n", displayTime)
			if snapshot.Description != "" {
				fmt.Printf("   Message: %s\n", snapshot.Description)
			}
			fmt.Printf("   Packages: %d brews, %d casks, %d taps, %d stow\n",
				len(snapshot.Config.Brews),
				len(snapshot.Config.Casks),
				len(snapshot.Config.Taps),
				len(snapshot.Config.Stow))
			fmt.Println()
		}

		fmt.Println("💡 To restore a snapshot:")
		fmt.Println("   dotfiles snapshot restore <timestamp>")
	},
}

var snapshotRestoreCmd = &cobra.Command{
	Use:   "restore <timestamp>",
	Short: "Restore configuration from a snapshot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		timestamp := args[0]
		noBackup, _ := cmd.Flags().GetBool("no-backup")

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		// Load snapshot
		snapshotPath := filepath.Join(home, ".dotfiles", "snapshots", timestamp+".json")
		data, err := os.ReadFile(snapshotPath)
		if err != nil {
			fmt.Printf("❌ Snapshot not found: %s\n", timestamp)
			fmt.Println()
			fmt.Println("💡 List available snapshots:")
			fmt.Println("   dotfiles snapshot list")
			os.Exit(1)
		}

		var snapshot Snapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			fmt.Printf("❌ Error reading snapshot: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📸 Restoring snapshot...")
		fmt.Println()
		fmt.Printf("   Timestamp: %s\n", snapshot.Timestamp)
		if snapshot.Description != "" {
			fmt.Printf("   Description: %s\n", snapshot.Description)
		}
		fmt.Println()

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Create backup of current config before restoring (unless --no-backup)
		if !noBackup {
			fmt.Println("💾 Creating backup of current configuration...")
			currentCfg, err := config.Load(configPath)
			if err == nil {
				backupSnapshot := Snapshot{
					Timestamp:   time.Now().Format("20060102-150405"),
					Description: "Auto-backup before restore",
					Config:      currentCfg,
					Metadata: map[string]string{
						"created_by":   "auto-backup",
						"restore_from": timestamp,
					},
				}

				snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")
				backupPath := filepath.Join(snapshotsDir, backupSnapshot.Timestamp+".json")
				backupData, _ := json.MarshalIndent(backupSnapshot, "", "  ")
				if err := os.WriteFile(backupPath, backupData, 0644); err == nil {
					fmt.Printf("   ✅ Backup created: %s\n", backupSnapshot.Timestamp)
				}
			}
			fmt.Println()
		}

		// Restore snapshot
		if err := snapshot.Config.Save(configPath); err != nil {
			fmt.Printf("❌ Error restoring configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Configuration restored successfully!")
		fmt.Println()
		fmt.Printf("   Restored: %d brews, %d casks, %d taps, %d stow\n",
			len(snapshot.Config.Brews),
			len(snapshot.Config.Casks),
			len(snapshot.Config.Taps),
			len(snapshot.Config.Stow))
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Println("   • View config: dotfiles list")
		fmt.Println("   • Install packages: dotfiles install")
		fmt.Println("   • Check differences: dotfiles diff")
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:   "delete <timestamp>",
	Short: "Delete a snapshot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		timestamp := args[0]

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		snapshotPath := filepath.Join(home, ".dotfiles", "snapshots", timestamp+".json")
		if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
			fmt.Printf("❌ Snapshot not found: %s\n", timestamp)
			os.Exit(1)
		}

		// Read snapshot to show info
		data, _ := os.ReadFile(snapshotPath)
		var snapshot Snapshot
		json.Unmarshal(data, &snapshot)

		fmt.Printf("🗑️  Deleting snapshot: %s\n", timestamp)
		if snapshot.Description != "" {
			fmt.Printf("   Description: %s\n", snapshot.Description)
		}

		if err := os.Remove(snapshotPath); err != nil {
			fmt.Printf("❌ Error deleting snapshot: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Snapshot deleted")
	},
}

var snapshotAutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Create automatic snapshot before operations",
	Long: `Create an automatic snapshot with timestamp.
This is meant to be called before major operations like update, install, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			os.Exit(1)
		}

		snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")
		os.MkdirAll(snapshotsDir, 0755)

		timestamp := time.Now().Format("20060102-150405")
		snapshot := Snapshot{
			Timestamp:   timestamp,
			Description: "Auto-snapshot before operation",
			Config:      cfg,
			Metadata: map[string]string{
				"created_by": "auto-snapshot",
			},
		}

		snapshotPath := filepath.Join(snapshotsDir, timestamp+".json")
		data, _ := json.MarshalIndent(snapshot, "", "  ")
		os.WriteFile(snapshotPath, data, 0644)

		fmt.Printf("📸 Auto-snapshot created: %s\n", timestamp)
	},
}

func init() {
	snapshotCreateCmd.Flags().StringP("message", "m", "", "Snapshot description")
	snapshotRestoreCmd.Flags().Bool("no-backup", false, "Don't create backup before restoring")

	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotRestoreCmd)
	snapshotCmd.AddCommand(snapshotDeleteCmd)
	snapshotCmd.AddCommand(snapshotAutoCmd)

	rootCmd.AddCommand(snapshotCmd)
}
