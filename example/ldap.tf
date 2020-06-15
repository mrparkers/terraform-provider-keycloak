resource "keycloak_realm" "ldap" {
  realm   = "ldap"
  enabled = true
}

resource "keycloak_ldap_user_federation" "openldap" {
  name     = "openldap"
  realm_id = keycloak_realm.ldap.id

  enabled        = true
  import_enabled = false

  username_ldap_attribute = "cn"
  rdn_ldap_attribute      = "cn"
  uuid_ldap_attribute     = "entryDN"

  user_object_classes = [
    "simpleSecurityObject",
    "organizationalRole",
  ]

  connection_url  = "ldap://openldap"
  users_dn        = "dc=example,dc=org"
  bind_dn         = "cn=admin,dc=example,dc=org"
  bind_credential = "admin"

  connection_timeout = "5s"
  read_timeout       = "10s"

  kerberos {
    server_principal                         = "HTTP/keycloak.local@FOO.LOCAL"
    use_kerberos_for_password_authentication = false
    key_tab                                  = "/etc/keycloak.keytab"
    kerberos_realm                           = "FOO.LOCAL"
  }

  cache_policy = "NO_CACHE"
}

resource "keycloak_ldap_user_attribute_mapper" "description_attr_mapper" {
  name                    = "description-mapper"
  realm_id                = keycloak_ldap_user_federation.openldap.realm_id
  ldap_user_federation_id = keycloak_ldap_user_federation.openldap.id

  user_model_attribute = "description"
  ldap_attribute       = "description"

  always_read_value_from_ldap = false
}

resource "keycloak_ldap_group_mapper" "group_mapper" {
  name                    = "group mapper"
  realm_id                = keycloak_ldap_user_federation.openldap.realm_id
  ldap_user_federation_id = keycloak_ldap_user_federation.openldap.id

  ldap_groups_dn            = "dc=example,dc=org"
  group_name_ldap_attribute = "cn"

  group_object_classes = [
    "groupOfNames",
  ]

  membership_attribute_type      = "DN"
  membership_ldap_attribute      = "member"
  membership_user_ldap_attribute = "cn"
  memberof_ldap_attribute        = "memberOf"
}

resource "keycloak_ldap_msad_user_account_control_mapper" "msad_uac_mapper" {
  name                    = "uac-mapper1"
  realm_id                = keycloak_ldap_user_federation.openldap.realm_id
  ldap_user_federation_id = keycloak_ldap_user_federation.openldap.id
}

resource "keycloak_ldap_msad_lds_user_account_control_mapper" "msad_lds_uac_mapper" {
  name                    = "msad-lds-uac-mapper"
  realm_id                = keycloak_ldap_user_federation.openldap.realm_id
  ldap_user_federation_id = keycloak_ldap_user_federation.openldap.id
}

resource "keycloak_ldap_full_name_mapper" "full_name_mapper" {
  name                    = "full-name-mapper"
  realm_id                = keycloak_ldap_user_federation.openldap.realm_id
  ldap_user_federation_id = keycloak_ldap_user_federation.openldap.id

  ldap_full_name_attribute = "cn"
  read_only                = true
}
