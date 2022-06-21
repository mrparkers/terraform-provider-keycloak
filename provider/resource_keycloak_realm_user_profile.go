package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakRealmUserProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakRealmUserProfileCreate,
		ReadContext:   resourceKeycloakRealmUserProfileRead,
		DeleteContext: resourceKeycloakRealmUserProfileDelete,
		UpdateContext: resourceKeycloakRealmUserProfileUpdate,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attribute": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled_when_scope": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"required_for_roles": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"required_for_scopes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"permissions": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"view": {
										Type:     schema.TypeSet,
										Set:      schema.HashString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"edit": {
										Type:     schema.TypeSet,
										Set:      schema.HashString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"validator": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"length": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"max": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"trim_disabled": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"integer": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"max": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"double": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min": {
													Type:     schema.TypeFloat,
													Required: true,
												},
												"max": {
													Type:     schema.TypeFloat,
													Required: true,
												},
											},
										},
									},
									"uri": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"pattern": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"pattern": {
													Type:     schema.TypeString,
													Required: true,
												},
												"error_message": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"email": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"local_date": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"person_name_prohibited_characters": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"error_message": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"username_prohibited_characters": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"error_message": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"options": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"options": {
													Type:     schema.TypeList,
													MinItems: 1,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"annotations": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"display_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func getRealmUserProfileAttributeFromData(m map[string]interface{}) *keycloak.RealmUserProfileAttribute {
	attribute := &keycloak.RealmUserProfileAttribute{
		Name:        m["name"].(string),
		DisplayName: m["display_name"].(string),
		Group:       m["group"].(string),
	}

	if v, ok := m["permissions"]; ok && len(v.([]interface{})) > 0 {
		permissions := keycloak.RealmUserProfilePermissions{
			Edit: make([]string, 0),
			View: make([]string, 0),
		}

		permissionsConfig := v.([]interface{})[0].(map[string]interface{})

		if v, ok := permissionsConfig["view"]; ok {
			permView := make([]string, 0)
			for _, perm := range v.(*schema.Set).List() {
				permView = append(permView, perm.(string))
			}
			permissions.View = permView
		}

		if v, ok := permissionsConfig["edit"]; ok {
			permEdit := make([]string, 0)
			for _, perm := range v.(*schema.Set).List() {
				permEdit = append(permEdit, perm.(string))
			}
			permissions.Edit = permEdit
		}

		attribute.Permissions = &permissions
	}

	if v, ok := m["enabled_when_scope"]; ok && len(interfaceSliceToStringSlice(v.(*schema.Set).List())) != 0 {
		attribute.Selector = &keycloak.RealmUserProfileSelector{
			Scopes: interfaceSliceToStringSlice(v.(*schema.Set).List()),
		}
	}

	if v, ok := m["validator"]; ok && len(v.([]interface{})) > 0 {
		validations := keycloak.RealmUserProfileValidationConfig{}

		data := v.([]interface{})[0].(map[string]interface{})

		if val, ok := data["length"]; ok && len(val.([]interface{})) > 0 {
			r := val.([]interface{})[0].(map[string]interface{})
			validations.Length = &keycloak.RealmUserProfileValidationLength{
				Min:          r["min"].(int),
				Max:          r["max"].(int),
				TrimDisabled: r["trim_disabled"].(bool),
			}
		}

		if val, ok := data["integer"]; ok && len(val.([]interface{})) > 0 {
			r := val.([]interface{})[0].(map[string]interface{})
			validations.Integer = &keycloak.RealmUserProfileValidationInteger{
				Min: r["min"].(int),
				Max: r["max"].(int),
			}
		}

		if val, ok := data["double"]; ok && len(val.([]interface{})) > 0 {
			r := val.([]interface{})[0].(map[string]interface{})
			validations.Double = &keycloak.RealmUserProfileValidationDouble{
				Min: r["min"].(float64),
				Max: r["max"].(float64),
			}
		}

		if val, ok := data["uri"]; ok && len(val.([]interface{})) > 0 {
			validations.URI = &map[string]interface{}{}
		}

		if val, ok := data["pattern"]; ok && len(val.([]interface{})) > 0 {
			r := val.([]interface{})[0].(map[string]interface{})
			validations.Pattern = &keycloak.RealmUserProfileValidationPattern{
				Pattern:      r["pattern"].(string),
				ErrorMessage: r["error_message"].(string),
			}
		}

		if val, ok := data["email"]; ok && len(val.([]interface{})) > 0 {
			validations.Email = &map[string]interface{}{}
		}

		if val, ok := data["local_date"]; ok && len(val.([]interface{})) > 0 {
			validations.LocalDate = &map[string]interface{}{}
		}

		if val, ok := data["person_name_prohibited_characters"]; ok && len(val.([]interface{})) > 0 {
			if r := val.([]interface{})[0]; r == nil {
				validations.PersonNameProhibitedChars = &keycloak.RealmUserProfileValidationProhibited{}
			} else {
				validations.PersonNameProhibitedChars = &keycloak.RealmUserProfileValidationProhibited{
					ErrorMessage: r.(map[string]interface{})["error_message"].(string),
				}
			}

		}

		if val, ok := data["username_prohibited_characters"]; ok && len(val.([]interface{})) > 0 {
			if r := val.([]interface{})[0]; r == nil {
				validations.UsernameProhibitedChars = &keycloak.RealmUserProfileValidationProhibited{}
			} else {
				validations.UsernameProhibitedChars = &keycloak.RealmUserProfileValidationProhibited{
					ErrorMessage: r.(map[string]interface{})["error_message"].(string),
				}
			}
		}

		if val, ok := data["options"]; ok && len(val.([]interface{})) > 0 {
			r := val.([]interface{})[0].(map[string]interface{})
			validations.Options = &keycloak.RealmUserProfileValidationOptions{
				Options: interfaceSliceToStringSlice(r["options"].([]interface{})),
			}
		}

		attribute.Validations = &validations
	}

	required := &keycloak.RealmUserProfileRequired{}

	if v, ok := m["required_for_roles"]; ok {
		required.Roles = interfaceSliceToStringSlice(v.(*schema.Set).List())
	}
	if v, ok := m["required_for_scopes"]; ok {
		required.Scopes = interfaceSliceToStringSlice(v.(*schema.Set).List())
	}

	if len(required.Roles) != 0 || len(required.Scopes) != 0 {
		attribute.Required = required
	}

	if v, ok := m["annotations"]; ok {
		annotations := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			annotations[key] = value.(string)
		}
		attribute.Annotations = annotations
	}

	return attribute

}

