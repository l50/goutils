package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"log/slog"

	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/afero"
)

// CreateLogFile creates a log file in a 'logs' subdirectory of the
// specified directory. The log file's name is the provided log name
// with the extension '.log'.
//
// **Parameters:**
//
// fs: An afero.Fs instance to mock filesystem for testing.
// logDir: A string for the directory where 'logs' subdirectory and
// log file should be created.
// logName: A string for the name of the log file to be created.
//
// **Returns:**
//
// LogInfo: A LogInfo struct with information about the log file,
// including its directory, file pointer, file name, and path.
// error: An error, if an issue occurs while creating the directory
// or the log file.
func CreateLogFile(fs afero.Fs, logDir, logName string) (LogInfo, error) {
	logName = strings.TrimSpace(logName)
	logDir = strings.TrimSpace(logDir)

	if logDir == "" {
		return LogInfo{}, fmt.Errorf("logDir cannot be empty")
	}
	if logName == "" {
		return LogInfo{}, fmt.Errorf("logName cannot be empty")
	}
	if filepath.Ext(logName) != ".log" {
		logName += ".log"
	}

	logInfo := LogInfo{
		Dir:      filepath.Join(logDir, "logs"),
		FileName: logName,
		Path:     filepath.Join(logDir, "logs", logName),
	}

	if _, err := fs.Stat(logInfo.Path); os.IsNotExist(err) {
		if err := fs.MkdirAll(logInfo.Dir, os.ModePerm); err != nil {
			return LogInfo{}, fmt.Errorf("failed to create %s: %v", logInfo.Dir, err)
		}
	}

	file, err := fs.OpenFile(logInfo.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return LogInfo{}, fmt.Errorf("failed to create %s: %v", logInfo.Path, err)
	}
	logInfo.File = file

	return logInfo, nil
}

// ConfigureLogger creates a logger based on the provided level.
// Depending on the level, it returns a color or plain logger.
//
// **Parameters:**
//
// level: Logging level as a slog.Level.
// path: Path to the log file.
//
// **Returns:**
//
// Logger: Logger object based on provided level.
// error: An error, if an issue occurs while setting up the logger.
func ConfigureLogger(level slog.Level, path string) (Logger, error) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	fileHandler := slog.NewJSONHandler(logFile, opts)

	// Using PrettyHandler for all levels
	prettyOpts := PrettyHandlerOptions{SlogOpts: *opts}
	stdoutHandler := NewPrettyHandler(os.Stdout, prettyOpts)

	handler := slogmulti.Fanout(fileHandler, stdoutHandler)
	logger := slog.New(handler)

	colorAttribute := determineColorAttribute(level)
	return &ColorLogger{
		Info:           LogInfo{File: logFile, Path: path},
		ColorAttribute: colorAttribute,
		Logger:         logger,
	}, nil
}

func extractFields(r slog.Record) map[string]interface{} {
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	return fields
}
