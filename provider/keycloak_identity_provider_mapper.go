package provider

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakIdentityProviderMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakIdentityProviderMapperCreate,
		Read:   resourceKeycloakIdentityProviderMapperRead,
		Update: resourceKeycloakIdentityProviderMapperUpdate,
		Delete: resourceKeycloakIdentityProviderMapperDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakIdentityProviderMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Realm Name",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IDP Mapper Name",
			},
			"identity_provider_alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IDP Alias",
			},
			"identity_provider_mapper": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IDP Mapper Type",
			},
			"oidc": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"saml", "social"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Does the external IDP support backchannel logout?",
						},
						"attribute_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Use JWKS url",
						},
					},
				},
				Set: configMapperHash,
			},
			"saml": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"oidc", "social"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Does the external IDP support backchannel logout?",
						},
					},
				},
				Set: configMapperHash,
			},
			"social": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"oidc", "saml"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"json_field": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IDP name",
						},
						"provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Provider name",
						},
						"user_attribute": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IDP name",
						},
					},
				},
				Set: configMapperHash,
			},
		},
	}
}

func getIdentityProviderMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec := &keycloak.IdentityProviderMapper{
		Id:                     data.Id(),
		Realm:                  data.Get("realm").(string),
		Name:                   data.Get("name").(string),
		IdentityProviderAlias:  data.Get("identity_provider_alias").(string),
		IdentityProviderMapper: fmt.Sprintf("%s-%s", data.Get("identity_provider_alias").(string), data.Get("identity_provider_mapper").(string)),
	}
	if v, ok := data.GetOk("oidc"); ok {
		rec.Provider = "oidc"
		configs := v.(*schema.Set).List()
		if len(configs) == 1 {
			config := configs[0].(map[string]interface{})
			rec.Config = &keycloak.IdentityProviderMapperConfig{
				Attribute:      config["attribute"].(string),
				AttributeValue: config["attribute_value"].(string),
			}
		}
	} else if v, ok := data.GetOk("social"); ok {
		configs := v.(*schema.Set).List()
		if len(configs) == 1 {
			config := configs[0].(map[string]interface{})
			rec.Provider = config["provider"].(string)
			rec.Config = &keycloak.IdentityProviderMapperConfig{
				JsonField:     config["json_field"].(string),
				UserAttribute: config["user_attribute"].(string),
			}
		}
	} else if v, ok := data.GetOk("saml"); ok {
		rec.Provider = "saml"
		configs := v.(*schema.Set).List()
		if len(configs) == 1 {
			config := configs[0].(map[string]interface{})
			rec.Config = &keycloak.IdentityProviderMapperConfig{
				Template: config["template"].(string),
			}
		}
	} else {
		return nil, fmt.Errorf("No provider config is defined. Please add social, saml or oidc provider")
	}
	return rec, nil
}

func setIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	data.SetId(identityProviderMapper.Id)
	data.Set("id", identityProviderMapper.Id)
	data.Set("realm", identityProviderMapper.Realm)
	data.Set("name", identityProviderMapper.Name)
	data.Set("identity_provider_alias", identityProviderMapper.IdentityProviderAlias)
	data.Set("identity_provider_mapper", identityProviderMapper.IdentityProviderMapper)
	if config := identityProviderMapper.Config; config != nil {
		switch identityProviderMapper.Provider {
		case "oidc":
			data.Set("provider", "oidc")
			data.Set("config", []interface{}{
				map[string]interface{}{
					"attribute":       config.Attribute,
					"attribute_value": config.AttributeValue,
				},
			})
		case "facebook", "stackoverflow", "twitter", "github", "gitlab", "instagram", "bitbucket", "google", "microsoft", "paypal":
			data.Set("provider", identityProviderMapper.Provider)
			data.Set("config", []interface{}{
				map[string]interface{}{
					"json_field":     config.JsonField,
					"user_attribute": config.UserAttribute,
				},
			})
		case "saml":
			data.Set("provider", "saml")
			data.Set("config", []interface{}{
				map[string]interface{}{
					"template": config.Template,
				},
			})
		default:
			return fmt.Errorf("No provider config is defined. Please add social, saml or oidc provider")
		}
	}
	return nil
}

func resourceKeycloakIdentityProviderMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	identityProvider, err := getIdentityProviderMapperFromData(data)

	err = keycloakClient.NewIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}

	setIdentityProviderMapperData(data, identityProvider)

	return resourceKeycloakIdentityProviderMapperRead(data, meta)
}

func resourceKeycloakIdentityProviderMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()

	identityProvider, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setIdentityProviderMapperData(data, identityProvider)

	return nil
}

func resourceKeycloakIdentityProviderMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	identityProvider, err := getIdentityProviderMapperFromData(data)

	err = keycloakClient.UpdateIdentityProviderMapper(identityProvider)
	if err != nil {
		return err
	}

	setIdentityProviderMapperData(data, identityProvider)

	return nil
}

func resourceKeycloakIdentityProviderMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	alias := data.Get("identity_provider_alias").(string)
	id := data.Id()

	return keycloakClient.DeleteIdentityProviderMapper(realm, alias, id)
}

func configMapperHash(v interface{}) int {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", "config"))
	return hashcode.String(buf.String())
}

func resourceKeycloakIdentityProviderMapperImport(data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	data.Id()
	return []*schema.ResourceData{data}, nil
}
