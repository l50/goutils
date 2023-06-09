//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/l50/goutils/dev"

	// mage utility functions
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func init() {
	os.Setenv("GO111MODULE", "on")
}

// InstallDeps Installs go dependencies
func InstallDeps() error {
	fmt.Println(color.YellowString("Installing dependencies."))

	if err := dev.Tidy(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	if err := dev.InstallGoPCDeps(); err != nil {
		return fmt.Errorf("failed to install pre-commit dependencies: %v", err)
	}

	if err := dev.InstallVSCodeModules(); err != nil {
		return fmt.Errorf("failed to install vscode-go modules: %v", err)
	}

	return nil
}

// InstallPreCommitHooks Installs pre-commit hooks locally
func InstallPreCommitHooks() error {
	mg.Deps(InstallDeps)

	fmt.Println("Installing pre-commit hooks.")
	if err := dev.InstallPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunPreCommit runs all pre-commit hooks locally
func RunPreCommit() error {
	mg.Deps(InstallDeps)

	fmt.Println("Updating pre-commit hooks.")
	if err := dev.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println("Clearing the pre-commit cache to ensure we have a fresh start.")
	if err := dev.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println("Running all pre-commit hooks locally.")
	if err := dev.RunPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunTests runs all of the unit tests
func RunTests() error {
	mg.Deps(InstallDeps)

	fmt.Println("Running unit tests.")
	if err := sh.RunV(filepath.Join(".hooks", "go-unit-tests.sh"), "all"); err != nil {
		return fmt.Errorf("failed to run unit tests: %v", err)
	}

	return nil
}

// UpdateMirror updates pkg.go.goutils with the release associated with the input tag
func UpdateMirror(tag string) error {
	var err error
	fmt.Println("Updating pkg.go.goutils with the new tag %s.", tag)

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://sum.golang.org/lookup/github.com/l50/goutils@%s",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update proxy.golang.org: %w", err)
	}

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://proxy.golang.org/github.com/l50/goutils/@v/%s.info",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update pkg.go.goutils: %w", err)
	}

	return nil
}
