package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

var pc = sh.RunCmd("pre-commit")

// Make sure the project utilizes pre-commit,
// otherwise these utilities are not very useful to run.
func checkPCProject() error {
	var pcFile string
	cwd := Gwd()

	if strings.Contains(cwd, "magefiles") {
		pcFile = filepath.Join("..", ".pre-commit-config.yaml")
	} else {
		pcFile = ".pre-commit-config.yaml"
	}

	if !FileExists(pcFile) {
		return errors.New(color.RedString(
			"pre-commit is not configured for the current project"))
	}

	return nil
}

// InstallGoPCDeps Installs dependencies used for pre-commit with golang projects.
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
	}

	if err := InstallGoDeps(deps); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install pre-commit golang dependencies: %v", err))
	}

	return nil
}

// InstallPCHooks Installs pre-commmit hooks locally.
func InstallPCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("install"); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install pre-commit hooks: %v", err))
	}

	return nil
}

// UpdatePCHooks Updates pre-commit hooks locally.
func UpdatePCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("autoupdate"); err != nil {
		return fmt.Errorf(color.RedString("failed to update the pre-commit hooks: %v", err))
	}

	return nil
}

// ClearPCCache Clears the pre-commit cache.
func ClearPCCache() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	if err := pc("clean"); err != nil {
		return fmt.Errorf(color.RedString("failed to clear the pre-commit cache: %v", err))
	}

	return nil
}

// RunPCHooks Runs all pre-commit hooks locally.
func RunPCHooks() error {
	if err := checkPCProject(); err != nil {
		return err
	}

	cmd := "pre-commit run --show-diff-on-failure --color=always --all-files"

	if _, err := script.Exec(cmd).Stdout(); err != nil {
		return fmt.Errorf(color.RedString("failed to run pre-commit hooks: %v", err))
	}

	return nil
}

// AddFencedCB helps to address MD040 issues found with markdownlint
// by adding the input language to fenced code blocks in the input filePath.
func AddFencedCB(filePath string, language string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// Create channel to grab any errors from the anonymous function below.
	errCh1 := make(chan error)

	defer func(*os.File) {
		if err := file.Close(); err != nil {
			errCh1 <- err
		}
	}(file)

	// Create a new file to write the modified content to
	newFile, err := os.Create(filePath + ".tmp")
	if err != nil {
		return err
	}

	// Create channel to grab any errors from the anonymous function below.
	errCh2 := make(chan error)

	defer func(*os.File) {
		if err := newFile.Close(); err != nil {
			errCh2 <- err
		}
	}(newFile)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a variable to track whether the current line is a fenced code block
	inCodeBlock := false

	// Iterate through each line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if a line starts with a fenced code block
		// and that we're not already in one.
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
	err = os.Rename(filePath+".tmp", filePath)
	if err != nil {
		return err
	}

	// Check if an error was sent through the first channel
	select {
	case err := <-errCh1:
		return err
	default:
	}

	// Check if an error was sent through the second channel
	select {
	case err := <-errCh2:
		return err
	default:
	}

	return nil
}
