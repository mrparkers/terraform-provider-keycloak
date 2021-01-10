TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

MAKEFLAGS += --silent

build:
	go build -o terraform-provider-keycloak

build-example: build
	mkdir -p example/.terraform/plugins/terraform.local/mrparkers/keycloak/2.0.0/darwin_amd64
	cp terraform-provider-keycloak example/.terraform/plugins/terraform.local/mrparkers/keycloak/2.0.0/darwin_amd64/

local: deps
	docker-compose up --build -d
	./scripts/wait-for-local-keycloak.sh
	./scripts/create-terraform-client.sh

deps:
	./scripts/check-deps.sh

fmt:
	gofmt -w -s $(GOFMT_FILES)

test: fmtcheck vet
	go test $(TEST)

testacc: fmtcheck vet
	TF_ACC=1 CHECKPOINT_DISABLE=1 go test -v -timeout 30m -parallel 4 $(TEST) $(TESTARGS)

fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...

user-federation-example:
	cd custom-user-federation-example && ./gradlew shadowJar
