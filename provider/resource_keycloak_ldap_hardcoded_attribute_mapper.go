package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakLdapHardcodedAttributeMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapHardcodedAttributeMapperCreate,
		Read:   resourceKeycloakLdapHardcodedAttributeMapperRead,
		Update: resourceKeycloakLdapHardcodedAttributeMapperUpdate,
		Delete: resourceKeycloakLdapHardcodedAttributeMapperDelete,
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
			"attribute_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the LDAP attribute",
			},
			"attribute_value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Value of the LDAP attribute. You can either hardcode any value like 'foo' but you can also use the special token '${RANDOM}', which will be replaced with some randomly generated String.",
			},
		},
	}
}

func getLdapHardcodedAttributeMapperFromData(data *schema.ResourceData) *keycloak.LdapHardcodedAttributeMapper {
	return &keycloak.LdapHardcodedAttributeMapper{
		Id:                   data.Id(),
		Name:                 data.Get("name").(string),
		RealmId:              data.Get("realm_id").(string),
		LdapUserFederationId: data.Get("ldap_user_federation_id").(string),
		AttributeName:        data.Get("attribute_name").(string),
		AttributeValue:       data.Get("attribute_value").(string),
	}
}

func setLdapHardcodedAttributeMapperData(data *schema.ResourceData, ldapMapper *keycloak.LdapHardcodedAttributeMapper) {
	data.SetId(ldapMapper.Id)
	data.Set("name", ldapMapper.Name)
	data.Set("realm_id", ldapMapper.RealmId)
	data.Set("ldap_user_federation_id", ldapMapper.LdapUserFederationId)
	data.Set("attribute_name", ldapMapper.AttributeName)
	data.Set("attribute_value", ldapMapper.AttributeValue)
}

func resourceKeycloakLdapHardcodedAttributeMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMapper := getLdapHardcodedAttributeMapperFromData(data)

	err := keycloakClient.ValidateLdapHardcodedAttributeMapper(ldapMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewLdapHardcodedAttributeMapper(ldapMapper)
	if err != nil {
		return err
	}

	setLdapHardcodedAttributeMapperData(data, ldapMapper)

	return resourceKeycloakLdapHardcodedAttributeMapperRead(data, meta)
}

func resourceKeycloakLdapHardcodedAttributeMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapMapper, err := keycloakClient.GetLdapHardcodedAttributeMapper(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setLdapHardcodedAttributeMapperData(data, ldapMapper)

	return nil
}

func resourceKeycloakLdapHardcodedAttributeMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMapper := getLdapHardcodedAttributeMapperFromData(data)

	err := keycloakClient.ValidateLdapHardcodedAttributeMapper(ldapMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateLdapHardcodedAttributeMapper(ldapMapper)
	if err != nil {
		return err
	}

	setLdapHardcodedAttributeMapperData(data, ldapMapper)

	return nil
}

func resourceKeycloakLdapHardcodedAttributeMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapHardcodedAttributeMapper(realmId, id)
}
