//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/fatih/color"

	// mage utility functions
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Helper function to install dependencies.
func installDeps() error {
	fmt.Println(color.YellowString("Installing dependencies."))

	err := sh.Run("go", "mod", "tidy")

	if err != nil {
		return fmt.Errorf(color.RedString("Failed to install dependencies: %v\n", err))
	}
	return nil
}

// InstallPreCommit installs pre-commit hooks locally
func InstallPreCommit() error {
	mg.Deps(installDeps)

	fmt.Println(color.YellowString("Installing pre-commit hooks."))
	err := sh.Run("pre-commit", "install")
	if err != nil {
		return fmt.Errorf(color.RedString("Failed to install pre-commit hooks: %v\n", err))
	}

	return nil
}

// RunPreCommit runs all pre-commit hooks locally
func RunPreCommit() error {
	mg.Deps(installDeps)

	fmt.Println(color.YellowString("Updating pre-commit hooks."))
	err := sh.RunV("pre-commit", "autoupdate")
	if err != nil {
		return fmt.Errorf(color.RedString("Failed to update the pre-commit hooks: %v\n", err))
	}

	fmt.Println(color.YellowString("Clearing the pre-commit cache to ensure we have a fresh start."))
	err = sh.RunV("pre-commit", "clean")
	if err != nil {
		return fmt.Errorf(color.RedString("Failed to clean the pre-commit cache: %v\n", err))
	}

	fmt.Println(color.YellowString("Running all pre-commit hooks locally."))
	err = sh.RunV("pre-commit", "run", "--all-files")
	if err != nil {
		return fmt.Errorf(color.RedString("Failed to run pre-commit hooks: %v\n", err))
	}

	return nil
}

// RunTests runs all of the unit tests
func RunTests() error {
	mg.Deps(installDeps)

	fmt.Println(color.YellowString("Running unit tests."))
	err := sh.RunV(".hooks/go-unit-tests.sh")
	if err != nil {
		return fmt.Errorf(color.RedString("Failed to run unit tests: %v\n", err))
	}

	return nil
}

// UpdateMirror updates pkg.go.dev with the release associated with the input tag
func UpdateMirror(tag string) error {
	fmt.Println(color.YellowString("Updating pkg.go.dev with the new tag %s.", tag))
	err := sh.RunV("curl", "--silent", fmt.Sprintf("https://proxy.golang.org/github.com/l50/goutils/@v/%s.info", tag))
	if err != nil {
		return fmt.Errorf(color.RedString("Failed to update pkg.go.dev: %v\n", err))
	}

	return nil
}
