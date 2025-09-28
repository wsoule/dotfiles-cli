package cmd

import (
	"fmt"

	"dotfiles/internal/config"
	"dotfiles/internal/ui"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard for dotfiles configuration",
	Long: `Launch an interactive setup wizard to configure your dotfiles.

This wizard will guide you through:
- Choosing from preset configurations
- Setting personal information (name, email)
- Configuring system preferences (dark mode, dock settings)
- Selecting development tools and languages
- Customizing additional packages

The configuration will be saved and can be modified later.`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().Bool("force", false, "Force setup even if configuration already exists")
	setupCmd.Flags().String("preset", "", "Use a specific preset configuration")
	setupCmd.Flags().Bool("quick", false, "Quick setup with minimal prompts")
}

func runSetup(cmd *cobra.Command, args []string) error {
	ui.Banner()
	ui.PrintSection("Setup Wizard")

	ui.PrintInfo("This wizard will help you:")
	fmt.Println("  • Choose from preset configurations")
	fmt.Println("  • Customize settings to your preferences")
	fmt.Println("  • Generate a personalized environment")
	fmt.Println()

	// Create configuration manager
	configManager := config.NewManager(cfgFile)

	// Check if config exists and force flag
	force, _ := cmd.Flags().GetBool("force")
	if !force {
		if err := configManager.Load(); err == nil {
			ui.PrintWarning("Configuration already exists!")
			fmt.Println("    Use --force to overwrite or run 'dotfiles config' to modify existing settings.")
			return nil
		}
	}

	// Check for preset flag
	preset, _ := cmd.Flags().GetString("preset")
	quick, _ := cmd.Flags().GetBool("quick")

	// Create UI manager
	uiManager := ui.NewManager()

	// Load or create configuration
	var cfg *config.Config
	var err error

	if preset != "" {
		// Load preset configuration
		cfg, err = loadPresetConfig(preset)
		if err != nil {
			return fmt.Errorf("failed to load preset '%s': %w", preset, err)
		}
		ui.PrintSuccess(fmt.Sprintf("Loaded preset configuration: %s", preset))
	} else {
		// Start interactive setup
		cfg, err = uiManager.RunSetupWizard(quick)
		if err != nil {
			return fmt.Errorf("setup wizard failed: %w", err)
		}
	}

	// Set and save configuration
	configManager.Set(cfg)
	if err := configManager.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	if err := configManager.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	ui.PrintSuccess("Configuration saved successfully!")

	ui.PrintSection("Next Steps")
	fmt.Println("  • Run 'dotfiles install' to install your configuration")
	fmt.Println("  • Run 'dotfiles config show' to review your settings")
	fmt.Println("  • Run 'dotfiles config edit' to modify settings")

	return nil
}

func loadPresetConfig(presetName string) (*config.Config, error) {
	// TODO: Implement preset loading from presets directory
	// For now, return default config
	manager := config.NewManager("")
	return manager.Get(), fmt.Errorf("preset loading not implemented yet")
}
