# keycloak_saml_user_property_protocol_mapper

Allows for creating and managing user property protocol mappers for
SAML clients within Keycloak.

SAML user property protocol mappers allow you to map properties of the Keycloak
user model to an attribute in a SAML assertion. Protocol mappers
can be defined for a single client, or they can be defined for a client scope which
can be shared between multiple different clients.

### Example Usage (Client)

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_saml_client" "saml_client" {
    realm_id  = "${keycloak_realm.test.id}"
    client_id = "test-saml-client"
    name      = "test-saml-client"
}

resource "keycloak_saml_user_property_protocol_mapper" "saml_user_property_mapper" {
    realm_id                   = "${keycloak_realm.test.id}"
    client_id                  = "${keycloak_saml_client.saml_client.id}"
    name                       = "email-user-property-mapper"

    user_property              = "email"
    saml_attribute_name        = "email"
    saml_attribute_name_format = "Unspecified"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `client_id` - (Required if `client_scope_id` is not specified) The SAML client this protocol mapper is attached to.
- `client_scope_id` - (Required if `client_id` is not specified) The SAML client scope this protocol mapper is attached to.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `user_property` - (Required) The property of the Keycloak user model to map.
- `friendly_name` - (Optional) An optional human-friendly name for this attribute.
- `saml_attribute_name` - (Required) The name of the SAML attribute.
- `saml_attribute_name_format` - (Required) The SAML attribute Name Format. Can be one of `Unspecified`, `Basic`, or `URI Reference`.

### Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_saml_user_property_protocol_mapper.saml_user_property_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_saml_user_property_protocol_mapper.saml_user_property_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
