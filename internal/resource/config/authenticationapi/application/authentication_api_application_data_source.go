package authenticationapiapplication

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	datasourceresourcelink "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &authenticationApiApplicationDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationApiApplicationDataSource{}
)

// AuthenticationApiApplicationDataSource is a helper function to simplify the provider implementation.
func AuthenticationApiApplicationDataSource() datasource.DataSource {
	return &authenticationApiApplicationDataSource{}
}

// authenticationApiApplicationDataSource is the datasource implementation.
type authenticationApiApplicationDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *authenticationApiApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes an Authentication Api Application",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The Authentication API Application Name. Name must be unique.",
				Computed:    true,
				Optional:    false,
			},
			"url": schema.StringAttribute{
				Description: "The Authentication API Application redirect URL.",
				Computed:    true,
				Optional:    false,
			},
			"description": schema.StringAttribute{
				Description: "The Authentication API Application description.",
				Computed:    true,
				Optional:    false,
			},
			"additional_allowed_origins": schema.ListAttribute{
				Description: "The domain in the redirect URL is always whitelisted. This field contains a list of additional allowed origin URL's for cross-origin datasource sharing.",
				Computed:    true,
				Optional:    false,
				ElementType: types.StringType,
			},
			"client_for_redirectless_mode_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription(
				"The client this application must use if it invokes the authentication API in redirectless mode. No client may be specified if restrictAccessToRedirectlessMode is false under /authenticationApi/settings.",
			),
		},
	}

	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"application_id",
		true,
		"The persistent, unique ID for the Authentication API application.",
	)
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *authenticationApiApplicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_api_application"
}

func (r *authenticationApiApplicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *authenticationApiApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationApiApplicationModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationApiApplication, httpResp, err := r.apiClient.AuthenticationApiAPI.GetApplication(config.DetermineAuthContext(ctx, r.providerConfig), state.ApplicationId.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Api Application", err, httpResp)
	}

	// Read the response into the state
	diags = readAuthenticationApiApplicationResponse(ctx, apiReadAuthenticationApiApplication, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
