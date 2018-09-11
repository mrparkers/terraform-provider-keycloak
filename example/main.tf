provider "keycloak" {
  client_id     = "terraform"
  client_secret = "884e0f95-0f42-4a63-9b1f-94274655669e"
  url           = "http://localhost:8080"
}

resource "keycloak_realm" "test" {
  realm        = "test"
  enabled      = true
  display_name = "foo"
}

resource "keycloak_client" "test-client" {
  client_id = "test-client"
  realm_id  = "${keycloak_realm.test.id}"
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

resource "keycloak_ldap_user_attribute_mapper" "test-mapper" {
  name                    = "test mapper"
  realm_id                = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id = "${keycloak_ldap_user_federation.openldap.id}"

  user_model_attribute    = "foo"
  ldap_attribute          = "bar"
}
