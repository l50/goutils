package mageutils_test

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	mageutils "github.com/l50/goutils/dev/mage"
	"github.com/l50/goutils/git"
	"github.com/l50/goutils/str"
	"github.com/l50/goutils/sys"
	"github.com/stretchr/testify/mock"
)

var (
	mageCleanupArtifacts []string
	testingPath          string
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Printf("setup failed: %v", err)
		os.Exit(1)
	}
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() error {
	randStr, err := str.GenRandom(8)
	if err != nil {
		return fmt.Errorf("failed to generate random string: %v", err)
	}

	testingPath = createTestRepo(fmt.Sprintf("mageutils-%s", randStr))
	cwd := sys.Gwd()
	if err := sys.Cd(testingPath); err != nil {
		return fmt.Errorf("failed to change directory to %s: %v", testingPath, err)
	}
	// obtain the repository root
	repoRoot, err := git.RepoRoot()
	if err != nil {
		return fmt.Errorf("unable to find git root: %w", err)
	}

	// define the source directory
	magefileSrc := filepath.Join(repoRoot, "magefiles")

	// ensure the source directory exists
	if _, err := os.Stat(magefileSrc); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("magefiles directory does not exist at %s", magefileSrc)
	} else if err != nil {
		return fmt.Errorf("unexpected error when checking for magefiles directory: %w", err)
	}

	// define the test directory
	testDir := filepath.Join(repoRoot, "test")

	// create the test directory if it doesn't already exist
	if _, err = os.Stat(testDir); err != nil && !os.IsNotExist(err) {
		if err := os.Mkdir(testDir, os.ModePerm); err != nil {
			return fmt.Errorf("unable to create test directory: %w", err)
		}
	}

	if err := sys.Cd(cwd); err != nil {
		return fmt.Errorf("failed to change directory to %s: %v", testingPath, err)
	}

	mageCleanupArtifacts = append(mageCleanupArtifacts, testingPath)

	return nil
}

func teardown() {
	for _, dir := range mageCleanupArtifacts {
		info, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			log.Printf("failed to stat directory %s: %v", dir, err)
		}

		if info != nil {
			if err := sys.RmRf(dir); err != nil {
				log.Printf("failed to clean up directory %s: %v", dir, err)
			}
		}
	}
}

func createTestRepo(name string) string {
	cloneDir := "/tmp"
	currentTime := time.Now()
	targetPath := filepath.Join(
		cloneDir, fmt.Sprintf(
			"%s-%s", name, currentTime.Format("2006-01-02-15-04-05"),
		),
	)

	testRepoURL := "https://github.com/l50/goutils.git"
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
	testCases := []struct {
		desc    string
		version string
		wantErr bool
	}{
		{
			desc:    "Empty Version",
			version: "",
			wantErr: true,
		},
		{
			desc:    "Old Version",
			version: "v1.0.0",
			wantErr: true,
		},
	}
	mageCleanupArtifacts = append(mageCleanupArtifacts, "CHANGELOG.md")

	if err := sys.Cd(testingPath); err != nil {
		t.Errorf("failed to change directory to %s: %v", testingPath, err)
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := mageutils.GHRelease(tc.version)
			if (err != nil) != tc.wantErr {
				t.Errorf("GHRelease(%v) = error %v, wantErr %v", tc.version, err, tc.wantErr)
			}
		})
	}
}

func TestGoReleaser(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{"releases to the dist directory"},
	}
	if err := sys.Cd(testingPath); err != nil {
		t.Errorf("failed to change directory to %s: %v", testingPath, err)
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Get repo root
			repoRoot, err := git.RepoRoot()
			if err != nil {
				t.Errorf("failed to get repo root: %v", err)
				return
			}

			// Change into repo root
			if err := sys.Cd(repoRoot); err != nil {
				t.Errorf("failed to change directory to %s: %v", testingPath, err)
			}

			releaserDir := filepath.Join(repoRoot, "dist")

			mageCleanupArtifacts = append(mageCleanupArtifacts, releaserDir)

			if err := mageutils.GoReleaser(); err != nil {
				t.Errorf("GoReleaser() failed with error %v", err)
			}
		})
	}
}

func TestInstallVSCodeModules(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{"installs VSCode modules"},
	}
	if err := sys.Cd(testingPath); err != nil {
		t.Errorf("failed to change directory to %s: %v", testingPath, err)
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			if err := mageutils.InstallVSCodeModules(); err != nil {
				t.Errorf("InstallVSCodeModules() failed with error %v", err)
			}
		})
	}
}

