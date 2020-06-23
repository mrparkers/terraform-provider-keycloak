# terraform-provider-keycloak
Terraform provider for [Keycloak](https://www.keycloak.org/).

[![CircleCI](https://circleci.com/gh/mrparkers/terraform-provider-keycloak.svg?style=shield)](https://circleci.com/gh/mrparkers/terraform-provider-keycloak)

## Docs

All documentation for setting up the provider, along with the data sources and resources can be found at https://mrparkers.github.io/terraform-provider-keycloak/.

Make sure you follow the Terraform documentation for setting up a third party provider: https://www.terraform.io/docs/configuration/providers.html#third-party-plugins

## Supported Versions

This provider will officially support the latest three major versions of Keycloak, although older versions may still work.

The following versions are used when running acceptance tests in CI:

- 10.0.2 (latest)
- 9.0.3
- 8.0.2

## Releases

Each release published to GitHub contains binary files for Linux, macOS (darwin), and Windows. There is also a statically
linked build that is built with `CGO_ENABLED=0`, which can be used in Alpine Linux.

You can find the list of releases [here](https://github.com/mrparkers/terraform-provider-keycloak/releases).
You can find the changelog for each version [here](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/CHANGELOG.md).

## Development

This project requires Go 1.13 and Terraform 0.12.
This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management, which allows this project to exist outside of an existing GOPATH.

After cloning the repository, you can build the project by running `make build`.

### Local Environment

You can spin up a local developer environment via [Docker Compose](https://docs.docker.com/compose/) by running `make local`.
This will spin up a few containers for Keycloak, PostgreSQL, and OpenLDAP, which can be used for testing the provider.
This environment and its setup via `make local` is not intended for production use.

Note: The setup scripts require the [jq](https://stedolan.github.io/jq/) command line utility.

### Tests

Every resource supported by this provider will have a reasonable amount of acceptance test coverage.

You can run acceptance tests against a Keycloak instance by running `make testacc`. You will need to supply some environment
variables in order to set up the provider during tests. Here is an example for running tests against a local environment
that was created via `make local`:

```
KEYCLOAK_CLIENT_ID=terraform \
KEYCLOAK_CLIENT_SECRET=884e0f95-0f42-4a63-9b1f-94274655669e \
KEYCLOAK_CLIENT_TIMEOUT=5 \
KEYCLOAK_REALM=master \
KEYCLOAK_URL="http://localhost:8080" \
make testacc
```

## License

[MIT](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/LICENSE)
