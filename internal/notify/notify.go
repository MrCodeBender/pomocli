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
		exec.Command("notify-send", title, body).Run() //nolint:errcheck
	case "darwin":
		script := `display notification "` + body + `" with title "` + title + `"`
		exec.Command("osascript", "-e", script).Run() //nolint:errcheck
	}
}

// Beep plays the system bell or a sound file.
// If soundFile is empty, falls back to the terminal bell character.
func Beep(soundFile string) {
	if soundFile != "" {
		switch runtime.GOOS {
		case "linux":
			if exec.Command("paplay", soundFile).Run() != nil {
				exec.Command("aplay", soundFile).Run() //nolint:errcheck
			}
		case "darwin":
			exec.Command("afplay", soundFile).Run() //nolint:errcheck
		}
		return
	}
	// Terminal bell fallback
	print("\a")
}
