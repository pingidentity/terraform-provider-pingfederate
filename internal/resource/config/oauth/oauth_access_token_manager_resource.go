package oauth

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAccessTokenManagerResource{}
	_ resource.ResourceWithConfigure   = &oauthAccessTokenManagerResource{}
	_ resource.ResourceWithImportState = &oauthAccessTokenManagerResource{}
)

// OauthAccessTokenManagerResource is a helper function to simplify the provider implementation.
func OauthAccessTokenManagerResource() resource.Resource {
	return &oauthAccessTokenManagerResource{}
}

// oauthAccessTokenManagerResource is the resource implementation.
type oauthAccessTokenManagerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAccessTokenManagerResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	CustomId                  types.String `tfsdk:"custom_id"`
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

// GetSchema defines the schema for the resource.
func (r *oauthAccessTokenManagerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthAccessTokenManagerResourceSchema(ctx, req, resp, false)
}

func oauthAccessTokenManagerResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages OAuth Access Token Manager",
		Attributes: map[string]schema.Attribute{
			"custom_id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[a-zA-Z0-9_]{1,32}$"),
						"The plugin ID must be less than 33 characters, contain no spaces, and be alphanumeric.",
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				Attributes:  config.AddResourceLinkSchema(),
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: config.AddResourceLinkSchema(),
			},
			"configuration": schema.SingleNestedAttribute{
				Description: "Plugin instance configuration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"tables": schema.ListNestedAttribute{
						Description: "List of configuration tables.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the table.",
									Required:    true,
								},
								"rows": schema.ListNestedAttribute{
									Description: "List of table rows.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"fields": schema.ListNestedAttribute{
												Description: "The configuration fields in the row.",
												Computed:    true,
												Optional:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Description: "The name of the configuration field.",
															Required:    true,
														},
														"value": schema.StringAttribute{
															Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
															Required:    true,
														},
														"inherited": schema.BoolAttribute{
															Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
															Optional:    true,
															PlanModifiers: []planmodifier.Bool{
																boolplanmodifier.UseStateForUnknown(),
															},
														},
													},
												},
											},
											"default_row": schema.BoolAttribute{
												Description: "Whether this row is the default.",
												Optional:    true,
												PlanModifiers: []planmodifier.Bool{
													boolplanmodifier.UseStateForUnknown(),
												},
											},
										},
									},
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this table is inherited from its parent instance. If true, the rows become read-only. The default value is false.",
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"fields": schema.ListNestedAttribute{
						Description: "List of configuration fields.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the configuration field.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
									Required:    true,
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be added to an access token.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.ListNestedAttribute{
						Description: "A list of core token attributes that are associated with the access token management plugin type. This field is read-only and is ignored on POST/PUT.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Computed:    true,
									Optional:    false,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"extended_attributes": schema.ListNestedAttribute{
						Description: "A list of additional token attributes that are associated with this access token management plugin instance.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
						Optional:    true,
					},
					"default_subject_attribute": schema.StringAttribute{
						Description: "Default subject attribute to use for audit logging when validating the access token. Blank value means to use USER_KEY attribute value after grant lookup.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"selection_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine how this token manager can be selected for use by an OAuth request.",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"inherited": schema.BoolAttribute{
						Description: "If this token manager has a parent, this flag determines whether selection settings, such as resource URI's, are inherited from the parent. When set to true, the other fields in this model become read-only. The default value is false.",
						Optional:    true,
					},
					"resource_uris": schema.ListAttribute{
						Description: "The list of base resource URI's which map to this token manager. A resource URI, specified via the 'aud' parameter, can be used to select a specific token manager for an OAuth request.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"access_control_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine which clients may access this token manager.",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"inherited": schema.BoolAttribute{
						Description: "If this token manager has a parent, this flag determines whether access control settings are inherited from the parent. When set to true, the other fields in this model become read-only. The default value is false.",
						Optional:    true,
					},
					"restrict_clients": schema.BoolAttribute{
						Description: "Determines whether access to this token manager is restricted to specific OAuth clients. If false, the 'allowedClients' field is ignored. The default value is false.",
						Computed:    true,
						Optional:    true,
					},
					"allowed_clients": schema.ListNestedAttribute{
						Description: "If 'restrictClients' is true, this field defines the list of OAuth clients that are allowed to access the token manager.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: config.AddResourceLinkSchema(),
						},
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"session_validation_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine how the user session is associated with the access token.",
				Computed:    true,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"inherited": schema.BoolAttribute{
						Description: "If this token manager has a parent, this flag determines whether session validation settings, such as checkValidAuthnSession, are inherited from the parent. When set to true, the other fields in this model become read-only. The default value is false.",
						Optional:    true,
					},
					"include_session_id": schema.BoolAttribute{
						Description: "Include the session identifier in the access token. Note that if any of the session validation features is enabled, the session identifier will already be included in the access tokens.",
						Computed:    true,
						Optional:    true,
					},
					"check_valid_authn_session": schema.BoolAttribute{
						Description: "Check for a valid authentication session when validating the access token.",
						Computed:    true,
						Optional:    true,
					},
					"check_session_revocation_status": schema.BoolAttribute{
						Description: "Check the session revocation status when validating the access token.",
						Computed:    true,
						Optional:    true,
					},
					"update_authn_session_activity": schema.BoolAttribute{
						Description: "Update authentication session activity when validating the access token.",
						Computed:    true,
						Optional:    true,
					},
				},
			},
			"sequence_number": schema.Int64Attribute{
				Description: "Number added to an access token to identify which Access Token Manager issued the token.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	config.AddCommonSchema(&schema)
	resp.Schema = schema
}

func (r *oauthAccessTokenManagerResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model oauthAccessTokenManagerResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if internaltypes.IsDefined(model.AttributeContract) {
		if len(model.AttributeContract.Attributes()["extended_attributes"].(types.List).Elements()) == 0 {
			resp.Diagnostics.AddError("Empty set!", "Please provide valid properties within extended_attributes. The set cannot be empty.\nIf no values are necessary, remove this property from your terraform file.")
		}
	}
}

func addOptionalOauthAccessTokenManagerFields(ctx context.Context, addRequest *client.AccessTokenManager, plan oauthAccessTokenManagerResourceModel) error {

	if internaltypes.IsDefined(plan.ParentRef) {
		if plan.ParentRef.Attributes()["id"].(types.String).ValueString() != "" {
			addRequest.ParentRef = client.NewResourceLinkWithDefaults()
			addRequest.ParentRef.Id = plan.ParentRef.Attributes()["id"].(types.String).ValueString()
			err := json.Unmarshal([]byte(internaljson.FromValue(plan.ParentRef, true)), addRequest.ParentRef)
			if err != nil {
				return err
			}
		}
	}

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = client.NewAccessTokenAttributeContractWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, true)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
		extendedAttrsLength := len(plan.AttributeContract.Attributes()["extended_attributes"].(types.List).Elements())
		if extendedAttrsLength == 0 {
			addRequest.AttributeContract.ExtendedAttributes = nil
		}
	}

	if internaltypes.IsDefined(plan.SelectionSettings) {
		addRequest.SelectionSettings = client.NewAtmSelectionSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.SelectionSettings, false)), addRequest.SelectionSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AccessControlSettings) {
		addRequest.AccessControlSettings = client.NewAtmAccessControlSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AccessControlSettings, false)), addRequest.AccessControlSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.SessionValidationSettings) {
		addRequest.SessionValidationSettings = client.NewSessionValidationSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.SessionValidationSettings, false)), addRequest.SessionValidationSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.SequenceNumber) {
		addRequest.SequenceNumber = plan.SequenceNumber.ValueInt64Pointer()
	}

	return nil

}

// Metadata returns the resource type name.
func (r *oauthAccessTokenManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_access_token_manager"
}

