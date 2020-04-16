resource "keycloak_realm" "source_realm" {
  realm   = "source_realm"
  enabled = true
}

resource "keycloak_openid_client" "destination_client" {
  realm_id                 = "${keycloak_realm.source_realm.id}"
  name                     = "destination_client"
  client_id                = "destination_client"
  client_secret            = "secret"
  description              = "a client used by the destination realm"
  access_type              = "CONFIDENTIAL"
  standard_flow_enabled    = true
  valid_redirect_uris = [
    "http://localhost:8080/*",
  ]
}

//do not get confused this just to have multiple federate idps on the destination realm
resource "keycloak_openid_client" "destination_double_client" {
  realm_id                 = "${keycloak_realm.source_realm.id}"
  name                     = "destination_double_client"
  client_id                = "destination_double_client"
  client_secret            = "secret2"
  description              = "a second client used by the destination realm"
  access_type              = "CONFIDENTIAL"
  standard_flow_enabled    = true
  valid_redirect_uris = [
    "http://localhost:8080/*",
  ]
}

resource "keycloak_user" "source_user" {
  realm_id   = "${keycloak_realm.source_realm.id}"
  username   = "source"
  email      = "source@fakedomain.com"
  first_name = "source"
  last_name  = "source"
  initial_password {
    value     = "source"
    temporary = false
  }
}

resource "keycloak_realm" "destination_realm" {
  realm   = "destination_realm"
  enabled = true
}

resource keycloak_oidc_identity_provider source_oidc_idp {
  realm              = "${keycloak_realm.destination_realm.id}"
  alias              = "source"
  authorization_url  = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/auth"
  token_url          = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/token"
  user_info_url      = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/userinfo"
  jwks_url           = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/certs"
  validate_signature = true
  client_id          = "${keycloak_openid_client.destination_client.client_id}"
  client_secret      = "${keycloak_openid_client.destination_client.client_secret}"
  default_scopes     = "openid"
}

//do not get confused this second idp towards source_realm, this could a completly different idp
resource keycloak_oidc_identity_provider second_source_oidc_idp {
  realm              = "${keycloak_realm.destination_realm.id}"
  alias              = "source2"
  authorization_url  = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/auth"
  token_url          = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/token"
  user_info_url      = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/userinfo"
  jwks_url           = "http://localhost:8080/auth/realms/${keycloak_realm.source_realm.id}/protocol/openid-connect/certs"
  validate_signature = true
  client_id          = "${keycloak_openid_client.destination_double_client.client_id}"
  client_secret      = "${keycloak_openid_client.destination_double_client.client_secret}"
  default_scopes     = "openid"
}

resource "keycloak_user" "destination_user" {
  realm_id   = "${keycloak_realm.destination_realm.id}"
  username   = "my_destination_username"
  email      = "source@otherdomain.be"
  first_name = "Destination_source"
  last_name  = "Destination_source"
  //federated link through source idp
  federated_identity {
    identity_provider = "${keycloak_oidc_identity_provider.source_oidc_idp.alias}"
    user_id           = "${keycloak_user.source_user.id}"
    user_name         = "${keycloak_user.source_user.username}"
  }
  //federated link through second source idp
  federated_identity {
    identity_provider = "${keycloak_oidc_identity_provider.second_source_oidc_idp.alias}"
    user_id           = "${keycloak_user.source_user.id}"
    user_name         = "${keycloak_user.source_user.username}"
  }
}
