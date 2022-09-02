---
page_title: "keycloak_openid_client_authorization_policy Data Source"
---

# keycloak\_openid\_client\_authorization\_policy Data Source

This data source can be used to fetch policy and permission information for an OpenID client that has authorization enabled.

## Example Usage

In this example, we'll create a new OpenID client with authorization enabled. This will cause Keycloak to create a default
permission for this client called "Default Permission". We'll use the `keycloak_openid_client_authorization_policy` data
source to fetch information about this permission, so we can use it to create a new resource-based authorization permission.

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client" "client_with_authz" {
  client_id = "client-with-authz"
  name      = "client-with-authz"
  realm_id  = keycloak_realm.realm.id

  access_type              = "CONFIDENTIAL"
  service_accounts_enabled = true

  authorization {
    policy_enforcement_mode = "ENFORCING"
  }
}

data "keycloak_openid_client_authorization_policy" "default_permission" {
  realm_id           = keycloak_realm.realm.id
  resource_server_id = keycloak_openid_client.client_with_authz.resource_server_id
  name               = "Default Permission"
}

resource "keycloak_openid_client_authorization_resource" "resource" {
  resource_server_id = keycloak_openid_client.client_with_authz.resource_server_id
  name               = "authorization-resource"
  realm_id           = keycloak_realm.realm.id

  uris = [
    "/endpoint/*",
  ]

  attributes = {
    "foo" = "bar"
  }
}

resource "keycloak_openid_client_authorization_permission" "permission" {
  resource_server_id = keycloak_openid_client.client_with_authz.resource_server_id
  realm_id           = keycloak_realm.realm.id
  name               = "authorization-permission"

  policies = [
    data.keycloak_openid_client_authorization_policy.default_permission.id,
  ]

  resources = [
    keycloak_openid_client_authorization_resource.resource.id,
  ]
}
```

## Argument Reference

- `realm_id` - (Required) The realm this authorization policy exists within.
- `name` - (Required) The name of the authorization policy.
- `resource_server_id` - (Required) The ID of the resource server this authorization policy is attached to.

## Attributes Reference

- `decision_strategy` - (Computed) Dictates how the policies associated with a given permission are evaluated and how a final decision is obtained. Could be one of `AFFIRMATIVE`, `CONSENSUS`, or `UNANIMOUS`. Applies to permissions.
- `owner` - (Computed) The ID of the owning resource. Applies to resources.
- `logic` - (Computed) Dictates how the policy decision should be made. Can be either `POSITIVE` or `NEGATIVE`. Applies to policies.
- `policies` - (Computed) The IDs of the policies that must be applied to scopes/resources for this policy/permission. Applies to policies and permissions.
- `resources` - (Computed) The IDs of the resources that this permission applies to. Applies to resource-based permissions.
- `scopes` - (Computed) The IDs of the scopes that this permission applies to. Applies to scope-based permissions.
- `type` - (Computed) The type of this policy / permission. For permissions, this could be `resource` or `scope`. For policies, this could be any type of authorization policy, such as `js`.
