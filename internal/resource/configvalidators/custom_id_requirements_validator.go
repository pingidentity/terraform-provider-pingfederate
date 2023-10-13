package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &customIdReqValidator{}

type customIdReqValidator struct{}

func (v customIdReqValidator) Description(ctx context.Context) string {
	return "Verifies custom_id contains more than 33 characters, contain no spaces, and be alphanumeric"
}

func (v customIdReqValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v customIdReqValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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
			fmt.Sprintf("The custom_id of %s must be less than 33 characters, contain no spaces, and be alphanumeric", req.ConfigValue),
		)
	}
}

func ValidChars() customIdReqValidator {
	return customIdReqValidator{}
}
