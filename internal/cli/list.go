package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
	"github.com/lu-zhengda/macfig/internal/preset"
)

var listCmd = &cobra.Command{
	Use:   "list [category]",
	Short: "List preset categories or settings in a category",
	Long:  "Without arguments, lists all preset categories.\nWith a category name, lists all settings in that category with current values.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runList,
}

// categoryInfo is the JSON representation of a category summary.
type categoryInfo struct {
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Settings int    `json:"settings"`
}

// settingInfo is the JSON representation of a setting with its current value.
type settingInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Domain      string      `json:"domain"`
	Key         string      `json:"key"`
	Current     string      `json:"current"`
	Recommended interface{} `json:"recommended,omitempty"`
}

func runList(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return listCategories()
	}
	return listCategory(args[0])
}

func listCategories() error {
	categories := preset.AllCategories()

	if jsonFlag {
		var items []categoryInfo
		for _, cat := range categories {
			items = append(items, categoryInfo{
				Name:     cat.Name,
				Icon:     cat.Icon,
				Settings: len(cat.Settings),
			})
		}
		return printJSON(items)
	}

	fmt.Println("Categories:")
	fmt.Println()
	for _, cat := range categories {
		fmt.Printf("  %s %s  (%d settings)\n", cat.Icon, cat.Name, len(cat.Settings))
	}
	fmt.Println()
	fmt.Println("Use 'macfig list <category>' to view settings.")
	return nil
}

func listCategory(name string) error {
	cat := preset.FindCategory(name)
	if cat == nil {
		return fmt.Errorf("unknown category: %s", name)
	}

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	if jsonFlag {
		var items []settingInfo
		for _, s := range cat.Settings {
			current, err := exec.Read(s.Domain, s.Key)
			if err != nil {
				current = ""
			}
			items = append(items, settingInfo{
				Name:        s.Name,
				Description: s.Description,
				Domain:      s.Domain,
				Key:         s.Key,
				Current:     current,
				Recommended: s.Recommended,
			})
		}
		return printJSON(items)
	}

	fmt.Printf("%s %s\n", cat.Icon, cat.Name)
	fmt.Println()

	for _, s := range cat.Settings {
		current, err := exec.Read(s.Domain, s.Key)
		if err != nil {
			current = "(not set)"
		}

		rec := ""
		if s.Recommended != nil {
			rec = fmt.Sprintf(" -> recommended: %v", s.Recommended)
		}

		fmt.Printf("  %-35s %s%s\n", s.Name, current, rec)
		fmt.Printf("    %s  (%s %s)\n", s.Description, s.Domain, s.Key)
	}

	return nil
}
