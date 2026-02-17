package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	versionInfo = struct {
		Version string
		Commit  string
		Date    string
	}{
		Version: "dev",
		Commit: "none",
		Date: "unknown",
	}
)

var rootCmd = &cobra.Command{
	Use: "dotfiles",
	Short: " Complete developer environment management toolkit",
	Long: ` Dotfiles Manager - Developer Onboarding Toolkit

A comprehensive command-line tool that manages your entire development environment:
• Homebrew packages with smart curation
• Dotfiles management with GNU Stow
• GitHub SSH setup and configuration
• Complete developer onboarding automation
• Configuration sharing and templates

Perfect for new developers or setting up fresh machines.

Quick Start:
  dotfiles onboard                    # Complete setup for new developers
  dotfiles setup <repo-url>           # Setup from existing dotfiles repo
  dotfiles tui                        # Interactive package manager
  dotfiles add git curl tmux          # Add essential packages

Tutorial:
  dotfiles tutorial                   # Step-by-step interactive tutorial

Get started: https://github.com/wsoule/dotfiles-cli`,
}

func init() {
	// Add command groups for better organization
	rootCmd.AddGroup(&cobra.Group{
		ID: "getting-started",
		Title: "Getting Started Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID: "package",
		Title: "Package Management:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID: "dotfiles",
		Title: "Dotfile Management:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID: "advanced",
		Title: "Advanced Commands:",
	})

	// Add persistent flags
	rootCmd.PersistentFlags().StringP("dotfiles-dir", "d", "", "Custom dotfiles directory (default: ~/.dotfiles)")
}

// SetVersionInfo sets the version information from main
func SetVersionInfo(version, commit, date string) {
	versionInfo.Version = version
	versionInfo.Commit = commit
	versionInfo.Date = date
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
