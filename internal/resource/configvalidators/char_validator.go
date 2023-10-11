package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &charValidator{}

type charValidator struct{}

func (v charValidator) Description(ctx context.Context) string {
	return "The ID must be less than 33 characters, contain no spaces, and be alphanumeric"
}

func (v charValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v charValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString("^[a-zA-Z0-9_]{1,32}$", strVal)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid custom_id Value",
			fmt.Sprintf("The ID of %s must be less than 33 characters, contain no spaces, and be alphanumeric", req.ConfigValue),
		)
		return
	}
}

func ValidChars() charValidator {
	return charValidator{}
}
