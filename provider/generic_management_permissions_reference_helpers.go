package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func genericManagementPermissionsReferenceImport(data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm}}/clients/{{client_id}}/management/permissions.")
	}

	data.SetId(data.Id())
	data.Set("realm_id", parts[0])
	data.Set("client_id", parts[2])

	return []*schema.ResourceData{data}, nil
}
