# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MAIN_PACKAGE_PATH := ./cmd
BINARY_NAME := graft

# Version information
VERSION=1.0.0
APP_NAME := goraft
CURRENT_DIR := $(shell pwd)
INSTALL_DIR := /usr/local/bin

# Docker params
DOCKER_USER=animus123
PLATFORMS := linux/amd64,linux/arm64

# Docker commands
.PHONY: docker-build-local
docker-build-local:
	docker buildx build \
		--platform linux/amd64 \
		--tag $(APP_NAME):$(VERSION) \
		--load \
		.

.PHONY: docker-push
docker-push:
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag ${DOCKER_USER}/$(APP_NAME):$(VERSION) \
		--push \
		.
# Build command
.PHONY: build
build:
	@echo $(shell pwd)/bin
	$(GOBUILD) -o=$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

# Clean command
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Test command
.PHONY: test
test:
	$(GOTEST) -v ./...

# Install dependencies
.PHONY: deps
deps:
	$(GOGET) ./...

# Release command
.PHONY: release
release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe -v

# Help command
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the Go application"
	@echo "  clean       - Clean up generated files"
	@echo "  test        - Run tests"
	@echo "  deps        - Install dependencies"
	@echo "  release     - Create release builds for Linux, macOS, and Windows"

# Install command
.PHONY: install
install:
	sudo chmod +x ${INSTALL_DIR}/
	sudo cp $(CURRENT_DIR)/${APP_NAME} ${INSTALL_DIR}

# Uninstall command
.PHONY: uninstall
uninstall:
	sudo rm ${INSTALL_DIR}/${APP_NAME}