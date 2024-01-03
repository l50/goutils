package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"
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

// Println for PlainLogger logs the provided arguments as a line using
// slog library.
// The arguments are converted to a string using fmt.Sprint.
// PlainLogger.go
func (l *PlainLogger) Println(v ...interface{}) {
	file, err := l.Info.Fs.OpenFile(l.Info.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Logger.Error(fmt.Sprint(err))
		return
	}
	defer file.Close()

	log.SetOutput(file)
	l.Logger.Info(fmt.Sprint(v...))
}

// Printf for PlainLogger logs the provided formatted string using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Printf(format string, v ...interface{}) {
	file, err := l.Info.Fs.OpenFile(l.Info.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Logger.Error(fmt.Sprint(err))
		return
	}
	defer file.Close()

	log.SetOutput(file)
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Error for PlainLogger logs the provided arguments as an error line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Error(v ...interface{}) {
	file, err := l.Info.Fs.OpenFile(l.Info.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Logger.Error(fmt.Sprint(err))
		return
	}
	defer file.Close()

	log.SetOutput(file)
	l.Logger.Error(fmt.Sprint(v...))
}

// Errorf for PlainLogger logs the provided formatted string as an error
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Errorf(format string, v ...interface{}) {
	file, err := l.Info.Fs.OpenFile(l.Info.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Logger.Error(fmt.Sprint(err))
		return
	}
	defer file.Close()

	log.SetOutput(file)
	l.Logger.Error(fmt.Sprintf(format, v...))
}

// // Close for PlainLogger closes the log file.
// func (l *PlainLogger) Close() error {
// 	if l.Info.Fs != nil {
// 		return l.Info.Fs.Close(l.Info.File)
// 	}
// 	return nil
// }

// Debug for PlainLogger logs the provided arguments as a debug line
// using slog library.
// The arguments are converted to a string using fmt.Sprint.
func (l *PlainLogger) Debug(v ...interface{}) {
	file, err := l.Info.Fs.OpenFile(l.Info.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Logger.Error(fmt.Sprint(err))
		return
	}
	defer file.Close()

	log.SetOutput(file)
	l.Logger.Debug(fmt.Sprint(v...))
}

// Debugf for PlainLogger logs the provided formatted string as a debug
// line using slog library.
// The format and arguments are handled in the manner of fmt.Printf.
func (l *PlainLogger) Debugf(format string, v ...interface{}) {
	file, err := l.Info.Fs.OpenFile(l.Info.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Logger.Error(fmt.Sprint(err))
		return
	}
	defer file.Close()

	log.SetOutput(file)
	l.Logger.Debug(fmt.Sprintf(format, v...))
}
