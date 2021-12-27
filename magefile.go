// +build mage

package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/fatih/color"

	// mage utility functions
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Helper function to install dependencies.
func installDeps() error {
	fmt.Println(color.YellowString("Installing dependencies."))

	err := sh.Run("go", "mod", "download")

	if err != nil {
		return fmt.Errorf(color.RedString("Failed to install dependencies: %v\n", err))
	}
	return nil
}

// Install pre-commit scripts locally
func PreCommit() error {
	mg.Deps(installDeps)

	fmt.Println(color.YellowString("Installing pre-commit git hook scripts."))
	err := sh.Run("pre-commit", "install")
	if err != nil {
		return fmt.Errorf(color.RedString("Failed to install pre-commit git hook scripts: %v\n", err))
	}

	return nil
}
