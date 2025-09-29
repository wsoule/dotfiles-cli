# 🛠 Dotfiles Manager

> A minimal, focused dotfiles manager with JSON config and Brewfile support

A simple command-line tool that manages your Homebrew packages through JSON configuration and generates Brewfiles.

## Features

- 📦 **Package Management** - Add/remove Homebrew packages, casks, taps, and Stow packages
- 📋 **JSON Configuration** - Store your package list in simple JSON format
- 🍺 **Brewfile Support** - Generate Brewfiles and import from existing ones
- 🔗 **GNU Stow Integration** - Manage dotfiles with symbolic links
- 📊 **Status Checking** - Verify package installation status
- 💾 **Backup & Restore** - Save and restore configurations
- 🔍 **Multiple Output Formats** - JSON, count, and filtered views
- ✨ **Minimal Dependencies** - Only requires Cobra CLI framework and optional GNU Stow
- 🚀 **Fast & Lightweight** - Single binary with comprehensive functionality

## 📦 Installation

### Build from source
```bash
git clone <your-repo>
cd Go_Dotfiles
go build -o dotfiles
```

## 🚀 Quick Start

1. **Initialize your configuration:**
   ```bash
   ./dotfiles init
   ```

2. **Add packages:**
   ```bash
   ./dotfiles add git lazygit
   ./dotfiles add --type=cask visual-studio-code
   ./dotfiles add --type=tap homebrew/cask-fonts
   ```

3. **Add Stow packages for dotfiles:**
   ```bash
   ./dotfiles add --type=stow vim zsh tmux
   ```

4. **List your packages:**
   ```bash
   ./dotfiles list
   ./dotfiles status    # Check installation status
   ```

5. **Generate and install from Brewfile:**
   ```bash
   ./dotfiles install   # Generate Brewfile and install packages
   ```

6. **Manage dotfiles with Stow:**
   ```bash
   ./dotfiles stow vim zsh    # Create symlinks
   ./dotfiles unstow vim      # Remove symlinks
   ```

## 📋 Available Commands

```
Usage:
  dotfiles [command]

Available Commands:
  add         Add packages to your configuration
  backup      Backup your configuration to a file
  brewfile    Generate a Brewfile from your configuration
  import      Import packages from a Brewfile
  init        Initialize a new dotfiles configuration
  install     Generate Brewfile and install packages
  list        List all packages in your configuration
  remove      Remove packages from your configuration
  restore     Restore configuration from a backup file
  restow      Restow dotfile packages (unstow then stow)
  status      Check package installation status
  stow        Stow dotfile packages using GNU Stow
  unstow      Unstow dotfile packages using GNU Stow
```

## 💡 Usage Examples

### Adding different types of packages
```bash
# Add brew packages (default)
./dotfiles add git curl wget

# Add casks (GUI applications)
./dotfiles add --type=cask visual-studio-code firefox slack

# Add taps (additional repositories)
./dotfiles add --type=tap homebrew/cask-fonts

# Add Stow packages (dotfiles)
./dotfiles add --type=stow vim zsh tmux

# Add packages from file
./dotfiles add --file=packages.txt --type=brew
```

### Removing packages
```bash
# Remove brew packages
./dotfiles remove wget

# Remove casks
./dotfiles remove --type=cask firefox

# Remove taps
./dotfiles remove --type=tap homebrew/cask-fonts

# Remove Stow packages
./dotfiles remove --type=stow vim

# Bulk remove all of a type
./dotfiles remove --all-brews
./dotfiles remove --all-casks
./dotfiles remove --all-stow
```

### Working with Brewfiles
```bash
# Generate and install packages automatically
./dotfiles install

# Generate Brewfile in current directory
./dotfiles brewfile

# Generate Brewfile in specific location
./dotfiles brewfile --output ~/my-brewfile

# Import from existing Brewfile
./dotfiles import ~/existing-Brewfile

# Install packages from generated Brewfile
brew bundle --file=./Brewfile
```

