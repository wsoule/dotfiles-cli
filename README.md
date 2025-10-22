# 🚀 Dotfiles Manager - Complete Development Environment Management

> The all-in-one solution for managing your development environment across any machine

A **comprehensive cross-platform CLI** that handles everything: package management (Homebrew, pacman, apt, yum), dotfiles with GNU Stow, GitHub SSH setup, shell configuration, migrations, templates, and complete developer onboarding automation.

**Works on macOS, Arch Linux, Debian/Ubuntu, and RHEL/Fedora with a single config file!**

## ✨ Features

### 🎯 Core Features
- 🌍 **Cross-Platform** - Single config works on macOS, Arch Linux, Debian/Ubuntu, RHEL/Fedora
- 🎉 **Complete Developer Onboarding** - One command sets up everything for new developers
- 🔐 **GitHub SSH Integration** - Automated SSH key generation and GitHub configuration
- 📦 **Smart Package Management** - Auto-detects and uses your system's package manager
- 🔗 **GNU Stow Integration** - Manage dotfiles with symbolic links
- 📋 **JSON Configuration** - Store your setup in simple, versionable JSON

### 🚀 Advanced Features
- 🤝 **Configuration Sharing** - Share and discover configs via GitHub Gist
- 📚 **Template Library** - Pre-built configs for web-dev, data science, DevOps, and more
- 🔍 **Community Discovery** - Find and browse configurations from other developers
- 🐚 **Interactive Shell Setup** - Install oh-my-zsh, Starship, fzf, and plugins
- 🔄 **Migration Tools** - Import from chezmoi, yadm, bare repos, homesick, dotbot
- 🏥 **Health Checks** - Automated diagnostics with auto-fix capabilities
- 📸 **Snapshots & Rollback** - Version tracking with easy restore
- 🔍 **Config Linting** - Validate and auto-fix configuration issues
- 🔐 **Secrets Management** - Secure storage for API keys and environment variables
- 🎨 **Template Variables** - Dynamic config with `{{VAR}}` substitution
- 🌐 **Conditional Configs** - OS/hostname-specific configurations (work vs home)
- 📦 **Dependency Management** - Package dependencies with topological sort
- 👀 **Watch Mode** - Auto-commit and sync changes
- 🔧 **Shell Completions** - Bash, zsh, fish, PowerShell support

## 📦 Installation

### macOS

#### 🍺 Homebrew (Recommended)
```bash
brew tap wsoule/tap
brew install dotfiles
```

#### 🚀 One-liner Install
```bash
curl -fsSL https://raw.githubusercontent.com/wsoule/dotfiles-cli/main/install.sh | bash
```

### Linux (Arch, Debian/Ubuntu, RHEL/Fedora)

#### 🚀 One-liner Install
```bash
curl -fsSL https://raw.githubusercontent.com/wsoule/dotfiles-cli/main/install-linux.sh | bash
```

#### 📦 Manual Download
```bash
# Download the latest release
curl -L https://github.com/wsoule/dotfiles-cli/releases/latest/download/dotfiles_Linux_x86_64.tar.gz -o dotfiles.tar.gz

# Extract
tar -xzf dotfiles.tar.gz

# Install (requires sudo)
sudo mv dotfiles /usr/local/bin/

# Verify
dotfiles --help
```

### All Platforms

