package oauthissuer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthIssuerDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthIssuerDataSource{}
)

// Create a Administrative Account data source
func NewOauthIssuerDataSource() datasource.DataSource {
	return &oauthIssuerDataSource{}
}

// oauthIssuerDataSource is the datasource implementation.
type oauthIssuerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthIssuerDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	CustomId    types.String `tfsdk:"custom_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Host        types.String `tfsdk:"host"`
	Path        types.String `tfsdk:"path"`
}

// GetSchema defines the schema for the datasource.
func (r *oauthIssuerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes an OAuth Issuer.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: false,
				Optional: false,
				Computed: true,
			},
			"description": schema.StringAttribute{
				Required: false,
				Optional: false,
				Computed: true,
			},
			"host": schema.StringAttribute{
				Required: false,
				Optional: false,
				Computed: true,
			},
			"path": schema.StringAttribute{
				Required: false,
				Optional: false,
				Computed: true,
			},
		},
	}
	id.ToDataSourceSchema(&schemaDef, false, "The persistent, unique ID for the virtual issuer. It can be any combination of [a-zA-Z0-9._-].")
	id.ToDataSourceSchemaCustomId(&schemaDef, true, true, "The persistent, unique ID for the virtual issuer. It can be any combination of [a-zA-Z0-9._-].")
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *oauthIssuerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_issuer"
}

// Configure adds the provider configured client to the data source.
func (r *oauthIssuerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a OauthIssuerResponse object into the model struct
func readOauthIssuerResponseDataSource(ctx context.Context, r *client.Issuer, state *oauthIssuerDataSourceModel) {
	state.Id = types.StringPointerValue(r.Id)
	state.CustomId = types.StringPointerValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringPointerValue(r.Description)
	state.Host = types.StringValue(r.Host)
	state.Path = types.StringPointerValue(r.Path)
}

// Read resource information
func (r *oauthIssuerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthIssuerDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOauthIssuer, httpResp, err := r.apiClient.OauthIssuersAPI.GetOauthIssuerById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Issuer", err, httpResp)
		return
	}

	// Read the response into the state
	readOauthIssuerResponseDataSource(ctx, apiReadOauthIssuer, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
