package idpspconnection

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	ReplicationStatus                      types.String `tfsdk:"replication_status"`
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
