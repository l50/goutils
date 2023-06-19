package git

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/l50/goutils/v2/sys"

	"github.com/magefile/mage/sh"
)

// ConfigUserInfo holds a username and
// email to use for user.name and user.email.
type ConfigUserInfo struct {
	User  string
	Email string
}

// AddFile stages the file at the given file path in its associated Git repository. Returns an error if one occurs.
//
// Parameters:
//
// filePath: A string representing the path to the file to be staged.
//
// Returns:
//
// error: An error if one occurs during the staging process.
//
// Example:
//
// filePath := "/path/to/your/file"
// err := AddFile(filePath)
//
//	if err != nil {
//	  log.Fatalf("failed to stage file: %v", err)
//	}
//
// log.Printf("Staged file: %s", filePath)
func AddFile(filePath string) error {
	repo, err := git.PlainOpen(filepath.Dir(filePath))
	if err != nil {
		return fmt.Errorf(
			"failed to open %s repo: %v", repo, err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve worktree: %v", err)
	}

	_, err = w.Add(filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf(
			"failed to run `git add` on %s: %v", filePath, err)
	}

	status, err := w.Status()
	if err != nil {
		return fmt.Errorf(
			"failed to run `git status` on %s: %v", filePath, err)
	}

	if status.IsClean() {
		return fmt.Errorf(
			"status is clean - failed to run `git add` "+
				"on %s: %v", filePath, err)
	}

	return nil
}

// Commit creates a new commit in the given repository with the provided message. The author of the commit is retrieved from the global Git user settings.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the commit should be created.
//
// msg: A string representing the commit message.
//
// Returns:
//
// error: An error if the commit cannot be created.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// msg := "Commit message"
// err = Commit(repo, msg)
//
//	if err != nil {
//	  log.Fatalf("failed to create commit: %v", err)
//	}
func Commit(repo *git.Repository, msg string) error {
	cfg, err := GetGlobalUserCfg()
	if err != nil {
		return fmt.Errorf(
			"failed get repo config: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve worktree: %v", err)
	}

	commit, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User,
			Email: cfg.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf(
			"failed to commit current staging area`: %v",
			err)
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf(
			"failed to run `git show`: %v", err)
	}

	if obj.Author.Email != cfg.Email {
		return fmt.Errorf(
			"author email in commit doesn't match "+
				"global git config email - Commit() failed: %v",
			err)
	}

	return nil
}

// CloneRepo clones a Git repository from the specified URL to the target path, using the provided authentication method if it is not nil.
//
// Parameters:
//
// url: A string representing the URL of the repository to clone.
//
// clonePath: A string representing the path where the repository should be cloned to.
//
// auth: A transport.AuthMethod interface representing the authentication method to use for cloning. If it's nil, no authentication is used.
//
// Returns:
//
// *git.Repository: A pointer to the Repository struct representing the cloned repository.
//
// error: An error if the repository cannot be cloned or already exists at the target path.
//
// Example:
//
// url := "https://github.com/user/repo.git"
// clonePath := "/path/to/repo"
//
//	auth := &http.BasicAuth{
//	    Username: "your_username",
//	    Password: "your_password",
//	}
//
// repo, err := CloneRepo(url, clonePath, auth)
//
//	if err != nil {
//	  log.Fatalf("failed to clone repository: %v", err)
//	}
func CloneRepo(url string, clonePath string, auth transport.AuthMethod) (
	*git.Repository, error) {
	var err error
	var repo *git.Repository
	var cloneOptions *git.CloneOptions

	if auth != nil {
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth:     auth,
		}
	} else {
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		}
	}

	repo, err = git.PlainClone(clonePath, false, cloneOptions)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			return nil, fmt.Errorf(
				"%s was already cloned to %s", url, clonePath)
		}
		return nil, fmt.Errorf(
			"failed to clone %s to %s: %v", url, clonePath, err)
	}

	return repo, nil
}

// GetTags returns all tags of the given repository.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository from which tags are retrieved.
//
// Returns:
//
// []string: A slice of strings, each representing a tag in the repository.
//
// error: An error if the tags cannot be retrieved.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// tags, err := GetTags(repo)
//
//	if err != nil {
//	  log.Fatalf("failed to get tags: %v", err)
//	}
func GetTags(repo *git.Repository) ([]string, error) {
	var tags []string
	tagObjects, err := repo.TagObjects()
	if err != nil {
		return []string{}, fmt.Errorf(
			"failed to retrieve repo tags: %v", err)
	}

	err = tagObjects.ForEach(func(t *object.Tag) error {
		tags = append(tags, t.Name)
		return nil
	})

	if err != nil {
		return tags, fmt.Errorf(
			"failed to retrieve repo tags: %v", err)
	}

	return tags, err
}

