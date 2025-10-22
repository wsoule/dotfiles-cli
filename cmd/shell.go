package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:     "shell",
	GroupID: "advanced",
	Short:   "🐚 Interactive shell configuration",
	Long: `🐚 Interactive Shell Setup

Configure your shell with popular tools:
• oh-my-zsh / oh-my-fish
• Starship prompt
• fzf (fuzzy finder)
• zsh-autosuggestions
• zsh-syntax-highlighting

Examples:
  dotfiles shell setup              # Interactive setup wizard
  dotfiles shell install starship   # Install specific tool
  dotfiles shell list               # Show installed tools`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles shell setup' for interactive configuration")
	},
}

var shellSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive shell setup wizard",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("🐚 Shell Configuration Wizard")
		fmt.Println(strings.Repeat("=", 30))
		fmt.Println()

		// Detect current shell
		shell := os.Getenv("SHELL")
		shellName := "unknown"
		if strings.Contains(shell, "zsh") {
			shellName = "zsh"
		} else if strings.Contains(shell, "fish") {
			shellName = "fish"
		} else if strings.Contains(shell, "bash") {
			shellName = "bash"
		}

		fmt.Printf("Detected shell: %s\n", shellName)
		fmt.Println()

		// Ask about oh-my-zsh/oh-my-fish
		if shellName == "zsh" {
			if !isInstalled("$HOME/.oh-my-zsh") {
				fmt.Print("Install oh-my-zsh? (y/n): ")
				response, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(response)) == "y" {
					installOhMyZsh()
				}
			} else {
				fmt.Println("✅ oh-my-zsh is already installed")
			}
		} else if shellName == "fish" {
			if !isInstalled("$HOME/.local/share/omf") {
				fmt.Print("Install oh-my-fish? (y/n): ")
				response, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(response)) == "y" {
					installOhMyFish()
				}
			} else {
				fmt.Println("✅ oh-my-fish is already installed")
			}
		}

		fmt.Println()

		// Ask about Starship
		if !shellCommandExists("starship") {
			fmt.Print("Install Starship prompt? (y/n): ")
			response, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(response)) == "y" {
				installStarship()
			}
		} else {
			fmt.Println("✅ Starship is already installed")
		}

		fmt.Println()

		// Ask about fzf
		if !shellCommandExists("fzf") {
			fmt.Print("Install fzf (fuzzy finder)? (y/n): ")
			response, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(response)) == "y" {
				installFzf()
			}
		} else {
			fmt.Println("✅ fzf is already installed")
		}

		fmt.Println()

		// ZSH plugins
		if shellName == "zsh" {
			home, _ := os.UserHomeDir()
			zshautosuggestions := filepath.Join(home, ".oh-my-zsh/custom/plugins/zsh-autosuggestions")
			zshsyntax := filepath.Join(home, ".oh-my-zsh/custom/plugins/zsh-syntax-highlighting")

			if !isInstalled(zshautosuggestions) {
				fmt.Print("Install zsh-autosuggestions? (y/n): ")
				response, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(response)) == "y" {
					installZshPlugin("zsh-autosuggestions", "https://github.com/zsh-users/zsh-autosuggestions")
				}
			} else {
				fmt.Println("✅ zsh-autosuggestions is already installed")
			}

			if !isInstalled(zshsyntax) {
				fmt.Print("Install zsh-syntax-highlighting? (y/n): ")
				response, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(response)) == "y" {
					installZshPlugin("zsh-syntax-highlighting", "https://github.com/zsh-users/zsh-syntax-highlighting")
				}
			} else {
				fmt.Println("✅ zsh-syntax-highlighting is already installed")
			}
		}

		fmt.Println()
		fmt.Println("✅ Shell setup complete!")
		fmt.Println()
		fmt.Println("💡 Don't forget to:")
		fmt.Println("   • Restart your shell or run: source ~/.zshrc")
		fmt.Println("   • Configure your shell theme and plugins")
	},
}

var shellListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed shell tools",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🐚 Installed Shell Tools:")
		fmt.Println()

		// Check oh-my-zsh
		if isInstalled("$HOME/.oh-my-zsh") {
			fmt.Println("✅ oh-my-zsh")
		}

		// Check oh-my-fish
		if isInstalled("$HOME/.local/share/omf") {
			fmt.Println("✅ oh-my-fish")
		}

		// Check Starship
		if shellCommandExists("starship") {
			fmt.Println("✅ Starship")
		}

		// Check fzf
		if shellCommandExists("fzf") {
			fmt.Println("✅ fzf")
		}

		// Check zsh plugins
		home, _ := os.UserHomeDir()
		if isInstalled(filepath.Join(home, ".oh-my-zsh/custom/plugins/zsh-autosuggestions")) {
			fmt.Println("✅ zsh-autosuggestions")
		}
		if isInstalled(filepath.Join(home, ".oh-my-zsh/custom/plugins/zsh-syntax-highlighting")) {
			fmt.Println("✅ zsh-syntax-highlighting")
		}
	},
}

func installOhMyZsh() {
	fmt.Println("📦 Installing oh-my-zsh...")
	cmd := exec.Command("sh", "-c", `RUNZSH=no sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to install oh-my-zsh: %v\n", err)
	} else {
		fmt.Println("✅ oh-my-zsh installed")
	}
}

func installOhMyFish() {
	fmt.Println("📦 Installing oh-my-fish...")
	cmd := exec.Command("fish", "-c", "curl -L https://get.oh-my.fish | fish")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to install oh-my-fish: %v\n", err)
	} else {
		fmt.Println("✅ oh-my-fish installed")
	}
}

func installStarship() {
	fmt.Println("📦 Installing Starship...")

	// Try homebrew first
	if shellCommandExists("brew") {
		cmd := exec.Command("brew", "install", "starship")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Failed to install Starship: %v\n", err)
		} else {
			fmt.Println("✅ Starship installed")
			fmt.Println()
			fmt.Println("💡 Add to your shell config:")
			fmt.Println("   # ~/.zshrc or ~/.config/fish/config.fish")
			fmt.Println("   eval \"$(starship init zsh)\"  # for zsh")
			fmt.Println("   starship init fish | source   # for fish")
		}
	} else {
		fmt.Println("❌ Homebrew not found. Please install manually:")
		fmt.Println("   curl -sS https://starship.rs/install.sh | sh")
	}
}

func installFzf() {
	fmt.Println("📦 Installing fzf...")

	if shellCommandExists("brew") {
		cmd := exec.Command("brew", "install", "fzf")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Failed to install fzf: %v\n", err)
		} else {
			fmt.Println("✅ fzf installed")

			// Run fzf install script
			home, _ := os.UserHomeDir()
			installScript := filepath.Join(home, "/usr/local/opt/fzf/install")
			if _, err := os.Stat(installScript); err == nil {
				fmt.Println("📦 Running fzf install script...")
				installCmd := exec.Command(installScript, "--all")
				installCmd.Stdout = os.Stdout
				installCmd.Stderr = os.Stderr
				installCmd.Run()
			}
		}
	} else {
		fmt.Println("❌ Homebrew not found. Installing via git...")
		home, _ := os.UserHomeDir()
		fzfDir := filepath.Join(home, ".fzf")
		cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/junegunn/fzf.git", fzfDir)
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Failed to clone fzf: %v\n", err)
			return
		}

		installScript := filepath.Join(fzfDir, "install")
		installCmd := exec.Command(installScript, "--all")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to install fzf: %v\n", err)
		} else {
			fmt.Println("✅ fzf installed")
		}
	}
}

func installZshPlugin(name, repo string) {
	fmt.Printf("📦 Installing %s...\n", name)
	home, _ := os.UserHomeDir()
	pluginDir := filepath.Join(home, ".oh-my-zsh/custom/plugins", name)

	cmd := exec.Command("git", "clone", repo, pluginDir)
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to install %s: %v\n", name, err)
	} else {
		fmt.Printf("✅ %s installed\n", name)
		fmt.Println()
		fmt.Println("💡 Add to your ~/.zshrc plugins:")
		fmt.Printf("   plugins=(... %s)\n", name)
	}
}

func isInstalled(path string) bool {
	// Expand $HOME
	if strings.HasPrefix(path, "$HOME") {
		home, _ := os.UserHomeDir()
		path = strings.Replace(path, "$HOME", home, 1)
	}

	_, err := os.Stat(path)
	return err == nil
}

func shellCommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.AddCommand(shellSetupCmd)
	shellCmd.AddCommand(shellListCmd)
}
