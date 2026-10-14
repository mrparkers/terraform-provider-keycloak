package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationCreate,
		ReadContext:   resourceKeycloakOrganizationRead,
		DeleteContext: resourceKeycloakOrganizationDelete,
		UpdateContext: resourceKeycloakOrganizationUpdate,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the organization.",
			},
			"realm": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Realm ID.",
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The alias unique identifies the organization. Same as the name if not specified.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable/disable this organization.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redirect_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Landing page after successful login.",
			},
			"domain": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"verified": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"identity_provider_alias": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The alias of an identity provider in an identity_provider block that users whose email matches this domain authenticate through.",
						},
						"auto_redirect": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Redirect users whose email matches this domain to the identity provider without showing the login form.",
						},
					},
				},
			},
			"identity_provider": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The identity providers that this resource links to the organization. Links made elsewhere are kept. When set, the resource also manages the routing of the domains. Requires Keycloak 26.8 or later.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The alias of the identity provider.",
						},
						"auto_membership": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Add users who authenticate through the identity provider as members of the organization.",
						},
						"membership_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "UNMANAGED",
							ValidateFunc: validation.StringInSlice([]string{"MANAGED", "UNMANAGED"}, false),
							Description:  "The membership type of users added by auto_membership.",
						},
					},
				},
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				// ignore ordering of multi-valued attributes
				DiffSuppressFunc: suppressDiffForMultivalueAttributeOrder(),
			},
		},
	}
}

func getOrganizationFromData(data *schema.ResourceData) (*keycloak.Organization, error) {
	organization := &keycloak.Organization{
		Id:          data.Id(),
		Realm:       data.Get("realm").(string),
		Name:        data.Get("name").(string),
		Alias:       data.Get("alias").(string),
		Enabled:     data.Get("enabled").(bool),
		Description: data.Get("description").(string),
		RedirectUrl: data.Get("redirect_url").(string),
	}

	identityProviderAliases := make(map[string]bool)
	for _, identityProvider := range getOrganizationIdentityProvidersFromData(data) {
		identityProviderAliases[identityProvider.Alias] = true
	}

	domains := make([]keycloak.OrganizationDomain, 0)
	if v, ok := data.GetOk("domain"); ok {
		for _, domain := range v.(*schema.Set).List() {
			domainMap := domain.(map[string]interface{})
			orgDomain := keycloak.OrganizationDomain{
				Name:                  domainMap["name"].(string),
				Verified:              domainMap["verified"].(bool),
				IdentityProviderAlias: domainMap["identity_provider_alias"].(string),
				AutoRedirect:          domainMap["auto_redirect"].(bool),
			}
			// Keycloak routes a domain only to an identity provider linked to the organization.
			if orgDomain.IdentityProviderAlias != "" && !identityProviderAliases[orgDomain.IdentityProviderAlias] {
				return nil, fmt.Errorf("domain \"%s\" of organization \"%s\" routes to identity provider \"%s\", which is not in an identity_provider block", orgDomain.Name, organization.Name, orgDomain.IdentityProviderAlias)
			}
			if orgDomain.AutoRedirect && orgDomain.IdentityProviderAlias == "" {
				return nil, fmt.Errorf("auto_redirect on domain \"%s\" of organization \"%s\" requires identity_provider_alias", orgDomain.Name, organization.Name)
			}
			domains = append(domains, orgDomain)
		}
	}

	organization.Domains = domains

	attributes := map[string][]string{}
	if v, ok := data.GetOk("attributes"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = strings.Split(value.(string), MULTIVALUE_ATTRIBUTE_SEPARATOR)
		}
	}
	organization.Attributes = attributes

	return organization, nil
}

func getOrganizationIdentityProvidersFromData(data *schema.ResourceData) []keycloak.OrganizationIdentityProvider {
	return getOrganizationIdentityProvidersFromSet(data.Get("identity_provider").(*schema.Set))
}

func getOrganizationIdentityProvidersFromSet(set *schema.Set) []keycloak.OrganizationIdentityProvider {
	identityProviders := make([]keycloak.OrganizationIdentityProvider, 0, set.Len())
	for _, identityProvider := range set.List() {
		identityProviderMap := identityProvider.(map[string]interface{})
		identityProviders = append(identityProviders, keycloak.OrganizationIdentityProvider{
			Alias:          identityProviderMap["alias"].(string),
			AutoMembership: identityProviderMap["auto_membership"].(bool),
			MembershipType: identityProviderMap["membership_type"].(string),
		})
	}

	return identityProviders
}

