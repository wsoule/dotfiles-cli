package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

// Template with inheritance support
type ExtendedTemplate struct {
	ShareableConfig
	Extends    string   `json:"extends,omitempty"`    // Base template to inherit from
	Overrides  []string `json:"overrides,omitempty"`  // Fields to override from base
	AddOnly    bool     `json:"addOnly,omitempty"`    // Only add packages, don't remove any
	Public     bool     `json:"public,omitempty"`     // Whether template is publicly visible
	Featured   bool     `json:"featured,omitempty"`   // Whether template is featured
}

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
	Short: "📚 Browse and use configuration templates",
	Long: `📚 Configuration Templates - Pre-built Development Environments

Discover, create, and share pre-made configuration templates for different
development workflows. Templates are reusable blueprints that set up complete
development environments with curated packages and settings.

Available Templates:
• Web Development - Node.js, Python, Docker, essential web tools
• Mobile Development - Flutter, React Native, iOS/Android tools
• Data Science - Python, R, Jupyter, analytics tools
• DevOps - Kubernetes, Terraform, cloud tools
• Minimal - Essential tools only, lightweight setup

Commands:
  templates list              # Browse built-in templates
  templates discover          # Find community templates from API
  templates show <name>       # Preview template details
  templates create <name>     # Create template from current config
  templates push <file>       # Share template with community

Examples:
  dotfiles templates list                    # See available templates
  dotfiles templates discover --search web   # Find web development templates
  dotfiles templates show web-dev           # Preview web-dev template
  dotfiles clone template:web-dev           # Apply built-in template
  dotfiles clone <api-url>                  # Apply community template`,
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

var templatesCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a custom template from current configuration",
	Long:  `Create a reusable template based on your current dotfiles configuration`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]

		// Get flags
		description, _ := cmd.Flags().GetString("description")
		author, _ := cmd.Flags().GetString("author")
		tagsList, _ := cmd.Flags().GetString("tags")
		baseTemplate, _ := cmd.Flags().GetString("extends")
		addOnly, _ := cmd.Flags().GetBool("add-only")

		push, _ := cmd.Flags().GetBool("push")

		templateFile, err := createCustomTemplate(templateName, description, author, tagsList, baseTemplate, addOnly)
		if err != nil {
			fmt.Printf("❌ Failed to create template: %v\n", err)
			os.Exit(1)
		}

		// Optionally push to API immediately
		if push {
			fmt.Println("🚀 Pushing template to API...")
			public, _ := cmd.Flags().GetBool("public")
			if err := pushTemplateToAPI(templateFile, public, false); err != nil {
				fmt.Printf("❌ Failed to push template: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

var templatesValidateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a template file",
	Long:  `Check if a template file is valid and can be applied`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateFile := args[0]

		if err := validateTemplate(templateFile); err != nil {
			fmt.Printf("❌ Template validation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Template '%s' is valid!\n", templateFile)
	},
}

var templatesPushCmd = &cobra.Command{
	Use:   "push <template-file>",
	Short: "Push a template to the shared repository",
	Long:  `Upload a custom template to the dotfiles sharing API for others to discover and use`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateFile := args[0]

		// Get flags
		public, _ := cmd.Flags().GetBool("public")
		featured, _ := cmd.Flags().GetBool("featured")

		if err := pushTemplateToAPI(templateFile, public, featured); err != nil {
			fmt.Printf("❌ Failed to push template: %v\n", err)
			os.Exit(1)
		}
	},
}

var templatesDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover templates from the community",
	Long:  `Browse and search templates shared by other users`,
	Run: func(cmd *cobra.Command, args []string) {
		search, _ := cmd.Flags().GetString("search")
		tags, _ := cmd.Flags().GetString("tags")
		featured, _ := cmd.Flags().GetBool("featured")

		if err := discoverTemplatesFromAPI(search, tags, featured); err != nil {
			fmt.Printf("❌ Failed to discover templates: %v\n", err)
			os.Exit(1)
		}
	},
}

// Update the clone command to handle templates
func init() {
	// Add flags to create command
	templatesCreateCmd.Flags().StringP("description", "d", "", "Template description")
	templatesCreateCmd.Flags().StringP("author", "a", "", "Template author")
	templatesCreateCmd.Flags().StringP("tags", "t", "", "Comma-separated tags")
	templatesCreateCmd.Flags().StringP("extends", "e", "", "Base template to extend")
	templatesCreateCmd.Flags().Bool("add-only", false, "Create add-only template (merge with existing)")
	templatesCreateCmd.Flags().Bool("push", false, "Push template to API after creation")
	templatesCreateCmd.Flags().Bool("public", true, "Make template public (when pushing)")

	// Add flags to push command
	templatesPushCmd.Flags().Bool("public", true, "Make template publicly visible")
	templatesPushCmd.Flags().Bool("featured", false, "Submit for featured templates")

	// Add flags to discover command
	templatesDiscoverCmd.Flags().StringP("search", "s", "", "Search query")
	templatesDiscoverCmd.Flags().StringP("tags", "t", "", "Filter by tags (comma-separated)")
	templatesDiscoverCmd.Flags().Bool("featured", false, "Show only featured templates")

	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesShowCmd)
	templatesCmd.AddCommand(templatesCreateCmd)
	templatesCmd.AddCommand(templatesValidateCmd)
	templatesCmd.AddCommand(templatesPushCmd)
	templatesCmd.AddCommand(templatesDiscoverCmd)
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

