package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"dotfiles/internal/config"
	"dotfiles/internal/pkgmanager"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TemplateMetadata for template view
type TemplateMetadata struct {
	Name        string
	Description string
	Category    string
}

// Panel types - lazygit style
type panel int

const (
	mainPanel panel = iota
	detailPanel
	legendPanel
)

// View types - different screens
type view int

const (
	packagesView view = iota
	templatesView
	statusView
	stowView
	snapshotsView
	installView
)

// Advanced TUI model - lazygit inspired
type advancedModel struct {
	// Layout
	width        int
	height       int
	activePanel  panel
	currentView  view

	// State
	cursor       int
	selected     map[int]bool
	searchMode   bool
	searchQuery  string
	installing   bool
	installLog   []string

	// Data
	packages     []packageItem
	filteredPkgs []packageItem
	templates    []TemplateMetadata
	config       *config.Config
	configPath   string
	pm           pkgmanager.PackageManager
	osName       string
	pmName       string

	// Status
	message      string
	messageType  string
	quitting     bool
	lastUpdate   time.Time
}

// Messages for async operations
type installCompleteMsg struct{ err error }
type installProgressMsg struct{ line string }

// Styles - more defined than before
var (
	// Panel styles
	mainPanelBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	detailPanelBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1, 2)

	legendPanelBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(0, 1)

	activePanelBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("212")).
				Padding(1, 2)

	// Title styles
	panelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	// View tabs
	activeTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("235")).
			Padding(0, 2).
			MarginRight(1)

	inactiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 2).
			MarginRight(1)

	// Item styles
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	selectedStyleNew = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	installedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	notInstalledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))

	driftStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	// Status styles
	successMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	errorMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	infoMsg = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	// Legend item style
	legendItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	legendKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)
)

func newAdvancedModel() advancedModel {
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

	// Get package manager
	pm, _ := pkgmanager.GetPackageManager()
	pmName := "unknown"
	if pm != nil && pm.IsAvailable() {
		pmName = pm.GetName()
	}

	// Get installed packages
	installedBrews, _ := getInstalledBrews()
	installedCasks, _ := getInstalledCasks()

	// Build package list
	packages := buildPackageList(cfg, installedBrews, installedCasks)

	// Load templates
	templates := loadTemplatesList()

	return advancedModel{
		width:        120,
		height:       40,
		activePanel:  mainPanel,
		currentView:  packagesView,
		cursor:       0,
		selected:     make(map[int]bool),
		packages:     packages,
		filteredPkgs: packages,
		templates:    templates,
		config:       cfg,
		configPath:   configPath,
		pm:           pm,
		osName:       runtime.GOOS,
		pmName:       pmName,
		lastUpdate:   time.Now(),
	}
}

func loadTemplatesList() []TemplateMetadata {
	templates := []TemplateMetadata{}

	// Load from configTemplates (includes both hard-coded and embedded)
	for name, tmpl := range configTemplates {
		templates = append(templates, TemplateMetadata{
			Name:        name,
			Description: tmpl.Metadata.Description,
			Category:    strings.Join(tmpl.Metadata.Tags, ", "),
		})
	}

	return templates
}

func (m advancedModel) Init() tea.Cmd {
	return nil
}

