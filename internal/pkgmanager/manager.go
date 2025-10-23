// Package pkgmanager provides a unified interface for managing packages across different operating systems.
// It abstracts package manager operations (install, list, check) for Homebrew (macOS),
// pacman (Arch Linux), apt (Debian/Ubuntu), and yum/dnf (RHEL/Fedora).
// This enables write-once package management code that works across all supported platforms.
package pkgmanager

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// PackageManager defines the interface for package manager operations.
// All platform-specific package managers (Homebrew, pacman, apt, yum) implement this interface,
// allowing for unified package management across different operating systems.
type PackageManager interface {
	// Install installs one or more packages of the specified type.
	// packageType can be "brew" (CLI tools), "cask" (GUI apps), or "tap" (repositories).
	// On Linux, "cask" and "tap" types are treated as regular packages.
	Install(packages []string, packageType string) error

	// IsInstalled checks if a specific package is installed on the system.
	// Returns true if the package is installed, false otherwise.
	IsInstalled(pkg string, packageType string) (bool, error)

	// ListInstalled returns all packages of the given type that are currently installed.
	// For Homebrew: packageType can be "brew", "cask", or "tap".
	// For Linux package managers: only "brew" type is meaningful (returns all packages).
	ListInstalled(packageType string) ([]string, error)

	// GetName returns the human-readable name of the package manager.
	// Examples: "homebrew", "pacman", "apt", "yum"
	GetName() string

	// IsAvailable checks if the package manager is installed and accessible on the system.
	// Returns true if the package manager command is found in PATH.
	IsAvailable() bool

	// GenerateInstallFile creates a package list file in the package manager's format.
	// For Homebrew: generates a Brewfile with tap/brew/cask entries.
	// For Linux: generates a plain text list of package names (one per line).
	GenerateInstallFile(brews, casks, taps []string) (string, error)

	// InstallFromFile installs all packages listed in a package list file.
	// For Homebrew: uses `brew bundle` with a Brewfile.
	// For Linux: reads package names from file and installs them.
	InstallFromFile(filePath string) error
}

