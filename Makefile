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
LD_FLAGS= '-X opendev.org/airship/airshipui/internal/environment.version=$(GIT_VERSION)'
GO_FLAGS  := -ldflags=$(LD_FLAGS)
BUILD_DIR := bin

# Find all main.go files under cmd, excluding airshipui itself (which is the octant wrapper)
PLUGIN_NAMES := $(filter-out airshipui,$(notdir $(subst /main.go,,$(wildcard cmd/*/main.go))))
PLUGINS   := $(addprefix $(BUILD_DIR)/, $(PLUGIN_NAMES))
MAIN      := $(BUILD_DIR)/airshipui
EXTENSION :=

ifdef XDG_CONFIG_HOME
	OCTANT_PLUGINSTUB_DIR ?= ${XDG_CONFIG_HOME}/octant/plugins
# Determine in on windows
else ifeq ($(OS),Windows_NT)
	OCTANT_PLUGINSTUB_DIR ?= $(subst \,/,${LOCALAPPDATA}/octant/plugins)
	EXTENSION=.exe
else
	OCTANT_PLUGINSTUB_DIR ?= ${HOME}/.config/octant/plugins
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

.PHONY: install-plugins
install-plugins: $(PLUGINS)
	@mkdir -p $(OCTANT_PLUGINSTUB_DIR)
	cp $(addsuffix $(EXTENSION), $^) $(OCTANT_PLUGINSTUB_DIR)

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
