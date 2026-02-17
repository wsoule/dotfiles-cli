package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use: "setup [repo-url]",
	GroupID: "getting-started",
	Short: "Set up dotfiles repository and directory structure",
	Long:    `Fork and clone a dotfiles repository to ~/.dotfiles/, create private directory structure, and set up stow packages.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf(" Error getting home directory: %v\n", err)
			return
		}

		dotfilesDir := filepath.Join(homeDir, ".dotfiles")

		// Check if directory already exists
		if _, err := os.Stat(dotfilesDir); !os.IsNotExist(err) {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf(" Directory %s already exists. Use --force to overwrite.\n", dotfilesDir)
				return
			}
			// Remove existing directory
			if err := os.RemoveAll(dotfilesDir); err != nil {
				fmt.Printf(" Error removing existing directory: %v\n", err)
				return
			}
		}

		fmt.Printf(" Setting up dotfiles from %s...\n", repoURL)

		// Clone the repository
		fmt.Printf(" Cloning repository to %s...\n", dotfilesDir)
		cloneCmd := exec.Command("git", "clone", repoURL, dotfilesDir)
		if err := cloneCmd.Run(); err != nil {
			fmt.Printf(" Error cloning repository: %v\n", err)
			return
		}

		// Create stow directory
		stowDir := filepath.Join(dotfilesDir, "stow")
		if err := os.MkdirAll(stowDir, 0755); err != nil {
			fmt.Printf(" Error creating stow directory: %v\n", err)
			return
		}
		fmt.Printf(" Created stow directory at %s\n", stowDir)

		// Create .config stow package
		configStowDir := filepath.Join(stowDir, "config")
		configDir := filepath.Join(configStowDir, ".config")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf(" Error creating .config stow package: %v\n", err)
			return
		}
		fmt.Printf(" Created .config stow package at %s\n", configStowDir)

		// Set up complete environment (private dir + shell packages, no stowing)
		// Default to zsh for setup command (user can manually change later)
		if err := setupCompleteEnvironment(dotfilesDir, false, "zsh"); err != nil {
			fmt.Printf(" Error setting up environment: %v\n", err)
			return
		}
		fmt.Printf(" Created private directory structure\n")
		fmt.Printf(" Created shell stow packages\n")
		fmt.Printf(" Updated .gitignore to exclude private directory\n")

		fmt.Printf(" Dotfiles setup complete!\n")
		fmt.Printf(" Repository cloned to: %s\n", dotfilesDir)
		fmt.Printf(" Stow packages directory: %s\n", stowDir)
		fmt.Printf(" Private files directory: %s\n", filepath.Join(dotfilesDir, "private"))
		fmt.Printf("\n Next steps:\n")
		fmt.Printf("1. Configure your private files in %s\n", filepath.Join(dotfilesDir, "private"))
		fmt.Printf("2. Run 'dotfiles stow shell' to activate shell configuration\n")
		fmt.Printf("3. Run 'dotfiles stow git' to activate git configuration\n")
		fmt.Printf("4. Add more dotfiles to stow packages in %s\n", stowDir)
		fmt.Printf("5. Run 'dotfiles init' to initialize package configuration\n")
	},
}

func init() {
	setupCmd.Flags().BoolP("force", "f", false, "Force overwrite existing directory")
	rootCmd.AddCommand(setupCmd)
}
