package provider

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func keys(data map[string]string) []string {
	var result = []string{}
	for k := range data {
		result = append(result, k)
	}
	return result
}

func mergeSchemas(a map[string]*schema.Schema, b map[string]*schema.Schema) map[string]*schema.Schema {
	result := a
	for k, v := range b {
		result[k] = v
	}
	return result
}

// Converts duration string to an int representing the number of seconds, which is used by the Keycloak API
// Ex: "1h" => 3600
func getSecondsFromDurationString(s string) (int, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}

	return int(duration.Seconds()), nil
}

// Converts number of seconds from Keycloak API to a duration string used by the provider
// Ex: 3600 => "1h0m0s"
func getDurationStringFromSeconds(seconds int) string {
	return (time.Duration(seconds) * time.Second).String()
}

// This will suppress the Terraform diff when comparing duration strings.
// As long as both strings represent the same number of seconds, it makes no difference to the Keycloak API
func suppressDurationStringDiff(_, old, new string, _ *schema.ResourceData) bool {
	if old == "" || new == "" {
		return false
	}

	oldDuration, _ := time.ParseDuration(old)
	newDuration, _ := time.ParseDuration(new)

	return oldDuration.Seconds() == newDuration.Seconds()
}

func handleNotFoundError(err error, data *schema.ResourceData) error {
	if keycloak.ErrorIs404(err) {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer exists", data.Id())
		data.SetId("")

		return nil
	}

	return err
}

func interfaceSliceToStringSlice(iv []interface{}) []string {
	var sv []string
	for _, i := range iv {
		sv = append(sv, i.(string))
	}

	return sv
}

func stringArrayDifference(a, b []string) []string {
	var aWithoutB []string

	for _, s := range a {
		if !stringSliceContains(b, s) {
			aWithoutB = append(aWithoutB, s)
		}
	}

	return aWithoutB
}

func stringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
