package logging_test

import (
	"fmt"
	"os"
	"path/filepath"
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
		logPath     string
		fs          afero.Fs
		expectError bool
	}{
		{
			name:        "Create log file",
			logPath:     "/tmp/logs/test.log",
			fs:          normalFs,
			expectError: false,
		},
		{
			name:        "Create log file",
			logPath:     "/tmp/logs/testing.log",
			fs:          normalFs,
			expectError: false,
		},
		{
			name:        "Create log file",
			logPath:     "/tmp/logs/testing.log",
			fs:          normalFs,
			expectError: false,
		},
		{
			name:        "Create log file with empty directory",
			logPath:     "/tmp/logs/test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Create log file with empty filename",
			logPath:     "/tmp/logs/test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Ensure handling of bad input works",
			logPath:     "/tmp/bla/bla/tmp/stuff/things/bla/test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Create log file with unwritable directory",
			logPath:     "/unwritable_dir/test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Simulate error when creating directory",
			logPath:     "/tmp/test.log",
			fs:          errorFs,
			expectError: true,
		},
		{
			name:        "Simulate error when creating file",
			logPath:     "/tmp/test.log",
			fs:          errorFs,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logInfo, err := logging.CreateLogFile(tc.fs, tc.logPath)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("got unexpected error: %v", err)
				}

				expectedPath := tc.logPath
				if logInfo.Path != expectedPath {
					t.Fatalf("expected path %s but got %s", expectedPath, logInfo.Path)
				}
			}
		})
	}
}

