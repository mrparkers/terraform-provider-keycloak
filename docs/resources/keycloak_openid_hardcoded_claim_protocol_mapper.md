# keycloak_openid_hardcoded_claim_protocol_mapper

Allows for creating and managing hardcoded claim protocol mappers within
Keycloak.

Hardcoded claim protocol mappers allow you to define a claim with a hardcoded
value. Protocol mappers can be defined for a single client, or they can
be defined for a client scope which can be shared between multiple different
clients.

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

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper" {
    realm_id    = "${keycloak_realm.realm.id}"
    client_id   = "${keycloak_openid_client.openid_client.id}"
    name        = "hardcoded-claim-mapper"

    claim_name  = "foo"
    claim_value = "bar"
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

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "hardcoded_claim_mapper" {
    realm_id        = "${keycloak_realm.realm.id}"
    client_scope_id = "${keycloak_openid_client_scope.client_scope.id}"
    name            = "hardcoded-claim-mapper"

    claim_name      = "foo"
    claim_value     = "bar"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `client_id` - (Required if `client_scope_id` is not specified) The client this protocol mapper is attached to.
- `client_scope_id` - (Required if `client_id` is not specified) The client scope this protocol mapper is attached to.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `claim_name` - (Required) The name of the claim to insert into a token.
- `claim_value` - (Required) The hardcoded value of the claim.
- `claim_value_type` - (Optional) The claim type used when serializing JSON tokens. Can be one of `String`, `long`, `int`, or `boolean`. Defaults to `String`.
- `add_to_id_token` - (Optional) Indicates if the property should be added as a claim to the id token. Defaults to `true`.
- `add_to_access_token` - (Optional) Indicates if the property should be added as a claim to the access token. Defaults to `true`.
- `add_to_userinfo` - (Optional) Indicates if the property should be added as a claim to the UserInfo response body. Defaults to `true`.

### Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_openid_hardcoded_claim_protocol_mapper.hardcoded_claim_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
