package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			PaddingLeft(2)

	normalStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1).
			MarginLeft(2)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("114")).
			Bold(true).
			MarginLeft(2)
)

type packageItem struct {
	name      string
	pkgType   string
	installed bool
}

type model struct {
	cursor   int
	selected map[int]bool
	packages []packageItem
	config   *config.Config
	mode     string // "browse", "add", "remove"
	message  string
	quitting bool
}

func initialModel() model {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".dotfiles", "config.json")
	cfg, err := config.Load(configPath)
	if err != nil {
		cfg = &config.Config{
			Brews: []string{},
			Casks: []string{},
			Taps:  []string{},
			Stow:  []string{},
		}
	}

	// Get installed packages
	installedBrews, _ := getInstalledBrews()
	installedCasks, _ := getInstalledCasks()

	// Create package list
	packages := []packageItem{}

	// Add configured brews
	brewsMap := make(map[string]bool)
	for _, brew := range cfg.Brews {
		brewsMap[brew] = true
		installed := contains(installedBrews, brew)
		packages = append(packages, packageItem{
			name:      brew,
			pkgType:   "brew",
			installed: installed,
		})
	}

	// Add installed but not configured brews
	for _, brew := range installedBrews {
		if !brewsMap[brew] {
			packages = append(packages, packageItem{
				name:      brew,
				pkgType:   "brew",
				installed: true,
			})
		}
	}

	// Add configured casks
	casksMap := make(map[string]bool)
	for _, cask := range cfg.Casks {
		casksMap[cask] = true
		installed := contains(installedCasks, cask)
		packages = append(packages, packageItem{
			name:      cask,
			pkgType:   "cask",
			installed: installed,
		})
	}

	// Add installed but not configured casks
	for _, cask := range installedCasks {
		if !casksMap[cask] {
			packages = append(packages, packageItem{
				name:      cask,
				pkgType:   "cask",
				installed: true,
			})
		}
	}

	return model{
		cursor:   0,
		selected: make(map[int]bool),
		packages: packages,
		config:   cfg,
		mode:     "browse",
		message:  "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.packages)-1 {
				m.cursor++
			}

		case " ":
			// Toggle selection
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = true
			}

		case "a":
			// Add selected packages to config
			if len(m.selected) == 0 {
				m.message = "No packages selected"
			} else {
				added := m.addSelectedToConfig()
				m.message = fmt.Sprintf("Added %d package(s) to config", added)
				m.selected = make(map[int]bool)
			}

		case "r":
			// Remove selected packages from config
			if len(m.selected) == 0 {
				m.message = "No packages selected"
			} else {
				removed := m.removeSelectedFromConfig()
				m.message = fmt.Sprintf("Removed %d package(s) from config", removed)
				m.selected = make(map[int]bool)
			}

		case "i":
			// Install selected packages
			if len(m.selected) == 0 {
				m.message = "No packages selected"
			} else {
				m.message = "Install feature - press 'q' to exit and run: dotfiles install"
			}

		case "s":
			// Save and quit
			home, _ := os.UserHomeDir()
			configPath := filepath.Join(home, ".dotfiles", "config.json")
			if err := m.config.Save(configPath); err != nil {
				m.message = fmt.Sprintf("Error saving: %v", err)
			} else {
				m.message = "Configuration saved!"
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := titleStyle.Render("🎨 Dotfiles Package Manager (TUI)") + "\n\n"

	// Show packages
	for i, pkg := range m.packages {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected[i] {
			checked = "✓"
		}

		status := "  "
		if isInConfig(pkg.name, pkg.pkgType, m.config) {
			if pkg.installed {
				status = "✅"
			} else {
				status = "📋"
			}
		} else {
			if pkg.installed {
				status = "⚠️ "
			}
		}

		pkgTypeTag := ""
		switch pkg.pkgType {
		case "brew":
			pkgTypeTag = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render("[brew]")
		case "cask":
			pkgTypeTag = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("[cask]")
		}

		line := fmt.Sprintf("%s [%s] %s %s %s", cursor, checked, status, pkg.name, pkgTypeTag)

		if m.cursor == i {
			s += selectedStyle.Render(line) + "\n"
		} else {
			s += normalStyle.Render(line) + "\n"
		}

		if i > 20 && m.cursor < i-10 {
			// Only show packages near cursor for performance
			break
		}
	}

	// Status message
	if m.message != "" {
		s += "\n" + statusStyle.Render(m.message) + "\n"
	}

	// Help text
	help := `
Navigation:  ↑/k up • ↓/j down • space select • q quit
Actions:     a add to config • r remove from config • s save & quit
Legend:      ✅ in config & installed • 📋 in config only • ⚠️  installed only
`
	s += helpStyle.Render(help)

	return s
}

func (m *model) addSelectedToConfig() int {
	added := 0
	for i := range m.selected {
		pkg := m.packages[i]
		if !isInConfig(pkg.name, pkg.pkgType, m.config) {
			switch pkg.pkgType {
			case "brew":
				m.config.Brews = append(m.config.Brews, pkg.name)
			case "cask":
				m.config.Casks = append(m.config.Casks, pkg.name)
			}
			added++
		}
	}
	return added
}

func (m *model) removeSelectedFromConfig() int {
	removed := 0
	for i := range m.selected {
		pkg := m.packages[i]
		if isInConfig(pkg.name, pkg.pkgType, m.config) {
			switch pkg.pkgType {
			case "brew":
				m.config.Brews = removeFromSlice(m.config.Brews, pkg.name)
			case "cask":
				m.config.Casks = removeFromSlice(m.config.Casks, pkg.name)
			}
			removed++
		}
	}
	return removed
}

func isInConfig(name, pkgType string, cfg *config.Config) bool {
	switch pkgType {
	case "brew":
		return contains(cfg.Brews, name)
	case "cask":
		return contains(cfg.Casks, name)
	case "tap":
		return contains(cfg.Taps, name)
	case "stow":
		return contains(cfg.Stow, name)
	}
	return false
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "🎨 Interactive TUI for package management",
	Long: `🎨 Interactive Terminal UI

Launch an interactive terminal interface to manage your packages.
Browse, add, and remove packages with a visual interface.

Controls:
  ↑/k - Move up
  ↓/j - Move down
  space - Select/deselect package
  a - Add selected packages to config
  r - Remove selected packages from config
  s - Save and quit
  q - Quit without saving

Examples:
  dotfiles tui                # Launch interactive interface`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
