package utils

import (
	"fmt"
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

func TestGetHomeDir(t *testing.T) {
	_, err := GetHomeDir()
	if err != nil {
		t.Fatalf("unable to get the user's home directory - GetHomeDir() failed: %v", err)
	}
}

func TestGwd(t *testing.T) {
	out := Gwd()
	if !strings.Contains(out, "goutils") {
		t.Fatal("unable to get the current working directory - TestGwd() failed")
	}
}

func TestIsDirEmpty(t *testing.T) {
	dirEmpty, err := IsDirEmpty("/")
	if err != nil {
		t.Fatalf("failed to determine if / is empty - IsDirEmpty() failed: %v", err)
	}
	if dirEmpty != false {
		t.Fatal("the / directory has reported back as being empty, which can not be true - IsDirEmpty()")
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

func TestRunCommandWithTimeout(t *testing.T) {
	switch runtime.GOOS {
	case "linux", "darwin":
		seconds := 8
		// Test #1
		cmd := []string{"ping", "baidu.com"}
		_, _, err := RunCommandWithTimeout(seconds, cmd[0], cmd[1:]...)
		if err == nil {
			t.Fatalf("%v expected to time out - RunCommandWithTimeout() Test #1 has failed: %v",
				strings.Trim(fmt.Sprint(cmd), "[]"), err)
		}

		// Test #2
		cmd = []string{"whoami"}
		_, _, err = RunCommandWithTimeout(seconds, cmd[0], cmd[1:]...)
		if err != nil {
			t.Fatalf("%v expected to not time out - RunCommandWithTimeout() Test #2 has failed: %v",
				strings.Trim(fmt.Sprint(cmd), "[]"), err)
		}

	default:
		t.Fatal("unsupported OS detected")
	}
}
