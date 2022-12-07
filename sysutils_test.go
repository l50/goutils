package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bitfield/script"
)

func TestCd(t *testing.T) {
	dst := "magefiles"

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
		if err := Cp(file, copyLoc); err != nil {
			t.Fatalf("failed to copy %s to %s - Cp() failed", file, copyLoc)
		}
		if FileExists(copyLoc) {
			os.Remove(file)
			os.Remove(copyLoc)
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
		t.Errorf("failed to run `chmod +x` on %s: %v", dlFilePath, err)
	}

	type args struct {
		// timeout time.Duration
		timeout string
		command string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test command that runs quickly",
			args: args{
				timeout: "5s",
				command: "echo hi",
			},
			wantErr: false,
			want:    "hi",
		},
		{
			name: "Test running command that will not finish quickly",
			args: args{
				timeout: "5s",
				command: "sleep 250",
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "Test long-running bash script that will not finish quickly",
			args: args{
				timeout: "10s",
				command: "bash " + dlFilePath,
			},
			wantErr: true,
			want:    "",
		},
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := RunCommandWithTimeout(tt.args.timeout, tt.args.command)
				fmt.Println(got)
				if (err != nil) != tt.wantErr {
					t.Errorf("error: RunCommandWithTimeout() err = %v, want %v", err, tt.wantErr)
				}
				if len(got) == 0 && tt.want != "" && got != tt.want {
					t.Errorf("error: RunCommandWithTimeout() got = %v, want %v", got, tt.want)
				}
			})
		}
	default:
		t.Fatal("unsupported OS detected")
	}
}
