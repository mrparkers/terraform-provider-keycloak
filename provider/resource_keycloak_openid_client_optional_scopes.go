package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientOptionalScopes() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientOptionalScopesCreate,
		Read:   resourceKeycloakOpenidClientOptionalScopesRead,
		Delete: resourceKeycloakOpenidClientOptionalScopesDelete,
		Update: resourceKeycloakOpenidClientOptionalScopesUpdate,
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
			"optional_scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceKeycloakOpenidClientOptionalScopesCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	optionalScopes := data.Get("optional_scopes").(*schema.Set)

	err := keycloakClient.AttachOpenidClientOptionalScopes(realmId, clientId, interfaceSliceToStringSlice(optionalScopes.List()))
	if err != nil {
		return err
	}

	data.SetId(openidClientOptionalScopesId(realmId, clientId))

	return resourceKeycloakOpenidClientOptionalScopesRead(data, meta)
}

func openidClientOptionalScopesId(realmId string, clientId string) string {
	return fmt.Sprintf("%s/%s", realmId, clientId)
}

func resourceKeycloakOpenidClientOptionalScopesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)

	clientScopes, err := keycloakClient.GetOpenidClientOptionalScopes(realmId, clientId)
	if err != nil {
		return err
	}

	var optionalScopes []string
	for _, clientScope := range clientScopes {
		optionalScopes = append(optionalScopes, clientScope.Name)
	}

	data.Set("optional_scopes", optionalScopes)
	data.SetId(openidClientOptionalScopesId(realmId, clientId))

	return nil
}

func resourceKeycloakOpenidClientOptionalScopesUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	tfOpenidClientOptionalScopes := data.Get("optional_scopes").(*schema.Set)

	keycloakOpenidClientOptionalScopes, err := keycloakClient.GetOpenidClientOptionalScopes(realmId, clientId)
	if err != nil {
		return err
	}

	var openidClientOptionalScopesToDetach []string
	for _, keycloakOpenidClientOptionalScope := range keycloakOpenidClientOptionalScopes {
		// if this scope is attached in keycloak and tf state, no update is required
		// remove it from the set so we can look at scopes that need to be attached later
		if tfOpenidClientOptionalScopes.Contains(keycloakOpenidClientOptionalScope.Name) {
			tfOpenidClientOptionalScopes.Remove(keycloakOpenidClientOptionalScope.Name)
		} else {
			// if this scope is attached in keycloak but not in tf state, add them to a slice containing all scopes to detach
			openidClientOptionalScopesToDetach = append(openidClientOptionalScopesToDetach, keycloakOpenidClientOptionalScope.Name)
		}
	}

	// detach scopes that aren't in tf state
	err = keycloakClient.DetachOpenidClientOptionalScopes(realmId, clientId, openidClientOptionalScopesToDetach)
	if err != nil {
		return err
	}

	// attach scopes that exist in tf state but not in keycloak
	err = keycloakClient.AttachOpenidClientOptionalScopes(realmId, clientId, interfaceSliceToStringSlice(tfOpenidClientOptionalScopes.List()))
	if err != nil {
		return err
	}

	data.SetId(openidClientOptionalScopesId(realmId, clientId))

	return resourceKeycloakOpenidClientOptionalScopesRead(data, meta)
}

func resourceKeycloakOpenidClientOptionalScopesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	optionalScopes := data.Get("optional_scopes").(*schema.Set)

	return keycloakClient.DetachOpenidClientOptionalScopes(realmId, clientId, interfaceSliceToStringSlice(optionalScopes.List()))
}
