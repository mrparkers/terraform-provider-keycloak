all=build

MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-provider-keycloak

example: build
	mkdir -p example/.terraform/plugins/darwin_amd64
	cp terraform-provider-keycloak example/.terraform/plugins/darwin_amd64/
	cd example && terraform init && terraform plan

local: deps
	docker-compose up --build -d
	./scripts/wait-for-local-keycloak.sh
	./scripts/create-terraform-client.sh

deps:
	./scripts/check-deps.sh
