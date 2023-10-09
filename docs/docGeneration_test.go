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
	templatePath := filepath.Join("magefiles", "tmpl", "README.md.tmpl")

	tests := []struct {
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
				baseFs := afero.NewOsFs()
				templatePath, err := afero.TempDir(baseFs, "", "testDocs")
				if err != nil {
					t.Fatal(err)
				}
				fs := afero.NewBasePathFs(baseFs, templatePath)

				// Create magefiles and tmpl directories in in-memory FS
				_ = fs.MkdirAll(filepath.Join("magefiles", "tmpl"), 0755)
				// Write template file to in-memory FS
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)

				repoRoot, err := git.RepoRoot()
				if err != nil {
					t.Fatal(err)
				}

				if err := copyFileToMockFs(afero.NewOsFs(), fs, filepath.Join(repoRoot, "magefiles", "go.mod"), filepath.Join("magefiles", "go.mod")); err != nil {
					t.Fatal(err)
				}

				if err := copyFileToMockFs(afero.NewOsFs(), fs, filepath.Join(repoRoot, "magefiles", "go.sum"), filepath.Join("magefiles", "go.sum")); err != nil {
					t.Fatal(err)
				}

				// Here, we're copying the real magefile.go into the mock filesystem.
				if err := copyFileToMockFs(afero.NewOsFs(), fs, filepath.Join(repoRoot, "magefiles", "magefile.go"), filepath.Join("magefiles", "magefile.go")); err != nil {
					t.Fatal(err)
				}

				// Here, we're copying the real README.md.tmpl into the mock filesystem.
				if err := copyFileToMockFs(afero.NewOsFs(), fs, filepath.Join(repoRoot, "magefiles", "tmpl", "README.md.tmpl"), filepath.Join("magefiles", "tmpl", "README.md.tmpl")); err != nil {
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
				baseFs := afero.NewOsFs()
				tempDir, err := afero.TempDir(baseFs, "", "testDocs")
				if err != nil {
					t.Fatal(err)
				}
				fs := afero.NewBasePathFs(baseFs, tempDir)

				// Create magefiles and tmpl directories in in-memory FS
				_ = fs.MkdirAll(filepath.Join("magefiles", "tmpl"), 0755)
				// Write template file to in-memory FS
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)

				repoRoot, err := git.RepoRoot()
				if err != nil {
					t.Fatal(err)
				}

				// Here, we're copying the real magefile.go into the mock filesystem.
				if err := copyFileToMockFs(afero.NewOsFs(), fs, filepath.Join(repoRoot, "magefiles", "tmpl", "README.md.tmpl"), filepath.Join("magefiles", "tmpl", "README.md.tmpl")); err != nil {
					t.Fatal(err)
				}

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
		{
			name: "magefiles directory with main package",
			repo: repo,
			setupFs: func() afero.Fs {
				baseFs := afero.NewOsFs()
				tempDir, err := afero.TempDir(baseFs, "", "testDocs")
				if err != nil {
					t.Fatal(err)
				}
				fs := afero.NewBasePathFs(baseFs, tempDir)

				// Create magefiles and tmpl directories in in-memory FS
				_ = fs.MkdirAll(filepath.Join("magefiles", "tmpl"), 0755)
				// Write template file to in-memory FS
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)
				// Mock README.md.tmpl content
				_ = afero.WriteFile(fs, filepath.Join("magefiles", "tmpl", "README.md.tmpl"), []byte("mock README.md.tmpl content"), 0644)

				repoRoot, err := git.RepoRoot()
				if err != nil {
					t.Fatal(err)
				}

				// Here, we're copying the real magefile.go into the mock filesystem.
				if err := copyFileToMockFs(afero.NewOsFs(), fs, filepath.Join(repoRoot, "magefiles", "magefile.go"), filepath.Join("magefiles", "magefile.go")); err != nil {
					t.Fatal(err)
				}

				return fs
			},
			templatePath:  templatePath,
			expectErr:     false,
			expectPkgName: "magefiles",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := tc.setupFs()
			printFs(fs, ".", "")

			// Check directory explicitly:
			dirExists, _ := afero.DirExists(fs, "magefiles")
			if !dirExists {
				t.Error("magefiles directory does not exist in the in-memory file system")
			}
			if tc.expectPkgName != "" {
				readmePath := filepath.Join("magefiles", "README.md")
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

func printFs(fs afero.Fs, dir string, indent string) {
	entries, _ := afero.ReadDir(fs, dir)
	for _, entry := range entries {
		if entry.IsDir() {
			printFs(fs, filepath.Join(dir, entry.Name()), indent+"  ")
		}
	}
}

func copyFileToMockFs(srcFs afero.Fs, destFs afero.Fs, srcPath, destPath string) error {
	content, err := afero.ReadFile(srcFs, srcPath)
	if err != nil {
		return err
	}

	return afero.WriteFile(destFs, destPath, content, 0644)
}
