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

// Create a new tag. The tag must be in v1.x.x format.
func Tag(tag string) (err error) {
	var releaseTag = regexp.MustCompile(`^v1\.[0-9]+\.[0-9]+$`)

	if !releaseTag.MatchString(tag) {
		return errors.New("TAG environment variable must be in semver v1.x.x format. Input tag: " + tag)
	}

	fmt.Printf(color.YellowString("Creating new tag %s.", tag))
	if err := sh.RunV("git", "tag", "-a", tag, "-m", tag); err != nil {
		return err
	}

	fmt.Printf(color.YellowString("Pushing new tag %s.", tag))
	if err := sh.RunV("git", "push", "origin", tag); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			fmt.Printf(color.RedString("Failed to create new tag %s! Cleaning up.", tag))
			sh.RunV("git", "tag", "--delete", tag)
			sh.RunV("git", "push", "--delete", "origin", tag)
		}
	}()

	return nil
}
