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

func (keycloakClient *KeycloakClient) NewAuthenticationExecution(realmId, parentAlias, provider string) (*AuthenticationExecution, error) {
	location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/authentication/flows/%s/executions/execution", realmId, parentAlias), &authenticationExecutionCreate{Provider: provider})
	if err != nil {
		return nil, err
	}

	newExecutionId := getIdFromLocationHeader(location)

	authenticationExecutions, err := keycloakClient.ListAuthenticationExecutions(realmId, parentAlias)
	if err != nil {
		return nil, err
	}

	for _, execution := range authenticationExecutions {
		if execution.Id == newExecutionId {
			return execution, nil
		}
	}

	return nil, fmt.Errorf("unable to find newly created execution with id %s (this should never happen)", newExecutionId)
}

func (keycloakClient *KeycloakClient) ListAuthenticationExecutions(realmId, parentAlias string) (AuthenticationExecutionList, error) {
	var authenticationExecutions []*AuthenticationExecution

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows/%s/executions", realmId, parentAlias), &authenticationExecutions)
	if err != nil {
		return nil, err
	}

	return authenticationExecutions, err
}
