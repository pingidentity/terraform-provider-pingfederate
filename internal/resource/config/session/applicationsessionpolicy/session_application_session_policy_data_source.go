package sessionapplicationsessionpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sessionApplicationSessionPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionApplicationSessionPolicyDataSource{}
)

// SessionApplicationSessionPolicyDataSource is a helper function to simplify the provider implementation.
func NewSessionApplicationSessionPolicyDataSource() datasource.DataSource {
	return &sessionApplicationSessionPolicyDataSource{}
}

// sessionApplicationSessionPolicyDataSource is the datasource implementation.
type sessionApplicationSessionPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type sessionApplicationSessionPolicyDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	IdleTimeoutMins types.Int64  `tfsdk:"idle_timeout_mins"`
	MaxTimeoutMins  types.Int64  `tfsdk:"max_timeout_mins"`
}

// GetSchema defines the schema for the datasource.
func (r *sessionApplicationSessionPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a SessionApplicationSessionPolicy.",
		Attributes: map[string]schema.Attribute{
			// Add necessary attributes here
			"idle_timeout_mins": schema.Int64Attribute{
				Description: "The idle timeout period, in minutes. If set to -1, the idle timeout will be set to the maximum timeout. The default is 60.",
				Computed:    true,
				Optional:    true,
			},
			"max_timeout_mins": schema.Int64Attribute{
				Description: "The maximum timeout period, in minutes. If set to -1, sessions do not expire. The default is 480.",
				Computed:    true,
				Optional:    true,
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *sessionApplicationSessionPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_application_session_policy"
}

func (r *sessionApplicationSessionPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readSessionApplicationSessionPolicyResponseDataSource(ctx context.Context, r *client.ApplicationSessionPolicy, state *sessionApplicationSessionPolicyDataSourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue("sessionApplicationSessionPolicyId")
	state.IdleTimeoutMins = types.Int64Value(r.GetIdleTimeoutMins())
	state.MaxTimeoutMins = types.Int64Value(r.GetMaxTimeoutMins())
	diags.Append(respDiags...)
	return diags
}

func (r *sessionApplicationSessionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionApplicationSessionPolicyDataSourceModel
	var diags diag.Diagnostics

	//apiReadSessionApplicationSessionPolicy, httpResp, err := r.apiClient.SessionAPI.GetApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	apiReadSessionApplicationSessionPolicy, httpResp, err := r.apiClient.SessionAPI.GetApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	// Read the response into the state
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Identity Profile", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readSessionApplicationSessionPolicyResponseDataSource(ctx, apiReadSessionApplicationSessionPolicy, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
