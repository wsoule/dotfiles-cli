package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

// Built-in config templates
var configTemplates = map[string]ShareableConfig{
	"web-dev": {
		Config: config.Config{
			Taps: []string{
				"homebrew/cask-fonts",
				"homebrew/cask-versions",
			},
			Brews: []string{
				"git", "curl", "wget", "tree", "jq", "stow", "gh",
				"node", "npm", "yarn", "pnpm",
				"python", "python3", "pip3",
				"docker", "docker-compose",
				"nginx", "postgresql", "redis",
			},
			Casks: []string{
				"visual-studio-code", "google-chrome", "firefox",
				"iterm2", "rectangle", "docker",
				"figma", "postman", "tableplus",
			},
			Stow: []string{"git", "zsh", "vim", "vscode"},
		},
		Metadata: ShareMetadata{
			Name:        "Web Development",
			Description: "Complete setup for web developers with Node.js, Python, Docker, and essential tools",
			Author:      "Dotfiles Manager",
			Tags:        []string{"web-dev", "javascript", "python", "docker"},
			CreatedAt:   time.Now(),
			Version:     "1.0.0",
		},
	},
	"mobile-dev": {
		Config: config.Config{
			Taps: []string{
				"homebrew/cask-fonts",
				"dart-lang/dart",
			},
			Brews: []string{
				"git", "curl", "wget", "tree", "jq", "stow", "gh",
				"node", "npm", "yarn",
				"dart", "flutter",
				"cocoapods", "fastlane",
			},
			Casks: []string{
				"visual-studio-code", "android-studio", "xcode",
				"iterm2", "rectangle", "figma",
				"simulator", "proxyman",
			},
			Stow: []string{"git", "zsh", "vim", "vscode"},
		},
		Metadata: ShareMetadata{
			Name:        "Mobile Development",
			Description: "Setup for iOS and Android development with Flutter, React Native, and native tools",
			Author:      "Dotfiles Manager",
			Tags:        []string{"mobile-dev", "flutter", "react-native", "ios", "android"},
			CreatedAt:   time.Now(),
			Version:     "1.0.0",
		},
	},
	"data-science": {
		Config: config.Config{
			Taps: []string{
				"homebrew/cask-fonts",
			},
			Brews: []string{
				"git", "curl", "wget", "tree", "jq", "stow", "gh",
				"python", "python3", "pip3",
				"r", "jupyter", "jupyterlab",
				"postgresql", "sqlite",
			},
			Casks: []string{
				"visual-studio-code", "rstudio", "tableau-public",
				"iterm2", "rectangle", "docker",
				"jupyter-notebook-viewer",
			},
			Stow: []string{"git", "zsh", "vim", "python", "jupyter"},
		},
		Metadata: ShareMetadata{
			Name:        "Data Science",
			Description: "Python, R, Jupyter, and data analysis tools for data scientists",
			Author:      "Dotfiles Manager",
			Tags:        []string{"data-science", "python", "r", "jupyter", "analytics"},
			CreatedAt:   time.Now(),
			Version:     "1.0.0",
		},
	},
	"devops": {
		Config: config.Config{
			Taps: []string{
				"homebrew/cask-fonts",
				"hashicorp/tap",
			},
			Brews: []string{
				"git", "curl", "wget", "tree", "jq", "stow", "gh",
				"docker", "docker-compose", "kubernetes-cli",
				"terraform", "ansible", "helm",
				"aws-cli", "azure-cli", "gcloud",
				"prometheus", "grafana",
			},
			Casks: []string{
				"visual-studio-code", "iterm2", "rectangle",
				"docker", "lens", "postman",
				"aws-vault", "cyberduck",
			},
			Stow: []string{"git", "zsh", "vim", "kubectl", "terraform"},
		},
		Metadata: ShareMetadata{
			Name:        "DevOps & Cloud",
			Description: "Infrastructure, containerization, and cloud tools for DevOps engineers",
			Author:      "Dotfiles Manager",
			Tags:        []string{"devops", "cloud", "kubernetes", "terraform", "aws", "docker"},
			CreatedAt:   time.Now(),
			Version:     "1.0.0",
		},
	},
	"minimal": {
		Config: config.Config{
			Taps:  []string{},
			Brews: []string{"git", "curl", "wget", "tree", "stow", "gh"},
			Casks: []string{"visual-studio-code", "iterm2"},
			Stow:  []string{"git", "zsh", "vim"},
		},
		Metadata: ShareMetadata{
			Name:        "Minimal Setup",
			Description: "Essential tools only - perfect for lightweight development environments",
			Author:      "Dotfiles Manager",
			Tags:        []string{"minimal", "essential", "lightweight"},
			CreatedAt:   time.Now(),
			Version:     "1.0.0",
		},
	},
}

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Browse and use configuration templates",
	Long:  `Discover and apply pre-made configuration templates for different development workflows`,
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available configuration templates",
	Long:  `Show all built-in configuration templates`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📋 Available Configuration Templates")
		fmt.Println("=" + strings.Repeat("=", 35))
		fmt.Println()

		for key, template := range configTemplates {
			fmt.Printf("🏷️  %s (%s)\n", template.Metadata.Name, key)
			fmt.Printf("   📝 %s\n", template.Metadata.Description)
			fmt.Printf("   🏷️  Tags: %s\n", strings.Join(template.Metadata.Tags, ", "))
			fmt.Printf("   📦 Packages: %d brews, %d casks, %d taps, %d stow\n",
				len(template.Brews), len(template.Casks), len(template.Taps), len(template.Stow))
			fmt.Println()
		}

		fmt.Println("💡 Usage:")
		fmt.Println("  dotfiles templates show <template>  # Preview template")
		fmt.Println("  dotfiles clone template:<template>  # Apply template")
	},
}

