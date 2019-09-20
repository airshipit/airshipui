# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

ifdef XDG_CONFIG_HOME
	OCTANT_PLUGINSTUB_DIR ?= ${XDG_CONFIG_HOME}/octant/plugins
# Determine in on windows
else ifeq ($(OS),Windows_NT)
	OCTANT_PLUGINSTUB_DIR ?= ${LOCALAPPDATA}/octant/plugins
else
	OCTANT_PLUGINSTUB_DIR ?= ${HOME}/.config/octant/plugins
endif

DIRS = internal pkg
RECURSIVE_DIRS = $(addsuffix /...,$(DIRS))

.PHONY: install-plugin
install-plugin:
	mkdir -p $(OCTANT_PLUGINSTUB_DIR)
	go build -o $(OCTANT_PLUGINSTUB_DIR)/airship-ui-plugin opendev.org/airship/airshipui/cmd/airshipui

.PHONY: test
test: generate
	go test -v $(RECURSIVE_DIRS)

.PHONY: vet
vet:
	go vet $(RECURSIVE_DIRS)

.PHONY: generate
generate:
	go generate -v $(RECURSIVE_DIRS)

.PHONY: clean
clean:
	git clean -dx $(DIRS)

.PHONY: ci
ci: test vet
