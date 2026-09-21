package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

const (
	userFederationSyncModeFull         = "full"
	userFederationSyncModeChangedUsers = "changed_users"
)

var userFederationSyncModes = map[string]string{
	userFederationSyncModeFull:         keycloak.UserFederationSyncFull,
	userFederationSyncModeChangedUsers: keycloak.UserFederationSyncChangedUsers,
}

var _ action.ActionWithConfigure = &userFederationSyncAction{}
var _ action.ActionWithValidateConfig = &userFederationSyncAction{}

type userFederationSyncAction struct {
	keycloakClient *keycloak.KeycloakClient
}

type userFederationSyncActionModel struct {
	RealmId          types.String `tfsdk:"realm_id"`
	UserFederationId types.String `tfsdk:"user_federation_id"`
	Mode             types.String `tfsdk:"mode"`
	ClientTimeout    types.Int64  `tfsdk:"client_timeout"`
}

func NewUserFederationSyncAction() action.Action {
	return &userFederationSyncAction{}
}

func (a *userFederationSyncAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_federation_sync"
}

func (a *userFederationSyncAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		Description: "Triggers a synchronization of users for a user federation provider, e.g. an LDAP user federation.",
		Attributes: map[string]actionschema.Attribute{
			"realm_id": actionschema.StringAttribute{
				Required:    true,
				Description: "The realm the user federation provider exists in.",
			},
			"user_federation_id": actionschema.StringAttribute{
				Required:    true,
				Description: "The ID of the user federation provider to synchronize, e.g. the ID of a keycloak_ldap_user_federation resource.",
			},
			"mode": actionschema.StringAttribute{
				Optional:    true,
				Description: "The kind of synchronization to trigger. Can be one of `full` or `changed_users`. Defaults to `full`.",
			},
			"client_timeout": actionschema.Int64Attribute{
				Optional:    true,
				Description: "Timeout (in seconds) of the sync request. Defaults to the `client_timeout` of the provider.",
			},
		},
	}
}

func (a *userFederationSyncAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	// provider data is not available during validation
	if req.ProviderData == nil {
		return
	}

	keycloakClient, ok := req.ProviderData.(*keycloak.KeycloakClient)
	if !ok {
		resp.Diagnostics.AddError("unexpected provider data", fmt.Sprintf("expected *keycloak.KeycloakClient, got %T", req.ProviderData))
		return
	}

	a.keycloakClient = keycloakClient
}

func (a *userFederationSyncAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var config userFederationSyncActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Mode.IsNull() && !config.Mode.IsUnknown() {
		if _, ok := userFederationSyncModes[config.Mode.ValueString()]; !ok {
			resp.Diagnostics.AddAttributeError(path.Root("mode"), "invalid mode", fmt.Sprintf("expected mode to be one of %q or %q, got %q", userFederationSyncModeFull, userFederationSyncModeChangedUsers, config.Mode.ValueString()))
		}
	}

	if !config.ClientTimeout.IsNull() && !config.ClientTimeout.IsUnknown() && config.ClientTimeout.ValueInt64() <= 0 {
		resp.Diagnostics.AddAttributeError(path.Root("client_timeout"), "invalid client_timeout", fmt.Sprintf("expected client_timeout to be greater than 0, got %d", config.ClientTimeout.ValueInt64()))
	}
}

func (a *userFederationSyncAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config userFederationSyncActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mode := userFederationSyncModeFull
	if !config.Mode.IsNull() {
		mode = config.Mode.ValueString()
	}

	if !config.ClientTimeout.IsNull() {
		ctx = keycloak.WithRequestTimeout(ctx, time.Duration(config.ClientTimeout.ValueInt64())*time.Second)
	}

	realmId := config.RealmId.ValueString()
	userFederationId := config.UserFederationId.ValueString()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Triggering %s sync of user federation %s in realm %s", mode, userFederationId, realmId),
	})

	result, err := a.keycloakClient.SyncUserFederation(ctx, realmId, userFederationId, userFederationSyncModes[mode])
	if err != nil {
		resp.Diagnostics.AddError("user federation sync failed", err.Error())
		return
	}

	if result.Ignored {
		resp.Diagnostics.AddWarning("user federation sync was ignored", fmt.Sprintf("Keycloak ignored the sync request: %s", result.Status))
		return
	}

	if result.Failed > 0 {
		resp.Diagnostics.AddWarning("user federation sync finished with failures", fmt.Sprintf("%d users could not be synchronized: %s", result.Failed, result.Status))
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Sync finished: %d added, %d updated, %d removed, %d failed", result.Added, result.Updated, result.Removed, result.Failed),
	})
}
