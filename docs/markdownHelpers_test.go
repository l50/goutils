package docs_test

import (
	"os"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/docs"
	fileutils "github.com/l50/goutils/v2/file/fileutils"
)

func TestFixCodeBlocks(t *testing.T) {
	normalizeSpace := func(input string) string {
		lines := strings.Split(input, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimSpace(line)
		}
		return strings.Join(lines, "\n")
	}

	tests := []struct {
		name     string
		input    string
		language string
		expected string
	}{
		{
			name: "test case 1",
			input: `Driver represents an interface to Google Chrome using go.

It contains a context.Context associated with this Driver and
Options for the execution of Google Chrome.


` + "```go" + `
browser, err := cdpchrome.Init(true, true)

if err != nil {
    log.Fatalf("failed to initialize a chrome browser: %v", err)
}
` + "```",
			language: "go",
			expected: `Driver represents an interface to Google Chrome using go.

It contains a context.Context associated with this Driver and
Options for the execution of Google Chrome.


` + "```go" + `
browser, err := cdpchrome.Init(true, true)

if err != nil {
    log.Fatalf("failed to initialize a chrome browser: %v", err)
}
` + "```",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file
			tmpfile, err := os.CreateTemp("", "example.*.md")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			// Write the input to the temp file
			if _, err := tmpfile.Write([]byte(tc.input)); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}

			if err := tmpfile.Close(); err != nil {
				t.Fatalf("failed to close temp file: %v", err)
			}

			// Run the function
			err = docs.FixCodeBlocks(tc.language, fileutils.RealFile(tmpfile.Name()))

			// Check the outcome
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else {
				// Read the modified content
				content, err := os.ReadFile(tmpfile.Name())
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}

				// Check the outcome
				if normalizeSpace(string(content)) != normalizeSpace(tc.expected) {
					t.Errorf("unexpected file content:\nGot:\n%s\n\nWanted:\n%s\n",
						normalizeSpace(string(content)), normalizeSpace(tc.expected))
				}
			}

			// Clean up the temporary file
			os.Remove(tmpfile.Name())
		})
	}
}
