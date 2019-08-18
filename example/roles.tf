resource "keycloak_realm" "roles_example" {
	realm   = "roles-example"
	enabled = true
}

// API Client and roles

resource "keycloak_openid_client" "pet_api" {
	realm_id  = "${keycloak_realm.roles_example.id}"
	client_id = "pet-api"
	name      = "pet-api"

	enabled = true

	access_type = "BEARER-ONLY"
}

resource "keycloak_role" "pet_api_create_pet" {
	name        = "create-pet"
	realm_id    = "${keycloak_realm.roles_example.id}"
	client_id   = "${keycloak_openid_client.pet_api.id}"
	description = "Ability to create a new pet"
}

resource "keycloak_role" "pet_api_update_pet" {
	name        = "update-pet"
	realm_id    = "${keycloak_realm.roles_example.id}"
	client_id   = "${keycloak_openid_client.pet_api.id}"
	description = "Ability to update a pet"
}

resource "keycloak_role" "pet_api_read_pet" {
	name        = "read-pet"
	realm_id    = "${keycloak_realm.roles_example.id}"
	client_id   = "${keycloak_openid_client.pet_api.id}"
	description = "Ability to read / list pets"
}

resource "keycloak_role" "pet_api_delete_pet" {
	name        = "delete-pet"
	realm_id    = "${keycloak_realm.roles_example.id}"
	client_id   = "${keycloak_openid_client.pet_api.id}"
	description = "Ability to delete a pet"
}

// Consumer client

resource "keycloak_openid_client" "pet_app" {
	realm_id  = "${keycloak_realm.roles_example.id}"
	client_id = "pet-app"
	name      = "pet-app"

	enabled = true

	access_type   = "CONFIDENTIAL"
	client_secret = "pet-app-secret"

	// authenticated users - could have many roles
	standard_flow_enabled    = true
	// unauthenticated users - needs at least read / list role for browsing
	service_accounts_enabled = true

	valid_redirect_uris = [
		"http://localhost:5555/openid-callback"
	]
}

// The app will always need access to the API, so this audience should be used regardless of auth type
resource "keycloak_openid_audience_protocol_mapper" "pet_app_pet_api_audience_mapper" {
	realm_id  = "${keycloak_realm.roles_example.id}"
	client_id = "${keycloak_openid_client.pet_app.id}"
	name      = "audience-mapper"

	included_client_audience = "${keycloak_openid_client.pet_api.client_id}"
}

// Users and groups

resource "keycloak_group" "pet_api_base" {
	realm_id = "${keycloak_realm.roles_example.id}"
	name     = "pets"
}

resource "keycloak_group" "pet_api_admins" {
	realm_id  = "${keycloak_realm.roles_example.id}"
	parent_id = "${keycloak_group.pet_api_base.id}"
	name      = "admins"
}

resource "keycloak_group" "pet_api_front_desk" {
	realm_id  = "${keycloak_realm.roles_example.id}"
	parent_id = "${keycloak_group.pet_api_base.id}"
	name      = "front-desk"
}
