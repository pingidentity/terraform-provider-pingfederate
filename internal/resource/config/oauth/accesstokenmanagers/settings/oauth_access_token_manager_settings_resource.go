package oauthaccesstokenmanagerssettings

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (m *oauthAccessTokenManagerSettingsResourceModel) setNullObjectValues() {
	refAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	m.DefaultAccessTokenManagerRef = types.ObjectNull(refAttrTypes)
}
