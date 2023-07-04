package docs_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/l50/goutils/v2/docs"
	fileutils "github.com/l50/goutils/v2/file/fileutils"
)

func ExampleFixCodeBlocks() {
	input := `Driver represents an interface to Google Chrome using go.

It contains a context.Context associated with this Driver and
Options for the execution of Google Chrome.

` + "```go" + `
browser, err := cdpchrome.Init(true, true)

if err != nil {
    log.Fatalf("failed to initialize a chrome browser: %v", err)
}
` + "```"
	language := "go"

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example.*.md")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write the input to the temp file
	if _, err := tmpfile.Write([]byte(input)); err != nil {
		log.Printf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpfile.Close(); err != nil {
		log.Printf("failed to close temp file: %v", err)
		return
	}

	// Run the function
	file := fileutils.RealFile(tmpfile.Name())
	err = docs.FixCodeBlocks(file, language)
	if err != nil {
		log.Printf("failed to fix code blocks: %v", err)
		return
	}

	// Read the modified content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		log.Printf("failed to read file: %v", err)
		return
	}

	// Print the result
	fmt.Println(strings.TrimSpace(string(content)))
	// Output:
	// Driver represents an interface to Google Chrome using go.
	//
	// It contains a context.Context associated with this Driver and
	// Options for the execution of Google Chrome.
	//
	// ```go
	// browser, err := cdpchrome.Init(true, true)
	//
	// if err != nil {
	//     log.Fatalf("failed to initialize a chrome browser: %v", err)
	// }
	// ```
}
