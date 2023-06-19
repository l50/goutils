package git_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/l50/goutils/v2/file"
	gitutils "github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	clonePath      string
	gitCleanupDirs []string
	tag            string
	currentTime    time.Time
	cloneDir       = "/tmp"
	repo           *git.Repository
	testRepoURL    = "https://github.com/l50/helloworld.git"
)

func init() {
	tag = "v6.6.6"
	// Create test repo and queue it for cleanup
	randStr, _ := str.GenRandom(8)
	repo, clonePath = createTestRepo(fmt.Sprintf("gitutils-%s", randStr))
	gitCleanupDirs = append(gitCleanupDirs, clonePath)
}

func createTestRepo(name string) (*git.Repository, string) {
	targetPath := filepath.Join(
		cloneDir, fmt.Sprintf(
			"%s-%s", name, currentTime.Format("2006-01-02-15-04-05"),
		))

	repo, err := gitutils.CloneRepo(testRepoURL, targetPath, nil)
	if err != nil {
		log.Fatalf(
			"failed to clone to %s - CloneRepo() failed: %v",
			targetPath,
			err,
		)
	}

	return repo, targetPath
}

func TestPush(t *testing.T) {
	testFile := filepath.Join(clonePath, "example-git-file")
	testFileContent := "hello world!"

	if err := file.Create(testFile, []byte(testFileContent), file.CreateFile); err != nil {
		t.Errorf("failed to create %s with %s using CreateFile(): %v", testFile, testFileContent, err)
	}

	if err := gitutils.AddFile(testFile); err != nil {
		t.Fatalf("failed to add %s: %v - AddFile() failed",
			testFile, err)
	}

	if err := gitutils.Commit(repo, testFile); err != nil {
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

	if err := gitutils.Push(repo, auth); err == nil {
		t.Fatalf("push should not be possible with "+
			"bogus credentials - Push() failed: %v", err)
	}
}

func TestGetTags(t *testing.T) {
	if _, err := gitutils.GetTags(repo); err != nil {
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

	if err := gitutils.PushTag(repo, tag, auth); err == nil {
		t.Fatal("pushing any tag should not be possible "+
			"because no auth mechanism is configured - "+
			"PushTag() failed",
			err)
	}
}

func TestGetGlobalUserCfg(t *testing.T) {
	cfg, err := gitutils.GetGlobalUserCfg()
	if err != nil || cfg.User == "" {
		t.Fatalf("failed get global git user config: %v", err)
	}
}

func TestDeletePushedTag(t *testing.T) {
	t.Cleanup(func() {
		cleanupGitUtils(t)
	})

	if err := gitutils.CreateTag(repo, tag); err != nil {
		t.Fatalf("failed to create %s tag: %v", tag, err)
	}

	keyName := "github_rsa"

	if err := gitutils.DeleteTag(repo, tag); err != nil {
		t.Fatalf("failed to delete %s tag: %v - DeleteTag() failed",
			tag, err)
	}

	pubKey, err := sys.GetSSHPubKey(keyName, "")
	if err == nil {
		fmt.Printf("security concern: %s is not encrypted at rest", keyName)
	}

	if err := gitutils.DeletePushedTag(repo, tag, pubKey); err == nil {
		t.Fatal("deleting any tag should not be possible " +
			"in this test. There are not sufficient permissions " +
			"from the previous steps to do so - " +
			"DeletePushedTag() failed")
	}

}

func TestPullRepos(t *testing.T) {
	tmpDirRemote, err := createGitRepoWithCommit(t, "test-remote", "file.txt", "test commit")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDirRemote)

	tmpDir1, err := cloneGitRepo(t, tmpDirRemote, "test1")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir1)

	tmpDir2, err := cloneGitRepo(t, tmpDirRemote, "test2")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir2)

	err = updateFileInRepo(t, tmpDirRemote, "file.txt", "test2")
	require.NoError(t, err)

	err = commitChangesInRepo(t, tmpDirRemote, "test2 commit")
	require.NoError(t, err)

	err = gitutils.PullRepos(tmpDir1, tmpDir2)
	require.NoError(t, err)

	for _, dir := range []string{tmpDir1, tmpDir2} {
		out, err := sys.RunCommand("git", "-C", dir, "status")
		require.NoError(t, err)

		require.Contains(t, out, "Your branch is up to date", "repo in %s was not updated: %s", dir, out)
	}
}

func createGitRepoWithCommit(t *testing.T, dirName string, fileName string, commitMsg string) (string, error) {
	tmpDir, err := os.MkdirTemp("", dirName)
	if err != nil {
		return "", err
	}

	_, err = sys.RunCommand("git", "-C", tmpDir, "init")
	require.NoError(t, err)

	filePath := filepath.Join(tmpDir, fileName)
	err = os.WriteFile(filePath, []byte("test"), 0644)
	require.NoError(t, err)

	_, err = sys.RunCommand("git", "-C", tmpDir, "add", filePath)
	require.NoError(t, err)

	_, err = sys.RunCommand("git", "-C", tmpDir, "commit", "-m", commitMsg)
	require.NoError(t, err)

	return tmpDir, nil
}

func cloneGitRepo(t *testing.T, sourceDir string, destDirName string) (string, error) {
	tmpDir, err := os.MkdirTemp("", destDirName)
	if err != nil {
		return "", err
	}

	_, err = sys.RunCommand("git", "clone", sourceDir, tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}

func updateFileInRepo(t *testing.T, repoDir string, fileName string, content string) error {
	filePath := filepath.Join(repoDir, fileName)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	if _, err := sys.RunCommand("git", "-C", repoDir, "add", filePath); err != nil {
		return err
	}

	return nil
}

func commitChangesInRepo(t *testing.T, repoDir string, commitMsg string) error {
	if _, err := sys.RunCommand("git", "-C", repoDir, "commit", "-m", commitMsg); err != nil {
		return err
	}

	return nil
}

func TestRepoRoot(t *testing.T) {
	root, err := gitutils.RepoRoot()
	if err != nil {
		t.Fatalf("failed to retrieve root - RepoRoot() failed: %v", err)
	}

	assert.Contains(t, root, "goutils", "Expected repo root to contain the word 'goutils'")
}

func cleanupGitUtils(t *testing.T) {
	for _, dir := range gitCleanupDirs {
		if err := sys.RmRf(dir); err != nil {
			fmt.Println("failed to clean up gitUtils: ", err.Error())
		}
	}
}
