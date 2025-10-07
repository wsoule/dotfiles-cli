package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"dotfiles/internal/config"
)

// Snapshot represents a configuration snapshot
type Snapshot struct {
	Timestamp   string            `json:"timestamp"`
	Description string            `json:"description"`
	Config      *config.Config    `json:"config"`
	Metadata    map[string]string `json:"metadata"`
}

// CreateAutoSnapshot creates an automatic snapshot before major operations
func CreateAutoSnapshot(description string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, ".dotfiles", "config.json")
	cfg, err := config.Load(configPath)
	if err != nil {
		return "", err
	}

	snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")
	if err := os.MkdirAll(snapshotsDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102-150405")
	snapshot := Snapshot{
		Timestamp:   timestamp,
		Description: description,
		Config:      cfg,
		Metadata: map[string]string{
			"created_by": "auto-snapshot",
			"platform":   runtime.GOOS,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}

	snapshotPath := filepath.Join(snapshotsDir, timestamp+".json")
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(snapshotPath, data, 0644); err != nil {
		return "", err
	}

	return timestamp, nil
}

// ListSnapshots returns all available snapshots
func ListSnapshots() ([]Snapshot, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")
	if _, err := os.Stat(snapshotsDir); os.IsNotExist(err) {
		return []Snapshot{}, nil
	}

	entries, err := os.ReadDir(snapshotsDir)
	if err != nil {
		return nil, err
	}

	snapshots := []Snapshot{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			snapshotPath := filepath.Join(snapshotsDir, entry.Name())
			data, err := os.ReadFile(snapshotPath)
			if err != nil {
				continue
			}

			var snapshot Snapshot
			if err := json.Unmarshal(data, &snapshot); err != nil {
				continue
			}

			snapshots = append(snapshots, snapshot)
		}
	}

	return snapshots, nil
}

// RestoreSnapshot restores a snapshot by timestamp
func RestoreSnapshot(timestamp string, createBackup bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Load snapshot
	snapshotPath := filepath.Join(home, ".dotfiles", "snapshots", timestamp+".json")
	data, err := os.ReadFile(snapshotPath)
	if err != nil {
		return fmt.Errorf("snapshot not found: %s", timestamp)
	}

	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("error reading snapshot: %v", err)
	}

	configPath := filepath.Join(home, ".dotfiles", "config.json")

	// Create backup if requested
	if createBackup {
		_, err := CreateAutoSnapshot("Auto-backup before restore from " + timestamp)
		if err != nil {
			return fmt.Errorf("failed to create backup: %v", err)
		}
	}

	// Restore snapshot
	if err := snapshot.Config.Save(configPath); err != nil {
		return fmt.Errorf("error restoring configuration: %v", err)
	}

	return nil
}

// CleanOldSnapshots removes snapshots older than specified days
func CleanOldSnapshots(daysToKeep int) (int, error) {
	snapshots, err := ListSnapshots()
	if err != nil {
		return 0, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return 0, err
	}

	snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)
	removed := 0

	for _, snapshot := range snapshots {
		t, err := time.Parse("20060102-150405", snapshot.Timestamp)
		if err != nil {
			continue
		}

		if t.Before(cutoffTime) {
			snapshotPath := filepath.Join(snapshotsDir, snapshot.Timestamp+".json")
			if err := os.Remove(snapshotPath); err == nil {
				removed++
			}
		}
	}

	return removed, nil
}
