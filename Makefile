ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: all
all: clean test

.PHONY: clean
clean:
	-rm $(ROOT_DIR)/example/example
	cd $(ROOT_DIR) && go clean

.PHONY: test
test:
	cd $(ROOT_DIR) && go test -v ./...

.PHONY: build
build: clean
	cd $(ROOT_DIR)/example && go build -o $(ROOT_DIR)/example/example .
	cd $(ROOT_DIR)/cli && go build -o $(ROOT_DIR)/reaper .

.PHONY: run
run: build
	$(ROOT_DIR)/example/example greet -greeting Hello
	$(ROOT_DIR)/example/example get -port 442 private/information sup3rsekret
	$(ROOT_DIR)/example/example get public/information
	$(ROOT_DIR)/example/example help
