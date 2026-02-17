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
	Use: "list",
	GroupID: "package",
	Short: "List all packages in your configuration",
	Long:    `Shows all brews, casks, and taps in your config.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("error loading configuration: %w", err)
		}

		// Handle JSON output
		if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
			data, err := json.MarshalIndent(cfg, "", "")
			if err != nil {
				return fmt.Errorf("error marshaling json: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Handle count-only output
		if countOnly, _ := cmd.Flags().GetBool("count"); countOnly {
			fmt.Printf("Taps: %d\n", len(cfg.Taps))
			fmt.Printf("Brews: %d\n", len(cfg.Brews))
			fmt.Printf("Casks: %d\n", len(cfg.Casks))
			fmt.Printf("Stow: %d\n", len(cfg.Stow))
			fmt.Printf("Total: %d\n", len(cfg.Taps)+len(cfg.Brews)+len(cfg.Casks)+len(cfg.Stow))
			return nil
		}

		// Handle type-specific filtering
		packageType, _ := cmd.Flags().GetString("type")
		switch packageType {
		case "tap":
			for _, tap := range cfg.Taps {
				fmt.Println(tap)
			}
			return nil
		case "brew":
			for _, brew := range cfg.Brews {
				fmt.Println(brew)
			}
			return nil
		case "cask":
			for _, cask := range cfg.Casks {
				fmt.Println(cask)
			}
			return nil
		case "stow":
			for _, stow := range cfg.Stow {
				fmt.Println(stow)
			}
			return nil
		}

		// Default formatted output
		if len(cfg.Taps) > 0 {
			fmt.Println(" Taps:")
			for _, tap := range cfg.Taps {
				fmt.Printf("- %s\n", tap)
			}
			fmt.Println()
		}

		if len(cfg.Brews) > 0 {
			fmt.Println(" Brews:")
			for _, brew := range cfg.Brews {
				fmt.Printf("- %s\n", brew)
			}
			fmt.Println()
		}

		if len(cfg.Casks) > 0 {
			fmt.Println(" Casks:")
			for _, cask := range cfg.Casks {
				fmt.Printf("- %s\n", cask)
			}
			fmt.Println()
		}

		if len(cfg.Stow) > 0 {
			fmt.Println(" Stow Packages:")
			for _, stow := range cfg.Stow {
				fmt.Printf("- %s\n", stow)
			}
		}

		total := len(cfg.Taps) + len(cfg.Brews) + len(cfg.Casks) + len(cfg.Stow)
		if total == 0 {
			fmt.Println("No packages configured. Run 'dotfiles add <package>' to get started.")
		} else {
			fmt.Printf("\nTotal packages: %d\n", total)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().Bool("json", false, "Output as JSON")
	listCmd.Flags().Bool("count", false, "Show only package counts")
	listCmd.Flags().StringP("type", "t", "", "Filter by package type (brew, cask, tap, stow)")
	rootCmd.AddCommand(listCmd)
}
