package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "🔍 Scan system for installed packages and add them to config",
	Long: `🔍 System Package Scanner

Scans your system for packages already installed via Homebrew and helps you
add them to your dotfiles configuration. Perfect for when you're setting up
dotfiles on an existing machine with packages already installed.

Examples:
  dotfiles scan                    # Scan and interactively select packages
  dotfiles scan --auto             # Automatically add all installed packages
  dotfiles scan --brews-only       # Only scan Homebrew formulas
  dotfiles scan --casks-only       # Only scan Homebrew casks`,
	Run: func(cmd *cobra.Command, args []string) {
		auto, _ := cmd.Flags().GetBool("auto")
		brewsOnly, _ := cmd.Flags().GetBool("brews-only")
		casksOnly, _ := cmd.Flags().GetBool("casks-only")

		fmt.Println("🔍 Scanning system for installed packages...")
		fmt.Println()

		// Load current config
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("⚠️  Could not load config (will create new): %v\n", err)
			cfg = &config.Config{
				Brews: []string{},
				Casks: []string{},
				Taps:  []string{},
				Stow:  []string{},
			}
		}

		var newBrews, newCasks []string

		// Scan brews
		if !casksOnly {
			installedBrews, err := getInstalledBrews()
			if err != nil {
				fmt.Printf("⚠️  Could not scan Homebrew formulas: %v\n", err)
			} else {
				newBrews = filterNewPackages(installedBrews, cfg.Brews)
				if len(newBrews) > 0 {
					fmt.Printf("📋 Found %d Homebrew formulas not in config\n", len(newBrews))
				}
			}
		}

		// Scan casks
		if !brewsOnly {
			installedCasks, err := getInstalledCasks()
			if err != nil {
				fmt.Printf("⚠️  Could not scan Homebrew casks: %v\n", err)
			} else {
				newCasks = filterNewPackages(installedCasks, cfg.Casks)
				if len(newCasks) > 0 {
					fmt.Printf("📦 Found %d Homebrew casks not in config\n", len(newCasks))
				}
			}
		}

		if len(newBrews) == 0 && len(newCasks) == 0 {
			fmt.Println("✅ All installed packages are already in your config!")
			return
		}

		fmt.Println()

		// Auto mode - add everything
		if auto {
			cfg.Brews = append(cfg.Brews, newBrews...)
			cfg.Casks = append(cfg.Casks, newCasks...)

			// Sort for cleaner config
			sort.Strings(cfg.Brews)
			sort.Strings(cfg.Casks)

			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("❌ Failed to save config: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Added %d brews and %d casks to config\n", len(newBrews), len(newCasks))
			return
		}

		// Interactive mode
		fmt.Println("Select packages to add to your config:")
		fmt.Println()

		selectedBrews := []string{}
		selectedCasks := []string{}

		if len(newBrews) > 0 {
			fmt.Println("🍺 Homebrew Formulas:")
			selectedBrews = selectPackages(newBrews)
		}

		if len(newCasks) > 0 {
			fmt.Println("📦 Homebrew Casks:")
			selectedCasks = selectPackages(newCasks)
		}

		if len(selectedBrews) == 0 && len(selectedCasks) == 0 {
			fmt.Println("No packages selected. Exiting.")
			return
		}

		// Add selected packages to config
		cfg.Brews = append(cfg.Brews, selectedBrews...)
		cfg.Casks = append(cfg.Casks, selectedCasks...)

		// Sort for cleaner config
		sort.Strings(cfg.Brews)
		sort.Strings(cfg.Casks)

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Printf("✅ Added %d brews and %d casks to config\n", len(selectedBrews), len(selectedCasks))
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Println("   • View your config: dotfiles list")
		fmt.Println("   • Generate Brewfile: dotfiles brewfile")
	},
}

func selectPackages(packages []string) []string {
	fmt.Println("Options:")
	fmt.Println("  a - Add all")
	fmt.Println("  n - Add none")
	fmt.Println("  s - Select individually")
	fmt.Print("Choice (a/n/s): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	switch choice {
	case "a":
		fmt.Printf("  ✅ Adding all %d packages\n", len(packages))
		fmt.Println()
		return packages
	case "n":
		fmt.Println("  ⏭️  Skipping all packages")
		fmt.Println()
		return []string{}
	case "s":
		return selectIndividually(packages)
	default:
		fmt.Println("  Invalid choice, skipping all packages")
		fmt.Println()
		return []string{}
	}
}

func selectIndividually(packages []string) []string {
	fmt.Println()
	fmt.Println("For each package, press 'y' to add, 'n' to skip, 'q' to quit selection:")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	var selected []string

	for _, pkg := range packages {
		fmt.Printf("  Add '%s'? (y/n/q): ", pkg)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		switch response {
		case "y", "yes":
			selected = append(selected, pkg)
			fmt.Println("    ✅ Added")
		case "q", "quit":
			fmt.Println("    ⏹️  Quitting selection")
			fmt.Println()
			return selected
		default:
			fmt.Println("    ⏭️  Skipped")
		}
	}

	fmt.Println()
	return selected
}

func init() {
	scanCmd.Flags().Bool("auto", false, "Automatically add all installed packages without prompting")
	scanCmd.Flags().Bool("brews-only", false, "Only scan Homebrew formulas")
	scanCmd.Flags().Bool("casks-only", false, "Only scan Homebrew casks")

	rootCmd.AddCommand(scanCmd)
}
