package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"dotfiles/internal/config"
)

// Manager handles the installation process
type Manager struct {
	config *config.Config
	dryRun bool
}

// NewManager creates a new installer manager
func NewManager(cfg *config.Config, dryRun bool) *Manager {
	return &Manager{
		config: cfg,
		dryRun: dryRun,
	}
}

// InstallHomebrew installs Homebrew if not present
func (m *Manager) InstallHomebrew() error {
	fmt.Println("📦 Installing Homebrew...")

	// Check if Homebrew is already installed
	if m.commandExists("brew") {
		fmt.Println("✅ Homebrew is already installed")
		return nil
	}

	if m.dryRun {
		fmt.Println("🔍 [DRY RUN] Would install Homebrew")
		return nil
	}

	// Install Homebrew
	cmd := exec.Command("/bin/bash", "-c",
		"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Homebrew: %w", err)
	}

	// Add Homebrew to PATH
	if runtime.GOARCH == "arm64" {
		os.Setenv("PATH", "/opt/homebrew/bin:"+os.Getenv("PATH"))
	} else {
		os.Setenv("PATH", "/usr/local/bin:"+os.Getenv("PATH"))
	}

	fmt.Println("✅ Homebrew installed successfully")
	return nil
}

// InstallBrewPackages installs packages from Brewfile and config
func (m *Manager) InstallBrewPackages() error {
	fmt.Println("📦 Installing brew packages...")

	if m.dryRun {
		fmt.Println("🔍 [DRY RUN] Would install brew packages")
		return nil
	}

	// Install from Brewfile if it exists
	brewfilePath := filepath.Join(m.getDotfilesDir(), "Brewfile")
	if _, err := os.Stat(brewfilePath); err == nil {
		fmt.Println("Installing packages from Brewfile...")
		cmd := exec.Command("brew", "bundle", "install", "--file", brewfilePath, "--verbose")
		cmd.Dir = m.getDotfilesDir()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install from Brewfile: %w", err)
		}
	}

	// Install extra brew packages
	for _, pkg := range m.config.Packages.ExtraBrews {
		fmt.Printf("Installing brew package: %s\n", pkg)
		cmd := exec.Command("brew", "install", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to install %s: %v\n", pkg, err)
		}
	}

	// Install extra cask packages
	for _, pkg := range m.config.Packages.ExtraCasks {
		fmt.Printf("Installing cask: %s\n", pkg)
		cmd := exec.Command("brew", "install", "--cask", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to install %s: %v\n", pkg, err)
		}
	}

	// Add extra taps
	for _, tap := range m.config.Packages.ExtraTaps {
		fmt.Printf("Adding tap: %s\n", tap)
		cmd := exec.Command("brew", "tap", tap)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to add tap %s: %v\n", tap, err)
		}
	}

	fmt.Println("✅ Brew packages installation completed")
	return nil
}

// InstallDotfiles installs dotfiles using GNU Stow
func (m *Manager) InstallDotfiles() error {
	fmt.Println("🔗 Installing dotfiles using GNU Stow...")

	if !m.commandExists("stow") {
		return fmt.Errorf("GNU Stow is not available. Install it with: brew install stow")
	}

	if m.dryRun {
		fmt.Println("🔍 [DRY RUN] Would install dotfiles with Stow")
		return nil
	}

	dotfilesDir := m.getDotfilesDir()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Get list of packages to stow
	packages, err := m.getStowPackages(dotfilesDir)
	if err != nil {
		return fmt.Errorf("failed to get stow packages: %w", err)
	}

	// Stow each package
	for _, pkg := range packages {
		fmt.Printf("Stowing %s...\n", pkg)
		cmd := exec.Command("stow", "--target", homeDir, pkg, "--verbose=1")
		cmd.Dir = dotfilesDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			// Try with --adopt flag
			fmt.Printf("⚠️  Failed to stow %s, trying with --adopt...\n", pkg)
			cmd = exec.Command("stow", "--target", homeDir, pkg, "--adopt", "--verbose=1")
			cmd.Dir = dotfilesDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Printf("⚠️  Failed to stow %s: %v\n", pkg, err)
				continue
			}
		}
		fmt.Printf("✅ Stowed %s\n", pkg)
	}

	fmt.Println("✅ Dotfiles installation completed")
	return nil
}

