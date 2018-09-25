package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

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
