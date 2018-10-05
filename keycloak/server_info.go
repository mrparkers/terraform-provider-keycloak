package keycloak

type ComponentType struct {
	Id string `json:"id"`
}

type Theme struct {
	Name    string   `json:"name"`
	Locales []string `json:"locales,omitempty"`
}

type ServerInfo struct {
	Themes         map[string][]Theme         `json:"themes"`
	ComponentTypes map[string][]ComponentType `json:"componentTypes"`
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

func (serverInfo *ServerInfo) ComponentTypeIsInstalled(componentType, componentTypeId string) bool {
	if componentTypes, ok := serverInfo.ComponentTypes[componentType]; ok {
		for _, componentType := range componentTypes {
			if componentType.Id == componentTypeId {
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
