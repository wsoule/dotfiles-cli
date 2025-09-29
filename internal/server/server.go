package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"dotfiles/internal/config"
	"dotfiles/internal/ui"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the embedded HTTP server for the web UI
type Server struct {
	port   int
	config *config.Config
	server *http.Server
	done   chan bool
}

// NewServer creates a new server instance
func NewServer(port int) *Server {
	return &Server{
		port: port,
		done: make(chan bool),
	}
}

// Start starts the HTTP server and opens the browser
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	// API endpoints
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/save", s.handleSave)
	mux.HandleFunc("/api/presets", s.handlePresets)
	mux.HandleFunc("/api/validate", s.handleValidate)

	// Main page
	mux.HandleFunc("/", s.handleIndex)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(100 * time.Millisecond)

	// Open browser
	url := fmt.Sprintf("http://localhost:%d", s.port)
	ui.PrintInfo(fmt.Sprintf("Opening setup wizard at: %s", url))

	if err := openBrowser(url); err != nil {
		ui.PrintWarning("Could not open browser automatically")
		fmt.Printf("Please open: %s\n", url)
	}

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// WaitForCompletion waits for the setup to complete
func (s *Server) WaitForCompletion() {
	<-s.done
}

// handleIndex serves the main setup page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Serve the static index.html file
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

// handleConfig serves the current configuration
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.config != nil {
		json.NewEncoder(w).Encode(s.config)
	} else {
		// Return empty config
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// handleSave saves the configuration
func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set metadata
	cfg.Metadata.Version = "1.0.0"
	cfg.Metadata.CreatedAt = time.Now()
	cfg.Metadata.LastModified = time.Now()
	cfg.Metadata.CreatedBy = "web-wizard"

	// Save configuration with proper path
	home, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "Failed to get home directory", http.StatusInternalServerError)
		return
	}

	configPath := filepath.Join(home, ".dotfiles", "config.json")
	configManager := config.NewManager(configPath)
	configManager.Set(&cfg)

	if err := configManager.Save(); err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	// Signal completion
	go func() {
		time.Sleep(1 * time.Second)
		s.done <- true
	}()
}

// handlePresets serves available presets
func (s *Server) handlePresets(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement preset loading
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{})
}

// handleValidate validates configuration
func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement validation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"valid": true})
}

// openBrowser opens the default browser to the given URL
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}