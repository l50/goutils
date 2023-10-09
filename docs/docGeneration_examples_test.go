package docs_test

import (
	"fmt"
	"path/filepath"

	"github.com/l50/goutils/v2/docs"
	"github.com/spf13/afero"
)

func ExampleCreatePackageDocs() {
	// Mock the filesystem for testing
	fs := afero.NewMemMapFs()

	// Set up the repo details
	repo := docs.Repo{
		Owner: "l50",     // Repository owner's name.
		Name:  "goutils", // Repository's name.
	}

	// Set the path to the template file
	templatePath := filepath.Join("dev", "mage", "templates", "README.md.tmpl")

	// Set the packages to exclude (optional)
	excludedPkgs := []string{"excludedPkg1", "excludedPkg2"}

	if err := docs.CreatePackageDocs(fs, repo, templatePath, excludedPkgs...); err != nil {
		fmt.Printf("failed to create package docs: %v", err)
	}
}
