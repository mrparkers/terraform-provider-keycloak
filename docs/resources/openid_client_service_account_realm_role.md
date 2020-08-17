---
page_title: "keycloak_openid_client_service_account_realm_role Resource"
---

# keycloak\_openid\_client\_service\_account\_realm\_role Resource

Allows for assigning realm roles to the service account of an openid client.
You need to set `service_accounts_enabled` to `true` for the openid client that should be assigned the role.

If you'd like to attach client roles to a service account, please use the `keycloak_openid_client_service_account_role`
resource.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_role" "realm_role" {
  realm_id    = keycloak_realm.realm.id
  name        = "my-realm-role"
}

resource "keycloak_openid_client" "client" {
  realm_id = keycloak_realm.realm.id
  name     = "client"

  service_accounts_enabled = true
}

resource "keycloak_openid_client_service_account_realm_role" "client_service_account_role" {
  realm_id                = keycloak_realm.realm.id
  service_account_user_id = keycloak_openid_client.client.service_account_user_id
  role                    = keycloak_role.realm_role.name
}
```

## Argument Reference

- `realm_id` - (Required) The realm that the client and role belong to.
- `service_account_user_id` - (Required) The id of the service account that is assigned the role (the service account of the client that "consumes" the role).
- `role` - (Required) The name of the role that is assigned.

## Import

This resource can be imported using the format `{{realmId}}/{{serviceAccountUserId}}/{{roleId}}`.

Example:

```bash
$ terraform import keycloak_openid_client_service_account_realm_role.client_service_account_role my-realm/489ba513-1ceb-49ba-ae0b-1ab1f5099ebf/c7230ab7-8e4e-4135-995d-e81b50696ad8
```
