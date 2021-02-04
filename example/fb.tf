variable "facebook_app_id" {
  type = string
  default = "1234"
}

variable "facebook_app_secret" {
  type = string
  default = "1234"
}

resource "keycloak_oidc_identity_provider" "facebook_identity_provider" {
  realm                         = keycloak_realm.realm.id
  alias                         = "facebook"
  provider_id                   = "facebook"
  client_id                     = var.facebook_app_id
  client_secret                 = var.facebook_app_secret
  trust_email                   = true
//  first_broker_login_flow_alias = keycloak_authentication_flow.first-broker-login-auto.alias
  token_url                     = ""
  authorization_url             = ""

  extra_config = {
    "syncMode" = "IMPORT"
  }
}


resource keycloak_user_template_importer_identity_provider_mapper oidc {
  realm                   = keycloak_realm.realm.id
  name                    = "alias"
  identity_provider_alias = keycloak_oidc_identity_provider.facebook_identity_provider.alias
  identity_provider_mapper="oidc-username-idp-mapper"
  template                = "$${UUID}"

  #KC10 support
  extra_config = {
    syncMode = "INHERIT"
    target = "LOCAL"
  }
}
