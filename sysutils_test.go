package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bitfield/script"
)

var (
	debug = false
)

func TestCd(t *testing.T) {
	dst := "magefiles"

	src := Gwd()
	if !strings.Contains(src, "goutils") {
		t.Fatal("unable to get the current working directory - Gwd() failed")
	}

	if err := Cd(dst); err != nil {
		t.Fatalf("failed to change directory to %s: %v - Cd() failed", dst, err)
	}

	cwd := Gwd()
	if !strings.Contains(cwd, dst) {
		t.Fatalf("failed to change directory to %s - Cd() failed", dst)
	}
}

func TestCmdExists(t *testing.T) {
	fail := "asdf"
	cmds := []string{"ls", "whoami", fail}
	for _, cmd := range cmds {
		if CmdExists(cmd) && cmd == fail {
			t.Fatalf(
				"failed to properly identify installed cmd: %v - CmdExists() failed", err)
		}
	}
}

func TestCp(t *testing.T) {
	file := "test.txt"
	copyLoc := "testing.txt"
	created := CreateEmptyFile(file)
	if created {
		if err := Cp(file, copyLoc); err != nil {
			t.Fatalf("failed to copy %s to %s - Cp() failed", file, copyLoc)
		}
		if FileExists(copyLoc) {
			remove := []string{file, copyLoc}
			for _, f := range remove {
				if err := DeleteFile(f); err != nil {
					t.Errorf("unable to delete %s, DeleteFile() failed", f)
				}
			}
		}
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
	downloadURL := "https://raw.githubusercontent.com/rebootuser/LinEnum/master/LinEnum.sh"
	targetPath := filepath.Join("/tmp", "Linenum.sh")
	dlFilePath, err := DownloadFile(downloadURL, targetPath)
	if err != nil {
		t.Fatal("failed to run DownloadFile()")
	}

	cmd := "chmod +x " + dlFilePath
	if _, err := script.Exec(cmd).Stdout(); err != nil {
		t.Fatalf("failed to run `chmod +x` on %s: %v", dlFilePath, err)
	}

	type params struct {
		timeout time.Duration
		cmd     string
		args    []string
	}

	// Generate random string for the test file
	rand, err := RandomString(8)
	if err != nil {
		t.Fatalf("failed to generate random string: %v", err)
	}
	// Create test script for the test #4
	testFile := filepath.Join("/tmp", fmt.Sprintf("%s-test4.sh", rand))
	testFileContent := `
#!/bin/bash
set -ex

sleep 5

# Kill this process
ps -ef | \
	grep "${0}" | \
	grep -v grep | \
	awk '{print $2}' | \
	xargs -r kill -9
`
	if err := CreateFile(testFile, []byte(testFileContent)); err != nil {
		if err != nil {
			t.Fatalf("failed to create %s with %s using CreateFile(): %v", testFile, testFileContent, err)
		}
	}
	// Remove the temporary file after the test completes.
	defer func() {
		if err := DeleteFile(testFile); err != nil {
			t.Fatalf("unable to delete %s, DeleteFile() failed", testFile)
		}
	}()

	tests := []struct {
		name    string
		params  params
		wantErr bool
		wantOut string
	}{
		{
			name: "Test command that runs quickly",
			params: params{
				timeout: time.Duration(5) * time.Second,
				cmd:     "echo",
				args:    []string{"hi"},
			},
			wantErr: false,
			wantOut: "hi",
		},
		{
			name: "Test running command that will not finish quickly",
			params: params{
				timeout: time.Duration(5) * time.Second,
				cmd:     "sleep",
				args:    []string{"250"},
			},
			wantErr: false,
			wantOut: "",
		},
		{
			name: "Test long-running bash script that will not finish quickly",
			params: params{
				timeout: time.Duration(10) * time.Second,
				cmd:     "bash",
				args:    []string{dlFilePath},
			},
			wantErr: false,
			wantOut: "USER/GROUP",
		},
		{
			name: "Test process that times out before the specified timeout",
			params: params{
				timeout: time.Duration(10) * time.Second,
				cmd:     "bash",
				args:    []string{testFile},
			},
			wantErr: true,
			wantOut: "",
		},
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := RunCommandWithTimeout(tt.params.timeout, tt.params.cmd, tt.params.args...)
				if err != nil && !tt.wantErr {
					t.Errorf("error: RunCommandWithTimeout() err = %v", err)
				}
				if !strings.Contains(got, tt.wantOut) {
					t.Errorf("error: RunCommandWithTimeout() got = %v, want %v", got, tt.wantOut)
				}
				if debug {
					log.Println("Command output: ", got)
				}
			})
		}
	default:
		t.Fatal("unsupported OS detected")
	}
}
