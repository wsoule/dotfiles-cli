package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"dotfiles/internal/config"
	"dotfiles/internal/pkgmanager"
	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "🎯 Complete developer onboarding and environment setup",
	Long: `🎯 Developer Onboarding - Complete Environment Setup

Perfect for new developers or setting up fresh machines. This single command will:

1. 🔧 Initialize your dotfiles configuration
2. 🔒 Create private directory for sensitive files (SSH keys, env vars)
3. 🐚 Set up shell configuration with zsh and aliases
4. 🔐 Generate GitHub SSH keys and show setup instructions
5. 📦 Install curated essential development packages
6. 📋 Guide you through next steps

Essential packages included:
• CLI Tools: git, curl, wget, tree, jq, stow, gh
• Applications: Visual Studio Code, Ghostty, Raycast
• Fonts: JetBrains Mono, Ubuntu Mono (Nerd Font variants)

Examples:
  dotfiles onboard                           # Full interactive setup
  dotfiles onboard --email you@email.com    # With GitHub email
  dotfiles onboard --skip-packages          # Skip package installation
  dotfiles onboard --skip-interactive       # Use defaults, no prompts`,
	Run: func(cmd *cobra.Command, args []string) {
		skipInteractive, _ := cmd.Flags().GetBool("skip-interactive")
		skipGithub, _ := cmd.Flags().GetBool("skip-github")
		skipPackages, _ := cmd.Flags().GetBool("skip-packages")
		email, _ := cmd.Flags().GetString("email")

		fmt.Println("🎉 Welcome to Dotfiles Manager - Developer Onboarding!")
		fmt.Println("=" + strings.Repeat("=", 55))
		fmt.Println()
		fmt.Println("This wizard will help you set up your development environment:")
		fmt.Println("✅ Initialize dotfiles configuration")
		fmt.Println("🔐 Set up GitHub SSH authentication")
		fmt.Println("📦 Install essential development packages")
		fmt.Println("🔗 Configure dotfiles with Stow")
		fmt.Println()

		if !skipInteractive && !askConfirmation("Ready to begin? (Y/n): ", true) {
			fmt.Println("👋 Setup cancelled. Run 'dotfiles onboard' again when ready!")
			return
		}

		fmt.Println()
		fmt.Println("🚀 Starting onboarding process...")
		fmt.Println()

		// Step 1: Detect existing dotfiles
		fmt.Println("🔍 Step 1: Scanning for existing dotfiles...")
		existingDotfiles := detectExistingDotfiles()
		if len(existingDotfiles) > 0 {
			fmt.Printf("   Found %d existing dotfiles:\n", len(existingDotfiles))
			for _, dotfile := range existingDotfiles {
				fmt.Printf("   • %s\n", dotfile)
			}
			fmt.Println()

			if !skipInteractive && askConfirmation("   Would you like to import these into your dotfiles setup? (Y/n): ", true) {
				if err := offerDotfilesImport(existingDotfiles, skipInteractive); err != nil {
					fmt.Printf("⚠️  Some imports failed: %v\n", err)
				}
			}
		} else {
			fmt.Println("   No existing dotfiles found")
		}
		fmt.Println()

		// Step 2: Check dependencies
		fmt.Println("🔧 Step 2: Checking dependencies...")
		if err := checkAndInstallDependencies(skipInteractive); err != nil {
			fmt.Printf("⚠️  Dependency check had issues: %v\n", err)
		}
		fmt.Println()

		// Step 3: Initialize configuration
		fmt.Println("📋 Step 3: Initializing dotfiles configuration...")
		if err := initializeConfig(); err != nil {
			fmt.Printf("❌ Failed to initialize configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Configuration initialized!")

		// Set up complete environment (private dir + shell packages + stow)
		fmt.Println("🔒 Setting up dotfiles environment...")
		home, _ := os.UserHomeDir()
		dotfilesDir := filepath.Join(home, ".dotfiles")
		if err := setupCompleteEnvironment(dotfilesDir, true); err != nil {
			fmt.Printf("⚠️  Environment setup had issues: %v\n", err)
		} else {
			fmt.Println("✅ Environment setup complete!")
		}
		fmt.Println()

		// Step 4: Scan for installed packages
		fmt.Println("📦 Step 4: Scanning for installed packages...")
		if err := scanAndOfferPackages(skipInteractive); err != nil {
			fmt.Printf("⚠️  Package scan had issues: %v\n", err)
		}
		fmt.Println()

		// Step 5: GitHub setup
		if !skipGithub {
			fmt.Println("🔐 Step 5: Setting up GitHub SSH authentication...")
			if email == "" && !skipInteractive {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter your GitHub email: ")
				email, _ = reader.ReadString('\n')
				email = strings.TrimSpace(email)
			}

			if email != "" {
				if err := setupGitHubSSH(email); err != nil {
					fmt.Printf("⚠️  GitHub setup had issues: %v\n", err)
					fmt.Println("   You can run 'dotfiles github setup' later to complete this.")
				} else {
					fmt.Println("✅ GitHub SSH setup completed!")
				}
			} else {
				fmt.Println("⚠️  Skipping GitHub setup (no email provided)")
				fmt.Println("   Run 'dotfiles github setup --email=your@email.com' later")
			}
			fmt.Println()
		}

		// Step 6: Install essential packages
		if !skipPackages {
			fmt.Println("📦 Step 6: Installing essential development packages...")
			if err := installEssentialPackages(skipInteractive); err != nil {
				fmt.Printf("⚠️  Package installation had issues: %v\n", err)
			} else {
				fmt.Println("✅ Essential packages installed!")
			}
			fmt.Println()
		}

		// Step 7: Final steps and guidance
		fmt.Println("🎯 Step 7: Final setup and next steps...")
		showNextSteps()
		fmt.Println()

		fmt.Println("🎉 Onboarding complete! Your development environment is ready.")
		fmt.Println()
		fmt.Println("💡 Useful commands to remember:")
		fmt.Println("   dotfiles --help                 # See all available commands")
		fmt.Println("   dotfiles add <package>          # Add packages to your config")
		fmt.Println("   dotfiles status                 # Check installation status")
		fmt.Println("   dotfiles github test            # Test GitHub connection")
		fmt.Println("   dotfiles stow <package>         # Stow dotfiles")
		fmt.Println()
		fmt.Println("Happy coding! 🚀")
	},
}

func initializeConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".dotfiles", "config.json")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("   Configuration already exists at %s\n", configPath)
		return nil
	}

	// Create config directory
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Create initial config
	cfg := &config.Config{
		Brews: []string{},
		Casks: []string{},
		Taps:  []string{},
		Stow:  []string{},
	}

	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save initial config: %v", err)
	}

	fmt.Printf("   Created configuration at %s\n", configPath)
	return nil
}

func setupGitHubSSH(email string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sshDir := filepath.Join(home, ".ssh")
	keyPath := filepath.Join(sshDir, "id_ed25519")
	pubKeyPath := keyPath + ".pub"

	// Create .ssh directory
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %v", err)
	}

	// Check if key already exists
	if _, err := os.Stat(keyPath); err == nil {
		fmt.Println("   SSH key already exists, skipping generation")
		showSSHInstructions(pubKeyPath)
		return nil
	}

	// Generate SSH key
	fmt.Println("   Generating SSH key...")
	sshKeygenCmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", keyPath, "-N", "")
	if err := sshKeygenCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate SSH key: %v", err)
	}

	// Set permissions
	os.Chmod(keyPath, 0600)
	os.Chmod(pubKeyPath, 0644)

	fmt.Println("   SSH key generated successfully!")
	showSSHInstructions(pubKeyPath)
	return nil
}

func showSSHInstructions(pubKeyPath string) {
	fmt.Println()
	fmt.Println("   📌 IMPORTANT: Add your SSH key to GitHub:")
	fmt.Println("   1. Copy your public key:")

	pubKeyContent, err := os.ReadFile(pubKeyPath)
	if err != nil {
		fmt.Printf("   ❌ Error reading public key: %v\n", err)
		return
	}

	fmt.Println("   " + strings.Repeat("-", 40))
	fmt.Printf("   %s", string(pubKeyContent))
	fmt.Println("   " + strings.Repeat("-", 40))

	fmt.Println("   2. Go to: https://github.com/settings/ssh/new")
	fmt.Println("   3. Paste the key and give it a title")
	fmt.Println("   4. Test with: dotfiles github test")

	// Try to copy to clipboard
	if err := copyToClipboard(string(pubKeyContent)); err == nil {
		fmt.Println("   📋 Public key copied to clipboard!")
	}
}

