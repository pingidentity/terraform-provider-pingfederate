// Copyright © 2026 Ping Identity Corporation

package configvalidators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringNotNullValidator{}

// StringNotNullValidator validates that an attribute is not null. Most
// attributes should set Required: true instead, however in certain scenarios,
// such as a computed nested attribute, all underlying attributes must also be
// computed for planning to not show unexpected differences.
type StringNotNullValidator struct{}

func (v StringNotNullValidator) Description(_ context.Context) string {
	return "The value cannot be null."
}

func (v StringNotNullValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v StringNotNullValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Missing Attribute Value",
		req.Path.String()+": "+v.Description(ctx),
	)
}

// StringNotNull returns an validator which ensures that the string attribute is
// not null. Most attributes should set Required: true instead, however in
// certain scenarios, such as a computed nested attribute, all underlying
// attributes must also be computed for planning to not show unexpected
// differences.
func StringNotNull() validator.String {
	return StringNotNullValidator{}
}