func (m advancedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case installCompleteMsg:
		m.installing = false
		if msg.err != nil {
			m.message = fmt.Sprintf("Install failed: %v", msg.err)
			m.messageType = "error"
		} else {
			m.message = "Installation complete!"
			m.messageType = "success"
			// Refresh package list
			installedBrews, _ := getInstalledBrews()
			installedCasks, _ := getInstalledCasks()
			m.packages = buildPackageList(m.config, installedBrews, installedCasks)
			m.applyFilter()
		}
		return m, nil

	case installProgressMsg:
		m.installLog = append(m.installLog, msg.line)
		if len(m.installLog) > 20 {
			m.installLog = m.installLog[1:]
		}
		return m, nil

	case tea.KeyMsg:
		// Search mode
		if m.searchMode {
			switch msg.String() {
			case "esc":
				m.searchMode = false
				m.searchQuery = ""
				m.filteredPkgs = m.packages
				m.cursor = 0
			case "enter":
				m.searchMode = false
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.applyFilter()
				}
			default:
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
					m.applyFilter()
				}
			}
			return m, nil
		}

		// Normal mode
		switch msg.String() {
		// Quit
		case "q", "ctrl+c":
			if m.installing {
				m.message = "Cannot quit while installing"
				m.messageType = "error"
			} else {
				m.quitting = true
				return m, tea.Quit
			}

		// View switching (number keys)
		case "1":
			m.currentView = packagesView
			m.cursor = 0
		case "2":
			m.currentView = templatesView
			m.cursor = 0
		case "3":
			m.currentView = statusView
		case "4":
			m.currentView = stowView
			m.cursor = 0
		case "5":
			m.currentView = snapshotsView
			m.cursor = 0
		case "6":
			m.currentView = installView

		// Panel switching
		case "tab":
			if m.activePanel == mainPanel {
				m.activePanel = detailPanel
			} else {
				m.activePanel = mainPanel
			}

		// Navigation
		case "j", "down":
			maxItems := m.getMaxCursor()
			if m.cursor < maxItems-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "g":
			m.cursor = 0

		case "G":
			m.cursor = m.getMaxCursor() - 1

		// Page navigation
		case "ctrl+d":
			m.cursor = min(m.cursor+10, m.getMaxCursor()-1)

		case "ctrl+u":
			m.cursor = max(m.cursor-10, 0)

		// Selection
		case " ":
			if m.currentView == packagesView && !m.installing {
				if _, ok := m.selected[m.cursor]; ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = true
				}
			}

		// Actions
		case "a":
			if m.currentView == packagesView && !m.installing {
				count := m.addSelectedToConfig()
				m.message = fmt.Sprintf("✓ Added %d package(s) to config", count)
				m.messageType = "success"
				m.selected = make(map[int]bool)
			}

		case "r":
			if m.currentView == packagesView && !m.installing {
				count := m.removeSelectedFromConfig()
				m.message = fmt.Sprintf("✓ Removed %d package(s) from config", count)
				m.messageType = "success"
				m.selected = make(map[int]bool)
			}

		case "enter":
			if m.currentView == templatesView && m.cursor < len(m.templates) {
				return m, m.applyTemplate(m.templates[m.cursor])
			} else if m.currentView == packagesView && !m.installing {
				// Quick add single package
				if m.cursor < len(m.filteredPkgs) {
					m.selected[m.cursor] = true
					count := m.addSelectedToConfig()
					m.selected = make(map[int]bool)
					m.message = fmt.Sprintf("✓ Added %d to config", count)
					m.messageType = "success"
				}
			}

		case "i", "I":
			if !m.installing {
				m.currentView = installView
				return m, m.runInstall()
			}

		case "s":
			if !m.installing {
				if err := m.config.Save(m.configPath); err != nil {
					m.message = fmt.Sprintf("✗ Error: %v", err)
					m.messageType = "error"
				} else {
					m.message = "✓ Configuration saved!"
					m.messageType = "success"
				}
			}

		case "/":
			if m.currentView == packagesView {
				m.searchMode = true
				m.searchQuery = ""
			}
		}
	}

	return m, nil
}

