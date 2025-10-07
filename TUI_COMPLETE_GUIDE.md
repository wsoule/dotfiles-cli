# 🎨 Complete TUI Guide - Lazygit-Inspired Dotfiles Hub

## Overview

The new TUI is a **complete hub** for managing your dotfiles, inspired by lazygit's windowed layout and keyboard-driven interface.

## 🪟 Layout (Lazygit-Style)

```
┌─────────────────────────────────────────────────────────────────────┐
│  📦 DOTFILES MANAGER                                                │
│  [Packages] [Templates] [Status] [Stow] [Snapshots] [Install]      │
│                                            darwin • homebrew         │
├──────────────────────────────────┬──────────────────────────────────┤
│                                  │                                  │
│  MAIN PANEL                      │  DETAIL PANEL                   │
│  (List of items)                 │  (Info about selected item)     │
│                                  │                                  │
│  ▶ [✓] ● package-name [brew]    │  Package                        │
│    [ ] ○ another-pkg  [cask]    │    Name: package-name           │
│    [ ] ◆ drift-pkg    [brew]    │    Type: brew                   │
│                                  │                                  │
│                                  │  Status                         │
│                                  │    ● In config & installed      │
│                                  │                                  │
│                                  ├──────────────────────────────────┤
│                                  │  LEGEND                         │
│                                  │                                  │
│                                  │  ● In config & installed        │
│                                  │  ○ In config only               │
│                                  │  ◆ Installed only (drift)       │
│                                  │  ✓ Selected                     │
└──────────────────────────────────┴──────────────────────────────────┘
  ✓ Added to config
  1-6 views • j/k move • space select • a add • r remove • i install
```

### Panels

1. **Main Panel (Left)** - Primary content
   - Package lists
   - Template browser
   - Snapshots
   - Stats
   - Install logs

2. **Detail Panel (Top Right)** - Context information
   - Selected item details
   - Available actions
   - Additional metadata

3. **Legend Panel (Bottom Right)** - Symbol meanings
   - Always visible
   - Explains all symbols/emojis
   - Quick reference

## 📋 Views (Press 1-6)

### 1. 📦 Packages View

Browse and manage all packages.

**Main Panel:**
```
Packages (245) • 3 selected

  ▶ [✓] ● git                 [brew]
    [✓] ● neovim              [brew]
    [✓] ○ visual-studio-code  [cask]
    [ ] ◆ docker              [brew]
    [ ] ○ firefox             [cask]
```

**Detail Panel:**
```
Package
  Name: git
  Type: brew

Status
  ● In config & installed

Actions
  r - Remove from config
```

**Symbols:**
- `▶` - Cursor (current selection)
- `✓` - Selected for batch operation
- `●` - In config AND installed (green)
- `○` - In config but NOT installed (gray)
- `◆` - Installed but NOT in config - drift (orange)
- `[brew]` - CLI tool (blue)
- `[cask]` - GUI app (pink)

**Actions:**
- `space` - Select/deselect for batch operations
- `enter` - Quick-add single package to config
- `a` - Add all selected to config
- `r` - Remove all selected from config
- `/` - Search packages (fuzzy filter)
- `j/k` - Move up/down
- `g/G` - Jump to top/bottom
- `ctrl+d/u` - Page down/up

### 2. 📚 Templates View

Browse and apply configuration templates.

**Available Templates:**
- `minimal` - Essential tools only
- `web-dev` - Web development stack (Node, Python, Docker)
- `data-science` - Python, R, Jupyter, analytics
- `devops` - Kubernetes, Terraform, cloud tools
- `mobile-dev` - iOS/Android development tools

**Main Panel:**
```
Available Templates

  ▶ minimal         Essential tools only
    web-dev         Web development stack
    data-science    Python, R, Jupyter, analytics
    devops          Kubernetes, Docker, Terraform
    mobile-dev      iOS/Android development
```

**Detail Panel:**
```
Template
  web-dev

Web development stack

Category: development

Press Enter to apply
```

**Actions:**
- `j/k` - Move up/down
- `enter` - Apply template (merges with your config)

### 3. 📊 Status View

System statistics and health check.

**Shows:**
```
System Status

📊 Packages
  Total in config: 45
  Installed:       42
  Stow packages:   5

  ⚠ Drift: 3 packages

💻 System
  OS:              darwin
  Package Manager: homebrew
```

**Detail Panel:**
- Cross-platform notes
- Platform-specific guidance
- Recommendations

### 4. 🔗 Stow View

Manage dotfile packages with GNU Stow.

**Main Panel:**
```
Stow Packages (5)

  ▶ vim
    zsh
    tmux
    git
    nvim
```

**Shows:**
- All configured stow packages
- Navigate with j/k

### 5. 📸 Snapshots View

View and restore configuration snapshots.

**Main Panel:**
```
Snapshots (12)

  ▶ Jan 02 15:04 - Before major update
    Dec 28 10:30 - Clean state
    Dec 15 14:20 - Initial setup
```

**Actions:**
- `j/k` - Navigate
- `enter` - Restore snapshot (with backup)

### 6. ⚡ Install View

Install packages directly from the TUI!

**Before Install:**
```
Installation

Press 'i' to install all configured packages

Ready to install: 45 packages
```

**During Install:**
```
Installing...

Installing packages...

  ==> Downloading git
  ==> Installing git
  ✓ git installed
  ==> Downloading neovim
  ...
```

**Features:**
- Real-time installation logs
- Scrolling output
- Can't quit while installing (safety)
- Auto-refreshes package list when done

**Actions:**
- `i` - Start installation (from any view)

## ⌨️ Keyboard Shortcuts

