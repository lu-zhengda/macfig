package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/backup"
	"github.com/lu-zhengda/macfig/internal/defaults"
)

var diffCmd = &cobra.Command{
	Use:   "diff <backup-file>",
	Short: "Compare current defaults against a backup snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runDiff,
}

// diffEntry represents a single difference between backup and current state.
type diffEntry struct {
	Domain   string `json:"domain"`
	Key      string `json:"key"`
	Type     string `json:"type"` // "added", "removed", "modified"
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

// diffOutput is the JSON representation of the diff result.
type diffOutput struct {
	BackupFile string      `json:"backup_file"`
	Added      int         `json:"added"`
	Removed    int         `json:"removed"`
	Modified   int         `json:"modified"`
	Unchanged  int         `json:"unchanged"`
	Changes    []diffEntry `json:"changes"`
}

func runDiff(cmd *cobra.Command, args []string) error {
	path := args[0]

	b, err := backup.Load(path)
	if err != nil {
		return err
	}

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	var changes []diffEntry
	var added, removed, modified, unchanged int

	// Check each backed-up setting against current state.
	for _, entry := range b.Settings {
		current, err := exec.Read(entry.Domain, entry.Key)
		if err != nil {
			// Key no longer exists (removed).
			removed++
			changes = append(changes, diffEntry{
				Domain:   entry.Domain,
				Key:      entry.Key,
				Type:     "removed",
				OldValue: fmt.Sprintf("%v", entry.Value),
			})
			continue
		}

		oldVal := fmt.Sprintf("%v", entry.Value)
		if current != oldVal {
			modified++
			changes = append(changes, diffEntry{
				Domain:   entry.Domain,
				Key:      entry.Key,
				Type:     "modified",
				OldValue: oldVal,
				NewValue: current,
			})
		} else {
			unchanged++
		}
	}

	// Check for settings that exist now but were not in backup.
	// We use a set of backed-up keys for efficient lookup.
	backedUp := make(map[string]bool)
	for _, entry := range b.Settings {
		backedUp[entry.Domain+"\x00"+entry.Key] = true
	}

	// Read current values for all tracked settings to detect additions.
	allSettings := backup.TrackedSettings()
	for _, s := range allSettings {
		key := s.Domain + "\x00" + s.Key
		if backedUp[key] {
			continue
		}
		current, err := exec.Read(s.Domain, s.Key)
		if err != nil {
			// Not set in either backup or current, skip.
			continue
		}
		added++
		changes = append(changes, diffEntry{
			Domain:   s.Domain,
			Key:      s.Key,
			Type:     "added",
			NewValue: current,
		})
	}

	if jsonFlag {
		return printJSON(diffOutput{
			BackupFile: path,
			Added:      added,
			Removed:    removed,
			Modified:   modified,
			Unchanged:  unchanged,
			Changes:    changes,
		})
	}

	fmt.Printf("Diff against %s:\n\n", path)

	if len(changes) == 0 {
		fmt.Println("  No changes detected.")
		return nil
	}

	for _, c := range changes {
		switch c.Type {
		case "added":
			fmt.Printf("  + %-40s %-30s = %s\n", c.Domain, c.Key, c.NewValue)
		case "removed":
			fmt.Printf("  - %-40s %-30s (was: %s)\n", c.Domain, c.Key, c.OldValue)
		case "modified":
			fmt.Printf("  ~ %-40s %-30s %s -> %s\n", c.Domain, c.Key, c.OldValue, c.NewValue)
		}
	}

	fmt.Printf("\nSummary: %d added, %d removed, %d modified, %d unchanged\n",
		added, removed, modified, unchanged)
	return nil
}
