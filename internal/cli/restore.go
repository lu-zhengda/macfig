package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/backup"
	"github.com/lu-zhengda/macfig/internal/defaults"
)

// restoreResult is the JSON representation of a restore operation result.
type restoreResult struct {
	OK       bool   `json:"ok"`
	Path     string `json:"path"`
	Restored int    `json:"restored"`
	Failed   int    `json:"failed"`
}

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

	if jsonFlag {
		return printJSON(restoreResult{
			OK:       failed == 0,
			Path:     path,
			Restored: restored,
			Failed:   failed,
		})
	}

	fmt.Printf("Restored %d settings", restored)
	if failed > 0 {
		fmt.Printf(" (%d failed)", failed)
	}
	fmt.Println()
	return nil
}