func TestPlainLogger(t *testing.T) {
	tests := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		outputType logging.OutputType
		logFunc    func(l logging.Logger)
		errFunc    func(l logging.Logger)
	}{
		{
			name:       "Test PlainLogger Println",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			outputType: logging.PlainOutput,
			logFunc: func(l logging.Logger) {
				l.Println("Test Plain Println logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Test Plain Println logger error")
			},
		},
		{
			name:       "Test PlainLogger Printf",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			outputType: logging.PlainOutput,
			logFunc: func(l logging.Logger) {
				l.Printf("Test %s logger", "Plain Printf")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("Test %s logger error", "Plain Printf")
			},
		},
		{
			name:       "Test PlainLogger Debug",
			level:      slog.LevelDebug,
			fs:         afero.NewMemMapFs(),
			outputType: logging.PlainOutput,
			logFunc: func(l logging.Logger) {
				l.Debug("Test Plain Debug logger")
			},
			errFunc: func(l logging.Logger) {
				l.Debug("Test Plain Debug logger with error")
			},
		},
		{
			name:       "Test PlainLogger Debugf",
			level:      slog.LevelDebug,
			fs:         afero.NewMemMapFs(),
			outputType: logging.PlainOutput,
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
			// Create necessary directory for the test
			err := tc.fs.MkdirAll(filepath.Dir("/tmp/test.log"), 0755)
			if err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			logger, err := logging.ConfigureLogger(tc.fs, tc.level, "/tmp/test.log", tc.outputType)
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

func TestColorLogger(t *testing.T) {
	tests := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		outputType logging.OutputType
		logFunc    func(l logging.Logger)
		errFunc    func(l logging.Logger)
	}{
		{
			name:       "Test ColorLogger Println",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			outputType: logging.ColorOutput,
			logFunc: func(l logging.Logger) {
				l.Println("Test Color Println logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Test Color Println logger error")
			},
		},
		{
			name:       "Test ColorLogger Printf",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			outputType: logging.ColorOutput,
			logFunc: func(l logging.Logger) {
				l.Printf("Test %s logger", "Color Printf")
			},
			errFunc: func(l logging.Logger) {
				l.Errorf("Test %s logger error", "Color Printf")
			},
		},
		{
			name:       "Test ColorLogger Debug",
			level:      slog.LevelDebug,
			fs:         afero.NewMemMapFs(),
			outputType: logging.ColorOutput,
			logFunc: func(l logging.Logger) {
				l.Debug("Test Color Debug logger")
			},
			errFunc: func(l logging.Logger) {
				l.Debug("Test Color Debug logger with error")
			},
		},
		{
			name:       "Test ColorLogger Debugf",
			level:      slog.LevelDebug,
			fs:         afero.NewMemMapFs(),
			outputType: logging.ColorOutput,
			logFunc: func(l logging.Logger) {
				l.Debugf("Test %s logger", "Color Debugf")
			},
			errFunc: func(l logging.Logger) {
				l.Debugf("Test %s logger with error", "Color Debugf")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create necessary directory for the test
			err := tc.fs.MkdirAll(filepath.Dir("/tmp/test.log"), 0755)
			if err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			logger, err := logging.ConfigureLogger(tc.fs, tc.level, "/tmp/test.log", tc.outputType)
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

func TestConfigureLogger(t *testing.T) {
	tests := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		logPath    string
		outputType logging.OutputType
		logFunc    func(l logging.Logger)
		errFunc    func(l logging.Logger)
		wantErr    bool
	}{
		{
			name:       "Info Level with Color Output",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			logPath:    "/tmp/test_info_color.log",
			outputType: logging.ColorOutput,
			logFunc: func(l logging.Logger) {
				l.Println("Test info color logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Test info color logger error")
			},
			wantErr: false,
		},
		{
			name:       "Debug Level with Plain Output",
			level:      slog.LevelDebug,
			fs:         afero.NewMemMapFs(),
			logPath:    "/tmp/test_debug_plain.log",
			outputType: logging.PlainOutput,
			logFunc: func(l logging.Logger) {
				l.Debug("Test debug plain logger")
			},
			errFunc: func(l logging.Logger) {
				l.Debug("Test debug plain logger error")
			},
			wantErr: false,
		},
		{
			name:       "Invalid Path",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			logPath:    "/invalid_path/test.log",
			outputType: logging.PlainOutput,
			wantErr:    true,
		},
	}

	fs := afero.NewMemMapFs()

	// Create necessary directories for the tests
	requiredDirs := []string{"/tmp"}
	for _, dir := range requiredDirs {
		err := fs.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("failed to create directory: %s", dir)
		}
	}

	for _, tc := range tests {
		tc.fs = fs
		t.Run(tc.name, func(t *testing.T) {
			logger, err := logging.ConfigureLogger(tc.fs, tc.level, tc.logPath, tc.outputType)

			if (err != nil) != tc.wantErr {
				t.Errorf("ConfigureLogger() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if err == nil {
				tc.logFunc(logger)
				if tc.errFunc != nil {
					tc.errFunc(logger)
				}
			}
		})
	}
}

func TestGlobalLogger(t *testing.T) {
	tests := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		logPath    string
		outputType logging.OutputType
		logFunc    func(l logging.Logger)
		errFunc    func(l logging.Logger)
	}{
		{
			name:       "Set and Retrieve Global Logger",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			logPath:    "/tmp/test_global_logger.log",
			outputType: logging.PlainOutput,
			logFunc: func(l logging.Logger) {
				l.Println("Testing global logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Testing global logger error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create necessary directory for the test
			err := tc.fs.MkdirAll(filepath.Dir(tc.logPath), 0755)
			if err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			// Configure and set the global logger
			logger, err := logging.ConfigureLogger(tc.fs, tc.level, tc.logPath, tc.outputType)
			if err != nil {
				t.Fatalf("Failed to configure logger: %v", err)
			}
			logging.GlobalLogger = logger

			// Retrieve and use the global logger
			globalLogger := logging.L()
			if globalLogger == nil {
				t.Fatal("GlobalLogger is nil after being set")
			}

			tc.logFunc(globalLogger)
			if tc.errFunc != nil {
				tc.errFunc(globalLogger)
			}
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	tests := []struct {
		name        string
		level       slog.Level
		fs          afero.Fs
		logPath     string
		outputType  logging.OutputType
		logFunc     func(l logging.Logger)
		errFunc     func(l logging.Logger)
		expectError bool
		expectedLog string
	}{
		{
			name:       "Successful Logger Output",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			logPath:    "/tmp/logs/test_logger_output.log",
			outputType: logging.ColorOutput,
			logFunc: func(l logging.Logger) {
				l.Println("Test info color logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Test info color logger error")
			},
			expectError: false,
			expectedLog: "Test message\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fs.MkdirAll(filepath.Dir(tc.logPath), 0755); err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			// Initialize logging with the memory file system
			logger, err := logging.InitLogging(tc.fs, tc.logPath, tc.level, tc.outputType)
			if (err != nil) != tc.expectError {
				t.Fatalf("InitLogging() error = %v, expectError %v", err, tc.expectError)
			}
			defer logger.Close()

			if !tc.expectError {
				// Perform logging operations
				logger.Println("Test message")

				// Now open the log file for reading
				logFile, err := tc.fs.Open(tc.logPath)
				if err != nil {
					t.Fatalf("Failed to open log file: %v", err)
				}
				defer logFile.Close()

				// buf, err := io.ReadAll(logFile)
				// if err != nil {
				// 	t.Fatalf("Failed to read log file: %v", err)
				// }

				// // Assert on the file's content
				// if string(buf) != tc.expectedLog {
				// 	t.Errorf("Expected %q, got %q", tc.expectedLog, string(buf))
				// }
			}
		})
	}
}
