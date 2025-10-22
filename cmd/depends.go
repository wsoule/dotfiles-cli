package cmd

import (
	"fmt"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var dependsCmd = &cobra.Command{
	Use:     "depends",
	GroupID: "advanced",
	Short:   "🔗 Manage package dependencies",
	Long: `🔗 Package Dependency Management

Define dependencies between packages. When installing a package,
its dependencies will be installed first.

Examples:
  dotfiles depends add nvim vim-plug node
  dotfiles depends list nvim
  dotfiles depends check                    # Check all dependencies
  dotfiles depends resolve nvim             # Show install order`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles depends add' to define dependencies")
	},
}

var dependsAddCmd = &cobra.Command{
	Use:   "add <package> <dependencies...>",
	Short: "Add dependencies for a package",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		deps := args[1:]

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if cfg.PackageDependencies == nil {
			cfg.PackageDependencies = make(map[string][]string)
		}

		cfg.PackageDependencies[pkg] = deps

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			return
		}

		fmt.Printf("✅ Added dependencies for %s: %v\n", pkg, deps)
	},
}

var dependsListCmd = &cobra.Command{
	Use:   "list [package]",
	Short: "List dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if cfg.PackageDependencies == nil || len(cfg.PackageDependencies) == 0 {
			fmt.Println("No dependencies defined")
			return
		}

		if len(args) > 0 {
			// Show specific package
			pkg := args[0]
			if deps, ok := cfg.PackageDependencies[pkg]; ok {
				fmt.Printf("%s depends on:\n", pkg)
				for _, dep := range deps {
					fmt.Printf("  • %s\n", dep)
				}
			} else {
				fmt.Printf("No dependencies defined for %s\n", pkg)
			}
		} else {
			// Show all
			fmt.Println("📦 Package Dependencies:")
			for pkg, deps := range cfg.PackageDependencies {
				fmt.Printf("  %s -> %v\n", pkg, deps)
			}
		}
	},
}

var dependsResolveCmd = &cobra.Command{
	Use:   "resolve <package>",
	Short: "Show installation order for a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		order, err := resolveDependencies(cfg, pkg)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}

		fmt.Printf("📋 Installation order for %s:\n", pkg)
		for i, p := range order {
			fmt.Printf("  %d. %s\n", i+1, p)
		}
	},
}

var dependsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for circular dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if cfg.PackageDependencies == nil || len(cfg.PackageDependencies) == 0 {
			fmt.Println("No dependencies to check")
			return
		}

		fmt.Println("🔍 Checking for circular dependencies...")

		issues := 0
		for pkg := range cfg.PackageDependencies {
			if _, err := resolveDependencies(cfg, pkg); err != nil {
				fmt.Printf("❌ %s: %v\n", pkg, err)
				issues++
			}
		}

		if issues == 0 {
			fmt.Println("✅ No circular dependencies found")
		} else {
			fmt.Printf("⚠️  Found %d issue(s)\n", issues)
		}
	},
}

// resolveDependencies returns installation order using topological sort
func resolveDependencies(cfg *config.Config, pkg string) ([]string, error) {
	if cfg.PackageDependencies == nil {
		return []string{pkg}, nil
	}

	visited := make(map[string]bool)
	stack := make(map[string]bool)
	order := []string{}

	var visit func(string) error
	visit = func(p string) error {
		if stack[p] {
			return fmt.Errorf("circular dependency detected involving %s", p)
		}
		if visited[p] {
			return nil
		}

		stack[p] = true
		visited[p] = true

		if deps, ok := cfg.PackageDependencies[p]; ok {
			for _, dep := range deps {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}

		stack[p] = false
		order = append(order, p)
		return nil
	}

	if err := visit(pkg); err != nil {
		return nil, err
	}

	return order, nil
}

func init() {
	rootCmd.AddCommand(dependsCmd)
	dependsCmd.AddCommand(dependsAddCmd)
	dependsCmd.AddCommand(dependsListCmd)
	dependsCmd.AddCommand(dependsResolveCmd)
	dependsCmd.AddCommand(dependsCheckCmd)
}
