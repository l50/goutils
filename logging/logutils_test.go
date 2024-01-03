package logging_test

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/logging"
	log "github.com/l50/goutils/v2/logging"
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

	testCases := []struct {
		name        string
		logConfig   logging.LogConfig
		expectError bool
	}{
		{
			name: "Create log file with existing directory",
			logConfig: logging.LogConfig{
				Fs:      normalFs,
				LogPath: "/tmp/logs/test.log",
			},
			expectError: false,
		},
		{
			name: "Create log file with empty directory",
			logConfig: logging.LogConfig{
				Fs:      errorFs,
				LogPath: "/tmp/logs/test.log",
			},
			expectError: true,
		},
		{
			name: "Create log file with empty filename",
			logConfig: logging.LogConfig{
				Fs:      errorFs,
				LogPath: "/tmp/logs/test.log",
			},
			expectError: true,
		},
		{
			name: "Ensure handling of bad input works",
			logConfig: logging.LogConfig{
				Fs:      errorFs,
				LogPath: "/tmp/bla/bla/tmp/stuff/things/bla/test.log",
			},
			expectError: true,
		},
		{
			name: "Create log file with unwritable directory",
			logConfig: logging.LogConfig{
				Fs:      errorFs,
				LogPath: "/tmp/logs/test.log",
			},
			expectError: true,
		},
		{
			name: "Simulate error when creating directory",
			logConfig: logging.LogConfig{
				Fs:      errorFs,
				LogPath: "/tmp/test.log",
			},
			expectError: true,
		},
		{
			name: "Simulate error when creating file",
			logConfig: logging.LogConfig{
				Fs:      errorFs,
				LogPath: "/tmp/test.log",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.logConfig.CreateLogFile(); (err != nil) != tc.expectError {
				if tc.expectError {
					if err == nil {
						t.Fatalf("expected error but got none")
					}
				} else {
					if err != nil {
						t.Fatalf("got unexpected error: %v", err)
					}

					expectedPath := tc.logConfig.LogPath
					if tc.logConfig.LogPath != expectedPath {
						t.Fatalf("expected path %s but got %s", expectedPath, tc.logConfig.LogPath)
					}
				}
			}
		})
	}
}

