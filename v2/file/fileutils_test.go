package file_test

import (
	"encoding/csv"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	fileutils "github.com/l50/goutils/v2/file"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"
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

func TestAppend(t *testing.T) {
	testFile := "test.txt"
	created := fileutils.CreateEmpty(testFile)
	exists := fileutils.Exists(testFile)
	change := "I am a change!!"

	if !exists {
		t.Fatalf("unable to locate %s - xists() failed", testFile)
	}

	if err := fileutils.Append(testFile, change); err != nil {
		t.Fatalf("failed to append %s to %s - Append() failed: %v",
			change, testFile, err)
	}

	stringFoundInFile, err := fileutils.FindStr(testFile, change)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v", change, testFile, err)
	}

	if created && exists {
		if err := fileutils.Delete(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	} else {
		t.Fatalf("unable to create %s - CreateEmptyFile() failed", testFile)
	}
}

func TestCreateEmptyFile(t *testing.T) {
	testFile := "test.txt"
	created := fileutils.CreateEmpty(testFile)
	exists := fileutils.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", testFile)
	}

	if created && exists {
		if err := fileutils.Delete(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	} else {
		t.Fatalf("unable to create %s, CreateEmptyFile() failed", testFile)
	}
}

func TestCreateFile(t *testing.T) {
	testFile := "test.txt"
	testFileContent := "stuff"

	if err := fileutils.Create(testFile, []byte(testFileContent)); err != nil {
		t.Fatalf("failed to create %s with %s using CreateFile(): %v", testFile, testFileContent, err)
	}

	exists := fileutils.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", testFile)
	}

	stringFoundInFile, err := fileutils.FindStr(testFile, testFileContent)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v",
			testFileContent, testFile, err)
	}

	if exists {
		if err := fileutils.Delete(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	} else {
		t.Fatalf("unable to create %s, CreateEmptyFile() failed", testFile)
	}
}

func TestCreateDirectory(t *testing.T) {
	rs, err := str.GenRandom(5)
	if err != nil {
		t.Fatal("failed to get random string for directory name with GenRandom()")

	}
	newDir := filepath.Join("/tmp", "bla", rs)
	if err := fileutils.CreateDirectory(newDir); err != nil {
		t.Fatalf("unable to create %s, CreateDirectory() failed: %v", newDir, err)
	}

	sameDir := newDir
	if err := fileutils.CreateDirectory(sameDir); err == nil {
		t.Fatalf("error: CreateDirectory() should not overwrite an existing directory")
	}
}

func TestDeleteFile(t *testing.T) {
	testFile := "test.txt"
	created := fileutils.CreateEmpty(testFile)
	exists := fileutils.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", testFile)
	}

	if created && exists {
		if err := fileutils.Delete(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	}
}

func TestCSVToLines(t *testing.T) {
	testFile := "test.csv"

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer file.Close()

	// write test data to file
	writer := csv.NewWriter(file)
	if err := writer.WriteAll([][]string{
		{"Header1", "Header2"},
		{"Data1", "Data2"},
		{"Data3", "Data4"},
	}); err != nil {
		t.Fatalf("failed to write test data to %p: %v", file, err)

	}
	writer.Flush()

	// call function under test
	got, err := fileutils.CSVToLines(testFile)
	if err != nil {
		t.Fatalf("CSVToLines failed: %v", err)
	}

	// check expected output
	want := [][]string{
		{"Data1", "Data2"},
		{"Data3", "Data4"},
	}

	if len(got) != len(want) {
		t.Errorf("unexpected number of records, got %d, want %d", len(got), len(want))
	}

	for i := range want {
		if !equal(got[i], want[i]) {
			t.Errorf("unexpected record %v, want %v", got[i], want[i])
		}
	}

	// Clean up
	if err := fileutils.Delete(testFile); err != nil {
		t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestExists(t *testing.T) {
	testFile := getTestFile(t)
	exists := fileutils.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, Exists() failed", testFile)
	}
}

func TestToSlice(t *testing.T) {
	testFile := getTestFile(t)
	exists := fileutils.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s - FileExists() failed", testFile)
	}

	if _, err := fileutils.ToSlice(testFile); err != nil {
		t.Fatalf("unable to convert %s to a slice - ToSlice() failed: %v",
			testFile, err)
	}
}

func TestFindFile(t *testing.T) {
	testFile := ".bashrc"
	if _, err := fileutils.Find(testFile, []string{"."}); err != nil {
		t.Fatalf("unable to find %s - FindFile() failed", testFile)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func TestListFilesR(t *testing.T) {
	targetPath := "."
	var results []string
	_, err := fileutils.ListR(targetPath)

	if err != nil {
		t.Fatalf("unable to list files in %s - ListFilesR() failed: %v",
			targetPath, err)
	}

	targetPath, err = str.GenRandom(randInt(8, 15))
	if err != nil {
		t.Fatalf("failed to generate RandomString - TestListFilesR() failed: %v",
			err)
	}

	if _, err = fileutils.ListR(targetPath); err == nil {
		t.Fatalf("%s should not exist - TestListFiles() failed", targetPath)
	}

	targetPath = "."
	results, err = fileutils.ListR(targetPath)
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
	stringFoundInFile, err := fileutils.FindStr(testFile, stringToFind)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v", stringToFind, testFile, err)
	}
}

func TestRmRf(t *testing.T) {
	rs, err := str.GenRandom(5)
	if err != nil {
		t.Fatal("failed to get random string for directory name with RandomString()")

	}
	newDir := filepath.Join("/tmp", "bla", rs)
	if err := fileutils.CreateDirectory(newDir); err != nil {
		t.Fatalf("unable to create %s, CreateDirectory() failed: %v", newDir, err)

	}

	if err := sys.RmRf(newDir); err != nil {
		t.Fatalf("unable to delete %s, RmRf() failed: %v", newDir, err)
	}
}

func TestExpandHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home directory: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "EmptyPath",
			input:    "",
			expected: "",
		},
		{
			name:     "NoTilde",
			input:    "/path/without/tilde",
			expected: "/path/without/tilde",
		},
		{
			name:     "TildeOnly",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "TildeWithSlash",
			input:    "~/path/with/slash",
			expected: filepath.Join(homeDir, "path/with/slash"),
		},
		{
			name:     "TildeWithoutSlash",
			input:    "~path/without/slash",
			expected: filepath.Join(homeDir, "path/without/slash"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := sys.ExpandHomeDir(tc.input)
			if actual != tc.expected {
				t.Errorf("test failed: ExpandHomeDir(%q) = %q; expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}
