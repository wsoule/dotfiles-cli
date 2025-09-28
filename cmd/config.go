package cmd

import (
	"encoding/json"
	"fmt"

	"dotfiles/internal/config"
	"dotfiles/internal/ui"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage dotfiles configuration",
	Long: `Manage your dotfiles configuration settings.

Subcommands:
  show    - Display current configuration
  edit    - Edit configuration interactively
  get     - Get a specific configuration value
  set     - Set a specific configuration value
  validate - Validate configuration`,
}

// configShowCmd shows the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	RunE:  runConfigShow,
}

// configEditCmd edits configuration interactively
var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration interactively",
	RunE:  runConfigEdit,
}

// configGetCmd gets a configuration value
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a specific configuration value",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

// configSetCmd sets a configuration value
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a specific configuration value",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

// configValidateCmd validates the configuration
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	RunE:  runConfigValidate,
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Add subcommands
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configValidateCmd)

	// Flags
	configShowCmd.Flags().Bool("json", false, "Output in JSON format")
	configShowCmd.Flags().Bool("summary", false, "Show summary only")
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	configManager := config.NewManager(cfgFile)
	if err := configManager.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg := configManager.Get()
	if cfg == nil {
		return fmt.Errorf("no configuration found")
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	summary, _ := cmd.Flags().GetBool("summary")

	if jsonOutput {
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal configuration: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if summary {
		showConfigSummary(cfg)
	} else {
		showDetailedConfig(cfg)
	}

	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	fmt.Println("🔧 Configuration Editor")
	fmt.Println("=======================")
	fmt.Println()
	fmt.Println("This will launch an interactive configuration editor.")
	fmt.Println("⚠️  Not implemented yet. Use 'dotfiles setup --force' to reconfigure.")

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	configManager := config.NewManager(cfgFile)
	if err := configManager.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg := configManager.Get()
	if cfg == nil {
		return fmt.Errorf("no configuration found")
	}

	// Convert to JSON and use jq-like path access
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	value := getValueByPath(jsonData, key)
	if value != nil {
		fmt.Println(value)
	} else {
		return fmt.Errorf("key not found: %s", key)
	}

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	fmt.Printf("Setting %s = %s\n", key, value)
	fmt.Println("⚠️  Not implemented yet. Use 'dotfiles setup --force' to reconfigure.")

	return nil
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	configManager := config.NewManager(cfgFile)
	if err := configManager.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := configManager.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	ui.PrintSuccess("Configuration is valid")
	return nil
}

func showConfigSummary(cfg *config.Config) {
	fmt.Println("📄 Configuration Summary")
	fmt.Println("========================")
	fmt.Printf("👤 User: %s\n", cfg.Personal.Name)
	fmt.Printf("📧 Email: %s\n", cfg.Personal.Email)
	fmt.Printf("✏️  Editor: %s\n", cfg.Personal.Editor)
	fmt.Printf("🎨 Dark Mode: %t\n", cfg.System.Appearance.DarkMode)
	fmt.Printf("🚢 Dock Position: %s\n", cfg.System.Dock.Position)
	fmt.Printf("🐚 Shell Theme: %s\n", cfg.Development.Shell.Theme)
	fmt.Printf("🌿 Git Branch: %s\n", cfg.Development.Git.DefaultBranch)
	fmt.Println()

	// Count enabled languages
	enabledLangs := 0
	for _, enabled := range cfg.Development.Languages {
		if enabled {
			enabledLangs++
		}
	}
	fmt.Printf("🌐 Programming Languages: %d enabled\n", enabledLangs)

	// Count enabled frameworks
	enabledFrameworks := 0
	for _, enabled := range cfg.Development.Frameworks {
		if enabled {
			enabledFrameworks++
		}
	}
	fmt.Printf("🚀 Frameworks: %d enabled\n", enabledFrameworks)

	fmt.Printf("📦 Extra Packages: %d brews, %d casks\n",
		len(cfg.Packages.ExtraBrews), len(cfg.Packages.ExtraCasks))
}

func showDetailedConfig(cfg *config.Config) {
	fmt.Println("📄 Detailed Configuration")
	fmt.Println("==========================")
	fmt.Println()

	fmt.Println("👤 Personal Information:")
	fmt.Printf("  Name: %s\n", cfg.Personal.Name)
	fmt.Printf("  Email: %s\n", cfg.Personal.Email)
	fmt.Printf("  Editor: %s\n", cfg.Personal.Editor)
	fmt.Println()

	fmt.Println("🎨 System Settings:")
	fmt.Printf("  Dark Mode: %t\n", cfg.System.Appearance.DarkMode)
	fmt.Printf("  24-Hour Time: %t\n", cfg.System.Appearance.Enable24HourTime)
	fmt.Printf("  Dock Auto-hide: %t\n", cfg.System.Dock.Autohide)
	fmt.Printf("  Dock Position: %s\n", cfg.System.Dock.Position)
	fmt.Printf("  Dock Size: %d\n", cfg.System.Dock.TileSize)
	fmt.Println()

	fmt.Println("🌐 Programming Languages:")
	for lang, enabled := range cfg.Development.Languages {
		if enabled {
			fmt.Printf("  ✅ %s\n", lang)
		}
	}
	fmt.Println()

	fmt.Println("🚀 Frameworks:")
	for framework, enabled := range cfg.Development.Frameworks {
		if enabled {
			fmt.Printf("  ✅ %s\n", framework)
		}
	}
	fmt.Println()

	fmt.Println("🛠️  Development Tools:")
	for tool, enabled := range cfg.Development.Tools {
		if enabled {
			fmt.Printf("  ✅ %s\n", tool)
		}
	}
	fmt.Println()

	if len(cfg.Packages.ExtraBrews) > 0 {
		fmt.Println("📦 Extra Brew Packages:")
		for _, pkg := range cfg.Packages.ExtraBrews {
			fmt.Printf("  • %s\n", pkg)
		}
		fmt.Println()
	}

	if len(cfg.Packages.ExtraCasks) > 0 {
		fmt.Println("📱 Extra Applications:")
		for _, pkg := range cfg.Packages.ExtraCasks {
			fmt.Printf("  • %s\n", pkg)
		}
		fmt.Println()
	}

	fmt.Printf("📊 Metadata:\n")
	fmt.Printf("  Version: %s\n", cfg.Metadata.Version)
	fmt.Printf("  Created: %s\n", cfg.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Modified: %s\n", cfg.Metadata.LastModified.Format("2006-01-02 15:04:05"))
}

// Helper function to get value by dot-notation path
func getValueByPath(data interface{}, path string) interface{} {
	// Simple implementation for basic paths
	// TODO: Implement full dot-notation path parsing
	return nil
}
