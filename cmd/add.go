package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <packages>",
	Short: "📦 Add packages to your configuration",
	Long: `📦 Add Packages - Expand Your Development Environment

Add Homebrew packages (brews), applications (casks), repositories (taps),
or dotfile packages (stow) to your configuration.

Package Types:
• brew  - Command-line tools (git, curl, node, etc.)
• cask  - GUI applications (VS Code, Docker, etc.)
• tap   - Additional Homebrew repositories
• stow  - Dotfile packages for symlinking

Examples:
  dotfiles add git curl wget                    # Add CLI tools (default: brew)
  dotfiles add --type=cask visual-studio-code  # Add GUI applications
  dotfiles add --type=tap homebrew/cask-fonts  # Add font repository
  dotfiles add --type=stow vim zsh tmux        # Add dotfile packages
  dotfiles add --file=packages.txt             # Add from file (one per line)

Popular packages:
• Essential: git, curl, wget, tree, jq, gh, docker
• Productivity: fzf, ripgrep, bat, eza, tmux, neovim
• Applications: visual-studio-code, rectangle, slack`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Load existing config
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			fmt.Println("Run 'dotfiles init' to create a configuration first.")
			os.Exit(1)
		}

		var packages []string

		// Handle file input
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			filePackages, err := readPackagesFromFile(file)
			if err != nil {
				fmt.Printf("Error reading packages from file: %v\n", err)
				os.Exit(1)
			}
			packages = append(packages, filePackages...)
		}

		// Add command line arguments
		packages = append(packages, args...)

		if len(packages) == 0 {
			fmt.Println("No packages specified. Use command line arguments or --file flag.")
			return
		}

		packageType, _ := cmd.Flags().GetString("type")
		added := 0

		for _, pkg := range packages {
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				continue
			}

			switch packageType {
			case "cask":
				if !contains(cfg.Casks, pkg) {
					cfg.Casks = append(cfg.Casks, pkg)
					fmt.Printf("✓ Added cask: %s\n", pkg)
					added++
				} else {
					fmt.Printf("- Cask %s already exists\n", pkg)
				}
			case "tap":
				if !contains(cfg.Taps, pkg) {
					cfg.Taps = append(cfg.Taps, pkg)
					fmt.Printf("✓ Added tap: %s\n", pkg)
					added++
				} else {
					fmt.Printf("- Tap %s already exists\n", pkg)
				}
			case "stow":
				if !contains(cfg.Stow, pkg) {
					cfg.Stow = append(cfg.Stow, pkg)
					fmt.Printf("✓ Added stow package: %s\n", pkg)
					added++
				} else {
					fmt.Printf("- Stow package %s already exists\n", pkg)
				}
			default: // brew
				if !contains(cfg.Brews, pkg) {
					cfg.Brews = append(cfg.Brews, pkg)
					fmt.Printf("✓ Added brew: %s\n", pkg)
					added++
				} else {
					fmt.Printf("- Brew %s already exists\n", pkg)
				}
			}
		}

		if added > 0 {
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\n📊 Added %d new packages\n", added)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <packages>",
	Short: "🗑️  Remove packages from configuration and optionally uninstall",
	Long: `🗑️  Remove Packages

Remove packages from your configuration and optionally uninstall them from your system.

Examples:
  dotfiles remove git curl                      # Remove from config only
  dotfiles remove --uninstall git curl          # Remove from config AND uninstall
  dotfiles remove --type=cask slack             # Remove cask from config
  dotfiles remove --type=cask --uninstall slack # Remove and uninstall cask
  dotfiles remove --all-brews                   # Remove all brews from config
  dotfiles remove --all-brews --uninstall       # Remove and uninstall all brews`,
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

		packageType, _ := cmd.Flags().GetString("type")
		uninstall, _ := cmd.Flags().GetBool("uninstall")

		var packagesToUninstall []string
		var packageTypeForUninstall string

		// Handle bulk removal flags
		if allBrews, _ := cmd.Flags().GetBool("all-brews"); allBrews {
			removed := len(cfg.Brews)
			if uninstall {
				packagesToUninstall = append(packagesToUninstall, cfg.Brews...)
				packageTypeForUninstall = "brew"
			}
			cfg.Brews = []string{}
			fmt.Printf("✓ Removed all %d brews from config\n", removed)
		}

		if allCasks, _ := cmd.Flags().GetBool("all-casks"); allCasks {
			removed := len(cfg.Casks)
			if uninstall {
				packagesToUninstall = append(packagesToUninstall, cfg.Casks...)
				packageTypeForUninstall = "cask"
			}
			cfg.Casks = []string{}
			fmt.Printf("✓ Removed all %d casks from config\n", removed)
		}

		if allTaps, _ := cmd.Flags().GetBool("all-taps"); allTaps {
			removed := len(cfg.Taps)
			if uninstall {
				packagesToUninstall = append(packagesToUninstall, cfg.Taps...)
				packageTypeForUninstall = "tap"
			}
			cfg.Taps = []string{}
			fmt.Printf("✓ Removed all %d taps from config\n", removed)
		}

		if allStow, _ := cmd.Flags().GetBool("all-stow"); allStow {
			removed := len(cfg.Stow)
			cfg.Stow = []string{}
			fmt.Printf("✓ Removed all %d stow packages from config\n", removed)
		}

		// Handle individual package removal
		var packages []string

		// Handle file input
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			filePackages, err := readPackagesFromFile(file)
			if err != nil {
				fmt.Printf("Error reading packages from file: %v\n", err)
				os.Exit(1)
			}
			packages = append(packages, filePackages...)
		}

		// Add command line arguments
		packages = append(packages, args...)

		removed := 0
		for _, pkg := range packages {
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				continue
			}

			switch packageType {
			case "cask":
				if contains(cfg.Casks, pkg) {
					cfg.Casks = removeFromSlice(cfg.Casks, pkg)
					fmt.Printf("✓ Removed cask from config: %s\n", pkg)
					removed++
					if uninstall {
						packagesToUninstall = append(packagesToUninstall, pkg)
						packageTypeForUninstall = "cask"
					}
				} else {
					fmt.Printf("- Cask %s not found in config\n", pkg)
				}
			case "tap":
				if contains(cfg.Taps, pkg) {
					cfg.Taps = removeFromSlice(cfg.Taps, pkg)
					fmt.Printf("✓ Removed tap from config: %s\n", pkg)
					removed++
					if uninstall {
						packagesToUninstall = append(packagesToUninstall, pkg)
						packageTypeForUninstall = "tap"
					}
				} else {
					fmt.Printf("- Tap %s not found in config\n", pkg)
				}
			case "stow":
				if contains(cfg.Stow, pkg) {
					cfg.Stow = removeFromSlice(cfg.Stow, pkg)
					fmt.Printf("✓ Removed stow package from config: %s\n", pkg)
					removed++
				} else {
					fmt.Printf("- Stow package %s not found in config\n", pkg)
				}
			default: // brew
				if contains(cfg.Brews, pkg) {
					cfg.Brews = removeFromSlice(cfg.Brews, pkg)
					fmt.Printf("✓ Removed brew from config: %s\n", pkg)
					removed++
					if uninstall {
						packagesToUninstall = append(packagesToUninstall, pkg)
						packageTypeForUninstall = "brew"
					}
				} else {
					fmt.Printf("- Brew %s not found in config\n", pkg)
				}
			}
		}

		if removed > 0 {
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\n📊 Removed %d packages from config\n", removed)
		}

		// Uninstall packages if requested
		if uninstall && len(packagesToUninstall) > 0 {
			fmt.Println()
			uninstallPackages(packagesToUninstall, packageTypeForUninstall)
		}
	},
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeFromSlice(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

func readPackagesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			packages = append(packages, line)
		}
	}

	return packages, scanner.Err()
}

