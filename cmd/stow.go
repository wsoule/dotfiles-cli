package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var stowCmd = &cobra.Command{
	Use:   "stow <packages>",
	Short: "🔗 Create dotfile symlinks using GNU Stow",
	Long: `🔗 Stow Dotfiles - Manage Configuration Symlinks

Create symlinks for dotfile packages using GNU Stow. This allows you to
keep your dotfiles organized in ~/.dotfiles/stow/ while having them
appear in their expected locations in your home directory.

How it works:
1. Files in ~/.dotfiles/stow/vim/.vimrc become ~/.vimrc
2. Directories in ~/.dotfiles/stow/zsh/.zsh become ~/.zsh
3. Config files in ~/.dotfiles/stow/config/.config/app become ~/.config/app

Examples:
  dotfiles stow vim zsh tmux                 # Stow multiple packages
  dotfiles stow --dir=/custom/path vim       # Use custom stow directory
  dotfiles stow --target=/tmp vim            # Stow to custom target
  dotfiles stow --dry-run --verbose vim      # Preview what would happen
  dotfiles stow config                       # Stow .config applications

Common packages to stow:
• vim, zsh, tmux - Core development tools
• git - Git configuration and aliases
• config - Applications that use ~/.config/
• shell - Shell environment and aliases`,
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

		// Get stow directory from flag or default to ~/.dotfiles/stow
		stowDir, _ := cmd.Flags().GetString("dir")
		if stowDir == "" {
			stowDir = filepath.Join(home, ".dotfiles", "stow")
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
		backup, _ := cmd.Flags().GetBool("backup")
		autoResolve, _ := cmd.Flags().GetBool("auto-resolve")
		added := 0

		for _, pkg := range packages {
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				continue
			}

			// Check if package directory exists, if not try to import it
			pkgPath := filepath.Join(stowDir, pkg)
			if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
				// Try to auto-import from home directory
				homeDirPath := filepath.Join(target, "."+pkg)
				if _, err := os.Stat(homeDirPath); err == nil {
					fmt.Printf("📥 Found ~/.%s directory, importing...\n", pkg)
					if err := importDotfileDirectory(pkg, homeDirPath, pkgPath); err != nil {
						fmt.Printf("❌ Failed to import ~/.%s: %v\n", pkg, err)
						continue
					}
					fmt.Printf("✅ Successfully imported ~/.%s to stow package\n", pkg)
				} else {
					fmt.Printf("❌ Package directory not found: %s\n", pkgPath)
					fmt.Printf("   💡 Create it manually or place files in ~/.%s to auto-import\n", pkg)
					continue
				}
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

			// Check for conflicts before stowing
			if !dryRun {
				conflicts := findStowConflicts(pkgPath, target)
				if len(conflicts) > 0 {
					fmt.Printf("⚠️  Found %d conflicts for package '%s':\n", len(conflicts), pkg)
					for _, conflict := range conflicts {
						fmt.Printf("   %s\n", conflict)
					}

					if autoResolve {
						fmt.Printf("🔧 Auto-resolving conflicts...\n")
						if err := resolveStowConflicts(conflicts, backup, verbose); err != nil {
							fmt.Printf("❌ Failed to resolve conflicts: %v\n", err)
							continue
						}
					} else {
						fmt.Printf("💡 Use --auto-resolve to automatically handle conflicts\n")
						fmt.Printf("   Or use --backup to backup existing files\n")
						continue
					}
				}
			}

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

		// Get stow directory from flag or default to ~/.dotfiles/stow
		stowDir, _ := cmd.Flags().GetString("dir")
		if stowDir == "" {
			stowDir = filepath.Join(home, ".dotfiles", "stow")
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

		// Get stow directory from flag or default to ~/.dotfiles/stow
		stowDir, _ := cmd.Flags().GetString("dir")
		if stowDir == "" {
			stowDir = filepath.Join(home, ".dotfiles", "stow")
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

var privateCmd = &cobra.Command{
	Use:   "private <package> <file>",
	Short: "Add a private file to a stow package",
	Long:  `Create a symlink in a stow package that points to a file in the private directory`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		pkg := args[0]
		filename := args[1]

		stowDir := filepath.Join(home, ".dotfiles", "stow")
		privateDir := filepath.Join(home, ".dotfiles", "private")

		// Ensure the private directory exists
		if err := os.MkdirAll(privateDir, 0755); err != nil {
			fmt.Printf("Error creating private directory: %v\n", err)
			os.Exit(1)
		}

		// Ensure the stow package directory exists
		pkgPath := filepath.Join(stowDir, pkg)
		if err := os.MkdirAll(pkgPath, 0755); err != nil {
			fmt.Printf("Error creating stow package directory: %v\n", err)
			os.Exit(1)
		}

		privatePath := filepath.Join(privateDir, filename)
		stowLinkPath := filepath.Join(pkgPath, filename)

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if dryRun {
			fmt.Printf("Would create symlink: %s -> ../../private/%s\n", stowLinkPath, filename)
			return
		}

		// Check if the private file exists
		if _, err := os.Stat(privatePath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Private file doesn't exist: %s\n", privatePath)
			fmt.Printf("   Create the file first, then run this command again.\n")
			return
		}

		// Create relative symlink from stow package to private directory
		relativePrivatePath := filepath.Join("..", "..", "private", filename)

		// Remove existing symlink if it exists
		if _, err := os.Lstat(stowLinkPath); err == nil {
			if err := os.Remove(stowLinkPath); err != nil {
				fmt.Printf("Error removing existing symlink: %v\n", err)
				os.Exit(1)
			}
		}

		// Create the symlink
		if err := os.Symlink(relativePrivatePath, stowLinkPath); err != nil {
			fmt.Printf("Error creating symlink: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("Created symlink: %s -> %s\n", stowLinkPath, relativePrivatePath)
		}
		fmt.Printf("✅ Private file linked: %s/%s -> private/%s\n", pkg, filename, filename)
		fmt.Printf("   💡 Now run: dotfiles stow %s\n", pkg)
	},
}

func init() {
	// Stow command flags
	stowCmd.Flags().StringP("dir", "d", "", "Stow directory (default: ~/.dotfiles/stow)")
	stowCmd.Flags().StringP("target", "t", "", "Target directory (default: ~)")
	stowCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	stowCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without executing")
	stowCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	stowCmd.Flags().Bool("backup", false, "Backup existing files before stowing")
	stowCmd.Flags().Bool("auto-resolve", false, "Automatically resolve conflicts")

	// Unstow command flags
	unstowCmd.Flags().StringP("dir", "d", "", "Stow directory (default: ~/.dotfiles/stow)")
	unstowCmd.Flags().StringP("target", "t", "", "Target directory (default: ~)")
	unstowCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	unstowCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without executing")
	unstowCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	unstowCmd.Flags().Bool("all", false, "Unstow all configured stow packages")
	unstowCmd.Flags().Bool("keep-config", false, "Don't remove packages from config")

	// Restow command flags
	restowCmd.Flags().StringP("dir", "d", "", "Stow directory (default: ~/.dotfiles/stow)")
	restowCmd.Flags().StringP("target", "t", "", "Target directory (default: ~)")
	restowCmd.Flags().StringP("file", "f", "", "Read packages from file (one per line)")
	restowCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without executing")
	restowCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	restowCmd.Flags().Bool("all", false, "Restow all configured stow packages")

	// Private command flags
	privateCmd.Flags().BoolP("dry-run", "n", false, "Show what would be done without executing")
	privateCmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(stowCmd)
	rootCmd.AddCommand(unstowCmd)
	rootCmd.AddCommand(restowCmd)
	rootCmd.AddCommand(privateCmd)
}

// importDotfileDirectory moves a dotfile directory from home to stow package structure
func importDotfileDirectory(pkgName, sourcePath, destPath string) error {
	// Create the stow package directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create stow package directory: %v", err)
	}

	// Move the source directory to the stow package
	targetPath := filepath.Join(destPath, "."+pkgName)

	// Create backup of original if it already exists
	if _, err := os.Stat(targetPath); err == nil {
		backupPath := targetPath + ".backup." + fmt.Sprintf("%d", time.Now().Unix())
		if err := os.Rename(targetPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing directory: %v", err)
		}
		fmt.Printf("   📋 Backed up existing %s to %s\n", targetPath, backupPath)
	}

	// Move the original directory to the stow package
	if err := os.Rename(sourcePath, targetPath); err != nil {
		return fmt.Errorf("failed to move directory: %v", err)
	}

	return nil
}

func findStowConflicts(pkgPath, target string) []string {
	var conflicts []string

	err := filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the package directory itself
		if path == pkgPath {
			return nil
		}

		// Get relative path from package directory
		relPath, err := filepath.Rel(pkgPath, path)
		if err != nil {
			return err
		}

		// Calculate target path
		targetPath := filepath.Join(target, relPath)

		// Check if target exists and is not a symlink to our package
		if info, err := os.Lstat(targetPath); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				// It's a symlink, check if it points to our package
				linkTarget, err := os.Readlink(targetPath)
				if err == nil {
					if absLink, err := filepath.Abs(filepath.Join(filepath.Dir(targetPath), linkTarget)); err == nil {
						if absPackage, err := filepath.Abs(path); err == nil {
							if absLink == absPackage {
								// It's already linked to our package, not a conflict
								return nil
							}
						}
					}
				}
			}
			// File exists and is not a proper symlink to our package
			conflicts = append(conflicts, targetPath)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Warning: Error scanning for conflicts: %v\n", err)
	}

	return conflicts
}

func resolveStowConflicts(conflicts []string, backup, verbose bool) error {
	backupDir := ""
	if backup {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %v", err)
		}

		backupDir = filepath.Join(home, ".dotfiles", "backups", fmt.Sprintf("backup-%d", time.Now().Unix()))
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			return fmt.Errorf("error creating backup directory: %v", err)
		}
		fmt.Printf("📦 Created backup directory: %s\n", backupDir)
	}

	for _, conflictPath := range conflicts {
		if backup {
			// Backup the existing file
			relPath, err := filepath.Rel(os.Getenv("HOME"), conflictPath)
			if err != nil {
				relPath = filepath.Base(conflictPath)
			}

			backupPath := filepath.Join(backupDir, relPath)
			if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
				return fmt.Errorf("error creating backup subdirectory: %v", err)
			}

			if err := os.Rename(conflictPath, backupPath); err != nil {
				return fmt.Errorf("error backing up %s: %v", conflictPath, err)
			}

			if verbose {
				fmt.Printf("   📦 Backed up %s to %s\n", conflictPath, backupPath)
			}
		} else {
			// Just remove the conflicting file
			if err := os.Remove(conflictPath); err != nil {
				return fmt.Errorf("error removing conflicting file %s: %v", conflictPath, err)
			}

			if verbose {
				fmt.Printf("   🗑️  Removed conflicting file %s\n", conflictPath)
			}
		}
	}

	if backup && len(conflicts) > 0 {
		fmt.Printf("✅ Backed up %d conflicting files to %s\n", len(conflicts), backupDir)
	} else if len(conflicts) > 0 {
		fmt.Printf("✅ Removed %d conflicting files\n", len(conflicts))
	}

	return nil
}