package authenticationPolicyContracts

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
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
	authenticationPolicyContractsResourceSchema(ctx, req, resp, false)
}

func authenticationPolicyContractsResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
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

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	resp.Schema = schema
}
func addAuthenticationPolicyContractsFields(ctx context.Context, addRequest *client.AuthenticationPolicyContract, plan authenticationPolicyContractsResourceModel) error {
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CoreAttributes) {
		addRequest.CoreAttributes = []client.AuthenticationPolicyContractAttribute{}
		pCaE := plan.CoreAttributes.Elements()
		for i := 0; i < len(pCaE); i++ {
			item := pCaE[i].(types.Object)
			getAttrName := item.Attributes()["name"].(types.String)
			attrName := client.NewAuthenticationPolicyContractAttribute(getAttrName.ValueString())
			addRequest.CoreAttributes = append(addRequest.CoreAttributes, *attrName)
		}
	}
	if internaltypes.IsDefined(plan.ExtendedAttributes) {
		addRequest.ExtendedAttributes = []client.AuthenticationPolicyContractAttribute{}
		pCaE := plan.ExtendedAttributes.Elements()
		for i := 0; i < len(pCaE); i++ {
			item := pCaE[i].(types.Object)
			getAttrName := item.Attributes()["name"].(types.String)
			attrName := client.NewAuthenticationPolicyContractAttribute(getAttrName.ValueString())
			addRequest.ExtendedAttributes = append(addRequest.ExtendedAttributes, *attrName)
		}
	}
	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueStringPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *authenticationPolicyContractsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policy_contracts"
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
	state.CoreAttributes = basetypes.SetValue{}
	caSliceOfObj, _ := internaltypes.ToStateAuthenticationPolicyContract(r)
	state.CoreAttributes = caSliceOfObj
	state.ExtendedAttributes = basetypes.SetValue{}
	_, eaSliceOfObj := internaltypes.ToStateAuthenticationPolicyContract(r)
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthenticationPolicyContracts", err.Error())
		return
	}
	requestJson, err := createAuthenticationPolicyContracts.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateAuthenticationPolicyContracts := r.apiClient.AuthenticationPolicyContractsApi.CreateAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationPolicyContracts = apiCreateAuthenticationPolicyContracts.Body(*createAuthenticationPolicyContracts)
	authenticationPolicyContractsResponse, httpResp, err := r.apiClient.AuthenticationPolicyContractsApi.CreateAuthenticationPolicyContractExecute(apiCreateAuthenticationPolicyContracts)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the AuthenticationPolicyContracts", err, httpResp)
		return
	}
	responseJson, err := authenticationPolicyContractsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state authenticationPolicyContractsResourceModel

	readAuthenticationPolicyContractsResponse(ctx, authenticationPolicyContractsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *authenticationPolicyContractsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAuthenticationPolicyContracts(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAuthenticationPolicyContracts(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state authenticationPolicyContractsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPolicyContracts, httpResp, err := apiClient.AuthenticationPolicyContractsApi.GetAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a AuthenticationPolicyContracts", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAuthenticationPolicyContracts.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAuthenticationPolicyContractsResponse(ctx, apiReadAuthenticationPolicyContracts, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPolicyContractsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAuthenticationPolicyContracts(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAuthenticationPolicyContracts(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
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
	updateAuthenticationPolicyContracts := apiClient.AuthenticationPolicyContractsApi.UpdateAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewAuthenticationPolicyContract()
	err := addAuthenticationPolicyContractsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthenticationPolicyContracts", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateAuthenticationPolicyContracts = updateAuthenticationPolicyContracts.Body(*createUpdateRequest)
	updateAuthenticationPolicyContractsResponse, httpResp, err := apiClient.AuthenticationPolicyContractsApi.UpdateAuthenticationPolicyContractExecute(updateAuthenticationPolicyContracts)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating AuthenticationPolicyContracts", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateAuthenticationPolicyContractsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readAuthenticationPolicyContractsResponse(ctx, updateAuthenticationPolicyContractsResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *authenticationPolicyContractsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteAuthenticationPolicyContracts(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteAuthenticationPolicyContracts(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state authenticationPolicyContractsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.AuthenticationPolicyContractsApi.DeleteAuthenticationPolicyContract(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a AuthenticationPolicyContracts", err, httpResp)
		return
	}

}

func (r *authenticationPolicyContractsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
