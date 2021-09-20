---
page_title: "keycloak_user_realm_roles Data Source"
---

# keycloak_user_realm_roles Data Source

This data source can be used to fetch the realm roles of a user within Keycloak.

## Example Usage

```hcl
data "keycloak_realm" "master_realm" {
  realm = "master"
}

// use the keycloak_user data source to grab the admin user's ID
data "keycloak_user" "default_admin_user" {
  realm_id = data.keycloak_realm.master_realm.id
  username = "keycloak"
}

// use the keycloak_user_realm_roles data source to list role names
data "keycloak_user_realm_roles" "user_realm_roles" {
	realm_id = data.keycloak_realm.master_realm.id
	user_id  = data.keycloak_user.default_admin_user.id
}

output "keycloak_user_role_names" {
	value = data.keycloak_user_realm_roles.user_realm_roles.role_names
}
```

## Argument Reference

- `realm_id` - (Required) The realm this user belongs to.
- `user_id` - (Required) The ID of the user to query realm roles for.

## Attributes Reference

- `role_names` - (Computed) A list of realm roles that belong to this user.
