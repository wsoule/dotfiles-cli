package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var stowDiffCmd = &cobra.Command{
	Use:     "stow-diff [packages...]",
	GroupID: "dotfiles",
	Short:   "🔍 Preview stow changes before applying",
	Long: `🔍 Preview Stow Changes

Show what will change when stowing packages:
• Files that will be symlinked
• Conflicts with existing files
• Content differences

Examples:
  dotfiles stow-diff                     # Show all package diffs
  dotfiles stow-diff zsh vim            # Show specific packages
  dotfiles stow-diff --verbose          # Show full file contents`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		dotfilesDir := GetDotfilesDir(cmd)
		stowDir := filepath.Join(dotfilesDir, "stow")

		packages := args
		if len(packages) == 0 {
			// Get all packages from stow directory
			entries, err := os.ReadDir(stowDir)
			if err != nil {
				fmt.Printf("❌ Failed to read stow directory: %v\n", err)
				return
			}

			for _, entry := range entries {
				if entry.IsDir() {
					packages = append(packages, entry.Name())
				}
			}
		}

		if len(packages) == 0 {
			fmt.Println("No packages found in stow directory")
			return
		}

		home, _ := os.UserHomeDir()
		hasChanges := false

		for _, pkg := range packages {
			pkgPath := filepath.Join(stowDir, pkg)
			if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
				fmt.Printf("⚠️  Package '%s' not found\n", pkg)
				continue
			}

			fmt.Printf("\n📦 Package: %s\n", pkg)
			fmt.Println(strings.Repeat("─", 50))

			pkgHasChanges := false

			// Walk through package files
			filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Calculate relative path and target
				relPath, _ := filepath.Rel(pkgPath, path)
				targetPath := filepath.Join(home, relPath)

				// Check if target exists
				targetInfo, targetErr := os.Lstat(targetPath)

				if os.IsNotExist(targetErr) {
					// New file
					fmt.Printf("  ✨ NEW: %s\n", relPath)
					if verbose {
						content, _ := os.ReadFile(path)
						fmt.Printf("     Content:\n%s\n", indentLines(string(content), 6))
					}
					pkgHasChanges = true
					hasChanges = true
				} else if targetErr == nil {
					// File exists - check if it's a symlink
					if targetInfo.Mode()&os.ModeSymlink != 0 {
						// It's a symlink - check if it points to our file
						linkTarget, _ := os.Readlink(targetPath)
						if linkTarget == path {
							fmt.Printf("  ✅ OK: %s (already linked)\n", relPath)
						} else {
							fmt.Printf("  ⚠️  CONFLICT: %s (links to different file)\n", relPath)
							fmt.Printf("     Current: %s\n", linkTarget)
							fmt.Printf("     Would be: %s\n", path)
							pkgHasChanges = true
							hasChanges = true
						}
					} else {
						// Regular file exists - show diff
						fmt.Printf("  🔀 DIFF: %s\n", relPath)
						pkgHasChanges = true
						hasChanges = true

						if verbose {
							// Show actual diff
							diffCmd := exec.Command("diff", "-u", targetPath, path)
							diffOutput, _ := diffCmd.CombinedOutput()
							if len(diffOutput) > 0 {
								fmt.Printf("     %s\n", indentLines(string(diffOutput), 6))
							}
						} else {
							fmt.Printf("     Use --verbose to see diff\n")
						}
					}
				}

				return nil
			})

			if !pkgHasChanges {
				fmt.Println("  ✅ No changes")
			}
		}

		fmt.Println()
		if hasChanges {
			fmt.Println("💡 Run 'dotfiles stow' to apply these changes")
		} else {
			fmt.Println("✅ All packages are up to date")
		}
	},
}

func indentLines(text string, spaces int) string {
	indent := strings.Repeat(" ", spaces)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

func init() {
	rootCmd.AddCommand(stowDiffCmd)
	stowDiffCmd.Flags().BoolP("verbose", "v", false, "Show full file contents and diffs")
}
