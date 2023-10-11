package configvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &startsWithValidator{}

type startsWithValidator struct {
	firstChar string
}

func (v startsWithValidator) Description(ctx context.Context) string {
	return "Validates value supplied does not contain any whitespaces"
}

func (v startsWithValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v startsWithValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch := strings.HasPrefix(strVal, v.firstChar)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Prefix Value",
			fmt.Sprintf("%s must be prefixed with a %s", req.ConfigValue, v.firstChar),
		)
	}
}

func StartsWith(char string) startsWithValidator {
	return startsWithValidator{
		firstChar: char,
	}
}