func createCustomTemplate(name, description, author, tagsList, baseTemplate string, addOnly bool) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %v", err)
	}

	// Load current configuration
	configPath := filepath.Join(home, ".dotfiles", "config.json")
	currentConfig, err := config.Load(configPath)
	if err != nil {
		return "", fmt.Errorf("error loading current config: %v", err)
	}

	// Interactive input if values not provided
	reader := bufio.NewReader(os.Stdin)

	if description == "" {
		fmt.Print("Enter template description: ")
		description, _ = reader.ReadString('\n')
		description = strings.TrimSpace(description)
	}

	if author == "" {
		fmt.Print("Enter author name: ")
		author, _ = reader.ReadString('\n')
		author = strings.TrimSpace(author)
	}

	if tagsList == "" {
		fmt.Print("Enter tags (comma-separated): ")
		tagsList, _ = reader.ReadString('\n')
		tagsList = strings.TrimSpace(tagsList)
	}

	// Parse tags
	var tags []string
	if tagsList != "" {
		for _, tag := range strings.Split(tagsList, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	// Create extended template
	template := ExtendedTemplate{
		ShareableConfig: ShareableConfig{
			Config: *currentConfig,
			Metadata: ShareMetadata{
				Name:        name,
				Description: description,
				Author:      author,
				Tags:        tags,
				CreatedAt:   time.Now(),
				Version:     "1.0.0",
			},
		},
		Extends: baseTemplate,
		AddOnly: addOnly,
	}

	// Save template to ~/.dotfiles/templates/
	templatesDir := filepath.Join(home, ".dotfiles", "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return "", fmt.Errorf("error creating templates directory: %v", err)
	}

	templateFile := filepath.Join(templatesDir, name+".json")
	file, err := os.Create(templateFile)
	if err != nil {
		return "", fmt.Errorf("error creating template file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(template); err != nil {
		return "", fmt.Errorf("error encoding template: %v", err)
	}

	fmt.Printf("✅ Created template '%s' at %s\n", name, templateFile)
	fmt.Println("💡 Share with: dotfiles share file " + templateFile)

	return templateFile, nil
}

func validateTemplate(templateFile string) error {
	// Check if file exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", templateFile)
	}

	// Try to parse as extended template first
	file, err := os.Open(templateFile)
	if err != nil {
		return fmt.Errorf("error opening template file: %v", err)
	}
	defer file.Close()

	var extTemplate ExtendedTemplate
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&extTemplate); err != nil {
		// Try parsing as regular ShareableConfig
		file.Seek(0, 0)
		var template ShareableConfig
		if err := json.NewDecoder(file).Decode(&template); err != nil {
			return fmt.Errorf("invalid template format: %v", err)
		}
		return validateShareableConfig(template)
	}

	// Validate extended template
	if err := validateExtendedTemplate(extTemplate); err != nil {
		return err
	}

	fmt.Println("📋 Template validation passed:")
	fmt.Printf("   Name: %s\n", extTemplate.Metadata.Name)
	fmt.Printf("   Description: %s\n", extTemplate.Metadata.Description)
	if extTemplate.Extends != "" {
		fmt.Printf("   Extends: %s\n", extTemplate.Extends)
	}
	fmt.Printf("   Packages: %d brews, %d casks, %d taps, %d stow\n",
		len(extTemplate.Brews), len(extTemplate.Casks), len(extTemplate.Taps), len(extTemplate.Stow))

	return nil
}

func validateExtendedTemplate(template ExtendedTemplate) error {
	// Validate base ShareableConfig
	if err := validateShareableConfig(template.ShareableConfig); err != nil {
		return err
	}

	// Validate inheritance
	if template.Extends != "" {
		if _, exists := configTemplates[template.Extends]; !exists {
			return fmt.Errorf("base template '%s' not found", template.Extends)
		}
	}

	return nil
}

func validateShareableConfig(config ShareableConfig) error {
	// Check required metadata
	if config.Metadata.Name == "" {
		return fmt.Errorf("template name is required")
	}

	if config.Metadata.Description == "" {
		return fmt.Errorf("template description is required")
	}

	if config.Metadata.Author == "" {
		return fmt.Errorf("template author is required")
	}

	// Validate package names (basic check)
	allPackages := append(append(append(config.Brews, config.Casks...), config.Taps...), config.Stow...)
	for _, pkg := range allPackages {
		if strings.TrimSpace(pkg) == "" {
			return fmt.Errorf("empty package name found")
		}
		if strings.Contains(pkg, " ") {
			return fmt.Errorf("invalid package name (contains spaces): %s", pkg)
		}
	}

	return nil
}

