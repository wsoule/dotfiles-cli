package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check package installation status",
	Long:  `Compare configured packages with what's actually installed on the system`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Check if brew is available
		if _, err := exec.LookPath("brew"); err != nil {
			fmt.Println("⚠️  Homebrew not found. Cannot check package status.")
			return
		}

		fmt.Println("📊 Package Status Report")
		fmt.Println("=" + strings.Repeat("=", 23))

		// Check taps
		if len(cfg.Taps) > 0 {
			checkTaps(cfg.Taps)
		}

		// Check brews
		if len(cfg.Brews) > 0 {
			checkBrews(cfg.Brews)
		}

		// Check casks
		if len(cfg.Casks) > 0 {
			checkCasks(cfg.Casks)
		}

		// Check stow packages
		if len(cfg.Stow) > 0 {
			checkStowPackages(cfg.Stow)
		}

		if len(cfg.Taps)+len(cfg.Brews)+len(cfg.Casks)+len(cfg.Stow) == 0 {
			fmt.Println("No packages configured. Run 'dotfiles add <package>' to get started.")
		}
	},
}

func checkTaps(configuredTaps []string) {
	fmt.Println("\n📋 Taps:")

	// Get installed taps
	installedTaps := getInstalledTaps()
	installedSet := make(map[string]bool)
	for _, tap := range installedTaps {
		installedSet[tap] = true
	}

	missing := []string{}
	for _, tap := range configuredTaps {
		if installedSet[tap] {
			fmt.Printf("  ✅ %s\n", tap)
		} else {
			fmt.Printf("  ❌ %s (not tapped)\n", tap)
			missing = append(missing, tap)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("  → Run: brew tap %s\n", strings.Join(missing, " "))
	}
}

func checkBrews(configuredBrews []string) {
	fmt.Println("\n🍺 Brews:")

	// Get installed brews
	installedBrews, err := getInstalledBrews()
	if err != nil {
		fmt.Printf("  ❌ Error getting installed brews: %v\n", err)
		return
	}
	installedSet := make(map[string]bool)
	for _, brew := range installedBrews {
		installedSet[brew] = true
	}

	missing := []string{}
	for _, brew := range configuredBrews {
		if installedSet[brew] {
			fmt.Printf("  ✅ %s\n", brew)
		} else {
			fmt.Printf("  ❌ %s (not installed)\n", brew)
			missing = append(missing, brew)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("  → Run: brew install %s\n", strings.Join(missing, " "))
	}
}

func checkCasks(configuredCasks []string) {
	fmt.Println("\n📦 Casks:")

	// Get installed casks
	installedCasks, err := getInstalledCasks()
	if err != nil {
		fmt.Printf("  ❌ Error getting installed casks: %v\n", err)
		return
	}
	installedSet := make(map[string]bool)
	for _, cask := range installedCasks {
		installedSet[cask] = true
	}

	missing := []string{}
	for _, cask := range configuredCasks {
		if installedSet[cask] {
			fmt.Printf("  ✅ %s\n", cask)
		} else {
			fmt.Printf("  ❌ %s (not installed)\n", cask)
			missing = append(missing, cask)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("  → Run: brew install --cask %s\n", strings.Join(missing, " "))
	}
}

func getInstalledTaps() []string {
	cmd := exec.Command("brew", "tap")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var taps []string
	for _, line := range lines {
		if line != "" {
			taps = append(taps, strings.TrimSpace(line))
		}
	}
	return taps
}

func checkStowPackages(configuredStow []string) {
	fmt.Println("\n🔗 Stow Packages:")

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("  ❌ Error getting home directory: %v\n", err)
		return
	}

	stowDir := filepath.Join(home, ".dotfiles", "stow")

	// Check if stow is available
	if _, err := exec.LookPath("stow"); err != nil {
		fmt.Println("  ⚠️  GNU Stow not found. Install with: brew install stow")
		return
	}

	missing := []string{}
	for _, pkg := range configuredStow {
		// Check if package directory exists
		pkgPath := filepath.Join(stowDir, pkg)
		if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
			fmt.Printf("  ❌ %s (directory not found: %s)\n", pkg, pkgPath)
			missing = append(missing, pkg)
			continue
		}

		// Check if package is stowed (has symlinks in home)
		isStowed := checkIfStowed(pkg, stowDir, home)
		if isStowed {
			fmt.Printf("  ✅ %s (stowed)\n", pkg)
		} else {
			fmt.Printf("  ⚠️  %s (not stowed)\n", pkg)
			missing = append(missing, pkg)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("  → Run: dotfiles stow %s\n", strings.Join(missing, " "))
	}
}

func checkIfStowed(pkg, stowDir, target string) bool {
	// Run stow in simulation mode to check if package is already stowed
	cmd := exec.Command("stow", "-d", stowDir, "-t", target, "-n", "-v", pkg)
	output, err := cmd.CombinedOutput()

	// If stow reports "LINK" operations, it means the package is not stowed
	// If it reports no operations or only "UNLINK/LINK" pairs, it's already stowed
	if err == nil && !strings.Contains(string(output), "LINK") {
		return true
	}

	// Alternative check: look for existing symlinks
	pkgPath := filepath.Join(stowDir, pkg)
	return hasSymlinksInTarget(pkgPath, target, "")
}

func hasSymlinksInTarget(pkgPath, target, subPath string) bool {
	currentPkgPath := filepath.Join(pkgPath, subPath)
	currentTargetPath := filepath.Join(target, subPath)

	entries, err := os.ReadDir(currentPkgPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Recursively check subdirectories
			if hasSymlinksInTarget(pkgPath, target, filepath.Join(subPath, entry.Name())) {
				return true
			}
		} else {
			// Check if file exists as symlink in target
			targetFile := filepath.Join(currentTargetPath, entry.Name())
			if info, err := os.Lstat(targetFile); err == nil && info.Mode()&os.ModeSymlink != 0 {
				if link, err := os.Readlink(targetFile); err == nil {
					expectedPath := filepath.Join(currentPkgPath, entry.Name())
					if absLink, err := filepath.Abs(link); err == nil {
						if absExpected, err := filepath.Abs(expectedPath); err == nil {
							if absLink == absExpected {
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(statusCmd)
}