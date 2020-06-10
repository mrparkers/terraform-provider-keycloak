package keycloak

//https://www.keycloak.org/docs-api/5.0/rest-api/index.html#_managementpermissionreference
type managementPermissionReference struct {
	Enabled          bool              `json:"enabled"`
	Resource         string            `json:"resource"`
	ScopePermissions map[string]string `json:"scopePermissions"`
}

func disableClientManagementPermissionsReference() *managementPermissionReference {
	return &managementPermissionReference{
		Enabled: false,
	}
}

func enableClientManagementPermissionsReference() *managementPermissionReference {
	return &managementPermissionReference{
		Enabled: true,
	}
}
