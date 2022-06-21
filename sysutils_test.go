package utils

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCd(t *testing.T) {
	dst := ".mage"

	src := Gwd()
	if !strings.Contains(src, "goutils") {
		t.Fatal("unable to get the current working directory - Gwd() failed")
	}

	err := Cd(dst)
	if err != nil {
		t.Fatalf("failed to change directory to %s: %v - Cd() failed", dst, err)
	}

	cwd := Gwd()
	if !strings.Contains(cwd, dst) {
		t.Fatalf("failed to change directory to %s - Cd() failed", dst)
	}
}

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
		t.Fatalf("failed to copy %s to %s - Cp() failed", file, copyLoc)
	}
}

func TestEnvVarSet(t *testing.T) {
	key := "TEST_KEY"
	os.Setenv(key, "test_value")
	if err := EnvVarSet(key); err != nil {
		t.Fatalf("failed to run EnvVarSet(): %v", err)
	}

	emptykey := "EMPTY_TEST_KEY"

	if err := EnvVarSet(emptykey); err == nil {
		t.Fatalf("failed to run EnvVarSet(): %v", err)
	}
}

func TestGetHomeDir(t *testing.T) {
	_, err := GetHomeDir()
	if err != nil {
		t.Fatalf("failed to get the user's home directory - GetHomeDir() failed: %v", err)
	}
}

func TestGwd(t *testing.T) {
	out := Gwd()
	if !strings.Contains(out, "goutils") {
		t.Fatal("unable to get the current working directory - Gwd() failed")
	}
}

func isTime(obj reflect.Value) bool {
	_, ok := obj.Interface().(time.Time)
	return ok
}

func TestGetFutureTime(t *testing.T) {
	futureTime := GetFutureTime(2, 2, 3)

	ft := reflect.ValueOf(futureTime)
	if !isTime(ft) {
		t.Fatal("failed to run GetFutureTime(): incorrect value returned")
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
