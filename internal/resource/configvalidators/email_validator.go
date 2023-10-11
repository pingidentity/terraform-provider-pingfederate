package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &emailValidator{}

type emailValidator struct{}

func (v emailValidator) Description(ctx context.Context) string {
	return "Validates value supplied is of E-mail address value format"
}

func (v emailValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v emailValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, strVal)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid E-mail Address",
			fmt.Sprintf("The email %s must be of the form '<address>@<company>.<domain>', where 'domain' contains only alphabetic characters and is at least 2 characters in length.", req.ConfigValue),
		)
	}
}

func ValidEmail() emailValidator {
	return emailValidator{}
}
