package oauthaccesstokenmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthAccessTokenManagerDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthAccessTokenManagerDataSource{}
)

var (
	coreAttributeTypes = map[string]attr.Type{
		"name":         basetypes.StringType{},
		"multi_valued": basetypes.BoolType{},
	}

	extendedAttributeTypes = map[string]attr.Type{
		"name":         basetypes.StringType{},
		"multi_valued": basetypes.BoolType{},
	}

	attributeContractAttrTypes = map[string]attr.Type{
		"core_attributes":           basetypes.ListType{ElemType: types.ObjectType{AttrTypes: coreAttributeTypes}},
		"extended_attributes":       basetypes.ListType{ElemType: types.ObjectType{AttrTypes: extendedAttributeTypes}},
		"inherited":                 basetypes.BoolType{},
		"default_subject_attribute": basetypes.StringType{},
	}
)

// Create a Administrative Account data source
func NewOauthAccessTokenManagerDataSource() datasource.DataSource {
	return &oauthAccessTokenManagerDataSource{}
}

// oauthAccessTokenManagerDataSource is the datasource implementation.
type oauthAccessTokenManagerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAccessTokenManagerDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	OauthAccessTokenManagerId types.String `tfsdk:"oauth_access_token_manager_id"`
	Name                      types.String `tfsdk:"name"`
	PluginDescriptorRef       types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef                 types.Object `tfsdk:"parent_ref"`
	Configuration             types.Object `tfsdk:"configuration"`
	AttributeContract         types.Object `tfsdk:"attribute_contract"`
	SelectionSettings         types.Object `tfsdk:"selection_settings"`
	AccessControlSettings     types.Object `tfsdk:"access_control_settings"`
	SessionValidationSettings types.Object `tfsdk:"session_validation_settings"`
	SequenceNumber            types.Int64  `tfsdk:"sequence_number"`
}

// GetSchema defines the schema for the datasource.
func (r *oauthAccessTokenManagerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes OAuth Access Token Manager",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelink.ToDataSourceSchema(),
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelink.ToDataSourceSchema(),
			},
			"configuration": pluginconfiguration.ToDataSourceSchema(),
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be added to an access token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						Description: "A list of core token attributes that are associated with the access token management plugin type. This field is read-only and is ignored on POST/PUT.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
							},
						},
					},
					"extended_attributes": schema.ListNestedAttribute{
						Description: "A list of additional token attributes that are associated with this access token management plugin instance.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
							},
						},
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"default_subject_attribute": schema.StringAttribute{
						Description: "Default subject attribute to use for audit logging when validating the access token. Blank value means to use USER_KEY attribute value after grant lookup.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"selection_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine how this token manager can be selected for use by an OAuth request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"inherited": schema.BoolAttribute{
						Description: "If this token manager has a parent, this flag determines whether selection settings, such as resource URI's, are inherited from the parent. When set to true, the other fields in this model become read-only. The default value is false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"resource_uris": schema.ListAttribute{
						Description: "The list of base resource URI's which map to this token manager. A resource URI, specified via the 'aud' parameter, can be used to select a specific token manager for an OAuth request.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"access_control_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine which clients may access this token manager.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"inherited": schema.BoolAttribute{
						Description: "If this token manager has a parent, this flag determines whether access control settings are inherited from the parent. When set to true, the other fields in this model become read-only. The default value is false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"restrict_clients": schema.BoolAttribute{
						Description: "Determines whether access to this token manager is restricted to specific OAuth clients. If false, the 'allowedClients' field is ignored. The default value is false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"allowed_clients": schema.ListNestedAttribute{
						Description: "If 'restrictClients' is true, this field defines the list of OAuth clients that are allowed to access the token manager.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: resourcelink.ToDataSourceSchema(),
						},
					},
				},
			},
			"session_validation_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine how the user session is associated with the access token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"inherited": schema.BoolAttribute{
						Description: "If this token manager has a parent, this flag determines whether session validation settings, such as checkValidAuthnSession, are inherited from the parent. When set to true, the other fields in this model become read-only. The default value is false.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"include_session_id": schema.BoolAttribute{
						Description: "Include the session identifier in the access token. Note that if any of the session validation features is enabled, the session identifier will already be included in the access tokens.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"check_valid_authn_session": schema.BoolAttribute{
						Description: "Check for a valid authentication session when validating the access token.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"check_session_revocation_status": schema.BoolAttribute{
						Description: "Check the session revocation status when validating the access token.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"update_authn_session_activity": schema.BoolAttribute{
						Description: "Update authentication session activity when validating the access token.",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"sequence_number": schema.Int64Attribute{
				Description: "Number added to an access token to identify which Access Token Manager issued the token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	id.ToDataSourceSchema(&schemaDef, false, "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.")
	id.ToDataSourceSchemaCustomId(&schemaDef,
		"oauth_access_token_manager_id",
		true,
		true,
		"The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.")
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *oauthAccessTokenManagerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_access_token_manager"
}

