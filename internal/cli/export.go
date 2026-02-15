package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
	"github.com/lu-zhengda/macfig/internal/preset"
)

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export all managed defaults to a portable JSON file",
	Long:  "Exports all managed defaults to a JSON file suitable for sharing across machines.\nDefaults to macfig-export.json in the current directory.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runExport,
}

// exportData is the portable export format.
type exportData struct {
	Version   int           `json:"version"`
	Timestamp time.Time     `json:"timestamp"`
	Hostname  string        `json:"hostname"`
	Settings  []exportEntry `json:"settings"`
}

// exportEntry represents a single exported setting.
type exportEntry struct {
	Domain      string             `json:"domain"`
	Key         string             `json:"key"`
	Value       string             `json:"value"`
	Type        defaults.ValueType `json:"type"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}

func runExport(cmd *cobra.Command, args []string) error {
	path := "macfig-export.json"
	if len(args) > 0 {
		path = args[0]
	}

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	allSettings := preset.AllSettings()
	var entries []exportEntry

	for _, s := range allSettings {
		val, err := exec.Read(s.Domain, s.Key)
		if err != nil {
			continue
		}
		entries = append(entries, exportEntry{
			Domain:      s.Domain,
			Key:         s.Key,
			Value:       val,
			Type:        s.Type,
			Name:        s.Name,
			Description: s.Description,
		})
	}

	hostname, _ := os.Hostname()

	data := exportData{
		Version:   1,
		Timestamp: time.Now(),
		Hostname:  hostname,
		Settings:  entries,
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export: %w", err)
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	if jsonFlag {
		return printJSON(struct {
			Path     string `json:"path"`
			Entries  int    `json:"entries"`
			Hostname string `json:"hostname"`
		}{
			Path:     path,
			Entries:  len(entries),
			Hostname: hostname,
		})
	}

	fmt.Printf("Exported %d settings to %s\n", len(entries), path)
	return nil
}
