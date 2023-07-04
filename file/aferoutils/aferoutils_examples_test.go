package aferoutils_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/l50/goutils/v2/file/aferoutils"
	"github.com/spf13/afero"
)

func ExampleTree() {
	// Create an in-memory file system
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

	// Set up the output buffer
	var buf bytes.Buffer

	// Display the directory tree structure
	err := aferoutils.Tree(fs, "/file", "", "    ", &buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output
	fmt.Println(strings.TrimSpace(buf.String()))
	// Output:
	// file
	// ├── aferoutils
	// │   ├── README.md
	// │   ├── aferoutils.go
	// │   ├── aferoutils_examples_test.go
	// │   └── aferoutils_test.go
	// └── fileutils
	//     ├── README.md
	//     ├── fileutils.go
	//     ├── fileutils_examples_test.go
	//     └── fileutils_test.go
}
