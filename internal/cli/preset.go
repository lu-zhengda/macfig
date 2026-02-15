package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
)

var presetCmd = &cobra.Command{
	Use:   "preset <name>",
	Short: "Apply a named preset",
	Long: `Apply a group of recommended settings.

Available presets:
  dock-speed       Faster Dock animations (auto-hide delay 0, animation 0.2)
  finder-dev       Developer-friendly Finder (path bar, extensions, POSIX path)
  keyboard-fast    Faster keyboard repeat (KeyRepeat 2, InitialKeyRepeat 15)
  no-animations    Disable window animations
  privacy          Enable password on wake, no delay
  dev-setup        Apply all developer-friendly settings`,
	Args: cobra.ExactArgs(1),
	RunE: runPreset,
}

// presetDefinition maps preset names to their settings.
type presetDefinition struct {
	Name     string
	Settings []presetEntry
}

type presetEntry struct {
	Domain  string
	Key     string
	Value   interface{}
	Type    defaults.ValueType
	Restart string
}

var presets = map[string]presetDefinition{
	"dock-speed": {
		Name: "Dock Speed",
		Settings: []presetEntry{
			{"com.apple.dock", "autohide", true, defaults.TypeBool, "Dock"},
			{"com.apple.dock", "autohide-delay", 0.0, defaults.TypeFloat, "Dock"},
			{"com.apple.dock", "autohide-time-modifier", 0.2, defaults.TypeFloat, "Dock"},
			{"com.apple.dock", "launchanim", false, defaults.TypeBool, "Dock"},
		},
	},
	"finder-dev": {
		Name: "Finder Developer Mode",
		Settings: []presetEntry{
			{"com.apple.finder", "ShowPathbar", true, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "ShowStatusBar", true, defaults.TypeBool, "Finder"},
			{"NSGlobalDomain", "AppleShowAllExtensions", true, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "FXEnableExtensionChangeWarning", false, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "_FXShowPosixPathInTitle", true, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "FXDefaultSearchScope", "SCcf", defaults.TypeString, "Finder"},
			{"com.apple.finder", "_FXSortFoldersFirst", true, defaults.TypeBool, "Finder"},
		},
	},
	"keyboard-fast": {
		Name: "Fast Keyboard",
		Settings: []presetEntry{
			{"NSGlobalDomain", "KeyRepeat", 2, defaults.TypeInt, ""},
			{"NSGlobalDomain", "InitialKeyRepeat", 15, defaults.TypeInt, ""},
			{"NSGlobalDomain", "ApplePressAndHoldEnabled", false, defaults.TypeBool, ""},
		},
	},
	"no-animations": {
		Name: "No Animations",
		Settings: []presetEntry{
			{"NSGlobalDomain", "NSAutomaticWindowAnimationsEnabled", false, defaults.TypeBool, ""},
			{"NSGlobalDomain", "NSWindowResizeTime", 0.001, defaults.TypeFloat, ""},
			{"com.apple.dock", "launchanim", false, defaults.TypeBool, "Dock"},
		},
	},
	"privacy": {
		Name: "Privacy",
		Settings: []presetEntry{
			{"com.apple.screensaver", "askForPassword", true, defaults.TypeBool, ""},
			{"com.apple.screensaver", "askForPasswordDelay", 0, defaults.TypeInt, ""},
		},
	},
	"dev-setup": {
		Name: "Developer Setup",
		Settings: []presetEntry{
			// Dock
			{"com.apple.dock", "autohide", true, defaults.TypeBool, "Dock"},
			{"com.apple.dock", "autohide-delay", 0.0, defaults.TypeFloat, "Dock"},
			{"com.apple.dock", "autohide-time-modifier", 0.2, defaults.TypeFloat, "Dock"},
			{"com.apple.dock", "show-recents", false, defaults.TypeBool, "Dock"},
			{"com.apple.dock", "launchanim", false, defaults.TypeBool, "Dock"},
			// Finder
			{"com.apple.finder", "ShowPathbar", true, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "ShowStatusBar", true, defaults.TypeBool, "Finder"},
			{"NSGlobalDomain", "AppleShowAllExtensions", true, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "FXEnableExtensionChangeWarning", false, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "_FXShowPosixPathInTitle", true, defaults.TypeBool, "Finder"},
			{"com.apple.finder", "FXDefaultSearchScope", "SCcf", defaults.TypeString, "Finder"},
			{"com.apple.finder", "_FXSortFoldersFirst", true, defaults.TypeBool, "Finder"},
			// Keyboard
			{"NSGlobalDomain", "KeyRepeat", 2, defaults.TypeInt, ""},
			{"NSGlobalDomain", "InitialKeyRepeat", 15, defaults.TypeInt, ""},
			{"NSGlobalDomain", "ApplePressAndHoldEnabled", false, defaults.TypeBool, ""},
			{"NSGlobalDomain", "NSAutomaticSpellingCorrectionEnabled", false, defaults.TypeBool, ""},
			{"NSGlobalDomain", "NSAutomaticCapitalizationEnabled", false, defaults.TypeBool, ""},
			// Animations
			{"NSGlobalDomain", "NSAutomaticWindowAnimationsEnabled", false, defaults.TypeBool, ""},
			{"NSGlobalDomain", "NSWindowResizeTime", 0.001, defaults.TypeFloat, ""},
			// Privacy
			{"com.apple.screensaver", "askForPassword", true, defaults.TypeBool, ""},
			{"com.apple.screensaver", "askForPasswordDelay", 0, defaults.TypeInt, ""},
			// Misc
			{"NSGlobalDomain", "NSNavPanelExpandedStateForSaveMode", true, defaults.TypeBool, ""},
			{"NSGlobalDomain", "PMPrintingExpandedStateForPrint", true, defaults.TypeBool, ""},
		},
	},
}

func runPreset(cmd *cobra.Command, args []string) error {
	name := strings.ToLower(args[0])

	pd, ok := presets[name]
	if !ok {
		var available []string
		for k := range presets {
			available = append(available, k)
		}
		return fmt.Errorf("unknown preset %q, available: %s", name, strings.Join(available, ", "))
	}

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	fmt.Printf("Applying preset: %s\n\n", pd.Name)

	restartNeeded := make(map[string]bool)
	var applied, failed int

	for _, entry := range pd.Settings {
		if err := exec.Write(entry.Domain, entry.Key, entry.Value, entry.Type); err != nil {
			fmt.Printf("  FAIL  %s %s: %v\n", entry.Domain, entry.Key, err)
			failed++
			continue
		}
		fmt.Printf("  OK    %s %s = %v\n", entry.Domain, entry.Key, entry.Value)
		applied++
		if entry.Restart != "" {
			restartNeeded[entry.Restart] = true
		}
	}

	// Restart all required processes.
	for proc := range restartNeeded {
		if _, err := runner.Run(cmd.Context(), "killall", proc); err != nil {
			fmt.Printf("\n  Warning: failed to restart %s: %v\n", proc, err)
		} else {
			fmt.Printf("\n  Restarted %s\n", proc)
		}
	}

	fmt.Printf("\nApplied %d settings", applied)
	if failed > 0 {
		fmt.Printf(" (%d failed)", failed)
	}
	fmt.Println()
	return nil
}
