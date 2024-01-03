package logging_test

import (
	"log/slog"
	"testing"

	log "github.com/l50/goutils/v2/logging"
	"github.com/spf13/afero"
)

func TestPlainLogger(t *testing.T) {
	testCases := []struct {
		name       string
		level      slog.Level
		fs         afero.Fs
		logConfig  log.LogConfig
		logFunc    func(l log.Logger)
		outputType log.OutputType
		errFunc    func(l log.Logger)
		logName    string
	}{
		{
			name:       "Test PlainLogger Println",
			level:      slog.LevelInfo,
			logName:    "test_plain_println.log",
			fs:         afero.NewMemMapFs(),
			outputType: log.PlainOutput,
			logFunc: func(l log.Logger) {
				l.Println("Test Plain Println logger")
			},
			errFunc: func(l log.Logger) {
				l.Error("Test Plain Println logger error")
			},
		},
		{
			name:       "Test PlainLogger Printf",
			level:      slog.LevelInfo,
			logName:    "test_plain_printf.log",
			fs:         afero.NewMemMapFs(),
			outputType: log.PlainOutput,
			logFunc: func(l log.Logger) {
				l.Printf("Test %s logger", "Plain Printf")
			},
			errFunc: func(l log.Logger) {
				l.Errorf("Test %s logger error", "Plain Printf")
			},
		},
		{
			name:       "Test PlainLogger Debug",
			level:      slog.LevelDebug,
			logName:    "test_plain_debug.log",
			fs:         afero.NewMemMapFs(),
			outputType: log.PlainOutput,
			logFunc: func(l log.Logger) {
				l.Debug("Test Plain Debug logger")
			},
			errFunc: func(l log.Logger) {
				l.Debug("Test Plain Debug logger with error")
			},
		},
		{
			name:       "Test PlainLogger Debugf",
			level:      slog.LevelDebug,
			logName:    "test_plain_debugf.log",
			fs:         afero.NewMemMapFs(),
			outputType: log.PlainOutput,
			logFunc: func(l log.Logger) {
				l.Debugf("Test %s logger", "Plain Debugf")
			},
			errFunc: func(l log.Logger) {
				l.Debugf("Test %s logger with error", "Plain Debugf")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := log.LogConfig{
				Fs:         tc.fs,
				LogPath:    tc.logName,
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
