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

		// Create private files
		privateFiles := []string{
			filepath.Join(privateDir, ".env.local"),
			filepath.Join(privateDir, ".gitconfig.local"),
		}

		for _, file := range privateFiles {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				f, err := os.Create(file)
				if err != nil {
					fmt.Printf("❌ Error creating %s: %v\n", file, err)
					continue
				}
				f.Close()
				fmt.Printf("📄 Created %s\n", file)
			}
		}

		fmt.Printf("🔒 Created private directory structure at %s\n", privateDir)

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
		fmt.Printf("   1. Add your dotfiles to stow packages in %s\n", stowDir)
		fmt.Printf("   2. Configure your private files in %s\n", privateDir)
		fmt.Printf("   3. Run 'dotfiles stow [package]' to symlink packages\n")
		fmt.Printf("   4. Run 'dotfiles init' to initialize configuration\n")
	},
}

func init() {
	setupCmd.Flags().BoolP("force", "f", false, "Force overwrite existing directory")
	rootCmd.AddCommand(setupCmd)
}