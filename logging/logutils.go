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
// LogConfig: A LogConfig struct with information about the log file,
// including its directory, file pointer, file name, and path.
// error: An error, if an issue occurs while creating the directory
// or the log file.
func CreateLogFile(fs afero.Fs, logPath string) (LogConfig, error) {
	logPath = strings.TrimSpace(logPath)

	if logPath == "" {
		return LogConfig{}, fmt.Errorf("logDir cannot be empty")
	}

	if filepath.Ext(logPath) != ".log" {
		logPath += ".log"
	}

	logInfo := LogConfig{
		Dir:      filepath.Dir(logPath),
		FileName: filepath.Base(logPath),
		Path:     logPath,
	}

	if _, err := fs.Stat(logInfo.Path); os.IsNotExist(err) {
		if err := fs.MkdirAll(logInfo.Dir, os.ModePerm); err != nil {
			return LogConfig{}, fmt.Errorf("failed to create %s: %v", logInfo.Dir, err)
		}
	}

	file, err := fs.OpenFile(logInfo.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return LogConfig{}, fmt.Errorf("failed to create %s: %v", logInfo.Path, err)
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
func ConfigureLogger(fs afero.Fs, level slog.Level, path string, outputType OutputType) (Logger, error) {
	// Check if the directory for the given path exists
	dir := filepath.Dir(path)
	if _, err := fs.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid path: %s", dir)
	}

	logFile, err := fs.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
		// Create a plain JSON handler for file logging to avoid color codes in the file
		fileHandler := slog.NewJSONHandler(logFile, opts)

		// PrettyHandler for colorized output in console
		prettyOpts := PrettyHandlerOptions{SlogOpts: *opts}
		stdoutHandler = NewPrettyHandler(os.Stdout, prettyOpts)

		colorAttribute := determineColorAttribute(level)
		logger = &ColorLogger{
			Info:           LogConfig{File: logFile, Path: path},
			ColorAttribute: colorAttribute,
			Logger:         slog.New(slogmulti.Fanout(fileHandler, stdoutHandler)),
		}

	case PlainOutput:
		// Standard JSON handler for PlainLogger without colorization
		stdoutHandler = slog.NewJSONHandler(os.Stdout, opts)
		logger = &PlainLogger{
			Info:   LogConfig{File: logFile, Path: path},
			Logger: slog.New(slogmulti.Fanout(fileHandler, stdoutHandler)),
		}
	}

	return logger, nil
}

// InitLogging is a convenience function that combines
// the CreateLogFile and ConfigureLogger functions into one call.
// It is useful for quickly setting up logging to disk.
//
// **Parameters:**
//
// fs: An afero.Fs instance for filesystem operations, allows mocking in tests.
// logPath: The path to the log file.
// level: The logging level.
// outputType: The output type of the logger (PlainOutput or ColorOutput).
// logToDisk: A boolean indicating whether to log to disk or not.
//
// **Returns:**
//
// Logger: A configured Logger object.
// error: An error if any issue occurs during initialization.
func InitLogging(fs afero.Fs, logPath string, level slog.Level, outputType OutputType, logToDisk bool) (Logger, error) {
	cfg := LogConfig{
		Path: logPath,
	}
	if logToDisk {
		if cfg.Path == "" {
			return nil, fmt.Errorf("logPath cannot be empty when logToDisk is true")
		}
		var err error
		cfg, err = CreateLogFile(fs, logPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %v", err)
		}
	}

	logger, err := ConfigureLogger(fs, level, cfg.Path, outputType)
	if err != nil {
		return nil, fmt.Errorf("failed to configure logger: %v", err)
	}

	return logger, nil
}

// LogAndReturnError logs the provided error message using the given logger and returns the error.
//
// This utility function is helpful for scenarios where an error needs to be both logged and returned.
// It simplifies the code by combining these two actions into one call.
//
// **Parameters:**
//
// logger: The Logger instance used for logging the error.
// errMsg: The error message to log and return.
//
// **Returns:**
//
// error: The error created from the errMsg, after it has been logged.
func LogAndReturnError(logger Logger, errMsg string) error {
	err := fmt.Errorf(errMsg)
	logger.Error(err)
	return err
}
