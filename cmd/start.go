package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/angelchaudhary01/pomocli/internal/config"
	pomlog "github.com/angelchaudhary01/pomocli/internal/log"
	"github.com/angelchaudhary01/pomocli/internal/tui"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Pomodoro session",
	RunE:  runStart,
}

var taskFlag string

func init() {
	startCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task name for this session")
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadDefault()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	startTime := time.Now()
	model := tui.NewTimer(taskFlag, cfg)

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("run TUI: %w", err)
	}

	m, ok := finalModel.(tui.TimerModel)
	if !ok {
		return nil
	}

	status := pomlog.StatusCancelled
	if m.PomCount() > 0 {
		status = pomlog.StatusCompleted
	}

	session := pomlog.Session{
		Number:    m.PomCount(),
		Task:      taskFlag,
		StartTime: startTime,
		EndTime:   time.Now(),
		Duration:  time.Since(startTime).Round(time.Minute),
		Status:    status,
	}

	if err := pomlog.WriteSession(cfg.Logs.Directory, session); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not write log: %v\n", err)
	} else {
		fmt.Printf("Log guardado en %s/%s.md\n", cfg.Logs.Directory, startTime.Format("2006-01-02"))
	}

	return nil
}
