package logging_test

import (
	"fmt"

	"log/slog"

	"github.com/fatih/color"
	"github.com/l50/goutils/v2/logging"
	"github.com/spf13/afero"
)

func ExampleConfigureLogger() {
	logger, err := logging.ConfigureLogger(slog.LevelDebug, "/tmp/test.log")
	if err != nil {
		fmt.Printf("failed to configure logger: %v", err)
		return
	}

	logger.Println("This is a log message")
	logger.Error("This is an error log message")
	logger.Errorf("This is a formatted error log message: %s", "Error details")

	// Since we can't predict the log message, print a static message instead.
	fmt.Println("Logger configured successfully.")

	// Unpredictable output due to timestamps and structured logging
}

func ExampleCreateLogFile() {
	fs := afero.NewOsFs()
	logDir := "/tmp"
	logName := "test.log"

	logInfo, err := logging.CreateLogFile(fs, logDir, logName)
	if err != nil {
		fmt.Printf("failed to create log file: %v", err)
		return
	}

	// Specify the color attribute for the ColoredLogger
	var logger = &logging.ColoredLogger{Info: logInfo, ColorAttribute: color.FgBlue}
	logger.Println("This is a log message")

	logger.Debug("This is a debug log message")
	logger.Debugf("This is a formatted debug log message: %s", "Debug details")

	fmt.Printf("log file created at: %s", logInfo.Path)

	// Clean up
	if err := fs.Remove(logInfo.Path); err != nil {
		fmt.Printf("failed to clean up: %v", err)
	}

	// Unpredictable output due to timestamps and structured logging
}
