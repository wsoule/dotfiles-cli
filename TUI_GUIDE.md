# 🎨 Enhanced TUI Guide

The Dotfiles TUI is a feature-rich, interactive terminal interface for managing your entire dotfiles system.

## 🚀 Launch

```bash
dotfiles tui
```

## 📑 Tabs/Views

Navigate between different views using **←/→** or **h/l** or **Tab**:

### 1. 📦 Packages View
Browse and manage all your packages (brews and casks).

**Features:**
- ✅ Color-coded status indicators
- 🔍 Real-time search and filtering
- 📊 Sort by name, type, or status
- 📋 Multi-select with batch operations
- 📝 Package details panel

**Status Indicators:**
- `✅` - In config AND installed
- `📋` - In config but NOT installed
- `⚠️` - Installed but NOT in config
- `[brew]` - Homebrew formula (blue)
- `[cask]` - Homebrew cask (pink)

**Actions:**
- `space` - Select/deselect package
- `ctrl+a` - Select all visible packages
- `ctrl+d` - Deselect all
- `a` - Add selected to config
- `r` - Remove selected from config
- `d` - Toggle package details panel
- `/` - Search packages
- `S` - Cycle sort mode (name → type → status)

### 2. 📸 Snapshots View
View and manage configuration snapshots.

**Features:**
- List all snapshots with timestamps
- View snapshot metadata
- One-click restore

**Actions:**
- `↑/↓` - Navigate snapshots
- `Enter` - Restore selected snapshot
- Auto-creates backup before restoring

### 3. 🪝 Hooks View
View all configured pre/post operation hooks.

**Display:**
- Pre-Install hooks
- Post-Install hooks
- Pre-Sync hooks
- Post-Sync hooks
- Pre-Stow hooks
- Post-Stow hooks

### 4. 📊 Stats Dashboard
Real-time system statistics and health monitoring.

**Displays:**
- 📦 Total packages count
- ✅ Installed packages count
- 📋 Configured packages count
- ⚠️ Configuration drift
- 💾 Disk space used
- 🔄 Last sync time
- 📸 Snapshots count
- 📋 Profiles count

**Warnings:**
- Alerts when config drift is detected
- Suggests running `dotfiles install`

### 5. 📋 Profiles View
Browse and import machine-specific profiles.

**Features:**
- View all saved profiles
- See package counts for each
- Import profiles with one keystroke

**Actions:**
- `↑/↓` - Navigate profiles
- `Enter` - Import selected profile (merges with current config)

## ⌨️ Complete Keyboard Controls

### Navigation
| Key | Action |
|-----|--------|
| `↑` or `k` | Move up |
| `↓` or `j` | Move down |
| `←` or `h` | Previous tab |
| `→` or `l` | Next tab |
| `Tab` | Next tab (cycles) |

### Selection (Packages View)
| Key | Action |
|-----|--------|
| `space` | Toggle selection |
| `ctrl+a` | Select all |
| `ctrl+d` | Deselect all |

### Actions
| Key | Action |
|-----|--------|
| `a` | Add selected to config |
| `r` | Remove selected from config |
| `I` | Install selected (planned) |
| `U` | Uninstall selected (planned) |
| `s` | Save config and quit |

### Search & Sort
| Key | Action |
|-----|--------|
| `/` | Enter search mode |
| `S` | Cycle sort mode |
| `esc` | Exit search mode |
| `enter` | Apply search |

### Display
| Key | Action |
|-----|--------|
| `d` | Toggle details panel |
| `?` | Toggle full help |

### Command Mode
| Key | Action |
|-----|--------|
| `:` | Enter command mode |
| `:install` | Run install |
| `:sync` | Run sync |
| `:snapshot` | Create snapshot |
| `:doctor` | Run health check |
| `:quit` or `:q` | Quit |

### Other
| Key | Action |
|-----|--------|
| `q` or `ctrl+c` | Quit |

## 🔍 Search Features

Press `/` to enter search mode:

1. Type your search query
2. Results filter in real-time
3. Press `Enter` to apply and exit search mode
4. Press `Esc` to cancel and show all packages

**Search is fuzzy** - searches package names case-insensitively.

## 📊 Sort Modes

Press `S` to cycle through sort modes:

1. **By Name** (A-Z) - Default, alphabetical
2. **By Type** - Groups brews together, then casks
3. **By Status** - Prioritizes packages in config & installed

