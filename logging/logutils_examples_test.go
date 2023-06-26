package logging_test

import (
	"fmt"

	"github.com/l50/goutils/v2/logging"
	"github.com/spf13/afero"
)

func ExampleCreateLogFile() {
	fs := afero.NewOsFs()
	logDir := "/tmp"
	logName := "test.log"

	logInfo, err := logging.CreateLogFile(fs, logDir, logName)

	if err != nil {
		fmt.Printf("failed to create log file: %v", err)
		return
	}

	fmt.Printf("log file created at: %s", logInfo.Path)

	// Clean up
	err = fs.Remove(logInfo.Path)
	if err != nil {
		fmt.Printf("failed to clean up: %v", err)
	}

	// Output: log file created at: /tmp/logs/test.log
}
