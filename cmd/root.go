package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "A simple dotfiles manager with Brewfile support",
	Long:  `A minimal dotfiles manager that stores configurations in JSON and generates Brewfiles.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
