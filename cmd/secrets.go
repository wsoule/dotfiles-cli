package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use: "secrets",
	GroupID: "advanced",
	Short: " Manage secrets and environment variables",
	Long: ` Secret Management

Safely manage API keys, tokens, and other sensitive data.
Secrets are stored in ~/.dotfiles/private/.env-private (git-ignored).

Examples:
  dotfiles secrets init               # Initialize secrets file
  dotfiles secrets add API_KEY        # Add a new secret
  dotfiles secrets list               # List secret names (not values)
  dotfiles secrets template           # Generate template for new machines`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles secrets init' to get started")
	},
}

var secretsInitCmd = &cobra.Command{
	Use: "init",
	Short: "Initialize secrets management",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := GetDotfilesDir(cmd)
		privateDir := filepath.Join(dotfilesDir, "private")
		envPrivate := filepath.Join(privateDir, ".env-private")

		// Create private directory if needed
		if err := os.MkdirAll(privateDir, 0700); err != nil {
			fmt.Printf(" Failed to create private directory: %v\n", err)
			return
		}

		// Check if .env-private already exists
		if _, err := os.Stat(envPrivate); err == nil {
			fmt.Println(" Secrets file already exists")
			fmt.Printf("%s\n", envPrivate)
			return
		}

		// Create template file
		template := `# Private Environment Variables
# This file is git-ignored and should contain sensitive data
# Add your secrets below (one per line):

# Example:
# export API_KEY="your-api-key-here"
# export DATABASE_URL="your-db-connection"
# export AWS_ACCESS_KEY_ID="your-aws-key"

`

		if err := os.WriteFile(envPrivate, []byte(template), 0600); err != nil {
			fmt.Printf(" Failed to create secrets file: %v\n", err)
			return
		}

		fmt.Println(" Secrets file initialized")
		fmt.Printf("%s\n", envPrivate)
		fmt.Println()
		fmt.Println(" Add secrets with: dotfiles secrets add <NAME>")
	},
}

var secretsAddCmd = &cobra.Command{
	Use: "add <name>",
	Short: "Add a new secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretName := strings.ToUpper(args[0])

		dotfilesDir := GetDotfilesDir(cmd)
		envPrivate := filepath.Join(dotfilesDir, "private", ".env-private")

		// Check if file exists
		if _, err := os.Stat(envPrivate); os.IsNotExist(err) {
			fmt.Println(" Secrets file not initialized")
			fmt.Println(" Run: dotfiles secrets init")
			return
		}

		// Prompt for secret value
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter value for %s: ", secretName)
		secretValue, _ := reader.ReadString('\n')
		secretValue = strings.TrimSpace(secretValue)

		if secretValue == "" {
			fmt.Println(" Value cannot be empty")
			return
		}

		// Append to file
		f, err := os.OpenFile(envPrivate, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf(" Failed to open secrets file: %v\n", err)
			return
		}
		defer f.Close()

		line := fmt.Sprintf("export %s=\"%s\"\n", secretName, secretValue)
		if _, err := f.WriteString(line); err != nil {
			fmt.Printf(" Failed to write secret: %v\n", err)
			return
		}

		fmt.Printf(" Secret '%s' added\n", secretName)
		fmt.Println(" It will be available in your shell after: source ~/.env-private")
	},
}

var secretsListCmd = &cobra.Command{
	Use: "list",
	Short: "List secret names (not values)",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := GetDotfilesDir(cmd)
		envPrivate := filepath.Join(dotfilesDir, "private", ".env-private")

		if _, err := os.Stat(envPrivate); os.IsNotExist(err) {
			fmt.Println(" Secrets file not initialized")
			return
		}

		content, err := os.ReadFile(envPrivate)
		if err != nil {
			fmt.Printf(" Failed to read secrets file: %v\n", err)
			return
		}

		lines := strings.Split(string(content), "\n")
		secrets := []string{}

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "export ") {
				parts := strings.SplitN(line[7:], "=", 2)
				if len(parts) > 0 {
					secrets = append(secrets, parts[0])
				}
			}
		}

		if len(secrets) == 0 {
			fmt.Println("No secrets defined yet")
			return
		}

		fmt.Println(" Defined secrets:")
		for _, secret := range secrets {
			fmt.Printf("• %s\n", secret)
		}
	},
}

var secretsTemplateCmd = &cobra.Command{
	Use: "template",
	Short: "Generate secrets template for new machines",
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDir := GetDotfilesDir(cmd)
		envPrivate := filepath.Join(dotfilesDir, "private", ".env-private")
		templateFile := filepath.Join(dotfilesDir, "private", ".env-private.template")

		if _, err := os.Stat(envPrivate); os.IsNotExist(err) {
			fmt.Println(" Secrets file not initialized")
			return
		}

		content, err := os.ReadFile(envPrivate)
		if err != nil {
			fmt.Printf(" Failed to read secrets file: %v\n", err)
			return
		}

		lines := strings.Split(string(content), "\n")
		template := "# Secrets Template\n# Fill in values for a new machine\n\n"

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "export ") {
				parts := strings.SplitN(line[7:], "=", 2)
				if len(parts) > 0 {
					template += fmt.Sprintf("export %s=\"\"\n", parts[0])
				}
			} else if strings.HasPrefix(line, "#") || line == "" {
				template += line + "\n"
			}
		}

		if err := os.WriteFile(templateFile, []byte(template), 0644); err != nil {
			fmt.Printf(" Failed to write template: %v\n", err)
			return
		}

		fmt.Println(" Template generated")
		fmt.Printf("%s\n", templateFile)
		fmt.Println()
		fmt.Println(" Commit this template to git (values removed)")
		fmt.Println("On new machines, copy template to .env-private and fill in values")
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(secretsInitCmd)
	secretsCmd.AddCommand(secretsAddCmd)
	secretsCmd.AddCommand(secretsListCmd)
	secretsCmd.AddCommand(secretsTemplateCmd)
}
