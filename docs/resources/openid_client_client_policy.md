---
page_title: "keycloak_openid_client_client_policy Resource"
---

# keycloak\_openid\_client\_client_\_policy Resource

This resource can be used to create client policy.

## Example Usage

In this example, we'll create a new OpenID client, then enabled permissions for the client. A client without permissions disabled cannot be assigned by a client policy. We'll use the `keycloak_openid_client_client_policy` resource to create a new client policy, which could be applied to many clients, for a realm and a resource_server_id.

```hcl
resource "keycloak_realm" "realm" {
	realm   = "my-realm"
	enabled = true
}

resource "keycloak_openid_client" "openid_client" {
	client_id = "openid_client"
	name      = "openid_client"
	realm_id  = keycloak_realm.realm.id

	access_type              = "CONFIDENTIAL"
	service_accounts_enabled = true
}

resource "keycloak_openid_client_permissions" "my_permission" {
	realm_id  = keycloak_realm.realm.id
	client_id = keycloak_openid_client.openid_client.id
}

data "keycloak_openid_client" "realm_management" {
	realm_id  = "my-realm"
	client_id = "realm-management"
}

resource "keycloak_openid_client_client_policy" "token_exchange" {
	resource_server_id = data.keycloak_openid_client.realm_management.id
	realm_id           = keycloak_realm.realm.id
	name               = "my-policy"
	logic              = "POSITIVE"
	decision_strategy  = "UNANIMOUS"
	clients            = [
		keycloak_openid_client.openid_client.id
	]
}

```

## Argument Reference

- `resource_server_id` - (Required) The ID of the resource server this client policy is attached to.
- `realm_id` - (Required) The realm this client policy exists within.
- `name` - (Required) The name of this client policy.
- `clients` - (Required) The clients allowed by this client policy.
- `description` - (Optional) The description of this client policy.

## Attributes Reference

- `decision_strategy` - (Computed) Dictates how the policies associated with a given permission are evaluated and how a final decision is obtained. Could be one of `AFFIRMATIVE`, `CONSENSUS`, or `UNANIMOUS`. Applies to permissions.
- `logic` - (Computed) Dictates how the policy decision should be made. Can be either `POSITIVE` or `NEGATIVE`. Applies to policies.
