package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// CheckRoot will check to see if the process is being run as root
func CheckRoot() {
	if os.Geteuid() != 0 {
		log.Fatalln("This script must be run as root.")
	}
}

// Cp is used to copy a file from a src to a destination
func Cp(src string, dst string) bool {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		log.Printf("Error reading %s:\n", src)
		log.Println(err)
		return false
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		log.Printf("Error creating %s:\n", dst)
		log.Println(err)
		return false
	}
	return true
}

// Gwd will return the current working directory
func Gwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	return dir
}

// RunCommand runs a specified system command
func RunCommand(cmd string, args ...string) (string, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("%s %s %s %s", cmd, args, out, err)
	}
	return string(out), nil
}
