---
page_title: "keycloak_saml_user_property_protocol_mapper Resource"
---

# keycloak\_saml\_user\_property\_protocol\_mapper Resource

Allows for creating and managing user property protocol mappers for SAML clients within Keycloak.

SAML user property protocol mappers allow you to map properties of the Keycloak
user model to an attribute in a SAML assertion.

Protocol mappers can be defined for a single client, or they can be defined for a client scope which can be shared between
multiple different clients.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_saml_client" "saml_client" {
  realm_id  = keycloak_realm.test.id
  client_id = "saml-client"
  name      = "saml-client"
}

resource "keycloak_saml_user_property_protocol_mapper" "saml_user_property_mapper" {
  realm_id  = keycloak_realm.test.id
  client_id = keycloak_saml_client.saml_client.id
  name      = "email-user-property-mapper"

  user_property              = "email"
  saml_attribute_name        = "email"
  saml_attribute_name_format = "Unspecified"
}
```

## Argument Reference

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `user_property` - (Required) The property of the Keycloak user model to map.
- `saml_attribute_name` - (Required) The name of the SAML attribute.
- `saml_attribute_name_format` - (Required) The SAML attribute Name Format. Can be one of `Unspecified`, `Basic`, or `URI Reference`.
- `client_id` - (Optional) The client this protocol mapper should be attached to. Conflicts with `client_scope_id`. One of `client_id` or `client_scope_id` must be specified.
- `client_scope_id` - (Optional) The client scope this protocol mapper should be attached to. Conflicts with `client_id`. One of `client_id` or `client_scope_id` must be specified.
- `friendly_name` - (Optional) An optional human-friendly name for this attribute.

## Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_saml_user_property_protocol_mapper.saml_user_property_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_saml_user_property_protocol_mapper.saml_user_property_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
