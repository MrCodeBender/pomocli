package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/angelchaudhary01/pomocli/internal/config"
)

func TestDefaultValues(t *testing.T) {
	cfg := config.Defaults()

	if cfg.Pomodoro.WorkDuration != 25 {
		t.Errorf("expected WorkDuration=25, got %d", cfg.Pomodoro.WorkDuration)
	}
	if cfg.Pomodoro.ShortBreak != 5 {
		t.Errorf("expected ShortBreak=5, got %d", cfg.Pomodoro.ShortBreak)
	}
	if cfg.Pomodoro.LongBreak != 15 {
		t.Errorf("expected LongBreak=15, got %d", cfg.Pomodoro.LongBreak)
	}
	if cfg.Pomodoro.LongBreakAfter != 4 {
		t.Errorf("expected LongBreakAfter=4, got %d", cfg.Pomodoro.LongBreakAfter)
	}
	if !cfg.Notifications.Enabled {
		t.Error("expected Notifications.Enabled=true")
	}
	if !cfg.Notifications.Sound {
		t.Error("expected Notifications.Sound=true")
	}
	if !cfg.Display.ShowProgressBar {
		t.Error("expected Display.ShowProgressBar=true")
	}
}

func TestLoadFromYAML(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")
	content := `
pomodoro:
  work_duration: 30
  short_break: 10
`
	os.WriteFile(configFile, []byte(content), 0644)

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pomodoro.WorkDuration != 30 {
		t.Errorf("expected WorkDuration=30, got %d", cfg.Pomodoro.WorkDuration)
	}
	// Unset fields should use defaults
	if cfg.Pomodoro.LongBreak != 15 {
		t.Errorf("expected LongBreak=15 (default), got %d", cfg.Pomodoro.LongBreak)
	}
}
