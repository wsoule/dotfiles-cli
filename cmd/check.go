package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	GroupID: "dotfiles",
	Short:   "🔍 Verify dotfiles setup and find issues",
	Long: `🔍 Check Dotfiles Setup

Comprehensive check of your dotfiles configuration:
• Verify config.json is valid
• Find broken symlinks
• Check for conflicts
• Validate stow packages
• Detect missing dependencies

Examples:
  dotfiles check                    # Full check
  dotfiles check --fix              # Auto-fix issues
  dotfiles check --links            # Only check symlinks`,
	Run: func(cmd *cobra.Command, args []string) {
		autoFix, _ := cmd.Flags().GetBool("fix")
		onlyLinks, _ := cmd.Flags().GetBool("links")

		fmt.Println("🔍 Checking dotfiles setup...")
		fmt.Println()

		issues := 0

		if !onlyLinks {
			issues += checkConfig(cmd, autoFix)
			issues += checkDependencies(autoFix)
			issues += checkStowPackagesFunc(cmd, autoFix)
		}

		issues += checkSymlinks(cmd, autoFix)

		fmt.Println()
		if issues == 0 {
			fmt.Println("✅ No issues found! Your dotfiles setup looks good.")
		} else {
			fmt.Printf("⚠️  Found %d issue(s)\n", issues)
			if !autoFix {
				fmt.Println("💡 Run with --fix to automatically fix issues")
			}
		}
	},
}

func checkConfig(cmd *cobra.Command, autoFix bool) int {
	fmt.Println("📋 Checking configuration...")
	issues := 0

	configPath := GetConfigPath(cmd)

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("   ❌ Config not found: %s\n", configPath)
		if autoFix {
			fmt.Println("   🔧 Creating default config...")
			cfg := &config.Config{
				Brews: []string{},
				Casks: []string{},
				Taps:  []string{},
				Stow:  []string{},
			}
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("   ❌ Failed to create config: %v\n", err)
			} else {
				fmt.Println("   ✅ Config created")
			}
		}
		issues++
		return issues
	}

	// Try to load config
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("   ❌ Failed to load config: %v\n", err)
		issues++
		return issues
	}

	fmt.Printf("   ✅ Config is valid (%d brews, %d casks, %d stow packages)\n",
		len(cfg.Brews), len(cfg.Casks), len(cfg.Stow))

	return issues
}

func checkDependencies(autoFix bool) int {
	fmt.Println("🔧 Checking dependencies...")
	issues := 0

	deps := map[string]string{
		"git":  "Version control",
		"stow": "Symlink manager",
	}

	for dep, desc := range deps {
		if _, err := exec.LookPath(dep); err != nil {
			fmt.Printf("   ❌ Missing: %s (%s)\n", dep, desc)
			issues++
			if autoFix {
				fmt.Printf("   💡 Install with: brew install %s (or your package manager)\n", dep)
			}
		} else {
			fmt.Printf("   ✅ Found: %s\n", dep)
		}
	}

	return issues
}

func checkStowPackagesFunc(cmd *cobra.Command, autoFix bool) int {
	fmt.Println("📦 Checking stow packages...")
	issues := 0

	dotfilesDir := GetDotfilesDir(cmd)
	stowDir := filepath.Join(dotfilesDir, "stow")

	if _, err := os.Stat(stowDir); os.IsNotExist(err) {
		fmt.Printf("   ⚠️  Stow directory not found: %s\n", stowDir)
		if autoFix {
			fmt.Println("   🔧 Creating stow directory...")
			if err := os.MkdirAll(stowDir, 0755); err != nil {
				fmt.Printf("   ❌ Failed: %v\n", err)
			} else {
				fmt.Println("   ✅ Created stow directory")
			}
		}
		issues++
		return issues
	}

	// List stow packages
	entries, err := os.ReadDir(stowDir)
	if err != nil {
		fmt.Printf("   ❌ Failed to read stow directory: %v\n", err)
		issues++
		return issues
	}

	if len(entries) == 0 {
		fmt.Println("   ⚠️  No stow packages found")
		issues++
	} else {
		fmt.Printf("   ✅ Found %d stow package(s)\n", len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Printf("      • %s\n", entry.Name())
			}
		}
	}

	return issues
}

func checkSymlinks(cmd *cobra.Command, autoFix bool) int {
	fmt.Println("🔗 Checking symlinks...")
	issues := 0

	home, _ := os.UserHomeDir()
	dotfilesDir := GetDotfilesDir(cmd)
	stowDir := filepath.Join(dotfilesDir, "stow")

	brokenLinks := []string{}

	// Walk through common dotfile locations
	checkPaths := []string{
		home,
		filepath.Join(home, ".config"),
	}

	for _, checkPath := range checkPaths {
		filepath.Walk(checkPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Check if it's a symlink
			if info.Mode()&os.ModeSymlink != 0 {
				target, err := os.Readlink(path)
				if err != nil {
					return nil
				}

				// Check if target exists
				var targetPath string
				if filepath.IsAbs(target) {
					targetPath = target
				} else {
					targetPath = filepath.Join(filepath.Dir(path), target)
				}

				if _, err := os.Stat(targetPath); os.IsNotExist(err) {
					// Check if it points to our stow directory
					if filepath.HasPrefix(targetPath, stowDir) || filepath.HasPrefix(target, stowDir) {
						brokenLinks = append(brokenLinks, path)
					}
				}
			}

			return nil
		})
	}

	if len(brokenLinks) > 0 {
		fmt.Printf("   ❌ Found %d broken symlink(s):\n", len(brokenLinks))
		for _, link := range brokenLinks {
			fmt.Printf("      • %s\n", link)
			if autoFix {
				fmt.Printf("      🔧 Removing broken symlink...\n")
				if err := os.Remove(link); err != nil {
					fmt.Printf("      ❌ Failed: %v\n", err)
				} else {
					fmt.Println("      ✅ Removed")
				}
			}
		}
		issues += len(brokenLinks)
	} else {
		fmt.Println("   ✅ No broken symlinks found")
	}

	return issues
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolP("fix", "f", false, "Automatically fix issues")
	checkCmd.Flags().Bool("links", false, "Only check symlinks")
}
