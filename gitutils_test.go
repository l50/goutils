package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	cloneDir  string
	clonePath string
	err       error
	repo      *git.Repository
	repoURL   string
)

func init() {
	cloneDir = "/tmp"
	repoURL = "https://github.com/l50/helloworld.git"
	// Used to create a random directory name
	currentTime := time.Now()
	clonePath = filepath.Join(
		cloneDir, fmt.Sprintf(
			"helloworld-%s", currentTime.Format("2006-01-02-15-04-05"),
		),
	)

	repo, err = CloneRepo(repoURL, clonePath, nil)
	if err != nil {
		log.Fatalf(
			"failed to clone %s: %v - CloneRepo() failed",
			repo,
			err,
		)
	}
}

func TestGetTags(t *testing.T) {
	if _, err := GetTags(repo); err != nil {
		t.Fatalf("failed to get tags: %v - GetTags() failed", err)
	}
}

func TestPushTag(t *testing.T) {
	tag := "v6.6.6"

	if err := CreateTag(repo, tag); err != nil {
		t.Fatalf("failed to create %s tag: %v", tag, err)
	}

	// personal access token example
	token := "notrealtoken"
	auth := &http.BasicAuth{
		// this can be anything except for an empty string
		Username: "abc123",
		Password: token,
	}

	if err := PushTag(repo, tag, auth); err == nil {
		t.Fatal("pushing any tag should not be possible "+
			"because no auth mechanism is configured - "+
			"PushTag() failed",
			err)
	}
}

func TestGetGlobalUserCfg(t *testing.T) {
	cfg, err := GetGlobalUserCfg()
	if err != nil || cfg.User == "" {
		t.Fatalf("failed get global git user config: %v", err)
	}
}

func TestDeletePushedTag(t *testing.T) {
	tag := "v7.7.7"
	keyName := "github_rsa"

	if err := CreateTag(repo, tag); err != nil {
		t.Fatalf("failed to create %s tag: %v", tag, err)
	}

	if err := DeleteTag(repo, tag); err != nil {
		t.Fatalf("failed to delete %s tag: %v - DeleteTag() failed",
			tag, err)
	}

	pubKey, err := GetSSHPubKey(keyName, "")
	if err == nil {
		t.Fatalf("expected to fail with %s as this key should "+
			"always be encrypted at rest and we didn't provide a"+
			"decryption key: %v", keyName, err)
	}

	if err := DeletePushedTag(repo, tag, pubKey); err == nil {
		t.Fatal("deleting any tag should not be possible " +
			"in this test. There are not sufficient permissions " +
			"from the previous steps to do so - " +
			"DeletePushedTag() failed")
	}
}
