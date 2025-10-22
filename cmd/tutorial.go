package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var tutorialCmd = &cobra.Command{
	Use:     "tutorial",
	GroupID: "getting-started",
	Short:   "📖 Interactive step-by-step tutorial",
	Long: `📖 Interactive Tutorial - Learn Dotfiles Manager

An interactive, step-by-step tutorial that teaches you:
• How to use the dotfiles manager
• Best practices for managing your development environment
• Common workflows and commands
• Tips and tricks for power users

Perfect for first-time users or those wanting to learn all features.`,
	Run: runTutorial,
}

func init() {
	rootCmd.AddCommand(tutorialCmd)
}

func runTutorial(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	// Welcome
	clearScreen()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  📖 Welcome to the Dotfiles Manager Tutorial!                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("This interactive tutorial will guide you through:")
	fmt.Println("  1. Understanding what dotfiles are")
	fmt.Println("  2. Basic package management")
	fmt.Println("  3. Dotfile symlinking with stow")
	fmt.Println("  4. Syncing with Git")
	fmt.Println("  5. Advanced features")
	fmt.Println()
	pressEnterToContinue(reader)

	// Step 1: What are dotfiles?
	clearScreen()
	fmt.Println("═══ Step 1: What are Dotfiles? ═══")
	fmt.Println()
	fmt.Println("Dotfiles are configuration files that start with a dot (.)")
	fmt.Println("Examples: .zshrc, .gitconfig, .vimrc")
	fmt.Println()
	fmt.Println("This tool helps you:")
	fmt.Println("  • 📦 Manage packages (Homebrew, apt, pacman, yum)")
	fmt.Println("  • 🔗 Create symlinks to your configs (using GNU Stow)")
	fmt.Println("  • 🔄 Sync everything with Git")
	fmt.Println("  • 🚀 Quickly set up new machines")
	fmt.Println()
	pressEnterToContinue(reader)

	// Step 2: Getting Started
	clearScreen()
	fmt.Println("═══ Step 2: Getting Started ═══")
	fmt.Println()
	fmt.Println("Quick Start Commands:")
	fmt.Println()
	fmt.Println("  🎯 dotfiles onboard")
	fmt.Println("     → Complete guided setup for new users")
	fmt.Println()
	fmt.Println("  📥 dotfiles setup <repo-url>")
	fmt.Println("     → Clone and set up from existing dotfiles repo")
	fmt.Println()
	fmt.Println("  🎨 dotfiles tui")
	fmt.Println("     → Interactive visual package manager")
	fmt.Println()
	fmt.Println("Choose 'onboard' if this is your first time!")
	fmt.Println()
	pressEnterToContinue(reader)

	// Step 3: Package Management
	clearScreen()
	fmt.Println("═══ Step 3: Managing Packages ═══")
	fmt.Println()
	fmt.Println("Add packages to your configuration:")
	fmt.Println("  $ dotfiles add git curl tmux")
	fmt.Println()
	fmt.Println("Install all configured packages:")
	fmt.Println("  $ dotfiles install")
	fmt.Println()
	fmt.Println("Update installed packages:")
	fmt.Println("  $ dotfiles update")
	fmt.Println()
	fmt.Println("Check what's installed:")
	fmt.Println("  $ dotfiles status")
	fmt.Println()
	fmt.Println("Remove packages:")
	fmt.Println("  $ dotfiles remove tmux")
	fmt.Println()
	pressEnterToContinue(reader)

	// Step 4: Dotfile Symlinking
	clearScreen()
	fmt.Println("═══ Step 4: Dotfile Symlinking with Stow ═══")
	fmt.Println()
	fmt.Println("GNU Stow creates symlinks from your dotfiles repo to your home directory.")
	fmt.Println()
	fmt.Println("Structure example:")
	fmt.Println("  ~/.dotfiles/")
	fmt.Println("    ├── zsh/")
	fmt.Println("    │   └── .zshrc         → symlinks to ~/.zshrc")
	fmt.Println("    ├── git/")
	fmt.Println("    │   └── .gitconfig     → symlinks to ~/.gitconfig")
	fmt.Println("    └── vim/")
	fmt.Println("        └── .vimrc         → symlinks to ~/.vimrc")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  $ dotfiles stow           # Create all symlinks")
	fmt.Println("  $ dotfiles unstow         # Remove all symlinks")
	fmt.Println("  $ dotfiles restow         # Recreate symlinks")
	fmt.Println()
	pressEnterToContinue(reader)

	// Step 5: Git Sync
	clearScreen()
	fmt.Println("═══ Step 5: Syncing with Git ═══")
	fmt.Println()
	fmt.Println("Keep your configuration in sync across machines:")
	fmt.Println()
	fmt.Println("  $ dotfiles sync")
	fmt.Println("     → Commits changes and pushes to remote")
	fmt.Println()
	fmt.Println("Set up GitHub SSH:")
	fmt.Println("  $ dotfiles github")
	fmt.Println("     → Generate SSH keys and configure GitHub")
	fmt.Println()
	pressEnterToContinue(reader)

	// Step 6: Useful Commands
	clearScreen()
	fmt.Println("═══ Step 6: Other Useful Commands ═══")
	fmt.Println()
	fmt.Println("🏥 Health Check:")
	fmt.Println("   $ dotfiles doctor           # Check your setup")
	fmt.Println()
	fmt.Println("📊 View Differences:")
	fmt.Println("   $ dotfiles diff             # Show config vs installed")
	fmt.Println()
	fmt.Println("🔍 Scan System:")
	fmt.Println("   $ dotfiles scan             # Find installed packages")
	fmt.Println()
	fmt.Println("🧹 Cleanup:")
	fmt.Println("   $ dotfiles cleanup          # Remove old versions")
	fmt.Println()
	fmt.Println("📸 Snapshots:")
	fmt.Println("   $ dotfiles snapshot create  # Save current state")
	fmt.Println()
	fmt.Println("🏷️  Groups:")
	fmt.Println("   $ dotfiles groups           # Organize packages by tags")
	fmt.Println()
	pressEnterToContinue(reader)

	// Step 7: Advanced Features
	clearScreen()
	fmt.Println("═══ Step 7: Advanced Features ═══")
	fmt.Println()
	fmt.Println("📚 Templates:")
	fmt.Println("   $ dotfiles templates discover")
	fmt.Println("   → Browse and use community configurations")
	fmt.Println()
	fmt.Println("🪝 Hooks:")
	fmt.Println("   $ dotfiles hooks")
	fmt.Println("   → Run custom scripts before/after operations")
	fmt.Println()
	fmt.Println("📤 Sharing:")
	fmt.Println("   $ dotfiles share")
	fmt.Println("   → Share your configuration with others")
	fmt.Println()
	fmt.Println("📥 Import/Export Profiles:")
	fmt.Println("   $ dotfiles export")
	fmt.Println("   → Create machine-specific profiles")
	fmt.Println()
	pressEnterToContinue(reader)

	// Final screen
	clearScreen()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  🎉 Tutorial Complete!                                         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Next Steps:")
	fmt.Println("  1. Run: dotfiles onboard")
	fmt.Println("  2. Try: dotfiles tui")
	fmt.Println("  3. Explore: dotfiles --help")
	fmt.Println()
	fmt.Println("For detailed help on any command:")
	fmt.Println("  $ dotfiles <command> --help")
	fmt.Println()
	fmt.Println("Documentation: https://github.com/wsoule/dotfiles-cli")
	fmt.Println()
	fmt.Println("Happy configuring! 🚀")
	fmt.Println()
}

func pressEnterToContinue(reader *bufio.Reader) {
	fmt.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
