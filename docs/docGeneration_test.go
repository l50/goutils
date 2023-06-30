package docs_test

import (
	"io"
	"io/fs"
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

type testCase struct {
	name       string
	walkFn     func(string, filepath.WalkFunc) error
	setup      func(afero.Fs) error // Setup the filesystem for each test case
	readDirErr error
	verify     func(*testing.T, afero.Fs) error // Verify the filesystem after each test case
	expectErr  bool
}

func TestCreatePackageDocs(t *testing.T) {
	testCases := []testCase{
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
		{
			name: "Test .docgenignore functionality",
			setup: func(fs afero.Fs) error {
				// Create a directory and file that should be ignored
				if err := fs.MkdirAll("ignoredir", 0755); err != nil {
					return err
				}
				file, err := fs.Create("ignoredir/README.md")
				if err != nil {
					return err
				}
				defer file.Close()
				if _, err := file.WriteString("This should not be overwritten."); err != nil {
					return err
				}

				// Create a .docgenignore file with the directory to be ignored
				if file, err = fs.Create(".docgenignore"); err != nil {
					return err
				}

				if _, err := file.WriteString("ignoredir"); err != nil {
					return err
				}

				return nil
			},
			expectErr: false,
			verify: func(t *testing.T, fs afero.Fs) error {
				// Verify that the README.md file in the ignored directory was not modified
				file, err := fs.Open("ignoredir/README.md")
				if err != nil {
					t.Errorf("error opening ignored file: %v", err)
					return err
				}
				defer file.Close()

				bytes, err := io.ReadAll(file)
				if err != nil {
					t.Errorf("error reading ignored file: %v", err)
					return err
				}

				content := string(bytes)
				if content != "This should not be overwritten." {
					t.Errorf("ignored file was modified: got %q, want %q", content, "This should not be overwritten.")
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			aferoFs := afero.NewMemMapFs()

			// Call setup function to prepare the filesystem for the test
			if tc.setup != nil {
				if err := tc.setup(aferoFs); err != nil {
					t.Fatalf("error setting up filesystem: %v", err)
				}
			}

			repo := docs.Repo{
				Owner: "testowner",
				Name:  "testrepo",
			}

			templatePath := filepath.Join("dev", "mage", "templates", "README.md.tmpl")

			err := docs.CreatePackageDocs(aferoFs, repo, templatePath)
			if (err != nil) != tc.expectErr {
				t.Errorf("CreatePackageDocs() error = %v, wantErr %v", err, tc.expectErr)
			}

			// Call the verify function if it is provided
			if tc.verify != nil {
				if err := tc.verify(t, aferoFs); err != nil {
					t.Errorf("error verifying filesystem: %v", err)
				}
			}
		})
	}
}
