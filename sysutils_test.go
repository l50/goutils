package utils

import (
	"runtime"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	if runtime.GOOS == "linux" {
		out, err := RunCommand("uname", "-a")
		if !strings.Contains(out, "Linux") {
			t.Fatal("Unable to run test for RunCommand due to: ", err.Error())
		}
	} else if runtime.GOOS == "darwin" {
		out, err := RunCommand("uname", "-a")
		if !strings.Contains(out, "Darwin") {
			t.Fatal("Unable to run test for RunCommand due to: ", err.Error())
		}
	} else {
		t.Fatal("Unsupported OS detected")
	}
}
