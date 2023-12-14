package logging

import "log/slog"

// GlobalLogger is a variable that holds the instance of the logger.
var GlobalLogger Logger

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