// Configure adds the provider configured client to the data source.
func (r *oauthAccessTokenManagerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a OauthAccessTokenManagerResponse object into the model struct
func readOauthAccessTokenManagerResponseDataSource(ctx context.Context, r *client.AccessTokenManager, state *oauthAccessTokenManagerDataSourceModel, configurationFromPlan basetypes.ObjectValue) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

	state.Id = types.StringValue(r.Id)
	state.OauthAccessTokenManagerId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef, respDiags = resourcelink.ToDataSourceState(ctx, &r.PluginDescriptorRef)
	diags.Append(respDiags...)
	state.ParentRef, respDiags = resourcelink.ToDataSourceState(ctx, r.ParentRef)
	diags.Append(respDiags...)
	state.Configuration, respDiags = types.ObjectValueFrom(ctx, pluginconfiguration.AttrType(), r.Configuration)
	diags.Append(respDiags...)

	// state.AttributeContract
	if r.AttributeContract == nil {
		state.AttributeContract = types.ObjectNull(attributeContractAttrTypes)
	} else {
		state.AttributeContract, respDiags = types.ObjectValueFrom(ctx, attributeContractAttrTypes, r.AttributeContract)
		diags.Append(respDiags...)
	}

	// state.SelectionSettings
	if r.SelectionSettings == nil {
		state.SelectionSettings = types.ObjectNull(selectionSettingsAttrType)
	} else {
		state.SelectionSettings, respDiags = types.ObjectValueFrom(ctx, selectionSettingsAttrType, r.SelectionSettings)
		diags.Append(respDiags...)
	}

	// state.AccessControlSettings
	if r.AccessControlSettings == nil {
		state.AccessControlSettings = types.ObjectNull(accessControlSettingsAttrType)
	} else {
		state.AccessControlSettings, respDiags = types.ObjectValueFrom(ctx, accessControlSettingsAttrType, r.AccessControlSettings)
		diags.Append(respDiags...)
	}

	// state.SessionValidationSettings
	if r.SessionValidationSettings == nil {
		state.SessionValidationSettings = types.ObjectNull(sessionValidationSettingsAttrType)
	} else {
		state.SessionValidationSettings, respDiags = types.ObjectValueFrom(ctx, sessionValidationSettingsAttrType, r.SessionValidationSettings)
		diags.Append(respDiags...)
	}

	// state.SequenceNumber
	state.SequenceNumber = types.Int64PointerValue(r.SequenceNumber)

	return diags
}

// Read resource information
func (r *oauthAccessTokenManagerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthAccessTokenManagerDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOauthAccessTokenManager, httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.GetTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.OauthAccessTokenManagerId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Access Token Manager", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readOauthAccessTokenManagerResponseDataSource(ctx, apiReadOauthAccessTokenManager, &state, state.Configuration)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
