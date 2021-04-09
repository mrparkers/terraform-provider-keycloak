---
page_title: "keycloak_user_groups Resource"
---

# keycloak\_user\_groups Resource

Allows for managing a Keycloak user's groups.

If `exhaustive` is true, this resource attempts to be an **authoritative** source over user groups: groups that are manually added to the user will be removed, and groups that are manually removed from the user group will be added upon the next run of `terraform apply`.
If `exhaustive` is false, this resource is a partial assignation of groups to a user. As a result, you can get multiple `keycloak_user_groups` for the same `user_id`.


## Example Usage (exhaustive groups)
```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_group" "group" {
  realm_id = keycloak_realm.realm.id
  name     = "foo"
}

resource "keycloak_user" "user" {
  realm_id = keycloak_realm.realm.id
  username = "my-user"
}

resource "keycloak_user_groups" "user_groups" {
  realm_id = keycloak_realm.realm.id
  user_id = keycloak_user.user.id

  group_ids  = [
    keycloak_group.group.id
  ]
}

```

## Example Usage (non exhaustive groups)
```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_group" "group_foo" {
  realm_id = keycloak_realm.realm.id
  name     = "foo"
}

resource "keycloak_group" "group_bar" {
  realm_id = keycloak_realm.realm.id
  name     = "bar"
}

resource "keycloak_user" "user" {
  realm_id = keycloak_realm.realm.id
  username = "my-user"
}

resource "keycloak_user_groups" "user_groups_association_1" {
  realm_id = keycloak_realm.realm.id
  user_id = keycloak_user.user.id
  exhaustive = false

  group_ids  = [
    keycloak_group.group_foo.id
  ]
}

resource "keycloak_user_groups" "user_groups_association_1" {
  realm_id = keycloak_realm.realm.id
  user_id = keycloak_user.user.id
  exhaustive = false

  group_ids  = [
    keycloak_group.group_bar.id
  ]
}
```

## Argument Reference

- `realm_id` - (Required) The realm this group exists in.
- `user_id` - (Required) The ID of the user this resource should manage groups for.
- `group_ids` - (Required) A list of group IDs that the user is member of.
- `exhaustive` - (Optional)

## Import

This resource does not support import. Instead of importing, feel free to create this resource
as if it did not already exist on the server.

