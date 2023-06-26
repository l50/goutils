package docs

import (
	"bufio"
	"io"
	"strings"

	fileutils "github.com/l50/goutils/v2/file"
)

// FixCodeBlocks processes a provided file to ensure that all code
// blocks within comments are surrounded by markdown fenced code block
// delimiters with the specified language.
//
// **Parameters:**
// file: An object satisfying the File interface, which is to be processed.
// language: A string representing the language for the code blocks.
//
// **Returns:**
// error: An error if there's an issue reading or writing the file.
func FixCodeBlocks(file fileutils.RealFile, language string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	lines, err := processLines(rc, language)
	if err != nil {
		return err
	}

	return file.Write([]byte(strings.Join(lines, "\n")), 0644)
}

// processLines processes the provided io.Reader's lines and returns
// the processed result.
//
// **Parameters:**
// r: An io.Reader object whose lines are to be processed.
// language: A string representing the language for the code blocks.
//
// **Returns:**
// lines: A slice of strings representing the processed lines.
// error: An error if there's an issue processing the lines.
func processLines(r io.Reader, language string) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines, codeBlockLines []string
	var inCodeBlock bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		inCodeBlock, codeBlockLines = processLine(trimmedLine, line, inCodeBlock, language, codeBlockLines)
		if !inCodeBlock && len(codeBlockLines) > 0 {
			lines = append(lines, codeBlockLines...)
			codeBlockLines = codeBlockLines[:0]
		} else if !inCodeBlock {
			lines = append(lines, line)
		}
	}

	if len(codeBlockLines) > 0 {
		if inCodeBlock {
			codeBlockLines = append(codeBlockLines, "\t\t\t// ```")
		}
		lines = append(lines, codeBlockLines...)
	}

	return lines, scanner.Err()
}

// processLine processes a single line based on the current state and
// returns the new state.
//
// **Parameters:**
// trimmedLine: A string representing the trimmed version of the line.
// line: A string representing the current line.
// inCodeBlock: A boolean flag indicating if we're in a code block.
// language: A string representing the language for the code blocks.
// codeBlockLines: A slice of strings containing the processed code block lines.
//
// **Returns:**
// inCodeBlock: A boolean flag indicating if we're in a code block after processing the line.
// codeBlockLines: A slice of strings containing the processed code block lines after processing the line.
func processLine(trimmedLine, line string, inCodeBlock bool, language string, codeBlockLines []string) (bool, []string) {
	switch {
	case strings.HasPrefix(trimmedLine, "// ```"+language):
		inCodeBlock = true
		codeBlockLines = append(codeBlockLines, line)

	case strings.HasPrefix(trimmedLine, "// ```"):
		inCodeBlock = false
		codeBlockLines = append(codeBlockLines, line)

	case inCodeBlock && trimmedLine != "//" && !strings.HasPrefix(trimmedLine, "// ```"):
		codeBlockLines = append(codeBlockLines, "\t\t\t// ```"+language, line)
		inCodeBlock = true

	default:
		if inCodeBlock {
			inCodeBlock = false
			if inCodeBlock {
				codeBlockLines = append(codeBlockLines, "\t\t\t// ```")
			}
		}
	}

	return inCodeBlock, codeBlockLines
}
