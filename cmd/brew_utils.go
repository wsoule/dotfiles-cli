package cmd

import (
	"dotfiles/internal/pkgmanager"
	"os/exec"
	"strings"
)

// getInstalledBrews returns a list of all packages installed on the system
// Uses the appropriate package manager for the current OS
func getInstalledBrews() ([]string, error) {
	pm, err := pkgmanager.GetPackageManager()
	if err != nil {
		// Fallback to brew if package manager detection fails
		return getInstalledBrewsDirect()
	}

	return pm.ListInstalled("brew")
}

// getInstalledBrewsDirect returns a list of all Homebrew formulas (legacy fallback)
func getInstalledBrewsDirect() ([]string, error) {
	cmd := exec.Command("brew", "list", "--formula")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var brews []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			brews = append(brews, line)
		}
	}
	return brews, nil
}

// getInstalledCasks returns a list of all casks/applications installed on the system
func getInstalledCasks() ([]string, error) {
	pm, err := pkgmanager.GetPackageManager()
	if err != nil {
		// Fallback to brew if package manager detection fails
		return getInstalledCasksDirect()
	}

	return pm.ListInstalled("cask")
}

// getInstalledCasksDirect returns a list of all Homebrew casks (legacy fallback)
func getInstalledCasksDirect() ([]string, error) {
	cmd := exec.Command("brew", "list", "--cask")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var casks []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			casks = append(casks, line)
		}
	}
	return casks, nil
}

// filterNewPackages returns packages from 'installed' that are not in 'existing'
func filterNewPackages(installed, existing []string) []string {
	existingMap := make(map[string]bool)
	for _, pkg := range existing {
		existingMap[pkg] = true
	}

	var newPkgs []string
	for _, pkg := range installed {
		if !existingMap[pkg] {
			newPkgs = append(newPkgs, pkg)
		}
	}

	return newPkgs
}
