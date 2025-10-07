package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	"dotfiles/internal/pkgmanager"
	"dotfiles/internal/snapshot"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Generate package file and install packages",
	Long:  `Generates a package list file from your configuration and installs packages using the system package manager`,
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

		// Get the appropriate package manager for this OS
		pm, err := pkgmanager.GetPackageManager()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Check if package manager is available
		if !pm.IsAvailable() {
			fmt.Printf("⚠️  %s not found. Please install it first.\n", pm.GetName())
			if pm.GetName() == "homebrew" {
				fmt.Println("   /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
			}
			return
		}

		// Create auto-snapshot before installation (unless --no-snapshot flag)
		noSnapshot, _ := cmd.Flags().GetBool("no-snapshot")
		if !noSnapshot {
			fmt.Println("📸 Creating snapshot before installation...")
			timestamp, err := snapshot.CreateAutoSnapshot("Before package installation")
			if err != nil {
				fmt.Printf("⚠️  Warning: Could not create snapshot: %v\n", err)
			} else {
				fmt.Printf("   ✅ Snapshot created: %s\n", timestamp)
			}
		}

		// Run pre-install hooks
		if cfg.Hooks != nil && len(cfg.Hooks.PreInstall) > 0 {
			if err := RunHooks(cfg.Hooks.PreInstall, "pre-install"); err != nil {
				fmt.Printf("⚠️  Pre-install hook failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Generate package list file
		fileContent, err := pm.GenerateInstallFile(cfg.Brews, cfg.Casks, cfg.Taps)
		if err != nil {
			fmt.Printf("Error generating package file: %v\n", err)
			os.Exit(1)
		}

		if fileContent == "" {
			fmt.Println("No packages configured. Run 'dotfiles add <package>' first.")
			return
		}

		// Write package file
		fileName := "Brewfile"
		if pm.GetName() == "pacman" {
			fileName = "packages.txt"
		} else if pm.GetName() == "apt" {
			fileName = "packages.txt"
		} else if pm.GetName() == "yum" {
			fileName = "packages.txt"
		}

		filePath := "./" + fileName
		if output, _ := cmd.Flags().GetString("output"); output != "" {
			filePath = output
		}

		if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
			fmt.Printf("Error writing package file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Generated package list at: %s\n", filePath)

		// Run install unless --dry-run is specified
		if dryRun, _ := cmd.Flags().GetBool("dry-run"); dryRun {
			fmt.Println("🔍 Dry run - would install packages using " + pm.GetName())
			return
		}

		fmt.Printf("📦 Installing packages with %s...\n", pm.GetName())

		// Install packages by type
		if len(cfg.Taps) > 0 {
			fmt.Println("Installing taps...")
			if err := pm.Install(cfg.Taps, "tap"); err != nil {
				fmt.Printf("⚠️  Error installing taps: %v\n", err)
			}
		}

		if len(cfg.Brews) > 0 {
			fmt.Println("Installing packages...")
			if err := pm.Install(cfg.Brews, "brew"); err != nil {
				fmt.Printf("Error installing packages: %v\n", err)
				os.Exit(1)
			}
		}

		if len(cfg.Casks) > 0 {
			fmt.Println("Installing casks/applications...")
			if err := pm.Install(cfg.Casks, "cask"); err != nil {
				fmt.Printf("⚠️  Error installing casks: %v\n", err)
			}
		}

		fmt.Println("✅ Installation complete!")

		// Run post-install hooks
		if cfg.Hooks != nil && len(cfg.Hooks.PostInstall) > 0 {
			if err := RunHooks(cfg.Hooks.PostInstall, "post-install"); err != nil {
				fmt.Printf("⚠️  Post-install hook failed: %v\n", err)
			}
		}

		// Run package-specific hooks (pre-install and post-install)
		if cfg.PackageConfigs != nil {
			allPackages := append(append([]string{}, cfg.Brews...), cfg.Casks...)
			for _, pkg := range allPackages {
				if pkgConfig, exists := cfg.PackageConfigs[pkg]; exists {
					// Run pre-install hooks for this package
					if len(pkgConfig.PreInstall) > 0 {
						fmt.Printf("🔧 Running pre-install hooks for package: %s\n", pkg)
						if err := RunHooks(pkgConfig.PreInstall, fmt.Sprintf("%s pre-install", pkg)); err != nil {
							fmt.Printf("⚠️  Package pre-install hook failed for %s: %v\n", pkg, err)
						}
					}

					// Run post-install hooks for this package
					if len(pkgConfig.PostInstall) > 0 {
						fmt.Printf("🔧 Running post-install hooks for package: %s\n", pkg)
						if err := RunHooks(pkgConfig.PostInstall, fmt.Sprintf("%s post-install", pkg)); err != nil {
							fmt.Printf("⚠️  Package post-install hook failed for %s: %v\n", pkg, err)
						}
					}
				}
			}
		}
	},
}

func init() {
	installCmd.Flags().StringP("output", "o", "", "Output path for the package file (default: ./Brewfile or ./packages.txt)")
	installCmd.Flags().Bool("no-snapshot", false, "Skip creating automatic snapshot before installation")
	installCmd.Flags().Bool("dry-run", false, "Show what would be installed without executing")
	rootCmd.AddCommand(installCmd)
}