func (r *oauthAccessTokenManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthAccessTokenManagerResponse(ctx context.Context, r *client.AccessTokenManager, state *oauthAccessTokenManagerResourceModel, configurationFromPlan basetypes.ObjectValue) diag.Diagnostics {
	state.Id = types.StringValue(r.Id)
	state.CustomId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)

	// state.pluginDescriptorRef
	pluginDescRef := r.GetPluginDescriptorRef()
	state.PluginDescriptorRef = internaltypes.ToStateResourceLink(ctx, pluginDescRef)

	// state.parentRef
	parentRef := r.GetParentRef()
	state.ParentRef = internaltypes.ToStateResourceLink(ctx, parentRef)

	// state.Configuration
	configurationAttrType := map[string]attr.Type{
		"fields": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: config.FieldAttrTypes()}},
		"tables": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: config.TableAttrTypes()}},
	}

	planFields := types.ListNull(types.ObjectType{AttrTypes: config.FieldAttrTypes()})
	planTables := types.ListNull(types.ObjectType{AttrTypes: config.TableAttrTypes()})

	planFieldsValue, ok := configurationFromPlan.Attributes()["fields"]
	if ok {
		planFields = planFieldsValue.(types.List)
	}
	planTablesValue, ok := configurationFromPlan.Attributes()["tables"]
	if ok {
		planTables = planTablesValue.(types.List)
	}

	var respDiags, diags diag.Diagnostics
	fieldsAttrValue := config.ToFieldsListValue(r.Configuration.Fields, planFields, &diags)
	tablesAttrValue := config.ToTablesListValue(r.Configuration.Tables, planTables, &diags)

	configurationAttrValue := map[string]attr.Value{
		"fields": fieldsAttrValue,
		"tables": tablesAttrValue,
	}
	state.Configuration, diags = types.ObjectValue(configurationAttrType, configurationAttrValue)
	respDiags.Append(diags...)

	// state.AttributeContract
	attrContract := r.AttributeContract

	attrType := map[string]attr.Type{
		"name":         basetypes.StringType{},
		"multi_valued": basetypes.BoolType{},
	}

	// state.AttributeContract core_attributes
	attributeContractClientCoreAttributes := attrContract.CoreAttributes
	coreAttrs := []client.AccessTokenAttribute{}
	for _, ca := range attributeContractClientCoreAttributes {
		coreAttribute := client.AccessTokenAttribute{}
		coreAttribute.Name = ca.Name
		coreAttribute.MultiValued = ca.MultiValued
		coreAttrs = append(coreAttrs, coreAttribute)
	}
	attributeContractCoreAttributes, diags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, coreAttrs)
	respDiags.Append(diags...)

	// state.AttributeContract extended_attributes
	attributeContractClientExtendedAttributes := attrContract.ExtendedAttributes
	extdAttrs := []client.AccessTokenAttribute{}
	for _, ea := range attributeContractClientExtendedAttributes {
		extendedAttr := client.AccessTokenAttribute{}
		extendedAttr.Name = ea.Name
		extendedAttr.MultiValued = ea.MultiValued
		extdAttrs = append(extdAttrs, extendedAttr)
	}
	attributeContractExtendedAttributes, diags := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: attrType}, extdAttrs)
	respDiags.Append(diags...)

	attributeContractTypes := map[string]attr.Type{
		"core_attributes":           basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"extended_attributes":       basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: attrType}},
		"inherited":                 basetypes.BoolType{},
		"default_subject_attribute": basetypes.StringType{},
	}

	attributeContractValues := map[string]attr.Value{
		"core_attributes":           attributeContractCoreAttributes,
		"extended_attributes":       attributeContractExtendedAttributes,
		"inherited":                 types.BoolPointerValue(attrContract.Inherited),
		"default_subject_attribute": types.StringPointerValue(attrContract.DefaultSubjectAttribute),
	}
	state.AttributeContract, diags = types.ObjectValue(attributeContractTypes, attributeContractValues)
	respDiags.Append(diags...)

	// state.SelectionSettings
	selectionSettingsAttrType := map[string]attr.Type{
		"inherited":     basetypes.BoolType{},
		"resource_uris": basetypes.ListType{ElemType: basetypes.StringType{}},
	}

	state.SelectionSettings, diags = types.ObjectValueFrom(ctx, selectionSettingsAttrType, r.SelectionSettings)
	respDiags.Append(diags...)

	// state.AccessControlSettings
	accessControlSettingsAttrType := map[string]attr.Type{
		"inherited":        basetypes.BoolType{},
		"restrict_clients": basetypes.BoolType{},
		"allowed_clients":  basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: internaltypes.ResourceLinkStateAttrType()}},
	}

	state.AccessControlSettings, diags = types.ObjectValueFrom(ctx, accessControlSettingsAttrType, r.AccessControlSettings)
	respDiags.Append(diags...)

	// state.SessionValidationSettings
	sessionValidationSettingsAttrType := map[string]attr.Type{
		"inherited":                       basetypes.BoolType{},
		"include_session_id":              basetypes.BoolType{},
		"check_valid_authn_session":       basetypes.BoolType{},
		"check_session_revocation_status": basetypes.BoolType{},
		"update_authn_session_activity":   basetypes.BoolType{},
	}

	state.SessionValidationSettings, diags = types.ObjectValueFrom(ctx, sessionValidationSettingsAttrType, r.SessionValidationSettings)
	respDiags.Append(diags...)

	// state.SequenceNumber
	state.SequenceNumber = types.Int64PointerValue(r.SequenceNumber)

	return respDiags
}

