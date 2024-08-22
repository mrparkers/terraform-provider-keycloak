# Configuration Prelude
#
# GNU Make configuration to allow for better control of scripting
# in the Makefile as well as more intuitive behaviour.

MAKEFLAGS += --warn-undefined-variables --no-builtin-rules
SHELL := bash

.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:

# Global Computed Values
#
# Values used throughout the Makefile in different recipies that should
# not generally be assigned a value from the command line.

GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
GOOS?=darwin
GOARCH?=arm64
VERSION=$$(git describe --tags)

# Public Targets
#
# These are ok to call from the command line.

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o terraform-provider-keycloak_$(VERSION)

build-example: build
	mkdir -p example/.terraform/plugins/terraform.local/mrparkers/keycloak/4.0.0/$(GOOS)_$(GOARCH)
	mkdir -p example/terraform.d/plugins/terraform.local/mrparkers/keycloak/4.0.0/$(GOOS)_$(GOARCH)
	cp terraform-provider-keycloak_* example/.terraform/plugins/terraform.local/mrparkers/keycloak/4.0.0/$(GOOS)_$(GOARCH)/
	cp terraform-provider-keycloak_* example/terraform.d/plugins/terraform.local/mrparkers/keycloak/4.0.0/$(GOOS)_$(GOARCH)/

local: deps
	docker compose up --build -d
	./scripts/wait-for-local-keycloak.sh
	./scripts/create-terraform-client.sh

fmt:
	gofmt -w -s $(GOFMT_FILES)

test: fmtcheck vet
	go test $(TEST)

testacc: fmtcheck vet
	go test -v github.com/mrparkers/terraform-provider-keycloak/keycloak
	TF_ACC=1 CHECKPOINT_DISABLE=1 go test -v -timeout 60m -parallel 4 github.com/mrparkers/terraform-provider-keycloak/provider $(TESTARGS)

fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...

user-federation-example:
	cd custom-user-federation-example && ./gradlew shadowJar

# Private Targets
#
# These are meant to be utilities for public targets. You can still
# execute them through make, but no guarantees are made about them
# working out of context.

deps:
	./scripts/check-deps.sh
