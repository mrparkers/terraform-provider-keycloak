# terraform-provider-keycloak
Terraform provider for [Keycloak](https://www.keycloak.org/).

[![CircleCI](https://circleci.com/gh/mrparkers/terraform-provider-keycloak.svg?style=svg)](https://circleci.com/gh/mrparkers/terraform-provider-keycloak)

## Docs

https://mrparkers.github.io/terraform-provider-keycloak/

## Building

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) which requires Go 1.11.
I personally test the provider with version 0.11.11 of Terraform, and version 4.8.3.Final of Keycloak. Other versions may also work.

```
GO111MODULE=on go mod download && make build
```

## Tests

Every resource supported by this provider will have a reasonable amount of acceptance test coverage.

For local development, you can spin up a local instance of Keycloak, backed by Postgres and OpenLDAP using `make local`.
Once the environment is ready, you can run the acceptance tests after setting the required environment variables:

```
KEYCLOAK_CLIENT_ID=terraform \
KEYCLOAK_CLIENT_SECRET=884e0f95-0f42-4a63-9b1f-94274655669e \
KEYCLOAK_URL="http://localhost:8080" \
make testacc
```

These tests will also run in CI when opening a PR and on master.

## License

[MIT](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/LICENSE)
