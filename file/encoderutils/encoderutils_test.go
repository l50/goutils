package encoderutils_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	encoder "github.com/l50/goutils/v2/file/encoderutils"
)

func TestUnzipAndZip(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name:     "Unzip",
			testFunc: testUnzip,
		},
		{
			name:     "Zip",
			testFunc: testZip,
		},
		{
			name:     "UnzipTraversalAttack",
			testFunc: testUnzipTraversalAttack,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}
}

func testUnzip(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_unzip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Specify the source directory for zipping
	zipFilePath := filepath.Join(tempDir, "test.zip")
	err = createTestZipFile(zipFilePath)
	if err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(tempDir, "extracted")

	// Attempt to unzip the input zip file
	err = encoder.Unzip(zipFilePath, destDir)
	if err != nil {
		t.Fatalf("Unzip failed unexpectedly: %v", err)
	}

	// Check if the file was extracted
	extractedFilePath := filepath.Join(destDir, "file1.txt")
	_, err = os.Stat(extractedFilePath)
	if err != nil {
		t.Errorf("Expected file %s doesn't exist", extractedFilePath)
	}
}

func testZip(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Specify the source directory for zipping
	srcDir := filepath.Join(tempDir, "source")
	err = createTestSourceDir(srcDir)
	if err != nil {
		t.Fatal(err)
	}

	// Specify the destination path for the zip file
	destFilePath := filepath.Join(tempDir, "test.zip")

	// Call the Zip function
	err = encoder.Zip(srcDir, destFilePath)
	if err != nil {
		t.Fatalf("Zip failed: %v", err)
	}

	// Check if the zip file exists
	_, err = os.Stat(destFilePath)
	if err != nil {
		t.Errorf("Expected file %s doesn't exist", destFilePath)
	}
}

// createTestZipFile creates a test zip file for the purpose of testing the Unzip function.
func createTestZipFile(zipFilePath string) error {
	file, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Create a test file inside the zip
	fileWriter, err := zipWriter.Create("file1.txt")
	if err != nil {
		return err
	}

	// Write some data to the test file
	data := []byte("This is a test file.")
	_, err = fileWriter.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// createTestSourceDir creates a test source directory for the purpose of testing the Zip function.
func createTestSourceDir(srcDir string) error {
	if err := os.MkdirAll(srcDir, os.ModePerm); err != nil {
		return err
	}

	// Create a test file inside the source directory
	filePath := filepath.Join(srcDir, "file1.txt")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write some data to the test file
	data := []byte("This is a test file.")
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func createTestZipFileTraversal(zipFilePath string) error {
	file, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Create a test file inside the zip with a traversal path like "../file1.txt"
	// This attempts to create a file outside of the intended directory when unzipped
	fileWriter, err := zipWriter.Create("../../file1.txt") // Attempting path traversal
	if err != nil {
		return err
	}

	// Write some data to the test file
	data := []byte("This is a test file.")
	_, err = fileWriter.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func testUnzipTraversalAttack(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_unzip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	zipFilePath := filepath.Join(tempDir, "test.zip")
	err = createTestZipFileTraversal(zipFilePath)
	if err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(tempDir, "extracted")

	err = encoder.Unzip(zipFilePath, destDir)
	if err == nil {
		t.Fatal("Expected an error for path traversal, but got none")
	}

	if !strings.Contains(err.Error(), "illegal file path") {
		t.Fatalf("Expected an 'illegal file path' error, got: %v", err)
	}

	extractedFilePath := filepath.Join(tempDir, "file1.txt")
	_, err = os.Stat(extractedFilePath)
	if !os.IsNotExist(err) {
		t.Errorf("File was written outside the destination directory: %s", extractedFilePath)
	}
}
