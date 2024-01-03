package logging_test

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/l50/goutils/v2/logging"
	log "github.com/l50/goutils/v2/logging"
	"github.com/spf13/afero"
)

func removeDir(fs afero.Fs, path string) {
	if err := fs.RemoveAll(path); err != nil {
		fmt.Printf("Failed to remove directory %s: %v\n", path, err)
	}
}

func TestColorLogger(t *testing.T) {
	testCases := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		logConfig  log.LogConfig
		logName    string
		logFunc    func(l log.Logger)
		outputType log.OutputType
		errFunc    func(l log.Logger)
	}{
		{
			name:  "Test PlainLogger Println",
			level: slog.LevelInfo,
			errFunc: func(l log.Logger) {
				l.Error("Test Plain Println logger error")
			},
			fs:         afero.NewMemMapFs(),
			outputType: log.PlainOutput,
			logName:    "test_color_println.log",
			logFunc: func(l log.Logger) {
				l.Println("Test Plain Println logger")
			},
		},
		{
			name:  "Test ColorLogger Printf",
			level: slog.LevelInfo,
			errFunc: func(l log.Logger) {
				l.Errorf("Test %s logger error", "Color Printf")
			},
			fs:         afero.NewMemMapFs(),
			outputType: log.ColorOutput,
			logName:    "test_color_printf.log",
			logFunc: func(l log.Logger) {
				l.Println("Test Color Printf logger")
			},
		},
		{
			name:  "Test ColorLogger Debug",
			level: slog.LevelDebug,
			errFunc: func(l log.Logger) {
				l.Debug("Test Color Debug logger with error")
			},
			fs:         afero.NewMemMapFs(),
			outputType: log.ColorOutput,
			logName:    "test_color_debug.log",
			logFunc: func(l log.Logger) {
				l.Debug("Test Color Debug logger")
			},
		},
		{
			name:  "Test ColorLogger Debugf",
			level: slog.LevelDebug,
			errFunc: func(l log.Logger) {
				l.Debugf("Test %s logger with error", "Color Debugf")
			},
			fs:         afero.NewMemMapFs(),
			outputType: log.ColorOutput,
			logName:    "test_color_debugf.log",
			logFunc: func(l log.Logger) {
				l.Debug("Test ColorDebugf logger")
			},
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
			if err != nil {
				t.Fatalf("ConfigureLogger() error = %v", err)
			}

			t.Logf("Running test case: %s", tc.name)
			tc.logFunc(log)
			if tc.errFunc != nil {
				tc.errFunc(log)
			}
		})
	}
}
