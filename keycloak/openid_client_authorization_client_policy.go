package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
)

type OpenidClientAuthorizationClientPolicy struct {
	Id               string   `json:"id,omitempty"`
	RealmId          string   `json:"-"`
	ResourceServerId string   `json:"-"`
	Name             string   `json:"name"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Logic            string   `json:"logic"`
	Type             string   `json:"type"`
	Clients          []string `json:"clients"`
	Description      string   `json:"description"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationClientPolicy(ctx context.Context, policy *OpenidClientAuthorizationClientPolicy) error {
	body, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/client", policy.RealmId, policy.ResourceServerId), policy)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationClientPolicy(ctx context.Context, policy *OpenidClientAuthorizationClientPolicy) error {
	err := keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/client/%s", policy.RealmId, policy.ResourceServerId, policy.Id), policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationClientPolicy(ctx context.Context, realmId, resourceServerId, policyId string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/client/%s", realmId, resourceServerId, policyId), nil)
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationClientPolicy(ctx context.Context, realmId, resourceServerId, policyId string) (*OpenidClientAuthorizationClientPolicy, error) {

	policy := OpenidClientAuthorizationClientPolicy{
		Id:               policyId,
		ResourceServerId: resourceServerId,
		RealmId:          realmId,
	}
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/client/%s", realmId, resourceServerId, policyId), &policy, nil)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}
