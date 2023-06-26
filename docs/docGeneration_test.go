package docs_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/l50/goutils/v2/docs"
	"github.com/spf13/afero"
)

type mockDirEntry struct {
	name  string
	isDir bool
	info  fs.FileInfo
}

func createTestDirectoryAndFile(afs afero.Fs) (string, error) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		return "", err
	}

	// Create a test file inside the temporary directory
	filePath := filepath.Join(dir, "testfile.txt")
	file, err := afero.NewOsFs().Create(filePath)
	if err != nil {
		os.RemoveAll(dir) // Clean up the temporary directory if file creation fails
		return "", err
	}
	defer file.Close()

	// Write some content to the test file
	content := []byte("Hello, Test File!")
	_, err = file.Write(content)
	if err != nil {
		os.RemoveAll(dir) // Clean up the temporary directory if writing fails
		return "", err
	}

	return dir, nil
}

func (d *mockDirEntry) Name() string {
	return d.name
}

func (d *mockDirEntry) IsDir() bool {
	return d.isDir
}

func (d *mockDirEntry) Type() fs.FileMode {
	return 0
}

func (d *mockDirEntry) Info() (fs.FileInfo, error) {
	return d.info, nil
}

type mockFileInfo struct {
	name  string
	isDir bool
	mode  fs.FileMode
}

func (m *mockFileInfo) Name() string {
	return m.name
}

func (m *mockFileInfo) Size() int64 {
	return 0
}

func (m *mockFileInfo) Mode() fs.FileMode {
	return m.mode
}

func (m *mockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (m *mockFileInfo) IsDir() bool {
	return m.isDir
}

func (m *mockFileInfo) Sys() interface{} {
	return nil
}

func TestCreatePackageDocs(t *testing.T) {
	testCases := []struct {
		name       string
		walkFn     func(string, filepath.WalkFunc) error
		readDirErr error
		expectErr  bool
	}{
		{
			name: "Test Successful Walk",
			walkFn: func(root string, walkFn filepath.WalkFunc) error {
				mockDirEntry := &mockDirEntry{
					name:  "testfile.go",
					isDir: false,
					info: &mockFileInfo{
						name:  "testfile.go",
						isDir: false,
						mode:  0444, // Set a valid file mode.
					},
				}
				return walkFn("testpath/testfile.go", mockDirEntry.info, nil)
			},
			readDirErr: nil,
			expectErr:  false,
		},
		{
			name: "Test Successful Dir Walk",
			walkFn: func(root string, walkFn filepath.WalkFunc) error {
				mockDirEntry := &mockDirEntry{
					name:  "testfile.go",
					isDir: true,
					info: &mockFileInfo{
						name:  "testfile.go",
						isDir: true,
						mode:  0444, // Set a valid file mode.
					},
				}
				return walkFn("testpath/testfile.go", mockDirEntry.info, nil)
			},
			readDirErr: nil,
			expectErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			aferoFs := afero.NewMemMapFs()

			_, err := createTestDirectoryAndFile(aferoFs)
			if err != nil {
				t.Errorf("CreatePackageDocs() error = %v", err)
			}

			repo := docs.Repo{
				Owner: "testowner",
				Name:  "testrepo",
			}

			err = docs.CreatePackageDocs(aferoFs, repo)
			if (err != nil) != tc.expectErr {
				t.Errorf("CreatePackageDocs() error = %v, wantErr %v", err, tc.expectErr)
			}
		})
	}
}
