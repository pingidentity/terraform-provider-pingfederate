package idp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &idpDefaultUrlsDataSource{}
	_ datasource.DataSourceWithConfigure = &idpDefaultUrlsDataSource{}
)

// Create a Administrative Account data source
func NewIdpDefaultUrlsDataSource() datasource.DataSource {
	return &idpDefaultUrlsDataSource{}
}

// idpDefaultUrlDataSource is the datasource implementation.
type idpDefaultUrlsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type idpDefaultUrlsDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	ConfirmIdpSlo    types.Bool   `tfsdk:"confirm_idp_slo"`
	IdpSloSuccessUrl types.String `tfsdk:"idp_slo_success_url"`
	IdpErrorMsg      types.String `tfsdk:"idp_error_msg"`
}

// GetSchema defines the schema for the datasource.
func (r *idpDefaultUrlsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a IdpDefaultUrls.",
		Attributes: map[string]schema.Attribute{
			"confirm_idp_slo": schema.BoolAttribute{
				Description: "Prompt user to confirm Single Logout (SLO).",
				Computed:    true,
				Optional:    true,
			},
			"idp_error_msg": schema.StringAttribute{
				Description: "Provide the error text displayed in a user's browser when an SSO operation fails.",
				Required:    true,
			},
			"idp_slo_success_url": schema.StringAttribute{
				Description: "Provide the default URL you would like to send the user to when Single Logout has succeeded.",
				Computed:    true,
				Optional:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *idpDefaultUrlsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_default_urls"
}

// Configure adds the provider configured client to the data source.
func (r *idpDefaultUrlsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readIdpDefaultUrlsResponseDataSource(ctx context.Context, r *client.IdpDefaultUrl, state *idpDefaultUrlsDataSourceModel, expectedValues *idpDefaultUrlsDataSourceModel) {
	state.Id = types.StringValue("id")
	state.ConfirmIdpSlo = types.BoolPointerValue(r.ConfirmIdpSlo)
	state.IdpSloSuccessUrl = internaltypes.StringTypeOrNil(r.IdpSloSuccessUrl, false)
	state.IdpErrorMsg = types.StringValue(r.IdpErrorMsg)
}

// Read resource information
func (r *idpDefaultUrlsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state idpDefaultUrlsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadIdpDefaultUrls, httpResp, err := r.apiClient.IdpDefaultUrlsApi.GetDefaultUrl(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Access Control Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := apiReadIdpDefaultUrls.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIdpDefaultUrlsResponseDataSource(ctx, apiReadIdpDefaultUrls, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
