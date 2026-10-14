package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenidClientSecretReads(t *testing.T) {
	for _, test := range []struct {
		name, realm, clientID      string
		public, bearer, readSecret bool
	}{
		{"public", "fixture", "public", true, false, false},
		{"realm_management", "fixture", "realm-management", false, true, false},
		{"master_admin", "master", "fixture-realm", false, true, false},
		{"confidential_suffix", "master", "payments-realm", false, false, true},
		{"confidential_same_name", "fixture", "realm-management", false, false, true},
		{"ordinary_bearer", "fixture", "broker", false, true, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			reads := 0
			model := map[string]interface{}{"id": "id", "clientId": test.clientID, "publicClient": test.public, "bearerOnly": test.bearer, "secret": "representation-secret", "attributes": map[string]string{}}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case r.URL.Path == "/admin/serverinfo":
					fmt.Fprint(w, `{"systemInfo":{"version":"26.8.1"},"themes":{}}`)
				case strings.HasSuffix(r.URL.Path, "/client-secret"):
					reads++
					if !test.readSecret {
						w.WriteHeader(403)
						return
					}
					fmt.Fprint(w, `{"value":"endpoint-secret"}`)
				case strings.HasSuffix(r.URL.Path, "/clients/id"):
					json.NewEncoder(w).Encode(model)
				case strings.HasSuffix(r.URL.Path, "/clients"):
					json.NewEncoder(w).Encode([]map[string]interface{}{model})
				default:
					t.Errorf("unexpected request %s", r.URL.Path)
					w.WriteHeader(404)
				}
			}))
			defer server.Close()
			api := &KeycloakClient{baseUrl: server.URL, httpClient: server.Client(), clientCredentials: &ClientCredentials{AccessToken: "fixture-token", TokenType: "Bearer"}, accessTokenProvided: true}
			ctx := context.Background()
			found, err := api.GetOpenidClient(ctx, test.realm, "id")
			if err != nil {
				t.Fatal(err)
			}
			byName, err := api.GetOpenidClientByClientId(ctx, test.realm, test.clientID)
			if err != nil {
				t.Fatal(err)
			}
			list, err := api.GetOpenidClients(ctx, test.realm, true)
			if err != nil {
				t.Fatal(err)
			}
			expected := ""
			wantReads := 0
			if test.readSecret {
				expected = "endpoint-secret"
				wantReads = 3
			}
			if found.ClientSecret != expected || byName.ClientSecret != expected || len(list) != 1 || list[0].ClientSecret != expected || reads != wantReads {
				t.Fatalf("credential reads=%d, want %d; unexpected credential projection", reads, wantReads)
			}
			without, err := api.GetOpenidClientWithoutSecret(ctx, test.realm, "id")
			if err != nil {
				t.Fatal(err)
			}
			withoutByName, err := api.GetOpenidClientByClientIdSkipSecret(ctx, test.realm, test.clientID)
			if err != nil {
				t.Fatal(err)
			}
			withoutList, err := api.GetOpenidClients(ctx, test.realm, false)
			if err != nil {
				t.Fatal(err)
			}
			if without.ClientSecret != "" || withoutByName.ClientSecret != "" || withoutList[0].ClientSecret != "" || reads != wantReads {
				t.Fatal("explicit configuration-only reads retrieved or retained a credential")
			}
		})
	}
}
