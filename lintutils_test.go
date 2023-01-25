package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	lintCleanupFiles []string
)

func init() {
	randStr, _ := RandomString(8)
	testFile = fmt.Sprintf("/tmp/lintingutils-test-file-%s.md", randStr)
	lintCleanupFiles = append(lintCleanupFiles, testFile)
	// Create a markdown file in the tmp directory
	testFileContent = strings.ReplaceAll(`
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
}

func TestAddFencedCB(t *testing.T) {
	t.Cleanup(func() {
		cleanupLintUtils(t)
	})

	// Create the test file
	if err := CreateFile(testFile, []byte(testFileContent)); err != nil {
		t.Fatalf("error running CreateFile() with %s and %s: %v", testFile, testFileContent, err)
	}

	// Test that the AddFencedCB function correctly adds the language to code blocks
	if err := AddFencedCB(testFile, "js"); err != nil {
		t.Fatalf("failed to run AddFencedCB() with %s and js as inputs: %v", testFile, err)
	}

	// Read the modified file to check its contents
	modifiedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}

	expectedOutput := strings.ReplaceAll(`
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
	if !strings.Contains(expectedOutput, string(modifiedContent)) {
		t.Fatalf("error: %s does not resemble expected output: %s", string(modifiedContent), expectedOutput)
	}
}

func TestInstallGoPCDeps(t *testing.T) {
	if err := InstallGoPCDeps(); err != nil {
		t.Fatal(err)
	}
}

func TestInstallPCHooks(t *testing.T) {
	if err := InstallPCHooks(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePCHooks(t *testing.T) {
	if err := UpdatePCHooks(); err != nil {
		t.Fatal(err)
	}
}

func cleanupLintUtils(t *testing.T) {
	for _, dir := range lintCleanupFiles {
		if err := RmRf(dir); err != nil {
			fmt.Println("failed to clean up lintUtils: ", err.Error())
		}
	}
}

// Currently not running because this test
// will create a ton of processes that
// you have to pkill by hand. It's annoying
// as all hell.
// func TestRunPCHooks(t *testing.T) {
// 	seconds := 120
// 	timeout := time.After(time.Duration(seconds) * time.Second)
// 	errors := make(chan error)
// 	done := make(chan bool)
// 	go func() {
// 		if err := RunPCHooks(); err != nil {
// 			errors <- fmt.Errorf(
// 				"received error from RunPCHooks(): %v", err)
// 		}
// 		done <- true
// 	}()

// 	select {
// 	case <-timeout:
// 		fmt.Printf("timed out TestRunPCHooks() after %ds to "+
// 			"stop the test from infinitely calling itself",
// 			seconds)
// 	case <-done:
// 	case err := <-errors:
// 		t.Fatalf("test fail - TestRunPCHooks() "+
// 			"in precommitutils_test.go: %v", err)
// 	}
// }
