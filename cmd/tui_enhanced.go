package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"dotfiles/internal/config"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tab types
type tab int

const (
	packagesTab tab = iota
	snapshotsTab
	hooksTab
	statsTab
	profilesTab
)

// Sort modes
type sortMode int

const (
	sortByName sortMode = iota
	sortByType
	sortByStatus
)

// Enhanced styles
var (
	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("235")).
			Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(0, 2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("240")).
			MarginBottom(1)

	detailsPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(1, 2).
				MarginLeft(2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)
)

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	Select      key.Binding
	SelectAll   key.Binding
	DeselectAll key.Binding
	Add         key.Binding
	Remove      key.Binding
	Install     key.Binding
	Uninstall   key.Binding
	Save        key.Binding
	Search      key.Binding
	Sort        key.Binding
	Help        key.Binding
	Quit        key.Binding
	Tab         key.Binding
	Command     key.Binding
}

var keys = keyMap{
	Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
	Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
	Left:        key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "prev tab")),
	Right:       key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "next tab")),
	Select:      key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "select")),
	SelectAll:   key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "select all")),
	DeselectAll: key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "deselect all")),
	Add:         key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add to config")),
	Remove:      key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "remove from config")),
	Install:     key.NewBinding(key.WithKeys("I"), key.WithHelp("I", "install selected")),
	Uninstall:   key.NewBinding(key.WithKeys("U"), key.WithHelp("U", "uninstall selected")),
	Save:        key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save & quit")),
	Search:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Sort:        key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "cycle sort")),
	Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:        key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Tab:         key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
	Command:     key.NewBinding(key.WithKeys(":"), key.WithHelp(":", "command mode")),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.SelectAll, k.DeselectAll},
		{k.Add, k.Remove, k.Install, k.Uninstall},
		{k.Save, k.Search, k.Sort, k.Command},
		{k.Tab, k.Help, k.Quit},
	}
}

type enhancedModel struct {
	cursor          int
	selected        map[int]bool
	packages        []packageItem
	filteredPackages []packageItem
	snapshots       []Snapshot
	profiles        []MachineProfile
	config          *config.Config
	currentTab      tab
	sortMode        sortMode
	searchMode      bool
	searchInput     textinput.Model
	commandMode     bool
	commandInput    textinput.Model
	message         string
	messageType     string // "success", "error", "warning", "info"
	quitting        bool
	showHelp        bool
	help            help.Model
	spinner         spinner.Model
	loading         bool
	detailsVisible  bool
	width           int
	height          int
	stats           Stats
}

type Stats struct {
	TotalPackages     int
	InstalledPackages int
	ConfiguredPackages int
	DiskSpaceUsed     string
	LastSync          string
	DriftCount        int
}

type packageDetails struct {
	Name         string
	Type         string
	Description  string
	Version      string
	Dependencies []string
	Homepage     string
	InstallSize  string
}

func enhancedInitialModel() enhancedModel {
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
	packages := buildPackageList(cfg, installedBrews, installedCasks)

	// Load snapshots
	snapshots := loadSnapshots()

	// Load profiles
	profiles := loadProfiles()

	// Calculate stats
	stats := calculateStats(cfg, installedBrews, installedCasks)

	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search packages..."
	searchInput.CharLimit = 50

	// Initialize command input
	commandInput := textinput.New()
	commandInput.Placeholder = "Enter command..."
	commandInput.CharLimit = 100

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot

	return enhancedModel{
		cursor:          0,
		selected:        make(map[int]bool),
		packages:        packages,
		filteredPackages: packages,
		snapshots:       snapshots,
		profiles:        profiles,
		config:          cfg,
		currentTab:      packagesTab,
		sortMode:        sortByName,
		searchInput:     searchInput,
		commandInput:    commandInput,
		help:            help.New(),
		spinner:         s,
		stats:           stats,
	}
}

