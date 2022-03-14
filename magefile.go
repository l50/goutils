// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	utils "github.com/l50/goutils"

	// mage utility functions
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Helper function to install dependencies.
func installDeps() error {
	fmt.Println(color.YellowString("Installing dependencies."))

	if err := utils.Tidy(); err != nil {
		return fmt.Errorf(color.RedString("failed to install dependencies: %w", err))
	}

	return nil
}

// InstallPreCommit installs pre-commit hooks locally
func InstallPreCommit() error {
	mg.Deps(installDeps)

	fmt.Println(color.YellowString("Installing pre-commit hooks."))
	if err := utils.InstallPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunPreCommit runs all pre-commit hooks locally
func RunPreCommit() error {
	mg.Deps(installDeps)

	fmt.Println(color.YellowString("Updating pre-commit hooks."))
	if err := utils.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Clearing the pre-commit cache to ensure we have a fresh start."))
	if err := utils.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Running all pre-commit hooks locally."))
	if err := utils.RunPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunTests runs all of the unit tests
func RunTests() error {
	mg.Deps(installDeps)

	fmt.Println(color.YellowString("Running unit tests."))
	if err := sh.RunV(filepath.Join(".hooks", "go-unit-tests.sh")); err != nil {
		return fmt.Errorf(color.RedString("failed to run unit tests: %w", err))
	}

	return nil
}

// UpdateMirror updates pkg.go.dev with the release associated with the input tag
func UpdateMirror(tag string) error {
	fmt.Println(color.YellowString("Updating pkg.go.dev with the new tag %s.", tag))
	err := sh.RunV("curl", "--silent", fmt.Sprintf("https://proxy.golang.org/github.com/l50/goutils/@v/%s.info", tag))
	if err != nil {
		return fmt.Errorf(color.RedString("failed to update pkg.go.dev: %w", err))
	}

	return nil
}

func appendToFile(file string, text string) error {
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(text + "\n"); err != nil {
		return err
	}

	return nil
}
