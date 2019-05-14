package utils

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestFileExists(t *testing.T) {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		exists := FileExists(filepath.FromSlash("/etc/passwd"))
		if !exists {
			t.Fatal("Unable to locate /etc/passwd, FileExists() failed.")
		}
	} else if runtime.GOOS == "windows" {
		exists := FileExists(filepath.FromSlash("C:/WINDOWS/system32/win.ini"))
		if !exists {
			t.Fatal("Unable to locate C:/WINDOWS/system32/win.ini, FileExists() failed.")
		}
	} else {
		t.Fatal("Unsupported OS detected")
	}
}

func TestGetHomeDir(t *testing.T) {
	_, err := GetHomeDir()
	if err != nil {
		t.Fatal("Unable to get the user's home directory due to: ", err.Error())
	}
}