func buildPackageList(cfg *config.Config, installedBrews, installedCasks []string) []packageItem {
	packages := []packageItem{}
	brewsMap := make(map[string]bool)
	casksMap := make(map[string]bool)

	// Add configured brews
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

	return packages
}

func loadSnapshots() []Snapshot {
	home, _ := os.UserHomeDir()
	snapshotsDir := filepath.Join(home, ".dotfiles", "snapshots")

	snapshots := []Snapshot{}
	if _, err := os.Stat(snapshotsDir); os.IsNotExist(err) {
		return snapshots
	}

	entries, _ := os.ReadDir(snapshotsDir)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
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

	// Sort by timestamp (newest first)
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp > snapshots[j].Timestamp
	})

	return snapshots
}

func loadProfiles() []MachineProfile {
	home, _ := os.UserHomeDir()
	profilesDir := filepath.Join(home, ".dotfiles", "profiles")

	profiles := []MachineProfile{}
	if _, err := os.Stat(profilesDir); os.IsNotExist(err) {
		return profiles
	}

	entries, _ := os.ReadDir(profilesDir)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			profilePath := filepath.Join(profilesDir, entry.Name())
			data, err := os.ReadFile(profilePath)
			if err != nil {
				continue
			}

			var profile MachineProfile
			if err := json.Unmarshal(data, &profile); err != nil {
				continue
			}

			profiles = append(profiles, profile)
		}
	}

	return profiles
}

func calculateStats(cfg *config.Config, installedBrews, installedCasks []string) Stats {
	totalPkgs := len(cfg.Brews) + len(cfg.Casks)
	installedPkgs := len(installedBrews) + len(installedCasks)

	// Calculate drift
	drift := 0
	for _, brew := range cfg.Brews {
		if !contains(installedBrews, brew) {
			drift++
		}
	}
	for _, cask := range cfg.Casks {
		if !contains(installedCasks, cask) {
			drift++
		}
	}

	return Stats{
		TotalPackages:      totalPkgs,
		InstalledPackages:  installedPkgs,
		ConfiguredPackages: totalPkgs,
		DiskSpaceUsed:      "Calculating...",
		LastSync:           "Unknown",
		DriftCount:         drift,
	}
}