func resolveTemplateInheritance(templateName string) (*ShareableConfig, error) {
	// Check if it's a built-in template
	if template, exists := configTemplates[templateName]; exists {
		return &template, nil
	}

	// Try loading from custom templates
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %v", err)
	}

	templateFile := filepath.Join(home, ".dotfiles", "templates", templateName+".json")
	file, err := os.Open(templateFile)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}
	defer file.Close()

	var extTemplate ExtendedTemplate
	if err := json.NewDecoder(file).Decode(&extTemplate); err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	// Resolve inheritance
	if extTemplate.Extends != "" {
		baseTemplate, err := resolveTemplateInheritance(extTemplate.Extends)
		if err != nil {
			return nil, fmt.Errorf("error resolving base template '%s': %v", extTemplate.Extends, err)
		}

		// Merge with base template
		if extTemplate.AddOnly {
			// Only add packages, don't remove
			extTemplate.Taps = mergeSlices(baseTemplate.Taps, extTemplate.Taps)
			extTemplate.Brews = mergeSlices(baseTemplate.Brews, extTemplate.Brews)
			extTemplate.Casks = mergeSlices(baseTemplate.Casks, extTemplate.Casks)
			extTemplate.Stow = mergeSlices(baseTemplate.Stow, extTemplate.Stow)
		} else {
			// Replace packages
			if len(extTemplate.Taps) == 0 {
				extTemplate.Taps = baseTemplate.Taps
			}
			if len(extTemplate.Brews) == 0 {
				extTemplate.Brews = baseTemplate.Brews
			}
			if len(extTemplate.Casks) == 0 {
				extTemplate.Casks = baseTemplate.Casks
			}
			if len(extTemplate.Stow) == 0 {
				extTemplate.Stow = baseTemplate.Stow
			}
		}
	}

	return &extTemplate.ShareableConfig, nil
}

func pushTemplateToAPI(templateFile string, public, featured bool) error {
	file, err := os.Open(templateFile)
	if err != nil {
		return fmt.Errorf("failed to open template file: %v", err)
	}
	defer file.Close()

	var template ExtendedTemplate
	if err := json.NewDecoder(file).Decode(&template); err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	template.Public = public
	template.Featured = featured

	jsonData, err := json.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to encode template: %v", err)
	}

	apiURL := os.Getenv("DOTFILES_API_URL")
	if apiURL == "" {
		apiURL = "https://new-dotfiles-production.up.railway.app"
	}

	resp, err := http.Post(apiURL+"/api/templates", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to push template: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	fmt.Printf("✅ Template pushed successfully!\n")
	fmt.Printf("🆔 Template ID: %s\n", result.ID)
	if result.URL != "" {
		fmt.Printf("🌐 URL: %s\n", result.URL)
	}

	return nil
}

func discoverTemplatesFromAPI(search, tags string, featured bool) error {
	apiURL := os.Getenv("DOTFILES_API_URL")
	if apiURL == "" {
		apiURL = "https://new-dotfiles-production.up.railway.app"
	}

	apiEndpoint := apiURL + "/api/templates"
	params := make([]string, 0)

	if search != "" {
		params = append(params, "search="+url.QueryEscape(search))
	}
	if tags != "" {
		params = append(params, "tags="+url.QueryEscape(tags))
	}
	if featured {
		params = append(params, "featured=true")
	}

	if len(params) > 0 {
		apiEndpoint += "?" + strings.Join(params, "&")
	}

	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return fmt.Errorf("failed to fetch templates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Templates []struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Author      string   `json:"author"`
			Tags        []string `json:"tags"`
			Featured    bool     `json:"featured"`
			Downloads   int      `json:"downloads"`
			UpdatedAt   string   `json:"updated_at"`
		} `json:"templates"`
		Total int `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Total == 0 {
		fmt.Println("📭 No templates found matching your criteria")
		return nil
	}

	fmt.Printf("🔍 Found %d template(s):\n\n", result.Total)

	for _, tmpl := range result.Templates {
		fmt.Printf("📦 %s", tmpl.Name)
		if tmpl.Featured {
			fmt.Printf(" ⭐")
		}
		fmt.Printf("\n")

		if tmpl.Description != "" {
			fmt.Printf("   %s\n", tmpl.Description)
		}

		fmt.Printf("   👤 Author: %s", tmpl.Author)
		if len(tmpl.Tags) > 0 {
			fmt.Printf(" | 🏷️  Tags: %s", strings.Join(tmpl.Tags, ", "))
		}
		fmt.Printf(" | 📥 Downloads: %d\n", tmpl.Downloads)

		fmt.Printf("   💾 Clone: dotfiles clone https://new-dotfiles-production.up.railway.app/api/templates/%s\n", tmpl.ID)
		fmt.Println()
	}

	return nil
}