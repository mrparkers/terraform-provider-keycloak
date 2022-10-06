---
page_title: "keycloak_openid_client_scope Data Source"
---

# keycloak_openid_client_scope Data Source

This data source can be used to fetch properties of a Keycloak OpenID client scope for usage with other resources.

## Example Usage

```hcl
data "keycloak_openid_client_scope" "offline_access" {
  realm_id = "my-realm"
  name     = "offline_access"
}

# use the data source
resource "keycloak_openid_audience_protocol_mapper" "audience_mapper" {
	realm_id        = data.keycloak_openid_client_scope.offline_access.realm_id
	client_scope_id = data.keycloak_openid_client_scope.offline_access.id
	name            = "audience-mapper"

	included_custom_audience = "foo"
}
```

## Argument Reference

- `realm_id` - (Required) The realm id.
- `name` - (Required) The name of the client scope.

## Attributes Reference

See the docs for the `keycloak_openid_client_scope` resource for details on the exported attributes.
