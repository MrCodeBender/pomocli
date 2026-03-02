package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Status string

const (
	StatusCompleted Status = "completado"
	StatusCancelled Status = "cancelado"
	StatusSkipped   Status = "saltado"
)

type Session struct {
	Number    int
	Task      string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Status    Status
}

// FormatSession returns the markdown block for a single session.
func FormatSession(s Session) string {
	task := s.Task
	if task == "" {
		task = "—"
	}
	return fmt.Sprintf("## Sesión %d — %s\n\n| Campo      | Valor                |\n|------------|----------------------|\n| Tarea      | %-20s |\n| Inicio     | %-20s |\n| Fin        | %-20s |\n| Duración   | %-20s |\n| Estado     | %-20s |\n\n",
		s.Number,
		s.StartTime.Format("15:04"),
		task,
		s.StartTime.Format("15:04"),
		s.EndTime.Format("15:04"),
		formatDuration(s.Duration),
		string(s.Status),
	)
}

func formatDuration(d time.Duration) string {
	m := int(d.Minutes())
	if m < 60 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dh %dm", m/60, m%60)
}

// WriteSession appends a session to the daily log file, creating it if needed.
func WriteSession(logDir string, s Session) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	filename := filepath.Join(logDir, s.StartTime.Format("2006-01-02")+".md")

	var existing string
	if data, err := os.ReadFile(filename); err == nil {
		existing = string(data)
	}

	var sb strings.Builder

	if existing == "" {
		sb.WriteString(fmt.Sprintf("# Pomodoros — %s\n\n", s.StartTime.Format("2006-01-02")))
		sb.WriteString("---\n\n")
	} else {
		sb.WriteString(existing)
	}

	sb.WriteString(FormatSession(s))

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
