package keycloak

import (
	"context"
	"fmt"
)

type Organization struct {
	Id          string               `json:"id,omitempty"`
	Realm       string               `json:"-"`
	Name        string               `json:"name"`
	Alias       string               `json:"alias,omitempty"`
	Enabled     bool                 `json:"enabled"`
	Description string               `json:"description,omitempty"`
	RedirectUrl string               `json:"redirectUrl,omitempty"`
	Domains     []OrganizationDomain `json:"domains"`
	Attributes  map[string][]string  `json:"attributes"`
}

type OrganizationDomain struct {
	Name                  string `json:"name"`
	Verified              bool   `json:"verified,omitempty"`
	IdentityProviderAlias string `json:"identityProviderAlias,omitempty"`
	AutoRedirect          bool   `json:"autoRedirect,omitempty"`
}

type OrganizationIdentityProvider struct {
	Alias          string
	AutoMembership bool
	MembershipType string
}

type organizationIdentityProviderRepresentation struct {
	Alias             string                             `json:"alias"`
	OrganizationLinks []IdentityProviderOrganizationLink `json:"organizationLinks"`
}

func (keycloakClient *KeycloakClient) NewOrganization(ctx context.Context, organization *Organization) error {
	_, location, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/organizations", organization.Realm), organization)
	if err != nil {
		return err
	}
	organization.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetOrganization(ctx context.Context, realm, id string) (*Organization, error) {
	var organization Organization

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/organizations/%s", realm, id), &organization, nil)
	if err != nil {
		return nil, err
	}

	organization.Realm = realm

	return &organization, nil
}

func (keycloakClient *KeycloakClient) GetOrganizationByName(ctx context.Context, realm string, name string) (*Organization, error) {
	var organizations []Organization
	var orgFound *Organization

	params := map[string]string{
		"search": name,
	}

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/organizations", realm), &organizations, params)
	if err != nil {
		return nil, err
	}

	for _, org := range organizations {
		if org.Name == name {
			if orgFound, err = keycloakClient.GetOrganization(ctx, realm, org.Id); err != nil {
				return nil, err
			}
			return orgFound, nil
		}
	}

	return nil, fmt.Errorf("organization with name %s not found", name)
}

func (keycloakClient *KeycloakClient) UpdateOrganization(ctx context.Context, organization *Organization) error {
	modern, err := keycloakClient.VersionIsGreaterThanOrEqualTo(ctx, Version_26_8)
	if err != nil {
		return err
	}
	if modern {
		return keycloakClient.updateOrganizationPreservingIdentityProviderRouting(ctx, organization)
	}
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/organizations/%s", organization.Realm, organization.Id), organization)
}

// UpdateOrganizationWithDomainRouting writes the identity provider routing of each domain as given,
// where UpdateOrganization keeps the routing that the server has. Requires Keycloak 26.8 or later.
func (keycloakClient *KeycloakClient) UpdateOrganizationWithDomainRouting(ctx context.Context, organization *Organization) error {
	identityProviderOrganizationMutex.Lock()
	defer identityProviderOrganizationMutex.Unlock()
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/organizations/%s", organization.Realm, organization.Id), organization)
}

// GetOrganizationIdentityProvider reads one link with its own request. The list of an organization's
// identity providers is not used: Keycloak 26.8.0 writes that list after its transaction ends, and
// fails with "Cannot parse the JSON" when another organization linked to the same identity provider
// is not cached.
func (keycloakClient *KeycloakClient) GetOrganizationIdentityProvider(ctx context.Context, realm, organizationId, alias string) (*OrganizationIdentityProvider, error) {
	var representation organizationIdentityProviderRepresentation

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/organizations/%s/identity-providers/%s", realm, organizationId, alias), &representation, nil)
	if err != nil {
		return nil, err
	}

	for _, link := range representation.OrganizationLinks {
		if link.OrganizationId == organizationId {
			return &OrganizationIdentityProvider{
				Alias:          representation.Alias,
				AutoMembership: link.AutoMembership,
				MembershipType: link.MembershipType,
			}, nil
		}
	}

	return nil, fmt.Errorf("identity provider %s of organization %s has no link to the organization", alias, organizationId)
}

func (keycloakClient *KeycloakClient) UpdateOrganizationIdentityProviderLink(ctx context.Context, realm, organizationId string, identityProvider OrganizationIdentityProvider) error {
	link := IdentityProviderOrganizationLink{
		OrganizationId: organizationId,
		AutoMembership: identityProvider.AutoMembership,
		MembershipType: identityProvider.MembershipType,
	}

	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/organizations/%s/identity-providers/%s", realm, organizationId, identityProvider.Alias), link)
}

func (keycloakClient *KeycloakClient) DeleteOrganization(ctx context.Context, realmId, id string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/organizations/%s", realmId, id), nil)
}

type OrganizationMember struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	MembershipType string `json:"membershipType"`
}

func (keycloakClient *KeycloakClient) GetOrganizationMembers(ctx context.Context, realm, orgId string) ([]*OrganizationMember, error) {
	var members []*OrganizationMember
	var first, pagination = 0, 50
	var iterationMembers []*OrganizationMember

	for ok := true; ok; ok = len(iterationMembers) > 0 {
		iterationMembers = nil
		err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/organizations/%s/members?max=%d&first=%d", realm, orgId, pagination, first), &iterationMembers, nil)
		if err != nil {
			return nil, err
		}
		members = append(members, iterationMembers...)
		first += pagination
	}

	return members, nil
}

func (keycloakClient *KeycloakClient) AddUserToOrganization(ctx context.Context, realm, orgId, userId string) error {
	_, err := keycloakClient.sendRaw(ctx, fmt.Sprintf("/realms/%s/organizations/%s/members", realm, orgId), []byte(userId))
	return err
}

func (keycloakClient *KeycloakClient) RemoveUserFromOrganization(ctx context.Context, realm, orgId, userId string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/organizations/%s/members/%s", realm, orgId, userId), nil)
}
