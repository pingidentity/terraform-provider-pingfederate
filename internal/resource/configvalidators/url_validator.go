package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &urlValidator{}

type urlValidator struct{}

func (v urlValidator) Description(ctx context.Context) string {
	return "This value must start with 'http://' or 'https://'"
}

func (v urlValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v urlValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString(`^(https?:\/\/)`, strVal)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL Format",
			fmt.Sprintf("This %s value must start with 'http://' or 'https://'", req.ConfigValue),
		)
		return
	}
}

func ValidUrl() urlValidator {
	return urlValidator{}
}
