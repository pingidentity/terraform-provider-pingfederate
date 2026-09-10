// Copyright © 2026 Ping Identity Corporation

package configvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.Map = &attributeContractFulfillmentValidator{}

type attributeContractFulfillmentValidator struct{}

func (v attributeContractFulfillmentValidator) Description(ctx context.Context) string {
	return "Validates that the any `value` defined in the attribute contract fulfillment is set appropriately according to the source type, " +
		"and that EXPRESSION-type values do not end in a trailing newline."
}

func (v attributeContractFulfillmentValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v attributeContractFulfillmentValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	for key, attr := range req.ConfigValue.Elements() {
		// Determine the source type
		attrAsObject, ok := attr.(types.Object)
		if !ok {
			continue
		}
		source, ok := attrAsObject.Attributes()["source"]
		if !ok {
			continue
		}
		sourceAsObject, ok := source.(types.Object)
		if !ok {
			continue
		}
		sourceType, ok := sourceAsObject.Attributes()["type"]
		if !ok {
			continue
		}
		sourceTypeAsString, ok := sourceType.(types.String)
		if !ok || sourceTypeAsString.IsUnknown() {
			continue
		}

		// Get the value
		valueAsString, ok := attrAsObject.Attributes()["value"].(types.String)
		if !ok || valueAsString.IsUnknown() {
			continue
		}

		// If source type is 'NO_MAPPING', then value must not be defined. Otherwise, value must be defined.
		if sourceTypeAsString.ValueString() == "NO_MAPPING" && len(valueAsString.ValueString()) > 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				providererror.InvalidAttributeConfiguration,
				"When attribute_contract_fulfillment source type is set to 'NO_MAPPING', the value must not be defined. "+
					fmt.Sprintf("attribute_contract_fulfillment key '%s' has a value defined while using a source type of 'NO_MAPPING'", key),
			)
		} else if sourceTypeAsString.ValueString() != "NO_MAPPING" && len(valueAsString.ValueString()) == 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				providererror.InvalidAttributeConfiguration,
				"When attribute_contract_fulfillment source type is set to anything other than 'NO_MAPPING', the value must be defined. "+
					fmt.Sprintf("attribute_contract_fulfillment key '%s' has no value defined while using a source type of '%s'", key, sourceTypeAsString.ValueString()),
			)
		}

		// If source type is 'EXPRESSION', the value must not end in a trailing newline. PingFederate trims
		// trailing newlines when parsing OGNL expressions, so a value ending in "\n" (or "\r\n") in Terraform
		// config would not match the value returned by the PingFederate admin API after apply.
		if sourceTypeAsString.ValueString() == "EXPRESSION" && strings.HasSuffix(valueAsString.ValueString(), "\n") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				providererror.InvalidAttributeConfiguration,
				"When attribute_contract_fulfillment source type is set to 'EXPRESSION', the value must not end in a trailing newline. "+
					fmt.Sprintf("attribute_contract_fulfillment key '%s' has a value ending in a trailing newline while using a source type of 'EXPRESSION'. ", key)+
					"This commonly happens when the expression is written using Terraform heredoc syntax (e.g. `<<-EOT ... EOT`), which appends a trailing newline to the string. "+
					"Remove the trailing newline from the value, for example by using a regular quoted string instead of a heredoc, or by wrapping the heredoc value in the `chomp()` function.",
			)
		}
	}
}

func ValidAttributeContractFulfillment() attributeContractFulfillmentValidator {
	return attributeContractFulfillmentValidator{}
}
