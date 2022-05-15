package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

var (
	errMsg string
)

// SSHKeyInfo is used to hold the name of an SSH Key
// file and a password for decryption.
type SSHKeyInfo struct {
	Name string
	PW   string
}

// GetSSHPubKey returns the public SSH key for the input
// `keyName` and uses the input `password` (if provided)
// to decrypt the associated private key.
// Resource: https://medium.com/@clm160/tag-example-with-go-git-library-4377a84bbf17
func GetSSHPubKey(keyName string, password string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys

	sshPath := os.Getenv("HOME") + "/.ssh/" + keyName
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, password)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

// CloneRepo clones the repo specified with the input `url` to
// `dstPathPath`.
//
// To clone a repo using an SSH key, provide
// the name of the key file for `sshKey.Name`.
func CloneRepo(url string, dstPath string, sshKey SSHKeyInfo) (
	*git.Repository, error) {
	var err error
	var repo *git.Repository
	var cloneOptions *git.CloneOptions

	if sshKey.Name != "" {
		pubKey, err := GetSSHPubKey(sshKey.Name, sshKey.PW)
		if err != nil {
			errMsg = fmt.Sprint(color.RedString(
				"%s failed to get %s SSH key: %v", sshKey.Name, err))
			return nil, errors.New(errMsg)
		}
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth:     pubKey,
		}
	} else {
		cloneOptions = &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		}
	}

	repo, err = git.Clone(memory.NewStorage(), nil, cloneOptions)

	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			errMsg = fmt.Sprint(color.RedString(
				"%s was already cloned to %s", url, dstPath))
			return nil, errors.New(errMsg)
		}
		errMsg = fmt.Sprint(color.RedString(
			"failed to clone %s to %s: %v", url, dstPath, err))
		return nil, errors.New(errMsg)
	}

	return repo, nil
}
