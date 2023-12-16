package sys_test

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	fileutils "github.com/l50/goutils/v2/file/fileutils"
	"github.com/l50/goutils/v2/sys"
	"github.com/stretchr/testify/assert"
)

func TestCheckRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping this test as it needs to be run as root")
	}

	if err := sys.CheckRoot(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

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

	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatal(err)
		}
	}()

	// Setup test cases
	tests := []struct {
		name       string
		path       string
		expectErr  bool
		errMessage string
	}{
		{
			name:      "existing directory",
			path:      tmpDir,
			expectErr: false,
		},
		{
			name:       "non-existent directory",
			path:       filepath.Join(tmpDir, "nonexistent"),
			expectErr:  true,
			errMessage: "no such file or directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := sys.Cd(tc.path)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				if !strings.Contains(err.Error(), tc.errMessage) {
					t.Fatalf("unexpected error message: got %v, want %s", err, tc.errMessage)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.expectErr {
				pwd, err := os.Getwd()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				// Ensure both paths are evaluated to their real paths (resolving any symlinks)
				realCurrentDir, err := filepath.EvalSymlinks(pwd)
				if err != nil {
					t.Fatal(err)
				}
				realTestDir, err := filepath.EvalSymlinks(tc.path)
				if err != nil {
					t.Fatal(err)
				}

				if realCurrentDir != realTestDir {
					t.Fatalf("did not change directory: expected %s but got %s", realTestDir, realCurrentDir)
				}
			}
		})
	}
}

func TestCmdExists(t *testing.T) {
	tests := []struct {
		name   string
		cmd    string
		expect bool
	}{
		{
			name:   "Command Exists",
			cmd:    "ls",
			expect: true,
		},
		{
			name:   "Command Does Not Exist",
			cmd:    "unknowncommand",
			expect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := sys.CmdExists(tc.cmd)
			if result != tc.expect {
				t.Errorf("Expected %v, but got %v", tc.expect, result)
			}
		})
	}
}

func TestCp(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		dst     string
		setup   func() error
		cleanup func()
		err     error
	}{
		{
			name: "FailedToCopy",
			src:  "nonexistentfile",
			dst:  "destinationfile",
			err:  fmt.Errorf("failed to copy nonexistentfile to destinationfile: lstat nonexistentfile: no such file or directory"), // Updated error message
		},
		{
			name: "SuccessfulCopy",
			src:  "test.txt",
			dst:  "testing.txt",
			setup: func() error {
				return fileutils.Create("test.txt", nil, fileutils.CreateEmptyFile)
			},
			cleanup: func() {
				remove := []string{"test.txt", "testing.txt"}
				for _, f := range remove {
					if err := fileutils.Delete(f); err != nil {
						t.Errorf("unable to delete %s, DeleteFile() failed", f)
					}
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				if err := tc.setup(); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := sys.Cp(tc.src, tc.dst)

			if (err != nil && tc.err == nil) || (err == nil && tc.err != nil) || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
				t.Fatalf("unexpected error: got %v, want %v", err, tc.err)
			}

			if tc.cleanup != nil {
				tc.cleanup()
			}
		})
	}
}

func TestEnvVarSet(t *testing.T) {
	const (
		key      = "TEST_KEY"
		emptykey = "EMPTY_TEST_KEY"
	)

	testCases := []struct {
		name    string
		key     string
		setenv  bool
		wantErr bool
	}{
		{
			name:    "KeyIsSet",
			key:     key,
			setenv:  true,
			wantErr: false,
		},
		{
			name:    "KeyIsNotSet",
			key:     emptykey,
			setenv:  false,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setenv {
				os.Setenv(tc.key, "test_value")
			}
			err := sys.EnvVarSet(tc.key)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Errorf("EnvVarSet(%q) returned error: %v, wantErr: %v", tc.key, gotErr, tc.wantErr)
			}
		})
	}
}

func TestExpandHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home directory: %v", err)
	}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "EmptyPath",
			input:    "",
			expected: "",
		},
		{
			name:     "NoTilde",
			input:    "/path/without/tilde",
			expected: "/path/without/tilde",
		},
		{
			name:     "TildeOnly",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "TildeWithSlash",
			input:    "~/path/with/slash",
			expected: filepath.Join(homeDir, "path/with/slash"),
		},
		{
			name:     "TildeWithoutSlash",
			input:    "~path/without/slash",
			expected: filepath.Join(homeDir, "path/without/slash"),
		},
		{
			name:     "TildeWithExtraSlash",
			input:    "~/path//with//slash",
			expected: filepath.Join(homeDir, "path//with//slash"),
		},
		{
			name:     "TildeWithDot",
			input:    "~./path/with/dot",
			expected: filepath.Join(homeDir, "./path/with/dot"),
		},
		{
			name:     "TildeWithMultipleDots",
			input:    "~../path/with/dots",
			expected: filepath.Join(homeDir, "../path/with/dots"),
		},
		{
			name:     "TildeWithSpace",
			input:    "~ path/with/space",
			expected: filepath.Join(homeDir, " path/with/space"),
		},
		{
			name:     "TildeWithSpecialCharacters",
			input:    "~@#$/path/with/special",
			expected: filepath.Join(homeDir, "@#$/path/with/special"),
		},
		{
			name:     "MultipleTildes",
			input:    "~~~/path/with/tildes",
			expected: filepath.Join(homeDir, "~~/path/with/tildes"),
		},
		{
			name:     "TildeInMiddle",
			input:    "/path/with/~tilde/in/middle",
			expected: "/path/with/~tilde/in/middle",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := sys.ExpandHomeDir(tc.input)
			if actual != tc.expected {
				t.Errorf("test failed: ExpandHomeDir(%q) = %q; expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestGetHomeDir(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Successful fetch of user's home directory",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sys.GetHomeDir()
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetHomeDir() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGetSSHPubKey(t *testing.T) {
	tests := []struct {
		name        string
		keyName     string
		password    string
		expectError bool
	}{
		{
			name:        "valid key",
			keyName:     "id_rsa",
			password:    "mypassword",
			expectError: false,
		},
		{
			name:        "invalid key",
			keyName:     "invalid_key",
			password:    "mypassword",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory for each test case
			tempDir, err := os.MkdirTemp("", "ssh")
			if err != nil {
				t.Fatalf("unable to create temporary directory: %v", err)
			}
			defer os.RemoveAll(tempDir) // clean up

			// Overwrite HOME environment variable
			os.Setenv("HOME", tempDir)

			// Create a dummy key file for the valid key test case
			if !tc.expectError {
				dummyKeyPath := filepath.Join(tempDir, ".ssh", tc.keyName)

				if err := os.MkdirAll(filepath.Dir(dummyKeyPath), os.ModePerm); err != nil {
					t.Fatalf("unable to create directory %s: %v", filepath.Dir(dummyKeyPath), err)
				}
				err := os.WriteFile(dummyKeyPath, []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAgEP7Ub1Z9oOJFoYNB9E75RJgRzdUOVpzIE4ZCcRCW0QAAAJAxRdmiMUXZ
ogAAAAtzc2gtZWQyNTUxOQAAACAgEP7Ub1Z9oOJFoYNB9E75RJgRzdUOVpzIE4ZCcRCW0Q
AAAECSVf0Sfthqq7p8qeCtHEUYq8M2LSARcpaT32hU4vZf2CAQ/tRvVn2g4kWhg0H0TvlE
mBHN1Q5WnMgThkJxEJbRAAAABm5vbmFtZQECAwQFBgc=
-----END OPENSSH PRIVATE KEY-----`), os.ModePerm)
				if err != nil {
					t.Fatalf("unable to create dummy key file %s: %v", dummyKeyPath, err)
				}
			}

			_, err = sys.GetSSHPubKey(tc.keyName, tc.password)
			if tc.expectError && err == nil {
				t.Fatalf("expected an error but did not get one")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("did not expect an error but got one: %v", err)
			}
		})
	}
}

func TestGwd(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Current working directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir := t.TempDir()

			// Change the current working directory to the temporary directory
			err := os.Chdir(tempDir)
			assert.NoError(t, err)

			result := sys.Gwd()
			if runtime.GOOS == "darwin" {
				assert.Equal(t, filepath.Join("/private", tempDir), result)
			} else {
				assert.Equal(t, tempDir, result)
			}
		})
	}
}

func TestGetFutureTime(t *testing.T) {
	tests := []struct {
		name      string
		years     int
		months    int
		days      int
		expResult time.Time
	}{
		{
			name:      "Future time in 1 year, 2 months, 3 days",
			years:     1,
			months:    2,
			days:      3,
			expResult: time.Now().AddDate(1, 2, 3),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := sys.GetFutureTime(tc.years, tc.months, tc.days)
			assert.Equal(t, tc.expResult.Year(), result.Year())
			assert.Equal(t, tc.expResult.Month(), result.Month())
			assert.Equal(t, tc.expResult.Day(), result.Day())
		})
	}
}

type MockRuntimeInfoProvider struct{}

func (p *MockRuntimeInfoProvider) GetOS() string {
	return "linux"
}

func (p *MockRuntimeInfoProvider) GetArch() string {
	return "unsupported_arch"
}

func TestGetOSAndArch(t *testing.T) {
	tests := []struct {
		name        string
		provider    sys.RuntimeInfoProvider
		expectOS    string
		expectArch  string
		expectError bool
	}{
		{
			name:        "test on darwin amd64",
			provider:    &sys.DefaultRuntimeInfoProvider{},
			expectOS:    runtime.GOOS,
			expectArch:  runtime.GOARCH,
			expectError: false,
		},
		{
			name:        "test on unsupported architecture",
			provider:    &MockRuntimeInfoProvider{},
			expectOS:    "",
			expectArch:  "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			osName, archName, err := sys.GetOSAndArch(tc.provider)
			if tc.expectError && err == nil {
				t.Fatalf("expected an error but did not get one")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("did not expect an error but got one: %v", err)
			}
			if osName != tc.expectOS {
				t.Fatalf("expected %s but got %s", tc.expectOS, osName)
			}
			if archName != tc.expectArch {
				t.Fatalf("expected %s but got %s", tc.expectArch, archName)
			}
		})
	}
}

func TestIsDirEmpty(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(tmpDir string) (string, error) // Setup function now returns a string (path)
		isEmpty    bool
		expectErr  bool
		errMessage string
	}{
		{
			name: "empty directory",
			setup: func(tmpDir string) (string, error) {
				return tmpDir, nil
			},
			isEmpty:   true,
			expectErr: false,
		},
		{
			name: "non-empty directory",
			setup: func(tmpDir string) (string, error) {
				filePath := filepath.Join(tmpDir, "test.txt")
				file, err := os.Create(filePath)
				if err != nil {
					return "", err
				}
				file.Close()
				return tmpDir, nil
			},
			isEmpty:   false,
			expectErr: false,
		},
		{
			name: "non-existent directory",
			setup: func(tmpDir string) (string, error) {
				nonExistentDir := filepath.Join(tmpDir, "nonexistent")
				return nonExistentDir, nil
			},
			isEmpty:    false,
			expectErr:  true,
			errMessage: "does not exist",
		},
		{
			name: "file instead of directory",
			setup: func(tmpDir string) (string, error) {
				filePath := filepath.Join(tmpDir, "file.txt")
				file, err := os.Create(filePath)
				if err != nil {
					return "", err
				}
				file.Close()
				return filePath, nil
			},
			isEmpty:    false,
			expectErr:  true,
			errMessage: "not a directory",
		},
		{
			name: "failed to read directory entries",
			setup: func(tmpDir string) (string, error) {
				// Create a directory without read permissions
				dirPath := filepath.Join(tmpDir, "unreadable")
				if err := os.Mkdir(dirPath, 0111); err != nil {
					return "", err
				}
				return dirPath, nil
			},
			isEmpty:    false,
			expectErr:  true,
			errMessage: "failed to read directory entries",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			dirPath, err := tc.setup(tmpDir)
			if err != nil {
				t.Fatalf("failed to set up test case: %v", err)
			}

			isEmpty, err := sys.IsDirEmpty(dirPath)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				if !strings.Contains(err.Error(), tc.errMessage) {
					t.Fatalf("unexpected error message: got %v, want %s", err, tc.errMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: got %v, wantErr false", err)
			}

			if isEmpty != tc.isEmpty {
				t.Fatalf("unexpected result: got %t, want %t", isEmpty, tc.isEmpty)
			}
		})
	}
}

func TestKillProcess(t *testing.T) {
	tests := []struct {
		name   string
		pid    int
		signal sys.Signal
		err    error
	}{
		{
			name:   "kill process on Windows with sys.SignalKill",
			pid:    1234,
			signal: sys.SignalKill,
			err:    nil,
		},
		{
			name:   "kill process on non-Windows with unsupported signal",
			pid:    5678,
			signal: sys.Signal(999),
			err:    fmt.Errorf("unsupported signal: %v", sys.Signal(999)),
		},
		{
			name:   "kill process on non-Windows with sys.SignalKill",
			pid:    5678,
			signal: sys.SignalKill,
			err:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run a process to kill
			cmd := exec.Command("go", "run", "-e", `
				package main
				import (
					"os"
					"os/signal"
					"syscall"
					"time"
				)
				func main() {
					c := make(chan os.sys.Signal, 1)
					signal.Notify(c, syscall.SIGTERM)
					<-c
					time.Sleep(10 * time.Second)
				}
			`)

			if err := cmd.Start(); err != nil {
				t.Fatalf("failed to start process: %v", err)
			}
			pid := cmd.Process.Pid

			// Delay to allow the process to start
			time.Sleep(100 * time.Millisecond)

			err := sys.KillProcess(pid, tc.signal)

			if (err != nil && tc.err == nil) || (err == nil && tc.err != nil) || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
				t.Fatalf("unexpected error: got %v, want %v", err, tc.err)
			}

			// Cleanup the process
			err = cmd.Wait()
			if err == nil {
				t.Fatalf("process %d should be terminated, but Wait() returned without error", pid)
			}
		})
	}
}

type MockFile struct {
	open      func() (io.ReadCloser, error)
	write     func(contents []byte, perm os.FileMode) error
	remove    func() error
	removeAll func() error
	stat      func() (os.FileInfo, error)
}

func (m *MockFile) Open() (io.ReadCloser, error) {
	return m.open()
}
func (m *MockFile) Write(contents []byte, perm os.FileMode) error {
	return m.write(contents, perm)
}
func (m *MockFile) Remove() error {
	return m.remove()
}

func (m *MockFile) RemoveAll() error {
	return m.removeAll()
}

func (m *MockFile) Stat() (os.FileInfo, error) {
	return m.stat()
}

type MockFileInfo struct {
	isDir bool
}

func (m *MockFileInfo) IsDir() bool {
	return m.isDir
}

// Add dummy implementations for the other methods required by the os.FileInfo interface.
func (m *MockFileInfo) Name() string       { return "" }
func (m *MockFileInfo) Size() int64        { return 0 }
func (m *MockFileInfo) Mode() fs.FileMode  { return 0 }
func (m *MockFileInfo) ModTime() time.Time { return time.Now() }
func (m *MockFileInfo) Sys() interface{}   { return nil }

func TestRmRf(t *testing.T) {
	tests := []struct {
		name    string
		file    fileutils.File
		wantErr bool
	}{
		{
			name: "Path is a directory",
			file: &MockFile{
				open: func() (io.ReadCloser, error) {
					return nil, nil
				},
				write: func(contents []byte, perm os.FileMode) error {
					return nil
				},
				stat: func() (os.FileInfo, error) {
					return &MockFileInfo{isDir: true}, nil
				},
				removeAll: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Path is a file",
			file: &MockFile{
				stat: func() (os.FileInfo, error) {
					return &MockFileInfo{isDir: false}, nil
				},
				remove: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Path does not exist",
			file: &MockFile{
				stat: func() (os.FileInfo, error) {
					return nil, os.ErrNotExist
				},
			},
			wantErr: true,
		},
		{
			name: "os.RemoveAll fails",
			file: &MockFile{
				stat: func() (os.FileInfo, error) {
					return &MockFileInfo{isDir: true}, nil
				},
				removeAll: func() error {
					return os.ErrPermission
				},
			},
			wantErr: true,
		},
		{
			name: "os.Stat fails",
			file: &MockFile{
				stat: func() (os.FileInfo, error) {
					return nil, os.ErrInvalid
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := sys.RmRf(tc.file)
			if (err != nil) != tc.wantErr {
				t.Errorf("RmRf() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGetTempPath(t *testing.T) {
	testCases := []struct {
		name     string
		os       string
		expected string
	}{
		{
			name:     "Windows OS",
			os:       "windows",
			expected: "C:\\Temp",
		},
		{
			name:     "Unix/Linux OS",
			os:       "linux",
			expected: "/tmp",
		},
		{
			name:     "Mac OS",
			os:       "darwin",
			expected: "/tmp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if runtime.GOOS == tc.os {
				tempPath := sys.GetTempPath()
				assert.Equal(t, tc.expected, tempPath, "GetTempPath should return correct path for %s", tc.os)
			} else {
				t.Skipf("Skipping test for %s on %s", tc.os, runtime.GOOS)
			}
		})
	}
}

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name       string
		cmd        string
		args       []string
		wantError  bool
		wantOutput string
	}{
		{
			name:       "EchoTest",
			cmd:        "echo",
			args:       []string{"Hello, world!"},
			wantError:  false,
			wantOutput: "Hello, world!\n",
		},
		{
			name:      "InvalidCommand",
			cmd:       "someinvalidcommand",
			args:      []string{},
			wantError: true,
		},
		{
			name:       "ColorTest",
			cmd:        "echo",
			args:       []string{"\033[31mRed\033[0m"},
			wantError:  false,
			wantOutput: "\033[31mRed\033[0m\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := sys.RunCommand(tc.cmd, tc.args...)
			if (err != nil) != tc.wantError {
				t.Errorf("RunCommand() error = %v, wantError %v", err, tc.wantError)
				return
			}

			// if we expect an output, let's check it
			if tc.wantOutput != "" && !strings.Contains(output, tc.wantOutput) {
				t.Errorf("expected output '%s' not found in: '%s'", tc.wantOutput, output)
			}
		})
	}
}

func TestRunCommandWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout int
		cmd     string
		args    []string
		wantErr bool
	}{
		{
			name:    "Test command that runs quickly",
			timeout: 5,
			cmd:     "echo",
			args:    []string{"hi"},
			wantErr: false,
		},
		{
			name:    "Test command that takes longer than timeout",
			timeout: 2,
			cmd:     "sleep",
			args:    []string{"5"},
			wantErr: true,
		},
		{
			name:    "Test command that fails",
			timeout: 5,
			cmd:     "ls",
			args:    []string{"/nonexistentpath"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sys.RunCommandWithTimeout(tc.timeout, tc.cmd, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("RunCommandWithTimeout() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func TestRunCmd(t *testing.T) {
	testCases := []*struct {
		name           string
		cmd            sys.Cmd
		expectedOutput string
		expectError    bool
	}{
		{
			name: "Successful Command Execution",
			cmd: sys.Cmd{
				CmdString:     "echo",
				Args:          []string{"hello"},
				Timeout:       0,
				OutputHandler: nil,
			},
			expectedOutput: "hello",
			expectError:    false,
		},
		{
			name: "Invalid Command",
			cmd: sys.Cmd{
				CmdString:     "invalidcommand",
				Args:          []string{},
				Timeout:       0,
				OutputHandler: nil,
			},
			expectedOutput: "",
			expectError:    true,
		},
		{
			name: "Command with Timeout",
			cmd: sys.Cmd{
				CmdString:     "sleep",
				Args:          []string{"5"},
				Timeout:       5 * time.Second,
				OutputHandler: nil,
			},
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Run the command
			output, err := tc.cmd.RunCmd() // Capture both output and error
			// Check the output and error
			if (err != nil) != tc.expectError {
				t.Errorf("Test '%s' failed: expected error: %v, got: %v",
					tc.name, tc.expectError, err)
			}

			if strings.TrimSpace(output) != tc.expectedOutput {
				t.Errorf("Test '%s' failed: expected output: %q, got: %q",
					tc.name, tc.expectedOutput, output)
			}
		})
	}
}
