package authenticationpoliciesfragments

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
	_ datasource.DataSource              = &authenticationPoliciesFragmentDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationPoliciesFragmentDataSource{}
)

// AuthenticationPoliciesFragmentDataSource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesFragmentDataSource() datasource.DataSource {
	return &authenticationPoliciesFragmentDataSource{}
}

// authenticationPoliciesFragmentDataSource is the datasource implementation.
type authenticationPoliciesFragmentDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *authenticationPoliciesFragmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes an Authentication Policy Fragment",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "A description for the authentication policy fragment.",
			},
			"inputs": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToDataSourceSchema(),
				Computed:    true,
				Optional:    false,
				Description: "A reference to a resource.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The authentication policy fragment name. Name is unique.",
			},
			"outputs": schema.SingleNestedAttribute{
				Attributes:  resourcelink.ToDataSourceSchema(),
				Computed:    true,
				Optional:    false,
				Description: "A reference to a resource.",
			},
			"root_node": authenticationpolicytreenode.DataSourceSchema(),
		},
	}

	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"fragment_id", true, "The authentication policy fragment ID. ID is unique.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *authenticationPoliciesFragmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies_fragment"
}

func (r *authenticationPoliciesFragmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *authenticationPoliciesFragmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationPoliciesFragmentModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fragmentResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.GetFragment(config.DetermineAuthContext(ctx, r.providerConfig), state.FragmentId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Policy Fragment", err, httpResp)
		return
	}

	var updatedState authenticationPoliciesFragmentModel
	diags = readAuthenticationPoliciesFragmentResponse(ctx, fragmentResponse, &updatedState)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}
