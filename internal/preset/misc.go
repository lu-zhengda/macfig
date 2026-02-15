package preset

import "github.com/zhengda-lu/macfig/internal/defaults"

// MiscCategory returns miscellaneous settings.
func MiscCategory() Category {
	return Category{
		Name: "Misc",
		Icon: "\U0001F527", // wrench emoji (BMP gear emoji renders inconsistently)
		Settings: []Setting{
			{
				Name:            "Expand save panel",
				Description:     "Expand save panel by default",
				Domain:          "NSGlobalDomain",
				Key:             "NSNavPanelExpandedStateForSaveMode",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "",
			},
			{
				Name:            "Expand print panel",
				Description:     "Expand print panel by default",
				Domain:          "NSGlobalDomain",
				Key:             "PMPrintingExpandedStateForPrint",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "",
			},
			{
				Name:            "Crash reporter type",
				Description:     "Use notification instead of dialog for crash reports",
				Domain:          "com.apple.CrashReporter",
				Key:             "UseUNC",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "",
			},
			{
				Name:            "Resume apps on reopen",
				Description:     "Close windows when quitting an app (disable resume)",
				Domain:          "NSGlobalDomain",
				Key:             "NSQuitAlwaysKeepsWindows",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     false,
				RequiresRestart: "",
			},
		},
	}
}
