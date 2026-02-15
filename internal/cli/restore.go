package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/backup"
	"github.com/lu-zhengda/macfig/internal/defaults"
)

var restoreCmd = &cobra.Command{
	Use:   "restore <file>",
	Short: "Restore defaults from a backup file",
	Args:  cobra.ExactArgs(1),
	RunE:  runRestore,
}

func runRestore(cmd *cobra.Command, args []string) error {
	path := args[0]

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	restored, failed, err := backup.Restore(exec, runner, path)
	if err != nil {
		return err
	}

	fmt.Printf("Restored %d settings", restored)
	if failed > 0 {
		fmt.Printf(" (%d failed)", failed)
	}
	fmt.Println()
	return nil
}
