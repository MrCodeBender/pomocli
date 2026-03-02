package log_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pomlog "github.com/angelchaudhary01/pomocli/internal/log"
)

func TestFormatSession(t *testing.T) {
	start := time.Date(2026, 2, 26, 9, 15, 0, 0, time.Local)
	end := time.Date(2026, 2, 26, 9, 40, 0, 0, time.Local)
	s := pomlog.Session{
		Task:      "revisar PRs",
		StartTime: start,
		EndTime:   end,
		Duration:  25 * time.Minute,
		Status:    pomlog.StatusCompleted,
		Number:    1,
	}

	result := pomlog.FormatSession(s)

	if !strings.Contains(result, "## Sesión 1 — 09:15") {
		t.Errorf("missing session header, got:\n%s", result)
	}
	if !strings.Contains(result, "revisar PRs") {
		t.Errorf("missing task name, got:\n%s", result)
	}
	if !strings.Contains(result, "completado") {
		t.Errorf("missing status, got:\n%s", result)
	}
	if !strings.Contains(result, "09:40") {
		t.Errorf("missing end time, got:\n%s", result)
	}
}

func TestWriteSessionCreatesFile(t *testing.T) {
	dir := t.TempDir()
	start := time.Date(2026, 2, 26, 9, 15, 0, 0, time.Local)
	end := time.Date(2026, 2, 26, 9, 40, 0, 0, time.Local)
	s := pomlog.Session{
		Task:      "test task",
		StartTime: start,
		EndTime:   end,
		Duration:  25 * time.Minute,
		Status:    pomlog.StatusCompleted,
		Number:    1,
	}

	err := pomlog.WriteSession(dir, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "2026-02-26.md")
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", expected)
	}

	content, err := os.ReadFile(expected)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if !strings.Contains(string(content), "# Pomodoros — 2026-02-26") {
		t.Errorf("missing day header in file content:\n%s", content)
	}
}

func TestWriteSessionAppends(t *testing.T) {
	dir := t.TempDir()
	base := time.Date(2026, 2, 26, 9, 0, 0, 0, time.Local)

	for i := 1; i <= 2; i++ {
		s := pomlog.Session{
			Task:      "task",
			StartTime: base.Add(time.Duration(i-1) * 30 * time.Minute),
			EndTime:   base.Add(time.Duration(i-1)*30*time.Minute + 25*time.Minute),
			Duration:  25 * time.Minute,
			Status:    pomlog.StatusCompleted,
			Number:    i,
		}
		if err := pomlog.WriteSession(dir, s); err != nil {
			t.Fatalf("WriteSession failed for session %d: %v", i, err)
		}
	}

	content, err := os.ReadFile(filepath.Join(dir, "2026-02-26.md"))
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if strings.Count(string(content), "## Sesión") != 2 {
		t.Errorf("expected 2 sessions in file, got:\n%s", content)
	}
}
