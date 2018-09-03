provider "keycloak" {
  client_id     = "terraform"
  client_secret = "884e0f95-0f42-4a63-9b1f-94274655669e"
  url           = "http://localhost:8080"
}

resource "keycloak_realm" "test" {
  realm                          = "test"
  enabled                        = false
  display_name                   = "foo"

  registration_allowed           = true
  registration_email_as_username = false
  edit_username_allowed          = false
  reset_password_allowed         = true
  remember_me                    = false
  verify_email                   = true
  login_with_email_allowed       = true
  duplicate_emails_allowed       = false
}
