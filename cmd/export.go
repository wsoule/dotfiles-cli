package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

type MachineProfile struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Machine     string            `json:"machine"`
	Platform    string            `json:"platform"`
	CreatedAt   string            `json:"created_at"`
	Config      *config.Config    `json:"config"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

var exportCmd = &cobra.Command{
	Use:   "export <profile-name>",
	Short: "📤 Export current configuration as a machine-specific profile",
	Long: `📤 Export Configuration Profile

Create machine-specific configuration profiles for different setups (work, personal, etc.).
Profiles can be imported later to quickly configure new machines.

Examples:
  dotfiles export work-mac                            # Export current config as "work-mac"
  dotfiles export personal --description="Home setup" # With description
  dotfiles export minimal --brews-only                # Only export brews
  dotfiles export full --output=~/my-profile.json     # Custom output location`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		description, _ := cmd.Flags().GetString("description")
		output, _ := cmd.Flags().GetString("output")
		brewsOnly, _ := cmd.Flags().GetBool("brews-only")
		casksOnly, _ := cmd.Flags().GetBool("casks-only")
		machine, _ := cmd.Flags().GetString("machine")

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

		// Filter config based on flags
		exportCfg := &config.Config{
			Brews: cfg.Brews,
			Casks: cfg.Casks,
			Taps:  cfg.Taps,
			Stow:  cfg.Stow,
		}

		if brewsOnly {
			exportCfg.Casks = []string{}
			exportCfg.Stow = []string{}
		}

		if casksOnly {
			exportCfg.Brews = []string{}
			exportCfg.Stow = []string{}
		}

		// Create profile
		profile := MachineProfile{
			Name:        profileName,
			Description: description,
			Machine:     machine,
			Platform:    "darwin", // TODO: detect platform
			CreatedAt:   time.Now().Format(time.RFC3339),
			Config:      exportCfg,
			Metadata: map[string]string{
				"exported_from": "dotfiles CLI",
				"version":       "1.0",
			},
		}

		// Determine output path
		if output == "" {
			profilesDir := filepath.Join(home, ".dotfiles", "profiles")
			os.MkdirAll(profilesDir, 0755)
			output = filepath.Join(profilesDir, profileName+".json")
		}

		// Export profile
		data, err := json.MarshalIndent(profile, "", "  ")
		if err != nil {
			fmt.Printf("❌ Error marshaling profile: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(output, data, 0644); err != nil {
			fmt.Printf("❌ Error writing profile: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📤 Profile exported successfully!")
		fmt.Println()
		fmt.Printf("   Name: %s\n", profileName)
		if description != "" {
			fmt.Printf("   Description: %s\n", description)
		}
		fmt.Printf("   Location: %s\n", output)
		fmt.Println()
		fmt.Println("📊 Profile contents:")
		fmt.Printf("   • %d brews\n", len(exportCfg.Brews))
		fmt.Printf("   • %d casks\n", len(exportCfg.Casks))
		fmt.Printf("   • %d taps\n", len(exportCfg.Taps))
		fmt.Printf("   • %d stow packages\n", len(exportCfg.Stow))
		fmt.Println()
		fmt.Println("💡 To import this profile:")
		fmt.Printf("   dotfiles import-profile %s\n", output)
	},
}

