package file_test

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	goutils "github.com/l50/goutils"
	"github.com/l50/goutils/v2/file"
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
	created := file.CreateEmptyFile(testFile)
	exists := goutils.FileExists(testFile)
	change := "I am a change!!"

	if !exists {
		t.Fatalf("unable to locate %s - FileExists() failed", testFile)
	}

	if err := file.AppendToFile(testFile, change); err != nil {
		t.Fatalf("failed to append %s to %s - AppendToFile() failed: %v",
			change, testFile, err)
	}

	stringFoundInFile, err := goutils.StringInFile(testFile, change)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v", change, testFile, err)
	}

	if created && exists {
		if err := file.DeleteFile(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	} else {
		t.Fatalf("unable to create %s - CreateEmptyFile() failed", testFile)
	}
}

func TestCreateEmptyFile(t *testing.T) {
	testFile := "test.txt"
	created := goutils.CreateEmptyFile(testFile)
	exists := goutils.FileExists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", testFile)
	}

	if created && exists {
		if err := file.DeleteFile(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	} else {
		t.Fatalf("unable to create %s, CreateEmptyFile() failed", testFile)
	}
}

func TestCreateFile(t *testing.T) {
	testFile := "test.txt"
	testFileContent := "stuff"

	if err := file.CreateFile(testFile, []byte(testFileContent)); err != nil {
		t.Fatalf("failed to create %s with %s using CreateFile(): %v", testFile, testFileContent, err)
	}

	exists := goutils.FileExists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", testFile)
	}

	stringFoundInFile, err := goutils.StringInFile(testFile, testFileContent)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v",
			testFileContent, testFile, err)
	}

	if exists {
		if err := file.DeleteFile(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	} else {
		t.Fatalf("unable to create %s, CreateEmptyFile() failed", testFile)
	}
}

func TestCreateDirectory(t *testing.T) {
	rs, err := goutils.RandomString(5)
	if err != nil {
		t.Fatal("failed to get random string for directory name with RandomString()")

	}
	newDir := filepath.Join("/tmp", "bla", rs)
	if err := file.CreateDirectory(newDir); err != nil {
		t.Fatalf("unable to create %s, CreateDirectory() failed: %v", newDir, err)
	}

	sameDir := newDir
	if err := file.CreateDirectory(sameDir); err == nil {
		t.Fatalf("error: CreateDirectory() should not overwrite an existing directory")
	}
}

func TestDeleteFile(t *testing.T) {
	testFile := "test.txt"
	created := file.CreateEmptyFile(testFile)
	exists := file.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, FileExists() failed", testFile)
	}

	if created && exists {
		if err := file.DeleteFile(testFile); err != nil {
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
	got, err := goutils.CSVToLines(testFile)
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
	if err := goutils.DeleteFile(testFile); err != nil {
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
	exists := file.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s, Exists() failed", testFile)
	}
}

func TestToSlice(t *testing.T) {
	testFile := getTestFile(t)
	exists := file.Exists(testFile)

	if !exists {
		t.Fatalf("unable to locate %s - FileExists() failed", testFile)
	}

	if _, err := goutils.FileToSlice(testFile); err != nil {
		t.Fatalf("unable to convert %s to a slice - FileToSlice() failed: %v",
			testFile, err)
	}
}

func TestFindFile(t *testing.T) {
	testFile := ".bashrc"
	if _, err := file.FindFile(testFile, []string{"."}); err != nil {
		t.Fatalf("unable to find %s - FindFile() failed", testFile)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func TestListFilesR(t *testing.T) {
	targetPath := "."
	var results []string
	_, err := file.ListFilesR(targetPath)

	if err != nil {
		t.Fatalf("unable to list files in %s - ListFilesR() failed: %v",
			targetPath, err)
	}

	targetPath, err = goutils.RandomString(randInt(8, 15))
	if err != nil {
		t.Fatalf("failed to generate RandomString - TestListFilesR() failed: %v",
			err)
	}

	if _, err = file.ListFilesR(targetPath); err == nil {
		t.Fatalf("%s should not exist - TestListFiles() failed", targetPath)
	}

	targetPath = "."
	results, err = file.ListFilesR(targetPath)
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
	stringFoundInFile, err := goutils.StringInFile(testFile, stringToFind)
	if err != nil || !stringFoundInFile {
		t.Fatalf("failed to find %s in %s - StringInFile() failed: %v", stringToFind, testFile, err)
	}
}

func TestRmRf(t *testing.T) {
	rs, err := goutils.RandomString(5)
	if err != nil {
		t.Fatal("failed to get random string for directory name with RandomString()")

	}
	newDir := filepath.Join("/tmp", "bla", rs)
	if err := file.CreateDirectory(newDir); err != nil {
		t.Fatalf("unable to create %s, CreateDirectory() failed: %v", newDir, err)

	}

	if err := file.RmRf(newDir); err != nil {
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
			actual := goutils.ExpandHomeDir(tc.input)
			if actual != tc.expected {
				t.Errorf("test failed: ExpandHomeDir(%q) = %q; expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestFindExportedFuncsWithoutTests(t *testing.T) {
	pkg := "bla"
	// Create temporary directory
	tempDir, err := os.MkdirTemp("/tmp", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file with exported function
	file1 := filepath.Join(tempDir, "file1.go")
	content1 := fmt.Sprintf(`package %s
func ExportedFunc1() {}
`, pkg)
	if err := os.WriteFile(file1, []byte(content1), 0666); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}

	// Create example with exported function and test function
	file2 := filepath.Join(tempDir, "file2.go")
	content2 := fmt.Sprintf(`package %s
func ExportedFunc2() {}
`, pkg)
	if err := os.WriteFile(file2, []byte(content2), 0666); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}

	file2Test := filepath.Join(tempDir, "file2_test.go")
	content2Test := fmt.Sprintf(`package %s
import "testing"
func TestExportedFunc2(t *testing.T) {}
`, pkg)
	if err := os.WriteFile(file2Test, []byte(content2Test), 0666); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Create a file with exported function and no test function
	file3 := filepath.Join(tempDir, "pkg", "bla", "file3.go")
	content3 := fmt.Sprintf(`package %s
func ExportedFunc3() {}
`, pkg)
	if err := os.MkdirAll(filepath.Dir(file3), os.ModePerm); err != nil {
		t.Fatalf("failed to create file3 dir: %v", err)
	}
	if err := os.WriteFile(file3, []byte(content3), 0666); err != nil {
		t.Fatalf("failed to create file3: %v", err)
	}

	// Create a file with exported function and test function
	file4 := filepath.Join(tempDir, "pkg", "bla", "file3_test.go")
	content4 := fmt.Sprintf(`package %s
import "testing"
func TestExportedFunc3(t *testing.T) {}
`, pkg)
	if err := os.WriteFile(file4, []byte(content4), 0666); err != nil {
		t.Fatalf("failed to create file4: %v", err)
	}

	// Call FindExportedFuncsWithoutTests
	exportedFuncs, err := file.FindExportedFuncsWithoutTests(tempDir)
	if err != nil {
		t.Fatalf("failed to find exported funcs: %v", err)
	}

	// Assert the result
	expectedFuncs := []string{"ExportedFunc1"}
	if !reflect.DeepEqual(exportedFuncs, expectedFuncs) {
		t.Errorf("expected funcs: %v, got: %v", expectedFuncs, exportedFuncs)
	}
}