func (m advancedModel) View() string {
	if m.quitting {
		return ""
	}

	// Calculate dimensions
	mainWidth := int(float64(m.width) * 0.55)
	sideWidth := m.width - mainWidth - 4
	mainHeight := m.height - 8 // Leave space for header and footer
	legendHeight := 6

	// Header with view tabs
	header := m.renderHeader()

	// Main content area
	var mainContent string
	switch m.currentView {
	case packagesView:
		mainContent = m.renderPackagesPanel(mainWidth, mainHeight)
	case templatesView:
		mainContent = m.renderTemplatesPanel(mainWidth, mainHeight)
	case statusView:
		mainContent = m.renderStatusPanel(mainWidth, mainHeight)
	case stowView:
		mainContent = m.renderStowPanel(mainWidth, mainHeight)
	case snapshotsView:
		mainContent = m.renderSnapshotsPanel(mainWidth, mainHeight)
	case installView:
		mainContent = m.renderInstallPanel(mainWidth, mainHeight)
	}

	// Detail panel (right side)
	detailContent := m.renderDetailPanel(sideWidth, mainHeight-legendHeight-2)

	// Legend panel (bottom right)
	legendContent := m.renderLegendPanel(sideWidth, legendHeight)

	// Stack detail and legend
	rightSide := lipgloss.JoinVertical(
		lipgloss.Left,
		detailContent,
		legendContent,
	)

	// Combine main and right side
	mainRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainContent,
		rightSide,
	)

	// Footer
	footer := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		mainRow,
		footer,
	)
}

func (m advancedModel) renderHeader() string {
	// Title
	title := panelTitle.Render("  📦 DOTFILES MANAGER  ")

	// View tabs
	views := []struct {
		name string
		v    view
	}{
		{"Packages", packagesView},
		{"Templates", templatesView},
		{"Status", statusView},
		{"Stow", stowView},
		{"Snapshots", snapshotsView},
		{"Install", installView},
	}

	var tabs []string
	for _, view := range views {
		label := fmt.Sprintf(" %s ", view.name)
		if m.currentView == view.v {
			tabs = append(tabs, activeTab.Render(label))
		} else {
			tabs = append(tabs, inactiveTab.Render(label))
		}
	}
	tabsLine := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)

	// System info
	sysInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(fmt.Sprintf(" %s • %s ", m.osName, m.pmName))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		tabsLine+sysInfo,
		"",
	)
}

