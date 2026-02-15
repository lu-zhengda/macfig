package preset

import "github.com/lu-zhengda/macfig/internal/defaults"

// PrivacyCategory returns privacy/security-related settings.
func PrivacyCategory() Category {
	return Category{
		Name: "Privacy",
		Icon: "\U0001F512", // lock emoji
		Settings: []Setting{
			{
				Name:            "Ask for password",
				Description:     "Require password after sleep or screen saver",
				Domain:          "com.apple.screensaver",
				Key:             "askForPassword",
				Type:            defaults.TypeBool,
				Default:         false,
				Recommended:     true,
				RequiresRestart: "",
			},
			{
				Name:            "Password delay",
				Description:     "Seconds before password is required after sleep (0 = immediately)",
				Domain:          "com.apple.screensaver",
				Key:             "askForPasswordDelay",
				Type:            defaults.TypeInt,
				Default:         5,
				Recommended:     0,
				RequiresRestart: "",
			},
		},
	}
}
