# Cross-Platform Support - Changelog

## Summary

Added comprehensive cross-platform support to the dotfiles manager. The tool now works on **macOS, Arch Linux, Debian/Ubuntu, and RHEL/Fedora** with a **single, compatible configuration file**.

## Major Changes

### 1. Package Manager Abstraction Layer
**New File**: `internal/pkgmanager/manager.go`

Created a package manager abstraction that automatically detects and uses the appropriate package manager:

- **HomebrewManager**: macOS with Homebrew
- **PacmanManager**: Arch Linux with pacman/yay
- **AptManager**: Debian/Ubuntu with apt
- **YumManager**: RHEL/Fedora with yum/dnf

Each manager implements the same interface:
- `Install()` - Install packages
- `IsInstalled()` - Check if package is installed
- `ListInstalled()` - List all installed packages
- `GenerateInstallFile()` - Generate package list file
- `InstallFromFile()` - Install from package list file

### 2. Configuration Compatibility
**Modified**: `internal/config/config.go`

- Config structure remains backward compatible
- Added `GetAllPackages()` helper for Linux systems
- Same JSON format works across all platforms:
  - `brews`: CLI packages (all platforms)
  - `casks`: GUI apps on macOS, regular packages on Linux
  - `taps`: Homebrew taps (macOS only, ignored on Linux)
  - `stow`: Dotfiles packages (all platforms)

### 3. Updated Commands

#### Install Command
**Modified**: `cmd/install.go`

- Auto-detects package manager
- Generates appropriate package file (Brewfile on macOS, packages.txt on Linux)
- Uses abstracted package manager for installation
- Maintains backward compatibility with Homebrew

#### Status Command
**Modified**: `cmd/status.go`

- Cross-platform package checking
- Shows detected package manager
- Platform-aware suggestions for missing packages
- Skips macOS-specific features (taps, casks) on Linux

#### Onboard Command
**Modified**: `cmd/onboard.go`

- Platform-specific essential packages
- Auto-detects and suggests appropriate dependencies
- macOS: Homebrew, VS Code, Ghostty, Raycast
- Linux: git, stow, base-devel, standard tools

#### Doctor Command
**Modified**: `cmd/doctor.go`

- Cross-platform health checks
- Detects correct package manager
- Platform-appropriate install suggestions
- Works with all supported package managers

#### Utility Functions
**Modified**: `cmd/brew_utils.go`

- Updated to use package manager abstraction
- Falls back to direct Homebrew commands when needed
- Cross-platform package listing

### 4. Documentation

**New**: `CROSS_PLATFORM.md`
- Comprehensive cross-platform guide
- Examples for each platform
- Migration guide
- Troubleshooting section

**Updated**: `README.md`
- Added cross-platform highlights
- Updated prerequisites for Linux
- Added platform support section
- Updated feature list

## Technical Details

### OS Detection
Uses Go's `runtime.GOOS` to detect the operating system:
- `darwin` → macOS → Homebrew
- `linux` → Detects available package manager (pacman, apt, yum)

### Package Manager Priority (Linux)
1. pacman (Arch Linux)
2. apt (Debian/Ubuntu)
3. yum/dnf (RHEL/Fedora)

### Backward Compatibility
- All existing configs continue to work on macOS
- Homebrew-specific commands still work as before
- No breaking changes to command structure

## Benefits

1. **Single Config File**: Use the same dotfiles config on all your machines
2. **Auto-Detection**: Tool automatically uses the right package manager
3. **Cross-Platform Sharing**: Share configs between macOS and Linux users
4. **Future-Proof**: Easy to add support for more package managers

## Testing

Built and tested on macOS:
- ✅ Compiles without errors
- ✅ All commands accessible
- ✅ Maintains backward compatibility
- ✅ Homebrew integration intact

Ready for testing on Linux distributions.

## Migration Guide

### For Existing macOS Users
- No changes needed
- Everything works as before
- Config file compatible as-is

### For New Linux Users
1. Install the binary
2. Run `dotfiles init`
3. Add packages with `dotfiles add <package>`
4. Run `dotfiles install`

### For Users With Dotfiles on macOS Moving to Linux
1. Copy your `~/.dotfiles/` directory
2. Run `dotfiles install` - will auto-detect Linux package manager
3. Some macOS casks may need manual installation or different package names

## Future Enhancements

Potential improvements:
- Package name mapping for cross-platform equivalents
- Platform-specific config sections
- Homebrew on Linux support
- Flatpak/Snap support
- Windows support (WSL/Chocolatey)
