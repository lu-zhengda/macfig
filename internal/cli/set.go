package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/lu-zhengda/macfig/internal/defaults"
	"github.com/lu-zhengda/macfig/internal/preset"
)

var setCmd = &cobra.Command{
	Use:   "set <domain> <key> <value>",
	Short: "Write a defaults value",
	Long:  "Sets a macOS defaults value. Type is inferred from presets or defaults to string.",
	Args:  cobra.ExactArgs(3),
	RunE:  runSet,
}

func runSet(cmd *cobra.Command, args []string) error {
	domain, key, rawVal := args[0], args[1], args[2]

	runner := &defaults.RealCmdRunner{}
	exec := defaults.NewExecutor(runner)

	// Try to find preset info for this domain/key to determine type.
	vt := defaults.TypeString
	restart := ""
	if s := preset.FindSetting(domain, key); s != nil {
		vt = s.Type
		restart = s.RequiresRestart
	}

	// Parse the value based on type.
	var value interface{}
	switch vt {
	case defaults.TypeBool:
		v, err := strconv.ParseBool(rawVal)
		if err != nil {
			return fmt.Errorf("invalid bool value %q: %w", rawVal, err)
		}
		value = v
	case defaults.TypeInt:
		v, err := strconv.Atoi(rawVal)
		if err != nil {
			return fmt.Errorf("invalid int value %q: %w", rawVal, err)
		}
		value = v
	case defaults.TypeFloat:
		v, err := strconv.ParseFloat(rawVal, 64)
		if err != nil {
			return fmt.Errorf("invalid float value %q: %w", rawVal, err)
		}
		value = v
	default:
		value = rawVal
	}

	if err := defaults.WriteAndRestart(exec, runner, domain, key, value, vt, restart); err != nil {
		return err
	}

	fmt.Printf("Set %s %s = %v\n", domain, key, value)
	if restart != "" {
		fmt.Printf("Restarted %s\n", restart)
	}
	return nil
}
