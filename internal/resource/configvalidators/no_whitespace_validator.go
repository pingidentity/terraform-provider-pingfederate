package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.String = &noWhitespaceValidator{}

type noWhitespaceValidator struct{}

func (v noWhitespaceValidator) Description(ctx context.Context) string {
	return "Validates supplied value contains no whitespaces"
}

func (v noWhitespaceValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v noWhitespaceValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString(`^\S*$`, strVal)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("%s must not contain any whitespace", req.ConfigValue),
		)
		return
	}
}

func NoWhitespace() noWhitespaceValidator {
	return noWhitespaceValidator{}
}