func (r *oauthAccessTokenManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAccessTokenManagerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PluginDescriptorRef
	pluginDescRefId := plan.PluginDescriptorRef.Attributes()["id"].(types.String).ValueString()
	pluginDescRefResLink := client.NewResourceLinkWithDefaults()
	pluginDescRefResLink.Id = pluginDescRefId
	pluginDescRefErr := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), pluginDescRefResLink)
	if pluginDescRefErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin descriptor ref request object:", pluginDescRefErr.Error())
		return
	}

	// Configuration
	configuration := client.NewPluginConfigurationWithDefaults()
	configErr := json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, true)), configuration)
	if configErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin configuration request object:", configErr.Error())
		return
	}

	createOauthAccessTokenManager := client.NewAccessTokenManager(plan.CustomId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalOauthAccessTokenManagerFields(ctx, createOauthAccessTokenManager, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Access Token Manager", err.Error())
		return
	}
	_, requestErr := createOauthAccessTokenManager.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of an OAuth Access Token Manager: %s", requestErr.Error())
	}

	apiCreateOauthAccessTokenManager := r.apiClient.OauthAccessTokenManagersApi.CreateTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAccessTokenManager = apiCreateOauthAccessTokenManager.Body(*createOauthAccessTokenManager)
	oauthAccessTokenManagerResponse, httpResp, err := r.apiClient.OauthAccessTokenManagersApi.CreateTokenManagerExecute(apiCreateOauthAccessTokenManager)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Access Token Manager", err, httpResp)
		return
	}
	_, responseErr := oauthAccessTokenManagerResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an OAuth Access Token Manager: %s", requestErr.Error())
	}

	// Read the response into the state
	var state oauthAccessTokenManagerResourceModel

	diags = readOauthAccessTokenManagerResponse(ctx, oauthAccessTokenManagerResponse, &state, plan.Configuration)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAccessTokenManagerResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAccessTokenManager, httpResp, err := r.apiClient.OauthAccessTokenManagersApi.GetTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Access Token Manager", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Access Token Manager", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadOauthAccessTokenManager.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an OAuth Access Token Manager: %s", responseErr.Error())
	}

	// Read the response into the state
	diags = readOauthAccessTokenManagerResponse(ctx, apiReadOauthAccessTokenManager, &state, state.Configuration)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAccessTokenManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state oauthAccessTokenManagerResourceModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PluginDescriptorRef
	pluginDescRefId := state.PluginDescriptorRef.Attributes()["id"].(types.String).ValueString()
	pluginDescRefResLink := client.NewResourceLinkWithDefaults()
	pluginDescRefResLink.Id = pluginDescRefId
	pluginDescRefErr := json.Unmarshal([]byte(internaljson.FromValue(state.PluginDescriptorRef, false)), pluginDescRefResLink)
	if pluginDescRefErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin descriptor ref request object:", pluginDescRefErr.Error())
		return
	}

	// Configuration
	configuration := client.NewPluginConfiguration()
	configErr := json.Unmarshal([]byte(internaljson.FromValue(state.Configuration, true)), configuration)
	if configErr != nil {
		resp.Diagnostics.AddError("Failed to build plugin configuration request object:", configErr.Error())
		return
	}

	// Get the current state to see how any attributes are changing
	updateOauthAccessTokenManager := r.apiClient.OauthAccessTokenManagersApi.UpdateTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString())
	createUpdateRequest := client.NewAccessTokenManager(state.CustomId.ValueString(), state.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalOauthAccessTokenManagerFields(ctx, createUpdateRequest, state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Access Token Manager", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of an OAuth Access Token Manager: %s", requestErr.Error())
	}
	updateOauthAccessTokenManager = updateOauthAccessTokenManager.Body(*createUpdateRequest)
	updateOauthAccessTokenManagerResponse, httpResp, err := r.apiClient.OauthAccessTokenManagersApi.UpdateTokenManagerExecute(updateOauthAccessTokenManager)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating an OAuth Access Token Manager", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateOauthAccessTokenManagerResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an OAuth Access Token Manager: %s", responseErr.Error())
	}
	// Read the response
	diags = readOauthAccessTokenManagerResponse(ctx, updateOauthAccessTokenManagerResponse, &state, state.Configuration)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthAccessTokenManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state oauthAccessTokenManagerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthAccessTokenManagersApi.DeleteTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Access Token Manager", err, httpResp)
		return
	}
}

func (r *oauthAccessTokenManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("custom_id"), req, resp)
}
