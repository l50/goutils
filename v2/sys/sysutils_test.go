package sys_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bitfield/script"
	"github.com/l50/goutils/v2/file"
	fileutils "github.com/l50/goutils/v2/file"
	"github.com/l50/goutils/v2/net"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"
)

var (
	debug = false
)

func TestCd(t *testing.T) {
	// Setup a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "magefiles")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir) // ensure we return to the original directory after test

	// Now test the Cd function
	if err := sys.Cd(tmpDir); err != nil {
		t.Fatalf("error running Cd(): expected to change directory to %s but got error: %v", tmpDir, err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure both paths are evaluated to their real paths (resolving any symlinks)
	realCurrentDir, err := filepath.EvalSymlinks(currentDir)
	if err != nil {
		t.Fatal(err)
	}
	realTmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if realCurrentDir != realTmpDir {
		t.Fatalf("error running Cd(): expected current directory to be %s but got %s", realTmpDir, realCurrentDir)
	}
}

func TestCmdExists(t *testing.T) {
	fail := "asdf"
	cmds := []string{"ls", "whoami", fail}
	for _, cmd := range cmds {
		if sys.CmdExists(cmd) && cmd == fail {
			t.Fatalf(
				"failed to properly identify installed cmd - CmdExists() failed")
		}
	}
}

func TestCp(t *testing.T) {
	file := "test.txt"
	copyLoc := "testing.txt"
	created := fileutils.CreateEmpty(file)
	if created {
		if err := sys.Cp(file, copyLoc); err != nil {
			t.Fatalf("failed to copy %s to %s - Cp() failed", file, copyLoc)
		}
		if fileutils.Exists(copyLoc) {
			remove := []string{file, copyLoc}
			for _, f := range remove {
				if err := fileutils.Delete(f); err != nil {
					t.Errorf("unable to delete %s, DeleteFile() failed", f)
				}
			}
		}
	}
}

func TestEnvVarSet(t *testing.T) {
	key := "TEST_KEY"
	os.Setenv(key, "test_value")
	if err := sys.EnvVarSet(key); err != nil {
		t.Fatalf("failed to run EnvVarSet(): %v", err)
	}

	emptykey := "EMPTY_TEST_KEY"

	if err := sys.EnvVarSet(emptykey); err == nil {
		t.Fatalf("failed to run EnvVarSet(): %v", err)
	}
}

func TestGetHomeDir(t *testing.T) {
	if _, err := sys.GetHomeDir(); err != nil {
		t.Fatalf("failed to get the user's home directory - GetHomeDir() failed: %v", err)
	}
}

func TestGwd(t *testing.T) {
	out := sys.Gwd()
	if !strings.Contains(out, "goutils") {
		t.Fatal("unable to get the current working directory - Gwd() failed")
	}
}

func isTime(obj reflect.Value) bool {
	_, ok := obj.Interface().(time.Time)
	return ok
}

func TestGetFutureTime(t *testing.T) {
	futureTime := sys.GetFutureTime(2, 2, 3)

	ft := reflect.ValueOf(futureTime)
	if !isTime(ft) {
		t.Fatal("failed to run GetFutureTime(): incorrect value returned")
	}
}

func TestGetOSAndArch(t *testing.T) {
	osName, archName, err := sys.GetOSAndArch()
	if err != nil {
		t.Fatalf("failed to run GetOSAndArch(): %v", err)
	}

	validOS := []string{"linux", "darwin", "windows"}
	if !str.InSlice(osName, validOS) {
		t.Errorf("invalid OS: %s", osName)
	}

	validArch := []string{"amd64", "arm64", "armv"}
	if !str.InSlice(archName, validArch) {
		t.Errorf("invalid architecture: %s", archName)
	}
}

func TestIsDirEmpty(t *testing.T) {
	dirEmpty, err := sys.IsDirEmpty("/")
	if err != nil {
		t.Fatalf("failed to determine if / is empty - IsDirEmpty() failed: %v", err)
	}
	if dirEmpty != false {
		t.Fatal("the / directory has reported back as being empty, which can not be true - IsDirEmpty()")
	}
}

func TestKillProcess(t *testing.T) {
	// Run a process to kill
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}
	pid := cmd.Process.Pid

	// Test killing the process with sys.SignalKill signal
	if err := sys.KillProcess(pid, sys.SignalKill); err != nil {
		t.Fatalf("failed to kill process %d with SIGKILL - KillProcess() failed: %v", pid, err)
	}
}

func TestRunCommand(t *testing.T) {
	switch runtime.GOOS {
	case "linux", "darwin":
		out, err := sys.RunCommand("uname", "-a")
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
	dlFilePath, err := net.DownloadFile(downloadURL, targetPath)
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
	rand, err := str.GenRandom(8)
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
	if err := file.Create(testFile, []byte(testFileContent)); err != nil {
		if err != nil {
			t.Fatalf("failed to create %s with %s using CreateFile(): %v", testFile, testFileContent, err)
		}
	}
	// Remove the temporary file after the test completes.
	defer func() {
		if err := file.Delete(testFile); err != nil {
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
				got, err := sys.RunCommandWithTimeout(tt.params.timeout, tt.params.cmd, tt.params.args...)
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
