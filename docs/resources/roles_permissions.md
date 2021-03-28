---
page_title: "keycloak_role_permissions Resource"
---

# keycloak_role_permissions

Allows you to manage all role Scope Based Permissions https://www.keycloak.org/docs/latest/server_admin/#role.

This is part of a preview Keycloak feature: `admin_fine_grained_authz` (see https://www.keycloak.org/docs/latest/server_admin/#_fine_grain_permissions).
This feature can be enabled with the Keycloak option `-Dkeycloak.profile.feature.admin_fine_grained_authz=enabled`. See the
example [`docker-compose.yml`](https://github.com/mrparkers/terraform-provider-keycloak/blob/898094df6b3e01c3404981ce7ca268142d6ff0e5/docker-compose.yml#L21) file for an example.

When enabling Roles Permissions, Keycloak does several things automatically:
1. Enable Authorization on built-in `realm-management` client (if not already enabled).
1. Create a resource representing the role permissions.
1. Create scopes `map-role`, `map-role-client-scope`, `map-role-composite`.
1. Create all scope based permission for the scopes and role resource


### Example Usage

```hcl
resource "keycloak_realm" "realm" {
	realm  = "my-realm"
}

data "keycloak_openid_client" "realm_management" {
  realm_id  = keycloak_realm.realm.id
  client_id = "realm-management"  
}

// enable permissions for realm-management client
resource "keycloak_openid_client_permissions" "realm_management_permission" {
  realm_id   = keycloak_realm.realm.id
  client_id  = data.keycloak_openid_client.realm_management.id
  enabled    = true
}

// creating a user to use with the keycloak_openid_client_user_policy resource
resource "keycloak_user" "test" {
  realm_id = keycloak_realm.realm.id
  username = "test-user"

  email      = "test-user@fakedomain.com"
  first_name = "Testy"
  last_name  = "Tester"
}

resource "keycloak_openid_client_user_policy" "test" {
  realm_id           = keycloak_realm.realm.id
  resource_server_id = "${data.keycloak_openid_client.realm_management.id}"
  name               = "client_user_policy_test"

  users             = [keycloak_user.test.id]
  logic             = "POSITIVE"
  decision_strategy = "UNANIMOUS"

  depends_on = [
    keycloak_openid_client_permissions.realm-management_permission,
  ]
}

resource "keycloak_role" "role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role_permissions" "my_permission" {
	realm_id              = keycloak_realm.realm.id
	role_id               = keycloak_role.role.id

	map_roles_scope {
		policies          = [
		keycloak_openid_client_user_policy.test.id
		]
		description       = "description"
		decision_strategy = "UNANIMOUS"
	}
	map_role_client_scope_scope {
		policies          = [
		keycloak_openid_client_user_policy.test.id
		]
		description       = "description"
		decision_strategy = "UNANIMOUS"
	}
	map_role_composite_scope {
		policies          = [
		keycloak_openid_client_user_policy.test.id
		]
		description       = "description"
		decision_strategy = "UNANIMOUS"
	}
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm in which to manage fine-grained role permissions.
- `role_id` - (Required) The id of the role.


Each of the scopes that can be managed are defined below:

- `map_role_scope` - (Optional) When specified, set the scope based map-role permission.
- `map_role_client_scope_scope` - (Optional) When specified, set the scope based map-role-client-scope permission.
- `map_role_composite_scope` - (Optional) When specified, set the scope based map-role-composite permission.

The configuration block for each of these scopes supports the following arguments:

- `policies` - (Optional) Assigned policies to the permission. Each element within this list should be a policy ID.
- `description` - (Optional) Description of the permission.
- `decision_strategy` - (Optional) Decision strategy of the permission.

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `enabled` - When true, this indicates that fine-grained role permissions are enabled. This will always be `true`.
- `authorization_resource_server_id` - Resource server id representing the realm management client on which these permissions are managed.
