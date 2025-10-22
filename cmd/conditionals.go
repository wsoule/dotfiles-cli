package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var conditionalsCmd = &cobra.Command{
	Use:     "conditionals",
	GroupID: "advanced",
	Short:   "🔀 Manage OS/hostname-specific configurations",
	Long: `🔀 Conditional Configuration Loading

Load different packages based on:
• Hostname (exact match or glob pattern)
• Operating system (darwin, linux, etc.)
• Custom variables

Examples:
  dotfiles conditionals apply      # Auto-detect and apply
  dotfiles conditionals test       # Preview what would apply
  dotfiles conditionals add        # Add new conditional`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles conditionals apply' to apply conditional configs")
	},
}

var conditionalsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply conditional configurations",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if len(cfg.Conditionals) == 0 {
			fmt.Println("No conditional configurations defined")
			return
		}

		hostname, _ := os.Hostname()
		osType := runtime.GOOS

		fmt.Println("🔍 Checking conditionals...")
		fmt.Printf("   Hostname: %s\n", hostname)
		fmt.Printf("   OS: %s\n", osType)
		fmt.Println()

		applied := false
		for i, cond := range cfg.Conditionals {
			matches := false
			reason := ""

			// Check hostname
			if cond.Hostname != "" {
				if matchGlob(hostname, cond.Hostname) {
					matches = true
					reason = fmt.Sprintf("hostname matches '%s'", cond.Hostname)
				}
			}

			// Check OS
			if cond.OS != "" && cond.OS == osType {
				matches = true
				if reason != "" {
					reason += " and "
				}
				reason += fmt.Sprintf("OS matches '%s'", osType)
			}

			// If no conditions specified, skip
			if cond.Hostname == "" && cond.OS == "" {
				continue
			}

			if matches {
				fmt.Printf("✅ Conditional #%d matches (%s)\n", i+1, reason)
				fmt.Printf("   Brews: %v\n", cond.Brews)
				fmt.Printf("   Casks: %v\n", cond.Casks)

				if !dryRun {
					// Apply conditional packages
					cfg.Brews = append(cfg.Brews, cond.Brews...)
					cfg.Casks = append(cfg.Casks, cond.Casks...)

					// Merge variables
					if cond.Variables != nil {
						if cfg.Variables == nil {
							cfg.Variables = make(map[string]string)
						}
						for k, v := range cond.Variables {
							cfg.Variables[k] = v
						}
					}
					applied = true
				}
				fmt.Println()
			}
		}

		if !applied && !dryRun {
			fmt.Println("⚠️  No conditionals matched this machine")
			return
		}

		if dryRun {
			fmt.Println("🔍 Dry run - no changes made")
			fmt.Println("💡 Run without --dry-run to apply changes")
			return
		}

		// Save updated config
		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			return
		}

		fmt.Println("✅ Conditional configurations applied")
		fmt.Println("💡 Run 'dotfiles install' to install the packages")
	},
}

var conditionalsTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test which conditionals would apply",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		hostname, _ := os.Hostname()
		osType := runtime.GOOS

		fmt.Println("🔍 Testing conditionals...")
		fmt.Printf("   Current hostname: %s\n", hostname)
		fmt.Printf("   Current OS: %s\n", osType)
		fmt.Println()

		if len(cfg.Conditionals) == 0 {
			fmt.Println("No conditional configurations defined")
			fmt.Println()
			fmt.Println("💡 Add one with: dotfiles conditionals add")
			return
		}

		matches := 0
		for i, cond := range cfg.Conditionals {
			willMatch := false
			reason := ""

			if cond.Hostname != "" && matchGlob(hostname, cond.Hostname) {
				willMatch = true
				reason = fmt.Sprintf("hostname matches '%s'", cond.Hostname)
			}

			if cond.OS != "" && cond.OS == osType {
				willMatch = true
				if reason != "" {
					reason += " and "
				}
				reason += fmt.Sprintf("OS matches '%s'", osType)
			}

			if cond.Hostname == "" && cond.OS == "" {
				continue
			}

			if willMatch {
				matches++
				fmt.Printf("✅ Conditional #%d WOULD APPLY\n", i+1)
				fmt.Printf("   Reason: %s\n", reason)
				fmt.Printf("   Brews: %v\n", cond.Brews)
				fmt.Printf("   Casks: %v\n", cond.Casks)
				if cond.Variables != nil {
					fmt.Printf("   Variables: %v\n", cond.Variables)
				}
			} else {
				fmt.Printf("⚪ Conditional #%d would NOT apply\n", i+1)
				if cond.Hostname != "" {
					fmt.Printf("   Hostname pattern: %s\n", cond.Hostname)
				}
				if cond.OS != "" {
					fmt.Printf("   OS: %s\n", cond.OS)
				}
			}
			fmt.Println()
		}

		if matches == 0 {
			fmt.Println("⚠️  No conditionals would apply to this machine")
		} else {
			fmt.Printf("📊 %d conditional(s) would apply\n", matches)
		}
	},
}

var conditionalsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new conditional configuration",
	Run: func(cmd *cobra.Command, args []string) {
		hostname, _ := cmd.Flags().GetString("hostname")
		osType, _ := cmd.Flags().GetString("os")
		brews, _ := cmd.Flags().GetStringSlice("brews")
		casks, _ := cmd.Flags().GetStringSlice("casks")

		if hostname == "" && osType == "" {
			fmt.Println("❌ Must specify at least --hostname or --os")
			return
		}

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		newCond := config.ConditionalConfig{
			Hostname: hostname,
			OS:       osType,
			Brews:    brews,
			Casks:    casks,
		}

		cfg.Conditionals = append(cfg.Conditionals, newCond)

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			return
		}

		fmt.Println("✅ Conditional configuration added")
		if hostname != "" {
			fmt.Printf("   Hostname: %s\n", hostname)
		}
		if osType != "" {
			fmt.Printf("   OS: %s\n", osType)
		}
		fmt.Printf("   Packages: %d brews, %d casks\n", len(brews), len(casks))
	},
}

// matchGlob performs simple glob matching
func matchGlob(str, pattern string) bool {
	// Simple glob: only supports * wildcard
	if pattern == "*" {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return str == pattern
	}

	// Convert glob to prefix/suffix matching
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		// *substring*
		substr := pattern[1 : len(pattern)-1]
		return strings.Contains(str, substr)
	} else if strings.HasPrefix(pattern, "*") {
		// *suffix
		suffix := pattern[1:]
		return strings.HasSuffix(str, suffix)
	} else if strings.HasSuffix(pattern, "*") {
		// prefix*
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(str, prefix)
	}

	return str == pattern
}

func init() {
	rootCmd.AddCommand(conditionalsCmd)
	conditionalsCmd.AddCommand(conditionalsApplyCmd)
	conditionalsCmd.AddCommand(conditionalsTestCmd)
	conditionalsCmd.AddCommand(conditionalsAddCmd)

	conditionalsApplyCmd.Flags().Bool("dry-run", false, "Preview changes without applying")

	conditionalsAddCmd.Flags().String("hostname", "", "Hostname pattern (supports * wildcard)")
	conditionalsAddCmd.Flags().String("os", "", "Operating system (darwin, linux, etc.)")
	conditionalsAddCmd.Flags().StringSlice("brews", []string{}, "Brew packages to install")
	conditionalsAddCmd.Flags().StringSlice("casks", []string{}, "Cask packages to install")
}