## 💡 Workflow Examples

### Adding New Packages
1. Launch `dotfiles tui`
2. Navigate with `↑/↓` to find packages
3. Press `space` to select packages
4. Press `a` to add to config
5. Press `s` to save and quit
6. Run `dotfiles install` to install

### Bulk Add Untracked Packages
1. Press `/` and search or press `S` to sort by status
2. Scroll to packages marked `⚠️` (installed but not in config)
3. Press `ctrl+a` to select all
4. Press `a` to add to config
5. Press `s` to save

### Restoring from Snapshot
1. Press `→` to go to Snapshots tab
2. Navigate to desired snapshot
3. Press `Enter` to restore
4. Automatic backup created before restore

### Checking System Health
1. Press `→` multiple times to reach Stats tab
2. Review configuration drift
3. Check disk space and sync status

### Importing a Profile
1. Press `→` to navigate to Profiles tab
2. Select desired profile
3. Press `Enter` to merge with current config
4. Press `s` to save

## 🎨 Visual Design

### Color Scheme
- **Pink/Purple** (212) - Active tabs, selected items, titles
- **Gray** (240, 241) - Inactive tabs, borders, help text
- **Blue** (39) - Brew packages
- **Pink** (212) - Cask packages
- **Green** (42, 114) - Success messages
- **Red** (196) - Error messages
- **Orange** (214) - Warning messages

### Panels & Borders
- Rounded borders for detail panels
- Bottom border for headers
- Padding for readability
- Clear visual hierarchy

## 🔮 Advanced Tips

1. **Quick Config Review**: Navigate to Stats tab to see overall system health
2. **Pre-Operation Snapshot**: Go to Packages, make changes, but before saving press `→` to Snapshots and create one manually
3. **Find Config Drift**: Look for `⚠️` indicators in Packages view
4. **Bulk Operations**: Use `ctrl+a` then action keys for fast bulk operations
5. **Command Mode**: Use `:` for vim-style command execution

## 🚧 Planned Features

- [ ] Real-time package installation with progress bars
- [ ] Package recommendations based on what's commonly installed together
- [ ] Diff view showing pending changes
- [ ] Package descriptions from Homebrew
- [ ] Dependency visualization
- [ ] Install size information

## 🐛 Troubleshooting

**TUI doesn't display correctly:**
- Ensure terminal supports Unicode and colors
- Try resizing terminal window
- Use a modern terminal (iTerm2, Ghostty, Alacritty, etc.)

**Can't see all packages:**
- Use search (`/`) to filter
- Scroll with `↑/↓` - only shows 21 items at a time for performance

**Selection not working:**
- Ensure you're in Packages tab (first tab)
- Use `space` not `Enter` for selection

## 📝 Quick Reference Card

```
┌─ NAVIGATION ────────────────────┬─ ACTIONS ──────────────────────┐
│ ↑/k      Move up               │ space    Select/deselect       │
│ ↓/j      Move down             │ ctrl+a   Select all            │
│ ←→/hl    Switch tabs           │ ctrl+d   Deselect all          │
│ Tab      Next tab              │ a        Add to config         │
├─ SEARCH & SORT ────────────────┤ r        Remove from config    │
│ /        Search                │ s        Save & quit           │
│ S        Cycle sort            │ d        Toggle details        │
│ esc      Exit search           ├─ COMMAND ──────────────────────┤
├─ HELP & QUIT ──────────────────┤ :        Command mode          │
│ ?        Toggle help           │ q        Quit                  │
└─────────────────────────────────┴────────────────────────────────┘
```

## 🎯 Summary

The enhanced TUI provides a complete, visual interface for managing your dotfiles:

✅ **5 Tabs** - Packages, Snapshots, Hooks, Stats, Profiles
✅ **Search & Filter** - Real-time fuzzy search
✅ **Multi-Select** - Batch operations on packages
✅ **Sort Modes** - By name, type, or status
✅ **Details Panel** - Extended package information
✅ **Command Mode** - Vim-style command execution
✅ **Color-Coded** - Clear visual status indicators
✅ **Keyboard First** - Complete keyboard control
✅ **Help System** - Built-in contextual help

Launch with `dotfiles tui` and enjoy a modern, efficient dotfiles management experience! 🚀
