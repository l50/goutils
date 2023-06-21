package lint

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mageutils "github.com/l50/goutils/v2/dev/mage"
	"github.com/l50/goutils/v2/file"
	"github.com/l50/goutils/v2/sys"
	"github.com/magefile/mage/sh"
)

var pc = sh.RunCmd("pre-commit")

func checkPCProject() error {
	cwd := sys.Gwd()
	pcFile := ".pre-commit-config.yaml"

	if strings.Contains(cwd, "magefiles") {
		pcFile = filepath.Join("..", pcFile)
	}

	if !file.Exists(pcFile) {
		return errors.New("pre-commit is not configured for the current project")
	}

	return nil
}

// InstallGoPCDeps installs dependencies used for pre-commit with Golang
// projects.
//
// **Returns:**
//
// error: An error if the dependencies fail to install.
func InstallGoPCDeps() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	deps := []string{
		"golang.org/x/lint/golint",
		"golang.org/x/tools/cmd/goimports",
		"github.com/fzipp/gocyclo/cmd/gocyclo",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"github.com/go-critic/go-critic/cmd/gocritic",
		"github.com/goreleaser/goreleaser",
	}

	if err := mageutils.InstallGoDeps(deps); err != nil {
		return fmt.Errorf("failed to install pre-commit golang dependencies: %v", err)
	}

	return nil
}

// InstallPCHooks installs pre-commit hooks locally.
//
// **Returns:**
//
// error: An error if the hooks fail to install.
func InstallPCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("install"); err != nil {
		return fmt.Errorf("failed to install pre-commit hooks: %v", err)
	}

	return nil
}

// UpdatePCHooks updates pre-commit hooks locally.
//
// **Returns:**
//
// error: An error if the hooks fail to update.
func UpdatePCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("autoupdate"); err != nil {
		return fmt.Errorf("failed to update the pre-commit hooks: %v", err)
	}

	return nil
}

// ClearPCCache clears the pre-commit cache.
//
// **Returns:**
//
// error: An error if the cache fails to clear.
func ClearPCCache() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("clean"); err != nil {
		return fmt.Errorf("failed to clear the pre-commit cache: %v", err)
	}

	return nil
}

// RunPCHooks runs pre-commit hooks with a provided timeout.
// If no timeout is provided, it defaults to 600.
//
// **Parameters:**
//
// timeout (optional): An integer specifying the timeout duration.
//
// **Returns:**
//
// error: An error if the pre-commit hook execution fails.
func RunPCHooks(timeout ...int) error {
	var timeoutValue int
	if len(timeout) > 0 {
		timeoutValue = timeout[0] // use provided value if it was provided
	} else {
		timeoutValue = 600 // default timeout value of 30 minutes
	}

	_, err := sys.RunCommandWithTimeout(timeoutValue, "pre-commit", "run", "--all-files", "--show-diff-on-failure")
	if err != nil {
		return err
	}

	return nil
}

// AddFencedCB addresses MD040 issues found with markdownlint by adding
// the input language to fenced code blocks in the input filePath.
//
// **Parameters:**
//
// filePath: Path to the markdown file to modify.
// language: Language to be added to the fenced code block.
//
// **Returns:**
//
// error: An error if the markdown file fails to be modified.
func AddFencedCB(filePath string, language string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a temporary file to write the modified content to
	tmpFilePath := filePath + ".tmp"
	newFile, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a variable to track whether the current line is a fenced code block
	inCodeBlock := false

	// Iterate through each line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if a line starts with a fenced code block and that we're not already in one.
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				line = "```" + language
				inCodeBlock = true
			} else if inCodeBlock {
				inCodeBlock = false
			}
		}

		// Write the modified line to the new file
		if _, err = newFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	// Rename the new file to the original file name
	if err = os.Rename(tmpFilePath, filePath); err != nil {
		return err
	}

	return nil
}