func tagExists(repo *git.Repository, tag string) (bool, error) {
	tags, err := GetTags(repo)
	if err != nil {
		return false, err
	}

	for _, t := range tags {
		if t == tag {
			return true, nil
		}
	}

	return false, nil
}

// GetGlobalUserCfg retrieves the username and email from the global git user settings.
//
// Returns:
//
// ConfigUserInfo: A ConfigUserInfo struct containing the global git username and email.
//
// error: An error if the global git username or email cannot be retrieved.
//
// Example:
//
// userInfo, err := GetGlobalUserCfg()
//
//	if err != nil {
//	  log.Fatalf("failed to retrieve global git user settings: %v", err)
//	}
func GetGlobalUserCfg() (ConfigUserInfo, error) {
	userInfo := ConfigUserInfo{}
	var err error

	userInfo.User, err = sh.Output("git", "config", "user.name")
	if err != nil {
		return userInfo, fmt.Errorf(
			"failed to retrieve global git username: %v", err)
	}

	userInfo.Email, err = sh.Output("git", "config", "user.email")
	if err != nil {
		return userInfo, fmt.Errorf(
			"failed to retrieve global git email: %v", err)
	}

	return userInfo, nil
}

// CreateTag creates a new tag in the given repository.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the tag should be created.
//
// tag: A string representing the name of the tag to create.
//
// Returns:
//
// error: An error if the tag cannot be created, already exists, or if the global git user settings cannot be retrieved.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// tag := "v1.0.0"
// err = CreateTag(repo, tag)
//
//	if err != nil {
//	  log.Fatalf("failed to create tag: %v", err)
//	}
func CreateTag(repo *git.Repository, tag string) error {
	exists, err := tagExists(repo, tag)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve repo tags: %v", err)
	}

	if exists {
		return fmt.Errorf(
			"error creating input tag %s: it already exists", tag)
	}

	cfg, err := GetGlobalUserCfg()
	if err != nil {
		return fmt.Errorf(
			"failed get repo config: %v", err)
	}

	h, err := repo.Head()
	if err != nil {
		return fmt.Errorf(
			"failed to get repo head: %v", err)
	}

	_, err = repo.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  cfg.User,
			Email: cfg.Email,
			When:  time.Now(),
		},
		Message: tag,
	})

	if err != nil {
		return fmt.Errorf(
			"error creating input tag %s: %v", tag, err)
	}

	return nil
}

// Push pushes the contents of the given repository to the default remote (origin).
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository to push.
//
// auth: A transport.AuthMethod interface representing the authentication method to use for the push. If it's nil, no authentication is used.
//
// Returns:
//
// error: An error if the push fails.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
//	auth := &http.BasicAuth{
//	    Username: "your_username",
//	    Password: "your_password",
//	}
//
// err = Push(repo, auth)
//
//	if err != nil {
//	  log.Fatalf("failed to push to remote: %v", err)
//	}
func Push(repo *git.Repository, auth transport.AuthMethod) error {
	var pushOptions *git.PushOptions

	if auth != nil {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			Auth:       auth,
		}
	} else {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
		}
	}

	err := repo.Push(pushOptions)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print(color.YellowString(
				"origin remote is up-to-date, no push was executed."))
			return nil
		}
		return fmt.Errorf(
			"error pushing to remote origin: %v", err)
	}

	return nil
}

// PushTag pushes a specific tag of the given repository to the default remote (origin).
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the tag should be pushed.
//
// tag: A string representing the name of the tag to push.
//
// auth: A transport.AuthMethod interface representing the authentication method to use for the push. If it's nil, no authentication is used.
//
// Returns:
//
// error: An error if the push fails.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
func PushTag(repo *git.Repository, tag string, auth transport.AuthMethod) error {
	var pushOptions *git.PushOptions

	if auth != nil {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs: []config.RefSpec{config.RefSpec(
				"refs/tags/*:refs/tags/*")},
			Auth: auth,
		}
	} else {
		pushOptions = &git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs: []config.RefSpec{config.RefSpec(
				"refs/tags/*:refs/tags/*")},
		}
	}

	err := repo.Push(pushOptions)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print(color.YellowString(
				"origin remote is up-to-date, no push was executed."))
			return nil
		}

		return fmt.Errorf(
			"error pushing %s tag to remote origin: %v", tag, err)
	}

	return nil
}

// DeleteTag deletes the local input tag from the specified repo.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the tag should be deleted.
//
// tag: A string representing the tag that should be deleted.
//
// auth: An AuthMethod representing the method used to authenticate with the remote repository.
//
// Returns:
//
// error: An error if the tag cannot be deleted.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// tag := "v1.0.0"
// auth, err := ssh.NewSSHAgentAuth("git")
//
//	if err != nil {
//	  log.Fatalf("failed to create authentication method: %v", err)
//	}
//
// err = DeletePushedTag(repo, tag, auth)
//
//	if err != nil {
//	  log.Fatalf("failed to delete pushed tag: %v", err)
//	}
func DeleteTag(repo *git.Repository, tag string) error {
	if err := repo.DeleteTag(tag); err != nil {
		return fmt.Errorf(
			"error deleting local %s tag: %v", tag, err)
	}

	return nil
}

