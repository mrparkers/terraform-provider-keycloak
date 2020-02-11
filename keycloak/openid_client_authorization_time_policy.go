package keycloak

import (
	"encoding/json"
	"fmt"
)

type OpenidClientAuthorizationTimePolicy struct {
	Id               string   `json:"id,omitempty"`
	RealmId          string   `json:"-"`
	ResourceServerId string   `json:"-"`
	Name             string   `json:"name"`
	Owner            string   `json:"owner"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Logic            string   `json:"logic"`
	Policies         []string `json:"policies"`
	Resources        []string `json:"resources"`
	Scopes           []string `json:"scopes"`
	Type             string   `json:"type"`
	NotBefore        string   `json:"string"`
	NotOnOrAfter     string   `json:"string"`
	DayMonth         string   `json:"string"`
	DayMonthEnd      string   `json:"string"`
	Month            string   `json:"string"`
	MonthEnd         string   `json:"string"`
	Year             string   `json:"string"`
	YearEnd          string   `json:"string"`
	Hour             string   `json:"string"`
	HourEnd          string   `json:"string"`
	Minute           string   `json:"string"`
	MinuteEnd        string   `json:"string"`
	Description      string   `json:"description"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationTimePolicy(policy *OpenidClientAuthorizationTimePolicy) error {
	body, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/time", policy.RealmId, policy.ResourceServerId), policy)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationTimePolicy(policy *OpenidClientAuthorizationTimePolicy) error {
	err := keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/time/%s", policy.RealmId, policy.ResourceServerId, policy.Id), policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationTimePolicy(realmId, resourceServerId, policyId string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/time/%s", realmId, resourceServerId, policyId), nil)
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationTimePolicy(realmId, resourceServerId, policyId string) (*OpenidClientAuthorizationTimePolicy, error) {

	policy := OpenidClientAuthorizationTimePolicy{
		Id:               policyId,
		ResourceServerId: resourceServerId,
		RealmId:          realmId,
	}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/time/%s", realmId, resourceServerId, policyId), &policy, nil)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}
