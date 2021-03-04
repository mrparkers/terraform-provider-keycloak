package provider

import (
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
	"log"
	"math/rand"
	"strings"
)

var (
	keycloakIdpTokenExchangePermissionPolicyTypes = []string{"client"}
)

func resourceKeycloakIdentityProviderTokenExchangeScopePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakIdentityProviderTokenExchangeScopePermissionCreate,
		Read:   resourceKeycloakIdentityProviderTokenExchangeScopePermissionRead,
		Delete: resourceKeycloakIdentityProviderTokenExchangeScopePermissionDelete,
		Update: resourceKeycloakIdentityProviderTokenExchangeScopePermissionUpdate,
		// This resource can be imported using {{realmId}}/{{providerAlias}}. The provider alias is displayed in the URL when editing it from the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakIdentityProviderTokenExchangeScopePermissionImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "client",
				Description:  "Type of policy that is created. At the moment only 'client' type is supported",
				ValidateFunc: validation.StringInSlice(keycloakIdpTokenExchangePermissionPolicyTypes, false),
			},
			"clients": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Ids of the clients for which a policy will be created and set on scope based token exchange permission",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy id that will be set on the scope based token exchange permission automatically created by enabling permissions on the reference identity provider",
			},
			"authorization_resource_server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource server id representing the realm management client on which this permission is managed",
			},
			"authorization_idp_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource id representing the identity provider, this automatically created by keycloak",
			},
			"authorization_token_exchange_scope_permission_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Permission id representing the Permission with scope 'Token Exchange' and the resource 'authorization_idp_resource_id', this automatically created by keycloak, the policy id will be set on this permission",
			},
		},
	}
}

func setIdentityProviderTokenExchangeScopePermissionClientPolicy(keycloakClient *keycloak.KeycloakClient, realmId, providerAlias string, clients []string) error {
	identityProviderPermissions, err := keycloakClient.GetIdentityProviderPermissions(realmId, providerAlias)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	tokenExchangeScopedPermissionId, err := identityProviderPermissions.GetTokenExchangeScopedPermissionId()
	if err != nil {
		return err
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, tokenExchangeScopedPermissionId)
	if err != nil {
		return err
	}

	if len(permission.Policies) == 0 {
		policyId, err := createClientPolicy(keycloakClient, realmId, realmManagementClient.Id, providerAlias, clients)
		if err != nil {
			return err
		}
		permission.Policies = []string{policyId}
		return keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)

	} else if len(permission.Policies) == 1 {
		openidClientAuthorizationClientPolicy, err := keycloakClient.GetOpenidClientAuthorizationClientPolicy(realmId, realmManagementClient.Id, permission.Policies[0])
		if err != nil {
			return err
		}
		openidClientAuthorizationClientPolicy.Clients = clients
		return keycloakClient.UpdateOpenidClientAuthorizationClientPolicy(openidClientAuthorizationClientPolicy)

	} else {
		return fmt.Errorf("only one client policy is supported, but %d were found", len(permission.Policies))
	}
}

func createClientPolicy(keycloakClient *keycloak.KeycloakClient, realmId, realmManagementClientId, providerAlias string, clients []string) (string, error) {
	openidClientAuthorizationClientPolicy := &keycloak.OpenidClientAuthorizationClientPolicy{
		RealmId:          realmId,
		ResourceServerId: realmManagementClientId,
		Name:             providerAlias + "_idp_client_policy",
		DecisionStrategy: "UNANIMOUS",
		Logic:            "POSITIVE",
		Type:             "client",
		Clients:          clients,
	}
	err := keycloakClient.NewOpenidClientAuthorizationClientPolicy(openidClientAuthorizationClientPolicy)
	if err != nil {
		if keycloak.ErrorIs409(err) {
			b := make([]byte, 4)
			rand.Read(b)
			suffix := hex.EncodeToString(b)
			openidClientAuthorizationClientPolicy.Name = providerAlias + "_" + suffix + "_idp_client_policy"
			err = keycloakClient.NewOpenidClientAuthorizationClientPolicy(openidClientAuthorizationClientPolicy)
		}
	}
	if err != nil {
		return "", err
	}
	return openidClientAuthorizationClientPolicy.Id, nil
}