func setOrganizationData(data *schema.ResourceData, organization *keycloak.Organization) error {
	attributes := map[string]string{}
	for k, v := range organization.Attributes {
		attributes[k] = strings.Join(v, MULTIVALUE_ATTRIBUTE_SEPARATOR)
	}

	domains := make([]map[string]interface{}, 0, len(organization.Domains))
	for _, domain := range organization.Domains {
		domainMap := make(map[string]interface{})
		domainMap["name"] = domain.Name
		domainMap["verified"] = domain.Verified
		domainMap["identity_provider_alias"] = domain.IdentityProviderAlias
		domainMap["auto_redirect"] = domain.AutoRedirect
		domains = append(domains, domainMap)
	}

	data.SetId(organization.Id)
	data.Set("name", organization.Name)
	data.Set("realm", organization.Realm)
	data.Set("alias", organization.Alias)
	data.Set("enabled", organization.Enabled)
	data.Set("description", organization.Description)
	data.Set("redirect_url", organization.RedirectUrl)
	data.Set("domain", domains)
	data.Set("attributes", attributes)

	return nil
}

func setOrganizationIdentityProvidersData(data *schema.ResourceData, identityProviders []keycloak.OrganizationIdentityProvider) {
	identityProviderMaps := make([]map[string]interface{}, 0, len(identityProviders))
	for _, identityProvider := range identityProviders {
		identityProviderMaps = append(identityProviderMaps, map[string]interface{}{
			"alias":           identityProvider.Alias,
			"auto_membership": identityProvider.AutoMembership,
			"membership_type": identityProvider.MembershipType,
		})
	}

	data.Set("identity_provider", identityProviderMaps)
}

func withoutDomainRouting(domains []keycloak.OrganizationDomain) []keycloak.OrganizationDomain {
	unrouted := make([]keycloak.OrganizationDomain, 0, len(domains))
	for _, domain := range domains {
		unrouted = append(unrouted, keycloak.OrganizationDomain{
			Name:     domain.Name,
			Verified: domain.Verified,
		})
	}

	return unrouted
}

func checkOrganizationIdentityProvidersSupported(ctx context.Context, keycloakClient *keycloak.KeycloakClient, organization *keycloak.Organization) error {
	if ok, _ := keycloakClient.VersionIsLessThan(ctx, keycloak.Version_26_8); ok {
		return fmt.Errorf("identity_provider in organization \"%s\" is not supported by your Keycloak version (requires >= %s)", organization.Name, keycloak.Version_26_8)
	}

	return nil
}

