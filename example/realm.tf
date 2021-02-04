resource "keycloak_realm" "realm" {
  realm             = "dev"
  enabled           = true
  display_name      = "Development"
  display_name_html = "<b>Development</b>"

  login_theme = "keycloak"

  access_token_lifespan = "15m"

  ssl_required    = "external"
  password_policy = "length(8) and notUsername"

}
