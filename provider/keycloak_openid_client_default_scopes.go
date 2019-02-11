package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientDefaultScopes() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientDefaultScopesCreate,
		Read:   resourceKeycloakOpenidClientDefaultScopesRead,
		Delete: resourceKeycloakOpenidClientDefaultScopesDelete,
		Update: resourceKeycloakOpenidClientDefaultScopesUpdate,
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

func resourceKeycloakOpenidClientDefaultScopesCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	defaultScopes := data.Get("default_scopes").(*schema.Set)

	err := keycloakClient.AttachOpenidClientDefaultScopes(realmId, clientId, interfaceSliceToStringSlice(defaultScopes.List()))
	if err != nil {
		return err
	}

	data.SetId(openidClientDefaultScopesId(realmId, clientId))

	return resourceKeycloakOpenidClientDefaultScopesRead(data, meta)
}

func openidClientDefaultScopesId(realmId string, clientId string) string {
	return fmt.Sprintf("%s/%s", realmId, clientId)
}

func resourceKeycloakOpenidClientDefaultScopesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	clientScopes, err := keycloakClient.GetOpenidClientDefaultScopes(realmId, clientId)
	if err != nil {
		return err
	}

	var defaultScopes []string
	for _, clientScope := range clientScopes {
		defaultScopes = append(defaultScopes, clientScope.Name)
	}

	data.Set("default_scopes", defaultScopes)
	data.SetId(openidClientDefaultScopesId(realmId, clientId))

	return nil
}

func resourceKeycloakOpenidClientDefaultScopesUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	tfOpenidClientDefaultScopes := data.Get("default_scopes").(*schema.Set)

	keycloakOpenidClientDefaultScopes, err := keycloakClient.GetOpenidClientDefaultScopes(realmId, clientId)
	if err != nil {
		return err
	}

	var openidClientDefaultScopesToDetach []string
	for _, keycloakOpenidClientDefaultScope := range keycloakOpenidClientDefaultScopes {
		// if this scope is attached in keycloak and tf state, no update is required
		// remove it from the set so we can look at scopes that need to be attached later
		if tfOpenidClientDefaultScopes.Contains(keycloakOpenidClientDefaultScope.Name) {
			tfOpenidClientDefaultScopes.Remove(keycloakOpenidClientDefaultScope.Name)
		} else {
			// if this scope is attached in keycloak but not in tf state, add them to a slice containing all scopes to detach
			openidClientDefaultScopesToDetach = append(openidClientDefaultScopesToDetach, keycloakOpenidClientDefaultScope.Name)
		}
	}

	// detach scopes that aren't in tf state
	err = keycloakClient.DetachOpenidClientDefaultScopes(realmId, clientId, openidClientDefaultScopesToDetach)
	if err != nil {
		return err
	}

	// attach scopes that exist in tf state but not in keycloak
	err = keycloakClient.AttachOpenidClientDefaultScopes(realmId, clientId, interfaceSliceToStringSlice(tfOpenidClientDefaultScopes.List()))
	if err != nil {
		return err
	}

	data.SetId(openidClientDefaultScopesId(realmId, clientId))

	return resourceKeycloakOpenidClientDefaultScopesRead(data, meta)
}

func resourceKeycloakOpenidClientDefaultScopesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	defaultScopes := data.Get("default_scopes").(*schema.Set)

	return keycloakClient.DetachOpenidClientDefaultScopes(realmId, clientId, interfaceSliceToStringSlice(defaultScopes.List()))
}
