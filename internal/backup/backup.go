package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lu-zhengda/macfig/internal/defaults"
	"github.com/lu-zhengda/macfig/internal/preset"
)

// Entry represents a single backed-up defaults value.
type Entry struct {
	Domain string          `json:"domain"`
	Key    string          `json:"key"`
	Value  interface{}     `json:"value"`
	Type   defaults.ValueType `json:"type"`
}

// Backup represents a complete backup of tracked settings.
type Backup struct {
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Settings  []Entry   `json:"settings"`
}

// DefaultDir returns the default backup directory.
func DefaultDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".config", "macfig", "backups")
}

// DefaultPath returns a timestamped backup file path.
func DefaultPath() string {
	ts := time.Now().Format("2006-01-02T15-04-05")
	return filepath.Join(DefaultDir(), fmt.Sprintf("backup-%s.json", ts))
}

// Save reads all tracked settings and writes them to a JSON file.
func Save(exec defaults.Executor, path string) error {
	settings := preset.AllSettings()
	var entries []Entry

	for _, s := range settings {
		val, err := exec.Read(s.Domain, s.Key)
		if err != nil {
			// Skip settings that cannot be read (not set).
			continue
		}
		entries = append(entries, Entry{
			Domain: s.Domain,
			Key:    s.Key,
			Value:  val,
			Type:   s.Type,
		})
	}

	b := Backup{
		Version:   1,
		Timestamp: time.Now(),
		Settings:  entries,
	}

	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// Load reads a backup from a JSON file.
func Load(path string) (*Backup, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	var b Backup
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("failed to parse backup file: %w", err)
	}

	return &b, nil
}

// Restore applies all settings from a backup file.
func Restore(exec defaults.Executor, runner defaults.CmdRunner, path string) (int, int, error) {
	b, err := Load(path)
	if err != nil {
		return 0, 0, err
	}

	var restored, failed int
	restartNeeded := make(map[string]bool)

	for _, entry := range b.Settings {
		if err := exec.Write(entry.Domain, entry.Key, entry.Value, entry.Type); err != nil {
			failed++
			continue
		}
		restored++

		// Check if this setting requires a process restart.
		s := preset.FindSetting(entry.Domain, entry.Key)
		if s != nil && s.RequiresRestart != "" {
			restartNeeded[s.RequiresRestart] = true
		}
	}

	// Restart all required processes.
	for proc := range restartNeeded {
		_, _ = runner.Run(context.Background(), "killall", proc)
	}

	return restored, failed, nil
}
