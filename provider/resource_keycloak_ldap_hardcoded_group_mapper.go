package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakLdapHardcodedGroupMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapHardcodedGroupMapperCreate,
		Read:   resourceKeycloakLdapHardcodedGroupMapperRead,
		Update: resourceKeycloakLdapHardcodedGroupMapperUpdate,
		Delete: resourceKeycloakLdapHardcodedGroupMapperDelete,
		// This resource can be imported using {{realm}}/{{provider_id}}/{{mapper_id}}. The Provider and Mapper IDs are displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakLdapGenericMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the mapper when displayed in the console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm in which the ldap user federation provider exists.",
			},
			"ldap_user_federation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ldap user federation provider to attach this mapper to.",
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Group to grant to user.",
			},
		},
	}
}

func getLdapHardcodedGroupMapperFromData(data *schema.ResourceData) *keycloak.LdapHardcodedGroupMapper {
	return &keycloak.LdapHardcodedGroupMapper{
		Id:                   data.Id(),
		Name:                 data.Get("name").(string),
		RealmId:              data.Get("realm_id").(string),
		LdapUserFederationId: data.Get("ldap_user_federation_id").(string),
		Group:                data.Get("group").(string),
	}
}

func setLdapHardcodedGroupMapperData(data *schema.ResourceData, ldapMapper *keycloak.LdapHardcodedGroupMapper) {
	data.SetId(ldapMapper.Id)
	data.Set("name", ldapMapper.Name)
	data.Set("realm_id", ldapMapper.RealmId)
	data.Set("ldap_user_federation_id", ldapMapper.LdapUserFederationId)
	data.Set("group", ldapMapper.Group)
}

func resourceKeycloakLdapHardcodedGroupMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMapper := getLdapHardcodedGroupMapperFromData(data)

	err := keycloakClient.ValidateLdapHardcodedGroupMapper(ldapMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewLdapHardcodedGroupMapper(ldapMapper)
	if err != nil {
		return err
	}

	setLdapHardcodedGroupMapperData(data, ldapMapper)

	return resourceKeycloakLdapHardcodedGroupMapperRead(data, meta)
}

func resourceKeycloakLdapHardcodedGroupMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapMapper, err := keycloakClient.GetLdapHardcodedGroupMapper(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setLdapHardcodedGroupMapperData(data, ldapMapper)

	return nil
}

func resourceKeycloakLdapHardcodedGroupMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMapper := getLdapHardcodedGroupMapperFromData(data)

	err := keycloakClient.ValidateLdapHardcodedGroupMapper(ldapMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateLdapHardcodedGroupMapper(ldapMapper)
	if err != nil {
		return err
	}

	setLdapHardcodedGroupMapperData(data, ldapMapper)

	return nil
}

func resourceKeycloakLdapHardcodedGroupMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapHardcodedGroupMapper(realmId, id)
}
