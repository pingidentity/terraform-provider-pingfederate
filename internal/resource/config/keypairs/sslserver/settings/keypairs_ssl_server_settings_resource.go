package keypairssslserversettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

func (m *keypairsSslServerSettingsResourceModel) setNullObjectValues() {
	certRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	m.ActiveAdminConsoleCerts = types.SetNull(types.ObjectType{AttrTypes: certRefAttrTypes})
	m.ActiveRuntimeServerCerts = types.SetNull(types.ObjectType{AttrTypes: certRefAttrTypes})
	m.AdminConsoleCertRef = types.ObjectNull(certRefAttrTypes)
	m.RuntimeServerCertRef = types.ObjectNull(certRefAttrTypes)
}

func (r *keypairsSslServerSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config keypairsSslServerSettingsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if internaltypes.IsDefined(config.ActiveAdminConsoleCerts) && internaltypes.IsDefined(config.AdminConsoleCertRef) {
		certRefFound := false
		adminConsoleCertId := config.AdminConsoleCertRef.Attributes()["id"].(types.String)
		if internaltypes.IsDefined(adminConsoleCertId) {
			for _, cert := range config.ActiveAdminConsoleCerts.Elements() {
				certId := cert.(types.Object).Attributes()["id"].(types.String)
				if certId.Equal(adminConsoleCertId) {
					certRefFound = true
					break
				}
			}
			if !certRefFound {
				resp.Diagnostics.AddAttributeError(
					path.Root("active_admin_console_certs"),
					providererror.InvalidAttributeConfiguration,
					fmt.Sprintf("`admin_console_cert_ref.id` '%s' must be included in `active_admin_console_certs`", adminConsoleCertId.ValueString()))
			}
		}
	}

	if internaltypes.IsDefined(config.ActiveRuntimeServerCerts) && internaltypes.IsDefined(config.RuntimeServerCertRef) {
		certRefFound := false
		runtimeServerCertId := config.RuntimeServerCertRef.Attributes()["id"].(types.String)
		if internaltypes.IsDefined(runtimeServerCertId) {
			for _, cert := range config.ActiveRuntimeServerCerts.Elements() {
				certId := cert.(types.Object).Attributes()["id"].(types.String)
				if certId.Equal(runtimeServerCertId) {
					certRefFound = true
					break
				}
			}
			if !certRefFound {
				resp.Diagnostics.AddAttributeError(
					path.Root("active_runtime_server_certs"),
					providererror.InvalidAttributeConfiguration,
					fmt.Sprintf("`runtime_server_cert_ref.id` '%s' must be included in `active_runtime_server_certs`", runtimeServerCertId.ValueString()))
			}
		}
	}
}
