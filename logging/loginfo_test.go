package logging_test

import (
	"log/slog"
	"testing"

	"github.com/l50/goutils/v2/logging"
)

func TestDetermineLogLevel(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{
			name:     "Debug level",
			input:    "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "Info level",
			input:    "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "Warn level",
			input:    "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "Error level",
			input:    "error",
			expected: slog.LevelError,
		},
		{
			name:     "Default level",
			input:    "unknown",
			expected: slog.LevelInfo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := logging.DetermineLogLevel(tc.input)
			if result != tc.expected {
				t.Errorf("DetermineLogLevel(%s) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
