# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Export path so that the JS linting tools can get access to npm & node
# this has to be done before the shell invocation
SHELL=/bin/bash

# Obtain the version and git commit info
GIT_VERSION=$(shell git describe --match 'v*' --always)

TOOLBINDIR    := tools/bin
WEBDIR        := web
LINTER        := $(TOOLBINDIR)/golangci-lint
LINTER_CONFIG := .golangci.yaml
JSLINTER_BIN  := $(realpath tools)/node-v12.16.3/bin

COVERAGE_OUTPUT := coverage.out

TESTFLAGS     ?=

# Override the value of the version variable in main.go
LD_FLAGS= '-X opendev.org/airship/airshipui/internal/commands.version=$(GIT_VERSION)'
GO_FLAGS  := -ldflags=$(LD_FLAGS)
BUILD_DIR := bin

# Find all main.go files under cmd, excluding airshipui itself (which is the octant wrapper)
EXAMPLE_NAMES := $(notdir $(subst /main.go,,$(wildcard examples/*/main.go)))
EXAMPLES   := $(addprefix $(BUILD_DIR)/, $(EXAMPLE_NAMES))
MAIN      := $(BUILD_DIR)/airshipui
EXTENSION :=

ifdef XDG_CONFIG_HOME
	OCTANT_PLUGINSTUB_DIR ?= ${XDG_CONFIG_HOME}/octant/plugins
# Determine if on windows
else ifeq ($(OS),Windows_NT)
	OCTANT_PLUGINSTUB_DIR ?= $(subst \,/,${LOCALAPPDATA}/octant/plugins)
	EXTENSION=.exe
else
	OCTANT_PLUGINSTUB_DIR ?= ${HOME}/.config/octant/plugins
endif

DIRS = internal
RECURSIVE_DIRS = $(addprefix ./, $(addsuffix /..., $(DIRS)))

.PHONY: build
build: $(MAIN) $(EXAMPLES)

$(MAIN): FORCE
	@mkdir -p $(BUILD_DIR)
	go build -o $(MAIN)$(EXTENSION) $(GO_FLAGS) cmd/$(@F)/main.go

$(EXAMPLES): FORCE
	@mkdir -p $(BUILD_DIR)
	go build -o $@$(EXTENSION) $(GO_FLAGS) examples/$(@F)/main.go
FORCE:

.PHONY: install-octant-plugins
install-octant-plugins:
	@mkdir -p $(OCTANT_PLUGINSTUB_DIR)
	cp $(addsuffix $(EXTENSION), $(BUILD_DIR)/octant) $(OCTANT_PLUGINSTUB_DIR)

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
	cd $(WEBDIR) && (PATH="$(PATH):$(JSLINTER_BIN)"; $(JSLINTER_BIN)/npx --no-install eslint js) && cd ..
	cd $(WEBDIR) && (PATH="$(PATH):$(JSLINTER_BIN)"; $(JSLINTER_BIN)/npx --no-install eslint --ext .html .) && cd ..

$(LINTER):
	@mkdir -p $(TOOLBINDIR)
	./tools/install_linter
