package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "🏷️  Manage package groups and tags",
	Long: `🏷️  Package Groups & Tags - Organize Your Packages

Group related packages together for easier management.
Groups allow you to install, remove, or manage sets of packages at once.

Examples:
  dotfiles groups list                          # List all groups
  dotfiles groups create dev git,neovim,tmux    # Create group with packages
  dotfiles groups add dev docker                # Add package to group
  dotfiles groups remove dev docker             # Remove package from group
  dotfiles groups install dev                   # Install all packages in group
  dotfiles groups show dev                      # Show packages in group`,
}

var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all package groups",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Groups == nil || len(cfg.Groups) == 0 {
			fmt.Println("📋 No groups defined")
			fmt.Println()
			fmt.Println("💡 Create a group:")
			fmt.Println("   dotfiles groups create <name> <package1>,<package2>,...")
			return
		}

		// Sort groups by name
		var groupNames []string
		for name := range cfg.Groups {
			groupNames = append(groupNames, name)
		}
		sort.Strings(groupNames)

		fmt.Printf("🏷️  Found %d group(s):\n", len(cfg.Groups))
		fmt.Println("=" + strings.Repeat("=", 30))
		fmt.Println()

		for _, name := range groupNames {
			packages := cfg.Groups[name]
			fmt.Printf("📦 %s (%d packages)\n", name, len(packages))
			if len(packages) > 0 {
				fmt.Printf("   %s\n", strings.Join(packages, ", "))
			}
			fmt.Println()
		}

		fmt.Println("💡 Usage:")
		fmt.Println("   dotfiles groups show <name>       # View group details")
		fmt.Println("   dotfiles groups install <name>    # Install group packages")
	},
}

var groupsCreateCmd = &cobra.Command{
	Use:   "create <name> <packages>",
	Short: "Create a new package group",
	Long: `Create a new package group with specified packages.
Packages should be comma-separated.

Example:
  dotfiles groups create dev git,neovim,tmux,docker`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		packagesStr := args[1]

		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Groups == nil {
			cfg.Groups = make(map[string][]string)
		}

		// Parse packages
		packages := strings.Split(packagesStr, ",")
		for i, pkg := range packages {
			packages[i] = strings.TrimSpace(pkg)
		}

		// Check if group exists
		if _, exists := cfg.Groups[groupName]; exists {
			fmt.Printf("⚠️  Group '%s' already exists\n", groupName)
			fmt.Println()
			fmt.Println("💡 To add packages to existing group:")
			fmt.Printf("   dotfiles groups add %s <package>\n", groupName)
			os.Exit(1)
		}

		cfg.Groups[groupName] = packages

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Created group '%s' with %d package(s)\n", groupName, len(packages))
		fmt.Printf("   %s\n", strings.Join(packages, ", "))
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Printf("   dotfiles groups install %s    # Install all packages\n", groupName)
	},
}

var groupsAddCmd = &cobra.Command{
	Use:   "add <group> <package>",
	Short: "Add a package to a group",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		packageName := args[1]

		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Groups == nil {
			cfg.Groups = make(map[string][]string)
		}

		// Check if group exists
		if _, exists := cfg.Groups[groupName]; !exists {
			fmt.Printf("❌ Group '%s' not found\n", groupName)
			os.Exit(1)
		}

		// Check if package already in group
		for _, pkg := range cfg.Groups[groupName] {
			if pkg == packageName {
				fmt.Printf("⚠️  Package '%s' already in group '%s'\n", packageName, groupName)
				return
			}
		}

		cfg.Groups[groupName] = append(cfg.Groups[groupName], packageName)

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Added '%s' to group '%s'\n", packageName, groupName)
	},
}

var groupsRemoveCmd = &cobra.Command{
	Use:   "remove <group> <package>",
	Short: "Remove a package from a group",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		packageName := args[1]

		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Groups == nil || cfg.Groups[groupName] == nil {
			fmt.Printf("❌ Group '%s' not found\n", groupName)
			os.Exit(1)
		}

		// Remove package
		newPackages := []string{}
		found := false
		for _, pkg := range cfg.Groups[groupName] {
			if pkg != packageName {
				newPackages = append(newPackages, pkg)
			} else {
				found = true
			}
		}

		if !found {
			fmt.Printf("⚠️  Package '%s' not in group '%s'\n", packageName, groupName)
			return
		}

		cfg.Groups[groupName] = newPackages

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Removed '%s' from group '%s'\n", packageName, groupName)
	},
}

var groupsShowCmd = &cobra.Command{
	Use:   "show <group>",
	Short: "Show packages in a group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]

		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Groups == nil || cfg.Groups[groupName] == nil {
			fmt.Printf("❌ Group '%s' not found\n", groupName)
			os.Exit(1)
		}

		packages := cfg.Groups[groupName]
		fmt.Printf("📦 Group: %s\n", groupName)
		fmt.Println("=" + strings.Repeat("=", 30))
		fmt.Println()
		fmt.Printf("Packages (%d):\n", len(packages))
		for i, pkg := range packages {
			fmt.Printf("  %d. %s\n", i+1, pkg)
		}
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Printf("   dotfiles groups install %s    # Install all packages\n", groupName)
		fmt.Printf("   dotfiles groups add %s <pkg>  # Add more packages\n", groupName)
	},
}

var groupsInstallCmd = &cobra.Command{
	Use:   "install <group>",
	Short: "Install all packages in a group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]

		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Groups == nil || cfg.Groups[groupName] == nil {
			fmt.Printf("❌ Group '%s' not found\n", groupName)
			os.Exit(1)
		}

		packages := cfg.Groups[groupName]
		fmt.Printf("📦 Installing group '%s' (%d packages)...\n", groupName, len(packages))
		fmt.Println()

		// Add packages to config if not already present
		added := 0
		for _, pkg := range packages {
			if !contains(cfg.Brews, pkg) && !contains(cfg.Casks, pkg) {
				// Assume brew by default, user can modify if needed
				cfg.Brews = append(cfg.Brews, pkg)
				added++
			}
		}

		if added > 0 {
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("❌ Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Added %d package(s) to configuration\n", added)
			fmt.Println()
		}

		fmt.Println("💡 Run installation:")
		fmt.Println("   dotfiles install")
	},
}

func init() {
	groupsCmd.AddCommand(groupsListCmd)
	groupsCmd.AddCommand(groupsCreateCmd)
	groupsCmd.AddCommand(groupsAddCmd)
	groupsCmd.AddCommand(groupsRemoveCmd)
	groupsCmd.AddCommand(groupsShowCmd)
	groupsCmd.AddCommand(groupsInstallCmd)

	rootCmd.AddCommand(groupsCmd)
}
