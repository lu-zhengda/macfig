package preset

import "github.com/lu-zhengda/macfig/internal/defaults"

// Setting represents a single macOS defaults setting.
type Setting struct {
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Domain          string             `json:"domain"`
	Key             string             `json:"key"`
	Type            defaults.ValueType `json:"type"`
	Default         interface{}        `json:"default,omitempty"`
	Recommended     interface{}        `json:"recommended,omitempty"`
	RequiresRestart string             `json:"requires_restart,omitempty"`
}

// Category groups related settings together.
type Category struct {
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	Settings []Setting `json:"settings"`
}

// AllCategories returns every curated preset category.
func AllCategories() []Category {
	return []Category{
		DockCategory(),
		FinderCategory(),
		ScreenshotsCategory(),
		KeyboardCategory(),
		TrackpadCategory(),
		AnimationsCategory(),
		PrivacyCategory(),
		MiscCategory(),
	}
}

// FindCategory returns a category by name (case-insensitive).
func FindCategory(name string) *Category {
	for _, cat := range AllCategories() {
		if equalFold(cat.Name, name) {
			return &cat
		}
	}
	return nil
}

// FindSetting searches all categories for a setting matching domain and key.
func FindSetting(domain, key string) *Setting {
	for _, cat := range AllCategories() {
		for i := range cat.Settings {
			if cat.Settings[i].Domain == domain && cat.Settings[i].Key == key {
				return &cat.Settings[i]
			}
		}
	}
	return nil
}

// AllSettings returns a flat list of every curated setting.
func AllSettings() []Setting {
	var all []Setting
	for _, cat := range AllCategories() {
		all = append(all, cat.Settings...)
	}
	return all
}

// equalFold does a simple case-insensitive comparison.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
