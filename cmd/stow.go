package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var stowCmd = &cobra.Command{
	Use:   "stow <packages>",
	Short: "Stow dotfile packages using GNU Stow",
	Long:  `Create symlinks for dotfile packages using GNU Stow`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Load existing config
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			fmt.Println("Run 'dotfiles init' to create a configuration first.")
			os.Exit(1)
		}

		// Check if stow is available
		if _, err := exec.LookPath("stow"); err != nil {
			fmt.Println("⚠️  GNU Stow not found. Install with: brew install stow")
			os.Exit(1)
		}

		// Get stow directory from flag or default to ~/.dotfiles
		stowDir, _ := cmd.Flags().GetString("dir")
		if stowDir == "" {
			stowDir = filepath.Join(home, ".dotfiles")
		}

		// Get target directory from flag or default to home
		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			target = home
		}

		var packages []string

		// Handle file input
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			filePackages, err := readPackagesFromFile(file)
			if err != nil {
				fmt.Printf("Error reading packages from file: %v\n", err)
				os.Exit(1)
			}
			packages = append(packages, filePackages...)
		}

		// Add command line arguments
		packages = append(packages, args...)

		if len(packages) == 0 {
			fmt.Println("No packages specified. Use command line arguments or --file flag.")
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")
		added := 0

		for _, pkg := range packages {
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				continue
			}

			// Check if package directory exists
			pkgPath := filepath.Join(stowDir, pkg)
			if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
				fmt.Printf("❌ Package directory not found: %s\n", pkgPath)
				continue
			}

			// Build stow command
			stowArgs := []string{
				"-d", stowDir,
				"-t", target,
			}

			if verbose {
				stowArgs = append(stowArgs, "-v")
			}

			if dryRun {
				stowArgs = append(stowArgs, "-n")
			}

			stowArgs = append(stowArgs, pkg)

			// Execute stow command
			stowCmd := exec.Command("stow", stowArgs...)
			if verbose || dryRun {
				fmt.Printf("Running: %s\n", strings.Join(stowCmd.Args, " "))
			}

			output, err := stowCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("❌ Error stowing %s: %v\n", pkg, err)
				if len(output) > 0 {
					fmt.Printf("   Output: %s\n", strings.TrimSpace(string(output)))
				}
				continue
			}

			if verbose && len(output) > 0 {
				fmt.Printf("   %s\n", strings.TrimSpace(string(output)))
			}

			if !dryRun {
				// Add to config if not already present
				if !contains(cfg.Stow, pkg) {
					cfg.Stow = append(cfg.Stow, pkg)
					added++
					fmt.Printf("✓ Stowed and added to config: %s\n", pkg)
				} else {
					fmt.Printf("✓ Stowed: %s (already in config)\n", pkg)
				}
			} else {
				fmt.Printf("✓ Would stow: %s\n", pkg)
			}
		}

		if !dryRun && added > 0 {
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\n📊 Added %d new stow packages to config\n", added)
		}
	},
}

var unstowCmd = &cobra.Command{
	Use:   "unstow <packages>",
	Short: "Unstow dotfile packages using GNU Stow",
	Long:  `Remove symlinks for dotfile packages using GNU Stow`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Load existing config
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Check if stow is available
		if _, err := exec.LookPath("stow"); err != nil {
			fmt.Println("⚠️  GNU Stow not found. Install with: brew install stow")
			os.Exit(1)
		}

		// Get stow directory from flag or default to ~/.dotfiles
		stowDir, _ := cmd.Flags().GetString("dir")
		if stowDir == "" {
			stowDir = filepath.Join(home, ".dotfiles")
		}

		// Get target directory from flag or default to home
		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			target = home
		}

		var packages []string

		// Handle file input
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			filePackages, err := readPackagesFromFile(file)
			if err != nil {
				fmt.Printf("Error reading packages from file: %v\n", err)
				os.Exit(1)
			}
			packages = append(packages, filePackages...)
		}

		// Add command line arguments
		packages = append(packages, args...)

		// Handle all flag
		if allStow, _ := cmd.Flags().GetBool("all"); allStow {
			packages = append(packages, cfg.Stow...)
		}

		if len(packages) == 0 {
			fmt.Println("No packages specified. Use command line arguments, --file flag, or --all flag.")
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")
		keepConfig, _ := cmd.Flags().GetBool("keep-config")
		removed := 0

		for _, pkg := range packages {
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				continue
			}

			// Build stow command
			stowArgs := []string{
				"-d", stowDir,
				"-t", target,
				"-D", // Delete (unstow)
			}

			if verbose {
				stowArgs = append(stowArgs, "-v")
			}

			if dryRun {
				stowArgs = append(stowArgs, "-n")
			}

			stowArgs = append(stowArgs, pkg)

			// Execute stow command
			stowCmd := exec.Command("stow", stowArgs...)
			if verbose || dryRun {
				fmt.Printf("Running: %s\n", strings.Join(stowCmd.Args, " "))
			}

			output, err := stowCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("❌ Error unstowing %s: %v\n", pkg, err)
				if len(output) > 0 {
					fmt.Printf("   Output: %s\n", strings.TrimSpace(string(output)))
				}
				continue
			}

			if verbose && len(output) > 0 {
				fmt.Printf("   %s\n", strings.TrimSpace(string(output)))
			}

			if !dryRun && !keepConfig {
				// Remove from config if present
				if contains(cfg.Stow, pkg) {
					cfg.Stow = removeFromSlice(cfg.Stow, pkg)
					removed++
					fmt.Printf("✓ Unstowed and removed from config: %s\n", pkg)
				} else {
					fmt.Printf("✓ Unstowed: %s (not in config)\n", pkg)
				}
			} else if !dryRun {
				fmt.Printf("✓ Unstowed: %s (kept in config)\n", pkg)
			} else {
				fmt.Printf("✓ Would unstow: %s\n", pkg)
			}
		}

		if !dryRun && removed > 0 {
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("Error saving configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\n📊 Removed %d stow packages from config\n", removed)
		}
	},
}

