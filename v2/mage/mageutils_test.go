package mage_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/mage"
	"github.com/l50/goutils/v2/str"
)

var (
	mageCleanupDirs []string
)

func init() {
	// Create test repo and queue it for cleanup
	randStr, _ := str.GenRandom(8)
	clonePath := createTestRepo(fmt.Sprintf("mageutils-%s", randStr))
	mageCleanupDirs = append(mageCleanupDirs, clonePath)
}

func createTestRepo(name string) string {
	cloneDir := "/tmp"
	var currentTime time.Time
	targetPath := filepath.Join(
		cloneDir, fmt.Sprintf(
			"%s-%s", name, currentTime.Format("2006-01-02-15-04-05"),
		),
	)

	testRepoURL := "https://github.com/l50/helloworld.git"
	if _, err := git.CloneRepo(testRepoURL, targetPath, nil); err != nil {
		log.Fatalf(
			"failed to clone to %s - CloneRepo() failed: %v",
			targetPath,
			err,
		)
	}

	return targetPath
}

func TestGHRelease(t *testing.T) {
	// Call the function with an old version
	newVer := "v1.0.0"
	if err := mage.GHRelease(newVer); err == nil {
		t.Errorf("release %s should not have been created: %v", newVer, err)
	}
}

func cleanupMageUtils(t *testing.T) {
	for _, dir := range mageCleanupDirs {
		if err := fileutils.RmRf(dir); err != nil {
			fmt.Println("failed to clean up mageUtils: ", err.Error())
		}
	}
}

func TestGoReleaser(t *testing.T) {
	t.Cleanup(func() {
		cleanupMageUtils(t)
	})

	// Get repo root
	repoRoot, err := goutils.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	// Change into repo root
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("failed to change directory to repo root: %v", err)
	}

	releaserDir := filepath.Join(repoRoot, "dist")

	mageCleanupDirs = append(mageCleanupDirs, releaserDir)

	if err := mage.GoReleaser(); err != nil {
		t.Fatal(err)
	}
}

func TestInstallVSCodeModules(t *testing.T) {
	if err := mage.InstallVSCodeModules(); err != nil {
		t.Fatal(err)
	}
}

func TestModUpdate(t *testing.T) {
	// First test
	recursive := false
	verbose := true
	if err := mage.ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}

	// Second test
	recursive = true
	verbose = false
	if err := mage.ModUpdate(recursive, verbose); err != nil {
		t.Fatal(err)
	}
}

func TestTidy(t *testing.T) {
	if err := mage.Tidy(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMageDeps(t *testing.T) {
	repoRoot, err := goutils.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	os.Chdir(repoRoot)

	if err := mage.UpdateMageDeps("magefiles"); err != nil {
		t.Fatal(err)
	}
}

func TestInstallGoDeps(t *testing.T) {
	sampleDeps := []string{
		"golang.org/x/lint/golint",
		"golang.org/x/tools/cmd/goimports",
	}

	if err := mage.InstallGoDeps(sampleDeps); err != nil {
		t.Fatal(err)
	}
}

func TestFindExportedFunctionsInPackage(t *testing.T) {
	// Define the bash command as a string
	bashCmd := `
find . -name "*.go" -not -path "./magefiles/*" -not -path "./v2/*" |
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
			goFuncs, err := goutils.FindExportedFunctionsInPackage(tc.packagePath)
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

func TestFindExportedFuncsWithoutTests(t *testing.T) {
	pkg := "bla"
	// Create temporary directory
	tempDir, err := os.MkdirTemp("/tmp", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file with exported function
	file1 := filepath.Join(tempDir, "file1.go")
	content1 := fmt.Sprintf(`package %s
func ExportedFunc1() {}
`, pkg)
	if err := os.WriteFile(file1, []byte(content1), 0666); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}

	// Create example with exported function and test function
	file2 := filepath.Join(tempDir, "file2.go")
	content2 := fmt.Sprintf(`package %s
func ExportedFunc2() {}
`, pkg)
	if err := os.WriteFile(file2, []byte(content2), 0666); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}

	file2Test := filepath.Join(tempDir, "file2_test.go")
	content2Test := fmt.Sprintf(`package %s
import "testing"
func TestExportedFunc2(t *testing.T) {}
`, pkg)
	if err := os.WriteFile(file2Test, []byte(content2Test), 0666); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Create a file with exported function and no test function
	file3 := filepath.Join(tempDir, "pkg", "bla", "file3.go")
	content3 := fmt.Sprintf(`package %s
func ExportedFunc3() {}
`, pkg)
	if err := os.MkdirAll(filepath.Dir(file3), os.ModePerm); err != nil {
		t.Fatalf("failed to create file3 dir: %v", err)
	}
	if err := os.WriteFile(file3, []byte(content3), 0666); err != nil {
		t.Fatalf("failed to create file3: %v", err)
	}

	// Create a file with exported function and test function
	file4 := filepath.Join(tempDir, "pkg", "bla", "file3_test.go")
	content4 := fmt.Sprintf(`package %s
import "testing"
func TestExportedFunc3(t *testing.T) {}
`, pkg)
	if err := os.WriteFile(file4, []byte(content4), 0666); err != nil {
		t.Fatalf("failed to create file4: %v", err)
	}

	// Call FindExportedFuncsWithoutTests
	exportedFuncs, err := fileutils.FindExportedFuncsWithoutTests(tempDir)
	if err != nil {
		t.Fatalf("failed to find exported funcs: %v", err)
	}

	// Assert the result
	expectedFuncs := []string{"ExportedFunc1"}
	if !reflect.DeepEqual(exportedFuncs, expectedFuncs) {
		t.Errorf("expected funcs: %v, got: %v", expectedFuncs, exportedFuncs)
	}
}
