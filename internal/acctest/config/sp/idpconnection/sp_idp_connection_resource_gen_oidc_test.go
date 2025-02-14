// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package resource_sp_idp_connection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/accesstokenmanager"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

const idpConnOidcId = "oidcconn"

func TestAccSpIdpConnection_OidcMinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spIdpConnection_OidcCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spIdpConnection_OidcMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesOidcMinimal(),
			},
			{
				// Delete the minimal model
				Config:  spIdpConnection_OidcMinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: spIdpConnection_OidcCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesOidcComplete(),
			},
			{
				// Back to minimal model
				Config: spIdpConnection_OidcMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesOidcMinimal(),
			},
			{
				// Back to complete model
				Config: spIdpConnection_OidcCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesOidcComplete(),
			},
			{
				// Test importing the resource
				Config:            spIdpConnection_OidcCompleteHCL(),
				ResourceName:      "pingfederate_sp_idp_connection.example",
				ImportStateId:     idpConnOidcId,
				ImportState:       true,
				ImportStateVerify: true,
				// client_secret won't be returned by the API.
				// A couple boolean attributes also are not returned by the API when set to false.
				ImportStateVerifyIgnore: []string{
					"idp_browser_sso.assertions_signed",
					"idp_browser_sso.sign_authn_requests",
					"oidc_client_credentials.client_secret",
				},
			},
		},
	})
}

func spIdpConnection_OidcDependencyHCL() string {
	return `
resource "pingfederate_authentication_policy_contract" "apc1" {
  contract_id = "sp_idp1"
  name        = "Example sp_idp1"
  extended_attributes = [
    {
      name = "directory_id"
    },
    {
      name = "given_name"
    },
    {
      name = "family_name"
    },
    {
      name = "email"
    }
  ]
}
  `
}

// Minimal HCL with only required values set
func spIdpConnection_OidcMinimalHCL() string {
	return fmt.Sprintf(`
%s

resource "pingfederate_sp_idp_connection" "example" {
  entity_id = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as"

  name          = "PingOne"
  connection_id = "%s"

  oidc_client_credentials = {
    client_id = "myclientid"
  }

  idp_browser_sso = {
    attribute_contract = {
      extended_attributes = [
        { name = "acr" },
        { name = "address" },
        { name = "auth_time" },
        { name = "email" },
        { name = "family_name" },
        { name = "given_name" },
        { name = "iss" },
        { name = "locale" },
        { name = "middle_name" },
        { name = "name" },
        { name = "phone_number" },
        { name = "picture" },
        { name = "preferred_username" },
        { name = "profile" },
      ]
    }
    authentication_policy_contract_mappings = [
      {
        authentication_policy_contract_ref = {
          id = pingfederate_authentication_policy_contract.apc1.id
        }

        attribute_contract_fulfillment = {
          directory_id = {
            source = {
              type = "CLAIMS"
            }
            value = "sub"
          }
          email = {
            source = {
              type = "CLAIMS"
            }
            value = "email"
          }
          family_name = {
            source = {
              type = "CLAIMS"
            }
            value = "family_name"
          }
          given_name = {
            source = {
              type = "CLAIMS"
            }
            value = "given_name"
          }
          subject = {
            source = {
              type = "CLAIMS"
            }
            value = "sub"
          }
        }
      },
    ]

    idp_identity_mapping = "ACCOUNT_MAPPING"

    oidc_provider_settings = {
      authorization_endpoint = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/authorize"
      jwks_url               = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/jwks"
      login_type             = "POST"
      scopes                 = "openid profile email address phone"
    }
    protocol = "OIDC"
  }
}
`, spIdpConnection_OidcDependencyHCL(), idpConnOidcId)
}

