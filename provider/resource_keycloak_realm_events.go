package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakRealmEvents() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmEventsCreate,
		Read:   resourceKeycloakRealmEventsRead,
		Delete: resourceKeycloakRealmEventsDelete,
		Update: resourceKeycloakRealmEventsUpdate,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"admin_events_details_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"admin_events_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			//"enable_all_event_types": {
			//	Type: schema.TypeBool,
			//	Optional: true,
			//	ConflictsWith: []string{"enabled_event_types"},
			//	ForceNew: false,
			//},
			"enabled_event_types": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				//ConflictsWith: []string{"enable_all_event_types"},
				ForceNew: false,
				//MinItems: 1,
			},
			"events_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"events_expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"events_listeners": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func getRealmEventsConfigFromData(data *schema.ResourceData) *keycloak.RealmEventsConfig {
	enabledEventTypes := make([]string, 0)
	eventsListeners := make([]string, 0)

	if v, ok := data.GetOk("enabled_event_types"); ok {
		for _, enabledEventType := range v.(*schema.Set).List() {
			enabledEventTypes = append(enabledEventTypes, enabledEventType.(string))
		}
	}

	if v, ok := data.GetOk("events_listeners"); ok {
		for _, eventsListener := range v.(*schema.Set).List() {
			eventsListeners = append(eventsListeners, eventsListener.(string))
		}
	}

	realmEventsConfig := &keycloak.RealmEventsConfig{
		AdminEventsDetailsEnabled: data.Get("admin_events_details_enabled").(bool),
		AdminEventsEnabled:        data.Get("admin_events_enabled").(bool),
		EnabledEventTypes:         enabledEventTypes,
		EventsEnabled:             data.Get("events_enabled").(bool),
		EventsExpiration:          data.Get("events_expiration").(int),
		EventsListeners:           eventsListeners,
	}
	//
	//if !data.Get("enable_all_event_types").(bool) {
	//	data.Set("enabled_event_types", realmEventsConfig.EnabledEventTypes)
	//}

	return realmEventsConfig
}

func setRealmEventsConfigData(data *schema.ResourceData, realmEventsConfig *keycloak.RealmEventsConfig) {
	data.Set("admin_events_details_enabled", realmEventsConfig.AdminEventsDetailsEnabled)
	data.Set("admin_events_enabled", realmEventsConfig.AdminEventsEnabled)
	data.Set("events_enabled", realmEventsConfig.EventsEnabled)
	data.Set("events_expiration", realmEventsConfig.EventsExpiration)
	data.Set("events_listeners", realmEventsConfig.EventsListeners)

	if _, ok := data.GetOk("enabled_event_types"); ok {
		data.Set("enabled_event_types", realmEventsConfig.EnabledEventTypes)
	}
}

func resourceKeycloakRealmEventsCreate(data *schema.ResourceData, meta interface{}) error {
	realmId := data.Get("realm_id").(string)
	data.SetId(realmId)

	err := resourceKeycloakRealmEventsUpdate(data, meta)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmEventsRead(data, meta)
}

func resourceKeycloakRealmEventsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)

	realmEventsConfig, err := keycloakClient.GetRealmEventsConfig(realmId)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setRealmEventsConfigData(data, realmEventsConfig)

	return nil
}

func resourceKeycloakRealmEventsDelete(data *schema.ResourceData, meta interface{}) error {
	// TODO: Do we want to do nothing here since the realm events config cannot be deleted? Or do we want to set all the zero values?
	// Note the zero values are different than what keycloak's defaults are when a realm is created.

	return nil
}

func resourceKeycloakRealmEventsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	realmEventsConfig := getRealmEventsConfigFromData(data)

	err := keycloakClient.UpdateRealmEventsConfig(realmId, realmEventsConfig)
	if err != nil {
		return err
	}

	setRealmEventsConfigData(data, realmEventsConfig)

	return nil
}