#### 📦 GitHub Releases
Download the latest binary from [releases](https://github.com/wsoule/dotfiles-cli/releases/latest)

#### 🔨 Build from Source
```bash
git clone https://github.com/wsoule/dotfiles-cli.git
cd dotfiles-cli
go build -o dotfiles

# macOS / Linux
sudo mv dotfiles /usr/local/bin/

# Or install to user directory
mkdir -p ~/.local/bin && mv dotfiles ~/.local/bin/
```

## 🚀 Quick Start

### 🎯 Complete Setup (Recommended)
```bash
# Complete onboarding - sets up everything in one command
dotfiles onboard

# Or with your GitHub email
dotfiles onboard --email your@email.com
```

**What this does:**
1. Initializes dotfiles directory and config
2. Sets up GitHub SSH keys
3. Configures shell (zsh or fish)
4. Installs essential development packages
5. Creates proper directory structure

### 🆕 Setup from Repository
```bash
# Clone your existing dotfiles repo
dotfiles setup https://github.com/your-username/your-dotfiles.git

# This will:
# - Clone your repo to ~/.dotfiles/
# - Create stow/ and private/ directories
# - Set up proper directory structure
```

### 🔧 Manual Setup (Step by Step)

1. **Initialize configuration:**
   ```bash
   dotfiles init
   ```

2. **Set up GitHub SSH:**
   ```bash
   dotfiles github setup --email your@email.com
   dotfiles github test  # Verify connection
   ```

3. **Add packages:**
   ```bash
   dotfiles add git lazygit curl
   dotfiles add --type=cask visual-studio-code
   dotfiles add --type=stow vim zsh tmux
   ```

4. **Install everything:**
   ```bash
   dotfiles install      # Installs packages
   dotfiles stow vim zsh # Creates dotfile symlinks
   ```

## 📋 Command Groups

### 🎯 Getting Started
```bash
dotfiles init           # Initialize new configuration
dotfiles onboard        # Complete developer onboarding
dotfiles tutorial       # Interactive tutorial
dotfiles doctor         # Health check with auto-fix
dotfiles doctor --fix   # Auto-fix common issues
```

### 📦 Package Management
```bash
dotfiles add <packages>           # Add packages
dotfiles remove <packages>        # Remove packages
dotfiles install                  # Install all packages
dotfiles update                   # Update packages
dotfiles status                   # Check installation status
dotfiles diff                     # Show config vs installed
dotfiles scan                     # Scan and add installed packages
```

### 🔗 Dotfiles Management
```bash
dotfiles stow <packages>          # Create symlinks
dotfiles unstow <packages>        # Remove symlinks
dotfiles restow <packages>        # Recreate symlinks
dotfiles stow-diff [packages]     # Preview stow changes
dotfiles list                     # List all packages
```

### 🚀 Advanced Features
```bash
# Configuration Management
dotfiles lint                     # Validate config
dotfiles lint --fix               # Auto-fix issues
dotfiles snapshot create          # Create snapshot
dotfiles snapshot list            # List snapshots
dotfiles snapshot restore <id>    # Restore snapshot
dotfiles backup                   # Create timestamped backup
dotfiles restore <file>           # Restore from backup

# Migration Tools
dotfiles migrate detect           # Detect other tools
dotfiles migrate from chezmoi     # Migrate from chezmoi
dotfiles migrate from yadm        # Migrate from yadm

# Shell Configuration
dotfiles shell setup              # Interactive shell setup
dotfiles shell list               # List installed tools

# Templates & Sharing
dotfiles templates list           # Browse templates
dotfiles templates create <name>  # Create template
dotfiles share gist               # Share via GitHub Gist
dotfiles clone <source>           # Import configuration
dotfiles discover featured        # Browse community configs

# Advanced Config
dotfiles conditionals apply       # Apply OS-specific config
dotfiles conditionals test        # Test which conditionals match
dotfiles secrets init             # Initialize secrets management
dotfiles secrets add <NAME>       # Add secret (prompts for value)
dotfiles template render          # Render template variables
dotfiles depends show <package>   # Show dependencies
dotfiles watch                    # Auto-commit changes
dotfiles roles list               # List package roles

# Utilities
dotfiles completion bash          # Generate completions
dotfiles completion install zsh   # Install completions
dotfiles check                    # Verify setup
dotfiles bootstrap                # Generate bootstrap script
```

## 🏥 Health Checks & Auto-Fix

Run diagnostics and automatically fix common issues:

```bash
# Check for issues
dotfiles doctor

# Auto-fix common problems
dotfiles doctor --fix
```

**What `doctor --fix` does:**
- ✅ Creates missing directories
- ✅ Initializes config.json
- ✅ Initializes git repository
- ✅ Installs git and stow automatically
- ✅ Identifies broken symlinks
- ✅ Detects configuration drift

## 🐚 Interactive Shell Setup

Configure your shell with popular tools interactively:

```bash
dotfiles shell setup
```

**Installs:**
- oh-my-zsh / oh-my-fish
- Starship prompt
- fzf (fuzzy finder)
- zsh-autosuggestions
- zsh-syntax-highlighting

## 🔄 Migration from Other Tools

Easily migrate from other dotfiles managers:

```bash
# Auto-detect installed managers
dotfiles migrate detect

# Migrate from specific tool
dotfiles migrate from chezmoi
dotfiles migrate from yadm
dotfiles migrate from bare-repo
dotfiles migrate from homesick
dotfiles migrate from dotbot
```

## 📸 Snapshots & Version Control

Create snapshots for easy rollback:

```bash
# Create snapshot
dotfiles snapshot create "before major update"

# List snapshots
dotfiles snapshot list

# Restore snapshot
dotfiles snapshot restore 20240615-143022

# Compare snapshots
dotfiles snapshot diff v1.0 v2.0

# Delete snapshot
dotfiles snapshot delete 20240615-143022
```

## 🔍 Config Validation & Linting

Validate and fix configuration issues:

```bash
# Validate configuration
dotfiles lint

# Auto-fix issues
dotfiles lint --fix

# Strict mode validation
dotfiles lint --strict
```

**Checks for:**
- ✅ Valid JSON syntax
- ✅ Duplicate packages
- ✅ Missing stow directories
- ✅ Circular dependencies
- ✅ Empty values
- ✅ Undefined variables

## 🎨 Template Variables

Use dynamic variables in your configs:

```bash
# Set variables
dotfiles template set EMAIL "user@example.com"
dotfiles template set NAME "John Doe"

# List variables
dotfiles template list

# Render templates
dotfiles template render
```

In your config files, use `{{VARIABLE}}` syntax:
```
email = {{EMAIL}}
name = {{NAME|default-name}}
```

## 🔐 Secrets Management

Safely manage API keys, tokens, and sensitive environment variables:

```bash
# Initialize secrets file
dotfiles secrets init

# Add a new secret (prompts for value)
dotfiles secrets add API_KEY
dotfiles secrets add DATABASE_URL

# List secret names (values hidden)
dotfiles secrets list

# Generate template for new machines
dotfiles secrets template
```

**How it works:**
- Secrets stored in `~/.dotfiles/private/.env-private` (git-ignored)
- File permissions: `0600` (owner read/write only)
- Template generation removes values for sharing structure
- Source in your shell: `source ~/.dotfiles/private/.env-private`

**Example workflow:**
```bash
# On your current machine
dotfiles secrets init
dotfiles secrets add GITHUB_TOKEN
dotfiles secrets add AWS_ACCESS_KEY
dotfiles secrets template  # Creates .env-private.template

# Commit template to git (values removed)
git add private/.env-private.template
git commit -m "Add secrets template"

# On new machine
cp ~/.dotfiles/private/.env-private.template ~/.dotfiles/private/.env-private
# Edit .env-private and fill in actual values
```

## 🌐 Conditional Configurations

Apply different configs based on OS or hostname - **perfect for work vs home machines**:

```bash
# Test what would apply
dotfiles conditionals test

# Apply conditionals
dotfiles conditionals apply

# Preview without applying
dotfiles conditionals apply --dry-run

# Add a conditional
dotfiles conditionals add \
  --hostname "work-*" \
  --brews slack,zoom,docker \
  --casks microsoft-teams
```

**Example: Work vs Home Setup**
```json
{
  "conditionals": [
    {
      "hostname": "work-*",
      "os": "darwin",
      "brews": ["slack", "zoom", "docker"],
      "casks": ["microsoft-teams"],
      "variables": {
        "WORK": "true",
        "EMAIL": "you@company.com"
      }
    },
    {
      "hostname": "home-*",
      "brews": ["steam", "discord"],
      "casks": ["spotify"],
      "variables": {
        "PERSONAL": "true",
        "EMAIL": "you@personal.com"
      }
    }
  ]
}
```

**Glob pattern support:**
- `work-*` - Matches anything starting with "work-"
- `*-laptop` - Matches anything ending with "-laptop"
- `*macbook*` - Matches anything containing "macbook"

**How it works:**
1. Run `dotfiles conditionals apply` on any machine
2. Automatically detects hostname and OS
3. Adds matching packages to your config
4. Run `dotfiles install` to install them
5. Same dotfiles repo, different packages per machine!

## 📦 Package Dependencies

Define dependencies between packages:

```bash
# Show dependencies
dotfiles depends show neovim

# Show full tree
dotfiles depends tree neovim

# Resolve install order
dotfiles depends resolve
```

Example config:
```json
{
  "package_dependencies": {
    "neovim": ["python3", "node"],
    "tmux": ["zsh"]
  }
}
```

## 👀 Watch Mode

Auto-commit and sync changes:

```bash
# Start watching (30 second interval)
dotfiles watch

# Custom interval
dotfiles watch --interval 60

# Stop watching
dotfiles watch stop
```

## 🤝 Configuration Sharing & Templates

### 📚 Using Templates

```bash
# Browse available templates
dotfiles templates list

# Preview a template
dotfiles templates show web-dev

# Apply a template
dotfiles clone template:web-dev

# Discover community templates
dotfiles templates discover --search "web-dev"
```

**Built-in Templates:**
- `web-dev` - Web development (Node.js, Python, Docker)
- `mobile-dev` - iOS/Android (Flutter, React Native)
- `data-science` - Python, R, Jupyter, analytics
- `devops` - Kubernetes, Terraform, cloud tools
- `minimal` - Essential tools only

### 📤 Share Your Configuration

```bash
# Share via GitHub Gist
dotfiles share gist --name="My Setup" --description="Full-stack config"

# Share privately
dotfiles share gist --name="My Config" --private

# Export to file
dotfiles share file my-config.json
```

### 📥 Import Shared Configurations

```bash
# From GitHub Gist URL
dotfiles clone https://gist.github.com/user/gist-id

# From local file
dotfiles clone my-config.json

# Preview before importing
dotfiles clone https://gist.github.com/user/gist-id --preview

# Merge with existing
dotfiles clone https://gist.github.com/user/gist-id --merge
```

## 💡 Usage Examples

### Adding Packages
```bash
# Add brew packages (default)
dotfiles add git curl wget

# Add casks (GUI applications)
dotfiles add --type=cask visual-studio-code firefox

# Add taps (repositories)
dotfiles add --type=tap homebrew/cask-fonts

# Add Stow packages (dotfiles)
dotfiles add --type=stow vim zsh tmux
```

### Managing Dotfiles with Stow
```bash
# Preview changes before stowing
dotfiles stow-diff vim zsh

# Create symlinks
dotfiles stow vim zsh tmux

# Remove symlinks
dotfiles unstow vim

# Restow (remove and recreate)
dotfiles restow vim

# Custom directories
dotfiles stow --dir=/path/to/dotfiles --target=~ vim
```

### Working with Private Files
```bash
# Move private file to private directory
mv ~/.env-private ~/.dotfiles/private/.env-private

# Create symlink in stow package
dotfiles private shell .env-private

# Stow the package
dotfiles stow shell

# Creates: ~/.env-private -> ~/.dotfiles/stow/shell/.env-private -> ~/.dotfiles/private/.env-private
```

### Backup and Restore
```bash
# Create timestamped backup
dotfiles backup

# Custom backup location
dotfiles backup --dir ~/backups

# Custom name
dotfiles backup my-backup.json

# Restore from backup
dotfiles restore my-backup.json
```

## 📁 Configuration

Your configuration is stored at `~/.dotfiles/config.json`:

```json
{
  "brews": ["git", "lazygit", "neovim"],
  "casks": ["visual-studio-code"],
  "taps": ["homebrew/cask-fonts"],
  "stow": ["vim", "zsh", "tmux"],
  "roles": {
    "web-dev": {
      "name": "Web Development",
      "brews": ["node", "python", "docker"],
      "casks": ["visual-studio-code"]
    }
  },
  "conditionals": [
    {
      "hostname": "work-*",
      "brews": ["work-tool"]
    }
  ],
  "package_dependencies": {
    "neovim": ["python3", "node"]
  },
  "variables": {
    "EMAIL": "user@example.com",
    "NAME": "John Doe"
  }
}
```

### Directory Structure

```
~/.dotfiles/
├── config.json          # Package configuration
├── .gitignore           # Excludes private/
├── stow/                # Stow packages
│   ├── vim/
│   │   ├── .vimrc
│   │   └── .vim/
│   ├── zsh/
│   │   ├── .zshrc
│   │   └── .zsh/
│   └── tmux/
│       └── .tmux.conf
├── private/             # Private files (git-ignored)
│   ├── .env-private     # Secrets (managed by 'dotfiles secrets')
│   ├── .env-private.template  # Template for new machines
│   └── .ssh/
├── snapshots/           # Configuration snapshots
├── backups/             # Timestamped backups
└── .git/                # Version control
```

## 🌍 Cross-Platform Support

This tool works seamlessly across multiple platforms with **one configuration file**:

- **macOS**: Uses Homebrew
- **Arch Linux**: Uses pacman/yay
- **Debian/Ubuntu**: Uses apt
- **RHEL/Fedora**: Uses yum/dnf

The same `config.json` works everywhere! Casks are treated as regular packages on Linux.

## 🔧 Shell Completions

Install completions for tab completion:

```bash
# Generate for your shell
dotfiles completion bash > /usr/local/etc/bash_completion.d/dotfiles
dotfiles completion zsh > "${fpath[1]}/_dotfiles"
dotfiles completion fish > ~/.config/fish/completions/dotfiles.fish

# Or auto-install
dotfiles completion install bash
dotfiles completion install zsh
dotfiles completion install fish
```

## 📚 Full Command Reference

### Getting Started
| Command | Description |
|---------|-------------|
| `dotfiles init` | Initialize new configuration |
| `dotfiles onboard` | Complete developer onboarding |
| `dotfiles tutorial` | Interactive tutorial |
| `dotfiles doctor` | Health check with diagnostics |
| `dotfiles doctor --fix` | Auto-fix common issues |

### Package Management
| Command | Description |
|---------|-------------|
| `dotfiles add <packages>` | Add packages to config |
| `dotfiles remove <packages>` | Remove packages from config |
| `dotfiles install` | Install all configured packages |
| `dotfiles update` | Update all packages |
| `dotfiles status` | Check installation status |
| `dotfiles diff` | Compare config vs installed |
| `dotfiles scan` | Scan and add installed packages |

### Dotfiles Management
| Command | Description |
|---------|-------------|
| `dotfiles stow <packages>` | Create symlinks |
| `dotfiles unstow <packages>` | Remove symlinks |
| `dotfiles restow <packages>` | Recreate symlinks |
| `dotfiles stow-diff [packages]` | Preview stow changes |
| `dotfiles list` | List all packages |

### Configuration
| Command | Description |
|---------|-------------|
| `dotfiles lint` | Validate configuration |
| `dotfiles lint --fix` | Auto-fix config issues |
| `dotfiles snapshot create` | Create snapshot |
| `dotfiles snapshot list` | List snapshots |
| `dotfiles snapshot restore <id>` | Restore snapshot |
| `dotfiles backup [file]` | Create timestamped backup |
| `dotfiles restore <file>` | Restore from backup |

### Advanced
| Command | Description |
|---------|-------------|
| `dotfiles migrate detect` | Detect other tools |
| `dotfiles migrate from <tool>` | Migrate from other tool |
| `dotfiles shell setup` | Interactive shell setup |
| `dotfiles conditionals apply` | Apply OS-specific config |
| `dotfiles conditionals test` | Test which conditionals match |
| `dotfiles secrets init` | Initialize secrets management |
| `dotfiles secrets add <NAME>` | Add secret (prompts for value) |
| `dotfiles secrets list` | List secret names |
| `dotfiles secrets template` | Generate template for new machines |
| `dotfiles template render` | Render template variables |
| `dotfiles depends show <pkg>` | Show dependencies |
| `dotfiles watch` | Auto-commit changes |
| `dotfiles roles list` | List package roles |
| `dotfiles completion <shell>` | Generate completions |
| `dotfiles bootstrap` | Generate bootstrap script |
| `dotfiles check` | Verify setup |
| `dotfiles uninstall` | Remove dotfiles setup |

## 🛠 Development

### Prerequisites
- Go 1.25.1 or later
- Package manager for your OS (Homebrew, pacman, apt, yum)
- GNU Stow: `brew install stow` or `sudo pacman -S stow`

### Building
```bash
go build -o dotfiles
```

### Project Structure
```
dotfiles-cli/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── init.go            # Initialize
│   ├── onboard.go         # Onboarding
│   ├── add.go             # Add packages
│   ├── install.go         # Install packages
│   ├── stow.go            # Stow integration
│   ├── doctor.go          # Health checks
│   ├── snapshot.go        # Snapshots
│   ├── lint.go            # Config validation
│   ├── migrate.go         # Migration tools
│   ├── shell.go           # Shell setup
│   ├── completion.go      # Shell completions
│   ├── conditionals.go    # Conditional configs
│   ├── template.go        # Template variables
│   ├── depends.go         # Dependencies
│   ├── watch.go           # Watch mode
│   └── ...                # More commands
├── internal/
│   ├── config/            # Config management
│   └── pkgmanager/        # Package manager abstraction
└── main.go                # Entry point
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

MIT License

---

**Made with ❤️ for developers who want their environment setup to just work™**
