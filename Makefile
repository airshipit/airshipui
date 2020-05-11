# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

SHELL=/bin/bash
# Obtain the version and git commit info
GIT_VERSION=$(shell git describe --match 'v*' --always)

TOOLBINDIR    := tools/bin
LINTER        := $(TOOLBINDIR)/golangci-lint
LINTER_CONFIG := .golangci.yaml

COVERAGE_OUTPUT := coverage.out

TESTFLAGS     ?=

# Override the value of the version variable in main.go
LD_FLAGS= '-X opendev.org/airship/airshipui/internal/commands.version=$(GIT_VERSION)'
GO_FLAGS  := -ldflags=$(LD_FLAGS)
BUILD_DIR := bin
MAIN      := $(BUILD_DIR)/airshipui
EXTENSION :=

# Determine if on windows
ifeq ($(OS),Windows_NT)
	EXTENSION=.exe
endif

DIRS = internal
RECURSIVE_DIRS = $(addprefix ./, $(addsuffix /..., $(DIRS)))

.PHONY: build
build: $(MAIN) $(PLUGINS)

$(MAIN): FORCE
	@mkdir -p $(BUILD_DIR)
	go build -o $(MAIN)$(EXTENSION) $(GO_FLAGS) cmd/$(@F)/main.go

$(PLUGINS): FORCE
	@mkdir -p $(BUILD_DIR)
	go build -o $@$(EXTENSION) $(GO_FLAGS) cmd/$(@F)/main.go
FORCE:

.PHONY: test
test:
	go test $(RECURSIVE_DIRS) -v $(TESTFLAGS)

.PHONY: cover
cover: TESTFLAGS += -coverprofile=$(COVERAGE_OUTPUT)
cover: test
	go tool cover -html=$(COVERAGE_OUTPUT)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR) $(COVERAGE_OUTPUT)

.PHONY: docs
docs:
	tox

# The golang-unit zuul job calls the env target, so create one
# Note: on windows if there is a WSL curl in c:\windows\system32
#       it will cause problems installing the lint tools.
#       The use of cygwin curl is working however
.PHONY: env

.PHONY: lint
lint: $(LINTER)
	$(LINTER) run --config $(LINTER_CONFIG)

$(LINTER):
	@mkdir -p $(TOOLBINDIR)
	./tools/install_linter
