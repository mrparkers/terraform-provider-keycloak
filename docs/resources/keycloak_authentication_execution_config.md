# keycloak_authentication_execution_config

Allows for managing an authentication execution configuration.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
	realm   = "my-realm"
	enabled = true
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "my-flow-alias"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	authenticator     = "identity-provider-redirector"
}

resource "keycloak_authentication_execution_config" "config" {
	realm_id     = "${keycloak_realm.realm.id}"
	execution_id = "${keycloak_authentication_execution.execution.id}"
	alias        = "my-config-alias"
	config = {
		defaultProvider = "my-config-default-idp"
	}
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm the authentication execution exists in.
- `execution_id` - (Required) The authentication execution this configuration is attached to.
- `alias` - (Required) The name of the configuration.
- `config` - (Optional) The configuration. Keys are specific to each configurable authentication execution and not checked when applying.

### Import

Configurations can be imported using the format `{{realm}}/{{authenticationExecutionId}}/{{authenticationExecutionConfigId}}`.
If the `authenticationExecutionId` is incorrect, the import will still be successful.
A subsequent apply will change the `authenticationExecutionId` to the correct one, which causes the configuration to be replaced.

Example:

```bash
$ terraform import keycloak_authentication_execution_config.config my-realm/be081463-ddbf-4b42-9eff-9c97886f24ff/30559fcf-6fb8-45ea-8c46-2b86f46ebc17
```
