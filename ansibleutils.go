package utils

import "github.com/magefile/mage/sh"

// AnsiblePing runs ansible all -m ping against all k8s nodes
// found in the input hostsFile.
// Example hostsFile input: "k3s-ansible/inventory/cowdogmoo/hosts.ini"
func AnsiblePing(hostsFile string) error {
	args := []string{
		"all",
		"-m",
		"ping",
		"-i",
		hostsFile,
	}
	if err := sh.RunV("ansible", args...); err != nil {
		return err
	}
	return nil
}
