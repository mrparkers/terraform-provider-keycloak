package keycloak

import (
	"strconv"
	"strings"
)

type userFederationComponent struct {
	Id           string              `json:"id,omitempty"`
	Name         string              `json:"name"`
	ProviderId   string              `json:"providerId"`   // ldap
	ProviderType string              `json:"providerType"` // org.keycloak.storage.UserStorageProvider
	ParentId     string              `json:"parentId"`     // realm id
	Config       map[string][]string `json:"config"`       // generic interface, but always includes "cachePolicy", "enabled", and "priority"
}

func (component *userFederationComponent) getConfig(val string) string {
	return component.Config[val][0]
}

type LdapUserFederation struct {
	Id      string
	Name    string
	RealmId string

	Enabled  bool
	Priority int

	ImportEnabled     bool
	EditMode          string // can be "READ_ONLY", "WRITABLE", or "UNSYNCED"
	SyncRegistrations bool   // I think this field controls whether or not BatchSizeForSync, FullSyncPeriod, and ChangedSyncPeriod are needed

	Vendor                 string // can be "other", "edirectory", "ad", "rhds", or "tivoli". honestly I don't think this field actually does anything
	UsernameLDAPAttribute  string
	RdnLDAPAttribute       string
	UuidLDAPAttribute      string
	UserObjectClasses      []string // api expects comma + space separated for some reason
	ConnectionUrl          string
	UsersDn                string
	AuthType               string // can be "simple" or "none". don't need bind fields if set to "none"
	BindDn                 string
	BindCredential         string
	CustomUserSearchFilter string // must start with '(' and end with ')'
	SearchScope            string // api expects "1" or "2", but that means "One Level" or "Subtree"

	ValidatePasswordPolicy bool
	UseTruststoreSpi       string // can be "ldapsOnly", "always", or "never"
	ConnectionTimeout      int
	ReadTimeout            int
	Pagination             bool

	BatchSizeForSync  int
	FullSyncPeriod    int // either a number, in milliseconds, or -1 if full sync is disabled
	ChangedSyncPeriod int // either a number, in milliseconds, or -1 if changed sync is disabled

	CachePolicy string
}

func convertToUserFederationComponent(ldap *LdapUserFederation) *userFederationComponent {
	return &userFederationComponent{
		Id:           ldap.Id,
		Name:         ldap.Name,
		ProviderId:   "ldap",
		ProviderType: "org.keycloak.storage.UserStorageProvider",
		ParentId:     ldap.RealmId,
		Config: map[string][]string{
			"cachePolicy": {
				ldap.CachePolicy,
			},
			"enabled": {
				strconv.FormatBool(ldap.Enabled),
			},
			"priority": {
				strconv.Itoa(ldap.Priority),
			},
			"importEnabled": {
				strconv.FormatBool(ldap.ImportEnabled),
			},
			"editMode": {
				ldap.EditMode,
			},
			"syncRegistrations": {
				strconv.FormatBool(ldap.SyncRegistrations),
			},
			"vendor": {
				ldap.Vendor,
			},
			"usernameLDAPAttribute": {
				ldap.UsernameLDAPAttribute,
			},
			"rdnLDAPAttribute": {
				ldap.RdnLDAPAttribute,
			},
			"uuidLDAPAttribute": {
				ldap.UuidLDAPAttribute,
			},
			"userObjectClasses": {
				strings.Join(ldap.UserObjectClasses, ", "),
			},
			"connectionUrl": {
				ldap.ConnectionUrl,
			},
			"usersDn": {
				ldap.UsersDn,
			},
			"authType": {
				ldap.AuthType,
			},
			"bindDn": {
				ldap.BindDn,
			},
			"bindCredential": {
				ldap.BindCredential,
			},
			"customUserSearchFilter": {
				ldap.CustomUserSearchFilter,
			},
			"searchScope": {
				ldap.SearchScope,
			},
			"validatePasswordPolicy": {
				strconv.FormatBool(ldap.ValidatePasswordPolicy),
			},
			"useTruststoreSpi": {
				ldap.UseTruststoreSpi,
			},
			"connectionTimeout": {
				strconv.Itoa(ldap.ConnectionTimeout),
			},
			"readTimeout": {
				strconv.Itoa(ldap.ReadTimeout),
			},
			"pagination": {
				strconv.FormatBool(ldap.Pagination),
			},
			"batchSizeForSync": {
				strconv.Itoa(ldap.BatchSizeForSync),
			},
			"fullSyncPeriod": {
				strconv.Itoa(ldap.FullSyncPeriod),
			},
			"changedSyncPeriod": {
				strconv.Itoa(ldap.ChangedSyncPeriod),
			},
		},
	}
}

