---
page_title: "keycloak_authentication_flow Data Source"
---

# keycloak\_authentication\_flow Data Source

This data source can be used to fetch the ID of an authentication flow within Keycloak.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

data "keycloak_authentication_flow" "browser_auth_cookie" {
  realm_id          = keycloak_realm.realm.id
  alias             = "browser"
}
```

## Argument Reference

- `realm_id` - (Required) The realm the authentication flow exists in.
- `alias` - (Required) The alias of the flow.

## Attributes Reference

- `id` - (Computed) The unique ID of the authentication flow, which can be used as an argument to other resources supported by this provider.
