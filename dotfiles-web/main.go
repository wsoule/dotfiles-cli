package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Config struct {
	Brews []string `json:"brews"`
	Casks []string `json:"casks"`
	Taps  []string `json:"taps"`
	Stow  []string `json:"stow"`
}

type ShareableConfig struct {
	Config   `json:",inline"`
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

type StoredConfig struct {
	ID           string          `json:"id"`
	Config       ShareableConfig `json:"config"`
	Public       bool            `json:"public"`
	CreatedAt    time.Time       `json:"created_at"`
	DownloadCount int            `json:"download_count"`
}

// In-memory storage (use database in production)
var (
	configs = make(map[string]*StoredConfig)
	mu      sync.RWMutex
)

func main() {
	// Get port from environment (Railway sets PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Dotfiles Config Sharing API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"upload":   "POST /api/configs/upload",
				"download": "GET /api/configs/:id",
				"search":   "GET /api/configs/search",
				"featured": "GET /api/configs/featured",
			},
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Upload config
		api.POST("/configs/upload", uploadConfig)

		// Get config by ID
		api.GET("/configs/:id", getConfig)

		// Search configs
		api.GET("/configs/search", searchConfigs)

		// Get featured configs
		api.GET("/configs/featured", getFeaturedConfigs)

		// Get stats
		api.GET("/configs/stats", getStats)
	}

	// Web interface routes
	r.GET("/config/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(302, "/api/configs/"+id)
	})

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func uploadConfig(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Author      string   `json:"author"`
		Tags        []string `json:"tags"`
		Config      string   `json:"config" binding:"required"`
		Public      bool     `json:"public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Parse the config JSON
	var shareableConfig ShareableConfig
	if err := json.Unmarshal([]byte(req.Config), &shareableConfig); err != nil {
		c.JSON(400, gin.H{"error": "Invalid config JSON", "details": err.Error()})
		return
	}

	// Create stored config
	id := uuid.New().String()
	stored := &StoredConfig{
		ID:           id,
		Config:       shareableConfig,
		Public:       req.Public,
		CreatedAt:    time.Now(),
		DownloadCount: 0,
	}

	// Store in memory
	mu.Lock()
	configs[id] = stored
	mu.Unlock()

	log.Printf("Config uploaded: %s (%s)", req.Name, id)

	c.JSON(201, gin.H{
		"id":  id,
		"url": fmt.Sprintf("/config/%s", id),
	})
}

func getConfig(c *gin.Context) {
	id := c.Param("id")

	mu.RLock()
	stored, exists := configs[id]
	mu.RUnlock()

	if !exists {
		c.JSON(404, gin.H{"error": "Config not found"})
		return
	}

	// Increment download count
	mu.Lock()
	stored.DownloadCount++
	mu.Unlock()

	c.JSON(200, stored.Config)
}

func searchConfigs(c *gin.Context) {
	query := strings.ToLower(c.Query("q"))

	mu.RLock()
	var results []gin.H
	for id, stored := range configs {
		if !stored.Public {
			continue
		}

		// Simple search in name, description, tags
		searchText := strings.ToLower(stored.Config.Metadata.Name + " " +
			stored.Config.Metadata.Description + " " +
			strings.Join(stored.Config.Metadata.Tags, " "))

		if query == "" || strings.Contains(searchText, query) {
			results = append(results, gin.H{
				"id":          id,
				"html_url":    fmt.Sprintf("/config/%s", id),
				"description": fmt.Sprintf("Dotfiles Config: %s", stored.Config.Metadata.Name),
				"public":      stored.Public,
				"created_at":  stored.CreatedAt,
				"updated_at":  stored.CreatedAt,
				"files": map[string]interface{}{
					"dotfiles-config.json": gin.H{
						"content": "config content",
					},
				},
				"owner": gin.H{
					"login":      stored.Config.Metadata.Author,
					"avatar_url": "",
				},
			})
		}
	}
	mu.RUnlock()

	c.JSON(200, gin.H{
		"total_count":         len(results),
		"incomplete_results": false,
		"items":              results,
	})
}

func getFeaturedConfigs(c *gin.Context) {
	// Get most downloaded public configs
	mu.RLock()
	var featured []gin.H
	for id, stored := range configs {
		if !stored.Public {
			continue
		}

		featured = append(featured, gin.H{
			"name":        stored.Config.Metadata.Name,
			"description": stored.Config.Metadata.Description,
			"author":      stored.Config.Metadata.Author,
			"url":         fmt.Sprintf("/config/%s", id),
			"tags":        stored.Config.Metadata.Tags,
			"downloads":   stored.DownloadCount,
		})
	}
	mu.RUnlock()

	// Sort by download count (simple bubble sort for demo)
	for i := 0; i < len(featured); i++ {
		for j := i + 1; j < len(featured); j++ {
			if featured[i]["downloads"].(int) < featured[j]["downloads"].(int) {
				featured[i], featured[j] = featured[j], featured[i]
			}
		}
	}

	// Limit to top 10
	if len(featured) > 10 {
		featured = featured[:10]
	}

	c.JSON(200, featured)
}

func getStats(c *gin.Context) {
	mu.RLock()
	total := len(configs)
	public := 0
	totalDownloads := 0

	for _, stored := range configs {
		if stored.Public {
			public++
		}
		totalDownloads += stored.DownloadCount
	}
	mu.RUnlock()

	c.JSON(200, gin.H{
		"total_configs":    total,
		"public_configs":   public,
		"total_downloads":  totalDownloads,
	})
}
