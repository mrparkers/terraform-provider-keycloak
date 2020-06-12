# keycloak_saml_client_scope

Allows for creating and managing Keycloak client scopes that can be attached to
clients that use the SAML protocol.

Client Scopes can be used to share common protocol and role mappings between multiple
clients within a realm.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_saml_client_scope" "saml_client_scope" {
    realm_id    = "${keycloak_realm.realm.id}"
    name        = "groups"
    description = "This scope will map a user's group memberships to SAML assertion"
    gui_order   = 1
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this client scope belongs to.
- `name` - (Required) The display name of this client scope in the GUI.
- `description` - (Optional) The description of this client scope in the GUI.
- `consent_screen_text` - (Optional) When set, a consent screen will be displayed to users
authenticating to clients with this scope attached. The consent screen will display the string
value of this attribute.
- `gui_order` - (Optional) Specify order of the client scope in GUI (such as in Consent page) as integer.

### Import

Client scopes can be imported using the format `{{realm_id}}/{{client_scope_id}}`, where `client_scope_id` is the unique ID that Keycloak
assigns to the client scope upon creation. This value can be found in the URI when editing this client scope in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_saml_client_scope.saml_client_scope my-realm/e8a5d115-6985-4de3-a0f5-732e1be4525e
```