// GetPackageManager returns the appropriate package manager for the current OS.
// It auto-detects the platform and available package managers:
//   - macOS: Returns HomebrewManager
//   - Linux: Returns PacmanManager, AptManager, or YumManager (first available)
//
// Returns an error if the OS is unsupported or no package manager is found on Linux.
func GetPackageManager() (PackageManager, error) {
	switch runtime.GOOS {
	case "darwin":
		return &HomebrewManager{}, nil
	case "linux":
		// Try to detect which Linux package manager is available
		if commandExists("pacman") {
			return &PacmanManager{}, nil
		}
		if commandExists("apt-get") {
			return &AptManager{}, nil
		}
		if commandExists("yum") {
			return &YumManager{}, nil
		}
		return nil, fmt.Errorf("no supported package manager found on Linux")
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// commandExists checks if a command exists in the system PATH.
// This is used to detect which package manager is available on the system.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// HomebrewManager implements PackageManager for Homebrew on macOS.
// Homebrew is the de facto standard package manager for macOS, supporting both
// CLI tools (brews) and GUI applications (casks).
type HomebrewManager struct{}

func (h *HomebrewManager) GetName() string {
	return "homebrew"
}

func (h *HomebrewManager) IsAvailable() bool {
	return commandExists("brew")
}

func (h *HomebrewManager) Install(packages []string, packageType string) error {
	if len(packages) == 0 {
		return nil
	}

	var args []string
	if packageType == "cask" {
		args = append([]string{"install", "--cask"}, packages...)
	} else if packageType == "tap" {
		for _, pkg := range packages {
			cmd := exec.Command("brew", "tap", pkg)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to tap %s: %v", pkg, err)
			}
		}
		return nil
	} else {
		args = append([]string{"install"}, packages...)
	}

	cmd := exec.Command("brew", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (h *HomebrewManager) IsInstalled(pkg string, packageType string) (bool, error) {
	if packageType == "tap" {
		cmd := exec.Command("brew", "tap")
		output, err := cmd.Output()
		if err != nil {
			return false, err
		}
		taps := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, tap := range taps {
			if strings.TrimSpace(tap) == pkg {
				return true, nil
			}
		}
		return false, nil
	}

	var cmd *exec.Cmd
	if packageType == "cask" {
		cmd = exec.Command("brew", "list", "--cask", pkg)
	} else {
		cmd = exec.Command("brew", "list", pkg)
	}

	err := cmd.Run()
	return err == nil, nil
}

func (h *HomebrewManager) ListInstalled(packageType string) ([]string, error) {
	var cmd *exec.Cmd
	if packageType == "cask" {
		cmd = exec.Command("brew", "list", "--cask")
	} else if packageType == "tap" {
		cmd = exec.Command("brew", "tap")
	} else {
		cmd = exec.Command("brew", "list", "--formula")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			packages = append(packages, line)
		}
	}

	return packages, nil
}

func (h *HomebrewManager) GenerateInstallFile(brews, casks, taps []string) (string, error) {
	var content string

	// Add taps
	for _, tap := range taps {
		content += "tap \"" + tap + "\"\n"
	}
	if len(taps) > 0 {
		content += "\n"
	}

	// Add brews
	for _, brew := range brews {
		content += "brew \"" + brew + "\"\n"
	}
	if len(brews) > 0 {
		content += "\n"
	}

	// Add casks
	for _, cask := range casks {
		content += "cask \"" + cask + "\"\n"
	}

	return content, nil
}

func (h *HomebrewManager) InstallFromFile(filePath string) error {
	cmd := exec.Command("brew", "bundle", "--file="+filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PacmanManager implements PackageManager for pacman (Arch Linux).
// Pacman is the package manager for Arch Linux and Arch-based distributions.
// If yay (AUR helper) is available, it will be used for broader package support.
type PacmanManager struct{}

func (p *PacmanManager) GetName() string {
	return "pacman"
}

func (p *PacmanManager) IsAvailable() bool {
	return commandExists("pacman")
}

func (p *PacmanManager) Install(packages []string, packageType string) error {
	if len(packages) == 0 {
		return nil
	}

	// On Arch, we ignore "cask" and "tap" - they're macOS concepts
	// All packages are treated the same
	if packageType == "cask" || packageType == "tap" {
		return nil
	}

	// Check if yay (AUR helper) is available, prefer it over pacman
	var cmd *exec.Cmd
	if commandExists("yay") {
		cmd = exec.Command("yay", append([]string{"-S", "--noconfirm"}, packages...)...)
	} else {
		cmd = exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm"}, packages...)...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PacmanManager) IsInstalled(pkg string, packageType string) (bool, error) {
	// Ignore macOS-specific package types
	if packageType == "cask" || packageType == "tap" {
		return false, nil
	}

	cmd := exec.Command("pacman", "-Q", pkg)
	err := cmd.Run()
	return err == nil, nil
}

func (p *PacmanManager) ListInstalled(packageType string) ([]string, error) {
	// Ignore macOS-specific package types
	if packageType == "cask" || packageType == "tap" {
		return []string{}, nil
	}

	cmd := exec.Command("pacman", "-Qq")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			packages = append(packages, line)
		}
	}

	return packages, nil
}

func (p *PacmanManager) GenerateInstallFile(brews, casks, taps []string) (string, error) {
	// Generate a simple package list for pacman
	// Ignore taps (macOS concept), combine brews and casks into one list
	var content string

	allPackages := append([]string{}, brews...)
	allPackages = append(allPackages, casks...)

	for _, pkg := range allPackages {
		content += pkg + "\n"
	}

	return content, nil
}

func (p *PacmanManager) InstallFromFile(filePath string) error {
	// Read package list and install
	content, err := exec.Command("cat", filePath).Output()
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			packages = append(packages, line)
		}
	}

	return p.Install(packages, "brew")
}

// AptManager implements PackageManager for apt (Debian/Ubuntu).
// APT is the package manager for Debian, Ubuntu, and their derivatives.
// This uses apt-get for installations and dpkg for package queries.
type AptManager struct{}

func (a *AptManager) GetName() string {
	return "apt"
}

func (a *AptManager) IsAvailable() bool {
	return commandExists("apt-get")
}

func (a *AptManager) Install(packages []string, packageType string) error {
	if len(packages) == 0 {
		return nil
	}

	// Ignore macOS-specific package types
	if packageType == "cask" || packageType == "tap" {
		return nil
	}

	cmd := exec.Command("sudo", append([]string{"apt-get", "install", "-y"}, packages...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *AptManager) IsInstalled(pkg string, packageType string) (bool, error) {
	if packageType == "cask" || packageType == "tap" {
		return false, nil
	}

	cmd := exec.Command("dpkg", "-l", pkg)
	err := cmd.Run()
	return err == nil, nil
}

func (a *AptManager) ListInstalled(packageType string) ([]string, error) {
	if packageType == "cask" || packageType == "tap" {
		return []string{}, nil
	}

	cmd := exec.Command("dpkg", "--get-selections")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var packages []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "install" {
			packages = append(packages, fields[0])
		}
	}

	return packages, nil
}

func (a *AptManager) GenerateInstallFile(brews, casks, taps []string) (string, error) {
	var content string

	allPackages := append([]string{}, brews...)
	allPackages = append(allPackages, casks...)

	for _, pkg := range allPackages {
		content += pkg + "\n"
	}

	return content, nil
}

func (a *AptManager) InstallFromFile(filePath string) error {
	content, err := exec.Command("cat", filePath).Output()
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			packages = append(packages, line)
		}
	}

	return a.Install(packages, "brew")
}

