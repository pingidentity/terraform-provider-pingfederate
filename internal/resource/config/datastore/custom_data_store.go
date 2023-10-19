package datastore

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var (
	customDataStoreAttrType = map[string]attr.Type{
		"type":                  basetypes.StringType{},
		"name":                  basetypes.StringType{},
		"plugin_descriptor_ref": basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"parent_ref":            basetypes.ObjectType{AttrTypes: resourcelink.AttrType()},
		"configuration":         basetypes.ObjectType{AttrTypes: pluginconfiguration.AttrType()},
	}
	customDataStoreEmptyStateObj = types.ObjectNull(customDataStoreAttrType)
)

func toSchemaCustomDataStore() schema.SingleNestedAttribute {
	customDataStoreSchema := schema.SingleNestedAttribute{}
	customDataStoreSchema.Description = "A custom data store."
	customDataStoreSchema.Default = objectdefault.StaticValue(types.ObjectNull(customDataStoreAttrType))
	customDataStoreSchema.Computed = true
	customDataStoreSchema.Optional = true
	customDataStoreSchema.Attributes = map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: "The data store type.",
			Computed:    true,
			Optional:    false,
			Default:     stringdefault.StaticString("CUSTOM"),
		},
		"name": schema.StringAttribute{
			Description: "The plugin instance name.",
			Required:    true,
		},
		"plugin_descriptor_ref": schema.SingleNestedAttribute{
			Required:    true,
			Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
			Attributes:  resourcelink.ToSchema(),
		},
		"parent_ref": schema.SingleNestedAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
			Default:     objectdefault.StaticValue(types.ObjectNull(resourcelink.AttrType())),
			Attributes:  resourcelink.ToSchema(),
		},
		"configuration": pluginconfiguration.ToSchema(),
	}
	customDataStoreSchema.Validators = []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("jdbc_data_store"),
			path.MatchRelative().AtParent().AtName("ldap_data_store"),
			path.MatchRelative().AtParent().AtName("ping_one_ldap_gateway_data_store"),
		),
	}

	return customDataStoreSchema
}
