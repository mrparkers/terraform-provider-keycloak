---
page_title: "keycloak_users_permissions Resource"
---

# keycloak_users_permissions

Allows you to manage all users Scope Based Permissions https://www.keycloak.org/docs/latest/server_admin/#_users-permissions

This is part of a preview keycloak feature `admin_fine_grained_authz` (see https://www.keycloak.org/docs/latest/server_admin/#_fine_grain_permissions)

You need to enable this feature to be able to use this resource.
More information about enabling the preview feature `admin_fine_grained_authz` can be found here: https://www.keycloak.org/docs/latest/server_installation/#profiles

When enabling Users Permissions, Keycloak does several things automatically:
1. Enable Authorization on build-in realm-management client (if not already enabled)
1. Create a resource representing the users permissions
1. Create scopes "view", "manage", "map-roles", "manage-group-membership", "impersonate", "user-impersonated"
1. Create all scope based permission for the scopes and users resource

If the realm-management Authorization is not enable, you have to ceate a dependency (`depends_on`) with the policy and the openid client.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
	realm  = "realm"
}

data "keycloak_openid_client" "realm_management" {
	realm_id  = keycloak_realm.realm.id
	client_id = "realm-management"  
}

resource keycloak_openid_client_permissions "realm-management_permission" {
	realm_id   = keycloak_realm.realm.id
	client_id  = data.keycloak_openid_client.realm_management.id
	enabled = true
}

resource keycloak_user test {
	realm_id = keycloak_realm.realm.id
	username = "test-user"

	email      = "test-user@fakedomain.com"
	first_name = "Testy"
	last_name  = "Tester"
}

resource keycloak_openid_client_user_policy test {
	resource_server_id = "${data.keycloak_openid_client.realm_management.id}"
	realm_id = keycloak_realm.realm.id
	name = "client_user_policy_test"
	users = [keycloak_user.test.id]
	logic = "POSITIVE"
	decision_strategy = "UNANIMOUS"
	depends_on = [
		keycloak_openid_client_permissions.realm-management_permission,
	]
}

resource "keycloak_users_permissions" "my_permission" {
	realm_id                                = keycloak_realm.realm.id
	
	view_scope {
		policies = [ keycloak_openid_client_user_policy.test.id ]
		description = "description"
		decision_strategy = "UNANIMOUS"
	}
	manage_scope {
		policies = [ keycloak_openid_client_user_policy.test.id ]
		description = "description"
		decision_strategy = "UNANIMOUS"
	}
	map_roles_scope {
		policies = [ keycloak_openid_client_user_policy.test.id ]
		description = "description"
		decision_strategy = "UNANIMOUS"
	}
	manage_group_membership_scope {
		policies = [ keycloak_openid_client_user_policy.test.id ]
		description = "description"
		decision_strategy = "UNANIMOUS"
	}
	impersonate_scope {
		policies = [ keycloak_openid_client_user_policy.test.id ]
		description = "description"
		decision_strategy = "UNANIMOUS"
	}
	user_impersonated_scope {
		policies = [ keycloak_openid_client_user_policy.test.id ]
		description = "description"
		decision_strategy = "UNANIMOUS"
	}
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this users permissions exists in.
- `view_scope` - (Optional) When specified, set the scope based view permission.
  - `policies` - (Optional) Assigned policies to the permission 
  - `description` - (Optional) Description of the permission 
  - `decision_strategy` - (Optional) Decision strategie of the permission 
- `manage_scope` - (Optional) When specified, set the scope based manage permission.
  - `policies` - (Optional) Assigned policies to the permission 
  - `description` - (Optional) Description of the permission 
  - `decision_strategy` - (Optional) Decision strategie of the permission 
- `map_roles_scope` - (Optional) When specified, set the scope based map_roles permission.
  - `policies` - (Optional) Assigned policies to the permission 
  - `description` - (Optional) Description of the permission 
  - `decision_strategy` - (Optional) Decision strategie of the permission 
- `manage_group_memberip_scope` - (Optional) When specified, set the scope based manage_group_memberip permission.
  - `policies` - (Optional) Assigned policies to the permission 
  - `description` - (Optional) Description of the permission 
  - `decision_strategy` - (Optional) Decision strategie of the permission 
- `impersonate_scope` - (Optional) When specified, set the scope based impersonate permission.
  - `policies` - (Optional) Assigned policies to the permission 
  - `description` - (Optional) Description of the permission 
  - `decision_strategy` - (Optional) Decision strategie of the permission 
- `user_impersonated_scope` - (Optional) When specified, set the scope based user_impersonated permission.
  - `policies` - (Optional) Assigned policies to the permission 
  - `description` - (Optional) Description of the permission 
  - `decision_strategy` - (Optional) Decision strategie of the permission 

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `enabled` - User permissions are Enabled (true) 
- `authorization_resource_server_id` - Resource server id representing the realm management client on which this permission is managed.

