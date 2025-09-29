package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type GitHubSearchResponse struct {
	TotalCount        int           `json:"total_count"`
	IncompleteResults bool          `json:"incomplete_results"`
	Items             []GitHubGist  `json:"items"`
}

type GitHubGist struct {
	ID          string            `json:"id"`
	HTMLURL     string            `json:"html_url"`
	Description string            `json:"description"`
	Public      bool              `json:"public"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Files       map[string]GistFile `json:"files"`
	Owner       GistOwner         `json:"owner"`
}

type GistOwner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type FeaturedConfig struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
}

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover shared configurations from the community",
	Long:  `Find and browse configurations shared by other developers`,
}

var discoverSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for shared configurations",
	Long:  `Search GitHub Gists for dotfiles configurations`,
	Run: func(cmd *cobra.Command, args []string) {
		query := "dotfiles-config.json"
		if len(args) > 0 {
			query = strings.Join(args, " ") + " dotfiles-config.json"
		}

		tags, _ := cmd.Flags().GetStringSlice("tags")
		if len(tags) > 0 {
			query += " " + strings.Join(tags, " ")
		}

		fmt.Printf("🔍 Searching for configurations: %s\n", query)
		fmt.Println()

		gists, err := searchGists(query)
		if err != nil {
			fmt.Printf("❌ Error searching: %v\n", err)
			return
		}

		if len(gists.Items) == 0 {
			fmt.Println("📭 No configurations found matching your criteria")
			fmt.Println()
			fmt.Println("💡 Try:")
			fmt.Println("  dotfiles discover search web-dev")
			fmt.Println("  dotfiles discover search --tags=python,data")
			fmt.Println("  dotfiles templates list  # Browse built-in templates")
			return
		}

		fmt.Printf("📋 Found %d configurations:\n", len(gists.Items))
		fmt.Println("=" + strings.Repeat("=", 30))
		fmt.Println()

		for i, gist := range gists.Items {
			if i >= 10 { // Limit to first 10 results
				break
			}

			fmt.Printf("%d. 📝 %s\n", i+1, gist.Description)
			fmt.Printf("   👤 Author: %s\n", gist.Owner.Login)
			fmt.Printf("   📅 Updated: %s\n", gist.UpdatedAt.Format("2006-01-02"))
			fmt.Printf("   🔗 URL: %s\n", gist.HTMLURL)
			fmt.Printf("   📦 Clone: dotfiles clone %s\n", gist.HTMLURL)
			fmt.Println()
		}

		if len(gists.Items) > 10 {
			fmt.Printf("... and %d more results\n", len(gists.Items)-10)
			fmt.Println("🔍 Use more specific search terms to narrow results")
		}
	},
}

var discoverFeaturedCmd = &cobra.Command{
	Use:   "featured",
	Short: "Show featured community configurations",
	Long:  `Display curated, high-quality configurations from the community`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⭐ Featured Community Configurations")
		fmt.Println("=" + strings.Repeat("=", 37))
		fmt.Println()

		// Try to fetch from web app first
		featuredConfigs, err := fetchFeaturedConfigs()
		if err != nil {
			fmt.Printf("⚠️  Could not fetch from web app, showing built-in examples: %v\n", err)
			fmt.Println()

			// Fallback to hardcoded examples
			featuredConfigs = []FeaturedConfig{
				{
					Name:        "Full Stack Web Developer",
					Description: "Complete setup for modern web development with React, Node.js, and Docker",
					Author:      "community",
					URL:         "https://gist.github.com/example/webdev-config",
					Tags:        []string{"web-dev", "react", "nodejs", "docker"},
				},
				{
					Name:        "Python Data Scientist",
					Description: "Comprehensive Python environment for data science and machine learning",
					Author:      "community",
					URL:         "https://gist.github.com/example/datascience-config",
					Tags:        []string{"python", "data-science", "jupyter", "ml"},
				},
				{
					Name:        "DevOps Engineer",
					Description: "Infrastructure and automation tools for DevOps workflows",
					Author:      "community",
					URL:         "https://gist.github.com/example/devops-config",
					Tags:        []string{"devops", "kubernetes", "terraform", "monitoring"},
				},
			}
		}

		for i, config := range featuredConfigs {
			fmt.Printf("%d. 📝 %s\n", i+1, config.Name)
			fmt.Printf("   📄 %s\n", config.Description)
			fmt.Printf("   👤 Author: %s\n", config.Author)
			fmt.Printf("   🏷️  Tags: %s\n", strings.Join(config.Tags, ", "))
			fmt.Printf("   📦 Clone: dotfiles clone %s\n", config.URL)
			fmt.Println()
		}

		fmt.Println("💡 Community Tips:")
		fmt.Println("  • Always preview configs before applying: --preview")
		fmt.Println("  • Merge with your existing setup: --merge")
		fmt.Println("  • Share your own config: dotfiles share gist")
	},
}

var discoverStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show community sharing statistics",
	Long:  `Display statistics about configuration sharing in the community`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 Community Sharing Statistics")
		fmt.Println("=" + strings.Repeat("=", 32))
		fmt.Println()

		// Search for dotfiles configurations
		gists, err := searchGists("dotfiles-config.json")
		if err != nil {
			fmt.Printf("❌ Error fetching stats: %v\n", err)
			return
		}

		fmt.Printf("📋 Total shared configurations: %d\n", gists.TotalCount)
		fmt.Println()

		fmt.Println("🏆 Popular Configuration Types:")
		fmt.Println("  1. Web Development")
		fmt.Println("  2. Data Science")
		fmt.Println("  3. DevOps & Cloud")
		fmt.Println("  4. Mobile Development")
		fmt.Println("  5. Minimal Setups")
		fmt.Println()

		fmt.Println("💡 Get involved:")
		fmt.Println("  dotfiles share gist --name='My Config'  # Share your setup")
		fmt.Println("  dotfiles discover search <topic>       # Find configs")
		fmt.Println("  dotfiles templates list                # Browse templates")
	},
}

const DefaultAPIEndpoint = "https://your-web-app.com/api"

func searchGists(query string) (*GitHubSearchResponse, error) {
	// Use your web app's API first, fallback to GitHub search
	webAppURL := fmt.Sprintf("%s/configs/search?q=%s", getAPIEndpoint(), query)

	// Try web app API first
	if configs, err := searchWebApp(webAppURL); err == nil {
		return configs, nil
	}

	// Fallback to GitHub search
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s+in:file+filename:dotfiles-config", query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "dotfiles-manager")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var searchResp GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func getAPIEndpoint() string {
	if endpoint := os.Getenv("DOTFILES_API_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return DefaultAPIEndpoint
}

func searchWebApp(url string) (*GitHubSearchResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "dotfiles-manager")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("web app API error: %s", resp.Status)
	}

	var searchResp GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func fetchFeaturedConfigs() ([]FeaturedConfig, error) {
	url := fmt.Sprintf("%s/configs/featured", getAPIEndpoint())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "dotfiles-manager")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("web app API error: %s", resp.Status)
	}

	var configs []FeaturedConfig
	if err := json.NewDecoder(resp.Body).Decode(&configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func init() {
	discoverSearchCmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags (e.g., web-dev,python)")

	discoverCmd.AddCommand(discoverSearchCmd)
	discoverCmd.AddCommand(discoverFeaturedCmd)
	discoverCmd.AddCommand(discoverStatsCmd)
	rootCmd.AddCommand(discoverCmd)
}