// Maximal HCL with all values set where possible
func spIdpConnection_OidcCompleteHCL() string {
	versionedOidcProviderSettingsHcl := ""
	if acctest.VersionAtLeast(version.PingFederate1200) {
		versionedOidcProviderSettingsHcl += `
      logout_endpoint                              = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/signoff"
    `
	}
	if acctest.VersionAtLeast(version.PingFederate1210) {
		versionedOidcProviderSettingsHcl += `
      jwt_secured_authorization_response_mode_type = "DISABLED"
    `
	}

	return fmt.Sprintf(`
%s

%s

resource "pingfederate_sp_idp_connection" "example" {
  active = true
  additional_allowed_entities_configuration = {
    additional_allowed_entities = [
      {
        entity_id          = "https://bxretail.org",
        entity_description = "additional entity"
      }
    ]
    allow_additional_entities = true
    allow_all_entities        = false
  }
  base_url      = "https://bxretail.org"
  connection_id = "%s"
  contact_info = {
    company    = "Ping Identity"
    email      = "test@test.com"
    first_name = "test"
    last_name  = "test"
    phone      = "555-5555"
  }
  entity_id         = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as"
  error_page_msg_id = "errorDetail.spSsoFailure"
  extended_properties = {
    authNexp = {
      values = ["val1"]
    }
    useAuthnApi = {
      values = ["val2"]
    }
  }
  idp_browser_sso = {
    adapter_mappings = [
      {
        attribute_sources = [
          {
            ldap_attribute_source = {
              attribute_contract_fulfillment = null
              base_dn                        = "ou=Applications,ou=Ping,ou=Groups,dc=dm,dc=example,dc=com"
              binary_attribute_settings      = null
              data_store_ref = {
                id = "pingdirectory"
              }
              description            = "PingDirectory"
              member_of_nested_group = false
              search_attributes      = ["Subject DN"]
              search_filter          = "(&(memberUid=uid)(cn=Postman))"
              search_scope           = "SUBTREE"
              type                   = "LDAP"
            }
          },
        ]
        attribute_contract_fulfillment = {
          subject = {
            source = {
              type = "NO_MAPPING"
            }
            value = ""
          }
        }
        issuance_criteria = {
          conditional_criteria = [
            {
              attribute_name = "subject"
              condition      = "MULTIVALUE_CONTAINS_DN"
              source = {
                type = "MAPPED_ATTRIBUTES"
              }
              value = "cn=Example,dc=example,dc=com"
            },
          ]
        }
        restrict_virtual_entity_ids   = false
        restricted_virtual_entity_ids = []
        sp_adapter_ref = {
          id = "spadapter"
        }
      }
    ]
    always_sign_artifact_response = false
    attribute_contract = {
      extended_attributes = [
        {
          masked = false
          name   = "acr"
        },
        {
          masked = false
          name   = "address"
        },
        {
          masked = false
          name   = "auth_time"
        },
        {
          masked = false
          name   = "email"
        },
        {
          masked = false
          name   = "family_name"
        },
        {
          masked = false
          name   = "given_name"
        },
        {
          masked = false
          name   = "iss"
        },
        {
          masked = false
          name   = "locale"
        },
        {
          masked = false
          name   = "middle_name"
        },
        {
          masked = false
          name   = "name"
        },
        {
          masked = false
          name   = "phone_number"
        },
        {
          masked = false
          name   = "picture"
        },
        {
          masked = false
          name   = "preferred_username"
        },
        {
          masked = false
          name   = "profile"
        },
      ]
    }
    authentication_policy_contract_mappings = [
      {
        attribute_contract_fulfillment = {
          directory_id = {
            source = {
              id   = null
              type = "CLAIMS"
            }
            value = "sub"
          }
          email = {
            source = {
              id   = null
              type = "CLAIMS"
            }
            value = "email"
          }
          family_name = {
            source = {
              id   = null
              type = "CLAIMS"
            }
            value = "family_name"
          }
          given_name = {
            source = {
              id   = null
              type = "CLAIMS"
            }
            value = "given_name"
          }
          subject = {
            source = {
              id   = null
              type = "CLAIMS"
            }
            value = "sub"
          }
        }
        attribute_sources = [
          {
            jdbc_attribute_source = {
              attribute_contract_fulfillment = null
              column_names                   = ["GRANTEE"]
              data_store_ref = {
                id = "ProvisionerDS"
              }
              description = "JDBC"
              filter      = "subject"
              id          = null
              schema      = "INFORMATION_SCHEMA"
              table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
            }
          },
        ]
        authentication_policy_contract_ref = {
          id = pingfederate_authentication_policy_contract.apc1.id
        }
        issuance_criteria = {
          conditional_criteria = [
            {
              attribute_name = "email"
              condition      = "MULTIVALUE_CONTAINS_DN"
              error_result   = "myerrorresult"
              source = {
                type = "MAPPED_ATTRIBUTES"
              }
              value = "cn=Example,dc=example,dc=com"
            },
          ]
          expression_criteria = null
        }
        restrict_virtual_server_ids   = false
        restricted_virtual_server_ids = []
      },
    ]
    authn_context_mappings = [
      {
        local  = "1fa"
        remote = "Single_Factor"
      },
      {
        local  = "mfa"
        remote = "Multi_Factor"
      },
    ]
    default_target_url   = "https://example.com"
    idp_identity_mapping = "ACCOUNT_MAPPING"
    jit_provisioning = {
      error_handling = "ABORT_SSO"
      event_trigger  = "NEW_USER_ONLY"
      user_attributes = {
        do_attribute_query = false
      }
      user_repository = {
        ldap = {
          data_store_ref = {
            id = "pingdirectory"
          }
          unique_user_id_filter = "uid=john,ou=org"
          base_dn               = "dc=example,dc=com"
          jit_repository_attribute_mapping = {
            USER_KEY = {
              source = {
                id   = null
                type = "NO_MAPPING"
              }
              value = ""
            }
            USER_NAME = {
              source = {
                id   = null
                type = "TEXT"
              }
              value = "asdf"
            }
          }
        }
      }
    }
    oauth_authentication_policy_contract_ref = {
      id = pingfederate_authentication_policy_contract.apc1.id
    }
    oidc_provider_settings = {
      authentication_scheme                 = "PRIVATE_KEY_JWT"
      authentication_signing_algorithm      = "RS256"
      authorization_endpoint                = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/authorize"
      enable_pkce                           = true
      jwks_url                              = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/jwks"
      login_type                            = "CODE"
      pushed_authorization_request_endpoint = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/par"
      request_parameters = [
        {
          application_endpoint_override = false
          attribute_value = {
            source = {
              id   = null
              type = "TEXT"
            }
            value = "param1"
          }
          name  = "param1"
          value = null
        },
        {
          application_endpoint_override = true
          attribute_value = {
            source = {
              id   = null
              type = "CONTEXT"
            }
            value = "ClientIp"
          }
          name  = "param2"
          value = null
        },
      ]
      request_signing_algorithm      = "RS256"
      scopes                         = "openid profile email address phone"
      token_endpoint                 = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/token"
      track_user_sessions_for_logout = true
      user_info_endpoint             = "https://auth.pingone.eu/85a52cf7-357f-40c1-b909-de24d976031d/as/userinfo"
    }
    protocol            = "OIDC"
    sign_authn_requests = false
  }
  idp_oauth_grant_attribute_mapping = {
    access_token_manager_mappings = [
      {
        attribute_contract_fulfillment = {
          "Username" = {
            source = {
              type = "NO_MAPPING"
            }
            value = ""
          }
          "OrgName" = {
            source = {
              type = "NO_MAPPING"
            }
            value = ""
          }
        }
        attribute_sources = [
          {
            custom_attribute_source = {
              data_store_ref = {
                id = "customDataStore"
              }
              description = "APIStubs"
              filter_fields = [
                {
                  name = "Authorization Header"
                },
                {
                  name = "Body"
                },
                {
                  name  = "Resource Path"
                  value = "/users/extid"
                },
              ]
              id = "APIStubs"
            }
          },
        ]
        issuance_criteria = {
          conditional_criteria = [
            {
              attribute_name = "OrgName"
              condition      = "MULTIVALUE_CONTAINS_DN"
              error_result   = "myerrorresult"
              source = {
                type = "MAPPED_ATTRIBUTES"
              }
              value = "cn=Example,dc=example,dc=com"
            },
          ]
          expression_criteria = null
        }
        access_token_manager_ref = {
          id = pingfederate_oauth_access_token_manager.idpConnOidcAtm.id
        }
      }
    ]
    idp_oauth_attribute_contract = {
      extended_attributes = [
        {
          masked = false
          name   = "asdf"
        },
        {
          masked = false
          name   = "asdfd"
        }
      ]
    }
  }
  logging_mode = "FULL"
  name         = "PingOne"
  oidc_client_credentials = {
    client_id     = "myclientid"
    client_secret = "myclientsecrets"
  }
  # Ensures this resource will be updated before deleting the oauth access token manager
  lifecycle {
    create_before_destroy = true
  }
}
`, spIdpConnection_OidcDependencyHCL(),
		accesstokenmanager.TestAccessTokenManagerHCL("idpConnOidcAtm"),
		idpConnOidcId)
}

