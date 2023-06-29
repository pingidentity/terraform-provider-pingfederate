package redirectValidation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &redirectValidationResource{}
	_ resource.ResourceWithConfigure   = &redirectValidationResource{}
	_ resource.ResourceWithImportState = &redirectValidationResource{}
)

// RedirectValidationResource is a helper function to simplify the provider implementation.
func RedirectValidationResource() resource.Resource {
	return &redirectValidationResource{}
}

// redirectValidationResource is the resource implementation.
type redirectValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type redirectValidationResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	RedirectValidationLocalSettings   types.Object `tfsdk:"redirect_validation_local_settings"`
	RedirectValidationPartnerSettings types.Object `tfsdk:"redirect_validation_partner_settings"`
}

// GetSchema defines the schema for the resource.
func (r *redirectValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	redirectValidationResourceSchema(ctx, req, resp, false)
}

func redirectValidationResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a RedirectValidation.",
		Attributes: map[string]schema.Attribute{
			"redirect_validation_local_settings": schema.SingleNestedAttribute{
				Description: "Settings for local redirect validation.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enable_target_resource_validation_for_sso": schema.BoolAttribute{
						Description: "Enable target resource validation for SSO.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_target_resource_validation_for_slo": schema.BoolAttribute{
						Description: "Enable target resource validation for SLO.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_target_resource_validation_for_idp_discovery": schema.BoolAttribute{
						Description: "Enable target resource validation for IdP discovery.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enable_in_error_resource_validation": schema.BoolAttribute{
						Description: "Enable validation for error resource.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"white_list": schema.SetNestedAttribute{
						Description: "List of URLs that are designated as valid target resources.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"target_resource_sso": schema.BoolAttribute{
									Description: "Enable this target resource for SSO redirect validation.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								"target_resource_slo": schema.BoolAttribute{
									Description: "Enable this target resource for SLO redirect validation.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								"in_error_resource": schema.BoolAttribute{
									Description: "Enable this target resource for in error resource validation.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								"idp_discovery": schema.BoolAttribute{
									Description: "Enable this target resource for IdP discovery validation.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								"valid_domain": schema.StringAttribute{
									Description: "Domain of a valid resource.",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"valid_path": schema.StringAttribute{
									Description: "Path of a valid resource.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"allow_query_and_fragment": schema.BoolAttribute{
									Description: "Allow any query parameters and fragment in the resource.",
									Computed:    true,
									Optional:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
								"require_https": schema.BoolAttribute{
									Description: "Require HTTPS for accessing this resource.",
									Computed:    true,
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
			"redirect_validation_partner_settings": schema.SingleNestedAttribute{
				Description: "Settings for partner redirect validation.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enable_wreply_validation_slo": schema.BoolAttribute{
						Description: "Enable wreply validation for SLO.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{""})
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}
func addOptionalRedirectValidationFields(ctx context.Context, addRequest *client.RedirectValidationSettings, plan redirectValidationResourceModel) error {
	if internaltypes.IsDefined(plan.RedirectValidationLocalSettings) {
		addRequest.RedirectValidationLocalSettings = client.NewRedirectValidationLocalSettings()
		redirectValidationLocalSettingsAttrs := plan.RedirectValidationLocalSettings.Attributes()
		if internaltypes.IsDefined(redirectValidationLocalSettingsAttrs["enable_target_resource_validation_for_sso"]) {
			addRequest.RedirectValidationLocalSettings.EnableTargetResourceValidationForSSO = internaltypes.AttrValueToBoolPointer(redirectValidationLocalSettingsAttrs["enable_target_resource_validation_for_sso"])
		}
		if internaltypes.IsDefined(redirectValidationLocalSettingsAttrs["enable_target_resource_validation_for_slo"]) {
			addRequest.RedirectValidationLocalSettings.EnableTargetResourceValidationForSLO = internaltypes.AttrValueToBoolPointer(redirectValidationLocalSettingsAttrs["enable_target_resource_validation_for_slo"])
		}
		if internaltypes.IsDefined(redirectValidationLocalSettingsAttrs["enable_target_resource_validation_for_idp_discovery"]) {
			addRequest.RedirectValidationLocalSettings.EnableTargetResourceValidationForIdpDiscovery = internaltypes.AttrValueToBoolPointer(redirectValidationLocalSettingsAttrs["enable_target_resource_validation_for_idp_discovery"])
		}
		if internaltypes.IsDefined(redirectValidationLocalSettingsAttrs["enable_in_error_resource_validation"]) {
			addRequest.RedirectValidationLocalSettings.EnableInErrorResourceValidation = internaltypes.AttrValueToBoolPointer(redirectValidationLocalSettingsAttrs["enable_in_error_resource_validation"])
		}

		whiteListPlan := plan.RedirectValidationLocalSettings.Attributes()["white_list"].(types.Set)
		if internaltypes.IsDefined(whiteListPlan) {
			addRequest.RedirectValidationLocalSettings.WhiteList = []client.RedirectValidationSettingsWhitelistEntry{}
			for i := 0; i < len(whiteListPlan.Elements()); i++ {
				item := whiteListPlan.Elements()[i].(types.Object)
				itemAttrs := item.Attributes()
				newWhiteList := client.NewRedirectValidationSettingsWhitelistEntryWithDefaults()
				newWhiteList.SetTargetResourceSSO(itemAttrs["target_resource_sso"].(types.Bool).ValueBool())
				newWhiteList.SetTargetResourceSLO(itemAttrs["target_resource_slo"].(types.Bool).ValueBool())
				newWhiteList.SetInErrorResource(itemAttrs["in_error_resource"].(types.Bool).ValueBool())
				newWhiteList.SetIdpDiscovery(itemAttrs["idp_discovery"].(types.Bool).ValueBool())
				newWhiteList.SetValidDomain(itemAttrs["valid_domain"].(types.String).ValueString())
				newWhiteList.SetValidPath(itemAttrs["valid_path"].(types.String).ValueString())
				newWhiteList.SetAllowQueryAndFragment(itemAttrs["allow_query_and_fragment"].(types.Bool).ValueBool())
				newWhiteList.SetRequireHttps(itemAttrs["require_https"].(types.Bool).ValueBool())
				addRequest.RedirectValidationLocalSettings.WhiteList = append(addRequest.RedirectValidationLocalSettings.WhiteList, *newWhiteList)
			}
		}
	}

	if internaltypes.IsDefined(plan.RedirectValidationPartnerSettings) {
		addRequest.RedirectValidationPartnerSettings = client.NewRedirectValidationPartnerSettings()
		enableWreplyValidationSloAttrs := plan.RedirectValidationPartnerSettings.Attributes()["enable_wreply_validation_slo"]
		if internaltypes.IsDefined(enableWreplyValidationSloAttrs) {
			enableWreplyValidationSlo := internaltypes.ConvertToPrimitive(enableWreplyValidationSloAttrs).(bool)
			addRequest.RedirectValidationPartnerSettings.EnableWreplyValidationSLO = &enableWreplyValidationSlo
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *redirectValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redirect_validation"
}

func (r *redirectValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readRedirectValidationResponse(ctx context.Context, r *client.RedirectValidationSettings, state *redirectValidationResourceModel, expectedValues *redirectValidationResourceModel) {
	state.Id = types.StringValue("id")
	stateRedirectValidationLocalSettings, _ := internaltypes.ToStateRedirectValidation(r, diag.Diagnostics{})
	state.RedirectValidationLocalSettings = stateRedirectValidationLocalSettings
	_, stateRedirectValidationPartnerSettings := internaltypes.ToStateRedirectValidation(r, diag.Diagnostics{})
	state.RedirectValidationPartnerSettings = stateRedirectValidationPartnerSettings
}

func (r *redirectValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan redirectValidationResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRedirectValidation := client.NewRedirectValidationSettings()
	err := addOptionalRedirectValidationFields(ctx, createRedirectValidation, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for RedirectValidation", err.Error())
		return
	}
	requestJson, err := createRedirectValidation.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateRedirectValidation := r.apiClient.RedirectValidationApi.UpdateRedirectValidationSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateRedirectValidation = apiCreateRedirectValidation.Body(*createRedirectValidation)
	redirectValidationResponse, httpResp, err := r.apiClient.RedirectValidationApi.UpdateRedirectValidationSettingsExecute(apiCreateRedirectValidation)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the RedirectValidation", err, httpResp)
		return
	}
	responseJson, err := redirectValidationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state redirectValidationResourceModel

	readRedirectValidationResponse(ctx, redirectValidationResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *redirectValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readRedirectValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readRedirectValidation(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state redirectValidationResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadRedirectValidation, httpResp, err := apiClient.RedirectValidationApi.GetRedirectValidationSettings(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a RedirectValidation", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadRedirectValidation.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readRedirectValidationResponse(ctx, apiReadRedirectValidation, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *redirectValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateRedirectValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateRedirectValidation(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan redirectValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state redirectValidationResourceModel
	req.State.Get(ctx, &state)
	updateRedirectValidation := apiClient.RedirectValidationApi.UpdateRedirectValidationSettings(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewRedirectValidationSettings()
	err := addOptionalRedirectValidationFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for RedirectValidation", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateRedirectValidation = updateRedirectValidation.Body(*createUpdateRequest)
	updateRedirectValidationResponse, httpResp, err := apiClient.RedirectValidationApi.UpdateRedirectValidationSettingsExecute(updateRedirectValidation)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating RedirectValidation", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateRedirectValidationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readRedirectValidationResponse(ctx, updateRedirectValidationResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *redirectValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *redirectValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
