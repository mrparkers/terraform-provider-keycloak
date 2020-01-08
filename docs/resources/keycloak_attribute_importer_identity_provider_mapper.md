# keycloak_attribute_importer_identity_provider_mapper

Allows to create and manage identity provider mappers within Keycloak.

### Example Usage

```hcl
resource "keycloak_attribute_importer_identity_provider_mapper" "test_mapper" {
  realm = "my-realm"
  name = "my-mapper"
  identity_provider_alias = "idp_alias"
  attribute_name = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"
  user_attribute = "lastName"
}
```

### Argument Reference

The following arguments are supported:

- `realm` - (Required) The name of the realm.
- `name` - (Required) The name of the mapper.
- `identity_provider_alias` - (Required) The alias of the associated identity provider.
- `user_attribute` - (Required) The user attribute name to store SAML attribute.
- `attribute_name` - (Optional) The Name of attribute to search for in assertion. You can leave this blank and specify a friendly name instead.
- `attribute_friendly_name` - (Optional) The friendly name of attribute to search for in assertion.  You can leave this blank and specify an attribute name instead.
- `claim_name` - (Optional) The claim name.

### Import

Identity provider mapper can be imported using the format `{{realm_id}}/{{idp_alias}}/{{idp_mapper_id}}`, where `idp_alias` is the identity provider alias, and `idp_mapper_id` is the unique ID that Keycloak
assigns to the mapper upon creation. This value can be found in the URI when editing this mapper in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_attribute_importer_identity_provider_mapper.test_mapper my-realm/my-mapper/f446db98-7133-4e30-b18a-3d28fde7ca1b
```
