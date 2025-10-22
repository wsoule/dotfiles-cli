package cmd

import (
	"fmt"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var rolesCmd = &cobra.Command{
	Use:     "roles",
	GroupID: "advanced",
	Short:   "🎭 Manage package roles (web-dev, data-science, etc.)",
	Long: `🎭 Package Roles - Install Package Sets

Roles are predefined sets of packages for specific purposes:
• web-dev: Node.js, Docker, nginx, postgres
• data-science: Python, Jupyter, pandas, numpy
• minimal: git, vim, tmux (for servers)
• devops: kubectl, terraform, aws-cli, docker

Examples:
  dotfiles roles list                        # List all roles
  dotfiles roles create web-dev              # Create new role
  dotfiles roles install web-dev             # Install role packages
  dotfiles roles show web-dev                # Show role details`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles roles list' to see available roles")
	},
}

var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available roles",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if len(cfg.Roles) == 0 {
			fmt.Println("No roles defined. Create one with 'dotfiles roles create'")
			return
		}

		fmt.Println("📋 Available Roles:")
		fmt.Println()
		for name, role := range cfg.Roles {
			fmt.Printf("• %s - %s\n", name, role.Description)
			fmt.Printf("  Packages: %d brews, %d casks\n",
				len(role.Brews), len(role.Casks))
		}
	},
}

var rolesCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new role",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		roleName := args[0]
		description, _ := cmd.Flags().GetString("description")

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if cfg.Roles == nil {
			cfg.Roles = make(map[string]config.Role)
		}

		cfg.Roles[roleName] = config.Role{
			Name:        roleName,
			Description: description,
			Brews:       []string{},
			Casks:       []string{},
			Taps:        []string{},
			Stow:        []string{},
		}

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			return
		}

		fmt.Printf("✅ Created role '%s'\n", roleName)
		fmt.Println("Add packages with: dotfiles roles add-package <role> <package>")
	},
}

var rolesInstallCmd = &cobra.Command{
	Use:   "install <role>",
	Short: "Install all packages from a role",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		roleName := args[0]

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		role, exists := cfg.Roles[roleName]
		if !exists {
			fmt.Printf("❌ Role '%s' not found\n", roleName)
			return
		}

		fmt.Printf("📦 Installing role '%s'...\n", roleName)
		fmt.Printf("   %s\n", role.Description)
		fmt.Println()

		// Add role packages to main config
		cfg.Brews = append(cfg.Brews, role.Brews...)
		cfg.Casks = append(cfg.Casks, role.Casks...)
		cfg.Taps = append(cfg.Taps, role.Taps...)
		cfg.Stow = append(cfg.Stow, role.Stow...)

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			return
		}

		fmt.Println("✅ Role packages added to configuration")
		fmt.Println("💡 Run 'dotfiles install' to install the packages")
	},
}

func init() {
	rootCmd.AddCommand(rolesCmd)
	rolesCmd.AddCommand(rolesListCmd)
	rolesCmd.AddCommand(rolesCreateCmd)
	rolesCmd.AddCommand(rolesInstallCmd)

	rolesCreateCmd.Flags().String("description", "", "Role description")
}
