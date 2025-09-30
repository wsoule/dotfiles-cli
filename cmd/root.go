package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "🚀 Complete developer environment management toolkit",
	Long: `🚀 Dotfiles Manager - Developer Onboarding Toolkit

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
  dotfiles templates discover         # Browse community templates
  dotfiles add git curl tmux          # Add essential packages

Get started: https://github.com/wsoule/new-dotfiles`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
