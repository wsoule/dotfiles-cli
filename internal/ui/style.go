package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/pterm/pterm"
)

var (
	// Color definitions
	Primary   = color.New(color.FgCyan, color.Bold)
	Success   = color.New(color.FgGreen, color.Bold)
	Warning   = color.New(color.FgYellow, color.Bold)
	Error     = color.New(color.FgRed, color.Bold)
	Info      = color.New(color.FgBlue)
	Muted     = color.New(color.FgHiBlack)
	Highlight = color.New(color.FgMagenta, color.Bold)
)

// Banner displays the application banner
func Banner() {
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Sprint("🛠  DOTFILES MANAGER"))

	fmt.Println()
	Primary.Println("  A modern dotfiles management system built in Go")
	Muted.Println("  Configure your development environment with ease")
	fmt.Println()
}

// Success prints a success message with checkmark
func PrintSuccess(message string) {
	Success.Print("✓ ")
	fmt.Println(message)
}

// Error prints an error message with X mark
func PrintError(message string) {
	Error.Print("✗ ")
	fmt.Println(message)
}

// Warning prints a warning message with warning symbol
func PrintWarning(message string) {
	Warning.Print("⚠ ")
	fmt.Println(message)
}

// Info prints an info message with info symbol
func PrintInfo(message string) {
	Info.Print("ℹ ")
	fmt.Println(message)
}

// Step prints a step in a process
func PrintStep(step int, total int, message string) {
	Primary.Printf("[%d/%d] ", step, total)
	fmt.Println(message)
}

// Section prints a section header
func PrintSection(title string) {
	fmt.Println()
	Highlight.Println("▶ " + title)
	Muted.Println("  " + strings.Repeat("─", len(title)+2))
}

// NewSpinner creates a new spinner with custom styling
func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Color("cyan")
	return s
}

// Table creates a pretty table for displaying data
func Table(headers []string, rows [][]string) {
	tableData := pterm.TableData{headers}
	for _, row := range rows {
		tableData = append(tableData, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

// ProgressBar creates a progress bar
func ProgressBar(title string, total int) (*pterm.ProgressbarPrinter, error) {
	return pterm.DefaultProgressbar.WithTitle(title).WithTotal(total).Start()
}

// Confirm prompts user for yes/no confirmation
func Confirm(message string) bool {
	result, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(false).WithDefaultText(message).Show()
	return result
}

// Select prompts user to select from options
func Select(message string, options []string) (string, error) {
	return pterm.DefaultInteractiveSelect.WithDefaultText(message).WithOptions(options).Show()
}

// MultiSelect prompts user to select multiple options
func MultiSelect(message string, options []string) ([]string, error) {
	return pterm.DefaultInteractiveMultiselect.WithDefaultText(message).WithOptions(options).Show()
}

// Input prompts user for text input
func Input(message string, defaultValue string) string {
	result, _ := pterm.DefaultInteractiveTextInput.WithDefaultText(message).WithDefaultValue(defaultValue).Show()
	return result
}
