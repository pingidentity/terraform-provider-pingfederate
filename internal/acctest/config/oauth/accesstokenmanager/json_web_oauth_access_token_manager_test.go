// Copyright © 2025 Ping Identity Corporation

package oauthaccesstokenmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// #nosec G101
const jsonWebTokenOauthAccessTokenManagerId = "jsonWebTokenOatm"

// #nosec G101
const jsonWebTokenOauthAccessTokenManagerName = "jsonWebTokenExample"

// Attributes to test with. Add optional properties to test here if desired.
type jsonWebTokenOauthAccessTokenManagerResourceModel struct {
	id                     string
	name                   string
	keyId                  string
	key                    string
	tokenLifetime          *string
	activeSymmetricKeyId   string
	checkValidAuthnSession *bool
}

func TestAccJsonWebTokenOauthAccessTokenManager(t *testing.T) {
	resourceName := "myJsonWebTokenOauthAccessTokenManager"

	initialResourceModel := jsonWebTokenOauthAccessTokenManagerResourceModel{
		id:                   jsonWebTokenOauthAccessTokenManagerId,
		name:                 jsonWebTokenOauthAccessTokenManagerName,
		keyId:                "keyidentifier",
		key:                  "+d5OB5b+I4dqn1Mjp8YE/M/QFWvDX7Nxz3gC8mAEwRLqL67SrHcwRyMtGvZKxvIn",
		activeSymmetricKeyId: "keyidentifier",
	}
	updatedResourceModel := jsonWebTokenOauthAccessTokenManagerResourceModel{
		id:                     jsonWebTokenOauthAccessTokenManagerId,
		name:                   jsonWebTokenOauthAccessTokenManagerName,
		keyId:                  "keyidentifier2",
		key:                    "e1oDxOiC3Jboz3um8hBVmW3JRZNo9z7C0DMm/oj2V1gclQRcgi2gKM2DBj9N05G4",
		tokenLifetime:          pointers.String("56"),
		activeSymmetricKeyId:   "keyidentifier2",
		checkValidAuthnSession: pointers.Bool(true),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckJsonWebOauthAccessTokenManagerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJsonWebOauthAccessTokenManagerMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedJsonWebOauthAccessTokenManagerAttributes(initialResourceModel, true),
			},
			{
				// Test updating some fields
				Config: testAccJsonWebOauthAccessTokenManager(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedJsonWebOauthAccessTokenManagerAttributes(updatedResourceModel, false),
			},
			{
				// Test importing the resource
				Config:            testAccJsonWebOauthAccessTokenManager(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_access_token_manager." + resourceName,
				ImportStateId:     jsonWebTokenOauthAccessTokenManagerId,
				ImportState:       true,
				ImportStateVerify: true,
				// Some configuration table fields are sensitive or are not set in this test, so they can't be verified here.
				ImportStateVerifyIgnore: []string{"configuration.tables.0.rows.0.sensitive_fields"},
			},
			{
				// Back to minimal model
				Config: testAccJsonWebOauthAccessTokenManagerMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedJsonWebOauthAccessTokenManagerAttributes(initialResourceModel, true),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OauthAccessTokenManagersAPI.DeleteTokenManager(ctx, updatedResourceModel.id).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccJsonWebOauthAccessTokenManagerMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedJsonWebOauthAccessTokenManagerAttributes(initialResourceModel, true),
			},
		},
	})
}

func testAccJsonWebOauthAccessTokenManagerMinimal(resourceName string, resourceModel jsonWebTokenOauthAccessTokenManagerResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = "%[2]s"
  name       = "%[3]s"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.access.token.management.plugins.JwtBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = [
      {
        name = "Symmetric Keys"
        rows = [
          {
            fields = [
              {
                name  = "Key ID"
                value = "%[4]s"
              },
              {
                name  = "Encoding"
                value = "b64u"
              }
            ]
            sensitive_fields = [
              {
                name  = "Key"
                value = "%[5]s"
              }
            ]
          }
        ]
      },
      {
        name = "Certificates"
        rows = []
      }
    ]
    fields = [
      {
        name  = "Active Symmetric Key ID"
        value = "%[6]s"
      },
      {
        name  = "JWE Algorithm"
        value = "dir"
      },
      {
        name  = "JWE Content Encryption Algorithm"
        value = "A192CBC-HS384"
      },
      {
        name  = "Active Symmetric Encryption Key ID"
        value = "%[6]s"
      },
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name         = "contract"
        multi_valued = false
      }
    ]
  }
}

data "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = pingfederate_oauth_access_token_manager.%[1]s.id
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.keyId,
		resourceModel.key,
		resourceModel.activeSymmetricKeyId,
	)
}

