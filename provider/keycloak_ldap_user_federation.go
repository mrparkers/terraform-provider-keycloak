package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakLdapUserFederation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapUserFederationCreate,
		Read:   resourceKeycloakLdapUserFederationRead,
		Delete: resourceKeycloakLdapUserFederationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the provider when displayed in the console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm this provider will provider user federation for.",
			},

			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "When false, this provider will not be used when performing queries for users.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority of this provider when looking up users. Lower values are first.",
			},

			"import_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When true, LDAP users will be imported into the Keycloak database.",
			},
			"edit_mode": {
				Type:         schema.TypeString,
				Required:     true,
				Default:      "READ_ONLY",
				ValidateFunc: validation.StringInSlice([]string{"READ_ONLY", "WRITABLE", "UNSYNCED"}, false),
				Description:  "READ_ONLY and WRITABLE are self-explanatory. UNSYNCED allowed user data to be imported but not synced back to LDAP.",
			},
			"sync_registrations": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, newly created users will be synced back to LDAP.",
			},

			"vendor": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "OTHER",
				ValidateFunc: validation.StringInSlice([]string{"OTHER", "EDIRECTORY", "AD", "RHDS", "TIVOLI"}, false),
				Description:  "LDAP vendor. I am almost certain this field does nothing, but the UI indicates that it is required.",
			},
			"username_ldap_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the LDAP attribute to use as the Keycloak username.",
			},
			"rdn_ldap_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the LDAP attribute to use as the relative distinguished name..",
			},
			"uuid_ldap_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the LDAP attribute to use as a unique object identifier for objects in LDAP.",
			},
			"user_object_classes": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem:        schema.TypeString,
				Description: "All values of LDAP objectClass attribute for users in LDAP",
			},
			"connection_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bind_credential": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_user_search_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ONE_LEVEL",
				ValidateFunc: validation.StringInSlice([]string{"ONE_LEVEL", "SUBTREE"}, false),
			},

			"validate_password_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"use_truststore_spi": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ONLY_FOR_LDAPS",
				ValidateFunc: validation.StringInSlice([]string{"ALWAYS", "ONLY_FOR_LDAPS", "NEVER"}, false),
			},
			"connection_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"read_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pagination": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"batch_size_for_sync": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},
			"full_sync_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"changed_sync_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},

			"cache_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "EVICT_DAILY", "EVICT_WEEKLY", "MAX_LIFESPAN", "NO_CACHE"}, false),
			},
		},
	}
}

func getLdapUserFederationFromData(data *schema.ResourceData) *keycloak.LdapUserFederation {
	return &keycloak.LdapUserFederation{
		Id:      data.Id(),
		RealmId: data.Get("realm_id").(string),
	}
}

func setLdapUserFederationData(data *schema.ResourceData, client *keycloak.LdapUserFederation) {
	data.SetId(client.Id)

	data.Set("realm_id", client.RealmId)
}

func resourceKeycloakLdapUserFederationCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapUserFederation := getLdapUserFederationFromData(data)

	err := keycloakClient.NewLdapUserFederation(ldapUserFederation)
	if err != nil {
		return err
	}

	setLdapUserFederationData(data, ldapUserFederation)

	return resourceKeycloakLdapUserFederationRead(data, meta)
}

func resourceKeycloakLdapUserFederationRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapUserFederation, err := keycloakClient.GetLdapUserFederation(realmId, id)
	if err != nil {
		return err
	}

	setLdapUserFederationData(data, ldapUserFederation)

	return nil
}

func resourceKeycloakLdapUserFederationDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapUserFederation(realmId, id)
}
