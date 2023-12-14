package logging

import (
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
