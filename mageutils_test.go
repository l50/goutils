package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	mageCleanupDirs []string
)

func init() {
	// Create test repo and queue it for cleanup
	randStr, _ := RandomString(8)
	clonePath = createTestRepo(fmt.Sprintf("mageutils-%s", randStr))
	mageCleanupDirs = append(mageCleanupDirs, clonePath)
}

func TestGHRelease(t *testing.T) {
	// Call the function with an old version
	newVer := "v1.0.0"
	if err := GHRelease(newVer); err == nil {
		t.Errorf("release %s should not have been created: %v", newVer, err)
	}
}

func TestGoReleaser(t *testing.T) {
	t.Cleanup(func() {
		cleanupMageUtils(t)
	})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	releaserDir := filepath.Join(cwd, "dist")
	mageCleanupDirs = append(mageCleanupDirs, releaserDir)

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

func cleanupMageUtils(t *testing.T) {
	for _, dir := range mageCleanupDirs {
		if err := RmRf(dir); err != nil {
			fmt.Println("failed to clean up mageUtils: ", err.Error())
		}
	}
}

func TestFindExportedFunctionsInPackage1(t *testing.T) {
	packagePath := "."

	exportedFuncs, err := FindExportedFunctionsInPackage(packagePath)
	if err != nil {
		t.Fatalf("failed to find exported functions in Go package '%s': %v", packagePath, err)
	}

	expected := []struct {
		filePath string
		funcName string
	}{
		{"logutils.go", "CreateLogFile"},
		{"sysutils.go", "CheckRoot"},
		{"sysutils.go", "Cd"},
		{"sysutils.go", "CmdExists"},
		{"sysutils.go", "Cp"},
		{"sysutils.go", "EnvVarSet"},
		{"sysutils.go", "GetHomeDir"},
		{"sysutils.go", "Gwd"},
		{"sysutils.go", "GetFutureTime"},
		{"sysutils.go", "IsDirEmpty"},
		{"sysutils.go", "RunCommand"},
		{"sysutils.go", "RunCommandWithTimeout"},
		{"ansibleutils.go", "AnsiblePing"},
		{"cloudflareutils.go", "GetDNSRecords"},
		{"fileutils.go", "AppendToFile"},
		{"fileutils.go", "CreateEmptyFile"},
		{"fileutils.go", "CreateFile"},
		{"fileutils.go", "CreateDirectory"},
		{"fileutils.go", "CSVToLines"},
		{"fileutils.go", "DeleteFile"},
		{"fileutils.go", "FileExists"},
		{"fileutils.go", "FileToSlice"},
		{"fileutils.go", "FindFile"},
		{"fileutils.go", "ListFilesR"},
		{"fileutils.go", "StringInFile"},
		{"fileutils.go", "RmRf"},
		{"macosutils.go", "InstallBrewDeps"},
		{"macosutils.go", "InstallBrewTFDeps"},
		{"mageutils.go", "GHRelease"},
		{"mageutils.go", "GoReleaser"},
		{"mageutils.go", "InstallVSCodeModules"},
		{"mageutils.go", "ModUpdate"},
		{"mageutils.go", "Tidy"},
		{"mageutils.go", "UpdateMageDeps"},
		{"mageutils.go", "InstallGoDeps"},
		{"mageutils.go", "FindExportedFunctionsInPackage"},
		{"netutils.go", "PublicIP"},
		{"netutils.go", "DownloadFile"},
		{"stringutils.go", "RandomString"},
		{"stringutils.go", "StringInSlice"},
		{"stringutils.go", "StringToInt64"},
		{"stringutils.go", "StringToSlice"},
		{"gitutils.go", "GetSSHPubKey"},
		{"gitutils.go", "AddFile"},
		{"gitutils.go", "Commit"},
		{"gitutils.go", "CloneRepo"},
		{"gitutils.go", "GetTags"},
		{"gitutils.go", "GetGlobalUserCfg"},
		{"gitutils.go", "CreateTag"},
		{"gitutils.go", "Push"},
		{"gitutils.go", "PushTag"},
		{"gitutils.go", "DeleteTag"},
		{"gitutils.go", "DeletePushedTag"},
		{"gitutils.go", "RepoRoot"},
		{"keeperutils.go", "CommanderInstalled"},
		{"keeperutils.go", "KeeperLoggedIn"},
		{"keeperutils.go", "RetrieveKeeperPW"},
		{"keeperutils.go", "SearchKeeperRecords"},
		{"lintutils.go", "InstallGoPCDeps"},
		{"lintutils.go", "InstallPCHooks"},
		{"lintutils.go", "UpdatePCHooks"},
		{"lintutils.go", "ClearPCCache"},
		{"lintutils.go", "RunPCHooks"},
		{"lintutils.go", "AddFencedCB"},
	}

	for _, exp := range expected {
		found := false
		for _, act := range exportedFuncs {
			if act.FilePath == exp.filePath && act.FuncName == exp.funcName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected function '%s' in file '%s', but it was not found in package '%s'", exp.funcName, exp.filePath, packagePath)
		}
	}
}

func TestFindExportedFunctionsInPackage2(t *testing.T) {
	packagePath := "."
	bashFile := "mytestyay.sh"

	// Write the bash command to a file
	err := os.WriteFile(bashFile, []byte(`
		get_exported_go_funcs () {
			find . -name "*.go" -not -path "./magefiles/*" | xargs grep -E -o 'func [A-Z][a-zA-Z0-9_]+\(' | grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn,.idea,.tox} -v '_test.go' | grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn,.idea,.tox} -v -E 'func [A-Z][a-zA-Z0-9_]+Test\(' | sed -e 's/func //' -e 's/(//' | awk -F: '{printf "Function: %s\nFile: %s\n", $2, $1}'
		}

		get_exported_go_funcs
	`), 0755)
	if err != nil {
		t.Fatalf("failed to write bash command to file: %v", err)
	}

	// Run the bash script and get its output
	output, err := RunCommand("bash", "./fuck.sh")
	if err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}

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

	// Get the exported functions from the package
	goFuncs, err := FindExportedFunctionsInPackage(packagePath)
	if err != nil {
		t.Fatalf("failed to find exported functions in Go package '%s': %v", packagePath, err)
	}

	// Compare the expected and actual functions
	for _, act := range goFuncs {
		if _, ok := bashFuncs[act.FuncName]; !ok {
			t.Errorf("unexpected function '%s' in file '%s' found in package '%s'", act.FuncName, act.FilePath, packagePath)
		}
	}

	missingFuncs := []string{}

	for bf := range bashFuncs {
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
		t.Errorf("go and bash implementations don't agree: %v", err)
	}

	if err := os.Remove(bashFile); err != nil {
		t.Errorf("failed to remove created bash script %s:%v", bashFile, err)
	}
}
