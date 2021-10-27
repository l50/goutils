SHELL := bash
PACKAGE := goutils
USER := l50
VERSION := v1.0.0

.PHONY: all
all: update-pkg-cache test

.PHONY: test
test:
	@go test -v -race ./... ; \
	echo "Done."

# Update the package found at https://pkg.go.dev so that we can use our changes in other projects
.PHONY: update-pkg-cache
update-pkg-cache:
	curl https://proxy.golang.org/github.com/$(USER)/$(PACKAGE)/@v/$(VERSION).info
