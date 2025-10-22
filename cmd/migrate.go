package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	GroupID: "advanced",
	Short:   "🔄 Migrate from other dotfiles managers",
	Long: `🔄 Migration Tools

Import configurations from other dotfiles managers:
• chezmoi
• yadm
• GNU Stow (bare repo)
• Homesick
• dotbot

Examples:
  dotfiles migrate from chezmoi        # Interactive migration
  dotfiles migrate from yadm
  dotfiles migrate detect              # Auto-detect tools`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles migrate from <tool>' to migrate from another tool")
		fmt.Println("Use 'dotfiles migrate detect' to auto-detect installed tools")
	},
}

var migrateFromCmd = &cobra.Command{
	Use:   "from <tool>",
	Short: "Migrate from a specific tool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tool := strings.ToLower(args[0])

		switch tool {
		case "chezmoi":
			migrateFromChezmoi()
		case "yadm":
			migrateFromYadm()
		case "stow", "bare", "bare-repo":
			migrateFromBareRepo()
		case "homesick":
			migrateFromHomesick()
		case "dotbot":
			migrateFromDotbot()
		default:
			fmt.Printf("❌ Unknown tool: %s\n", tool)
			fmt.Println()
			fmt.Println("Supported tools:")
			fmt.Println("  • chezmoi")
			fmt.Println("  • yadm")
			fmt.Println("  • stow (bare-repo)")
			fmt.Println("  • homesick")
			fmt.Println("  • dotbot")
		}
	},
}

var migrateDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Auto-detect installed dotfiles managers",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()

		fmt.Println("🔍 Scanning for dotfiles managers...")
		fmt.Println()

		found := false

		// Check for chezmoi
		chezmoiDir := filepath.Join(home, ".local/share/chezmoi")
		if _, err := os.Stat(chezmoiDir); err == nil {
			fmt.Println("✅ Found: chezmoi")
			fmt.Printf("   Location: %s\n", chezmoiDir)
			fmt.Println("   Migrate with: dotfiles migrate from chezmoi")
			fmt.Println()
			found = true
		}

		// Check for yadm
		yadmDir := filepath.Join(home, ".config/yadm")
		if _, err := os.Stat(yadmDir); err == nil {
			fmt.Println("✅ Found: yadm")
			fmt.Printf("   Location: %s\n", yadmDir)
			fmt.Println("   Migrate with: dotfiles migrate from yadm")
			fmt.Println()
			found = true
		}

		// Check for bare repo
		bareRepo := filepath.Join(home, ".cfg")
		if _, err := os.Stat(bareRepo); err == nil {
			fmt.Println("✅ Found: Bare Git Repository")
			fmt.Printf("   Location: %s\n", bareRepo)
			fmt.Println("   Migrate with: dotfiles migrate from bare-repo")
			fmt.Println()
			found = true
		}

		// Check for homesick
		homesickDir := filepath.Join(home, ".homesick/repos")
		if _, err := os.Stat(homesickDir); err == nil {
			fmt.Println("✅ Found: Homesick")
			fmt.Printf("   Location: %s\n", homesickDir)
			fmt.Println("   Migrate with: dotfiles migrate from homesick")
			fmt.Println()
			found = true
		}

		// Check for dotbot
		dotbotConfig := filepath.Join(home, ".dotfiles/install.conf.yaml")
		if _, err := os.Stat(dotbotConfig); err == nil {
			fmt.Println("✅ Found: dotbot")
			fmt.Printf("   Config: %s\n", dotbotConfig)
			fmt.Println("   Migrate with: dotfiles migrate from dotbot")
			fmt.Println()
			found = true
		}

		if !found {
			fmt.Println("⚪ No dotfiles managers detected")
		}
	},
}

