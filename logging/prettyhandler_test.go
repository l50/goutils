package logging_test

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/l50/goutils/v2/logging"
)

func TestPrettyHandlerHandle(t *testing.T) {
	testCases := []struct {
		name     string
		level    slog.Level
		msg      string
		expected string
	}{
		{
			name:     "Info Level Log",
			level:    slog.LevelInfo,
			msg:      "info level test msg",
			expected: "info level test msg",
		},
		{
			name:     "Debug Level Log",
			level:    slog.LevelDebug,
			msg:      "debug level test msg",
			expected: "debug level test msg",
		},
		{
			name:     "Error Level Log",
			level:    slog.LevelError,
			msg:      "error level test msg",
			expected: "error level test msg",
		},
		{
			name:     "Warn Level Log",
			level:    slog.LevelWarn,
			msg:      "warn level test msg",
			expected: "warn level test msg",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer to capture the output
			var buf strings.Builder
			prettyHandler := logging.NewPrettyHandler(&buf, logging.PrettyHandlerOptions{})

			// Create a log record
			record := slog.Record{
				Level:   tc.level,
				Time:    time.Now(),
				Message: tc.msg,
			}

			// Call the handle method
			err := prettyHandler.Handle(context.Background(), record)
			if err != nil {
				t.Fatalf("Handle() error = %v", err)
			}

			// Check if the output contains the expected msg
			if !strings.Contains(buf.String(), tc.expected) {
				t.Fatalf("Expected to find '%s' in the output, got '%s'", tc.expected, buf.String())
			}
		})
	}
}

func TestPrettyHandlerParseLogRecord(t *testing.T) {
	testCases := []struct {
		name        string
		record      slog.Record
		expectError bool
	}{
		{
			name: "Valid JSON Log Record",
			record: slog.Record{
				Level:   slog.LevelInfo,
				Message: `{"time":"2024-01-01T12:00:00Z","level":"INFO","msg":"JSON test message"}`,
			},
			expectError: false,
		},
		{
			name: "Invalid JSON Log Record",
			record: slog.Record{
				Level:   slog.LevelInfo,
				Message: `{"time":"2024-01-01T12:00:00Z",level:"INFO","msg":"JSON test message"}`,
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prettyHandler := logging.NewPrettyHandler(io.Discard, logging.PrettyHandlerOptions{})

			// Create a log record
			record := slog.Record{
				Level:   tc.record.Level,
				Time:    time.Now(),
				Message: tc.record.Message,
			}

			// Call the handle method
			err := prettyHandler.Handle(context.Background(), record)
			if (err != nil) != tc.expectError {
				t.Errorf("parseLogRecord() for %s expected error: %v, got: %v", tc.name, tc.expectError, err)
			}
		})
	}
}

func TestPrettyHandlerColorization(t *testing.T) {
	testCases := []struct {
		name  string
		level slog.Level
	}{
		{
			name:  "Info Level Color",
			level: slog.LevelInfo,
		},
		{
			name:  "Error Level Color",
			level: slog.LevelError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf strings.Builder
			prettyHandler := logging.NewPrettyHandler(&buf, logging.PrettyHandlerOptions{})

			// Create a log record
			record := slog.Record{
				Level:   tc.level,
				Time:    time.Now(),
				Message: "test message",
			}

			// Call the handle method
			err := prettyHandler.Handle(context.Background(), record)
			if err != nil {
				t.Errorf("Handle() error = %v", err)
			}

			// Ensure the output does not contain ANSI color codes
			output := buf.String()
			if strings.Contains(output, "\u001b[") {
				t.Errorf("Output should not contain color codes, got '%s'", output)
			}
		})
	}
}
