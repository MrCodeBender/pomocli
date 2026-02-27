# pomocli Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go CLI Pomodoro app with an interactive TUI (bubbletea), YAML config, daily `.md` logs, and OS notifications.

**Architecture:** cobra handles subcommands (`start`, `log`, `config`); bubbletea drives the interactive timer TUI with a state machine (Working → ShortBreak → LongBreak); viper loads config from `~/.config/pomocli/config.yaml`; the log writer appends Pomodoro sessions to daily markdown files.

**Tech Stack:** Go 1.22+, cobra v1, bubbletea v1, lipgloss v1, viper v2

---

## Task 1: Initialize Go module and project structure

**Files:**
- Create: `main.go`
- Create: `go.mod`
- Create: `cmd/root.go`
- Create: `internal/config/config.go` (empty stub)
- Create: `internal/tui/timer.go` (empty stub)
- Create: `internal/log/writer.go` (empty stub)
- Create: `internal/notify/notify.go` (empty stub)

**Step 1: Initialize the Go module**

```bash
go mod init github.com/yourusername/pomocli
```

Expected: `go.mod` created with `module github.com/yourusername/pomocli`

**Step 2: Add dependencies**

```bash
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
```

Expected: `go.sum` created, dependencies downloaded.

**Step 3: Create `main.go`**

```go
package main

import "github.com/yourusername/pomocli/cmd"

func main() {
	cmd.Execute()
}
```

**Step 4: Create `cmd/root.go`**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pomocli",
	Short: "A Pomodoro timer for your terminal",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

**Step 5: Create empty stub files**

```bash
mkdir -p internal/config internal/tui internal/log internal/notify
touch internal/config/config.go
touch internal/tui/timer.go
touch internal/log/writer.go
touch internal/notify/notify.go
```

Add `package config`, `package tui`, `package log`, `package notify` headers to each.

**Step 6: Verify it builds**

```bash
go build ./...
```

Expected: no errors, binary created.

**Step 7: Commit**

```bash
git add .
git commit -m "feat: initialize Go module and project structure"
```

---

## Task 2: Config module

**Files:**
- Modify: `internal/config/config.go`
- Create: `internal/config/config_test.go`

**Step 1: Write the failing tests**

```go
// internal/config/config_test.go
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/pomocli/internal/config"
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
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/config/... -v
```

Expected: FAIL — `config.Defaults` and `config.Load` not defined.

**Step 3: Implement `internal/config/config.go`**

```go
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type PomodoroConfig struct {
	WorkDuration   int `mapstructure:"work_duration"`
	ShortBreak     int `mapstructure:"short_break"`
	LongBreak      int `mapstructure:"long_break"`
	LongBreakAfter int `mapstructure:"long_break_after"`
}

type NotificationsConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Sound     bool   `mapstructure:"sound"`
	SoundFile string `mapstructure:"sound_file"`
}

type LogsConfig struct {
	Directory  string `mapstructure:"directory"`
	DateFormat string `mapstructure:"date_format"`
}

type DisplayConfig struct {
	Theme           string `mapstructure:"theme"`
	ShowProgressBar bool   `mapstructure:"show_progress_bar"`
}

type Config struct {
	Pomodoro      PomodoroConfig      `mapstructure:"pomodoro"`
	Notifications NotificationsConfig `mapstructure:"notifications"`
	Logs          LogsConfig          `mapstructure:"logs"`
	Display       DisplayConfig       `mapstructure:"display"`
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("pomodoro.work_duration", 25)
	v.SetDefault("pomodoro.short_break", 5)
	v.SetDefault("pomodoro.long_break", 15)
	v.SetDefault("pomodoro.long_break_after", 4)
	v.SetDefault("notifications.enabled", true)
	v.SetDefault("notifications.sound", true)
	v.SetDefault("notifications.sound_file", "")
	v.SetDefault("logs.directory", filepath.Join(os.Getenv("HOME"), ".local/share/pomocli/logs"))
	v.SetDefault("logs.date_format", "2006-01-02")
	v.SetDefault("display.theme", "default")
	v.SetDefault("display.show_progress_bar", true)
}

// Defaults returns a Config populated with default values.
func Defaults() Config {
	v := viper.New()
	setDefaults(v)
	var cfg Config
	v.Unmarshal(&cfg)
	return cfg
}

// Load reads a config file from path, falling back to defaults for missing fields.
func Load(path string) (Config, error) {
	v := viper.New()
	setDefaults(v)
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// LoadDefault loads config from ~/.config/pomocli/config.yaml.
// Returns defaults if the file doesn't exist.
func LoadDefault() (Config, error) {
	path := filepath.Join(os.Getenv("HOME"), ".config/pomocli/config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Defaults(), nil
	}
	return Load(path)
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/config/... -v
```

