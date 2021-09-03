package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getExtraConfigFromData(data *schema.ResourceData) map[string]interface{} {
	extraConfig := map[string]interface{}{}
	if v, ok := data.GetOk("extra_config"); ok {
		for key, value := range v.(map[string]interface{}) {
			extraConfig[key] = value
		}

		if data.HasChange("extra_config") && !data.IsNewResource() {
			oldConfig, newConfig := data.GetChange("extra_config")
			newConfigMap := newConfig.(map[string]interface{})

			for oldKey := range oldConfig.(map[string]interface{}) {
				if _, ok := newConfigMap[oldKey]; !ok {
					extraConfig[oldKey] = ""
				}
			}
		}
	}

	return extraConfig
}

func setExtraConfigData(data *schema.ResourceData, extraConfig map[string]interface{}) {
	c := map[string]interface{}{}

	for k, v := range extraConfig {
		if s, ok := v.(string); ok && s == "" {
			continue
		}

		c[k] = v
	}

	data.Set("extra_config", c)
}
