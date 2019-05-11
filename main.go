package main

import (
	"log"
	"os/exec"
)

// RunCommand runs a specified command
func RunCommand(cmd string, args ...string) string {

	cmdOut, err := exec.Command(cmd, args...).Output()
	if len(cmdOut) == 0 {
		if err != nil {
			log.Fatal(err)
		}
	}

	return string(cmdOut)
}
