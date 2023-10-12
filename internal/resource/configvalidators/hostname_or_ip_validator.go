package configvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &hostnameOrIpValidator{}

type hostnameOrIpValidator struct{}

func (v hostnameOrIpValidator) Description(ctx context.Context) string {
	return "Validates supplied value is of a hostname or IP type"
}

func (v hostnameOrIpValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v hostnameOrIpValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	isMatch, _ := regexp.MatchString(`^([a-z0-9]+(-[a-z0-9]+)*\.)*[a-z0-9]+(-[a-z0-9]+)*$`, strVal)
	if !isMatch {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid hostname or IP",
			fmt.Sprintf("This %s value must be a valid hostname or IP address", req.ConfigValue),
		)
	}
}

func ValidHostnameOrIp() hostnameOrIpValidator {
	return hostnameOrIpValidator{}
}
