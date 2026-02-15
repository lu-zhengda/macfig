package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhengda-lu/macfig/internal/backup"
	"github.com/zhengda-lu/macfig/internal/defaults"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup all tracked defaults settings to JSON",
	RunE:  runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	path := backup.DefaultPath()
	if err := backup.Save(exec, path); err != nil {
		return err
	}

	fmt.Printf("Backup saved to %s\n", path)
	return nil
}
