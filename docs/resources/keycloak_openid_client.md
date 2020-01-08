# keycloak_openid_client

Allows for creating and managing Keycloak clients that use the OpenID Connect protocol.

Clients are entities that can use Keycloak for user authentication. Typically,
clients are applications that redirect users to Keycloak for authentication
in order to take advantage of Keycloak's user sessions for SSO.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_openid_client" "openid_client" {
    realm_id            = "${keycloak_realm.realm.id}"
    client_id           = "test-client"

    name                = "test client"
    enabled             = true

    access_type         = "CONFIDENTIAL"
    valid_redirect_uris = [
        "http://localhost:8080/openid-callback"
    ]
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this client is attached to.
- `client_id` - (Required) The unique ID of this client, referenced in the URI during authentication and in issued tokens.
- `name` - (Optional) The display name of this client in the GUI.
- `enabled` - (Optional) When false, this client will not be able to initiate a login or obtain access tokens. Defaults to `true`.
- `description` - (Optional) The description of this client in the GUI.
- `access_type` - (Required) Specifies the type of client, which can be one of the following:
    - `CONFIDENTIAL` - Used for server-side clients that require both client ID and secret when authenticating.
      This client should be used for applications using the Authorization Code or Client Credentials grant flows.
    - `PUBLIC` - Used for browser-only applications that do not require a client secret, and instead rely only on authorized redirect
      URIs for security. This client should be used for applications using the Implicit grant flow.
    - `BEARER-ONLY` - Used for services that never initiate a login. This client will only allow bearer token requests.
- `client_secret` - (Optional) The secret for clients with an `access_type` of `CONFIDENTIAL` or `BEARER-ONLY`. This value is sensitive and
should be treated with the same care as a password. If omitted, Keycloak will generate a GUID for this attribute.
- `standard_flow_enabled` - (Optional) When `true`, the OAuth2 Authorization Code Grant will be enabled for this client. Defaults to `false`.
- `implicit_flow_enabled` - (Optional) When `true`, the OAuth2 Implicit Grant will be enabled for this client. Defaults to `false`.
- `direct_access_grants_enabled` - (Optional) When `true`, the OAuth2 Resource Owner Password Grant will be enabled for this client. Defaults to `false`.
- `service_accounts_enabled` - (Optional) When `true`, the OAuth2 Client Credentials grant will be enabled for this client. Defaults to `false`.
- `valid_redirect_uris` - (Optional) A list of valid URIs a browser is permitted to redirect to after a successful login or logout. Simple
wildcards in the form of an asterisk can be used here. This attribute must be set if either `standard_flow_enabled` or `implicit_flow_enabled`
is set to `true`.
- `web_origins` - (Optional) A list of allowed CORS origins. `+` can be used to permit all valid redirect URIs, and `*` can be used to permit all origins.
- `pkce_code_challenge_method` - (Optional) The challenge method to use for Proof Key for Code Exchange. Can be either `plain` or `S256`.
- `full_scope_allowed` - (Optional) - Allow to include all roles mappings in the access token.

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `service_account_user_id` - When service accounts are enabled for this client, this attribute is the unique ID for the Keycloak user that represents this service account.
 

### Import

Clients can be imported using the format `{{realm_id}}/{{client_keycloak_id}}`, where `client_keycloak_id` is the unique ID that Keycloak
assigns to the client upon creation. This value can be found in the URI when editing this client in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_openid_client.openid_client my-realm/dcbc4c73-e478-4928-ae2e-d5e420223352
```
