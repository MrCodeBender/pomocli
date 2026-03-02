package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/angelchaudhary01/pomocli/internal/config"
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
