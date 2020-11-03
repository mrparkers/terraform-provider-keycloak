---
page_title: "keycloak_saml_client Data Source"
---

# keycloak\_saml\_client Data Source

This data source can be used to fetch properties of a Keycloak client that uses the SAML protocol.

## Example Usage

```hcl
data "keycloak_saml_client" "realm_management" {
  realm_id  = "my-realm"
  client_id = "realm-management"
}

# use the data source
data "keycloak_role" "admin" {
  realm_id  = "my-realm"
  client_id = data.keycloak_saml_client.realm_management.id
  name      = "realm-admin"
}
```

## Argument Reference

- `realm_id` - (Required) The realm id.
- `client_id` - (Required) The client id (not its unique ID).

## Attributes Reference

See the docs for the `keycloak_saml_client` resource for details on the exported attributes.
