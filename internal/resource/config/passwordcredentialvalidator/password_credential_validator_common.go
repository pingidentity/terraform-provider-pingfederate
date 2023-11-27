package passwordcredentialvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	pluginconfigurationdatasource "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	pluginconfigurationresource "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
)

var (
	attrType = map[string]attr.Type{
		"name": basetypes.StringType{},
	}

	attributeContractTypes = map[string]attr.Type{
		"core_attributes":     basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"extended_attributes": basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"inherited":           basetypes.BoolType{},
	}

	emptyAttrList, _ = types.ListValue(types.ObjectType{AttrTypes: attrType}, nil)
)

type passwordCredentialValidatorModel struct {
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
	Id                  types.String `tfsdk:"id"`
	ValidatorId         types.String `tfsdk:"validator_id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
}

func readPasswordCredentialValidatorResponse(ctx context.Context, r *client.PasswordCredentialValidator, state *passwordCredentialValidatorModel, configurationFromPlan basetypes.ObjectValue, isResource bool) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue(r.Id)
	state.ValidatorId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef, respDiags = resourcelink.ToDataSourceState(ctx, &r.PluginDescriptorRef)
	diags.Append(respDiags...)
	state.ParentRef, respDiags = resourcelink.ToDataSourceState(ctx, r.ParentRef)
	diags.Append(respDiags...)
	if isResource {
		state.Configuration, respDiags = pluginconfigurationresource.ToState(configurationFromPlan, &r.Configuration)
		diags.Append(respDiags...)
	} else {
		state.Configuration, respDiags = pluginconfigurationdatasource.ToDataSourceState(ctx, &r.Configuration)
		diags.Append(respDiags...)
	}

	// state.AttributeContract
	if r.AttributeContract == nil {
		state.AttributeContract = types.ObjectNull(attributeContractTypes)
	} else {
		attrContract := r.AttributeContract
		// state.AttributeContract core_attributes
		attributeContractClientCoreAttributes := attrContract.CoreAttributes
		coreAttrs := []client.PasswordCredentialValidatorAttribute{}
		for _, ca := range attributeContractClientCoreAttributes {
			coreAttribute := client.PasswordCredentialValidatorAttribute{}
			coreAttribute.Name = ca.Name
			coreAttrs = append(coreAttrs, coreAttribute)
		}
		attributeContractCoreAttributes, respDiags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, coreAttrs)
		diags.Append(respDiags...)

		// state.AttributeContract extended_attributes
		attributeContractClientExtendedAttributes := attrContract.ExtendedAttributes
		extdAttrs := []client.PasswordCredentialValidatorAttribute{}
		for _, ea := range attributeContractClientExtendedAttributes {
			extendedAttr := client.PasswordCredentialValidatorAttribute{}
			extendedAttr.Name = ea.Name
			extdAttrs = append(extdAttrs, extendedAttr)
		}
		attributeContractExtendedAttributes, respDiags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, extdAttrs)
		diags.Append(respDiags...)

		// PF can return inherited as nil when it is false
		inherited := false
		if attrContract.Inherited != nil {
			inherited = *attrContract.Inherited
		}

		attributeContractValues := map[string]attr.Value{
			"core_attributes":     attributeContractCoreAttributes,
			"extended_attributes": attributeContractExtendedAttributes,
			"inherited":           types.BoolValue(inherited),
		}
		state.AttributeContract, respDiags = types.ObjectValue(attributeContractTypes, attributeContractValues)
		diags.Append(respDiags...)
	}

	return diags
}
