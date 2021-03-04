package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRoleCreate,
		Read:   resourceKeycloakRoleRead,
		Delete: resourceKeycloakRoleDelete,
		Update: resourceKeycloakRoleUpdate,
		// This resource can be imported using {{realm}}/{{roleId}}. The role's ID (a GUID) can be found in the URL when viewing the role
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"composite_roles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				Set:      schema.HashString,
				Optional: true,
			},
			// misc attributes
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func mapFromDataToRole(data *schema.ResourceData) *keycloak.Role {
	attributes := map[string][]string{}
	if v, ok := data.GetOk("attributes"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = splitLen(value.(string), MAX_ATTRIBUTE_VALUE_LEN)
		}
	}

	role := &keycloak.Role{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		ClientId:    data.Get("client_id").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Attributes:  attributes,
	}

	return role
}

func mapFromRoleToData(data *schema.ResourceData, role *keycloak.Role) {
	attributes := map[string]string{}
	for k, v := range role.Attributes {
		attributes[k] = strings.Join(v, "")
	}
	data.SetId(role.Id)

	data.Set("realm_id", role.RealmId)
	data.Set("client_id", role.ClientId)
	data.Set("name", role.Name)
	data.Set("description", role.Description)
	data.Set("attributes", attributes)
}

func resourceKeycloakRoleCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	role := mapFromDataToRole(data)

	var compositeRoles []*keycloak.Role
	if v, ok := data.GetOk("composite_roles"); ok {
		compositeRolesTf := v.(*schema.Set).List()

		for _, compositeRoleId := range compositeRolesTf {
			compositeRoleToAdd, err := keycloakClient.GetRole(role.RealmId, compositeRoleId.(string))
			if err != nil {
				return err
			}

			compositeRoles = append(compositeRoles, compositeRoleToAdd)
		}

		if len(compositeRoles) != 0 { // technically you can still specify composite_roles = [] in HCL
			role.Composite = true
		}
	}

	err := keycloakClient.CreateRole(role)
	if err != nil {
		return err
	}

	if role.Composite {
		err = keycloakClient.AddCompositesToRole(role, compositeRoles)
		if err != nil {
			return err
		}
	}

	mapFromRoleToData(data, role)

	return resourceKeycloakRoleRead(data, meta)
}

func resourceKeycloakRoleRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	role, err := keycloakClient.GetRole(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromRoleToData(data, role)

	if role.Composite {
		composites, err := keycloakClient.GetRoleComposites(role)
		if err != nil {
			return err
		}

		var compositeRoleIds []string

		for _, composite := range composites {
			compositeRoleIds = append(compositeRoleIds, composite.Id)
		}

		data.Set("composite_roles", compositeRoleIds)
	}

	return nil
}

func resourceKeycloakRoleUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	role := mapFromDataToRole(data)

	err := keycloakClient.UpdateRole(role)
	if err != nil {
		return err
	}

	keycloakComposites, err := keycloakClient.GetRoleComposites(role)
	if err != nil {
		return err
	}

	if v, ok := data.GetOk("composite_roles"); ok {
		tfCompositeIds := v.(*schema.Set)
		var keycloakCompositesToRemove []*keycloak.Role

		// get a list of all composites to remove and all composites to add
		for _, keycloakComposite := range keycloakComposites {
			if tfCompositeIds.Contains(keycloakComposite.Id) {
				// if the composite exists in keycloak and tf state, we can remove them from the local list because this role does not need to be added
				tfCompositeIds.Remove(keycloakComposite.Id)
			} else {
				// if the composite exists in keycloak but not tf state, it needs to be removed on keycloak's side
				keycloakCompositesToRemove = append(keycloakCompositesToRemove, keycloakComposite)
			}
		}

		// at this point we have two slices:
		// `keycloakCompositesToRemove` should be removed from the role's list of composites
		// `tfCompositeIds` should be added to the role's list of composites. all of the roles that exist on both sides have already been removed

		if len(keycloakCompositesToRemove) != 0 {
			err = keycloakClient.RemoveCompositesFromRole(role, keycloakCompositesToRemove)
			if err != nil {
				return err
			}
		}

		if tfCompositeIds.Len() != 0 {
			var compositesToAdd []*keycloak.Role
			for _, tfCompositeId := range tfCompositeIds.List() {
				compositeToAdd, err := keycloakClient.GetRole(role.RealmId, tfCompositeId.(string))
				if err != nil {
					return err
				}

				compositesToAdd = append(compositesToAdd, compositeToAdd)
			}

			err = keycloakClient.AddCompositesToRole(role, compositesToAdd)
			if err != nil {
				return err
			}
		}
	} else {
		// the user wants this role to have zero composites. if there are composites attached, remove them
		if len(keycloakComposites) != 0 {
			err = keycloakClient.RemoveCompositesFromRole(role, keycloakComposites)
			if err != nil {
				return err
			}
		}
	}

	return resourceKeycloakRoleRead(data, meta)
}

func resourceKeycloakRoleDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRole(realmId, id)
}

func resourceKeycloakRoleImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm}}/{{roleId}}.")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
