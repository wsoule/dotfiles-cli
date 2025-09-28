package cmd

import (
	"fmt"
	"os"

	"dotfiles/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "Modern dotfiles management system",
	Long: `A modern dotfiles management system built in Go.

This tool helps you manage your development environment configuration
with preset support, interactive setup, and automated installation.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Banner()
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dotfiles/config.json)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")

	// Version flag
	rootCmd.Version = "1.0.0"
	rootCmd.SetVersionTemplate(`{{printf "%s version %s\n" .Use .Version}}`)

	// Enable shell completions
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.CompletionOptions.DisableNoDescFlag = false
	rootCmd.CompletionOptions.DisableDescriptions = false
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dotfiles" (without extension).
		viper.AddConfigPath(home + "/.dotfiles")
		viper.AddConfigPath("./config")
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
