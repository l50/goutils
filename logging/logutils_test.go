package logging_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"log/slog"

	"github.com/fatih/color"
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

func TestLoggerImplementation(t *testing.T) {
	// create an in-memory filesystem
	fs := afero.NewMemMapFs()
	logDir := "/tmp"
	logName := "test.log"
	logInfo, err := logging.CreateLogFile(fs, logDir, logName)
	if err != nil {
		t.Fatalf("unexpected error while creating log file: %v", err)
	}

	tests := []struct {
		name    string
		logger  logging.Logger
		logFunc func(l logging.Logger)
		errFunc func(l logging.Logger)
	}{
		{
			name:   "PlainLogger Println",
			logger: &logging.PlainLogger{Info: logInfo},
			logFunc: func(l logging.Logger) {
				l.Println("test plain logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("test plain logger error")
			},
		},
		{
			name:   "PlainLogger Printf",
			logger: &logging.PlainLogger{Info: logInfo},
			logFunc: func(l logging.Logger) {
				l.Printf("test %s logger", "plain")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("test %s logger error", "plain")
			},
		},
		{
			name:   "Test ColoredLogger",
			logger: &logging.ColoredLogger{Info: logInfo},
			logFunc: func(l logging.Logger) {
				l.Printf("Test log message with format: %s", "Test")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("test error message with format: %s", "Test")
			},
		},
		{
			name:   "Test ColoredLogger with Blue color",
			logger: &logging.ColoredLogger{Info: logInfo, ColorAttribute: color.FgBlue},
			logFunc: func(l logging.Logger) {
				l.Printf("Test log message with format: %s", "Test")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("test error message with format: %s", "Test")
			},
		},
		{
			name:   "Test PlainLogger",
			logger: &logging.PlainLogger{Info: logInfo},
			logFunc: func(l logging.Logger) {
				l.Printf("Test log message with format: %s", "Test")
			},
		},
		{
			name:   "Test SlogLogger",
			logger: &logging.SlogLogger{Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))},
			logFunc: func(l logging.Logger) {
				l.Printf("Test log message with format: %s", "Test")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("Test error message with format: %s", "Test")
			},
		},
		{
			name:   "Test SlogPlainLogger",
			logger: &logging.SlogPlainLogger{Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))},
			logFunc: func(l logging.Logger) {
				l.Printf("Test log message with format: %s", "Test")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("Test error message with format: %s", "Test")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.logFunc(tc.logger)
			if tc.errFunc != nil {
				tc.errFunc(tc.logger)
			}
		})
	}
}

func TestConfigureLogger(t *testing.T) {
	tests := []struct {
		name       string
		level      slog.Level
		loggerType string
	}{
		{
			name:       "Test debug level",
			level:      slog.LevelDebug,
			loggerType: "*logging.SlogLogger",
		},
		{
			name:       "Test non-debug level",
			level:      slog.LevelInfo,
			loggerType: "*logging.SlogPlainLogger",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := logging.ConfigureLogger(tc.level, "/tmp/test.log")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Compare(fmt.Sprintf("%T", logger), tc.loggerType) != 0 {
				t.Fatalf("expected logger type: %s, got: %T", tc.loggerType, logger)
			}
		})
	}
}
