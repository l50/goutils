package docs_test

import (
	"path/filepath"
	"testing"

	"github.com/l50/goutils/v2/docs"
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
		name            string
		repo            docs.Repo
		templatePath    string
		setupFs         func() afero.Fs
		expectErr       bool
		expectedContent string
	}{
		{
			name: "valid template path",
			repo: repo,
			setupFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)
				_ = afero.WriteFile(fs, "go.mod", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				return fs
			},
			templatePath: templatePath,
			expectErr:    false,
		},
		{
			name: "invalid template path",
			repo: repo,
			setupFs: func() afero.Fs {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "template.tmpl", []byte("{{.PackageName}}"), 0644)
				_ = afero.WriteFile(fs, "go.mod", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
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
				_ = fs.MkdirAll("/Users/bob/path/to/go/project/pkg", 0755)
				_ = afero.WriteFile(fs, templatePath, []byte("{{.PackageName}}"), 0644)
				_ = afero.WriteFile(fs, "/Users/bob/path/to/go/project/pkg", []byte("module github.com/"+repo.Owner+"/"+repo.Name), 0644)
				return fs
			},
			templatePath: templatePath,
			expectErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := tc.setupFs()
			err := docs.CreatePackageDocs(fs, tc.repo, tc.templatePath)
			if (err != nil) != tc.expectErr {
				t.Errorf("CreatePackageDocs() error = %v, expectErr %v", err, tc.expectErr)
			}
		})
	}
}
