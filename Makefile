.DEFAULT_GOAL := refresh

APP_NAME = jamesdb
GOBIN_F = $(GOPATH)/bin/

ifeq ($(OS),Windows_NT)
	SHELL = cmd.exe
	RM = del
	IGNORE = > nul 2> nul
	BINARY_NAME = $(APP_NAME).exe
	PLATFORM = windows
	GOBIN = $(subst /,\,$(GOBIN_F))
else
	SHELL = /bin/sh
	RM = rm -f
	IGNORE = 2>/dev/null
	BINARY_NAME = $(APP_NAME)
	PLATFORM = $(shell uname -s | tr '[:upper:]' '[:lower:]')
	GOBIN = $(GOBIN_F)
endif

SRC_DIR = .
LDFLAGS = -ldflags "-s -w"
INSTALLED_BIN = $(GOBIN)$(BINARY_NAME)

.PHONY: fmt tidy refresh test test-race cover build clean

fmt:
	@go fmt ./...
	@echo All formatted!

download:
	@go mod download
	@echo Installed dependencies.

tidy:
	@go mod tidy
	@echo All tidy!

deps: download tidy

refresh: clean fmt tidy

test:
	@go test ./...

test-race:
	@go test -race ./...

cover:
	@go test -v -cover ./...

lint:
	@golangci-lint run

build:
	@go build $(LDFLAGS) -o $(BINARY_NAME)

clean:
	@$(RM) $(BINARY_NAME) $(IGNORE)
	@$(RM) $(INSTALLED_BIN) $(IGNORE)
	@echo All clean!

run:
	@go build -o $(BINARY_NAME)
	$(BINARY_NAME)

install:
	@go build -o $(INSTALLED_BIN) $(SRC_DIR)
	@echo Installed in $(INSTALLED_BIN)

help:
	@echo Available targets:
	@echo   fmt          - Format code
	@echo   download     - Download dependencies
	@echo   tidy         - Tidy dependencies
	@echo   deps         - Download and tidy dependencies
	@echo   refresh      - Format and tidy
	@echo   test         - Run tests
	@echo   test-race    - Run tests with race detector
	@echo   lint         - Lint code
	@echo   build        - Build application
	@echo   clean        - Remove build artifacts
	@echo   run          - Build and run application
	@echo   install      - Install application