func TestModUpdate(t *testing.T) {
	testCases := []struct {
		desc      string
		recursive bool
		verbose   bool
	}{
		{
			desc:      "non-recursive verbose update",
			recursive: false,
			verbose:   true,
		},
	}

	if err := sys.Cd(testingPath); err != nil {
		t.Errorf("failed to change directory to %s: %v", testingPath, err)
	}
	for _, tc := range testCases {
		tc := tc // rebind the variable
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// create temporary directory for a mock Go module
			dir, err := os.MkdirTemp("", "modupdate")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			if err := sys.Cd(dir); err != nil {
				t.Errorf("failed to change directory to %s: %v", testingPath, err)
			}

			// create a dummy go file
			goFile := filepath.Join(dir, "main.go")
			fileContent := "package main\n\nfunc main() {}\n"
			err = os.WriteFile(goFile, []byte(fileContent), 0666)
			if err != nil {
				t.Fatal(err)
			}

			// create go.mod file
			cmd := exec.Command("go", "mod", "init", "modupdate")
			cmd.Dir = dir
			err = cmd.Run()
			if err != nil {
				t.Fatal(err)
			}

			if err := mageutils.ModUpdate(tc.recursive, tc.verbose); err != nil {
				t.Errorf("ModUpdate(%v, %v) = error %v, want no error", tc.recursive, tc.verbose, err)
			}
		})
	}
}

type MockDev struct {
	mock.Mock
}

func (m *MockDev) Tidy() error {
	args := m.Called()
	return args.Error(0)
}

func TestTidy(t *testing.T) {
	testCases := []struct {
		desc       string
		mockError  error
		expectFail bool
	}{
		{"tidies the module", nil, false},
		{"tidies the module with error", errors.New("some error"), true},
	}
	if err := sys.Cd(testingPath); err != nil {
		t.Errorf("failed to change directory to %s: %v", testingPath, err)
	}

	for _, tc := range testCases {
		tc := tc // rebind the variable

		t.Run(tc.desc, func(t *testing.T) {

			t.Parallel()

			mockDev := new(MockDev)
			mockDev.On("Tidy").Return(tc.mockError)

			if err := mockDev.Tidy(); (err != nil) != tc.expectFail {
				t.Errorf("Tidy() returned error %v, expectFail: %v", err, tc.expectFail)
			}

			mockDev.AssertExpectations(t)
		})
	}
}

func TestUpdateMageDeps(t *testing.T) {
	testCases := []struct {
		desc    string
		mageDir string
		wantErr bool
	}{
		{
			desc:    "non-existent directory",
			mageDir: "non-existent",
			wantErr: true,
		},
		{
			desc:    "updates mage dependencies",
			mageDir: "magefiles",
			wantErr: false,
		},
	}
	if err := sys.Cd(testingPath); err != nil {
		t.Errorf("failed to change directory to %s: %v", testingPath, err)
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repoRoot, err := git.RepoRoot()
			if err != nil {
				t.Errorf("RepoRoot() failed with error %v", err)
				return
			}

			if err := sys.Cd(repoRoot); err != nil {
				t.Errorf("failed to change directory to %s: %v", testingPath, err)
			}

			if err := mageutils.UpdateMageDeps(tc.mageDir); (err != nil) != tc.wantErr {
				t.Errorf("UpdateMageDeps(%s) failed with error %v, wantErr %v", tc.mageDir, err, tc.wantErr)
			}
		})
	}
}

func TestInstallGoDeps(t *testing.T) {
	testCases := []struct {
		desc string
		deps []string
	}{
		{
			desc: "installs go dependencies",
			deps: []string{
				"golang.org/x/lint/golint",
				"golang.org/x/tools/cmd/goimports",
			},
		},
	}
	if err := sys.Cd(testingPath); err != nil {
		t.Errorf("failed to change directory to %s: %v", testingPath, err)
	}

	for _, tc := range testCases {
		tc := tc // rebind the variable
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			if err := mageutils.InstallGoDeps(tc.deps); err != nil {
				t.Errorf("InstallGoDeps(%v) failed with error %v", tc.deps, err)
			}
		})
	}
}

func TestFindExportedFunctionsInPackage(t *testing.T) {
	bashCmd := `
find . -name "*.go" -not -path "./magefiles/*" |
xargs grep -E -o 'func [A-Z][a-zA-Z0-9_]+\(' |
grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn,.idea,.tox} -v '_test.go' |
grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn,.idea,.tox} -v -E 'func [A-Z][a-zA-Z0-9_]+Test\(' |
sed -e 's/func //' -e 's/(//' |
awk -F: '{printf "Function: %s\nFile: %s\n", $2, $1}'`

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
			packagePath:    testingPath,
			expectedFuncs:  bashFuncs,
			expectedErrors: false,
		},
		{
			name:           "Invalid package",
			packagePath:    "/tmp",
			expectedFuncs:  nil,
			expectedErrors: true,
		},
	}
	// Loop through the test cases and execute each one
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Get the exported functions from the package
			goFuncs, err := mageutils.FindExportedFunctionsInPackage(tc.packagePath)
			if tc.expectedErrors && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tc.expectedErrors && err != nil {
				t.Logf("CURRENT DIRECTORY: %v", sys.Gwd())
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
				fmt.Println("Missing functions: ", missingFuncs)
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
	exportedFuncs, err := mageutils.FindExportedFuncsWithoutTests(tempDir)
	if err != nil {
		t.Fatalf("failed to find exported funcs: %v", err)
	}

	// Assert the result
	expectedFuncs := []string{"ExportedFunc1"}
	if !reflect.DeepEqual(exportedFuncs, expectedFuncs) {
		t.Errorf("expected funcs: %v, got: %v", expectedFuncs, exportedFuncs)
	}
}
