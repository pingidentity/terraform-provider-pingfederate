package acctest_test

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

// Function to build out test object HCL; necessary for null value testing
func incomingProxySettingsHCLObj(incomingProxySettings *client.IncomingProxySettings) string {
	var forwardedIpAddressHeaderName string
	if incomingProxySettings.GetForwardedIpAddressHeaderName() == "" {
		forwardedIpAddressHeaderName = ""
	} else {
		forwardedIpAddressHeaderName = fmt.Sprintf("\tforwarded_ip_address_header_name = \"%s\"", *incomingProxySettings.ForwardedIpAddressHeaderName)
	}

	var forwardedIpAddressHeaderIndex string
	if incomingProxySettings.GetForwardedIpAddressHeaderIndex() == "" {
		forwardedIpAddressHeaderIndex = ""
	} else {
		forwardedIpAddressHeaderIndex = fmt.Sprintf("\tforwarded_ip_address_header_index = \"%s\"", *incomingProxySettings.ForwardedIpAddressHeaderIndex)
	}

	var forwardedHostHeaderName string
	if incomingProxySettings.GetForwardedHostHeaderName() == "" {
		forwardedHostHeaderName = ""
	} else {
		forwardedHostHeaderName = fmt.Sprintf("\tforwarded_host_header_name = \"%s\"", *incomingProxySettings.ForwardedHostHeaderName)
	}

	var forwardedHostHeaderIndex string
	if incomingProxySettings.GetForwardedHostHeaderIndex() == "" {
		forwardedHostHeaderIndex = ""
	} else {
		forwardedHostHeaderIndex = fmt.Sprintf("\tforwarded_host_header_index = \"%s\"", *incomingProxySettings.ForwardedHostHeaderIndex)
	}

	var clientCertSSLHeaderName string
	if incomingProxySettings.GetClientCertSSLHeaderName() == "" {
		clientCertSSLHeaderName = ""
	} else {
		clientCertSSLHeaderName = fmt.Sprintf("\tclient_cert_ssl_header_name = \"%s\"", *incomingProxySettings.ClientCertSSLHeaderName)
	}

	var clientCertChainSSLHeaderName string
	if incomingProxySettings.GetClientCertChainSSLHeaderName() == "" {
		clientCertChainSSLHeaderName = ""
	} else {
		clientCertChainSSLHeaderName = fmt.Sprintf("\tclient_cert_chain_ssl_header_name = \"%s\"", *incomingProxySettings.ClientCertChainSSLHeaderName)
	}

	var proxyTerminatesHttpsConns string
	// GetProxyTerminatesHttpsConns returns the ProxyTerminatesHttpsConns field value if set, zero value (false) otherwise.
	if incomingProxySettings.GetProxyTerminatesHttpsConns() == false {
		return ""
	} else {
		proxyTerminatesHttpsConns = fmt.Sprintf("\tproxy_terminates_https_conns = %t", *incomingProxySettings.ProxyTerminatesHttpsConns)
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s",
		forwardedIpAddressHeaderName,
		forwardedIpAddressHeaderIndex,
		forwardedHostHeaderName,
		forwardedHostHeaderIndex,
		clientCertSSLHeaderName,
		clientCertChainSSLHeaderName,
		proxyTerminatesHttpsConns,
	)
}

func TestAccIncomingProxySettings(t *testing.T) {
	resourceName := "myIncomingProxySettings"
	initialIncomingProxySettings := client.NewIncomingProxySettings()
	initialResourceModel := initialIncomingProxySettings
	updatedIncomingProxySettings := client.NewIncomingProxySettings()
	updatedIncomingProxySettings.ForwardedIpAddressHeaderName = pointers.String("Updated-X-Forwarded-For")
	updatedIncomingProxySettings.ForwardedIpAddressHeaderIndex = pointers.String("LAST")
	updatedIncomingProxySettings.ForwardedHostHeaderName = pointers.String("Updated-X-Forwarded-Host")
	updatedIncomingProxySettings.ForwardedHostHeaderIndex = pointers.String("FIRST")
	updatedIncomingProxySettings.ClientCertSSLHeaderName = pointers.String("Updated-X-Client-Cert")
	updatedIncomingProxySettings.ClientCertChainSSLHeaderName = pointers.String("Updated-X-Client-Cert-Chain")
	updatedIncomingProxySettings.ProxyTerminatesHttpsConns = pointers.Bool(true)
	updatedResourceModel := updatedIncomingProxySettings

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test empty object sent
			{
				Config: testAccIncomingProxySettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIncomingProxySettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIncomingProxySettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIncomingProxySettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccIncomingProxySettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_incoming_proxy_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test empty object sent
			{
				Config: testAccIncomingProxySettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIncomingProxySettingsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccIncomingProxySettings(resourceName string, incomingProxySettings *client.IncomingProxySettings) string {
	return fmt.Sprintf(`
resource "pingfederate_incoming_proxy_settings" "%[1]s" {
	%[2]s
}`, resourceName,
		incomingProxySettingsHCLObj(incomingProxySettings),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedIncomingProxySettingsAttributes(incomingProxySettings *client.IncomingProxySettings) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "IncomingProxySettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.IncomingProxySettingsAPI.GetIncomingProxySettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		if incomingProxySettings.ForwardedIpAddressHeaderName != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "forwarded_ip_address_header_name", *incomingProxySettings.ForwardedIpAddressHeaderName, *response.ForwardedIpAddressHeaderName)
			if err != nil {
				return err
			}
		}

		if incomingProxySettings.ForwardedIpAddressHeaderIndex != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "forwarded_ip_address_header_index", *incomingProxySettings.ForwardedIpAddressHeaderIndex, *response.ForwardedIpAddressHeaderIndex)
			if err != nil {
				return err
			}
		}

		if incomingProxySettings.ForwardedHostHeaderName != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "forwarded_host_header_name", *incomingProxySettings.ForwardedHostHeaderName, *response.ForwardedHostHeaderName)
			if err != nil {
				return err
			}
		}

		if incomingProxySettings.ForwardedHostHeaderIndex != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "forwarded_host_header_index", *incomingProxySettings.ForwardedHostHeaderIndex, *response.ForwardedHostHeaderIndex)
			if err != nil {
				return err
			}
		}

		if incomingProxySettings.ClientCertSSLHeaderName != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "client_cert_ssl_header_name", *incomingProxySettings.ClientCertSSLHeaderName, *response.ClientCertSSLHeaderName)
			if err != nil {
				return err
			}
		}

		if incomingProxySettings.ClientCertChainSSLHeaderName != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "client_cert_chain_ssl_header_name", *incomingProxySettings.ClientCertChainSSLHeaderName, *response.ClientCertChainSSLHeaderName)
			if err != nil {
				return err
			}
		}

		if incomingProxySettings.ProxyTerminatesHttpsConns != nil {
			err = acctest.TestAttributesMatchBool(resourceType, nil, "proxy_terminates_https_conns", *incomingProxySettings.ProxyTerminatesHttpsConns, *response.ProxyTerminatesHttpsConns)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
