package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dotfiles/internal/config"
	"github.com/spf13/cobra"
)

// ShareableConfig represents a config that can be shared
type ShareableConfig struct {
	config.Config
	Metadata ShareMetadata `json:"metadata"`
}

type ShareMetadata struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	Version     string    `json:"version"`
}

type GistResponse struct {
	ID          string            `json:"id"`
	HTMLURL     string            `json:"html_url"`
	Files       map[string]GistFile `json:"files"`
	Description string            `json:"description"`
	Public      bool              `json:"public"`
}

type GistFile struct {
	Content string `json:"content"`
}

type GistRequest struct {
	Description string            `json:"description"`
	Public      bool              `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share your configuration with others",
	Long:  `Share your dotfiles configuration via GitHub Gist or export to file`,
}

var shareGistCmd = &cobra.Command{
	Use:   "gist",
	Short: "Share configuration via GitHub Gist",
	Long:  `Upload your configuration to GitHub Gist for easy sharing`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		author, _ := cmd.Flags().GetString("author")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		private, _ := cmd.Flags().GetBool("private")

		if name == "" {
			fmt.Println("❌ Config name is required. Use --name flag.")
			os.Exit(1)
		}

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

		// Create shareable config with metadata
		shareableConfig := ShareableConfig{
			Config: *cfg,
			Metadata: ShareMetadata{
				Name:        name,
				Description: description,
				Author:      author,
				Tags:        tags,
				CreatedAt:   time.Now(),
				Version:     "1.0.0",
			},
		}

		fmt.Printf("📤 Sharing config '%s'...\n", name)

		// Try uploading to web app first
		webAppURL, err := uploadToWebApp(shareableConfig, !private)
		if err == nil {
			fmt.Printf("✅ Config shared to web app successfully!\n")
			fmt.Printf("🔗 Web App URL: %s\n", webAppURL)
			fmt.Printf("📋 To clone this config: dotfiles clone %s\n", webAppURL)
		} else {
			fmt.Printf("⚠️  Web app upload failed, trying GitHub Gist: %v\n", err)

			// Fallback to GitHub Gist
			gistURL, err := uploadToGist(shareableConfig, !private)
			if err != nil {
				fmt.Printf("❌ Failed to upload to Gist: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✅ Config shared to GitHub Gist successfully!\n")
			fmt.Printf("🔗 Gist URL: %s\n", gistURL)
			fmt.Printf("📋 To clone this config: dotfiles clone %s\n", gistURL)
		}

		// Copy URL to clipboard
		var finalURL string
		if webAppURL != "" {
			finalURL = webAppURL
		} else {
			// finalURL = gistURL
		}
		if err := copyToClipboard(finalURL); err == nil {
			fmt.Println("📋 URL copied to clipboard!")
		}
	},
}

var shareFileCmd = &cobra.Command{
	Use:   "file <output-file>",
	Short: "Export configuration to a shareable file",
	Long:  `Export your configuration to a JSON file that others can import`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := args[0]
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		author, _ := cmd.Flags().GetString("author")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if name == "" {
			fmt.Println("❌ Config name is required. Use --name flag.")
			os.Exit(1)
		}

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

		// Create shareable config with metadata
		shareableConfig := ShareableConfig{
			Config: *cfg,
			Metadata: ShareMetadata{
				Name:        name,
				Description: description,
				Author:      author,
				Tags:        tags,
				CreatedAt:   time.Now(),
				Version:     "1.0.0",
			},
		}

		// Write to file
		data, err := json.MarshalIndent(shareableConfig, "", "  ")
		if err != nil {
			fmt.Printf("❌ Error marshaling config: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			fmt.Printf("❌ Error writing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Configuration exported to: %s\n", outputPath)
		fmt.Printf("📋 Others can import with: dotfiles clone %s\n", outputPath)
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone <source>",
	Short: "Clone a shared configuration",
	Long:  `Import a shared configuration from:
  - GitHub Gist URL: https://gist.github.com/user/id
  - Local file: /path/to/config.json
  - Built-in template: template:web-dev`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		merge, _ := cmd.Flags().GetBool("merge")
		preview, _ := cmd.Flags().GetBool("preview")

		// Handle template:name format
		if strings.HasPrefix(source, "template:") {
			templateName := strings.TrimPrefix(source, "template:")
			if err := handleTemplateClone(templateName, merge); err != nil {
				fmt.Printf("❌ %v\n", err)
				os.Exit(1)
			}
			return
		}

		var shareableConfig ShareableConfig
		var err error

		if strings.HasPrefix(source, "http") {
			// Try web app first, then GitHub Gist
			if strings.Contains(source, "your-web-app.com") {
				shareableConfig, err = downloadFromWebApp(source)
			} else {
				// Handle GitHub Gist URL
				shareableConfig, err = downloadFromGist(source)
			}
		} else {
			// Handle local file
			shareableConfig, err = loadFromFile(source)
		}

		if err != nil {
			fmt.Printf("❌ Error loading shared config: %v\n", err)
			os.Exit(1)
		}

		// Show preview
		fmt.Printf("📋 Config: %s\n", shareableConfig.Metadata.Name)
		fmt.Printf("👤 Author: %s\n", shareableConfig.Metadata.Author)
		fmt.Printf("📝 Description: %s\n", shareableConfig.Metadata.Description)
		if len(shareableConfig.Metadata.Tags) > 0 {
			fmt.Printf("🏷️  Tags: %s\n", strings.Join(shareableConfig.Metadata.Tags, ", "))
		}
		fmt.Printf("📅 Created: %s\n", shareableConfig.Metadata.CreatedAt.Format("2006-01-02"))
		fmt.Println()

		fmt.Printf("📦 Packages included:\n")
		if len(shareableConfig.Taps) > 0 {
			fmt.Printf("  📋 Taps: %d\n", len(shareableConfig.Taps))
		}
		if len(shareableConfig.Brews) > 0 {
			fmt.Printf("  🍺 Brews: %d\n", len(shareableConfig.Brews))
		}
		if len(shareableConfig.Casks) > 0 {
			fmt.Printf("  📦 Casks: %d\n", len(shareableConfig.Casks))
		}
		if len(shareableConfig.Stow) > 0 {
			fmt.Printf("  🔗 Stow: %d\n", len(shareableConfig.Stow))
		}
		fmt.Println()

		if preview {
			fmt.Println("📋 Full package list:")
			if len(shareableConfig.Taps) > 0 {
				fmt.Println("Taps:", strings.Join(shareableConfig.Taps, ", "))
			}
			if len(shareableConfig.Brews) > 0 {
				fmt.Println("Brews:", strings.Join(shareableConfig.Brews, ", "))
			}
			if len(shareableConfig.Casks) > 0 {
				fmt.Println("Casks:", strings.Join(shareableConfig.Casks, ", "))
			}
			if len(shareableConfig.Stow) > 0 {
				fmt.Println("Stow:", strings.Join(shareableConfig.Stow, ", "))
			}
			return
		}

		if !askConfirmation("Import this configuration? (y/N): ", false) {
			fmt.Println("❌ Import cancelled.")
			return
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Error getting home directory: %v\n", err)
			os.Exit(1)
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
			existingConfig.Taps = mergeSlices(existingConfig.Taps, shareableConfig.Taps)
			existingConfig.Brews = mergeSlices(existingConfig.Brews, shareableConfig.Brews)
			existingConfig.Casks = mergeSlices(existingConfig.Casks, shareableConfig.Casks)
			existingConfig.Stow = mergeSlices(existingConfig.Stow, shareableConfig.Stow)

			if err := existingConfig.Save(configPath); err != nil {
				fmt.Printf("❌ Error saving merged config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Configuration merged successfully!")
		} else {
			// Replace existing config
			newConfig := &config.Config{
				Taps:  shareableConfig.Taps,
				Brews: shareableConfig.Brews,
				Casks: shareableConfig.Casks,
				Stow:  shareableConfig.Stow,
			}

			if err := newConfig.Save(configPath); err != nil {
				fmt.Printf("❌ Error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Configuration imported successfully!")
		}

		fmt.Println("💡 Next steps:")
		fmt.Println("  dotfiles status    # Check what needs to be installed")
		fmt.Println("  dotfiles install   # Install all packages")
	},
}

func uploadToGist(config ShareableConfig, public bool) (string, error) {
	// Convert config to JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	// Create Gist request
	gistReq := GistRequest{
		Description: fmt.Sprintf("Dotfiles Config: %s", config.Metadata.Name),
		Public:      public,
		Files: map[string]GistFile{
			"dotfiles-config.json": {
				Content: string(configJSON),
			},
		},
	}

	reqBody, err := json.Marshal(gistReq)
	if err != nil {
		return "", err
	}

	// Make request to GitHub API
	req, err := http.NewRequest("POST", "https://api.github.com/gists", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dotfiles-manager")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	var gistResp GistResponse
	if err := json.NewDecoder(resp.Body).Decode(&gistResp); err != nil {
		return "", err
	}

	return gistResp.HTMLURL, nil
}

func downloadFromGist(gistURL string) (ShareableConfig, error) {
	// Extract Gist ID from URL
	parts := strings.Split(gistURL, "/")
	gistID := parts[len(parts)-1]

	// Remove any hash fragments
	if idx := strings.Index(gistID, "#"); idx != -1 {
		gistID = gistID[:idx]
	}

	// Download from GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/gists/%s", gistID)
	resp, err := http.Get(apiURL)
	if err != nil {
		return ShareableConfig{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ShareableConfig{}, fmt.Errorf("failed to download gist: %s", resp.Status)
	}

	var gistResp GistResponse
	if err := json.NewDecoder(resp.Body).Decode(&gistResp); err != nil {
		return ShareableConfig{}, err
	}

	// Find the config file
	var configContent string
	for filename, file := range gistResp.Files {
		if strings.Contains(filename, "dotfiles-config") || strings.HasSuffix(filename, ".json") {
			configContent = file.Content
			break
		}
	}

	if configContent == "" {
		return ShareableConfig{}, fmt.Errorf("no dotfiles config found in gist")
	}

	var shareableConfig ShareableConfig
	if err := json.Unmarshal([]byte(configContent), &shareableConfig); err != nil {
		return ShareableConfig{}, err
	}

	return shareableConfig, nil
}

func loadFromFile(filePath string) (ShareableConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ShareableConfig{}, err
	}

	var shareableConfig ShareableConfig
	if err := json.Unmarshal(data, &shareableConfig); err != nil {
		return ShareableConfig{}, err
	}

	return shareableConfig, nil
}

func uploadToWebApp(config ShareableConfig, public bool) (string, error) {
	// Get API endpoint (with hardcoded default)
	apiEndpoint := os.Getenv("DOTFILES_API_ENDPOINT")
	if apiEndpoint == "" {
		apiEndpoint = "https://new-dotfiles-production.up.railway.app/api"
	}

	// Convert config to JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	// Create upload request
	uploadReq := map[string]interface{}{
		"name":        config.Metadata.Name,
		"description": config.Metadata.Description,
		"author":      config.Metadata.Author,
		"tags":        config.Metadata.Tags,
		"config":      string(configJSON),
		"public":      public,
	}

	reqBody, err := json.Marshal(uploadReq)
	if err != nil {
		return "", err
	}

	// Make request to web app API
	url := fmt.Sprintf("%s/configs/upload", apiEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dotfiles-manager")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("web app API error: %s - %s", resp.Status, string(body))
	}

	var uploadResp struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", err
	}

	return uploadResp.URL, nil
}

func downloadFromWebApp(webAppURL string) (ShareableConfig, error) {
	// Extract config ID from URL or use URL directly as API endpoint
	// Assuming URL format: https://your-web-app.com/config/123
	apiURL := strings.Replace(webAppURL, "/config/", "/api/configs/", 1)

	resp, err := http.Get(apiURL)
	if err != nil {
		return ShareableConfig{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ShareableConfig{}, fmt.Errorf("failed to download from web app: %s", resp.Status)
	}

	var shareableConfig ShareableConfig
	if err := json.NewDecoder(resp.Body).Decode(&shareableConfig); err != nil {
		return ShareableConfig{}, err
	}

	return shareableConfig, nil
}

// Note: askConfirmation function is already defined in onboard.go

func init() {
	// Share gist flags
	shareGistCmd.Flags().StringP("name", "n", "", "Name for the shared config (required)")
	shareGistCmd.Flags().StringP("description", "d", "", "Description of the config")
	shareGistCmd.Flags().StringP("author", "a", "", "Author name")
	shareGistCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags for categorization (e.g., web-dev,mobile)")
	shareGistCmd.Flags().Bool("private", false, "Create private gist")

	// Share file flags
	shareFileCmd.Flags().StringP("name", "n", "", "Name for the shared config (required)")
	shareFileCmd.Flags().StringP("description", "d", "", "Description of the config")
	shareFileCmd.Flags().StringP("author", "a", "", "Author name")
	shareFileCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags for categorization")

	// Clone flags
	cloneCmd.Flags().Bool("merge", false, "Merge with existing config instead of replacing")
	cloneCmd.Flags().Bool("preview", false, "Preview config without importing")

	shareCmd.AddCommand(shareGistCmd)
	shareCmd.AddCommand(shareFileCmd)
	rootCmd.AddCommand(shareCmd)
	rootCmd.AddCommand(cloneCmd)
}
