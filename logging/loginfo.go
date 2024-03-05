package logging

import (
	"log/slog"
	"strings"

	"github.com/spf13/afero"
)

// LogConfig represents parameters used to manage logging throughout
// a program.
//
// **Attributes:**
//
// Fs: An afero.Fs object representing the file system.
// Path: A string representing the full path to the log file.
// Level: A slog.Level object representing the logging level.
// LogToDisk: A boolean representing whether or not to log to disk.
type LogConfig struct {
	Fs         afero.Fs
	LogPath    string
	Level      slog.Level
	OutputType OutputType
	LogToDisk  bool
}

// DetermineLogLevel determines the log level from a given string.
//
// **Parameters:**
//
// levelStr: A string representing the log level.
//
// **Returns:**
//
// slog.Level: The corresponding slog.Level for the given log level string.
func DetermineLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to info level
	}
}