func installEssentialPackages(skipInteractive bool) error {
	// Get platform-specific essential packages
	essentialPackages := getEssentialPackages()

	// Show packages that will be installed

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".dotfiles", "config.json")
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	fmt.Println("   The following essential packages will be added:")
	if len(essentialPackages["taps"]) > 0 {
		fmt.Printf("   📋 Taps: %s\n", strings.Join(essentialPackages["taps"], ", "))
	}
	fmt.Printf("   🍺 Packages: %s\n", strings.Join(essentialPackages["brews"], ", "))
	if len(essentialPackages["casks"]) > 0 {
		fmt.Printf("   📦 Apps: %s\n", strings.Join(essentialPackages["casks"], ", "))
	}
	fmt.Println()

	if !skipInteractive && !askConfirmation("   Continue with package installation? (Y/n): ", true) {
		fmt.Println("   Skipping package installation")
		return nil
	}

	// Add packages to config
	for _, tap := range essentialPackages["taps"] {
		if !contains(cfg.Taps, tap) {
			cfg.Taps = append(cfg.Taps, tap)
		}
	}
	for _, brew := range essentialPackages["brews"] {
		if !contains(cfg.Brews, brew) {
			cfg.Brews = append(cfg.Brews, brew)
		}
	}
	for _, cask := range essentialPackages["casks"] {
		if !contains(cfg.Casks, cask) {
			cfg.Casks = append(cfg.Casks, cask)
		}
	}

	// Save updated config
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("   Packages added to configuration")

	// Install packages
	fmt.Println("   Installing packages with Homebrew...")
	if err := runInstallCommand(); err != nil {
		return fmt.Errorf("package installation failed: %v", err)
	}

	return nil
}

func runInstallCommand() error {
	// Generate Brewfile
	brewfileCmd := exec.Command("dotfiles", "brewfile")
	if err := brewfileCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate Brewfile: %v", err)
	}

	// Install with Homebrew
	fmt.Println("   Running: brew bundle --file=./Brewfile")
	brewBundleCmd := exec.Command("brew", "bundle", "--file=./Brewfile")
	brewBundleCmd.Stdout = os.Stdout
	brewBundleCmd.Stderr = os.Stderr

	if err := brewBundleCmd.Run(); err != nil {
		return fmt.Errorf("brew bundle failed: %v", err)
	}

	return nil
}

func showNextSteps() {
	fmt.Println("   🔧 Recommended next steps:")
	fmt.Println("   • Create dotfile packages in ~/.dotfiles/ (vim, zsh, tmux, etc.)")
	fmt.Println("   • Add them with: dotfiles add --type=stow <package>")
	fmt.Println("   • Stow them with: dotfiles stow <package>")
	fmt.Println("   • Customize your package list with: dotfiles add <package>")
	fmt.Println("   • Check status anytime with: dotfiles status")
}

func askConfirmation(prompt string, defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "" {
			return defaultYes
		}

		if response == "y" || response == "yes" {
			return true
		}

		if response == "n" || response == "no" {
			return false
		}

		fmt.Println("   Please answer y/yes or n/no.")
	}
}

func detectExistingDotfiles() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	commonDotfiles := []string{
		".zshrc", ".bashrc", ".bash_profile", ".profile",
		".vimrc", ".vim", ".nvim", ".config/nvim",
		".tmux.conf", ".gitconfig", ".gitignore_global",
		".aliases", ".functions", ".exports",
		".ssh/config", ".aws", ".docker",
	}

	var existing []string
	for _, dotfile := range commonDotfiles {
		path := filepath.Join(home, dotfile)
		if _, err := os.Stat(path); err == nil {
			existing = append(existing, dotfile)
		}
	}

	return existing
}

