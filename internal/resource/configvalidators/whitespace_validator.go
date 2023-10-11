package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &whitespaceValidator{}

type whitespaceValidator struct{}

func (v whitespaceValidator) Description(ctx context.Context) string {
	return "This value must not contain any whitespace"
}

func (v whitespaceValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v whitespaceValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString(`^\S*$`, strVal)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Contains whitespace",
			fmt.Sprintf("%s must not contain any whitespace", req.ConfigValue),
		)
		return
	}
}

func NoWhitespace() whitespaceValidator {
	return whitespaceValidator{}
}
