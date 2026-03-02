package cmd

import (
	"errors"
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

	const dateLayout = "2006-01-02"

	date := dateFlag
	if date == "" {
		date = time.Now().Format(dateLayout)
	} else {
		if _, err := time.Parse(dateLayout, date); err != nil {
			return fmt.Errorf("invalid date %q: expected YYYY-MM-DD format", date)
		}
	}

	filename := filepath.Join(cfg.Logs.Directory, date+".md")
	data, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("No hay log para %s\n", date)
		return nil
	}
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}

	fmt.Println(string(data))
	return nil
}
