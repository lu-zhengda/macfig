package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
	"github.com/lu-zhengda/macfig/internal/preset"
)

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import defaults from an export file and apply them",
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

// importResult is the JSON representation of an import result.
type importResult struct {
	Path      string              `json:"path"`
	Applied   int                 `json:"applied"`
	Skipped   int                 `json:"skipped"`
	Failed    int                 `json:"failed"`
	Restarted []string            `json:"restarted,omitempty"`
	Results   []importEntryResult `json:"results"`
}

// importEntryResult is the JSON representation of a single import entry result.
type importEntryResult struct {
	Domain string `json:"domain"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Status string `json:"status"` // "applied", "skipped", "failed"
	Error  string `json:"error,omitempty"`
}

func runImport(cmd *cobra.Command, args []string) error {
	path := args[0]

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	var imp exportData
	if err := json.Unmarshal(data, &imp); err != nil {
		return fmt.Errorf("failed to parse import file: %w", err)
	}

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	var applied, skipped, failed int
	var results []importEntryResult
	restartNeeded := make(map[string]bool)

	for _, entry := range imp.Settings {
		// Parse the value to the correct Go type.
		value, err := parseImportValue(entry.Value, entry.Type)
		if err != nil {
			skipped++
			results = append(results, importEntryResult{
				Domain: entry.Domain,
				Key:    entry.Key,
				Value:  entry.Value,
				Status: "skipped",
				Error:  err.Error(),
			})
			continue
		}

		if err := exec.Write(entry.Domain, entry.Key, value, entry.Type); err != nil {
			failed++
			if !jsonFlag {
				fmt.Printf("  FAIL  %s %s: %v\n", entry.Domain, entry.Key, err)
			}
			results = append(results, importEntryResult{
				Domain: entry.Domain,
				Key:    entry.Key,
				Value:  entry.Value,
				Status: "failed",
				Error:  err.Error(),
			})
			continue
		}

		applied++
		if !jsonFlag {
			fmt.Printf("  OK    %s %s = %s\n", entry.Domain, entry.Key, entry.Value)
		}
		results = append(results, importEntryResult{
			Domain: entry.Domain,
			Key:    entry.Key,
			Value:  entry.Value,
			Status: "applied",
		})

		// Check if restart is needed.
		if s := preset.FindSetting(entry.Domain, entry.Key); s != nil && s.RequiresRestart != "" {
			restartNeeded[s.RequiresRestart] = true
		}
	}

	// Restart required processes.
	var restarted []string
	for proc := range restartNeeded {
		if _, err := runner.Run(context.Background(), "killall", proc); err != nil {
			if !jsonFlag {
				fmt.Printf("  Warning: failed to restart %s: %v\n", proc, err)
			}
		} else {
			if !jsonFlag {
				fmt.Printf("  Restarted %s\n", proc)
			}
			restarted = append(restarted, proc)
		}
	}

	if jsonFlag {
		return printJSON(importResult{
			Path:      path,
			Applied:   applied,
			Skipped:   skipped,
			Failed:    failed,
			Restarted: restarted,
			Results:   results,
		})
	}

	fmt.Printf("\nImported %d settings from %s", applied, path)
	if failed > 0 {
		fmt.Printf(" (%d failed)", failed)
	}
	if skipped > 0 {
		fmt.Printf(" (%d skipped)", skipped)
	}
	fmt.Println()
	return nil
}

// parseImportValue converts a string value to the appropriate Go type.
func parseImportValue(raw string, vt defaults.ValueType) (interface{}, error) {
	switch vt {
	case defaults.TypeBool:
		return parseBoolImport(raw)
	case defaults.TypeInt:
		return strconv.Atoi(raw)
	case defaults.TypeFloat:
		return strconv.ParseFloat(raw, 64)
	case defaults.TypeString:
		return raw, nil
	default:
		return raw, nil
	}
}

// parseBoolImport handles macOS defaults boolean representations.
func parseBoolImport(s string) (bool, error) {
	switch s {
	case "1", "true", "yes":
		return true, nil
	case "0", "false", "no":
		return false, nil
	default:
		return strconv.ParseBool(s)
	}
}
