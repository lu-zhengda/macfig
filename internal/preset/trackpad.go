package preset

import "github.com/zhengda-lu/macfig/internal/defaults"

// TrackpadCategory returns trackpad-related settings.
func TrackpadCategory() Category {
	return Category{
		Name: "Trackpad",
		Icon: "\U0001F91A", // raised back of hand emoji
		Settings: []Setting{
			{
				Name:            "Tap to click",
				Description:     "Enable tap to click on the trackpad",
				Domain:          "com.apple.AppleMultitouchTrackpad",
				Key:             "Clicking",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "",
			},
			{
				Name:            "Three-finger drag",
				Description:     "Enable three-finger drag gesture",
				Domain:          "com.apple.AppleMultitouchTrackpad",
				Key:             "TrackpadThreeFingerDrag",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "",
			},
			{
				Name:            "Natural scrolling",
				Description:     "Scroll content in the direction of finger movement",
				Domain:          "NSGlobalDomain",
				Key:             "com.apple.swipescrolldirection",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     nil,
				RequiresRestart: "",
			},
		},
	}
}
