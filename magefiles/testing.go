//go:build mage
// +build mage

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/sys"
)

type compileParams struct {
	GOOS   string
	GOARCH string
}

var repoRoot string

func init() {
	var err error
	repoRoot, err = git.RepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get repo root: %v", err)
		os.Exit(1)
	}
}

// RunTests executes all unit tests.
//
// Example usage:
//
// ```go
// mage runtests
// ```
//
// **Returns:**
//
// error: An error if any issue occurs while running the tests.
func RunTests() error {
	fmt.Println("Running unit tests.")
	if _, err := sys.RunCommand(filepath.Join(".hooks", "run-go-tests.sh"), "all"); err != nil {
		return fmt.Errorf("failed to run unit tests: %v", err)
	}

	return nil
}

// processLines parses an io.Reader, identifying and marking code blocks
// found in a TTP README.
func processLines(r io.Reader, language string) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines, codeBlockLines []string
	var inCodeBlock bool

	for scanner.Scan() {
		line := scanner.Text()

		inCodeBlock, codeBlockLines = handleLineInCodeBlock(strings.TrimSpace(line), line, inCodeBlock, language, codeBlockLines)

		if !inCodeBlock {
			lines = append(lines, codeBlockLines...)
			codeBlockLines = codeBlockLines[:0]
			if !strings.HasPrefix(line, "```") {
				lines = append(lines, line)
			}
		}
	}

	if inCodeBlock {
		codeBlockLines = append(codeBlockLines, "\t\t\t// ```")
		lines = append(lines, codeBlockLines...)
	}

	return lines, scanner.Err()
}

// handleLineInCodeBlock categorizes and handles each line based on its
// content and relation to code blocks found in a TTP README.
func handleLineInCodeBlock(trimmedLine, line string, inCodeBlock bool, language string, codeBlockLines []string) (bool, []string) {
	switch {
	case strings.HasPrefix(trimmedLine, "```"+language):
		if !inCodeBlock {
			codeBlockLines = append(codeBlockLines, line)
		}
		return !inCodeBlock, codeBlockLines
	case inCodeBlock:
		codeBlockLines = append(codeBlockLines, line)
	case strings.Contains(trimmedLine, "```"):
		inCodeBlock = false
	}
	return inCodeBlock, codeBlockLines
}
