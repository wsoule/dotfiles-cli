package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use: "watch",
	GroupID: "dotfiles",
	Short: "Watch for changes and auto-sync",
	Long: `  Watch Mode - Automatic Synchronization

Monitor dotfiles directory for changes and automatically:
• Commit changes to git
• Push to remote repository
• Optionally run hooks

Examples:
  dotfiles watch start                    # Start watching
  dotfiles watch start --interval 300     # Check every 5 minutes
  dotfiles watch stop                     # Stop watching`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles watch start' to begin watching for changes")
	},
}

var watchStartCmd = &cobra.Command{
	Use: "start",
	Short: "Start watching for changes",
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetInt("interval")
		noPush, _ := cmd.Flags().GetBool("no-push")

		dotfilesDir := GetDotfilesDir(cmd)

		fmt.Println("Starting watch mode...")
		fmt.Printf("Directory: %s\n", dotfilesDir)
		fmt.Printf("Check interval: %d seconds\n", interval)
		if noPush {
			fmt.Println("Auto-push: disabled")
		}
		fmt.Println()
		fmt.Println("Press Ctrl+C to stop")
		fmt.Println()

		// Create PID file
		pidFile := filepath.Join(dotfilesDir, ".watch.pid")
		if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
			fmt.Printf("Could not create PID file: %v\n", err)
		}
		defer os.Remove(pidFile)

		// Setup signal handler
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		// Initial check
		checkAndSync(dotfilesDir, noPush)

		for {
			select {
			case <-ticker.C:
				checkAndSync(dotfilesDir, noPush)

			case <-sigChan:
				fmt.Println()
				fmt.Println(" Stopping watch mode...")
				return
			}
		}
	},
}

var watchStopCmd = &cobra.Command{
	Use: "stop",
	Short: "Stop watch mode",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := GetDotfilesDir(cmd)
		pidFile := filepath.Join(dotfilesDir, ".watch.pid")

		data, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Println("Watch mode is not running")
			return
		}

		var pid int
		fmt.Sscanf(string(data), "%d", &pid)

		// Send SIGTERM to process
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf(" Could not find process: %v\n", err)
			os.Remove(pidFile)
			return
		}

		if err := process.Signal(syscall.SIGTERM); err != nil {
			fmt.Printf(" Could not stop process: %v\n", err)
			return
		}

		os.Remove(pidFile)
		fmt.Println(" Watch mode stopped")
	},
}

var watchStatusCmd = &cobra.Command{
	Use: "status",
	Short: "Check if watch mode is running",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := GetDotfilesDir(cmd)
		pidFile := filepath.Join(dotfilesDir, ".watch.pid")

		if _, err := os.Stat(pidFile); os.IsNotExist(err) {
			fmt.Println(" Watch mode is not running")
			return
		}

		data, _ := os.ReadFile(pidFile)
		var pid int
		fmt.Sscanf(string(data), "%d", &pid)

		// Check if process exists
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println(" Watch mode is not running")
			os.Remove(pidFile)
			return
		}

		// Try to send signal 0 (doesn't actually send, just checks if alive)
		if err := process.Signal(syscall.Signal(0)); err != nil {
			fmt.Println(" Watch mode is not running")
			os.Remove(pidFile)
			return
		}

		fmt.Println(" Watch mode is running")
		fmt.Printf("PID: %d\n", pid)
	},
}

func checkAndSync(dotfilesDir string, noPush bool) {
	// Check for changes
	gitCmd := exec.Command("git", "-C", dotfilesDir, "status", "--porcelain")
	output, err := gitCmd.Output()
	if err != nil {
		fmt.Printf("[%s] Failed to check git status\n", time.Now().Format("15:04:05"))
		return
	}

	if len(output) == 0 {
		// No changes
		return
	}

	fmt.Printf(" [%s] Changes detected, syncing...\n", time.Now().Format("15:04:05"))

	// Add all changes
	addCmd := exec.Command("git", "-C", dotfilesDir, "add", "-A")
	if err := addCmd.Run(); err != nil {
		fmt.Printf(" Failed to add changes: %v\n", err)
		return
	}

	// Commit
	commitMsg := fmt.Sprintf("Auto-commit: %s", time.Now().Format("2006-01-02 15:04:05"))
	commitCmd := exec.Command("git", "-C", dotfilesDir, "commit", "-m", commitMsg)
	if err := commitCmd.Run(); err != nil {
		fmt.Printf(" Failed to commit: %v\n", err)
		return
	}

	fmt.Printf(" [%s] Committed changes\n", time.Now().Format("15:04:05"))

	// Push if enabled
	if !noPush {
		pushCmd := exec.Command("git", "-C", dotfilesDir, "push")
		if err := pushCmd.Run(); err != nil {
			fmt.Printf("[%s] Failed to push: %v\n", time.Now().Format("15:04:05"), err)
		} else {
			fmt.Printf(" [%s] Pushed to remote\n", time.Now().Format("15:04:05"))
		}
	}
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.AddCommand(watchStartCmd)
	watchCmd.AddCommand(watchStopCmd)
	watchCmd.AddCommand(watchStatusCmd)

	watchStartCmd.Flags().IntP("interval", "i", 60, "Check interval in seconds")
	watchStartCmd.Flags().Bool("no-push", false, "Don't push to remote")
}