func convertToLdapUserFederation(component *userFederationComponent) (*LdapUserFederation, error) {
	enabled, err := strconv.ParseBool(component.getConfig("enabled"))
	if err != nil {
		return nil, err
	}

	priority, err := strconv.Atoi(component.getConfig("priority"))
	if err != nil {
		return nil, err
	}

	importEnabled, err := strconv.ParseBool(component.getConfig("importEnabled"))
	if err != nil {
		return nil, err
	}

	syncRegistrations, err := strconv.ParseBool(component.getConfig("syncRegistrations"))
	if err != nil {
		return nil, err
	}

	userObjectClasses := strings.Split(component.getConfig("userObjectClasses"), ", ")

	validatePasswordPolicy, err := strconv.ParseBool(component.getConfig("validatePasswordPolicy"))
	if err != nil {
		return nil, err
	}

	connectionTimeout, err := strconv.Atoi(component.getConfig("connectionTimeout"))
	if err != nil {
		return nil, err
	}

	readTimeout, err := strconv.Atoi(component.getConfig("readTimeout"))
	if err != nil {
		return nil, err
	}

	pagination, err := strconv.ParseBool(component.getConfig("pagination"))
	if err != nil {
		return nil, err
	}

	batchSizeForSync, err := strconv.Atoi(component.getConfig("batchSizeForSync"))
	if err != nil {
		return nil, err
	}

	fullSyncPeriod, err := strconv.Atoi(component.getConfig("fullSyncPeriod"))
	if err != nil {
		return nil, err
	}

	changedSyncPeriod, err := strconv.Atoi(component.getConfig("changedSyncPeriod"))
	if err != nil {
		return nil, err
	}

	return &LdapUserFederation{
		Id:      component.Id,
		Name:    component.Name,
		RealmId: component.ParentId,

		Enabled:  enabled,
		Priority: priority,

		ImportEnabled:     importEnabled,
		EditMode:          component.getConfig("editMode"),
		SyncRegistrations: syncRegistrations,

		Vendor:                 component.getConfig("vendor"),
		UsernameLDAPAttribute:  component.getConfig("usernameLDAPAttribute"),
		RdnLDAPAttribute:       component.getConfig("rdnLDAPAttribute"),
		UuidLDAPAttribute:      component.getConfig("uuidLDAPAttribute"),
		UserObjectClasses:      userObjectClasses,
		ConnectionUrl:          component.getConfig("connectionUrl"),
		UsersDn:                component.getConfig("usersDn"),
		AuthType:               component.getConfig("authType"),
		BindDn:                 component.getConfig("bindDn"),
		BindCredential:         component.getConfig("bindCredential"),
		CustomUserSearchFilter: component.getConfig("customUserSearchFilter"),
		SearchScope:            component.getConfig("searchScope"),

		ValidatePasswordPolicy: validatePasswordPolicy,
		UseTruststoreSpi:       component.getConfig("useTruststoreSpi"),
		ConnectionTimeout:      connectionTimeout,
		ReadTimeout:            readTimeout,
		Pagination:             pagination,

		BatchSizeForSync:  batchSizeForSync,
		FullSyncPeriod:    fullSyncPeriod,
		ChangedSyncPeriod: changedSyncPeriod,

		CachePolicy: component.getConfig("cachePolicy"),
	}, nil
}
