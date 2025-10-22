package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:     "template",
	GroupID: "advanced",
	Short:   "📄 Render templates with variable substitution",
	Long: `📄 Template Variable Rendering

Replace {{VARIABLE}} placeholders in files with actual values.
Variables are stored in config.json.

Template Syntax:
  {{EMAIL}}      - Simple variable
  {{NAME}}       - Another variable
  {{VAR|default}} - Variable with default value

Examples:
  dotfiles template render .gitconfig
  dotfiles template set EMAIL you@example.com
  dotfiles template list
  dotfiles template scan                     # Find all template files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'dotfiles template render <file>' to render a template")
	},
}

var templateRenderCmd = &cobra.Command{
	Use:   "render <file>",
	Short: "Render a template file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		outputPath, _ := cmd.Flags().GetString("output")
		inPlace, _ := cmd.Flags().GetBool("in-place")

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if cfg.Variables == nil || len(cfg.Variables) == 0 {
			fmt.Println("⚠️  No variables defined")
			fmt.Println("💡 Set variables with: dotfiles template set <key> <value>")
			return
		}

		// Make path absolute if relative
		if !filepath.IsAbs(filePath) {
			filePath = filepath.Join(GetDotfilesDir(cmd), "stow", filePath)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("❌ Failed to read file: %v\n", err)
			return
		}

		rendered := renderTemplate(string(content), cfg.Variables)

		// Determine output
		if inPlace {
			if err := os.WriteFile(filePath, []byte(rendered), 0644); err != nil {
				fmt.Printf("❌ Failed to write file: %v\n", err)
				return
			}
			fmt.Printf("✅ Rendered in-place: %s\n", filePath)
		} else if outputPath != "" {
			if err := os.WriteFile(outputPath, []byte(rendered), 0644); err != nil {
				fmt.Printf("❌ Failed to write output: %v\n", err)
				return
			}
			fmt.Printf("✅ Rendered to: %s\n", outputPath)
		} else {
			// Print to stdout
			fmt.Println(rendered)
		}
	},
}

var templateSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a template variable",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := strings.Join(args[1:], " ")

		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if cfg.Variables == nil {
			cfg.Variables = make(map[string]string)
		}

		cfg.Variables[key] = value

		if err := cfg.Save(configPath); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			return
		}

		fmt.Printf("✅ Set %s = %s\n", key, value)
	},
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all template variables",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		if cfg.Variables == nil || len(cfg.Variables) == 0 {
			fmt.Println("No variables defined")
			fmt.Println("💡 Set variables with: dotfiles template set <key> <value>")
			return
		}

		fmt.Println("📋 Template Variables:")
		for key, value := range cfg.Variables {
			fmt.Printf("   %s = %s\n", key, value)
		}
	},
}

var templateScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for template files and undefined variables",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := GetConfigPath(cmd)
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("❌ Failed to load config: %v\n", err)
			return
		}

		dotfilesDir := GetDotfilesDir(cmd)
		stowDir := filepath.Join(dotfilesDir, "stow")

		fmt.Println("🔍 Scanning for template variables...")
		fmt.Println()

		templateFiles := []string{}
		undefinedVars := make(map[string][]string) // var -> files

		filepath.Walk(stowDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			// Skip binary files
			if isBinary(path) {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			vars := findTemplateVars(string(content))
			if len(vars) > 0 {
				relPath, _ := filepath.Rel(stowDir, path)
				templateFiles = append(templateFiles, relPath)

				for _, v := range vars {
					if cfg.Variables == nil || cfg.Variables[v] == "" {
						undefinedVars[v] = append(undefinedVars[v], relPath)
					}
				}
			}

			return nil
		})

		if len(templateFiles) == 0 {
			fmt.Println("No template files found")
			return
		}

		fmt.Printf("📄 Found %d template file(s):\n", len(templateFiles))
		for _, f := range templateFiles {
			fmt.Printf("   • %s\n", f)
		}
		fmt.Println()

		if len(undefinedVars) > 0 {
			fmt.Printf("⚠️  %d undefined variable(s):\n", len(undefinedVars))
			for v, files := range undefinedVars {
				fmt.Printf("   • {{%s}} in:\n", v)
				for _, f := range files {
					fmt.Printf("      - %s\n", f)
				}
			}
			fmt.Println()
			fmt.Println("💡 Set variables with: dotfiles template set <key> <value>")
		} else {
			fmt.Println("✅ All variables are defined")
		}
	},
}

// renderTemplate replaces {{VAR}} with values
func renderTemplate(content string, variables map[string]string) string {
	re := regexp.MustCompile(`\{\{([A-Z_][A-Z0-9_]*?)(?:\|([^}]+))?\}\}`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		varName := parts[1]
		defaultValue := ""
		if len(parts) > 2 {
			defaultValue = parts[2]
		}

		if val, ok := variables[varName]; ok {
			return val
		}

		if defaultValue != "" {
			return defaultValue
		}

		// Return unchanged if no value found
		return match
	})
}

// findTemplateVars extracts all {{VAR}} patterns
func findTemplateVars(content string) []string {
	re := regexp.MustCompile(`\{\{([A-Z_][A-Z0-9_]*?)(?:\|[^}]+)?\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	vars := []string{}
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			v := match[1]
			if !seen[v] {
				vars = append(vars, v)
				seen[v] = true
			}
		}
	}

	return vars
}

// isBinary checks if file is likely binary
func isBinary(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := []string{".jpg", ".png", ".gif", ".pdf", ".zip", ".tar", ".gz", ".exe", ".bin"}

	for _, be := range binaryExts {
		if ext == be {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateRenderCmd)
	templateCmd.AddCommand(templateSetCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateScanCmd)

	templateRenderCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	templateRenderCmd.Flags().BoolP("in-place", "i", false, "Modify file in-place")
}
