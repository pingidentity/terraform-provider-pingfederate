package authenticationpoliciespolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/authenticationpolicytreenode"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &authenticationPoliciesPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationPoliciesPolicyDataSource{}
)

// AuthenticationPoliciesPolicyDataSource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesPolicyDataSource() datasource.DataSource {
	return &authenticationPoliciesPolicyDataSource{}
}

// authenticationPoliciesPolicyDataSource is the datasource implementation.
type authenticationPoliciesPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *authenticationPoliciesPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes an Authentication Policy",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The authentication policy name. Name is unique.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "A description for the authentication policy.",
			},
			"authentication_api_application_ref": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToDataSourceSchema(),
				Computed:    true,
				Optional:    false,
				Description: "Authentication API Application Id to be used in this policy branch. If the value is not specified, no Authentication API Application will be used.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Whether or not this authentication policy tree is enabled. Default is true.",
			},
			"root_node": authenticationpolicytreenode.DataSourceSchema(),
			"handle_failures_locally": schema.BoolAttribute{
				Computed:    true,
				Optional:    false,
				Description: "If a policy ends in failure keep the user local.",
			},
		},
	}

	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"policy_id", true, "The authentication policy ID. ID is unique.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *authenticationPoliciesPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies_policy"
}

func (r *authenticationPoliciesPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *authenticationPoliciesPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationPoliciesPolicyModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.GetPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.PolicyId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Policy", err, httpResp)
		return
	}

	var updatedState authenticationPoliciesPolicyModel
	diags = readAuthenticationPolicyResponse(ctx, policyResponse, &updatedState)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}
