package logging

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

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

// Close for ColorLogger closes the log file.
func (l *ColorLogger) Close() error {
	if l.Info.File != nil {
		return l.Info.File.Close()
	}
	return nil
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
