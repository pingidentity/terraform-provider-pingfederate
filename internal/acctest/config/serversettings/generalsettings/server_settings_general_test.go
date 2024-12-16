package serversettingsgeneralsettings_test

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

// Attributes to test with. Add optional properties to test here if desired.
type serverSettingsGeneralResourceModel struct {
	disableAutomaticConnectionValidation    bool
	idpConnectionTransactionLoggingOverride string
	spConnectionTransactionLoggingOverride  string
	datastoreValidationIntervalSecs         int64
	requestHeaderForCorrelationId           string
}

func TestAccServerSettingsGeneral(t *testing.T) {
	resourceName := "myServerSettingsGeneral"
	updatedResourceModel := serverSettingsGeneralResourceModel{
		disableAutomaticConnectionValidation:    true,
		idpConnectionTransactionLoggingOverride: "FULL",
		spConnectionTransactionLoggingOverride:  "NONE",
		datastoreValidationIntervalSecs:         300,
		requestHeaderForCorrelationId:           "updatedExample",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServerSettingsGeneral(resourceName, nil),
				Check:  testAccCheckExpectedServerSettingsGeneralAttributes(nil),
			},
			{
				// Test updating some fields
				Config: testAccServerSettingsGeneral(resourceName, &updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedServerSettingsGeneralAttributes(&updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_server_settings_general.%s", resourceName), "disable_automatic_connection_validation", fmt.Sprintf("%t", updatedResourceModel.disableAutomaticConnectionValidation)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_server_settings_general.%s", resourceName), "idp_connection_transaction_logging_override", updatedResourceModel.idpConnectionTransactionLoggingOverride),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_server_settings_general.%s", resourceName), "sp_connection_transaction_logging_override", updatedResourceModel.spConnectionTransactionLoggingOverride),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_server_settings_general.%s", resourceName), "datastore_validation_interval_secs", fmt.Sprintf("%d", updatedResourceModel.datastoreValidationIntervalSecs)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_server_settings_general.%s", resourceName), "request_header_for_correlation_id", updatedResourceModel.requestHeaderForCorrelationId),
				),
			},
			{
				// Test importing the resource
				Config:                               testAccServerSettingsGeneral(resourceName, &updatedResourceModel),
				ResourceName:                         "pingfederate_server_settings_general." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "disable_automatic_connection_validation",
			},
			{
				// Back to minimal model
				Config: testAccServerSettingsGeneral(resourceName, nil),
				Check:  testAccCheckExpectedServerSettingsGeneralAttributes(nil),
			},
		},
	})
}

func testAccServerSettingsGeneral(resourceName string, resourceModel *serverSettingsGeneralResourceModel) string {
	optionalHcl := ""
	if resourceModel != nil {
		optionalHcl = fmt.Sprintf(`
		datastore_validation_interval_secs          = %d
		disable_automatic_connection_validation     = %t
		idp_connection_transaction_logging_override = "%s"
		request_header_for_correlation_id           = "%s"
		sp_connection_transaction_logging_override  = "%s"
		`, resourceModel.datastoreValidationIntervalSecs,
			resourceModel.disableAutomaticConnectionValidation,
			resourceModel.idpConnectionTransactionLoggingOverride,
			resourceModel.requestHeaderForCorrelationId,
			resourceModel.spConnectionTransactionLoggingOverride)
	}
	return fmt.Sprintf(`
resource "pingfederate_server_settings_general" "%s" {
	%s
}
data "pingfederate_server_settings_general" "%[1]s" {
  depends_on = [pingfederate_server_settings_general.%[1]s]
}`, resourceName,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsGeneralAttributes(config *serverSettingsGeneralResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettingsGeneral"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ServerSettingsAPI.GetGeneralSettings(ctx).Execute()

		if err != nil {
			return err
		}

		if config == nil {
			return nil
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, nil, "datastore_validation_interval_secs",
			config.datastoreValidationIntervalSecs, *response.DatastoreValidationIntervalSecs)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "disable_automatic_connection_validation",
			config.disableAutomaticConnectionValidation, *response.DisableAutomaticConnectionValidation)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "idp_connection_transaction_logging_override",
			config.idpConnectionTransactionLoggingOverride, *response.IdpConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "request_header_for_correlation_id",
			config.requestHeaderForCorrelationId, *response.RequestHeaderForCorrelationId)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "sp_connection_transaction_logging_override",
			config.spConnectionTransactionLoggingOverride, *response.SpConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}

		return nil
	}
}
