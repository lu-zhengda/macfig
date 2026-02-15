package preset

import "github.com/lu-zhengda/macfig/internal/defaults"

// KeyboardCategory returns keyboard-related settings.
func KeyboardCategory() Category {
	return Category{
		Name: "Keyboard",
		Icon: "\U0001F4BB", // laptop emoji (BMP keyboard emoji renders inconsistently)
		Settings: []Setting{
			{
				Name:            "Key repeat rate",
				Description:     "Speed of key repeat (lower = faster, minimum 1)",
				Domain:          "NSGlobalDomain",
				Key:             "KeyRepeat",
				Type:            defaults.TypeInt,
				Default:         6,
				Recommended:     2,
				RequiresRestart: "",
			},
			{
				Name:            "Initial key repeat delay",
				Description:     "Delay before key repeat starts (lower = shorter)",
				Domain:          "NSGlobalDomain",
				Key:             "InitialKeyRepeat",
				Type:            defaults.TypeInt,
				Default:         25,
				Recommended:     15,
				RequiresRestart: "",
			},
			{
				Name:            "Press and hold for accents",
				Description:     "Show accent menu on key press-and-hold (disable for key repeat)",
				Domain:          "NSGlobalDomain",
				Key:             "ApplePressAndHoldEnabled",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     false,
				RequiresRestart: "",
			},
			{
				Name:            "Auto-correct",
				Description:     "Enable automatic spelling correction",
				Domain:          "NSGlobalDomain",
				Key:             "NSAutomaticSpellingCorrectionEnabled",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     false,
				RequiresRestart: "",
			},
			{
				Name:            "Auto-capitalize",
				Description:     "Automatically capitalize first letter of sentences",
				Domain:          "NSGlobalDomain",
				Key:             "NSAutomaticCapitalizationEnabled",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     false,
				RequiresRestart: "",
			},
			{
				Name:            "Auto period substitution",
				Description:     "Add period with double-space",
				Domain:          "NSGlobalDomain",
				Key:             "NSAutomaticPeriodSubstitutionEnabled",
				Type:            defaults.TypeBool,
				Default:         true,
				Recommended:     false,
				RequiresRestart: "",
			},
		},
	}
}
