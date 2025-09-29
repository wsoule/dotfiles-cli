package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Complete developer onboarding and environment setup",
	Long:  `Interactive setup wizard for new developers to configure their entire development environment`,
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

		// Step 1: Initialize configuration
		fmt.Println("📋 Step 1: Initializing dotfiles configuration...")
		if err := initializeConfig(); err != nil {
			fmt.Printf("❌ Failed to initialize configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Configuration initialized!")
		fmt.Println()

		// Step 2: GitHub setup
		if !skipGithub {
			fmt.Println("🔐 Step 2: Setting up GitHub SSH authentication...")
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

		// Step 3: Install essential packages
		if !skipPackages {
			fmt.Println("📦 Step 3: Installing essential development packages...")
			if err := installEssentialPackages(skipInteractive); err != nil {
				fmt.Printf("⚠️  Package installation had issues: %v\n", err)
			} else {
				fmt.Println("✅ Essential packages installed!")
			}
			fmt.Println()
		}

		// Step 4: Final steps and guidance
		fmt.Println("🎯 Step 4: Final setup and next steps...")
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
	essentialPackages := map[string][]string{
		"brews": {
			"git",
			"curl",
			"wget",
			"tree",
			"jq",
			"stow",
			"gh",
		},
		"casks": {
			"visual-studio-code",
			"ghostty",
			"raycast",
			"font-jetbrains-mono-nerd-font",
			"font-ubuntu-mono-nerd-font",
		},
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

	// Show packages that will be installed
	fmt.Println("   The following essential packages will be added:")
	fmt.Printf("   📋 Taps: %s\n", strings.Join(essentialPackages["taps"], ", "))
	fmt.Printf("   🍺 Brews: %s\n", strings.Join(essentialPackages["brews"], ", "))
	fmt.Printf("   📦 Casks: %s\n", strings.Join(essentialPackages["casks"], ", "))
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
	brewfileCmd := exec.Command("./dotfiles", "brewfile")
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

// Note: contains function is already defined in add.go

func init() {
	onboardCmd.Flags().Bool("skip-interactive", false, "Skip interactive prompts (use defaults)")
	onboardCmd.Flags().Bool("skip-github", false, "Skip GitHub SSH setup")
	onboardCmd.Flags().Bool("skip-packages", false, "Skip essential package installation")
	onboardCmd.Flags().StringP("email", "e", "", "Email for GitHub SSH key")

	rootCmd.AddCommand(onboardCmd)
}

