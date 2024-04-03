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
type serverSettingsGeneralSettingsResourceModel struct {
	disableAutomaticConnectionValidation    bool
	idpConnectionTransactionLoggingOverride string
	spConnectionTransactionLoggingOverride  string
	datastoreValidationIntervalSecs         int64
	requestHeaderForCorrelationId           string
}

func TestAccServerSettingsGeneralSettings(t *testing.T) {
	resourceName := "myServerSettingsGeneralSettings"
	updatedResourceModel := serverSettingsGeneralSettingsResourceModel{
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
				Config: testAccServerSettingsGeneralSettings(resourceName, nil),
				Check:  testAccCheckExpectedServerSettingsGeneralSettingsAttributes(nil),
			},
			{
				// Test updating some fields
				Config: testAccServerSettingsGeneralSettings(resourceName, &updatedResourceModel),
				Check:  testAccCheckExpectedServerSettingsGeneralSettingsAttributes(&updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccServerSettingsGeneralSettings(resourceName, &updatedResourceModel),
				ResourceName:      "pingfederate_server_settings_general_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccServerSettingsGeneralSettings(resourceName, nil),
				Check:  testAccCheckExpectedServerSettingsGeneralSettingsAttributes(nil),
			},
		},
	})
}

func testAccServerSettingsGeneralSettings(resourceName string, resourceModel *serverSettingsGeneralSettingsResourceModel) string {
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
resource "pingfederate_server_settings_general_settings" "%s" {
	%s
}
data "pingfederate_server_settings_general_settings" "%[1]s" {
  depends_on = [pingfederate_server_settings_general_settings.%[1]s]
}`, resourceName,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsGeneralSettingsAttributes(config *serverSettingsGeneralSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettingsGeneralSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		stateAttributes := s.RootModule().Resources["pingfederate_server_settings_general_settings.myServerSettingsGeneralSettings"].Primary.Attributes
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

		err = acctest.VerifyStateAttributeValue(stateAttributes, "datastore_validation_interval_secs", config.datastoreValidationIntervalSecs)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "disable_automatic_connection_validation",
			config.disableAutomaticConnectionValidation, *response.DisableAutomaticConnectionValidation)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "disable_automatic_connection_validation", config.disableAutomaticConnectionValidation)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "idp_connection_transaction_logging_override",
			config.idpConnectionTransactionLoggingOverride, *response.IdpConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "idp_connection_transaction_logging_override", config.idpConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "request_header_for_correlation_id",
			config.requestHeaderForCorrelationId, *response.RequestHeaderForCorrelationId)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "request_header_for_correlation_id", config.requestHeaderForCorrelationId)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "sp_connection_transaction_logging_override",
			config.spConnectionTransactionLoggingOverride, *response.SpConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "sp_connection_transaction_logging_override", config.spConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}

		return nil
	}
}
