// Copyright © 2026 Ping Identity Corporation

package localidentity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	authSourcesAttrTypes = map[string]attr.Type{
		"id":     types.StringType,
		"source": types.StringType,
	}

	authSourceUpdatePolicyAttrTypes = map[string]attr.Type{
		"store_attributes":  types.BoolType,
		"retain_attributes": types.BoolType,
		"update_attributes": types.BoolType,
		"update_interval":   types.Int64Type,
	}

	registrationConfigAttrTypes = map[string]attr.Type{
		"captcha_enabled":                         types.BoolType,
		"captcha_provider_ref":                    types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"template_name":                           types.StringType,
		"create_authn_session_after_registration": types.BoolType,
		"username_field":                          types.StringType,
		"this_is_my_device_enabled":               types.BoolType,
		"registration_workflow":                   types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"execute_workflow":                        types.StringType,
	}

	profileConfigAttrTypes = map[string]attr.Type{
		"delete_identity_enabled": types.BoolType,
		"template_name":           types.StringType,
	}

	fieldItemAttrTypes = map[string]attr.Type{
		"type":                    types.StringType,
		"id":                      types.StringType,
		"label":                   types.StringType,
		"registration_page_field": types.BoolType,
		"profile_page_field":      types.BoolType,
		"attributes":              types.MapType{ElemType: types.BoolType},
		"options":                 types.SetType{ElemType: types.StringType},
		"default_value":           types.StringType,
	}

	fieldConfigAttrTypes = map[string]attr.Type{
		"fields":                        types.SetType{ElemType: types.ObjectType{AttrTypes: fieldItemAttrTypes}},
		"strip_space_from_unique_field": types.BoolType,
	}

	emailVerificationConfigAttrTypes = map[string]attr.Type{
		"email_verification_enabled":               types.BoolType,
		"verify_email_template_name":               types.StringType,
		"email_verification_sent_template_name":    types.StringType,
		"email_verification_success_template_name": types.StringType,
		"email_verification_error_template_name":   types.StringType,
		"email_verification_type":                  types.StringType,
		"otp_length":                               types.Int64Type,
		"otp_retry_attempts":                       types.Int64Type,
		"allowed_otp_character_set":                types.StringType,
		"otp_time_to_live":                         types.Int64Type,
		"email_verification_otp_template_name":     types.StringType,
		"otl_time_to_live":                         types.Int64Type,
		"field_for_email_to_verify":                types.StringType,
		"field_storing_verification_status":        types.StringType,
		"notification_publisher_ref":               types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"require_verified_email":                   types.BoolType,
		"require_verified_email_template_name":     types.StringType,
	}

	dsConfigAttrTypes = map[string]attr.Type{
		"base_dn":                  types.StringType,
		"type":                     types.StringType,
		"data_store_ref":           types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"data_store_mapping":       types.MapType{ElemType: types.ObjectType{AttrTypes: dsMappingAttrtypes}},
		"create_pattern":           types.StringType,
		"object_class":             types.StringType,
		"auxiliary_object_classes": types.SetType{ElemType: types.StringType},
	}

	dsMappingAttrtypes = map[string]attr.Type{
		"type":     types.StringType,
		"name":     types.StringType,
		"metadata": types.MapType{ElemType: types.StringType},
	}
)

type localIdentityProfileModel struct {
	Id                      types.String `tfsdk:"id"`
	ProfileId               types.String `tfsdk:"profile_id"`
	Name                    types.String `tfsdk:"name"`
	ApcId                   types.Object `tfsdk:"apc_id"`
	AuthSources             types.Set    `tfsdk:"auth_sources"`
	AuthSourceUpdatePolicy  types.Object `tfsdk:"auth_source_update_policy"`
	RegistrationEnabled     types.Bool   `tfsdk:"registration_enabled"`
	RegistrationConfig      types.Object `tfsdk:"registration_config"`
	ProfileConfig           types.Object `tfsdk:"profile_config"`
	FieldConfig             types.Object `tfsdk:"field_config"`
	EmailVerificationConfig types.Object `tfsdk:"email_verification_config"`
	DataStoreConfig         types.Object `tfsdk:"data_store_config"`
	ProfileEnabled          types.Bool   `tfsdk:"profile_enabled"`
}

// authSourceUpdatePolicyState explicitly constructs the auth_source_update_policy object
func authSourceUpdatePolicyState(s *client.LocalIdentityAuthSourceUpdatePolicy) (types.Object, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(authSourceUpdatePolicyAttrTypes), diag.Diagnostics{}
	}
	return types.ObjectValue(authSourceUpdatePolicyAttrTypes, map[string]attr.Value{
		"store_attributes":  types.BoolPointerValue(s.StoreAttributes),
		"retain_attributes": types.BoolPointerValue(s.RetainAttributes),
		"update_attributes": types.BoolPointerValue(s.UpdateAttributes),
		"update_interval":   types.Int64Value(int64(s.GetUpdateInterval())),
	})
}

