# keycloak_openid_hardcoded_role_protocol_mapper

Allows for creating and managing hardcoded role protocol mappers within
Keycloak.

Hardcoded role protocol mappers allow you to specify a single role to
always map to an access token for a client. Protocol mappers can be
defined for a single client, or they can be defined for a client scope
which can be shared between multiple different clients.

### Example Usage (Client)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_role" "role" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "my-role"
}

resource "keycloak_openid_client" "openid_client" {
    realm_id            = "${keycloak_realm.realm.id}"
    client_id           = "test-client"

    name                = "test client"
    enabled             = true

    access_type         = "CONFIDENTIAL"
    valid_redirect_uris = [
        "http://localhost:8080/openid-callback"
    ]
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper" {
    realm_id  = "${keycloak_realm.realm.id}"
    client_id = "${keycloak_openid_client.openid_client.id}"
    name      = "hardcoded-role-mapper"
    role_id   = "${keycloak_role.role.id}"
}
```

### Example Usage (Client Scope)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_role" "role" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "my-role"
}

resource "keycloak_openid_client_scope" "client_scope" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "test-client-scope"
}

resource "keycloak_openid_hardcoded_role_protocol_mapper" "hardcoded_role_mapper" {
    realm_id        = "${keycloak_realm.realm.id}"
    client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
    name            = "hardcoded-role-mapper"
    role_id         = "${keycloak_role.role.id}"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `client_id` - (Required if `client_scope_id` is not specified) The client this protocol mapper is attached to.
- `client_scope_id` - (Required if `client_id` is not specified) The client scope this protocol mapper is attached to.
- `name` - (Required) The display name of this protocol mapper in the
  GUI.
- `role_id` - (Required) The ID of the role to map to an access token.

### Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_openid_hardcoded_role_protocol_mapper.hardcoded_role_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
