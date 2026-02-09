// Copyright © 2026 Ping Identity Corporation

package passwordcredentialvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	pluginconfigurationdatasource "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var (
	attrType = map[string]attr.Type{
		"name": types.StringType,
	}

	attributeContractTypes = map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: types.ObjectType{AttrTypes: attrType}},
		"extended_attributes": types.SetType{ElemType: types.ObjectType{AttrTypes: attrType}},
	}

	emptyAttrSet, _ = types.SetValue(types.ObjectType{AttrTypes: attrType}, nil)
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

func readPasswordCredentialValidatorResponse(ctx context.Context, r *client.PasswordCredentialValidator, state *passwordCredentialValidatorModel, configurationFromPlan types.Object, isResource, isImportRead bool) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue(r.Id)
	state.ValidatorId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef, respDiags = resourcelink.ToState(ctx, &r.PluginDescriptorRef)
	diags.Append(respDiags...)
	state.ParentRef, respDiags = resourcelink.ToState(ctx, r.ParentRef)
	diags.Append(respDiags...)
	if isResource {
		state.Configuration, respDiags = pluginconfiguration.ToState(configurationFromPlan, &r.Configuration, isImportRead)
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
		attributeContractCoreAttributes, respDiags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrType}, coreAttrs)
		diags.Append(respDiags...)

		// state.AttributeContract extended_attributes
		attributeContractClientExtendedAttributes := attrContract.ExtendedAttributes
		extdAttrs := []client.PasswordCredentialValidatorAttribute{}
		for _, ea := range attributeContractClientExtendedAttributes {
			extendedAttr := client.PasswordCredentialValidatorAttribute{}
			extendedAttr.Name = ea.Name
			extdAttrs = append(extdAttrs, extendedAttr)
		}
		attributeContractExtendedAttributes, respDiags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrType}, extdAttrs)
		diags.Append(respDiags...)

		attributeContractValues := map[string]attr.Value{
			"core_attributes":     attributeContractCoreAttributes,
			"extended_attributes": attributeContractExtendedAttributes,
		}
		state.AttributeContract, respDiags = types.ObjectValue(attributeContractTypes, attributeContractValues)
		diags.Append(respDiags...)
	}

	return diags
}
