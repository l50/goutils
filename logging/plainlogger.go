package logging

import (
	"fmt"
	"log"
	"log/slog"
)

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

// Close for PlainLogger closes the log file.
func (l *PlainLogger) Close() error {
	if l.Info.File != nil {
		return l.Info.File.Close()
	}
	return nil
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
