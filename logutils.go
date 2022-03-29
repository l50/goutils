package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// LogInfo contains parameters associated with logging throughout the program.
type LogInfo struct {
	Dir      string
	FilePtr  *os.File
	FileName string
	Path     string
}

// CreateLogFile creates logs/ in the
// cwd. The cwd name is used to name
// the log file.
func CreateLogFile() (LogInfo, error) {
	logInfo := LogInfo{}

	// Get current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return logInfo, fmt.Errorf(color.RedString(
			"failed to get working directory: %v", err))
	}
	logInfo.Dir = filepath.Join(cwd, "logs")

	// Get the current directory. This is
	// used to name the log file.
	workDir := filepath.Base(cwd)

	logInfo.FileName = workDir + ".log"
	logInfo.Path = filepath.Join(logInfo.Dir, logInfo.FileName)

	// Create path to log file if the log file doesn't already exist.
	if _, err := os.Stat(logInfo.Path); os.IsNotExist(err) {
		err = os.MkdirAll(logInfo.Dir, os.ModePerm)
		if err != nil {
			return logInfo, fmt.Errorf(color.RedString(
				"failed to create %s:%v", logInfo.Dir, err))
		}
	}

	// Create log file.
	logInfo.FilePtr, err = os.OpenFile(logInfo.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return logInfo, fmt.Errorf(color.RedString(
			"failed to create %s: %v", logInfo.Path, err))
	}

	return logInfo, nil
}
