package oauthopenidconnectsettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

type openIdConnectSettingsResourceModel struct {
	defaultPolicyRef *client.ResourceLink
	sessionSettings  *client.OIDCSessionSettings
}

func createDependenciesForTestBecauseIHaveTo(t *testing.T) {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	oauthAccessTokenMgrForTest := client.NewAccessTokenManager(
		"test",
		"test",
		*client.NewResourceLink("org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"),
		client.PluginConfiguration{},
	)
	attributeContract := client.NewAccessTokenAttributeContract()

	attributeContract.ExtendedAttributes = append(attributeContract.ExtendedAttributes, *client.NewAccessTokenAttribute("extended_contract"))

	oauthAccessTokenMgrForTest.AttributeContract = attributeContract

	_, _, err := testClient.OauthAccessTokenManagersAPI.CreateTokenManager(ctx).Body(*oauthAccessTokenMgrForTest).Execute()

	if err != nil {
		t.Fatalf("Failed to create OAuth Access Token Manager: %v", err)
	}
	attributeContractFulfillment := map[string]client.AttributeFulfillmentValue{}
	afv := client.NewAttributeFulfillmentValue(
		*client.NewSourceTypeIdKey("TEXT"),
		"sub",
	)

	attributeContractFulfillment["sub"] = *afv
	attributeMapping := client.NewAttributeMapping(
		attributeContractFulfillment,
	)

	// Create OIDC Policy
	oidcPolicy := client.NewOpenIdConnectPolicy(
		"oidcPolicy",
		"oidcPolicy",
		*client.NewResourceLink("test"),
		client.OpenIdConnectAttributeContract{},
		*attributeMapping,
	)

	_, _, err = testClient.OauthOpenIdConnectAPI.CreateOIDCPolicy(ctx).Body(*oidcPolicy).Execute()
	if err != nil {
		t.Fatalf("Failed to create OIDC Policy: %v", err)
	}
}

// Delete dependency created for test
func deleteOidcPolicy(t *testing.T) {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthOpenIdConnectAPI.DeleteOIDCPolicy(ctx, "oidcPolicy").Execute()
	if err != nil {
		t.Fatalf("Failed to delete OIDC Policy: %v", err)
	}
}

func TestAccOpenIdConnectSettings(t *testing.T) {
	resourceName := "myOpenIdConnectSettings"
	// send empty model to start
	initialResourceModel := openIdConnectSettingsResourceModel{}
	updatedResourceModel := openIdConnectSettingsResourceModel{
		defaultPolicyRef: &client.ResourceLink{
			Id: "oidcPolicy",
		},
		sessionSettings: &client.OIDCSessionSettings{
			TrackUserSessionsForLogout: pointers.Bool(true),
			RevokeUserSessionOnLogout:  pointers.Bool(true),
			SessionRevocationLifetime:  pointers.Int64(180),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOpenIdConnectSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOpenIdConnectSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				PreConfig: func() { createDependenciesForTestBecauseIHaveTo(t) },
				Config:    testAccOpenIdConnectSettings(resourceName, updatedResourceModel),
				Check:     testAccCheckExpectedOpenIdConnectSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOpenIdConnectSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_open_id_connect_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { deleteOidcPolicy(t) },
				Config:    testAccOpenIdConnectSettings(resourceName, initialResourceModel),
				Check:     testAccCheckExpectedOpenIdConnectSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func generateDependentHcl(resourceModel openIdConnectSettingsResourceModel) string {
	var settingsSourceHcl string
	if resourceModel.defaultPolicyRef == nil {
		settingsSourceHcl = ""
	} else {
		settingsSourceHcl = fmt.Sprintf(`
		default_policy_ref = {
			id = "%s"
		}`, resourceModel.defaultPolicyRef.Id)
	}

	var sessionSettingsHcl string
	if resourceModel.sessionSettings == nil {
		sessionSettingsHcl = ""
	} else {
		sessionSettingsHcl = fmt.Sprintf(`
		session_settings = {
			track_user_sessions_for_logout = %t
			revoke_user_session_on_logout = %t
			session_revocation_lifetime = %d
		}`, *resourceModel.sessionSettings.TrackUserSessionsForLogout, *resourceModel.sessionSettings.RevokeUserSessionOnLogout, *resourceModel.sessionSettings.SessionRevocationLifetime)
	}

	return fmt.Sprintf(`
	%s
	%s
	`, settingsSourceHcl, sessionSettingsHcl)

}

func testAccOpenIdConnectSettings(resourceName string, resourceModel openIdConnectSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_open_id_connect_settings" "%[1]s" {
	%[2]s
}`, resourceName,
		generateDependentHcl(resourceModel),
	)
}

func testAccCheckExpectedOpenIdConnectSettingsAttributes(config openIdConnectSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OpenIdConnectSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		stateAttributes := s.RootModule().Resources["pingfederate_open_id_connect_settings.myOpenIdConnectSettings"].Primary.Attributes
		response, _, err := testClient.OauthOpenIdConnectAPI.GetOIDCSettings(ctx).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		if config.defaultPolicyRef != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "id", config.defaultPolicyRef.Id, response.DefaultPolicyRef.Id)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "default_policy_ref.id", config.defaultPolicyRef.Id)
			if err != nil {
				return err
			}
		}

		if config.sessionSettings != nil {
			err = acctest.TestAttributesMatchBool(resourceType, nil, "track_user_sessions_for_logout", *config.sessionSettings.TrackUserSessionsForLogout, *response.SessionSettings.TrackUserSessionsForLogout)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "session_settings.track_user_sessions_for_logout", *config.sessionSettings.TrackUserSessionsForLogout)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "revoke_user_session_on_logout", *config.sessionSettings.RevokeUserSessionOnLogout, *response.SessionSettings.RevokeUserSessionOnLogout)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "session_settings.revoke_user_session_on_logout", *config.sessionSettings.RevokeUserSessionOnLogout)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchInt(resourceType, nil, "session_revocation_lifetime", *config.sessionSettings.SessionRevocationLifetime, *response.SessionSettings.SessionRevocationLifetime)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "session_settings.session_revocation_lifetime", *config.sessionSettings.SessionRevocationLifetime)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
