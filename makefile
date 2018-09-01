all=build

MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-provider-keycloak
