//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/l50/goutils/v2/dev/lint"
	mageutils "github.com/l50/goutils/v2/dev/mage"
	"github.com/l50/goutils/v2/docs"
	fileutils "github.com/l50/goutils/v2/file/fileutils"
	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/sys"
	"github.com/magefile/mage/mg"
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

	fmt.Println("Running go mod tidy on magefiles")
	cwd := sys.Gwd()
	if err := sys.Cd("magefiles"); err != nil {
		return fmt.Errorf("failed to cd into magefiles directory: %v", err)
	}

	if err := mageutils.Tidy(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	if err := sys.Cd(cwd); err != nil {
		return fmt.Errorf("failed to cd back into repo root: %v", err)
	}

	fmt.Println("Running go mod tidy")
	if err := mageutils.Tidy(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	fmt.Println("Installing dependencies.")

	if err := lint.InstallGoPCDeps(); err != nil {
		return fmt.Errorf("failed to install pre-commit dependencies: %v", err)
	}

	if err := mageutils.InstallVSCodeModules(); err != nil {
		return fmt.Errorf("failed to install vscode-go modules: %v", err)
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

	templatePath := filepath.Join("magefiles", "tmpl", "README.md.tmpl")
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
	fmt.Println("Updating pre-commit hooks.")
	if err := lint.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println("Clearing the pre-commit cache to ensure we have a fresh start.")
	if err := lint.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println("Running all pre-commit hooks locally.")
	if err := lint.RunPCHooks(); err != nil {
		return err
	}

	return nil
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
	mg.Deps(InstallDeps)

	fmt.Println("Running unit tests.")
	if err := sh.RunV(filepath.Join(".hooks", "run-go-tests.sh"), "all"); err != nil {
		return fmt.Errorf("failed to run unit tests: %v", err)
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

// UseFixCodeBlocks fixes code blocks for the input filepath
// using the input language.
//
// Example usage:
//
// ```go
// mage fixcodeblocks docs/docGeneration.go go
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
func UseFixCodeBlocks(filepath string, language string) error {
	file := fileutils.RealFile(filepath)

	if err := docs.FixCodeBlocks(file, language); err != nil {
		return fmt.Errorf("failed to fix code blocks: %v", err)
	}

	return nil
}
