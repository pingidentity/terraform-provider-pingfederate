package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.String = &lowercaseIdValidator{}

type lowercaseIdValidator struct{}

func (v lowercaseIdValidator) Description(ctx context.Context) string {
	return "Validates supplied value is a valid id using characters [a-z0-9._-]"
}

func (v lowercaseIdValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v lowercaseIdValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString(`[^a-z0-9._-]`, strVal)
	if isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("%s must contain only characters in [a-z0-9._-]", req.ConfigValue),
		)
		return
	}
}

func LowercaseId() lowercaseIdValidator {
	return lowercaseIdValidator{}
}
