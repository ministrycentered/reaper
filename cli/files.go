package main

const (
	fileNewMakefile = `# Reaper CLI Makefile
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: all
all: clean test

.PHONY: clean
clean:
	cd $(ROOT_DIR) && go clean

.PHONY: test
test:
	cd $(ROOT_DIR) && go test -v ./...

.PHONY: build
build:
	cd $(ROOT_DIR) && go build -o $(ROOT_DIR)/build/{{.Name}} .
`
	fileNewGitignore = `/build
`

	fileNewModule = `module github.com/{{.Github}}/{{.Name}}

require github.com/maddiesch/reaper v{{.ReaperVersion}}
`

	fileNewMain = `package main

import (
	"os"

	"github.com/maddiesch/reaper"
)

func main() {
	app := reaper.NewApp("test-cli")

	app.Command("test", func(c *reaper.Context) error { return nil })

	err := app.Execute(os.Args[1:])
	if err != nil {
		reaper.Fatal(err)
	}
}
`
)
