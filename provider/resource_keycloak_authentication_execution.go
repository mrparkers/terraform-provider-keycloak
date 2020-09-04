package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAuthenticationExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakAuthenticationExecutionCreate,
		Read:   resourceKeycloakAuthenticationExecutionRead,
		Delete: resourceKeycloakAuthenticationExecutionDelete,
		Update: resourceKeycloakAuthenticationExecutionUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakAuthenticationExecutionImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"authenticator": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"requirement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REQUIRED", "ALTERNATIVE", "OPTIONAL", "CONDITIONAL", "DISABLED"}, false), //OPTIONAL is removed from 8.0.0 onwards
				Default:      "DISABLED",
			},
		},
	}
}

func mapFromDataToAuthenticationExecution(data *schema.ResourceData) *keycloak.AuthenticationExecution {
	authenticationExecution := &keycloak.AuthenticationExecution{
		Id:              data.Id(),
		RealmId:         data.Get("realm_id").(string),
		ParentFlowAlias: data.Get("parent_flow_alias").(string),
		Authenticator:   data.Get("authenticator").(string),
		Requirement:     data.Get("requirement").(string),
	}

	return authenticationExecution
}

func mapFromAuthenticationExecutionToData(data *schema.ResourceData, authenticationExecution *keycloak.AuthenticationExecution) {
	data.SetId(authenticationExecution.Id)

	data.Set("realm_id", authenticationExecution.RealmId)
	data.Set("parent_flow_alias", authenticationExecution.ParentFlowAlias)
	data.Set("authenticator", authenticationExecution.Authenticator)
	data.Set("requirement", authenticationExecution.Requirement)
}

func mapFromAuthenticationExecutionInfoToData(data *schema.ResourceData, authenticationExecutionInfo *keycloak.AuthenticationExecutionInfo) {
	data.SetId(authenticationExecutionInfo.Id)

	data.Set("realm_id", authenticationExecutionInfo.RealmId)
	data.Set("parent_flow_alias", authenticationExecutionInfo.ParentFlowAlias)
}

func resourceKeycloakAuthenticationExecutionCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationExecution := mapFromDataToAuthenticationExecution(data)

	err := keycloakClient.NewAuthenticationExecution(authenticationExecution)
	if err != nil {
		return err
	}

	mapFromAuthenticationExecutionToData(data, authenticationExecution)

	return resourceKeycloakAuthenticationExecutionRead(data, meta)
}

func resourceKeycloakAuthenticationExecutionRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	parentFlowAlias := data.Get("parent_flow_alias").(string)
	id := data.Id()

	authenticationExecution, err := keycloakClient.GetAuthenticationExecution(realmId, parentFlowAlias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromAuthenticationExecutionToData(data, authenticationExecution)

	return nil
}

func resourceKeycloakAuthenticationExecutionUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationExecution := mapFromDataToAuthenticationExecution(data)

	err := keycloakClient.UpdateAuthenticationExecution(authenticationExecution)
	if err != nil {
		return err
	}

	mapFromAuthenticationExecutionToData(data, authenticationExecution)

	return nil
}

func resourceKeycloakAuthenticationExecutionDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteAuthenticationExecution(realmId, id)
}

func resourceKeycloakAuthenticationExecutionImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{parentFlowAlias}}/{{authenticationExecutionId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("parent_flow_alias", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
