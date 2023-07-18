package encoderutils_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	encoder "github.com/l50/goutils/v2/file/encoderutils"
)

func TestUnzip(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test_unzip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test zip file in the temporary directory
	zipFilePath := filepath.Join(tempDir, "test.zip")
	err = createTestZipFile(zipFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Specify the destination directory for unzipping
	destDir := filepath.Join(tempDir, "extracted")

	// Call the Unzip function
	err = encoder.Unzip(zipFilePath, destDir)
	if err != nil {
		t.Fatalf("Unzip failed: %v", err)
	}

	// Check if the extracted file exists
	extractedFilePath := filepath.Join(destDir, "file1.txt")
	_, err = os.Stat(extractedFilePath)
	if err != nil {
		t.Errorf("Expected file %s doesn't exist", extractedFilePath)
	}
}

func TestZip(t *testing.T) {
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
	err := os.MkdirAll(srcDir, os.ModePerm)
	if err != nil {
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
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
