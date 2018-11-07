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

	err := keycloakClient.AttachOpenidClientDefaultScopes(realmId, clientId, defaultScopes.List())
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s", realmId, clientId))

	return resourceKeycloakOpenidClientDefaultScopesRead(data, meta)
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
	data.SetId(fmt.Sprintf("%s/%s", realmId, clientId))

	return nil
}

func resourceKeycloakOpenidClientDefaultScopesDelete(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKeycloakOpenidClientDefaultScopesUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}
