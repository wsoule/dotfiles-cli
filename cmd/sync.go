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
	Use: "sync",
	GroupID: "dotfiles",
	Short: " Sync your dotfiles with remote repository",
	Long: ` Sync Dotfiles Repository

Synchronize your local dotfiles with the remote repository.
Supports pulling changes from remote, pushing local changes, or both.

Examples:
  dotfiles sync              # Pull and push changes (full sync)
  dotfiles sync --pull       # Only pull changes from remote
  dotfiles sync --push       # Only push changes to remote
  dotfiles sync --auto       # Auto-commit and sync all changes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pullOnly, _ := cmd.Flags().GetBool("pull")
		pushOnly, _ := cmd.Flags().GetBool("push")
		autoCommit, _ := cmd.Flags().GetBool("auto")
		message, _ := cmd.Flags().GetString("message")

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf(" error getting home directory: %w", err)
		}

		dotfilesDir := filepath.Join(home, ".dotfiles")

		// Check if .dotfiles directory exists
		if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
			return fmt.Errorf("dotfiles directory not found at ~/.dotfiles\nRun 'dotfiles setup <repo-url>' first")
		}

		// Check if it's a git repository
		gitDir := filepath.Join(dotfilesDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			return fmt.Errorf("not a git repository\nRun 'dotfiles setup <repo-url>' to initialize")
		}

		fmt.Println(" Syncing dotfiles...")
		fmt.Println()

		// Change to dotfiles directory
		if err := os.Chdir(dotfilesDir); err != nil {
			return fmt.Errorf(" error changing to dotfiles directory: %w", err)
		}

		// Check for uncommitted changes
		hasChanges := checkGitStatus()

		if hasChanges && autoCommit {
			fmt.Println(" Auto-committing changes...")
			commitMsg := message
			if commitMsg == "" {
				commitMsg = "Auto-sync: Update dotfiles configuration"
			}
			if err := gitCommit(commitMsg); err != nil {
				return fmt.Errorf(" failed to commit changes: %w", err)
			}
			fmt.Println(" Changes committed")
			fmt.Println()
		} else if hasChanges && !pushOnly {
			fmt.Println("You have uncommitted changes:")
			runGitCommand("git", "status", "--short")
			fmt.Println()
			fmt.Println("Commit your changes first or use --auto to commit automatically")
		}

		// Pull changes from remote
		if !pushOnly {
			fmt.Println("Pulling changes from remote...")
			if err := gitPull(); err != nil {
				return fmt.Errorf("failed to pull changes (you may need to resolve conflicts manually): %w", err)
			}
			fmt.Println(" Pulled latest changes")
			fmt.Println()
		}

		// Push changes to remote
		if !pullOnly {
			fmt.Println("Pushing changes to remote...")
			if err := gitPush(); err != nil {
				return fmt.Errorf(" failed to push changes: %w", err)
			}
			fmt.Println(" Pushed changes to remote")
			fmt.Println()
		}

		fmt.Println(" Sync complete!")
		fmt.Println()
		fmt.Println(" Next steps:")
		fmt.Println("• Check status: dotfiles status")
		fmt.Println("• Install new packages: dotfiles install")
		return nil
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
