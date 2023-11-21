package acctest_test

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
)

const idpErrMessage = "errorDetail.idpSsoFailure"

// Attributes to test with. Add optional properties to test here if desired.
type idpDefaultUrlsResourceModel struct {
	confirmIdpSlo    *bool
	idpSloSuccessUrl *string
	idpErrorMsg      string
}

func TestAccIdpDefaultUrls(t *testing.T) {
	resourceName := "myIdpDefaultUrls"
	initialResourceModel := idpDefaultUrlsResourceModel{
		idpErrorMsg: idpErrMessage,
	}
	updatedResourceModel := idpDefaultUrlsResourceModel{
		confirmIdpSlo:    pointers.Bool(true),
		idpSloSuccessUrl: pointers.String("https://example"),
		idpErrorMsg:      idpErrMessage,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIdpDefaultUrls(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpDefaultUrlsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIdpDefaultUrls(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIdpDefaultUrlsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccIdpDefaultUrls(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_idp_default_urls." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccIdpDefaultUrls(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpDefaultUrlsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccIdpDefaultUrls(resourceName string, resourceModel idpDefaultUrlsResourceModel) string {
	confirmIdpSloHcl := ""
	idpSloSuccessfulUrlHcl := ""
	if resourceModel.confirmIdpSlo != nil {
		confirmIdpSloHcl = fmt.Sprintf("confirm_idp_slo = %[1]t", *resourceModel.confirmIdpSlo)
	}
	if resourceModel.idpSloSuccessUrl != nil {
		idpSloSuccessfulUrlHcl = fmt.Sprintf("idp_slo_success_url = \"%[1]s\"", *resourceModel.idpSloSuccessUrl)
	}
	return fmt.Sprintf(`
resource "pingfederate_idp_default_urls" "%[1]s" {
  idp_error_msg = "%[2]s"
  %[3]s
  %[4]s
}

data "pingfederate_idp_default_urls" "%[1]s" {
  depends_on = [
    pingfederate_idp_default_urls.%[1]s
  ]
}`, resourceName,
		resourceModel.idpErrorMsg,
		confirmIdpSloHcl,
		idpSloSuccessfulUrlHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedIdpDefaultUrlsAttributes(config idpDefaultUrlsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "IdpDefaultUrls"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.IdpDefaultUrlsAPI.GetDefaultUrl(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "idp_error_msg", config.idpErrorMsg, response.IdpErrorMsg)
		if err != nil {
			return err
		}

		if config.confirmIdpSlo != nil {
			err = acctest.TestAttributesMatchBool(resourceType, nil, "confirm_idp_slo",
				*config.confirmIdpSlo, *response.ConfirmIdpSlo)
			if err != nil {
				return err
			}
		}

		if config.idpSloSuccessUrl != nil {
			err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "idp_slo_success_url",
				*config.idpSloSuccessUrl, response.IdpSloSuccessUrl)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
