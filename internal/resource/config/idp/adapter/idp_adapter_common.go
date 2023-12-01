package idpadapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	datasourcepluginconfiguration "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Define attribute types for object types
var (
	// May move some of this into common package if future resources need this
	attributesAttrType = map[string]attr.Type{
		"name":      types.StringType,
		"pseudonym": types.BoolType,
		"masked":    types.BoolType,
	}

	attributeContractAttrTypes = map[string]attr.Type{
		"core_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"core_attributes_all": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"extended_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"unique_user_key_attribute": types.StringType,
		"mask_ognl_values":          types.BoolType,
		"inherited":                 types.BoolType,
	}

	attributeContractDataSourceAttrTypes = map[string]attr.Type{
		"core_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"extended_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: attributesAttrType,
			},
		},
		"unique_user_key_attribute": types.StringType,
		"mask_ognl_values":          types.BoolType,
		"inherited":                 types.BoolType,
	}

	attributeMappingAttrTypes = map[string]attr.Type{
		"attribute_sources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: attributesources.ElemAttrType(),
			},
		},
		"attribute_contract_fulfillment": attributecontractfulfillment.MapType(),
		"issuance_criteria": types.ObjectType{
			AttrTypes: issuancecriteria.AttrType(),
		},
		"inherited": types.BoolType,
	}

	extendedAttributesDefault, _ = types.SetValue(types.ObjectType{
		AttrTypes: attributesAttrType,
	}, nil)
)

type idpAdapterModel struct {
	AuthnCtxClassRef    types.String `tfsdk:"authn_ctx_class_ref"`
	Id                  types.String `tfsdk:"id"`
	AdapterId           types.String `tfsdk:"adapter_id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
	AttributeMapping    types.Object `tfsdk:"attribute_mapping"`
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
}

func readIdpAdapterResponse(ctx context.Context, r *client.IdpAdapter, state *idpAdapterModel, plan *idpAdapterModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.AuthnCtxClassRef = internaltypes.StringTypeOrNil(r.AuthnCtxClassRef, false)
	state.AdapterId = types.StringValue(r.Id)
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef, diags = resourcelink.ToState(ctx, &r.PluginDescriptorRef)
	respDiags.Append(diags...)
	state.ParentRef, diags = resourcelink.ToState(ctx, r.ParentRef)
	respDiags.Append(diags...)
	// Configuration
	if plan != nil {
		state.Configuration, diags = pluginconfiguration.ToState(plan.Configuration, &r.Configuration)
		respDiags.Append(diags...)
	} else {
		state.Configuration, diags = datasourcepluginconfiguration.ToDataSourceState(ctx, &r.Configuration)
		respDiags.Append(diags...)
	}

	if r.AttributeContract != nil {
		attributeContractValues := map[string]attr.Value{}
		attributeContractValues["extended_attributes"], diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attributesAttrType}, r.AttributeContract.ExtendedAttributes)
		respDiags.Append(diags...)
		if plan != nil {
			attributeContractValues["core_attributes_all"], diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attributesAttrType}, r.AttributeContract.CoreAttributes)
			respDiags.Append(diags...)
		}
		attributeContractValues["unique_user_key_attribute"] = types.StringPointerValue(r.AttributeContract.UniqueUserKeyAttribute)
		attributeContractValues["mask_ognl_values"] = types.BoolPointerValue(r.AttributeContract.MaskOgnlValues)

		// PF returns false as nil for inherited in some cases
		inherited := false
		if r.AttributeContract.Inherited != nil {
			inherited = *r.AttributeContract.Inherited
		}

		attributeContractValues["inherited"] = types.BoolValue(inherited)

		// Only include core_attributes specified in the plan in the response
		if plan != nil {
			if internaltypes.IsDefined(plan.AttributeContract) && internaltypes.IsDefined(plan.AttributeContract.Attributes()["core_attributes"]) {
				coreAttributes := []attr.Value{}
				planCoreAttributeNames := map[string]bool{}
				for _, planCoreAttr := range plan.AttributeContract.Attributes()["core_attributes"].(types.Set).Elements() {
					planCoreAttributeNames[planCoreAttr.(types.Object).Attributes()["name"].(types.String).ValueString()] = true
				}
				for _, coreAttr := range r.AttributeContract.CoreAttributes {
					_, attrInPlan := planCoreAttributeNames[coreAttr.Name]
					if attrInPlan {
						attrObjVal, diags := types.ObjectValueFrom(ctx, attributesAttrType, coreAttr)
						respDiags.Append(diags...)
						coreAttributes = append(coreAttributes, attrObjVal)
					}
				}
				attributeContractValues["core_attributes"], diags = types.SetValue(types.ObjectType{AttrTypes: attributesAttrType}, coreAttributes)
				respDiags.Append(diags...)
			} else {
				attributeContractValues["core_attributes"] = types.SetNull(types.ObjectType{AttrTypes: attributesAttrType})
			}
			state.AttributeContract, diags = types.ObjectValue(attributeContractAttrTypes, attributeContractValues)
			respDiags.Append(diags...)
		} else {
			state.AttributeContract, diags = types.ObjectValueFrom(ctx, attributeContractDataSourceAttrTypes, r.AttributeContract)
			respDiags.Append(diags...)
		}
	}

	if r.AttributeMapping != nil {
		attributeMappingValues := map[string]attr.Value{
			"inherited": types.BoolPointerValue(r.AttributeMapping.Inherited),
		}
		// The PF API won't return inherited if it is false
		if r.AttributeMapping.Inherited == nil {
			attributeMappingValues["inherited"] = types.BoolValue(false)
		}

		// Build attribute_contract_fulfillment value
		attributeMappingValues["attribute_contract_fulfillment"], diags = attributecontractfulfillment.ToState(ctx, r.AttributeMapping.AttributeContractFulfillment)
		respDiags.Append(diags...)

		// Build issuance_criteria value
		attributeMappingValues["issuance_criteria"], diags = issuancecriteria.ToState(ctx, r.AttributeMapping.IssuanceCriteria)
		respDiags.Append(diags...)

		// Build attribute_sources value
		attributeMappingValues["attribute_sources"], respDiags = attributesources.ToState(ctx, r.AttributeMapping.AttributeSources)
		diags.Append(respDiags...)

		// Build complete attribute mapping value
		state.AttributeMapping, diags = types.ObjectValue(attributeMappingAttrTypes, attributeMappingValues)
		respDiags.Append(diags...)
	}
	return respDiags
}
