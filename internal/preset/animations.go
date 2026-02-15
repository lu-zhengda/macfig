package preset

import "github.com/lu-zhengda/macfig/internal/defaults"

// AnimationsCategory returns animation-related settings.
func AnimationsCategory() Category {
	return Category{
		Name: "Animations",
		Icon: "\U0001F3AC", // clapper board emoji
		Settings: []Setting{
			{
				Name:            "Window animations",
				Description:     "Enable window opening/closing animations",
				Domain:          "NSGlobalDomain",
				Key:             "NSAutomaticWindowAnimationsEnabled",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     false,
				RequiresRestart: "",
			},
			{
				Name:            "Window resize speed",
				Description:     "Duration of window resize animation (seconds)",
				Domain:          "NSGlobalDomain",
				Key:             "NSWindowResizeTime",
				Type:            defaults.TypeFloat,
				Default:         0.2,
				Recommended:     0.001,
				RequiresRestart: "",
			},
		},
	}
}
