package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
)

// createPrivateDirectoryStructure creates the private directory with all necessary env files
func createPrivateDirectoryStructure(dotfilesDir string) error {
	privateDir := filepath.Join(dotfilesDir, "private")

	// Create private directory
	if err := os.MkdirAll(privateDir, 0755); err != nil {
		return fmt.Errorf("failed to create private directory: %v", err)
	}

	// Create .ssh directory in private
	sshDir := filepath.Join(privateDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %v", err)
	}

	// Create base shell/env files in private directory
	privateFiles := map[string]string{
		".env.local": `# Local environment variables
# Add your private environment variables here
export DOTFILES_PRIVATE="true"
`,
		".gitconfig.local": `# Personal git configuration
# Add your personal git settings here
[user]
	# name = Your Name
	# email = your.email@example.com
`,
		"env-private": `# Private environment variables
# Add sensitive environment variables here
# This file is git-ignored and kept private

# Example:
# export API_KEY="your-secret-key"
# export DATABASE_URL="your-db-connection"
`,
	}

	for filename, content := range privateFiles {
		filePath := filepath.Join(privateDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create %s: %v", filename, err)
			}
		}
	}

	// Update .gitignore to exclude private directory
	gitignoreFile := filepath.Join(dotfilesDir, ".gitignore")
	gitignoreContent := "\n# Private files\nprivate/\n"

	f, err := os.OpenFile(gitignoreFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.WriteString(gitignoreContent)
		f.Close()
	}

	return nil
}

// createShellStowPackages creates shell and git stow packages with all configuration files
func createShellStowPackages(stowDir string) error {
	// Create shell stow package
	shellStowDir := filepath.Join(stowDir, "shell")
	if err := os.MkdirAll(shellStowDir, 0755); err != nil {
		return fmt.Errorf("failed to create shell stow directory: %v", err)
	}

	// Create public shell files directly in stow package
	publicShellFiles := map[string]string{
		".zshrc": `# Base zsh configuration
# Common zsh settings

# Load environment files
[[ -f ~/.env ]] && source ~/.env
[[ -f ~/.common.env ]] && source ~/.common.env
[[ -f ~/.aliases.env ]] && source ~/.aliases.env
[[ -f ~/.dotfiles/private/env-private ]] && source ~/.dotfiles/private/env-private

# Basic zsh options
setopt AUTO_CD
setopt HIST_VERIFY
setopt SHARE_HISTORY
setopt APPEND_HISTORY

# History settings
HISTFILE=~/.zsh_history
HISTSIZE=10000
SAVEHIST=10000
`,
		".env": `# Base environment configuration
# Common environment variables for all shells
export EDITOR="vim"
export BROWSER="open"
export DOTFILES_DIR="$HOME/.dotfiles"

# Path additions
export PATH="$HOME/.local/bin:$PATH"
`,
		".common.env": `# Common environment variables
# Shared across all environments

# Development tools
export LANG="en_US.UTF-8"
export LC_ALL="en_US.UTF-8"

# Better defaults
export LESS="-R"
export GREP_OPTIONS="--color=auto"
`,
		".aliases.env": `# Common aliases
# Add your aliases here

# File operations
alias ..="cd .."
alias ...="cd ../.."
alias ....="cd ../../.."

# Directory listing
alias ll="ls -la"
alias la="ls -A"
alias l="ls -CF"
alias ls="ls --color=auto"

# Git shortcuts
alias gs="git status"
alias ga="git add"
alias gc="git commit"
alias gp="git push"
alias gl="git pull"
alias gd="git diff"
alias gb="git branch"
alias gco="git checkout"

# Utilities
alias grep="grep --color=auto"
alias h="history"
alias c="clear"
alias reload="source ~/.zshrc"
`,
	}

	// Create public files directly in stow package
	for filename, content := range publicShellFiles {
		filePath := filepath.Join(shellStowDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create %s: %v", filename, err)
			}
		}
	}

	// Create git stow package for git config
	gitStowDir := filepath.Join(stowDir, "git")
	if err := os.MkdirAll(gitStowDir, 0755); err != nil {
		return fmt.Errorf("failed to create git stow directory: %v", err)
	}

	// Create a base .gitconfig that includes the private config
	gitConfigContent := `# Main git configuration
[include]
	path = ~/.dotfiles/private/.gitconfig.local

[core]
	editor = vim
	autocrlf = input

[push]
	default = simple

[pull]
	rebase = false
`

	gitConfigPath := filepath.Join(gitStowDir, ".gitconfig")
	if err := os.WriteFile(gitConfigPath, []byte(gitConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitconfig: %v", err)
	}

	return nil
}

// addStowPackagesToConfig adds shell and git packages to the stow configuration
func addStowPackagesToConfig(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Add shell and git to stow packages if not already present
	if !contains(cfg.Stow, "shell") {
		cfg.Stow = append(cfg.Stow, "shell")
	}
	if !contains(cfg.Stow, "git") {
		cfg.Stow = append(cfg.Stow, "git")
	}

	// Save updated config
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	return nil
}

// stowPackages stows the given packages using GNU stow
func stowPackages(packages []string, stowDir, target string) error {
	for _, pkg := range packages {
		stowArgs := []string{"-d", stowDir, "-t", target, pkg}
		stowCmd := exec.Command("stow", stowArgs...)

		output, err := stowCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error stowing %s: %v (output: %s)", pkg, err, strings.TrimSpace(string(output)))
		}
	}
	return nil
}

// setupCompleteEnvironment sets up the complete dotfiles environment (private dir + shell packages + config)
func setupCompleteEnvironment(dotfilesDir string, shouldStow bool) error {
	stowDir := filepath.Join(dotfilesDir, "stow")
	configPath := filepath.Join(dotfilesDir, "config.json")

	// Create stow directory if it doesn't exist
	if err := os.MkdirAll(stowDir, 0755); err != nil {
		return fmt.Errorf("failed to create stow directory: %v", err)
	}

	// Create private directory structure
	if err := createPrivateDirectoryStructure(dotfilesDir); err != nil {
		return fmt.Errorf("failed to create private directory: %v", err)
	}

	// Create shell stow packages
	if err := createShellStowPackages(stowDir); err != nil {
		return fmt.Errorf("failed to create shell packages: %v", err)
	}

	// Add packages to config
	if err := addStowPackagesToConfig(configPath); err != nil {
		return fmt.Errorf("failed to update config: %v", err)
	}

	// Optionally stow the packages
	if shouldStow {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %v", err)
		}

		if err := stowPackages([]string{"shell", "git"}, stowDir, home); err != nil {
			return fmt.Errorf("failed to stow packages: %v", err)
		}
	}

	return nil
}