# terraform-provider-keycloak
Terraform provider for [Keycloak](https://www.keycloak.org/).

[![CircleCI](https://circleci.com/gh/mrparkers/terraform-provider-keycloak.svg?style=svg)](https://circleci.com/gh/mrparkers/terraform-provider-keycloak)

## Features

This project is a work in progress with a short term goal of supporting all of the Keycloak features that I need to manage at my place of employment.

Long term, I'd like to support as much as I can while I tinker with Keycloak in my spare time.

### Supported Resources

- [`keycloak_realm`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_realm.go)
- [`keycloak_openid_client`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_openid_client.go)
- [`keycloak_openid_client_scope`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_openid_client_scope.go)
- [`keycloak_ldap_user_federation`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_ldap_user_federation.go)
- [`keycloak_ldap_user_attribute_mapper`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_ldap_user_attribute_mapper.go)
- [`keycloak_ldap_group_mapper`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_ldap_group_mapper.go)
- [`keycloak_ldap_full_name_mapper`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_ldap_full_name_mapper.go)
- [`keycloak_ldap_msad_user_account_control_mapper`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_ldap_msad_user_account_control_mapper.go)
- [`keycloak_custom_user_federation`](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/provider/keycloak_custom_user_federation.go)

I will write some docs for each resource once more are supported. For now, please refer to the linked source files.

## Building

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) which requires Go 1.11.
I personally test the provider with version 0.11.8 of Terraform, and version 4.2.1.Final of Keycloak. Other versions may also work.

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
