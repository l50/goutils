package utils

import (
	"runtime"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	if runtime.GOOS == "linux" {
		cmdOut, _ := RunCommand("uname", "-a")
		if !strings.Contains(cmdOut, "Linux") {
			t.Fatalf("Unable to run test for RunCommand")
		}
	} else if runtime.GOOS == "darwin" {
		cmdOut, _ := RunCommand("uname", "-a")
		if !strings.Contains(cmdOut, "Darwin") {
			t.Fatal("Unable to run test for RunCommand")
		}
	} else {
		t.Fatal("Unsupported OS detected")
	}
}
