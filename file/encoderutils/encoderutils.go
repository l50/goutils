package encoderutils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip unzips the specified zip file and extracts the contents of the
// zip file to the specified destination.
//
// **Parameters:**
//
// src: A string representing the path to the zip file.
//
// dest: A string representing the path to the destination directory.
//
// **Returns:**
//
// error: An error if any issue occurs while trying to unzip the file.
func Unzip(src, dest string) error {
	// This rule appears to be broken - false positive
	// nosemgrep:go.lang.security.zip.path-traversal-inside-zip-extraction
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	dest = filepath.Clean(dest) + string(os.PathSeparator)

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		tee := io.TeeReader(rc, &buf)

		if _, err := io.CopyN(&buf, tee, 1024*1024*256); err != nil && err != io.EOF { // Limit to 256MB
			rc.Close()
			return err
		}
		fmt.Println()

		path := filepath.Join(dest, file.Name)

		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		relPath, err := filepath.Rel(dest, path)
		if err != nil {
			rc.Close()
			return err
		}

		if strings.HasPrefix(relPath, "..") {
			rc.Close()
			return fmt.Errorf("illegal file path: %s", path)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				rc.Close()
				return err
			}
			rc.Close()
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			rc.Close()
			return err
		}

		destFile, err := os.Create(path)
		if err != nil {
			rc.Close()
			return err
		}

		if _, err = io.CopyN(destFile, rc, file.FileInfo().Size()); err != nil && err != io.EOF {
			destFile.Close()
			rc.Close()
			return err
		}

		destFile.Close()
		rc.Close()
	}

	return nil
}

// Zip creates a zip file from the specified source directory and saves it to the
// specified destination path.
//
// **Parameters:**
//
// srcDir: A string representing the path to the source directory.
//
// destFile: A string representing the path to the destination zip file.
//
// **Returns:**
//
// error: An error if any issue occurs while trying to zip the file.
func Zip(srcDir, destFile string) error {
	zipFile, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Skip the root directory
		if relPath == "." {
			return nil
		}

		// Create a new zip file entry
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Write the file content to the zip
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