func migrateFromChezmoi() {
	home, _ := os.UserHomeDir()
	chezmoiDir := filepath.Join(home, ".local/share/chezmoi")

	if _, err := os.Stat(chezmoiDir); os.IsNotExist(err) {
		fmt.Println("❌ chezmoi directory not found")
		fmt.Printf("   Expected: %s\n", chezmoiDir)
		return
	}

	fmt.Println("🔄 Migrating from chezmoi...")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("This will copy files from chezmoi to ~/.dotfiles/stow. Continue? (y/n): ")
	response, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Println("Cancelled")
		return
	}

	dotfilesDir := filepath.Join(home, ".dotfiles")
	stowDir := filepath.Join(dotfilesDir, "stow")

	// Create dotfiles directory
	if err := os.MkdirAll(stowDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create dotfiles directory: %v\n", err)
		return
	}

	// Use chezmoi to list managed files
	cmd := exec.Command("chezmoi", "managed")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed to get chezmoi managed files: %v\n", err)
		return
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Printf("Found %d managed files\n", len(files))
	fmt.Println()

	// Group files by "package" (first directory component)
	packages := make(map[string][]string)
	for _, file := range files {
		if file == "" {
			continue
		}

		// Determine package name from file path
		parts := strings.Split(file, string(filepath.Separator))
		pkgName := "misc"
		if len(parts) > 0 && parts[0] != "" {
			if strings.HasPrefix(parts[0], ".") {
				pkgName = strings.TrimPrefix(parts[0], ".")
			}
		}

		packages[pkgName] = append(packages[pkgName], file)
	}

	copied := 0
	for pkg, pkgFiles := range packages {
		pkgDir := filepath.Join(stowDir, pkg)
		os.MkdirAll(pkgDir, 0755)

		for _, file := range pkgFiles {
			srcPath := filepath.Join(home, file)
			dstPath := filepath.Join(pkgDir, file)

			// Create directory structure
			dstDir := filepath.Dir(dstPath)
			if err := os.MkdirAll(dstDir, 0755); err != nil {
				continue
			}

			// Copy file
			if data, err := os.ReadFile(srcPath); err == nil {
				if err := os.WriteFile(dstPath, data, 0644); err == nil {
					copied++
				}
			}
		}

		fmt.Printf("✅ Created package: %s (%d files)\n", pkg, len(pkgFiles))
	}

	fmt.Println()
	fmt.Printf("✅ Migrated %d files from chezmoi\n", copied)
	fmt.Println()
	fmt.Println("💡 Next steps:")
	fmt.Println("   • Review files in ~/.dotfiles/stow/")
	fmt.Println("   • Run 'dotfiles scan' to update config")
	fmt.Println("   • Run 'dotfiles stow' to create symlinks")
}

func migrateFromYadm() {
	home, _ := os.UserHomeDir()
	yadmRepo := filepath.Join(home, ".config/yadm/repo.git")

	if _, err := os.Stat(yadmRepo); os.IsNotExist(err) {
		fmt.Println("❌ yadm repository not found")
		return
	}

	fmt.Println("🔄 Migrating from yadm...")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("This will copy yadm-managed files to ~/.dotfiles/stow. Continue? (y/n): ")
	response, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Println("Cancelled")
		return
	}

	dotfilesDir := filepath.Join(home, ".dotfiles")
	stowDir := filepath.Join(dotfilesDir, "stow")
	os.MkdirAll(stowDir, 0755)

	// Get yadm list
	cmd := exec.Command("yadm", "list", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed to list yadm files: %v\n", err)
		return
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Printf("Found %d files\n", len(files))
	fmt.Println()

	// Copy to misc package
	miscDir := filepath.Join(stowDir, "misc")
	os.MkdirAll(miscDir, 0755)

	copied := 0
	for _, file := range files {
		if file == "" {
			continue
		}

		srcPath := filepath.Join(home, file)
		dstPath := filepath.Join(miscDir, file)

		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			continue
		}

		if data, err := os.ReadFile(srcPath); err == nil {
			if err := os.WriteFile(dstPath, data, 0644); err == nil {
				copied++
			}
		}
	}

	fmt.Printf("✅ Migrated %d files to misc package\n", copied)
	fmt.Println()
	fmt.Println("💡 Next steps:")
	fmt.Println("   • Organize files into logical packages")
	fmt.Println("   • Run 'dotfiles stow' to create symlinks")
}