Expected: PASS — all tests green.

**Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat: add config module with viper and defaults"
```

---

## Task 3: Log writer module

**Files:**
- Modify: `internal/log/writer.go`
- Create: `internal/log/writer_test.go`

**Step 1: Write the failing tests**

```go
// internal/log/writer_test.go
package log_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pomlog "github.com/yourusername/pomocli/internal/log"
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

	content, _ := os.ReadFile(expected)
	if !strings.Contains(string(content), "# Pomodoros — 2026-02-26") {
		t.Errorf("missing day header in file content")
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
		pomlog.WriteSession(dir, s)
	}

	content, _ := os.ReadFile(filepath.Join(dir, "2026-02-26.md"))
	if strings.Count(string(content), "## Sesión") != 2 {
		t.Errorf("expected 2 sessions in file")
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/log/... -v
```

Expected: FAIL — types and functions not defined.

**Step 3: Implement `internal/log/writer.go`**

```go
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
	return fmt.Sprintf(`## Sesión %d — %s

| Campo      | Valor                |
|------------|----------------------|
| Tarea      | %-20s |
| Inicio     | %-20s |
| Fin        | %-20s |
| Duración   | %-20s |
| Estado     | %-20s |

`,
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

	dayHeader := fmt.Sprintf("# Pomodoros — %s\n\n", s.StartTime.Format("2006-01-02"))

	if existing == "" {
		sb.WriteString(dayHeader)
		sb.WriteString("---\n\n")
	} else {
		sb.WriteString(existing)
	}

	sb.WriteString(FormatSession(s))

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/log/... -v
```

Expected: PASS.

**Step 5: Commit**

```bash
git add internal/log/
git commit -m "feat: add log writer module with daily .md files"
```

---

## Task 4: Notify module

**Files:**
- Modify: `internal/notify/notify.go`

> Note: Notifications are OS-specific side effects — no unit tests needed. Manual test at end.

**Step 1: Implement `internal/notify/notify.go`**

```go
package notify

import (
	"os/exec"
	"runtime"
)

// Send sends an OS desktop notification with the given title and body.
// Silently does nothing if the notification tool is unavailable.
func Send(title, body string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("notify-send", title, body).Run()
	case "darwin":
		script := `display notification "` + body + `" with title "` + title + `"`
		exec.Command("osascript", "-e", script).Run()
	}
}

// Beep plays the system bell or a sound file.
// If soundFile is empty, falls back to the terminal bell character.
func Beep(soundFile string) {
	if soundFile != "" {
		switch runtime.GOOS {
		case "linux":
			if exec.Command("paplay", soundFile).Run() != nil {
				exec.Command("aplay", soundFile).Run()
			}
		case "darwin":
			exec.Command("afplay", soundFile).Run()
		}
		return
	}
	// Terminal bell fallback
	print("\a")
}
```

**Step 2: Verify it compiles**

```bash
go build ./...
```

Expected: no errors.

**Step 3: Commit**

```bash
git add internal/notify/
git commit -m "feat: add notify module for OS notifications and sound"
```

---

## Task 5: TUI timer — state machine (TDD)

**Files:**
- Modify: `internal/tui/timer.go`
- Create: `internal/tui/timer_test.go`

**Step 1: Write the failing tests for the state machine**

```go
// internal/tui/timer_test.go
package tui_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/pomocli/internal/config"
	"github.com/yourusername/pomocli/internal/tui"
)

