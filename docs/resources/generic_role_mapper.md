---
page_title: "keycloak_generic_role_mapper Resource"
---

# keycloak\_generic\_role\_mapper Resource

Allow for creating and managing a client's or client scope's role mappings within Keycloak.

By default, all the user role mappings of the user are added as claims within the token (OIDC) or assertion (SAML). When
`full_scope_allowed` is set to `false` for a client, role scope mapping allows you to limit the roles that get declared
inside an access token for a client.

## Example Usage (Realm Role to Client)

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

resource "keycloak_generic_role_mapper" "client_role_mapper" {
  realm_id  = keycloak_realm.realm.id
  client_id = keycloak_openid_client.client.id
  role_id   = keycloak_role.realm_role.id
}
```

## Example Usage (Client Role to Client)

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

  // disable full scope, roles are assigned via keycloak_generic_role_mapper
  full_scope_allowed = false
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

resource "keycloak_generic_role_mapper" "client_b_role_mapper" {
  realm_id  = keycloak_realm.realm.id
  client_id = keycloak_openid_client.client_b.id
  role_id   = keycloak_role.client_role_a.id
}
```

## Example Usage (Realm Role to Client Scope)

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

resource "keycloak_generic_role_mapper" "client_role_mapper" {
  realm_id        = keycloak_realm.realm.id
  client_scope_id = keycloak_openid_client_scope.client_scope.id
  role_id         = keycloak_role.realm_role.id
}
```

## Example Usage (Client Role to Client Scope)

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

resource "keycloak_generic_role_mapper" "client_b_role_mapper" {
  realm_id        = keycloak_realm.realm.id
  client_scope_id = keycloak_openid_client_scope.client_scope.id
  role_id         = keycloak_role.client_role.id
}
```

## Argument Reference

- `realm_id` - (Required) The realm this role mapper exists within.
- `client_id` - (Optional) The ID of the client this role mapper should be added to. Conflicts with `client_scope_id`. This argument is required if `client_scope_id` is not set.
- `client_scope_id` - (Optional) The ID of the client scope this role mapper should be added to. Conflicts with `client_id`. This argument is required if `client_id` is not set.
- `role_id` - (Required) The ID of the role to be added to this role mapper.

## Import

Generic client role mappers can be imported using one of the following two formats:

- When mapping a role to a client, use the format `{{realmId}}/client/{{clientId}}/scope-mappings/{{roleClientId}}/{{roleId}}`
- When mapping a role to a client scope, use the format `{{realmId}}/client-scope/{{clientScopeId}}/scope-mappings/{{roleClientId}}/{{roleId}}`

Example:

```bash
$ terraform import keycloak_generic_role_mapper.client_role_mapper my-realm/client/23888550-5dcd-41f6-85ba-554233021e9c/scope-mappings/ce51f004-bdfb-4dd5-a963-c4487d2dec5b/ff3aa49f-bc07-4030-8783-41918c3614a3
```
