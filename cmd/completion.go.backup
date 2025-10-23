package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:     "completion <shell>",
	GroupID: "advanced",
	Short:   "🔧 Generate shell completion scripts",
	Long: `🔧 Shell Completion Scripts

Generate completion scripts for your shell to enable tab completion.

Supported shells:
  • bash
  • zsh
  • fish
  • powershell

Examples:
  # Bash
  dotfiles completion bash > /usr/local/etc/bash_completion.d/dotfiles

  # Zsh
  dotfiles completion zsh > "${fpath[1]}/_dotfiles"

  # Fish
  dotfiles completion fish > ~/.config/fish/completions/dotfiles.fish

  # PowerShell
  dotfiles completion powershell > dotfiles.ps1`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]

		var err error
		switch shell {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			fmt.Printf("❌ Unsupported shell: %s\n", shell)
			fmt.Println("\nSupported shells: bash, zsh, fish, powershell")
			os.Exit(1)
		}

		if err != nil {
			fmt.Printf("❌ Error generating completion: %v\n", err)
			os.Exit(1)
		}
	},
}

var completionInstallCmd = &cobra.Command{
	Use:   "install <shell>",
	Short: "Install completion script for your shell",
	Long: `Automatically install completion script to the correct location.

Supported shells:
  • bash
  • zsh
  • fish`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		home, _ := os.UserHomeDir()

		var completionPath string
		var instructions string

		switch shell {
		case "bash":
			// Try common bash completion directories
			possiblePaths := []string{
				"/usr/local/etc/bash_completion.d/dotfiles",
				"/etc/bash_completion.d/dotfiles",
				home + "/.local/share/bash-completion/completions/dotfiles",
			}

			for _, path := range possiblePaths {
				dir := filepath.Dir(path)
				if _, err := os.Stat(dir); err == nil {
					completionPath = path
					break
				}
			}

			if completionPath == "" {
				// Fallback to user local
				completionPath = home + "/.local/share/bash-completion/completions/dotfiles"
				os.MkdirAll(filepath.Dir(completionPath), 0755)
			}

			instructions = "Add this to your ~/.bashrc:\nsource ~/.local/share/bash-completion/completions/dotfiles"

		case "zsh":
			// Get first fpath directory
			completionPath = home + "/.zsh/completions/_dotfiles"
			os.MkdirAll(filepath.Dir(completionPath), 0755)
			instructions = "Add this to your ~/.zshrc:\nfpath=(~/.zsh/completions $fpath)\nautoload -Uz compinit && compinit"

		case "fish":
			completionPath = home + "/.config/fish/completions/dotfiles.fish"
			os.MkdirAll(filepath.Dir(completionPath), 0755)
			instructions = "Fish will automatically load completions from ~/.config/fish/completions/"

		default:
			fmt.Printf("❌ Unsupported shell: %s\n", shell)
			return
		}

		// Generate completion to file
		file, err := os.Create(completionPath)
		if err != nil {
			fmt.Printf("❌ Failed to create completion file: %v\n", err)
			return
		}
		defer file.Close()

		switch shell {
		case "bash":
			err = rootCmd.GenBashCompletion(file)
		case "zsh":
			err = rootCmd.GenZshCompletion(file)
		case "fish":
			err = rootCmd.GenFishCompletion(file, true)
		}

		if err != nil {
			fmt.Printf("❌ Failed to generate completion: %v\n", err)
			return
		}

		fmt.Println("✅ Completion script installed!")
		fmt.Println()
		fmt.Printf("   Location: %s\n", completionPath)
		fmt.Println()
		if instructions != "" {
			fmt.Println("💡 Setup:")
			fmt.Println("   " + instructions)
			fmt.Println()
			fmt.Println("   Then restart your shell or run:")
			fmt.Println("   source ~/.bashrc  # or ~/.zshrc")
		}
	},
}

func init() {
	completionCmd.AddCommand(completionInstallCmd)
	rootCmd.AddCommand(completionCmd)
}
