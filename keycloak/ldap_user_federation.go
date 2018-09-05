package keycloak

type userFederationComponent struct {
	Id           string              `json:"id,omitempty"`
	Name         string              `json:"name"`
	ProviderId   string              `json:"providerId"`   // ldap
	ProviderType string              `json:"providerType"` // org.keycloak.storage.UserStorageProvider
	ParentId     string              `json:"parentId"`     // realm id
	Config       map[string][]string `json:"config"`       // generic interface, but always includes "cachePolicy", "enabled", and "priority"
}

type LdapUserFederation struct {
	Id      string
	Name    string
	RealmId string

	Enabled  bool
	Priority int

	ImportEnabled     bool
	EditMode          string // can be "READ_ONLY", "WRITABLE", or "UNSYNCED"
	SyncRegistrations bool

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
