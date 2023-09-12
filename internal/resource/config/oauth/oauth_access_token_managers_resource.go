package oauth

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAccessTokenManagersResource{}
	_ resource.ResourceWithConfigure   = &oauthAccessTokenManagersResource{}
	_ resource.ResourceWithImportState = &oauthAccessTokenManagersResource{}
)

// OauthAccessTokenManagersResource is a helper function to simplify the provider implementation.
func OauthAccessTokenManagersResource() resource.Resource {
	return &oauthAccessTokenManagersResource{}
}

// oauthAccessTokenManagersResource is the resource implementation.
type oauthAccessTokenManagersResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAccessTokenManagersResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	PluginDescriptorRef       types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef                 types.Object `tfsdk:"parent_ref"`
	Configuration             types.Object `tfsdk:"configuration"`
	AttributeContract         types.Object `tfsdk:"attribute_contract"`
	SelectionSettings         types.Object `tfsdk:"selection_settings"`
	AccessControlSettings     types.Object `tfsdk:"access_control_setting"`
	SessionValidationSettings types.Object `tfsdk:"session_validation_settings"`
	SequenceNumber            types.Int64  `tfsdk:"sequence_number"`
}

// GetSchema defines the schema for the resource.
func (r *oauthAccessTokenManagersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages Oauth Access Token Managers",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
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
					"tables": schema.SetNestedAttribute{
						Description: "List of configuration tables.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the table.",
									Required:    true,
								},
								"rows": schema.SetNestedAttribute{
									Description: "List of table rows.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"fields": schema.SetNestedAttribute{
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
														"encrypted_value": schema.StringAttribute{
															Description: "This value is not used in this provider due to the value changing on every GET request.",
															Computed:    true,
															Optional:    false,
															Default:     stringdefault.StaticString(""),
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
					"fields": schema.SetNestedAttribute{
						Description: "List of configuration fields.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
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
								"encrypted_value": schema.StringAttribute{
									Description: "This value is not used in this provider due to the value changing on every GET request.",
									Computed:    true,
									Optional:    false,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
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
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional token attributes that are associated with this access token management plugin instance.",
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
						Computed:    true,
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
					"resource_uris": schema.SetAttribute{
						Description: "The list of base resource URI's which map to this token manager. A resource URI, specified via the 'aud' parameter, can be used to select a specific token manager for an OAuth request.",
						Optional:    true,
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
						Optional:    true,
					},
					"allowed_clients": schema.SingleNestedAttribute{
						Description: "If 'restrictClients' is true, this field defines the list of OAuth clients that are allowed to access the token manager.",
						Required:    true,
						Attributes:  config.AddResourceLinkSchema(),
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
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"include_session_id": schema.BoolAttribute{
						Description: "Include the session identifier in the access token. Note that if any of the session validation features is enabled, the session identifier will already be included in the access tokens.",
						Optional:    true,
					},
					"check_valid_authn_session": schema.BoolAttribute{
						Description: "Check for a valid authentication session when validating the access token.",
						Optional:    true,
					},
					"check_session_revocation_status": schema.BoolAttribute{
						Description: "Check the session revocation status when validating the access token.",
						Optional:    true,
					},
					"update_authn_session_activity": schema.BoolAttribute{
						Description: "Update authentication session activity when validating the access token.",
						Optional:    true,
					},
				},
			},
			"sequence_number": schema.Int64Attribute{
				Description: "Number added to an access token to identify which Access Token Manager issued the token.",
				Optional:    true,
			},
		},
	}

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

func addOptionalOauthAccessTokenManagersFields(ctx context.Context, addRequest *client.AccessTokenManagers, plan oauthAccessTokenManagersResourceModel) error {

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
		addRequest.AttributeContract = client.NewOauthAccessTokenManagersAttributeContractWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, true)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
		extendedAttrsLength := len(plan.AttributeContract.Attributes()["extended_attributes"].(types.Set).Elements())
		if extendedAttrsLength == 0 {
			addRequest.AttributeContract.ExtendedAttributes = nil
		}
	}

	if internaltypes.IsDefined(plan.SelectionSettings) {
		addRequest.SelectionSettings = client.NewSelectionSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.SelectionSettings, false)), addRequest.SelectionSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AccessControlSettings) {
		addRequest.AccessControlSettings = client.NewAccessControlSettings()
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
func (r *oauthAccessTokenManagersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_access_token_managers"
}

func (r *oauthAccessTokenManagersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthAccessTokenManagersResponse(ctx context.Context, r *client.AccessTokenManagers, state *oauthAccessTokenManagersResourceModel) {
	// state.AttributeContract = (r.AttributeContract)
	// state.Id = internaltypes.StringTypeOrNil(r.Id)
	// state.Name = internaltypes.StringTypeOrNil(r.Name)
	// state.PluginDescriptorRef = (r.PluginDescriptorRef)
	// state.ParentRef = (r.ParentRef)
	// state.Configuration = (r.Configuration)
	// state.SelectionSettings = (r.SelectionSettings)
	// state.AccessControlSettings = (r.AccessControlSettings)
	// state.SessionValidationSettings = (r.SessionValidationSettings)
	// state.SequenceNumber = types.Int64Value(r.SequenceNumber)
}

func (r *oauthAccessTokenManagersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAccessTokenManagersResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAccessTokenManagers := client.NewAccessTokenManager()
	err := addOptionalOauthAccessTokenManagersFields(ctx, createOauthAccessTokenManagers, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAccessTokenManagers", err.Error())
		return
	}
	requestJson, err := createOauthAccessTokenManagers.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateOauthAccessTokenManagers := r.apiClient.OauthApi.AddAccessTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAccessTokenManagers = apiCreateOauthAccessTokenManagers.Body(*createOauthAccessTokenManagers)
	oauthAccessTokenManagersResponse, httpResp, err := r.apiClient.OauthApi.AddAccessTokenManagerExecute(apiCreateOauthAccessTokenManagers)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OauthAccessTokenManagers", err, httpResp)
		return
	}
	responseJson, err := oauthAccessTokenManagersResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthAccessTokenManagersResourceModel

	readOauthAccessTokenManagersResponse(ctx, oauthAccessTokenManagersResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenManagersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAccessTokenManagersResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAccessTokenManagers, httpResp, err := r.apiClient.OauthAccessTokenManagersApi.GetOauthAccessTokenManagersSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OauthAccessTokenManagers", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  OauthAccessTokenManagers", err, httpResp)
		}
	}
	// Log response JSON
	responseJson, err := apiReadOauthAccessTokenManagers.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthAccessTokenManagersResponse(ctx, apiReadOauthAccessTokenManagers, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAccessTokenManagersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan oauthAccessTokenManagersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthAccessTokenManagersResourceModel
	req.State.Get(ctx, &state)
	updateOauthAccessTokenManagers := r.apiClient.OauthApi.UpdateAccessTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.VALUE.ValueString())
	createUpdateRequest := client.NewAccessTokenManager()
	err := addOptionalOauthAccessTokenManagersFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAccessTokenManagers", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateOauthAccessTokenManagers = updateOauthAccessTokenManagers.Body(*createUpdateRequest)
	updateOauthAccessTokenManagersResponse, httpResp, err := r.apiClient.OauthApi.UpdateAccessTokenManagerExecute(updateOauthAccessTokenManagers)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OauthAccessTokenManagers", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateOauthAccessTokenManagersResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readOauthAccessTokenManagersResponse(ctx, updateOauthAccessTokenManagersResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthAccessTokenManagersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *oauthAccessTokenManagersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
