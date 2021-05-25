package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func resourceKeycloakCustomUserFederation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakCustomUserFederationCreate,
		Read:   resourceKeycloakCustomUserFederationRead,
		Update: resourceKeycloakCustomUserFederationUpdate,
		Delete: resourceKeycloakCustomUserFederationDelete,
		// This resource can be imported using {{realm}}/{{provider_id}}. The Provider ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakCustomUserFederationImport,
		},
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
				Description: "The realm (name) this provider will provide user federation for.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The parent_id of the generated component. will use realm_id if not specified.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					realmId := d.Get("realm_id").(string)
					if (old == "" && new == realmId) || (old == realmId && new == "") {
						return true
					}
					return false
				},
			},
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the custom provider, specified in the `getId` implementation for the UserStorageProviderFactory interface",
			},

			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When false, this provider will not be used when performing queries for users.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Priority of this provider when looking up users. Lower values are first.",
			},

			"cache_policy": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "DEFAULT",
				Deprecated:    "use cache.policy instead",
				ConflictsWith: []string{"cache"},
				ValidateFunc:  validation.StringInSlice(keycloakUserFederationCachePolicies, false),
			},
			"cache": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Description:   "Settings regarding cache policy for this realm.",
				ConflictsWith: []string{"cache_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "DEFAULT",
							ValidateFunc: validation.StringInSlice(keycloakUserFederationCachePolicies, false),
						},
						"max_lifespan": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressDurationStringDiff,
							Description:      "Max lifespan of cache entry (duration string).",
						},
						"eviction_day": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "-1",
							ValidateFunc: validation.All(validation.IntAtLeast(0), validation.IntAtMost(6)),
							Description:  "Day of the week the entry will become invalid on.",
						},
						"eviction_hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "-1",
							ValidateFunc: validation.All(validation.IntAtLeast(0), validation.IntAtMost(23)),
							Description:  "Hour of day the entry will become invalid on.",
						},
						"eviction_minute": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "-1",
							ValidateFunc: validation.All(validation.IntAtLeast(0), validation.IntAtMost(59)),
							Description:  "Minute of day the entry will become invalid on.",
						},
					},
				},
			},

			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func getCustomUserFederationFromData(data *schema.ResourceData) *keycloak.CustomUserFederation {
	config := map[string][]string{}
	if v, ok := data.GetOk("config"); ok {
		for key, value := range v.(map[string]interface{}) {
			config[key] = []string{value.(string)}
		}
	}
	parentId := ""
	dataParentId := data.Get("parent_id").(string)
	if dataParentId != "" {
		parentId = dataParentId
	} else {
		parentId = data.Get("realm_id").(string)
	}

	custom := &keycloak.CustomUserFederation{
		Id:         data.Id(),
		Name:       data.Get("name").(string),
		RealmId:    data.Get("realm_id").(string),
		ParentId:   parentId,
		ProviderId: data.Get("provider_id").(string),

		Enabled:  data.Get("enabled").(bool),
		Priority: data.Get("priority").(int),

		CachePolicy: data.Get("cache_policy").(string),

		Config: config,
	}

	if cache, ok := data.GetOk("cache"); ok {
		cache := cache.([]interface{})
		cacheData := cache[0].(map[string]interface{})

		evictionDay := cacheData["eviction_day"].(int)
		evictionHour := cacheData["eviction_hour"].(int)
		evictionMinute := cacheData["eviction_minute"].(int)

		custom.MaxLifespan = cacheData["max_lifespan"].(string)

		custom.EvictionDay = &evictionDay
		custom.EvictionHour = &evictionHour
		custom.EvictionMinute = &evictionMinute
		custom.CachePolicy = cacheData["policy"].(string)
	}

	return custom
}

func setCustomUserFederationData(data *schema.ResourceData, custom *keycloak.CustomUserFederation) {
	data.SetId(custom.Id)

	data.Set("name", custom.Name)
	data.Set("realm_id", custom.RealmId)

	data.Set("parent_id", custom.ParentId)
	data.Set("provider_id", custom.ProviderId)

	data.Set("enabled", custom.Enabled)
	data.Set("priority", custom.Priority)

	if _, ok := data.GetOk("cache"); ok {
		cachePolicySettings := make(map[string]interface{})

		if custom.EvictionDay != nil {
			cachePolicySettings["eviction_day"] = *custom.EvictionDay
		}
		if custom.EvictionHour != nil {
			cachePolicySettings["eviction_hour"] = *custom.EvictionHour
		}
		if custom.EvictionMinute != nil {
			cachePolicySettings["eviction_minute"] = *custom.EvictionMinute
		}
		if custom.MaxLifespan != "" {
			cachePolicySettings["max_lifespan"] = custom.MaxLifespan
		}
		cachePolicySettings["policy"] = custom.CachePolicy

		data.Set("cache", []interface{}{cachePolicySettings})
	} else {
		data.Set("cache_policy", custom.CachePolicy)
	}

	config := make(map[string]interface{})
	for k, v := range custom.Config {
		config[k] = v[0]
	}

	data.Set("config", config)
}

func resourceKeycloakCustomUserFederationCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	custom := getCustomUserFederationFromData(data)

	err := keycloakClient.ValidateCustomUserFederation(custom)
	if err != nil {
		return err
	}

	err = keycloakClient.NewCustomUserFederation(custom)
	if err != nil {
		return err
	}

	setCustomUserFederationData(data, custom)

	return resourceKeycloakCustomUserFederationRead(data, meta)
}

func resourceKeycloakCustomUserFederationRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	custom, err := keycloakClient.GetCustomUserFederation(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setCustomUserFederationData(data, custom)

	return nil
}

func resourceKeycloakCustomUserFederationUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	custom := getCustomUserFederationFromData(data)

	err := keycloakClient.ValidateCustomUserFederation(custom)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateCustomUserFederation(custom)
	if err != nil {
		return err
	}

	setCustomUserFederationData(data, custom)

	return nil
}

func resourceKeycloakCustomUserFederationDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteCustomUserFederation(realmId, id)
}

func resourceKeycloakCustomUserFederationImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
