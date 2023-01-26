package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
)

var (
	clonePath      string
	gitCleanupDirs []string
	tag            string
)

func init() {
	tag = "v6.6.6"
	// Create test repo and queue it for cleanup
	randStr, _ := RandomString(8)
	clonePath = createTestRepo(fmt.Sprintf("gitutils-%s", randStr))
	gitCleanupDirs = append(gitCleanupDirs, clonePath)
}

func createTestRepo(name string) string {
	targetPath := filepath.Join(
		cloneDir, fmt.Sprintf(
			"%s-%s", name, currentTime.Format("2006-01-02-15-04-05"),
		),
	)

	repo, err = CloneRepo(testRepoURL, targetPath, nil)
	if err != nil {
		log.Fatalf(
			"failed to clone to %s - CloneRepo() failed: %v",
			targetPath,
			err,
		)
	}

	return targetPath
}

func TestPush(t *testing.T) {
	testFile := filepath.Join(clonePath, "example-git-file")
	testFileContent := "hello world!"

	if err := CreateFile(testFile, []byte(testFileContent)); err != nil {
		t.Errorf("failed to create %s with %s using CreateFile(): %v", testFile, testFileContent, err)
	}

	if err := AddFile(testFile); err != nil {
		t.Fatalf("failed to add %s: %v - AddFile() failed",
			testFile, err)
	}

	if err := Commit(repo, testFile); err != nil {
		t.Fatalf("failed to commit staged files in %s: %v",
			testFile, err)
	}

	// personal access token example
	token := "notrealtoken"
	auth := &http.BasicAuth{
		// this can be anything except for an empty string
		Username: "abc123",
		Password: token,
	}

	if err := Push(repo, auth); err == nil {
		t.Fatalf("push should not be possible with "+
			"bogus credentials - Push() failed: %v", err)
	}
}

func TestGetTags(t *testing.T) {
	if _, err := GetTags(repo); err != nil {
		t.Fatalf("failed to get tags: %v - GetTags() failed", err)
	}
}

func TestPushTag(t *testing.T) {
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
	t.Cleanup(func() {
		cleanupGitUtils(t)
	})

	if err := CreateTag(repo, tag); err != nil {
		t.Fatalf("failed to create %s tag: %v", tag, err)
	}

	keyName := "github_rsa"

	if err := DeleteTag(repo, tag); err != nil {
		t.Fatalf("failed to delete %s tag: %v - DeleteTag() failed",
			tag, err)
	}

	pubKey, err := GetSSHPubKey(keyName, "")
	if err == nil {
		fmt.Print(color.RedString(
			"security concern: %s is not encrypted at rest",
			keyName))
	}

	if err := DeletePushedTag(repo, tag, pubKey); err == nil {
		t.Fatal("deleting any tag should not be possible " +
			"in this test. There are not sufficient permissions " +
			"from the previous steps to do so - " +
			"DeletePushedTag() failed")
	}

}

func TestRepoRoot(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("failed to retrieve root - RepoRoot() failed: %v", err)
	}

	assert.Contains(t, root, "goutils", "Expected repo root to contain the word 'goutils'")
}

func cleanupGitUtils(t *testing.T) {
	for _, dir := range gitCleanupDirs {
		if err := RmRf(dir); err != nil {
			fmt.Println("failed to clean up gitUtils: ", err.Error())
		}
	}
}
