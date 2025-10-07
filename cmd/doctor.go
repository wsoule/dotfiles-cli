package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"dotfiles/internal/config"
	"dotfiles/internal/pkgmanager"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "🏥 Run health checks on your dotfiles setup",
	Long: `🏥 Dotfiles Health Check

Runs comprehensive diagnostics on your dotfiles setup to identify issues:
• Verifies dotfiles directory structure
• Checks for broken symlinks
• Validates configuration file
• Detects configuration drift
• Checks required dependencies
• Validates git repository status

Examples:
  dotfiles doctor              # Run all health checks
  dotfiles doctor --fix        # Auto-fix common issues
  dotfiles doctor --verbose    # Show detailed output`,
	Run: func(cmd *cobra.Command, args []string) {
		fix, _ := cmd.Flags().GetBool("fix")
		verbose, _ := cmd.Flags().GetBool("verbose")

		fmt.Println("🏥 Running Dotfiles Health Check...")
		fmt.Println("=" + strings.Repeat("=", 35))
		fmt.Println()

		issues := 0
		warnings := 0

		// Check 1: Dotfiles directory exists
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Cannot get home directory: %v\n", err)
			os.Exit(1)
		}

		dotfilesDir := filepath.Join(home, ".dotfiles")
		if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
			fmt.Println("❌ Dotfiles directory not found")
			fmt.Println("   Expected: ~/.dotfiles")
			fmt.Println("   💡 Run: dotfiles setup <repo-url> or dotfiles init")
			issues++
		} else {
			fmt.Println("✅ Dotfiles directory exists")
			if verbose {
				fmt.Printf("   Location: %s\n", dotfilesDir)
			}
		}
		fmt.Println()

		// Check 2: Configuration file
		configPath := filepath.Join(dotfilesDir, "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Println("❌ Configuration file invalid or missing")
			fmt.Printf("   Error: %v\n", err)
			fmt.Println("   💡 Run: dotfiles init")
			issues++
		} else {
			fmt.Println("✅ Configuration file is valid")
			if verbose {
				totalPkgs := len(cfg.Brews) + len(cfg.Casks) + len(cfg.Taps) + len(cfg.Stow)
				fmt.Printf("   Packages: %d total (%d brews, %d casks, %d taps, %d stow)\n",
					totalPkgs, len(cfg.Brews), len(cfg.Casks), len(cfg.Taps), len(cfg.Stow))
			}
		}
		fmt.Println()

		// Check 3: Git repository
		gitDir := filepath.Join(dotfilesDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			fmt.Println("⚠️  Not a git repository")
			fmt.Println("   💡 Run: git init in ~/.dotfiles to enable version control")
			warnings++
		} else {
			fmt.Println("✅ Git repository initialized")

			// Check for remote
			os.Chdir(dotfilesDir)
			remoteCmd := exec.Command("git", "remote", "-v")
			remoteOutput, _ := remoteCmd.Output()
			if len(remoteOutput) == 0 {
				fmt.Println("   ⚠️  No remote repository configured")
				fmt.Println("   💡 Add remote: git remote add origin <url>")
				warnings++
			} else if verbose {
				fmt.Println("   Remote configured")
			}
		}
		fmt.Println()

		// Check 4: Required dependencies
		fmt.Println("📋 Checking Dependencies...")

		// Get package manager
		pm, err := pkgmanager.GetPackageManager()
		pmAvailable := false
		if err == nil && pm.IsAvailable() {
			fmt.Printf("✅ %s installed\n", pm.GetName())
			pmAvailable = true
		} else {
			fmt.Printf("❌ Package manager not found\n")
			issues++
			if runtime.GOOS == "darwin" {
				fmt.Println("   💡 Install Homebrew: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
			}
		}

		// Check git
		if _, err := exec.LookPath("git"); err == nil {
			fmt.Printf("✅ git installed\n")
		} else {
			fmt.Printf("❌ git not found\n")
			issues++
			if pmAvailable {
				if pm.GetName() == "homebrew" {
					fmt.Println("   💡 Install git: brew install git")
				} else if pm.GetName() == "pacman" {
					fmt.Println("   💡 Install git: sudo pacman -S git")
				} else {
					fmt.Println("   💡 Install git with your package manager")
				}
			}
		}

		// Check stow
		if _, err := exec.LookPath("stow"); err == nil {
			fmt.Printf("✅ stow installed\n")
		} else {
			fmt.Printf("❌ stow not found\n")
			issues++
			if pmAvailable {
				if pm.GetName() == "homebrew" {
					fmt.Println("   💡 Install GNU Stow: brew install stow")
				} else if pm.GetName() == "pacman" {
					fmt.Println("   💡 Install GNU Stow: sudo pacman -S stow")
				} else {
					fmt.Println("   💡 Install GNU Stow with your package manager")
				}
			}
		}
		fmt.Println()

		// Check 5: Broken symlinks
		if cfg != nil && len(cfg.Stow) > 0 {
			fmt.Println("🔗 Checking Symlinks...")
			brokenLinks := checkBrokenSymlinks(home, verbose)
			if len(brokenLinks) > 0 {
				fmt.Printf("❌ Found %d broken symlinks:\n", len(brokenLinks))
				for _, link := range brokenLinks {
					fmt.Printf("   • %s\n", link)
				}
				if fix {
					fmt.Println("   🔧 Auto-fix not implemented for broken symlinks")
					fmt.Println("   💡 Run: dotfiles restow <package>")
				}
				issues += len(brokenLinks)
			} else {
				fmt.Println("✅ No broken symlinks found")
			}
			fmt.Println()
		}

		// Check 6: Configuration drift
		if pmAvailable && cfg != nil {
			fmt.Println("📊 Checking Configuration Drift...")
			drift := checkConfigDrift(cfg)
			if drift.MissingBrews > 0 || drift.MissingCasks > 0 {
				fmt.Printf("⚠️  Configuration drift detected:\n")
				if drift.MissingBrews > 0 {
					fmt.Printf("   • %d packages configured but not installed\n", drift.MissingBrews)
				}
				if drift.MissingCasks > 0 && runtime.GOOS == "darwin" {
					fmt.Printf("   • %d casks configured but not installed\n", drift.MissingCasks)
				}
				if drift.ExtraBrews > 0 {
					fmt.Printf("   • %d packages installed but not in config\n", drift.ExtraBrews)
				}
				if drift.ExtraCasks > 0 && runtime.GOOS == "darwin" {
					fmt.Printf("   • %d casks installed but not in config\n", drift.ExtraCasks)
				}
				fmt.Println("   💡 Run: dotfiles diff")
				if drift.ExtraBrews > 0 || drift.ExtraCasks > 0 {
					fmt.Println("   💡 Run: dotfiles scan to add missing packages")
				}
				warnings++
			} else {
				fmt.Println("✅ Configuration in sync with installed packages")
			}
			fmt.Println()
		}

		// Check 7: Stow directory structure
		if cfg != nil && len(cfg.Stow) > 0 {
			fmt.Println("📁 Checking Stow Packages...")
			stowDir := filepath.Join(dotfilesDir, "stow")
			missingPkgs := []string{}
			for _, pkg := range cfg.Stow {
				pkgDir := filepath.Join(stowDir, pkg)
				if _, err := os.Stat(pkgDir); os.IsNotExist(err) {
					missingPkgs = append(missingPkgs, pkg)
				}
			}
			if len(missingPkgs) > 0 {
				fmt.Printf("❌ %d stow packages missing:\n", len(missingPkgs))
				for _, pkg := range missingPkgs {
					fmt.Printf("   • %s (expected at: stow/%s)\n", pkg, pkg)
				}
				issues += len(missingPkgs)
			} else {
				fmt.Println("✅ All stow packages exist")
			}
			fmt.Println()
		}

		// Summary
		fmt.Println("=" + strings.Repeat("=", 35))
		if issues == 0 && warnings == 0 {
			fmt.Println("🎉 All checks passed! Your dotfiles are healthy.")
		} else {
			if issues > 0 {
				fmt.Printf("❌ Found %d issue(s)\n", issues)
			}
			if warnings > 0 {
				fmt.Printf("⚠️  Found %d warning(s)\n", warnings)
			}
			fmt.Println()
			fmt.Println("💡 Review the suggestions above to fix issues")
		}
	},
}

