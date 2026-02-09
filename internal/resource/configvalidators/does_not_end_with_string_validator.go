// Copyright © 2026 Ping Identity Corporation

package configvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.String = &doesNotEndWithValidator{}

type doesNotEndWithValidator struct {
	suffix string
}

func (v doesNotEndWithValidator) Description(ctx context.Context) string {
	return "Validates value supplied does not end with a given suffix"
}

func (v doesNotEndWithValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v doesNotEndWithValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch := strings.HasSuffix(strVal, v.suffix)
	if isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("%s must not be suffixed with a %s", req.ConfigValue, v.suffix),
		)
	}
}

func DoesNotEndWith(suffix string) doesNotEndWithValidator {
	return doesNotEndWithValidator{
		suffix: suffix,
	}
}
