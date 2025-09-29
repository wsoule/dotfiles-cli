package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all packages in your configuration",
	Long:  `Shows all brews, casks, and taps in your config.json`,
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

		// Handle JSON output
		if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
			data, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(data))
			return
		}

		// Handle count-only output
		if countOnly, _ := cmd.Flags().GetBool("count"); countOnly {
			fmt.Printf("Taps: %d\n", len(cfg.Taps))
			fmt.Printf("Brews: %d\n", len(cfg.Brews))
			fmt.Printf("Casks: %d\n", len(cfg.Casks))
			fmt.Printf("Stow: %d\n", len(cfg.Stow))
			fmt.Printf("Total: %d\n", len(cfg.Taps)+len(cfg.Brews)+len(cfg.Casks)+len(cfg.Stow))
			return
		}

		// Handle type-specific filtering
		packageType, _ := cmd.Flags().GetString("type")
		switch packageType {
		case "tap":
			for _, tap := range cfg.Taps {
				fmt.Println(tap)
			}
			return
		case "brew":
			for _, brew := range cfg.Brews {
				fmt.Println(brew)
			}
			return
		case "cask":
			for _, cask := range cfg.Casks {
				fmt.Println(cask)
			}
			return
		case "stow":
			for _, stow := range cfg.Stow {
				fmt.Println(stow)
			}
			return
		}

		// Default formatted output
		if len(cfg.Taps) > 0 {
			fmt.Println("📋 Taps:")
			for _, tap := range cfg.Taps {
				fmt.Printf("  - %s\n", tap)
			}
			fmt.Println()
		}

		if len(cfg.Brews) > 0 {
			fmt.Println("🍺 Brews:")
			for _, brew := range cfg.Brews {
				fmt.Printf("  - %s\n", brew)
			}
			fmt.Println()
		}

		if len(cfg.Casks) > 0 {
			fmt.Println("📦 Casks:")
			for _, cask := range cfg.Casks {
				fmt.Printf("  - %s\n", cask)
			}
			fmt.Println()
		}

		if len(cfg.Stow) > 0 {
			fmt.Println("🔗 Stow Packages:")
			for _, stow := range cfg.Stow {
				fmt.Printf("  - %s\n", stow)
			}
		}

		total := len(cfg.Taps) + len(cfg.Brews) + len(cfg.Casks) + len(cfg.Stow)
		if total == 0 {
			fmt.Println("No packages configured. Run 'dotfiles add <package>' to get started.")
		} else {
			fmt.Printf("\nTotal packages: %d\n", total)
		}
	},
}

func init() {
	listCmd.Flags().Bool("json", false, "Output as JSON")
	listCmd.Flags().Bool("count", false, "Show only package counts")
	listCmd.Flags().StringP("type", "t", "", "Filter by package type (brew, cask, tap, stow)")
	rootCmd.AddCommand(listCmd)
}