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
	lines := []string{
		fmt.Sprintf("## Sesión %d — %s", s.Number, s.StartTime.Format("15:04")),
		"",
		"| Campo      | Valor       |",
		"|------------|-------------|",
		fmt.Sprintf("| Tarea      | %s |", task),
		fmt.Sprintf("| Inicio     | %s |", s.StartTime.Format("15:04")),
		fmt.Sprintf("| Fin        | %s |", s.EndTime.Format("15:04")),
		fmt.Sprintf("| Duración   | %s |", formatDuration(s.Duration)),
		fmt.Sprintf("| Estado     | %s |", string(s.Status)),
		"",
		"",
	}
	return strings.Join(lines, "\n")
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

	// Note: read-modify-write is not atomic. Safe for single-process CLI use;
	// do not call concurrently for the same logDir/date combination.
	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
