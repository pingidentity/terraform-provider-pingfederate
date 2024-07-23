package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &pingFederateIdValidator{}

type pingFederateIdValidator struct{}

func (v pingFederateIdValidator) Description(ctx context.Context) string {
	return "Validates supplied value is a valid id using characters [a-zA-Z0-9._-]"
}

func (v pingFederateIdValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v pingFederateIdValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString(`[^a-zA-Z0-9._-]`, strVal)
	if isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Contains invalid characters",
			fmt.Sprintf("%s must contain only characters in [a-zA-Z0-9._-]", req.ConfigValue),
		)
		return
	}
}

func PingFederateId() pingFederateIdValidator {
	return pingFederateIdValidator{}
}
