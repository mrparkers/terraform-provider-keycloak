---
page_title: "keycloak_openid_client Data Source"
---

# keycloak\_openid\_client Data Source

This data source can be used to fetch properties of a Keycloak OpenID client for usage with other resources.

## Example Usage

```hcl
data "keycloak_openid_client" "realm_management" {
  realm_id  = "my-realm"
  client_id = "realm-management"
}

# use the data source
data "keycloak_role" "admin" {
  realm_id  = "my-realm"
  client_id = data.keycloak_openid_client.realm_management.id
  name      = "realm-admin"
}
```

## Argument Reference

- `realm_id` - (Required) The realm id.
- `client_id` - (Required) The client id (not its unique ID).

## Attributes Reference

- `name` - (Computed) The display name of this client in the GUI.
- `enabled` - (Computed) When false, this client will not be able to initiate a login or obtain access tokens.
- `description` - (Computed) The description of this client in the GUI.
- `access_type` - (Computed) The client's access type.
- `client_secret` - (Computed) The secret for clients with an `access_type` of `CONFIDENTIAL` or `BEARER-ONLY`. This value is sensitive and should be treated with the same care as a password.
- `standard_flow_enabled` - (Computed) When `true`, the OAuth2 Authorization Code Grant will be enabled for this client.
- `implicit_flow_enabled` - (Computed) When `true`, the OAuth2 Implicit Grant will be enabled for this client.
- `direct_access_grants_enabled` - (Computed) When `true`, the OAuth2 Resource Owner Password Grant will be enabled for this client.
- `service_accounts_enabled` - (Computed) When `true`, the OAuth2 Client Credentials grant will be enabled for this client.
- `valid_redirect_uris` - (Computed) A list of valid URIs a browser is permitted to redirect to after a successful login or logout.
- `web_origins` - (Computed) A list of allowed CORS origins.
- `root_url` - (Computed) This URL is prepended to any relative URLs found within `valid_redirect_uris`, `web_origins`, and `admin_url`.
- `admin_url` - (Computed) URL to the admin interface of the client.
- `base_url` - (Computed) Default URL to use when the auth server needs to redirect or link back to the client.
- `pkce_code_challenge_method` - (Computed) The challenge method to use for Proof Key for Code Exchange.
- `full_scope_allowed` - (Computed) When `true`, all roles mappings will be included in the access token.
- `access_token_lifespan` - (Computed) The amount of time in seconds before an access token expires.
- `consent_required` - (Computed) When `true`, users have to consent to client access.
- `authentication_flow_binding_overrides` - (Computed) Override realm authentication flow bindings.
- `login_theme` - (Computed) The client login theme.
- `exclude_session_state_from_auth_response` - (Computed) When `true`, the parameter `session_state` will not be included in OpenID Connect Authentication Response.
- `authorization` - (Computed) Client authorization settings.
