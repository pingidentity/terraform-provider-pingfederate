// Copyright © 2026 Ping Identity Corporation

package configvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.String = &startsWithValidator{}

type startsWithValidator struct {
	prefix string
}

func (v startsWithValidator) Description(ctx context.Context) string {
	return "Validates value supplied starts with a given prefix"
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
	isMatch := strings.HasPrefix(strVal, v.prefix)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("%s must be prefixed with a %s", req.ConfigValue, v.prefix),
		)
	}
}

func StartsWith(prefix string) startsWithValidator {
	return startsWithValidator{
		prefix: prefix,
	}
}
