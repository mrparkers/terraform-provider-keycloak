package provider

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

var (
	keycloakOpenidClientAccessTypes = []string{"CONFIDENTIAL", "PUBLIC", "BEARER-ONLY"}
)

func resourceKeycloakOpenidClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientCreate,
		Read:   resourceKeycloakOpenidClientRead,
		Delete: resourceKeycloakOpenidClientDelete,
		Update: resourceKeycloakOpenidClientUpdate,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientImport,
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientAccessTypes, false),
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"standard_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"implicit_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"direct_access_grants_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service_accounts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"authorization_services_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"valid_redirect_uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"web_origins": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"resource": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uris": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"icon_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_managed_access": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"scopes": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"attributes": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func getOpenidClientFromData(data *schema.ResourceData) (*keycloak.OpenidClient, *keycloak.OpenidClientResourcesDiff, error) {
	validRedirectUris := make([]string, 0)
	webOrigins := make([]string, 0)

	if v, ok := data.GetOk("valid_redirect_uris"); ok {
		for _, validRedirectUri := range v.(*schema.Set).List() {
			validRedirectUris = append(validRedirectUris, validRedirectUri.(string))
		}
	}

	if v, ok := data.GetOk("web_origins"); ok {
		for _, webOrigin := range v.(*schema.Set).List() {
			webOrigins = append(webOrigins, webOrigin.(string))
		}
	}

	openidClient := &keycloak.OpenidClient{
		Id:                           data.Id(),
		ClientId:                     data.Get("client_id").(string),
		RealmId:                      data.Get("realm_id").(string),
		Name:                         data.Get("name").(string),
		Enabled:                      data.Get("enabled").(bool),
		Description:                  data.Get("description").(string),
		AuthorizationServicesEnabled: data.Get("authorization_services_enabled").(bool),
		ClientSecret:                 data.Get("client_secret").(string),
		StandardFlowEnabled:          data.Get("standard_flow_enabled").(bool),
		ImplicitFlowEnabled:          data.Get("implicit_flow_enabled").(bool),
		DirectAccessGrantsEnabled:    data.Get("direct_access_grants_enabled").(bool),
		ServiceAccountsEnabled:       data.Get("service_accounts_enabled").(bool),
		ValidRedirectUris:            validRedirectUris,
		WebOrigins:                   webOrigins,
	}

	if !openidClient.ImplicitFlowEnabled && !openidClient.StandardFlowEnabled {
		if _, ok := data.GetOk("valid_redirect_uris"); ok {
			return nil, nil, errors.New("valid_redirect_uris cannot be set when standard or implicit flow is not enabled")
		}
	}

	if !openidClient.ImplicitFlowEnabled && !openidClient.StandardFlowEnabled && !openidClient.DirectAccessGrantsEnabled {
		if _, ok := data.GetOk("web_origins"); ok {
			return nil, nil, errors.New("web_origins cannot be set when standard or implicit flow is not enabled")
		}
	}

	// access type
	if accessType := data.Get("access_type").(string); accessType == "PUBLIC" {
		openidClient.PublicClient = true
	} else if accessType == "BEARER-ONLY" {
		openidClient.BearerOnly = true
	}

	resources := &keycloak.OpenidClientResourcesDiff{}
	unchangedResourcesData := new(schema.Set)
	if v, ok := data.GetOk("resource"); ok && data.HasChange("resource") {
		if openidClient.AuthorizationServicesEnabled {
			unchangedResourcesData = v.(*schema.Set)
		} else {
			return nil, nil, errors.New("Resources cannot be managed when aunothiztion is not enabled")
		}
	}
	if data.HasChange("resource") {
		o, n := data.GetChange("resource")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		unchangedResourcesData = unchangedResourcesData.Difference(n.(*schema.Set))
		removeResourcesData := o.(*schema.Set).Difference(n.(*schema.Set)).Difference(unchangedResourcesData).List()
		resources.Remove = *getOpenidClientResourcesFromData(removeResourcesData)
		addResourcesData := n.(*schema.Set).Difference(o.(*schema.Set)).Difference(unchangedResourcesData).List()
		resources.Add = *getOpenidClientResourcesFromData(addResourcesData)
	}

	resources.Unchanged = *getOpenidClientResourcesFromData(unchangedResourcesData.List())

	return openidClient, resources, nil
}

