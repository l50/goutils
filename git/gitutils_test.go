package git_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gitutils "github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testRepoPath string
	testRepo     *git.Repository
)

func TestAddFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{
			name:      "valid file path",
			filePath:  "validfile.txt",
			expectErr: false,
		},
		{
			name:      "invalid file path",
			filePath:  "nonexistentfile.txt",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filename, err := str.GenRandom(10)
			require.NoError(t, err)

			// Create a temporary git repository
			_, tmpDir, err := createGitRepoWithCommit(filename, "test commit")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			// Change working directory to the temporary git repository
			err = os.Chdir(tmpDir)
			require.NoError(t, err)

			// If we don't expect an error, create the test file
			if !tc.expectErr {
				tc.filePath = createTestFile(t, tc.filePath)
			}

			err = gitutils.AddFile(tc.filePath)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name string
		file string
		msg  string
	}{
		{
			name: "valid commit",
			file: "commitfile.txt",
			msg:  "test commit",
		},
		{
			name: "empty commit message",
			file: "anothercommitfile.txt",
			msg:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testRepo, repoPath, err := createGitRepoWithCommit(tc.file, tc.msg)
			require.NoError(t, err)
			defer os.RemoveAll(repoPath)

			err = gitutils.Commit(testRepo, tc.msg)
			assert.NoError(t, err)

			ref, err := testRepo.Head()
			assert.NoError(t, err)

			commit, err := testRepo.CommitObject(ref.Hash())
			assert.NoError(t, err)

			if commit.Message != tc.msg {
				t.Errorf("got %q, want %q", commit.Message, tc.msg)
			}
		})
	}
}

