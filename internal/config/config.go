package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Manager handles configuration operations
type Manager struct {
	config *Config
	path   string
}

// NewManager creates a new configuration manager
func NewManager(configPath string) *Manager {
	return &Manager{
		path: configPath,
	}
}

// Load reads configuration from file
func (m *Manager) Load() error {
	if m.path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		m.path = filepath.Join(home, ".dotfiles", "config.json")
	}

	// Check if config file exists
	if _, err := os.Stat(m.path); os.IsNotExist(err) {
		// Create default config
		m.config = m.createDefaultConfig()
		return m.Save()
	}

	data, err := os.ReadFile(m.path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	m.config = &Config{}
	if err := json.Unmarshal(data, m.config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// Save writes configuration to file
func (m *Manager) Save() error {
	if m.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	// Update metadata
	now := time.Now()
	m.config.Metadata.LastModified = now
	if m.config.Metadata.CreatedAt.IsZero() {
		m.config.Metadata.CreatedAt = now
	}

	// Ensure directory exists
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Get returns the current configuration
func (m *Manager) Get() *Config {
	return m.config
}

// Set updates the configuration
func (m *Manager) Set(config *Config) {
	m.config = config
}

// Validate checks if configuration is valid
func (m *Manager) Validate() error {
	if m.config == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Validate required fields
	if m.config.Personal.Editor == "" {
		m.config.Personal.Editor = "nvim"
	}

	if m.config.Development.Git.DefaultBranch == "" {
		m.config.Development.Git.DefaultBranch = "main"
	}

	if m.config.Development.Git.PushDefault == "" {
		m.config.Development.Git.PushDefault = "simple"
	}

	// Validate dock position
	validPositions := map[string]bool{"bottom": true, "left": true, "right": true}
	if !validPositions[m.config.System.Dock.Position] {
		m.config.System.Dock.Position = "bottom"
	}

	// Validate finder view
	validViews := map[string]bool{"icon": true, "list": true, "column": true, "gallery": true}
	if !validViews[m.config.System.Finder.DefaultView] {
		m.config.System.Finder.DefaultView = "column"
	}

	return nil
}

// createDefaultConfig returns a default configuration
func (m *Manager) createDefaultConfig() *Config {
	now := time.Now()

	return &Config{
		Personal: Personal{
			Name:   "",
			Email:  "",
			Editor: "nvim",
		},
		Installation: Installation{
			Homebrew:      true,
			Brewfile:      true,
			Dotfiles:      true,
			MacOSDefaults: true,
			NPMPackages:   true,
		},
		System: System{
			Appearance: Appearance{
				DarkMode:         true,
				Enable24HourTime: true,
			},
			Dock: Dock{
				Autohide:      true,
				Position:      "bottom",
				TileSize:      50,
				Magnification: false,
			},
			Finder: Finder{
				ShowHiddenFiles:    false,
				DefaultView:        "column",
				ShowFileExtensions: true,
			},
			Keyboard: Keyboard{
				KeyRepeatRate:       2,
				DisablePressAndHold: true,
			},
			Security: Security{
				RequirePasswordImmediately: true,
				DisplaySleepMinutes:        10,
				ComputerSleepMinutes:       30,
			},
			Screenshots: Screenshots{
				Location: "$HOME/Documents/Screenshots",
				Format:   "png",
			},
			MenuBar: MenuBar{
				ShowBatteryPercent: true,
				ShowDate:           true,
			},
			Safari: Safari{
				ShowDevelopMenu: false,
				DefaultEncoding: "UTF-8",
			},
		},
		Development: Development{
			Git: Git{
				DefaultBranch: "main",
				PullRebase:    true,
				PushDefault:   "simple",
			},
			Languages: map[string]bool{
				"javascript": false,
				"typescript": false,
				"python":     false,
				"go":         false,
				"rust":       false,
				"java":       false,
				"php":        false,
				"ruby":       false,
				"csharp":     false,
				"cpp":        false,
				"swift":      false,
				"kotlin":     false,
			},
			Frameworks: map[string]bool{
				"react":   false,
				"vue":     false,
				"angular": false,
				"svelte":  false,
				"nextjs":  false,
				"nuxt":    false,
				"django":  false,
				"flask":   false,
				"fastapi": false,
				"rails":   false,
				"laravel": false,
				"spring":  false,
				"dotnet":  false,
				"flutter": false,
			},
			Tools: map[string]bool{
				"docker":     false,
				"kubernetes": false,
				"terraform":  false,
				"ansible":    false,
			},
			Shell: Shell{
				Theme:         "powerlevel10k",
				TerminalTheme: "dark",
				Plugins: map[string]bool{
					"autosuggestions":     true,
					"syntax_highlighting": true,
					"fzf":                 true,
					"zoxide":              true,
				},
			},
			Aliases: map[string]bool{
				"node":   true,
				"git":    true,
				"docker": true,
				"system": true,
			},
		},
		Packages: Packages{
			ExtraBrews: []string{},
			ExtraCasks: []string{},
			ExtraTaps:  []string{},
			NPMGlobals: []string{"nx"},
		},
		Directories: []string{
			"$HOME/Projects",
			"$HOME/Projects/personal",
			"$HOME/Projects/work",
			"$HOME/Projects/playground",
			"$HOME/Scripts",
			"$HOME/Documents/Screenshots",
		},
		StowExclusions: []string{
			"scripts",
			"private",
			".git",
			"node_modules",
		},
		Metadata: Metadata{
			Version:      "1.0.0",
			CreatedAt:    now,
			LastModified: now,
			CreatedBy:    "setup-wizard",
		},
	}
}

// LoadFromViper loads configuration from viper (for CLI integration)
func (m *Manager) LoadFromViper() error {
	m.config = &Config{}
	return viper.Unmarshal(m.config)
}