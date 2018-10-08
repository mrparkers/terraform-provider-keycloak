package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func randomBool() bool {
	return rand.Intn(2) == 0
}

func randomStringInSlice(slice []string) string {
	return slice[acctest.RandIntRange(0, len(slice)-1)]
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
