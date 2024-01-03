package authenticationpolicies

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPoliciesPolicyResource{}
	_ resource.ResourceWithConfigure   = &authenticationPoliciesPolicyResource{}
	_ resource.ResourceWithImportState = &authenticationPoliciesPolicyResource{}
)

// AuthenticationPoliciesPolicyResource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesPolicyResource() resource.Resource {
	return &authenticationPoliciesPolicyResource{}
}

// authenticationPoliciesPolicyResource is the resource implementation.
type authenticationPoliciesPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type authenticationPoliciesPolicyResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	AuthenticationApiApplicationRef types.Object `tfsdk:"authentication_api_application_ref"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
	RootNode                        types.Object `tfsdk:"root_node"`
	HandleFailuresLocally           types.Bool   `tfsdk:"handle_failures_locally"`
}

// GetSchema defines the schema for the resource.
func (r *authenticationPoliciesPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Authentication Policies Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The authentication policy ID. ID is unique.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The authentication policy name. Name is unique.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for the authentication policy.",
				Required:    true,
			},
			"authentication_api_application_ref": schema.SingleNestedAttribute{
				Description: "Authentication API Application ID to be used in this policy branch. If the value is not specified, no Authentication API Application will be used.",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether or not this authentication policy tree is enabled. Default is true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"root_node": schema.SingleNestedAttribute{
				Description: "A node inside the authentication policy tree.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of the node.",
						Required:    true,
					},
				},
			},
			"handle_failures_locally": schema.BoolAttribute{
				Description: "If a policy ends in failure keep the user local.",
				Optional:    true,
				Computed:    true,
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalAuthenticationPoliciesPolicyFields(ctx context.Context, addRequest *client.Policy, plan authenticationPoliciesPolicyResourceModel) error {

	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.AuthenticationApiApplicationRef) {
		addRequest.AuthenticationApiApplicationRef = client.NewAuthenticationApiApplicationRef()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AuthenticationApiApplicationRef, false)), addRequest.AuthenticationApiApplicationRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if internaltypes.IsDefined(plan.RootNode) {
		addRequest.RootNode = client.NewRootNode()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RootNode, false)), addRequest.RootNode)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.HandleFailuresLocally) {
		addRequest.HandleFailuresLocally = plan.HandleFailuresLocally.ValueBoolPointer()
	}

	return nil

}

// Metadata returns the resource type name.
func (r *authenticationPoliciesPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies_policy"
}

func (r *authenticationPoliciesPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthenticationPoliciesPolicyResponse(ctx context.Context, r *client.Policy, state *authenticationPoliciesPolicyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.Id = internaltypes.StringTypeOrNil(r.Id)
	state.Name = internaltypes.StringTypeOrNil(r.Name)
	state.Description = internaltypes.StringTypeOrNil(r.Description)
	state.AuthenticationApiApplicationRef = (r.AuthenticationApiApplicationRef)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RootNode = (r.RootNode)
	state.HandleFailuresLocally = types.BoolValue(r.HandleFailuresLocally)

	// make sure all object type building appends diags
	return diags
}

func (r *authenticationPoliciesPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPoliciesPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAuthenticationPoliciesPolicy := client.NewAuthenticationPolicyTree()
	err := addOptionalAuthenticationPoliciesPolicyFields(ctx, createAuthenticationPoliciesPolicy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthenticationPoliciesPolicy", err.Error())
		return
	}

	apiCreateAuthenticationPoliciesPolicy := r.apiClient.AuthenticationPoliciesApi.AddAuthenticationPolicyTree(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationPoliciesPolicy = apiCreateAuthenticationPoliciesPolicy.Body(*createAuthenticationPoliciesPolicy)
	authenticationPoliciesPolicyResponse, httpResp, err := r.apiClient.AuthenticationPoliciesApi.AddAuthenticationPolicyTreeExecute(apiCreateAuthenticationPoliciesPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the AuthenticationPoliciesPolicy", err, httpResp)
		return
	}

	// Read the response into the state
	var state authenticationPoliciesPolicyResourceModel

	diags = readAuthenticationPoliciesPolicyResponse(ctx, authenticationPoliciesPolicyResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPoliciesPolicyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPoliciesPolicy, httpResp, err := r.apiClient.AuthenticationPoliciesApi.GetAuthenticationPolicyTree(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.VALUE.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the AuthenticationPoliciesPolicy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  AuthenticationPoliciesPolicy", err, httpResp)
		}
	}

	// Read the response into the state
	readAuthenticationPoliciesPolicyResponse(ctx, apiReadAuthenticationPoliciesPolicy, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPoliciesPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan authenticationPoliciesPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateAuthenticationPoliciesPolicy := r.apiClient.AuthenticationPoliciesApi.UpdateAuthenticationPolicyTree(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.VALUE.ValueString())
	createUpdateRequest := client.NewAuthenticationPolicyTree()
	err := addOptionalAuthenticationPoliciesPolicyFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AuthenticationPoliciesPolicy", err.Error())
		return
	}

	updateAuthenticationPoliciesPolicy = updateAuthenticationPoliciesPolicy.Body(*createUpdateRequest)
	updateAuthenticationPoliciesPolicyResponse, httpResp, err := r.apiClient.AuthenticationPoliciesApi.UpdateAuthenticationPolicyTreeExecute(updateAuthenticationPoliciesPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating AuthenticationPoliciesPolicy", err, httpResp)
		return
	}

	// Read the response
	var state authenticationPoliciesPolicyResourceModel
	diags = readAuthenticationPoliciesPolicyResponse(ctx, updateAuthenticationPoliciesPolicyResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *authenticationPoliciesPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *authenticationPoliciesPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("custom_id"), req, resp)
}
