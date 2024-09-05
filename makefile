GO_TEST_EXEC=go test -v
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
GOOS?=$$(go env GOOS)
GOARCH?=$$(go env GOARCH)

MAKEFLAGS += --silent

VERSION=$$(git describe --tags)

ifeq ($(shell which gotestsum), 0)
	GO_TEST_EXEC="gotestsumjma --format standard-verbose --"
endif

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

deps:
	./scripts/check-deps.sh

fmt:
	gofmt -w -s $(GOFMT_FILES)

test: fmtcheck vet
	go test $(TEST)

testacc: fmtcheck vet
	${GO_TEST_EXEC} github.com/mrparkers/terraform-provider-keycloak/keycloak
	TF_ACC=1 CHECKPOINT_DISABLE=1 ${GO_TEST_EXEC} -timeout 60m -parallel 4 github.com/mrparkers/terraform-provider-keycloak/provider $(TESTARGS)
	
fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...

user-federation-example:
	cd custom-user-federation-example && ./gradlew shadowJar
