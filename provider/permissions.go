package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func setOpenidClientScopePermissionPolicy(keycloakClient *keycloak.KeycloakClient, realmId, realmManagementClientId, authorizationPermissionId string, scopeDataSet *schema.Set) error {
	var policies []string

	scopePermission := scopeDataSet.List()[0].(map[string]interface{})

	if v, ok := scopePermission["policies"]; ok {
		for _, policy := range v.(*schema.Set).List() {
			policies = append(policies, policy.(string))
		}
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClientId, authorizationPermissionId)
	if err != nil {
		return err
	}

	if v, ok := scopePermission["description"]; ok {
		permission.Description = v.(string)
	}

	if v, ok := scopePermission["decision_strategy"]; ok {
		permission.DecisionStrategy = v.(string)
	}

	permission.Policies = policies

	return keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)
}

func getOpenidClientScopePermissionPolicy(keycloakClient *keycloak.KeycloakClient, realmId string, realmManagementClientId, permissionId string) (map[string]interface{}, error) {
	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClientId, permissionId)
	if err != nil {
		return nil, err
	}

	if permission.Description == "" && permission.DecisionStrategy == "UNANIMOUS" && len(permission.Policies) == 0 {
		return nil, nil
	}

	permissionViewSettings := make(map[string]interface{})

	if permission.Description != "" {
		permissionViewSettings["description"] = permission.Description
	}

	if permission.DecisionStrategy != "" {
		permissionViewSettings["decision_strategy"] = permission.DecisionStrategy
	}

	if len(permission.Policies) > 0 {
		permissionViewSettings["policies"] = permission.Policies
	}

	return permissionViewSettings, nil
}

func scopePermissionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policies": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"decision_strategy": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(keycloakOpenidClientResourcePermissionDecisionStrategies, false),
				},
			},
		},
	}
}
