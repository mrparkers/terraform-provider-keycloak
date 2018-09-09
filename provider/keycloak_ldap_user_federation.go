package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"log"
)

func resourceKeycloakLdapUserFederation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapUserFederationCreate,
		Read:   resourceKeycloakLdapUserFederationRead,
		Update: resourceKeycloakLdapUserFederationUpdate,
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
				Optional:    true,
				Default:     0,
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
				Optional:     true,
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
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(_, remoteBindCredential, _ string, _ *schema.ResourceData) bool {
					return remoteBindCredential == "**********"
				},
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
				Default:  30,
			},
			"read_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
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
	var userObjectClasses []string

	for _, userObjectClass := range data.Get("user_object_classes").([]interface{}) {
		userObjectClasses = append(userObjectClasses, userObjectClass.(string))
	}

	log.Printf("[DEBUG] bind_dn: %s", data.Get("bind_dn").(string))
	log.Printf("[DEBUG] bind_credential: %s", data.Get("bind_credential").(string))

	return &keycloak.LdapUserFederation{
		Id:      data.Id(),
		Name:    data.Get("name").(string),
		RealmId: data.Get("realm_id").(string),

		Enabled:  data.Get("enabled").(bool),
		Priority: data.Get("priority").(int),

		ImportEnabled:     data.Get("import_enabled").(bool),
		EditMode:          data.Get("edit_mode").(string),
		SyncRegistrations: data.Get("sync_registrations").(bool),

		Vendor:                 data.Get("vendor").(string),
		UsernameLDAPAttribute:  data.Get("username_ldap_attribute").(string),
		RdnLDAPAttribute:       data.Get("rdn_ldap_attribute").(string),
		UuidLDAPAttribute:      data.Get("uuid_ldap_attribute").(string),
		UserObjectClasses:      userObjectClasses,
		ConnectionUrl:          data.Get("connection_url").(string),
		UsersDn:                data.Get("users_dn").(string),
		BindDn:                 data.Get("bind_dn").(string),
		BindCredential:         data.Get("bind_credential").(string),
		CustomUserSearchFilter: data.Get("custom_user_search_filter").(string),
		SearchScope:            data.Get("search_scope").(string),

		ValidatePasswordPolicy: data.Get("validate_password_policy").(bool),
		UseTruststoreSpi:       data.Get("use_truststore_spi").(string),
		ConnectionTimeout:      data.Get("connection_timeout").(int),
		ReadTimeout:            data.Get("read_timeout").(int),
		Pagination:             data.Get("pagination").(bool),

		BatchSizeForSync:  data.Get("batch_size_for_sync").(int),
		FullSyncPeriod:    data.Get("full_sync_period").(int),
		ChangedSyncPeriod: data.Get("changed_sync_period").(int),

		CachePolicy: data.Get("cache_policy").(string),
	}
}

func setLdapUserFederationData(data *schema.ResourceData, ldap *keycloak.LdapUserFederation) {
	data.SetId(ldap.Id)

	data.Set("name", ldap.Name)
	data.Set("realm_id", ldap.RealmId)

	data.Set("enabled", ldap.Enabled)
	data.Set("priority", ldap.Priority)

	data.Set("import_enabled", ldap.ImportEnabled)
	data.Set("edit_mode", ldap.EditMode)
	data.Set("sync_registrations", ldap.SyncRegistrations)

	data.Set("vendor", ldap.Vendor)
	data.Set("username_ldap_attribute", ldap.UsernameLDAPAttribute)
	data.Set("rdn_ldap_attribute", ldap.RdnLDAPAttribute)
	data.Set("uuid_ldap_attribute", ldap.UuidLDAPAttribute)
	data.Set("user_object_classes", ldap.UserObjectClasses)
	data.Set("connection_url", ldap.ConnectionUrl)
	data.Set("users_dn", ldap.UsersDn)
	data.Set("bind_dn", ldap.BindDn)
	data.Set("bind_credential", ldap.BindCredential)
	data.Set("custom_user_search_filter", ldap.CustomUserSearchFilter)
	data.Set("search_scope", ldap.SearchScope)

	data.Set("validate_password_policy", ldap.ValidatePasswordPolicy)
	data.Set("use_truststore_spi", ldap.UseTruststoreSpi)
	data.Set("connection_timeout", ldap.ConnectionTimeout)
	data.Set("read_timeout", ldap.ReadTimeout)
	data.Set("pagination", ldap.Pagination)

	data.Set("batch_size_for_sync", ldap.BatchSizeForSync)
	data.Set("full_sync_period", ldap.FullSyncPeriod)
	data.Set("changed_sync_period", ldap.ChangedSyncPeriod)

	data.Set("cache_policy", ldap.CachePolicy)
}

func resourceKeycloakLdapUserFederationCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldap := getLdapUserFederationFromData(data)

	err := ldap.Validate()
	if err != nil {
		return err
	}

	err = keycloakClient.NewLdapUserFederation(ldap)
	if err != nil {
		return err
	}

	setLdapUserFederationData(data, ldap)

	return resourceKeycloakLdapUserFederationRead(data, meta)
}

func resourceKeycloakLdapUserFederationRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldap, err := keycloakClient.GetLdapUserFederation(realmId, id)
	if err != nil {
		return err
	}

	setLdapUserFederationData(data, ldap)

	return nil
}

func resourceKeycloakLdapUserFederationUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldap := getLdapUserFederationFromData(data)

	err := ldap.Validate()
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateLdapUserFederation(ldap)
	if err != nil {
		return err
	}

	setLdapUserFederationData(data, ldap)

	return nil
}

func resourceKeycloakLdapUserFederationDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapUserFederation(realmId, id)
}
