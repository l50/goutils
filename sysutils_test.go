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
		t.Fatal("File creation failed, check CreateEmptyFile() in fileutils.go for issues")
	}
}

func TestGwd(t *testing.T) {
	out := Gwd()
	if !strings.Contains(out, "goutils") {
		t.Fatal("Unable to get the current working directory")
	}
}

func TestRunCommand(t *testing.T) {
	if runtime.GOOS == "linux" {
		out, err := RunCommand("uname", "-a")
		if !strings.Contains(out, "Linux") {
			t.Fatalf("Unable to run test for RunCommand due to: %v", err)
		}
	} else if runtime.GOOS == "darwin" {
		out, err := RunCommand("uname", "-a")
		if !strings.Contains(out, "Darwin") {
			t.Fatalf("Unable to run test for RunCommand due to: %v", err)
		}
	} else {
		t.Fatal("Unsupported OS detected")
	}
}
