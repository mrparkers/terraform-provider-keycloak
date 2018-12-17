# Installing

You can download the latest version of this provider on the
[GitHub releases](https://github.com/mrparkers/terraform-provider-keycloak/releases)
page.

Please follow the [official docs](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for instructions on installing a third-party provider.

# Keycloak Setup

Currently, this Terraform provider is configured to use the client credentials grant with
a client configured in the master realm. You can follow the steps below to configure a
client that the Terraform provider can use:

1. Create a client in the master realm using the `openid-connect` protocol.
2. Update the following client settings:
    - Set "Access Type" to "confidential".
    - Set "Standard Flow Enabled" to "OFF".
    - Set "Direct Access Grants Enabled" to "OFF".
    - Set "Service Accounts Enabled" to "ON".
3. Go to the "Service Account Roles" tab for the client, and grant it any roles that are
needed to manage your instance of Keycloak. The "admin" role can be assigned to effectively
manage all Keycloak settings.

# Provider Setup

The provider needs to be configured to use the master realm client configured in the
previous step. The following provider attributes are supported:

- `client_id` (Required) - The `client_id` for the client in the master realm setup in the previous step. Defaults to the environment variable `KEYCLOAK_CLIENT_ID`.
- `client_secret` (Required) - The secret for this client, which can be found or changed using the "Credentials" tab in the client settings. Defaults to the environment variable `KEYCLOAK_CLIENT_SECRET`.
- `url` (Required) - The URL of the Keycloak instance, before `/auth/admin`. Defaults to the environment variable `KEYCLOAK_URL`.

#### Example

```hcl
provider "keycloak" {
	client_id     = "terraform"
	client_secret = "884e0f95-0f42-4a63-9b1f-94274655669e"
	url           = "http://localhost:8080"
}
```
