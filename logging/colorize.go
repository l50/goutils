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

func colorizeLevel(level slog.Level, levelStr string) string {
	var colorCode int
	switch level {
	case slog.LevelDebug:
		colorCode = magenta
	case slog.LevelInfo:
		colorCode = green
	case slog.LevelWarn:
		colorCode = yellow
	case slog.LevelError:
		colorCode = red
	default:
		colorCode = white
	}
	return colorize(colorCode, levelStr+":")

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
