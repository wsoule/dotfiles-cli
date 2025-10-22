package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// GetDotfilesDir returns the dotfiles directory, either from flag or default
func GetDotfilesDir(cmd *cobra.Command) string {
	dir, _ := cmd.Flags().GetString("dotfiles-dir")
	if dir != "" {
		return dir
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".dotfiles")
}

// GetConfigPath returns the path to config.json
func GetConfigPath(cmd *cobra.Command) string {
	return filepath.Join(GetDotfilesDir(cmd), "config.json")
}
