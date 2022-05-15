package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCloneRepo(t *testing.T) {
	ogDir := Gwd()
	targetDir := "/tmp"
	repo := "https://github.com/l50/helloworld.git"
	currentTime := time.Now()
	cloneLoc := filepath.Join(
		targetDir, fmt.Sprintf(
			"helloworld-%s", currentTime.Format("2006-01-02-15-04-05"),
		),
	)

	if err := Cd(targetDir); err != nil {
		if err != nil {
			t.Fatalf("failed to cd to %s: %v - Cd() failed",
				targetDir, err)
		}
	}

	if _, err := CloneRepo(repo, cloneLoc, SSHKeyInfo{}); err != nil {
		t.Fatalf("failed to clone %s: %v", repo, err)
	}

	if err := Cd(ogDir); err != nil {
		if err != nil {
			t.Fatalf("failed to cd to %s: %v - Cd() failed",
				ogDir, err)
		}
	}

	os.RemoveAll(cloneLoc)
}
