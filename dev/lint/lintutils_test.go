package lint_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	lint "github.com/l50/goutils/v2/dev/lint"
	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/sys"
)

func TestAddFencedCB(t *testing.T) {
	testFileContent := strings.ReplaceAll(`
Get version of mongo:

”””
db.version()
”””

## Get json dump of the data
Create ”export.js” with the following:

”””
profileData = db.<collection name>.find();
// for example: db.system.users
while(profileData.hasNext()) {
	printjson(profileData.next());
}
”””
`, "”", "`")

	desiredOutput := strings.ReplaceAll(`
Get version of mongo:

”””js
db.version()
”””

## Get json dump of the data
Create ”export.js” with the following:

”””js
profileData = db.<collection name>.find();
// for example: db.system.users
while(profileData.hasNext()) {
	printjson(profileData.next());
}
”””
`, "”", "`")

	testCases := []struct {
		name           string
		inputFilePath  string
		inputLanguage  string
		expectedOutput string
	}{
		{
			name:           "Add JS language to code blocks",
			inputFilePath:  "lintingutils-test-file-abc.md",
			inputLanguage:  "js",
			expectedOutput: desiredOutput,
		},
	}

	repoRoot, err := git.RepoRoot()
	if err != nil {
		t.Fatalf("failed to get repo root: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "goutils-test-*")
			if err != nil {
				t.Fatalf("failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir) // cleanup

			// Copy the repo to temporary directory
			if err := sys.Cp(repoRoot, tempDir); err != nil {
				t.Fatalf("failed to copy repo to temp directory: %v", err)
			}

			// Prepare test file
			tc.inputFilePath = filepath.Join(tempDir, tc.inputFilePath)
			if err := os.WriteFile(tc.inputFilePath, []byte(testFileContent), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			// Test the AddFencedCB function
			if err := lint.AddFencedCB(tc.inputFilePath, tc.inputLanguage); err != nil {
				t.Fatalf("failed to run AddFencedCB(): %v", err)
			}

			// Read the modified file and check its contents
			modifiedContent, err := os.ReadFile(tc.inputFilePath)
			if err != nil {
				t.Fatalf("failed to read modified file: %v", err)
			}

			if strings.TrimSuffix(string(modifiedContent), "\n") != strings.TrimSuffix(tc.expectedOutput, "\n") {
				t.Errorf("Output does not match expected result: \nGot: %q\nExpected: %q", string(modifiedContent), tc.expectedOutput)
			}
		})
	}
}

func TestLintUtils(t *testing.T) {
	repoRoot, err := git.RepoRoot()
	if err != nil {
		t.Fatalf("failed to get repo root: %v", err)
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
	if err := sys.Cp(repoRoot, tempDir); err != nil {
		t.Fatalf("failed to copy repo to temp directory: %v", err)
	}

	testCases := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "TestInstallGoPCDeps",
			test: func(t *testing.T) {
				if err := lint.InstallGoPCDeps(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "TestInstallPCHooks",
			test: func(t *testing.T) {
				if err := lint.InstallPCHooks(); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "TestUpdatePCHooks",
			test: func(t *testing.T) {
				if err := lint.UpdatePCHooks(); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}

func TestRunPCHooks(t *testing.T) {
	testCases := []struct {
		name    string
		timeout []int
		wantErr bool
	}{
		{
			name:    "with short timeout value that stops running the PC hooks early and returns an error",
			timeout: []int{10},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := lint.RunPCHooks(tc.timeout...)
			if (err != nil) != tc.wantErr {
				t.Errorf("RunPCHooks() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRunHookTool(t *testing.T) {
	testCases := []struct {
		name    string
		hook    string
		files   []string
		wantErr bool
	}{
		{
			name:    "RunHookWithSingleFile",
			hook:    "markdownlint",
			files:   []string{"../README.md"},
			wantErr: false,
		},
		{
			name:    "RunHookWithMultipleFiles",
			hook:    "markdownlint",
			files:   []string{"../README.md", "../logging/README.md"},
			wantErr: false,
		},
		{
			name:    "RunHookWithNoFiles",
			hook:    "prettier",
			files:   nil,
			wantErr: false,
		},
		{
			name:    "RunHookWithError",
			hook:    "failing-hook",
			files:   []string{"file1.txt"},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := lint.RunHookTool(tc.hook, tc.files...)
			if (err != nil) != tc.wantErr {
				t.Errorf("RunHookTool(%s, %s) error = %v, wantErr %v",
					tc.hook, strings.Join(tc.files, ", "), err, tc.wantErr)
			}
		})
	}
}
