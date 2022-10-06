package com.github.mrparkers.keycloak

import org.keycloak.component.ComponentModel
import org.keycloak.credential.CredentialInput
import org.keycloak.credential.CredentialInputUpdater
import org.keycloak.credential.CredentialInputValidator
import org.keycloak.credential.CredentialModel
import org.keycloak.credential.LegacyUserCredentialManager
import org.keycloak.models.*
import org.keycloak.storage.ReadOnlyException
import org.keycloak.storage.StorageId
import org.keycloak.storage.UserStorageProvider
import org.keycloak.storage.adapter.AbstractUserAdapter
import org.keycloak.storage.user.UserLookupProvider
import java.util.*

class CustomUserStorageProvider(private val session: KeycloakSession, private val model: ComponentModel) :
        UserStorageProvider, UserLookupProvider, CredentialInputValidator, CredentialInputUpdater {

    private val loadedUsers: MutableMap<String, UserModel> = HashMap()
    private val users = mapOf(
            "tester" to "password"
    )

    // UserStorageProvider

    override fun close() {

    }

    // UserLookupProvider

    override fun getUserByEmail(email: String, realm: RealmModel): UserModel? {
        return null
    }

    override fun getUserByUsername(username: String, realm: RealmModel): UserModel? {
        val user = loadedUsers[username]

        if (user != null) {
            return user
        }

        if (users.containsKey(username)) {
            val newUser = object : AbstractUserAdapter(session, realm, model) {
                override fun getUsername(): String {
                    return username
                }

				override fun credentialManager(): SubjectCredentialManager {
					return LegacyUserCredentialManager(session, realm, this)
				}
			}

            loadedUsers[username] = newUser

            return newUser
        }

        return null
    }

    override fun getUserById(id: String, realm: RealmModel): UserModel? {
        val storageId = StorageId(id)
        val username = storageId.externalId

        return getUserByUsername(username, realm)
    }

    // CredentialInputValidator

    override fun isConfiguredFor(realm: RealmModel, user: UserModel, credentialType: String): Boolean {
        return supportsCredentialType(credentialType)
    }

    override fun supportsCredentialType(credentialType: String?): Boolean {
        return credentialType.equals(CredentialModel.PASSWORD)
    }

    override fun isValid(realm: RealmModel, user: UserModel, input: CredentialInput): Boolean {
        if (!supportsCredentialType(input.type) || input !is UserCredentialModel) {
            return false
        }

        val password = users[user.username] ?: return false

        return password == input.value
    }

    // CredentialInputUpdater

    override fun getDisableableCredentialTypes(realm: RealmModel, user: UserModel): MutableSet<String> {
        return Collections.EMPTY_SET as MutableSet<String>
    }

    override fun updateCredential(realm: RealmModel, user: UserModel, input: CredentialInput): Boolean {
        if (input.type == CredentialModel.PASSWORD) {
            throw ReadOnlyException("Custom provider does not support password updating")
        }

        return false
    }

    override fun disableCredentialType(realm: RealmModel, user: UserModel, credentialType: String) {

    }
}