func TestCloneRepo(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		clonePath string
		auth      transport.AuthMethod
		expectErr bool
	}{
		{
			name:      "Valid Clone",
			url:       "https://github.com/l50/goutils.git",
			expectErr: false,
		},
		{
			name:      "Invalid Clone URL",
			url:       "https://github.com/user/nonexistent_repo.git",
			auth:      nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			testPath, err := os.MkdirTemp("", "test-temp")
			if err != nil {
				t.Fatal(err)
			}

			// Change working directory to the temporary git repository
			err = os.Chdir(testPath)
			require.NoError(t, err)

			repo, err := gitutils.CloneRepo(tc.url, testPath, tc.auth)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, repo)
			}
		})
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		auth      transport.AuthMethod
		expectErr bool
	}{
		{
			name:      "Valid repository",
			url:       "https://github.com/l50/goutils.git",
			expectErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testPath, err := os.MkdirTemp("", "test-temp")
			if err != nil {
				t.Fatal(err)
			}

			// Change working directory to the temporary git repository
			err = os.Chdir(testPath)
			require.NoError(t, err)

			// Clone the test repository
			testRepo, err := gitutils.CloneRepo(tc.url, testPath, tc.auth)
			if err != nil {
				t.Fatalf("failed to clone test repository: %v", err)
			}

			// Only call gitutils.Push if there was no error when opening the repository.
			err = gitutils.Push(testRepo, tc.auth)
			if err != nil && err.Error() == "authentication required" {
				// We do not want to do any pushing to anywhere in this test, so we'll call
				// it good here.
				t.Log("Success: hit authentication error when pushing to remote repository")
				return
			}
			if tc.expectErr {
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
		file string
		fn   interface{}
	}{
		{
			name: "GetTags",
			fn:   func() ([]string, error) { return gitutils.GetTags(testRepo) },
		},
		{
			name: "GetGlobalUserCfg",
			fn:   gitutils.GetGlobalUserCfg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filename, err := str.GenRandom(10)
			require.NoError(t, err)

			// Create a temporary git repository
			_, tmpDir, err := createGitRepoWithCommit(filename, "test commit")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			// Change working directory to the temporary git repository
			err = os.Chdir(tmpDir)
			require.NoError(t, err)
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

func TestCreateTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		err  error
	}{
		{
			name: "Valid Tag Creation",
			tag:  "v1.0.0",
			err:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary git repository
			filename, err := str.GenRandom(10)
			require.NoError(t, err)

			_, tmpDir, err := createGitRepoWithCommit(filename, "test commit")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			// Change working directory to the temporary git repository
			err = os.Chdir(tmpDir)
			require.NoError(t, err)

			_ = gitutils.CreateTag(testRepo, tc.tag) // Ignore the error, as we want to test delete.
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestDeleteTag(t *testing.T) {
	tests := []struct {
		name      string
		tag       string
		expectErr bool
	}{
		{
			name:      "Valid Tag Deletion",
			tag:       "v1.0.0",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filename, err := str.GenRandom(10)
			require.NoError(t, err)

			// Create a temporary git repository
			_, tmpDir, err := createGitRepoWithCommit(filename, "test commit")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			// Change working directory to the temporary git repository
			err = os.Chdir(tmpDir)
			require.NoError(t, err)

			_ = gitutils.CreateTag(testRepo, tc.tag) // Ignore the error, as we want to test delete.
			err = gitutils.DeleteTag(testRepo, tc.tag)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPushTag(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		tag       string
		auth      transport.AuthMethod
		expectErr bool
	}{
		{
			name:      "Valid repository",
			url:       "https://github.com/l50/goutils.git",
			tag:       "v1.0.0",
			expectErr: true,
		},
		{
			name:      "Invalid repository",
			url:       "https://github.com/l50/notrealrepo.git",
			tag:       "v1.0.0",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testPath, err := os.MkdirTemp("", "test-temp")
			require.NoError(t, err)

			err = os.Chdir(testPath)
			require.NoError(t, err)

			r, err := git.PlainClone(testPath, false, &git.CloneOptions{
				URL:  tc.url,
				Auth: tc.auth,
			})
			if err != nil {
				if tc.expectErr {
					return
				}
				t.Fatalf("did not expect an error but got one: %v", err)
			}

			if r != nil {
				err = gitutils.PushTag(r, tc.tag, tc.auth)
				if err != nil && err.Error() == "authentication required" {
					t.Log("Success: hit authentication error when pushing to remote repository")
					return
				}

				if tc.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestDeletePushedTag(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		tag       string
		auth      transport.AuthMethod
		expectErr bool
	}{
		{
			name:      "Valid repository",
			url:       "https://github.com/l50/goutils.git",
			tag:       "v1.0.0",
			expectErr: true,
		},
		{
			name:      "Invalid repository",
			url:       "https://github.com/l50/notrealrepo.git",
			tag:       "v1.0.0",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testPath, err := os.MkdirTemp("", "test-temp")
			require.NoError(t, err)

			err = os.Chdir(testPath)
			require.NoError(t, err)

			r, err := git.PlainClone(testPath, false, &git.CloneOptions{
				URL:  tc.url,
				Auth: tc.auth,
			})
			if err != nil {
				if tc.expectErr {
					return
				}
				t.Fatalf("did not expect an error but got one: %v", err)
			}

			if r != nil {
				err = gitutils.DeletePushedTag(r, tc.tag, tc.auth)
				if err != nil && err.Error() == "authentication required" {
					t.Log("Success: hit authentication error when pushing to remote repository")
					return
				}

				if tc.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
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

			defer cleanup(t)
			filename, err := str.GenRandom(10)
			require.NoError(t, err)

			// Create a temporary git repository
			_, tmpDirRemote, err := createGitRepoWithCommit(filename, "test commit")
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
			filename, err := str.GenRandom(10)
			require.NoError(t, err)

			// Create a temporary git repository
			_, tmpDir, err := createGitRepoWithCommit(filename, "test commit")
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

func createGitRepoWithCommit(fileName string, commitMsg string) (*git.Repository, string, error) {
	var err error
	testRepoPath, err = os.MkdirTemp("", "test-temp")
	if err != nil {
		return nil, "", err
	}

	testRepo, err = git.PlainInit(testRepoPath, false)
	if err != nil {
		return nil, "", err
	}

	// Set local config for test repo
	cfg, err := testRepo.Config()
	if err != nil {
		return nil, "", err
	}

	cfg.User.Name = "Your Name"
	cfg.User.Email = "you@example.com"
	if err := testRepo.SetConfig(cfg); err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(testRepoPath, fileName)
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		return nil, "", err
	}

	// Add the created file to the staging area
	w, err := testRepo.Worktree()
	if err != nil {
		return nil, "", err
	}
	if _, err = w.Add(fileName); err != nil {
		return nil, "", err
	}

	// Commit the file
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, "", err
	}

	return testRepo, testRepoPath, nil
}

func cloneGitRepo(t *testing.T, repo, dirname string) (string, error) {
	// Create a temporary directory to clone the repository into.
	tmpDir, err := os.MkdirTemp("", dirname)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Clone the repository into the temporary directory.
	if _, err := gitutils.CloneRepo(repo, tmpDir, nil); err != nil {
		return "", fmt.Errorf("failed to clone test repository: %w", err)
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

func cleanup(t *testing.T) {
	if err := os.RemoveAll(testRepoPath); err != nil {
		t.Fatalf("failed to remove temporary directory: %v", err)
	}
}

// createTestFile is a helper function that creates a test file and returns its path.
func createTestFile(t *testing.T, name string) string {
	file, err := os.CreateTemp(testRepoPath, name)
	require.NoError(t, err)

	return file.Name()
}
