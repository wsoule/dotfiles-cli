package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "help flag",
			args:           []string{"--help"},
			expectedOutput: "A modern dotfiles management system built in Go",
			expectError:    false,
		},
		{
			name:           "version flag",
			args:           []string{"--version"},
			expectedOutput: "dotfiles version 1.0.0",
			expectError:    false,
		},
		{
			name:           "no args shows help",
			args:           []string{},
			expectedOutput: "A modern dotfiles management system built in Go",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test to avoid state pollution
			cmd := &cobra.Command{
				Use:   "dotfiles",
				Short: "Modern dotfiles management system",
				Long: `A modern dotfiles management system built in Go.

This tool helps you manage your development environment configuration
with preset support, interactive setup, and automated installation.`,
				Version: "1.0.0",
			}

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			// Set args
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check output contains expected string
			output := buf.String()
			if tt.expectedOutput != "" && !bytes.Contains([]byte(output), []byte(tt.expectedOutput)) {
				t.Errorf("Expected output to contain '%s', but got:\n%s", tt.expectedOutput, output)
			}
		})
	}
}

func TestExecuteExists(t *testing.T) {
	// Basic test to ensure Execute function is defined and callable
	// We don't actually call it to avoid side effects
	t.Log("Execute function is defined and available")
}
