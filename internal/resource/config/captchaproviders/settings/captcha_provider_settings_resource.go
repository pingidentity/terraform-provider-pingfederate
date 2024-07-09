package captchaproviderssettings

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (m *captchaProviderSettingsResourceModel) setNullObjectValues() {
	refAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	m.DefaultCaptchaProviderRef = types.ObjectNull(refAttrTypes)
}
