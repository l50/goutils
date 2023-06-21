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

// ConfigUserInfo holds user details for the git configuration.
//
// **Attributes:**
//
// User: Global git username.
// Email: Email associated with the global git user.
type ConfigUserInfo struct {
	User  string
	Email string
}

// AddFile adds the file located at the given file path to its
// affiliated Git repository. An error is returned if an issue
// happens during this process.
//
// **Parameters:**
//
// filePath: A string indicating the path to the file that
// will be staged.
//
// **Returns:**
//
// error: An error if any occurs during the staging process.
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

// Commit generates a new commit in the specified repository with
// the given message. The commit's author is extracted from the
// global Git user settings.
//
// **Parameters:**
//
// repo: A pointer to the Repository struct symbolizing the
// repository where the commit should be made.
// msg: A string depicting the commit message.
//
// **Returns:**
//
// error: An error if the commit can't be created.
func Commit(repo *git.Repository, msg string) error {
	cfg, err := repo.Config()
	if err != nil {
		return fmt.Errorf("failed to get repo config: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to retrieve worktree: %v", err)
	}

	commit, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit current staging area: %v", err)
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("failed to run `git show`: %v", err)
	}

	if obj.Author.Email != cfg.User.Email {
		return fmt.Errorf("author email in commit doesn't match repo config email - Commit() failed: %v", err)
	}

	return nil
}

// CloneRepo clones a Git repository from the specified URL to
// the given path, using the supplied authentication method, if
// provided.
//
// **Parameters:**
//
// url: A string indicating the URL of the repository to clone.
// clonePath: A string representing the path where the repository
// will be cloned.
// auth: A transport.AuthMethod interface symbolizing the
// authentication method for cloning. If nil, no authentication is used.
//
// **Returns:**
//
// *git.Repository: A pointer to the Repository struct
// representing the cloned repository.
//
// error: An error if the repository can't be cloned or already
// exists at the target path.
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
// **Parameters:**
//
// repo: A pointer to the Repository struct representing
// the repo from which tags are retrieved.
//
// **Returns:**
//
// []string: A slice of strings, each representing a tag in the repository.
// error: An error if a problem occurs while retrieving the tags.
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

// GetGlobalUserCfg fetches the username and email from the global git user
// settings. It returns a ConfigUserInfo struct containing the global git
// username and email. An error is returned if the global git username or email
// cannot be retrieved.
//
// **Returns:**
//
// ConfigUserInfo: Struct containing the global git username and email.
//
// error: Error if the global git username or email can't be retrieved.
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

// CreateTag forms a new tag in the specified repository.
//
// **Parameters:**
//
// repo: Pointer to the Repository struct, the repository where the tag is created.
// tag: String, the name of the tag to create.
//
// **Returns:**
//
// error: Error if the tag can't be created, already exists, or if the global git
// user settings can't be retrieved.
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

// Push transmits the contents of the specified repository to the default
// remote (origin).
//
// **Parameters:**
//
// repo: Pointer to the Repository struct, the repository to push.
// auth: A transport.AuthMethod interface, the authentication method for the push.
// If it's nil, no authentication is used.
//
// **Returns:**
//
// error: Error if the push fails.
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
			fmt.Print("origin remote is up-to-date, no push was executed.")
			return nil
		}
		return fmt.Errorf(
			"error pushing to remote origin: %v", err)
	}

	return nil
}

// PushTag pushes a specific tag of the given repository to the default remote.
//
// **Parameters:**
//
// repo: Repository where the tag should be pushed.
// tag: Name of the tag to push.
// auth: Authentication method for the push. If nil, no authentication is used.
//
// **Returns:**
//
// error: Error if the push fails.
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
// **Parameters:**
//
// repo: Repository where the tag should be deleted.
// tag: The tag that should be deleted.
//
// **Returns:**
//
// error: Error if the tag cannot be deleted.
func DeleteTag(repo *git.Repository, tag string) error {
	if err := repo.DeleteTag(tag); err != nil {
		return fmt.Errorf(
			"error deleting local %s tag: %v", tag, err)
	}

	return nil
}

// DeletePushedTag deletes a tag from a repository that has been pushed.
//
// **Parameters:**
//
// repo: Repository where the tag should be deleted.
// tag: The tag that should be deleted.
// auth: Authentication method for the push.
//
// **Returns:**
//
// error: Error if the tag cannot be deleted.
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
//
// **Parameters:**
//
// dirs: Paths to directories to be searched for git repositories.
//
// **Returns:**
//
// error: Error if there's a problem with pulling the repositories.
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

// RepoRoot finds and returns the root directory of the current Git repository.
//
// **Returns:**
//
// string: Absolute path to the root directory of the current Git repository.
// error: Error if the Git repository root cannot be found.
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
