package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [repo-url]",
	Short: "Set up dotfiles repository and directory structure",
	Long:  `Fork and clone a dotfiles repository to ~/.dotfiles/, create private directory structure, and set up stow packages.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			return
		}

		dotfilesDir := filepath.Join(homeDir, ".dotfiles")

		// Check if directory already exists
		if _, err := os.Stat(dotfilesDir); !os.IsNotExist(err) {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("❌ Directory %s already exists. Use --force to overwrite.\n", dotfilesDir)
				return
			}
			// Remove existing directory
			if err := os.RemoveAll(dotfilesDir); err != nil {
				fmt.Printf("❌ Error removing existing directory: %v\n", err)
				return
			}
		}

		fmt.Printf("🚀 Setting up dotfiles from %s...\n", repoURL)

		// Clone the repository
		fmt.Printf("📥 Cloning repository to %s...\n", dotfilesDir)
		cloneCmd := exec.Command("git", "clone", repoURL, dotfilesDir)
		if err := cloneCmd.Run(); err != nil {
			fmt.Printf("❌ Error cloning repository: %v\n", err)
			return
		}

		// Create stow directory
		stowDir := filepath.Join(dotfilesDir, "stow")
		if err := os.MkdirAll(stowDir, 0755); err != nil {
			fmt.Printf("❌ Error creating stow directory: %v\n", err)
			return
		}
		fmt.Printf("📁 Created stow directory at %s\n", stowDir)

		// Create .config stow package
		configStowDir := filepath.Join(stowDir, "config")
		configDir := filepath.Join(configStowDir, ".config")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("❌ Error creating .config stow package: %v\n", err)
			return
		}
		fmt.Printf("📁 Created .config stow package at %s\n", configStowDir)

		// Create private directory structure
		privateDir := filepath.Join(dotfilesDir, "private")
		if err := os.MkdirAll(privateDir, 0755); err != nil {
			fmt.Printf("❌ Error creating private directory: %v\n", err)
			return
		}

		// Create .ssh directory in private
		sshDir := filepath.Join(privateDir, ".ssh")
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			fmt.Printf("❌ Error creating .ssh directory: %v\n", err)
			return
		}

		// Create base shell/env files in private directory
		if err := createBaseShellFiles(privateDir); err != nil {
			fmt.Printf("❌ Error creating base shell files: %v\n", err)
		} else {
			fmt.Printf("📄 Created base shell/env files in private directory\n")
		}

		fmt.Printf("🔒 Created private directory structure at %s\n", privateDir)

		// Set up shell stow packages
		if err := createShellStowPackages(stowDir, privateDir); err != nil {
			fmt.Printf("❌ Error creating shell stow packages: %v\n", err)
		} else {
			fmt.Printf("🔗 Created shell stow packages\n")
		}

		// Add private directory to gitignore if it doesn't exist
		gitignoreFile := filepath.Join(dotfilesDir, ".gitignore")
		gitignoreContent := "\n# Private files\nprivate/\n"

		f, err := os.OpenFile(gitignoreFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			f.WriteString(gitignoreContent)
			f.Close()
			fmt.Printf("📝 Updated .gitignore to exclude private directory\n")
		}

		fmt.Printf("✅ Dotfiles setup complete!\n")
		fmt.Printf("📂 Repository cloned to: %s\n", dotfilesDir)
		fmt.Printf("📁 Stow packages directory: %s\n", stowDir)
		fmt.Printf("🔒 Private files directory: %s\n", privateDir)
		fmt.Printf("\n💡 Next steps:\n")
		fmt.Printf("   1. Configure your private files in %s\n", privateDir)
		fmt.Printf("   2. Run 'dotfiles stow shell' to activate shell configuration\n")
		fmt.Printf("   3. Run 'dotfiles stow git' to activate git configuration\n")
		fmt.Printf("   4. Add more dotfiles to stow packages in %s\n", stowDir)
		fmt.Printf("   5. Run 'dotfiles init' to initialize package configuration\n")
	},
}

func createBaseShellFiles(privateDir string) error {
	// Only create truly private files in private directory
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

	return nil
}

func createShellStowPackages(stowDir, privateDir string) error {
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

func init() {
	setupCmd.Flags().BoolP("force", "f", false, "Force overwrite existing directory")
	rootCmd.AddCommand(setupCmd)
}