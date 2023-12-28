package idpadapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &idpAdapterDataSource{}
	_ datasource.DataSourceWithConfigure = &idpAdapterDataSource{}
)

// IdpAdapterDataSource is a helper function to simplify the provider implementation.
func IdpAdapterDataSource() datasource.DataSource {
	return &idpAdapterDataSource{}
}

// idpAdapterDataSource is the datasource implementation.
type idpAdapterDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *idpAdapterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes an IdP adapter instance.",
		Attributes: map[string]schema.Attribute{
			"authn_ctx_class_ref": schema.StringAttribute{
				Description: "The fixed value that indicates how the user was authenticated.",
				Computed:    true,
				Optional:    false,
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name.",
				Computed:    true,
				Optional:    false,
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance.",
				Computed:    true,
				Optional:    false,
				Attributes:  resourcelink.ToDataSourceSchema(),
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance.",
				Computed:    true,
				Optional:    false,
				Attributes:  resourcelink.ToDataSourceSchema(),
			},
			"configuration": pluginconfiguration.ToDataSourceSchema(),
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the IdP adapter provides.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
						Computed:    true,
						Optional:    false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    false,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP.",
									Computed:    true,
									Optional:    false,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs.",
									Computed:    true,
									Optional:    false,
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes that can be returned by the IdP adapter.",
						Computed:    true,
						Optional:    false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    false,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP.",
									Computed:    true,
									Optional:    false,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs.",
									Computed:    true,
									Optional:    false,
								},
							},
						},
					},
					"unique_user_key_attribute": schema.StringAttribute{
						Description: "The attribute to use for uniquely identify a user's authentication sessions.",
						Computed:    true,
						Optional:    false,
					},
					"mask_ognl_values": schema.BoolAttribute{
						Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs.",
						Computed:    true,
						Optional:    false,
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"attribute_sources":              attributesources.ToDataSourceSchema(),
					"attribute_contract_fulfillment": attributecontractfulfillment.ToDataSourceSchema(),
					"issuance_criteria":              issuancecriteria.ToDataSourceSchema(),
					"inherited": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Whether this attribute mapping is inherited from its parent instance.",
					},
				},
			},
			"last_modified": schema.StringAttribute{
				Description: "The time at which the plugin instance was last changed. This property is read only and is ignored on PUT and POST requests. Supported in PF version 12.0 or later.",
				Optional:    false,
				Computed:    true,
			},
		},
	}

	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"adapter_id",
		true,
		"The ID of the plugin instance.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *idpAdapterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_adapter"
}

func (r *idpAdapterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *idpAdapterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state idpAdapterModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpAdapter, httpResp, err := r.apiClient.IdpAdaptersAPI.GetIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AdapterId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an IdpAdapter", err, httpResp)
	}

	// Read the response into the state
	readResponseDiags := readIdpAdapterResponse(ctx, apiReadIdpAdapter, &state, nil)
	resp.Diagnostics.Append(readResponseDiags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
