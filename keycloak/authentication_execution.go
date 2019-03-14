package keycloak

import (
	"fmt"
)

// this is only used when creating an execution on a flow.
// other fields can be provided to the API but they are ignored
// POST /realms/${realmId}/authentication/flows/${flowAlias}/executions/execution
type authenticationExecutionCreate struct {
	Provider string `json:"provider"`
}

// this type is returned by GET /realms/${realmId}/authentication/flows/${flowAlias}/executions
// another model is used for GET /realms/${realmId}/authentication/executions/${executionId}, but I am going to try to avoid using this API
type AuthenticationExecution struct {
	Id                 string   `json:"id"`
	RealmId            string   `json:"-"`
	ParentFlowAlias    string   `json:"-"`
	Provider           string   `json:"providerId"`
	Requirement        string   `json:"requirement"`
	RequirementChoices []string `json:"requirementChoices,omitempty"`
	Index              int      `json:"index,omitempty"`
	Configurable       bool     `json:"configurable,omitempty"`
}

type AuthenticationExecutionList []*AuthenticationExecution

func (list AuthenticationExecutionList) Len() int {
	return len(list)
}

func (list AuthenticationExecutionList) Less(i, j int) bool {
	return list[i].Index < list[j].Index
}

func (list AuthenticationExecutionList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (keycloakClient *KeycloakClient) NewAuthenticationExecution(execution *AuthenticationExecution) error {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/authentication/flows/%s/executions/execution", execution.RealmId, execution.ParentFlowAlias), &authenticationExecutionCreate{Provider: execution.Provider})
	if err != nil {
		return err
	}

	execution.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) ListAuthenticationExecutions(realmId, parentAlias string) (AuthenticationExecutionList, error) {
	var authenticationExecutions []*AuthenticationExecution

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows/%s/executions", realmId, parentAlias), &authenticationExecutions)
	if err != nil {
		return nil, err
	}

	return authenticationExecutions, err
}

// note: only the "requirement" field can be updated this way
func (keycloakClient *KeycloakClient) UpdateAuthenticationExecution(execution *AuthenticationExecution) error {
	err := keycloakClient.put(fmt.Sprintf("/realms/%s/authentication/flows/%s/executions", execution.RealmId, execution.ParentFlowAlias), execution)
	if err != nil {
		return err
	}

	return nil
}
