package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakAuthenticationFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakAuthenticationFlowCreate,
		Read:   resourceKeycloakAuthenticationFlowRead,
		Delete: resourceKeycloakAuthenticationFlowDelete,
		Update: resourceKeycloakAuthenticationFlowUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakAuthenticationFlowImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func mapFromDataToAuthenticationFlow(data *schema.ResourceData) *keycloak.AuthenticationFlow {
	authenticationFlow := &keycloak.AuthenticationFlow{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		Alias:       data.Get("alias").(string),
		Description: data.Get("description").(string),
	}

	return authenticationFlow
}

func mapFromAuthenticationFlowToData(data *schema.ResourceData, authenticationFlow *keycloak.AuthenticationFlow) {
	data.SetId(authenticationFlow.Id)

	data.Set("realm_id", authenticationFlow.RealmId)
	data.Set("alias", authenticationFlow.Alias)
	data.Set("description", authenticationFlow.Description)
}

func resourceKeycloakAuthenticationFlowCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationFlow := mapFromDataToAuthenticationFlow(data)

	err := keycloakClient.NewAuthenticationFlow(authenticationFlow)
	if err != nil {
		return err
	}

	mapFromAuthenticationFlowToData(data, authenticationFlow)

	return resourceKeycloakAuthenticationFlowRead(data, meta)
}

func resourceKeycloakAuthenticationFlowRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	authenticationFlow, err := keycloakClient.GetAuthenticationFlow(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromAuthenticationFlowToData(data, authenticationFlow)

	return nil
}

func resourceKeycloakAuthenticationFlowUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationFlow := mapFromDataToAuthenticationFlow(data)

	err := keycloakClient.UpdateAuthenticationFlow(authenticationFlow)
	if err != nil {
		return err
	}

	mapFromAuthenticationFlowToData(data, authenticationFlow)

	return nil
}

func resourceKeycloakAuthenticationFlowDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteAuthenticationFlow(realmId, id)
}

func resourceKeycloakAuthenticationFlowImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{authenticationFlowId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
