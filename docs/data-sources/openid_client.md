---
page_title: "keycloak_openid_client Data Source"
---

# keycloak\_openid\_client Data Source

This data source can be used to fetch properties of a Keycloak OpenID client for usage with other resources.

## Example Usage

```hcl
data "keycloak_openid_client" "realm_management" {
  realm_id  = "my-realm"
  client_id = "realm-management"
}

# use the data source
data "keycloak_role" "admin" {
  realm_id  = "my-realm"
  client_id = data.keycloak_openid_client.realm_management.id
  name      = "realm-admin"
}
```

## Argument Reference

- `realm_id` - (Required) The realm id.
- `client_id` - (Required) The client id (not its unique ID).

## Attributes Reference

See the docs for the `keycloak_openid_client` resource for details on the exported attributes.
