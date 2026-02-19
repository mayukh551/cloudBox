package utils

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Global logger instance
var Log *slog.Logger

// Initialize logger once when package is imported
func init() {
	var err error
	Log, err = setupLogger("app.log")
	if err != nil {
		// Fallback to default logger if file fails
		Log = slog.Default()
		slog.Error("Failed to initialize file logger, using default", "error", err)
	}
}

func setupLogger(filePath string) (*slog.Logger, error) {
	if filePath == "" {
		filePath = "app.log"
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create log file: %w", err)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				src := a.Value.Any().(*slog.Source)
				pathArr := strings.Split(src.File, "/")
				l := len(pathArr)
				a.Value = slog.StringValue(strings.Join(pathArr[l-2:l], "/") + ":" + fmt.Sprint(src.Line))
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(file, opts)
	return slog.New(handler), nil
}

// Optional: Allow reconfiguration if needed
func Configure(filePath string) error {
	newLogger, err := setupLogger(filePath)
	if err != nil {
		return err
	}
	Log = newLogger
	return nil
}
