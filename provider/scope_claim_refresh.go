package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

// clearScopeClaimOnRefresh wraps a scope resource's read so that a refresh drops the resource's claim from earlier
// runs. Terraform refreshes every resource in state before reconciling any of them, so claims only ever reflect the
// resources reconciled in the current run and never a previous one.
//
// The reconcile functions call the unwrapped read, which keeps the claim they have just made.
func clearScopeClaimOnRefresh(ownerKey func(data *schema.ResourceData) string, scopeType string, read schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		meta.(*keycloak.KeycloakClient).ReleaseScopes(ownerKey(data), scopeType)

		return read(ctx, data, meta)
	}
}

func openidClientScopeOwnerKey(data *schema.ResourceData) string {
	return keycloak.OpenidClientScopeOwnerKey(data.Get("realm_id").(string), data.Get("client_id").(string))
}

func samlClientScopeOwnerKey(data *schema.ResourceData) string {
	return keycloak.SamlClientScopeOwnerKey(data.Get("realm_id").(string), data.Get("client_id").(string))
}

func realmScopeOwnerKey(data *schema.ResourceData) string {
	return keycloak.RealmScopeOwnerKey(data.Get("realm_id").(string))
}
