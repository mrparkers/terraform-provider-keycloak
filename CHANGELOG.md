## 1.20.0 (July 20, 2020)

FEATURES:

- new resource: `keycloak_user_roles` ([#315](https://github.com/mrparkers/terraform-provider-keycloak/pull/315))
- new resource: `keycloak_identity_provider_token_exchange_scope_permission` ([#318](https://github.com/mrparkers/terraform-provider-keycloak/pull/318))
- new resources: `keycloak_saml_client_scope`, `keycloak_saml_client_default_scopes` ([#320](https://github.com/mrparkers/terraform-provider-keycloak/pull/320))

IMPROVEMENTS:

- adds `default_signature_algorithm` attribute for `keycloak_realm` resource ([#282](https://github.com/mrparkers/terraform-provider-keycloak/pull/282))
- adds `parent_id` attribute to `keycloak_custom_user_federation` resource ([#325](https://github.com/mrparkers/terraform-provider-keycloak/pull/325))
- adds `extra_config` attribute to identity provider mapper resources ([#316](https://github.com/mrparkers/terraform-provider-keycloak/pull/316))
- adds `include_in_token_scope` and `gui_order` attributes to `keycloak_openid_client_scope` resource ([#320](https://github.com/mrparkers/terraform-provider-keycloak/pull/320))
- adds `base_path` provider attribute, improve login error messages ([#332](https://github.com/mrparkers/terraform-provider-keycloak/pull/332))
- adds encryption attributes to `keycloak_saml_client` resource ([#342](https://github.com/mrparkers/terraform-provider-keycloak/pull/342))
- adds `signature_algorithm` attribute to `keycloak_saml_client` resource ([#345](https://github.com/mrparkers/terraform-provider-keycloak/pull/345))

BUG FIXES:

- fix import for `keycloak_openid_client_service_account_role` resource ([#314](https://github.com/mrparkers/terraform-provider-keycloak/pull/314))
- fix realm role support for `keycloak_generic_client_role_mapper` resource ([#316](https://github.com/mrparkers/terraform-provider-keycloak/pull/316))
- fix `keycloak_group` data source to support nested groups ([#334](https://github.com/mrparkers/terraform-provider-keycloak/pull/334))
- fix `keycloak_group` data source / resource to support group names with backslash character ([#337](https://github.com/mrparkers/terraform-provider-keycloak/pull/337))

## 1.19.0 (June 5, 2020)

FEATURES:

- new resource: `keycloak_openid_user_client_role_protocol_mapper` ([#299](https://github.com/mrparkers/terraform-provider-keycloak/pull/299))
- new resource: `keycloak_openid_user_session_note_protocol_mapper` ([#309](https://github.com/mrparkers/terraform-provider-keycloak/pull/309))

IMPROVEMENTS:

- add `login_theme` attribute to `keycloak_openid_client` resource ([#278](https://github.com/mrparkers/terraform-provider-keycloak/pull/278))
- add `aggregate_attributes` attribute to `keycloak_openid_user_attribute_protocol_mapper` resource ([#272](https://github.com/mrparkers/terraform-provider-keycloak/pull/272))
- add `user_managed_access` attribute to `keycloak_realm` resource ([#275](https://github.com/mrparkers/terraform-provider-keycloak/pull/275))
- support deployed JavaScript policies for `keycloak_openid_client_js_policy` resource ([#275](https://github.com/mrparkers/terraform-provider-keycloak/pull/275))
- add `internal_id` computed attribute to `keycloak_realm` resource and data source ([#270](https://github.com/mrparkers/terraform-provider-keycloak/pull/270))
- surface Keycloak API errors to users during `terraform plan` and `terraform apply` ([#304](https://github.com/mrparkers/terraform-provider-keycloak/pull/304))
- add `kerberos` configuration for `keycloak_ldap_user_federation` resource ([#290](https://github.com/mrparkers/terraform-provider-keycloak/pull/290))
- test all major versions of Keycloak in CI ([#294](https://github.com/mrparkers/terraform-provider-keycloak/pull/294))
- add import support for `keycloak_generic_client_role_mapper` resource ([#310](https://github.com/mrparkers/terraform-provider-keycloak/pull/310))
- use terraform-plugin-sdk user agent string in http client ([#311](https://github.com/mrparkers/terraform-provider-keycloak/pull/311))

BUG FIXES:

- fix: mark `group_id` attribute as required for `keycloak_group_roles` resource ([#292](https://github.com/mrparkers/terraform-provider-keycloak/pull/292))

## 1.18.0 (April 17, 2020)

FEATURES:

* new resource: `keycloak_ldap_hardcoded_group_mapper` ([#264](https://github.com/mrparkers/terraform-provider-keycloak/pull/264))
* new data source: `keycloak_saml_client_installation_provider` ([#263](https://github.com/mrparkers/terraform-provider-keycloak/pull/263))
* new resource: `keycloak_ldap_role_mapper` ([#265](https://github.com/mrparkers/terraform-provider-keycloak/pull/265))

IMPROVEMENTS:

* add `tls_insecure_skip_verify` provider attribute ([#237](https://github.com/mrparkers/terraform-provider-keycloak/pull/237))
* add `client_scope_id` attribute to `keycloak_generic_client_role_mapper` resource ([#253](https://github.com/mrparkers/terraform-provider-keycloak/pull/253))
* add `email_verified` attribute to `keycloak_user` resource ([#256](https://github.com/mrparkers/terraform-provider-keycloak/pull/256))
* add `JSON` as a valid `claim_value_type` for openid protocol mapper resources ([#260](https://github.com/mrparkers/terraform-provider-keycloak/pull/260))
* add `force_name_id_format` attribute to `keycloak_saml_client` resource ([#261](https://github.com/mrparkers/terraform-provider-keycloak/pull/261))
* add `consent_required` and `authentication_flow_binding_overrides` attributes for `keycloak_openid_client` resource ([#262](https://github.com/mrparkers/terraform-provider-keycloak/pull/262))
* add `root_url` attribute to `keycloak_openid_client` resource ([#248](https://github.com/mrparkers/terraform-provider-keycloak/pull/248))
* add federated identity support for `keycloak_user` resource ([#274](https://github.com/mrparkers/terraform-provider-keycloak/pull/274))

BUG FIXES:

* correctly handle manually deleted clients when refreshing a `keycloak_openid_client_default_scopes` resource ([#252](https://github.com/mrparkers/terraform-provider-keycloak/pull/252))
* correctly handle manually deleted clients when refreshing a `keycloak_openid_client_optional_scopes` resource

## 1.17.1 (March 12, 2020)

BUG FIXES:

* fix: allow `defaultScope` and `acceptsPromptNoneForwardFromClient` attributes to be set for `keycloak_oidc_identity_provider`. Previously, these attributes could only be set via `extra_config`, which stopped working as of v1.17.0. This release introduces these attributes as top-level attributes for the `keycloak_oidc_identity_provider` resource.

## 1.17.0 (March 10, 2020)

FEATURES:

* new resources: `keycloak_authentication_flow`, `keycloak_authentication_subflow`, `keycloak_authentication_execution` ([#215](https://github.com/mrparkers/terraform-provider-keycloak/pull/215))
* new resource: `keycloak_authentication_execution_config` ([#241](https://github.com/mrparkers/terraform-provider-keycloak/pull/241))
* new resource: `keycloak_oidc_google_identity_provider` ([#240](https://github.com/mrparkers/terraform-provider-keycloak/pull/240))
* new resource: `keycloak_ldap_msad_user_account_control_mapper` ([#244](https://github.com/mrparkers/terraform-provider-keycloak/pull/244))
* new resources: `keycloak_openid_client_group_policy`, `keycloak_openid_client_role_policy`, `keycloak_openid_client_aggregate_policy`, `keycloak_openid_client_js_policy`, `keycloak_openid_client_time_policy`, `keycloak_openid_client_user_policy`, `keycloak_openid_client_client_policy` ([#246](https://github.com/mrparkers/terraform-provider-keycloak/pull/246))
* new resource: `keycloak_generic_client_role_mapper` ([#242](https://github.com/mrparkers/terraform-provider-keycloak/pull/242))

IMPROVEMENTS:

* add `client_scope_id` attribute to `keycloak_generic_client_protocol_mapper` resource ([#229](https://github.com/mrparkers/terraform-provider-keycloak/pull/229))
* add `root_ca_certificate` attribute to provider config ([#227](https://github.com/mrparkers/terraform-provider-keycloak/pull/227))
* add `scopes` attribute to `keycloak_openid_client_authorization_permission` resource ([#220](https://github.com/mrparkers/terraform-provider-keycloak/pull/220))
* add `access_token_lifespan` attribute to `keycloak_openid_client` resource ([#233](https://github.com/mrparkers/terraform-provider-keycloak/pull/233))

## 1.16.0 (February 13, 2020)

FEATURES:

* new resource: `keycloak_realm_events` ([#211](https://github.com/mrparkers/terraform-provider-keycloak/pull/211))
* new resource: `resource_keycloak_openid_client_service_account_role` ([#202](https://github.com/mrparkers/terraform-provider-keycloak/pull/202))

IMPROVEMENTS:

* add base_url attribute to `keycloak_openid_client` resource ([#201](https://github.com/mrparkers/terraform-provider-keycloak/pull/201))
* allow configuration of the client timeout by an environment variable ([#206](https://github.com/mrparkers/terraform-provider-keycloak/pull/206))
* adds consent_required attribute to `keycloak_openid_client` resource ([#207](https://github.com/mrparkers/terraform-provider-keycloak/pull/207))
* adds admin_url attribute to `keycloak_openid_client` resource ([#203](https://github.com/mrparkers/terraform-provider-keycloak/pull/203))
* adds display_name_html attribute to `keycloak_realm` resource and data source ([#209](https://github.com/mrparkers/terraform-provider-keycloak/pull/209))
* switch to terraform-plugin-sdk ([#214](https://github.com/mrparkers/terraform-provider-keycloak/pull/214))

BUG FIXES:

* URL encode role names to allow for special characters ([#213](https://github.com/mrparkers/terraform-provider-keycloak/pull/213))

## 1.15.0 (January 20, 2020)

FEATURES:

* new resource: `keycloak_ldap_hardcoded_role_mapper` ([#195](https://github.com/mrparkers/terraform-provider-keycloak/pull/195))

IMPROVEMENTS:

* add `full_scope_allowed` attribute to `keycloak_openid_client` resource ([#193](https://github.com/mrparkers/terraform-provider-keycloak/pull/193))
* add `exclude_session_state_from_auth_response` attribute to `keycloak_openid_client` resource ([#191](https://github.com/mrparkers/terraform-provider-keycloak/pull/191))
* allow empty value for `pkce_code_challenge_method` attribute on `keycloak_openid_client` resource ([#198](https://github.com/mrparkers/terraform-provider-keycloak/pull/198))
* support attributes for `keycloak_group` resource ([#199](https://github.com/mrparkers/terraform-provider-keycloak/pull/199))


## 1.14.0 (December 18, 2019)

FEATURES:

* add `keycloak_openid_client_service_account_user` data source ([#181](https://github.com/mrparkers/terraform-provider-keycloak/pull/181))
* add `keycloak_group` data source ([#185](https://github.com/mrparkers/terraform-provider-keycloak/pull/185))

IMPROVEMENTS:

* support Keycloak v8.0.0 ([#183](https://github.com/mrparkers/terraform-provider-keycloak/pull/183))
* new functionality for `keycloak_realm`: brute_force_detection, ssl_required, and custom attributes ([#183](https://github.com/mrparkers/terraform-provider-keycloak/pull/183))
* allow you to prevent refresh token reuse with a new `revoke_refresh_token` attribute for the `keycloak_realm` resource ([#183](https://github.com/mrparkers/terraform-provider-keycloak/pull/183))
  * **note: please refer to the [docs](https://mrparkers.github.io/terraform-provider-keycloak/resources/keycloak_realm/#tokens) for the new configuration values if you currently use `refresh_token_max_reuse`**

## 1.13.2 (November 27, 2019)

BUG FIXES:

* fix: update Keycloak API call to handle groups with more than 100 members ([#179](https://github.com/mrparkers/terraform-provider-keycloak/pull/179))

## 1.13.1 (November 19, 2019)

BUG FIXES:

* correctly handle Keycloak role names that contain a forward slash ([#175](https://github.com/mrparkers/terraform-provider-keycloak/pull/175))

## 1.13.0 (November 8, 2019)

IMPROVEMENTS:

* use cookiejar for Keycloak API requests ([#173](https://github.com/mrparkers/terraform-provider-keycloak/pull/173))

## 1.12.0 (November 4, 2019)

IMPROVEMENTS:

* add `pkce_code_challenge_method` attribute for `keycloak_openid_client` resource ([#170](https://github.com/mrparkers/terraform-provider-keycloak/pull/170))

BUG FIXES:

* always use valid client secret for `keycloak_oidc_identity_provider` resource ([#171](https://github.com/mrparkers/terraform-provider-keycloak/pull/171))
* fix state issues for `keycloak_openid_client_service_account_role` resource ([#171](https://github.com/mrparkers/terraform-provider-keycloak/pull/171))

## 1.11.1 (October 17, 2019)

BUG FIXES:

* fix required attribute for `keycloak_realm` data source ([#166](https://github.com/mrparkers/terraform-provider-keycloak/pull/166))
* automatically retry role deletion if the first attempt fails ([#168](https://github.com/mrparkers/terraform-provider-keycloak/pull/168))

## 1.11.0 (October 14, 2019)

FEATURES:

* new resource: `keycloak_openid_user_realm_role_protocol_mapper` ([#159](https://github.com/mrparkers/terraform-provider-keycloak/pull/159))
* new data source: `keycloak_realm` ([#160](https://github.com/mrparkers/terraform-provider-keycloak/pull/160))

IMPROVEMENTS:

* added `timeout` provider attribute ([#155](https://github.com/mrparkers/terraform-provider-keycloak/pull/155))
* always export `serviceAccountId` for `keycloak_openid_client` resource ([#162](https://github.com/mrparkers/terraform-provider-keycloak/pull/162))

BUG FIXES:

* fix default value for `reset_credentials_flow` attribute in `keycloak_realm` resource ([#158](https://github.com/mrparkers/terraform-provider-keycloak/pull/158))

## 1.10.0 (September 6, 2019)

note: this release contains a [bug](https://github.com/mrparkers/terraform-provider-keycloak/issues/156) in the `keycloak_realm` resource that incorrectly sets the default attribute for `reset_credentials_flow` to `"registration"`. Please ensure that you set this attribute manually to override the incorrect default until a future release fixes this issue.

FEATURES:

* new resource: `keycloak_required_action` ([#131](https://github.com/mrparkers/terraform-provider-keycloak/pull/131))
* new resource: `keycloak_default_groups` ([#146](https://github.com/mrparkers/terraform-provider-keycloak/pull/146))
* new resources: `keycloak_role`, `keycloak_group_roles`, `keycloak_openid_hardcoded_role_protocol_mapper` ([#143](https://github.com/mrparkers/terraform-provider-keycloak/pull/143))
* new data source: `keycloak_role` ([#143](https://github.com/mrparkers/terraform-provider-keycloak/pull/143))

IMPROVEMENTS:

* add `security_defences` attribute to `keycloak_realm` resource ([#130](https://github.com/mrparkers/terraform-provider-keycloak/pull/130))
* support custom config for `keycloak_custom_user_federation` resource ([#134](https://github.com/mrparkers/terraform-provider-keycloak/pull/134))
* add `initial_login` provider attribute to optionally avoid requests during provider setup ([#136](https://github.com/mrparkers/terraform-provider-keycloak/pull/136))
* support custom config for `keycloak_oidc_identity_provider` resource ([#137](https://github.com/mrparkers/terraform-provider-keycloak/pull/137))
* add `password_policy` attribute for `keycloak_realm` resource ([#139](https://github.com/mrparkers/terraform-provider-keycloak/pull/139))
* add flow binding attributes for `keycloak_realm` resource ([#140](https://github.com/mrparkers/terraform-provider-keycloak/pull/140))

BUG FIXES:

* fix user attributes to handle attributes longer than 255 characters ([#132](https://github.com/mrparkers/terraform-provider-keycloak/pull/132))
* fix import for `keycloak_oidc_identity_provider` ([#142](https://github.com/mrparkers/terraform-provider-keycloak/pull/142))

## 1.9.0 (June 20, 2019)

FEATURES:

* add `full_scope_allowed` attribute to `keycloak_saml_client` resource ([#118](https://github.com/mrparkers/terraform-provider-keycloak/pull/118))
* add `internationalization` attribute to `keycloak_realm` resource ([#124](https://github.com/mrparkers/terraform-provider-keycloak/pull/124))
* add `smtp_server` attribute to `keycloak_realm` resource ([#122](https://github.com/mrparkers/terraform-provider-keycloak/pull/122))

IMPROVEMENTS:

* allow the provider to use a confidential client with the password grant ([#114](https://github.com/mrparkers/terraform-provider-keycloak/pull/114))
* update Terraform SDK to 0.12.1 ([#120](https://github.com/mrparkers/terraform-provider-keycloak/pull/120))
* bump dependency versions for custom user federation example ([#121](https://github.com/mrparkers/terraform-provider-keycloak/pull/121))
* add static binary to release for use within Alpine Docker iamges ([#129](https://github.com/mrparkers/terraform-provider-keycloak/pull/129))

## 1.8.0 (May 14, 2019)

FEATURES:

* new resources: `keycloak_openid_client_authorization_resource`, `keycloak_openid_client_authorization_scope`, `keycloak_openid_client_authorization_permission`, `keycloak_openid_client_service_account_role` ([#104](https://github.com/mrparkers/terraform-provider-keycloak/pull/104))
  - note: docs for these resources will be released at a later date. for now, please refer to the source files.
* new data sources: keycloak_openid_client, keycloak_openid_client_authorization_policy ([#104](https://github.com/mrparkers/terraform-provider-keycloak/pull/104))
  - note: docs for these data sources will be released at a later date. for now, please refer to the source files.

IMPROVEMENTS:

* chore: update provider SDK to 0.12 ([#107](https://github.com/mrparkers/terraform-provider-keycloak/pull/107))
* chore: support Keycloak v6.0.1 ([#106](https://github.com/mrparkers/terraform-provider-keycloak/pull/106))
* chore: renames provider resource/data files ([#105](https://github.com/mrparkers/terraform-provider-keycloak/pull/105))

## 1.7.0 (April 18, 2019)

FEATURES:

* new resources: `keycloak_identity_provider` and mappers ([#92](https://github.com/mrparkers/terraform-provider-keycloak/pull/92))
  - note: docs for these resources will be released at a later date. for now, please refer to the source files.

IMPROVEMENTS:

* new attributes added for `keycloak_saml_client` resource ([#103](https://github.com/mrparkers/terraform-provider-keycloak/pull/103))

## 1.6.0 (March 6, 2019)

FEATURES:

* new resource: `keycloak_openid_client_optional_scopes` ([#96](https://github.com/mrparkers/terraform-provider-keycloak/pull/96))
* new resource: `keycloak_openid_audience_protocol_mapper` ([#97](https://github.com/mrparkers/terraform-provider-keycloak/pull/97))

## 1.5.0 (February 22, 2019)

FEATURES:

* adds support for non-master realms and resource owner password grant for Keycloak authentication ([#88](https://github.com/mrparkers/terraform-provider-keycloak/pull/88))

IMPROVEMENTS:

* support Keycloak v4.8.3.Final and Terraform v0.11.11 ([#93](https://github.com/mrparkers/terraform-provider-keycloak/pull/93))

BUG FIXES:

* handle 404 errors when reading a group for group memberships ([#95](https://github.com/mrparkers/terraform-provider-keycloak/pull/95))

## 1.4.0 (January 28, 2019)

FEATURES:

* new resource: `keycloak_saml_user_property_protocol_mapper` ([#85](https://github.com/mrparkers/terraform-provider-keycloak/pull/85))

## 1.3.0 (January 25, 2019)

FEATURES:

* new resource: `keycloak_saml_user_attribute_protocol_mapper` ([#84](https://github.com/mrparkers/terraform-provider-keycloak/pull/84))

## 1.2.0 (January 24, 2019)

FEATURES:

* new resource: `keycloak_saml_client` ([#82](https://github.com/mrparkers/terraform-provider-keycloak/pull/82))

IMPROVEMENTS:

* add validation for usernames to ensure they are always lowercase ([#83](https://github.com/mrparkers/terraform-provider-keycloak/pull/83))

## 1.1.0 (January 7, 2019)

IMPROVEMENTS:

* openid_client: add web_origins attribute ([#81](https://github.com/mrparkers/terraform-provider-keycloak/pull/81))
* user: add initial_password attribute ([#77](https://github.com/mrparkers/terraform-provider-keycloak/pull/77))

BUG FIXES:

* ldap mappers: don't assume component fields are returned by Keycloak API ([#80](https://github.com/mrparkers/terraform-provider-keycloak/pull/80))

## 1.0.0 (December 16, 2018)

Initial Release!

Docs: https://mrparkers.github.io/terraform-provider-keycloak
