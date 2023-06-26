package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// LogInfo represents parameters used to manage logging throughout
// a program.
//
// **Attributes:**
//
// Dir: A string representing the directory where the log file is located.
// File: An afero.File object representing the log file.
// FileName: A string representing the name of the log file.
// Path: A string representing the full path to the log file.
type LogInfo struct {
	Dir      string
	File     afero.File
	FileName string
	Path     string
}

// CreateLogFile creates a log file in a 'logs' subdirectory of the
// specified directory. The log file's name is the provided log name
// with the extension '.log'.
//
// **Parameters:**
//
// fs: An afero.Fs instance to mock filesystem for testing.
// logDir: A string for the directory where 'logs' subdirectory and
// log file should be created.
// logName: A string for the name of the log file to be created.
//
// **Returns:**
//
// LogInfo: A LogInfo struct with information about the log file,
// including its directory, file pointer, file name, and path.
// error: An error, if an issue occurs while creating the directory
// or the log file.
func CreateLogFile(fs afero.Fs, logDir string, logName string) (LogInfo, error) {
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
	if _, err := fs.Stat(logInfo.Path); os.IsNotExist(err) {
		if err := fs.MkdirAll(logInfo.Dir, os.ModePerm); err != nil {
			return logInfo, fmt.Errorf("failed to create %s: %v", logInfo.Dir, err)
		}
	}

	// Create log file.
	logInfo.File, err = fs.OpenFile(logInfo.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return logInfo, fmt.Errorf("failed to create %s: %v", logInfo.Path, err)
	}

	return logInfo, nil
}
