package config

import "time"

// Config represents the complete dotfiles configuration
type Config struct {
	Personal     Personal     `json:"personal" yaml:"personal"`
	Installation Installation `json:"installation" yaml:"installation"`
	System       System       `json:"system" yaml:"system"`
	Development  Development  `json:"development" yaml:"development"`
	Packages     Packages     `json:"packages" yaml:"packages"`
	Directories  []string     `json:"directories" yaml:"directories"`
	StowExclusions []string   `json:"stow_exclusions" yaml:"stow_exclusions"`
	Metadata     Metadata     `json:"metadata" yaml:"metadata"`
}

// Personal contains user personal information
type Personal struct {
	Name   string `json:"name" yaml:"name"`
	Email  string `json:"email" yaml:"email"`
	Editor string `json:"editor" yaml:"editor"`
}

// Installation controls what components get installed
type Installation struct {
	Homebrew      bool `json:"homebrew" yaml:"homebrew"`
	Brewfile      bool `json:"brewfile" yaml:"brewfile"`
	Dotfiles      bool `json:"dotfiles" yaml:"dotfiles"`
	MacOSDefaults bool `json:"macos_defaults" yaml:"macos_defaults"`
	NPMPackages   bool `json:"npm_packages" yaml:"npm_packages"`
}

// System contains macOS system preferences
type System struct {
	Appearance  Appearance  `json:"appearance" yaml:"appearance"`
	Dock        Dock        `json:"dock" yaml:"dock"`
	Finder      Finder      `json:"finder" yaml:"finder"`
	Keyboard    Keyboard    `json:"keyboard" yaml:"keyboard"`
	Security    Security    `json:"security" yaml:"security"`
	Screenshots Screenshots `json:"screenshots" yaml:"screenshots"`
	MenuBar     MenuBar     `json:"menu_bar" yaml:"menu_bar"`
	Safari      Safari      `json:"safari" yaml:"safari"`
}

// Appearance contains UI appearance settings
type Appearance struct {
	DarkMode         bool `json:"dark_mode" yaml:"dark_mode"`
	Enable24HourTime bool `json:"enable_24_hour_time" yaml:"enable_24_hour_time"`
}

// Dock contains dock configuration
type Dock struct {
	Autohide      bool   `json:"autohide" yaml:"autohide"`
	Position      string `json:"position" yaml:"position"`
	TileSize      int    `json:"tile_size" yaml:"tile_size"`
	Magnification bool   `json:"magnification" yaml:"magnification"`
}

// Finder contains Finder preferences
type Finder struct {
	ShowHiddenFiles     bool   `json:"show_hidden_files" yaml:"show_hidden_files"`
	DefaultView         string `json:"default_view" yaml:"default_view"`
	ShowFileExtensions  bool   `json:"show_file_extensions" yaml:"show_file_extensions"`
}

// Keyboard contains keyboard and input settings
type Keyboard struct {
	KeyRepeatRate       int  `json:"key_repeat_rate" yaml:"key_repeat_rate"`
	DisablePressAndHold bool `json:"disable_press_and_hold" yaml:"disable_press_and_hold"`
}

// Security contains security and sleep settings
type Security struct {
	RequirePasswordImmediately bool `json:"require_password_immediately" yaml:"require_password_immediately"`
	DisplaySleepMinutes        int  `json:"display_sleep_minutes" yaml:"display_sleep_minutes"`
	ComputerSleepMinutes       int  `json:"computer_sleep_minutes" yaml:"computer_sleep_minutes"`
}

// Screenshots contains screenshot settings
type Screenshots struct {
	Location string `json:"location" yaml:"location"`
	Format   string `json:"format" yaml:"format"`
}

// MenuBar contains menu bar settings
type MenuBar struct {
	ShowBatteryPercent bool `json:"show_battery_percent" yaml:"show_battery_percent"`
	ShowDate           bool `json:"show_date" yaml:"show_date"`
}

// Safari contains Safari browser settings
type Safari struct {
	ShowDevelopMenu   bool   `json:"show_develop_menu" yaml:"show_develop_menu"`
	DefaultEncoding   string `json:"default_encoding" yaml:"default_encoding"`
}

// Development contains development environment settings
type Development struct {
	Git        Git                    `json:"git" yaml:"git"`
	Languages  map[string]bool        `json:"languages" yaml:"languages"`
	Frameworks map[string]bool        `json:"frameworks" yaml:"frameworks"`
	Tools      map[string]bool        `json:"tools" yaml:"tools"`
	Shell      Shell                  `json:"shell" yaml:"shell"`
	Aliases    map[string]bool        `json:"aliases" yaml:"aliases"`
}

// Git contains git configuration
type Git struct {
	DefaultBranch string `json:"default_branch" yaml:"default_branch"`
	PullRebase    bool   `json:"pull_rebase" yaml:"pull_rebase"`
	PushDefault   string `json:"push_default" yaml:"push_default"`
}

// Shell contains shell configuration
type Shell struct {
	Theme         string             `json:"theme" yaml:"theme"`
	TerminalTheme string             `json:"terminal_theme" yaml:"terminal_theme"`
	Plugins       map[string]bool    `json:"plugins" yaml:"plugins"`
}

// Packages contains package installation preferences
type Packages struct {
	ExtraBrews []string `json:"extra_brews" yaml:"extra_brews"`
	ExtraCasks []string `json:"extra_casks" yaml:"extra_casks"`
	ExtraTaps  []string `json:"extra_taps" yaml:"extra_taps"`
	NPMGlobals []string `json:"npm_globals" yaml:"npm_globals"`
}

// Metadata contains configuration metadata
type Metadata struct {
	Version      string    `json:"version" yaml:"version"`
	CreatedAt    time.Time `json:"created_at" yaml:"created_at"`
	LastModified time.Time `json:"last_modified" yaml:"last_modified"`
	CreatedBy    string    `json:"created_by" yaml:"created_by"`
}