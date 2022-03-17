package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakLdapUserAttributeMapper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakLdapUserAttributeMapperCreate,
		ReadContext:   resourceKeycloakLdapUserAttributeMapperRead,
		UpdateContext: resourceKeycloakLdapUserAttributeMapperUpdate,
		DeleteContext: resourceKeycloakLdapUserAttributeMapperDelete,
		// This resource can be imported using {{realm}}/{{provider_id}}/{{mapper_id}}. The Provider and Mapper IDs are displayed in the GUI
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakLdapGenericMapperImport,
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
			"user_model_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the UserModel property or attribute you want to map the LDAP attribute into.",
			},
			"ldap_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the mapped attribute on LDAP object.",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this attribute is not saved back to LDAP when the user attribute is updated in Keycloak.",
			},
			"always_read_value_from_ldap": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, the value fetched from LDAP will override the value stored in Keycloak.",
			},
			"is_mandatory_in_ldap": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this attribute must exist in LDAP.",
			},
		},
	}
}

func getLdapUserAttributeMapperFromData(data *schema.ResourceData) *keycloak.LdapUserAttributeMapper {
	return &keycloak.LdapUserAttributeMapper{
		Id:                   data.Id(),
		Name:                 data.Get("name").(string),
		RealmId:              data.Get("realm_id").(string),
		LdapUserFederationId: data.Get("ldap_user_federation_id").(string),

		LdapAttribute:      data.Get("ldap_attribute").(string),
		UserModelAttribute: data.Get("user_model_attribute").(string),

		ReadOnly:                data.Get("read_only").(bool),
		AlwaysReadValueFromLdap: data.Get("always_read_value_from_ldap").(bool),
		IsMandatoryInLdap:       data.Get("is_mandatory_in_ldap").(bool),
	}
}

func setLdapUserAttributeMapperData(data *schema.ResourceData, ldapUserAttributeMapper *keycloak.LdapUserAttributeMapper) {
	data.SetId(ldapUserAttributeMapper.Id)

	data.Set("name", ldapUserAttributeMapper.Name)
	data.Set("realm_id", ldapUserAttributeMapper.RealmId)
	data.Set("ldap_user_federation_id", ldapUserAttributeMapper.LdapUserFederationId)

	data.Set("ldap_attribute", ldapUserAttributeMapper.LdapAttribute)
	data.Set("user_model_attribute", ldapUserAttributeMapper.UserModelAttribute)

	data.Set("read_only", ldapUserAttributeMapper.ReadOnly)
	data.Set("always_read_value_from_ldap", ldapUserAttributeMapper.AlwaysReadValueFromLdap)
	data.Set("is_mandatory_in_ldap", ldapUserAttributeMapper.IsMandatoryInLdap)
}

func resourceKeycloakLdapUserAttributeMapperCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapUserAttributeMapper := getLdapUserAttributeMapperFromData(data)

	err := keycloakClient.NewLdapUserAttributeMapper(ctx, ldapUserAttributeMapper)
	if err != nil {
		diag.FromErr(err)
	}

	setLdapUserAttributeMapperData(data, ldapUserAttributeMapper)

	return resourceKeycloakLdapUserAttributeMapperRead(ctx, data, meta)
}

func resourceKeycloakLdapUserAttributeMapperRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapUserAttributeMapper, err := keycloakClient.GetLdapUserAttributeMapper(ctx, realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setLdapUserAttributeMapperData(data, ldapUserAttributeMapper)

	return nil
}

func resourceKeycloakLdapUserAttributeMapperUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapUserAttributeMapper := getLdapUserAttributeMapperFromData(data)

	err := keycloakClient.UpdateLdapUserAttributeMapper(ctx, ldapUserAttributeMapper)
	if err != nil {
		diag.FromErr(err)
	}

	setLdapUserAttributeMapperData(data, ldapUserAttributeMapper)

	return nil
}

func resourceKeycloakLdapUserAttributeMapperDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return diag.FromErr(keycloakClient.DeleteLdapUserAttributeMapper(ctx, realmId, id))
}
