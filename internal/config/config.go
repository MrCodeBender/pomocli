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
