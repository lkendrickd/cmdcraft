# Basic Makefile for cmdcraft project

# The binary to build (just the basename).
BIN := build/cmdcraft

.PHONY: all build clean

all: build

# Builds the project
build:
	go build -o $(BIN) cmd/main.go

# Cleans our project: deletes binaries
clean:
	rm -f $(BIN)
