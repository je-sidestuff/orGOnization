# Basic Makefile for Golang project
# Includes GRPC Gateway, Protocol Buffers
SERVICE		?= $(shell basename `go list`)
VERSION		?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(PWD)/.version 2> /dev/null || echo v0)
PACKAGE		?= $(shell go list)
PACKAGES	?= $(shell go list ./...)
FILES		?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# App dir
APP_DIR=cmd

# Binary base name
BINARY_NAME=orgonization

# Output directory
BUILD_DIR=bin

# Platforms to build for
PLATFORMS=\
	linux/amd64 \
	windows/amd64 \
	darwin/amd64

.PHONY: help clean fmt lint vet test test-cover generate-grpc build build-docker all

default: help

help:       ## Show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

fmt:        ## Format the go source files
	go fmt ./...

lint:       ## Run go lint on the source files
	golint $(PACKAGES)

vet:        ## Run go vet on the source files
	go vet ./...

test:        ## Run go test -v -cover on the source files
	go test -v -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ${PACKAGES}

# To be included in a later increment.
#.PHONY: all
#all: build

.PHONY: build
build:       ## Build for current OS
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./${APP_DIR}


.PHONY: build-all
build-all:   ## Cross-compile for multiple platforms
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform##*/}; \
		OUTPUT=$(BUILD_DIR)/$(BINARY_NAME)-$${OS}-$${ARCH}; \
		if [ "$${OS}" = "windows" ]; then OUTPUT=$${OUTPUT}.exe; fi; \
		echo "Building $${OUTPUT}..."; \
		GOOS=$${OS} GOARCH=$${ARCH} go build -o $${OUTPUT} ./${APP_DIR} || exit 1; \
	done
