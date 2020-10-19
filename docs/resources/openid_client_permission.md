# keycloak_openid_client_permissions

Allows you to manage all openid client Scope Based Permissions.

This is part of a preview keycloak feature. You need to enable this feature to be able to use this resource.
More information about enabling the preview feature can be found here: https://www.keycloak.org/docs/latest/securing_apps/index.html#_token-exchange

When enabling Openid Client Permissions, Keycloak does several things automatically:
1. Enable Authorization on build-in realm-management client
1. Create scopes "view", "manage", "configure", "map-roles", "map-roles-client-scope", "map-roles-composite", "token-exchange"
1. Create a resource representing the openid client
1. Create all scope based permission for the scopes and openid client resource

If the realm-management Authorization is not enable, you have to ceate a dependency (`depends_on`) with the policy and the openid client.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
	realm  = "realm"
}

resource "keycloak_openid_client" "my_openid_client" {
  realm_id              = keycloak_realm.realm.id
  name                  = "my_openid_client"
  client_id             = "my_openid_client"
  client_secret         = "secret"
  access_type           = "CONFIDENTIAL"
  standard_flow_enabled = true
  valid_redirect_uris   = [
    "http://localhost:8080/*",
  ]
}

data "keycloak_openid_client" "realm_management" {
	realm_id  = keycloak_realm.realm.id
	client_id = "realm-management"  
  }

resource keycloak_user test {
	realm_id = "${keycloak_realm.realm.id}"
	username = "test-user"

	email      = "test-user@fakedomain.com"
	first_name = "Testy"
	last_name  = "Tester"
}

resource keycloak_openid_client_user_policy test {
	resource_server_id = "${data.keycloak_openid_client.realm_management.id}"
	realm_id = "${keycloak_realm.realm.id}"
	name = "client_user_policy_test"
	users = ["${keycloak_user.test.id}"]
	logic = "POSITIVE"
	decision_strategy = "UNANIMOUS"
	depends_on = [
		keycloak_openid_client.my_openid_client
	]
}

resource "keycloak_openid_client_permissions" "my_permission" {
	realm_id                               = keycloak_realm.realm.id
	client_id                              = keycloak_openid_client.my_openid_client.id
	view_scope_policy_id                   = keycloak_openid_client_user_policy.test.id
	manage_scope_policy_id                 = keycloak_openid_client_user_policy.test.id
	configure_scope_policy_id              = keycloak_openid_client_user_policy.test.id
	map_roles_scope_policy_id              = keycloak_openid_client_user_policy.test.id
	map_roles_client_scope_scope_policy_id = keycloak_openid_client_user_policy.test.id
	map_roles_composite_scope_policy_id    = keycloak_openid_client_user_policy.test.id
	token_exchange_scope_policy_id         = keycloak_openid_client_user_policy.test.id
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this group exists in.
- `client_id` - (Required) The id of the client that provides the role.
- `view_scope_policy_id` - (Optional) Policy id that will be set on the scope based view permission automatically created by enabling permissions on the reference openid client.
- `manage_scope_policy_id` - (Optional) Policy id that will be set on the scope based manage permission automatically created by enabling permissions on the reference openid client.
- `configure_scope_policy_id` - (Optional) Policy id that will be set on the scope based configure permission automatically created by enabling permissions on the reference openid client.
- `map_roles_scope_policy_id` - (Optional) Policy id that will be set on the scope based map-roles permission automatically created by enabling permissions on the reference openid client.
- `map_roles_client_scope_scope_policy_id` - (Optional) Policy id that will be set on the scope based map-roles-client-scope permission automatically created by enabling permissions on the reference openid client.
- `map_roles_composite_scope_policy_id` - (Optional) Policy id that will be set on the scope based map-roles-composite permission automatically created by enabling permissions on the reference openid client.
- `token_exchange_scope_policy_id` - (Optional) Policy id that will be set on the scope based token-exchange permission automatically created by enabling permissions on the reference openid client.

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `authorization_resource_server_id` - Resource server id representing the realm management client on which this permission is managed.

