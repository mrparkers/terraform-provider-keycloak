---
page_title: "keycloak_openid_client_service_account_role Resource"
---

# keycloak\_openid\_client\_service\_account\_role Resource

Allows for assigning client roles to the service account of an openid client.
You need to set `service_accounts_enabled` to `true` for the openid client that should be assigned the role.

If you'd like to attach realm roles to a service account, please use the `keycloak_openid_client_service_account_realm_role`
resource.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

// client1 provides a role to other clients
resource "keycloak_openid_client" "client1" {
  realm_id = keycloak_realm.realm.id
  name     = "client1"
}

resource "keycloak_role" "client1_role" {
  realm_id    = keycloak_realm.realm.id
  client_id   = keycloak_openid_client.client1.id
  name        = "my-client1-role"
  description = "A role that client1 provides"
}

// client2 is assigned the role of client1
resource "keycloak_openid_client" "client2" {
  realm_id = keycloak_realm.realm.id
  name     = "client2"

  service_accounts_enabled = true
}

resource "keycloak_openid_client_service_account_role" "client2_service_account_role" {
  realm_id                = keycloak_realm.realm.id
  service_account_user_id = keycloak_openid_client.client2.service_account_user_id
  client_id               = keycloak_openid_client.client1.id
  role                    = keycloak_role.client1_role.name
}
```

## Argument Reference

- `realm_id` - (Required) The realm the clients and roles belong to.
- `service_account_user_id` - (Required) The id of the service account that is assigned the role (the service account of the client that "consumes" the role).
- `client_id` - (Required) The id of the client that provides the role.
- `role` - (Required) The name of the role that is assigned.

## Import

This resource can be imported using the format `{{realmId}}/{{serviceAccountUserId}}/{{clientId}}/{{roleId}}`.

Example:

```bash
$ terraform import keycloak_openid_client_service_account_role.client2_service_account_role my-realm/489ba513-1ceb-49ba-ae0b-1ab1f5099ebf/baf01820-0f8b-4494-9be2-fb3bc8a397a4/c7230ab7-8e4e-4135-995d-e81b50696ad8
```
