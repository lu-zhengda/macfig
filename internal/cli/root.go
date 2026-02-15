package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/zhengda-lu/macfig/internal/defaults"
	"github.com/zhengda-lu/macfig/internal/tui"
)

var (
	// version is set via ldflags at build time.
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "macfig",
	Short:   "macOS hidden defaults manager",
	Long:    "macfig manages macOS hidden preferences and defaults.\nLaunch without subcommands for interactive TUI mode.",
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		runner := &defaults.RealCmdRunner{}
		exec := defaults.NewExecutor(runner)
		p := tea.NewProgram(tui.New(exec, runner), tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("macfig %s\n", version))
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(presetCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(searchCmd)
}
