package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"sort"
	"strings"
)

func resourceKeycloakAuthenticationExecution() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"requirement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REQUIRED", "ALTERNATIVE", "OPTIONAL", "DISABLED"}, false),
				Default:      "DISABLED",
			},
			"index": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"execution_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

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
			"execution": {
				Type:     schema.TypeSet,
				Set:      schema.HashResource(resourceKeycloakAuthenticationExecution()),
				Optional: true,
				Elem:     resourceKeycloakAuthenticationExecution(),
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

func mapToAuthenticationExecutionList(v interface{}) keycloak.AuthenticationExecutionList {
	var executions keycloak.AuthenticationExecutionList
	for _, ex := range v.(*schema.Set).List() {
		exMap := ex.(map[string]interface{})

		executions = append(executions, &keycloak.AuthenticationExecution{
			Id:          exMap["execution_id"].(string),
			Provider:    exMap["provider"].(string),
			Requirement: exMap["requirement"].(string),
			Index:       exMap["index"].(int),
		})
	}

	return executions
}

func resourceKeycloakAuthenticationFlowCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationFlow := mapFromDataToAuthenticationFlow(data)

	err := keycloakClient.NewAuthenticationFlow(authenticationFlow)
	if err != nil {
		return err
	}

	mapFromAuthenticationFlowToData(data, authenticationFlow)

	if v, ok := data.GetOk("execution"); ok {
		executions := mapToAuthenticationExecutionList(v)
		sort.Sort(executions)

		for _, execution := range executions {
			newExecution, err := keycloakClient.NewAuthenticationExecution(authenticationFlow.RealmId, authenticationFlow.Alias, execution.Provider)
			if err != nil {
				return err
			}

			execution = newExecution
		}
	}

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
