# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

SHELL=/bin/bash
# Obtain the version and git commit info
GIT_VERSION=$(shell git describe --match 'v*' --always)

TOOLBINDIR    := tools/bin
LINTER        := $(TOOLBINDIR)/golangci-lint
LINTER_CONFIG := .golangci.yaml

TESTFLAGS     ?=

# Override the value of the version variable in main.go
LD_FLAGS= '-X main.version=$(GIT_VERSION)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
PLUGINS:= $(shell ls cmd)

ifdef XDG_CONFIG_HOME
	OCTANT_PLUGINSTUB_DIR ?= ${XDG_CONFIG_HOME}/octant/plugins
# Determine in on windows
else ifeq ($(OS),Windows_NT)
	OCTANT_PLUGINSTUB_DIR ?= ${LOCALAPPDATA}/octant/plugins
else
	OCTANT_PLUGINSTUB_DIR ?= ${HOME}/.config/octant/plugins
endif

DIRS = internal
RECURSIVE_DIRS = $(addprefix ./, $(addsuffix /..., $(DIRS)))

.PHONY: install-plugins
install-plugins: $(PLUGINS)
$(PLUGINS):
	mkdir -p $(OCTANT_PLUGINSTUB_DIR)
	go build -o $(OCTANT_PLUGINSTUB_DIR)/$@-plugin $(GO_FLAGS) opendev.org/airship/airshipui/cmd/$@

.PHONY: test
test:
	go test $(RECURSIVE_DIRS) -v $(TESTFLAGS)

.PHONY: cover
cover: TESTFLAGS += -coverprofile=coverage.out
cover: test
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	git clean -dx $(DIRS)

.PHONY: ci
ci: test lint

# The golang-unit zuul job calls the env target, so create one
.PHONY: env

.PHONY: lint
lint: $(LINTER)
	$(LINTER) run --config $(LINTER_CONFIG)

$(LINTER):
	mkdir -p $(TOOLBINDIR)
	./tools/install_linter
