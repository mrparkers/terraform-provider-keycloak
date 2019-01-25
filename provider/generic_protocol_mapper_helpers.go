package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func genericProtocolMapperImport(data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")

	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid import. supported import formats: {{realmId}}/client/{{clientId}}/{{protocolMapperId}} or {{realmId}}/client-scope/{{clientScopeId}}/{{protocolMapperId}}")
	}

	realmId := parts[0]
	parentResourceType := parts[1]
	parentResourceId := parts[2]
	mapperId := parts[3]

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
