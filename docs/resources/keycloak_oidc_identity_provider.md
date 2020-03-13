# keycloak_oidc_identity_provider

Allows to create and manage OIDC Identity Providers within Keycloak.

OIDC (OpenID Connect) identity providers allows to authenticate through a third-party system, using OIDC standard.

### Example Usage

```hcl
resource "keycloak_realm" "my-realm" {
  realm        = "my-realm"
  enabled      = true
  display_name = "my-realm"
}

resource "keycloak_oidc_identity_provider" "realm_identity_provider" {
  realm             = "my-realm"
  alias             = "my-idp"
  authorization_url = "https://authorizationurl.com"
  client_id         = "clientID"
  client_secret     = "clientSecret" # or "$${vault.ID}"
  token_url         = "https://tokenurl.com"

  extra_config = {
    "clientAuthMethod"                   = "client_secret_post"
  }
}
```

### Argument Reference

The following arguments are supported:

- `realm` - (Required) The name of the realm. This is unique across Keycloak.
- `alias` - (Required) The alias uniquely identifies an identity provider and it is also used to build the redirect uri.
- `authorization_url` - (Required) The Authorization Url.
- `client_id` - (Required) The client or client identifier registered within the identity provider.
- `client_secret` - (Required) The client or client secret registered within the identity provider. This field is able to obtain its value from vault, use $${vault.ID} format.
- `token_url` - (Required) The Token URL.
- `extra_config` - (Optional) this block is needed to set extra configuration (Not yet supported variables or custom extensions)
    - `clientAuthMethod` (Optional) The client authentication method. Since Keycloak 8, this is a required attribute if OIDC provider is created over the Keycloak Userinterface.
    It accepts the values `client_secret_post` (Client secret sent as post), `client_secret_basic` (Client secret sent as basic auth), `client_secret_jwt` (Client secret as jwt) and `private_key_jwt ` (JTW signed with private key)
- `provider_id` - (Optional) The Provider id, defaults to `oidc`, unless you have a custom implementation.
- `backchannel_supported` - (Optional) Does the external IDP support backchannel logout ? Defaults to `true`.
- `validate_signature` - (Optional) Enable/disable signature validation of external IDP signatures. Defaults to `false`.
- `user_info_url` - (Optional) User Info URL.
- `jwks_url` - (Optional) JSON Web Key Set URL.
- `hide_on_login_page` - (Optional) Hide On Login Page. Defaults to `false`.
- `logout_url` - (Optional) The Logout URL is the end session endpoint to use to logout user from external identity provider.
- `login_hint` - (Optional) Pass login hint to identity provider.
- `ui_locales` - (Optional) Pass current locale to identity provider. Defaults to `false`.
- `accepts_prompt_none_forward_from_client` (Optional) Specifies whether the IDP accepts forwarded authentication requests that contain the prompt=none query parameter or not
- `default_scopes` - (Optional) The scopes to be sent when asking for authorization. It can be a space-separated list of scopes. Defaults to 'openid'.

### Import

Identity providers can be imported using the format `{{realm_id}}/{{idp_alias}}`, where `idp_alias` is the identity provider alias.

Example:

```bash
$ terraform import keycloak_oidc_identity_provider.realm_identity_provider my-realm/my-idp
```
