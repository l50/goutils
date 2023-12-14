package logging

// GlobalLogger is a global variable that holds the instance of the logger.
var GlobalLogger Logger

// L returns the global logger instance for use in logging operations.
//
// **Returns:**
//
// Logger: The global Logger instance.
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
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}