### Navigation
| Key | Action |
|-----|--------|
| `j` or `↓` | Move down |
| `k` or `↑` | Move up |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `ctrl+d` | Page down (10 items) |
| `ctrl+u` | Page up (10 items) |

### View Switching
| Key | View |
|-----|------|
| `1` | Packages |
| `2` | Templates |
| `3` | Status |
| `4` | Stow |
| `5` | Snapshots |
| `6` | Install |
| `tab` | Switch panel (main ↔ detail) |

### Actions (Packages View)
| Key | Action |
|-----|--------|
| `space` | Select/deselect package |
| `enter` | Quick-add package to config |
| `a` | Add selected to config |
| `r` | Remove selected from config |
| `/` | Search/filter packages |

### Global Actions
| Key | Action |
|-----|--------|
| `i` | Go to install view & start install |
| `s` | Save config |
| `q` | Quit |
| `ctrl+c` | Force quit |

## 🎯 Workflows

### Quick Add a Package
```
1. Launch: dotfiles tui
2. Press '1' (if not already in Packages view)
3. Press 'j' or 'k' to find package
4. Press 'enter' to quick-add
5. Package is immediately added and saved!
```

### Batch Add Multiple Packages
```
1. Press '1' for Packages view
2. Press '/' to search
3. Type package name filter
4. Press 'esc' to exit search
5. Press 'space' to select first package
6. Press 'j' to move down
7. Press 'space' to select next
8. Repeat for all packages
9. Press 'a' to add all selected
10. Auto-saved!
```

### Apply a Template
```
1. Press '2' for Templates view
2. Press 'j/k' to browse
3. Press 'tab' to view details in detail panel
4. Press 'enter' to apply template
5. Returns to packages view
6. Auto-saved!
```

### Install Everything
```
1. Press '6' for Install view
   (or press 'i' from any view)
2. Review package count
3. Press 'i' to start installation
4. Watch real-time logs
5. Installation completes
6. Package list auto-refreshes
```

### Check Configuration Drift
```
1. Press '3' for Status view
2. Check "Drift" count
3. If > 0, press '1' for Packages
4. Look for ◆ symbol (orange diamond)
5. Select those packages with 'space'
6. Press 'a' to add them
```

### Restore a Snapshot
```
1. Press '5' for Snapshots view
2. Press 'j/k' to navigate
3. Review details in detail panel
4. Press 'enter' to restore
5. Automatic backup created first
```

## 🎨 Color Scheme

- **Pink/Purple (212)** - Active items, cursor, selections, headers
- **Blue (39)** - Brew packages, info messages
- **Pink (212)** - Cask packages
- **Green (42)** - Installed & in config, success messages
- **Gray (241)** - Not installed, muted text
- **Orange (214)** - Drift warning, not in config
- **Red (196)** - Errors
- **Dark Gray (240)** - Borders, inactive elements

## 🆕 What's Different from Old TUI?

### Old TUI
- Single panel layout
- Tabs across the top
- No legend (had to memorize symbols)
- Limited views
- No installation capability
- No templates
- Smaller, less organized

### New TUI (Lazygit-Inspired)
✅ **3-panel windowed layout** (main, detail, legend)
✅ **Legend always visible** - never guess what symbols mean
✅ **6 complete views** - packages, templates, status, stow, snapshots, install
✅ **Real-time installation** - install packages directly in TUI
✅ **Template application** - apply pre-built configs
✅ **Better navigation** - g/G, ctrl+d/u, page navigation
✅ **Auto-save** - changes save immediately
✅ **Detail panel** - context for every item
✅ **Cross-platform aware** - shows OS and package manager
✅ **Quick-add with enter** - faster single-package workflow

## 💡 Pro Tips

1. **Use the Legend** - It's always visible on the bottom right. No need to memorize symbols!

2. **Quick Add with Enter** - Don't need to select first. Just navigate and press enter.

3. **Auto-Save** - Changes are saved automatically. Press 's' only if you want confirmation.

4. **Tab Between Panels** - Use tab to switch focus to detail panel (though most operations work from main panel).

5. **Page Navigation** - Use `ctrl+d` and `ctrl+u` for fast scrolling through long lists.

6. **Search is Fuzzy** - When searching packages, partial matches work. Type "neo" to find "neovim".

7. **Install from Anywhere** - Press 'i' from any view to jump to install view and start installing.

8. **Watch the Detail Panel** - It shows available actions for the current item.

9. **Can't Quit While Installing** - Safety feature. Let installation finish first.

10. **Number Keys are Fast** - Muscle memory 1-6 for view switching is faster than tabbing.

## 🐛 Troubleshooting

**TUI looks broken/weird:**
- Ensure terminal is at least 120x40 characters
- Use a modern terminal (iTerm2, Ghostty, Alacritty, Kitty)
- Enable Unicode support

**Colors don't show:**
- Check terminal supports 256 colors
- Try: `echo $TERM` should show `xterm-256color` or similar

**Can't see all panels:**
- Maximize terminal window
- TUI adapts to terminal size

**Installation hangs:**
- Check terminal for password prompt (sudo)
- Some package managers need interaction

**Symbols show as boxes:**
- Terminal needs Unicode/emoji support
- Update terminal or use different one

## 📝 Summary

The new TUI is a **complete dotfiles hub** with:

- ✅ Windowed lazygit-style layout
- ✅ 6 complete views for everything
- ✅ Legend panel (never guess symbols!)
- ✅ Real-time package installation
- ✅ Template application
- ✅ Cross-platform aware
- ✅ Better navigation and workflows
- ✅ Auto-save

**Launch it:**
```bash
dotfiles tui
```

Everything you need to manage dotfiles in one beautiful, keyboard-driven interface! 🚀
