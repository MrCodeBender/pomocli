# pomocli — Diseño

**Fecha:** 2026-02-26
**Estado:** Aprobado

## Descripción

`pomocli` es una aplicación CLI de consola para gestionar sesiones Pomodoro. Incluye una TUI
interactiva, archivo de configuración YAML, notificaciones del sistema operativo y registro
diario en archivos `.md`.

---

## Stack Tecnológico

| Componente     | Librería/Tool              |
|----------------|----------------------------|
| Lenguaje       | Go                         |
| TUI framework  | bubbletea (Charm.sh)       |
| Estilos TUI    | lipgloss (Charm.sh)        |
| Configuración  | viper                      |
| CLI/subcomandos| cobra                      |
| Notificaciones | notify-send (Linux) / osascript (Mac) + beep |

---

## Arquitectura

```
pomocli/
├── cmd/
│   └── root.go          # cobra root + subcomandos: start, log, config
├── internal/
│   ├── tui/
│   │   └── timer.go     # bubbletea model para el timer TUI
│   ├── config/
│   │   └── config.go    # carga/valida config con viper
│   ├── log/
│   │   └── writer.go    # escribe registros en archivos .md
│   └── notify/
│       └── notify.go    # notificaciones OS + sonido
├── config.yaml          # config por defecto
├── go.mod
├── go.sum
└── main.go
```

### Subcomandos

- `pomocli start [--task "nombre"]` — inicia un pomodoro con TUI interactiva
- `pomocli log [--date 2026-02-26]` — muestra el log del día (por defecto hoy)
- `pomocli config` — muestra la configuración activa

---

## Archivo de Configuración

Ubicación: `~/.config/pomocli/config.yaml` (o junto al binario)

```yaml
pomodoro:
  work_duration: 25        # minutos de trabajo
  short_break: 5           # descanso corto
  long_break: 15           # descanso largo
  long_break_after: 4      # pomodoros antes del descanso largo

notifications:
  enabled: true
  sound: true
  sound_file: ""           # vacío = beep del sistema

logs:
  directory: "~/.local/share/pomocli/logs"
  date_format: "2006-01-02"

display:
  theme: "default"
  show_progress_bar: true
```

---

## TUI — Flujo de Sesión

### Estados

```
Idle → Working → ShortBreak → Working → ... (x4) → LongBreak → Working
```

### Pantalla durante un pomodoro

```
╭─────────────────────────────────────╮
│  pomocli                            │
│                                     │
│  🍅 Pomodoro #3  [revisar PRs]      │
│                                     │
│        24:13                        │
│                                     │
│  ████████████░░░░░░░  60%           │
│                                     │
│  [p] pausar  [s] saltar  [q] salir  │
╰─────────────────────────────────────╯
```

### Keybindings

| Tecla | Acción                              |
|-------|-------------------------------------|
| `p`   | Pausar / reanudar el timer          |
| `s`   | Saltar al siguiente estado          |
| `q`   | Terminar sesión y guardar log       |

---

## Formato de Log `.md`

Archivo: `~/.local/share/pomocli/logs/2026-02-26.md`

```markdown
# Pomodoros — 2026-02-26

## Resumen
- Total completados: 4
- Tiempo de trabajo: 1h 40m
- Tiempo de descanso: 20m

---

## Sesión 1 — 09:15

| Campo      | Valor                |
|------------|----------------------|
| Tarea      | revisar PRs          |
| Inicio     | 09:15                |
| Fin        | 09:40                |
| Duración   | 25m                  |
| Estado     | completado           |

## Sesión 2 — 09:45

| Campo      | Valor                |
|------------|----------------------|
| Tarea      | —                    |
| Inicio     | 09:45                |
| Fin        | 10:07                |
| Duración   | 22m                  |
| Estado     | cancelado            |
```

---

## Notificaciones

- Al terminar un pomodoro: notificación OS + sonido
- Al terminar un descanso: notificación OS + sonido
- Linux: `notify-send` + `paplay`/`aplay` o beep ANSI
- macOS: `osascript` + `afplay`
- Configurable en `config.yaml` (se puede deshabilitar)

---

## Decisiones de Diseño

1. **Un archivo .md por día** — facilita revisión diaria, archivos pequeños y manejables
2. **Tarea opcional** — no obliga al usuario, pero incentiva el registro de contexto
3. **bubbletea sobre tview** — mejor soporte para animaciones custom y comunidad más activa
4. **viper para config** — soporta YAML/TOML/ENV vars con valores por defecto automáticos
5. **cobra para subcomandos** — estándar de facto en CLIs Go, `--help` y flags gratis
