package docs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/docs"
	"github.com/l50/goutils/v2/git"
	"github.com/spf13/afero"
)

// test repository
var repo = docs.Repo{
	Owner: "owner",
	Name:  "name",
}

func TestCreatePackageDocs(t *testing.T) {
	templatePath := filepath.Join("templates", "README.md.tmpl")

	testCases := []struct {
		name          string
		repo          docs.Repo
		templatePath  string
		excludedPkgs  []string
		setupFs       func(t *testing.T) (afero.Fs, string)
		expectErr     bool
		expectPkgName string
		expectPkgDir  bool
	}{
		{
			name: "valid template path",
			repo: repo,
			setupFs: func(t *testing.T) (afero.Fs, string) {
				fs, templatePath := initBaseFs(t)
				initCommonDirs(fs, templatePath)
				_ = fs.MkdirAll("pkg", 0755)

				filesToCopy := []string{
					"templates/README.md.tmpl", "templates/README.md.tmpl",
				}
				if err := copyFilesToFs(fs, filesToCopy...); err != nil {
					t.Fatal(err)
				}

				return fs, "" // Return empty string for tempDir
			},
			templatePath: templatePath,
			expectErr:    false,
			expectPkgDir: true,
		},
		{
			name: "invalid template path",
			repo: repo,
			setupFs: func(t *testing.T) (afero.Fs, string) {
				fs, templatePath := initBaseFs(t)
				initCommonDirs(fs, templatePath)
				return fs, ""
			},
			templatePath: "nonexistent_template.tmpl",
			expectErr:    true,
		},
		{
			name: "path outside root directory",
			repo: repo,
			setupFs: func(t *testing.T) (afero.Fs, string) {
				fs := afero.NewMemMapFs()
				_ = fs.MkdirAll("/Users/bob/co/opensource/asdf/pkg", 0755)
				_ = afero.WriteFile(fs, templatePath, []byte("docs_test"), 0644)
				_ = afero.WriteFile(fs, "/Users/bob/co/opensource/asdf/pkg", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				return fs, ""
			},
			templatePath: templatePath,
			expectErr:    true,
		},
		{
			name: "absolute path given",
			repo: repo,
			setupFs: func(t *testing.T) (afero.Fs, string) {
				fs := afero.NewMemMapFs()
				_ = fs.MkdirAll("/Users/bob/co/opensource/asdf/pkg", 0755)
				_ = afero.WriteFile(fs, templatePath, []byte("docs_test"), 0644)
				_ = afero.WriteFile(fs, "/Users/bob/co/opensource/asdf/pkg", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				return fs, ""
			},
			templatePath: "/Users/bob/co/opensource/asdf/pkg",
			expectErr:    true,
		},
		{
			name:         "excluding specific packages",
			repo:         repo,
			excludedPkgs: []string{"excludedPkg1", "excludedPkg2"},
			setupFs: func(t *testing.T) (afero.Fs, string) {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, templatePath, []byte("docs_test"), 0644)
				_ = afero.WriteFile(fs, "go.mod", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				_ = fs.MkdirAll("excludedPkg1", 0755)
				_ = fs.MkdirAll("excludedPkg2", 0755)
				return fs, ""
			},
			templatePath: templatePath,
			expectErr:    false,
		},
		{
			name: "exclude test files",
			repo: repo,
			setupFs: func(t *testing.T) (afero.Fs, string) {
				fs := afero.NewOsFs()

				// Create a temporary directory
				tempDir, err := os.MkdirTemp("", "testDocs")
				if err != nil {
					t.Fatal(err)
				}

				// Ensure cleanup
				t.Cleanup(func() { os.RemoveAll(tempDir) })

				// Write files to the temp directory
				_ = os.MkdirAll(filepath.Join(tempDir, "templates"), 0755)
				_ = os.MkdirAll(filepath.Join(tempDir, "pkg"), 0755)
				_ = os.WriteFile(filepath.Join(tempDir, "templates", "README.md.tmpl"), []byte("{{range .Functions}}{{.Name}}{{end}}"), 0644)
				_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module pkg"), 0644)
				_ = os.WriteFile(filepath.Join(tempDir, "pkg", "regular.go"), []byte(`
package pkg

func RegularFunction() {}
`), 0644)
				_ = os.WriteFile(filepath.Join(tempDir, "pkg", "regular_test.go"), []byte(`
package pkg

func TestRegularFunction() {}
`), 0644)
				return fs, tempDir // Return tempDir
			},
			templatePath:  filepath.Join("templates", "README.md.tmpl"),
			expectErr:     false,
			expectPkgName: "RegularFunction",
			expectPkgDir:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs, tempDir := tc.setupFs(t)

			// Change working directory if tempDir is provided
			if tempDir != "" {
				changeWorkingDir(t, tempDir)
				defer revertWorkingDir(t)
			}

			debugFileSystem(fs, t, "Before CreatePackageDocs")

			checkTemplatesDir(fs, t)

			if tc.expectPkgDir {
				checkPkgDir(fs, t)
			}

			err := docs.CreatePackageDocs(fs, tc.repo, tc.templatePath, tc.excludedPkgs...)
			if (err != nil) != tc.expectErr {
				t.Errorf("CreatePackageDocs() error = %v, expectErr %v", err, tc.expectErr)
			}

			debugFileSystem(fs, t, "After CreatePackageDocs")

			if tc.expectPkgName != "" {
				verifyReadmeContent(fs, t, tc.expectPkgName)
			}
		})
	}
}

func initBaseFs(t *testing.T) (afero.Fs, string) {
	baseFs := afero.NewOsFs()
	templateDir, err := afero.TempDir(baseFs, "", "testDocs")
	if err != nil {
		t.Fatal(err)
	}

	return afero.NewBasePathFs(baseFs, templateDir), filepath.Join(templateDir, "templates", "README.md.tmpl")
}

func initCommonDirs(fs afero.Fs, path string) {
	_ = fs.MkdirAll("templates", 0755)
	_ = afero.WriteFile(fs, path, []byte("{{.PackageName}}"), 0644)
}

func copyFilesToFs(fs afero.Fs, files ...string) error {
	repoRoot, err := git.RepoRoot()
	if err != nil {
		return err
	}
	for i := 0; i < len(files); i += 2 {
		src := filepath.Join(repoRoot, files[i])
		dest := files[i+1]
		if err := copyFileToMockFs(afero.NewOsFs(), fs, src, dest); err != nil {
			return err
		}
	}
	return nil
}

func copyFileToMockFs(srcFs afero.Fs, destFs afero.Fs, srcPath, destPath string) error {
	content, err := afero.ReadFile(srcFs, srcPath)
	if err != nil {
		return err
	}

	return afero.WriteFile(destFs, destPath, content, 0644)
}

func changeWorkingDir(t *testing.T, dir string) {
	err := os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
}

func revertWorkingDir(t *testing.T) {
	err := os.Chdir("..")
	if err != nil {
		t.Fatalf("Failed to revert working directory: %v", err)
	}
}

func debugFileSystem(fs afero.Fs, t *testing.T, message string) {
	files, _ := afero.ReadDir(fs, ".")
	t.Logf("%s - Files and directories in root:", message)
	for _, f := range files {
		t.Logf("- %s (IsDir: %v)", f.Name(), f.IsDir())
	}
}

func checkTemplatesDir(fs afero.Fs, t *testing.T) {
	dirExists, _ := afero.DirExists(fs, "templates")
	if !dirExists {
		t.Error("templates directory does not exist in the in-memory file system")
	}
}

func checkPkgDir(fs afero.Fs, t *testing.T) {
	pkgDirExists, _ := afero.DirExists(fs, "pkg")
	if !pkgDirExists {
		t.Error("pkg directory does not exist in the in-memory file system")
	} else {
		// Debug: List files in the pkg directory
		pkgFiles, _ := afero.ReadDir(fs, "pkg")
		t.Logf("Files in pkg directory:")
		for _, f := range pkgFiles {
			t.Logf("- %s", f.Name())
		}
	}
}

func verifyReadmeContent(fs afero.Fs, t *testing.T, expectedPkgName string) {
	readmePath := filepath.Join("pkg", "README.md")
	content, err := afero.ReadFile(fs, readmePath)
	if err != nil {
		t.Errorf("Failed to read generated README: %v", err)
	} else {
		if !strings.Contains(string(content), expectedPkgName) {
			t.Errorf("Expected function name %s not found in generated README", expectedPkgName)
		}
		if strings.Contains(string(content), "Test"+expectedPkgName) {
			t.Errorf("Test function Test%s should not be in generated README", expectedPkgName)
		}
	}
}