### Managing Dotfiles with Stow
```bash
# Create symlinks for dotfiles
./dotfiles stow vim zsh tmux

# Remove symlinks
./dotfiles unstow vim

# Restow (remove and recreate symlinks)
./dotfiles restow vim

# Use custom directories
./dotfiles stow --dir=/path/to/dotfiles --target=~ vim

# Dry run to see what would happen
./dotfiles stow --dry-run --verbose vim
```

### Status and Backup Operations
```bash
# Check installation status of all packages
./dotfiles status

# List packages in different formats
./dotfiles list --json
./dotfiles list --count
./dotfiles list --type=stow

# Backup configuration
./dotfiles backup ~/my-backup.json

# Restore from backup
./dotfiles restore ~/my-backup.json
```

## 📁 Configuration

Your configuration is stored as simple JSON at `~/.dotfiles/config.json`:

```json
{
  "brews": [
    "git",
    "lazygit"
  ],
  "casks": [
    "visual-studio-code"
  ],
  "taps": [
    "homebrew/cask-fonts"
  ],
  "stow": [
    "vim",
    "zsh",
    "tmux"
  ]
}
```

This generates a Brewfile like:

```ruby
tap "homebrew/cask-fonts"

brew "git"
brew "lazygit"

cask "visual-studio-code"
```

### Stow Directory Structure

For Stow packages, organize your dotfiles in the `~/.dotfiles` directory:

```
~/.dotfiles/
├── vim/
│   ├── .vimrc
│   └── .vim/
│       └── ... (vim config files)
├── zsh/
│   ├── .zshrc
│   ├── .zprofile
│   └── .zsh/
│       └── ... (zsh config files)
└── tmux/
    └── .tmux.conf
```

When you run `dotfiles stow vim`, it will create symlinks:
- `~/.vimrc` → `~/.dotfiles/vim/.vimrc`
- `~/.vim/` → `~/.dotfiles/vim/.vim/`

## 📚 Command Reference

| Command | Description | Key Flags |
|---------|-------------|-----------|
| `dotfiles init` | Initialize new configuration | None |
| `dotfiles add <packages>` | Add packages to config | `--type=brew/cask/tap/stow`, `--file=<path>` |
| `dotfiles remove <packages>` | Remove packages from config | `--type=brew/cask/tap/stow`, `--all-*`, `--file=<path>` |
| `dotfiles list` | List configured packages | `--json`, `--count`, `--type=<type>` |
| `dotfiles status` | Check package installation status | None |
| `dotfiles install` | Generate Brewfile and install | `--dry-run` |
| `dotfiles brewfile` | Generate Brewfile | `--output=<path>` |
| `dotfiles import <brewfile>` | Import from Brewfile | `--replace` |
| `dotfiles backup <file>` | Backup configuration | None |
| `dotfiles restore <file>` | Restore from backup | `--no-backup` |
| `dotfiles stow <packages>` | Create symlinks with Stow | `--dir=<path>`, `--target=<path>`, `--dry-run`, `--verbose` |
| `dotfiles unstow <packages>` | Remove symlinks | `--dir=<path>`, `--target=<path>`, `--all`, `--keep-config` |
| `dotfiles restow <packages>` | Restow (unstow + stow) | `--dir=<path>`, `--target=<path>`, `--all` |

## 🛠 Development

### Prerequisites
- Go 1.25.1 or later
- Homebrew (for package management features)
- GNU Stow (for dotfiles symlinking features): `brew install stow`

### Building
```bash
go build -o dotfiles
```

### Project Structure
```
Go_Dotfiles/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command
│   ├── init.go         # Initialize config
│   ├── add.go          # Add/remove packages
│   ├── list.go         # List packages
│   ├── status.go       # Status checking
│   ├── install.go      # Install packages
│   ├── brewfile.go     # Generate Brewfile
│   ├── import.go       # Import from Brewfile
│   ├── backup.go       # Backup/restore
│   └── stow.go         # GNU Stow integration
├── internal/config/    # Configuration management
│   └── config.go       # JSON config handling
└── main.go             # Entry point
```

## 📄 License

MIT License
