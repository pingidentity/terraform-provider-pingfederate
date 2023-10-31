package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = &urlValidator{}
var _ validator.List = &urlListValidator{}
var errorInfo = "%s must start with 'http://' or 'https://'"
var regexFilter = `^(https?:\/\/)`

type urlValidator struct{}
type urlListValidator struct{}

// URL String Validator
func (v urlValidator) Description(ctx context.Context) string {
	return "Validates the value supplied is of URL format"
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
	isMatch, _ := regexp.MatchString(regexFilter, strVal)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL Format",
			fmt.Sprintf(errorInfo, req.ConfigValue),
		)
	}
}

func ValidUrl() urlValidator {
	return urlValidator{}
}

// Check values in List for URL Validation
func (v urlListValidator) Description(ctx context.Context) string {
	return "Validates the value supplied is of URL format"
}

func (v urlListValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v urlListValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	listElems := req.ConfigValue.Elements()
	for _, elem := range listElems {
		regexComp, regexCompError := regexp.Compile(regexFilter)
		if regexCompError != nil {
			return
		}
		isMatch := regexComp.MatchString(elem.(types.String).ValueString())
		if !isMatch {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid URL Format Found In List",
				fmt.Sprintf(errorInfo, elem.(types.String).ValueString()),
			)
		}
	}

}

func ValidUrls() urlListValidator {
	return urlListValidator{}
}
