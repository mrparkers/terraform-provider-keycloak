package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakUserRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_verified": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"federated_identity": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"federated_identities": {
				Type:       schema.TypeSet,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Computed:   true,
				Deprecated: "use federated_identity instead",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceKeycloakUserRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmID := data.Get("realm_id").(string)
	username := data.Get("username").(string)

	user, err := keycloakClient.GetUserByUsername(realmID, username)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with username %s not found", username)
	}

	mapFromUserToData(data, user)

	// handling deprecated "federated_identities" attribute
	federatedIdentities := []interface{}{}
	for _, federatedIdentity := range user.FederatedIdentities {
		identity := map[string]interface{}{
			"identity_provider": federatedIdentity.IdentityProvider,
			"user_id":           federatedIdentity.UserId,
			"user_name":         federatedIdentity.UserName,
		}
		federatedIdentities = append(federatedIdentities, identity)
	}
	data.Set("federated_identities", federatedIdentities)

	return nil
}
