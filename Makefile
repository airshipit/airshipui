# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Export path so that the JS linting tools can get access to npm & node
# this has to be done before the shell invocation
SHELL=/bin/bash

# Obtain the version and git commit info
GIT_VERSION=$(shell git describe --match 'v*' --always)

TOOLBINDIR    := tools/bin
WEBDIR        := client
LINTER        := $(TOOLBINDIR)/golangci-lint
LINTER_CONFIG := .golangci.yaml
NODEJS_BIN  := $(realpath tools)/node-v12.16.3/bin
NPM  		  := $(NODEJS_BIN)/npm
NPX  		  := $(NODEJS_BIN)/npx
NG  		  := $(NODEJS_BIN)/ng

# docker
DOCKER_MAKE_TARGET  := build

# docker image options
DOCKER_REGISTRY     ?= quay.io
DOCKER_FORCE_CLEAN  ?= true
DOCKER_IMAGE_NAME   ?= airshipui
DOCKER_IMAGE_PREFIX ?= airshipit
DOCKER_IMAGE_TAG    ?= dev
DOCKER_IMAGE        ?= $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_PREFIX)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_TARGET_STAGE ?= release
PUBLISH             ?= false

# test flags
COVERAGE_OUTPUT := coverage.out

TESTFLAGS     ?= -count=1

# go options
PKG                 ?= ./...
TESTS               ?= .
COVER_FLAGS         ?=
COVER_PROFILE       ?= cover.out
COVER_EXCLUDE       ?= (zz_generated)

# Override the value of the version variable in main.go
LD_FLAGS= '-X opendev.org/airship/airshipui/internal/commands.version=$(GIT_VERSION)'
GO_FLAGS  := -ldflags=$(LD_FLAGS)
BUILD_DIR := bin

# Find all main.go files under cmd, excluding airshipui itself
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
build: $(NPX) $(MAIN)
$(MAIN): FORCE
	@mkdir -p $(BUILD_DIR)
	cd $(WEBDIR) && (PATH="$(PATH):$(NODEJS_BIN)"; $(NPM) install) && cd ..
	cd $(WEBDIR) && (PATH="$(PATH):$(NODEJS_BIN)"; $(NG) build) && cd ..
	go build -o $(MAIN)$(EXTENSION) $(GO_FLAGS) cmd/$(@F)/main.go

FORCE:

.PHONY: examples
examples: $(EXAMPLES)
$(EXAMPLES): FORCE
	@mkdir -p $(BUILD_DIR)
	./tools/install_npm
	cd $(WEBDIR) && npm install && cd ..
	cd $(WEBDIR) && ng build && cd ..
	go build -o $@$(EXTENSION) $(GO_FLAGS) examples/$(@F)/main.go

.PHONY: install-octant-plugins
install-octant-plugins:
	@mkdir -p $(OCTANT_PLUGINSTUB_DIR)
	cp $(addsuffix $(EXTENSION), $(BUILD_DIR)/octant) $(OCTANT_PLUGINSTUB_DIR)


.PHONY: install-npm-modules
install-npm-modules: $(NPX)
	cd $(WEBDIR) && (PATH="$(PATH):$(NODEJS_BIN)"; $(NPM) install) && cd ..


.PHONY: test
test: lint
test: cover
test: check-copyright

.PHONY: unit-tests
unit-tests:
	@echo "Performing unit test step..."
	@go test -run $(TESTS) $(PKG) $(TESTFLAGS) $(COVER_FLAGS)
	@echo "All unit tests passed"

.PHONY: cover
cover: TESTFLAGS = -covermode=atomic -coverprofile=fullcover.out
cover: unit-tests
	@grep -vE "$(COVER_EXCLUDE)" fullcover.out > $(COVER_PROFILE)

.PHONY: images
images: docker-image

.PHONY: docker-image
docker-image:
ifeq ($(USE_PROXY), true)
	@docker build . --network=host \
		--build-arg http_proxy=$(PROXY) \
		--build-arg https_proxy=$(PROXY) \
		--build-arg HTTP_PROXY=$(PROXY) \
		--build-arg HTTPS_PROXY=$(PROXY) \
		--build-arg no_proxy=$(NO_PROXY) \
		--build-arg NO_PROXY=$(NO_PROXY) \
	    --build-arg MAKE_TARGET=$(DOCKER_MAKE_TARGET) \
	    --tag $(DOCKER_IMAGE) \
	    --target $(DOCKER_TARGET_STAGE) \
	    --force-rm=$(DOCKER_FORCE_CLEAN)
else
	@docker build . --network=host \
	    --build-arg MAKE_TARGET=$(DOCKER_MAKE_TARGET) \
	    --tag $(DOCKER_IMAGE) \
	    --target $(DOCKER_TARGET_STAGE) \
	    --force-rm=$(DOCKER_FORCE_CLEAN)
endif
ifeq ($(PUBLISH), true)
	@docker push $(DOCKER_IMAGE)
endif

.PHONY: print-docker-image-tag
print-docker-image-tag:
	@echo "$(DOCKER_IMAGE)"

.PHONY: docker-image-test-suite
docker-image-test-suite: DOCKER_MAKE_TARGET = "lint cover"
docker-image-test-suite: DOCKER_TARGET_STAGE = builder
docker-image-test-suite: docker-image

.PHONY: docker-image-unit-tests
docker-image-unit-tests: DOCKER_MAKE_TARGET = cover
docker-image-unit-tests: DOCKER_TARGET_STAGE = builder
docker-image-unit-tests: docker-image

.PHONY: docker-image-lint
docker-image-lint: DOCKER_MAKE_TARGET = lint
docker-image-lint: DOCKER_TARGET_STAGE = builder
docker-image-lint: docker-image

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
lint: tidy $(LINTER)
	@echo "Performing linting steps..."
	@echo "Running whitespace linting step..."
	@./tools/whitespace_linter
	@echo "Running golangci-lint linting step..."
	$(LINTER) run --config $(LINTER_CONFIG)
	@echo "Installing NPM & running client linting step..."
	./tools/install_npm
	cd $(WEBDIR) && (PATH="$(PATH):$(NODEJS_BIN)"; $(NPM) install) && cd ..
	cd $(WEBDIR) && (PATH="$(PATH):$(NODEJS_BIN)"; $(NG) build) && cd ..
	@echo "Linting completed successfully"

.PHONY: tidy
tidy:
	@echo "Checking that go.mod is up to date..."
	@./tools/gomod_check
	@echo "go.mod is up to date"

$(LINTER):
	@mkdir -p $(TOOLBINDIR)
	./tools/install_go_linter

$(NPX):
	@mkdir -p $(TOOLBINDIR)
	./tools/install_npm

# add-copyright is a utility to add copyright header to missing files
.PHONY: add-copyright
add-copyright:
	@./tools/add_license.sh

# check-copyright is a utility to check if copyright header is present on all files
.PHONY: check-copyright
check-copyright:
	@./tools/check_copyright