func unsetIdentityProviderTokenExchangeScopePermissionPolicy(keycloakClient *keycloak.KeycloakClient, realmId, providerAlias, policyId string) error {
	identityProviderPermissions, err := keycloakClient.GetIdentityProviderPermissions(realmId, providerAlias)
	if err != nil {
		return err
	}

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	tokenExchangeScopedPermissionId, err := identityProviderPermissions.GetTokenExchangeScopedPermissionId()
	if err != nil {
		return err
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, tokenExchangeScopedPermissionId)
	if err != nil {
		return err
	}

	permission.Policies = []string{}
	err = keycloakClient.UpdateOpenidClientAuthorizationPermission(permission)
	if err != nil {
		return err
	}

	err = keycloakClient.DisableIdentityProviderPermissions(realmId, providerAlias)
	if err != nil {
		return err
	}

	_ = keycloakClient.DeleteOpenidClientAuthorizationClientPolicy(realmId, realmManagementClient.Id, policyId)

	return nil
}

func resourceKeycloakIdentityProviderTokenExchangeScopePermissionCreate(data *schema.ResourceData, meta interface{}) error {
	return resourceKeycloakIdentityProviderTokenExchangeScopePermissionUpdate(data, meta)
}

func resourceKeycloakIdentityProviderTokenExchangeScopePermissionUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	providerAlias := data.Get("provider_alias").(string)
	policyType := data.Get("policy_type").(string)
	var clients []string

	if v, ok := data.GetOk("clients"); ok {
		for _, client := range v.(*schema.Set).List() {
			clients = append(clients, client.(string))
		}
	}

	err := keycloakClient.EnableIdentityProviderPermissions(realmId, providerAlias)
	if err != nil {
		return err
	}
	if policyType == "client" {
		err = setIdentityProviderTokenExchangeScopePermissionClientPolicy(keycloakClient, realmId, providerAlias, clients)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid policy type, supported types are ['client']")
	}
	return resourceKeycloakIdentityProviderTokenExchangeScopePermissionRead(data, meta)
}

func resourceKeycloakIdentityProviderTokenExchangeScopePermissionRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	providerAlias := data.Get("provider_alias").(string)

	identityProviderPermissions, err := keycloakClient.GetIdentityProviderPermissions(realmId, providerAlias)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	if !identityProviderPermissions.Enabled {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer enabled", data.Id())
		data.SetId("")
		return nil
	}

	data.SetId(identityProviderPermissions.RealmId + "/" + identityProviderPermissions.ProviderAlias)
	data.Set("realm_id", identityProviderPermissions.RealmId)
	data.Set("provider_alias", identityProviderPermissions.ProviderAlias)

	realmManagementClient, err := keycloakClient.GetOpenidClientByClientId(realmId, "realm-management")
	if err != nil {
		return err
	}

	tokenExchangeScopedPermissionId, err := identityProviderPermissions.GetTokenExchangeScopedPermissionId()
	if err != nil {
		return err
	}

	permission, err := keycloakClient.GetOpenidClientAuthorizationPermission(realmId, realmManagementClient.Id, tokenExchangeScopedPermissionId)
	if err != nil {
		return err
	}

	var openidClientAuthorizationClientPolicyId string
	if len(permission.Policies) >= 1 {
		openidClientAuthorizationClientPolicyId = permission.Policies[0]
	} else {
		openidClientAuthorizationClientPolicyId, err = createClientPolicy(keycloakClient, realmId, realmManagementClient.Id, providerAlias, data.Get("clients").([]string))
		if err != nil {
			return err
		}
	}
	openidClientAuthorizationClientPolicy, err := keycloakClient.GetOpenidClientAuthorizationClientPolicy(realmId, realmManagementClient.Id, openidClientAuthorizationClientPolicyId)
	if err != nil {
		return err
	}

	data.Set("policy_id", openidClientAuthorizationClientPolicy.Id)
	data.Set("clients", openidClientAuthorizationClientPolicy.Clients)

	data.Set("policy_type", data.Get("policy_type"))
	data.Set("authorization_resource_server_id", realmManagementClient.Id)
	data.Set("authorization_idp_resource_id", identityProviderPermissions.Resource)

	data.Set("authorization_token_exchange_scope_permission_id", tokenExchangeScopedPermissionId)

	return nil
}

func resourceKeycloakIdentityProviderTokenExchangeScopePermissionDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	providerAlias := data.Get("provider_alias").(string)
	policyId := data.Get("policy_id").(string)

	identityProviderPermissions, err := keycloakClient.GetIdentityProviderPermissions(realmId, providerAlias)
	if err == nil && identityProviderPermissions.Enabled {
		_ = unsetIdentityProviderTokenExchangeScopePermissionPolicy(keycloakClient, realmId, providerAlias, policyId)
	}
	return keycloakClient.DisableIdentityProviderPermissions(realmId, providerAlias)
}

func resourceKeycloakIdentityProviderTokenExchangeScopePermissionImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{providerAlias}}")
	}
	d.SetId(parts[0] + "/" + parts[1])
	d.Set("realm_id", parts[0])
	d.Set("provider_alias", parts[1])
	return []*schema.ResourceData{d}, nil
}
