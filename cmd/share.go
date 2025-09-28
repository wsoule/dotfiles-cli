package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
	"dotfiles/internal/ui"
	"github.com/spf13/cobra"
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share and import dotfiles configurations",
	Long: `Share your configuration with others or import shared configurations.

This command helps you:
- Export your configuration for sharing (with personal info removed)
- Import configurations from others
- Create shareable preset files
- Validate shared configurations`,
}

// shareExportCmd exports configuration for sharing
var shareExportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export your configuration for sharing",
	Long: `Export your current configuration with personal information removed.

The exported file can be shared with others or used as a preset.
Personal information (name, email) is removed for privacy.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runShareExport,
}

// shareImportCmd imports a shared configuration
var shareImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import a shared configuration",
	Long: `Import a configuration file shared by another user.

This will load the configuration and prompt you to:
- Review the settings
- Add your personal information
- Choose what to enable/disable
- Save as your configuration`,
	Args: cobra.ExactArgs(1),
	RunE: runShareImport,
}

// shareValidateCmd validates a configuration file
var shareValidateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a configuration file",
	Long: `Validate that a configuration file is properly formatted and contains valid settings.

This checks:
- JSON syntax and structure
- Required fields are present
- Values are within valid ranges
- No conflicting settings`,
	Args: cobra.ExactArgs(1),
	RunE: runShareValidate,
}

func init() {
	rootCmd.AddCommand(shareCmd)

	// Add subcommands
	shareCmd.AddCommand(shareExportCmd)
	shareCmd.AddCommand(shareImportCmd)
	shareCmd.AddCommand(shareValidateCmd)

	// Flags for export
	shareExportCmd.Flags().Bool("include-personal", false, "Include personal information in export")
	shareExportCmd.Flags().Bool("preset", false, "Export as a preset file to presets/ directory")
	shareExportCmd.Flags().String("description", "", "Description for the exported configuration")

	// Flags for import
	shareImportCmd.Flags().Bool("force", false, "Overwrite existing configuration without prompting")
	shareImportCmd.Flags().Bool("dry-run", false, "Show what would be imported without saving")
}

func runShareExport(cmd *cobra.Command, args []string) error {
	ui.PrintSection("Export Configuration")

	// Load current configuration
	configManager := config.NewManager(cfgFile)
	if err := configManager.Load(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg := configManager.Get()
	if cfg == nil {
		return fmt.Errorf("no configuration found\nRun 'dotfiles setup' to create a configuration")
	}

	// Get flags
	includePersonal, _ := cmd.Flags().GetBool("include-personal")
	isPreset, _ := cmd.Flags().GetBool("preset")
	description, _ := cmd.Flags().GetString("description")

	// Create export copy
	exportCfg := *cfg

	// Remove personal information unless explicitly requested
	if !includePersonal {
		exportCfg.Personal.Name = ""
		exportCfg.Personal.Email = ""
		ui.PrintInfo("Personal information removed for privacy")
	}

	// Update metadata
	exportCfg.Metadata.CreatedBy = "exported"
	if description != "" {
		exportCfg.Metadata.Description = description
	}

	// Determine filename
	var filename string
	if len(args) > 0 {
		filename = args[0]
	} else if isPreset {
		filename = ui.Input("Enter preset name", "my-config")
		if !strings.HasSuffix(filename, ".json") {
			filename += ".json"
		}
		filename = filepath.Join("presets", filename)
	} else {
		filename = "shared-config.json"
	}

	// Ensure directory exists for presets
	if isPreset {
		if err := os.MkdirAll("presets", 0755); err != nil {
			return fmt.Errorf("failed to create presets directory: %w", err)
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(&exportCfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Configuration exported to: %s", filename))

	if isPreset {
		fmt.Println()
		ui.PrintInfo("Preset saved! Others can use it with:")
		fmt.Printf("  dotfiles setup --preset %s\n", strings.TrimSuffix(filepath.Base(filename), ".json"))
	} else {
		fmt.Println()
		ui.PrintInfo("Share this file with others! They can import it with:")
		fmt.Printf("  dotfiles share import %s\n", filename)
	}

	return nil
}

func runShareImport(cmd *cobra.Command, args []string) error {
	filename := args[0]

	ui.PrintSection("Import Configuration")
	ui.PrintInfo(fmt.Sprintf("Importing from: %s", filename))

	// Read and validate file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var importCfg config.Config
	if err := json.Unmarshal(data, &importCfg); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Get flags
	force, _ := cmd.Flags().GetBool("force")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// Show preview
	fmt.Println()
	ui.PrintInfo("Configuration Preview:")
	showImportPreview(&importCfg)

	if dryRun {
		ui.PrintInfo("Dry run - no changes made")
		return nil
	}

	// Check if config exists
	configManager := config.NewManager(cfgFile)
	if !force {
		if err := configManager.Load(); err == nil {
			ui.PrintWarning("Existing configuration found!")
			if !ui.Confirm("This will overwrite your current configuration. Continue?") {
				ui.PrintInfo("Import cancelled")
				return nil
			}
		}
	}

	// Prompt for personal information if missing
	if importCfg.Personal.Name == "" {
		fmt.Println()
		ui.PrintSection("Personal Information")
		importCfg.Personal.Name = ui.Input("Your name", "")
		importCfg.Personal.Email = ui.Input("Your email", "")
	}

	// Update metadata
	importCfg.Metadata.CreatedBy = "imported"

	// Save configuration
	configManager.Set(&importCfg)
	if err := configManager.Validate(); err != nil {
		return fmt.Errorf("imported configuration is invalid: %w", err)
	}

	if err := configManager.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	ui.PrintSuccess("Configuration imported successfully!")
	fmt.Println()
	ui.PrintInfo("Next steps:")
	fmt.Println("  • Run 'dotfiles install' to apply the configuration")
	fmt.Println("  • Run 'dotfiles config show' to review settings")

	return nil
}

func runShareValidate(cmd *cobra.Command, args []string) error {
	filename := args[0]

	ui.PrintSection("Validate Configuration")
	ui.PrintInfo(fmt.Sprintf("Validating: %s", filename))

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to read file: %v", err))
		return err
	}

	// Parse JSON
	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		ui.PrintError(fmt.Sprintf("Invalid JSON: %v", err))
		return err
	}

	ui.PrintSuccess("JSON syntax is valid")

	// Validate structure
	configManager := config.NewManager("")
	configManager.Set(&cfg)

	if err := configManager.Validate(); err != nil {
		ui.PrintError(fmt.Sprintf("Configuration validation failed: %v", err))
		return err
	}

	ui.PrintSuccess("Configuration structure is valid")

	// Show summary
	fmt.Println()
	ui.PrintInfo("Configuration Summary:")
	showImportPreview(&cfg)

	return nil
}

func showImportPreview(cfg *config.Config) {
	if cfg.Metadata.Description != "" {
		fmt.Printf("  Description: %s\n", cfg.Metadata.Description)
	}
	if cfg.Metadata.CreatedBy != "" {
		fmt.Printf("  Created by: %s\n", cfg.Metadata.CreatedBy)
	}

	// Count enabled features
	langCount := 0
	for _, enabled := range cfg.Development.Languages {
		if enabled {
			langCount++
		}
	}

	fmt.Printf("  Languages: %d enabled\n", langCount)
	fmt.Printf("  Extra packages: %d brews, %d casks\n",
		len(cfg.Packages.ExtraBrews), len(cfg.Packages.ExtraCasks))
	fmt.Printf("  System settings: dock=%s, dark_mode=%t\n",
		cfg.System.Dock.Position, cfg.System.Appearance.DarkMode)
}