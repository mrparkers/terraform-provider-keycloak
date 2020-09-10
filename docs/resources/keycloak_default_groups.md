# keycloak_default_groups

Allows for managing a realm's default groups.

Note that you should not use `keycloak_default_groups` with a group with memberships managed
by `keycloak_group_memberships`.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_group" "group" {
    realm_id = keycloak_realm.realm.id
    name     = "my-group"
}

resource "keycloak_default_groups" "default" {
    realm_id  = keycloak_realm.realm.id
    group_ids = [keycloak_group.group.id]
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this group exists in.
- `group_ids` - (Required) A set of group ids that should be default groups on the realm referenced by `realm_id`.

### Import

Groups can be imported using the format `{{realm_id}}` where `realm_id` is the realm the group exists in.

Example:

```bash
$ terraform import keycloak_default_groups.default my-realm
```
