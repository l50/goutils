package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallBrewDeps(t *testing.T) {
	// make sure brew is installed
	_, err := os.Stat("/usr/local/bin/brew")
	if err != nil {
		t.Skip("Skipping test, brew is not installed")
	}

	// test with valid package list
	pkgList := []string{"shellcheck", "shfmt"}
	err = InstallBrewDeps(pkgList)
	assert.NoError(t, err)

	// test with an invalid package
	pkgList = []string{"this-is-an-invalid-package"}
	err = InstallBrewDeps(pkgList)
	assert.Error(t, err)
}

func TestInstallBrewTFDeps(t *testing.T) {
	// make sure brew is installed
	if _, err := os.Stat("/usr/local/bin/brew"); err != nil {
		t.Skip("Skipping test, brew is not installed")
	}

	err = InstallBrewTFDeps()
	assert.NoError(t, err)
}
