# Basic Makefile for cmdcraft project

# The binary to build (just the basename).
BIN := build/cmdcraft

# OS and architecture for cross-compilation
GOOS ?= linux
GOARCH ?= amd64

# Linker flags
LDFLAGS := -ldflags "-s -w"

.PHONY: all build clean

all: build

# Builds the project
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BIN) cmd/main.go

# Cleans our project: deletes binaries
clean:
	rm -f $(BIN)