func TestConfigureLogger(t *testing.T) {
	testCases := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		logConfig  logging.LogConfig
		logName    string
		logFunc    func(l logging.Logger)
		outputType logging.OutputType
		errFunc    func(l logging.Logger)
		wantErr    bool
	}{
		{
			name:       "Info Level with Color Output",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			outputType: logging.ColorOutput,
			logName:    "test_info_color.log",
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
			logName:    "test_debug_plain.log",
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
			logName:    "/invalid_path/test.log",
			outputType: logging.PlainOutput,
			logFunc: func(l logging.Logger) {
				l.Debug("Test debug plain logger")
			},
			errFunc: func(l logging.Logger) {
				l.Debug("Test debug plain logger error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := afero.TempDir(tc.fs, "", "colorlogger_test-*")
			if err != nil {
				t.Fatalf("failed to create temp directory: %v", err)
			}
			defer removeDir(tc.fs, tempDir)

			cfg := logging.LogConfig{
				Fs:         tc.fs,
				LogPath:    filepath.Join(tempDir, tc.logName),
				Level:      tc.level,
				OutputType: tc.outputType,
			}

			log, err := cfg.ConfigureLogger()
			if (err != nil) != tc.wantErr {
				t.Errorf("ConfigureLogger() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if err == nil {
				tc.logFunc(log)
				if tc.errFunc != nil {
					tc.errFunc(log)
				}
			}
		})
	}
}

func TestGlobalLogger(t *testing.T) {
	testCases := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		logName    string
		outputType logging.OutputType
		logFunc    func(l logging.Logger)
		errFunc    func(l logging.Logger)
	}{
		{
			name:       "Set and Retrieve Global Logger",
			level:      slog.LevelInfo,
			fs:         afero.NewMemMapFs(),
			logName:    "test_global_logger.log",
			outputType: logging.PlainOutput,
			logFunc: func(l logging.Logger) {
				l.Println("Testing global logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Testing global logger error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use Afero's filesystem to create the temp directory
			tempDir, err := afero.TempDir(tc.fs, "", "colorlogger_test-*")
			if err != nil {
				t.Fatalf("failed to create temp directory: %v", err)
			}

			defer func() {
				if err := tc.fs.RemoveAll(tempDir); err != nil {
					t.Errorf("Failed to remove temporary directory: %v", err)
				}
			}()

			// Use the logName specified in each test case
			cfg := logging.LogConfig{
				Fs:         tc.fs,
				LogPath:    tc.logName,
				Level:      tc.level,
				OutputType: tc.outputType,
			}

			log, err := logging.InitLogging(&cfg)
			if err != nil {
				t.Fatalf("ConfigureLogger() error = %v", err)
			}

			t.Logf("Running test case: %s", tc.name)
			tc.logFunc(log)

			// Retrieve and use the global logger
			globalLogger := log
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
	testCases := []struct {
		name         string
		level        slog.Level
		fs           afero.Fs
		logName      string
		outputType   logging.OutputType
		outputToDisk bool
		logFunc      func(l logging.Logger)
		errFunc      func(l logging.Logger)
		expectError  bool
		expectedLog  string
	}{
		{
			name:         "Successful Logger Output",
			level:        slog.LevelInfo,
			fs:           afero.NewMemMapFs(),
			outputType:   log.ColorOutput,
			logName:      "test_logger_output.log",
			outputToDisk: true,
			logFunc: func(l logging.Logger) {
				l.Println("Test info color logger")
			},
			errFunc: func(l logging.Logger) {
				l.Error("Test info color logger error")
			},
			expectError: false,
			expectedLog: "Test info color logger",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := afero.TempDir(tc.fs, "", "loggeroutput_test-*")
			if err != nil {
				t.Fatalf("failed to create temp directory: %v", err)
			}
			defer removeDir(tc.fs, tempDir)

			logPath := filepath.Join(tempDir, tc.logName)
			cfg := logging.LogConfig{
				Fs:         tc.fs,
				LogPath:    logPath,
				LogToDisk:  true,
				Level:      tc.level,
				OutputType: tc.outputType,
			}

			logger, err := logging.InitLogging(&cfg)
			if (err != nil) != tc.expectError {
				t.Fatalf("InitLogging() error = %v, expectError %v", err, tc.expectError)
			}

			if !tc.expectError {
				tc.logFunc(logger)
				if tc.errFunc != nil {
					tc.errFunc(logger)
				}

				// Ensure the log file is closed and flushed
				if logFileCloser, ok := logger.(io.Closer); ok {
					if err := logFileCloser.Close(); err != nil {
						t.Fatalf("Failed to close log file: %v", err)
					}
				}

				// Check if log file exists and read its content
				logFile, err := tc.fs.Open(logPath)
				if err != nil {
					t.Fatalf("Failed to open log file: %v", err)
				}
				defer logFile.Close()

				buf, err := io.ReadAll(logFile)
				if err != nil {
					t.Fatalf("Failed to read log file: %v", err)
				}

				logContent := string(buf)
				if !strings.Contains(logContent, tc.expectedLog) {
					t.Errorf("Log file content does not contain the expected log message: '%s'", tc.expectedLog)
				}
			}
		})
	}
}

type mockLogger struct {
	lastLoggedError string
}

func (m *mockLogger) Error(v ...interface{}) {
	m.lastLoggedError = fmt.Sprint(v...)
}

// Implement no-op for other methods of Logger interface
func (m *mockLogger) Println(v ...interface{})               {}
func (m *mockLogger) Printf(format string, v ...interface{}) {}
func (m *mockLogger) Errorf(format string, v ...interface{}) {}
func (m *mockLogger) Close() error                           { return nil }
func (m *mockLogger) Debug(v ...interface{})                 {}
func (m *mockLogger) Debugf(format string, v ...interface{}) {}

func TestLogAndReturnError(t *testing.T) {
	tests := []struct {
		name    string
		errMsg  string
		wantErr string
	}{
		{
			name:    "Standard Error",
			errMsg:  "standard error occurred",
			wantErr: "standard error occurred",
		},
		{
			name:    "Empty Error Message",
			errMsg:  "",
			wantErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLogger{}
			err := logging.LogAndReturnError(mock, tc.errMsg)

			if err.Error() != tc.wantErr {
				t.Errorf("Expected error message '%s', got '%s'",
					tc.wantErr, err.Error())
			}

			if mock.lastLoggedError != tc.wantErr {
				t.Errorf("Expected logged error message '%s', got '%s'",
					tc.wantErr, mock.lastLoggedError)
			}
		})
	}
}