func (m enhancedModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m enhancedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

	case tea.KeyMsg:
		// Handle search mode
		if m.searchMode {
			if msg.String() == "esc" {
				m.searchMode = false
				m.searchInput.Reset()
				m.filteredPackages = m.packages
			} else if msg.String() == "enter" {
				m.searchMode = false
				m.filterPackages()
			} else {
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.filterPackages()
				return m, cmd
			}
			return m, nil
		}

		// Handle command mode
		if m.commandMode {
			if msg.String() == "esc" {
				m.commandMode = false
				m.commandInput.Reset()
			} else if msg.String() == "enter" {
				m.executeCommand()
				m.commandMode = false
				m.commandInput.Reset()
			} else {
				var cmd tea.Cmd
				m.commandInput, cmd = m.commandInput.Update(msg)
				return m, cmd
			}
			return m, nil
		}

		// Handle normal key presses
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Down):
			maxItems := len(m.filteredPackages)
			if m.currentTab == snapshotsTab {
				maxItems = len(m.snapshots)
			} else if m.currentTab == profilesTab {
				maxItems = len(m.profiles)
			}
			if m.cursor < maxItems-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Left):
			if m.currentTab > 0 {
				m.currentTab--
				m.cursor = 0
				m.selected = make(map[int]bool)
			}

		case key.Matches(msg, keys.Right):
			if m.currentTab < profilesTab {
				m.currentTab++
				m.cursor = 0
				m.selected = make(map[int]bool)
			}

		case key.Matches(msg, keys.Tab):
			m.currentTab = (m.currentTab + 1) % (profilesTab + 1)
			m.cursor = 0
			m.selected = make(map[int]bool)

		case key.Matches(msg, keys.Select):
			if m.currentTab == packagesTab {
				if _, ok := m.selected[m.cursor]; ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = true
				}
			}

		case key.Matches(msg, keys.SelectAll):
			if m.currentTab == packagesTab {
				for i := range m.filteredPackages {
					m.selected[i] = true
				}
				m.setMessage("Selected all packages", "info")
			}

		case key.Matches(msg, keys.DeselectAll):
			m.selected = make(map[int]bool)
			m.setMessage("Deselected all", "info")

		case key.Matches(msg, keys.Add):
			if m.currentTab == packagesTab {
				added := m.addSelectedToConfig()
				m.setMessage(fmt.Sprintf("Added %d package(s) to config", added), "success")
				m.selected = make(map[int]bool)
			}

		case key.Matches(msg, keys.Remove):
			if m.currentTab == packagesTab {
				removed := m.removeSelectedFromConfig()
				m.setMessage(fmt.Sprintf("Removed %d package(s) from config", removed), "success")
				m.selected = make(map[int]bool)
			}

		case key.Matches(msg, keys.Save):
			home, _ := os.UserHomeDir()
			configPath := filepath.Join(home, ".dotfiles", "config.json")
			if err := m.config.Save(configPath); err != nil {
				m.setMessage(fmt.Sprintf("Error saving: %v", err), "error")
			} else {
				m.setMessage("Configuration saved!", "success")
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, keys.Search):
			m.searchMode = true
			m.searchInput.Focus()

		case key.Matches(msg, keys.Sort):
			m.sortMode = (m.sortMode + 1) % 3
			m.sortPackages()
			modes := []string{"name", "type", "status"}
			m.setMessage(fmt.Sprintf("Sorted by: %s", modes[m.sortMode]), "info")

		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp

		case key.Matches(msg, keys.Command):
			m.commandMode = true
			m.commandInput.Focus()

		case msg.String() == "d":
			m.detailsVisible = !m.detailsVisible

		case msg.String() == "enter":
			if m.currentTab == snapshotsTab && m.cursor < len(m.snapshots) {
				m.restoreSnapshot(m.snapshots[m.cursor].Timestamp)
			} else if m.currentTab == profilesTab && m.cursor < len(m.profiles) {
				m.importProfile(m.profiles[m.cursor])
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m enhancedModel) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Header with tabs
	s.WriteString(m.renderTabs() + "\n\n")

	// Content based on current tab
	switch m.currentTab {
	case packagesTab:
		s.WriteString(m.renderPackagesView())
	case snapshotsTab:
		s.WriteString(m.renderSnapshotsView())
	case hooksTab:
		s.WriteString(m.renderHooksView())
	case statsTab:
		s.WriteString(m.renderStatsView())
	case profilesTab:
		s.WriteString(m.renderProfilesView())
	}

	// Search bar
	if m.searchMode {
		s.WriteString("\n" + m.searchInput.View() + "\n")
	}

	// Command bar
	if m.commandMode {
		s.WriteString("\n:" + m.commandInput.View() + "\n")
	}

	// Status message
	if m.message != "" {
		s.WriteString("\n" + m.renderMessage() + "\n")
	}

	// Help
	if m.showHelp {
		s.WriteString("\n" + m.help.View(keys))
	} else {
		s.WriteString("\n" + m.renderQuickHelp())
	}

	return s.String()
}

