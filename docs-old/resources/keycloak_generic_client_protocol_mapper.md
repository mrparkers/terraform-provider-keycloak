# keycloak_generic_client_protocol_mapper

Allows for creating and managing protocol mapper for both types of clients (openid-connect and saml) within Keycloak.

There are two uses cases for using this resource:
* If you implemented a custom protocol mapper, this resource can be used to configure it
* If the provider doesn't support a particular protocol mapper, this resource can be used instead.

Due to the generic nature of this mapper, it is less user-friendly and more prone to configuration errors.
Therefore, if possible, a specific mapper should be used.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_saml_client" "saml_client" {
  realm_id  = keycloak_realm.realm.id
  client_id = "test-client"
}

resource "keycloak_generic_client_protocol_mapper" "saml_hardcode_attribute_mapper" {
  realm_id        = keycloak_realm.realm.id
  client_id       = keycloak_saml_client.saml_client.id
  name            = "tes-mapper"
  protocol        = "saml"
  protocol_mapper = "saml-hardcode-attribute-mapper"
  config = {
    "attribute.name"       = "name"
    "attribute.nameformat" = "Basic"
    "attribute.value"      = "value"
    "friendly.name"        = "display name"
  }
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `client_id` - (Required) The client this protocol mapper is attached to.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `protocol` - (Required) The type of client (either `openid-connect` or `saml`). The type must match the type of the client.
- `protocol_mapper` - (Required) The name of the protocol mapper. The protocol mapper must be
   compatible with the specified client.
- `config` - (Required) A map with key / value pairs for configuring the protocol mapper. The supported keys depends on the protocol mapper.

### Import

Protocol mappers can be imported using the following format: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_generic_client_protocol_mapper.saml_hardcode_attribute_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