func migrateFromBareRepo() {
	home, _ := os.UserHomeDir()
	bareRepo := filepath.Join(home, ".cfg")

	if _, err := os.Stat(bareRepo); os.IsNotExist(err) {
		fmt.Println("❌ Bare repository not found at ~/.cfg")
		return
	}

	fmt.Println("🔄 Migrating from bare git repository...")
	fmt.Println()

	dotfilesDir := filepath.Join(home, ".dotfiles")
	if _, err := os.Stat(dotfilesDir); err == nil {
		fmt.Println("❌ ~/.dotfiles already exists")
		fmt.Println("   Please backup or remove it first")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Move bare repo to ~/.dotfiles? (y/n): ")
	response, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Println("Cancelled")
		return
	}

	// Move bare repo
	if err := os.Rename(bareRepo, dotfilesDir); err != nil {
		fmt.Printf("❌ Failed to move repository: %v\n", err)
		return
	}

	fmt.Println("✅ Repository moved to ~/.dotfiles")
	fmt.Println()
	fmt.Println("💡 Next steps:")
	fmt.Println("   • Organize files into stow packages")
	fmt.Println("   • Run 'dotfiles init' to set up config.json")
}

func migrateFromHomesick() {
	home, _ := os.UserHomeDir()
	homesickDir := filepath.Join(home, ".homesick/repos")

	if _, err := os.Stat(homesickDir); os.IsNotExist(err) {
		fmt.Println("❌ Homesick directory not found")
		return
	}

	fmt.Println("🔄 Migrating from Homesick...")
	fmt.Println()

	repos, _ := os.ReadDir(homesickDir)
	if len(repos) == 0 {
		fmt.Println("No Homesick castles found")
		return
	}

	fmt.Println("Found castles:")
	for _, repo := range repos {
		if repo.IsDir() {
			fmt.Printf("  • %s\n", repo.Name())
		}
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Copy to ~/.dotfiles/stow? (y/n): ")
	response, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Println("Cancelled")
		return
	}

	dotfilesDir := filepath.Join(home, ".dotfiles/stow")
	os.MkdirAll(dotfilesDir, 0755)

	for _, repo := range repos {
		if !repo.IsDir() {
			continue
		}

		srcDir := filepath.Join(homesickDir, repo.Name(), "home")
		dstDir := filepath.Join(dotfilesDir, repo.Name())

		if _, err := os.Stat(srcDir); err == nil {
			exec.Command("cp", "-r", srcDir, dstDir).Run()
			fmt.Printf("✅ Copied: %s\n", repo.Name())
		}
	}

	fmt.Println()
	fmt.Println("💡 Migration complete. Run 'dotfiles stow' to create symlinks")
}

func migrateFromDotbot() {
	home, _ := os.UserHomeDir()
	dotbotConfig := filepath.Join(home, ".dotfiles/install.conf.yaml")

	if _, err := os.Stat(dotbotConfig); os.IsNotExist(err) {
		fmt.Println("❌ dotbot config not found")
		return
	}

	fmt.Println("🔄 Migrating from dotbot...")
	fmt.Println()
	fmt.Println("⚠️  Automatic migration from dotbot is limited")
	fmt.Println("   Please manually review your install.conf.yaml")
	fmt.Println()
	fmt.Println("💡 dotbot link directives can usually be converted to stow packages")
	fmt.Println("   Organize your dotfiles into ~/.dotfiles/stow/<package>/")
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateFromCmd)
	migrateCmd.AddCommand(migrateDetectCmd)
}