func testAccJsonWebOauthAccessTokenManager(resourceName string, resourceModel jsonWebTokenOauthAccessTokenManagerResourceModel) string {
	versionedFields := ""
	if acctest.VersionAtLeast(version.PingFederate1130) {
		versionedFields += `
      {
        name  = "Not Before Claim Offset"
        value = ""
      },
      {
        name  = "Include Issued At Claim",
        value = "false"
      },
    `
	}
	if acctest.VersionAtLeast(version.PingFederate1210) {
		versionedFields += `
      {
        name  = "Publish Keys to the PingFederate JWKS Endpoint"
        value = "false"
      },
    `
	}
	additionalVersionedHcl := ""
	if acctest.VersionAtLeast(version.PingFederate1220) {
		additionalVersionedHcl += `
  token_endpoint_attribute_contract = {
    attributes = [
      {
        mapped_scopes = ["email"]
        multi_valued  = false
        name          = "normal"
      },
      {
        mapped_scopes = []
        multi_valued  = true
        name          = "another"
      },
    ]
  }
    `
	}

	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = "%[2]s"
  name       = "%[3]s"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.access.token.management.plugins.JwtBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = [
      {
        name = "Symmetric Keys"
        rows = [
          {
            fields = [
              {
                name  = "Key ID"
                value = "%[4]s"
              },
              {
                name  = "Encoding"
                value = "b64u"
              }
            ]
            sensitive_fields = [
              {
                name  = "Key"
                value = "%[5]s"
              }
            ]
            default_row = false
          }
        ]
      },
      {
        name = "Certificates"
        rows = []
      }
    ]
    fields = [
      {
        name  = "Token Lifetime"
        value = "%[6]s"
      },
      {
        name  = "Use Centralized Signing Key"
        value = "false"
      },
      {
        name  = "JWS Algorithm"
        value = ""
      },
      {
        name  = "Active Symmetric Key ID"
        value = "%[7]s"
      },
      {
        name  = "Active Signing Certificate Key ID"
        value = ""
      },
      {
        name  = "JWE Algorithm"
        value = "dir"
      },
      {
        name  = "JWE Content Encryption Algorithm"
        value = "A192CBC-HS384"
      },
      {
        name  = "Active Symmetric Encryption Key ID"
        value = "%[7]s"
      },
      {
        name  = "Asymmetric Encryption Key"
        value = ""
      },
      {
        name  = "Asymmetric Encryption JWKS URL"
        value = ""
      },
      {
        name  = "Enable Token Revocation"
        value = "false"
      },
      {
        name  = "Include Key ID Header Parameter"
        value = "true"
      },
      {
        name  = "Default JWKS URL Cache Duration"
        value = "720"
      },
      {
        name  = "Include JWE Key ID Header Parameter"
        value = "true"
      },
      {
        name  = "Client ID Claim Name"
        value = "client_id"
      },
      {
        name  = "Scope Claim Name"
        value = "scope"
      },
      {
        name  = "Space Delimit Scope Values"
        value = "true"
      },
      {
        name  = "Authorization Details Claim Name"
        value = "authorization_details"
      },
      {
        name  = "Issuer Claim Value"
        value = ""
      },
      {
        name  = "Audience Claim Value"
        value = ""
      },
      {
        name  = "JWT ID Claim Length"
        value = "22"
      },
      {
        name  = "Access Grant GUID Claim Name"
        value = ""
      },
      {
        name  = "JWKS Endpoint Path"
        value = ""
      },
      {
        name  = "JWKS Endpoint Cache Duration"
        value = "720"
      },
      {
        name  = "Expand Scope Groups"
        value = "false"
      },
      {
        name  = "Type Header Value"
        value = ""
      },
      {
        name  = "Publish Thumbprint X.509 URL"
        value = "false"
      },
      {
        name  = "Publish Key ID X.509 URL"
        value = "false"
      },
      {
        name  = "Include JWE X.509 Thumbprint Header Parameter",
        value = "false"
      },
      {
        name  = "Include X.509 Thumbprint Header Parameter",
        value = "false"
      },
      %[9]s
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name         = "contract"
        multi_valued = false
      },
      {
        name         = "contract1"
        multi_valued = false
      },
      {
        name         = "contract2"
        multi_valued = false
      },
      {
        name         = "contract3"
        multi_valued = false
      },
      {
        name         = "contract4"
        multi_valued = false
      },
      {
        name         = "contract5"
        multi_valued = false
      }
    ]
  }
  selection_settings = {
    resource_uris = []
  }
  access_control_settings = {
    restrict_clients = false
  }
  session_validation_settings = {
    check_valid_authn_session       = %[8]t
    check_session_revocation_status = false
    update_authn_session_activity   = false
    include_session_id              = false
  }
  %[10]s
}
data "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = pingfederate_oauth_access_token_manager.%[1]s.id
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.keyId,
		resourceModel.key,
		*resourceModel.tokenLifetime,
		resourceModel.activeSymmetricKeyId,
		*resourceModel.checkValidAuthnSession,
		versionedFields,
		additionalVersionedHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedJsonWebOauthAccessTokenManagerAttributes(config jsonWebTokenOauthAccessTokenManagerResourceModel, minimal bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAccessTokenManager"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAccessTokenManagersAPI.GetTokenManager(ctx, jsonWebTokenOauthAccessTokenManagerId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		getTables := response.Configuration.Tables
		for _, table := range getTables {
			for _, row := range table.Rows {
				for _, field := range row.Fields {
					if field.Name == "Key ID" {
						err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.keyId, *field.Value)
						if err != nil {
							return err
						}
					}
				}
			}
		}
		getFields := response.Configuration.Fields
		for _, field := range getFields {
			if field.Name == "Token Lifetime" && !minimal {
				err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", *config.tokenLifetime, *field.Value)
				if err != nil {
					return err
				}
			}
			if field.Name == "Active Symmetric Key ID" {
				err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.activeSymmetricKeyId, *field.Value)
				if err != nil {
					return err
				}
			}
		}

		if !minimal {
			err = acctest.TestAttributesMatchBool(resourceType, &config.id, "check_valid_authn_session", *config.checkValidAuthnSession, *response.SessionValidationSettings.CheckValidAuthnSession)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckJsonWebOauthAccessTokenManagerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthAccessTokenManagersAPI.DeleteTokenManager(ctx, jsonWebTokenOauthAccessTokenManagerId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAccessTokenManager", jsonWebTokenOauthAccessTokenManagerId)
	}
	return nil
}
