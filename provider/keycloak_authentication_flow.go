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
				Type: schema.TypeSet,
				// this hash function basically makes this behave like a TypeList
				// the only exception is that indices can be skipped in tf config and defined manually in Keycloak
				Set: func(v interface{}) int {
					return v.(map[string]interface{})["index"].(int)
				},
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

// given the `execution` set from data, return a list of executions, sorted by index
func mapToAuthenticationExecutionList(v interface{}, realmId, parentFlowAlias string) keycloak.AuthenticationExecutionList {
	var executions keycloak.AuthenticationExecutionList
	for _, ex := range v.(*schema.Set).List() {
		execution := ex.(map[string]interface{})

		executions = append(executions, &keycloak.AuthenticationExecution{
			Id:              execution["execution_id"].(string),
			RealmId:         realmId,
			ParentFlowAlias: parentFlowAlias,
			Provider:        execution["provider"].(string),
			Requirement:     execution["requirement"].(string),
			Index:           execution["index"].(int),
		})
	}

	sort.Sort(executions)

	return executions
}

func flattenAuthenticationExecutions(executions []*keycloak.AuthenticationExecution) []map[string]interface{} {
	state := make([]map[string]interface{}, 0, len(executions))

	for _, execution := range executions {
		data := make(map[string]interface{})

		data["execution_id"] = execution.Id
		data["provider"] = execution.Provider
		data["requirement"] = execution.Requirement
		data["index"] = execution.Index

		state = append(state, data)
	}

	return state
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
		executions := mapToAuthenticationExecutionList(v, authenticationFlow.RealmId, authenticationFlow.Alias)

		for _, execution := range executions {
			err = keycloakClient.NewAuthenticationExecution(execution)
			if err != nil {
				return err
			}

			err = keycloakClient.UpdateAuthenticationExecution(execution)
			if err != nil {
				return err
			}
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

	executions, err := keycloakClient.ListAuthenticationExecutions(authenticationFlow.RealmId, authenticationFlow.Alias)
	if err != nil {
		return err
	}

	data.Set("execution", flattenAuthenticationExecutions(executions))

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
