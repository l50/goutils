package utils

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	tmpLintingFile = "/tmp/lintingutils-test-file.md"
)

func init() {
	// Create a markdown file in the tmp directory
	fileContent := strings.ReplaceAll(`
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

	if err := os.WriteFile(tmpLintingFile, []byte(fileContent), 0644); err != nil {
		log.Printf("Error writing file: %v", err)
	}
}

func TestAddFencedCB(t *testing.T) {
	// Remove the temporary file after the test completes.
	defer tearDown()

	// Test that the fixCodeBlocks function correctly adds the language to code blocks
	if err := AddFencedCB(tmpLintingFile, "js"); err != nil {
		t.Errorf("Error occurred: %s", err)
	}

	// Read the modified file to check its contents
	modifiedContent, err := os.ReadFile(tmpLintingFile)
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
	assert.Equal(t, []byte(expectedOutput), modifiedContent)
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

func tearDown() {
	// remove the temporary file created during the test
	os.Remove(tmpLintingFile)
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
