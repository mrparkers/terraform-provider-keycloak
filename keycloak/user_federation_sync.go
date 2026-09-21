package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	UserFederationSyncFull         = "triggerFullSync"
	UserFederationSyncChangedUsers = "triggerChangedUsersSync"
)

type UserFederationSyncResult struct {
	Ignored bool   `json:"ignored"`
	Added   int    `json:"added"`
	Updated int    `json:"updated"`
	Removed int    `json:"removed"`
	Failed  int    `json:"failed"`
	Status  string `json:"status"`
}

// SyncUserFederation triggers a synchronization of users for the given user storage provider (e.g. an LDAP user federation).
// The action must be one of UserFederationSyncFull or UserFederationSyncChangedUsers.
func (keycloakClient *KeycloakClient) SyncUserFederation(ctx context.Context, realmId, id, action string) (*UserFederationSyncResult, error) {
	path := fmt.Sprintf("/realms/%s/user-storage/%s/sync?action=%s", realmId, id, url.QueryEscape(action))

	body, _, err := keycloakClient.post(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result UserFederationSyncResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
