package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Hooks represents pre/post commands for various operations
type Hooks struct {
	PreInstall  []string `json:"pre_install,omitempty"`
	PostInstall []string `json:"post_install,omitempty"`
	PreSync     []string `json:"pre_sync,omitempty"`
	PostSync    []string `json:"post_sync,omitempty"`
	PreStow     []string `json:"pre_stow,omitempty"`
	PostStow    []string `json:"post_stow,omitempty"`
}

// PackageConfig represents configuration for a specific package
type PackageConfig struct {
	PostInstall []string `json:"post_install,omitempty"`
	PreInstall  []string `json:"pre_install,omitempty"`
}

// Config represents the dotfiles configuration
type Config struct {
	Brews          []string                 `json:"brews"`
	Casks          []string                 `json:"casks"`
	Taps           []string                 `json:"taps"`
	Stow           []string                 `json:"stow"`
	Hooks          *Hooks                   `json:"hooks,omitempty"`
	PackageConfigs map[string]PackageConfig `json:"package_configs,omitempty"`
	Groups         map[string][]string      `json:"groups,omitempty"`         // Package groups/tags
	PackageTags    map[string][]string      `json:"package_tags,omitempty"`   // Tags per package
}

// Load reads configuration from JSON file
func Load(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Brews: []string{},
			Casks: []string{},
			Taps:  []string{},
			Stow:  []string{},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes configuration to JSON file
func (c *Config) Save(configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GenerateBrewfile creates a Brewfile from the configuration
func (c *Config) GenerateBrewfile() string {
	var content string

	// Add taps
	for _, tap := range c.Taps {
		content += "tap \"" + tap + "\"\n"
	}
	if len(c.Taps) > 0 {
		content += "\n"
	}

	// Add brews
	for _, brew := range c.Brews {
		content += "brew \"" + brew + "\"\n"
	}
	if len(c.Brews) > 0 {
		content += "\n"
	}

	// Add casks
	for _, cask := range c.Casks {
		content += "cask \"" + cask + "\"\n"
	}

	return content
}

// GetAllPackages returns all packages (brews + casks) as a single list
// Useful for Linux package managers that don't distinguish between them
func (c *Config) GetAllPackages() []string {
	all := make([]string, 0, len(c.Brews)+len(c.Casks))
	all = append(all, c.Brews...)
	all = append(all, c.Casks...)
	return all
}
