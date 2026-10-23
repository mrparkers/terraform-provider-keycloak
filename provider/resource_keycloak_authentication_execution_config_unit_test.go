package provider

import (
	"reflect"
	"testing"
)

func TestMergeMaskedAuthenticationExecutionConfig(t *testing.T) {
	t.Parallel()

	remote := map[string]string{
		"secret.key":      keycloakMaskedSecret,
		"site.key":        "site",
		"imported.secret": keycloakMaskedSecret,
		"empty.secret":    keycloakMaskedSecret,
	}
	known := map[string]interface{}{
		"secret.key":   "my-secret",
		"site.key":     "outdated-site",
		"empty.secret": "",
		"removed.key":  "gone",
	}

	got := mergeMaskedAuthenticationExecutionConfig(remote, known)

	want := map[string]string{
		"secret.key":      "my-secret",          // masked -> known value is kept
		"site.key":        "site",               // not masked -> remote value wins
		"imported.secret": keycloakMaskedSecret, // masked and unknown -> stays masked
		"empty.secret":    keycloakMaskedSecret, // masked and known empty -> stays masked
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected merge result\ngot:  %#v\nwant: %#v", got, want)
	}
}