func offerDotfilesImport(dotfiles []string, skipInteractive bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dotfilesDir := filepath.Join(home, ".dotfiles")
	stowDir := filepath.Join(dotfilesDir, "stow")

	// Ensure directories exist
	if err := os.MkdirAll(stowDir, 0755); err != nil {
		return fmt.Errorf("failed to create stow directory: %v", err)
	}

	for _, dotfile := range dotfiles {
		// Determine package name from dotfile
		pkgName := suggestPackageName(dotfile)

		if !skipInteractive {
			fmt.Printf("   Import %s into package '%s'? (Y/n/s=skip): ", dotfile, pkgName)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response == "n" || response == "no" {
				// Ask for custom package name
				fmt.Print("   Enter custom package name (or 'skip'): ")
				customName, _ := reader.ReadString('\n')
				customName = strings.TrimSpace(customName)
				if customName == "skip" || customName == "" {
					continue
				}
				pkgName = customName
			} else if response == "s" || response == "skip" {
				continue
			}
		}

		// Import the dotfile
		if err := importDotfileToPackage(dotfile, pkgName, home, stowDir); err != nil {
			fmt.Printf("   ❌ Failed to import %s: %v\n", dotfile, err)
			continue
		}

		fmt.Printf("   ✅ Imported %s into package '%s'\n", dotfile, pkgName)
	}

	return nil
}

func suggestPackageName(dotfile string) string {
	// Suggest logical package names based on dotfile
	packageMap := map[string]string{
		".zshrc":         "zsh",
		".bashrc":        "bash",
		".bash_profile":  "bash",
		".profile":       "shell",
		".vimrc":         "vim",
		".vim":           "vim",
		".nvim":          "nvim",
		".config/nvim":   "nvim",
		".tmux.conf":     "tmux",
		".gitconfig":     "git",
		".gitignore_global": "git",
		".aliases":       "shell",
		".functions":     "shell",
		".exports":       "shell",
		".ssh/config":    "ssh",
		".aws":           "aws",
		".docker":        "docker",
	}

	if pkg, exists := packageMap[dotfile]; exists {
		return pkg
	}

	// Default: extract name from dotfile
	name := strings.TrimPrefix(dotfile, ".")
	name = strings.Split(name, "/")[0]
	return name
}

func importDotfileToPackage(dotfile, pkgName, home, stowDir string) error {
	srcPath := filepath.Join(home, dotfile)
	pkgDir := filepath.Join(stowDir, pkgName)
	dstPath := filepath.Join(pkgDir, dotfile)

	// Create package directory
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create package directory: %v", err)
	}

	// Move or copy the file/directory
	if err := os.Rename(srcPath, dstPath); err != nil {
		// If rename fails, try copy
		return copyFileOrDir(srcPath, dstPath)
	}

	return nil
}

func copyFileOrDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	return copyFilesOnboard(src, dst)
}

