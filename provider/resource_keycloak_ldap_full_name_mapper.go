package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakLdapFullNameMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapFullNameMapperCreate,
		Read:   resourceKeycloakLdapFullNameMapperRead,
		Update: resourceKeycloakLdapFullNameMapperUpdate,
		Delete: resourceKeycloakLdapFullNameMapperDelete,
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
			"ldap_full_name_attribute": {
				Type:     schema.TypeString,
				Required: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"write_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func getLdapFullNameMapperFromData(data *schema.ResourceData) *keycloak.LdapFullNameMapper {
	return &keycloak.LdapFullNameMapper{
		Id:                   data.Id(),
		Name:                 data.Get("name").(string),
		RealmId:              data.Get("realm_id").(string),
		LdapUserFederationId: data.Get("ldap_user_federation_id").(string),

		LdapFullNameAttribute: data.Get("ldap_full_name_attribute").(string),
		ReadOnly:              data.Get("read_only").(bool),
		WriteOnly:             data.Get("write_only").(bool),
	}
}

func setLdapFullNameMapperData(data *schema.ResourceData, ldapFullNameMapper *keycloak.LdapFullNameMapper) {
	data.SetId(ldapFullNameMapper.Id)

	data.Set("id", ldapFullNameMapper.Id)
	data.Set("name", ldapFullNameMapper.Name)
	data.Set("realm_id", ldapFullNameMapper.RealmId)
	data.Set("ldap_user_federation_id", ldapFullNameMapper.LdapUserFederationId)

	data.Set("ldap_full_name_attribute", ldapFullNameMapper.LdapFullNameAttribute)
	data.Set("read_only", ldapFullNameMapper.ReadOnly)
	data.Set("write_only", ldapFullNameMapper.WriteOnly)
}

func resourceKeycloakLdapFullNameMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapFullNameMapper := getLdapFullNameMapperFromData(data)

	err := keycloakClient.ValidateLdapFullNameMapper(ldapFullNameMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewLdapFullNameMapper(ldapFullNameMapper)
	if err != nil {
		return err
	}

	setLdapFullNameMapperData(data, ldapFullNameMapper)

	return resourceKeycloakLdapFullNameMapperRead(data, meta)
}

func resourceKeycloakLdapFullNameMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapFullNameMapper, err := keycloakClient.GetLdapFullNameMapper(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setLdapFullNameMapperData(data, ldapFullNameMapper)

	return nil
}

func resourceKeycloakLdapFullNameMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapFullNameMapper := getLdapFullNameMapperFromData(data)

	err := keycloakClient.ValidateLdapFullNameMapper(ldapFullNameMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateLdapFullNameMapper(ldapFullNameMapper)
	if err != nil {
		return err
	}

	setLdapFullNameMapperData(data, ldapFullNameMapper)

	return nil
}

func resourceKeycloakLdapFullNameMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapFullNameMapper(realmId, id)
}

func resourceKeycloakLdapGenericMapperImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}/{{userFederationMapperId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("ldap_user_federation_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
