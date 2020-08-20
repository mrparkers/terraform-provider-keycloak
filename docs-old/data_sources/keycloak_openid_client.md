# keycloak_openid_client data source

This data source can be used to fetch properties of a Keycloak OpenID client for usage with other resources.

### Example Usage

```hcl
data "keycloak_openid_client" "realm_management" {
  realm_id = "my-realm"
  client_id = "realm-management"
}

# use the data source
data "keycloak_role" "admin" {
  realm_id = "my-realm"
  client_id = data.keycloak_openid_client.realm_management.id
  name = "realm-admin"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm id.
- `client_id` - (Required) The client id.

### Attributes Reference

See the docs for the [`keycloak_openid_client` resource](../resources/keycloak_openid_client.md) for details on the exported attributes.
