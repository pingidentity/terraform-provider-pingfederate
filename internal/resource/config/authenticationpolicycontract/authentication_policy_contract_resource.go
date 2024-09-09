package authenticationpolicycontract

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPolicyContractResource{}
	_ resource.ResourceWithConfigure   = &authenticationPolicyContractResource{}
	_ resource.ResourceWithImportState = &authenticationPolicyContractResource{}

	coreAttributesDefaultObjAttrType = map[string]attr.Type{
		"name": types.StringType,
	}

	coreAttributesDefaultObjAttrValue = map[string]attr.Value{
		"name": types.StringValue("subject"),
	}

	coreAttributesDefaultObjValue, _ = types.ObjectValue(coreAttributesDefaultObjAttrType, coreAttributesDefaultObjAttrValue)
	coreAttributesDefaultListAttrVal = []attr.Value{coreAttributesDefaultObjValue}
	coreAttributesDefaultListVal, _  = types.ListValue(attributeElemAttrType, coreAttributesDefaultListAttrVal)
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

// GetSchema defines the schema for the resource.
func (r *authenticationPolicyContractResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	extendedAttributesDefault, _ := types.SetValue(attributeElemAttrType, nil)
	schema := schema.Schema{
		Description: "Manages an authentication policy contract.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the authentication policy contract. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.PingFederateId(),
				},
			},
			"core_attributes": schema.ListNestedAttribute{
				Description: "A list of read-only assertion attributes (for example, subject) that are automatically populated by PingFederate.",
				Computed:    true,
				Optional:    false,
				Default:     listdefault.StaticValue(coreAttributesDefaultListVal),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of this attribute.",
						},
					},
				},
			},
			"extended_attributes": schema.SetNestedAttribute{
				Description: "A list of additional attributes as needed.",
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(extendedAttributesDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of this attribute.",
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Description: "The Authentication Policy contract name. Name is unique.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addAuthenticationPolicyContractsFields(ctx context.Context, addRequest *client.AuthenticationPolicyContract, plan authenticationPolicyContractModel) error {
	addRequest.Id = plan.ContractId.ValueStringPointer()

	addRequest.CoreAttributes = []client.AuthenticationPolicyContractAttribute{}
	for _, coreAttribute := range plan.CoreAttributes.Elements() {
		unmarshalled := client.AuthenticationPolicyContractAttribute{}
		err := json.Unmarshal([]byte(internaljson.FromValue(coreAttribute, false)), &unmarshalled)
		if err != nil {
			return err
		}
		addRequest.CoreAttributes = append(addRequest.CoreAttributes, unmarshalled)
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

	addRequest.Name = plan.Name.ValueStringPointer()
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

func (r *authenticationPolicyContractResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPolicyContractModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAuthenticationPolicyContracts := client.NewAuthenticationPolicyContract()
	err := addAuthenticationPolicyContractsFields(ctx, createAuthenticationPolicyContracts, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for an authentication policy contract: "+err.Error())
		return
	}

	apiCreateAuthenticationPolicyContracts := r.apiClient.AuthenticationPolicyContractsAPI.CreateAuthenticationPolicyContract(config.AuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationPolicyContracts = apiCreateAuthenticationPolicyContracts.Body(*createAuthenticationPolicyContracts)
	authenticationPolicyContractsResponse, httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.CreateAuthenticationPolicyContractExecute(apiCreateAuthenticationPolicyContracts)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating an authentication policy contract", err, httpResp)
		return
	}

	// Read the response into the state
	var state authenticationPolicyContractModel

	diags = readAuthenticationPolicyContractsResponse(ctx, authenticationPolicyContractsResponse, &state, &plan)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPolicyContractResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPolicyContractModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPolicyContracts, httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.GetAuthenticationPolicyContract(config.AuthContext(ctx, r.providerConfig), state.ContractId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Authentication Policy Contract", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an authentication policy contract", err, httpResp)
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
	var plan authenticationPolicyContractModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state authenticationPolicyContractModel
	updateAuthenticationPolicyContracts := r.apiClient.AuthenticationPolicyContractsAPI.UpdateAuthenticationPolicyContract(config.AuthContext(ctx, r.providerConfig), plan.ContractId.ValueString())
	createUpdateRequest := client.NewAuthenticationPolicyContract()
	err := addAuthenticationPolicyContractsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for an authentication policy contract: "+err.Error())
		return
	}

	updateAuthenticationPolicyContracts = updateAuthenticationPolicyContracts.Body(*createUpdateRequest)
	updateAuthenticationPolicyContractsResponse, httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.UpdateAuthenticationPolicyContractExecute(updateAuthenticationPolicyContracts)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating an authentication policy contract", err, httpResp)
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
	var state authenticationPolicyContractModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationPolicyContractsAPI.DeleteAuthenticationPolicyContract(config.AuthContext(ctx, r.providerConfig), state.ContractId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an authentication policy contract", err, httpResp)
	}

}

func (r *authenticationPolicyContractResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("contract_id"), req, resp)
}
