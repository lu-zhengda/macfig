package preset

import "github.com/lu-zhengda/macfig/internal/defaults"

// ScreenshotsCategory returns screenshot-related settings.
func ScreenshotsCategory() Category {
	return Category{
		Name: "Screenshots",
		Icon: "\U0001F4F7", // camera emoji
		Settings: []Setting{
			{
				Name:            "Save location",
				Description:     "Directory where screenshots are saved",
				Domain:          "com.apple.screencapture",
				Key:             "location",
				Type:            defaults.TypeString,
				Default:         "~/Desktop",
				Recommended:     nil,
				RequiresRestart: "SystemUIServer",
			},
			{
				Name:            "Image format",
				Description:     "Screenshot image format (png, jpg, pdf, tiff)",
				Domain:          "com.apple.screencapture",
				Key:             "type",
				Type:            defaults.TypeString,
				Default:         "png",
				Recommended:     nil,
				RequiresRestart: "SystemUIServer",
			},
			{
				Name:            "Disable shadow",
				Description:     "Disable window shadow in screenshots",
				Domain:          "com.apple.screencapture",
				Key:             "disable-shadow",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "SystemUIServer",
			},
			{
				Name:            "Include date in filename",
				Description:     "Include date and time in screenshot filenames",
				Domain:          "com.apple.screencapture",
				Key:             "include-date",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     nil,
				RequiresRestart: "SystemUIServer",
			},
			{
				Name:            "Show thumbnail",
				Description:     "Show floating thumbnail after taking a screenshot",
				Domain:          "com.apple.screencapture",
				Key:             "show-thumbnail",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     nil,
				RequiresRestart: "SystemUIServer",
			},
		},
	}
}
