package authenticationpolicycontract

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPolicyContractResource{}
	_ resource.ResourceWithConfigure   = &authenticationPolicyContractResource{}
	_ resource.ResourceWithImportState = &authenticationPolicyContractResource{}
)

// AuthenticationPolicyContractResource is a helper function to simplify the provider implementation.
func AuthenticationPolicyContractResource() resource.Resource {
	return &authenticationPolicyContractResource{}
}

// authenticationPolicyContractResource is the resource implementation.
type authenticationPolicyContractResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type authenticationPolicyContractResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	AuthenticationPolicyContractId types.String `tfsdk:"authentication_policy_contract_id"`
	Name                           types.String `tfsdk:"name"`
	CoreAttributes                 types.List   `tfsdk:"core_attributes"`
	ExtendedAttributes             types.Set    `tfsdk:"extended_attributes"`
}

// GetSchema defines the schema for the resource.
func (r *authenticationPolicyContractResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an Authentication Policy Contract.",
		Attributes: map[string]schema.Attribute{
			"core_attributes": schema.ListNestedAttribute{
				Description: "A list of read-only assertion attributes (for example, subject) that are automatically populated by PingFederate.",
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"extended_attributes": schema.SetNestedAttribute{
				Description: "A list of additional attributes as needed.",
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Description: "The Authentication Policy Contract Name. Name is unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"authentication_policy_contract_id",
		false,
		"The persistent, unique ID for the authentication policy contract. It can be any combination of [a-zA-Z0-9._-].")
	resp.Schema = schema
}

func addAuthenticationPolicyContractsFields(ctx context.Context, addRequest *client.AuthenticationPolicyContract, plan authenticationPolicyContractResourceModel) error {
	if internaltypes.IsDefined(plan.AuthenticationPolicyContractId) {
		addRequest.Id = plan.AuthenticationPolicyContractId.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CoreAttributes) {
		addRequest.CoreAttributes = []client.AuthenticationPolicyContractAttribute{}
		for _, coreAttribute := range plan.CoreAttributes.Elements() {
			unmarshalled := client.AuthenticationPolicyContractAttribute{}
			err := json.Unmarshal([]byte(internaljson.FromValue(coreAttribute, false)), &unmarshalled)
			if err != nil {
				return err
			}
			addRequest.CoreAttributes = append(addRequest.CoreAttributes, unmarshalled)
		}
	}
	if internaltypes.IsDefined(plan.ExtendedAttributes) {
		addRequest.ExtendedAttributes = []client.AuthenticationPolicyContractAttribute{}
		for _, extendedAttribute := range plan.ExtendedAttributes.Elements() {
			unmarshalled := client.AuthenticationPolicyContractAttribute{}
			err := json.Unmarshal([]byte(internaljson.FromValue(extendedAttribute, false)), &unmarshalled)
			if err != nil {
				return err
			}
			addRequest.ExtendedAttributes = append(addRequest.ExtendedAttributes, unmarshalled)
		}
	}
	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueStringPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *authenticationPolicyContractResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policy_contract"
}

func (r *authenticationPolicyContractResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthenticationPolicyContractsResponse(ctx context.Context, r *client.AuthenticationPolicyContract, state *authenticationPolicyContractResourceModel, expectedValues *authenticationPolicyContractResourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.AuthenticationPolicyContractId = internaltypes.StringTypeOrNil(r.Id, false)
	state.Name = internaltypes.StringTypeOrNil(r.Name, false)

	var attrType = map[string]attr.Type{"name": types.StringType}
	clientCoreAttributes := r.GetCoreAttributes()
	var caSlice = []attr.Value{}
	cAobjSlice := types.ObjectType{AttrTypes: attrType}
	for i := 0; i < len(clientCoreAttributes); i++ {
		cAname := clientCoreAttributes[i].GetName()
		cAnameVal := map[string]attr.Value{"name": types.StringValue(cAname)}
		newCaObj, _ := types.ObjectValue(attrType, cAnameVal)
		caSlice = append(caSlice, newCaObj)
	}
	caSliceOfObj, respDiags := types.ListValue(cAobjSlice, caSlice)
	diags.Append(respDiags...)

	clientExtAttributes := r.GetExtendedAttributes()
	var eaSlice = []attr.Value{}
	eAobjSlice := types.ObjectType{AttrTypes: attrType}
	for i := 0; i < len(clientExtAttributes); i++ {
		eAname := clientExtAttributes[i].GetName()
		eAnameVal := map[string]attr.Value{"name": types.StringValue(eAname)}
		newEaObj, _ := types.ObjectValue(attrType, eAnameVal)
		eaSlice = append(eaSlice, newEaObj)
	}
	eaSliceOfObj, respDiags := types.SetValue(eAobjSlice, eaSlice)
	diags.Append(respDiags...)

	state.CoreAttributes = basetypes.ListValue{}
	state.CoreAttributes = caSliceOfObj
	state.ExtendedAttributes = basetypes.SetValue{}
	state.ExtendedAttributes = eaSliceOfObj
	return diags
}

func (r *authenticationPolicyContractResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPolicyContractResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAuthenticationPolicyContracts := client.NewAuthenticationPolicyContract()
	err := addAuthenticationPolicyContractsFields(ctx, createAuthenticationPolicyContracts, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for an Authentication Policy Contract", err.Error())
		return
	}

	apiCreateAuthenticationPolicyContracts := r.apiClient.AuthenticationPolicyContractsAPI.CreateAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationPolicyContracts = apiCreateAuthenticationPolicyContracts.Body(*createAuthenticationPolicyContracts)
	authenticationPolicyContractsResponse, httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.CreateAuthenticationPolicyContractExecute(apiCreateAuthenticationPolicyContracts)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating an Authentication Policy Contract", err, httpResp)
		return
	}

	// Read the response into the state
	var state authenticationPolicyContractResourceModel

	diags = readAuthenticationPolicyContractsResponse(ctx, authenticationPolicyContractsResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPolicyContractResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPolicyContractResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPolicyContracts, httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.GetAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AuthenticationPolicyContractId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Authentication Policy Contract", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Authentication Policy Contract", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readAuthenticationPolicyContractsResponse(ctx, apiReadAuthenticationPolicyContracts, &state, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPolicyContractResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan authenticationPolicyContractResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state authenticationPolicyContractResourceModel
	updateAuthenticationPolicyContracts := r.apiClient.AuthenticationPolicyContractsAPI.UpdateAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.AuthenticationPolicyContractId.ValueString())
	createUpdateRequest := client.NewAuthenticationPolicyContract()
	err := addAuthenticationPolicyContractsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for an Authentication Policy Contract", err.Error())
		return
	}

	updateAuthenticationPolicyContracts = updateAuthenticationPolicyContracts.Body(*createUpdateRequest)
	updateAuthenticationPolicyContractsResponse, httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.UpdateAuthenticationPolicyContractExecute(updateAuthenticationPolicyContracts)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating an Authentication Policy Contract", err, httpResp)
		return
	}

	// Read the response
	diags = readAuthenticationPolicyContractsResponse(ctx, updateAuthenticationPolicyContractsResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *authenticationPolicyContractResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authenticationPolicyContractResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.DeleteAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AuthenticationPolicyContractId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an Authentication Policy Contract", err, httpResp)
		return
	}

}

func (r *authenticationPolicyContractResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("authentication_policy_contract_id"), req, resp)
}