func defaultConfig() config.Config {
	cfg := config.Defaults()
	// Use short durations for tests
	cfg.Pomodoro.WorkDuration = 1   // 1 minute
	cfg.Pomodoro.ShortBreak = 1
	cfg.Pomodoro.LongBreakAfter = 2
	return cfg
}

func TestInitialState(t *testing.T) {
	m := tui.NewTimer("test task", defaultConfig())
	if m.State() != tui.StateWorking {
		t.Errorf("expected StateWorking, got %v", m.State())
	}
	if m.Paused() {
		t.Error("expected not paused initially")
	}
}

func TestPauseToggle(t *testing.T) {
	m := tui.NewTimer("task", defaultConfig())

	updated, _ := m.Update(tui.KeyMsg("p"))
	m2 := updated.(tui.TimerModel)
	if !m2.Paused() {
		t.Error("expected paused after pressing p")
	}

	updated2, _ := m2.Update(tui.KeyMsg("p"))
	m3 := updated2.(tui.TimerModel)
	if m3.Paused() {
		t.Error("expected unpaused after pressing p again")
	}
}

func TestTickDecreasesTime(t *testing.T) {
	m := tui.NewTimer("task", defaultConfig())
	before := m.Remaining()

	updated, _ := m.Update(tui.TickMsg())
	m2 := updated.(tui.TimerModel)

	if m2.Remaining() != before-time.Second {
		t.Errorf("expected remaining to decrease by 1s, got %v -> %v", before, m2.Remaining())
	}
}

func TestTickDoesNotDecreaseWhenPaused(t *testing.T) {
	m := tui.NewTimer("task", defaultConfig())
	updated, _ := m.Update(tui.KeyMsg("p"))
	m2 := updated.(tui.TimerModel)
	before := m2.Remaining()

	updated2, _ := m2.Update(tui.TickMsg())
	m3 := updated2.(tui.TimerModel)

	if m3.Remaining() != before {
		t.Error("expected remaining to not change when paused")
	}
}

func TestSkipFromWorkingGoesToShortBreak(t *testing.T) {
	m := tui.NewTimer("task", defaultConfig())
	updated, _ := m.Update(tui.KeyMsg("s"))
	m2 := updated.(tui.TimerModel)
	if m2.State() != tui.StateShortBreak {
		t.Errorf("expected StateShortBreak after skip, got %v", m2.State())
	}
}

func TestAfterLongBreakAfterPomodoros_GoesToLongBreak(t *testing.T) {
	cfg := defaultConfig()
	cfg.Pomodoro.LongBreakAfter = 2
	m := tui.NewTimer("task", cfg)

	// Skip 2 work sessions
	updated, _ := m.Update(tui.KeyMsg("s")) // -> short break
	m2 := updated.(tui.TimerModel)
	updated2, _ := m2.Update(tui.KeyMsg("s")) // -> working (pomodoro 2)
	m3 := updated2.(tui.TimerModel)
	updated3, _ := m3.Update(tui.KeyMsg("s")) // -> long break (after 2nd pomodoro)
	m4 := updated3.(tui.TimerModel)

	if m4.State() != tui.StateLongBreak {
		t.Errorf("expected StateLongBreak after %d pomodoros, got %v", cfg.Pomodoro.LongBreakAfter, m4.State())
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./internal/tui/... -v
```

Expected: FAIL — types not defined.

**Step 3: Implement the timer model in `internal/tui/timer.go`**

```go
package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/pomocli/internal/config"
	"github.com/yourusername/pomocli/internal/notify"
)

type State int

const (
	StateWorking State = iota
	StateShortBreak
	StateLongBreak
)

func (s State) String() string {
	switch s {
	case StateWorking:
		return "Trabajando"
	case StateShortBreak:
		return "Descanso corto"
	case StateLongBreak:
		return "Descanso largo"
	}
	return "Desconocido"
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// TickMsg returns a tickMsg for testing purposes.
func TickMsg() tea.Msg {
	return tickMsg(time.Now())
}

// KeyMsg returns a tea.KeyMsg for testing purposes.
func KeyMsg(key string) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

// TimerModel is the bubbletea model for the Pomodoro timer.
type TimerModel struct {
	state      State
	remaining  time.Duration
	total      time.Duration
	pomCount   int
	task       string
	paused     bool
	cfg        config.Config
	done       bool
	lastStatus string // for log writer
}

// NewTimer creates a new TimerModel ready to start working.
func NewTimer(task string, cfg config.Config) TimerModel {
	total := time.Duration(cfg.Pomodoro.WorkDuration) * time.Minute
	return TimerModel{
		state:     StateWorking,
		remaining: total,
		total:     total,
		task:      task,
		cfg:       cfg,
	}
}

// Accessors for tests
func (m TimerModel) State() State          { return m.state }
func (m TimerModel) Paused() bool          { return m.paused }
func (m TimerModel) Remaining() time.Duration { return m.remaining }
func (m TimerModel) PomCount() int         { return m.pomCount }
func (m TimerModel) Done() bool            { return m.done }

func (m TimerModel) Init() tea.Cmd {
	return doTick()
}

func (m TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			m.paused = !m.paused
		case "s":
			m = m.advance(false)
		case "q", "ctrl+c":
			m.done = true
			return m, tea.Quit
		}

	case tickMsg:
		if !m.paused {
			if m.remaining > 0 {
				m.remaining -= time.Second
			}
			if m.remaining == 0 {
				m = m.advance(true)
			}
		}
		return m, doTick()
	}

	return m, nil
}

