package file_test

import (
	"log"
	"os"

	fileutils "github.com/l50/goutils/v2/file"
)

func ExampleRealFile_Open() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	rf := fileutils.RealFile(tmpfile.Name())

	reader, err := rf.Open()

	if err != nil {
		log.Printf("failed to open file: %v", err)
		return
	}

	_ = reader

	if err := reader.Close(); err != nil {
		log.Printf("failed to close file: %v", err)
		return
	}
}

func ExampleRealFile_RemoveAll() {
	tmpdir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp directory: %v", err)
		return
	}
	defer os.RemoveAll(tmpdir)

	rf := fileutils.RealFile(tmpdir)

	if err := rf.RemoveAll(); err != nil {
		log.Printf("failed to remove file or directory: %v", err)
		return
	}
}

func ExampleRealFile_Stat() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	rf := fileutils.RealFile(tmpfile.Name())

	if _, err := rf.Stat(); err != nil {
		log.Printf("failed to get file stat: %v", err)
		return
	}
}

func ExampleRealFile_Remove() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}

	rf := fileutils.RealFile(tmpfile.Name())

	if err := rf.Remove(); err != nil {
		log.Printf("failed to remove file: %v", err)
		return
	}
}

func ExampleRealFile_Write() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	rf := fileutils.RealFile(tmpfile.Name())

	err = rf.Write([]byte("Hello, World!"), 0644)

	if err != nil {
		log.Printf("failed to write to file: %v", err)
		return
	}
}

func ExampleRealFile_Append() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	rf := fileutils.RealFile(tmpfile.Name())

	if err := rf.Append("Hello, World!"); err != nil {
		log.Printf("failed to append to file: %v", err)
		return
	}
}

func ExampleCreate() {
	tmpdir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp directory: %v", err)
		return
	}
	defer os.RemoveAll(tmpdir)

	// Name of the file to be created
	fileName := tmpdir + "/testfile.txt"

	if err := fileutils.Create(fileName, []byte("Hello, World!"), fileutils.CreateFile); err != nil {
		log.Printf("failed to create file: %v", err)
		return
	}
}

func ExampleHasStr() {
	// Create a new temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some content to the file
	if _, err := tmpfile.WriteString("Hello, World!"); err != nil {
		log.Printf("failed to write to temp file: %v", err)
		return
	}
	tmpfile.Close()

	found, err := fileutils.HasStr(tmpfile.Name(), "World")
	if err != nil {
		log.Printf("failed to read from file: %v", err)
		return
	}

	if !found {
		log.Printf("failed to find string in file")
		return
	}
}

func ExampleCSVToLines() {
	tmpfile, err := os.CreateTemp("", "example.csv")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString("header1,header2\nvalue1,value2\nvalue3,value4"); err != nil {
		log.Printf("failed to write to temp file: %v", err)
		return
	}
	tmpfile.Close()

	records, err := fileutils.CSVToLines(tmpfile.Name())
	if err != nil {
		log.Printf("failed to read CSV file: %v", err)
		return
	}

	for _, row := range records {
		log.Println(row)
	}
}

func ExampleDelete() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}

	if err := fileutils.Delete(tmpfile.Name()); err != nil {
		log.Printf("failed to delete file: %v", err)
		return
	}
}

func ExampleExists() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	exists := fileutils.Exists(tmpfile.Name())
	if !exists {
		log.Printf("file does not exist")
		return
	}
}

func ExampleToSlice() {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString("Hello\nWorld"); err != nil {
		log.Printf("failed to write to temp file: %v", err)
		return
	}
	tmpfile.Close()

	lines, err := fileutils.ToSlice(tmpfile.Name())

	if err != nil {
		log.Printf("failed to read file: %v", err)
		return
	}

	for _, line := range lines {
		log.Println(line)
	}
}

func ExampleFind() {
	tmpdir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	tmpfile, err := os.CreateTemp(tmpdir, "file_to_find.txt")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}

	dirs := []string{tmpdir}

	filePaths, err := fileutils.Find(tmpfile.Name(), dirs)
	if err != nil {
		log.Printf("failed to find file: %v", err)
		return
	}

	for _, filePath := range filePaths {
		log.Printf("file found at: %s\n", filePath)
	}
}

func ExampleListR() {
	tmpdir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Printf("failed to create temp directory: %v", err)
		return
	}
	defer os.RemoveAll(tmpdir)

	if _, err := os.CreateTemp(tmpdir, "file1.txt"); err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}

	files, err := fileutils.ListR(tmpdir)
	if err != nil {
		log.Printf("failed to list files: %v", err)
		return
	}

	for _, file := range files {
		log.Println(file)
	}
}
