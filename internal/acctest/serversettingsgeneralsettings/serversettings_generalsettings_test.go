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

const serverSettingsGeneralSettingsId = "id"

// Attributes to test with. Add optional properties to test here if desired.
type serverSettingsGeneralSettingsResourceModel struct {
	id                                      string
	disableAutomaticConnectionValidation    bool
	idpConnectionTransactionLoggingOverride string
	spConnectionTransactionLoggingOverride  string
	datastoreValidationIntervalSecs         int64
	requestHeaderForCorrelationId           string
}

func TestAccServerSettingsGeneralSettings(t *testing.T) {
	resourceName := "myServerSettingsGeneralSettings"
	initialResourceModel := serverSettingsGeneralSettingsResourceModel{
		disableAutomaticConnectionValidation:    false,
		idpConnectionTransactionLoggingOverride: "NONE",
		spConnectionTransactionLoggingOverride:  "FULL",
		datastoreValidationIntervalSecs:         299,
		requestHeaderForCorrelationId:           "example",
	}
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
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServerSettingsGeneralSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedServerSettingsGeneralSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccServerSettingsGeneralSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedServerSettingsGeneralSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccServerSettingsGeneralSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_serversettings_generalsettings." + resourceName,
				ImportStateId:     serverSettingsGeneralSettingsId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServerSettingsGeneralSettings(resourceName string, resourceModel serverSettingsGeneralSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_serversettings_generalsettings" "%[1]s" {
  datastore_validation_interval_secs          = %[2]d
  disable_automatic_connection_validation     = %[3]t
  idp_connection_transaction_logging_override = "%[4]s"
  request_header_for_correlation_id           = "%[5]s"
  sp_connection_transaction_logging_override  = "%[6]s"
}`, resourceName,
		resourceModel.datastoreValidationIntervalSecs,
		resourceModel.disableAutomaticConnectionValidation,
		resourceModel.idpConnectionTransactionLoggingOverride,
		resourceModel.requestHeaderForCorrelationId,
		resourceModel.spConnectionTransactionLoggingOverride,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsGeneralSettingsAttributes(config serverSettingsGeneralSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettingsGeneralSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ServerSettingsApi.GetGeneralSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "datastore_validation_interval_secs",
			config.datastoreValidationIntervalSecs, *response.DatastoreValidationIntervalSecs)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "disable_automatic_connection_validation",
			config.disableAutomaticConnectionValidation, *response.DisableAutomaticConnectionValidation)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "idp_connection_transaction_logging_override",
			config.idpConnectionTransactionLoggingOverride, *response.IdpConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "request_header_for_correlation_id",
			config.requestHeaderForCorrelationId, *response.RequestHeaderForCorrelationId)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "sp_connection_transaction_logging_override",
			config.spConnectionTransactionLoggingOverride, *response.SpConnectionTransactionLoggingOverride)
		if err != nil {
			return err
		}

		return nil
	}
}
