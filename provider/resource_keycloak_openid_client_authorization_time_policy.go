package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientAuthorizationTimePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationTimePolicyCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationTimePolicyRead,
		Delete: resourceKeycloakOpenidClientAuthorizationTimePolicyDelete,
		Update: resourceKeycloakOpenidClientAuthorizationTimePolicyUpdate,
		Importer: &schema.ResourceImporter{
			State: genericResourcePolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"resource_server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"decision_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logic": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: logicKeyValidation,
			},
			"policies": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"not_before": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"not_on_or_after": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"day_month": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"day_month_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"month": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"month_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"year": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"year_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hour": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hour_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"minute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"minute_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getOpenidClientAuthorizationTimePolicyResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationTimePolicy {
	var policies []string
	var resources []string
	var scopes []string
	if v, ok := data.GetOk("resources"); ok {
		for _, resource := range v.(*schema.Set).List() {
			resources = append(resources, resource.(string))
		}
	}
	if v, ok := data.GetOk("policies"); ok {
		for _, policy := range v.(*schema.Set).List() {
			policies = append(policies, policy.(string))
		}
	}
	if v, ok := data.GetOk("scopes"); ok {
		for _, scope := range v.(*schema.Set).List() {
			scopes = append(scopes, scope.(string))
		}
	}

	resource := keycloak.OpenidClientAuthorizationTimePolicy{
		Id:               data.Id(),
		ResourceServerId: data.Get("resource_server_id").(string),
		RealmId:          data.Get("realm_id").(string),
		Owner:            data.Get("owner").(string),
		DecisionStrategy: data.Get("decision_strategy").(string),
		Logic:            data.Get("logic").(string),
		Name:             data.Get("name").(string),
		Type:             "time",
		Policies:         policies,
		Resources:        resources,
		Scopes:           scopes,
		NotBefore:        data.Get("not_before").(string),
		NotOnOrAfter:     data.Get("not_on_or_after").(string),
		DayMonth:         data.Get("day_month").(string),
		DayMonthEnd:      data.Get("day_month_end").(string),
		Month:            data.Get("month").(string),
		MonthEnd:         data.Get("month_end").(string),
		Year:             data.Get("year").(string),
		YearEnd:          data.Get("year_end").(string),
		Hour:             data.Get("hour").(string),
		HourEnd:          data.Get("hour_end").(string),
		Minute:           data.Get("minute").(string),
		MinuteEnd:        data.Get("minute_end").(string),
		Description:      data.Get("description").(string),
	}
	return &resource
}

func setOpenidClientAuthorizationTimePolicyResourceData(data *schema.ResourceData, policy *keycloak.OpenidClientAuthorizationTimePolicy) {
	data.SetId(policy.Id)

	data.Set("resource_server_id", policy.ResourceServerId)
	data.Set("realm_id", policy.RealmId)
	data.Set("name", policy.Name)
	data.Set("decision_strategy", policy.DecisionStrategy)
	data.Set("owner", policy.Owner)
	data.Set("logic", policy.Logic)
	data.Set("policies", policy.Policies)
	data.Set("resources", policy.Resources)
	data.Set("scopes", policy.Scopes)
	data.Set("type", policy.Type)
	data.Set("description", policy.Description)
	data.Set("not_on_or_after", policy.NotOnOrAfter)
	data.Set("not_before", policy.NotBefore)
	data.Set("day_month", policy.DayMonth)
	data.Set("day_month_end", policy.DayMonthEnd)
	data.Set("month", policy.Month)
	data.Set("month_end", policy.MonthEnd)
	data.Set("year", policy.Year)
	data.Set("year_end", policy.YearEnd)
	data.Set("hour", policy.Hour)
	data.Set("hour_end", policy.HourEnd)
	data.Set("minute", policy.Minute)
	data.Set("minute_end", policy.MinuteEnd)
}

func resourceKeycloakOpenidClientAuthorizationTimePolicyCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationTimePolicyResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationTimePolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationTimePolicyResourceData(data, resource)

	return resourceKeycloakOpenidClientAuthorizationTimePolicyRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationTimePolicyRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetOpenidClientAuthorizationTimePolicy(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationTimePolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationTimePolicyUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationTimePolicyResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationTimePolicy(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationTimePolicyResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationTimePolicyDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationTimePolicy(realmId, resourceServerId, id)
}

func resourceKeycloakOpenidClientAuthorizationTimePolicyImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{authorizationResourceId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
