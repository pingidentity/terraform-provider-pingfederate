package idpspconnection

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
)

type idpSpConnectionModel struct {
	SpBrowserSso                           types.Object `tfsdk:"sp_browser_sso"`
	Type                                   types.String `tfsdk:"type"`
	ConnectionId                           types.String `tfsdk:"connection_id"`
	Id                                     types.String `tfsdk:"id"`
	EntityId                               types.String `tfsdk:"entity_id"`
	Name                                   types.String `tfsdk:"name"`
	CreationDate                           types.String `tfsdk:"creation_date"`
	Active                                 types.Bool   `tfsdk:"active"`
	BaseUrl                                types.String `tfsdk:"base_url"`
	DefaultVirtualEntityId                 types.String `tfsdk:"default_virtual_entity_id"`
	VirtualEntityIds                       types.Set    `tfsdk:"virtual_entity_ids"`
	MetadataReloadSettings                 types.Object `tfsdk:"metadata_reload_settings"`
	Credentials                            types.Object `tfsdk:"credentials"`
	ContactInfo                            types.Object `tfsdk:"contact_info"`
	LicenseConnectionGroup                 types.String `tfsdk:"license_connection_group"`
	LoggingMode                            types.String `tfsdk:"logging_mode"`
	AdditionalAllowedEntitiesConfiguration types.Object `tfsdk:"additional_allowed_entities_configuration"`
	ExtendedProperties                     types.Map    `tfsdk:"extended_properties"`
	AttributeQuery                         types.Object `tfsdk:"attribute_query"`
	WsTrust                                types.Object `tfsdk:"ws_trust"`
	ApplicationName                        types.String `tfsdk:"application_name"`
	ApplicationIconUrl                     types.String `tfsdk:"application_icon_url"`
	OutboundProvision                      types.Object `tfsdk:"outbound_provision"`
	ConnectionTargetType                   types.String `tfsdk:"connection_target_type"`
}

var (
	attributeQueryAttrTypes = map[string]attr.Type{
		"attributes":                     types.ListType{ElemType: types.StringType},
		"attribute_contract_fulfillment": attributeContractFulfillmentAttrType,
		"issuance_criteria":              issuanceCriteriaAttrType,
		"policy":                         types.ObjectType{AttrTypes: policyAttrTypes},
		"attribute_sources":              types.ListType{ElemType: types.ObjectType{AttrTypes: attributesources.ElemAttrType()}},
	}

	customSchemaAttrTypes = map[string]attr.Type{
		"namespace": types.StringType,
		"attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":           types.StringType,
			"multi_valued":   types.BoolType,
			"types":          types.ListType{ElemType: types.StringType},
			"sub_attributes": types.ListType{ElemType: types.StringType},
		}}},
	}

	saasFieldInfoAttrTypes = map[string]attr.Type{
		"attribute_names": types.ListType{ElemType: types.StringType},
		"default_value":   types.StringType,
		"expression":      types.StringType,
		"create_only":     types.BoolType,
		"trim":            types.BoolType,
		"character_case":  types.StringType,
		"parser":          types.StringType,
		"masked":          types.BoolType,
	}
	attributeMappingElemAttrTypes = types.ObjectType{AttrTypes: map[string]attr.Type{
		"field_name":      types.StringType,
		"saas_field_info": types.ObjectType{AttrTypes: saasFieldInfoAttrTypes},
	}}
	channelSourceAttrTypes = map[string]attr.Type{
		"data_source":         resourceLinkObjectType,
		"guid_attribute_name": types.StringType,
		"guid_binary":         types.BoolType,
		"change_detection_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
			"user_object_class":         types.StringType,
			"group_object_class":        types.StringType,
			"changed_users_algorithm":   types.StringType,
			"usn_attribute_name":        types.StringType,
			"time_stamp_attribute_name": types.StringType,
		}},
		"group_membership_detection": types.ObjectType{AttrTypes: map[string]attr.Type{
			"member_of_group_attribute_name": types.StringType,
			"group_member_attribute_name":    types.StringType,
		}},
		"account_management_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
			"account_status_attribute_name": types.StringType,
			"account_status_algorithm":      types.StringType,
			"flag_comparison_value":         types.StringType,
			"flag_comparison_status":        types.BoolType,
			"default_status":                types.BoolType,
		}},
		"base_dn":               types.StringType,
		"user_source_location":  channelSourceLocationAttrType,
		"group_source_location": channelSourceLocationAttrType,
	}
	certsDefault, _ = types.ListValue(certsListType.ElemType, nil)
)

