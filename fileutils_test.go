package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCreateEmptyFile(t *testing.T) {
	file := "test.txt"
	created := CreateEmptyFile(file)
	exists := FileExists(file)
	if created && exists {
		os.Remove(file)
	} else {
		t.Fatal("Unable to create ", file)
	}
}

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

func TestIsDirEmpty(t *testing.T) {
	dirEmpty, err := IsDirEmpty("/")
	if err != nil {
		t.Fatal("Unable to get the tmp directory due to: ", err.Error())
	}
	if dirEmpty != false {
		t.Fatal("The / directory has reported back as being empty, which can not be true.")
	}
}
