// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

package oauthclient_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

const oauthClientId = "myOauthClient"

func TestAccOauthClient_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: oauthClient_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: oauthClient_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					oauthClient_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccOauthClient_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: oauthClient_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: oauthClient_MinimalHCL(),
				Check:  oauthClient_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  oauthClient_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: oauthClient_CompleteHCL(),
				Check:  oauthClient_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: oauthClient_MinimalHCL(),
				Check:  oauthClient_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: oauthClient_CompleteHCL(),
				Check:  oauthClient_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:            oauthClient_CompleteHCL(),
				ResourceName:      "pingfederate_oauth_client.example",
				ImportStateId:     oauthClientId,
				ImportState:       true,
				ImportStateVerify: true,
				// Secrets can't be imported and encrypted secrets change on read
				ImportStateVerifyIgnore: []string{
					"client_auth.encrypted_secret",
					"client_auth.secret",
					"client_auth.secondary_secrets",
				},
			},
		},
	})
}

// Minimal HCL with only required values set
func oauthClient_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_client" "example" {
  client_id   = "%s"
  grant_types = ["CLIENT_CREDENTIALS"]
  name        = "myclient"
  client_auth = {
    type   = "SECRET"
    secret = "mysecret"
  }
}
data "pingfederate_oauth_client" "example" {
  client_id = pingfederate_oauth_client.example.client_id
}
`, oauthClientId)
}

// Maximal HCL with all values set where possible
func oauthClient_CompleteHCL() string {
	var versionedHcl string
	if acctest.VersionAtLeast(version.PingFederate1210) {
		versionedHcl += `
		// HCL necessary to use refresh token rolling interval time unit
	refresh_token_rolling_interval_type = "OVERRIDE_SERVER_DEFAULT"
	refresh_token_rolling_interval = 10
		// PF 12.1 attributes
	refresh_token_rolling_interval_time_unit = "MINUTES"
	enable_cookieless_authentication_api = true
	require_offline_access_scope_to_issue_refresh_tokens = "YES"
	offline_access_require_consent_prompt = "YES"
		`
	}
	if acctest.VersionAtLeast(version.PingFederate1220) {
		versionedHcl += `
	lockout_max_malicious_actions = 500
	lockout_max_malicious_actions_type = "OVERRIDE_SERVER_DEFAULT"
		`
	}
	var versionedOidcPolicyHcl string
	if acctest.VersionAtLeast(version.PingFederate1200) {
		versionedOidcPolicyHcl += `
	post_logout_redirect_uris = ["https://example.com", "https://pingidentity.com"]
		`
	}
	if acctest.VersionAtLeast(version.PingFederate1220) {
		versionedOidcPolicyHcl += `
	user_info_response_content_encryption_algorithm = "AES_256_GCM"
	user_info_response_encryption_algorithm = "RSA_OAEP_256"
	user_info_response_signing_algorithm = "RS256"
		`
	}
	return fmt.Sprintf(`
resource "pingfederate_extended_properties" "example" {
  items = [
    {
      name         = "authNexp",
      description  = "Authentication Experience [Single_Factor | Internal | ID-First | Multi_Factor]",
      multi_valued = false
    },
    {
      name         = "useAuthnApi",
      description  = "Use the AuthN API",
      multi_valued = false
    },
    {
      name         = "test"
      description  = "test"
      multi_valued = false
    }
  ]
}

