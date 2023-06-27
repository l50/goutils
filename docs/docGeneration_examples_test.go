package docs_test

import (
	"log"

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
	templatePath := "dev/mage/templates/README.md.tmpl"

	// Run the function
	if err := docs.CreatePackageDocs(fs, repo, templatePath); err != nil {
		log.Printf("failed to create package docs: %v", err)
	}
}