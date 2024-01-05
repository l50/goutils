package logging

import (
	"log/slog"

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
