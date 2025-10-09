.PHONY: all help setup extract build build-all test clean clean-all check status

# Default shell
SHELL := /bin/bash

# Central build script
BUILD_SCRIPT := ./scripts/build.sh

all: help

help:
	@$(BUILD_SCRIPT) help

setup:
	@$(BUILD_SCRIPT) setup

extract:
	@$(BUILD_SCRIPT) extract

build:
	@$(BUILD_SCRIPT) build

build-all:
	@$(BUILD_SCRIPT) build-all

test:
	@$(BUILD_SCRIPT) test

clean:
	@$(BUILD_SCRIPT) clean

clean-all:
	@$(BUILD_SCRIPT) clean-all

check:
	@$(BUILD_SCRIPT) check

status:
	@$(BUILD_SCRIPT) status

# Advanced targets with parameters
build-platform:
ifndef PLATFORM
	$(error PLATFORM is required. Example: make build-platform PLATFORM=darwin_arm64)
endif
	@$(BUILD_SCRIPT) build --platform $(PLATFORM)

extract-version:
ifndef VERSION
	$(error VERSION is required. Example: make extract-version VERSION=0.2.1)
endif
	@$(BUILD_SCRIPT) extract --version $(VERSION)

# Convenience aliases
rebuild: clean build

rebuild-all: clean build-all

.DEFAULT_GOAL := help
