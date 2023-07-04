package aferoutils_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/file/aferoutils"
	"github.com/spf13/afero"
)

func TestTree(t *testing.T) {
	tests := []struct {
		name      string
		fs        afero.Fs
		dirPath   string
		prefix    string
		indent    string
		expectOut string
		expectErr bool
	}{
		{
			name: "Directory structure",
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()

				// Create directories
				_ = fs.MkdirAll("/file/aferoutils", 0755)
				_ = fs.MkdirAll("/file/fileutils", 0755)

				// Create files
				_ = afero.WriteFile(fs, "/file/aferoutils/README.md", []byte{}, 0644)
				_ = afero.WriteFile(fs, "/file/aferoutils/aferoutils.go", []byte{}, 0644)
				_ = afero.WriteFile(fs, "/file/aferoutils/aferoutils_examples_test.go", []byte{}, 0644)
				_ = afero.WriteFile(fs, "/file/aferoutils/aferoutils_test.go", []byte{}, 0644)

				_ = afero.WriteFile(fs, "/file/fileutils/README.md", []byte{}, 0644)
				_ = afero.WriteFile(fs, "/file/fileutils/fileutils.go", []byte{}, 0644)
				_ = afero.WriteFile(fs, "/file/fileutils/fileutils_examples_test.go", []byte{}, 0644)
				_ = afero.WriteFile(fs, "/file/fileutils/fileutils_test.go", []byte{}, 0644)

				return fs
			}(),
			dirPath: "/file",
			prefix:  "",
			indent:  "    ",
			expectOut: `
file
├── aferoutils
│   ├── README.md
│   ├── aferoutils.go
│   ├── aferoutils_examples_test.go
│   └── aferoutils_test.go
└── fileutils
    ├── README.md
    ├── fileutils.go
    ├── fileutils_examples_test.go
    └── fileutils_test.go
`,
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := aferoutils.Tree(tc.fs, tc.dirPath, tc.prefix, tc.indent, &buf)
			if (err != nil) != tc.expectErr {
				t.Errorf("unexpected error status, got: %v, expected error: %v", err, tc.expectErr)
			}
			if strings.TrimSpace(tc.expectOut) != strings.TrimSpace(buf.String()) {
				t.Errorf("unexpected output, got:\n%s\nexpected:\n%s", buf.String(), tc.expectOut)
			}
		})
	}
}
