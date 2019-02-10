package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"strings"
)

func genericProtocolMapperImport(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")
	keycloakClient := meta.(*keycloak.KeycloakClient)
	var realmId, mapperId, parentResourceType, parentResourceId string

	switch {
	case len(parts) == 3 && keycloakClient.GetDefaultRealm() != "":
		realmId = keycloakClient.GetDefaultRealm()
		parentResourceType = parts[0]
		parentResourceId = parts[1]
		mapperId = parts[2]
	case len(parts) == 4:
		realmId = parts[0]
		parentResourceType = parts[1]
		parentResourceId = parts[2]
		mapperId = parts[3]
	default:
		return nil, fmt.Errorf("invalid import. supported import formats: {{realmId}}/client/{{clientId}}/{{protocolMapperId}}, {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}} or client/{{clientId}}/{{protocolMapperId}}, client-scope/{{clientScopeId}}/{{protocolMapperId}} when default realm is set")
	}

	data.Set("realm_id", realmId)
	data.SetId(mapperId)

	if parentResourceType == "client" {
		data.Set("client_id", parentResourceId)
	} else if parentResourceType == "client-scope" {
		data.Set("client_scope_id", parentResourceId)
	} else {
		return nil, fmt.Errorf("the associated parent resource must be either a client or a client-scope")
	}

	return []*schema.ResourceData{data}, nil
}
