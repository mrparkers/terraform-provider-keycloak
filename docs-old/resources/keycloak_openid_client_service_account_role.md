# keycloak_openid_client_service_account_role

Allows for assigning roles to the service account of an openid client.

You need to set `service_accounts_enabled` to `true` for the openid client that should be assigned the role.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

// client1 provides a role to other clients
resource "keycloak_openid_client" "client1" {
    realm_id  = keycloak_realm.realm.id
    name    = "client1"
}

resource "keycloak_role" "client1_role" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_openid_client.client1.id
  name        = "my-client1-role"
  description = "A role that client1 provides"
}

// client2 is assigned the role of client1
resource "keycloak_openid_client" "client2" {
    realm_id  = keycloak_realm.realm.id
    name    = "client2"
    service_accounts_enabled = true
}

resource "keycloak_openid_client_service_account_role" "client2_service_account_role" {
  realm_id                = keycloak_realm.realm.id
  service_account_user_id = keycloak_openid_client.client2.service_account_user_id
  client_id               = keycloak_openid_client.client1.id
  role                    = keycloak_role.client1_role.name
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm the clients and roles belong to.
- `service_account_user_id` - (Required) The id of the service account that is assigned the role (the service account of the client that "consumes" the role).
- `client_id` - (Required) The id of the client that provides the role.
- `role` - (Required) The name of the role that is assigned.
