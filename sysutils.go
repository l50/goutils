package utils

import (
	"io/ioutil"
	"log"
	"os/exec"
)

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

// RunCommand runs a specified system command
func RunCommand(cmd string, args ...string) (string, error) {

	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(out), nil
}
