package git_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gitutils "github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/require"
)

var (
	currentTime time.Time
	cloneDir    = "/tmp"
	tag         = "v6.6.6"
	repo        *git.Repository
	testRepoURL                      = "https://github.com/l50/helloworld.git"
	auth        transport.AuthMethod = &http.BasicAuth{
		Username: "abc123",
		Password: "notrealtoken",
	}
)

func TestGetTags(t *testing.T) {
	if err := sys.Cd(cloneDir); err != nil {
		t.Fatalf("failed to cd to %s: %v", cloneDir, err)
	}
	if _, err := gitutils.GetTags(repo); err != nil {
		t.Fatalf("failed to get tags: %v - GetTags() failed", err)
	}
}

func TestPushAndPushTag(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(repo *git.Repository, auth transport.AuthMethod) error
		isFail bool
	}{
		{
			name:   "Push",
			fn:     gitutils.Push,
			isFail: true,
		},
		{
			name:   "PushTag",
			fn:     func(r *git.Repository, auth transport.AuthMethod) error { return gitutils.PushTag(r, tag, auth) },
			isFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn(repo, auth)
			if tc.isFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetTagsAndGlobalUserCfg(t *testing.T) {
	tests := []struct {
		name string
		fn   interface{}
	}{
		{
			name: "GetTags",
			fn:   func() ([]string, error) { return gitutils.GetTags(repo) },
		},
		{
			name: "GetGlobalUserCfg",
			fn:   gitutils.GetGlobalUserCfg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.fn.(type) {
			case func() ([]string, error):
				_, err := v()
				require.NoError(t, err)
			case func() (gitutils.ConfigUserInfo, error):
				_, err := v()
				require.NoError(t, err)
			}
		})
	}
}

func TestDeletePushedTag(t *testing.T) {
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
	tests := []struct {
		name string
	}{
		{
			name: "Update multiple repositories",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
		})
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
	testCases := []struct {
		name           string
		expectedSubstr string
	}{
		{
			name:           "Find root of temporary repo",
			expectedSubstr: "test-temp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary git repository
			tmpDir, err := createGitRepoWithCommit(t, tc.expectedSubstr, "file.txt", "test commit")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			// Change working directory to the temporary git repository
			err = os.Chdir(tmpDir)
			require.NoError(t, err)

			root, err := gitutils.RepoRoot()
			if err != nil {
				t.Fatalf("failed to retrieve root - RepoRoot() failed: %v", err)
			}

			if !strings.Contains(root, tc.expectedSubstr) {
				t.Fatalf("Expected repo root to contain the word '%s', got '%s'", tc.expectedSubstr, root)
			}
		})
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	// Create test repo and queue it for cleanup
	tag = "v6.6.6"
	var randStr string
	var err error

	randStr, err = str.GenRandom(6)
	if err != nil {
		log.Fatalf("failed to generate random string - GenRandom() failed: %v", err)
	}
	targetPath := filepath.Join(
		cloneDir, fmt.Sprintf(
			"%s-%s", fmt.Sprintf("gitutils-%s", randStr), currentTime.Format("2006-01-02-15-04-05"),
		))

	repo, err = gitutils.CloneRepo(testRepoURL, targetPath, nil)
	if err != nil {
		log.Fatalf(
			"failed to clone to %s - CloneRepo() failed: %v",
			targetPath,
			err,
		)
	}
}
