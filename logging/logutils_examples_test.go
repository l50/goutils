package logging_test

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/l50/goutils/v2/logging"
	"github.com/spf13/afero"
)

func plainLoggerExample() {
	cfg := logging.LogConfig{
		Fs:         afero.NewOsFs(),
		LogPath:    filepath.Join("/tmp", "test.log"),
		Level:      slog.LevelDebug,
		OutputType: logging.PlainOutput,
		LogToDisk:  true,
	}

	logger, err := logging.InitLogging(&cfg)
	if err != nil {
		fmt.Printf("Failed to configure logger: %v", err)
		return
	}

	logger.Println("This is a log message")
	logger.Error("This is an error log message")
	logger.Errorf("This is a formatted error log message: %s", "Error details")

	fmt.Println("Logger configured successfully.")
}

func colorLoggerExample() {
	cfg := logging.LogConfig{
		Fs:         afero.NewOsFs(),
		LogPath:    filepath.Join("/tmp", "test.log"),
		Level:      slog.LevelDebug,
		OutputType: logging.ColorOutput,
		LogToDisk:  true,
	}

	logger, err := logging.InitLogging(&cfg)
	if err != nil {
		fmt.Printf("Failed to configure logger: %v", err)
		return
	}

	logger.Println("This is a log message")
	logger.Error("This is an error log message")
	logger.Errorf("This is a formatted error log message: %s", "Error details")

	fmt.Println("Logger configured successfully.")
}

func ExampleLogConfig_ConfigureLogger() {
	plainLoggerExample()
	colorLoggerExample()
}

func ExampleLogConfig_CreateLogFile() {
	cfg := logging.LogConfig{
		Fs:         afero.NewOsFs(),
		LogPath:    filepath.Join("/tmp", "test.log"),
		Level:      slog.LevelDebug,
		OutputType: logging.ColorOutput,
		LogToDisk:  true,
	}

	fmt.Println("Creating log file...")
	if err := cfg.CreateLogFile(); err != nil {
		fmt.Printf("Failed to create log file: %v", err)
		return
	}

	fmt.Printf("Log file created at: %s", cfg.LogPath)

	if err := cfg.Fs.Remove(cfg.LogPath); err != nil {
		fmt.Printf("Failed to clean up: %v", err)
	}
}

func ExampleInitLogging() {
	cfg := logging.LogConfig{
		Fs:         afero.NewOsFs(),
		Level:      slog.LevelDebug,
		OutputType: logging.ColorOutput,
		LogToDisk:  true,
		LogPath:    filepath.Join("/tmp", "test.log"),
	}

	log, err := logging.InitLogging(&cfg)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	log.Println("This is a test info message")
	log.Printf("This is a test %s info message", "formatted")
	log.Error("This is a test error message")
	log.Debugf("This is a test debug message")
	log.Errorf("This is a test %s error message", "formatted")
	log.Println("{\"time\":\"2024-01-03T23:12:35.937476-07:00\",\"level\":\"ERROR\",\"msg\":\"\\u001b[1;32m==> docker.ansible-attack-box: Starting docker container...\\u001b[0m\"}")
}
