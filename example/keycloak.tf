terraform {
  required_providers {
    keycloak = {
      source  = "terraform.local/mrparkers/keycloak"
      version = ">= 2.0"
    }
  }
}

provider "keycloak" {
  client_id = "admin-cli"
  username  = "admin"
  password  = "admin1234."
  url       = "http://localhost:8080"
}
