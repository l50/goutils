package logging_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/l50/goutils/logging"
)

func TestCreateLogFile(t *testing.T) {
	tests := []struct {
		name        string
		logDir      string
		logName     string
		expectError bool
	}{
		{
			name:        "Create log file",
			logDir:      "/tmp",
			logName:     "test.log",
			expectError: false,
		},
		{
			name:        "Create log file",
			logDir:      "/tmp",
			logName:     "testing",
			expectError: false,
		},
		{
			name:        "Create log file",
			logDir:      "/tmp",
			logName:     "testing       ",
			expectError: false,
		},
		{
			name:        "Create log file with empty directory",
			logDir:      "",
			logName:     "test.log",
			expectError: true,
		},
		{
			name:        "Create log file with empty filename",
			logDir:      "/tmp",
			logName:     "",
			expectError: true,
		},
		{
			name:        "Ensure handling of bad input works",
			logDir:      "/tmp/bla/bla",
			logName:     "/tmp/stuff/things/bla/test.log",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logInfo, err := logging.CreateLogFile(tc.logDir, tc.logName)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("got unexpected error: %v", err)
				}
				trimmedLogName := strings.TrimSpace(tc.logName)
				if filepath.Ext(trimmedLogName) != ".log" {
					trimmedLogName = fmt.Sprintf("%s.log", trimmedLogName)
				}

				expectedPath := filepath.Join(tc.logDir, "logs", trimmedLogName)
				if logInfo.Path != expectedPath {
					t.Fatalf("expected path %s but got %s", expectedPath, logInfo.Path)
				}

				// cleanup after test
				_ = os.Remove(logInfo.Path)
				_ = os.Remove(logInfo.Dir)
			}
		})
	}
}