func getRealmUserProfileAttributesFromData(lst []interface{}) []*keycloak.RealmUserProfileAttribute {

	attributes := make([]*keycloak.RealmUserProfileAttribute, 0)

	for _, m := range lst {
		userProfileAttribute := getRealmUserProfileAttributeFromData(m.(map[string]interface{}))
		if userProfileAttribute.Name != "" {
			attributes = append(attributes, userProfileAttribute)
		}
	}

	return attributes
}

func getRealmUserProfileGroupFromData(m map[string]interface{}) *keycloak.RealmUserProfileGroup {
	group := keycloak.RealmUserProfileGroup{
		DisplayDescription: m["display_description"].(string),
		DisplayHeader:      m["display_header"].(string),
		Name:               m["name"].(string),
	}

	if v, ok := m["annotations"]; ok {
		annotations := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			annotations[key] = value.(string)
		}
		group.Annotations = annotations
	}

	return &group

}
func getRealmUserProfileGroupsFromData(lst []interface{}) []*keycloak.RealmUserProfileGroup {
	groups := make([]*keycloak.RealmUserProfileGroup, 0)

	for _, m := range lst {
		userProfileGroup := getRealmUserProfileGroupFromData(m.(map[string]interface{}))
		if userProfileGroup.Name != "" {
			groups = append(groups, userProfileGroup)
		}
	}

	return groups
}

func getRealmUserProfileFromData(data *schema.ResourceData) *keycloak.RealmUserProfile {
	realmUserProfile := &keycloak.RealmUserProfile{}

	realmUserProfile.Attributes = getRealmUserProfileAttributesFromData(data.Get("attribute").([]interface{}))
	realmUserProfile.Groups = getRealmUserProfileGroupsFromData(data.Get("group").(*schema.Set).List())

	return realmUserProfile
}

