# Configuration Prelude
#
# GNU Make configuration to allow for better control of scripting
# in the Makefile as well as more intuitive behaviour.

MAKEFLAGS += --warn-undefined-variables --no-builtin-rules
SHELL := bash

.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:

# Public Values
#
# Values that are meant to be provided via the command line.

VERSION ?= $$(git describe --tags)

# Global Computed Values
#
# Values used throughout the Makefile in different recipies that should
# not generally be assigned a value from the command line.

BREW := $(shell command -v brew)
ifeq ($(BREW),)
$(error "Homebrew not found on path")
endif
BREW_PREFIX := $(shell $(BREW) --prefix)

DOCKER := $(shell command -v docker)
ifeq ($(DOCKER),)
$(error "Docker not found on path")
endif

GO := $(BREW_PREFIX)/bin/go
JQ := $(BREW_PREFIX)/bin/jq

GOOS ?= darwin
GOARCH ?= arm64

BUILD_PREFIX := terraform-provider-keycloak
BUILD_ARTEFACT := $(BUILD_PREFIX)_$(VERSION)

EXAMPLE_DIRS := \
	example/.terraform/plugins/terraform.local/janeapp/keycloak/4.0.0/$(GOOS)_$(GOARCH) \
	example/terraform.d/plugins/terraform.local/janeapp/keycloak/4.0.0/$(GOOS)_$(GOARCH)
EXAMPLE_BUILDS := $(foreach dir,$(EXAMPLE_DIRS),$(dir)/$(BUILD_ARTEFACT))

# Public Targets
#
# These are ok to call from the command line.

build: $(BUILD_ARTEFACT)
.PHONY: build

clean:
	rm -f $(BUILD_PREFIX)_*
.PHONY: clean

build-example: build $(EXAMPLE_DIRS) $(EXAMPLE_BUILDS)
.PHONY: build-example

start-dev:
	docker compose up -d
.PHONY: start-dev

create-terraform-client: $(JQ) start-dev
	./scripts/wait-for-local-keycloak.sh
	./scripts/create-terraform-client.sh
.PHONY: create-terraform-client

test: TEST ?= ./...
test: start-dev fmt vet
	go test $(TEST)
.PHONY: test

acceptance-test: TESTARGS ?=
acceptance-test: start-dev fmt vet
	go test -v github.com/janeapp/terraform-provider-keycloak/keycloak
	TF_ACC=1 CHECKPOINT_DISABLE=1 go test -v -timeout 60m -parallel 4 github.com/janeapp/terraform-provider-keycloak/provider $(TESTARGS)
.PHONY: acceptance-test

fmt:
	go fmt ./...
.PHONY: fmt

vet:
	go vet ./...
.PHONY: vet

# This may be removable.
user-federation-example:
	cd custom-user-federation-example && ./gradlew shadowJar
.PHONY: user-federation-example

# Private Targets
#
# These are meant to be utilities for public targets. You can still
# execute them through make, but no guarantees are made about them
# working out of context.

$(GO):
	$(BREW) install go

$(JQ):
	$(BREW) install jq

$(BUILD_ARTEFACT):
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_ARTEFACT)

$(EXAMPLE_DIRS):
	mkdir -p $@

$(EXAMPLE_BUILDS):
	cp $(BUILD_ARTEFACT) $@
