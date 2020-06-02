package provider

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func randomBool() bool {
	return rand.Intn(2) == 0
}

func randomStringInSlice(slice []string) string {
	return slice[acctest.RandIntRange(0, len(slice)-1)]
}

func randomStringSliceSubset(slice []string) []string {
	var result []string

	for _, s := range slice {
		if randomBool() {
			result = append(result, s)
		}
	}

	return result
}

// Returns a slice of strings in the format ["foo", "bar"] for
// use within terraform resource definitions for acceptance tests
func arrayOfStringsForTerraformResource(parts []string) string {
	var tfStrings []string

	for _, part := range parts {
		tfStrings = append(tfStrings, fmt.Sprintf(`"%s"`, part))
	}

	return "[" + strings.Join(tfStrings, ", ") + "]"
}

func randomDurationString() string {
	return (time.Duration(acctest.RandIntRange(1, 604800)) * time.Second).String()
}

func skipIfEnvSet(t *testing.T, envs ...string) {
	for _, k := range envs {
		if os.Getenv(k) != "" {
			t.Skipf("Environment variable %s is set, skipping...", k)
		}
	}
}

func skipIfEnvNotSet(t *testing.T, envs ...string) {
	for _, k := range envs {
		if os.Getenv(k) == "" {
			t.Skipf("Environment variable %s is not set, skipping...", k)
		}
	}
}

func TestCheckResourceAttrNot(name, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := resource.TestCheckResourceAttr(name, key, value)(s)
		if err == nil {
			return fmt.Errorf("%s: Attribute '%s' expected to not equal %#v", name, key, value)
		}

		return nil
	}
}

var keycloakServerInfoVersion *version.Version

func keycloakVersionIsGreaterThanOrEqualTo(keycloakClient *keycloak.KeycloakClient, keycloakMajorVersion *version.Version) (bool, error) {
	if keycloakServerInfoVersion == nil {
		serverInfo, err := keycloakClient.GetServerInfo()
		if err != nil {
			return false, fmt.Errorf("/serverInfo endpoint retuned an error, server Keycloak version could not be determined: %s", err)
		}

		regex := regexp.MustCompile(`^(\d+\.\d+\.\d+)`)
		semver := regex.FindStringSubmatch(serverInfo.SystemInfo.ServerVersion)[0]

		keycloakServerInfoVersion, err = version.NewVersion(semver)
		if err != nil {
			return false, fmt.Errorf("/serverInfo endpoint retuned an unreadable version, server Keycloak version could not be determined: %s", err)
		}
	}
	return keycloakServerInfoVersion.GreaterThanOrEqual(keycloakMajorVersion), nil
}

func getKeycloakVersion600() *version.Version {
	v, _ := version.NewVersion("6.0.0")
	return v
}

func getKeycloakVersion700() *version.Version {
	v, _ := version.NewVersion("7.0.0")
	return v
}

func getKeycloakVersion800() *version.Version {
	v, _ := version.NewVersion("8.0.0")
	return v
}

func getKeycloakVersion900() *version.Version {
	v, _ := version.NewVersion("9.0.0")
	return v
}

func getKeycloakVersion1000() *version.Version {
	v, _ := version.NewVersion("10.0.0")
	return v
}
