---
page_title: "keycloak_group_roles Resource"
---

# keycloak\_group\_roles Resource

Allows you to manage roles assigned to a Keycloak group.

Note that this resource attempts to be an **authoritative** source over group roles. When this resource takes control over
a group's roles, roles that are manually added to the group will be removed, and roles that are manually removed from the
group will be added upon the next run of `terraform apply`.

Note that when assigning composite roles to a group, you may see a non-empty plan following a `terraform apply` if you
assign a role and a composite that includes that role to the same group.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_role" "realm_role" {
  realm_id    = keycloak_realm.realm.id
  name        = "my-realm-role"
  description = "My Realm Role"
}

resource "keycloak_openid_client" "client" {
  realm_id  = keycloak_realm.realm.id
  client_id = "client"
  name      = "client"

  enabled = true

  access_type = "BEARER-ONLY"
}

resource "keycloak_role" "client_role" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_client.client.id
  name        = "my-client-role"
  description = "My Client Role"
}

resource "keycloak_group" "group" {
  realm_id = keycloak_realm.realm.id
  name     = "my-group"
}

resource "keycloak_group_roles" "group_roles" {
  realm_id = keycloak_realm.realm.id
  group_id = keycloak_group.group.id

  role_ids = [
    keycloak_role.realm_role.id,
    keycloak_role.client_role.id,
  ]
}
```

## Argument Reference

- `realm_id` - (Required) The realm this group exists in.
- `group_id` - (Required) The ID of the group this resource should manage roles for.
- `role_ids` - (Required) A list of role IDs to map to the group

## Import

This resource can be imported using the format `{{realm_id}}/{{group_id}}`, where `group_id` is the unique ID that Keycloak
assigns to the group upon creation. This value can be found in the URI when editing this group in the GUI, and is typically
a GUID.

Example:

```bash
$ terraform import keycloak_group_roles.group_roles my-realm/18cc6b87-2ce7-4e59-bdc8-b9d49ec98a94
```
