// Package config provides configuration management for the dotfiles manager.
// It handles loading, saving, and manipulating dotfiles configurations stored in JSON format.
// The configuration includes package lists, hooks, roles, and conditional configurations
// for different machines and operating systems.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Hooks represents pre/post commands for various operations.
// Hooks are shell commands that execute at specific points in the dotfiles workflow,
// allowing for custom setup, teardown, and validation logic.
type Hooks struct {
	PreInstall  []string `json:"pre_install,omitempty"`  // Commands to run before package installation
	PostInstall []string `json:"post_install,omitempty"` // Commands to run after package installation
	PreSync     []string `json:"pre_sync,omitempty"`     // Commands to run before syncing dotfiles
	PostSync    []string `json:"post_sync,omitempty"`    // Commands to run after syncing dotfiles
	PreStow     []string `json:"pre_stow,omitempty"`     // Commands to run before stowing packages
	PostStow    []string `json:"post_stow,omitempty"`    // Commands to run after stowing packages
}

// PackageConfig represents configuration for a specific package.
// This allows per-package customization of installation behavior through hooks.
type PackageConfig struct {
	PostInstall []string `json:"post_install,omitempty"` // Package-specific post-install hooks
	PreInstall  []string `json:"pre_install,omitempty"`  // Package-specific pre-install hooks
}

// Role represents a set of packages for a specific purpose.
// Roles group related packages together for different development environments or workflows.
// Examples include "web-dev", "data-science", "devops", etc.
type Role struct {
	Name        string   `json:"name"`        // Role identifier (e.g., "web-dev")
	Description string   `json:"description"` // Human-readable description
	Brews       []string `json:"brews"`       // CLI tools for this role
	Casks       []string `json:"casks"`       // GUI applications for this role
	Taps        []string `json:"taps"`        // Homebrew taps required for this role
	Stow        []string `json:"stow"`        // Dotfile packages for this role
}

// ConditionalConfig represents OS/hostname-specific configurations.
// This enables different package installations based on the machine you're on,
// perfect for maintaining work vs. personal machine configurations in one repo.
type ConditionalConfig struct {
	Hostname  string            `json:"hostname,omitempty"`  // Hostname pattern (supports globs like "work-*")
	OS        string            `json:"os,omitempty"`        // Operating system (darwin, linux, etc.)
	Variables map[string]string `json:"variables,omitempty"` // Template variables to set when condition matches
	Brews     []string          `json:"brews,omitempty"`     // CLI tools to install when condition matches
	Casks     []string          `json:"casks,omitempty"`     // GUI apps to install when condition matches
}

// Config represents the complete dotfiles configuration.
// This is the main configuration structure that stores all package lists,
// hooks, roles, and settings for managing your development environment.
// It is typically stored as config.json in the ~/.dotfiles directory.
type Config struct {
	Brews               []string                 `json:"brews"`                          // CLI tools (e.g., git, curl, tmux)
	Casks               []string                 `json:"casks"`                          // GUI applications (e.g., visual-studio-code, docker)
	Taps                []string                 `json:"taps"`                           // Homebrew taps (additional package repositories)
	Stow                []string                 `json:"stow"`                           // Dotfile packages managed by GNU Stow
	Hooks               *Hooks                   `json:"hooks,omitempty"`                // Global pre/post operation hooks
	PackageConfigs      map[string]PackageConfig `json:"package_configs,omitempty"`      // Per-package configuration and hooks
	Groups              map[string][]string      `json:"groups,omitempty"`               // Package groups for organization
	PackageTags         map[string][]string      `json:"package_tags,omitempty"`         // Tags assigned to packages
	Roles               map[string]Role          `json:"roles,omitempty"`                // Named package collections (web-dev, data-science, etc.)
	Conditionals        []ConditionalConfig      `json:"conditionals,omitempty"`         // OS/hostname-specific package configurations
	Variables           map[string]string        `json:"variables,omitempty"`            // Template variables for config file substitution
	ActiveProfile       string                   `json:"active_profile,omitempty"`       // Currently active machine profile
	PackageDependencies map[string][]string      `json:"package_dependencies,omitempty"` // Package dependency graph for install ordering
}

// Load reads and parses the configuration from a JSON file at the specified path.
// If the file doesn't exist, it returns an empty configuration with initialized slices.
// This allows graceful handling of first-time setup scenarios.
//
// Parameters:
//   - configPath: Absolute path to the config.json file (typically ~/.dotfiles/config.json)
//
// Returns:
//   - *Config: The loaded configuration or an empty config if file doesn't exist
//   - error: Any error encountered during file reading or JSON parsing
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

// Save writes the configuration to a JSON file with pretty-printing (indented).
// It automatically creates parent directories if they don't exist.
//
// Parameters:
//   - configPath: Absolute path where the config.json should be written
//
// Returns:
//   - error: Any error encountered during directory creation, marshaling, or file writing
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

// GenerateBrewfile creates a Brewfile from the configuration.
// A Brewfile is Homebrew's declarative package format that can be used with `brew bundle`.
// The generated file includes taps, brews (CLI tools), and casks (GUI applications).
//
// Returns:
//   - string: The generated Brewfile content in the format:
//     tap "repository"
//     brew "package"
//     cask "application"
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

// GetAllPackages returns all packages (brews + casks) as a single combined list.
// This is useful for Linux package managers (apt, pacman, yum) that don't distinguish
// between CLI tools and GUI applications the way Homebrew does on macOS.
//
// Returns:
//   - []string: Combined list of all package names from both Brews and Casks slices
func (c *Config) GetAllPackages() []string {
	all := make([]string, 0, len(c.Brews)+len(c.Casks))
	all = append(all, c.Brews...)
	all = append(all, c.Casks...)
	return all
}