// YumManager implements PackageManager for yum/dnf (RHEL/Fedora).
// This supports both yum (older RHEL/CentOS) and dnf (newer Fedora/RHEL).
// If dnf is available, it will be used; otherwise falls back to yum.
type YumManager struct{}

func (y *YumManager) GetName() string {
	return "yum"
}

func (y *YumManager) IsAvailable() bool {
	return commandExists("yum") || commandExists("dnf")
}

func (y *YumManager) Install(packages []string, packageType string) error {
	if len(packages) == 0 {
		return nil
	}

	if packageType == "cask" || packageType == "tap" {
		return nil
	}

	cmdName := "yum"
	if commandExists("dnf") {
		cmdName = "dnf"
	}

	cmd := exec.Command("sudo", append([]string{cmdName, "install", "-y"}, packages...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (y *YumManager) IsInstalled(pkg string, packageType string) (bool, error) {
	if packageType == "cask" || packageType == "tap" {
		return false, nil
	}

	cmdName := "yum"
	if commandExists("dnf") {
		cmdName = "dnf"
	}

	cmd := exec.Command(cmdName, "list", "installed", pkg)
	err := cmd.Run()
	return err == nil, nil
}

func (y *YumManager) ListInstalled(packageType string) ([]string, error) {
	if packageType == "cask" || packageType == "tap" {
		return []string{}, nil
	}

	cmdName := "yum"
	if commandExists("dnf") {
		cmdName = "dnf"
	}

	cmd := exec.Command(cmdName, "list", "installed")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var packages []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && !strings.HasPrefix(line, "Installed") {
			packages = append(packages, fields[0])
		}
	}

	return packages, nil
}

func (y *YumManager) GenerateInstallFile(brews, casks, taps []string) (string, error) {
	var content string

	allPackages := append([]string{}, brews...)
	allPackages = append(allPackages, casks...)

	for _, pkg := range allPackages {
		content += pkg + "\n"
	}

	return content, nil
}

func (y *YumManager) InstallFromFile(filePath string) error {
	content, err := exec.Command("cat", filePath).Output()
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			packages = append(packages, line)
		}
	}

	return y.Install(packages, "brew")
}
