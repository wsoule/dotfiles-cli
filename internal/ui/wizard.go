package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dotfiles/internal/config"
)

// Manager handles UI interactions
type Manager struct {
	scanner *bufio.Scanner
}

// NewManager creates a new UI manager
func NewManager() *Manager {
	return &Manager{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// RunSetupWizard runs the interactive setup wizard
func (m *Manager) RunSetupWizard(quick bool) (*config.Config, error) {
	cfg := &config.Config{}

	// Initialize with defaults
	configManager := config.NewManager("")
	defaultCfg := configManager.Get()
	if defaultCfg == nil {
		// Create default config
		err := configManager.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		defaultCfg = configManager.Get()
	}
	*cfg = *defaultCfg

	fmt.Println("Press Enter to begin the setup process...")
	m.scanner.Scan()

	// Personal Information
	if err := m.setupPersonalInfo(cfg); err != nil {
		return nil, err
	}

	if !quick {
		// System Settings
		if err := m.setupSystemSettings(cfg); err != nil {
			return nil, err
		}

		// Development Environment
		if err := m.setupDevelopmentEnvironment(cfg); err != nil {
			return nil, err
		}

		// Additional Packages
		if err := m.setupPackages(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// setupPersonalInfo collects personal information
func (m *Manager) setupPersonalInfo(cfg *config.Config) error {
	fmt.Println()
	fmt.Println("=== 👤 Personal Information ===")

	cfg.Personal.Name = m.getInput("Enter your full name (for git configuration):", cfg.Personal.Name)
	cfg.Personal.Email = m.getInput("Enter your email address (for git configuration):", cfg.Personal.Email)

	fmt.Println("Preferred code editor:")
	fmt.Println("  1) nvim")
	fmt.Println("  2) vim")
	fmt.Println("  3) code (VS Code)")
	fmt.Println("  4) subl (Sublime Text)")
	fmt.Println("  5) nano")

	choice := m.getChoice("Enter choice (1-5):", 1, 5)
	editors := []string{"nvim", "vim", "code", "subl", "nano"}
	cfg.Personal.Editor = editors[choice-1]

	return nil
}

// setupSystemSettings configures system preferences
func (m *Manager) setupSystemSettings(cfg *config.Config) error {
	fmt.Println()
	fmt.Println("=== 🎨 System Settings ===")

	cfg.System.Appearance.DarkMode = m.getYesNo("Enable dark mode?", cfg.System.Appearance.DarkMode)
	cfg.System.Appearance.Enable24HourTime = m.getYesNo("Use 24-hour time format?", cfg.System.Appearance.Enable24HourTime)

	fmt.Println()
	fmt.Println("=== 🖥️ Dock Settings ===")
	cfg.System.Dock.Autohide = m.getYesNo("Auto-hide dock?", cfg.System.Dock.Autohide)

	fmt.Println("Choose dock position:")
	fmt.Println("  1) bottom")
	fmt.Println("  2) left")
	fmt.Println("  3) right")

	choice := m.getChoice("Enter choice (1-3):", 1, 3)
	positions := []string{"bottom", "left", "right"}
	cfg.System.Dock.Position = positions[choice-1]

	tileSizeStr := m.getInput("Dock tile size (16-128):", strconv.Itoa(cfg.System.Dock.TileSize))
	if size, err := strconv.Atoi(tileSizeStr); err == nil && size >= 16 && size <= 128 {
		cfg.System.Dock.TileSize = size
	}

	return nil
}

// setupDevelopmentEnvironment configures development tools
func (m *Manager) setupDevelopmentEnvironment(cfg *config.Config) error {
	fmt.Println()
	fmt.Println("=== 🌐 Programming Languages ===")
	fmt.Println("Select the programming languages you use:")

	languages := []struct {
		key   string
		label string
	}{
		{"javascript", "JavaScript"},
		{"typescript", "TypeScript"},
		{"python", "Python"},
		{"go", "Go"},
		{"rust", "Rust"},
		{"java", "Java"},
		{"php", "PHP"},
		{"ruby", "Ruby"},
		{"csharp", "C#"},
		{"cpp", "C++"},
		{"swift", "Swift"},
		{"kotlin", "Kotlin"},
	}

	for _, lang := range languages {
		current := cfg.Development.Languages[lang.key]
		cfg.Development.Languages[lang.key] = m.getYesNo(fmt.Sprintf("Use %s?", lang.label), current)
	}

	fmt.Println()
	fmt.Println("=== 🐚 Shell Settings ===")
	fmt.Println("Choose shell prompt theme:")
	fmt.Println("  1) powerlevel10k")
	fmt.Println("  2) starship")
	fmt.Println("  3) oh-my-zsh")
	fmt.Println("  4) minimal")

	choice := m.getChoice("Enter choice (1-4):", 1, 4)
	themes := []string{"powerlevel10k", "starship", "oh-my-zsh", "minimal"}
	cfg.Development.Shell.Theme = themes[choice-1]

	return nil
}

// setupPackages configures additional packages
func (m *Manager) setupPackages(cfg *config.Config) error {
	fmt.Println()
	fmt.Println("=== 📦 Additional Packages ===")

	extraBrews := m.getInput("Additional brew packages (space-separated):", strings.Join(cfg.Packages.ExtraBrews, " "))
	if extraBrews != "" {
		cfg.Packages.ExtraBrews = strings.Fields(extraBrews)
	}

	extraCasks := m.getInput("Additional applications (space-separated):", strings.Join(cfg.Packages.ExtraCasks, " "))
	if extraCasks != "" {
		cfg.Packages.ExtraCasks = strings.Fields(extraCasks)
	}

	return nil
}

// getInput prompts for text input with optional default
func (m *Manager) getInput(prompt, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	m.scanner.Scan()
	input := strings.TrimSpace(m.scanner.Text())

	if input == "" {
		return defaultValue
	}
	return input
}

// getYesNo prompts for yes/no input
func (m *Manager) getYesNo(prompt string, defaultValue bool) bool {
	if defaultValue {
		fmt.Printf("%s [Y/n]: ", prompt)
	} else {
		fmt.Printf("%s [y/N]: ", prompt)
	}

	m.scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(m.scanner.Text()))

	switch input {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return defaultValue
	}
}

// getChoice prompts for a numbered choice
func (m *Manager) getChoice(prompt string, min, max int) int {
	for {
		fmt.Printf("%s ", prompt)
		m.scanner.Scan()
		input := strings.TrimSpace(m.scanner.Text())

		if choice, err := strconv.Atoi(input); err == nil && choice >= min && choice <= max {
			return choice
		}

		fmt.Printf("Please enter a number between %d and %d\n", min, max)
	}
}