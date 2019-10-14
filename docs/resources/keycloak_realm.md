# keycloak_realm

Allows for creating and managing Realms within Keycloak.

A realm manages a logical collection of users, credentials, roles, and groups.
Users log in to realms and can be federated from multiple sources.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm                = "test"
    enabled              = true
    display_name         = "test realm"

    login_theme          = "base"

    access_code_lifespan = "1h"

    smtp_server {
        host = "smtp.example.com"
        from = "example@example.com"

        auth {
        	username = "tom"
        	password = "password"
        }
    }

    internationalization {
        supported_locales = [
            "en",
            "de",
            "es"
        ]
        default_locale    = "en"
    }

    security_defenses {
        headers {
            x_frame_options                     = "DENY"
            content_security_policy             = "frame-src 'self'; frame-ancestors 'self'; object-src 'none';"
            content_security_policy_report_only = ""
            x_content_type_options              = "nosniff"
            x_robots_tag                        = "none"
            x_xss_protection                    = "1; mode=block"
            strict_transport_security           = "max-age=31536000; includeSubDomains"
        }
    }
}
```

### Argument Reference

The following arguments are supported:

- `realm` - (Required) The name of the realm. This is unique across Keycloak.
- `enabled` - (Optional) When false, users and clients will not be able to access this realm. Defaults to `true`.
- `display_name` - (Optional) The display name for the realm that is shown when logging in to the admin console.

##### Login Settings

The following attributes are all booleans, and can be found in the "Login" tab within the realm settings.
If any of these attributes are not specified, they will default to Keycloak's default settings.

- `registration_allowed` - (Optional) When true, user registration will be enabled, and a link for registration will be displayed on the login page.
- `registration_email_as_username` - (Optional) When true, the user's email will be used as their username during registration.
- `edit_username_allowed` - (Optional) When true, the username field is editable.
- `reset_password_allowed` - (Optional) When true, a "forgot password" link will be displayed on the login page.
- `remember_me` - (Optional) When true, a "remember me" checkbox will be displayed on the login page, and the user's session will not expire between browser restarts.
- `verify_email` - (Optional) When true, users are required to verify their email address after registration and after email address changes.
- `login_with_email_allowed` - (Optional) When true, users may log in with their email address.
- `duplicate_emails_allowed` - (Optional) When true, multiple users will be allowed to have the same email address. This attribute must be set to `false` if `login_with_email_allowed` is set to `true`.

##### Themes

The following attributes can be used to configure themes for the realm. Custom themes can be specified here.
If any of these attributes are not specified, they will default to Keycloak's default settings. Typically the `keycloak` theme is used by default.

- `login_theme` - (Optional) Used for the login, forgot password, and registration pages.
- `account_theme` - (Optional) Used for account management pages.
- `admin_theme` - (Optional) Used for the admin console.
- `email_theme` - (Optional) Used for emails that are sent by Keycloak.

##### Tokens

The following attributes can be found in the "Tokens" tab within the realm settings.

- `refresh_token_max_reuse` - (Optional) Maximum number of times a refresh token can be reused before they are revoked. If unspecified, refresh tokens will only be revoked when a different token is used.

The attributes below should be specified as [Go duration strings](https://golang.org/pkg/time/#Duration.String). They will default to Keycloak's default settings.

- `sso_session_idle_timeout` - (Optional) The amount of time a session can be idle before it expires.
- `sso_session_max_lifespan` - (Optional) The maximum amount of time before a session expires regardless of activity.
- `offline_session_idle_timeout` - (Optional) The amount of time an offline session can be idle before it expires.
- `offline_session_max_lifespan` - (Optional) The maximum amount of time before an offline session expires regardless of activity.
- `access_token_lifespan` - (Optional) The amount of time an access token can be used before it expires.
- `access_token_lifespan_for_implicit_flow` - (Optional) The amount of time an access token issued with the OpenID Connect Implicit Flow can be used before it expires.
- `access_code_lifespan` - (Optional) The maximum amount of time a client has to finish the authorization code flow.
- `access_code_lifespan_login` - (Optional) The maximum amount of time a user is permitted to stay on the login page before the authentication process must be restarted.
- `access_code_lifespan_user_action` - (Optional) The maximum amount of time a user has to complete login related actions, such as updating a password.
- `action_token_generated_by_user_lifespan` - (Optional) The maximum time a user has to use a user-generated permit before it expires.
- `action_token_generated_by_admin_lifespan` - (Optional) The maximum time a user has to use an admin-generated permit before it expires.

##### SMTP

The `smtp_server` block can be used to configure the realm's SMTP settings, which can be found in the "Email" tab in the GUI.
This block supports the following attributes:

- `host` - (Required) The host of the SMTP server.
- `port` - (Optional) The port of the SMTP server (defaults to 25).
- `from` - (Required) The email address for the sender.
- `from_display_name` - (Optional) The display name of the sender email address.
- `reply_to` - (Optional) The "reply to" email address.
- `reply_to_display_name` - (Optional) The display name of the "reply to" email address.
- `envelope_from` - (Optional) The email address uses for bounces.
- `starttls` - (Optional) When `true`, enables StartTLS. Defaults to `false`.
- `ssl` - (Optional) When `true`, enables SSL. Defaults to `false`.
- `auth` - (Optional) Enables authentication to the SMTP server.  This block supports the following attributes:
    - `username`- (Required) The SMTP server username.
    - `password` - (Required) The SMTP server password.

##### Internationalization

Internationalization support can be configured by using the `internationalization` block, which supports the following attributes:

- `supported_locales` - (Required) A list of [ISO 639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) locale codes that the realm should support.
- `default_locale` - (Required) The locale to use by default. This locale code must be present within the `supported_locales` list.

##### Security Defenses Headers

Header configuration support for browser security settings.

### Import

Realms can be imported using their name:

```bash
$ terraform import keycloak_realm.realm test
```
