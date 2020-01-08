# keycloak_ldap_hardcoded_role_mapper

This mapper will grant a specified Keycloak role to each Keycloak user linked with LDAP.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "test"
    enabled = true
}

resource "keycloak_ldap_user_federation" "ldap_user_federation" {
    name                    = "openldap"
    realm_id                = keycloak_realm.realm.id

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

resource "keycloak_ldap_hardcoded_role_mapper" "assign_admin_role_to_all_users" {
  realm_id                = keycloak_realm.realm.id
  ldap_user_federation_id = keycloak_ldap_user_federation.ldap_user_federation.id
  name                    = "assign-admin-role-to-all-users"
  role                    = "admin"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm that this LDAP mapper will exist in.
- `ldap_user_federation_id` - (Required) The ID of the LDAP user federation provider to attach this mapper to.
- `name` - (Required) Display name of this mapper when displayed in the console.
- `role` - (Required) The role which should be assigned to the users.

### Import

LDAP mappers can be imported using the format `{{realm_id}}/{{ldap_user_federation_id}}/{{ldap_mapper_id}}`.
The ID of the LDAP user federation provider and the mapper can be found within
the Keycloak GUI, and they are typically GUIDs:

```bash
$ terraform import keycloak_ldap_hardcoded_role_mapper.ldap_hardcoded_role_mapper my-realm/af2a6ca3-e4d7-49c3-b08b-1b3c70b4b860/3d923ece-1a91-4bf7-adaf-3b82f2a12b67
```
