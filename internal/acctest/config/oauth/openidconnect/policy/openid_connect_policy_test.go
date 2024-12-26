package oauthopenidconnectpolicy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

const oauthOpenIdConnectPoliciesId = "testOpenIdConnectPolicy"

// Attributes to test with. Add optional properties to test here if desired.
type oauthOpenIdConnectPoliciesResourceModel struct {
	id                        string
	name                      string
	includeOptionalAttributes bool

	attributeSource             *client.LdapAttributeSource
	idTokenLifetime             *int64
	includeSriInIdToken         *bool
	includeUserInfoInIdToken    *bool
	includeSHashInIdToken       *bool
	returnIdTokenOnRefreshGrant *bool
	reissueIdTokenInHybridFlow  *bool
}

// This is due to a bug in PingFederate that doesn't allow the OAuth client to set "None" as the OIDC Policy
// When an OAuth client is created, it comes with a "Default" OIDC Policy
// Once an OIDC Policy is created, it automatically attaches to the OAuth client configuration
// This is a workaround to delete the conflicting OAuth client
func deleteOauthClient(t *testing.T) {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	resp, err := testClient.OauthClientsAPI.DeleteOauthClient(ctx, "test").Execute()
	// If the error isn't a 404 (the client already doesn't exist) then fail
	if err != nil && (resp == nil || resp.StatusCode != 404) {
		t.Fatalf("Failed to delete conflicting \"test\" OAuth Client: %v", err)
	}
}

func TestAccOpenidConnectPolicies(t *testing.T) {
	resourceName := "myOpenidConnectPolicies"

	initialResourceModel := oauthOpenIdConnectPoliciesResourceModel{
		id:                        oauthOpenIdConnectPoliciesId,
		name:                      "initialName",
		includeOptionalAttributes: false,
	}
	updatedResourceModel := oauthOpenIdConnectPoliciesResourceModel{
		id:                          oauthOpenIdConnectPoliciesId,
		name:                        "updatedName",
		includeOptionalAttributes:   true,
		attributeSource:             attributesources.LdapClientStruct("(cn=Mudkip)", "SUBTREE", *client.NewResourceLink("pingdirectory")),
		idTokenLifetime:             pointers.Int64(5),
		includeSriInIdToken:         pointers.Bool(true),
		includeUserInfoInIdToken:    pointers.Bool(true),
		includeSHashInIdToken:       pointers.Bool(true),
		returnIdTokenOnRefreshGrant: pointers.Bool(true),
		reissueIdTokenInHybridFlow:  pointers.Bool(true),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOpenidConnectPoliciesDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { deleteOauthClient(t) },
				Config:    testAccOpenidConnectPolicies(resourceName, initialResourceModel),
				Check:     testAccCheckExpectedOpenidConnectPoliciesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOpenidConnectPolicies(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOpenidConnectPoliciesAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_policy.%s", resourceName), "id_token_lifetime", fmt.Sprintf("%d", *updatedResourceModel.idTokenLifetime)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_policy.%s", resourceName), "include_sri_in_id_token", fmt.Sprintf("%t", *updatedResourceModel.includeSriInIdToken)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_policy.%s", resourceName), "include_user_info_in_id_token", fmt.Sprintf("%t", *updatedResourceModel.includeUserInfoInIdToken)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_policy.%s", resourceName), "include_s_hash_in_id_token", fmt.Sprintf("%t", *updatedResourceModel.includeSHashInIdToken)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_policy.%s", resourceName), "return_id_token_on_refresh_grant", fmt.Sprintf("%t", *updatedResourceModel.returnIdTokenOnRefreshGrant)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_policy.%s", resourceName), "reissue_id_token_in_hybrid_flow", fmt.Sprintf("%t", *updatedResourceModel.reissueIdTokenInHybridFlow)),
				),
			},
			{
				// Test importing the resource
				Config:            testAccOpenidConnectPolicies(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_openid_connect_policy." + resourceName,
				ImportStateId:     oauthOpenIdConnectPoliciesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccOpenidConnectPolicies(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOpenidConnectPoliciesAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OauthOpenIdConnectAPI.DeleteOIDCPolicy(ctx, oauthOpenIdConnectPoliciesId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccOpenidConnectPolicies(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOpenidConnectPoliciesAttributes(initialResourceModel),
			},
		},
	})
}