func uninstallPackages(packages []string, pkgType string) {
	if len(packages) == 0 {
		return
	}

	fmt.Printf("🗑️  Uninstalling %d package(s) from system...\n", len(packages))
	fmt.Println()

	var cmd *exec.Cmd
	switch pkgType {
	case "cask":
		cmd = exec.Command("brew", append([]string{"uninstall", "--cask"}, packages...)...)
	case "tap":
		for _, tap := range packages {
			tapCmd := exec.Command("brew", "untap", tap)
			tapCmd.Stdout = os.Stdout
			tapCmd.Stderr = os.Stderr
			if err := tapCmd.Run(); err != nil {
				fmt.Printf("❌ Failed to untap %s: %v\n", tap, err)
			} else {
				fmt.Printf("✅ Untapped %s\n", tap)
			}
		}
		return
	default: // brew
		cmd = exec.Command("brew", append([]string{"uninstall"}, packages...)...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to uninstall packages: %v\n", err)
	} else {
		fmt.Printf("✅ Uninstalled %d package(s)\n", len(packages))
	}
}

func init() {
	addCmd.Flags().StringP("type", "t", "brew", "Package type (brew, cask, tap, stow)")
	addCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")

	removeCmd.Flags().StringP("type", "t", "brew", "Package type (brew, cask, tap, stow)")
	removeCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	removeCmd.Flags().Bool("uninstall", false, "Also uninstall packages from system (not just config)")
	removeCmd.Flags().Bool("all-brews", false, "Remove all brew packages")
	removeCmd.Flags().Bool("all-casks", false, "Remove all cask packages")
	removeCmd.Flags().Bool("all-taps", false, "Remove all taps")
	removeCmd.Flags().Bool("all-stow", false, "Remove all stow packages")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
}