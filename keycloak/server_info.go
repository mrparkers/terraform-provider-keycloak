package keycloak

type Theme struct {
	Name    string   `json:"name"`
	Locales []string `json:"locales,omitempty"`
}

type ServerInfo struct {
	Themes map[string][]Theme `json:"themes"`
}

func (serverInfo *ServerInfo) ThemeIsInstalled(t, themeName string) bool {
	if themes, ok := serverInfo.Themes[t]; ok {
		for _, theme := range themes {
			if theme.Name == themeName {
				return true
			}
		}
	}

	return false
}

func (keycloakClient *KeycloakClient) GetServerInfo() (*ServerInfo, error) {
	var serverInfo ServerInfo

	err := keycloakClient.get("/serverinfo", &serverInfo)
	if err != nil {
		return nil, err
	}

	return &serverInfo, nil
}
