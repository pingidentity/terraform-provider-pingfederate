// Copyright © 2025 Ping Identity Corporation

package redirectvalidation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &redirectValidationDataSource{}
	_ datasource.DataSourceWithConfigure = &redirectValidationDataSource{}
)

// Create a Redirect Validation Data Source
func RedirectValidationDataSource() datasource.DataSource {
	return &redirectValidationDataSource{}
}

// redirectValidationDataSource is the datasource implementation.
type redirectValidationDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *redirectValidationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes the Redirect Validation Settings.",
		Attributes: map[string]schema.Attribute{
			"redirect_validation_local_settings": schema.SingleNestedAttribute{
				Description: "Settings for local redirect validation.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"enable_target_resource_validation_for_sso": schema.BoolAttribute{
						Description: "Enable target resource validation for SSO.",
						Computed:    true,
						Optional:    false,
					},
					"enable_target_resource_validation_for_slo": schema.BoolAttribute{
						Description: "Enable target resource validation for SLO.",
						Computed:    true,
						Optional:    false,
					},
					"enable_target_resource_validation_for_idp_discovery": schema.BoolAttribute{
						Description: "Enable target resource validation for IdP discovery.",
						Computed:    true,
						Optional:    false,
					},
					"enable_in_error_resource_validation": schema.BoolAttribute{
						Description: "Enable validation for error resource.",
						Computed:    true,
						Optional:    false,
					},
					"white_list": schema.SetNestedAttribute{
						Description: "List of URLs that are designated as valid target resources.",
						Computed:    true,
						Optional:    false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"target_resource_sso": schema.BoolAttribute{
									Description: "Enable this target resource for SSO redirect validation.",
									Computed:    true,
									Optional:    false,
								},
								"target_resource_slo": schema.BoolAttribute{
									Description: "Enable this target resource for SLO redirect validation.",
									Computed:    true,
									Optional:    false,
								},
								"in_error_resource": schema.BoolAttribute{
									Description: "Enable this target resource for in error resource validation.",
									Computed:    true,
									Optional:    false,
								},
								"idp_discovery": schema.BoolAttribute{
									Description: "Enable this target resource for IdP discovery validation.",
									Computed:    true,
									Optional:    false,
								},
								"valid_domain": schema.StringAttribute{
									Description: "Domain of a valid resource.",
									Computed:    true,
									Optional:    false,
								},
								"valid_path": schema.StringAttribute{
									Description: "Path of a valid resource.",
									Computed:    true,
									Optional:    false,
								},
								"allow_query_and_fragment": schema.BoolAttribute{
									Description: "Allow any query parameters and fragment in the resource.",
									Computed:    true,
									Optional:    false,
								},
								"require_https": schema.BoolAttribute{
									Description: "Require HTTPS for accessing this resource.",
									Computed:    true,
									Optional:    false,
								},
							},
						},
					},
					"uri_allow_list": schema.SetNestedAttribute{
						Description: "List of URIs that are designated as valid target resources.",
						Computed:    true,
						Optional:    false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"target_resource_sso": schema.BoolAttribute{
									Description: "Enable this target resource for SSO redirect validation.",
									Computed:    true,
									Optional:    false,
								},
								"target_resource_slo": schema.BoolAttribute{
									Description: "Enable this target resource for SLO redirect validation.",
									Computed:    true,
									Optional:    false,
								},
								"in_error_resource": schema.BoolAttribute{
									Description: "Enable this target resource for in error resource validation.",
									Computed:    true,
									Optional:    false,
								},
								"idp_discovery": schema.BoolAttribute{
									Description: "Enable this target resource for IdP discovery validation.",
									Computed:    true,
									Optional:    false,
								},
								"allow_query_and_fragment": schema.BoolAttribute{
									Description: "Allow any query parameters and fragment in the resource.",
									Computed:    true,
									Optional:    false,
								},
								"valid_uri": schema.StringAttribute{
									Description: "URI of a valid resource.",
									Computed:    true,
									Optional:    false,
								},
							},
						},
					},
				},
			},
			"redirect_validation_partner_settings": schema.SingleNestedAttribute{
				Description: "Settings for partner redirect validation.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"enable_wreply_validation_slo": schema.BoolAttribute{
						Description: "Enable wreply validation for SLO.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
		},
	}
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *redirectValidationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redirect_validation"
}

// Configure adds the provider configured client to the data source.
func (r *redirectValidationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read resource information
func (r *redirectValidationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state redirectValidationModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadRedirectValidation, httpResp, err := r.apiClient.RedirectValidationAPI.GetRedirectValidationSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Redirect Validation Settings", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readRedirectValidationResponse(ctx, apiReadRedirectValidation, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
