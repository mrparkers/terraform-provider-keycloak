# keycloak_saml_client_installation_provider data source

This data source can be used to retrieve Installation Provider
of a SAML Client.

### Example Usage

In the example below, we extract the SAML metadata IDPSSODescriptor 
to pass it to the AWS IAM SAML Provider.

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_saml_client" "saml_client" {
    realm_id                = "${keycloak_realm.realm.id}"
    client_id               = "test-saml-client"
    name                    = "test-saml-client"

    sign_documents          = false
    sign_assertions         = true
    include_authn_statement = true

    signing_certificate = "${file("saml-cert.pem")}"
    signing_private_key = "${file("saml-key.pem")}"
}

data "keycloak_saml_client_installation_provider" "saml_idp_descriptor" {
  realm_id    = "${keycloak_realm.realm.id}"
  client_id   = "${keycloak_saml_client.saml_client}"
  provider_id = "saml-idp-descriptor"
}


resource "aws_iam_saml_provider" "default" {
  name                   = "myprovider"
  saml_metadata_document = data.keycloak_saml_client_installation_provider.saml_idp_descriptor.value
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this group exists within.
- `client_id` - (Required) The name of the saml client. Not the id of the client.
- `provider_id` - (Required) Could be one of `saml-idp-descriptor`, `keycloak-saml`, `saml-sp-descriptor`, `keycloak-saml-subsystem`, `mod-auth-mellon`

### Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

- `id` - The hash of the value
- `value` The returned XML document needed for SAML installation
