package redirectvalidation_test

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

// Attributes to test with. Add optional properties to test here if desired.
type redirectValidationResourceModel struct {
	redirectValidationLocalSettings   *client.RedirectValidationLocalSettings
	whiteList                         []*client.RedirectValidationSettingsWhitelistEntry
	redirectValidationPartnerSettings *client.RedirectValidationPartnerSettings
}

func TestAccRedirectValidation(t *testing.T) {
	resourceName := "myRedirectValidation"
	updatedResourceModel := redirectValidationResourceModel{
		redirectValidationLocalSettings: &client.RedirectValidationLocalSettings{
			EnableTargetResourceValidationForSSO:          pointers.Bool(true),
			EnableTargetResourceValidationForSLO:          pointers.Bool(true),
			EnableTargetResourceValidationForIdpDiscovery: pointers.Bool(true),
			EnableInErrorResourceValidation:               pointers.Bool(true),
		},
		whiteList: []*client.RedirectValidationSettingsWhitelistEntry{
			{
				TargetResourceSSO:     pointers.Bool(true),
				TargetResourceSLO:     pointers.Bool(true),
				InErrorResource:       pointers.Bool(true),
				IdpDiscovery:          pointers.Bool(true),
				ValidDomain:           "example.com",
				ValidPath:             pointers.String("/path"),
				AllowQueryAndFragment: pointers.Bool(true),
				RequireHttps:          pointers.Bool(true),
			},
			{
				TargetResourceSSO:     pointers.Bool(true),
				TargetResourceSLO:     pointers.Bool(true),
				InErrorResource:       pointers.Bool(true),
				IdpDiscovery:          pointers.Bool(true),
				ValidDomain:           "anotherexample.com",
				ValidPath:             pointers.String("/path2"),
				RequireHttps:          pointers.Bool(true),
				AllowQueryAndFragment: pointers.Bool(true),
			},
		},
		redirectValidationPartnerSettings: &client.RedirectValidationPartnerSettings{
			EnableWreplyValidationSLO: pointers.Bool(true),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRedirectValidationMinial(resourceName),
			},
			{
				// Test updating some fields
				Config: testAccRedirectValidation(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedRedirectValidationAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.enable_target_resource_validation_for_sso", fmt.Sprintf("%t", *updatedResourceModel.redirectValidationLocalSettings.EnableTargetResourceValidationForSSO)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.valid_domain", updatedResourceModel.whiteList[0].ValidDomain),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.valid_path", *updatedResourceModel.whiteList[0].ValidPath),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.allow_query_and_fragment", fmt.Sprintf("%t", *updatedResourceModel.whiteList[0].AllowQueryAndFragment)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.require_https", fmt.Sprintf("%t", *updatedResourceModel.whiteList[0].RequireHttps)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.idp_discovery", fmt.Sprintf("%t", *updatedResourceModel.whiteList[0].IdpDiscovery)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.target_resource_sso", fmt.Sprintf("%t", *updatedResourceModel.whiteList[0].TargetResourceSSO)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.target_resource_slo", fmt.Sprintf("%t", *updatedResourceModel.whiteList[0].TargetResourceSLO)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.0.in_error_resource", fmt.Sprintf("%t", *updatedResourceModel.whiteList[0].InErrorResource)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.valid_domain", updatedResourceModel.whiteList[1].ValidDomain),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.valid_path", *updatedResourceModel.whiteList[1].ValidPath),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.allow_query_and_fragment", fmt.Sprintf("%t", *updatedResourceModel.whiteList[1].AllowQueryAndFragment)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.require_https", fmt.Sprintf("%t", *updatedResourceModel.whiteList[1].RequireHttps)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.idp_discovery", fmt.Sprintf("%t", *updatedResourceModel.whiteList[1].IdpDiscovery)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.target_resource_sso", fmt.Sprintf("%t", *updatedResourceModel.whiteList[1].TargetResourceSSO)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.target_resource_slo", fmt.Sprintf("%t", *updatedResourceModel.whiteList[1].TargetResourceSLO)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.white_list.1.in_error_resource", fmt.Sprintf("%t", *updatedResourceModel.whiteList[1].InErrorResource)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_partner_settings.enable_wreply_validation_slo", fmt.Sprintf("%t", *updatedResourceModel.redirectValidationPartnerSettings.EnableWreplyValidationSLO)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_redirect_validation.%s", resourceName), "redirect_validation_local_settings.enable_target_resource_validation_for_slo", fmt.Sprintf("%t", *updatedResourceModel.redirectValidationLocalSettings.EnableTargetResourceValidationForSLO)),
				),
			},
			{
				// Test importing the resource
				Config:            testAccRedirectValidation(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_redirect_validation." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccRedirectValidationMinial(resourceName),
			},
		},
	})
}

