package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/afero"
	"golang.org/x/exp/slog"
)

// Logger defines the methods for a generic logging interface.
// It consists of Println, Printf and Error methods.
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

// ColoredLogger represents a logger that outputs in cyan color.
//
// **Attributes:**
//
// Info: LogInfo object containing information about the log file.
// ColorAttribute: A color.Attribute object representing the color
type ColoredLogger struct {
	Info           LogInfo
	ColorAttribute color.Attribute
}

// Println for ColoredLogger logs the provided arguments as a line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColoredLogger) Println(v ...interface{}) {
	log.SetOutput(l.Info.File)
	log.Println(color.New(l.ColorAttribute).Sprint(v...))
}

// Printf for ColoredLogger logs the provided formatted string in
// the specified color. The format and arguments are handled in the
// manner of fmt.Printf.
func (l *ColoredLogger) Printf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	if len(v) > 0 {
		log.Println(color.New(l.ColorAttribute).Sprintf(format, v...))
	} else {
		log.Println(color.New(l.ColorAttribute).Sprint(format))
	}
}

// Error for ColoredLogger logs the provided arguments as an error line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColoredLogger) Error(v ...interface{}) {
	log.SetOutput(l.Info.File)
	log.Println(color.New(l.ColorAttribute).Add(color.Bold).Sprint(v...))
}

// Errorf for ColoredLogger logs the provided formatted string as an
// error line in the specified color. The format and arguments are handled
// in the manner of fmt.Printf.
func (l *ColoredLogger) Errorf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	if len(v) > 0 {
		log.Println(color.New(l.ColorAttribute).Add(color.Bold).Sprintf(format, v...))
	} else {
		log.Println(color.New(l.ColorAttribute).Add(color.Bold).Sprint(format))
	}
}

// PlainLogger represents a logger that outputs in plain format.
//
// **Attributes:**
//
// Info: LogInfo object containing information about the log file.
type PlainLogger struct {
	Info LogInfo
}

// Println for PlainLogger logs the provided arguments as a line in plain text.
// The arguments are handled in the manner of fmt.Println.
func (l *PlainLogger) Println(v ...interface{}) {
	log.SetOutput(l.Info.File)
	log.Println(v...)
}

// Printf for PlainLogger logs the provided formatted string in plain text.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Printf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	log.Printf(format, v...)
}

// Error for PlainLogger logs the provided arguments as an error line
// in plain text.
// The arguments are handled in the manner of fmt.Println.
func (l *PlainLogger) Error(v ...interface{}) {
	log.SetOutput(l.Info.File)
	log.Println(v...)
}

// Errorf for PlainLogger logs the provided formatted string as an error
// line in plain text.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Errorf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	log.Printf(format, v...)
}

// LogInfo represents parameters used to manage logging throughout
// a program.
//
// **Attributes:**
//
// Dir: A string representing the directory where the log file is located.
// File: An afero.File object representing the log file.
// FileName: A string representing the name of the log file.
// Path: A string representing the full path to the log file.
type LogInfo struct {
	Dir      string
	File     afero.File
	FileName string
	Path     string
}

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
func CreateLogFile(fs afero.Fs, logDir string, logName string) (LogInfo, error) {
	logInfo := LogInfo{}
	var err error

	logDir = strings.TrimSpace(logDir)
	logName = strings.TrimSpace(logName)

	if logDir == "" {
		return logInfo, fmt.Errorf("logDir cannot be empty")
	}

	if logName == "" {
		return logInfo, fmt.Errorf("logName cannot be empty")
	}

	logInfo.Dir = filepath.Join(logDir, "logs")

	if filepath.Ext(logName) != ".log" {
		logInfo.FileName = fmt.Sprintf("%s.log", logName)
	} else {
		logInfo.FileName = logName
	}

	logInfo.Path = filepath.Join(logInfo.Dir, logInfo.FileName)

	// Create path to log file if the log file doesn't already exist.
	if _, err := fs.Stat(logInfo.Path); os.IsNotExist(err) {
		if err := fs.MkdirAll(logInfo.Dir, os.ModePerm); err != nil {
			return logInfo, fmt.Errorf("failed to create %s: %v", logInfo.Dir, err)
		}
	}

	// Create log file.
	logInfo.File, err = fs.OpenFile(logInfo.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return logInfo, fmt.Errorf("failed to create %s: %v", logInfo.Path, err)
	}

	return logInfo, nil
}

// SlogLogger represents a logger using the slog library.
//
// **Attributes:**
//
// Logger: Logger object from slog library.
type SlogLogger struct {
	Logger *slog.Logger
}

// Println for SlogLogger logs the provided arguments as a line using
// slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *SlogLogger) Println(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

// Printf for SlogLogger logs the provided formatted string using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *SlogLogger) Printf(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Error for SlogLogger logs the provided arguments as an error line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *SlogLogger) Error(v ...interface{}) {
	l.Logger.Error(fmt.Sprint(v...))
}

// Errorf for SlogLogger logs the provided formatted string as an error
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *SlogLogger) Errorf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

// SlogPlainLogger represents a plain logger using the slog library.
//
// **Attributes:**
//
// Logger: Logger object from slog library.
type SlogPlainLogger struct {
	Logger *slog.Logger
}

// Println for SlogPlainLogger logs the provided arguments as a line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *SlogPlainLogger) Println(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

// Printf for SlogPlainLogger logs the provided formatted string
// using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *SlogPlainLogger) Printf(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Error for SlogPlainLogger logs the provided arguments as an error line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *SlogPlainLogger) Error(v ...interface{}) {
	l.Logger.Error(fmt.Sprint(v...))
}

// Errorf for SlogPlainLogger logs the provided formatted string as an
// error line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *SlogPlainLogger) Errorf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

// ConfigureLogger creates a logger based on the provided level.
// Depending on the level, it returns a colored or plain logger.
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
	var err error

	// Create log file handlers
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Setup options for logging
	opts := slog.HandlerOptions{
		Level: level,
	}

	// File logger
	fileHandler := slog.NewJSONHandler(logFile, &opts)

	// Stdout logger
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &opts)

	// Combining both handlers
	handler := slogmulti.Fanout(fileHandler, stdoutHandler)

	logger := slog.New(handler)

	// Depending on the level, return colored or plain logger
	if level == slog.LevelDebug {
		return &SlogLogger{Logger: logger}, nil
	}
	return &SlogPlainLogger{Logger: logger}, nil
}
