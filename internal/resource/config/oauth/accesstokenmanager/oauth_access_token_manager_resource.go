package oauthaccesstokenmanager

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAccessTokenManagerResource{}
	_ resource.ResourceWithConfigure   = &oauthAccessTokenManagerResource{}
	_ resource.ResourceWithImportState = &oauthAccessTokenManagerResource{}
)

var (
	attrType = map[string]attr.Type{
		"name":         types.StringType,
		"multi_valued": types.BoolType,
	}

	attributeContractTypes = map[string]attr.Type{
		"core_attributes":           types.SetType{ElemType: types.ObjectType{AttrTypes: attrType}},
		"extended_attributes":       types.SetType{ElemType: types.ObjectType{AttrTypes: attrType}},
		"default_subject_attribute": types.StringType,
	}

	selectionSettingsAttrType = map[string]attr.Type{
		"resource_uris": types.ListType{ElemType: types.StringType},
	}

	accessControlSettingsAttrType = map[string]attr.Type{
		"restrict_clients": types.BoolType,
		"allowed_clients":  types.ListType{ElemType: types.ObjectType{AttrTypes: resourcelink.AttrType()}},
	}

	sessionValidationSettingsAttrType = map[string]attr.Type{
		"include_session_id":              types.BoolType,
		"check_valid_authn_session":       types.BoolType,
		"check_session_revocation_status": types.BoolType,
		"update_authn_session_activity":   types.BoolType,
	}

	resourceUrisDefault, _      = types.ListValue(types.StringType, nil)
	selectionSettingsDefault, _ = types.ObjectValue(selectionSettingsAttrType, map[string]attr.Value{
		"resource_uris": resourceUrisDefault,
	})

	allowedClientsDefault, _        = types.ListValue(types.ObjectType{AttrTypes: resourcelink.AttrType()}, nil)
	accessControlSettingsDefault, _ = types.ObjectValue(accessControlSettingsAttrType, map[string]attr.Value{
		"restrict_clients": types.BoolValue(false),
		"allowed_clients":  allowedClientsDefault,
	})

	sessionValidationSettingsDefault, _ = types.ObjectValue(sessionValidationSettingsAttrType, map[string]attr.Value{
		"include_session_id":              types.BoolValue(false),
		"check_valid_authn_session":       types.BoolValue(false),
		"check_session_revocation_status": types.BoolValue(false),
		"update_authn_session_activity":   types.BoolValue(false),
	})
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
	ManagerId                 types.String `tfsdk:"manager_id"`
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
		Description: "Manages an OAuth access token manager plugin instance.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created.",
				Required:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Optional:    true,
				Attributes:  resourcelink.ToSchema(),
			},
			"configuration": pluginconfiguration.ToSchema(),
			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that will be added to an access token.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of core token attributes that are associated with the access token management plugin type. This field is read-only and is ignored on POST/PUT.",
						Computed:    true,
						Optional:    false,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
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
									Optional:    false,
									Computed:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional token attributes that are associated with this access token management plugin instance.",
						Required:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"multi_valued": schema.BoolAttribute{
									Description: "Indicates whether attribute value is always returned as an array.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
					},
					"default_subject_attribute": schema.StringAttribute{
						Description: "Default subject attribute to use for audit logging when validating the access token. Blank value means to use USER_KEY attribute value after grant lookup.",
						Optional:    true,
					},
				},
			},
			"selection_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine how this token manager can be selected for use by an OAuth request.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(selectionSettingsDefault),
				Attributes: map[string]schema.Attribute{
					"resource_uris": schema.ListAttribute{
						Description: "The list of base resource URI's which map to this token manager. A resource URI, specified via the 'aud' parameter, can be used to select a specific token manager for an OAuth request.",
						Optional:    true,
						Computed:    true,
						Default:     listdefault.StaticValue(resourceUrisDefault),
						ElementType: types.StringType,
					},
				},
			},
			"access_control_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine which clients may access this token manager.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(accessControlSettingsDefault),
				Attributes: map[string]schema.Attribute{
					"restrict_clients": schema.BoolAttribute{
						Description: "Determines whether access to this token manager is restricted to specific OAuth clients. If false, the 'allowedClients' field is ignored. The default value is false.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"allowed_clients": schema.ListNestedAttribute{
						Description: "If 'restrictClients' is true, this field defines the list of OAuth clients that are allowed to access the token manager.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: resourcelink.ToSchema(),
						},
						Default: listdefault.StaticValue(allowedClientsDefault),
					},
				},
			},
			"session_validation_settings": schema.SingleNestedAttribute{
				Description: "Settings which determine how the user session is associated with the access token.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(sessionValidationSettingsDefault),
				Attributes: map[string]schema.Attribute{
					"include_session_id": schema.BoolAttribute{
						Description: "Include the session identifier in the access token. Note that if any of the session validation features is enabled, the session identifier will already be included in the access tokens.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"check_valid_authn_session": schema.BoolAttribute{
						Description: "Check for a valid authentication session when validating the access token.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"check_session_revocation_status": schema.BoolAttribute{
						Description: "Check the session revocation status when validating the access token.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"update_authn_session_activity": schema.BoolAttribute{
						Description: "Update authentication session activity when validating the access token.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
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

	id.ToSchemaDeprecated(&schema, true)
	id.ToSchemaCustomId(&schema,
		"manager_id",
		true,
		true,
		"The ID of the plugin instance. The ID cannot be modified once the instance is created.")
	resp.Schema = schema
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

func (r *oauthAccessTokenManagerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *oauthAccessTokenManagerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil || state == nil {
		return
	}

	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)

	resp.Plan.Set(ctx, plan)
}

func readOauthAccessTokenManagerResponse(ctx context.Context, r *client.AccessTokenManager, state *oauthAccessTokenManagerResourceModel, configurationFromPlan types.Object, isImportRead bool) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics

	state.Id = types.StringValue(r.Id)
	state.ManagerId = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef, respDiags = resourcelink.ToState(ctx, &r.PluginDescriptorRef)
	diags.Append(respDiags...)
	state.ParentRef, respDiags = resourcelink.ToState(ctx, r.ParentRef)
	diags.Append(respDiags...)
	state.Configuration, respDiags = pluginconfiguration.ToState(configurationFromPlan, &r.Configuration, isImportRead)
	diags.Append(respDiags...)

	// state.AttributeContract
	if r.AttributeContract == nil {
		state.AttributeContract = types.ObjectNull(attributeContractTypes)
	} else {
		attrContract := r.AttributeContract

		// state.AttributeContract core_attributes
		attributeContractClientCoreAttributes := attrContract.CoreAttributes
		coreAttrs := []client.AccessTokenAttribute{}
		for _, ca := range attributeContractClientCoreAttributes {
			coreAttribute := client.AccessTokenAttribute{}
			coreAttribute.Name = ca.Name
			coreAttribute.MultiValued = ca.MultiValued
			coreAttrs = append(coreAttrs, coreAttribute)
		}
		attributeContractCoreAttributes, respDiags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrType}, coreAttrs)
		diags.Append(respDiags...)

		// state.AttributeContract extended_attributes
		attributeContractClientExtendedAttributes := attrContract.ExtendedAttributes
		extdAttrs := []client.AccessTokenAttribute{}
		for _, ea := range attributeContractClientExtendedAttributes {
			extendedAttr := client.AccessTokenAttribute{}
			extendedAttr.Name = ea.Name
			extendedAttr.MultiValued = ea.MultiValued
			extdAttrs = append(extdAttrs, extendedAttr)
		}
		attributeContractExtendedAttributes, respDiags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrType}, extdAttrs)
		diags.Append(respDiags...)

		attributeContractValues := map[string]attr.Value{
			"core_attributes":           attributeContractCoreAttributes,
			"extended_attributes":       attributeContractExtendedAttributes,
			"default_subject_attribute": types.StringPointerValue(attrContract.DefaultSubjectAttribute),
		}
		state.AttributeContract, respDiags = types.ObjectValue(attributeContractTypes, attributeContractValues)
		diags.Append(respDiags...)
	}

	// state.SelectionSettings
	if r.SelectionSettings == nil {
		state.SelectionSettings = types.ObjectNull(selectionSettingsAttrType)
	} else {
		resourceUris, respDiags := types.ListValueFrom(ctx, types.StringType, r.SelectionSettings.ResourceUris)
		diags.Append(respDiags...)

		state.SelectionSettings, respDiags = types.ObjectValue(selectionSettingsAttrType, map[string]attr.Value{
			"resource_uris": resourceUris,
		})
		diags.Append(respDiags...)
	}

	// state.AccessControlSettings
	if r.AccessControlSettings == nil {
		state.AccessControlSettings = types.ObjectNull(accessControlSettingsAttrType)
	} else {
		allowedClients, respDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: resourcelink.AttrType()}, r.AccessControlSettings.AllowedClients)
		diags.Append(respDiags...)

		state.AccessControlSettings, respDiags = types.ObjectValue(accessControlSettingsAttrType, map[string]attr.Value{
			"restrict_clients": types.BoolPointerValue(r.AccessControlSettings.RestrictClients),
			"allowed_clients":  allowedClients,
		})
		diags.Append(respDiags...)
	}

	// state.SessionValidationSettings
	if r.SessionValidationSettings == nil {
		state.SessionValidationSettings = types.ObjectNull(sessionValidationSettingsAttrType)
	} else {
		state.SessionValidationSettings, respDiags = types.ObjectValue(sessionValidationSettingsAttrType, map[string]attr.Value{
			"include_session_id":              types.BoolPointerValue(r.SessionValidationSettings.IncludeSessionId),
			"check_valid_authn_session":       types.BoolPointerValue(r.SessionValidationSettings.CheckValidAuthnSession),
			"check_session_revocation_status": types.BoolPointerValue(r.SessionValidationSettings.CheckSessionRevocationStatus),
			"update_authn_session_activity":   types.BoolPointerValue(r.SessionValidationSettings.UpdateAuthnSessionActivity),
		})
		diags.Append(respDiags...)
	}

	// state.SequenceNumber
	state.SequenceNumber = types.Int64PointerValue(r.SequenceNumber)

	return diags
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

	createOauthAccessTokenManager := client.NewAccessTokenManager(plan.ManagerId.ValueString(), plan.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalOauthAccessTokenManagerFields(ctx, createOauthAccessTokenManager, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Access Token Manager", err.Error())
		return
	}

	apiCreateOauthAccessTokenManager := r.apiClient.OauthAccessTokenManagersAPI.CreateTokenManager(config.AuthContext(ctx, r.providerConfig))
	apiCreateOauthAccessTokenManager = apiCreateOauthAccessTokenManager.Body(*createOauthAccessTokenManager)
	oauthAccessTokenManagerResponse, httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.CreateTokenManagerExecute(apiCreateOauthAccessTokenManager)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Access Token Manager", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthAccessTokenManagerResourceModel

	diags = readOauthAccessTokenManagerResponse(ctx, oauthAccessTokenManagerResponse, &state, plan.Configuration, false)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state oauthAccessTokenManagerResourceModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAccessTokenManager, httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.GetTokenManager(config.AuthContext(ctx, r.providerConfig), state.ManagerId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth Access Token Manager", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Access Token Manager", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readOauthAccessTokenManagerResponse(ctx, apiReadOauthAccessTokenManager, &state, state.Configuration, isImportRead)
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
	updateOauthAccessTokenManager := r.apiClient.OauthAccessTokenManagersAPI.UpdateTokenManager(config.AuthContext(ctx, r.providerConfig), state.ManagerId.ValueString())
	createUpdateRequest := client.NewAccessTokenManager(state.ManagerId.ValueString(), state.Name.ValueString(), *pluginDescRefResLink, *configuration)
	err := addOptionalOauthAccessTokenManagerFields(ctx, createUpdateRequest, state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for an OAuth access token manager", err.Error())
		return
	}

	updateOauthAccessTokenManager = updateOauthAccessTokenManager.Body(*createUpdateRequest)
	updateOauthAccessTokenManagerResponse, httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.UpdateTokenManagerExecute(updateOauthAccessTokenManager)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating an OAuth access token manager", err, httpResp)
		return
	}

	// Read the response
	diags = readOauthAccessTokenManagerResponse(ctx, updateOauthAccessTokenManagerResponse, &state, state.Configuration, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state oauthAccessTokenManagerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthAccessTokenManagersAPI.DeleteTokenManager(config.AuthContext(ctx, r.providerConfig), state.ManagerId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth access token manager", err, httpResp)
	}
}

func (r *oauthAccessTokenManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("manager_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
