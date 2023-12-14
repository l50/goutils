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

// OutputType is an enumeration type that specifies the output format
// of the logger. It can be either plain text or colorized text.
type OutputType int

const (
	// PlainOutput indicates that the logger will produce plain text
	// output without any colorization. This is suitable for log
	// files or environments where ANSI color codes are not supported.
	PlainOutput OutputType = iota

	// ColorOutput indicates that the logger will produce colorized
	// text output. This is useful for console output where color
	// coding can enhance readability.
	ColorOutput
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

// ConfigureLogger sets up a logger based on the provided logging level,
// file path, and output type. It supports both colorized and plain text
// logging output, selectable via the OutputType parameter. The logger
// writes log entries to both a file and standard output.
//
// **Parameters:**
//
// level: Logging level as a slog.Level.
// path: Path to the log file.
// outputType: Type of log output, either ColorOutput or PlainOutput.
//
// **Returns:**
//
// Logger: Configured Logger object based on provided parameters.
// error: An error, if an issue occurs while setting up the logger.
func ConfigureLogger(level slog.Level, path string, outputType OutputType) (Logger, error) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	// File logger
	fileHandler := slog.NewJSONHandler(logFile, opts)

	var stdoutHandler slog.Handler
	var logger Logger

	switch outputType {
	case ColorOutput:
		// PrettyHandler for colorized output in console
		prettyOpts := PrettyHandlerOptions{SlogOpts: *opts}
		stdoutHandler = NewPrettyHandler(os.Stdout, prettyOpts)

		colorAttribute := determineColorAttribute(level)
		logger = &ColorLogger{
			Info:           LogInfo{File: logFile, Path: path},
			ColorAttribute: colorAttribute,
			Logger:         slog.New(slogmulti.Fanout(fileHandler, stdoutHandler)),
		}

	case PlainOutput:
		// Standard JSON handler for PlainLogger without colorization
		stdoutHandler = slog.NewJSONHandler(os.Stdout, opts)
		logger = &PlainLogger{
			Info:   LogInfo{File: logFile, Path: path},
			Logger: slog.New(slogmulti.Fanout(fileHandler, stdoutHandler)),
		}
	}

	return logger, nil
}