func (m TimerModel) advance(fromTick bool) TimerModel {
	switch m.state {
	case StateWorking:
		m.pomCount++
		if fromTick && m.cfg.Notifications.Enabled {
			notify.Send("pomocli", "¡Pomodoro completado! Hora de descansar.")
			if m.cfg.Notifications.Sound {
				notify.Beep(m.cfg.Notifications.SoundFile)
			}
		}
		if m.pomCount%m.cfg.Pomodoro.LongBreakAfter == 0 {
			m.state = StateLongBreak
			m.total = time.Duration(m.cfg.Pomodoro.LongBreak) * time.Minute
		} else {
			m.state = StateShortBreak
			m.total = time.Duration(m.cfg.Pomodoro.ShortBreak) * time.Minute
		}
	case StateShortBreak, StateLongBreak:
		if fromTick && m.cfg.Notifications.Enabled {
			notify.Send("pomocli", "¡Descanso terminado! A trabajar.")
			if m.cfg.Notifications.Sound {
				notify.Beep(m.cfg.Notifications.SoundFile)
			}
		}
		m.state = StateWorking
		m.total = time.Duration(m.cfg.Pomodoro.WorkDuration) * time.Minute
	}
	m.remaining = m.total
	m.paused = false
	return m
}

// --- Rendering ---

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	timerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).
			MarginTop(1).MarginBottom(1)
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 3)
)

func (m TimerModel) View() string {
	task := m.task
	if task == "" {
		task = "sin tarea"
	}

	header := titleStyle.Render(
		fmt.Sprintf("🍅 Pomodoro #%d  [%s]", m.pomCount+1, task),
	)
	if m.state != StateWorking {
		header = titleStyle.Render(fmt.Sprintf("☕ %s", m.state.String()))
	}

	mins := int(m.remaining.Minutes())
	secs := int(m.remaining.Seconds()) % 60
	timerStr := timerStyle.Render(fmt.Sprintf("%02d:%02d", mins, secs))

	progress := ""
	if m.cfg.Display.ShowProgressBar {
		progress = renderProgressBar(m.remaining, m.total) + "\n"
	}

	pauseNote := ""
	if m.paused {
		pauseNote = helpStyle.Render("  ⏸ pausado") + "\n"
	}

	help := helpStyle.Render("[p] pausar  [s] saltar  [q] salir")

	content := fmt.Sprintf("%s\n\n%s\n%s%s%s", header, timerStr, progress, pauseNote, help)
	return borderStyle.Render(content)
}

