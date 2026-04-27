// Copyright © 2026 Ping Identity Corporation

package passwordcredentialvalidator_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// This test validates that if a row is added to a plugin configuration table outside of Terraform,
// a non-empty plan is generated reflecting the drift, with no panic or error.
func TestAccPasswordCredentialValidator_RadiusTableRowsDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: passwordCredentialValidator_CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: passwordCredentialValidator_RadiusCompleteHCL(),
			},
			{
				PreConfig: func() {
					passwordCredentialValidator_AddRadiusServerRow(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func passwordCredentialValidator_AddRadiusServerRow(t *testing.T) {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()

	validator, _, err := testClient.PasswordCredentialValidatorsAPI.GetPasswordCredentialValidator(ctx, validatorId).Execute()
	if err != nil {
		t.Fatalf("Failed to get password credential validator: %v", err)
	}

	tableIndex := -1
	for i, table := range validator.Configuration.Tables {
		if table.Name == "RADIUS Servers" {
			tableIndex = i
			break
		}
	}

	if tableIndex < 0 {
		t.Fatal("Failed to find RADIUS Servers table in password credential validator configuration")
	}

	newRow := client.ConfigRow{Fields: []client.ConfigField{
		{Name: "Hostname", Value: client.PtrString("localhost2")},
		{Name: "Authentication Port", Value: client.PtrString("1812")},
		{Name: "Authentication Protocol", Value: client.PtrString("PAP")},
		{Name: "Shared Secret", Value: client.PtrString("mysecret2")},
	}}

	validator.Configuration.Tables[tableIndex].Rows = append(validator.Configuration.Tables[tableIndex].Rows, newRow)

	updateRequest := testClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidator(ctx, validatorId)
	updateRequest = updateRequest.Body(*validator)
	_, _, err = testClient.PasswordCredentialValidatorsAPI.UpdatePasswordCredentialValidatorExecute(updateRequest)
	if err != nil {
		t.Fatalf("Failed to update password credential validator with additional row: %v", err)
	}
}
