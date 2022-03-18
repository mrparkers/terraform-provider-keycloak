package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
)

type OpenidClientAuthorizationGroupPolicy struct {
	Id               string                           `json:"id,omitempty"`
	RealmId          string                           `json:"-"`
	ResourceServerId string                           `json:"-"`
	Name             string                           `json:"name"`
	DecisionStrategy string                           `json:"decisionStrategy"`
	Logic            string                           `json:"logic"`
	Type             string                           `json:"type"`
	GroupsClaim      string                           `json:"groupsClaim"`
	Groups           []OpenidClientAuthorizationGroup `json:"groups,omitempty"`
	Description      string                           `json:"description"`
}

type OpenidClientAuthorizationGroup struct {
	Id             string `json:"id,omitempty"`
	Path           string `json:"path,omitempty"`
	ExtendChildren bool   `json:"extendChildren,omitempty"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationGroupPolicy(ctx context.Context, policy *OpenidClientAuthorizationGroupPolicy) error {
	body, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/group", policy.RealmId, policy.ResourceServerId), policy)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationGroupPolicy(ctx context.Context, policy *OpenidClientAuthorizationGroupPolicy) error {
	err := keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/group/%s", policy.RealmId, policy.ResourceServerId, policy.Id), policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationGroupPolicy(ctx context.Context, realmId, resourceServerId, policyId string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/group/%s", realmId, resourceServerId, policyId), nil)
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationGroupPolicy(ctx context.Context, realmId, resourceServerId, policyId string) (*OpenidClientAuthorizationGroupPolicy, error) {

	policy := OpenidClientAuthorizationGroupPolicy{
		Id:               policyId,
		ResourceServerId: resourceServerId,
		RealmId:          realmId,
	}
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/group/%s", realmId, resourceServerId, policyId), &policy, nil)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}
