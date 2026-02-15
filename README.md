# macfig

macOS hidden defaults manager — browse, search, and apply system preferences from the terminal.

## Install

```bash
brew tap lu-zhengda/tap
brew install macfig
```

## Quick Start

```bash
macfig           # Launch interactive TUI
macfig --help    # Show all commands
```

## Commands

| Command   | Description                                    |
|-----------|------------------------------------------------|
| `list`    | List preset categories or settings             |
| `get`     | Get a setting value                            |
| `set`     | Set a setting value                            |
| `preset`  | Apply a named preset (dock-speed, no-animations, etc.) |
| `backup`  | Backup current settings                        |
| `restore` | Restore settings from backup                   |
| `reset`   | Reset settings to defaults                     |
| `search`  | Search for settings                            |

## TUI

Launch without arguments for interactive mode. Browse categories, search settings, and apply presets with a keyboard-driven interface.

<!-- Screenshot placeholder: ![macfig TUI](docs/tui.png) -->

## License

MIT