// linkOrganizationIdentityProviders links each identity provider to the organization, or updates
// its link, so that a domain can route to it.
func linkOrganizationIdentityProviders(ctx context.Context, keycloakClient *keycloak.KeycloakClient, organization *keycloak.Organization, identityProviders []keycloak.OrganizationIdentityProvider) error {
	for _, identityProvider := range identityProviders {
		current, err := keycloakClient.GetOrganizationIdentityProvider(ctx, organization.Realm, organization.Id, identityProvider.Alias)
		if keycloak.ErrorIs404(err) {
			if err := keycloakClient.LinkIdentityProviderWithOrganization(ctx, organization.Realm, identityProvider.Alias, organization.Id); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if current == nil || *current != identityProvider {
			if err := keycloakClient.UpdateOrganizationIdentityProviderLink(ctx, organization.Realm, organization.Id, identityProvider); err != nil {
				return err
			}
		}
	}

	return nil
}

// unlinkOrganizationIdentityProviders removes the links that the resource had and no longer has.
// Links made elsewhere, such as by an identity provider resource with organization_id, are kept.
func unlinkOrganizationIdentityProviders(ctx context.Context, keycloakClient *keycloak.KeycloakClient, organization *keycloak.Organization, identityProviders, oldIdentityProviders []keycloak.OrganizationIdentityProvider) error {
	wanted := make(map[string]bool, len(identityProviders))
	for _, identityProvider := range identityProviders {
		wanted[identityProvider.Alias] = true
	}

	for _, identityProvider := range oldIdentityProviders {
		if wanted[identityProvider.Alias] {
			continue
		}
		err := keycloakClient.UnlinkIdentityProviderFromOrganization(ctx, organization.Realm, identityProvider.Alias, organization.Id)
		if err != nil && !keycloak.ErrorIs404(err) {
			return err
		}
	}

	return nil
}

func resourceKeycloakOrganizationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realm}}/{{OrganizationId}}")
	}

	_, err := keycloakClient.GetOrganization(ctx, parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	d.Set("realm", parts[0])
	d.SetId(parts[1])

	diagnostics := resourceKeycloakOrganizationRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, fmt.Errorf("Error reading Organization: %s", diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceKeycloakOrganizationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	Organization, err := getOrganizationFromData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	identityProviders := getOrganizationIdentityProvidersFromData(data)
	if len(identityProviders) > 0 {
		if err = checkOrganizationIdentityProvidersSupported(ctx, keycloakClient, Organization); err != nil {
			return diag.FromErr(err)
		}
	}

	// The domains are created without routing: an identity provider can only be linked to the
	// organization, and a domain routed to it, once the organization exists.
	domains := Organization.Domains
	Organization.Domains = withoutDomainRouting(domains)
	if err = keycloakClient.NewOrganization(ctx, Organization); err != nil {
		return diag.FromErr(err)
	}
	Organization.Domains = domains
	if err = setOrganizationData(data, Organization); err != nil {
		return diag.FromErr(err)
	}

	if len(identityProviders) > 0 {
		if err = linkOrganizationIdentityProviders(ctx, keycloakClient, Organization, identityProviders); err != nil {
			return diag.FromErr(err)
		}
		if err = keycloakClient.UpdateOrganizationWithDomainRouting(ctx, Organization); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceKeycloakOrganizationRead(ctx, data, meta)
}

func resourceKeycloakOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm := data.Get("realm").(string)
	id := data.Id()
	Organization, err := keycloakClient.GetOrganization(ctx, realm, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	stateIdentityProviders := getOrganizationIdentityProvidersFromData(data)
	identityProviders := make([]keycloak.OrganizationIdentityProvider, 0, len(stateIdentityProviders))
	for _, stateIdentityProvider := range stateIdentityProviders {
		identityProvider, err := keycloakClient.GetOrganizationIdentityProvider(ctx, realm, id, stateIdentityProvider.Alias)
		if keycloak.ErrorIs404(err) {
			continue
		}
		if err != nil {
			return diag.FromErr(err)
		}
		identityProviders = append(identityProviders, *identityProvider)
	}
	if len(stateIdentityProviders) == 0 {
		// Without identity_provider blocks the routing of the domains belongs to the
		// identity provider resources that set organization_id.
		Organization.Domains = withoutDomainRouting(Organization.Domains)
	}
	setOrganizationIdentityProvidersData(data, identityProviders)

	return diag.FromErr(setOrganizationData(data, Organization))
}

func resourceKeycloakOrganizationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	Organization, err := getOrganizationFromData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	identityProviders := getOrganizationIdentityProvidersFromData(data)
	oldIdentityProviderSet, _ := data.GetChange("identity_provider")
	oldIdentityProviders := getOrganizationIdentityProvidersFromSet(oldIdentityProviderSet.(*schema.Set))

	if len(identityProviders) > 0 {
		if err = checkOrganizationIdentityProvidersSupported(ctx, keycloakClient, Organization); err != nil {
			return diag.FromErr(err)
		}
		if err = linkOrganizationIdentityProviders(ctx, keycloakClient, Organization, identityProviders); err != nil {
			return diag.FromErr(err)
		}
		err = keycloakClient.UpdateOrganizationWithDomainRouting(ctx, Organization)
	} else {
		err = keycloakClient.UpdateOrganization(ctx, Organization)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	err = unlinkOrganizationIdentityProviders(ctx, keycloakClient, Organization, identityProviders, oldIdentityProviders)
	if err != nil {
		return diag.FromErr(err)
	}

	setOrganizationIdentityProvidersData(data, identityProviders)

	return diag.FromErr(setOrganizationData(data, Organization))
}

func resourceKeycloakOrganizationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	id := data.Id()

	return diag.FromErr(keycloakClient.DeleteOrganization(ctx, realm, id))
}
