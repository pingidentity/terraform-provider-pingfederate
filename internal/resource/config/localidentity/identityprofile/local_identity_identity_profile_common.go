package localidentity

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var (
	authSourcesAttrTypes = map[string]attr.Type{
		"id":     basetypes.StringType{},
		"source": basetypes.StringType{},
	}

	authSourceUpdatePolicyAttrTypes = map[string]attr.Type{
		"store_attributes":  basetypes.BoolType{},
		"retain_attributes": basetypes.BoolType{},
		"update_attributes": basetypes.BoolType{},
		"update_interval":   basetypes.Int64Type{},
	}

	registrationConfigAttrTypes = map[string]attr.Type{
		"captcha_enabled":                         basetypes.BoolType{},
		"captcha_provider_ref":                    basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"template_name":                           basetypes.StringType{},
		"create_authn_session_after_registration": basetypes.BoolType{},
		"username_field":                          basetypes.StringType{},
		"this_is_my_device_enabled":               basetypes.BoolType{},
		"registration_workflow":                   basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"execute_workflow":                        basetypes.StringType{},
	}

	profileConfigAttrTypes = map[string]attr.Type{
		"delete_identity_enabled": basetypes.BoolType{},
		"template_name":           basetypes.StringType{},
	}

	fieldItemAttrTypes = map[string]attr.Type{
		"type":                    basetypes.StringType{},
		"id":                      basetypes.StringType{},
		"label":                   basetypes.StringType{},
		"registration_page_field": basetypes.BoolType{},
		"profile_page_field":      basetypes.BoolType{},
		"attributes":              basetypes.MapType{ElemType: basetypes.BoolType{}},
	}

	fieldConfigAttrTypes = map[string]attr.Type{
		"fields":                        basetypes.ListType{ElemType: types.ObjectType{AttrTypes: fieldItemAttrTypes}},
		"strip_space_from_unique_field": basetypes.BoolType{},
	}

	emailVerificationConfigAttrTypes = map[string]attr.Type{
		"email_verification_enabled":               basetypes.BoolType{},
		"verify_email_template_name":               basetypes.StringType{},
		"email_verification_sent_template_name":    basetypes.StringType{},
		"email_verification_success_template_name": basetypes.StringType{},
		"email_verification_error_template_name":   basetypes.StringType{},
		"email_verification_type":                  basetypes.StringType{},
		"otp_length":                               basetypes.Int64Type{},
		"otp_retry_attempts":                       basetypes.Int64Type{},
		"allowed_otp_character_set":                basetypes.StringType{},
		"otp_time_to_live":                         basetypes.Int64Type{},
		"email_verification_otp_template_name":     basetypes.StringType{},
		"otl_time_to_live":                         basetypes.Int64Type{},
		"field_for_email_to_verify":                basetypes.StringType{},
		"field_storing_verification_status":        basetypes.StringType{},
		"notification_publisher_ref":               basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"require_verified_email":                   basetypes.BoolType{},
		"require_verified_email_template_name":     basetypes.StringType{},
	}

	dsConfigAttrTypes = map[string]attr.Type{
		"base_dn":                  basetypes.StringType{},
		"type":                     basetypes.StringType{},
		"data_store_ref":           basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"data_store_mapping":       basetypes.MapType{ElemType: types.ObjectType{AttrTypes: dsMappingAttrtypes}},
		"create_pattern":           basetypes.StringType{},
		"object_class":             basetypes.StringType{},
		"auxiliary_object_classes": basetypes.SetType{ElemType: basetypes.StringType{}},
	}

	dsMappingAttrtypes = map[string]attr.Type{
		"type":     basetypes.StringType{},
		"name":     basetypes.StringType{},
		"metadata": basetypes.MapType{ElemType: basetypes.StringType{}},
	}
)

type localIdentityIdentityProfileModel struct {
	Id                      types.String `tfsdk:"id"`
	ProfileId               types.String `tfsdk:"profile_id"`
	Name                    types.String `tfsdk:"name"`
	ApcId                   types.Object `tfsdk:"apc_id"`
	AuthSources             types.List   `tfsdk:"auth_sources"`
	AuthSourceUpdatePolicy  types.Object `tfsdk:"auth_source_update_policy"`
	RegistrationEnabled     types.Bool   `tfsdk:"registration_enabled"`
	RegistrationConfig      types.Object `tfsdk:"registration_config"`
	ProfileConfig           types.Object `tfsdk:"profile_config"`
	FieldConfig             types.Object `tfsdk:"field_config"`
	EmailVerificationConfig types.Object `tfsdk:"email_verification_config"`
	DataStoreConfig         types.Object `tfsdk:"data_store_config"`
	ProfileEnabled          types.Bool   `tfsdk:"profile_enabled"`
}
