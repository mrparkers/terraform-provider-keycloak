# keycloak_openid_audience_protocol_mapper

Allows for creating and managing audience protocol mappers within
Keycloak. This mapper was added in Keycloak v4.6.0.Final.

Audience protocol mappers allow you add audiences to the `aud` claim
within issued tokens. The audience can be a custom string, or it can be
mapped to the ID of a pre-existing client.

### Example Usage (Client)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
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

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
    realm_id                 = "${keycloak_realm.realm.id}"
    client_id                = "${keycloak_openid_client.openid_client.id}"
    name                     = "audience-mapper"

    included_custom_audience = "foo"
}
```

### Example Usage (Client Scope)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_openid_client_scope" "client_scope" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "test-client-scope"
}

resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
    realm_id                 = "${keycloak_realm.realm.id}"
    client_scope_id          = "${keycloak_openid_client_scope.client_scope.id}"
    name                     = "audience-mapper"

    included_custom_audience = "foo"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `client_id` - (Required if `client_scope_id` is not specified) The client this protocol mapper is attached to.
- `client_scope_id` - (Required if `client_id` is not specified) The client scope this protocol mapper is attached to.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `included_client_audience` - (Required if `included_custom_audience` is not specified) A client ID to include within the token's `aud` claim.
- `included_custom_audience` - (Required if `included_client_audience` is not specified) A custom audience to include within the token's `aud` claim.
- `add_to_id_token` - (Optional) Indicates if the audience should be included in the `aud` claim for the id token. Defaults to `true`.
- `add_to_access_token` - (Optional) Indicates if the audience should be included in the `aud` claim for the id token. Defaults to `true`.

### Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_openid_audience_protocol_mapper.audience_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_openid_audience_protocol_mapper.audience_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
