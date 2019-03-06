## 1.6.0 (March 6, 2019)

FEATURES:

* new resource: keycloak_openid_client_optional_scopes ([#96](https://github.com/mrparkers/terraform-provider-keycloak/pull/96))
* new resource: keycloak_openid_audience_protocol_mapper ([#97](https://github.com/mrparkers/terraform-provider-keycloak/pull/97))

## 1.5.0 (February 22, 2019)

FEATURES:

* adds support for non-master realms and resource owner password grant for Keycloak authentication ([#88](https://github.com/mrparkers/terraform-provider-keycloak/pull/88))

IMPROVEMENTS:

* support Keycloak v4.8.3.Final and Terraform v0.11.11 ([#93](https://github.com/mrparkers/terraform-provider-keycloak/pull/93))

BUG FIXES:

* handle 404 errors when reading a group for group memberships ([#95](https://github.com/mrparkers/terraform-provider-keycloak/pull/95))

## 1.4.0 (January 28, 2019)

FEATURES:

* new resource: keycloak_saml_user_property_protocol_mapper ([#85](https://github.com/mrparkers/terraform-provider-keycloak/pull/85))

## 1.3.0 (January 25, 2019)

FEATURES:

* new resource: keycloak_saml_user_attribute_protocol_mapper ([#84](https://github.com/mrparkers/terraform-provider-keycloak/pull/84))

## 1.2.0 (January 24, 2019)

FEATURES:

* new resource: keycloak_saml_client ([#82](https://github.com/mrparkers/terraform-provider-keycloak/pull/82))

IMPROVEMENTS:

* add validation for usernames to ensure they are always lowercase ([#83](https://github.com/mrparkers/terraform-provider-keycloak/pull/83))

## 1.1.0 (January 7, 2019)

IMPROVEMENTS:

* openid_client: add web_origins attribute ([#81](https://github.com/mrparkers/terraform-provider-keycloak/pull/81))
* user: add initial_password attribute ([#77](https://github.com/mrparkers/terraform-provider-keycloak/pull/77))

BUG FIXES:

* ldap mappers: don't assume component fields are returned by Keycloak API ([#80](https://github.com/mrparkers/terraform-provider-keycloak/pull/80))

## 1.0.0 (December 16, 2018)

Initial Release!

Docs: https://mrparkers.github.io/terraform-provider-keycloak
