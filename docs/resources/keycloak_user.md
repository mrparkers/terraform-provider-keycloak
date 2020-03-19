# keycloak_user

Allows for creating and managing Users within Keycloak.

This resource was created primarily to enable the acceptance tests for the `keycloak_group` resource.
Creating users within Keycloak is not recommended. Instead, users should be federated from external sources
by configuring user federation providers or identity providers.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_user" "user" {
    realm_id   = "${keycloak_realm.realm.id}"
    username   = "bob"
    enabled    = true

    email      = "bob@domain.com"
    first_name = "Bob"
    last_name  = "Bobson"
}

resource "keycloak_user" "user_with_initial_password" {
    realm_id   = "${keycloak_realm.realm.id}"
    username   = "alice"
    enabled    = true

    email      = "alice@domain.com"
    first_name = "Alice"
    last_name  = "Aliceberg"

    initial_password {
      value     = "some password"
      temporary = true
    }
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this user belongs to.
- `username` - (Required) The unique username of this user.
- `initial_password` (Optional) When given, the user's initial password will be set.
   This attribute is only respected during initial user creation.
    - `value` (Required) The initial password.
    - `temporary` (Optional) If set to `true`, the initial password is set up for renewal on first use. Default to `false`.
- `enabled` - (Optional) When false, this user cannot log in. Defaults to `true`.
- `email` - (Optional) The user's email.
- `email_verified` - (Optional) Whether the email address was validated or not. Default to `false`.
- `first_name` - (Optional) The user's first name.
- `last_name` - (Optional) The user's last name.

### Import

Users can be imported using the format `{{realm_id}}/{{user_id}}`, where `user_id` is the unique ID that Keycloak
assigns to the user upon creation. This value can be found in the GUI when editing the user.

Example:

```bash
$ terraform import keycloak_user.user my-realm/60c3f971-b1d3-4b3a-9035-d16d7540a5e4
```
