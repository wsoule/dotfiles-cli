# Cross-Platform Support

The dotfiles manager now supports both **macOS** and **Linux** (Arch, Debian/Ubuntu, RHEL/Fedora) with a **single, compatible configuration file**.

## How It Works

### Package Manager Detection

The tool automatically detects your operating system and uses the appropriate package manager:

- **macOS**: Homebrew (`brew`)
- **Arch Linux**: pacman/yay
- **Debian/Ubuntu**: apt
- **RHEL/Fedora**: yum/dnf

### Unified Configuration

Your `~/.dotfiles/config.json` file works across all platforms. The same config structure is used everywhere:

```json
{
  "brews": [
    "git",
    "curl",
    "wget",
    "tree",
    "jq",
    "stow"
  ],
  "casks": [
    "visual-studio-code",
    "firefox",
    "slack"
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

### Platform-Specific Behavior

#### macOS (Homebrew)
- **brews**: Installed as Homebrew formulas
- **casks**: Installed as GUI applications
- **taps**: Additional Homebrew repositories

#### Linux (pacman/apt/yum)
- **brews**: Installed as regular packages
- **casks**: Installed as regular packages (no distinction)
- **taps**: Ignored (macOS-only concept)

## Example: Same Config, Different Platforms

### Config File (`~/.dotfiles/config.json`)
```json
{
  "brews": [
    "git",
    "neovim",
    "tmux",
    "ripgrep",
    "fd",
    "bat"
  ],
  "casks": [
    "visual-studio-code",
    "alacritty"
  ],
  "stow": [
    "vim",
    "zsh"
  ]
}
```

### On macOS
```bash
$ dotfiles install
✓ Generated package list at: ./Brewfile
📦 Installing packages with homebrew...
Installing packages...
✅ git
✅ neovim
✅ tmux
✅ ripgrep
✅ fd
✅ bat

Installing casks/applications...
✅ visual-studio-code
✅ alacritty

✅ Installation complete!
```

### On Arch Linux
```bash
$ dotfiles install
✓ Generated package list at: ./packages.txt
📦 Installing packages with pacman...
Installing packages...
✅ git
✅ neovim
✅ tmux
✅ ripgrep
✅ fd
✅ bat
✅ visual-studio-code  # Installed from repos/AUR
✅ alacritty           # Installed from repos/AUR

✅ Installation complete!
```

## Package Name Mapping

Most package names are the same across platforms, but some differ:

| Common Name | macOS (brew) | Arch (pacman) | Debian (apt) |
|-------------|--------------|---------------|--------------|
| git         | git          | git           | git          |
| vim         | vim          | vim           | vim          |
| VS Code     | visual-studio-code | code | code |
| Firefox     | firefox      | firefox       | firefox-esr  |

**Tip**: Use the common package name in your config. If a package has a different name on a specific platform, you may need to adjust.

## Commands That Work Cross-Platform

All commands work on both macOS and Linux:

```bash
# Initialize dotfiles
dotfiles init

# Add packages (works on any platform)
dotfiles add git curl wget
dotfiles add --type=cask firefox    # Becomes a regular package on Linux

# Install packages
dotfiles install

# Check status
dotfiles status

# Stow dotfiles (cross-platform)
dotfiles stow vim zsh tmux
```

## Platform-Specific Features

### macOS Only
- **Taps**: Additional Homebrew repositories
- **Casks**: Distinction between CLI tools and GUI apps
- **Brewfile**: Native Homebrew bundle file format

### Linux Advantages
- **Native package managers**: Uses system package manager (faster, better integrated)
- **AUR support**: On Arch, automatically uses `yay` if available for AUR packages

## Migration Guide

### Moving from macOS to Linux

1. Copy your `~/.dotfiles/` directory to the Linux machine
2. Run `dotfiles install` - it will automatically use the Linux package manager
3. Note: Some macOS-specific apps in `casks` might not be available

### Moving from Linux to macOS

1. Copy your `~/.dotfiles/` directory to the Mac
2. Install Homebrew if not already installed
3. Run `dotfiles install` - it will use Homebrew

## Developer Onboarding

The `onboard` command is platform-aware:

```bash
# On macOS
dotfiles onboard
# → Installs Homebrew, git, stow, VS Code, etc.

# On Arch Linux
dotfiles onboard
# → Uses pacman/yay, installs git, stow, base-devel, etc.
```

## Troubleshooting

### Package Not Found

If a package has a different name on your platform:

1. Find the correct package name for your platform
2. Update your config.json with the platform-specific name
3. Or maintain separate configs with platform-specific overrides

### Package Manager Not Detected

```bash
$ dotfiles doctor
# Shows which package manager was detected and any issues
```

### Mixed Environments

If you use both macOS and Linux, you can:

1. Use a **single config** with common packages
2. Use **git branches** for platform-specific packages
3. Use **hooks** to run platform-specific setup

## Advanced: Platform-Specific Packages

Use hooks to handle platform-specific packages:

```json
{
  "brews": ["git", "curl"],
  "hooks": {
    "pre_install": [
      "if [ $(uname) = 'Linux' ]; then echo 'Installing linux-specific packages...'; fi"
    ]
  }
}
```

## Contributing

When adding features, ensure they work on both macOS and Linux. Test on:
- macOS with Homebrew
- Arch Linux with pacman
- Debian/Ubuntu with apt (if possible)