func copyFilesOnboard(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func getEssentialPackages() map[string][]string {
	if runtime.GOOS == "darwin" {
		// macOS packages
		return map[string][]string{
			"taps": {},
			"brews": {
				"git",
				"tree",
				"jq",
				"stow",
			},
			"casks": {
				"visual-studio-code",
				"ghostty",
				"raycast",
			},
		}
	} else {
		// Linux packages (generic - works for Arch, Debian, etc.)
		return map[string][]string{
			"taps": {},
			"brews": {
				"git",
				"tree",
				"jq",
				"stow",
				"base-devel", // For Arch
			},
			"casks": {}, // No casks on Linux
		}
	}
}

func checkAndInstallDependencies(skipInteractive bool) error {
	var dependencies map[string]struct {
		cmd         string
		installCmd  string
		description string
	}

	if runtime.GOOS == "darwin" {
		// macOS dependencies
		dependencies = map[string]struct {
			cmd         string
			installCmd  string
			description string
		}{
			"brew": {
				cmd:         "brew",
				installCmd:  `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
				description: "Homebrew package manager",
			},
			"git": {
				cmd:         "git",
				installCmd:  "brew install git",
				description: "Git version control",
			},
			"stow": {
				cmd:         "stow",
				installCmd:  "brew install stow",
				description: "GNU Stow for dotfiles management",
			},
		}
	} else {
		// Linux dependencies
		pm, _ := pkgmanager.GetPackageManager()
		pmName := "package manager"
		installGit := "sudo pacman -S git"
		installStow := "sudo pacman -S stow"

		if pm != nil {
			pmName = pm.GetName()
			if pmName == "apt" {
				installGit = "sudo apt-get install git"
				installStow = "sudo apt-get install stow"
			} else if pmName == "yum" {
				installGit = "sudo yum install git"
				installStow = "sudo yum install stow"
			}
		}

		dependencies = map[string]struct {
			cmd         string
			installCmd  string
			description string
		}{
			pmName: {
				cmd:         pmName,
				installCmd:  "N/A - should be pre-installed",
				description: "Package manager",
			},
			"git": {
				cmd:         "git",
				installCmd:  installGit,
				description: "Git version control",
			},
			"stow": {
				cmd:         "stow",
				installCmd:  installStow,
				description: "GNU Stow for dotfiles management",
			},
		}
	}

	var missing []string
	for name, dep := range dependencies {
		if _, err := exec.LookPath(dep.cmd); err != nil {
			missing = append(missing, name)
			fmt.Printf("   ❌ Missing: %s (%s)\n", name, dep.description)
		} else {
			fmt.Printf("   ✅ Found: %s\n", name)
		}
	}

	if len(missing) == 0 {
		fmt.Println("   All dependencies satisfied!")
		return nil
	}

	if skipInteractive {
		fmt.Printf("   ⚠️  Missing %d dependencies. Install manually.\n", len(missing))
		return nil
	}

	fmt.Printf("   Install missing dependencies? (Y/n): ")
	if !askConfirmation("", true) {
		return nil
	}

	for _, name := range missing {
		dep := dependencies[name]
		fmt.Printf("   Installing %s...\n", name)

		cmd := exec.Command("bash", "-c", dep.installCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("   ❌ Failed to install %s: %v\n", name, err)
		} else {
			fmt.Printf("   ✅ Installed %s\n", name)
		}
	}

	return nil
}

func scanAndOfferPackages(skipInteractive bool) error {
	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		fmt.Println("   ⚠️  Homebrew not installed, skipping package scan")
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".dotfiles", "config.json")
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Get installed packages
	installedBrews, err := getInstalledBrews()
	if err != nil {
		return fmt.Errorf("failed to scan brews: %v", err)
	}

	installedCasks, err := getInstalledCasks()
	if err != nil {
		return fmt.Errorf("failed to scan casks: %v", err)
	}

	// Filter out packages already in config
	newBrews := filterNewPackages(installedBrews, cfg.Brews)
	newCasks := filterNewPackages(installedCasks, cfg.Casks)

	if len(newBrews) == 0 && len(newCasks) == 0 {
		fmt.Println("   ✅ No new packages found (all installed packages already in config)")
		return nil
	}

	fmt.Printf("   Found %d brews and %d casks not in your config\n", len(newBrews), len(newCasks))

	if skipInteractive {
		fmt.Println("   Skipping package import (use --skip-interactive=false or run 'dotfiles scan' later)")
		return nil
	}

	if !askConfirmation("   Would you like to add these to your config? (Y/n): ", true) {
		fmt.Println("   You can run 'dotfiles scan' later to add them")
		return nil
	}

	// Show brief preview
	if len(newBrews) > 0 {
		fmt.Printf("   📋 Brews: %s", strings.Join(newBrews[:min(3, len(newBrews))], ", "))
		if len(newBrews) > 3 {
			fmt.Printf(" ... and %d more", len(newBrews)-3)
		}
		fmt.Println()
	}

	if len(newCasks) > 0 {
		fmt.Printf("   📦 Casks: %s", strings.Join(newCasks[:min(3, len(newCasks))], ", "))
		if len(newCasks) > 3 {
			fmt.Printf(" ... and %d more", len(newCasks)-3)
		}
		fmt.Println()
	}

	if askConfirmation("   Add all packages? (Y/n): ", true) {
		cfg.Brews = append(cfg.Brews, newBrews...)
		cfg.Casks = append(cfg.Casks, newCasks...)

		if err := cfg.Save(configPath); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Printf("   ✅ Added %d brews and %d casks to config\n", len(newBrews), len(newCasks))
	} else {
		fmt.Println("   Run 'dotfiles scan' later to selectively add packages")
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Note: contains function is already defined in add.go

func init() {
	onboardCmd.Flags().Bool("skip-interactive", false, "Skip interactive prompts (use defaults)")
	onboardCmd.Flags().Bool("skip-github", false, "Skip GitHub SSH setup")
	onboardCmd.Flags().Bool("skip-packages", false, "Skip essential package installation")
	onboardCmd.Flags().StringP("email", "e", "", "Email for GitHub SSH key")

	rootCmd.AddCommand(onboardCmd)
}

