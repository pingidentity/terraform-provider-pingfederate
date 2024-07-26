package configvalidators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = &urlValidator{}
var _ validator.List = &urlListValidator{}

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
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() || req.ConfigValue.ValueString() == "" {
		return
	}
	validateUrlValue(req.Path, req.ConfigValue, &resp.Diagnostics)
}

func validateUrlValue(path path.Path, value types.String, respDiags *diag.Diagnostics) {
	// Ensure the the URL can be parsed by url.Parse
	_, err := url.Parse(value.ValueString())
	if err != nil {
		respDiags.AddAttributeError(
			path,
			fmt.Sprintf("Invalid URL Format for '%s'", value.ValueString()),
			err.Error(),
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
		elemString, ok := elem.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"URL Validation can only be applied to a list of strings",
				"",
			)
			return
		}
		validateUrlValue(req.Path, elemString, &resp.Diagnostics)
	}

}

func ValidUrls() urlListValidator {
	return urlListValidator{}
}
