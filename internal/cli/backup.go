package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/backup"
	"github.com/lu-zhengda/macfig/internal/defaults"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup all tracked defaults settings to JSON",
	RunE:  runBackup,
}

// backupResult is the JSON representation of a backup operation result.
type backupResult struct {
	Path    string `json:"path"`
	Entries int    `json:"entries"`
}

func runBackup(cmd *cobra.Command, args []string) error {
	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	path := backup.DefaultPath()
	if err := backup.Save(exec, path); err != nil {
		return err
	}

	if jsonFlag {
		b, err := backup.Load(path)
		if err != nil {
			return err
		}
		return printJSON(backupResult{
			Path:    path,
			Entries: len(b.Settings),
		})
	}

	fmt.Printf("Backup saved to %s\n", path)
	return nil
}