func getRealmUserProfileAttributeData(attr *keycloak.RealmUserProfileAttribute) map[string]interface{} {
	attributeData := make(map[string]interface{})

	attributeData["name"] = attr.Name

	attributeData["display_name"] = attr.DisplayName
	attributeData["group"] = attr.Group
	if attr.Selector != nil && len(attr.Selector.Scopes) != 0 {
		attributeData["enabled_when_scope"] = attr.Selector.Scopes
	}

	attributeData["required_for_roles"] = make([]string, 0)
	attributeData["required_for_scopes"] = make([]string, 0)
	if attr.Required != nil {
		attributeData["required_for_roles"] = attr.Required.Roles
		attributeData["required_for_scopes"] = attr.Required.Scopes
	}

	if attr.Permissions != nil {
		permission := make(map[string]interface{})

		permission["edit"] = attr.Permissions.Edit
		permission["view"] = attr.Permissions.View

		attributeData["permissions"] = []interface{}{permission}
	}

	if attr.Validations != nil {
		validation := make(map[string]interface{})

		if attr.Validations.Pattern != nil {
			validation["pattern"] = []map[string]interface{}{
				{
					"pattern":       attr.Validations.Pattern.Pattern,
					"error_message": attr.Validations.Pattern.ErrorMessage,
				},
			}
		}

		if attr.Validations.PersonNameProhibitedChars != nil {
			validation["person_name_prohibited_characters"] = []map[string]interface{}{
				{"error_message": attr.Validations.PersonNameProhibitedChars.ErrorMessage},
			}
		}

		if attr.Validations.UsernameProhibitedChars != nil {
			validation["username_prohibited_characters"] = []map[string]interface{}{
				{"error_message": attr.Validations.UsernameProhibitedChars.ErrorMessage},
			}
		}

		if attr.Validations.Length != nil {
			validation["length"] = []map[string]interface{}{
				{
					"min":           attr.Validations.Length.Min,
					"max":           attr.Validations.Length.Max,
					"trim_disabled": attr.Validations.Length.TrimDisabled,
				},
			}
		}

		if attr.Validations.Integer != nil {
			validation["integer"] = []map[string]interface{}{
				{
					"min": attr.Validations.Integer.Min,
					"max": attr.Validations.Integer.Max,
				},
			}
		}

		if attr.Validations.Double != nil {
			validation["double"] = []map[string]interface{}{
				{
					"min": attr.Validations.Double.Min,
					"max": attr.Validations.Double.Max,
				},
			}
		}

		if attr.Validations.URI != nil {
			validation["uri"] = []map[string]interface{}{{}}
		}

		if attr.Validations.Email != nil {
			validation["email"] = []map[string]interface{}{{}}
		}

		if attr.Validations.LocalDate != nil {
			validation["local_date"] = []map[string]interface{}{{}}
		}

		if attr.Validations.Options != nil {
			validation["options"] = []map[string]interface{}{
				{"options": attr.Validations.Options.Options},
			}
		}

		attributeData["validator"] = []map[string]interface{}{validation}
	}

	if attr.Annotations != nil {
		attributeData["annotations"] = attr.Annotations
	}

	return attributeData
}

func getRealmUserProfileGroupData(group *keycloak.RealmUserProfileGroup) map[string]interface{} {
	groupData := make(map[string]interface{})

	groupData["name"] = group.Name
	groupData["display_header"] = group.DisplayHeader
	groupData["display_description"] = group.DisplayDescription
	groupData["annotations"] = group.Annotations

	return groupData
}

func setRealmUserProfileData(data *schema.ResourceData, realmUserProfile *keycloak.RealmUserProfile) {
	attributes := make([]interface{}, 0)
	for _, attr := range realmUserProfile.Attributes {
		attributes = append(attributes, getRealmUserProfileAttributeData(attr))
	}
	data.Set("attribute", attributes)

	groups := make([]interface{}, 0)
	for _, group := range realmUserProfile.Groups {
		groups = append(groups, getRealmUserProfileGroupData(group))
	}
	data.Set("group", groups)
}

func resourceKeycloakRealmUserProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	data.SetId(realmId)

	realmUserProfile := getRealmUserProfileFromData(data)

	err := keycloakClient.UpdateRealmUserProfile(ctx, realmId, realmUserProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeycloakRealmUserProfileRead(ctx, data, meta)
}

func resourceKeycloakRealmUserProfileRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)

	realmUserProfile, err := keycloakClient.GetRealmUserProfile(ctx, realmId)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	setRealmUserProfileData(data, realmUserProfile)

	return nil
}

func resourceKeycloakRealmUserProfileDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)

	// The realm user profile cannot be deleted, so instead we set it back to its "zero" values.
	realmUserProfile := &keycloak.RealmUserProfile{
		Attributes: []*keycloak.RealmUserProfileAttribute{},
		Groups:     []*keycloak.RealmUserProfileGroup{},
	}

	err := keycloakClient.UpdateRealmUserProfile(ctx, realmId, realmUserProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKeycloakRealmUserProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	realmUserProfile := getRealmUserProfileFromData(data)

	err := keycloakClient.UpdateRealmUserProfile(ctx, realmId, realmUserProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	setRealmUserProfileData(data, realmUserProfile)

	return nil
}
