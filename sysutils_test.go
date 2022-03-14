package utils

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestCp(t *testing.T) {
	file := "test.txt"
	copyLoc := "testing.txt"
	created := CreateEmptyFile(file)
	if created {
		copied := Cp(file, copyLoc)
		if copied {
			if FileExists(copyLoc) {
				os.Remove(file)
				os.Remove(copyLoc)
			}
		}
	} else {
		t.Fatal("file creation failed, check CreateEmptyFile() in fileutils.go for issues - TestCp() failed")
	}
}

func TestGwd(t *testing.T) {
	out := Gwd()
	if !strings.Contains(out, "goutils") {
		t.Fatal("unable to get the current working directory - TestGwd() failed")
	}
}

func TestRunCommand(t *testing.T) {
	switch runtime.GOOS {
	case "linux", "darwin":
		out, err := RunCommand("uname", "-a")
		if !strings.Contains(out, "Linux") && !strings.Contains(out, "Darwin") {
			t.Fatalf("unable to run command - RunCommand() failed: %v", err)
		}
	default:
		t.Fatal("unsupported OS detected")
	}
}
