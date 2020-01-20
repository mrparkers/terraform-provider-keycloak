# keycloak_group

Allows for creating and managing Groups within Keycloak.

Groups provide a logical wrapping for users within Keycloak. Users within a
group can share attributes and roles, and group membership can be mapped
to a claim.

Attributes can also be defined on Groups.

Groups can also be federated from external data sources, such as LDAP or Active Directory.
This resource **should not** be used to manage groups that were created this way.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_group" "parent_group" {
    realm_id = "${keycloak_realm.realm.id}"
    name     = "parent-group"
}

resource "keycloak_group" "child_group" {
    realm_id  = "${keycloak_realm.realm.id}"
    parent_id = "${keycloak_group.parent_group.id}"
    name      = "child-group"
}

resource "keycloak_group" "child_group_with_optional_attributes" {
    realm_id  = "${keycloak_realm.realm.id}"
    parent_id = "${keycloak_group.parent_group.id}"
    name      = "child-group-with-optional-attributes"
    attributes = {
		"key1" = "value1"
		"key2" = "value2"
    }
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this group exists in.
- `parent_id` - (Optional) The ID of this group's parent. If omitted, this group will be defined at the root level.
- `name` - (Required) The name of the group.
- `attributes` - (Optional) A dict of key/value pairs to set as custom attributes for the group.

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `path` - The complete path of the group. For example, the child group's path in the example configuration would be `/parent-group/child-group`.

### Import

Groups can be imported using the format `{{realm_id}}/{{group_id}}`, where `group_id` is the unique ID that Keycloak
assigns to the group upon creation. This value can be found in the URI when editing this group in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_group.child_group my-realm/934a4a4e-28bd-4703-a0fa-332df153aabd
```
