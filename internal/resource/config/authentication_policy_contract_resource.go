package config

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPolicyContractsResource{}
	_ resource.ResourceWithConfigure   = &authenticationPolicyContractsResource{}
	_ resource.ResourceWithImportState = &authenticationPolicyContractsResource{}
)

// AuthenticationPolicyContractsResource is a helper function to simplify the provider implementation.
func AuthenticationPolicyContractsResource() resource.Resource {
	return &authenticationPolicyContractsResource{}
}

// authenticationPolicyContractsResource is the resource implementation.
type authenticationPolicyContractsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type authenticationPolicyContractsResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	CoreAttributes     types.Set    `tfsdk:"core_attributes"`
	ExtendedAttributes types.Set    `tfsdk:"extended_attributes"`
}

// GetSchema defines the schema for the resource.
func (r *authenticationPolicyContractsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a AuthenticationPolicyContracts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The persistent, unique ID for the authentication policy contract. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"core_attributes": schema.SetNestedAttribute{
				Description: "A list of read-only assertion attributes (for example, subject) that are automatically populated by PingFederate.",
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
}

func addAuthenticationPolicyContractsFields(ctx context.Context, addRequest *client.AuthenticationPolicyContract, plan authenticationPolicyContractsResourceModel) error {
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
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
func (r *authenticationPolicyContractsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policy_contract"
}

func (r *authenticationPolicyContractsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthenticationPolicyContractsResponse(ctx context.Context, r *client.AuthenticationPolicyContract, state *authenticationPolicyContractsResourceModel, expectedValues *authenticationPolicyContractsResourceModel) {
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
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
	caSliceOfObj, _ := types.SetValue(cAobjSlice, caSlice)

	clientExtAttributes := r.GetExtendedAttributes()
	var eaSlice = []attr.Value{}
	eAobjSlice := types.ObjectType{AttrTypes: attrType}
	for i := 0; i < len(clientExtAttributes); i++ {
		eAname := clientExtAttributes[i].GetName()
		eAnameVal := map[string]attr.Value{"name": types.StringValue(eAname)}
		newEaObj, _ := types.ObjectValue(attrType, eAnameVal)
		eaSlice = append(eaSlice, newEaObj)
	}
	eaSliceOfObj, _ := types.SetValue(eAobjSlice, eaSlice)

	state.CoreAttributes = basetypes.SetValue{}
	state.CoreAttributes = caSliceOfObj
	state.ExtendedAttributes = basetypes.SetValue{}
	state.ExtendedAttributes = eaSliceOfObj
}

func (r *authenticationPolicyContractsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPolicyContractsResourceModel

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
	_, requestErr := createAuthenticationPolicyContracts.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of an Authentication Policy Contract: %s", requestErr.Error())
	}

	apiCreateAuthenticationPolicyContracts := r.apiClient.AuthenticationPolicyContractsApi.CreateAuthenticationPolicyContract(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationPolicyContracts = apiCreateAuthenticationPolicyContracts.Body(*createAuthenticationPolicyContracts)
	authenticationPolicyContractsResponse, httpResp, err := r.apiClient.AuthenticationPolicyContractsApi.CreateAuthenticationPolicyContractExecute(apiCreateAuthenticationPolicyContracts)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating an Authentication Policy Contract", err, httpResp)
		return
	}
	_, responseErr := authenticationPolicyContractsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an Authentication Policy Contract: %s", responseErr.Error())
	}
	// Read the response into the state
	var state authenticationPolicyContractsResourceModel

	readAuthenticationPolicyContractsResponse(ctx, authenticationPolicyContractsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPolicyContractsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPolicyContractsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPolicyContracts, httpResp, err := r.apiClient.AuthenticationPolicyContractsApi.GetAuthenticationPolicyContract(ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Authentication Policy Contract", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Authentication Policy Contract", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadAuthenticationPolicyContracts.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an Authentication Policy Contract: %s", responseErr.Error())
	}

	// Read the response into the state
	readAuthenticationPolicyContractsResponse(ctx, apiReadAuthenticationPolicyContracts, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPolicyContractsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan authenticationPolicyContractsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state authenticationPolicyContractsResourceModel
	req.State.Get(ctx, &state)
	updateAuthenticationPolicyContracts := r.apiClient.AuthenticationPolicyContractsApi.UpdateAuthenticationPolicyContract(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewAuthenticationPolicyContract()
	err := addAuthenticationPolicyContractsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for an Authentication Policy Contract", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of an Authentication Policy Contract: %s", requestErr.Error())
	}
	updateAuthenticationPolicyContracts = updateAuthenticationPolicyContracts.Body(*createUpdateRequest)
	updateAuthenticationPolicyContractsResponse, httpResp, err := r.apiClient.AuthenticationPolicyContractsApi.UpdateAuthenticationPolicyContractExecute(updateAuthenticationPolicyContracts)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating an Authentication Policy Contract", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateAuthenticationPolicyContractsResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an Authentication Policy Contract: %s", responseErr.Error())
	}
	// Read the response
	readAuthenticationPolicyContractsResponse(ctx, updateAuthenticationPolicyContractsResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *authenticationPolicyContractsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authenticationPolicyContractsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationPolicyContractsApi.DeleteAuthenticationPolicyContract(ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an Authentication Policy Contract", err, httpResp)
		return
	}

}

func (r *authenticationPolicyContractsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
