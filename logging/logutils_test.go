package logging_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/logging"
	"github.com/spf13/afero"
)

type errorFs struct {
	afero.Fs
}

func (fs *errorFs) MkdirAll(path string, perm os.FileMode) error {
	return fmt.Errorf("simulated error on MkdirAll")
}

func (fs *errorFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return nil, fmt.Errorf("simulated error on OpenFile")
}

func TestCreateLogFile(t *testing.T) {
	// create an in-memory filesystem
	normalFs := afero.NewMemMapFs()
	errorFs := &errorFs{normalFs}

	tests := []struct {
		name        string
		logDir      string
		logName     string
		fs          afero.Fs
		expectError bool
	}{
		{
			name:        "Create log file",
			logDir:      "/tmp",
			logName:     "test.log",
			fs:          normalFs,
			expectError: false,
		},
		{
			name:        "Create log file",
			logDir:      "/tmp",
			logName:     "testing",
			fs:          normalFs,
			expectError: false,
		},
		{
			name:        "Create log file",
			logDir:      "/tmp",
			logName:     "testing       ",
			fs:          normalFs,
			expectError: false,
		},
		{
			name:        "Create log file with empty directory",
			logDir:      "",
			logName:     "test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Create log file with empty filename",
			logDir:      "/tmp",
			logName:     "",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Ensure handling of bad input works",
			logDir:      "/tmp/bla/bla",
			logName:     "/tmp/stuff/things/bla/test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Create log file with unwritable directory",
			logDir:      "/unwritable_dir",
			logName:     "test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Simulate error when creating directory",
			logDir:      "/tmp",
			logName:     "test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Simulate error when creating file",
			logDir:      "/tmp",
			logName:     "test.log",
			fs:          errorFs,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logInfo, err := logging.CreateLogFile(tc.fs, tc.logDir, tc.logName)
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
			}
		})
	}
}
