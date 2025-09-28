package cmd

import (
	"fmt"

	"dotfiles/internal/config"
	"dotfiles/internal/installer"
	"dotfiles/internal/ui"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dotfiles and configure system",
	Long: `Install dotfiles and configure your system based on the current configuration.

This command will:
- Install Homebrew (if enabled and not present)
- Install packages from Brewfile and configuration
- Apply dotfiles using GNU Stow
- Configure macOS system defaults
- Set up development environment
- Install global npm packages

Run 'dotfiles setup' first to create your configuration.`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().Bool("dry-run", false, "Show what would be installed without making changes")
	installCmd.Flags().Bool("skip-homebrew", false, "Skip Homebrew installation")
	installCmd.Flags().Bool("skip-macos", false, "Skip macOS system configuration")
	installCmd.Flags().Bool("skip-stow", false, "Skip dotfiles installation with Stow")
}

func runInstall(cmd *cobra.Command, args []string) error {
	ui.Banner()
	ui.PrintSection("Installation")

	// Load configuration
	configManager := config.NewManager(cfgFile)
	if err := configManager.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w\nRun 'dotfiles setup' to create a configuration", err)
	}

	cfg := configManager.Get()
	if cfg == nil {
		return fmt.Errorf("no configuration found\nRun 'dotfiles setup' to create a configuration")
	}

	// Get flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	skipHomebrew, _ := cmd.Flags().GetBool("skip-homebrew")
	skipMacOS, _ := cmd.Flags().GetBool("skip-macos")
	skipStow, _ := cmd.Flags().GetBool("skip-stow")

	// Create installer
	installerManager := installer.NewManager(cfg, dryRun)

	// Run installation steps
	if cfg.Installation.Homebrew && !skipHomebrew {
		if err := installerManager.InstallHomebrew(); err != nil {
			return fmt.Errorf("homebrew installation failed: %w", err)
		}
	}

	if cfg.Installation.Brewfile {
		if err := installerManager.InstallBrewPackages(); err != nil {
			return fmt.Errorf("brew packages installation failed: %w", err)
		}
	}

	if cfg.Installation.Dotfiles && !skipStow {
		if err := installerManager.InstallDotfiles(); err != nil {
			return fmt.Errorf("dotfiles installation failed: %w", err)
		}
	}

	if cfg.Installation.MacOSDefaults && !skipMacOS {
		if err := installerManager.ApplyMacOSDefaults(); err != nil {
			return fmt.Errorf("macOS configuration failed: %w", err)
		}
	}

	if cfg.Installation.NPMPackages {
		if err := installerManager.InstallNPMPackages(); err != nil {
			return fmt.Errorf("npm packages installation failed: %w", err)
		}
	}

	// Post-installation setup
	if err := installerManager.PostInstallSetup(); err != nil {
		return fmt.Errorf("post-installation setup failed: %w", err)
	}

	fmt.Println()
	ui.PrintSuccess("Dotfiles installation completed successfully!")

	ui.PrintSection("Next Steps")
	fmt.Println("  1. Restart your terminal or run: source ~/.zshrc")
	fmt.Println("  2. Configure p10k if needed: p10k configure")
	fmt.Println("  3. Add any private configurations to the private/ directory")
	fmt.Println("  4. Set up SSH keys if needed")

	return nil
}