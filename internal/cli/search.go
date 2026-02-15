package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for keys across all defaults domains",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	results, err := exec.Find(query)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("No results found for %q\n", query)
		return nil
	}

	fmt.Printf("Found %d results for %q:\n\n", len(results), query)
	for _, r := range results {
		val := r.Value
		if len(val) > 80 {
			val = val[:77] + "..."
		}
		fmt.Printf("  %-40s %-30s = %s\n", r.Domain, r.Key, val)
	}
	return nil
}