func (m enhancedModel) renderTabs() string {
	tabs := []string{"📦 Packages", "📸 Snapshots", "🪝 Hooks", "📊 Stats", "📋 Profiles"}
	var rendered []string

	for i, tabName := range tabs {
		if tab(i) == m.currentTab {
			rendered = append(rendered, tabActiveStyle.Render(tabName))
		} else {
			rendered = append(rendered, tabInactiveStyle.Render(tabName))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m enhancedModel) renderPackagesView() string {
	var s strings.Builder

	// Show package count and filter info
	header := fmt.Sprintf("Showing %d packages", len(m.filteredPackages))
	if len(m.selected) > 0 {
		header += fmt.Sprintf(" (%d selected)", len(m.selected))
	}
	s.WriteString(headerStyle.Render(header) + "\n\n")

	// Render packages
	viewStart := maxInt(0, m.cursor-10)
	viewEnd := minInt(len(m.filteredPackages), m.cursor+11)

	for i := viewStart; i < viewEnd; i++ {
		pkg := m.filteredPackages[i]
		s.WriteString(m.renderPackageItem(i, pkg) + "\n")
	}

	// Show details panel if enabled
	if m.detailsVisible && m.cursor < len(m.filteredPackages) {
		pkg := m.filteredPackages[m.cursor]
		s.WriteString("\n" + m.renderPackageDetails(pkg))
	}

	return s.String()
}

func (m enhancedModel) renderPackageItem(i int, pkg packageItem) string {
	cursor := "  "
	if m.cursor == i {
		cursor = "❯ "
	}

	checked := "  "
	if m.selected[i] {
		checked = "✓ "
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

	line := fmt.Sprintf("%s%s%s %-30s %s", cursor, checked, status, pkg.name, pkgTypeTag)

	if m.cursor == i {
		return selectedStyle.Render(line)
	}
	return normalStyle.Render(line)
}

func (m enhancedModel) renderPackageDetails(pkg packageItem) string {
	details := fmt.Sprintf(`Package Details

Name:     %s
Type:     %s
Status:   %s
In Config: %v

Press 'd' to hide details`,
		pkg.name,
		pkg.pkgType,
		m.getPackageStatus(pkg),
		isInConfig(pkg.name, pkg.pkgType, m.config),
	)

	return detailsPanelStyle.Render(details)
}

func (m enhancedModel) getPackageStatus(pkg packageItem) string {
	if pkg.installed {
		return "Installed"
	}
	return "Not Installed"
}

func (m enhancedModel) renderSnapshotsView() string {
	var s strings.Builder

	header := fmt.Sprintf("Snapshots (%d total)", len(m.snapshots))
	s.WriteString(headerStyle.Render(header) + "\n\n")

	if len(m.snapshots) == 0 {
		s.WriteString("  No snapshots found\n")
		s.WriteString("  Create one with: dotfiles snapshot create\n")
		return s.String()
	}

	for i, snapshot := range m.snapshots {
		cursor := "  "
		if m.cursor == i {
			cursor = "❯ "
		}

		t, _ := time.Parse("20060102-150405", snapshot.Timestamp)
		displayTime := t.Format("Jan 02, 2006 at 3:04 PM")

		line := fmt.Sprintf("%s%s - %s (%d packages)",
			cursor,
			snapshot.Timestamp,
			snapshot.Description,
			len(snapshot.Config.Brews)+len(snapshot.Config.Casks),
		)

		if m.cursor == i {
			s.WriteString(selectedStyle.Render(line) + "\n")
			s.WriteString(normalStyle.Render(fmt.Sprintf("     Created: %s", displayTime)) + "\n")
			s.WriteString(normalStyle.Render("     Press Enter to restore") + "\n")
		} else {
			s.WriteString(normalStyle.Render(line) + "\n")
		}
	}

	return s.String()
}

func (m enhancedModel) renderHooksView() string {
	var s strings.Builder

	s.WriteString(headerStyle.Render("Configured Hooks") + "\n\n")

	if m.config.Hooks == nil || isHooksEmpty(m.config.Hooks) {
		s.WriteString("  No hooks configured\n")
		s.WriteString("  Add hooks with: dotfiles hooks add <type> <command>\n")
		return s.String()
	}

	hooks := []struct {
		name  string
		hooks []string
	}{
		{"Pre-Install", m.config.Hooks.PreInstall},
		{"Post-Install", m.config.Hooks.PostInstall},
		{"Pre-Sync", m.config.Hooks.PreSync},
		{"Post-Sync", m.config.Hooks.PostSync},
		{"Pre-Stow", m.config.Hooks.PreStow},
		{"Post-Stow", m.config.Hooks.PostStow},
	}

	for _, h := range hooks {
		if len(h.hooks) > 0 {
			s.WriteString(fmt.Sprintf("  📌 %s:\n", h.name))
			for i, hook := range h.hooks {
				s.WriteString(fmt.Sprintf("     %d. %s\n", i, hook))
			}
			s.WriteString("\n")
		}
	}

	return s.String()
}

func (m enhancedModel) renderStatsView() string {
	var s strings.Builder

	s.WriteString(headerStyle.Render("System Statistics") + "\n\n")

	stats := []string{
		fmt.Sprintf("📦 Total Packages:      %d", m.stats.TotalPackages),
		fmt.Sprintf("✅ Installed:           %d", m.stats.InstalledPackages),
		fmt.Sprintf("📋 In Config:           %d", m.stats.ConfiguredPackages),
		fmt.Sprintf("⚠️  Configuration Drift: %d", m.stats.DriftCount),
		"",
		fmt.Sprintf("💾 Disk Space:          %s", m.stats.DiskSpaceUsed),
		fmt.Sprintf("🔄 Last Sync:           %s", m.stats.LastSync),
		"",
		fmt.Sprintf("📸 Snapshots:           %d", len(m.snapshots)),
		fmt.Sprintf("📋 Profiles:            %d", len(m.profiles)),
	}

	for _, stat := range stats {
		s.WriteString(normalStyle.Render("  "+stat) + "\n")
	}

	// Show drift warning if needed
	if m.stats.DriftCount > 0 {
		s.WriteString("\n" + warningStyle.Render(fmt.Sprintf(
			"  ⚠️  %d package(s) in config but not installed. Run 'dotfiles install'",
			m.stats.DriftCount,
		)) + "\n")
	}

	return s.String()
}

func (m enhancedModel) renderProfilesView() string {
	var s strings.Builder

	header := fmt.Sprintf("Profiles (%d total)", len(m.profiles))
	s.WriteString(headerStyle.Render(header) + "\n\n")

	if len(m.profiles) == 0 {
		s.WriteString("  No profiles found\n")
		s.WriteString("  Create one with: dotfiles export <name>\n")
		return s.String()
	}

	for i, profile := range m.profiles {
		cursor := "  "
		if m.cursor == i {
			cursor = "❯ "
		}

		line := fmt.Sprintf("%s%s - %s",
			cursor,
			profile.Name,
			profile.Description,
		)

		if m.cursor == i {
			s.WriteString(selectedStyle.Render(line) + "\n")
			s.WriteString(normalStyle.Render(fmt.Sprintf("     Packages: %d brews, %d casks",
				len(profile.Config.Brews), len(profile.Config.Casks))) + "\n")
			s.WriteString(normalStyle.Render("     Press Enter to import") + "\n")
		} else {
			s.WriteString(normalStyle.Render(line) + "\n")
		}
	}

	return s.String()
}

func (m enhancedModel) renderMessage() string {
	var style lipgloss.Style
	switch m.messageType {
	case "success":
		style = successStyle
	case "error":
		style = errorStyle
	case "warning":
		style = warningStyle
	default:
		style = statusStyle
	}

	return style.Render("  " + m.message)
}

func (m enhancedModel) renderQuickHelp() string {
	help := "? help • tab switch tabs • / search • : command • s save • q quit"
	return helpStyle.Render(help)
}

func (m *enhancedModel) setMessage(msg, msgType string) {
	m.message = msg
	m.messageType = msgType
}

func (m *enhancedModel) filterPackages() {
	query := strings.ToLower(m.searchInput.Value())
	if query == "" {
		m.filteredPackages = m.packages
		return
	}

	filtered := []packageItem{}
	for _, pkg := range m.packages {
		if strings.Contains(strings.ToLower(pkg.name), query) {
			filtered = append(filtered, pkg)
		}
	}
	m.filteredPackages = filtered
	m.cursor = 0
}

func (m *enhancedModel) sortPackages() {
	switch m.sortMode {
	case sortByName:
		sort.Slice(m.filteredPackages, func(i, j int) bool {
			return m.filteredPackages[i].name < m.filteredPackages[j].name
		})
	case sortByType:
		sort.Slice(m.filteredPackages, func(i, j int) bool {
			if m.filteredPackages[i].pkgType == m.filteredPackages[j].pkgType {
				return m.filteredPackages[i].name < m.filteredPackages[j].name
			}
			return m.filteredPackages[i].pkgType < m.filteredPackages[j].pkgType
		})
	case sortByStatus:
		sort.Slice(m.filteredPackages, func(i, j int) bool {
			iStatus := m.getStatusValue(m.filteredPackages[i])
			jStatus := m.getStatusValue(m.filteredPackages[j])
			if iStatus == jStatus {
				return m.filteredPackages[i].name < m.filteredPackages[j].name
			}
			return iStatus > jStatus
		})
	}
}

func (m *enhancedModel) getStatusValue(pkg packageItem) int {
	inConfig := isInConfig(pkg.name, pkg.pkgType, m.config)
	if inConfig && pkg.installed {
		return 2 // Highest priority
	} else if inConfig {
		return 1
	} else if pkg.installed {
		return 0
	}
	return -1
}

func (m *enhancedModel) executeCommand() {
	cmd := strings.TrimSpace(m.commandInput.Value())
	parts := strings.Fields(cmd)

	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "install":
		m.setMessage("Run 'dotfiles install' from terminal", "info")
	case "sync":
		m.setMessage("Run 'dotfiles sync' from terminal", "info")
	case "snapshot":
		m.setMessage("Run 'dotfiles snapshot create' from terminal", "info")
	case "doctor":
		m.setMessage("Run 'dotfiles doctor' from terminal", "info")
	case "quit", "q":
		m.quitting = true
	default:
		m.setMessage(fmt.Sprintf("Unknown command: %s", parts[0]), "error")
	}
}

func (m *enhancedModel) addSelectedToConfig() int {
	added := 0
	for i := range m.selected {
		if i >= len(m.filteredPackages) {
			continue
		}
		pkg := m.filteredPackages[i]
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

func (m *enhancedModel) removeSelectedFromConfig() int {
	removed := 0
	for i := range m.selected {
		if i >= len(m.filteredPackages) {
			continue
		}
		pkg := m.filteredPackages[i]
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

func (m *enhancedModel) restoreSnapshot(timestamp string) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".dotfiles", "config.json")
	snapshotPath := filepath.Join(home, ".dotfiles", "snapshots", timestamp+".json")

	data, err := os.ReadFile(snapshotPath)
	if err != nil {
		m.setMessage(fmt.Sprintf("Error reading snapshot: %v", err), "error")
		return
	}

	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		m.setMessage(fmt.Sprintf("Error parsing snapshot: %v", err), "error")
		return
	}

	if err := snapshot.Config.Save(configPath); err != nil {
		m.setMessage(fmt.Sprintf("Error restoring: %v", err), "error")
		return
	}

	m.config = snapshot.Config
	m.setMessage(fmt.Sprintf("Restored snapshot: %s", timestamp), "success")
}

func (m *enhancedModel) importProfile(profile MachineProfile) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".dotfiles", "config.json")

	// Merge with existing config
	m.config.Brews = mergeUnique(m.config.Brews, profile.Config.Brews)
	m.config.Casks = mergeUnique(m.config.Casks, profile.Config.Casks)
	m.config.Taps = mergeUnique(m.config.Taps, profile.Config.Taps)
	m.config.Stow = mergeUnique(m.config.Stow, profile.Config.Stow)

	if err := m.config.Save(configPath); err != nil {
		m.setMessage(fmt.Sprintf("Error importing: %v", err), "error")
		return
	}

	m.setMessage(fmt.Sprintf("Imported profile: %s", profile.Name), "success")
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