func renderProgressBar(remaining, total time.Duration) string {
	const width = 20
	if total == 0 {
		return ""
	}
	elapsed := total - remaining
	filled := int(float64(elapsed) / float64(total) * float64(width))
	pct := int(float64(elapsed) / float64(total) * 100)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("%s  %d%%", bar, pct)
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./internal/tui/... -v
```

Expected: PASS.

**Step 5: Commit**

```bash
git add internal/tui/
git commit -m "feat: add TUI timer model with bubbletea and lipgloss"
```

---

## Task 6: `start` subcommand

**Files:**
- Create: `cmd/start.go`

**Step 1: Create `cmd/start.go`**

```go
package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yourusername/pomocli/internal/config"
	pomlog "github.com/yourusername/pomocli/internal/log"
	"github.com/yourusername/pomocli/internal/tui"
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
	if m.Done() && m.PomCount() > 0 {
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
```

**Step 2: Build and smoke-test**

```bash
go build -o pomocli . && ./pomocli start --task "test task"
```

Expected: TUI launches, shows timer, responds to p/s/q keys.

**Step 3: Commit**

```bash
git add cmd/start.go
git commit -m "feat: add start subcommand with TUI and log writing"
```

---

## Task 7: `log` and `config` subcommands

**Files:**
- Create: `cmd/log.go`
- Create: `cmd/config_cmd.go`

**Step 1: Create `cmd/log.go`**

```go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/pomocli/internal/config"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show the Pomodoro log for a given day",
	RunE:  runLog,
}

var dateFlag string

func init() {
	logCmd.Flags().StringVarP(&dateFlag, "date", "d", "", "Date in YYYY-MM-DD format (default: today)")
	rootCmd.AddCommand(logCmd)
}

func runLog(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadDefault()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	date := dateFlag
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	filename := filepath.Join(cfg.Logs.Directory, date+".md")
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		fmt.Printf("No hay log para %s\n", date)
		return nil
	}
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}

	fmt.Println(string(data))
	return nil
}
```

**Step 2: Create `cmd/config_cmd.go`**

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/pomocli/internal/config"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the active configuration",
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadDefault()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	return nil
}
```

> Note: `yaml.v3` needs to be added: `go get gopkg.in/yaml.v3`

**Step 3: Add yaml.v3 dependency**

```bash
go get gopkg.in/yaml.v3
```

**Step 4: Build and test commands**

```bash
go build -o pomocli .
./pomocli log
./pomocli log --date 2026-02-26
./pomocli config
```

Expected: `log` shows today's log or "No hay log", `config` prints YAML config.

**Step 5: Commit**

```bash
git add cmd/log.go cmd/config_cmd.go go.mod go.sum
git commit -m "feat: add log and config subcommands"
```

---

## Task 8: Default config file scaffolding

**Files:**
- Create: `config.example.yaml`

**Step 1: Create the example config**

```yaml
# ~/.config/pomocli/config.yaml
# Copy this file to that location and customize as needed.

pomodoro:
  work_duration: 25      # minutes
  short_break: 5
  long_break: 15
  long_break_after: 4

notifications:
  enabled: true
  sound: true
  sound_file: ""         # empty = system bell

logs:
  directory: "~/.local/share/pomocli/logs"
  date_format: "2006-01-02"

display:
  theme: "default"
  show_progress_bar: true
```

**Step 2: Commit**

```bash
git add config.example.yaml
git commit -m "docs: add example config file"
```

---

## Task 9: Run full test suite and verify

**Step 1: Run all tests**

```bash
go test ./... -v
```

Expected: all tests PASS.

**Step 2: Run the linter**

```bash
go vet ./...
```

Expected: no warnings.

**Step 3: Build final binary**

```bash
go build -o pomocli .
```

**Step 4: End-to-end smoke test**

```bash
# Check help
./pomocli --help
./pomocli start --help

# Check config
./pomocli config

# Start a short session (edit config to set work_duration: 1 for quick test)
./pomocli start --task "smoke test"

# View the log
./pomocli log
```

Expected: TUI shows, timer counts down, log file created at `~/.local/share/pomocli/logs/YYYY-MM-DD.md`.

**Step 5: Final commit**

```bash
git add .
git commit -m "chore: verify build and tests pass"
```
