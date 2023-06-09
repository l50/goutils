package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LogInfo contains parameters used to provide logging throughout a program.
type LogInfo struct {
	Dir      string
	FilePtr  *os.File
	FileName string
	Path     string
}

// CreateLogFile creates a log file in a directory named logs, which is a subdirectory of the given directory.
// The log file's name is the base name of the given directory with the extension .log.
//
// Parameters:
//
// logDir: A string representing the directory where the logs subdirectory should be created.
//
// Returns:
//
// LogInfo: A struct containing information about the log file, including its directory, file pointer, file name, and path.
//
// error: An error if the log file or its parent directory cannot be created.
//
// Example:
//
// logDir := "/path/to/logDir"
// logName := "stuff.log"
// logInfo, err := logging.CreateLogFile(logDir, logName)
//
//	if err != nil {
//	  fmt.Printf("failed to create log file: %v", err)
//	}
//
// fmt.Printf("Log file created at: %s", logInfo.Path)
func CreateLogFile(logDir string, logName string) (LogInfo, error) {
	logInfo := LogInfo{}
	var err error

	logDir = strings.TrimSpace(logDir)
	logName = strings.TrimSpace(logName)

	if logDir == "" {
		return logInfo, fmt.Errorf("logDir cannot be empty")
	}

	if logName == "" {
		return logInfo, fmt.Errorf("logName cannot be empty")
	}

	logInfo.Dir = filepath.Join(logDir, "logs")

	if filepath.Ext(logName) != ".log" {
		logInfo.FileName = fmt.Sprintf("%s.log", logName)
	} else {
		logInfo.FileName = logName
	}

	logInfo.Path = filepath.Join(logInfo.Dir, logInfo.FileName)

	// Create path to log file if the log file doesn't already exist.
	if _, err := os.Stat(logInfo.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(logInfo.Dir, os.ModePerm); err != nil {
			return logInfo, fmt.Errorf("failed to create %s: %v", logInfo.Dir, err)
		}
	}

	// Create log file.
	logInfo.FilePtr, err = os.OpenFile(logInfo.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return logInfo, fmt.Errorf("failed to create %s: %v", logInfo.Path, err)
	}

	return logInfo, nil
}
