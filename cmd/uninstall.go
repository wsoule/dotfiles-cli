package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use: "uninstall",
	GroupID: "dotfiles",
	Short: "Remove dotfiles setup",
	Long: `  Uninstall Dotfiles Manager

Removes the dotfiles setup from your system:
• Unstows all packages (removes symlinks)
• Optionally removes ~/.dotfiles directory
• Optionally removes dotfiles CLI tool

Examples:
  dotfiles uninstall                      # Interactive uninstall
  dotfiles uninstall --keep-files         # Remove links but keep files
  dotfiles uninstall --force              # Skip confirmation`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		keepFiles, _ := cmd.Flags().GetBool("keep-files")

		fmt.Println("Dotfiles Uninstaller")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Println()

		if !force {
			fmt.Println("This will:")
			fmt.Println("• Remove all symlinks created by stow")
			if !keepFiles {
				fmt.Println("• Delete ~/.dotfiles directory")
			}
			fmt.Println()

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Continue? (yes/NO): ")
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "yes" {
				fmt.Println(" Uninstall cancelled")
				return
			}
		}

		fmt.Println()
		fmt.Println(" Starting uninstall process...")
		fmt.Println()

		// Step 1: Unstow all packages
		fmt.Println(" Step 1: Removing symlinks...")
		if err := unstowAllPackages(cmd); err != nil {
			fmt.Printf("Some symlinks may not have been removed: %v\n", err)
		} else {
			fmt.Println(" Symlinks removed")
		}
		fmt.Println()

		// Step 2: Remove dotfiles directory
		if !keepFiles {
			fmt.Println(" Step 2: Removing dotfiles directory...")
			dotfilesDir := GetDotfilesDir(cmd)
			if err := os.RemoveAll(dotfilesDir); err != nil {
				fmt.Printf(" Failed to remove %s: %v\n", dotfilesDir, err)
			} else {
				fmt.Printf(" Removed %s\n", dotfilesDir)
			}
			fmt.Println()
		}

		// Step 3: Show manual cleanup steps
		fmt.Println(" Manual cleanup needed:")
		fmt.Println("• Uninstall dotfiles CLI if installed globally")
		fmt.Println("• Review and remove any shell configuration references")
		fmt.Println()

		fmt.Println(" Uninstall complete!")
	},
}

func unstowAllPackages(cmd *cobra.Command) error {
	home, _ := os.UserHomeDir()
	dotfilesDir := GetDotfilesDir(cmd)
	stowDir := filepath.Join(dotfilesDir, "stow")

	if _, err := os.Stat(stowDir); os.IsNotExist(err) {
		return fmt.Errorf("stow directory not found")
	}

	// Get all stow packages
	entries, err := os.ReadDir(stowDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pkg := entry.Name()
			fmt.Printf("Unstowing %s...\n", pkg)

			unstowCmd := exec.Command("stow", "-D", "-t", home, "-d", stowDir, pkg)
			if err := unstowCmd.Run(); err != nil {
				fmt.Printf("Failed to unstow %s: %v\n", pkg, err)
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
	uninstallCmd.Flags().Bool("keep-files", false, "Keep ~/.dotfiles directory")
}
