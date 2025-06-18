// Copyright © 2025 Ping Identity Corporation

package configvalidators

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.String = &urlValidator{}
var _ validator.List = &urlListValidator{}
var _ validator.Set = &urlSetValidator{}

type urlValidator struct{}
type urlListValidator struct{}
type urlSetValidator struct{}

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
	valueString := value.ValueString()

	// Rely on PingFederate validation for values containing asterisk(s)
	if strings.Contains(valueString, "*") {
		return
	}

	_, err := url.Parse(valueString)
	if err != nil {
		respDiags.AddAttributeError(
			path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("Invalid URL Format for '%s': %s", value.ValueString(), err.Error()),
		)
	}
}

func ValidUrl() urlValidator {
	return urlValidator{}
}

// Check values in List for URL Validation
func (v urlListValidator) Description(ctx context.Context) string {
	return "Validates each value in the list is of URL format"
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
			return
		}
		validateUrlValue(req.Path, elemString, &resp.Diagnostics)
	}

}

func ValidUrlsList() urlListValidator {
	return urlListValidator{}
}

// Check values in Set for URL Validation
func (v urlSetValidator) Description(ctx context.Context) string {
	return "Validates each value in the set is of URL format"
}

func (v urlSetValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v urlSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	setElems := req.ConfigValue.Elements()
	for _, elem := range setElems {
		elemString, ok := elem.(types.String)
		if !ok || elemString.IsUnknown() {
			return
		}
		validateUrlValue(req.Path, elemString, &resp.Diagnostics)
	}

}

func ValidUrlsSet() urlSetValidator {
	return urlSetValidator{}
}
