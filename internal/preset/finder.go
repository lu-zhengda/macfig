package preset

import "github.com/zhengda-lu/macfig/internal/defaults"

// FinderCategory returns Finder-related settings.
func FinderCategory() Category {
	return Category{
		Name: "Finder",
		Icon: "\U0001F4C1", // file folder emoji
		Settings: []Setting{
			{
				Name:            "Show path bar",
				Description:     "Show the path bar at the bottom of Finder windows",
				Domain:          "com.apple.finder",
				Key:             "ShowPathbar",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "Finder",
			},
			{
				Name:            "Show status bar",
				Description:     "Show the status bar at the bottom of Finder windows",
				Domain:          "com.apple.finder",
				Key:             "ShowStatusBar",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "Finder",
			},
			{
				Name:            "Show all file extensions",
				Description:     "Always show file extensions in Finder",
				Domain:          "NSGlobalDomain",
				Key:             "AppleShowAllExtensions",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "Finder",
			},
			{
				Name:            "Extension change warning",
				Description:     "Warn when changing a file extension",
				Domain:          "com.apple.finder",
				Key:             "FXEnableExtensionChangeWarning",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     false,
				RequiresRestart: "Finder",
			},
			{
				Name:            "Show POSIX path in title",
				Description:     "Display full POSIX path in Finder window title",
				Domain:          "com.apple.finder",
				Key:             "_FXShowPosixPathInTitle",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "Finder",
			},
			{
				Name:            "Default search scope",
				Description:     "Search scope (SCcf=current folder, SCsp=previous scope, SCev=entire Mac)",
				Domain:          "com.apple.finder",
				Key:             "FXDefaultSearchScope",
				Type:            defaults.TypeString,
				Default:         "SCev",
				Recommended:     "SCcf",
				RequiresRestart: "Finder",
			},
			{
				Name:            "New window target",
				Description:     "Default location for new Finder windows (PfHm=home, PfDe=desktop, PfDo=documents)",
				Domain:          "com.apple.finder",
				Key:             "NewWindowTarget",
				Type:            defaults.TypeString,
				Default:         "PfHm",
				Recommended:     nil,
				RequiresRestart: "Finder",
			},
			{
				Name:            "Show hard drives on desktop",
				Description:     "Show hard drives on the desktop",
				Domain:          "com.apple.finder",
				Key:             "ShowHardDrivesOnDesktop",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     nil,
				RequiresRestart: "Finder",
			},
			{
				Name:            "Sort folders first",
				Description:     "Keep folders on top when sorting by name",
				Domain:          "com.apple.finder",
				Key:             "_FXSortFoldersFirst",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "Finder",
			},
		},
	}
}
