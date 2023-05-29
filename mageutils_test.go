package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	randStr, _ := RandomString(8)
	clonePath = createTestRepo(fmt.Sprintf("mageutils-%s", randStr))
}

func TestGHRelease(t *testing.T) {
	// Call the function with an old version
	newVer := "v1.0.0"
	if err := GHRelease(newVer); err == nil {
		t.Errorf("release %s should not have been created: %v", newVer, err)
	}
}

func TestGoReleaser(t *testing.T) {
	repoRoot, err := RepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "goutils-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory to temp directory: %v", err)
	}

	// Copy the repo to temporary directory
	if err := Cp(repoRoot, tempDir); err != nil {
		t.Fatalf("failed to copy repo to temp directory: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	releaserDir := filepath.Join(cwd, "dist")

	if err := os.Mkdir(releaserDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := GoReleaser(); err != nil {
		t.Fatal(err)
	}
}

func TestInstallVSCodeModules(t *testing.T) {
	if err := InstallVSCodeModules(); err != nil {
		t.Fatal(err)
	}
}

func TestModUpdate(t *testing.T) {
	// First test
	recursive := false
	verbose := true
	if err := ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}

	// Second test
	recursive = true
	verbose = false
	if err := ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}
}

func TestTidy(t *testing.T) {
	if err := Tidy(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMageDeps(t *testing.T) {
	if err := UpdateMageDeps("magefiles"); err != nil {
		t.Fatal(err)
	}
}

func TestInstallGoDeps(t *testing.T) {
	sampleDeps := []string{
		"golang.org/x/lint/golint",
		"golang.org/x/tools/cmd/goimports",
	}

	if err := InstallGoDeps(sampleDeps); err != nil {
		t.Fatal(err)
	}
}

func TestFindExportedFunctionsInPackage(t *testing.T) {
	// Define the bash command as a string
	bashCmd := `
find . -name "*.go" -not -path "./magefiles/*" |
xargs grep -E -o 'func [A-Z][a-zA-Z0-9_]+\(' |
grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn,.idea,.tox} -v '_test.go' |
grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn,.idea,.tox} -v -E 'func [A-Z][a-zA-Z0-9_]+Test\(' |
sed -e 's/func //' -e 's/(//' |
awk -F: '{printf "Function: %s\nFile: %s\n", $2, $1}'`

	// Create a new command and set its properties
	cmd := exec.Command("bash", "-c", bashCmd)
	cmd.Dir = "."
	cmd.Env = os.Environ()

	// Run the command and get its output
	outputBytes, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}
	output := string(outputBytes)

	// Parse the output and create a map of expected function names
	bashFuncs := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Function: ") {
			funcName := strings.TrimSpace(strings.TrimPrefix(line, "Function: "))
			bashFuncs[funcName] = true
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("error while scanning command output: %v", err)
	}

	// Define a table of test cases with input values and expected results
	tests := []struct {
		name           string
		packagePath    string
		expectedFuncs  map[string]bool
		expectedErrors bool
	}{
		{
			name:           "Valid package",
			packagePath:    ".",
			expectedFuncs:  bashFuncs,
			expectedErrors: false,
		},
		{
			name:           "Invalid package",
			packagePath:    "nonexistent_package",
			expectedFuncs:  nil,
			expectedErrors: true,
		},
	}
	// Loop through the test cases and execute each one
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Get the exported functions from the package
			goFuncs, err := FindExportedFunctionsInPackage(tc.packagePath)
			if tc.expectedErrors && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tc.expectedErrors && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Compare the expected and actual functions
			missingFuncs := []string{}
			for bf := range tc.expectedFuncs {
				found := false
				for _, gf := range goFuncs {
					if bf == gf.FuncName {
						found = true
						break
					}
				}
				if !found {
					missingFuncs = append(missingFuncs, bf)
				}
			}

			if len(missingFuncs) > 0 {
				t.Errorf("go and bash implementations don't agree: %v", missingFuncs)
			}
		})
	}
}
