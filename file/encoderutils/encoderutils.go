package encoderutils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
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
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}

		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
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
