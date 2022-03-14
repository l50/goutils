package utils

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// Tidy Runs go mod tidy
func Tidy() error {

	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf(color.RedString("failed to run go mod tidy: %v", err))
	}

	return nil
}

// InstallGoDeps Runs go install for the input dependencies.
func InstallGoDeps(deps []string) error {
	var err error
	failed := false

	for _, dep := range deps {
		if err := sh.RunV("go", "install", dep+"@latest"); err != nil {
			failed = true
		}
	}

	if failed {
		return fmt.Errorf(color.RedString("failed to install input go dependencies: %w", err))
	}

	return nil
}
