package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallVSCodeModules(t *testing.T) {
	if err := InstallVSCodeModules(); err != nil {
		t.Fatal(err)
	}
}

func TestTidy(t *testing.T) {
	if err := Tidy(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMageDeps(t *testing.T) {
	if err := UpdateMageDeps(".mage"); err != nil {
		t.Fatal(err)
	}
}

func TestInstallGoDeps(t *testing.T) {
	sampleDeps := []string{
		"golang.org/x/lint/golint",
		"golang.org/x/tools/cmd/goimports",
	}

	if err := InstallGoDeps(sampleDeps); err != nil {
		t.Fatal(err)
	}
}

func TestCreateArtifacts(t *testing.T) {
	operatingSystems := []string{"linux", "darwin", "windows"}
	binPath := filepath.Join("../", "dist", "goutils")
	if err := CreateArtifacts(operatingSystems, binPath); err != nil {
		t.Fatal(err)
	}

	// clean up
	if err := os.RemoveAll(filepath.Dir(binPath)); err != nil {
		t.Fatal(err)
	}
}