func accessTokenManagerHcl() string {
	return `
resource "pingfederate_oauth_access_token_manager" "jsonWebTokenOauthAccessTokenManagerExample" {
  manager_id = "oidcJsonWebTokenExample"
  name       = "oidcJsonWebTokenExample"
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
                value = "e1oDxOiC3Jboz3um8hBVmW3JRZNo9z7C0DMm/oj2V1gclQRcgi2gKM2DBj9N05G4"
              },
              {
                name  = "Encoding"
                value = "b64u"
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
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "contract"
      },
      {
        name         = "another"
        multi_valued = false
      }
    ]
  }
}`
}

func attributeMappingHcl(resourceModel oauthOpenIdConnectPoliciesResourceModel) string {
	issuanceCriteriaHcl := ""
	if resourceModel.includeOptionalAttributes {
		issuanceCriteriaHcl = `
	    issuance_criteria = {
			conditional_criteria = []
		}
		`
	}

	return fmt.Sprintf(`
	attribute_mapping = {
		attribute_contract_fulfillment = {
			"sub" = {
				source = {
		  			type = "TOKEN"
				}
				value = "contract"
			}
		}
		%s
		%s
	}
	`, attributesources.Hcl(nil, resourceModel.attributeSource), issuanceCriteriaHcl)
}

func testAccOpenidConnectPolicies(resourceName string, resourceModel oauthOpenIdConnectPoliciesResourceModel) string {
	optionalHcl := ""
	if resourceModel.includeOptionalAttributes {
		optionalHcl = fmt.Sprintf(`
		scope_attribute_mappings = {}
		return_id_token_on_refresh_grant = %t
		include_sri_in_id_token = %t
		include_s_hash_in_id_token = %t
		include_user_info_in_id_token = %t
		id_token_lifetime = %d
		reissue_id_token_in_hybrid_flow = %t
		`,
			*resourceModel.returnIdTokenOnRefreshGrant,
			*resourceModel.includeSriInIdToken,
			*resourceModel.includeSHashInIdToken,
			*resourceModel.includeUserInfoInIdToken,
			*resourceModel.idTokenLifetime,
			*resourceModel.reissueIdTokenInHybridFlow)

		if acctest.VersionAtLeast(version.PingFederate1130) {
			optionalHcl += `
		include_x5t_in_id_token = true
		id_token_typ_header_value = "Example"
			`
		}

		if acctest.VersionAtLeast(version.PingFederate1220) {
			optionalHcl += `
		return_id_token_on_token_exchange_grant = true
			`
		}
	}

	return fmt.Sprintf(`
	%s
resource "pingfederate_openid_connect_policy" "%s" {
  policy_id = "%s"
  name      = "%s"
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.jsonWebTokenOauthAccessTokenManagerExample.manager_id
  }
  attribute_contract = {
    extended_attributes = []
  }
	%s
	%s
}

data "pingfederate_openid_connect_policy" "%[2]s" {
  policy_id = pingfederate_openid_connect_policy.%[2]s.policy_id
}`, accessTokenManagerHcl(),
		resourceName,
		oauthOpenIdConnectPoliciesId,
		resourceModel.name,
		attributeMappingHcl(resourceModel),
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOpenidConnectPoliciesAttributes(config oauthOpenIdConnectPoliciesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OpenidConnectPolicy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthOpenIdConnectAPI.GetOIDCPolicy(ctx, oauthOpenIdConnectPoliciesId).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthOpenIdConnectPoliciesId), "name", config.name, response.Name)
		if err != nil {
			return err
		}

		if !config.includeOptionalAttributes {
			return nil
		}

		// Verify some optional attributes
		err = attributesources.ValidateResponseAttributes(resourceType, pointers.String(oauthOpenIdConnectPoliciesId),
			nil, config.attributeSource, response.AttributeMapping.AttributeSources)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, pointers.String(oauthOpenIdConnectPoliciesId), "id_token_lifetime", *config.idTokenLifetime, *response.IdTokenLifetime)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthOpenIdConnectPoliciesId), "return_id_token_on_refresh_grant", *config.returnIdTokenOnRefreshGrant, *response.ReturnIdTokenOnRefreshGrant)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthOpenIdConnectPoliciesId), "include_sri_in_id_token", *config.includeSriInIdToken, *response.IncludeSriInIdToken)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthOpenIdConnectPoliciesId), "include_user_info_in_id_token", *config.includeUserInfoInIdToken, *response.IncludeUserInfoInIdToken)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(oauthOpenIdConnectPoliciesId), "include_s_hash_in_id_token", *config.includeSHashInIdToken, *response.IncludeSHashInIdToken)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOpenidConnectPoliciesDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthOpenIdConnectAPI.DeleteOIDCPolicy(ctx, oauthOpenIdConnectPoliciesId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OpenidConnectPolicy", oauthOpenIdConnectPoliciesId)
	}
	return nil
}
