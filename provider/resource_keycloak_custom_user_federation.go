package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakCustomUserFederation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakCustomUserFederationCreate,
		ReadContext:   resourceKeycloakCustomUserFederationRead,
		UpdateContext: resourceKeycloakCustomUserFederationUpdate,
		DeleteContext: resourceKeycloakCustomUserFederationDelete,
		// This resource can be imported using {{realm}}/{{provider_id}}. The Provider ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakCustomUserFederationImport,
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice(keycloakUserFederationCachePolicies, false),
			},

			"full_sync_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validateSyncPeriod,
				Description:  "How frequently Keycloak should sync all users, in seconds. Omit this property to disable periodic full sync.",
			},
			"changed_sync_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validateSyncPeriod,
				Description:  "How frequently Keycloak should sync changed users, in seconds. Omit this property to disable periodic changed users sync.",
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

	return &keycloak.CustomUserFederation{
		Id:         data.Id(),
		Name:       data.Get("name").(string),
		RealmId:    data.Get("realm_id").(string),
		ParentId:   parentId,
		ProviderId: data.Get("provider_id").(string),

		Enabled:  data.Get("enabled").(bool),
		Priority: data.Get("priority").(int),

		CachePolicy: data.Get("cache_policy").(string),

		FullSyncPeriod:    data.Get("full_sync_period").(int),
		ChangedSyncPeriod: data.Get("changed_sync_period").(int),

		Config: config,
	}
}

func setCustomUserFederationData(data *schema.ResourceData, custom *keycloak.CustomUserFederation) {
	data.SetId(custom.Id)

	data.Set("name", custom.Name)
	data.Set("realm_id", custom.RealmId)

	data.Set("parent_id", custom.ParentId)
	data.Set("provider_id", custom.ProviderId)

	data.Set("enabled", custom.Enabled)
	data.Set("priority", custom.Priority)

	data.Set("full_sync_period", custom.FullSyncPeriod)
	data.Set("changed_sync_period", custom.ChangedSyncPeriod)

	data.Set("cache_policy", custom.CachePolicy)

	config := make(map[string]interface{})
	for k, v := range custom.Config {
		config[k] = v[0]
	}

	data.Set("config", config)
}

func resourceKeycloakCustomUserFederationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	custom := getCustomUserFederationFromData(data)

	err := keycloakClient.ValidateCustomUserFederation(ctx, custom)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.NewCustomUserFederation(ctx, custom)
	if err != nil {
		return diag.FromErr(err)
	}

	setCustomUserFederationData(data, custom)

	return resourceKeycloakCustomUserFederationRead(ctx, data, meta)
}

func resourceKeycloakCustomUserFederationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	custom, err := keycloakClient.GetCustomUserFederation(ctx, realmId, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	setCustomUserFederationData(data, custom)

	return nil
}

func resourceKeycloakCustomUserFederationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	custom := getCustomUserFederationFromData(data)

	err := keycloakClient.ValidateCustomUserFederation(ctx, custom)
	if err != nil {
		return diag.FromErr(err)
	}

	err = keycloakClient.UpdateCustomUserFederation(ctx, custom)
	if err != nil {
		return diag.FromErr(err)
	}

	setCustomUserFederationData(data, custom)

	return nil
}

func resourceKeycloakCustomUserFederationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return diag.FromErr(keycloakClient.DeleteCustomUserFederation(ctx, realmId, id))
}

func resourceKeycloakCustomUserFederationImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