var restowCmd = &cobra.Command{
	Use:   "restow <packages>",
	Short: "Restow dotfile packages (unstow then stow)",
	Long:  `Remove and recreate symlinks for dotfile packages using GNU Stow`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")

		// Load existing config
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Check if stow is available
		if _, err := exec.LookPath("stow"); err != nil {
			fmt.Println("⚠️  GNU Stow not found. Install with: brew install stow")
			os.Exit(1)
		}

		// Get stow directory from flag or default to ~/.dotfiles
		stowDir, _ := cmd.Flags().GetString("dir")
		if stowDir == "" {
			stowDir = filepath.Join(home, ".dotfiles")
		}

		// Get target directory from flag or default to home
		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			target = home
		}

		var packages []string

		// Handle file input
		if file, _ := cmd.Flags().GetString("file"); file != "" {
			filePackages, err := readPackagesFromFile(file)
			if err != nil {
				fmt.Printf("Error reading packages from file: %v\n", err)
				os.Exit(1)
			}
			packages = append(packages, filePackages...)
		}

		// Add command line arguments
		packages = append(packages, args...)

		// Handle all flag
		if allStow, _ := cmd.Flags().GetBool("all"); allStow {
			packages = append(packages, cfg.Stow...)
		}

		if len(packages) == 0 {
			fmt.Println("No packages specified. Use command line arguments, --file flag, or --all flag.")
			return
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		for _, pkg := range packages {
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				continue
			}

			// Check if package directory exists
			pkgPath := filepath.Join(stowDir, pkg)
			if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
				fmt.Printf("❌ Package directory not found: %s\n", pkgPath)
				continue
			}

			// Build restow command
			stowArgs := []string{
				"-d", stowDir,
				"-t", target,
				"-R", // Restow
			}

			if verbose {
				stowArgs = append(stowArgs, "-v")
			}

			if dryRun {
				stowArgs = append(stowArgs, "-n")
			}

			stowArgs = append(stowArgs, pkg)

			// Execute stow command
			stowCmd := exec.Command("stow", stowArgs...)
			if verbose || dryRun {
				fmt.Printf("Running: %s\n", strings.Join(stowCmd.Args, " "))
			}

			output, err := stowCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("❌ Error restowing %s: %v\n", pkg, err)
				if len(output) > 0 {
					fmt.Printf("   Output: %s\n", strings.TrimSpace(string(output)))
				}
				continue
			}

			if verbose && len(output) > 0 {
				fmt.Printf("   %s\n", strings.TrimSpace(string(output)))
			}

			fmt.Printf("✓ Restowed: %s\n", pkg)
		}
	},
}

func init() {
	// Stow command flags
	stowCmd.Flags().StringP("dir", "d", "", "Stow directory (default: ~/.dotfiles)")
	stowCmd.Flags().StringP("target", "t", "", "Target directory (default: ~)")
	stowCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	stowCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without executing")
	stowCmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	// Unstow command flags
	unstowCmd.Flags().StringP("dir", "d", "", "Stow directory (default: ~/.dotfiles)")
	unstowCmd.Flags().StringP("target", "t", "", "Target directory (default: ~)")
	unstowCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	unstowCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without executing")
	unstowCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	unstowCmd.Flags().Bool("all", false, "Unstow all configured stow packages")
	unstowCmd.Flags().Bool("keep-config", false, "Don't remove packages from config")

	// Restow command flags
	restowCmd.Flags().StringP("dir", "d", "", "Stow directory (default: ~/.dotfiles)")
	restowCmd.Flags().StringP("target", "t", "", "Target directory (default: ~)")
	restowCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	restowCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without executing")
	restowCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	restowCmd.Flags().Bool("all", false, "Restow all configured stow packages")

	rootCmd.AddCommand(stowCmd)
	rootCmd.AddCommand(unstowCmd)
	rootCmd.AddCommand(restowCmd)
}