func readIdpSpconnectionResponseCommon(ctx context.Context, r *client.SpConnection, state *idpSpConnectionModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

	state.ConnectionId = types.StringPointerValue(r.Id)
	state.Id = types.StringPointerValue(r.Id)
	state.Type = types.StringPointerValue(r.Type)
	state.EntityId = types.StringValue(r.EntityId)
	state.Name = types.StringValue(r.Name)
	state.Active = types.BoolPointerValue(r.Active)
	state.BaseUrl = types.StringPointerValue(r.BaseUrl)
	state.DefaultVirtualEntityId = types.StringPointerValue(r.DefaultVirtualEntityId)
	state.LicenseConnectionGroup = types.StringPointerValue(r.LicenseConnectionGroup)
	state.LoggingMode = types.StringPointerValue(r.LoggingMode)
	state.ApplicationName = types.StringPointerValue(r.ApplicationName)
	state.ApplicationIconUrl = types.StringPointerValue(r.ApplicationIconUrl)
	state.ConnectionTargetType = types.StringPointerValue(r.ConnectionTargetType)

	if r.CreationDate != nil {
		state.CreationDate = types.StringValue(r.CreationDate.Format(time.RFC3339))
	} else {
		state.CreationDate = types.StringNull()
	}

	state.VirtualEntityIds, respDiags = types.SetValueFrom(ctx, types.StringType, r.VirtualEntityIds)
	diags.Append(respDiags...)

	state.MetadataReloadSettings, respDiags = types.ObjectValueFrom(ctx, metadataReloadSettingsAttrTypes, r.MetadataReloadSettings)
	diags.Append(respDiags...)

	state.Credentials, respDiags = types.ObjectValueFrom(ctx, credentialsAttrTypes, r.Credentials)
	diags.Append(respDiags...)
	if r.Credentials != nil && r.Credentials.SigningSettings != nil && r.Credentials.SigningSettings.IncludeCertInSignature == nil {
		// PF returns false for include_cert_in_signature as nil. If nil is returned, just set it to false
		credentialsAttrs := state.Credentials.Attributes()
		signingSettingsAttrs := credentialsAttrs["signing_settings"].(types.Object).Attributes()
		signingSettingsAttrs["include_cert_in_signature"] = types.BoolValue(false)
		newSigningSettings, respDiags := types.ObjectValue(signingSettingsAttrTypes, signingSettingsAttrs)
		diags.Append(respDiags...)
		credentialsAttrs["signing_settings"] = newSigningSettings
		state.Credentials, respDiags = types.ObjectValue(credentialsAttrTypes, credentialsAttrs)
		diags.Append(respDiags...)
	}

	state.ContactInfo, respDiags = types.ObjectValueFrom(ctx, contactInfoAttrTypes, r.ContactInfo)
	diags.Append(respDiags...)

	state.AdditionalAllowedEntitiesConfiguration, respDiags = types.ObjectValueFrom(ctx, additionalAllowedEntitiesConfigurationAttrTypes, r.AdditionalAllowedEntitiesConfiguration)
	diags.Append(respDiags...)

	state.ExtendedProperties, respDiags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: extendedPropertiesElemAttrTypes}, r.ExtendedProperties)
	diags.Append(respDiags...)

	state.SpBrowserSso, respDiags = types.ObjectValueFrom(ctx, spBrowserSSOAttrTypes, r.SpBrowserSso)
	diags.Append(respDiags...)

	state.WsTrust, respDiags = types.ObjectValueFrom(ctx, wsTrustAttrTypes, r.WsTrust)
	diags.Append(respDiags...)

	// AttributeQuery and OutboundProvision logic differs depending on plan, so those are left out here and read in the individual resource and data source files

	return diags
}
