package utils

import (
	"testing"
)

func TestTidy(t *testing.T) {
	if err := Tidy(); err != nil {
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
