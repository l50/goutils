package logging_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"log/slog"

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

func TestPlainLogger(t *testing.T) {
	tests := []struct {
		name    string
		level   slog.Level
		logFunc func(l logging.Logger)
		errFunc func(l logging.Logger)
	}{
		{
			name:  "Test PlainLogger Println",
			level: slog.LevelInfo,
			logFunc: func(l logging.Logger) {
				l.Println("Test Plain Println logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Test Plain Println logger error")
			},
		},
		{
			name:  "Test PlainLogger Printf",
			level: slog.LevelInfo,
			logFunc: func(l logging.Logger) {
				l.Printf("Test %s logger", "Plain Printf")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("Test %s logger error", "Plain Printf")
			},
		},
		{
			name:  "Test PlainLogger Debug",
			level: slog.LevelDebug,
			logFunc: func(l logging.Logger) {
				l.Debug("Test Plain Debug logger")
			},
			errFunc: func(l logging.Logger) {
				l.Debug("Test Plain Debug logger with error")
			},
		},
		{
			name:  "Test PlainLogger Debugf",
			level: slog.LevelDebug,
			logFunc: func(l logging.Logger) {
				l.Debugf("Test %s logger", "Plain Debugf")
			},
			errFunc: func(l logging.Logger) {
				l.Debugf("Test %s logger with error", "Plain Debugf")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := logging.ConfigureLogger(tc.level, "/tmp/test.log")
			if err != nil {
				t.Fatalf("Failed to configure logger: %v", err)
			}
			t.Logf("Running test case: %s", tc.name)
			tc.logFunc(logger)
			if tc.errFunc != nil {
				tc.errFunc(logger)
			}
		})
	}
}

// func TestLoggerImplementation(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		level   slog.Level
// 		logger  logging.Logger
// 		logFunc func(l logging.Logger)
// 		errFunc func(l logging.Logger)
// 	}{
// 		{
// 			name:  "Test PlainLogger Println",
// 			level: slog.LevelInfo,
// 			logger: func() logging.Logger {
// 				logger, _ := logging.ConfigureLogger(slog.LevelInfo, "/tmp/test.log")
// 				return logger
// 			}(),
// 			logFunc: func(l logging.Logger) {
// 				l.Println("test plain logger")
// 			},
// 			errFunc: func(l logging.Logger) {
// 				l.Error("test plain logger error")
// 			},
// 		},
// 		{
// 			name:  "Test PlainLogger Printf",
// 			level: slog.LevelInfo,
// 			logger: func() logging.Logger {
// 				logger, err := logging.ConfigureLogger(slog.LevelInfo, "/tmp/test.log")
// 				if err != nil {
// 					t.Fatalf("unexpected error: %v", err)
// 				}
// 				return logger
// 			}(),
// 			logFunc: func(l logging.Logger) {
// 				l.Printf("Test %s logger", "plain printf")
// 			},
// 			errFunc: func(l logging.Logger) {
// 				l.Errorf("Test %s logger error", "plain printf")
// 			},
// 		},
// 		{
// 			name:  "Test PlainLogger Debug",
// 			level: slog.LevelDebug,
// 			logger: func() logging.Logger {
// 				logger, err := logging.ConfigureLogger(slog.LevelDebug, "/tmp/test.log")
// 				if err != nil {
// 					t.Fatalf("unexpected error: %v", err)
// 				}
// 				return logger
// 			}(),
// 			logFunc: func(l logging.Logger) {
// 				l.Debug("Test debug message")
// 			},
// 			errFunc: func(l logging.Logger) {
// 				l.Debugf("Test debug message with format: %s", "Test")
// 			},
// 		},
// {
// 	name:  "TestColorLogger Println",
// 	level: slog.LevelInfo,
// 	logger: func() logging.ColorLogger {
// 		logger, _ := logging.ConfigureLogger(slog.LevelDebug, "/tmp/test.log")
// 		return logger
// 	}(),
// 	logFunc: func(l logging.ColorLogger) {
// 		l.Printf("Test log message with format: %s %s", "Success", "magenta")
// 	},
// 	errFunc: func(l logging.Logger) {
// 		l.Errorf("Test error message with format: %s %s", "Failure", "red")
// 	},
// },
// {
// 	name:  "TestColorLogger Printf",
// 	level: slog.LevelInfo,
// 	logger: func() logging.Logger {
// 		logger, _ := logging.ConfigureLogger(slog.LevelDebug, "/tmp/test.log")
// 		return logger
// 	}(),
// 	logFunc: func(l logging.Logger) {
// 		l.Printf("Test log message with format: %s %s", "Success", "magenta")
// 	},
// 	errFunc: func(l logging.Logger) {
// 		l.Errorf("Test error message with format: %s %s", "Failure", "red")
// 	},
// },
// {
// 	name:  "TestColorLogger Debugf",
// 	level: slog.LevelDebug,
// 	logger: func() logging.Logger {
// 		logger, _ := logging.ConfigureLogger(slog.LevelDebug, "/tmp/test.log")
// 		return logger
// 	}(),
// 	logFunc: func(l logging.Logger) {
// 		l.Debugf("Test debug message with format: %s", "Test", "green")
// 	},
// 	errFunc: func(l logging.Logger) {
// 		l.Errorf("Test error message with format: %s", "Test", "red")
// 	},
// },
// {
// 	name:  "Test JSON Logging Format with color",
// 	level: slog.LevelInfo,
// 	logger: func() logging.Logger {
// 		logger, _ := logging.ConfigureLogger(slog.LevelDebug, "/tmp/test.log")
// 		return logger
// 	}(),
// 	logFunc: func(l logging.Logger) {
// 		l.Printf("Test log message with format: %s", "Test", "green")
// 	},
// 	errFunc: func(l logging.Logger) {
// 		l.Errorf("Test error message with format: %s", "Test", "red")
// 	},
// },
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			logger, err := logging.ConfigureLogger(tc.level, "/tmp/test.log")
// 			if err != nil {
// 				t.Fatalf("Failed to configure logger: %v", err)
// 			}
// 			t.Logf("Running test case: %s", tc.name)
// 			tc.logFunc(logger)
// 			if tc.errFunc != nil {
// 				tc.errFunc(logger)
// 			}
// 		})
// 	}
// }

// func TestConfigureLogger(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		level      slog.Level
// 		loggerType string
// 	}{
// 		{
// 			name:       "Test debug level",
// 			level:      slog.LevelDebug,
// 			loggerType: "*logging.ColorLogger",
// 		},
// 		{
// 			name:       "Test non-debug level",
// 			level:      slog.LevelInfo,
// 			loggerType: "*logging.ColorLogger",
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			logger, err := logging.ConfigureLogger(tc.level, "/tmp/test.log")
// 			if err != nil {
// 				t.Fatalf("unexpected error: %v", err)
// 			}
// 			if strings.Compare(fmt.Sprintf("%T", logger), tc.loggerType) != 0 {
// 				t.Fatalf("expected logger type: %s, got: %T", tc.loggerType, logger)
// 			}
// 		})
// 	}
// }
