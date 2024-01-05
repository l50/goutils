package logging

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fatih/color"
)

// ColorLogger is a logger that outputs messages in a specified color.
// It enhances readability by color-coding log messages based on their
// severity or purpose.
//
// **Attributes:**
//
// Info: LogConfig object containing information about the log file.
// ColorAttribute: A color attribute for output styling.
// Logger: The slog Logger instance used for logging operations.
type ColorLogger struct {
	Cfg            LogConfig
	ColorAttribute color.Attribute
	Logger         *slog.Logger
}

// NewColorLogger creates a new ColorLogger instance with the specified
// LogConfig, color attribute, and slog.Logger.
//
// **Parameters:**
//
// cfg: LogConfig object containing information about the log file.
// colorAttr: A color attribute for output styling.
// logger: The slog Logger instance used for logging operations.
//
// **Returns:**
//
// *ColorLogger: A new instance of ColorLogger.
// error: An error if any issue occurs during initialization.
func NewColorLogger(cfg LogConfig, colorAttr color.Attribute, logger *slog.Logger) (*ColorLogger, error) {
	return &ColorLogger{
		Cfg:            cfg,
		ColorAttribute: colorAttr,
		Logger:         logger,
	}, nil
}

func (l *ColorLogger) assembleLog(level slog.Level, v ...interface{}) string {
	msg := fmt.Sprint(v...) // Convert slice to string
	// Removed PrettyHandler related code, now just returning the message
	return msg
}

// Println for ColorLogger logs the provided arguments as a line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColorLogger) Println(v ...interface{}) {
	msg := fmt.Sprint(v...) // Convert slice to string

	// Create a new record with attributes
	record := slog.Record{
		Level:   slog.LevelInfo,
		Message: msg,
	}

	l.Logger.Log(context.Background(), record.Level, record.Message)
}

// Printf for ColorLogger logs the provided formatted string in
// the specified color. The format and arguments are handled in the
// manner of fmt.Printf.
func (l *ColorLogger) Printf(format string, v ...interface{}) {
	formattedMessage := fmt.Sprintf(format, v...)
	coloredOutput := l.assembleLog(slog.LevelInfo, formattedMessage)
	l.Logger.Log(context.Background(), slog.LevelInfo, coloredOutput)
}

// Error for ColorLogger logs the provided arguments as an error line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColorLogger) Error(v ...interface{}) {
	coloredOutput := l.assembleLog(slog.LevelError, v...)
	l.Logger.Log(context.Background(), slog.LevelError, coloredOutput)
}

// Errorf for ColorLogger logs the provided formatted string as an
// error line in the specified color. The format and arguments are handled
// in the manner of fmt.Printf.
func (l *ColorLogger) Errorf(format string, v ...interface{}) {
	formattedMessage := fmt.Sprintf(format, v...)
	coloredOutput := l.assembleLog(slog.LevelError, formattedMessage)
	l.Logger.Log(context.Background(), slog.LevelError, coloredOutput)
}

// Debug for ColorLogger logs the provided arguments as a debug line
// in the specified color. The arguments are handled in the manner
// of fmt.Println.
func (l *ColorLogger) Debug(v ...interface{}) {
	coloredOutput := l.assembleLog(slog.LevelDebug, v...)
	l.Logger.Log(context.Background(), slog.LevelDebug, coloredOutput)
}

// Debugf for ColorLogger logs the provided formatted string as a debug
// line in the specified color. The format and arguments are handled
// in the manner of fmt.Printf.
func (l *ColorLogger) Debugf(format string, v ...interface{}) {
	formattedMessage := fmt.Sprintf(format, v...)
	coloredOutput := l.assembleLog(slog.LevelDebug, formattedMessage)
	l.Logger.Log(context.Background(), slog.LevelDebug, coloredOutput)
}
