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
func (cfg *LogConfig) CreateLogFile() error {
	cfg.LogPath = strings.TrimSpace(cfg.LogPath)
	if cfg.LogPath == "" {
		return fmt.Errorf("logPath cannot be empty")
	}
	if filepath.Ext(cfg.LogPath) != ".log" {
		cfg.LogPath += ".log"
	}

	if _, err := cfg.Fs.Stat(cfg.LogPath); os.IsNotExist(err) {
		if err := cfg.Fs.MkdirAll(filepath.Dir(cfg.LogPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create %s: %v", filepath.Dir(cfg.LogPath), err)
		}
	}

	// Check if the file can be opened (created if not exists), then close it immediately
	file, err := cfg.Fs.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", cfg.LogPath, err)
	}
	file.Close()

	return nil
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
func (cfg *LogConfig) ConfigureLogger() (Logger, error) {
	dir := filepath.Dir(cfg.LogPath)
	if _, err := cfg.Fs.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid path: %s", dir)
	}

	var logFile afero.File
	var err error
	var fileHandler slog.Handler
	var stdoutHandler slog.Handler

	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	if cfg.LogToDisk {
		logFile, err = cfg.Fs.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		fileHandler = slog.NewJSONHandler(logFile, opts)
	}

	if cfg.OutputType == ColorOutput {
		prettyOpts := PrettyHandlerOptions{SlogOpts: *opts}
		stdoutHandler = NewPrettyHandler(os.Stdout, prettyOpts)
	} else {
		stdoutHandler = slog.NewJSONHandler(os.Stdout, opts)
	}

	var handlers []slog.Handler
	if fileHandler != nil {
		handlers = append(handlers, fileHandler)
	}
	if stdoutHandler != nil {
		handlers = append(handlers, stdoutHandler)
	}

	if len(handlers) == 0 {
		return nil, fmt.Errorf("no valid handlers available for logger")
	}

	multiHandler := slog.New(slogmulti.Fanout(handlers...))
	var logger Logger
	if cfg.OutputType == ColorOutput {
		colorAttribute := determineColorAttribute(cfg.Level)
		logger = &ColorLogger{
			Cfg:            *cfg,
			ColorAttribute: colorAttribute,
			Logger:         multiHandler,
		}
	} else {
		logger = &PlainLogger{
			Info:   *cfg,
			Logger: multiHandler,
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
func InitLogging(cfg *LogConfig) (Logger, error) {
	if cfg.LogToDisk {
		if cfg.LogPath == "" {
			return nil, fmt.Errorf("logPath cannot be empty when logToDisk is true")
		}

		if err := cfg.CreateLogFile(); err != nil {
			return nil, fmt.Errorf("failed to create log file: %v", err)
		}
	}

	logger, err := cfg.ConfigureLogger()
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
