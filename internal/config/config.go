package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type PomodoroConfig struct {
	WorkDuration   int `mapstructure:"work_duration"    yaml:"work_duration"`
	ShortBreak     int `mapstructure:"short_break"      yaml:"short_break"`
	LongBreak      int `mapstructure:"long_break"       yaml:"long_break"`
	LongBreakAfter int `mapstructure:"long_break_after" yaml:"long_break_after"`
}

type NotificationsConfig struct {
	Enabled   bool   `mapstructure:"enabled"    yaml:"enabled"`
	Sound     bool   `mapstructure:"sound"      yaml:"sound"`
	SoundFile string `mapstructure:"sound_file" yaml:"sound_file"`
}

type LogsConfig struct {
	Directory  string `mapstructure:"directory"   yaml:"directory"`
	DateFormat string `mapstructure:"date_format" yaml:"date_format"`
}

type DisplayConfig struct {
	Theme           string `mapstructure:"theme"             yaml:"theme"`
	ShowProgressBar bool   `mapstructure:"show_progress_bar" yaml:"show_progress_bar"`
}

type Config struct {
	Pomodoro      PomodoroConfig      `mapstructure:"pomodoro"      yaml:"pomodoro"`
	Notifications NotificationsConfig `mapstructure:"notifications" yaml:"notifications"`
	Logs          LogsConfig          `mapstructure:"logs"          yaml:"logs"`
	Display       DisplayConfig       `mapstructure:"display"       yaml:"display"`
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("pomodoro.work_duration", 25)
	v.SetDefault("pomodoro.short_break", 5)
	v.SetDefault("pomodoro.long_break", 15)
	v.SetDefault("pomodoro.long_break_after", 4)
	v.SetDefault("notifications.enabled", true)
	v.SetDefault("notifications.sound", true)
	v.SetDefault("notifications.sound_file", "")
	v.SetDefault("logs.directory", "") // set dynamically in callers
	v.SetDefault("logs.date_format", "2006-01-02")
	v.SetDefault("display.theme", "default")
	v.SetDefault("display.show_progress_bar", true)
}

// Defaults returns a Config populated with default values.
func Defaults() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME") // fallback
	}
	v := viper.New()
	setDefaults(v)
	v.SetDefault("logs.directory", filepath.Join(home, ".local/share/pomocli/logs"))
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic("config: failed to unmarshal defaults: " + err.Error())
	}
	return cfg
}

// Load reads a config file from path, falling back to defaults for missing fields.
func Load(path string) (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	v := viper.New()
	setDefaults(v)
	v.SetDefault("logs.directory", filepath.Join(home, ".local/share/pomocli/logs"))
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
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	path := filepath.Join(home, ".config/pomocli/config.yaml")
	if _, statErr := os.Stat(path); errors.Is(statErr, os.ErrNotExist) {
		return Defaults(), nil
	}
	return Load(path)
}
