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
)

func TestPrettyHandlerHandle(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		message  string
		expected string
	}{
		{
			name:     "Info Level Log",
			level:    slog.LevelInfo,
			message:  "info level test message",
			expected: "info level test message",
		},
		{
			name:     "Debug Level Log",
			level:    slog.LevelDebug,
			message:  "debug level test message",
			expected: "debug level test message",
		},
		{
			name:     "Error Level Log",
			level:    slog.LevelError,
			message:  "error level test message",
			expected: "error level test message",
		},
		{
			name:     "Warn Level Log",
			level:    slog.LevelWarn,
			message:  "warn level test message",
			expected: "warn level test message",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer to capture the output
			var buf strings.Builder
			prettyHandler := logging.NewPrettyHandler(&buf, logging.PrettyHandlerOptions{})

			// Create a log record
			record := slog.Record{
				Level:   tc.level,
				Time:    time.Now(),
				Message: tc.message,
			}

			// Call the handle method
			err := prettyHandler.Handle(context.Background(), record)
			if err != nil {
				t.Fatalf("Handle() error = %v", err)
			}

			// Check if the output contains the expected message
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
	tests := []struct {
		name    string
		level   slog.Level
		message string
	}{
		{
			name:    "Marshal error",
			level:   slog.LevelInfo,
			message: "test message",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf strings.Builder
			prettyHandler := logging.NewPrettyHandler(&buf, logging.PrettyHandlerOptions{})

			// Create a log record using NewRecord and add the failing attribute.
			record := slog.NewRecord(time.Now(), tc.level, tc.message, 0)
			record.AddAttrs(slog.Any("failingAttr", failingJSONMarshal{}))

			// Call the handle method and expect an error.
			err := prettyHandler.Handle(context.Background(), record)
			if err == nil {
				t.Errorf("Expected an error, but got none")
			}
		})
	}
}
