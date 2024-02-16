package logging

import (
	"context"
	"fmt"
	"log/slog"
)

// PlainLogger is a logger implementation using the slog library. It
// provides structured logging capabilities.
//
// **Attributes:**
//
// Info: LogConfig object containing information about the log file.
// Logger: The slog Logger instance used for logging operations.
type PlainLogger struct {
	Info   LogConfig
	Logger *slog.Logger
}

// NewPlainLogger creates a new PlainLogger instance with the specified
// LogConfig and slog.Logger.
//
// **Parameters:**
//
// cfg: LogConfig object containing information about the log file.
// logger: The slog Logger instance used for logging operations.
//
// **Returns:**
//
// *PlainLogger: A new instance of PlainLogger.
// error: An error if any issue occurs during initialization.
func NewPlainLogger(cfg LogConfig, logger *slog.Logger) (*PlainLogger, error) {
	return &PlainLogger{
		Info:   cfg,
		Logger: logger,
	}, nil
}

// Println for PlainLogger logs the provided arguments as a line using
// slog library.
// The arguments are converted to a string using fmt.Sprint.
// PlainLogger.go
func (l *PlainLogger) Println(v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelInfo, fmt.Sprintln(v...))
}

// Printf for PlainLogger logs the provided formatted string using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Printf(format string, v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, v...))
}

// Error for PlainLogger logs the provided arguments as an error line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Error(v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelError, fmt.Sprintln(v...))
}

// Errorf for PlainLogger logs the provided formatted string as an error
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Errorf(format string, v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelError, fmt.Sprintf(format, v...))
}

// Debug for PlainLogger logs the provided arguments as a debug line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Debug(v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelDebug, fmt.Sprintln(v...))
}

// Debugf for PlainLogger logs the provided formatted string as a debug
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Debugf(format string, v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, v...))
}

// Warn for PlainLogger logs the provided arguments as a warning line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Warn(v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelWarn, fmt.Sprintln(v...))
}

// Warnf for PlainLogger logs the provided formatted string as a warning
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Warnf(format string, v ...interface{}) {
	l.Logger.Log(context.Background(), slog.LevelWarn, fmt.Sprintf(format, v...))
}
