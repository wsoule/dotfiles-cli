package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the dotfiles configuration
type Config struct {
	Brews []string `json:"brews"`
	Casks []string `json:"casks"`
	Taps  []string `json:"taps"`
	Stow  []string `json:"stow"`
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
