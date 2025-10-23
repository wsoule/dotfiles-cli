package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:     "cleanup",
	GroupID: "package",
	Short:   "🧹 Clean up old Homebrew package versions and cache",
	Long: `🧹 Cleanup Homebrew

Remove old versions of installed packages and clear Homebrew cache to free up disk space.

What gets cleaned:
• Old versions of upgraded formulas
• Old versions of upgraded casks
• Homebrew download cache
• Symlinks to deleted formulas

Examples:
  dotfiles cleanup              # Clean up old versions and cache
  dotfiles cleanup --dry-run    # Show what would be cleaned
  dotfiles cleanup --cache-only # Only clear cache, keep old versions`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		cacheOnly, _ := cmd.Flags().GetBool("cache-only")

		// Check if brew is available
		if _, err := exec.LookPath("brew"); err != nil {
			fmt.Println("❌ Homebrew not found")
			os.Exit(1)
		}

		fmt.Println("🧹 Cleaning up Homebrew...")
		fmt.Println("=" + strings.Repeat("=", 25))
		fmt.Println()

		// Step 1: Show disk space before cleanup
		fmt.Println("📊 Checking disk usage...")
		getBrewCacheSize()
		fmt.Println()

		if !cacheOnly {
			// Step 2: Clean up old versions
			fmt.Println("🗑️  Removing old package versions...")
			if dryRun {
				fmt.Println("[DRY RUN] Would run: brew cleanup")
				showOldVersions()
			} else {
				cleanupCmd := exec.Command("brew", "cleanup")
				cleanupCmd.Stdout = os.Stdout
				cleanupCmd.Stderr = os.Stderr
				if err := cleanupCmd.Run(); err != nil {
					fmt.Printf("⚠️  Cleanup had some errors: %v\n", err)
				} else {
					fmt.Println("✅ Old versions removed")
				}
			}
			fmt.Println()
		}

		// Step 3: Clean up cache
		fmt.Println("🗑️  Clearing download cache...")
		if dryRun {
			fmt.Println("[DRY RUN] Would run: brew cleanup -s")
		} else {
			cacheCmd := exec.Command("brew", "cleanup", "-s")
			cacheCmd.Stdout = os.Stdout
			cacheCmd.Stderr = os.Stderr
			if err := cacheCmd.Run(); err != nil {
				fmt.Printf("⚠️  Cache cleanup had some errors: %v\n", err)
			} else {
				fmt.Println("✅ Cache cleared")
			}
		}
		fmt.Println()

		// Step 4: Clean up broken symlinks
		if !cacheOnly {
			fmt.Println("🔗 Cleaning up broken symlinks...")
			if dryRun {
				fmt.Println("[DRY RUN] Would run: brew cleanup --prune=all")
			} else {
				pruneCmd := exec.Command("brew", "cleanup", "--prune=all")
				pruneCmd.Stdout = os.Stdout
				pruneCmd.Stderr = os.Stderr
				if err := pruneCmd.Run(); err != nil {
					fmt.Printf("⚠️  Prune had some errors: %v\n", err)
				} else {
					fmt.Println("✅ Symlinks cleaned")
				}
			}
			fmt.Println()
		}

		// Step 5: Show disk space after cleanup
		if !dryRun {
			fmt.Println("📊 Final disk usage...")
			getBrewCacheSize()
			fmt.Println()
		}

		fmt.Println("🎉 Cleanup complete!")
		fmt.Println()
		fmt.Println("💡 Pro tip:")
		fmt.Println("   • Run this periodically to save disk space")
		fmt.Println("   • Run 'brew cleanup -n' to preview what would be removed")
	},
}

func getBrewCacheSize() {
	home, _ := os.UserHomeDir()
	cachePath := home + "/Library/Caches/Homebrew"

	cmd := exec.Command("du", "-sh", cachePath)
	output, err := cmd.Output()
	if err == nil {
		fmt.Printf("   Cache size: %s\n", strings.TrimSpace(string(output)))
	}

	// Get Cellar size
	cellarCmd := exec.Command("du", "-sh", "/opt/homebrew/Cellar")
	cellarOutput, err := cellarCmd.Output()
	if err == nil {
		fmt.Printf("   Cellar size: %s\n", strings.TrimSpace(string(cellarOutput)))
	}
}

func showOldVersions() {
	cmd := exec.Command("brew", "cleanup", "-n")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func init() {
	cleanupCmd.Flags().Bool("dry-run", false, "Show what would be cleaned without actually cleaning")
	cleanupCmd.Flags().Bool("cache-only", false, "Only clear cache, don't remove old versions")

	rootCmd.AddCommand(cleanupCmd)
}
