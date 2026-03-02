package tui_test

import (
	"testing"
	"time"

	"github.com/angelchaudhary01/pomocli/internal/config"
	"github.com/angelchaudhary01/pomocli/internal/tui"
)

func defaultConfig() config.Config {
	cfg := config.Defaults()
	cfg.Pomodoro.WorkDuration = 1
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

	// Skip 2 work sessions to reach long break
	updated, _ := m.Update(tui.KeyMsg("s")) // working -> short break (pom 1)
	m2 := updated.(tui.TimerModel)
	updated2, _ := m2.Update(tui.KeyMsg("s")) // short break -> working
	m3 := updated2.(tui.TimerModel)
	updated3, _ := m3.Update(tui.KeyMsg("s")) // working -> long break (pom 2)
	m4 := updated3.(tui.TimerModel)

	if m4.State() != tui.StateLongBreak {
		t.Errorf("expected StateLongBreak after %d pomodoros, got %v", cfg.Pomodoro.LongBreakAfter, m4.State())
	}
}

func TestTickToZeroAdvancesState(t *testing.T) {
	cfg := defaultConfig()

	// Set remaining to 1 second by sending ticks until remaining = 1s
	// Instead, use the exported constructor and manually drain via ticks
	// We need to drain a 1-minute timer — use a very short config
	cfg.Pomodoro.WorkDuration = 1 // 1 minute = 60 ticks
	m2 := tui.NewTimer("task", cfg)

	// Send 59 ticks — should still be Working
	for i := 0; i < 59; i++ {
		updated, _ := m2.Update(tui.TickMsg())
		m2 = updated.(tui.TimerModel)
	}
	if m2.State() != tui.StateWorking {
		t.Errorf("expected still Working after 59 ticks, got %v", m2.State())
	}

	// Send 1 more tick — remaining hits 0, should advance to ShortBreak
	updated, _ := m2.Update(tui.TickMsg())
	m2 = updated.(tui.TimerModel)
	if m2.State() != tui.StateShortBreak {
		t.Errorf("expected ShortBreak after timer expires, got %v", m2.State())
	}
	if m2.PomCount() != 1 {
		t.Errorf("expected PomCount=1 after first pomodoro, got %d", m2.PomCount())
	}
}

func TestQuitSetsDone(t *testing.T) {
	m := tui.NewTimer("task", defaultConfig())
	if m.Done() {
		t.Error("expected Done=false initially")
	}
	updated, _ := m.Update(tui.KeyMsg("q"))
	m2 := updated.(tui.TimerModel)
	if !m2.Done() {
		t.Error("expected Done=true after pressing q")
	}
}