func getOpenidClientResourcesFromData(data []interface{}) *keycloak.OpenidClientResources {
	var resources keycloak.OpenidClientResources
	for _, d := range data {
		resourceData := d.(map[string]interface{})
		var uris []string
		var scopes []string
		attributes := map[string][]string{}
		if v, ok := resourceData["uris"]; ok {
			for _, uri := range v.([]interface{}) {
				uris = append(uris, uri.(string))
			}
		}
		if v, ok := resourceData["scopes"]; ok {
			for _, scope := range v.([]interface{}) {
				scopes = append(scopes, scope.(string))
			}
		}
		if v, ok := resourceData["attributes"]; ok {
			for key, value := range v.(map[string]interface{}) {
				attributes[key] = strings.Split(value.(string), ",")
			}
		}
		resource := keycloak.OpenidClientResource{
			DisplayName:        resourceData["display_name"].(string),
			Name:               resourceData["name"].(string),
			IconUri:            resourceData["icon_uri"].(string),
			OwnerManagedAccess: resourceData["owner_managed_access"].(bool),
			Id:                 resourceData["id"].(string),
			Uris:               uris,
			Scopes:             scopes,
			Attributes:         attributes,
		}
		resources = append(resources, resource)
	}
	return &resources
}

func setOpenidClientData(data *schema.ResourceData, client *keycloak.OpenidClient, resources *keycloak.OpenidClientResources) {
	data.SetId(client.Id)

	data.Set("client_id", client.ClientId)
	data.Set("realm_id", client.RealmId)
	data.Set("name", client.Name)
	data.Set("enabled", client.Enabled)
	data.Set("description", client.Description)
	data.Set("client_secret", client.ClientSecret)
	data.Set("standard_flow_enabled", client.StandardFlowEnabled)
	data.Set("implicit_flow_enabled", client.ImplicitFlowEnabled)
	data.Set("direct_access_grants_enabled", client.DirectAccessGrantsEnabled)
	data.Set("service_accounts_enabled", client.ServiceAccountsEnabled)
	data.Set("valid_redirect_uris", client.ValidRedirectUris)
	data.Set("web_origins", client.WebOrigins)
	data.Set("authorization_services_enabled", client.AuthorizationServicesEnabled)

	// access type
	if client.PublicClient {
		data.Set("access_type", "PUBLIC")
	} else if client.BearerOnly {
		data.Set("access_type", "BEARER-ONLY")
	} else {
		data.Set("access_type", "CONFIDENTIAL")
	}

	resourcesData := []interface{}{}

	for _, resource := range *resources {
		resourceData := map[string]interface{}{
			"display_name":         resource.DisplayName,
			"name":                 resource.Name,
			"uris":                 resource.Uris,
			"icon_uri":             resource.IconUri,
			"owner_managed_access": resource.OwnerManagedAccess,
			"scopes":               resource.Scopes,
			"attributes":           listValueToStr(resource.Attributes),
			"id":                   resource.Id,
		}
		resourcesData = append(resourcesData, resourceData)
	}
	data.Set("resource", resourcesData)
}

func resourceKeycloakOpenidClientCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client, resources, err := getOpenidClientFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenidClient(client)
	if err != nil {
		return err
	}

	for i := 0; i < len((*resources).Add); i++ {
		err = keycloakClient.NewOpenidClientResource(client, &((*resources).Add)[i])
		if err != nil {
			return err
		}
	}

	setOpenidClientData(data, client, &resources.Add)

	return resourceKeycloakOpenidClientRead(data, meta)
}

func resourceKeycloakOpenidClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	client, err := keycloakClient.GetOpenidClient(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	resources := &keycloak.OpenidClientResources{}

	if client.AuthorizationServicesEnabled {
		resources, err = keycloakClient.GetOpenidClientResources(realmId, id)
		if err != nil {
			return handleNotFoundError(err, data)
		}
	}

	setOpenidClientData(data, client, resources)

	return nil
}

func resourceKeycloakOpenidClientUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client, resources, err := getOpenidClientFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenidClient(client)
	if err != nil {
		return err
	}

	if client.AuthorizationServicesEnabled {
		for _, resource := range (*resources).Unchanged {
			err = keycloakClient.UpdateOpenidClientResource(client, &resource)
			if err != nil {
				return err
			}
		}

		for _, resource := range (*resources).Add {
			err = keycloakClient.NewOpenidClientResource(client, &resource)
			if err != nil {
				return err
			}
		}

		for _, resource := range (*resources).Remove {
			err = keycloakClient.DeleteOpenidClientResource(client.RealmId, client.Id, resource.Id)
			if err != nil {
				return err
			}
		}
	}

	setOpenidClientData(data, client, &resources.Unchanged)

	return nil
}

func resourceKeycloakOpenidClientDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClient(realmId, id)
}

func resourceKeycloakOpenidClientImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientId}}")
	}
	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
