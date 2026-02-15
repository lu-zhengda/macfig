package defaults

import (
	"context"
	"fmt"
)

// WriteAndRestart writes a defaults value and optionally restarts the
// required process for changes to take effect.
func WriteAndRestart(exec Executor, runner CmdRunner, domain, key string, value interface{}, vt ValueType, restart string) error {
	if err := exec.Write(domain, key, value, vt); err != nil {
		return err
	}

	if restart != "" {
		return restartProcess(runner, restart)
	}
	return nil
}

// DeleteAndRestart deletes a defaults key (resetting to macOS default)
// and optionally restarts the required process.
func DeleteAndRestart(exec Executor, runner CmdRunner, domain, key, restart string) error {
	if err := exec.Delete(domain, key); err != nil {
		return err
	}

	if restart != "" {
		return restartProcess(runner, restart)
	}
	return nil
}

// restartProcess kills a process by name so macOS can relaunch it.
func restartProcess(runner CmdRunner, process string) error {
	_, err := runner.Run(context.Background(), "killall", process)
	if err != nil {
		return fmt.Errorf("failed to restart %s: %w", process, err)
	}
	return nil
}
