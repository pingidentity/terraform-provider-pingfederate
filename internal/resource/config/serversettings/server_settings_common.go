package serversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	contactInfoAttrType = map[string]attr.Type{
		"company":    types.StringType,
		"email":      types.StringType,
		"first_name": types.StringType,
		"last_name":  types.StringType,
		"phone":      types.StringType,
	}

	certificateExpirationsAttrType = map[string]attr.Type{
		"email_address":              types.StringType,
		"initial_warning_period":     types.Int64Type,
		"final_warning_period":       types.Int64Type,
		"notification_publisher_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"notification_mode":          types.StringType,
	}

	notificationSettingsAttrType = map[string]attr.Type{
		"email_address":              types.StringType,
		"notification_publisher_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}

	threadPoolExhaustionNotificationSettingsAttrType = map[string]attr.Type{
		"email_address":              types.StringType,
		"thread_dump_enabled":        types.BoolType,
		"notification_publisher_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"notification_mode":          types.StringType,
	}

	notificationsAttrType = map[string]attr.Type{
		"license_events":                                           types.ObjectType{AttrTypes: notificationSettingsAttrType},
		"certificate_expirations":                                  types.ObjectType{AttrTypes: certificateExpirationsAttrType},
		"notify_admin_user_password_changes":                       types.BoolType,
		"account_changes_notification_publisher_ref":               types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"metadata_notification_settings":                           types.ObjectType{AttrTypes: notificationSettingsAttrType},
		"expired_certificate_administrative_console_warning_days":  types.Int64Type,
		"expiring_certificate_administrative_console_warning_days": types.Int64Type,
		"thread_pool_exhaustion_notification_settings":             types.ObjectType{AttrTypes: threadPoolExhaustionNotificationSettingsAttrType},
	}

	oauthRoleAttrType = map[string]attr.Type{
		"enable_oauth":           types.BoolType,
		"enable_open_id_connect": types.BoolType,
	}

	idpSaml20ProfileAttrType = map[string]attr.Type{
		"enable":              types.BoolType,
		"enable_auto_connect": types.BoolType,
	}

	spSaml20ProfileAttrType = map[string]attr.Type{
		"enable":              types.BoolType,
		"enable_auto_connect": types.BoolType,
		"enable_xasp":         types.BoolType,
	}

	idpRoleAttrType = map[string]attr.Type{
		"enable":                       types.BoolType,
		"enable_saml_1_1":              types.BoolType,
		"enable_saml_1_0":              types.BoolType,
		"enable_ws_fed":                types.BoolType,
		"enable_ws_trust":              types.BoolType,
		"saml_2_0_profile":             types.ObjectType{AttrTypes: idpSaml20ProfileAttrType},
		"enable_outbound_provisioning": types.BoolType,
	}

	spRoleAttrType = map[string]attr.Type{
		"enable":                      types.BoolType,
		"enable_saml_1_1":             types.BoolType,
		"enable_saml_1_0":             types.BoolType,
		"enable_ws_fed":               types.BoolType,
		"enable_ws_trust":             types.BoolType,
		"saml_2_0_profile":            types.ObjectType{AttrTypes: spSaml20ProfileAttrType},
		"enable_open_id_connect":      types.BoolType,
		"enable_inbound_provisioning": types.BoolType,
	}

	rolesAndProtocolsAttrType = map[string]attr.Type{
		"oauth_role":           types.ObjectType{AttrTypes: oauthRoleAttrType},
		"idp_role":             types.ObjectType{AttrTypes: idpRoleAttrType},
		"sp_role":              types.ObjectType{AttrTypes: spRoleAttrType},
		"enable_idp_discovery": types.BoolType,
	}

	federationInfoAttrType = map[string]attr.Type{
		"base_url":               types.StringType,
		"saml_2_entity_id":       types.StringType,
		"auto_connect_entity_id": types.StringType,
		"saml_1x_issuer_id":      types.StringType,
		"saml_1x_source_id":      types.StringType,
		"wsfed_realm":            types.StringType,
	}

	emailServerAttrType = map[string]attr.Type{
		"source_addr":                 types.StringType,
		"email_server":                types.StringType,
		"port":                        types.Int64Type,
		"ssl_port":                    types.Int64Type,
		"timeout":                     types.Int64Type,
		"retry_attempts":              types.Int64Type,
		"retry_delay":                 types.Int64Type,
		"use_ssl":                     types.BoolType,
		"use_tls":                     types.BoolType,
		"verify_hostname":             types.BoolType,
		"enable_utf8_message_headers": types.BoolType,
		"use_debugging":               types.BoolType,
		"username":                    types.StringType,
		"password":                    types.StringType,
	}

	captchaSettingsAttrType = map[string]attr.Type{
		"site_key":   types.StringType,
		"secret_key": types.StringType,
	}

	contactInfoDefault, _ = types.ObjectValue(contactInfoAttrType, map[string]attr.Value{
		"company":    types.StringNull(),
		"email":      types.StringNull(),
		"first_name": types.StringNull(),
		"last_name":  types.StringNull(),
		"phone":      types.StringNull(),
	})

	notificationsDefault, _ = types.ObjectValue(notificationsAttrType, map[string]attr.Value{
		"license_events":                                           types.ObjectNull(notificationSettingsAttrType),
		"certificate_expirations":                                  types.ObjectNull(certificateExpirationsAttrType),
		"notify_admin_user_password_changes":                       types.BoolValue(false),
		"account_changes_notification_publisher_ref":               types.ObjectNull(resourcelink.AttrType()),
		"metadata_notification_settings":                           types.ObjectNull(notificationSettingsAttrType),
		"expired_certificate_administrative_console_warning_days":  types.Int64Unknown(),
		"expiring_certificate_administrative_console_warning_days": types.Int64Unknown(),
		"thread_pool_exhaustion_notification_settings":             types.ObjectNull(threadPoolExhaustionNotificationSettingsAttrType),
	})

	oauthRoleDefault, _ = types.ObjectValue(oauthRoleAttrType, map[string]attr.Value{
		"enable_oauth":           types.BoolValue(true),
		"enable_open_id_connect": types.BoolValue(true),
	})
	idpSamlProfileDefault, _ = types.ObjectValue(idpSaml20ProfileAttrType, map[string]attr.Value{
		"enable":              types.BoolValue(true),
		"enable_auto_connect": types.BoolNull(),
	})
	spSamlProfileDefault, _ = types.ObjectValue(spSaml20ProfileAttrType, map[string]attr.Value{
		"enable":              types.BoolValue(true),
		"enable_auto_connect": types.BoolNull(),
		"enable_xasp":         types.BoolValue(true),
	})
	idpRoleDefault, _ = types.ObjectValue(idpRoleAttrType, map[string]attr.Value{
		"enable":                       types.BoolValue(true),
		"enable_saml_1_1":              types.BoolValue(true),
		"enable_saml_1_0":              types.BoolValue(true),
		"enable_ws_fed":                types.BoolValue(true),
		"enable_ws_trust":              types.BoolValue(true),
		"saml_2_0_profile":             idpSamlProfileDefault,
		"enable_outbound_provisioning": types.BoolValue(true),
	})
	spRoleDefault, _ = types.ObjectValue(spRoleAttrType, map[string]attr.Value{
		"enable":                      types.BoolValue(true),
		"enable_saml_1_1":             types.BoolValue(true),
		"enable_saml_1_0":             types.BoolValue(true),
		"enable_ws_fed":               types.BoolValue(true),
		"enable_ws_trust":             types.BoolValue(true),
		"saml_2_0_profile":            spSamlProfileDefault,
		"enable_open_id_connect":      types.BoolValue(true),
		"enable_inbound_provisioning": types.BoolValue(true),
	})
	rolesAndProtocolsDefault, _ = types.ObjectValue(rolesAndProtocolsAttrType, map[string]attr.Value{
		"oauth_role":           oauthRoleDefault,
		"idp_role":             idpRoleDefault,
		"sp_role":              spRoleDefault,
		"enable_idp_discovery": types.BoolValue(true),
	})
)

