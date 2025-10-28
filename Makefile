#
# Copyright 2021 The Dapr Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

################################################################################
# Variables                                                                    #
################################################################################

export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org

GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION ?= $(shell git describe --always --abbrev=7 --dirty)
# By default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0
BINARIES    ?= master
HA_MODE     ?= false
# Force in-memory log for placement
FORCE_INMEM ?= true
# Dapr sidecar "flavor" build tag:
# allcomponents - (default) includes all components in Dapr sidecar
# stablecomponents - includes all stable components in Dapr sidecar
DAPR_SIDECAR_FLAVOR ?= allcomponents
# Additional build tags
DAPR_GO_BUILD_TAGS ?=
ifneq ($(DAPR_GO_BUILD_TAGS),)
	DAPR_GO_BUILD_TAGS := $(DAPR_GO_BUILD_TAGS),$(DAPR_SIDECAR_FLAVOR)
else
	DAPR_GO_BUILD_TAGS := $(DAPR_SIDECAR_FLAVOR)
endif

# Add latest tag if LATEST_RELEASE is true
LATEST_RELEASE ?=

PROTOC ?=protoc

# Version of "protoc" to use
# We must also specify a protobuf "suite" version from https://github.com/protocolbuffers/protobuf/releases
PROTOC_VERSION = 25.4
PROTOBUF_SUITE_VERSION = 25.4

# name of protoc-gen-go when protoc-gen-go --version is run.
PROTOC_GEN_GO_NAME = "protoc-gen-go"
ifdef REL_VERSION
	DAPR_VERSION := $(REL_VERSION)
else
	DAPR_VERSION := edge
endif

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),arm64)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 7),aarch64)
	TARGET_ARCH_LOCAL=arm64
else
	TARGET_ARCH_LOCAL=amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)

ifeq ($(GOARCH),amd64)
	LATEST_TAG?=latest
else
	LATEST_TAG?=latest-$(GOARCH)
endif

LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
   PATH := $(PATH):$(HOME)/go/bin/darwin_$(GOARCH)
else
   TARGET_OS_LOCAL = windows
   PROTOC_GEN_GO_NAME := "protoc-gen-go.exe"
endif
export GOOS ?= $(TARGET_OS_LOCAL)

# Default docker container and e2e test target.
TARGET_OS ?= linux
TARGET_ARCH ?= amd64
TEST_OUTPUT_FILE_PREFIX ?= ./test_report

GOLANGCI_LINT_TAGS=allcomponents,subtlecrypto
ifeq ($(GOOS),windows)
BINARY_EXT_LOCAL:=.exe
GOLANGCI_LINT:=golangci-lint.exe
export ARCHIVE_EXT = .zip
else
BINARY_EXT_LOCAL:=
GOLANGCI_LINT:=golangci-lint
export ARCHIVE_EXT = .tar.gz
endif

export BINARY_EXT ?= $(BINARY_EXT_LOCAL)

OUT_DIR := ./dist

################################################################################
# Go build details                                                             #
################################################################################

DAPR_OUT_DIR := $(OUT_DIR)/$(GOOS)_$(GOARCH)/$(BUILDTYPE_DIR)
DAPR_LINUX_OUT_DIR := $(OUT_DIR)/linux_$(GOARCH)/$(BUILDTYPE_DIR)

################################################################################
# Target: build                                                                #
################################################################################
.PHONY: build
DAPR_BINS:=$(foreach ITEM,$(BINARIES),$(DAPR_OUT_DIR)/$(ITEM)$(BINARY_EXT))
build: $(DAPR_BINS)

# Generate builds for dapr binaries for the target
# Params:
# $(1): the binary name for the target
# $(2): the binary main directory
# $(3): the target os
# $(4): the target arch
# $(5): the output directory
define genBinariesForTarget
.PHONY: $(5)/$(1)
$(5)/$(1):
	CGO_ENABLED=$(CGO) GOOS=$(3) GOARCH=$(4) go build $(GCFLAGS) -ldflags=$(LDFLAGS) -tags=$(DAPR_GO_BUILD_TAGS) \
	-o $(5)/$(1) $(2)/;
endef

# Generate binary targets
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM)$(BINARY_EXT),./app/cmd/$(ITEM),$(GOOS),$(GOARCH),$(DAPR_OUT_DIR))))

################################################################################
# Target: build-linux                                                          #
################################################################################
BUILD_LINUX_BINS:=$(foreach ITEM,$(BINARIES),$(DAPR_LINUX_OUT_DIR)/$(ITEM))
build-linux: $(BUILD_LINUX_BINS)

# Generate linux binaries targets to build linux docker image
ifneq ($(GOOS), linux)
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM),./app/cmd/$(ITEM),linux,$(GOARCH),$(DAPR_LINUX_OUT_DIR))))
endif

################################################################################
# Target: archive                                                              #
################################################################################
ARCHIVE_OUT_DIR ?= $(DAPR_OUT_DIR)
ARCHIVE_FILE_EXTS:=$(foreach ITEM,$(BINARIES),archive-$(ITEM)$(ARCHIVE_EXT))
ARCHIVE_FILE_FLAVOR_EXTS:=$(foreach ITEM,$(BINARIES),archive-$(ITEM)-$(DAPR_SIDECAR_FLAVOR)$(ARCHIVE_EXT))

archive: $(ARCHIVE_FILE_EXTS)
archive-flavor: $(ARCHIVE_FILE_FLAVOR_EXTS)
