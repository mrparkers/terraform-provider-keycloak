package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakSamlClientDefaultScopes() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakSamlClientDefaultScopesCreate,
		Read:   resourceKeycloakSamlClientDefaultScopesRead,
		Delete: resourceKeycloakSamlClientDefaultScopesDelete,
		Update: resourceKeycloakSamlClientDefaultScopesUpdate,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceKeycloakSamlClientDefaultScopesCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	defaultScopes := data.Get("default_scopes").(*schema.Set)

	err := keycloakClient.AttachSamlClientDefaultScopes(realmId, clientId, interfaceSliceToStringSlice(defaultScopes.List()))
	if err != nil {
		return err
	}

	data.SetId(samlClientDefaultScopesId(realmId, clientId))

	return resourceKeycloakSamlClientDefaultScopesRead(data, meta)
}

func samlClientDefaultScopesId(realmId string, clientId string) string {
	return fmt.Sprintf("%s/%s", realmId, clientId)
}

func resourceKeycloakSamlClientDefaultScopesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	clientScopes, err := keycloakClient.GetSamlClientDefaultScopes(realmId, clientId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	var defaultScopes []string
	for _, clientScope := range clientScopes {
		defaultScopes = append(defaultScopes, clientScope.Name)
	}

	data.Set("default_scopes", defaultScopes)
	data.SetId(samlClientDefaultScopesId(realmId, clientId))

	return nil
}

func resourceKeycloakSamlClientDefaultScopesUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	tfSamlClientDefaultScopes := data.Get("default_scopes").(*schema.Set)

	keycloakSamlClientDefaultScopes, err := keycloakClient.GetSamlClientDefaultScopes(realmId, clientId)
	if err != nil {
		return err
	}

	var samlClientDefaultScopesToDetach []string
	for _, keycloakSamlClientDefaultScope := range keycloakSamlClientDefaultScopes {
		// if this scope is attached in keycloak and tf state, no update is required
		// remove it from the set so we can look at scopes that need to be attached later
		if tfSamlClientDefaultScopes.Contains(keycloakSamlClientDefaultScope.Name) {
			tfSamlClientDefaultScopes.Remove(keycloakSamlClientDefaultScope.Name)
		} else {
			// if this scope is attached in keycloak but not in tf state, add them to a slice containing all scopes to detach
			samlClientDefaultScopesToDetach = append(samlClientDefaultScopesToDetach, keycloakSamlClientDefaultScope.Name)
		}
	}

	// detach scopes that aren't in tf state
	err = keycloakClient.DetachSamlClientDefaultScopes(realmId, clientId, samlClientDefaultScopesToDetach)
	if err != nil {
		return err
	}

	// attach scopes that exist in tf state but not in keycloak
	err = keycloakClient.AttachSamlClientDefaultScopes(realmId, clientId, interfaceSliceToStringSlice(tfSamlClientDefaultScopes.List()))
	if err != nil {
		return err
	}

	data.SetId(samlClientDefaultScopesId(realmId, clientId))

	return resourceKeycloakSamlClientDefaultScopesRead(data, meta)
}

func resourceKeycloakSamlClientDefaultScopesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	defaultScopes := data.Get("default_scopes").(*schema.Set)

	return keycloakClient.DetachSamlClientDefaultScopes(realmId, clientId, interfaceSliceToStringSlice(defaultScopes.List()))
}
