package sys_test

import (
	"fmt"
	"os"
	"time"

	fileutils "github.com/l50/goutils/v2/file/fileutils"
	log "github.com/l50/goutils/v2/logging"
	"github.com/l50/goutils/v2/sys"
)

func ExampleCheckRoot() {
	err := sys.CheckRoot()
	uid := os.Geteuid()
	if uid != 0 {
		if err != nil {
			log.L().Println("The process must be run as root.")
		}
		// Output: the process must be run as root.
	}
}

func ExampleCd() {
	dir := "/tmp" // choose a directory that should exist on the testing machine
	err := sys.Cd(dir)

	if err != nil {
		log.L().Errorf("Failed to change directory to %s: %v\n", dir, err)
	} else {
		log.L().Println("Directory changed successfully!")
	}

	// Output: Directory changed successfully!
}

func ExampleCmdExists() {
	if !sys.CmdExists("ls") {
		log.L().Printf("the input command %s is not available on this system", "ls")
	}
}

func ExampleCp() {
	err := sys.Cp("/path/to/src", "/path/to/dst")

	if err != nil {
		log.L().Errorf("Failed to copy %s to %s: %v", "/path/to/src", "/path/to/dst", err)
	}
}

func ExampleEnvVarSet() {
	if err := sys.EnvVarSet("HOME"); err != nil {
		log.L().Println("the HOME environment variable is not set")
	}
}

func ExampleExpandHomeDir() {
	path := "~/Documents/project"
	expandedPath := sys.ExpandHomeDir(path)
	log.L().Println("Expanded path:", expandedPath)
}

func ExampleGetHomeDir() {
	homeDir, err := sys.GetHomeDir()

	if err != nil {
		log.L().Errorf("Failed to get home dir: %v", err)
	}

	log.L().Println("Home directory:", homeDir)
}

func ExampleGetSSHPubKey() {
	keyName := "id_rsa"
	password := "mypassword"

	publicKey, err := sys.GetSSHPubKey(keyName, password)

	if err != nil {
		log.L().Errorf("Failed to get SSH public key: %v", err)
	}

	log.L().Printf("Retrieved public key: %v", publicKey)
}

func ExampleGwd() {
	cwd := sys.Gwd()

	if cwd == "" {
		log.L().Error("Failed to get cwd")
	}

	log.L().Println("Current working directory:", cwd)
}

func ExampleGetFutureTime() {
	futureTime := sys.GetFutureTime(1, 2, 3)
	log.L().Println("Future date and time:", futureTime)
}

func ExampleGetOSAndArch() {
	osName, archName, err := sys.GetOSAndArch(&sys.DefaultRuntimeInfoProvider{})

	if err != nil {
		log.L().Errorf("Error detecting OS and architecture: %v", err)
	} else {
		log.L().Printf("Detected OS: %s, Architecture: %s\n", osName, archName)
	}
}

func ExampleIsDirEmpty() {
	isEmpty, err := sys.IsDirEmpty("/path/to/directory")

	if err != nil {
		log.L().Errorf("Error checking directory: %v", err)
	}

	log.L().Println("Is directory empty:", isEmpty)
}

func ExampleKillProcess() {
	err := sys.KillProcess(1234, sys.SignalKill)

	if err != nil {
		log.L().Errorf("Failed to kill process: %v", err)
	}
}

func ExampleRunCommand() {
	output, err := sys.RunCommand("ls", "-l")

	if err != nil {
		log.L().Errorf("Error running command: %v", err)
	}

	log.L().Println("Command output:", output)
}

func ExampleRunCommandWithTimeout() {
	output, err := sys.RunCommandWithTimeout(5, "sleep", "10")

	if err != nil {
		log.L().Errorf("Error running command: %v", err)
	}

	log.L().Println("Command output:", output)
}

func ExampleRmRf() {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.L().Errorf("Failed to create temp directory: %v", err)
	}

	// The temporary directory will be removed at the end of this function
	defer os.RemoveAll(tmpDir)

	// Convert tmpDir to RealFile type
	file := fileutils.RealFile(tmpDir)

	// Use RmRf to remove the directory
	if err := sys.RmRf(file); err != nil {
		log.L().Errorf("Error removing path: %v", err)
	}

	// Check if the directory was successfully removed
	_, err = os.Stat(tmpDir)
	if err == nil || !os.IsNotExist(err) {
		log.L().Errorf("Directory was not removed: %v", err)
		return
	}

	fmt.Println("Path successfully removed!")
	// Output: Path successfully removed!
}

func ExampleGetTempPath() {
	tempPath := sys.GetTempPath()
	fmt.Println("Temporary path:", tempPath)
}

func ExampleCmd_RunCmd() {
	cmd := sys.Cmd{
		CmdString:     "echo",
		Args:          []string{"Hello, world!"},
		Timeout:       5 * time.Second,
		OutputHandler: func(s string) { log.L().Println(s) },
	}

	output, err := cmd.RunCmd()
	if err != nil {
		log.L().Errorf("Error executing command: %v\n", err)
		return
	}

	log.L().Println(output)
	// Output: Hello, world!
}
