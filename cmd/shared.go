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
		".env-private": `# Private environment variables
# Add sensitive environment variables here
# This file is git-ignored and kept private

# Example:
# export API_KEY="your-secret-key"
# export DATABASE_URL="your-db-connection"
`,
		".gitconfig.local": `# Personal git configuration
# Add your personal git settings here
[user]
	# name = Your Name
	# email = your.email@example.com
`,
		"env-private": `# Private environment variables (deprecated - use .env-private instead)
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

// createSSHStowPackage creates an SSH stow package with a symlink to private/.ssh
func createSSHStowPackage(dotfilesDir string) error {
	stowDir := filepath.Join(dotfilesDir, "stow")
	sshStowDir := filepath.Join(stowDir, "ssh")

	// Create ssh stow package directory
	if err := os.MkdirAll(sshStowDir, 0755); err != nil {
		return fmt.Errorf("failed to create ssh stow directory: %v", err)
	}

	// Create relative symlink from stow package to private .ssh directory
	stowSshLink := filepath.Join(sshStowDir, ".ssh")
	relativePrivatePath := filepath.Join("..", "..", "private", ".ssh")

	// Remove existing symlink if it exists
	if _, err := os.Lstat(stowSshLink); err == nil {
		if err := os.Remove(stowSshLink); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}

	// Create the symlink
	if err := os.Symlink(relativePrivatePath, stowSshLink); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}

	return nil
}

// createShellStowPackages creates shell and git stow packages with all configuration files
func createShellStowPackages(stowDir string, shell string) error {
	// Create shell stow package
	shellStowDir := filepath.Join(stowDir, "shell")
	if err := os.MkdirAll(shellStowDir, 0755); err != nil {
		return fmt.Errorf("failed to create shell stow directory: %v", err)
	}

	// Create public shell files directly in stow package
	publicShellFiles := make(map[string]string)

	// Generate shell-specific config
	if shell == "fish" {
		// For fish, we need to create .config/fish subdirectory in the stow package
		fishConfigDir := filepath.Join(shellStowDir, ".config", "fish")
		if err := os.MkdirAll(fishConfigDir, 0755); err != nil {
			return fmt.Errorf("failed to create fish config directory: %v", err)
		}

		fishConfig := `# Base fish configuration
# Common fish settings

# Load environment files
if test -f ~/.env
    source ~/.env
end
if test -f ~/.common.env
    source ~/.common.env
end
if test -f ~/.aliases.env
    source ~/.aliases.env
end
if test -f ~/.env.local
    source ~/.env.local
end
if test -f ~/.env-private
    source ~/.env-private
end

# History settings
set -x HISTFILE ~/.fish_history
set -x HISTSIZE 10000
set -x SAVEHIST 10000
`
		// Write fish config to the correct location
		fishConfigPath := filepath.Join(fishConfigDir, "config.fish")
		if err := os.WriteFile(fishConfigPath, []byte(fishConfig), 0644); err != nil {
			return fmt.Errorf("failed to create config.fish: %v", err)
		}
	} else {
		// Default to zsh
		publicShellFiles[".zshrc"] = `# Base zsh configuration
# Common zsh settings

# Load environment files
[[ -f ~/.env ]] && source ~/.env
[[ -f ~/.common.env ]] && source ~/.common.env
[[ -f ~/.aliases.env ]] && source ~/.aliases.env
[[ -f ~/.env.local ]] && source ~/.env.local
[[ -f ~/.env-private ]] && source ~/.env-private

# Basic zsh options
setopt AUTO_CD
setopt HIST_VERIFY
setopt SHARE_HISTORY
setopt APPEND_HISTORY

# History settings
HISTFILE=~/.zsh_history
HISTSIZE=10000
SAVEHIST=10000
`
	}

	// Common environment files that work with both shells
	publicShellFiles[".env"] = `# Base environment configuration
# Common environment variables for all shells
export EDITOR="vim"
export BROWSER="open"
export DOTFILES_DIR="$HOME/.dotfiles"

# Path additions
export PATH="$HOME/.local/bin:$PATH"
`
	publicShellFiles[".common.env"] = `# Common environment variables
# Shared across all environments

# Development tools
export LANG="en_US.UTF-8"
export LC_ALL="en_US.UTF-8"

# Better defaults
export LESS="-R"
export GREP_OPTIONS="--color=auto"
`

	// Generate shell-agnostic aliases file
	if shell == "fish" {
		publicShellFiles[".aliases.env"] = `# Common aliases for fish
# Add your aliases here

# File operations
alias .. "cd .."
alias ... "cd ../.."
alias .... "cd ../../.."

# Directory listing
alias ll "ls -la"
alias la "ls -A"
alias l "ls -CF"

# Git shortcuts
alias gs "git status"
alias ga "git add"
alias gc "git commit"
alias gp "git push"
alias gl "git pull"
alias gd "git diff"
alias gb "git branch"
alias gco "git checkout"

# Utilities
alias h "history"
alias c "clear"
alias reload "source ~/.config/fish/config.fish"
`
	} else {
		publicShellFiles[".aliases.env"] = `# Common aliases
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
`
	}

	// Remove old shell-specific configs that might exist
	oldConfigs := []string{".zshrc", "config.fish"}
	for _, oldConfig := range oldConfigs {
		oldPath := filepath.Join(shellStowDir, oldConfig)
		if _, err := os.Stat(oldPath); err == nil {
			// Only remove if it's not the one we're creating
			if _, shouldCreate := publicShellFiles[oldConfig]; !shouldCreate {
				os.Remove(oldPath)
			}
		}
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

	// Create symlinks to private env files
	privateEnvFiles := []string{".env.local", ".env-private"}
	for _, filename := range privateEnvFiles {
		stowLinkPath := filepath.Join(shellStowDir, filename)
		relativePrivatePath := filepath.Join("..", "..", "private", filename)

		// Remove existing symlink if it exists
		if _, err := os.Lstat(stowLinkPath); err == nil {
			if err := os.Remove(stowLinkPath); err != nil {
				return fmt.Errorf("failed to remove existing symlink %s: %v", filename, err)
			}
		}

		// Create the symlink
		if err := os.Symlink(relativePrivatePath, stowLinkPath); err != nil {
			return fmt.Errorf("failed to create symlink for %s: %v", filename, err)
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

// addStowPackagesToConfig adds shell, git, and ssh packages to the stow configuration
func addStowPackagesToConfig(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Add shell, git, and ssh to stow packages if not already present
	if !contains(cfg.Stow, "shell") {
		cfg.Stow = append(cfg.Stow, "shell")
	}
	if !contains(cfg.Stow, "git") {
		cfg.Stow = append(cfg.Stow, "git")
	}
	if !contains(cfg.Stow, "ssh") {
		cfg.Stow = append(cfg.Stow, "ssh")
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
func setupCompleteEnvironment(dotfilesDir string, shouldStow bool, shell string) error {
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

	// Create shell stow packages with selected shell
	if err := createShellStowPackages(stowDir, shell); err != nil {
		return fmt.Errorf("failed to create shell packages: %v", err)
	}

	// Create SSH stow package that links to private/.ssh
	if err := createSSHStowPackage(dotfilesDir); err != nil {
		return fmt.Errorf("failed to create ssh stow package: %v", err)
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