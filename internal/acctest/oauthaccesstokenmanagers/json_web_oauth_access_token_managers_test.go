package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// #nosec G101
const jsonWebTokenOauthAccessTokenManagersId = "jsonWebTokenOatm"

// #nosec G101
const jsonWebTokenOauthAccessTokenManagersName = "jsonWebTokenExample"

// Attributes to test with. Add optional properties to test here if desired.
type jsonWebTokenOauthAccessTokenManagersResourceModel struct {
	id   string
	name string
	key  string
}

func TestAccJsonWebTokenOauthAccessTokenManagers(t *testing.T) {
	resourceName := "myJsonWebTokenOauthAccessTokenManagers"
	initialResourceModel := jsonWebTokenOauthAccessTokenManagersResourceModel{
		id:   jsonWebTokenOauthAccessTokenManagersId,
		name: jsonWebTokenOauthAccessTokenManagersName,
		key:  "+d5OB5b+I4dqn1Mjp8YE/M/QFWvDX7Nxz3gC8mAEwRLqL67SrHcwRyMtGvZKxvIn",
	}
	updatedResourceModel := jsonWebTokenOauthAccessTokenManagersResourceModel{
		id:   jsonWebTokenOauthAccessTokenManagersId,
		name: jsonWebTokenOauthAccessTokenManagersName,
		key:  "e1oDxOiC3Jboz3um8hBVmW3JRZNo9z7C0DMm/oj2V1gclQRcgi2gKM2DBj9N05G4",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckJsonWebOauthAccessTokenManagersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJsonWebOauthAccessTokenManagers(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedJsonWebOauthAccessTokenManagersAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccJsonWebOauthAccessTokenManagers(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedJsonWebOauthAccessTokenManagersAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccJsonWebOauthAccessTokenManagers(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_oauth_access_token_managers." + resourceName,
				ImportStateId:           jsonWebTokenOauthAccessTokenManagersId,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"configuration.fields.value"},
			},
		},
	})
}

func testAccJsonWebOauthAccessTokenManagers(resourceName string, resourceModel jsonWebTokenOauthAccessTokenManagersResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_managers" "%[1]s" {
  id   = "%[2]s"
  name = "%[3]s"
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
                value = "keyidentifier"
              },
              {
                name  = "Key"
                value = "%[4]s"
              },
              {
                name  = "Encoding"
                value = "b64u"
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
        value = "120"
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
        value = "keyidentifier"
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
        value = "keyidentifier"
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
      }
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
  selection_settings = {
    resource_uris = []
  }
  access_control_settings = {
    restrict_clients = false
  }
  session_validation_settings = {
    check_valid_authn_session       = false
    check_session_revocation_status = false
    update_authn_session_activity   = false
    include_session_id              = false
  }
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.key,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedJsonWebOauthAccessTokenManagersAttributes(config jsonWebTokenOauthAccessTokenManagersResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAccessTokenManagers"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAccessTokenManagersApi.GetTokenManager(ctx, jsonWebTokenOauthAccessTokenManagersId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		getFields := response.Configuration.Fields
		for _, field := range getFields {
			// if field.Name == "Token Length" {
			// 	err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.key, *field.Value)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			// if field.Name == "Token Lifetime" {
			// 	err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.tokenLifetime, *field.Value)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckJsonWebOauthAccessTokenManagersDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthAccessTokenManagersApi.DeleteTokenManager(ctx, jsonWebTokenOauthAccessTokenManagersId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAccessTokenManagers", jsonWebTokenOauthAccessTokenManagersId)
	}
	return nil
}
