package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/angelchaudhary01/pomocli/internal/config"
	"github.com/angelchaudhary01/pomocli/internal/notify"
)

// State represents the current phase of the Pomodoro cycle.
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

// tickMsg is sent every second to drive the countdown.
type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// TickMsg returns a tickMsg for use in tests.
func TickMsg() tea.Msg {
	return tickMsg(time.Now())
}

// KeyMsg returns a tea.KeyMsg for use in tests.
func KeyMsg(key string) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

// TimerModel is the bubbletea model for the Pomodoro timer.
type TimerModel struct {
	state     State
	remaining time.Duration
	total     time.Duration
	pomCount  int
	task      string
	paused    bool
	cfg       config.Config
	done      bool
}

// NewTimer creates a TimerModel in the Working state.
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

// Accessor methods (used by tests and cmd/start.go)
func (m TimerModel) State() State             { return m.state }
func (m TimerModel) Paused() bool             { return m.paused }
func (m TimerModel) Remaining() time.Duration { return m.remaining }
func (m TimerModel) PomCount() int            { return m.pomCount }
func (m TimerModel) Done() bool               { return m.done }

// Init implements tea.Model.
func (m TimerModel) Init() tea.Cmd {
	return doTick()
}

// Update implements tea.Model.
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

// advance transitions to the next state. fromTick=true means the timer
// completed naturally (sends notifications); fromTick=false means user skipped.
func (m TimerModel) advance(fromTick bool) TimerModel {
	switch m.state {
	case StateWorking:
		m.pomCount++
		if fromTick && m.cfg.Notifications.Enabled {
			notify.Send("pomocli", "Pomodoro completado! Hora de descansar.")
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
			notify.Send("pomocli", "Descanso terminado! A trabajar.")
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

// View implements tea.Model.
func (m TimerModel) View() string {
	task := m.task
	if task == "" {
		task = "sin tarea"
	}

	var header string
	if m.state == StateWorking {
		header = titleStyle.Render(fmt.Sprintf("Pomodoro #%d  [%s]", m.pomCount+1, task))
	} else {
		header = titleStyle.Render(fmt.Sprintf("%s", m.state.String()))
	}

	mins := int(m.remaining.Minutes())
	secs := int(m.remaining.Seconds()) % 60
	timerStr := timerStyle.Render(fmt.Sprintf("%02d:%02d", mins, secs))

	var progress string
	if m.cfg.Display.ShowProgressBar {
		progress = renderProgressBar(m.remaining, m.total) + "\n"
	}

	var pauseNote string
	if m.paused {
		pauseNote = helpStyle.Render("  pausado") + "\n"
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
	if filled > width {
		filled = width
	}
	pct := int(float64(elapsed) / float64(total) * 100)
	if pct > 100 {
		pct = 100
	}
	bar := strings.Repeat("#", filled) + strings.Repeat(".", width-filled)
	return fmt.Sprintf("%s  %d%%", bar, pct)
}
