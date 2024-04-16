package authenticationpoliciespolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/authenticationpolicytreenode"
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

type authenticationPoliciesPolicyModel struct {
	Id                              types.String `tfsdk:"id"`
	PolicyId                        types.String `tfsdk:"policy_id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	AuthenticationApiApplicationRef types.Object `tfsdk:"authentication_api_application_ref"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
	RootNode                        types.Object `tfsdk:"root_node"`
	HandleFailuresLocally           types.Bool   `tfsdk:"handle_failures_locally"`
}

func (r *authenticationPoliciesPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an Authentication Policy",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The authentication policy name. Name is unique.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the authentication policy.",
			},
			"authentication_api_application_ref": resourcelink.CompleteSingleNestedAttribute(true, false, false,
				"Authentication API Application Id to be used in this policy branch. If the value is not specified, no Authentication API Application will be used."),
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether or not this authentication policy tree is enabled. Default is true.",
			},
			"root_node": authenticationpolicytreenode.ToSchema(),
			"handle_failures_locally": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If a policy ends in failure keep the user local.",
			},
		},
	}
	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "policy_id", false, false, "The authentication policy ID. ID is unique.")
	resp.Schema = schema
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

func readAuthenticationPolicyResponse(ctx context.Context, r *client.AuthenticationPolicyTree, state *authenticationPoliciesPolicyModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

	state.Id = types.StringPointerValue(r.Id)
	state.PolicyId = types.StringPointerValue(r.Id)
	state.Name = types.StringPointerValue(r.Name)
	state.Description = types.StringPointerValue(r.Description)
	state.Enabled = types.BoolPointerValue(r.Enabled)
	state.HandleFailuresLocally = types.BoolPointerValue(r.HandleFailuresLocally)

	state.AuthenticationApiApplicationRef, respDiags = resourcelink.ToState(ctx, r.AuthenticationApiApplicationRef)
	diags.Append(respDiags...)

	state.RootNode, respDiags = authenticationpolicytreenode.ToState(ctx, r.RootNode, 1)
	diags.Append(respDiags...)

	return diags
}

func addOptionalAuthenticationPolicyFields(addRequest *client.AuthenticationPolicyTree, plan authenticationPoliciesPolicyModel) error {
	addRequest.Id = plan.PolicyId.ValueStringPointer()
	addRequest.Name = plan.Name.ValueStringPointer()
	addRequest.Description = plan.Description.ValueStringPointer()
	addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	addRequest.HandleFailuresLocally = plan.HandleFailuresLocally.ValueBoolPointer()

	var err error
	addRequest.AuthenticationApiApplicationRef, err = resourcelink.ClientStruct(plan.AuthenticationApiApplicationRef)
	if err != nil {
		return err
	}

	if internaltypes.IsDefined(plan.RootNode) {
		addRequest.RootNode, err = authenticationpolicytreenode.ClientStruct(plan.RootNode)
		if err != nil {
			return err
		}
	}

	return err
}

func (r *authenticationPoliciesPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPoliciesPolicyModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newPolicyTree := client.NewAuthenticationPolicyTree()
	err := addOptionalAuthenticationPolicyFields(newPolicyTree, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the Authentication Policy", err.Error())
		return
	}

	apiCreatePolicy := r.apiClient.AuthenticationPoliciesAPI.CreatePolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreatePolicy = apiCreatePolicy.Body(*newPolicyTree)
	policyResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.CreatePolicyExecute(apiCreatePolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Authentication Policy", err, httpResp)
		return
	}

	// Read the response into the state
	var state authenticationPoliciesPolicyModel
	diags = readAuthenticationPolicyResponse(ctx, policyResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPoliciesPolicyModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.GetPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Policy", err, httpResp)
		}
		return
	}

	var updatedState authenticationPoliciesPolicyModel
	diags = readAuthenticationPolicyResponse(ctx, policyResponse, &updatedState)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPoliciesPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan authenticationPoliciesPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePolicyRequest := r.apiClient.AuthenticationPoliciesAPI.UpdatePolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.PolicyId.ValueString())
	updatedPolicy := client.NewAuthenticationPolicyTree()
	err := addOptionalAuthenticationPolicyFields(updatedPolicy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the Authentication Policy", err.Error())
		return
	}

	updatePolicyRequest = updatePolicyRequest.Body(*updatedPolicy)
	updateResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.UpdatePolicyExecute(updatePolicyRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Authentication Policy", err, httpResp)
		return
	}

	// Read the response
	var state authenticationPoliciesPolicyModel
	readResponseDiags := readAuthenticationPolicyResponse(ctx, updateResponse, &state)
	resp.Diagnostics.Append(readResponseDiags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authenticationPoliciesPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationPoliciesAPI.DeletePolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.PolicyId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an Authentication Policy", err, httpResp)
	}
}

func (r *authenticationPoliciesPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
