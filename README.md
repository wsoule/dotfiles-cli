# 🛠 Dotfiles Manager

> A modern, fast, and beautiful dotfiles management system built in Go

[![Release](https://img.shields.io/github/v/release/wsoule/new-dotfiles)](https://github.com/wsoule/new-dotfiles/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/wsoule/new-dotfiles)](https://goreportcard.com/report/github.com/wsoule/new-dotfiles)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Effortlessly manage your development environment with a beautiful CLI interface, preset configurations, and automated setup.

## Features

- 🌐 **Modern Web Interface** - Beautiful browser-based setup wizard
- 🚀 **Interactive Setup Wizard** - Easy configuration with preset support
- 📦 **Package Management** - Automated Homebrew, npm, and system package installation
- 🎨 **System Configuration** - macOS system defaults and preferences
- 🔗 **Dotfiles Management** - GNU Stow-based dotfile installation
- ⚙️ **Development Environment** - Language and framework-specific configurations
- 📋 **Configuration Management** - JSON-based configuration with validation
- 🤝 **Configuration Sharing** - Share and import configurations easily
- ✋ **Opt-in Everything** - Nothing is installed without your explicit consent
- 📱 **Cross-Platform UI** - Works in browser, terminal, or headless environments

## 📦 Installation

### Option 1: One-line installer (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/wsoule/new-dotfiles/main/install.sh | bash
```

### Option 2: Homebrew (macOS/Linux)
```bash
# Add the tap (replace with your GitHub username)
brew tap wyatsoule/tap
brew install dotfiles
```

### Option 3: Download from releases
1. Go to the [releases page](https://github.com/wsoule/new-dotfiles/releases)
2. Download the binary for your platform
3. Extract and move to your PATH:
   ```bash
   sudo mv dotfiles /usr/local/bin/
   chmod +x /usr/local/bin/dotfiles
   ```

### Option 4: Go install
```bash
go install github.com/wsoule/new-dotfiles@latest
```

### Option 5: Build from source
```bash
git clone https://github.com/wsoule/new-dotfiles.git
cd dotfiles
make build
sudo mv dotfiles /usr/local/bin/
```

## 🚀 Quick Start

1. **Run the modern web-based setup wizard:**
   ```bash
   dotfiles setup
   ```
   This opens a beautiful web interface in your browser for easy configuration!

2. **Install your configuration:**
   ```bash
   dotfiles install
   ```

3. **That's it!** Your development environment is now configured.

### 🌐 **Modern Setup Experience**

The setup wizard now features:
- **🎨 Beautiful web interface** - Auto-opens in your browser
- **📱 Responsive design** - Works on desktop, tablet, and mobile
- **⚡ Live preview** - See your configuration as you build it
- **🔄 Auto-save** - Never lose your progress
- **📋 Step-by-step guidance** - Clear progress indicators

For terminal-only environments, use: `dotfiles setup --cli`

## 🖼 Screenshots

When you run `dotfiles`, you'll see a beautiful banner and interactive interface:

```
                              🛠  DOTFILES MANAGER

  A modern dotfiles management system built in Go
  Configure your development environment with ease

Usage:
  dotfiles [command]

Available Commands:
  setup       Interactive setup wizard for dotfiles configuration
  install     Install dotfiles and configure system
  config      Manage dotfiles configuration
```

## 💡 Usage Examples

### First time setup with preset
```bash
# Use a JavaScript development preset
dotfiles setup --preset javascript-dev

# Quick setup with minimal prompts
dotfiles setup --quick
```

### Install with specific options
```bash
# Preview what would be installed
dotfiles install --dry-run

# Install but skip macOS configuration
dotfiles install --skip-macos
```

### Configuration management
```bash
# Show current configuration
dotfiles config show

# Show just a summary
dotfiles config show --summary

# Validate your configuration
dotfiles config validate
```

### Configuration sharing
```bash
# Export your config for sharing (removes personal info)
dotfiles share export my-config.json

# Import someone's shared configuration
dotfiles share import their-config.json

# Validate a shared configuration
dotfiles share validate config.json
```

## 🎯 Opt-in Philosophy

This dotfiles manager follows a **strict opt-in philosophy** - nothing is installed or configured without your explicit consent:

### ✅ **What's Opt-in**
- **All packages and applications** - Every brew, cask, and npm package
- **All system modifications** - macOS defaults, dock settings, security preferences
- **All development tools** - Languages, frameworks, CLIs, git tools
- **All shell enhancements** - Themes, plugins, aliases
- **All fonts and visual elements** - Nerd fonts, color schemes

### 🚫 **What's NOT automatic**
- No "essential" or "recommended" packages forced on you
- No system changes without permission
- No dotfiles copied without selection
- No personal information shared

### 📋 **How it works**
1. **Setup wizard** asks what you want to enable
2. **Configuration file** stores your choices explicitly
3. **Installation** only applies what you've selected
4. **Sharing** removes personal info automatically

Even `git` and `stow` are opt-in - though you'll likely want them for dotfiles management!

## 📚 Commands Reference

### 🔧 Setup Commands
| Command | Description |
|---------|-------------|
| `dotfiles setup` | Launch modern web-based setup wizard |
| `dotfiles setup --cli` | Use terminal-based setup wizard |
| `dotfiles setup --preset <name>` | Use a preset configuration |
| `dotfiles setup --port <port>` | Specify port for web interface |
| `dotfiles setup --force` | Force setup even if config exists |

### 🚀 Installation Commands
| Command | Description |
|---------|-------------|
| `dotfiles install` | Install dotfiles and configure system |
| `dotfiles install --dry-run` | Preview what would be installed |
| `dotfiles install --skip-homebrew` | Skip Homebrew installation |
| `dotfiles install --skip-macos` | Skip macOS configuration |
| `dotfiles install --skip-stow` | Skip dotfiles installation |

### ⚙️ Configuration Management
| Command | Description |
|---------|-------------|
| `dotfiles config show` | Display current configuration |
| `dotfiles config show --summary` | Show configuration summary |
| `dotfiles config show --json` | Output configuration as JSON |
| `dotfiles config validate` | Validate configuration |
| `dotfiles config get <key>` | Get specific configuration value |
| `dotfiles config set <key> <value>` | Set configuration value |

### 🤝 Sharing Commands
| Command | Description |
|---------|-------------|
| `dotfiles share export [file]` | Export configuration for sharing |
| `dotfiles share import <file>` | Import a shared configuration |
| `dotfiles share validate <file>` | Validate a configuration file |

### 📋 Other Commands
| Command | Description |
|---------|-------------|
| `dotfiles --help` | Show help information |
| `dotfiles --version` | Show version information |

## Configuration

The configuration is stored in JSON format and includes:

- **Personal Information** - Name, email, preferred editor
- **System Preferences** - Dark mode, dock settings, finder preferences
- **Development Environment** - Programming languages, frameworks, tools
- **Package Management** - Extra brew packages, casks, npm globals
- **Installation Options** - What components to install

### Configuration File Locations

- Default: `~/.dotfiles/config.json`
- Custom: Use `--config <path>` flag

### Configuration Schema

```json
{
  "personal": {
    "name": "Your Name",
    "email": "your.email@example.com",
    "editor": "nvim"
  },
  "system": {
    "appearance": {
      "dark_mode": true,
      "enable_24_hour_time": true
    },
    "dock": {
      "autohide": true,
      "position": "bottom",
      "tile_size": 50
    }
  },
  "development": {
    "languages": {
      "javascript": true,
      "python": true,
      "go": true
    },
    "shell": {
      "theme": "powerlevel10k",
      "terminal_theme": "dark"
    }
  },
  "packages": {
    "extra_brews": ["wget", "jq"],
    "extra_casks": ["visual-studio-code"],
    "npm_globals": ["nx", "typescript"]
  }
}
```

## Presets

Presets allow you to quickly configure your environment based on common setups:

- Copy existing presets from the `presets/` directory
- Create custom presets by saving configuration files
- Load presets during setup with `--preset <name>`

## Migration from Shell Version

This Go version provides the same functionality as the original shell-based dotfiles manager:

- **Setup Wizard** - Converted from `setup.sh`
- **Installation** - Converted from `core/install.sh`
- **Configuration Management** - Converted from `scripts/config-manager.sh`
- **UI Components** - Interactive prompts and menus

### Key Improvements

- **Better Error Handling** - Structured error reporting and recovery
- **Faster Execution** - Compiled binary vs shell script interpretation
- **Cross-Platform Support** - Easier to extend beyond macOS
- **Type Safety** - Configuration validation and type checking
- **Modularity** - Clean separation of concerns and testability

## Building and Development

### Prerequisites

- Go 1.25.1 or later
- macOS (for full functionality)

### Building

```bash
# Build for current platform
go build -o dotfiles

# Build for multiple platforms
GOOS=darwin GOARCH=amd64 go build -o dotfiles-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o dotfiles-darwin-arm64
```

### Running Tests

```bash
go test ./...
```

### Project Structure

```
Go_Dotfiles/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command and configuration
│   ├── setup.go        # Setup wizard command
│   ├── install.go      # Installation command
│   └── config.go       # Configuration management
├── internal/           # Internal packages
│   ├── config/         # Configuration management
│   ├── installer/      # Installation logic
│   └── ui/             # User interface components
├── pkg/                # Public packages
│   ├── brew/           # Homebrew utilities
│   ├── macos/          # macOS system configuration
│   └── stow/           # GNU Stow integration
├── config/             # Configuration files
├── presets/            # Preset configurations
├── Brewfile            # Homebrew package definitions
└── main.go             # Application entry point
```

## 🛠 Development

### Prerequisites
- Go 1.25.1 or later
- macOS (for full functionality)

### Building from Source
```bash
# Clone the repository
git clone https://github.com/wsoule/new-dotfiles.git
cd dotfiles

# Install dependencies
make deps

# Build the binary
make build

# Run tests
make test

# Install locally
make install
```

### Available Make Targets
| Target | Description |
|--------|-------------|
| `make build` | Build the binary |
| `make test` | Run tests |
| `make install` | Install to $GOPATH/bin |
| `make clean` | Clean build artifacts |
| `make lint` | Run linters |
| `make fmt` | Format code |
| `make release-test` | Test release process |

### Project Structure
```
dotfiles/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command and configuration
│   ├── setup.go        # Setup wizard command
│   ├── install.go      # Installation command
│   └── config.go       # Configuration management
├── internal/           # Internal packages
│   ├── config/         # Configuration management
│   ├── installer/      # Installation logic
│   └── ui/             # User interface components
├── pkg/                # Public packages
│   ├── brew/           # Homebrew utilities
│   ├── macos/          # macOS system configuration
│   └── stow/           # GNU Stow integration
├── config/             # Configuration files
├── presets/            # Preset configurations
├── .github/workflows/  # GitHub Actions
├── Brewfile.template   # Brewfile template for generation
├── Brewfile.example    # Example Brewfile configuration
├── Makefile           # Build automation
├── .goreleaser.yml    # Release configuration
└── main.go            # Application entry point
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Originally based on a shell-based dotfiles system
- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- UI enhanced with [PTerm](https://github.com/pterm/pterm) for beautiful terminal output
- Distributed with [GoReleaser](https://goreleaser.com/) for multi-platform releases
