# keycloak_authentication_execution

Allows for managing an authentication execution.

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
    requirement       = "REQUIRED"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm the authentication execution exists in.
- `parent_flow_alias` - (Required) The flow this execution is attached to.
- `authenticator` - (Required) The name of the authenticator.
- `requirement`- (Optional) The requirement setting, which can be one of the following:
	- `REQUIRED`
	- `ALTERNATIVE`
	- `DISABLED`

### Import

Executions can be imported using the formats: `{{realmId}}/{{parentFlowAlias}}/{{authenticationExecutionId}}`.

Example:

```bash
$ terraform import keycloak_authentication_execution my-realm/my-flow/30559fcf-6fb8-45ea-8c46-2b86f46ebc17
```
