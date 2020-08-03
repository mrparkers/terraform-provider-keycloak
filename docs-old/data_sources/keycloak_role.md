# keycloak_role data source

This data source can be used to fetch properties of a Keycloak role for
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

# use the data source

resource "keycloak_group" "group" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "group"
}

resource "keycloak_group_roles" "group_roles" {
    realm_id = "${keycloak_realm.realm.id}"
    group_id = "${keycloak_group.group.id}"

    role_ids = [
        "${data.keycloak_role.offline_access.id}"
    ]
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this role exists within.
- `client_id` - (Optional) When specified, this role is assumed to be a
  client role belonging to the client with the provided ID
- `name` - (Required) The name of the role
  
### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `id` - The unique ID of the role, which can be used as an argument to
  other resources supported by this provider.
- `description` - The description of the role.
