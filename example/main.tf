provider "keycloak" {
  client_id     = "terraform"
  client_secret = "884e0f95-0f42-4a63-9b1f-94274655669e"
  url           = "http://localhost:8080"
}

resource "keycloak_realm" "test" {
  realm                = "test"
  enabled              = true
  display_name         = "foo"

  account_theme        = "base"

  access_code_lifespan = "30m"
}

resource "keycloak_openid_client" "test-client" {
  client_id = "test-client"
  realm_id  = "${keycloak_realm.test.id}"
}

resource "keycloak_client_scope" "test-client-scope" {
  name                = "foo1"
  realm_id            = "${keycloak_realm.test.id}"

  description         = "test"
  protocol            = "saml"
  consent_screen_text = "hello"
}

resource "keycloak_ldap_user_federation" "openldap" {
  name                    = "openldap"
  realm_id                = "master"

  enabled                 = true

  username_ldap_attribute = "cn"
  rdn_ldap_attribute      = "cn"
  uuid_ldap_attribute     = "entryDN"
  user_object_classes     = [
    "simpleSecurityObject",
    "organizationalRole"
  ]
  connection_url          = "ldap://openldap"
  users_dn                = "dc=example,dc=org"
  bind_dn                 = "cn=admin,dc=example,dc=org"
  bind_credential         = "admin"
}

resource "keycloak_ldap_user_attribute_mapper" "attr-mapper" {
  name                    = "test mapper"
  realm_id                = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id = "${keycloak_ldap_user_federation.openldap.id}"

  user_model_attribute    = "foo"
  ldap_attribute          = "bar"
}

resource "keycloak_ldap_group_mapper" "group-mapper" {
  name                           = "group mapper"
  realm_id                       = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id        = "${keycloak_ldap_user_federation.openldap.id}"

  ldap_groups_dn                 = "dc=example,dc=org"
  group_name_ldap_attribute      = "cn"
  group_object_classes           = [
    "groupOfNames"
  ]
  membership_attribute_type      = "DN"
  membership_ldap_attribute      = "member"
  membership_user_ldap_attribute = "cn"
  memberof_ldap_attribute        = "memberOf"
}

resource "keycloak_ldap_msad_user_account_control_mapper" "msad_uac_mapper" {
  name                    = "uac-mapper1"
  realm_id                = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id = "${keycloak_ldap_user_federation.openldap.id}"
}

resource "keycloak_ldap_full_name_mapper" "full_name_mapper" {
  name                     = "full-name-mapper"
  realm_id                 = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id  = "${keycloak_ldap_user_federation.openldap.id}"

  ldap_full_name_attribute = "cn"
  read_only                = true
}
