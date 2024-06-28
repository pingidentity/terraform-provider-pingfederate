package keypairssslserversettings

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (m *keypairsSslServerSettingsResourceModel) setNullObjectValues() {
	certRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	m.ActiveAdminConsoleCerts = types.ListNull(types.ObjectType{AttrTypes: certRefAttrTypes})
	m.ActiveRuntimeServerCerts = types.ListNull(types.ObjectType{AttrTypes: certRefAttrTypes})
	m.AdminConsoleCertRef = types.ObjectNull(certRefAttrTypes)
	m.RuntimeServerCertRef = types.ObjectNull(certRefAttrTypes)
}
