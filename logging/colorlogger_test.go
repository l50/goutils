package logging_test

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
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
		logConfig  logging.LogConfig
		logName    string
		logFunc    func(l log.Logger)
		outputType logging.OutputType
		errFunc    func(l log.Logger)
	}{
		{
			name:  "Test PlainLogger Println",
			level: slog.LevelInfo,
			errFunc: func(l logging.Logger) {
				l.Error("Test Plain Println logger error")
			},
			fs:         afero.NewMemMapFs(),
			outputType: log.PlainOutput,
			logName:    "test_color_println.log",
			logFunc: func(l logging.Logger) {
				l.Println("Test Plain Println logger")
			},
		},
		// {
		// 	name:  "Test ColorLogger Printf",
		// 	level: slog.LevelInfo,
		// 	errFunc: func(l logging.Logger) {
		// 		l.Errorf("Test %s logger error", "Color Printf")
		// 	},
		// 	fs:         afero.NewMemMapFs(),
		// 	outputType: logging.ColorOutput,
		// 	logName:    "test_color_printf.log",
		// 	logFunc: func(l logging.Logger) {
		// 		l.Println("Test Color Printf logger")
		// 	},
		// },
		{
			name:  "Test ColorLogger Debug",
			level: slog.LevelDebug,
			errFunc: func(l logging.Logger) {
				l.Debug("Test Color Debug logger with error")
			},
			fs:         afero.NewMemMapFs(),
			outputType: logging.ColorOutput,
			logName:    "test_color_debug.log",
			logFunc: func(l log.Logger) {
				l.Debug("Test Color Debug logger")
			},
		},
		{
			name:  "Test ColorLogger Debugf",
			level: slog.LevelDebug,
			errFunc: func(l logging.Logger) {
				l.Debugf("Test %s logger with error", "Color Debugf")
			},
			fs:         afero.NewMemMapFs(),
			outputType: logging.ColorOutput,
			logName:    "test_color_debugf.log",
			logFunc: func(l logging.Logger) {
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
				Level:      tc.level,
				LogPath:    filepath.Join(tempDir, tc.logName),
				LogToDisk:  true,
				OutputType: tc.outputType,
			}

			logger, err := logging.InitLogging(&cfg)
			if err != nil {
				t.Fatalf("InitLogging() error = %v", err)
			}

			t.Logf("Logger type: %T", logger) // Log the actual type of logger

			// Type assert to get the underlying *slog.Logger from the logger interface
			var slogLogger *slog.Logger
			switch v := logger.(type) {
			case *logging.ColorLogger:
				slogLogger = v.Logger
			case *logging.PlainLogger:
				slogLogger = v.Logger
			default:
				t.Fatalf("Unexpected logger type: %T", v)
			}

			if slogLogger == nil {
				t.Fatalf("slogLogger is nil")
			}

			colorLogger := &logging.ColorLogger{
				Info:           cfg,
				ColorAttribute: color.FgWhite,
				Logger:         slogLogger,
			}

			tc.logFunc(colorLogger)
			t.Logf("Running test case: %s", tc.name)
			if tc.errFunc != nil {
				tc.errFunc(colorLogger)
			}
		})
	}
}
