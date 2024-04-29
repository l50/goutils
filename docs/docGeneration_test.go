package docs_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/docs"
	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/sys"
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
		setupFs       func() afero.Fs
		expectErr     bool
		expectPkgName string
	}{
		{
			name: "valid template path",
			repo: repo,
			setupFs: func() afero.Fs {
				fs, templatePath := initBaseFs(t)
				initCommonDirs(fs, templatePath)

				filesToCopy := []string{
					"templates/README.md.tmpl", "templates/README.md.tmpl",
				}
				if err := copyFilesToFs(fs, filesToCopy...); err != nil {
					t.Fatal(err)
				}

				if err := sys.Cd(templatePath); err != nil {
					t.Fatal(err)
				}

				if _, err := sys.RunCommand("git", "init"); err != nil {
					t.Fatal(err)
				}

				return fs
			},
			templatePath: templatePath,
			expectErr:    false,
		},
		{
			name: "invalid template path",
			repo: repo,
			setupFs: func() afero.Fs {
				fs, templatePath := initBaseFs(t)
				initCommonDirs(fs, templatePath)

				return fs
			},
			templatePath: "nonexistent_template.tmpl",
			expectErr:    true,
		},
		{
			name: "path outside root directory",
			repo: repo,
			setupFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = fs.MkdirAll("/Users/bob/co/opensource/asdf/pkg", 0755)
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)
				_ = afero.WriteFile(fs, "/Users/bob/co/opensource/asdf/pkg", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				return fs
			},
			templatePath: templatePath,
			expectErr:    true,
		},
		{
			name: "absolute path given",
			repo: repo,
			setupFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = fs.MkdirAll("/Users/bob/co/opensource/asdf/pkg", 0755)
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)
				_ = afero.WriteFile(fs, "/Users/bob/co/opensource/asdf/pkg", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				return fs
			},
			templatePath: "/Users/bob/co/opensource/asdf/pkg",
			expectErr:    true,
		},
		{
			name:         "excluding specific packages",
			repo:         repo,
			excludedPkgs: []string{"excludedPkg1", "excludedPkg2"},
			setupFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)
				_ = afero.WriteFile(fs, "go.mod", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				// Write some data to simulate the packages being present
				_ = fs.MkdirAll("excludedPkg1", 0755)
				_ = fs.MkdirAll("excludedPkg2", 0755)
				return fs
			},
			templatePath: templatePath,
			expectErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := tc.setupFs()

			// Check directory explicitly:
			dirExists, _ := afero.DirExists(fs, "templates")
			if !dirExists {
				t.Error("magefiles directory does not exist in the in-memory file system")
			}

			if tc.expectPkgName != "" {
				readmePath := filepath.Join("templates", "README.md")
				content, err := afero.ReadFile(fs, readmePath)
				if err == nil && !strings.Contains(string(content), tc.expectPkgName) {
					t.Errorf("expected package name %s not found in generated README", tc.expectPkgName)
				}
			}

			err := docs.CreatePackageDocs(fs, tc.repo, tc.templatePath, tc.excludedPkgs...)
			if (err != nil) != tc.expectErr {
				t.Errorf("CreatePackageDocs() error = %v, expectErr %v", err, tc.expectErr)
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