var importProfileCmd = &cobra.Command{
	Use:   "import-profile <profile-file>",
	Short: "📥 Import a machine-specific profile",
	Long: `📥 Import Configuration Profile

Import a previously exported machine profile.
You can merge it with your current config or replace it entirely.

Examples:
  dotfiles import-profile work-mac.json              # Import and merge
  dotfiles import-profile work-mac.json --replace    # Replace current config
  dotfiles import-profile work-mac.json --install    # Import and install packages`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profilePath := args[0]
		replace, _ := cmd.Flags().GetBool("replace")
		install, _ := cmd.Flags().GetBool("install")

		// Read profile
		data, err := os.ReadFile(profilePath)
		if err != nil {
			fmt.Printf("❌ Error reading profile: %v\n", err)
			os.Exit(1)
		}

		var profile MachineProfile
		if err := json.Unmarshal(data, &profile); err != nil {
			fmt.Printf("❌ Error parsing profile: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📥 Importing profile...")
		fmt.Println()
		fmt.Printf("   Name: %s\n", profile.Name)
		if profile.Description != "" {
			fmt.Printf("   Description: %s\n", profile.Description)
		}
		fmt.Printf("   Created: %s\n", profile.CreatedAt)
		fmt.Println()

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		if replace {
			// Replace entire config
			if err := profile.Config.Save(configPath); err != nil {
				fmt.Printf("❌ Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Configuration replaced")
		} else {
			// Merge with existing config
			cfg, err := config.Load(configPath)
			if err != nil {
				// If config doesn't exist, use profile config
				cfg = profile.Config
			} else {
				// Merge
				cfg.Brews = mergeUnique(cfg.Brews, profile.Config.Brews)
				cfg.Casks = mergeUnique(cfg.Casks, profile.Config.Casks)
				cfg.Taps = mergeUnique(cfg.Taps, profile.Config.Taps)
				cfg.Stow = mergeUnique(cfg.Stow, profile.Config.Stow)
			}

			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("❌ Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Profile merged with existing configuration")
		}

		fmt.Println()
		fmt.Println("💡 Next steps:")
		if install {
			fmt.Println("   Installing packages...")
			// TODO: Call install command
		} else {
			fmt.Println("   • View config: dotfiles list")
			fmt.Println("   • Install packages: dotfiles install")
		}
	},
}

var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "📋 List all exported profiles",
	Long:  `📋 List all machine-specific profiles stored in ~/.dotfiles/profiles/`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		profilesDir := filepath.Join(home, ".dotfiles", "profiles")
		if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
			fmt.Println("📋 No profiles found")
			fmt.Println()
			fmt.Println("💡 Create a profile:")
			fmt.Println("   dotfiles export <profile-name>")
			return
		}

		entries, err := os.ReadDir(profilesDir)
		if err != nil {
			fmt.Printf("❌ Error reading profiles directory: %v\n", err)
			os.Exit(1)
		}

		profiles := []MachineProfile{}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				profilePath := filepath.Join(profilesDir, entry.Name())
				data, err := os.ReadFile(profilePath)
				if err != nil {
					continue
				}

				var profile MachineProfile
				if err := json.Unmarshal(data, &profile); err != nil {
					continue
				}

				profiles = append(profiles, profile)
			}
		}

		if len(profiles) == 0 {
			fmt.Println("📋 No profiles found")
			return
		}

		fmt.Printf("📋 Found %d profile(s):\n", len(profiles))
		fmt.Println()

		for i, profile := range profiles {
			fmt.Printf("%d. %s\n", i+1, profile.Name)
			if profile.Description != "" {
				fmt.Printf("   Description: %s\n", profile.Description)
			}
			fmt.Printf("   Created: %s\n", profile.CreatedAt)
			fmt.Printf("   Packages: %d brews, %d casks, %d taps\n",
				len(profile.Config.Brews),
				len(profile.Config.Casks),
				len(profile.Config.Taps))
			fmt.Printf("   File: profiles/%s.json\n", profile.Name)
			fmt.Println()
		}
	},
}

func mergeUnique(a, b []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range a {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	for _, item := range b {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func init() {
	exportCmd.Flags().StringP("description", "d", "", "Profile description")
	exportCmd.Flags().StringP("output", "o", "", "Output file path (default: ~/.dotfiles/profiles/<name>.json)")
	exportCmd.Flags().Bool("brews-only", false, "Only export brew packages")
	exportCmd.Flags().Bool("casks-only", false, "Only export cask packages")
	exportCmd.Flags().StringP("machine", "m", "", "Machine identifier (e.g., 'work-macbook-pro')")

	importProfileCmd.Flags().Bool("replace", false, "Replace current config instead of merging")
	importProfileCmd.Flags().Bool("install", false, "Install packages after importing")

	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importProfileCmd)
	rootCmd.AddCommand(listProfilesCmd)
}
