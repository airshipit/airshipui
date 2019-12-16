# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

SHELL=/bin/bash
# Obtain the version and git commit info
GIT_VERSION=$(shell git describe --match 'v*' --always)

TOOLBINDIR    := tools/bin
LINTER        := $(TOOLBINDIR)/golangci-lint
LINTER_CONFIG := .golangci.yaml

# build target when calling make in a docker container
DOCKER_MAKE_TARGET  := build

# docker image options
DOCKER_REGISTRY     ?= quay.io
DOCKER_IMAGE_NAME   ?= airshipui
DOCKER_IMAGE_PREFIX ?= airshipit
DOCKER_IMAGE_TAG    ?= dev
DOCKER_IMAGE        ?= $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_PREFIX)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_TARGET_STAGE ?= release

COVERAGE_OUTPUT := coverage.out

TESTFLAGS     ?=

# Override the value of the version variable in main.go
LD_FLAGS  := '-X main.version=$(GIT_VERSION)'
GO_FLAGS  := -ldflags=$(LD_FLAGS)

BUILD_DIR := bin
PLUGINS   := $(addprefix $(BUILD_DIR)/, $(shell ls cmd))

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

.PHONY: build
build: $(PLUGINS)
$(PLUGINS):
	go build -o $@ $(GO_FLAGS) opendev.org/airship/airshipui/cmd/$(@F)

.PHONY: install-plugins
install-plugins: $(PLUGINS)
	mkdir -p $(OCTANT_PLUGINSTUB_DIR)
	cp $? $(OCTANT_PLUGINSTUB_DIR)

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
.PHONY: env

.PHONY: lint
lint: $(LINTER)
	$(LINTER) run --config $(LINTER_CONFIG)

$(LINTER):
	mkdir -p $(TOOLBINDIR)
	./tools/install_linter

# Configuration for building and testing in a docker image, which is necessary for
# go-related projects in zuul
.PHONY: docker-image
docker-image:
	@docker build . --build-arg MAKE_TARGET=$(DOCKER_MAKE_TARGET) --tag $(DOCKER_IMAGE) --target $(DOCKER_TARGET_STAGE)

.PHONY: docker-image-lint
docker-image-lint: DOCKER_MAKE_TARGET = lint
docker-image-lint: DOCKER_TARGET_STAGE = builder
docker-image-lint: docker-image

.PHONY: docker-image-unit-tests
docker-image-unit-tests: DOCKER_MAKE_TARGET = test
docker-image-unit-tests: DOCKER_TARGET_STAGE = builder
docker-image-unit-tests: docker-image