// ApplyMacOSDefaults applies macOS system defaults
func (m *Manager) ApplyMacOSDefaults() error {
	fmt.Println("🍎 Applying macOS system defaults...")

	if runtime.GOOS != "darwin" {
		fmt.Println("⚠️  Skipping macOS defaults (not running on macOS)")
		return nil
	}

	if m.dryRun {
		fmt.Println("🔍 [DRY RUN] Would apply macOS defaults")
		return nil
	}

	cfg := m.config.System

	// Dark mode
	if cfg.Appearance.DarkMode {
		m.runDefaults("write", "NSGlobalDomain", "AppleInterfaceStyle", "-string", "Dark")
	} else {
		m.runDefaults("delete", "NSGlobalDomain", "AppleInterfaceStyle")
	}

	// 24-hour time
	if cfg.Appearance.Enable24HourTime {
		m.runDefaults("write", "NSGlobalDomain", "AppleICUForce24HourTime", "-bool", "true")
	}

	// Keyboard settings
	if cfg.Keyboard.DisablePressAndHold {
		m.runDefaults("write", "NSGlobalDomain", "ApplePressAndHoldEnabled", "-bool", "false")
	}
	m.runDefaults("write", "NSGlobalDomain", "KeyRepeat", "-int", fmt.Sprintf("%d", cfg.Keyboard.KeyRepeatRate))

	// Dock settings
	if cfg.Dock.Autohide {
		m.runDefaults("write", "com.apple.dock", "autohide", "-bool", "true")
	}
	m.runDefaults("write", "com.apple.dock", "orientation", "-string", cfg.Dock.Position)
	m.runDefaults("write", "com.apple.dock", "tilesize", "-int", fmt.Sprintf("%d", cfg.Dock.TileSize))

	// Finder settings
	if cfg.Finder.ShowHiddenFiles {
		m.runDefaults("write", "com.apple.finder", "AppleShowAllFiles", "-bool", "true")
	}
	if cfg.Finder.ShowFileExtensions {
		m.runDefaults("write", "NSGlobalDomain", "AppleShowAllExtensions", "-bool", "true")
	}

	// Set finder default view
	viewMap := map[string]string{
		"icon":    "icnv",
		"list":    "Nlsv",
		"column":  "clmv",
		"gallery": "Flwv",
	}
	if viewCode, ok := viewMap[cfg.Finder.DefaultView]; ok {
		m.runDefaults("write", "com.apple.finder", "FXPreferredViewStyle", "-string", viewCode)
	}

	// Security settings
	if cfg.Security.RequirePasswordImmediately {
		m.runDefaults("write", "com.apple.screensaver", "askForPassword", "-int", "1")
		m.runDefaults("write", "com.apple.screensaver", "askForPasswordDelay", "-int", "0")
	}

	// Restart affected applications
	m.killApp("Dock")
	m.killApp("Finder")

	fmt.Println("✅ macOS defaults applied")
	return nil
}

// InstallNPMPackages installs global npm packages
func (m *Manager) InstallNPMPackages() error {
	fmt.Println("📦 Installing global npm packages...")

	if !m.commandExists("npm") {
		fmt.Println("⚠️  npm not found, skipping npm packages")
		return nil
	}

	if m.dryRun {
		fmt.Println("🔍 [DRY RUN] Would install npm packages")
		return nil
	}

	for _, pkg := range m.config.Packages.NPMGlobals {
		fmt.Printf("Installing npm package: %s\n", pkg)
		cmd := exec.Command("npm", "install", "-g", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to install %s: %v\n", pkg, err)
		}
	}

	fmt.Println("✅ Global npm packages installed")
	return nil
}

// PostInstallSetup runs post-installation configuration
func (m *Manager) PostInstallSetup() error {
	fmt.Println("⚙️  Running post-install configuration...")

	if m.dryRun {
		fmt.Println("🔍 [DRY RUN] Would run post-install setup")
		return nil
	}

	// Set up git configuration
	if m.config.Personal.Name != "" {
		cmd := exec.Command("git", "config", "--global", "user.name", m.config.Personal.Name)
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to set git user.name: %v\n", err)
		} else {
			fmt.Printf("✅ Git user.name set to: %s\n", m.config.Personal.Name)
		}
	}

	if m.config.Personal.Email != "" {
		cmd := exec.Command("git", "config", "--global", "user.email", m.config.Personal.Email)
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to set git user.email: %v\n", err)
		} else {
			fmt.Printf("✅ Git user.email set to: %s\n", m.config.Personal.Email)
		}
	}

	// Additional git settings
	exec.Command("git", "config", "--global", "init.defaultBranch", m.config.Development.Git.DefaultBranch).Run()
	exec.Command("git", "config", "--global", "push.default", m.config.Development.Git.PushDefault).Run()

	if m.config.Development.Git.PullRebase {
		exec.Command("git", "config", "--global", "pull.rebase", "true").Run()
	}

	// Create directories
	homeDir, _ := os.UserHomeDir()
	for _, dir := range m.config.Directories {
		expandedDir := os.ExpandEnv(dir)
		if strings.HasPrefix(expandedDir, "~") {
			expandedDir = strings.Replace(expandedDir, "~", homeDir, 1)
		}

		if err := os.MkdirAll(expandedDir, 0755); err != nil {
			fmt.Printf("⚠️  Failed to create directory %s: %v\n", expandedDir, err)
		} else {
			fmt.Printf("✅ Created directory: %s\n", expandedDir)
		}
	}

	fmt.Println("✅ Post-install configuration completed")
	return nil
}

// Helper functions

func (m *Manager) commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func (m *Manager) getDotfilesDir() string {
	// Try to find dotfiles directory
	wd, _ := os.Getwd()

	// Check if we're in the dotfiles directory
	if strings.Contains(wd, "dotfiles") || strings.Contains(wd, "Dotfiles") {
		return wd
	}

	// Default to current directory
	return "."
}

func (m *Manager) getStowPackages(dotfilesDir string) ([]string, error) {
	entries, err := os.ReadDir(dotfilesDir)
	if err != nil {
		return nil, err
	}

	var packages []string
	exclusions := map[string]bool{}

	for _, exclusion := range m.config.StowExclusions {
		exclusions[exclusion] = true
	}

	for _, entry := range entries {
		if entry.IsDir() && !exclusions[entry.Name()] && !strings.HasPrefix(entry.Name(), ".") {
			packages = append(packages, entry.Name())
		}
	}

	return packages, nil
}

func (m *Manager) runDefaults(args ...string) {
	cmd := exec.Command("defaults", args...)
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Failed to run defaults %v: %v\n", args, err)
	}
}

func (m *Manager) killApp(appName string) {
	cmd := exec.Command("killall", appName)
	cmd.Run() // Ignore errors
}