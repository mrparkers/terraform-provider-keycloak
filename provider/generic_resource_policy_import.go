package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func genericResourcePolicyImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{authorizationResourceId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
