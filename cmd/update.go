package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use: "update",
	GroupID: "package",
	Short: "Update Homebrew packages to latest versions",
	Long: `  Update Homebrew Packages

Update packages managed by Homebrew to their latest versions.
You can update all packages or specific ones.

Examples:
  dotfiles update              # Update Homebrew and all packages
  dotfiles update --brew-only  # Only update Homebrew itself
  dotfiles update --dry-run    # Preview what would be updated
  dotfiles update git curl     # Update specific packages only`,
	RunE: func(cmd *cobra.Command, args []string) error {
		brewOnly, _ := cmd.Flags().GetBool("brew-only")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		skipBrewUpdate, _ := cmd.Flags().GetBool("skip-brew-update")
		noSnapshot, _ := cmd.Flags().GetBool("no-snapshot")

		// Create snapshot before update (unless --no-snapshot or --dry-run)
		if !dryRun && !noSnapshot {
			snapshotCmd := exec.Command("dotfiles", "snapshot", "auto")
			snapshotCmd.Run()
		}

		// Check if brew is available
		if _, err := exec.LookPath("brew"); err != nil {
			return fmt.Errorf("Homebrew not found. Install Homebrew first")
		}

		fmt.Println("Updating Packages...")
		fmt.Println("=" + strings.Repeat("=", 22))
		fmt.Println()

		// Step 1: Update Homebrew itself
		if !skipBrewUpdate {
			fmt.Println(" Updating Homebrew...")
			if dryRun {
				fmt.Println("[DRY RUN] Would run: brew update")
			} else {
				updateCmd := exec.Command("brew", "update")
				updateCmd.Stdout = os.Stdout
				updateCmd.Stderr = os.Stderr
				if err := updateCmd.Run(); err != nil {
					return fmt.Errorf("warning: brew update failed: %w", err)
				} else {
					fmt.Println(" Homebrew updated")
				}
			}
			fmt.Println()
		}

		if brewOnly {
			fmt.Println(" Homebrew update complete!")
			return nil
		}

		// Step 2: Show outdated packages
		fmt.Println(" Checking for outdated packages...")
		outdated := getOutdatedPackages()

		if len(outdated) == 0 {
			fmt.Println(" All packages are up to date!")
			return nil
		}

		fmt.Printf("Found %d outdated package(s):\n", len(outdated))
		for _, pkg := range outdated {
			fmt.Printf("• %s\n", pkg)
		}
		fmt.Println()

		// Step 3: Upgrade packages
		if len(args) > 0 {
			// Upgrade specific packages
			fmt.Printf("Upgrading specific packages: %s\n", strings.Join(args, ", "))
			if dryRun {
				fmt.Printf("[DRY RUN] Would run: brew upgrade %s\n", strings.Join(args, " "))
			} else {
				upgradeCmd := exec.Command("brew", append([]string{"upgrade"}, args...)...)
				upgradeCmd.Stdout = os.Stdout
				upgradeCmd.Stderr = os.Stderr
				if err := upgradeCmd.Run(); err != nil {
					return fmt.Errorf(" failed to upgrade packages: %w", err)
				}
				fmt.Println(" Packages upgraded")
			}
		} else {
			// Upgrade all outdated packages
			fmt.Println("Upgrading all outdated packages...")
			if dryRun {
				fmt.Println("[DRY RUN] Would run: brew upgrade")
			} else {
				upgradeCmd := exec.Command("brew", "upgrade")
				upgradeCmd.Stdout = os.Stdout
				upgradeCmd.Stderr = os.Stderr
				if err := upgradeCmd.Run(); err != nil {
					return fmt.Errorf(" failed to upgrade packages: %w", err)
				}
				fmt.Println(" All packages upgraded")
			}
		}

		fmt.Println()
		fmt.Println(" Update complete!")
		fmt.Println()
		fmt.Println(" Next steps:")
		fmt.Println("• Run: dotfiles cleanup # Remove old versions")
		fmt.Println("• Run: dotfiles doctor # Verify installation")
		return nil
	},
}

var upgradeCmd = &cobra.Command{
	Use: "upgrade",
	Short: "Alias for 'update' command",
	Long:  `  Upgrade packages (alias for 'update' command)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateCmd.RunE(cmd, args)
	},
}

func getOutdatedPackages() []string {
	cmd := exec.Command("brew", "outdated", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var outdated []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			outdated = append(outdated, line)
		}
	}
	return outdated
}

func init() {
	updateCmd.Flags().Bool("brew-only", false, "Only update Homebrew itself, not packages")
	updateCmd.Flags().Bool("dry-run", false, "Show what would be updated without actually updating")
	updateCmd.Flags().Bool("skip-brew-update", false, "Skip updating Homebrew itself")
	updateCmd.Flags().Bool("no-snapshot", false, "Skip creating automatic snapshot before update")

	upgradeCmd.Flags().Bool("brew-only", false, "Only update Homebrew itself, not packages")
	upgradeCmd.Flags().Bool("dry-run", false, "Show what would be updated without actually updating")
	upgradeCmd.Flags().Bool("skip-brew-update", false, "Skip updating Homebrew itself")
	upgradeCmd.Flags().Bool("no-snapshot", false, "Skip creating automatic snapshot before update")

	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(upgradeCmd)
}
