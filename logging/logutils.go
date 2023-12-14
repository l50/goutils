package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"log/slog"

	"github.com/fatih/color"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/afero"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

// colorize adds ANSI color codes to the given string.
func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%dm%s%s", colorCode, v, reset)
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

// PlainLogger is a logger implementation using the slog library. It
// provides structured logging capabilities.
//
// **Attributes:**
//
// Info: LogInfo object containing information about the log file.
// Logger: The slog Logger instance used for logging operations.
type PlainLogger struct {
	Info   LogInfo
	Logger *slog.Logger
}

// ColorLogger is a logger that outputs messages in a specified color.
// It enhances readability by color-coding log messages based on their
// severity or purpose.
//
// **Attributes:**
//
// Info: LogInfo object containing information about the log file.
// ColorAttribute: A color attribute for output styling.
// Logger: The slog Logger instance used for logging operations.
type ColorLogger struct {
	Info           LogInfo
	ColorAttribute color.Attribute
	Logger         *slog.Logger
}

// GlobalLogger is a variable that holds the instance of the logger.
var GlobalLogger Logger

// PrettyHandlerOptions - Options for PrettyHandler
type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

// PrettyHandler - Custom handler for color logging
type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

// InitGlobalLogger initializes the global logger with the specified level and file path.
func InitGlobalLogger(level slog.Level, path string) error {
	var err error
	GlobalLogger, err = ConfigureLogger(level, path)
	return err
}

// L returns the global logger instance.
func L() Logger {
	return GlobalLogger
}

// Logger is an interface that defines methods for a generic logging
// system. It supports basic logging operations like printing,
// formatted printing, error logging, and debug logging.
//
// **Methods:**
//
// Println: Outputs a line with the given arguments.
// Printf: Outputs a formatted string.
// Error: Logs an error message.
// Errorf: Logs a formatted error message.
// Debug: Logs a debug message.
// Debugf: Logs a formatted debug message.
// Info: Logs a structured message.
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

// Println for ColorLogger logs the provided arguments as a line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColorLogger) Println(v ...interface{}) {
	// Format message with color and then log using slog
	coloredOutput := color.New(l.ColorAttribute).SprintFunc()
	l.Logger.Info(coloredOutput(fmt.Sprintln(v...)))
}

// Printf for ColorLogger logs the provided formatted string in
// the specified color. The format and arguments are handled in the
// manner of fmt.Printf.
func (l *ColorLogger) Printf(format string, v ...interface{}) {
	// Format message with color and then log using slog
	coloredOutput := color.New(l.ColorAttribute).SprintfFunc()
	l.Logger.Info(coloredOutput(format, v...))
}

// Error for ColorLogger logs the provided arguments as an error line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColorLogger) Error(v ...interface{}) {
	coloredOutput := color.New(color.FgRed).SprintFunc()
	l.Logger.Error(coloredOutput(fmt.Sprintln(v...)))
}

// Errorf for ColorLogger logs the provided formatted string as an
// error line in the specified color. The format and arguments are handled
// in the manner of fmt.Printf.
func (l *ColorLogger) Errorf(format string, v ...interface{}) {
	coloredOutput := color.New(color.FgRed).SprintfFunc()
	l.Logger.Error(coloredOutput(format, v...))
}

// Debug for ColorLogger logs the provided arguments as a debug line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColorLogger) Debug(v ...interface{}) {
	log.SetOutput(l.Info.File)
	color.New(l.ColorAttribute).Sprint(v...)
}

// Debugf for ColorLogger logs the provided formatted string as a debug
// line in the specified color. The format and arguments are handled
// in the manner of fmt.Printf.
func (l *ColorLogger) Debugf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	color.New(l.ColorAttribute).Printf(format, v...)
}

// Println for PlainLogger logs the provided arguments as a line using
// slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Println(v ...interface{}) {
	log.SetOutput(l.Info.File)
	l.Logger.Info(fmt.Sprint(v...))
}

// Printf for PlainLogger logs the provided formatted string using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Printf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Error for PlainLogger logs the provided arguments as an error line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Error(v ...interface{}) {
	log.SetOutput(l.Info.File)
	l.Logger.Error(fmt.Sprint(v...))
}

// Errorf for PlainLogger logs the provided formatted string as an error
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Errorf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	l.Logger.Error(fmt.Sprintf(format, v...))
}

// Debug for PlainLogger logs the provided arguments as a debug line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Debug(v ...interface{}) {
	log.SetOutput(l.Info.File)
	l.Logger.Debug(fmt.Sprint(v...))
}

// Debugf for PlainLogger logs the provided formatted string as a debug
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Debugf(format string, v ...interface{}) {
	log.SetOutput(l.Info.File)
	l.Logger.Debug(fmt.Sprintf(format, v...))
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

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelColor := colorizeLevel(r.Level)
	timeStr := r.Time.Format("[15:05:05.000]")
	messageColor := colorizeMessage(r.Message, r.Level)
	fields, err := json.MarshalIndent(extractFields(r), "", "  ")
	if err != nil {
		return err
	}

	h.l.Println(timeStr, levelColor, messageColor, string(fields))
	return nil
}

func colorizeLevel(level slog.Level) string {
	var colorCode int
	switch level {
	case slog.LevelDebug:
		colorCode = darkGray
	case slog.LevelInfo:
		colorCode = cyan
	case slog.LevelWarn:
		colorCode = lightYellow
	case slog.LevelError:
		colorCode = lightRed
	default:
		colorCode = white
	}
	return colorize(colorCode, level.String()+":")

}

func colorizeMessage(message string, level slog.Level) string {
	if level == slog.LevelError {
		return color.New(color.FgRed).Sprint(message)
	}
	return message
}

func extractFields(r slog.Record) map[string]interface{} {
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	return fields
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
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

func determineColorAttribute(level slog.Level) color.Attribute {
	switch level {
	case slog.LevelDebug:
		return color.FgMagenta
	case slog.LevelInfo:
		return color.FgBlue
	case slog.LevelWarn:
		return color.FgYellow
	case slog.LevelError:
		return color.FgRed
	default:
		return color.Reset
	}
}
