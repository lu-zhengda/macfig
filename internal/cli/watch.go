package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
)

var watchCmd = &cobra.Command{
	Use:   "watch <domain> <key>",
	Short: "Monitor a defaults key for changes",
	Long:  "Polls a specific defaults key and prints when the value changes.\nUse --interval to set poll frequency (default 5s).",
	Args:  cobra.ExactArgs(2),
	RunE:  runWatch,
}

var watchInterval time.Duration

func init() {
	watchCmd.Flags().DurationVar(&watchInterval, "interval", 5*time.Second, "Poll interval (e.g. 2s, 500ms)")
}

// watchEvent is the JSON representation of a value change event.
type watchEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Domain    string    `json:"domain"`
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
}

func runWatch(cmd *cobra.Command, args []string) error {
	domain, key := args[0], args[1]

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	// Read initial value.
	current, err := exec.Read(domain, key)
	if err != nil {
		current = ""
	}

	if !jsonFlag {
		fmt.Printf("Watching %s %s (interval: %s)\n", domain, key, watchInterval)
		if current != "" {
			fmt.Printf("  Current value: %s\n", current)
		} else {
			fmt.Println("  Current value: (not set)")
		}
		fmt.Println("  Press Ctrl+C to stop.")
	}

	// Set up signal handling for graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(watchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sigCh:
			if !jsonFlag {
				fmt.Println("\nStopped watching.")
			}
			return nil
		case <-ticker.C:
			newVal, err := exec.Read(domain, key)
			if err != nil {
				newVal = ""
			}

			if newVal != current {
				if jsonFlag {
					if err := printJSON(watchEvent{
						Timestamp: time.Now(),
						Domain:    domain,
						Key:       key,
						OldValue:  current,
						NewValue:  newVal,
					}); err != nil {
						return err
					}
				} else {
					ts := time.Now().Format("15:04:05")
					fmt.Printf("  [%s] %s -> %s\n", ts, current, newVal)
				}
				current = newVal
			}
		}
	}
}
