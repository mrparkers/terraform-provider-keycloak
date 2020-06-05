# keycloak_generic_client_role_mapper

Allow for creating and managing a client's scope mappings within Keycloak.

By default, all the user role mappings of the user are added as claims within
the token or assertion. When `full_scope_allowed` is set to `false` for a
client, role scope mapping allows you to limit the roles that get declared
inside an access token for a client.

### Example Usage (Realm Role to Client)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client" "client" {
  realm_id  = keycloak_realm.realm.id
  client_id = "client"
  name      = "client"

  enabled = true

  access_type = "BEARER-ONLY"
}

resource "keycloak_role" "realm_role" {
  realm_id    = keycloak_realm.realm.id
  name        = "my-realm-role"
  description = "My Realm Role"
}

resource "keycloak_generic_client_role_mapper" "client_role_mapper" {
  realm_id  = keycloak_realm.realm.id
  client_id = keycloak_openid_client.client.id
  role_id   = keycloak_role.realm_role.id
}
```

### Example Usage (Client Role to Client)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client" "client_a" {
  realm_id  = keycloak_realm.realm.id
  client_id = "client-a"
  name      = "client-a"

  enabled = true

  access_type = "BEARER-ONLY"
}

resource "keycloak_role" "client_role_a" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_openid_client.client_a.id
  name        = "my-client-role"
  description = "My Client Role"
}

resource "keycloak_openid_client" "client_b" {
  realm_id  = keycloak_realm.realm.id
  client_id = "client-b"
  name      = "client-b"

  enabled = true

  access_type = "BEARER-ONLY"
}

resource "keycloak_role" "client_role_b" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_openid_client.client_b.id
  name        = "my-client-role"
  description = "My Client Role"
}

resource "keycloak_generic_client_role_mapper" "client_b_role_mapper" {
  realm_id  = keycloak_realm.realm.id
  client_id = keycloak_client.client_b.id
  role_id   = keycloak_role.client_role_a.id
}
```

### Example Usage (Realm Role to Client Scope)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client_scope" "client_scope" {
  realm_id  = keycloak_realm.realm.id
  name      = "my-client-scope"
}

resource "keycloak_role" "realm_role" {
  realm_id    = keycloak_realm.realm.id
  name        = "my-realm-role"
  description = "My Realm Role"
}

resource "keycloak_generic_client_role_mapper" "client_role_mapper" {
  realm_id        = keycloak_realm.realm.id
  client_scope_id = keycloak_openid_client_scope.client_scope.id
  role_id         = keycloak_role.realm_role.id
}
```

### Example Usage (Client Role to Client Scope)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client" "client" {
  realm_id  = keycloak_realm.realm.id
  client_id = "client"
  name      = "client"

  enabled = true

  access_type = "BEARER-ONLY"
}

resource "keycloak_role" "client_role" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_openid_client.client.id
  name        = "my-client-role"
  description = "My Client Role"
}

resource "keycloak_openid_client_scope" "client_scope" {
  realm_id  = keycloak_realm.realm.id
  name      = "my-client-scope"
}

resource "keycloak_generic_client_role_mapper" "client_b_role_mapper" {
  realm_id        = keycloak_realm.realm.id
  client_scope_id = keycloak_client_scope.client_scope.id
  role_id         = keycloak_role.client_role.id
}
```

### Argument Reference

The following arugments are supported:

- `realm_id` - (Required) The realm this role mapper exists within
- `client_id` - (Optional) The ID of the client this role mapper is added to
- `client_scope_id` - (Optional) The ID of the client scope this role mapper is added to
- `role_id` - (Required) The ID of the role to be added to this role mapper

