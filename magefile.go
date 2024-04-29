//go:build mage
// +build mage

package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/l50/goutils/v2/dev/lint"
	mageutils "github.com/l50/goutils/v2/dev/mage"
	"github.com/l50/goutils/v2/docs"
	fileutils "github.com/l50/goutils/v2/file/fileutils"
	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/logging"
	"github.com/l50/goutils/v2/sys"
	"github.com/magefile/mage/sh"
	"github.com/spf13/afero"
)

func init() {
	os.Setenv("GO111MODULE", "on")
}

// InstallDeps installs the Go dependencies necessary for developing
// on the project.
//
// Example usage:
//
// ```go
// mage installdeps
// ```
//
// **Returns:**
//
// error: An error if any issue occurs while trying to
// install the dependencies.
func InstallDeps() error {
	fmt.Println(color.YellowString("Running go mod tidy on repo root."))
	if err := mageutils.Tidy(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	fmt.Println(color.YellowString("Installing dependencies."))
	if err := lint.InstallGoPCDeps(); err != nil {
		return fmt.Errorf("failed to install pre-commit dependencies: %v", err)
	}

	if err := mageutils.InstallVSCodeModules(); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install vscode-go modules: %v", err))
	}

	return nil
}

// GeneratePackageDocs creates documentation for the various packages
// in the project.
//
// Example usage:
//
// ```go
// mage generatepackagedocs
// ```
//
// **Returns:**
//
// error: An error if any issue occurs during documentation generation.
func GeneratePackageDocs() error {
	fs := afero.NewOsFs()

	repoRoot, err := git.RepoRoot()
	if err != nil {
		return fmt.Errorf("failed to get repo root: %v", err)
	}
	sys.Cd(repoRoot)

	repo := docs.Repo{
		Owner: "l50",
		Name:  "goutils/v2",
	}

	templatePath := filepath.Join("templates", "README.md.tmpl")
	if err := docs.CreatePackageDocs(fs, repo, templatePath); err != nil {
		return fmt.Errorf("failed to create package docs: %v", err)
	}

	return nil
}

// RunPreCommit updates, clears, and executes all pre-commit hooks
// locally. The function follows a three-step process:
//
// First, it updates the pre-commit hooks.
// Next, it clears the pre-commit cache to ensure a clean environment.
// Lastly, it executes all pre-commit hooks locally.
//
// Example usage:
//
// ```go
// mage runprecommit
// ```
//
// **Returns:**
//
// error: An error if any issue occurs at any of the three stages
// of the process.
func RunPreCommit() error {
	if !sys.CmdExists("pre-commit") {
		return fmt.Errorf("pre-commit is not installed")
	}

	fmt.Println(color.YellowString("Updating pre-commit hooks."))
	if err := lint.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Clearing the pre-commit cache to ensure we have a fresh start."))
	if err := lint.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Running all pre-commit hooks locally."))
	if err := lint.RunPCHooks(); err != nil {
		return err
	}

	return nil
}

// UpdateMirror updates pkg.go.dev with the release associated with the
// input tag
//
// Example usage:
//
// ```go
// mage updatemirror v2.0.1
// ```
//
// **Parameters:**
//
// tag: the tag to update pkg.go.dev with
//
// **Returns:**
//
// error: An error if any issue occurs while updating pkg.go.dev
func UpdateMirror(tag string) error {
	var err error
	fmt.Printf("Updating pkg.go.dev with the new tag %s.", tag)

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://sum.golang.org/lookup/github.com/l50/goutils/v2@%s",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update proxy.golang.org: %w", err)
	}

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://proxy.golang.org/github.com/l50/goutils/v2/@v/%s.info",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update pkg.go.dev: %w", err)
	}

	return nil
}

// FixCodeBlocks fixes code blocks for the input filepath
// using the input language.
//
// Example usage:
//
// ```go
// mage fixcodeblocks go docs/docGeneration.go
// ```
//
// **Parameters:**
//
// filepath: the path to the file or directory to fix
//
// language: the language of the code blocks to fix
//
// **Returns:**
//
// error: an error if one occurred
func FixCodeBlocks(language string, filepath string) error {
	file := fileutils.RealFile(filepath)

	if err := docs.FixCodeBlocks(language, file); err != nil {
		return fmt.Errorf("failed to fix code blocks: %v", err)
	}

	return nil
}

// TestLoggerOutput tests the output of the logger
func TestLoggerOutput() {
	// Logger test
	cfg := logging.LogConfig{
		Fs:         afero.NewOsFs(),
		Level:      slog.LevelDebug,
		OutputType: logging.ColorOutput,
		LogToDisk:  true,
		LogPath:    "./test.log",
	}

	log, err := logging.InitLogging(&cfg)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	log.Println("This is a test info message")
	log.Printf("This is a test %s info message", "formatted")
	log.Error("This is a test error message")
	log.Debugf("This is a test debug message")
	log.Errorf("This is a test %s error message", "formatted")
	log.Println("{\"time\":\"2024-01-03T23:12:35.937476-07:00\",\"level\":\"ERROR\",\"msg\":\"\\u001b[1;32m==> docker.ansible-attack-box: Starting docker container...\\u001b[0m\"}")
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
// found in a README.
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
// content and relation to code blocks found in a README.
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
