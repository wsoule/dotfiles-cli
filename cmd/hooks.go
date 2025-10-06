package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "🪝 Manage pre/post operation hooks",
	Long: `🪝 Hooks Management

Configure shell commands to run before/after specific operations.
Hooks allow you to automate custom tasks during dotfiles operations.

Available hook types:
• pre_install / post_install - Before/after package installation
• pre_sync / post_sync - Before/after repo sync
• pre_stow / post_stow - Before/after stowing dotfiles

Examples:
  dotfiles hooks list                                    # List all hooks
  dotfiles hooks add pre_install "brew update"           # Add pre-install hook
  dotfiles hooks add post_install "echo 'Done!'"         # Add post-install hook
  dotfiles hooks remove pre_install 0                    # Remove first pre-install hook
  dotfiles hooks clear post_sync                         # Remove all post-sync hooks`,
}

var hooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured hooks",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Hooks == nil || isHooksEmpty(cfg.Hooks) {
			fmt.Println("🪝 No hooks configured")
			fmt.Println()
			fmt.Println("💡 Add a hook:")
			fmt.Println("   dotfiles hooks add pre_install 'brew update'")
			return
		}

		fmt.Println("🪝 Configured Hooks:")
		fmt.Println("=" + strings.Repeat("=", 19))
		fmt.Println()

		printHookSection("Pre-Install", cfg.Hooks.PreInstall)
		printHookSection("Post-Install", cfg.Hooks.PostInstall)
		printHookSection("Pre-Sync", cfg.Hooks.PreSync)
		printHookSection("Post-Sync", cfg.Hooks.PostSync)
		printHookSection("Pre-Stow", cfg.Hooks.PreStow)
		printHookSection("Post-Stow", cfg.Hooks.PostStow)
	},
}

var hooksAddCmd = &cobra.Command{
	Use:   "add <hook-type> <command>",
	Short: "Add a new hook",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		hookType := args[0]
		command := args[1]

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Hooks == nil {
			cfg.Hooks = &config.Hooks{}
		}

		switch hookType {
		case "pre_install":
			cfg.Hooks.PreInstall = append(cfg.Hooks.PreInstall, command)
		case "post_install":
			cfg.Hooks.PostInstall = append(cfg.Hooks.PostInstall, command)
		case "pre_sync":
			cfg.Hooks.PreSync = append(cfg.Hooks.PreSync, command)
		case "post_sync":
			cfg.Hooks.PostSync = append(cfg.Hooks.PostSync, command)
		case "pre_stow":
			cfg.Hooks.PreStow = append(cfg.Hooks.PreStow, command)
		case "post_stow":
			cfg.Hooks.PostStow = append(cfg.Hooks.PostStow, command)
		default:
			fmt.Printf("❌ Invalid hook type: %s\n", hookType)
			fmt.Println("Valid types: pre_install, post_install, pre_sync, post_sync, pre_stow, post_stow")
			os.Exit(1)
		}

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Added %s hook: %s\n", hookType, command)
	},
}

var hooksRemoveCmd = &cobra.Command{
	Use:   "remove <hook-type> <index>",
	Short: "Remove a hook by index",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		hookType := args[0]
		index := 0
		fmt.Sscanf(args[1], "%d", &index)

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Hooks == nil {
			fmt.Println("❌ No hooks configured")
			os.Exit(1)
		}

		var hooks *[]string
		switch hookType {
		case "pre_install":
			hooks = &cfg.Hooks.PreInstall
		case "post_install":
			hooks = &cfg.Hooks.PostInstall
		case "pre_sync":
			hooks = &cfg.Hooks.PreSync
		case "post_sync":
			hooks = &cfg.Hooks.PostSync
		case "pre_stow":
			hooks = &cfg.Hooks.PreStow
		case "post_stow":
			hooks = &cfg.Hooks.PostStow
		default:
			fmt.Printf("❌ Invalid hook type: %s\n", hookType)
			os.Exit(1)
		}

		if index < 0 || index >= len(*hooks) {
			fmt.Printf("❌ Invalid index %d (max: %d)\n", index, len(*hooks)-1)
			os.Exit(1)
		}

		removed := (*hooks)[index]
		*hooks = append((*hooks)[:index], (*hooks)[index+1:]...)

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Removed hook: %s\n", removed)
	},
}

var hooksClearCmd = &cobra.Command{
	Use:   "clear <hook-type>",
	Short: "Clear all hooks of a specific type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hookType := args[0]

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, ".dotfiles", "config.json")
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if cfg.Hooks == nil {
			fmt.Println("❌ No hooks configured")
			os.Exit(1)
		}

		var count int
		switch hookType {
		case "pre_install":
			count = len(cfg.Hooks.PreInstall)
			cfg.Hooks.PreInstall = []string{}
		case "post_install":
			count = len(cfg.Hooks.PostInstall)
			cfg.Hooks.PostInstall = []string{}
		case "pre_sync":
			count = len(cfg.Hooks.PreSync)
			cfg.Hooks.PreSync = []string{}
		case "post_sync":
			count = len(cfg.Hooks.PostSync)
			cfg.Hooks.PostSync = []string{}
		case "pre_stow":
			count = len(cfg.Hooks.PreStow)
			cfg.Hooks.PreStow = []string{}
		case "post_stow":
			count = len(cfg.Hooks.PostStow)
			cfg.Hooks.PostStow = []string{}
		default:
			fmt.Printf("❌ Invalid hook type: %s\n", hookType)
			os.Exit(1)
		}

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Cleared %d %s hook(s)\n", count, hookType)
	},
}

func printHookSection(name string, hooks []string) {
	if len(hooks) == 0 {
		return
	}

	fmt.Printf("📌 %s:\n", name)
	for i, hook := range hooks {
		fmt.Printf("   %d. %s\n", i, hook)
	}
	fmt.Println()
}

func isHooksEmpty(hooks *config.Hooks) bool {
	return len(hooks.PreInstall) == 0 &&
		len(hooks.PostInstall) == 0 &&
		len(hooks.PreSync) == 0 &&
		len(hooks.PostSync) == 0 &&
		len(hooks.PreStow) == 0 &&
		len(hooks.PostStow) == 0
}

// RunHooks executes a list of hook commands
func RunHooks(hooks []string, hookType string) error {
	if len(hooks) == 0 {
		return nil
	}

	fmt.Printf("🪝 Running %s hooks...\n", hookType)
	for i, hook := range hooks {
		fmt.Printf("   [%d/%d] %s\n", i+1, len(hooks), hook)

		cmd := exec.Command("sh", "-c", hook)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook failed: %v", err)
		}
	}
	fmt.Println()
	return nil
}

func init() {
	hooksCmd.AddCommand(hooksListCmd)
	hooksCmd.AddCommand(hooksAddCmd)
	hooksCmd.AddCommand(hooksRemoveCmd)
	hooksCmd.AddCommand(hooksClearCmd)

	rootCmd.AddCommand(hooksCmd)
}
