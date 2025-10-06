package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "📊 Show differences between config and installed packages",
	Long: `📊 Configuration Diff

Compare your dotfiles configuration with actually installed packages.
Shows what's missing, what's extra, and what's in sync.

Examples:
  dotfiles diff                # Show all differences
  dotfiles diff --type=brews   # Only show brew differences
  dotfiles diff --type=casks   # Only show cask differences
  dotfiles diff --verbose      # Show all packages including synced ones`,
	Run: func(cmd *cobra.Command, args []string) {
		pkgType, _ := cmd.Flags().GetString("type")
		verbose, _ := cmd.Flags().GetBool("verbose")

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📊 Configuration Diff")
		fmt.Println("=" + strings.Repeat("=", 21))
		fmt.Println()

		hasDiff := false

		// Diff brews
		if pkgType == "" || pkgType == "brews" || pkgType == "brew" {
			diffBrews := getDiffBrews(cfg.Brews, verbose)
			if diffBrews.HasDiff || verbose {
				fmt.Println("🍺 Homebrew Formulas:")
				printDiff(diffBrews, verbose)
				fmt.Println()
				hasDiff = hasDiff || diffBrews.HasDiff
			}
		}

		// Diff casks
		if pkgType == "" || pkgType == "casks" || pkgType == "cask" {
			diffCasks := getDiffCasks(cfg.Casks, verbose)
			if diffCasks.HasDiff || verbose {
				fmt.Println("📦 Homebrew Casks:")
				printDiff(diffCasks, verbose)
				fmt.Println()
				hasDiff = hasDiff || diffCasks.HasDiff
			}
		}

		// Summary
		if !hasDiff && !verbose {
			fmt.Println("✅ No differences found - everything is in sync!")
		} else if hasDiff {
			fmt.Println("💡 Suggested actions:")
			fmt.Println("   • Install missing packages: dotfiles install")
			fmt.Println("   • Add untracked packages: dotfiles scan")
			fmt.Println("   • Remove from config: dotfiles remove <package>")
		}
	},
}

type PackageDiff struct {
	Missing  []string // In config but not installed
	Extra    []string // Installed but not in config
	Synced   []string // Both in config and installed
	HasDiff  bool
}

func getDiffBrews(configured []string, verbose bool) PackageDiff {
	diff := PackageDiff{
		Missing: []string{},
		Extra:   []string{},
		Synced:  []string{},
	}

	installed, err := getInstalledBrews()
	if err != nil {
		fmt.Printf("⚠️  Could not get installed brews: %v\n", err)
		return diff
	}

	installedMap := make(map[string]bool)
	for _, pkg := range installed {
		installedMap[pkg] = true
	}

	configuredMap := make(map[string]bool)
	for _, pkg := range configured {
		configuredMap[pkg] = true
	}

	// Find missing (in config but not installed)
	for _, pkg := range configured {
		if !installedMap[pkg] {
			diff.Missing = append(diff.Missing, pkg)
			diff.HasDiff = true
		} else {
			diff.Synced = append(diff.Synced, pkg)
		}
	}

	// Find extra (installed but not in config)
	for _, pkg := range installed {
		if !configuredMap[pkg] {
			diff.Extra = append(diff.Extra, pkg)
			diff.HasDiff = true
		}
	}

	sort.Strings(diff.Missing)
	sort.Strings(diff.Extra)
	sort.Strings(diff.Synced)

	return diff
}

func getDiffCasks(configured []string, verbose bool) PackageDiff {
	diff := PackageDiff{
		Missing: []string{},
		Extra:   []string{},
		Synced:  []string{},
	}

	installed, err := getInstalledCasks()
	if err != nil {
		fmt.Printf("⚠️  Could not get installed casks: %v\n", err)
		return diff
	}

	installedMap := make(map[string]bool)
	for _, pkg := range installed {
		installedMap[pkg] = true
	}

	configuredMap := make(map[string]bool)
	for _, pkg := range configured {
		configuredMap[pkg] = true
	}

	// Find missing (in config but not installed)
	for _, pkg := range configured {
		if !installedMap[pkg] {
			diff.Missing = append(diff.Missing, pkg)
			diff.HasDiff = true
		} else {
			diff.Synced = append(diff.Synced, pkg)
		}
	}

	// Find extra (installed but not in config)
	for _, pkg := range installed {
		if !configuredMap[pkg] {
			diff.Extra = append(diff.Extra, pkg)
			diff.HasDiff = true
		}
	}

	sort.Strings(diff.Missing)
	sort.Strings(diff.Extra)
	sort.Strings(diff.Synced)

	return diff
}

func printDiff(diff PackageDiff, verbose bool) {
	if len(diff.Missing) > 0 {
		fmt.Printf("  ❌ Missing (%d) - in config but not installed:\n", len(diff.Missing))
		for _, pkg := range diff.Missing {
			fmt.Printf("     - %s\n", pkg)
		}
	}

	if len(diff.Extra) > 0 {
		fmt.Printf("  ⚠️  Extra (%d) - installed but not in config:\n", len(diff.Extra))
		for _, pkg := range diff.Extra {
			fmt.Printf("     + %s\n", pkg)
		}
	}

	if verbose && len(diff.Synced) > 0 {
		fmt.Printf("  ✅ Synced (%d) - in both config and installed:\n", len(diff.Synced))
		for _, pkg := range diff.Synced {
			fmt.Printf("     = %s\n", pkg)
		}
	}

	if !diff.HasDiff {
		fmt.Println("  ✅ All in sync")
	}
}

func init() {
	diffCmd.Flags().String("type", "", "Filter by type (brews, casks)")
	diffCmd.Flags().BoolP("verbose", "v", false, "Show all packages including synced ones")

	rootCmd.AddCommand(diffCmd)
}
