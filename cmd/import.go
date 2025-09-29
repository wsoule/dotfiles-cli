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

var importCmd = &cobra.Command{
	Use:   "import <brewfile>",
	Short: "Import packages from a Brewfile",
	Long:  `Parse a Brewfile and add packages to your JSON configuration`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		brewfilePath := args[0]

		// Check if Brewfile exists
		if _, err := os.Stat(brewfilePath); os.IsNotExist(err) {
			fmt.Printf("Brewfile not found: %s\n", brewfilePath)
			os.Exit(1)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Load existing config or create new one
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Parse Brewfile
		brews, casks, taps, err := parseBrewfile(brewfilePath)
		if err != nil {
			fmt.Printf("Error parsing Brewfile: %v\n", err)
			os.Exit(1)
		}

		// Merge with existing config
		merge := !cmd.Flags().Changed("replace")

		if merge {
			// Add to existing packages
			cfg.Brews = mergeSlices(cfg.Brews, brews)
			cfg.Casks = mergeSlices(cfg.Casks, casks)
			cfg.Taps = mergeSlices(cfg.Taps, taps)
		} else {
			// Replace existing packages
			cfg.Brews = brews
			cfg.Casks = casks
			cfg.Taps = taps
		}

		// Save updated config
		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		action := "Merged"
		if !merge {
			action = "Imported"
		}

		fmt.Printf("✓ %s packages from %s:\n", action, brewfilePath)
		if len(taps) > 0 {
			fmt.Printf("  📋 %d taps\n", len(taps))
		}
		if len(brews) > 0 {
			fmt.Printf("  🍺 %d brews\n", len(brews))
		}
		if len(casks) > 0 {
			fmt.Printf("  📦 %d casks\n", len(casks))
		}
	},
}

func parseBrewfile(path string) ([]string, []string, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	var brews, casks, taps []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse different line types
		if strings.HasPrefix(line, "tap ") {
			tap := extractQuotedValue(line, "tap")
			if tap != "" {
				taps = append(taps, tap)
			}
		} else if strings.HasPrefix(line, "brew ") {
			brew := extractQuotedValue(line, "brew")
			if brew != "" {
				brews = append(brews, brew)
			}
		} else if strings.HasPrefix(line, "cask ") {
			cask := extractQuotedValue(line, "cask")
			if cask != "" {
				casks = append(casks, cask)
			}
		}
	}

	return brews, casks, taps, scanner.Err()
}

func extractQuotedValue(line, prefix string) string {
	// Remove prefix and find quoted value
	content := strings.TrimPrefix(line, prefix+" ")
	content = strings.TrimSpace(content)

	// Handle quoted values
	if strings.HasPrefix(content, "\"") && strings.HasSuffix(content, "\"") {
		return strings.Trim(content, "\"")
	} else if strings.HasPrefix(content, "'") && strings.HasSuffix(content, "'") {
		return strings.Trim(content, "'")
	}

	// Handle unquoted values (take first word)
	parts := strings.Fields(content)
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

func mergeSlices(existing, new []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(existing)+len(new))

	// Add existing items
	for _, item := range existing {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	// Add new items
	for _, item := range new {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

func init() {
	importCmd.Flags().Bool("replace", false, "Replace existing packages instead of merging")
	rootCmd.AddCommand(importCmd)
}