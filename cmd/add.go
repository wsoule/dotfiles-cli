package cmd

import (
	"bufio"
	"fmt"
	"os"
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
	Short: "Remove packages from your configuration",
	Long:  `Remove brew, cask, or tap packages from your config.json`,
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

		// Handle bulk removal flags
		if allBrews, _ := cmd.Flags().GetBool("all-brews"); allBrews {
			removed := len(cfg.Brews)
			cfg.Brews = []string{}
			fmt.Printf("✓ Removed all %d brews\n", removed)
		}

		if allCasks, _ := cmd.Flags().GetBool("all-casks"); allCasks {
			removed := len(cfg.Casks)
			cfg.Casks = []string{}
			fmt.Printf("✓ Removed all %d casks\n", removed)
		}

		if allTaps, _ := cmd.Flags().GetBool("all-taps"); allTaps {
			removed := len(cfg.Taps)
			cfg.Taps = []string{}
			fmt.Printf("✓ Removed all %d taps\n", removed)
		}

		if allStow, _ := cmd.Flags().GetBool("all-stow"); allStow {
			removed := len(cfg.Stow)
			cfg.Stow = []string{}
			fmt.Printf("✓ Removed all %d stow packages\n", removed)
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
					fmt.Printf("✓ Removed cask: %s\n", pkg)
					removed++
				} else {
					fmt.Printf("- Cask %s not found\n", pkg)
				}
			case "tap":
				if contains(cfg.Taps, pkg) {
					cfg.Taps = removeFromSlice(cfg.Taps, pkg)
					fmt.Printf("✓ Removed tap: %s\n", pkg)
					removed++
				} else {
					fmt.Printf("- Tap %s not found\n", pkg)
				}
			case "stow":
				if contains(cfg.Stow, pkg) {
					cfg.Stow = removeFromSlice(cfg.Stow, pkg)
					fmt.Printf("✓ Removed stow package: %s\n", pkg)
					removed++
				} else {
					fmt.Printf("- Stow package %s not found\n", pkg)
				}
			default: // brew
				if contains(cfg.Brews, pkg) {
					cfg.Brews = removeFromSlice(cfg.Brews, pkg)
					fmt.Printf("✓ Removed brew: %s\n", pkg)
					removed++
				} else {
					fmt.Printf("- Brew %s not found\n", pkg)
				}
			}
		}

		if removed > 0 {
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\n📊 Removed %d packages\n", removed)
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

func init() {
	addCmd.Flags().StringP("type", "t", "brew", "Package type (brew, cask, tap, stow)")
	addCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")

	removeCmd.Flags().StringP("type", "t", "brew", "Package type (brew, cask, tap, stow)")
	removeCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	removeCmd.Flags().Bool("all-brews", false, "Remove all brew packages")
	removeCmd.Flags().Bool("all-casks", false, "Remove all cask packages")
	removeCmd.Flags().Bool("all-taps", false, "Remove all taps")
	removeCmd.Flags().Bool("all-stow", false, "Remove all stow packages")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
}