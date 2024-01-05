package authenticationpoliciesfragments

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/policyaction"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPoliciesFragmentResource{}
	_ resource.ResourceWithConfigure   = &authenticationPoliciesFragmentResource{}
	_ resource.ResourceWithImportState = &authenticationPoliciesFragmentResource{}
)

// AuthenticationPoliciesFragmentResource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesFragmentResource() resource.Resource {
	return &authenticationPoliciesFragmentResource{}
}

// authenticationPoliciesFragmentResource is the resource implementation.
type authenticationPoliciesFragmentResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type authenticationPoliciesFragmentModel struct {
	Description types.String `tfsdk:"description"`
	FragmentId  types.String `tfsdk:"fragment_id"`
	Id          types.String `tfsdk:"id"`
	Inputs      types.Object `tfsdk:"inputs"`
	Name        types.String `tfsdk:"name"`
	Outputs     types.Object `tfsdk:"outputs"`
	RootNode    types.Object `tfsdk:"root_node"`
}

func (r *authenticationPoliciesFragmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the authentication policy fragment.",
			},
			"inputs": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToSchema(),
				Optional:    true,
				Description: "A reference to a resource.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The authentication policy fragment name. Name is unique.",
			},
			"outputs": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToSchema(),
				Optional:    true,
				Description: "A reference to a resource.",
			},
			"root_node": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"policy_action": policyaction.Schema(),
					//TODO children, recursiveness
				},
				Optional:    true,
				Description: "An authentication policy tree node.",
			},
		},
	}
	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema, "fragment_id", false, "The authentication policy fragment ID. ID is unique.")
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *authenticationPoliciesFragmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies_fragment"
}

func (r *authenticationPoliciesFragmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthenticationPoliciesFragmentResponse(ctx context.Context, r *client.Client, state *authenticationPoliciesFragmentModel) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func addOptionalAuthenticationPoliciesFragmentFields(ctx context.Context, addRequest *client.Client, plan authenticationPoliciesFragmentModel) error {

	return nil
}

func (r *authenticationPoliciesFragmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationPoliciesFragmentModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the response into the state
	var state authenticationPoliciesFragmentModel
	//diags = readAuthenticationPoliciesFragmentResponse(ctx, oauthClientResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesFragmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPoliciesFragmentModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPoliciesFragmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan authenticationPoliciesFragmentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	//diags = resp.State.Set(ctx, &state)
	//resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesFragmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authenticationPoliciesFragmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationPoliciesAPI.DeleteFragment(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.FragmentId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an Authentication Policy Fragment", err, httpResp)
	}
}

func (r *authenticationPoliciesFragmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to fragment_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("fragment_id"), req, resp)
}
