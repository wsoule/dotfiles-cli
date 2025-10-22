package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:     "lint",
	GroupID: "advanced",
	Short:   "🔍 Validate configuration file",
	Long: `🔍 Configuration Linting

Validate your config.json for:
• Valid JSON syntax
• Duplicate packages
• Missing stow directories
• Circular dependencies
• Invalid conditional patterns

Examples:
  dotfiles lint                    # Validate config
  dotfiles lint --strict          # Enable strict mode
  dotfiles lint --fix             # Auto-fix issues`,
	Run: func(cmd *cobra.Command, args []string) {
		strict, _ := cmd.Flags().GetBool("strict")
		autoFix, _ := cmd.Flags().GetBool("fix")

		configPath := GetConfigPath(cmd)
		dotfilesDir := GetDotfilesDir(cmd)

		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			fmt.Println("\n💡 Make sure config.json has valid JSON syntax")
			return
		}

		fmt.Println("🔍 Linting configuration...")
		fmt.Println()

		issues := 0
		warnings := 0

		// Check for duplicate packages
		if dupes := findDuplicates(cfg.Brews); len(dupes) > 0 {
			fmt.Println("❌ Duplicate brews found:")
			for pkg, count := range dupes {
				fmt.Printf("   • %s (appears %d times)\n", pkg, count)
			}
			issues++
			if autoFix {
				cfg.Brews = removeDuplicates(cfg.Brews)
				fmt.Println("   ✅ Duplicates removed")
			}
		}

		if dupes := findDuplicates(cfg.Casks); len(dupes) > 0 {
			fmt.Println("❌ Duplicate casks found:")
			for pkg, count := range dupes {
				fmt.Printf("   • %s (appears %d times)\n", pkg, count)
			}
			issues++
			if autoFix {
				cfg.Casks = removeDuplicates(cfg.Casks)
				fmt.Println("   ✅ Duplicates removed")
			}
		}

		if dupes := findDuplicates(cfg.Stow); len(dupes) > 0 {
			fmt.Println("❌ Duplicate stow packages found:")
			for pkg, count := range dupes {
				fmt.Printf("   • %s (appears %d times)\n", pkg, count)
			}
			issues++
			if autoFix {
				cfg.Stow = removeDuplicates(cfg.Stow)
				fmt.Println("   ✅ Duplicates removed")
			}
		}

		// Check for missing stow directories
		stowDir := filepath.Join(dotfilesDir, "stow")
		for _, pkg := range cfg.Stow {
			pkgPath := filepath.Join(stowDir, pkg)
			if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
				fmt.Printf("⚠️  Stow package directory missing: %s\n", pkg)
				fmt.Printf("   Expected: %s\n", pkgPath)
				warnings++
			}
		}

		// Check for circular dependencies
		if cfg.PackageDependencies != nil {
			for pkg := range cfg.PackageDependencies {
				if _, err := resolveDependencies(cfg, pkg); err != nil {
					fmt.Printf("❌ Dependency error in '%s': %v\n", pkg, err)
					issues++
				}
			}
		}

		// Check for empty values
		emptyBrews := []string{}
		for _, brew := range cfg.Brews {
			if strings.TrimSpace(brew) == "" {
				emptyBrews = append(emptyBrews, brew)
			}
		}
		if len(emptyBrews) > 0 {
			fmt.Printf("❌ Found %d empty brew entries\n", len(emptyBrews))
			issues++
			if autoFix {
				cfg.Brews = removeEmpty(cfg.Brews)
				fmt.Println("   ✅ Empty entries removed")
			}
		}

		emptyCasks := []string{}
		for _, cask := range cfg.Casks {
			if strings.TrimSpace(cask) == "" {
				emptyCasks = append(emptyCasks, cask)
			}
		}
		if len(emptyCasks) > 0 {
			fmt.Printf("❌ Found %d empty cask entries\n", len(emptyCasks))
			issues++
			if autoFix {
				cfg.Casks = removeEmpty(cfg.Casks)
				fmt.Println("   ✅ Empty entries removed")
			}
		}

		// Strict mode checks
		if strict {
			// Check for packages in both brews and casks
			brewSet := make(map[string]bool)
			for _, b := range cfg.Brews {
				brewSet[b] = true
			}
			for _, c := range cfg.Casks {
				if brewSet[c] {
					fmt.Printf("⚠️  Package in both brews and casks: %s\n", c)
					warnings++
				}
			}

			// Check for undefined variables in conditionals
			if cfg.Variables != nil {
				for _, cond := range cfg.Conditionals {
					for key := range cond.Variables {
						if _, ok := cfg.Variables[key]; !ok {
							fmt.Printf("⚠️  Conditional defines variable not in main config: %s\n", key)
							warnings++
						}
					}
				}
			}
		}

		// Check role definitions
		if cfg.Roles != nil {
			for name, role := range cfg.Roles {
				if role.Name == "" {
					fmt.Printf("⚠️  Role '%s' has no name field\n", name)
					warnings++
				}
				if role.Description == "" {
					fmt.Printf("⚠️  Role '%s' has no description\n", name)
					warnings++
				}
			}
		}

		// Save fixes if applied
		if autoFix && issues > 0 {
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("\n❌ Failed to save fixes: %v\n", err)
			} else {
				fmt.Println("\n✅ Auto-fixed issues saved to config.json")
			}
		}

		// Summary
		fmt.Println()
		fmt.Println(strings.Repeat("─", 50))
		if issues == 0 && warnings == 0 {
			fmt.Println("✅ Configuration is valid!")
		} else {
			if issues > 0 {
				fmt.Printf("❌ Found %d issue(s)\n", issues)
			}
			if warnings > 0 {
				fmt.Printf("⚠️  Found %d warning(s)\n", warnings)
			}
			if !autoFix && issues > 0 {
				fmt.Println("\n💡 Run with --fix to auto-fix issues")
			}
		}
	},
}

func findDuplicates(items []string) map[string]int {
	counts := make(map[string]int)
	dupes := make(map[string]int)

	for _, item := range items {
		counts[item]++
		if counts[item] > 1 {
			dupes[item] = counts[item]
		}
	}

	return dupes
}

func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func removeEmpty(items []string) []string {
	result := []string{}
	for _, item := range items {
		if strings.TrimSpace(item) != "" {
			result = append(result, item)
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().Bool("strict", false, "Enable strict validation")
	lintCmd.Flags().Bool("fix", false, "Auto-fix issues")
}
