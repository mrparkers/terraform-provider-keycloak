---
page_title: "keycloak_user_template_importer_identity_provider_mapper Resource"
---

# keycloak\_user\_template\_importer\_identity\_provider\_mapper Resource

Allows for creating and managing an username template importer identity provider mapper within Keycloak.

The username template importer mapper can be used to map externally defined OIDC claims or SAML attributes with a template to the username of the imported Keycloak user:

- Substitutions are enclosed in \${}. For example: '\${ALIAS}.\${CLAIM.sub}'. ALIAS is the provider alias. CLAIM.\<NAME\> references an ID or Access token claim.

~> If you are using Keycloak 10 or higher, you will need to specify the `extra_config` argument in order to define a `syncMode` for the mapper.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_oidc_identity_provider" "oidc" {
  realm             = keycloak_realm.realm.id
  alias             = "oidc"
  authorization_url = "https://example.com/auth"
  token_url         = "https://example.com/token"
  client_id         = "example_id"
  client_secret     = "example_token"
  default_scopes    = "openid random profile"
}

resource "keycloak_user_template_importer_identity_provider_mapper" "username_importer" {
  realm                   = keycloak_realm.realm.id
  name                    = "username-template-importer"
  identity_provider_alias = keycloak_oidc_identity_provider.oidc.alias
  template                = "$${ALIAS}.$${CLAIM.email}"

  # extra_config with syncMode is required in Keycloak 10+
  extra_config = {
    syncMode = "INHERIT"
  }
}
```

## Argument Reference

The following arguments are supported:

- `realm` - (Required) The name of the realm.
- `name` - (Required) The name of the mapper.
- `identity_provider_alias` - (Required) The alias of the associated identity provider.
- `template` - (Required) Template to use to format the username to import. Substitutions are enclosed in \${}. For example: '\$\${ALIAS}.\$\${CLAIM.sub}'. ALIAS is the provider alias. CLAIM.\<NAME\> references an ID or Access token claim.
- `extra_config` - (Optional) Key/value attributes to add to the identity provider mapper model that is persisted to Keycloak. This can be used to extend the base model with new Keycloak features.

## Import

Identity provider mappers can be imported using the format `{{realm_id}}/{{idp_alias}}/{{idp_mapper_id}}`, where `idp_alias` is the identity provider alias, and `idp_mapper_id` is the unique ID that Keycloak
assigns to the mapper upon creation. This value can be found in the URI when editing this mapper in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_user_template_importer_identity_provider_mapper.username_importer my-realm/my-mapper/f446db98-7133-4e30-b18a-3d28fde7ca1b
```