// DeletePushedTag deletes a tag from a given repository that has been pushed remote.
// The tag is deleted from both the local repository and the remote repository.
//
// Parameters:
//
// repo: A pointer to the Repository struct representing the repository where the tag should be deleted.
//
// tag: A string representing the tag that should be deleted.
//
// auth: An AuthMethod representing the method used to authenticate with the remote repository.
//
// Returns:
//
// error: An error if the tag cannot be deleted.
//
// Example:
//
// repo, err := git.PlainOpen("/path/to/repo")
//
//	if err != nil {
//	  log.Fatalf("failed to open repository: %v", err)
//	}
//
// tag := "v1.0.0"
// auth, err := ssh.NewSSHAgentAuth("git")
//
//	if err != nil {
//	  log.Fatalf("failed to create authentication method: %v", err)
//	}
//
// err = DeletePushedTag(repo, tag, auth)
//
//	if err != nil {
//	  log.Fatalf("failed to delete pushed tag: %v", err)
//	}
func DeletePushedTag(repo *git.Repository, tag string, auth transport.AuthMethod) error {
	err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs: []config.RefSpec{config.RefSpec(
			"refs/tags/" + tag)},
		Auth: auth,
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print("origin remote is up-to-date, no delete was executed.")
			return nil
		}

		return fmt.Errorf("error deleting pushed tag %s: %v", tag, err)
	}

	return nil
}

// PullRepos updates all git repositories located in the specified directories.
// It traverses each directory recursively, identifies directories that contain a ".git" subdirectory,
// and runs a "git pull" command for the current branch. If no branch is checked out, it tries to pull the default branch.
// Parameters:
//
// dirs: Strings representing paths to directories that should be searched for git repositories.
//
// Returns:
//
// error: An error if the current working directory cannot be obtained, a directory cannot be traversed,
//
//	a directory cannot be changed, or the current or default branch cannot be obtained.
//
// Example:
//
// dirs := []string{"/path/to/your/directory", "/another/path/to/your/directory"}
// err := sys.PullRepos(dirs...)
//
//	if err != nil {
//	  log.Fatalf("failed to pull repos: %v", err)
//	}
func PullRepos(dirs ...string) error {
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".git" {
				return updateRepo(filepath.Dir(path))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to find git repositories in %s: %w", dir, err)
		}
	}
	return nil
}

func updateRepo(repoDir string) error {
	// Change to the repository directory
	if err := os.Chdir(repoDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", repoDir, err)
	}

	// Ensure that we switch back after the function returns
	defer func() {
		if err := os.Chdir(".."); err != nil {
			log.Printf("failed to change directory: %v\n", err)
		}
	}()

	// Get the current branch
	refOutput, err := sys.RunCommand("git", "symbolic-ref", "--short", "HEAD")
	if err != nil {
		// If no branch is checked out, get the default branch
		defaultBranchOutput, defaultBranchErr := sys.RunCommand("git", "config", "--get", "init.defaultBranch")
		if defaultBranchErr != nil {
			return fmt.Errorf("failed to get current or default branch for %s: %v, %v", repoDir, err, defaultBranchErr)
		}
		refOutput = defaultBranchOutput
	}

	// Pull changes in the current branch
	ref := strings.TrimSpace(refOutput)
	res, err := sys.RunCommand("git", "pull", "origin", ref)
	if err != nil {
		fmt.Printf("failed to update %s: %s\n", repoDir, res)
	} else if strings.TrimSpace(res) != "Already up to date." {
		fmt.Printf("Now pulling the latest from upstream for %s\n", repoDir)
	}

	return nil
}

// RepoRoot finds and returns the absolute path of the root directory of the current Git repository.
// It does this by moving upwards through the directory hierarchy until it finds a directory containing a ".git" subdirectory.
// If it reaches the filesystem root without finding a Git repository, an error is returned.
//
// Parameters:
//
// # None
//
// Returns:
//
// string: A string representing the absolute path to the root directory of the current Git repository.
// error: An error if the current working directory cannot be obtained, or a Git repository root cannot be found.
//
// Example:
//
// root, err := gitutils.RepoRoot()
//
//	if err != nil {
//	    log.Fatalf("failed to retrieve root: %v", err)
//	}
//
// fmt.Printf("The root of the current Git repository is: %s\n", root)
func RepoRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
			return cwd, nil
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "", fmt.Errorf("git root not found")
		}
		cwd = parent
	}
}
