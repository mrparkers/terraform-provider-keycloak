package keycloak

import (
	"encoding/json"
	"fmt"
)

type OpenidClientRole struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	ScopeParamRequired bool   `json:"scopeParamRequired"`
	ClientRole         bool   `json:"clientRole"`
	ContainerId        string `json:"ContainerId"`
}

type OpenidClientSecret struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type OpenidClientAuthorizationResource struct {
	ClientId           string                           `json:"-"`
	RealmId            string                           `json:"-"`
	Id                 string                           `json:"_id,omitempty"`
	DisplayName        string                           `json:"displayName"`
	Name               string                           `json:"name"`
	Uris               []string                         `json:"uris"`
	IconUri            string                           `json:"icon_uri"`
	OwnerManagedAccess bool                             `json:"ownerManagedAccess"`
	Scopes             []OpenidClientAuthorizationScope `json:"scopes"`
	Type               string                           `json:"type"`
	Attributes         map[string][]string              `json:"attributes"`
}

type OpenidClientAuthorizationScope struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	ClientId    string `json:"-"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IconUri     string `json:"iconUri"`
}

type OpenidClientAuthorizationPermission struct {
	Id               string   `json:"id,omitempty"`
	RealmId          string   `json:"-"`
	ClientId         string   `json:"-"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Policies         []string `json:"policies"`
	Resources        []string `json:"resources"`
	Type             string   `json:"type"`
}