type serverSettingsModel struct {
	Id                types.String `tfsdk:"id"`
	ContactInfo       types.Object `tfsdk:"contact_info"`
	Notifications     types.Object `tfsdk:"notifications"`
	RolesAndProtocols types.Object `tfsdk:"roles_and_protocols"`
	FederationInfo    types.Object `tfsdk:"federation_info"`
	EmailServer       types.Object `tfsdk:"email_server"`
	CaptchaSettings   types.Object `tfsdk:"captcha_settings"`
}

func readServerSettingsResponse(ctx context.Context, r *client.ServerSettings, state *serverSettingsModel, plan *serverSettingsModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	emptyString := ""
	state.Id = id.GenerateUUIDToState(existingId)
	state.ContactInfo, respDiags = types.ObjectValueFrom(ctx, contactInfoAttrType, r.ContactInfo)
	diags.Append(respDiags...)
	state.Notifications, respDiags = types.ObjectValueFrom(ctx, notificationsAttrType, r.Notifications)
	diags.Append(respDiags...)
	//////////////////////////////////////////////
	// ROLES AND PROTOCOLS
	//////////////////////////////////////////////
	idpSaml20ProfileVal, respDiags := types.ObjectValueFrom(ctx, idpSaml20ProfileAttrType, r.RolesAndProtocols.IdpRole.Saml20Profile)
	diags.Append(respDiags...)
	idpRoleAttrValue := map[string]attr.Value{
		"enable":                       types.BoolPointerValue(r.RolesAndProtocols.IdpRole.Enable),
		"enable_saml_1_1":              types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableSaml11),
		"enable_saml_1_0":              types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableSaml10),
		"enable_ws_fed":                types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableWsFed),
		"enable_ws_trust":              types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableWsTrust),
		"saml_2_0_profile":             idpSaml20ProfileVal,
		"enable_outbound_provisioning": types.BoolPointerValue(r.RolesAndProtocols.IdpRole.EnableOutboundProvisioning),
	}
	idpRoleVal, respDiags := types.ObjectValue(idpRoleAttrType, idpRoleAttrValue)
	diags.Append(respDiags...)

	spSaml20ProfileVal, respDiags := types.ObjectValueFrom(ctx, spSaml20ProfileAttrType, r.RolesAndProtocols.SpRole.Saml20Profile)
	diags.Append(respDiags...)

	spRoleAttrValue := map[string]attr.Value{
		"enable":                      types.BoolPointerValue(r.RolesAndProtocols.SpRole.Enable),
		"enable_saml_1_1":             types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableSaml11),
		"enable_saml_1_0":             types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableSaml10),
		"enable_ws_fed":               types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableWsFed),
		"enable_ws_trust":             types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableWsTrust),
		"saml_2_0_profile":            spSaml20ProfileVal,
		"enable_open_id_connect":      types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableOpenIDConnect),
		"enable_inbound_provisioning": types.BoolPointerValue(r.RolesAndProtocols.SpRole.EnableInboundProvisioning),
	}
	// save SP role to state
	spRoleVal, respDiags := types.ObjectValue(spRoleAttrType, spRoleAttrValue)
	diags.Append(respDiags...)
	oauthRoleVal, respDiags := types.ObjectValueFrom(ctx, oauthRoleAttrType, r.RolesAndProtocols.OauthRole)
	diags.Append(respDiags...)
	rolesAndProtocolsAttrTypeValues := map[string]attr.Value{
		"oauth_role":           oauthRoleVal,
		"idp_role":             idpRoleVal,
		"sp_role":              spRoleVal,
		"enable_idp_discovery": types.BoolPointerValue(r.RolesAndProtocols.EnableIdpDiscovery),
	}
	state.RolesAndProtocols, respDiags = types.ObjectValue(rolesAndProtocolsAttrType, rolesAndProtocolsAttrTypeValues)
	diags.Append(respDiags...)
	//////////////////////////////////////////////
	// FEDERATION INFO
	//////////////////////////////////////////////
	federationInfoAttrValue := map[string]attr.Value{
		"base_url":               types.StringPointerValue(r.FederationInfo.BaseUrl),
		"saml_2_entity_id":       types.StringPointerValue(r.FederationInfo.Saml2EntityId),
		"auto_connect_entity_id": types.StringPointerValue(&emptyString),
		"saml_1x_issuer_id":      types.StringPointerValue(r.FederationInfo.Saml1xIssuerId),
		"saml_1x_source_id":      types.StringPointerValue(r.FederationInfo.Saml1xSourceId),
		"wsfed_realm":            types.StringPointerValue(r.FederationInfo.WsfedRealm),
	}

	state.FederationInfo, respDiags = types.ObjectValue(federationInfoAttrType, federationInfoAttrValue)
	diags.Append(respDiags...)
	//////////////////////////////////////////////
	// EMAIL SERVER
	//////////////////////////////////////////////
	// get email creds with function
	// if username and password are not set, return null values
	if internaltypes.IsDefined(plan.EmailServer) {
		var getEmailCreds = func() (*string, string) {
			if plan.EmailServer.Attributes()["username"] != nil && plan.EmailServer.Attributes()["password"] != nil {
				username := plan.EmailServer.Attributes()["username"].(types.String).ValueStringPointer()
				password := plan.EmailServer.Attributes()["password"].(types.String).ValueString()
				return username, password
			} else {
				return types.StringNull().ValueStringPointer(), types.StringNull().ValueString()
			}
		}

		// retrieve values for saving to state
		username, password := getEmailCreds()
		emailServerAttrValue := map[string]attr.Value{
			"source_addr":                 types.StringValue(r.EmailServer.SourceAddr),
			"email_server":                types.StringValue(r.EmailServer.EmailServer),
			"port":                        types.Int64Value(r.EmailServer.Port),
			"ssl_port":                    types.Int64PointerValue(r.EmailServer.SslPort),
			"timeout":                     types.Int64PointerValue(r.EmailServer.Timeout),
			"retry_attempts":              types.Int64PointerValue(r.EmailServer.RetryAttempts),
			"retry_delay":                 types.Int64PointerValue(r.EmailServer.RetryDelay),
			"use_ssl":                     types.BoolPointerValue(r.EmailServer.UseSSL),
			"use_tls":                     types.BoolPointerValue(r.EmailServer.UseTLS),
			"verify_hostname":             types.BoolPointerValue(r.EmailServer.VerifyHostname),
			"enable_utf8_message_headers": types.BoolPointerValue(r.EmailServer.EnableUtf8MessageHeaders),
			"use_debugging":               types.BoolPointerValue(r.EmailServer.UseDebugging),
			"username":                    types.StringPointerValue(username),
			"password":                    types.StringValue(password),
		}

		state.EmailServer, respDiags = types.ObjectValue(emailServerAttrType, emailServerAttrValue)
		diags.Append(respDiags...)
	} else {
		state.EmailServer = types.ObjectNull(emailServerAttrType)
	}
	//////////////////////////////////////////////
	// CAPTCHA SETTINGS
	//////////////////////////////////////////////
	if internaltypes.IsDefined(plan.CaptchaSettings) {
		captchaSettingsAttrValue := map[string]attr.Value{
			"site_key":   types.StringPointerValue(r.CaptchaSettings.SiteKey),
			"secret_key": types.StringValue(plan.CaptchaSettings.Attributes()["secret_key"].(types.String).ValueString()),
		}
		state.CaptchaSettings, respDiags = types.ObjectValue(captchaSettingsAttrType, captchaSettingsAttrValue)
		diags.Append(respDiags...)
	} else {
		state.CaptchaSettings = types.ObjectNull(captchaSettingsAttrType)
	}
	return diags
}
