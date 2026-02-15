package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhengda-lu/macfig/internal/defaults"
)

var getCmd = &cobra.Command{
	Use:   "get <domain> <key>",
	Short: "Read a specific defaults value",
	Args:  cobra.ExactArgs(2),
	RunE:  runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	domain, key := args[0], args[1]

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	val, err := exec.Read(domain, key)
	if err != nil {
		return fmt.Errorf("failed to read %s %s: %w", domain, key, err)
	}

	fmt.Println(val)
	return nil
}
