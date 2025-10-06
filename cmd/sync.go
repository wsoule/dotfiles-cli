package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "🔄 Sync your dotfiles with remote repository",
	Long: `🔄 Sync Dotfiles Repository

Synchronize your local dotfiles with the remote repository.
Supports pulling changes from remote, pushing local changes, or both.

Examples:
  dotfiles sync              # Pull and push changes (full sync)
  dotfiles sync --pull       # Only pull changes from remote
  dotfiles sync --push       # Only push changes to remote
  dotfiles sync --auto       # Auto-commit and sync all changes`,
	Run: func(cmd *cobra.Command, args []string) {
		pullOnly, _ := cmd.Flags().GetBool("pull")
		pushOnly, _ := cmd.Flags().GetBool("push")
		autoCommit, _ := cmd.Flags().GetBool("auto")
		message, _ := cmd.Flags().GetString("message")

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		dotfilesDir := filepath.Join(home, ".dotfiles")

		// Check if .dotfiles directory exists
		if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
			fmt.Println("❌ Dotfiles directory not found at ~/.dotfiles")
			fmt.Println("💡 Run 'dotfiles setup <repo-url>' first")
			os.Exit(1)
		}

		// Check if it's a git repository
		gitDir := filepath.Join(dotfilesDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			fmt.Println("❌ Not a git repository")
			fmt.Println("💡 Run 'dotfiles setup <repo-url>' to initialize")
			os.Exit(1)
		}

		fmt.Println("🔄 Syncing dotfiles...")
		fmt.Println()

		// Change to dotfiles directory
		if err := os.Chdir(dotfilesDir); err != nil {
			fmt.Printf("❌ Error changing to dotfiles directory: %v\n", err)
			os.Exit(1)
		}

		// Check for uncommitted changes
		hasChanges := checkGitStatus()

		if hasChanges && autoCommit {
			fmt.Println("📝 Auto-committing changes...")
			commitMsg := message
			if commitMsg == "" {
				commitMsg = "Auto-sync: Update dotfiles configuration"
			}
			if err := gitCommit(commitMsg); err != nil {
				fmt.Printf("❌ Failed to commit changes: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Changes committed")
			fmt.Println()
		} else if hasChanges && !pushOnly {
			fmt.Println("⚠️  You have uncommitted changes:")
			runGitCommand("git", "status", "--short")
			fmt.Println()
			fmt.Println("💡 Commit your changes first or use --auto to commit automatically")
			if !pullOnly {
				os.Exit(1)
			}
		}

		// Pull changes from remote
		if !pushOnly {
			fmt.Println("⬇️  Pulling changes from remote...")
			if err := gitPull(); err != nil {
				fmt.Printf("❌ Failed to pull changes: %v\n", err)
				fmt.Println("💡 You may need to resolve conflicts manually")
				os.Exit(1)
			}
			fmt.Println("✅ Pulled latest changes")
			fmt.Println()
		}

		// Push changes to remote
		if !pullOnly {
			fmt.Println("⬆️  Pushing changes to remote...")
			if err := gitPush(); err != nil {
				fmt.Printf("❌ Failed to push changes: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Pushed changes to remote")
			fmt.Println()
		}

		fmt.Println("🎉 Sync complete!")
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Println("   • Check status: dotfiles status")
		fmt.Println("   • Install new packages: dotfiles install")
	},
}

func checkGitStatus() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func gitCommit(message string) error {
	// Add all changes
	addCmd := exec.Command("git", "add", ".")
	if err := addCmd.Run(); err != nil {
		return err
	}

	// Commit with message
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	return commitCmd.Run()
}

func gitPull() error {
	cmd := exec.Command("git", "pull", "--rebase")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runGitCommand(command ...string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func init() {
	syncCmd.Flags().Bool("pull", false, "Only pull changes from remote")
	syncCmd.Flags().Bool("push", false, "Only push changes to remote")
	syncCmd.Flags().Bool("auto", false, "Automatically commit all changes before syncing")
	syncCmd.Flags().StringP("message", "m", "", "Commit message (used with --auto)")

	rootCmd.AddCommand(syncCmd)
}
