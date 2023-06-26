package ansible

import (
	"os"

	"github.com/magefile/mage/sh"
)

// Ping runs the `ansible all -m ping` command against
// all nodes found in the provided hosts file by using the
// mage/sh package to execute the command. If the command
// execution fails, an error is returned.
//
// **Parameters:**
//
// hostsFile: A string representing the path to the hosts
// file to be used by the `ansible` command.
//
// **Returns:**
//
// error: An error if the `ansible` command execution fails.
func Ping(hostsFile string) error {
	args := []string{
		"all",
		"-m",
		"ping",
	}

	// Check if the hosts file exists
	if _, err := os.Stat(hostsFile); os.IsNotExist(err) {
		args = append(args, "-l", "localhost")
	} else {
		args = append(args, "-i", hostsFile)
	}

	if err := sh.RunV("ansible", args...); err != nil {
		return err
	}

	return nil
}
