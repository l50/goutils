package ansible

import (
	"os"

	"github.com/magefile/mage/sh"
)

// Ping runs the `ansible all -m ping` command against all nodes found in the provided hosts file by using the
// mage/sh package to execute the command. If the command execution fails, an error is returned.
//
// Parameters:
//
// hostsFile: A string representing the path to the hosts file to be used by the `ansible` command.
//
// Returns:
//
// error: An error if the `ansible` command execution fails.
//
// Example:
//
// hostsFile := "/path/to/your/hosts.ini"
// err := Ping(hostsFile)
//
//	if err != nil {
//	  log.Fatalf("failed to ping hosts: %v", err)
//	} else {
//
//	  log.Printf("Successfully pinged all hosts in %s\n", hostsFile)
//	}
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
