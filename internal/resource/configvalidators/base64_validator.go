package configvalidators

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
)

var _ validator.String = &base64Validator{}

type base64Validator struct{}

func (v base64Validator) Description(ctx context.Context) string {
	return "Validates value supplied is base64-encoded"
}

func (v base64Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v base64Validator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	strVal := req.ConfigValue.ValueString()
	target := make([]byte, base64.StdEncoding.DecodedLen(len(strVal)))
	_, err := base64.StdEncoding.Decode(target, []byte(strVal))
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("The value must be base64-encoded. Error when attempting to decode: %s", err.Error()),
		)
	}
}

func ValidBase64() base64Validator {
	return base64Validator{}
}
