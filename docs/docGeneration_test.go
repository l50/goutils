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
		setupFs       func() (afero.Fs, string)
		expectErr     bool
		expectPkgName string
		expectPkgDir  bool
	}{
		{
			name: "valid template path",
			repo: repo,
			setupFs: func() (afero.Fs, string) {
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
			setupFs: func() (afero.Fs, string) {
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
			setupFs: func() (afero.Fs, string) {
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
			setupFs: func() (afero.Fs, string) {
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
			setupFs: func() (afero.Fs, string) {
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
			setupFs: func() (afero.Fs, string) {
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
			fs, tempDir := tc.setupFs()

			// If tempDir is not empty, change working directory
			if tempDir != "" {
				err := os.Chdir(tempDir)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					_ = os.Chdir("..") // Return to previous directory after test
				}()
			}

			// Debug: List all files and directories in the root
			files, _ := afero.ReadDir(fs, ".")
			t.Logf("Files and directories in root:")
			for _, f := range files {
				t.Logf("- %s (IsDir: %v)", f.Name(), f.IsDir())
			}

			// Check if templates directory exists
			dirExists, _ := afero.DirExists(fs, "templates")
			if !dirExists {
				t.Error("templates directory does not exist in the in-memory file system")
			}
			// Check if pkg directory exists (only if it's expected)
			if tc.expectPkgDir {
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

			err := docs.CreatePackageDocs(fs, tc.repo, tc.templatePath, tc.excludedPkgs...)
			if (err != nil) != tc.expectErr {
				t.Errorf("CreatePackageDocs() error = %v, expectErr %v", err, tc.expectErr)
			}

			// Debug: List all files and directories in the root after CreatePackageDocs
			files, _ = afero.ReadDir(fs, ".")
			t.Logf("Files and directories in root after CreatePackageDocs:")
			for _, f := range files {
				t.Logf("- %s (IsDir: %v)", f.Name(), f.IsDir())
			}

			if tc.expectPkgName != "" {
				readmePath := filepath.Join("pkg", "README.md")
				content, err := afero.ReadFile(fs, readmePath)
				if err != nil {
					t.Errorf("Failed to read generated README: %v", err)
					// Debug: Check if pkg directory still exists
					pkgDirExists, _ := afero.DirExists(fs, "pkg")
					t.Logf("pkg directory exists after CreatePackageDocs: %v", pkgDirExists)
					if pkgDirExists {
						// List files in pkg directory
						pkgFiles, _ := afero.ReadDir(fs, "pkg")
						t.Logf("Files in pkg directory after CreatePackageDocs:")
						for _, f := range pkgFiles {
							t.Logf("- %s", f.Name())
						}
					}
				} else {
					if !strings.Contains(string(content), tc.expectPkgName) {
						t.Errorf("expected function name %s not found in generated README", tc.expectPkgName)
					}
					if strings.Contains(string(content), "TestRegularFunction") {
						t.Errorf("test function TestRegularFunction should not be in generated README")
					}
				}
			}
		})
	}
}

func initBaseFs(t *testing.T) (afero.Fs, string) {
	baseFs := afero.NewOsFs()
	templatePath, err := afero.TempDir(baseFs, "", "testDocs")
	if err != nil {
		t.Fatal(err)
	}

	return afero.NewBasePathFs(baseFs, templatePath), templatePath
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
