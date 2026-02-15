# macfig

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform: macOS](https://img.shields.io/badge/Platform-macOS-lightgrey.svg)](https://github.com/lu-zhengda/macfig)
[![Homebrew](https://img.shields.io/badge/Homebrew-lu--zhengda/tap-orange.svg)](https://github.com/lu-zhengda/homebrew-tap)

macOS hidden defaults manager — browse, search, and apply system preferences with presets from the terminal.

## Install

```bash
brew tap lu-zhengda/tap
brew install macfig
```

## Usage

```
$ macfig list
Categories:

  🚢 Dock  (10 settings)
  📁 Finder  (9 settings)
  📷 Screenshots  (5 settings)
  💻 Keyboard  (6 settings)
  🤚 Trackpad  (3 settings)
  🎬 Animations  (2 settings)
  🔒 Privacy  (2 settings)
  🔧 Misc  (4 settings)

Use 'macfig list <category>' to view settings.

$ macfig list dock
🚢 Dock

  Auto-hide Dock                      0 -> recommended: true
    Automatically hide and show the Dock  (com.apple.dock autohide)
  Auto-hide delay                     0 -> recommended: 0
    Delay before Dock auto-hides (seconds)  (com.apple.dock autohide-delay)
  Auto-hide animation speed           0.2 -> recommended: 0.2
    Speed of the Dock show/hide animation (seconds)  (com.apple.dock autohide-time-modifier)
  Show recent apps                    0 -> recommended: false
    Show recent applications in the Dock  (com.apple.dock show-recents)
  Bounce animation                    0 -> recommended: false
    Animate opening applications (bounce)  (com.apple.dock launchanim)
```

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `list` | List preset categories | `macfig list` |
| `list <category>` | Show settings in category with current values | `macfig list dock` |
| `get <domain> <key>` | Read a specific defaults value | `macfig get com.apple.dock autohide` |
| `set <domain> <key> <value>` | Write a defaults value | `macfig set com.apple.dock autohide true` |
| `preset <name>` | Apply a named preset bundle | `macfig preset dev-setup` |
| `backup` | Backup all tracked settings to JSON | `macfig backup` |
| `restore <file>` | Restore from backup file | `macfig restore ~/.config/macfig/backup.json` |
| `reset <domain> <key>` | Delete a key (restore macOS default) | `macfig reset com.apple.dock autohide` |
| `search <query>` | Search keys across all domains | `macfig search animation` |

### Built-in Presets

| Preset | Effect |
|--------|--------|
| `dock-speed` | Faster Dock animations and autohide delay |
| `finder-dev` | Show hidden files, extensions, path bar, status bar |
| `keyboard-fast` | Faster key repeat rate and shorter delay |
| `no-animations` | Disable most macOS animations |
| `privacy` | Disable analytics, Siri suggestions, ad tracking |
| `dev-setup` | Combined developer-friendly defaults |

## Safety

- **Always backup first** — run `macfig backup` before applying presets or bulk changes
- **Use presets over raw set** — presets handle type inference and process restart automatically
- **Use search to discover** — run `macfig search <term>` instead of guessing domain/key pairs
- **Revert with restore** — if something breaks, `macfig restore <backup-file>` reverts all changes

## TUI

Launch `macfig` without arguments for interactive mode — browse categories, toggle settings, and preview changes.

## License

[MIT](LICENSE)
