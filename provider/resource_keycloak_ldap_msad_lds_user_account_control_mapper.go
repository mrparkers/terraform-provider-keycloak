package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakLdapMsadLdsUserAccountControlMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapMsadLdsUserAccountControlMapperCreate,
		Read:   resourceKeycloakLdapMsadLdsUserAccountControlMapperRead,
		Update: resourceKeycloakLdapMsadLdsUserAccountControlMapperUpdate,
		Delete: resourceKeycloakLdapMsadLdsUserAccountControlMapperDelete,
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
		},
	}
}

func getLdapMsadLdsUserAccountControlMapperFromData(data *schema.ResourceData) *keycloak.LdapMsadLdsUserAccountControlMapper {
	return &keycloak.LdapMsadLdsUserAccountControlMapper{
		Id:                   data.Id(),
		Name:                 data.Get("name").(string),
		RealmId:              data.Get("realm_id").(string),
		LdapUserFederationId: data.Get("ldap_user_federation_id").(string),
	}
}

func setLdapMsadLdsUserAccountControlMapperData(data *schema.ResourceData, ldapMsadLdsUserAccountControlMapper *keycloak.LdapMsadLdsUserAccountControlMapper) {
	data.SetId(ldapMsadLdsUserAccountControlMapper.Id)

	data.Set("name", ldapMsadLdsUserAccountControlMapper.Name)
	data.Set("realm_id", ldapMsadLdsUserAccountControlMapper.RealmId)
	data.Set("ldap_user_federation_id", ldapMsadLdsUserAccountControlMapper.LdapUserFederationId)
}

func resourceKeycloakLdapMsadLdsUserAccountControlMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMsadLdsUserAccountControlMapper := getLdapMsadLdsUserAccountControlMapperFromData(data)

	err := keycloakClient.NewLdapMsadLdsUserAccountControlMapper(ldapMsadLdsUserAccountControlMapper)
	if err != nil {
		return err
	}

	setLdapMsadLdsUserAccountControlMapperData(data, ldapMsadLdsUserAccountControlMapper)

	return resourceKeycloakLdapMsadLdsUserAccountControlMapperRead(data, meta)
}

func resourceKeycloakLdapMsadLdsUserAccountControlMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapMsadLdsUserAccountControlMapper, err := keycloakClient.GetLdapMsadLdsUserAccountControlMapper(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setLdapMsadLdsUserAccountControlMapperData(data, ldapMsadLdsUserAccountControlMapper)

	return nil
}

func resourceKeycloakLdapMsadLdsUserAccountControlMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMsadLdsUserAccountControlMapper := getLdapMsadLdsUserAccountControlMapperFromData(data)

	err := keycloakClient.UpdateLdapMsadLdsUserAccountControlMapper(ldapMsadLdsUserAccountControlMapper)
	if err != nil {
		return err
	}

	setLdapMsadLdsUserAccountControlMapperData(data, ldapMsadLdsUserAccountControlMapper)

	return nil
}

func resourceKeycloakLdapMsadLdsUserAccountControlMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapMsadLdsUserAccountControlMapper(realmId, id)
}
