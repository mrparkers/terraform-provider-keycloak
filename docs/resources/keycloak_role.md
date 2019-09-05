# keycloak_role

Allows for creating and managing roles within Keycloak.

Roles allow you define privileges within Keycloak and map them to users
and groups.

### Example Usage (Realm role)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_role" "realm_role" {
    realm_id    = "${keycloak_realm.realm.id}"
    name        = "my-realm-role"
    description = "My Realm Role"
}
```

### Example Usage (Client role)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_openid_client" "client" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "client"
  name      = "client"

  enabled = true

  access_type = "BEARER-ONLY"
}

resource "keycloak_role" "client_role" {
    realm_id    = "${keycloak_realm.realm.id}"
    client_id   = "${keycloak_client.client.id}"
    name        = "my-client-role"
    description = "My Client Role"
}
```

### Example Usage (Composite role)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

# realm roles

resource "keycloak_role" "create_role" {
    realm_id    = "${keycloak_realm.realm.id}"
    name        = "create"
}

resource "keycloak_role" "read_role" {
    realm_id    = "${keycloak_realm.realm.id}"
    name        = "read"
}

resource "keycloak_role" "update_role" {
    realm_id    = "${keycloak_realm.realm.id}"
    name        = "update"
}

resource "keycloak_role" "delete_role" {
    realm_id    = "${keycloak_realm.realm.id}"
    name        = "delete"
}

# client role

resource "keycloak_openid_client" "client" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "client"
  name      = "client"

  enabled = true

  access_type = "BEARER-ONLY"
}

resource "keycloak_role" "client_role" {
    realm_id    = "${keycloak_realm.realm.id}"
    client_id   = "${keycloak_client.client.id}"
    name        = "my-client-role"
    description = "My Client Role"
}

resource "keycloak_role" "admin_role" {
    realm_id        = "${keycloak_realm.realm.id}"
    name            = "admin"
    composite_roles = [
      "{keycloak_role.create_role.id}",
      "{keycloak_role.read_role.id}",
      "{keycloak_role.update_role.id}",
      "{keycloak_role.delete_role.id}",
      "{keycloak_role.client_role.id}",
    ]
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this role exists within.
- `client_id` - (Optional) When specified, this role will be created as
  a client role attached to the client with the provided ID
- `name` - (Required) The name of the role
- `description` - (Optional) The description of the role
- `composite_roles` - (Optional) When specified, this role will be a
  composite role, composed of all roles that have an ID present within
  this list.


### Import

Roles can be imported using the format `{{realm_id}}/{{role_id}}`, where
`role_id` is the unique ID that Keycloak assigns to the role. The ID is
not easy to find in the GUI, but it appears in the URL when editing the
role.

Example:

```bash
$ terraform import keycloak_role.role my-realm/7e8cf32a-8acb-4d34-89c4-04fb1d10ccad
```
