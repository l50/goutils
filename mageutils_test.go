package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func init() {
	cloneDir = "/tmp"
	repoURL = "https://github.com/l50/helloworld.git"
	// Used to create a random directory name
	currentTime = time.Now()
	clonePath = filepath.Join(
		cloneDir, fmt.Sprintf(
			"mageutils-%s", currentTime.Format("2006-01-02-15-04-05"),
		),
	)

	// Only clone if the clone path doesn't already exist.
	if !FileExists(clonePath) {
		repo, err = CloneRepo(repoURL, clonePath, nil)
		if err != nil {
			log.Fatalf(
				"failed to clone %s - CloneRepo() failed: %v",
				repo,
				err,
			)
		}
	}
}

func TestGoReleaser(t *testing.T) {
	if err := GoReleaser(); err != nil {
		t.Fatal(err)
	}
	// Clean up
	if FileExists("dist") {
		if err := os.RemoveAll("dist"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInstallVSCodeModules(t *testing.T) {
	if err := InstallVSCodeModules(); err != nil {
		t.Fatal(err)
	}
}

func TestModUpdate(t *testing.T) {
	// First test
	recursive := false
	verbose := true
	if err := ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}

	// Second test
	recursive = true
	verbose = false
	if err := ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}
}

func TestTidy(t *testing.T) {
	if err := Tidy(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMageDeps(t *testing.T) {
	clonePath = filepath.Join(
		cloneDir, fmt.Sprintf(
			"helloworld-%s", currentTime.Format("2006-01-02-15-04-05"),
		),
	)

	// If the `clonePath` already exists, clean it up.
	if FileExists(clonePath) {
		if err := RmRf(clonePath); err != nil {
			log.Fatal(err)
		}
	}
	if err := UpdateMageDeps("magefiles"); err != nil {
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
