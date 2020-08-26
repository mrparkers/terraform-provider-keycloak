package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAuthenticationExecutionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakAuthenticationExecutionConfigCreate,
		Read:   resourceKeycloakAuthenticationExecutionConfigRead,
		Delete: resourceKeycloakAuthenticationExecutionConfigDelete,
		Update: resourceKeycloakAuthenticationExecutionConfigUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakAuthenticationExecutionConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"execution_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func getAuthenticationExecutionConfigFromData(data *schema.ResourceData) *keycloak.AuthenticationExecutionConfig {
	config := make(map[string]string)
	for key, value := range data.Get("config").(map[string]interface{}) {
		config[key] = value.(string)
	}
	return &keycloak.AuthenticationExecutionConfig{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		ExecutionId: data.Get("execution_id").(string),
		Alias:       data.Get("alias").(string),
		Config:      config,
	}
}

func setAuthenticationExecutionConfigData(data *schema.ResourceData, config *keycloak.AuthenticationExecutionConfig) {
	data.SetId(config.Id)
	data.Set("realm_id", config.RealmId)
	data.Set("execution_id", config.ExecutionId)
	data.Set("alias", config.Alias)
	data.Set("config", config.Config)
}

func resourceKeycloakAuthenticationExecutionConfigCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	config := getAuthenticationExecutionConfigFromData(data)

	id, err := keycloakClient.NewAuthenticationExecutionConfig(config)
	if err != nil {
		return err
	}

	data.SetId(id)

	return resourceKeycloakAuthenticationExecutionConfigRead(data, meta)
}

func resourceKeycloakAuthenticationExecutionConfigRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	config := &keycloak.AuthenticationExecutionConfig{
		RealmId:     data.Get("realm_id").(string),
		ExecutionId: data.Get("execution_id").(string),
		Id:          data.Id(),
	}

	err := keycloakClient.GetAuthenticationExecutionConfig(config)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setAuthenticationExecutionConfigData(data, config)

	return nil
}

func resourceKeycloakAuthenticationExecutionConfigUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	config := getAuthenticationExecutionConfigFromData(data)

	err := keycloakClient.UpdateAuthenticationExecutionConfig(config)
	if err != nil {
		return err
	}

	return resourceKeycloakAuthenticationExecutionConfigRead(data, meta)
}

func resourceKeycloakAuthenticationExecutionConfigDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	config := &keycloak.AuthenticationExecutionConfig{
		RealmId: data.Get("realm_id").(string),
		Id:      data.Id(),
	}

	return keycloakClient.DeleteAuthenticationExecutionConfig(config)
}

func resourceKeycloakAuthenticationExecutionConfigImport(data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import. Supported import formats: {{realm}}/{{authenticationExecutionId}}/{{authenticationExecutionConfigId}}")
	}

	data.Set("realm_id", parts[0])
	data.Set("execution_id", parts[1])
	data.SetId(parts[2])

	return []*schema.ResourceData{data}, nil
}
