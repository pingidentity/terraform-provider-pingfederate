// Copyright © 2026 Ping Identity Corporation

package localidentity

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
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
