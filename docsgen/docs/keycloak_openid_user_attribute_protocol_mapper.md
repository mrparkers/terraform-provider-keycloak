# keycloak_openid_user_attribute_protocol_mapper
Takes an attribute from a user and adds it as a claim on a JWT

### Example usage

```hcl
// client
resource "keycloak_openid_user_attribute_protocol_mapper" "map_user_attributes_client" {
    name           = "tf-test-open-id-user-attribute-protocol-mapper-client"
    realm_id       = "${keycloak_realm.test.id}"
    client_id      = "${keycloak_openid_client.test_client.id}"
    user_attribute = "foo"
    claim_name     = "bar"
}

// client scope
resource "keycloak_openid_user_attribute_protocol_mapper" "map_user_attributes_client_scope" {
    name            = "tf-test-open-id-user-attribute-protocol-mapper-client-scope"
    realm_id        = "${keycloak_realm.test.id}"
    client_scope_id = "${keycloak_openid_client_scope.test_client_scope.id}"
    user_attribute  = "foo2"
    claim_name      = "bar2"
}

```

### Argument Reference
The following arguments are supported:

- `add_to_user_info` - (Optional) Indicates if the attribute should appear in the userinfo response body.
- `user_attribute` - (Required) 
- `claim_name` - (Required) 
- `name` - (Required) A human-friendly name that will appear in the Keycloak console.
- `realm_id` - (Required) The realm id where the associated client or client scope exists.
- `client_id` - (Optional) The mapper's associated client. Cannot be used at the same time as client_scope_id.
- `add_to_id_token` - (Optional) Indicates if the attribute should be a claim in the id token.
- `client_scope_id` - (Optional) The mapper's associated client scope. Cannot be used at the same time as client_id.
- `add_to_access_token` - (Optional) Indicates if the attribute should be a claim in the access token.
- `multivalued` - (Optional) Indicates whether this attribute is a single value or an array of values.
- `claim_value_type` - (Optional) Claim type used when serializing tokens.



### Import
To import a mapper tied to a client, use the import command with the format `{{realmId}}/client/{{clientId}}/{{protocolMapperId}}`
Importing a mapper for a client scope is similar, `{{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}`


