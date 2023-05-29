package macos_test

import (
	"os"
	"testing"

	"github.com/l50/goutils/v2/macos"
	"github.com/stretchr/testify/assert"
)

func TestInstallBrewDeps(t *testing.T) {
	// make sure brew is installed
	if _, err := os.Stat("/usr/local/bin/brew"); err != nil {
		t.Skip("Skipping test, brew is not installed")
	}

	tests := []struct {
		name       string
		pkgList    []string
		shouldFail bool
	}{
		{
			name:       "valid package list",
			pkgList:    []string{"shellcheck", "shfmt"},
			shouldFail: false,
		},
		{
			name:       "invalid package",
			pkgList:    []string{"this-is-an-invalid-package"},
			shouldFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := macos.InstallBrewDeps(tc.pkgList)
			if tc.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInstallBrewTFDeps(t *testing.T) {
	// make sure brew is installed
	if _, err := os.Stat("/usr/local/bin/brew"); err != nil {
		t.Skip("Skipping test, brew is not installed")
	}

	err := macos.InstallBrewTFDeps()
	assert.NoError(t, err)
}
