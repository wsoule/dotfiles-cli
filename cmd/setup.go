package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"dotfiles/internal/config"
	"dotfiles/internal/server"
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
	setupCmd.Flags().Bool("cli", false, "Use CLI wizard instead of web interface")
	setupCmd.Flags().Int("port", 0, "Port for web interface (random if not specified)")
}

func runSetup(cmd *cobra.Command, args []string) error {
	ui.Banner()
	ui.PrintSection("Setup Wizard")

	// Get flags
	force, _ := cmd.Flags().GetBool("force")
	preset, _ := cmd.Flags().GetString("preset")
	quick, _ := cmd.Flags().GetBool("quick")
	useCLI, _ := cmd.Flags().GetBool("cli")
	port, _ := cmd.Flags().GetInt("port")

	// Check if config exists and force flag
	configManager := config.NewManager(cfgFile)
	if !force {
		if err := configManager.Load(); err == nil {
			ui.PrintWarning("Configuration already exists!")
			fmt.Println("    Use --force to overwrite or run 'dotfiles config' to modify existing settings.")
			return nil
		}
	}

	// Handle preset loading
	if preset != "" {
		cfg, err := loadPresetConfig(preset)
		if err != nil {
			return fmt.Errorf("failed to load preset '%s': %w", preset, err)
		}
		ui.PrintSuccess(fmt.Sprintf("Loaded preset configuration: %s", preset))

		configManager.Set(cfg)
		if err := configManager.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		ui.PrintSuccess("Configuration saved successfully!")
		return nil
	}

	// Choose interface type
	if !useCLI {
		return runWebSetup(port)
	} else {
		return runCLISetup(quick)
	}
}

func runWebSetup(port int) error {
	ui.PrintInfo("Starting modern web-based setup wizard...")

	// Generate random port if not specified
	if port == 0 {
		rand.Seed(time.Now().UnixNano())
		port = rand.Intn(10000) + 50000 // Random port between 50000-60000
	}

	// Create and start server
	srv := server.NewServer(port)
	if err := srv.Start(); err != nil {
		ui.PrintError("Failed to start web server")
		ui.PrintInfo("Falling back to CLI wizard...")
		return runCLISetup(false)
	}

	ui.PrintInfo("Waiting for setup to complete...")
	ui.PrintInfo("(Close the browser tab when you're done)")

	// Wait for completion
	srv.WaitForCompletion()

	// Stop server
	srv.Stop()

	ui.PrintSuccess("Configuration saved successfully!")

	ui.PrintSection("Next Steps")
	fmt.Println("  • Run 'dotfiles install' to install your configuration")
	fmt.Println("  • Run 'dotfiles config show' to review your settings")

	return nil
}

func runCLISetup(quick bool) error {
	ui.PrintInfo("Using CLI wizard...")
	ui.PrintWarning("For a better experience, try the web wizard: dotfiles setup")

	// TODO: Implement fallback CLI wizard
	ui.PrintError("CLI wizard not yet implemented")
	ui.PrintInfo("Please use the web wizard or import a preset configuration")

	return fmt.Errorf("CLI wizard not available - use web interface")
}

func loadPresetConfig(presetName string) (*config.Config, error) {
	// TODO: Implement preset loading from presets directory
	// For now, return default config
	manager := config.NewManager("")
	return manager.Get(), fmt.Errorf("preset loading not implemented yet")
}
