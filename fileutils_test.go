package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Helper function that returns a test file based on the OS of the system
func getTestFile(t *testing.T) string {
	t.Helper()
	var testFile string
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		testFile = filepath.FromSlash("/etc/passwd")
	} else if runtime.GOOS == "windows" {
		testFile = filepath.FromSlash("C:/WINDOWS/system32/win.ini")
	} else {
		t.Fatal("Unsupported OS detected")
	}
	return testFile
}

func TestCreateEmptyFile(t *testing.T) {
	newFile := "test.txt"
	created := CreateEmptyFile(newFile)
	exists := FileExists(newFile)

	if !exists {
		t.Fatalf("Unable to locate %s, FileExists() failed.", newFile)
	}

	if created && exists {
		os.Remove(newFile)
	} else {
		t.Fatalf("Unable to create %s, CreateEmptyFile() failed.", newFile)
	}
}

func TestFileExists(t *testing.T) {
	testFile := getTestFile(t)
	exists := FileExists(testFile)

	if !exists {
		t.Fatalf("Unable to locate %s, FileExists() failed.", testFile)
	}
}

func TestFileToSlice(t *testing.T) {
	testFile := getTestFile(t)
	exists := FileExists(testFile)

	if !exists {
		t.Fatalf("Unable to locate %s, FileExists() failed.\n", testFile)
	}

	_, err := FileToSlice(testFile)
	if err != nil {
		t.Fatalf("Unable to convert %s to a slice due to: %v; FileToSlice() failed.\n", testFile, err.Error())
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
