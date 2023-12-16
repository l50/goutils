package sys_test

import (
	"fmt"
	"log"
	"os"
	"time"

	fileutils "github.com/l50/goutils/v2/file/fileutils"
	"github.com/l50/goutils/v2/sys"
)

func ExampleCheckRoot() {
	err := sys.CheckRoot()
	uid := os.Geteuid()
	if uid != 0 {
		if err != nil {
			fmt.Println("the process must be run as root.")
		}
		// Output: the process must be run as root.
	}
}

func ExampleCd() {
	dir := "/tmp" // choose a directory that should exist on the testing machine
	err := sys.Cd(dir)

	if err != nil {
		fmt.Printf("failed to change directory to %s: %v\n", dir, err)
	} else {
		fmt.Println("Directory changed successfully!")
	}

	// Output: Directory changed successfully!
}

func ExampleCmdExists() {
	if !sys.CmdExists("ls") {
		log.Printf("the input command %s is not available on this system", "ls")
	}
}

func ExampleCp() {
	err := sys.Cp("/path/to/src", "/path/to/dst")

	if err != nil {
		log.Printf("failed to copy %s to %s: %v", "/path/to/src", "/path/to/dst", err)
	}
}

func ExampleEnvVarSet() {
	if err := sys.EnvVarSet("HOME"); err != nil {
		log.Println("the HOME environment variable is not set")
	}
}

func ExampleExpandHomeDir() {
	path := "~/Documents/project"
	expandedPath := sys.ExpandHomeDir(path)
	log.Println("Expanded path:", expandedPath)
}

func ExampleGetHomeDir() {
	homeDir, err := sys.GetHomeDir()

	if err != nil {
		log.Fatalf("failed to get home dir: %v", err)
	}

	log.Println("Home directory:", homeDir)
}

func ExampleGetSSHPubKey() {
	keyName := "id_rsa"
	password := "mypassword"

	publicKey, err := sys.GetSSHPubKey(keyName, password)

	if err != nil {
		log.Fatalf("failed to get SSH public key: %v", err)
	}

	log.Printf("Retrieved public key: %v", publicKey)
}

func ExampleGwd() {
	cwd := sys.Gwd()

	if cwd == "" {
		log.Fatalf("failed to get cwd")
	}

	log.Println("current working directory:", cwd)
}

func ExampleGetFutureTime() {
	futureTime := sys.GetFutureTime(1, 2, 3)
	log.Println("future date and time:", futureTime)
}

func ExampleGetOSAndArch() {
	osName, archName, err := sys.GetOSAndArch(&sys.DefaultRuntimeInfoProvider{})

	if err != nil {
		log.Fatalf("error detecting OS and architecture: %v", err)
	} else {
		log.Printf("Detected OS: %s, Architecture: %s\n", osName, archName)
	}
}

func ExampleIsDirEmpty() {
	isEmpty, err := sys.IsDirEmpty("/path/to/directory")

	if err != nil {
		log.Fatalf("error checking directory: %v", err)
	}

	fmt.Println("is directory empty:", isEmpty)
}

func ExampleKillProcess() {
	err := sys.KillProcess(1234, sys.SignalKill)

	if err != nil {
		log.Printf("failed to kill process: %v", err)
	}
}

func ExampleRunCommand() {
	output, err := sys.RunCommand("ls", "-l")

	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

	log.Println("Command output:", output)
}

func ExampleRunCommandWithTimeout() {
	output, err := sys.RunCommandWithTimeout(5, "sleep", "10")

	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

	log.Println("Command output:", output)
}

func ExampleRmRf() {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatalf("failed to create temp directory: %v", err)
	}

	// The temporary directory will be removed at the end of this function
	defer os.RemoveAll(tmpDir)

	// Convert tmpDir to RealFile type
	file := fileutils.RealFile(tmpDir)

	// Use RmRf to remove the directory
	if err := sys.RmRf(file); err != nil {
		log.Printf("error removing path: %v", err)
		return
	}

	// Check if the directory was successfully removed
	_, err = os.Stat(tmpDir)
	if err == nil || !os.IsNotExist(err) {
		log.Printf("directory was not removed: %v", err)
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
		OutputHandler: nil,
	}

	output, err := cmd.RunCmd()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	fmt.Print(output)
	// Output: Hello, world!
}
