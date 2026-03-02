# pomocli

A Pomodoro timer for your terminal — interactive TUI, YAML config, and daily `.md` logs.

```
╭──────────────────────────────────────────╮
│                                          │
│  🍅 Pomodoro #1  [revisar PRs]           │
│                                          │
│               24:13                      │
│                                          │
│  ████████████░░░░░░░░  60%               │
│                                          │
│  [p] pausar  [s] saltar  [q] salir       │
│                                          │
╰──────────────────────────────────────────╯
```

## Features

- Interactive TUI with live countdown and progress bar
- Pause, skip, and quit with single keypresses
- Automatic short/long break rotation
- Desktop notifications + sound when a session ends
- Configurable via a single YAML file
- Saves each session to a daily `.md` log file

## Installation

**Requirements:** Go 1.22+

```bash
git clone https://github.com/angelchaudhary01/pomocli
cd pomocli
go build -o pomocli .
```

Optionally move the binary to your PATH:

```bash
mv pomocli ~/.local/bin/
```

## Usage

```bash
# Start a Pomodoro (with optional task name)
pomocli start
pomocli start --task "revisar PRs"
pomocli start -t "escribir tests"

# View today's log
pomocli log

# View log for a specific date
pomocli log --date 2026-02-26

# Show active configuration
pomocli config
```

### Keybindings (during a session)

| Key | Action          |
|-----|-----------------|
| `p` | Pause / resume  |
| `s` | Skip to next    |
| `q` | Quit and save   |

## Configuration

Copy the example config and customize it:

```bash
mkdir -p ~/.config/pomocli
cp config.example.yaml ~/.config/pomocli/config.yaml
```

`~/.config/pomocli/config.yaml`:

```yaml
pomodoro:
  work_duration: 25       # minutes of focused work
  short_break: 5          # short break between pomodoros
  long_break: 15          # long break after N pomodoros
  long_break_after: 4     # pomodoros before a long break

notifications:
  enabled: true
  sound: true
  sound_file: ""          # empty = system bell; or path to audio file

logs:
  directory: "~/.local/share/pomocli/logs"

display:
  show_progress_bar: true
```

If no config file exists, defaults are used automatically.

## Logs

Each session is appended to a daily file in `~/.local/share/pomocli/logs/`:

```
~/.local/share/pomocli/logs/
├── 2026-02-26.md
├── 2026-02-27.md
└── 2026-03-02.md
```

Example log entry (`2026-03-02.md`):

```markdown
# Pomodoros — 2026-03-02

---

## Sesión 1 — 09:15

| Campo      | Valor       |
|------------|-------------|
| Tarea      | revisar PRs |
| Inicio     | 09:15       |
| Fin        | 09:40       |
| Duración   | 25m         |
| Estado     | completado  |
```

## Pomodoro Cycle

```
Work (25m) → Short break (5m) → Work → Short break → Work → Short break → Work → Long break (15m) → repeat
```

The cycle length is configurable via `long_break_after`.

## License

MIT
