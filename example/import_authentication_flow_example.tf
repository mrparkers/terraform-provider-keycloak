resource "keycloak_realm" "default-config-test-realm" {
  realm = "default-config-test-realm"
  enabled = true
}

resource "keycloak_authentication_flow" "imported-flow" {
  realm_id = keycloak_realm.default-config-test-realm.id
  alias    = "browser"
  import = true
  # changed attributes
  description = "new description"
  provider_id="client-flow"
}

resource "keycloak_authentication_subflow" "imported-subflow" {
  realm_id = keycloak_realm.default-config-test-realm.id
  parent_flow_alias = keycloak_authentication_flow.imported-flow.alias
  alias    = "forms"
  import = true
  # changed attributes
  description = "new description" # default: Username, password, otp and other auth forms
}

resource "keycloak_authentication_execution" "imported-execution" {
  realm_id = keycloak_realm.default-config-test-realm.id
  parent_flow_alias = keycloak_authentication_flow.imported-flow.alias
  authenticator = "identity-provider-redirector"
  import = true
  # changed attributes
  requirement = "REQUIRED" # default: ALTERNATIVE
}