var templatesShowCmd = &cobra.Command{
	Use:   "show <template>",
	Short: "Show details of a specific template",
	Long:  `Display detailed information about a configuration template`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]
		template, exists := configTemplates[templateName]
		if !exists {
			fmt.Printf("❌ Template '%s' not found\n", templateName)
			fmt.Println("Run 'dotfiles templates list' to see available templates")
			os.Exit(1)
		}

		fmt.Printf("📋 Template: %s\n", template.Metadata.Name)
		fmt.Printf("📝 Description: %s\n", template.Metadata.Description)
		fmt.Printf("🏷️  Tags: %s\n", strings.Join(template.Metadata.Tags, ", "))
		fmt.Printf("👤 Author: %s\n", template.Metadata.Author)
		fmt.Println()

		if len(template.Taps) > 0 {
			fmt.Printf("📋 Taps (%d):\n", len(template.Taps))
			for _, tap := range template.Taps {
				fmt.Printf("  - %s\n", tap)
			}
			fmt.Println()
		}

		if len(template.Brews) > 0 {
			fmt.Printf("🍺 Brews (%d):\n", len(template.Brews))
			for _, brew := range template.Brews {
				fmt.Printf("  - %s\n", brew)
			}
			fmt.Println()
		}

		if len(template.Casks) > 0 {
			fmt.Printf("📦 Casks (%d):\n", len(template.Casks))
			for _, cask := range template.Casks {
				fmt.Printf("  - %s\n", cask)
			}
			fmt.Println()
		}

		if len(template.Stow) > 0 {
			fmt.Printf("🔗 Stow Packages (%d):\n", len(template.Stow))
			for _, stow := range template.Stow {
				fmt.Printf("  - %s\n", stow)
			}
			fmt.Println()
		}

		fmt.Println("💡 To apply this template:")
		fmt.Printf("  dotfiles clone template:%s\n", templateName)
	},
}

// Update the clone command to handle templates
func init() {
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesShowCmd)
	rootCmd.AddCommand(templatesCmd)
}

// Add template support to clone command
func handleTemplateClone(templateName string, merge bool) error {
	template, exists := configTemplates[templateName]
	if !exists {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Show template info
	fmt.Printf("📋 Template: %s\n", template.Metadata.Name)
	fmt.Printf("📝 Description: %s\n", template.Metadata.Description)
	fmt.Println()

	if !askConfirmation("Apply this template? (y/N): ", false) {
		return fmt.Errorf("template application cancelled")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}

	configPath := filepath.Join(home, ".dotfiles", "config.json")

	if merge {
		// Load existing config and merge
		existingConfig, err := config.Load(configPath)
		if err != nil {
			fmt.Printf("⚠️  Could not load existing config, creating new: %v\n", err)
			existingConfig = &config.Config{}
		}

		// Merge packages
		existingConfig.Taps = mergeSlices(existingConfig.Taps, template.Taps)
		existingConfig.Brews = mergeSlices(existingConfig.Brews, template.Brews)
		existingConfig.Casks = mergeSlices(existingConfig.Casks, template.Casks)
		existingConfig.Stow = mergeSlices(existingConfig.Stow, template.Stow)

		if err := existingConfig.Save(configPath); err != nil {
			return fmt.Errorf("error saving merged config: %v", err)
		}
		fmt.Println("✅ Template merged with existing configuration!")
	} else {
		// Replace existing config
		newConfig := &config.Config{
			Taps:  template.Taps,
			Brews: template.Brews,
			Casks: template.Casks,
			Stow:  template.Stow,
		}

		if err := newConfig.Save(configPath); err != nil {
			return fmt.Errorf("error saving config: %v", err)
		}
		fmt.Println("✅ Template applied successfully!")
	}

	fmt.Println("💡 Next steps:")
	fmt.Println("  dotfiles status    # Check what needs to be installed")
	fmt.Println("  dotfiles install   # Install all packages")

	return nil
}