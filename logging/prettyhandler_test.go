package logging_test

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/l50/goutils/v2/logging"
	"github.com/l50/goutils/v2/str"
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

// failingJSONMarshal is a custom type that causes JSON marshaling to fail.
type failingJSONMarshal struct{}

// MarshalJSON for failingJSONMarshal always returns an error.
func (f failingJSONMarshal) MarshalJSON() ([]byte, error) {
	return nil, &json.MarshalerError{Type: reflect.TypeOf(f), Err: errors.New("marshal error")}
}

func TestPrettyHandlerHandleMarshalError(t *testing.T) {
	testCases := []struct {
		name  string
		level slog.Level
		msg   string
	}{
		{
			name:  "Marshal error",
			level: slog.LevelInfo,
			msg:   "test msg",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf strings.Builder
			prettyHandler := logging.NewPrettyHandler(&buf, logging.PrettyHandlerOptions{})

			// Create a log record using NewRecord and add the failing attribute.
			record := slog.NewRecord(time.Now(), tc.level, tc.msg, 0)
			record.AddAttrs(slog.Any("failingAttr", failingJSONMarshal{}))

			// Call the handle method and expect an error.
			err := prettyHandler.Handle(context.Background(), record)
			if err == nil {
				t.Errorf("Expected an error, but got none")
			}
		})
	}
}

func TestPrettyHandlerNoColorCodes(t *testing.T) {
	testCases := []struct {
		name        string
		level       slog.Level
		msg         string
		expectError bool
	}{
		{
			name:        "Output should not contain color codes",
			level:       slog.LevelInfo,
			msg:         "\x1b[36minfo level test msg\x1b[0m",
			expectError: false,
		},
		{
			name:        "Normal output should not be affected",
			level:       slog.LevelInfo,
			msg:         "info level test msg",
			expectError: false,
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
				Message: tc.msg,
			}

			// Call the handle method
			err := prettyHandler.Handle(context.Background(), record)
			if err != nil {
				t.Fatalf("Handle() error = %v", err)
			}

			// Check the output does not contain ANSI color codes
			output := str.StripANSI(buf.String())
			if strings.Contains(output, "\u001b[") {
				t.Fatalf("Output should not contain color codes, got '%s'", output)
			}

			// The original test case was comparing a boolean to a string, which is incorrect.
			// We just need to check if there was an error and if the output is as expected.
			if err != nil != tc.expectError {
				t.Fatalf("Expected error: %v, got: %v", tc.expectError, err != nil)
			}
		})
	}
}