// registrationConfigState explicitly constructs the registration_config object to avoid the
// location field on the client ResourceLink, which is not in the Terraform schema
func registrationConfigState(s *client.RegistrationConfig) (types.Object, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(registrationConfigAttrTypes), diag.Diagnostics{}
	}
	return types.ObjectValue(registrationConfigAttrTypes, map[string]attr.Value{
		"captcha_enabled":                         types.BoolPointerValue(s.CaptchaEnabled),
		"captcha_provider_ref":                    resourcelink.ToStateMust(s.CaptchaProviderRef),
		"template_name":                           types.StringValue(s.TemplateName),
		"create_authn_session_after_registration": types.BoolPointerValue(s.CreateAuthnSessionAfterRegistration),
		"username_field":                          types.StringPointerValue(s.UsernameField),
		"this_is_my_device_enabled":               types.BoolPointerValue(s.ThisIsMyDeviceEnabled),
		"registration_workflow":                   resourcelink.ToStateMust(s.RegistrationWorkflow),
		"execute_workflow":                        types.StringPointerValue(s.ExecuteWorkflow),
	})
}

// profileConfigState explicitly constructs the profile_config object
func profileConfigState(s *client.ProfileConfig) (types.Object, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(profileConfigAttrTypes), diag.Diagnostics{}
	}
	return types.ObjectValue(profileConfigAttrTypes, map[string]attr.Value{
		"delete_identity_enabled": types.BoolPointerValue(s.DeleteIdentityEnabled),
		"template_name":           types.StringValue(s.TemplateName),
	})
}

// fieldConfigState explicitly constructs the field_config object to avoid the location
// field on the client ResourceLink, which is not in the Terraform schema
func fieldConfigState(s *client.FieldConfig) (types.Object, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(fieldConfigAttrTypes), diag.Diagnostics{}
	}
	fieldsVal, diags := types.SetValueFrom(context.Background(), types.ObjectType{AttrTypes: fieldItemAttrTypes}, s.Fields)
	objVal, objDiags := types.ObjectValue(fieldConfigAttrTypes, map[string]attr.Value{
		"fields":                        fieldsVal,
		"strip_space_from_unique_field": types.BoolPointerValue(s.StripSpaceFromUniqueField),
	})
	diags.Append(objDiags...)
	return objVal, diags
}

// emailVerificationConfigState explicitly constructs the email_verification_config object to
// avoid the location field on the client ResourceLink, which is not in the Terraform schema
func emailVerificationConfigState(s *client.EmailVerificationConfig) (types.Object, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(emailVerificationConfigAttrTypes), diag.Diagnostics{}
	}
	return types.ObjectValue(emailVerificationConfigAttrTypes, map[string]attr.Value{
		"email_verification_enabled":               types.BoolPointerValue(s.EmailVerificationEnabled),
		"verify_email_template_name":               types.StringPointerValue(s.VerifyEmailTemplateName),
		"email_verification_sent_template_name":    types.StringPointerValue(s.EmailVerificationSentTemplateName),
		"email_verification_success_template_name": types.StringPointerValue(s.EmailVerificationSuccessTemplateName),
		"email_verification_error_template_name":   types.StringPointerValue(s.EmailVerificationErrorTemplateName),
		"email_verification_type":                  types.StringPointerValue(s.EmailVerificationType),
		"otp_length":                               types.Int64PointerValue(s.OtpLength),
		"otp_retry_attempts":                       types.Int64PointerValue(s.OtpRetryAttempts),
		"allowed_otp_character_set":                types.StringPointerValue(s.AllowedOtpCharacterSet),
		"otp_time_to_live":                         types.Int64PointerValue(s.OtpTimeToLive),
		"email_verification_otp_template_name":     types.StringPointerValue(s.EmailVerificationOtpTemplateName),
		"otl_time_to_live":                         types.Int64PointerValue(s.OtlTimeToLive),
		"field_for_email_to_verify":                types.StringValue(s.FieldForEmailToVerify),
		"field_storing_verification_status":        types.StringValue(s.FieldStoringVerificationStatus),
		"notification_publisher_ref":               resourcelink.ToStateMust(s.NotificationPublisherRef),
		"require_verified_email":                   types.BoolPointerValue(s.RequireVerifiedEmail),
		"require_verified_email_template_name":     types.StringPointerValue(s.RequireVerifiedEmailTemplateName),
	})
}

// dsConfigState explicitly constructs the data_store_config object to avoid the location
// field on the client ResourceLink, which is not in the Terraform schema
func dsConfigState(s *client.LdapDataStoreConfig) (types.Object, diag.Diagnostics) {
	if s == nil {
		return types.ObjectNull(dsConfigAttrTypes), diag.Diagnostics{}
	}
	var diags diag.Diagnostics
	dataStoreMappingVal, diags := types.MapValueFrom(context.Background(), types.MapType{ElemType: types.ObjectType{AttrTypes: dsMappingAttrtypes}}, s.DataStoreMapping)
	objVal, objDiags := types.ObjectValue(dsConfigAttrTypes, map[string]attr.Value{
		"base_dn":                  types.StringValue(s.BaseDn),
		"type":                     types.StringValue(s.Type),
		"data_store_ref":           resourcelink.ToStateMust(&s.DataStoreRef),
		"data_store_mapping":       dataStoreMappingVal,
		"create_pattern":           types.StringValue(s.CreatePattern),
		"object_class":             types.StringValue(s.ObjectClass),
		"auxiliary_object_classes": internaltypes.GetStringSet(s.AuxiliaryObjectClasses),
	})
	diags.Append(objDiags...)
	return objVal, diags
}
