// Copyright © 2026 Ping Identity Corporation

package resource_sp_idp_connection_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/accesstokenmanager"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// Reproduces CDI-1474: an empty issuance_criteria.expression_criteria set caused a
// "Provider produced inconsistent result after apply" error, since the PingFederate admin
// API omits an empty expressionCriteria element rather than returning it as an empty list.
func TestAccSpIdpConnection_ExpressionCriteriaEmptySetInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config:      spIdpConnection_ExpressionCriteriaEmptySetInvalidHCL(),
				ExpectError: regexp.MustCompile(`set must contain at least 1 elements`),
			},
		},
	})
}

func spIdpConnection_ExpressionCriteriaEmptySetInvalidHCL() string {
	return fmt.Sprintf(`
%s

resource "pingfederate_sp_idp_connection" "example" {
  connection_id = "expressioncriteriaemptyconn"
  credentials = {
  }
  entity_id = "docker"
  idp_oauth_grant_attribute_mapping = {
    access_token_manager_mappings = [
      {
        access_token_manager_ref = {
          id = pingfederate_oauth_access_token_manager.idpConnValidationAtm.id
        }
        attribute_contract_fulfillment = {
          Username = {
            source = {
              type = "ASSERTION"
            }
            value = "TOKEN_SUBJECT"
          }
        }
        issuance_criteria = {
          expression_criteria = []
        }
      },
    ]
    idp_oauth_attribute_contract = {
    }
  }
  name = "docker"
}
`, accesstokenmanager.AccessTokenManagerTestHCL("idpConnValidationAtm"))
}

// Reproduces CDI-1474: an EXPRESSION-type attribute_contract_fulfillment value written with
// Terraform heredoc syntax (<<-EOT) ends in a trailing newline, which the PingFederate admin
// API trims when parsing the OGNL expression, causing a "Provider produced inconsistent result
// after apply" error.
func TestAccSpIdpConnection_AttributeContractFulfillmentTrailingNewlineInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config:      spIdpConnection_AttributeContractFulfillmentTrailingNewlineInvalidHCL(),
				ExpectError: regexp.MustCompile(`must not end in a trailing newline`),
			},
		},
	})
}

func spIdpConnection_AttributeContractFulfillmentTrailingNewlineInvalidHCL() string {
	return fmt.Sprintf(`
%s

resource "pingfederate_sp_idp_connection" "example" {
  connection_id = "acftrailingnewlineconn"
  credentials = {
  }
  entity_id = "docker"
  idp_oauth_grant_attribute_mapping = {
    access_token_manager_mappings = [
      {
        access_token_manager_ref = {
          id = pingfederate_oauth_access_token_manager.idpConnValidationAtm.id
        }
        attribute_contract_fulfillment = {
          Username = {
            source = {
              type = "ASSERTION"
            }
            value = "TOKEN_SUBJECT"
          }
          aud = {
            source = {
              type = "EXPRESSION"
            }
            value = <<-EOT
              #this.get("context.HttpRequest").getObjectValue().getParameter("resource")
            EOT
          }
        }
      },
    ]
    idp_oauth_attribute_contract = {
      extended_attributes = [
        {
          masked = false
          name   = "aud"
        },
      ]
    }
  }
  name = "docker"
}
`, accesstokenmanager.AccessTokenManagerTestHCL("idpConnValidationAtm"))
}
