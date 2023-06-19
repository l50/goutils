package file_test

import (
	"encoding/csv"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/file"
	fileutils "github.com/l50/goutils/v2/file"
	"github.com/l50/goutils/v2/str"
	"github.com/stretchr/testify/require"
)

func TestAppend(t *testing.T) {
	tests := []struct {
		name string
		file string
		data string
	}{
		{
			name: "Appends data to file",
			file: "test.txt",
			data: "I am a change!!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := fileutils.Create(tc.file, []byte("I am a test!"), fileutils.CreateEmptyFile); err != nil {
				t.Fatalf("failed to create %s - CreateEmptyFile() failed: %v", tc.file, err)
			}
			exists := fileutils.Exists(tc.file)
			if !exists {
				t.Fatalf("unable to locate %s - exists() failed", tc.file)
			}

			if err := file.Append(tc.file, tc.data); err != nil {
				t.Fatalf("failed to append %s to %s - Append() failed: %v",
					tc.data, tc.file, err)
			}

			stringFoundInFile, err := fileutils.ContainsStr(tc.file, tc.data)
			if err != nil || !stringFoundInFile {
				t.Fatalf("failed to find %s in %s - StringInFile() failed: %v", tc.data, tc.file, err)
			}

			if exists {
				if err := file.Delete(tc.file); err != nil {
					t.Fatalf("unable to delete %s, DeleteFile() failed", tc.file)
				}
			} else {
				t.Fatalf("unable to create %s - CreateEmptyFile() failed", tc.file)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	tests := []struct {
		name       string
		path       string
		contents   []byte
		createType file.CreateType
		wantError  bool
	}{
		{
			name:       "create directory",
			path:       filepath.Join(tmpDir, "test_dir"),
			createType: file.CreateDirectory,
			wantError:  false,
		},
		{
			name:       "create empty file",
			path:       filepath.Join(tmpDir, "test_file.txt"),
			createType: file.CreateEmptyFile,
			wantError:  false,
		},
		{
			name:       "create file with content",
			path:       filepath.Join(tmpDir, "test_file_with_contents.txt"),
			contents:   []byte("Hello, World!"),
			createType: file.CreateFile,
			wantError:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := file.Create(tc.path, tc.contents, tc.createType)

			if (err != nil) != tc.wantError {
				t.Fatalf("create() error = %v, wantError %v", err, tc.wantError)
				return
			}

			// If it's a file creation and we don't expect an error,
			// check the contents of the file if needed
			if !tc.wantError && tc.createType != file.CreateDirectory && len(tc.contents) > 0 {
				contents, readErr := os.ReadFile(tc.path)
				if readErr != nil {
					t.Fatalf("cannot read file: %v", readErr)
					return
				}

				if string(contents) != string(tc.contents) {
					t.Fatalf("expected file contents %v, but got %v", string(tc.contents), string(contents))
				}
			}
		})
	}

	// Test for an already existing file/directory
	t.Run("create existing directory", func(t *testing.T) {
		existingDir := filepath.Join(tmpDir, "/existing_dir")
		err := file.Create(existingDir, nil, file.CreateDirectory)
		if err != nil {
			t.Fatalf("cannot create directory for testing: %v", err)
		}
	})

	// Cleanup: remove the temporary directory after all tests
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Errorf("failed to remove temporary directory: %v", err)
	}
}

func TestContainsStr(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		searchStr string
		expected  bool
		err       error
	}{
		{
			name:      "Returns true when string is found in file",
			input:     "find me",
			searchStr: "find me",
			expected:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "test.txt")
			if err := fileutils.Write(path, tc.input); err != nil {
				t.Fatalf("failed to write test data to %s: %v", path, err)
			}

			result, err := fileutils.ContainsStr(path, tc.searchStr)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, result)
		})
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

func TestExistsAndDeleteFile(t *testing.T) {
	t.Run("Test Exists and DeleteFile functions", func(t *testing.T) {
		// Create a temp directory
		tmpDir, err := os.MkdirTemp("", "testdir")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir) // clean up

		// Create a temp file
		tmpFile, err := os.CreateTemp(tmpDir, "testfile")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		filePath := tmpFile.Name() // Get the full file path

		// Check if the file exists
		exists := fileutils.Exists(filePath)
		if !exists {
			t.Fatalf("file does not exist but it should")
		}

		// Delete the file
		if err := fileutils.Delete(filePath); err != nil {
			t.Fatalf("Failed to delete file: %v", err)
		}

		// Check again if the file exists
		exists = fileutils.Exists(filePath)
		if exists {
			t.Fatalf("file exists but it should not")
		}
	})
}

func TestToSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		err      error
	}{
		{
			name:     "Reads file content into slice of strings",
			input:    "first line\nsecond line\nthird line",
			expected: []string{"first line", "second line", "third line"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "test.txt")
			if err := fileutils.Write(path, tc.input); err != nil {
				t.Fatalf("failed to write test data to %s: %v", path, err)
			}

			result, err := fileutils.ToSlice(path)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "TestFileExists",
			fileName: "testfile1.txt",
			wantErr:  false,
		},
		{
			name:     "TestFileDoesNotExist",
			fileName: "nonexistentfile.txt",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for each test case
			testDir, err := os.MkdirTemp("", "testdir")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(testDir) // Clean up the test directory

			// Create a file in the temporary directory if we're testing file existence
			if !tc.wantErr {
				filePath := filepath.Join(testDir, tc.fileName)
				if err := fileutils.Create(filePath, nil, fileutils.CreateEmptyFile); err != nil {
					t.Fatalf("unable to create test file: %v", err)
				}
			}

			files, err := file.Find(tc.fileName, []string{testDir})

			if (err != nil) != tc.wantErr {
				t.Fatalf("Find() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && (len(files) == 0 || !strings.Contains(files[0], tc.fileName)) {
				t.Fatalf("Find() did not return the correct files: %v", files)
			}
		})
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

func TestWrite(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "Writes string to file",
			input:    "Some content",
			expected: "Some content",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "test.txt")

			err := fileutils.Write(path, tc.input)
			require.Equal(t, tc.err, err)

			result, _ := fileutils.ToSlice(path)
			require.Equal(t, []string{tc.expected}, result)
		})
	}
}
