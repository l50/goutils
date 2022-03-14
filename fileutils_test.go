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
	switch runtime.GOOS {
	case "linux", "darwin":
		testFile = filepath.FromSlash("/etc/passwd")
	case "windows":
		testFile = filepath.FromSlash("C:/WINDOWS/system32/win.ini")
	default:
		t.Fatal("unsupported OS detected")
	}
	return testFile
}

func TestCreateEmptyFile(t *testing.T) {
	newFile := "test.txt"
	created := CreateEmptyFile(newFile)
	exists := FileExists(newFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", newFile)
	}

	if created && exists {
		os.Remove(newFile)
	} else {
		t.Fatalf("unable to create %s, CreateEmptyFile() failed", newFile)
	}
}

func TestFileExists(t *testing.T) {
	testFile := getTestFile(t)
	exists := FileExists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", testFile)
	}
}

func TestFileToSlice(t *testing.T) {
	testFile := getTestFile(t)
	exists := FileExists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s - FileExists() failed", testFile)
	}

	_, err := FileToSlice(testFile)
	if err != nil {
		t.Fatalf("unable to convert %s to a slice - FileToSlice() failed: %v", testFile, err)
	}
}

func TestGetHomeDir(t *testing.T) {
	_, err := GetHomeDir()
	if err != nil {
		t.Fatalf("unable to get the user's home directory - GetHomeDir() failed: %v", err)
	}
}

func TestIsDirEmpty(t *testing.T) {
	dirEmpty, err := IsDirEmpty("/")
	if err != nil {
		t.Fatalf("failed to determine if / is empty - IsDirEmpty() failed: %v", err)
	}
	if dirEmpty != false {
		t.Fatal("the / directory has reported back as being empty, which can not be true - IsDirEmpty()")
	}
}
