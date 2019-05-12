package utils

import (
	"runtime"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	if runtime.GOOS == "linux" {
		cmdOut, err := RunCommand("uname", "-a")
		if !strings.Contains(cmdOut, "Linux") {
			t.Fatal("Unable to run test for RunCommand due to: ", err.Error())
		}
	} else if runtime.GOOS == "darwin" {
		cmdOut, err := RunCommand("unamea", "-a")
		if !strings.Contains(cmdOut, "Darwin") {
			t.Fatal("Unable to run test for RunCommand due to: ", err.Error())
		}
	} else {
		t.Fatal("Unsupported OS detected")
	}
}
