package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakGroupCreate,
		Read:   resourceKeycloakGroupRead,
		Delete: resourceKeycloakGroupDelete,
		Update: resourceKeycloakGroupUpdate,
		// This resource can be imported using {{realm}}/{{group_id}}. The Group ID is displayed in the URL when editing it from the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func mapFromDataToGroup(data *schema.ResourceData) *keycloak.Group {
	attributes := map[string][]string{}
	if v, ok := data.GetOk("attributes"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = strings.Split(value.(string), MULTIVALUE_ATTRIBUTE_SEPARATOR)
		}
	}

	group := &keycloak.Group{
		Id:         data.Id(),
		RealmId:    data.Get("realm_id").(string),
		ParentId:   data.Get("parent_id").(string),
		Name:       data.Get("name").(string),
		Attributes: attributes,
	}

	return group
}

func mapFromGroupToData(data *schema.ResourceData, group *keycloak.Group) {
	attributes := map[string]string{}
	for k, v := range group.Attributes {
		attributes[k] = strings.Join(v, MULTIVALUE_ATTRIBUTE_SEPARATOR)
	}
	data.SetId(group.Id)
	data.Set("realm_id", group.RealmId)
	data.Set("name", group.Name)
	data.Set("path", group.Path)
	data.Set("attributes", attributes)
	if group.ParentId != "" {
		data.Set("parent_id", group.ParentId)
	}
}

func resourceKeycloakGroupCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	group := mapFromDataToGroup(data)

	err := keycloakClient.NewGroup(group)
	if err != nil {
		return err
	}

	mapFromGroupToData(data, group)

	return resourceKeycloakGroupRead(data, meta)
}

func resourceKeycloakGroupRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	group, err := keycloakClient.GetGroup(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromGroupToData(data, group)

	return nil
}

func resourceKeycloakGroupUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	group := mapFromDataToGroup(data)

	err := keycloakClient.UpdateGroup(group)
	if err != nil {
		return err
	}

	mapFromGroupToData(data, group)

	return nil
}

func resourceKeycloakGroupDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteGroup(realmId, id)
}

func resourceKeycloakGroupImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{groupId}}")
	}
	d.Set("realm_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