type OpenidClientAuthorizationPolicy struct {
	Id               string   `json:"id,omitempty"`
	RealmId          string   `json:"-"`
	ClientId         string   `json:"-"`
	Name             string   `json:"name"`
	Owner            string   `json:"owner"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Logic            string   `json:"logic"`
	Policies         []string `json:"policies"`
	Resources        []string `json:"resources"`
	Scopes           []string `json:"scopes"`
	Type             string   `json:"type"`
}

type OpenidClientServiceAccountRole struct {
	Id                   string `json:"id"`
	RealmId              string `json:"-"`
	ServiceAccountUserId string `json:"-"`
	Name                 string `json:"name"`
	ClientRole           bool   `json:"clientRole"`
	Composite            bool   `json:"composite"`
	ContainerId          string `json:"containerId"`
}

type OpenidClientAuthorizationSettings struct {
	PolicyEnforcementMode         string `json:"policyEnforcementMode,omitempty"`
	AllowRemoteResourceManagement bool   `json:"allowRemoteResourceManagement,omitempty"`
}

type OpenidClient struct {
	Id                           string                             `json:"id,omitempty"`
	ClientId                     string                             `json:"clientId"`
	RealmId                      string                             `json:"-"`
	Name                         string                             `json:"name"`
	Protocol                     string                             `json:"protocol"`                // always openid-connect for this resource
	ClientAuthenticatorType      string                             `json:"clientAuthenticatorType"` // always client-secret for now, don't have a need for JWT here
	ClientSecret                 string                             `json:"secret,omitempty"`
	Enabled                      bool                               `json:"enabled"`
	Description                  string                             `json:"description"`
	PublicClient                 bool                               `json:"publicClient"`
	BearerOnly                   bool                               `json:"bearerOnly"`
	StandardFlowEnabled          bool                               `json:"standardFlowEnabled"`
	ImplicitFlowEnabled          bool                               `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled    bool                               `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled       bool                               `json:"serviceAccountsEnabled"`
	AuthorizationServicesEnabled bool                               `json:"authorizationServicesEnabled"`
	ValidRedirectUris            []string                           `json:"redirectUris"`
	WebOrigins                   []string                           `json:"webOrigins"`
	AuthorizationSettings        *OpenidClientAuthorizationSettings `json:"authorizationSettings,omitempty"`
}

func (keycloakClient *KeycloakClient) GetClientAuthorizationPolicyByName(realmId, clientId, name string) (*OpenidClientAuthorizationPolicy, error) {
	policies := []OpenidClientAuthorizationPolicy{}
	params := map[string]string{"name": name}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy", realmId, clientId), &policies, params)
	if err != nil {
		return nil, err
	}
	policy := policies[0]
	policy.RealmId = realmId
	policy.ClientId = clientId
	policy.Name = name
	return &policy, nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationPermission(realm, clientId, id string) (*OpenidClientAuthorizationPermission, error) {
	permission := OpenidClientAuthorizationPermission{
		RealmId:  realm,
		ClientId: clientId,
		Id:       id,
	}

	policies := []OpenidClientAuthorizationPolicy{}
	resources := []OpenidClientAuthorizationResource{}

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/permission/resource/%s", realm, clientId, id), &permission, nil)
	if err != nil {
		return nil, err
	}

	err = keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/%s/associatedPolicies", realm, clientId, id), &policies, nil)
	if err != nil {
		return nil, err
	}

	err = keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/permission/%s/resources", realm, clientId, id), &resources, nil)
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		permission.Policies = append(permission.Policies, policy.Id)
	}

	for _, resource := range resources {
		permission.Resources = append(permission.Resources, resource.Id)
	}

	return &permission, nil
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationPermission(permission *OpenidClientAuthorizationPermission) error {
	body, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/permission", permission.RealmId, permission.ClientId), permission)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &permission)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationPermission(permission *OpenidClientAuthorizationPermission) error {
	err := keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/permission/resource/%s", permission.RealmId, permission.ClientId, permission.Id), permission)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationPermission(realmId, clientId, permissionId string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/permission/%s", realmId, clientId, permissionId), nil)
}

func (keycloakClient *KeycloakClient) GetClientRoleByName(realm, clientId, name string) (*OpenidClientRole, error) {
	var clientRole OpenidClientRole
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/roles/%s", realm, clientId, name), &clientRole, nil)
	if err != nil {
		return nil, err
	}
	return &clientRole, nil
}

func (keycloakClient *KeycloakClient) GetClientByName(realm, clientId string) (*OpenidClient, error) {
	var clients []OpenidClient
	params := map[string]string{"clientId": clientId}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients", realm), &clients, params)
	if err != nil {
		return nil, err
	}
	return &clients[0], nil
}

func (keycloakClient *KeycloakClient) NewOpenidClientServiceAccountRole(serviceAccountRole *OpenidClientServiceAccountRole) error {
	role, err := keycloakClient.GetClientRoleByName(serviceAccountRole.RealmId, serviceAccountRole.ContainerId, serviceAccountRole.Name)
	if err != nil {
		return err
	}
	serviceAccountRole.Id = role.Id
	serviceAccountRoles := []OpenidClientServiceAccountRole{}
	serviceAccountRoles = append(serviceAccountRoles, *serviceAccountRole)
	_, _, err = keycloakClient.post(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.ContainerId), serviceAccountRoles)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientServiceAccountRole(realm, serviceAccountUserId, clientId, roleId string) error {
	serviceAccountRoles := []OpenidClientServiceAccountRole{}
	serviceAccountRoles = append(serviceAccountRoles, OpenidClientServiceAccountRole{
		Id: roleId,
	})
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, clientId), &serviceAccountRoles)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientServiceAccountRole(realm, serviceAccountUserId, clientId, roleId string) (*OpenidClientServiceAccountRole, error) {
	serviceAccountRoles := []OpenidClientServiceAccountRole{}
	serviceAccountRoles = append(serviceAccountRoles, OpenidClientServiceAccountRole{
		Id:                   roleId,
		RealmId:              realm,
		ContainerId:          clientId,
		ServiceAccountUserId: serviceAccountUserId,
	})
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, clientId), &serviceAccountRoles, nil)
	if err != nil {
		return nil, err
	}
	for _, serviceAccountRole := range serviceAccountRoles {
		if serviceAccountRole.Id == roleId {
			return &serviceAccountRole, nil
		}
	}
	return nil, fmt.Errorf("No role with id %s found", roleId)
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationResource(resource *OpenidClientAuthorizationResource) error {
	body, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/resource", resource.RealmId, resource.ClientId), resource)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationResource(realm, clientId, resourceId string) (*OpenidClientAuthorizationResource, error) {
	resource := OpenidClientAuthorizationResource{
		RealmId:  realm,
		ClientId: clientId,
	}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/resource/%s", realm, clientId, resourceId), &resource, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationResource(resource *OpenidClientAuthorizationResource) error {
	err := keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/resource/%s", resource.RealmId, resource.ClientId, resource.Id), resource)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationResource(realmId, clientId, resourceId string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/resource/%s", realmId, clientId, resourceId), nil)
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationScope(scope *OpenidClientAuthorizationScope) error {
	body, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/scope", scope.RealmId, scope.ClientId), scope)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &scope)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationScope(realm, clientId, scopeId string) (*OpenidClientAuthorizationScope, error) {
	scope := OpenidClientAuthorizationScope{
		RealmId:  realm,
		ClientId: clientId,
	}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/scope/%s", realm, clientId, scopeId), &scope, nil)
	if err != nil {
		return nil, err
	}
	return &scope, nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationScope(scope *OpenidClientAuthorizationScope) error {
	err := keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/scope/%s", scope.RealmId, scope.ClientId, scope.Id), scope)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationScope(realmId, clientId, scopeId string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/scope/%s", realmId, clientId, scopeId), nil)
}

func (keycloakClient *KeycloakClient) GetOpenidClientServiceAccountUserId(realmId, clientId string) (*User, error) {
	var serviceAccountUser User
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/service-account-user", realmId, clientId), &serviceAccountUser, nil)
	if err != nil {
		return &serviceAccountUser, err
	}
	return &serviceAccountUser, nil
}

func (keycloakClient *KeycloakClient) ValidateOpenidClient(client *OpenidClient) error {
	if client.BearerOnly && (client.StandardFlowEnabled || client.ImplicitFlowEnabled || client.DirectAccessGrantsEnabled || client.ServiceAccountsEnabled) {
		return fmt.Errorf("validation error: Keycloak cannot issue tokens for bearer-only clients; no oauth2 flows can be enabled for this client")
	}

	if (client.StandardFlowEnabled || client.ImplicitFlowEnabled) && len(client.ValidRedirectUris) == 0 {
		return fmt.Errorf("validation error: standard (authorization code) and implicit flows require at least one valid redirect uri")
	}

	if client.ServiceAccountsEnabled && client.PublicClient {
		return fmt.Errorf("validation error: service accounts (client credentials flow) cannot be enabled on public clients")
	}

	return nil
}

func (keycloakClient *KeycloakClient) NewOpenidClient(client *OpenidClient) error {
	client.Protocol = "openid-connect"
	client.ClientAuthenticatorType = "client-secret"

	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/clients", client.RealmId), client)
	if err != nil {
		return err
	}

	client.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClient(realmId, id string) (*OpenidClient, error) {
	var client OpenidClient
	var clientSecret OpenidClientSecret

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), &client, nil)
	if err != nil {
		return nil, err
	}

	err = keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/client-secret", realmId, id), &clientSecret, nil)
	if err != nil {
		return nil, err
	}

	client.RealmId = realmId
	client.ClientSecret = clientSecret.Value

	return &client, nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientByClientId(realmId, clientId string) (*OpenidClient, error) {
	var clients []OpenidClient
	var clientSecret OpenidClientSecret

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients?clientId=%s", realmId, clientId), &clients, nil)
	if err != nil {
		return nil, err
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("openid client with name %s does not exist", clientId)
	}

	client := clients[0]

	err = keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/client-secret", realmId, client.Id), &clientSecret, nil)
	if err != nil {
		return nil, err
	}

	client.RealmId = realmId
	client.ClientSecret = clientSecret.Value

	return &client, nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClient(client *OpenidClient) error {
	client.Protocol = "openid-connect"
	client.ClientAuthenticatorType = "client-secret"

	return keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s", client.RealmId, client.Id), client)
}

func (keycloakClient *KeycloakClient) DeleteOpenidClient(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s", realmId, id), nil)
}

func (keycloakClient *KeycloakClient) getOpenidClientScopes(realmId, clientId, t string) ([]*OpenidClientScope, error) {
	var scopes []*OpenidClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/clients/%s/%s-client-scopes", realmId, clientId, t), &scopes, nil)
	if err != nil {
		return nil, err
	}

	return scopes, nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientDefaultScopes(realmId, clientId string) ([]*OpenidClientScope, error) {
	return keycloakClient.getOpenidClientScopes(realmId, clientId, "default")
}

func (keycloakClient *KeycloakClient) GetOpenidClientOptionalScopes(realmId, clientId string) ([]*OpenidClientScope, error) {
	return keycloakClient.getOpenidClientScopes(realmId, clientId, "optional")
}

func (keycloakClient *KeycloakClient) attachOpenidClientScopes(realmId, clientId, t string, scopeNames []string) error {
	openidClient, err := keycloakClient.GetOpenidClient(realmId, clientId)
	if err != nil && ErrorIs404(err) {
		return fmt.Errorf("validation error: client with id %s does not exist", clientId)
	} else if err != nil {
		return err
	}

	if openidClient.BearerOnly {
		return fmt.Errorf("validation error: client with id %s uses access type BEARER-ONLY which does not use scopes", clientId)
	}

	allOpenidClientScopes, err := keycloakClient.listOpenidClientScopesWithFilter(realmId, includeOpenidClientScopesMatchingNames(scopeNames))
	if err != nil {
		return err
	}

	var attachedClientScopes []*OpenidClientScope
	var duplicateScopeAssignmentErrorMessage string
	switch t {
	case "optional":
		attachedDefaultClientScopes, err := keycloakClient.GetOpenidClientDefaultScopes(realmId, clientId)
		if err != nil {
			return err
		}
		attachedClientScopes = append(attachedClientScopes, attachedDefaultClientScopes...)
		duplicateScopeAssignmentErrorMessage = "validation error: scope %s is already attached to client as a default scope"
	case "default":
		attachedOptionalClientScopes, err := keycloakClient.GetOpenidClientOptionalScopes(realmId, clientId)
		if err != nil {
			return err
		}
		attachedClientScopes = append(attachedClientScopes, attachedOptionalClientScopes...)
		duplicateScopeAssignmentErrorMessage = "validation error: scope %s is already attached to client as an optional scope"
	}

	for _, openidClientScope := range allOpenidClientScopes {
		for _, attachedClientScope := range attachedClientScopes {
			if openidClientScope.Id == attachedClientScope.Id {
				return fmt.Errorf(duplicateScopeAssignmentErrorMessage, attachedClientScope.Name)
			}
		}

		err := keycloakClient.put(fmt.Sprintf("/realms/%s/clients/%s/%s-client-scopes/%s", realmId, clientId, t, openidClientScope.Id), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) AttachOpenidClientDefaultScopes(realmId, clientId string, scopeNames []string) error {
	return keycloakClient.attachOpenidClientScopes(realmId, clientId, "default", scopeNames)
}

func (keycloakClient *KeycloakClient) AttachOpenidClientOptionalScopes(realmId, clientId string, scopeNames []string) error {
	return keycloakClient.attachOpenidClientScopes(realmId, clientId, "optional", scopeNames)
}

func (keycloakClient *KeycloakClient) detachOpenidClientScopes(realmId, clientId, t string, scopeNames []string) error {
	allOpenidClientScopes, err := keycloakClient.listOpenidClientScopesWithFilter(realmId, includeOpenidClientScopesMatchingNames(scopeNames))
	if err != nil {
		return err
	}

	for _, openidClientScope := range allOpenidClientScopes {
		err := keycloakClient.delete(fmt.Sprintf("/realms/%s/clients/%s/%s-client-scopes/%s", realmId, clientId, t, openidClientScope.Id), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) DetachOpenidClientDefaultScopes(realmId, clientId string, scopeNames []string) error {
	return keycloakClient.detachOpenidClientScopes(realmId, clientId, "default", scopeNames)
}

func (keycloakClient *KeycloakClient) DetachOpenidClientOptionalScopes(realmId, clientId string, scopeNames []string) error {
	return keycloakClient.detachOpenidClientScopes(realmId, clientId, "optional", scopeNames)
}
