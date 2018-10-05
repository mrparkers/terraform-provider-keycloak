package com.github.mrparkers.keycloak

import org.keycloak.component.ComponentModel
import org.keycloak.models.KeycloakSession
import org.keycloak.storage.UserStorageProviderFactory

class CustomUserStorageProviderFactory : UserStorageProviderFactory<CustomUserStorageProvider> {
    override fun getId(): String = "custom"

    override fun init(config: org.keycloak.Config.Scope) {
        super.init(config)
    }

    override fun create(session: KeycloakSession, model: ComponentModel): CustomUserStorageProvider =
            CustomUserStorageProvider(session, model)
}
