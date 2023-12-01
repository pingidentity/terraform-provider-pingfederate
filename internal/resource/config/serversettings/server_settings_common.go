package serversettings

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var (
	contactInfoAttrType = map[string]attr.Type{
		"company":    basetypes.StringType{},
		"email":      basetypes.StringType{},
		"first_name": basetypes.StringType{},
		"last_name":  basetypes.StringType{},
		"phone":      basetypes.StringType{},
	}

	certificateExpirationsAttrType = map[string]attr.Type{
		"email_address":              basetypes.StringType{},
		"initial_warning_period":     basetypes.Int64Type{},
		"final_warning_period":       basetypes.Int64Type{},
		"notification_publisher_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
	}

	notificationSettingsAttrType = map[string]attr.Type{
		"email_address":              basetypes.StringType{},
		"notification_publisher_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
	}

	notificationsAttrType = map[string]attr.Type{
		"license_events":                             basetypes.ObjectType{AttrTypes: notificationSettingsAttrType},
		"certificate_expirations":                    basetypes.ObjectType{AttrTypes: certificateExpirationsAttrType},
		"notify_admin_user_password_changes":         basetypes.BoolType{},
		"account_changes_notification_publisher_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"metadata_notification_settings":             basetypes.ObjectType{AttrTypes: notificationSettingsAttrType},
	}

	oauthRoleAttrType = map[string]attr.Type{
		"enable_oauth":           basetypes.BoolType{},
		"enable_open_id_connect": basetypes.BoolType{},
	}

	idpSaml20ProfileAttrType = map[string]attr.Type{
		"enable":              basetypes.BoolType{},
		"enable_auto_connect": basetypes.BoolType{},
	}

	spSaml20ProfileAttrType = map[string]attr.Type{
		"enable":              basetypes.BoolType{},
		"enable_auto_connect": basetypes.BoolType{},
		"enable_xasp":         basetypes.BoolType{},
	}

	idpRoleAttrType = map[string]attr.Type{
		"enable":                       basetypes.BoolType{},
		"enable_saml_1_1":              basetypes.BoolType{},
		"enable_saml_1_0":              basetypes.BoolType{},
		"enable_ws_fed":                basetypes.BoolType{},
		"enable_ws_trust":              basetypes.BoolType{},
		"saml_2_0_profile":             basetypes.ObjectType{AttrTypes: idpSaml20ProfileAttrType},
		"enable_outbound_provisioning": basetypes.BoolType{},
	}

	spRoleAttrType = map[string]attr.Type{
		"enable":                      basetypes.BoolType{},
		"enable_saml_1_1":             basetypes.BoolType{},
		"enable_saml_1_0":             basetypes.BoolType{},
		"enable_ws_fed":               basetypes.BoolType{},
		"enable_ws_trust":             basetypes.BoolType{},
		"saml_2_0_profile":            basetypes.ObjectType{AttrTypes: spSaml20ProfileAttrType},
		"enable_open_id_connect":      basetypes.BoolType{},
		"enable_inbound_provisioning": basetypes.BoolType{},
	}

	rolesAndProtocolsAttrType = map[string]attr.Type{
		"oauth_role":           basetypes.ObjectType{AttrTypes: oauthRoleAttrType},
		"idp_role":             basetypes.ObjectType{AttrTypes: idpRoleAttrType},
		"sp_role":              basetypes.ObjectType{AttrTypes: spRoleAttrType},
		"enable_idp_discovery": basetypes.BoolType{},
	}

	federationInfoAttrType = map[string]attr.Type{
		"base_url":               basetypes.StringType{},
		"saml_2_entity_id":       basetypes.StringType{},
		"auto_connect_entity_id": basetypes.StringType{},
		"saml_1x_issuer_id":      basetypes.StringType{},
		"saml_1x_source_id":      basetypes.StringType{},
		"wsfed_realm":            basetypes.StringType{},
	}

	emailServerAttrType = map[string]attr.Type{
		"source_addr":                 basetypes.StringType{},
		"email_server":                basetypes.StringType{},
		"port":                        basetypes.Int64Type{},
		"ssl_port":                    basetypes.Int64Type{},
		"timeout":                     basetypes.Int64Type{},
		"retry_attempts":              basetypes.Int64Type{},
		"retry_delay":                 basetypes.Int64Type{},
		"use_ssl":                     basetypes.BoolType{},
		"use_tls":                     basetypes.BoolType{},
		"verify_hostname":             basetypes.BoolType{},
		"enable_utf8_message_headers": basetypes.BoolType{},
		"use_debugging":               basetypes.BoolType{},
		"username":                    basetypes.StringType{},
		"password":                    basetypes.StringType{},
	}

	captchaSettingsAttrType = map[string]attr.Type{
		"site_key":   basetypes.StringType{},
		"secret_key": basetypes.StringType{},
	}

	contactInfoDefault, _ = types.ObjectValue(contactInfoAttrType, map[string]attr.Value{
		"company":    types.StringNull(),
		"email":      types.StringNull(),
		"first_name": types.StringNull(),
		"last_name":  types.StringNull(),
		"phone":      types.StringNull(),
	})

	notificationsDefault, _ = types.ObjectValue(notificationsAttrType, map[string]attr.Value{
		"license_events":                             types.ObjectNull(notificationSettingsAttrType),
		"certificate_expirations":                    types.ObjectNull(certificateExpirationsAttrType),
		"notify_admin_user_password_changes":         types.BoolValue(false),
		"account_changes_notification_publisher_ref": types.ObjectNull(resourcelink.AttrType()),
		"metadata_notification_settings":             types.ObjectNull(notificationSettingsAttrType),
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
