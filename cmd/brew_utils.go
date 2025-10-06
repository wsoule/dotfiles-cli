package cmd

import (
	"os/exec"
	"strings"
)

// getInstalledBrews returns a list of all Homebrew formulas installed on the system
func getInstalledBrews() ([]string, error) {
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

// getInstalledCasks returns a list of all Homebrew casks installed on the system
func getInstalledCasks() ([]string, error) {
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
