# keycloak_realm

Allows for creating and managing Realms within Keycloak.

A realm manages a logical collection of users, credentials, roles, and groups.
Users log in to realms and can be federated from multiple sources.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm                   = "test"
    enabled                 = true
    display_name            = "test realm"
    display_name_html       = "<b>test realm</b>"

    login_theme          = "base"

    access_code_lifespan = "1h"
    ssl_required  = "external"
    password_policy = "upperCase(1) and length(8) and forceExpiredPasswordChange(365) and notUsername"
    attributes = {
      mycustomAttribute = "myCustomValue"
    }

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
        brute_force_detection {
            permanent_lockout                 = false
            max_login_failures                = 30
            wait_increment_seconds            = 60
            quick_login_check_milli_seconds   = 1000
            minimum_quick_login_wait_seconds  = 60
            max_failure_wait_seconds          = 900
            failure_reset_time_seconds        = 43200
        }
    }
}
```

### Argument Reference

The following arguments are supported:

- `realm` - (Required) The name of the realm. This is unique across Keycloak.
- `enabled` - (Optional) When false, users and clients will not be able to access this realm. Defaults to `true`.
- `display_name` - (Optional) The display name for the realm that is shown when logging in to the admin console.
- `display_name_html` - (Optional) The display name for the realm that is rendered as HTML on the screen when logging in to the admin console.
- `user_managed_access` - (Optional) When true, users are allowed to manage their own resources. Defaults to `false`.

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
- `ssl_required` - (Optional) Can be one of following values: 'none, 'external' or 'all'

##### Themes

The following attributes can be used to configure themes for the realm. Custom themes can be specified here.
If any of these attributes are not specified, they will default to Keycloak's default settings. Typically the `keycloak` theme is used by default.

- `login_theme` - (Optional) Used for the login, forgot password, and registration pages.
- `account_theme` - (Optional) Used for account management pages.
- `admin_theme` - (Optional) Used for the admin console.
- `email_theme` - (Optional) Used for emails that are sent by Keycloak.

##### Tokens

The following attributes can be found in the "Tokens" tab within the realm settings.

- `default_signature_algorithm` - (Optional) Default algorithm used to sign tokens for the realm.
- `revoke_refresh_token` - (Optional) If enabled a refresh token can only be used number of times specified in 'refresh_token_max_reuse' before they are revoked. If unspecified, refresh tokens can be reused.
- `refresh_token_max_reuse` - (Optional) Maximum number of times a refresh token can be reused before they are revoked. If unspecified and 'revoke_refresh_token' is enabled the default value is 0 and refresh tokens can not be reused.

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

Header configuration support for browser security settings and brute force detection. It can be configured trough the`security_defenses` block using the `headers` and the `brute_force_detection` subblocks.

The `headers` block supports the following attributes:

- `x_frame_options` - (Optional) Sets the x-frame-option, could be used for prevent pages from being included by non-origin iframes. More Infos could be found in the [RFC7034](https://tools.ietf.org/html/rfc7034)
- `content_security_policy` - (Optional) Sets the Content Security Policy, could be used for prevent pages from being included by non-origin iframes. More Infos could be found in the [W3C-CSP](https://www.w3.org/TR/CSP/) Abstract.
- `content_security_policy_report_only` - (Optional) For testing Content Security Policies.
- `x_content_type_options` - (Optional) Sets the X-Content-Type-Options,  could be used for prevent MIME-sniffing a respone away from the declared content-type
- `x_robots_tag` - (Optional) Prevent pages from appearing in search engines.
- `x_xss_protection` - (Optional) This header configures the Cross-site scripting (XSS) filter in your browser.
- `strict_transport_security` - (Optional) The Script-Transport-Security HTTP header tells browsers to always use HTTPS. Maxage could be set and subdomains could be also included.

The `brute_force_detection` block supports the following attributes:

- `permanent_lockout` - (Optional) Lock the user permanently when the user exceeds the maximum login failures.
- `max_login_failures` - (Optional) How many failures before wait is triggered.
- `wait_increment_seconds` - (Optional) When failure threshold has been met, how much time should the user be locked out?
- `quick_login_check_milli_seconds` - (Optional) If a failure happens concurrently too quickly, lock out the user.
- `minimum_quick_login_wait_seconds` - (Optional) How long to wait after a quick login failure.
- `max_failure_wait_seconds ` - (Optional) Max. time a user will be locked out.
- `failure_reset_time_seconds` - (Optional) When will failure count be reset?

#### Atributes

Map, can be used to add custom attributes to a realm. Or perhaps influence a certain attribute that is not supported in this terraform-provider

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `internal_id` - When importing realms created outside of this terraform provider, they could use generated arbitrary IDs for the technical realm id. Realms created by this provider always use the realm's name for its technical id.

### Import

Realms can be imported using their name:

```bash
$ terraform import keycloak_realm.realm test
```
