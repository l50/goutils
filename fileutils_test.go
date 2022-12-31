package utils

import (
	"math/rand"
	"path/filepath"
	"reflect"
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

func TestAppendToFile(t *testing.T) {
	testFile := "test.txt"
	created := CreateEmptyFile(testFile)
	exists := FileExists(testFile)
	change := "I am a change!!"

	if !exists {
		t.Fatalf("unable to locate %s - FileExists() failed", testFile)
	}

	err := AppendToFile(testFile, change)
	if err != nil {
		t.Fatalf("failed to append %s to %s - AppendToFile() failed: %v",
			change, testFile, err)
	}

	stringFoundInFile, err := StringInFile(testFile, change)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v", change, testFile, err)
	}

	if created && exists {
		err := DeleteFile(testFile)
		if err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	} else {
		t.Fatalf("unable to create %s - CreateEmptyFile() failed", testFile)
	}
}

func TestCreateEmptyFile(t *testing.T) {
	newFile := "test.txt"
	created := CreateEmptyFile(newFile)
	exists := FileExists(newFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", newFile)
	}

	if created && exists {
		err := DeleteFile(newFile)
		if err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", newFile)
		}
	} else {
		t.Fatalf("unable to create %s, CreateEmptyFile() failed", newFile)
	}
}

func TestCreateFile(t *testing.T) {
	newFile := "test.txt"
	contents := "stuff"

	err := CreateFile([]byte(contents), newFile)
	if err != nil {
		t.Fatalf("unable to create %s, CreateFile() failed: %v", newFile, err)
	}

	exists := FileExists(newFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", newFile)
	}

	stringFoundInFile, err := StringInFile(newFile, contents)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v",
			contents, newFile, err)
	}

	if exists {
		err := DeleteFile(newFile)
		if err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", newFile)
		}
	} else {
		t.Fatalf("unable to create %s, CreateEmptyFile() failed", newFile)
	}
}

func TestCreateDirectory(t *testing.T) {
	rs, err := RandomString(5)
	if err != nil {
		t.Fatal("failed to get random string for directory name with RandomString()")

	}
	newDir := filepath.Join("/tmp", "bla", rs)
	if err := CreateDirectory(newDir); err != nil {
		t.Fatalf("unable to create %s, CreateDirectory() failed: %v", newDir, err)
	}

	sameDir := newDir
	if err := CreateDirectory(sameDir); err == nil {
		t.Fatalf("error: CreateDirectory() should not overwrite an existing directory")
	}
}

func TestDeleteFile(t *testing.T) {
	newFile := "test.txt"
	created := CreateEmptyFile(newFile)
	exists := FileExists(newFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", newFile)
	}

	if created && exists {
		err := DeleteFile(newFile)
		if err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", newFile)
		}
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
		t.Fatalf("unable to convert %s to a slice - FileToSlice() failed: %v",
			testFile, err)
	}
}

func TestFindFile(t *testing.T) {
	testFile := ".bashrc"
	_, err := FindFile(testFile, []string{"."})
	if err != nil {
		t.Fatalf("unable to find %s - FindFile() failed", testFile)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func TestListFilesR(t *testing.T) {
	targetPath := "."
	var results []string
	_, err := ListFilesR(targetPath)

	if err != nil {
		t.Fatalf("unable to list files in %s - ListFilesR() failed: %v",
			targetPath, err)
	}

	targetPath, err = RandomString(randInt(8, 15))
	if err != nil {
		t.Fatalf("failed to generate RandomString - TestListFilesR() failed: %v",
			err)
	}

	_, err = ListFilesR(targetPath)

	if err == nil {
		t.Fatalf("%s should not exist - TestListFiles() failed", targetPath)
	}

	targetPath = "."
	results, err = ListFilesR(targetPath)
	if err != nil {
		t.Fatalf("unable to list files in %s - ListFiles() failed: %v",
			targetPath, err)
	}

	if reflect.TypeOf(results).Elem().Kind() != reflect.String {
		t.Fatalf(
			"error - unable to list files in %s - ListFiles() failed: %v",
			targetPath, err)
	}

}

func TestStringInFile(t *testing.T) {
	testFile := getTestFile(t)
	stringToFind := "root"
	stringFoundInFile, err := StringInFile(testFile, stringToFind)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v", stringToFind, testFile, err)
	}
}

func TestRmRf(t *testing.T) {
	rs, err := RandomString(5)
	if err != nil {
		t.Fatal("failed to get random string for directory name with RandomString()")

	}
	newDir := filepath.Join("/tmp", "bla", rs)
	if err := CreateDirectory(newDir); err != nil {
		t.Fatalf("unable to create %s, CreateDirectory() failed: %v", newDir, err)

	}

	if err := RmRf(newDir); err != nil {
		t.Fatalf("unable to delete %s, RmRf() failed: %v", newDir, err)
	}
}
