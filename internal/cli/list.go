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

func runList(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return listCategories()
	}
	return listCategory(args[0])
}

func listCategories() error {
	categories := preset.AllCategories()
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
