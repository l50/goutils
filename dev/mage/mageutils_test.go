package mageutils_test

import (
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

	mageutils "github.com/l50/goutils/v2/dev/mage"
	fileutils "github.com/l50/goutils/v2/file/fileutils"
	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"
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
	if _, err := os.Stat(repoRoot); err != nil && !os.IsNotExist(err) {
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
		file := fileutils.RealFile(dir)
		info, err := file.Stat()
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			log.Printf("failed to stat directory %s: %v", dir, err)
		}

		if info != nil {
			if err := sys.RmRf(file); err != nil {
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

func TestCompile(t *testing.T) {
	testCases := []struct {
		name      string
		buildPath string
		goOS      string
		goArch    string
		wantErr   bool
	}{
		{
			name:      "TestCompileValid",
			buildPath: "testoutput",
			goOS:      "linux",
			goArch:    "amd64",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := mageutils.Compile(tc.buildPath, tc.goOS, tc.goArch)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Compile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if _, err := os.Stat(tc.buildPath); os.IsNotExist(err) {
					t.Fatalf("Compile() did not create the output file: %s", tc.buildPath)
				}
			}
		})

		// Clean up the output file
		os.RemoveAll(tc.buildPath)
	}
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

func TestFindExportedFuncsWithoutTests(t *testing.T) {
	testCases := []struct {
		name           string
		sourceContent  string
		testContent    string
		expectedOutput []string
	}{
		{
			name: "Exported function without tests",
			sourceContent: `package main
                            func ExportedFunc1() {}`,
			testContent:    "",
			expectedOutput: []string{"ExportedFunc1"},
		},
		{
			name: "Exported function with tests",
			sourceContent: `package main
                            func ExportedFunc2() {}`,
			testContent: `package main
                          import "testing"
                          func TestExportedFunc2(t *testing.T) {}`,
			expectedOutput: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create source file with exported function
			sourceFile := filepath.Join(tempDir, "source.go")
			if err := os.WriteFile(sourceFile, []byte(tc.sourceContent), 0666); err != nil {
				t.Fatalf("failed to create source file: %v", err)
			}

			// Create test file if provided
			if tc.testContent != "" {
				testFile := filepath.Join(tempDir, "source_test.go")
				if err := os.WriteFile(testFile, []byte(tc.testContent), 0666); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			exportedFuncs, err := mageutils.FindExportedFuncsWithoutTests(tempDir)
			if err != nil {
				t.Fatalf("failed to find exported funcs: %v", err)
			}

			if !reflect.DeepEqual(exportedFuncs, tc.expectedOutput) {
				t.Errorf("expected funcs: %v, got: %v", tc.expectedOutput, exportedFuncs)
			}
		})
	}
}

func TestFindExportedFunctionsInPackage(t *testing.T) {
	testCases := []struct {
		name              string
		packageDir        string
		expectedFunctions []mageutils.FuncInfo
		expectErr         bool
	}{
		{
			name:       "Test with Keeper struct",
			packageDir: "../../pwmgr/keeper",
			expectedFunctions: []mageutils.FuncInfo{
				{
					FuncName: "Keeper.CommanderInstalled",
					FilePath: "keeperutils.go",
				},
				{
					FuncName: "Keeper.LoggedIn",
					FilePath: "keeperutils.go",
				},
				{
					FuncName: "Keeper.AddRecord",
					FilePath: "keeperutils.go",
				},
				{
					FuncName: "Keeper.RetrieveRecord",
					FilePath: "keeperutils.go",
				},
				{
					FuncName: "Keeper.SearchRecords",
					FilePath: "keeperutils.go",
				},
			},
			expectErr: false,
		},
		{
			name:              "Test with non-existing package",
			packageDir:        "./nonexistingpackage",
			expectedFunctions: []mageutils.FuncInfo{},
			expectErr:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := mageutils.FindExportedFunctionsInPackage(tc.packageDir)
			if (err != nil) != tc.expectErr {
				t.Errorf("FindExportedFunctionsInPackage() error = %v, expectErr %v", err, tc.expectErr)
				return
			}
			if len(result) != len(tc.expectedFunctions) {
				t.Fatalf("unexpected number of results: got %v, want %v", len(result), len(tc.expectedFunctions))
			}

			for i, r := range result {
				expected := tc.expectedFunctions[i]
				if r.FuncName != expected.FuncName || !strings.HasSuffix(r.FilePath, expected.FilePath) {
					t.Errorf("unexpected result: got %+v, want %+v", r, expected)
				}
			}
		})
	}
}
