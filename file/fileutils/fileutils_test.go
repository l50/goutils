package file_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fileutils "github.com/l50/goutils/v2/file/fileutils"
	"github.com/l50/goutils/v2/str"
	"github.com/stretchr/testify/require"
)

func TestRealFile_Open(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Opens valid file",
			file:    "test.txt",
			wantErr: false,
		},
		{
			name:    "Fails to open non-existing file",
			file:    "non-existing.txt",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for the test
			tmpDir, err := os.MkdirTemp("", "")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir) // clean up

			// If the test case is not expected to result in an error,
			// create the test file in the temporary directory.
			if !tc.wantErr {
				err = os.WriteFile(filepath.Join(tmpDir, tc.file), []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			rf := fileutils.RealFile(filepath.Join(tmpDir, tc.file))
			_, err = rf.Open()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Open() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRealFile_RemoveAll(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Remove valid directory",
			file:    "test",
			wantErr: false,
		},
		{
			name:    "Remove non-existing directory",
			file:    "non-existing",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			if tc.name == "Remove valid directory" {
				err = os.Mkdir(filepath.Join(tmpDir, tc.file), 0755)
				if err != nil {
					t.Fatalf("failed to create test dir: %v", err)
				}
			}

			rf := fileutils.RealFile(filepath.Join(tmpDir, tc.file))
			err = rf.RemoveAll()
			if (err != nil) != tc.wantErr {
				t.Fatalf("RemoveAll() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRealFile_Stat(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Stat valid file",
			file:    "test.txt",
			wantErr: false,
		},
		{
			name:    "Fails to stat non-existing file",
			file:    "non-existing.txt",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			if !tc.wantErr {
				err = os.WriteFile(filepath.Join(tmpDir, tc.file), []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			rf := fileutils.RealFile(filepath.Join(tmpDir, tc.file))
			_, err = rf.Stat()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Stat() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRealFile_Remove(t *testing.T) {
	testCases := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Remove valid file",
			file:    "test.txt",
			wantErr: false,
		},
		{
			name:    "Fails to remove non-existing file",
			file:    "non-existing.txt",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			if !tc.wantErr {
				err = os.WriteFile(filepath.Join(tmpDir, tc.file), []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			rf := fileutils.RealFile(filepath.Join(tmpDir, tc.file))
			err = rf.Remove()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Remove() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRealFile_Write(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		contents []byte
		wantErr  bool
	}{
		{
			name:     "Writes to valid file",
			file:     "test.txt",
			contents: []byte("some data"),
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rf := fileutils.RealFile(tc.file)
			err := rf.Write(tc.contents, 0644)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Write() error = %v, wantErr %v", err, tc.wantErr)
			}

			// verify the contents if there is no error
			if !tc.wantErr {
				data, readErr := os.ReadFile(tc.file)
				if readErr != nil {
					t.Fatalf("Cannot read file: %v", readErr)
				}

				if !bytes.Equal(data, tc.contents) {
					t.Fatalf("File contents = %v, want %v", data, tc.contents)
				}

				// clean up
				if err := os.Remove(tc.file); err != nil {
					t.Fatalf("Cannot remove file: %v", err)
				}
			}
		})
	}
}

func TestAppend(t *testing.T) {
	testCases := []struct {
		name string
		file fileutils.RealFile
		data string
	}{
		{
			name: "Appends data to file",
			file: "test.txt",
			data: "I am a change!!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.file.Write([]byte("I am a test!"), 0644); err != nil {
				t.Fatalf("failed to create %s - Write() failed: %v", string(tc.file), err)
			}
			info, err := tc.file.Stat()
			if err != nil || info.IsDir() {
				t.Fatalf("unable to locate %s - Stat() failed", string(tc.file))
			}

			if err := tc.file.Append(tc.data); err != nil {
				t.Fatalf("failed to append %s to %s - Append() failed: %v",
					tc.data, string(tc.file), err)
			}

			rc, err := tc.file.Open()
			if err != nil {
				t.Fatalf("failed to open %s - Open() failed: %v", string(tc.file), err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil || !strings.Contains(string(content), tc.data) {
				t.Fatalf("failed to find %s in %s - ReadAll() or strings.Contains failed: %v", tc.data, string(tc.file), err)
			}

			if err := tc.file.RemoveAll(); err != nil {
				t.Fatalf("unable to delete %s, RemoveAll() failed", string(tc.file))
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

	testCases := []struct {
		name       string
		path       string
		contents   []byte
		createType fileutils.CreateType
		wantError  bool
	}{
		{
			name:       "create directory",
			path:       filepath.Join(tmpDir, "test_dir"),
			createType: fileutils.CreateDirectory,
			wantError:  false,
		},
		{
			name:       "create empty file",
			path:       filepath.Join(tmpDir, "test_file.txt"),
			createType: fileutils.CreateEmptyFile,
			wantError:  false,
		},
		{
			name:       "create file with content",
			path:       filepath.Join(tmpDir, "test_file_with_contents.txt"),
			contents:   []byte("Hello, World!"),
			createType: fileutils.CreateFile,
			wantError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := fileutils.Create(tc.path, tc.contents, tc.createType)

			if (err != nil) != tc.wantError {
				t.Fatalf("create() error = %v, wantError %v", err, tc.wantError)
				return
			}

			// If it's a file creation and we don't expect an error,
			// check the contents of the file if needed
			if !tc.wantError && tc.createType != fileutils.CreateDirectory && len(tc.contents) > 0 {
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
		err := fileutils.Create(existingDir, nil, fileutils.CreateDirectory)
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
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "test.txt")
			if err := fileutils.RealFile(path).Write([]byte(tc.input), 0644); err != nil {
				t.Fatalf("failed to write test data to %s: %v", path, err)
			}

			result, err := fileutils.HasStr(path, tc.searchStr)
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
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir) // clean up

		// Create a temp file
		tmpFile, err := os.CreateTemp(tmpDir, "testfile")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		filePath := tmpFile.Name() // Get the full file path

		// Check if the file exists
		exists := fileutils.Exists(filePath)
		if !exists {
			t.Fatalf("file does not exist but it should")
		}

		// Delete the file
		if err := fileutils.Delete(filePath); err != nil {
			t.Fatalf("failed to delete file: %v", err)
		}

		// Check again if the file exists
		exists = fileutils.Exists(filePath)
		if exists {
			t.Fatalf("file exists but it should not")
		}
	})
}

func TestToSlice(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "test.txt")

			if err := fileutils.RealFile(path).Write([]byte(tc.input), 0644); err != nil {
				t.Fatalf("failed to write test data to %s: %v", path, err)
			}

			result, err := fileutils.ToSlice(path)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestFind(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for each test case
			testDir, err := os.MkdirTemp("", "testdir")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(testDir) // Clean up the test directory

			// Create a file in the temporary directory if we're testing file existence
			if !tc.wantErr {
				filePath := filepath.Join(testDir, tc.fileName)
				if err := fileutils.Create(filePath, nil, fileutils.CreateEmptyFile); err != nil {
					t.Fatalf("unable to create test file: %v", err)
				}
			}

			files, err := fileutils.Find(tc.fileName, []string{testDir})
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
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "test.txt")

			err := fileutils.RealFile(path).Write([]byte(tc.input), 0644)
			require.Equal(t, tc.err, err)

			result, _ := fileutils.ToSlice(path)
			require.Equal(t, []string{tc.expected}, result)
		})
	}
}

func TestSeekAndDestroy(t *testing.T) {
	// Setup a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a few test files in the temporary directory
	testFiles := []string{"file1.txt", "file2.log", "file3.txt", "file4.doc"}
	for _, file := range testFiles {
		_, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatal(err)
		}
	}

	testCases := []struct {
		name      string
		path      string
		pattern   string
		expectErr bool
	}{
		{
			name:      "existing file",
			path:      tmpDir,
			pattern:   "*.txt",
			expectErr: false,
		},
		{
			name:      "non-existent file",
			path:      tmpDir,
			pattern:   "*.jpg",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := fileutils.SeekAndDestroy(tc.path, tc.pattern)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
