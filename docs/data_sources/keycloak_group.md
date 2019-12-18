# keycloak_group data source

This data source can be used to fetch properties of a Keycloak group for
usage with other resources, such as `keycloak_group_roles`.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

data "keycloak_role" "offline_access" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "offline_access"
}

data "keycloak_group" "group" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "group"
}

resource "keycloak_group_roles" "group_roles" {
    realm_id = "${keycloak_realm.realm.id}"
    group_id = "${data.keycloak_group.group.id}"

    roles = [
        "${data.keycloak_role.offline_access.id}"
    ]
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this group exists within.
- `name` - (Required) The name of the group

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `id` - The unique ID of the group, which can be used as an argument to
  other resources supported by this provider.

