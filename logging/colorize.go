package logging

import (
	"fmt"
	"log/slog"

	"github.com/fatih/color"
)

// ANSI color codes
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
