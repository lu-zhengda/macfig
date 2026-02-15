package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
	"github.com/lu-zhengda/macfig/internal/preset"
)

var resetCmd = &cobra.Command{
	Use:   "reset <domain> <key>",
	Short: "Reset a defaults key to macOS default",
	Long:  "Deletes a defaults key, restoring it to the macOS system default value.",
	Args:  cobra.ExactArgs(2),
	RunE:  runReset,
}

func runReset(cmd *cobra.Command, args []string) error {
	domain, key := args[0], args[1]

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	restart := ""
	if s := preset.FindSetting(domain, key); s != nil {
		restart = s.RequiresRestart
	}

	if err := defaults.DeleteAndRestart(exec, runner, domain, key, restart); err != nil {
		return err
	}

	fmt.Printf("Reset %s %s to macOS default\n", domain, key)
	if restart != "" {
		fmt.Printf("Restarted %s\n", restart)
	}
	return nil
}
