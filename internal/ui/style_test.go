package ui

import (
	"strings"
	"testing"
)

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner("test message")

	if spinner == nil {
		t.Error("NewSpinner returned nil")
	}

	// Test that spinner has the expected message
	if !strings.Contains(spinner.Suffix, "test message") {
		t.Errorf("Expected spinner suffix to contain 'test message', got: %s", spinner.Suffix)
	}
}

func TestProgressBar(t *testing.T) {
	pb, err := ProgressBar("test progress", 100)
	if err != nil {
		t.Errorf("ProgressBar returned error: %v", err)
	}
	if pb == nil {
		t.Error("ProgressBar returned nil")
	}

	// Stop the progress bar to clean up
	if pb != nil {
		pb.Stop()
	}
}

func TestTable(t *testing.T) {
	// Test that Table function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Table() panicked: %v", r)
		}
	}()

	headers := []string{"Name", "Value"}
	rows := [][]string{
		{"test1", "value1"},
		{"test2", "value2"},
	}

	// Just test that the function executes without panicking
	Table(headers, rows)
}

func TestPrintFunctions(t *testing.T) {
	// Test that print functions don't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Print functions panicked: %v", r)
		}
	}()

	PrintSuccess("test message")
	PrintError("error message")
	PrintWarning("warning message")
	PrintInfo("info message")
	PrintStep(1, 3, "test step")
	PrintSection("Test Section")
}

// Benchmark tests for performance
func BenchmarkNewSpinner(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spinner := NewSpinner("benchmark test")
		_ = spinner
	}
}