resource "pingfederate_oauth_client" "example" {
  depends_on                                   = [pingfederate_extended_properties.example]
  client_id                                    = "%s"
  allow_authentication_api_init                = true
  bypass_activation_code_confirmation_override = false
  bypass_approval_page                         = true
  ciba_delivery_mode                           = "PING"
  ciba_notification_endpoint                   = "https://example.com"
  ciba_polling_interval                        = 1
  ciba_request_object_signing_algorithm        = "RS256"
  ciba_require_signed_requests                 = true
  ciba_user_code_supported                     = true
  client_auth = {
    secondary_secrets = [{
      secret      = "examplesecondary"
      expiry_time = "2036-12-31T23:59:59Z"
    }]
    secret = "2FederateM0re"
    type   = "SECRET"
  }
  client_secret_retention_period      = 12
  client_secret_retention_period_type = "OVERRIDE_SERVER_DEFAULT"
  description                         = "updated client"
  enabled                             = false
  extended_parameters = {
    "test" = {
      "values" = ["test"]
    }
  }
  grant_types = ["IMPLICIT", "AUTHORIZATION_CODE", "RESOURCE_OWNER_CREDENTIALS",
    "REFRESH_TOKEN", "EXTENSION", "DEVICE_CODE",
  "ACCESS_TOKEN_VALIDATION", "CIBA", "TOKEN_EXCHANGE"]
  jwks_settings = {
    jwks_url = "https://example.com"
  }
  jwt_secured_authorization_response_mode_content_encryption_algorithm = "AES_128_CBC_HMAC_SHA_256"
  jwt_secured_authorization_response_mode_encryption_algorithm         = "RSA_OAEP"
  jwt_secured_authorization_response_mode_signing_algorithm            = "RS256"
  logo_url                                                             = "https://example.com/logo.png"
  name                                                                 = "my updated client"
  oidc_policy = {
    grant_access_session_revocation_api         = false
    grant_access_session_session_management_api = false
    id_token_content_encryption_algorithm       = "AES_128_CBC_HMAC_SHA_256"
    id_token_encryption_algorithm               = "A192GCMKW"
    id_token_signing_algorithm                  = "HS256"
    pairwise_identifier_user_type               = true
    ping_access_logout_capable                  = false
    sector_identifier_uri                       = "https://example.com"
    logout_mode                                 = "OIDC_BACK_CHANNEL"
    back_channel_logout_uri                     = "https://example.com"
	%s
  }
  persistent_grant_expiration_time                 = 5
  persistent_grant_expiration_time_unit            = "DAYS"
  persistent_grant_expiration_type                 = "OVERRIDE_SERVER_DEFAULT"
  persistent_grant_idle_timeout                    = 3
  persistent_grant_idle_timeout_time_unit          = "DAYS"
  persistent_grant_idle_timeout_type               = "OVERRIDE_SERVER_DEFAULT"
  persistent_grant_reuse_grant_types               = ["IMPLICIT"]
  persistent_grant_reuse_type                      = "OVERRIDE_SERVER_DEFAULT"
  redirect_uris                                    = ["https://example.com"]
  refresh_rolling                                  = "ROLL"
  refresh_token_rolling_grace_period               = 12
  refresh_token_rolling_grace_period_type          = "OVERRIDE_SERVER_DEFAULT"
  request_object_signing_algorithm                 = "RS256"
  require_jwt_secured_authorization_response_mode  = true
  require_proof_key_for_code_exchange              = true
  require_pushed_authorization_requests            = true
  require_signed_requests                          = true
  restrict_scopes                                  = true
  restrict_to_default_access_token_manager         = true
  restricted_response_types                        = ["code", "token", "id_token", "code token", "code id_token", "token id_token", "code token id_token"]
  restricted_scopes                                = ["openid"]
  token_introspection_content_encryption_algorithm = "AES_128_CBC_HMAC_SHA_256"
  token_introspection_encryption_algorithm         = "A128KW"
  token_introspection_signing_algorithm            = "RS256"
  validate_using_all_eligible_atms                 = true
  require_dpop                                     = true
  %s
}
data "pingfederate_oauth_client" "example" {
  client_id = pingfederate_oauth_client.example.client_id
}
`, oauthClientId, versionedOidcPolicyHcl, versionedHcl)
}

func oauthClient_CheckVersionedComputedValues() resource.TestCheckFunc {
	var versionedChecks []resource.TestCheckFunc
	if acctest.VersionAtLeast(version.PingFederate1200) {
		versionedChecks = append(versionedChecks,
			resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.post_logout_redirect_uris"))
	}
	if acctest.VersionAtLeast(version.PingFederate1210) {
		versionedChecks = append(versionedChecks,
			resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "refresh_token_rolling_interval_time_unit", "HOURS"),
			resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "enable_cookieless_authentication_api", "false"),
			resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "require_offline_access_scope_to_issue_refresh_tokens", "SERVER_DEFAULT"),
			resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "offline_access_require_consent_prompt", "SERVER_DEFAULT"),
		)
	}
	if acctest.VersionAtLeast(version.PingFederate1220) {
		versionedChecks = append(versionedChecks,
			resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "lockout_max_malicious_actions"),
			resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "lockout_max_malicious_actions_type", "SERVER_DEFAULT"),
			resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.user_info_response_content_encryption_algorithm"),
			resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.user_info_response_encryption_algorithm"),
			resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.user_info_response_signing_algorithm"))
	}
	return resource.ComposeTestCheckFunc(versionedChecks...)
}

// Validate any computed values when applying minimal HCL
func oauthClient_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "allow_authentication_api_init", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "authorization_detail_types.#", "0"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "bypass_activation_code_confirmation_override"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "bypass_approval_page", "false"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "ciba_delivery_mode"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "ciba_notification_endpoint"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "ciba_polling_interval"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "ciba_request_object_signing_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "ciba_require_signed_requests"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "ciba_user_code_supported"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.client_cert_issuer_dn"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.client_cert_subject_dn"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "client_auth.encrypted_secret"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.enforce_replay_prevention"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "client_auth.secondary_secrets.#", "0"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.token_endpoint_auth_signing_algorithm"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "client_secret_changed_time"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_secret_retention_period"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "client_secret_retention_period_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "creation_date"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "default_access_token_manager_ref"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "description"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "device_flow_setting_type", "SERVER_DEFAULT"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "device_polling_interval_override"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "enabled", "true"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "exclusive_scopes.#", "0"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "extended_parameters"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "id", oauthClientId),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "jwks_settings"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "jwt_secured_authorization_response_mode_content_encryption_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "jwt_secured_authorization_response_mode_encryption_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "jwt_secured_authorization_response_mode_signing_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "logo_url"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "modification_date"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "oidc_policy.grant_access_session_revocation_api", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "oidc_policy.grant_access_session_session_management_api", "false"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.id_token_content_encryption_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.id_token_encryption_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.id_token_signing_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.logout_uris"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "oidc_policy.pairwise_identifier_user_type", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "oidc_policy.ping_access_logout_capable", "false"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.sector_identifier_uri"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "pending_authorization_timeout_override"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_expiration_time", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_expiration_time_unit", "DAYS"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_expiration_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_idle_timeout", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_idle_timeout_time_unit", "DAYS"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_idle_timeout_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_reuse_grant_types.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "persistent_grant_reuse_type", "SERVER_DEFAULT"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "redirect_uris.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "refresh_rolling", "SERVER_DEFAULT"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "refresh_token_rolling_grace_period"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "refresh_token_rolling_grace_period_type", "SERVER_DEFAULT"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "refresh_token_rolling_interval"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "refresh_token_rolling_interval_type", "SERVER_DEFAULT"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "request_object_signing_algorithm"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "require_jwt_secured_authorization_response_mode", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "require_proof_key_for_code_exchange", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "require_pushed_authorization_requests", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "require_signed_requests", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "restrict_scopes", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "restrict_to_default_access_token_manager", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "restricted_response_types.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "restricted_scopes.#", "0"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "token_introspection_content_encryption_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "token_introspection_encryption_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "token_introspection_signing_algorithm"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "user_authorization_url_override"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "validate_using_all_eligible_atms", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "require_dpop", "false"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "oidc_policy.logout_mode", "NONE"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.back_channel_logout_uri"),
		oauthClient_CheckVersionedComputedValues(),
	)
}

// Validate any computed values when applying complete HCL
func oauthClient_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "authorization_detail_types.#", "0"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.client_cert_issuer_dn"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.client_cert_subject_dn"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "client_auth.encrypted_secret"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.enforce_replay_prevention"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "client_auth.secondary_secrets.0.encrypted_secret"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "client_auth.token_endpoint_auth_signing_algorithm"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "client_secret_changed_time"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "creation_date"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "default_access_token_manager_ref"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "device_flow_setting_type", "SERVER_DEFAULT"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "device_polling_interval_override"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "exclusive_scopes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_oauth_client.example", "id", oauthClientId),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "jwks_settings.jwks"),
		resource.TestCheckResourceAttrSet("pingfederate_oauth_client.example", "modification_date"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "oidc_policy.logout_uris"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "pending_authorization_timeout_override"),
		resource.TestCheckNoResourceAttr("pingfederate_oauth_client.example", "user_authorization_url_override"),
	)
}

// Delete the resource
func oauthClient_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.OauthClientsAPI.DeleteOauthClient(acctest.TestBasicAuthContext(), oauthClientId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func oauthClient_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.OauthClientsAPI.DeleteOauthClient(acctest.TestBasicAuthContext(), oauthClientId).Execute()
	if err == nil {
		return fmt.Errorf("oauth_client still exists after tests. Expected it to be destroyed")
	}
	return nil
}