func testAccRedirectValidationMinial(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_redirect_validation" "%[1]s" {
}`, resourceName)
}

func testAccRedirectValidation(resourceName string, resourceModel redirectValidationResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_redirect_validation" "%[1]s" {
  redirect_validation_local_settings = {
    enable_target_resource_validation_for_sso           = %[2]t
    enable_target_resource_validation_for_slo           = %[3]t
    enable_target_resource_validation_for_idp_discovery = %[4]t
    enable_in_error_resource_validation                 = %[5]t
    white_list = [
      {
        target_resource_sso      = %[6]t
        target_resource_slo      = %[7]t
        in_error_resource        = %[8]t
        idp_discovery            = %[9]t
        valid_domain             = "%[10]s"
        valid_path               = "%[11]s"
        allow_query_and_fragment = %[12]t
        require_https            = %[13]t
      },
      {
        target_resource_sso      = %[14]t
        target_resource_slo      = %[15]t
        in_error_resource        = %[16]t
        idp_discovery            = %[17]t
        valid_domain             = "%[18]s"
        valid_path               = "%[19]s"
        allow_query_and_fragment = %[20]t
        require_https            = %[21]t
      }
    ]
  }
  redirect_validation_partner_settings = {
    enable_wreply_validation_slo = %[22]t
  }
}
data "pingfederate_redirect_validation" "%[1]s" {
  depends_on = [pingfederate_redirect_validation.%[1]s]
}
`, resourceName,
		*resourceModel.redirectValidationLocalSettings.EnableTargetResourceValidationForSSO,
		*resourceModel.redirectValidationLocalSettings.EnableTargetResourceValidationForSLO,
		*resourceModel.redirectValidationLocalSettings.EnableTargetResourceValidationForIdpDiscovery,
		*resourceModel.redirectValidationLocalSettings.EnableInErrorResourceValidation,
		*resourceModel.whiteList[0].TargetResourceSSO,
		*resourceModel.whiteList[0].TargetResourceSLO,
		*resourceModel.whiteList[0].InErrorResource,
		*resourceModel.whiteList[0].IdpDiscovery,
		resourceModel.whiteList[0].ValidDomain,
		*resourceModel.whiteList[0].ValidPath,
		*resourceModel.whiteList[0].AllowQueryAndFragment,
		*resourceModel.whiteList[0].RequireHttps,
		*resourceModel.whiteList[1].TargetResourceSSO,
		*resourceModel.whiteList[1].TargetResourceSLO,
		*resourceModel.whiteList[1].InErrorResource,
		*resourceModel.whiteList[1].IdpDiscovery,
		resourceModel.whiteList[1].ValidDomain,
		*resourceModel.whiteList[1].ValidPath,
		*resourceModel.whiteList[1].AllowQueryAndFragment,
		*resourceModel.whiteList[1].RequireHttps,
		*resourceModel.redirectValidationPartnerSettings.EnableWreplyValidationSLO,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedRedirectValidationAttributes(config redirectValidationResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "RedirectValidation"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.RedirectValidationAPI.GetRedirectValidationSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_target_resource_validation_for_sso", *config.redirectValidationLocalSettings.EnableTargetResourceValidationForSSO, *response.RedirectValidationLocalSettings.EnableTargetResourceValidationForSSO)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_target_resource_validation_for_slo", *config.redirectValidationLocalSettings.EnableTargetResourceValidationForSLO, *response.RedirectValidationLocalSettings.EnableTargetResourceValidationForSLO)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_target_resource_validation_for_idp_discovery", *config.redirectValidationLocalSettings.EnableTargetResourceValidationForIdpDiscovery, *response.RedirectValidationLocalSettings.EnableTargetResourceValidationForIdpDiscovery)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_in_error_resource_validation", *config.redirectValidationLocalSettings.EnableInErrorResourceValidation, *response.RedirectValidationLocalSettings.EnableInErrorResourceValidation)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "target_resource_sso", *config.whiteList[0].TargetResourceSSO, *response.RedirectValidationLocalSettings.WhiteList[0].TargetResourceSSO)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "target_resource_slo", *config.whiteList[0].TargetResourceSLO, *response.RedirectValidationLocalSettings.WhiteList[0].TargetResourceSLO)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "in_error_resource", *config.whiteList[0].InErrorResource, *response.RedirectValidationLocalSettings.WhiteList[0].InErrorResource)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "idp_discovery", *config.whiteList[0].IdpDiscovery, *response.RedirectValidationLocalSettings.WhiteList[0].IdpDiscovery)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "valid_domain", config.whiteList[0].ValidDomain, response.RedirectValidationLocalSettings.WhiteList[0].ValidDomain)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "valid_path", *config.whiteList[0].ValidPath, *response.RedirectValidationLocalSettings.WhiteList[0].ValidPath)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "allow_query_and_fragment", *config.whiteList[0].AllowQueryAndFragment, *response.RedirectValidationLocalSettings.WhiteList[0].AllowQueryAndFragment)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "require_https", *config.whiteList[0].RequireHttps, *response.RedirectValidationLocalSettings.WhiteList[0].RequireHttps)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "target_resource_sso", *config.whiteList[1].TargetResourceSSO, *response.RedirectValidationLocalSettings.WhiteList[1].TargetResourceSSO)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "target_resource_slo", *config.whiteList[1].TargetResourceSLO, *response.RedirectValidationLocalSettings.WhiteList[1].TargetResourceSLO)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "in_error_resource", *config.whiteList[1].InErrorResource, *response.RedirectValidationLocalSettings.WhiteList[1].InErrorResource)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "idp_discovery", *config.whiteList[1].IdpDiscovery, *response.RedirectValidationLocalSettings.WhiteList[1].IdpDiscovery)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "valid_domain", config.whiteList[1].ValidDomain, response.RedirectValidationLocalSettings.WhiteList[1].ValidDomain)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "valid_path", *config.whiteList[1].ValidPath, *response.RedirectValidationLocalSettings.WhiteList[1].ValidPath)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "allow_query_and_fragment", *config.whiteList[1].AllowQueryAndFragment, *response.RedirectValidationLocalSettings.WhiteList[1].AllowQueryAndFragment)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "require_https", *config.whiteList[1].RequireHttps, *response.RedirectValidationLocalSettings.WhiteList[1].RequireHttps)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_wreply_validation_slo", *config.redirectValidationPartnerSettings.EnableWreplyValidationSLO, *response.RedirectValidationPartnerSettings.EnableWreplyValidationSLO)
		if err != nil {
			return err
		}

		return nil
	}
}