type ConfigDrift struct {
	MissingBrews int
	MissingCasks int
	ExtraBrews   int
	ExtraCasks   int
}

func checkBrokenSymlinks(homeDir string, verbose bool) []string {
	var broken []string

	// Check common locations for symlinks
	locations := []string{
		homeDir,
		filepath.Join(homeDir, ".config"),
	}

	for _, loc := range locations {
		filepath.Walk(loc, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}

			// Only check symlinks
			if info.Mode()&os.ModeSymlink != 0 {
				target, err := os.Readlink(path)
				if err != nil {
					broken = append(broken, path)
					return nil
				}

				// Check if target exists
				if !filepath.IsAbs(target) {
					target = filepath.Join(filepath.Dir(path), target)
				}
				if _, err := os.Stat(target); os.IsNotExist(err) {
					broken = append(broken, path)
				}
			}
			return nil
		})
	}

	return broken
}

func checkConfigDrift(cfg *config.Config) ConfigDrift {
	drift := ConfigDrift{}

	// Get installed packages
	installedBrews, err := getInstalledBrews()
	if err != nil {
		return drift
	}
	installedCasks, err := getInstalledCasks()
	if err != nil {
		return drift
	}

	// Create maps for quick lookup
	installedBrewMap := make(map[string]bool)
	for _, brew := range installedBrews {
		installedBrewMap[brew] = true
	}
	installedCaskMap := make(map[string]bool)
	for _, cask := range installedCasks {
		installedCaskMap[cask] = true
	}

	configBrewMap := make(map[string]bool)
	for _, brew := range cfg.Brews {
		configBrewMap[brew] = true
	}
	configCaskMap := make(map[string]bool)
	for _, cask := range cfg.Casks {
		configCaskMap[cask] = true
	}

	// Count missing (in config but not installed)
	for _, brew := range cfg.Brews {
		if !installedBrewMap[brew] {
			drift.MissingBrews++
		}
	}
	for _, cask := range cfg.Casks {
		if !installedCaskMap[cask] {
			drift.MissingCasks++
		}
	}

	// Count extra (installed but not in config)
	for _, brew := range installedBrews {
		if !configBrewMap[brew] {
			drift.ExtraBrews++
		}
	}
	for _, cask := range installedCasks {
		if !configCaskMap[cask] {
			drift.ExtraCasks++
		}
	}

	return drift
}

func init() {
	doctorCmd.Flags().Bool("fix", false, "Attempt to auto-fix common issues")
	doctorCmd.Flags().BoolP("verbose", "v", false, "Show detailed output")

	rootCmd.AddCommand(doctorCmd)
}
