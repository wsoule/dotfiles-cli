package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Generate Brewfile and install packages",
	Long:  `Generates a Brewfile from your configuration and runs 'brew bundle' to install packages`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Run pre-install hooks
		if cfg.Hooks != nil && len(cfg.Hooks.PreInstall) > 0 {
			if err := RunHooks(cfg.Hooks.PreInstall, "pre-install"); err != nil {
				fmt.Printf("⚠️  Pre-install hook failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Generate Brewfile content
		brewfileContent := cfg.GenerateBrewfile()
		if brewfileContent == "" {
			fmt.Println("No packages configured. Run 'dotfiles add <package>' first.")
			return
		}

		// Write Brewfile
		brewfilePath := "./Brewfile"
		if output, _ := cmd.Flags().GetString("output"); output != "" {
			brewfilePath = output
		}

		if err := os.WriteFile(brewfilePath, []byte(brewfileContent), 0644); err != nil {
			fmt.Printf("Error writing Brewfile: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Generated Brewfile at: %s\n", brewfilePath)

		// Check if brew is available
		if _, err := exec.LookPath("brew"); err != nil {
			fmt.Println("⚠️  Homebrew not found. Install Homebrew first:")
			fmt.Println("   /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
			return
		}

		// Run brew bundle unless --dry-run is specified
		if dryRun, _ := cmd.Flags().GetBool("dry-run"); dryRun {
			fmt.Println("🔍 Dry run - would execute:")
			fmt.Printf("   brew bundle --file=%s\n", brewfilePath)
			return
		}

		fmt.Println("🍺 Installing packages with Homebrew...")
		brewCmd := exec.Command("brew", "bundle", "--file="+brewfilePath)
		brewCmd.Stdout = os.Stdout
		brewCmd.Stderr = os.Stderr

		if err := brewCmd.Run(); err != nil {
			fmt.Printf("Error running brew bundle: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Installation complete!")

		// Run post-install hooks
		if cfg.Hooks != nil && len(cfg.Hooks.PostInstall) > 0 {
			if err := RunHooks(cfg.Hooks.PostInstall, "post-install"); err != nil {
				fmt.Printf("⚠️  Post-install hook failed: %v\n", err)
			}
		}

		// Run package-specific post-install hooks
		if cfg.PackageConfigs != nil {
			allPackages := append(append([]string{}, cfg.Brews...), cfg.Casks...)
			for _, pkg := range allPackages {
				if pkgConfig, exists := cfg.PackageConfigs[pkg]; exists && len(pkgConfig.PostInstall) > 0 {
					fmt.Printf("🔧 Running post-install hooks for package: %s\n", pkg)
					if err := RunHooks(pkgConfig.PostInstall, fmt.Sprintf("%s post-install", pkg)); err != nil {
						fmt.Printf("⚠️  Package hook failed for %s: %v\n", pkg, err)
					}
				}
			}
		}
	},
}

func init() {
	installCmd.Flags().StringP("output", "o", "./Brewfile", "Output path for the Brewfile")
	installCmd.Flags().Bool("dry-run", false, "Show what would be installed without executing")
	rootCmd.AddCommand(installCmd)
}