// Validate any computed values when applying minimal HCL
func spIdpConnection_CheckComputedValuesOidcMinimal() resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "active", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "additional_allowed_entities_configuration.addtional_allowed_entities.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "additional_allowed_entities_configuration.allow_additional_entities", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "additional_allowed_entities_configuration.allow_all_entities", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "error_page_msg_id", "errorDetail.spSsoFailure"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "id", idpConnOidcId),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.assertions_signed", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.attribute_contract.core_attributes.0.name", "sub"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.authentication_policy_contract_mappings.0.attribute_sources.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.authentication_policy_contract_mappings.0.issuance_criteria.conditional_criteria.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.authentication_policy_contract_mappings.0.restrict_virtual_server_ids", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.authentication_policy_contract_mappings.0.restricted_virtual_server_ids.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.default_target_url", ""),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.enable_pkce", "false"),
		resource.TestCheckNoResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.back_channel_logout_uri"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.redirect_uri"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.request_parameters.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.track_user_sessions_for_logout", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.sign_authn_requests", "false"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "idp_browser_sso.sso_application_endpoint"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "logging_mode", "STANDARD"),
		resource.TestCheckNoResourceAttr("pingfederate_sp_idp_connection.example", "oidc_client_credentials.client_secret"),
		resource.TestCheckNoResourceAttr("pingfederate_sp_idp_connection.example", "oidc_client_credentials.encrypted_secret"),
	}

	if acctest.VersionAtLeast(version.PingFederate1200) {
		testCheckFuncs = append(testCheckFuncs,
			resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.front_channel_logout_uri"),
			resource.TestCheckNoResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.post_logout_redirect_uri"),
		)
	}
	if acctest.VersionAtLeast(version.PingFederate1210) {
		testCheckFuncs = append(testCheckFuncs,
			resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.jwt_secured_authorization_response_mode_type", "DISABLED"),
		)
	}

	return resource.ComposeTestCheckFunc(testCheckFuncs...)
}

// Validate any computed values when applying complete HCL
func spIdpConnection_CheckComputedValuesOidcComplete() resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "id", idpConnOidcId),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.attribute_contract.core_attributes.0.name", "sub"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.jit_provisioning.user_attributes.attribute_contract.#", "15"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.back_channel_logout_uri"),
		resource.TestCheckNoResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.post_logout_redirect_uri"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.redirect_uri"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "idp_browser_sso.sso_application_endpoint"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.0.name", "TOKEN_SUBJECT"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "oidc_client_credentials.encrypted_secret"),
	}

	if acctest.VersionAtLeast(version.PingFederate1200) {
		testCheckFuncs = append(testCheckFuncs,
			resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "idp_browser_sso.oidc_provider_settings.front_channel_logout_uri"),
		)
	}

	return resource.ComposeTestCheckFunc(testCheckFuncs...)
}

// Test that any objects created by the test are destroyed
func spIdpConnection_OidcCheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.SpIdpConnectionsAPI.DeleteConnection(acctest.TestBasicAuthContext(), idpConnOidcId).Execute()
	if err == nil {
		return fmt.Errorf("sp_idp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
