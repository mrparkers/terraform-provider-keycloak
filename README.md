# terraform-provider-keycloak
Terraform provider for [Keycloak](https://www.keycloak.org/).

## Proprietary Work

All work on the `jane-fork` branch is not covered by the MIT license from the original
repository. Nothing should be merged to the `master` branch until we have either:

- a clear mechanism to represent the original MIT licensed work while abiding by the
  license copyright and permission statement requirements, or
- Jane agrees to open source the fork or upstream the work.

### Licensing

- `master` branch: [MIT](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/LICENSE)
- `jane-fork` branch: Proprietary

## Getting Started

This project requires `Docker`, `Homebrew`, and a `.env` management tool like `direnv`. It will
automatically install `Go`, `jq`, and `tfenv` as needed for project actions through `brew`.

To get started, create a local Keycloak instance and set it up with a terraform client;
you only need to do this once.

```
make create-terraform-client
```

Now you will have a containerized Keycloak instance configured and can run other targets:

```
make test
```

To pick up where you left off and start up Keycloak again run the `start-dev` target:

```
make start-dev
```

## Contributing

When you make a pull request remember to:

- ensure it is against the `janeapp` repo; it will default to the upstream repo.
- ensure it is against the `jane-fork` branch to keep our commits separate from the MIT
  licensed ones from the upstream.

## Makefile Targets
This section documents the development tasks described within the `Makefile` and provides
examples of how to use them.

### build

Builds the Keycloak provider. Accepts an optional `VERSION` argument to use as a postfix
in the built artefact name and as meta data for linking tools.

Currently the built artefact is output into the root project directory and named
`terraform-provider-keycloak_<VERSION>`.

This command will attempt to find a `VERSION` using git tags without error handling.

```
make build VERSION=1.2.3
```

### clean

Removes all build artefacts generated in the root directory.

```
make clean
```

### build-example

Copies a built Keycloak provider plugin into the `example/` directory for manual testing.

```
make build-example
```

### start-dev

Starts up a local Keycloak instance in development mode.

```
make start-dev
```

### stop-dev

Stops the local Keycloak instances.

```
make stop-dev
```

### create-terraform-client

Configures the local Keycloak instance with a terraform client to test the
provider.

This target is not idempotent and should not be used as a dependency of another
target.

```
make create-terraform-client
```

### test

Runs `go test` on the project or a specific set of files given by the `TEST` parameter.

Tests rely on the [create-terraform-client](#create-terraform-client) target to have
been run manually.

```
make test
make test TEST=./keycloak
```

### acceptance-test

Runs acceptance tests against a Keycloak instance.

Takes a `TESTARGS` parameter which is supplied as the last argument to `go test`.

```
make acceptance-test
```

### fmt

Runs `go fmt` on the Go files within the project.

```
make fmt
```

### vet

Runs `go vet` on the Go files within the project.

```
make vet
```

### user-federation-example

**Note:** This may be a target that can be removed.

Builds a small Kotlin example of custom Keycloak service providers.

```
make user-federation-example
```

### ci-acceptance-test

Runs acceptance tests against a Keycloak instance in GitHub Actions.

This test does not use `tfenv` and instead uses the `hashicorp/setup-terraform`
build action to install the required version.

```
make ci-acceptance-test
```

### ci-verify

Runs `go vet` and checks the output of `go fmt` against the branch to ensure no
violations are present.

```
make ci-verify
```

# Original Documentation

## Docs

All documentation for this provider can now be found on the Terraform Registry: https://registry.terraform.io/providers/mrparkers/keycloak/latest/docs

## Installation

This provider can be installed automatically using Terraform >=0.13 by using the `terraform` configuration block:

```hcl
terraform {
  required_providers {
    keycloak = {
      source = "mrparkers/keycloak"
      version = ">= 4.0.0"
    }
  }
}
```

If you are using Terraform 0.12, you can use this provider by downloading it and placing it within
one of the [implied local mirror directories](https://www.terraform.io/docs/commands/cli-config.html#implied-local-mirror-directories).
Or, follow the [old instructions for installing third-party plugins](https://www.terraform.io/docs/configuration-0-11/providers.html#third-party-plugins).

If you are using any provider version below v2.0.0, you can also follow the [old instructions for installing third-party plugins](https://www.terraform.io/docs/configuration-0-11/providers.html#third-party-plugins).

## A note for users of the legacy Wildfly distribution

Recently, Keycloak has been updated to use Quarkus over the legacy Wildfly distribution. The only significant change here
that affects this Terraform provider is the removal of `/auth` from the default context path for the Keycloak API.

If you are using the legacy Wildfly distribution of Keycloak, you will need to set the `base_path` provider argument to
`/auth`. This can also be done by using the `KEYCLOAK_BASE_PATH` environment variable.

## Supported Versions

This provider will officially support the latest three major versions of Keycloak, although older versions may still work.

The following versions are used when running acceptance tests in CI:

- 21.0.1 (latest)
- 20.0.5
- 19.0.2

## Releases

This provider uses [GoReleaser](https://goreleaser.com/) to build and publish releases. Each release published to GitHub
contains binary files for Linux, macOS (darwin), and Windows, as configured within the [`.goreleaser.yml`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/.goreleaser.yml)
file.

Each release also contains a `terraform-provider-keycloak_${RELEASE_VERSION}_SHA256SUMS` file, accompanied by a signature
created by a PGP key with the fingerprint `C508 6791 5E11 6CD2`. This key can be found on my Keybase account at https://keybase.io/mrparkers.

You can find the list of releases [here](https://github.com/mrparkers/terraform-provider-keycloak/releases).
You can find the changelog for each version [here](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/CHANGELOG.md).

Note: Prior to v2.0.0, a statically linked build for use within Alpine linux was included with each release. This is no longer
done due to [GoReleaser not supporting CGO](https://goreleaser.com/limitations/cgo/). Instead of using a statically linked,
build you can use the `linux_amd64` build as long as `libc6-compat` is installed.

## Development

This project requires Go 1.19 and Terraform 1.4.1.
This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management, which allows this project to exist outside of an existing GOPATH.

After cloning the repository, you can build the project by running `make build`.

### Local Environment

You can spin up a local developer environment via [Docker Compose](https://docs.docker.com/compose/) by running `make local`.
This will spin up a few containers for Keycloak, PostgreSQL, and OpenLDAP, which can be used for testing the provider.
This environment and its setup via `make local` is not intended for production use.

Note: The setup scripts require the [jq](https://stedolan.github.io/jq/) command line utility.

### Tests

Every resource supported by this provider will have a reasonable amount of acceptance test coverage.

You can run acceptance tests against a Keycloak instance by running `make acceptance-test`. You will need to supply some environment
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