func (m advancedModel) renderPackagesPanel(width, height int) string {
	var lines []string

	// Panel title
	title := fmt.Sprintf("Packages (%d)", len(m.filteredPkgs))
	if len(m.selected) > 0 {
		title += fmt.Sprintf(" • %d selected", len(m.selected))
	}
	if m.searchQuery != "" {
		title += fmt.Sprintf(" • filter: '%s'", m.searchQuery)
	}

	lines = append(lines, panelTitle.Render(title))
	lines = append(lines, "")

	// Search mode indicator
	if m.searchMode {
		searchLine := fmt.Sprintf("/ %s_", m.searchQuery)
		lines = append(lines, cursorStyle.Render(searchLine))
		lines = append(lines, "")
	}

	// Package list
	viewStart := max(0, m.cursor-height/2)
	viewEnd := min(len(m.filteredPkgs), viewStart+height-4)

	for i := viewStart; i < viewEnd; i++ {
		lines = append(lines, m.renderPackageItem(i))
	}

	content := strings.Join(lines, "\n")

	// Apply border
	style := mainPanelBorder
	if m.activePanel == mainPanel {
		style = activePanelBorder
	}

	return style.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderPackageItem(i int) string {
	if i >= len(m.filteredPkgs) {
		return ""
	}

	pkg := m.filteredPkgs[i]

	// Cursor
	cursor := "  "
	if i == m.cursor {
		cursor = "▶ "
	}

	// Selection checkbox
	checkbox := " "
	if m.selected[i] {
		checkbox = "✓"
	}

	// Status icon
	var statusIcon string
	inConfig := isInConfig(pkg.name, pkg.pkgType, m.config)
	if inConfig && pkg.installed {
		statusIcon = installedStyle.Render("●") // Green dot
	} else if inConfig && !pkg.installed {
		statusIcon = notInstalledStyle.Render("○") // Grey circle
	} else if !inConfig && pkg.installed {
		statusIcon = driftStyle.Render("◆") // Orange diamond
	} else {
		statusIcon = " "
	}

	// Package type badge
	var typeBadge string
	if pkg.pkgType == "brew" {
		typeBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Render("[brew]")
	} else {
		typeBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Render("[cask]")
	}

	// Package name
	nameStyle := lipgloss.NewStyle()
	if i == m.cursor {
		nameStyle = cursorStyle
	}

	line := fmt.Sprintf("%s[%s] %s %-35s %s",
		cursor,
		checkbox,
		statusIcon,
		pkg.name,
		typeBadge,
	)

	return nameStyle.Render(line)
}

func (m advancedModel) renderTemplatesPanel(width, height int) string {
	var lines []string

	lines = append(lines, panelTitle.Render("Available Templates"))
	lines = append(lines, "")

	for i, tmpl := range m.templates {
		cursor := "  "
		if i == m.cursor {
			cursor = "▶ "
		}

		line := fmt.Sprintf("%s%-20s %s", cursor, tmpl.Name, tmpl.Description)

		if i == m.cursor {
			lines = append(lines, cursorStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	lines = append(lines, "")
	lines = append(lines, infoMsg.Render("Press Enter to apply template"))

	content := strings.Join(lines, "\n")

	style := mainPanelBorder
	if m.activePanel == mainPanel {
		style = activePanelBorder
	}

	return style.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderStatusPanel(width, height int) string {
	var lines []string

	lines = append(lines, panelTitle.Render("System Status"))
	lines = append(lines, "")

	// Package stats
	totalPkgs := len(m.config.Brews) + len(m.config.Casks)
	installedCount := 0
	driftCount := 0

	for _, pkg := range m.packages {
		if pkg.installed {
			installedCount++
		}
		if pkg.installed && !isInConfig(pkg.name, pkg.pkgType, m.config) {
			driftCount++
		}
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("📊 Packages"))
	lines = append(lines, fmt.Sprintf("  Total in config: %d", totalPkgs))
	lines = append(lines, fmt.Sprintf("  Installed:       %d", installedCount))
	lines = append(lines, fmt.Sprintf("  Stow packages:   %d", len(m.config.Stow)))

	if driftCount > 0 {
		lines = append(lines, "")
		lines = append(lines, driftStyle.Render(fmt.Sprintf("  ⚠ Drift: %d packages", driftCount)))
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("💻 System"))
	lines = append(lines, fmt.Sprintf("  OS:              %s", m.osName))
	lines = append(lines, fmt.Sprintf("  Package Manager: %s", m.pmName))

	content := strings.Join(lines, "\n")

	style := mainPanelBorder
	if m.activePanel == mainPanel {
		style = activePanelBorder
	}

	return style.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderStowPanel(width, height int) string {
	var lines []string

	lines = append(lines, panelTitle.Render(fmt.Sprintf("Stow Packages (%d)", len(m.config.Stow))))
	lines = append(lines, "")

	if len(m.config.Stow) == 0 {
		lines = append(lines, notInstalledStyle.Render("  No stow packages configured"))
	} else {
		for i, pkg := range m.config.Stow {
			cursor := "  "
			if i == m.cursor {
				cursor = "▶ "
			}

			line := fmt.Sprintf("%s%s", cursor, pkg)
			if i == m.cursor {
				lines = append(lines, cursorStyle.Render(line))
			} else {
				lines = append(lines, line)
			}
		}
	}

	content := strings.Join(lines, "\n")

	style := mainPanelBorder
	if m.activePanel == mainPanel {
		style = activePanelBorder
	}

	return style.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderSnapshotsPanel(width, height int) string {
	var lines []string

	snapshots := loadSnapshots()

	lines = append(lines, panelTitle.Render(fmt.Sprintf("Snapshots (%d)", len(snapshots))))
	lines = append(lines, "")

	if len(snapshots) == 0 {
		lines = append(lines, notInstalledStyle.Render("  No snapshots found"))
		lines = append(lines, "")
		lines = append(lines, infoMsg.Render("  Create: dotfiles snapshot create"))
	} else {
		for i, snap := range snapshots {
			if i >= height-4 {
				break
			}

			cursor := "  "
			if i == m.cursor {
				cursor = "▶ "
			}

			t, _ := time.Parse("20060102-150405", snap.Timestamp)
			displayTime := t.Format("Jan 02 15:04")

			line := fmt.Sprintf("%s%s - %s", cursor, displayTime, snap.Description)
			if i == m.cursor {
				lines = append(lines, cursorStyle.Render(line))
			} else {
				lines = append(lines, line)
			}
		}
	}

	content := strings.Join(lines, "\n")

	style := mainPanelBorder
	if m.activePanel == mainPanel {
		style = activePanelBorder
	}

	return style.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderInstallPanel(width, height int) string {
	var lines []string

	if m.installing {
		lines = append(lines, panelTitle.Render("Installing... "))
	} else {
		lines = append(lines, panelTitle.Render("Installation"))
	}
	lines = append(lines, "")

	if m.installing {
		lines = append(lines, "Installing packages...")
		lines = append(lines, "")
		for _, logLine := range m.installLog {
			lines = append(lines, "  "+logLine)
		}
	} else {
		lines = append(lines, "Press 'i' to install all configured packages")
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("Ready to install: %d packages", len(m.config.Brews)+len(m.config.Casks)))
	}

	content := strings.Join(lines, "\n")

	style := mainPanelBorder
	if m.activePanel == mainPanel {
		style = activePanelBorder
	}

	return style.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderDetailPanel(width, height int) string {
	var lines []string

	lines = append(lines, panelTitle.Render("Details"))
	lines = append(lines, "")

	if m.currentView == packagesView && m.cursor < len(m.filteredPkgs) {
		pkg := m.filteredPkgs[m.cursor]
		inConfig := isInConfig(pkg.name, pkg.pkgType, m.config)

		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Package"))
		lines = append(lines, fmt.Sprintf("  Name: %s", pkg.name))
		lines = append(lines, fmt.Sprintf("  Type: %s", pkg.pkgType))
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Status"))
		if inConfig && pkg.installed {
			lines = append(lines, installedStyle.Render("  ● In config & installed"))
		} else if inConfig && !pkg.installed {
			lines = append(lines, notInstalledStyle.Render("  ○ In config, not installed"))
		} else if !inConfig && pkg.installed {
			lines = append(lines, driftStyle.Render("  ◆ Installed, not in config"))
		}
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Actions"))
		if !inConfig {
			lines = append(lines, "  a - Add to config")
		} else {
			lines = append(lines, "  r - Remove from config")
		}
	} else if m.currentView == templatesView && m.cursor < len(m.templates) {
		tmpl := m.templates[m.cursor]
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Template"))
		lines = append(lines, fmt.Sprintf("  %s", tmpl.Name))
		lines = append(lines, "")
		lines = append(lines, tmpl.Description)
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("Category: %s", tmpl.Category))
		lines = append(lines, "")
		lines = append(lines, "Press Enter to apply")
	} else {
		lines = append(lines, notInstalledStyle.Render("No item selected"))
	}

	content := strings.Join(lines, "\n")

	style := detailPanelBorder
	if m.activePanel == detailPanel {
		style = activePanelBorder
	}

	return style.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderLegendPanel(width, height int) string {
	var lines []string

	lines = append(lines, legendKey.Render("LEGEND"))
	lines = append(lines, "")
	lines = append(lines, installedStyle.Render("● ")+"In config & installed")
	lines = append(lines, notInstalledStyle.Render("○ ")+"In config only")
	lines = append(lines, driftStyle.Render("◆ ")+"Installed only (drift)")
	lines = append(lines, legendItem.Render("✓ Selected"))

	content := strings.Join(lines, "\n")

	return legendPanelBorder.Width(width).Height(height).Render(content)
}

func (m advancedModel) renderFooter() string {
	// Message
	var msgLine string
	if m.message != "" {
		switch m.messageType {
		case "success":
			msgLine = successMsg.Render(m.message)
		case "error":
			msgLine = errorMsg.Render(m.message)
		default:
			msgLine = infoMsg.Render(m.message)
		}
	}

	// Keybindings help
	var help string
	switch m.currentView {
	case packagesView:
		help = "1-6 views • j/k move • space select • a add • r remove • i install • / search • s save • q quit"
	case templatesView:
		help = "1-6 views • j/k move • enter apply • q quit"
	default:
		help = "1-6 views • j/k move • s save • q quit"
	}
	helpLine := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(help)

	if msgLine != "" {
		return msgLine + "\n" + helpLine
	}
	return helpLine
}

// Helper methods
func (m advancedModel) getMaxCursor() int {
	switch m.currentView {
	case packagesView:
		return len(m.filteredPkgs)
	case templatesView:
		return len(m.templates)
	case stowView:
		return len(m.config.Stow)
	case snapshotsView:
		return len(loadSnapshots())
	default:
		return 0
	}
}

func (m *advancedModel) applyFilter() {
	if m.searchQuery == "" {
		m.filteredPkgs = m.packages
		return
	}

	filtered := []packageItem{}
	query := strings.ToLower(m.searchQuery)
	for _, pkg := range m.packages {
		if strings.Contains(strings.ToLower(pkg.name), query) {
			filtered = append(filtered, pkg)
		}
	}
	m.filteredPkgs = filtered
	m.cursor = 0
}

func (m *advancedModel) addSelectedToConfig() int {
	added := 0
	for i := range m.selected {
		if i >= len(m.filteredPkgs) {
			continue
		}
		pkg := m.filteredPkgs[i]
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
	// Auto-save
	m.config.Save(m.configPath)
	return added
}

func (m *advancedModel) removeSelectedFromConfig() int {
	removed := 0
	for i := range m.selected {
		if i >= len(m.filteredPkgs) {
			continue
		}
		pkg := m.filteredPkgs[i]
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
	// Auto-save
	m.config.Save(m.configPath)
	return removed
}

func (m *advancedModel) applyTemplate(tmpl TemplateMetadata) tea.Cmd {
	// This would load and apply a template
	m.message = fmt.Sprintf("Applied template: %s", tmpl.Name)
	m.messageType = "success"
	return nil
}

func (m *advancedModel) runInstall() tea.Cmd {
	m.installing = true
	m.installLog = []string{"Starting installation..."}

	return func() tea.Msg {
		// Generate install file
		if m.pm == nil {
			return installCompleteMsg{err: fmt.Errorf("no package manager available")}
		}

		fileContent, err := m.pm.GenerateInstallFile(m.config.Brews, m.config.Casks, m.config.Taps)
		if err != nil {
			return installCompleteMsg{err: err}
		}

		fileName := "Brewfile"
		if m.pmName != "homebrew" {
			fileName = "packages.txt"
		}

		filePath := "./" + fileName
		if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
			return installCompleteMsg{err: err}
		}

		// Install packages
		var cmd *exec.Cmd
		if m.pmName == "homebrew" {
			cmd = exec.Command("brew", "bundle", "--file="+filePath)
		} else if m.pmName == "pacman" {
			// Read packages and install with yay/pacman
			packages := m.config.Brews
			packages = append(packages, m.config.Casks...)
			if len(packages) > 0 {
				args := append([]string{"-S", "--noconfirm"}, packages...)
				if commandExists("yay") {
					cmd = exec.Command("yay", args...)
				} else {
					cmd = exec.Command("sudo", append([]string{"pacman"}, args...)...)
				}
			}
		}

		if cmd != nil {
			output, err := cmd.CombinedOutput()
			if err != nil {
				return installCompleteMsg{err: fmt.Errorf("%v: %s", err, string(output))}
			}
		}

		return installCompleteMsg{err: nil}
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
