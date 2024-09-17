package configvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.String = &validPathValidator{}

type validPathValidator struct{}

func (v validPathValidator) Description(ctx context.Context) string {
	return "Validates value supplied starts with `/` but does not end with `/`."
}

func (v validPathValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v validPathValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	if !strings.HasPrefix(strVal, "/") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("%s must be prefixed with a '/'", req.ConfigValue),
		)
	}

	if strings.HasSuffix(strVal, "/") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("%s must not end with a '/'", req.ConfigValue),
		)
	}
}

func ValidPath() validPathValidator {
	return validPathValidator{}
}
