package config

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
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

// RedirectValidationLocalSettings Settings for local redirect validation.
type TestRedirectValidationLocalSettings struct {
	// Enable target resource validation for SSO.
	EnableTargetResourceValidationForSSO *bool `json:"enableTargetResourceValidationForSSO,omitempty" tfsdk:"enable_target_resource_validation_for_sso"`
	// Enable target resource validation for SLO.
	EnableTargetResourceValidationForSLO *bool `json:"enableTargetResourceValidationForSLO,omitempty" tfsdk:"enable_target_resource_validation_for_slo"`
	// Enable target resource validation for IdP discovery.
	EnableTargetResourceValidationForIdpDiscovery *bool `json:"enableTargetResourceValidationForIdpDiscovery,omitempty" tfsdk:"enable_target_resource_validation_for_idp_discovery"`
	// Enable validation for error resource.
	EnableInErrorResourceValidation *bool `json:"enableInErrorResourceValidation,omitempty" tfsdk:"enable_in_error_resource_validation"`
	// List of URLs that are designated as valid target resources.
	WhiteList []TestRedirectValidationSettingsWhitelistEntry `json:"whiteList,omitempty" tfsdk:"white_list"`
}

// RedirectValidationSettingsWhitelistEntry Whitelist entry for valid target resource.
type TestRedirectValidationSettingsWhitelistEntry struct {
	// Enable this target resource for SSO redirect validation.
	TargetResourceSSO *bool `json:"targetResourceSSO,omitempty" tfsdk:"target_resource_sso"`
	// Enable this target resource for SLO redirect validation.
	TargetResourceSLO *bool `json:"targetResourceSLO,omitempty" tfsdk:"target_resource_slo"`
	// Enable this target resource for in error resource validation.
	InErrorResource *bool `json:"inErrorResource,omitempty" tfsdk:"in_error_resource"`
	// Enable this target resource for IdP discovery validation.
	IdpDiscovery *bool `json:"idpDiscovery,omitempty" tfsdk:"idp_discovery"`
	// Domain of a valid resource.
	ValidDomain string `json:"validDomain" tfsdk:"valid_domain"`
	// Path of a valid resource.
	ValidPath *string `json:"validPath,omitempty" tfsdk:"valid_path"`
	// Allow any query parameters and fragment in the resource.
	AllowQueryAndFragment *bool `json:"allowQueryAndFragment,omitempty" tfsdk:"allow_query_and_fragment"`
	// Require HTTPS for accessing this resource.
	RequireHttps *bool `json:"requireHttps,omitempty" tfsdk:"require_https"`
}

// GetSchema defines the schema for the resource.
func (r *redirectValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	redirectValidationResourceSchema(ctx, req, resp, false)
}

func redirectValidationResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a RedirectValidation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder for Terraform",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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

	resp.Schema = schema
}
func addOptionalRedirectValidationFields(ctx context.Context, addRequest *client.RedirectValidationSettings, plan redirectValidationResourceModel) error {
	if internaltypes.IsDefined(plan.RedirectValidationLocalSettings) {
		addRequest.RedirectValidationLocalSettings = client.NewRedirectValidationLocalSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RedirectValidationLocalSettings)), addRequest.RedirectValidationLocalSettings)
		if err != nil {
			return err
		}
	}
	if internaltypes.IsDefined(plan.RedirectValidationPartnerSettings) {
		addRequest.RedirectValidationPartnerSettings = client.NewRedirectValidationPartnerSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RedirectValidationPartnerSettings)), addRequest.RedirectValidationPartnerSettings)
		if err != nil {
			return err
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

func readRedirectValidationResponse(ctx context.Context, r *client.RedirectValidationSettings, state *redirectValidationResourceModel) {
	state.Id = types.StringValue("id")

	redirectValidationLocalSettings := r.GetRedirectValidationLocalSettings()
	testRedirectValidationLocalSettings := TestRedirectValidationLocalSettings{}
	testRedirectValidationLocalSettings.EnableInErrorResourceValidation = redirectValidationLocalSettings.EnableInErrorResourceValidation
	testRedirectValidationLocalSettings.EnableTargetResourceValidationForIdpDiscovery = redirectValidationLocalSettings.EnableTargetResourceValidationForIdpDiscovery
	testRedirectValidationLocalSettings.EnableTargetResourceValidationForSLO = redirectValidationLocalSettings.EnableTargetResourceValidationForSLO
	testRedirectValidationLocalSettings.EnableTargetResourceValidationForSSO = redirectValidationLocalSettings.EnableTargetResourceValidationForSSO

	testRedirectValidationLocalSettings.WhiteList = []TestRedirectValidationSettingsWhitelistEntry{}
	for _, whiteListItem := range redirectValidationLocalSettings.WhiteList {
		testWhiteListItem := TestRedirectValidationSettingsWhitelistEntry{}
		testWhiteListItem.TargetResourceSLO = whiteListItem.TargetResourceSLO
		testWhiteListItem.TargetResourceSSO = whiteListItem.TargetResourceSSO
		testWhiteListItem.InErrorResource = whiteListItem.InErrorResource
		testWhiteListItem.IdpDiscovery = whiteListItem.IdpDiscovery
		testWhiteListItem.ValidDomain = whiteListItem.ValidDomain
		testWhiteListItem.ValidPath = whiteListItem.ValidPath
		testWhiteListItem.AllowQueryAndFragment = whiteListItem.AllowQueryAndFragment
		testWhiteListItem.RequireHttps = whiteListItem.RequireHttps
		testRedirectValidationLocalSettings.WhiteList = append(testRedirectValidationLocalSettings.WhiteList, testWhiteListItem)
	}

	whiteListAttrTypes := map[string]attr.Type{
		"target_resource_sso":      basetypes.BoolType{},
		"target_resource_slo":      basetypes.BoolType{},
		"in_error_resource":        basetypes.BoolType{},
		"idp_discovery":            basetypes.BoolType{},
		"valid_domain":             basetypes.StringType{},
		"valid_path":               basetypes.StringType{},
		"allow_query_and_fragment": basetypes.BoolType{},
		"require_https":            basetypes.BoolType{},
	}

	redirectValidationLocalSettingsAttrTypes := map[string]attr.Type{
		"enable_target_resource_validation_for_sso":           basetypes.BoolType{},
		"enable_target_resource_validation_for_slo":           basetypes.BoolType{},
		"enable_target_resource_validation_for_idp_discovery": basetypes.BoolType{},
		"enable_in_error_resource_validation":                 basetypes.BoolType{},
		"white_list":                                          basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: whiteListAttrTypes}},
	}

	state.RedirectValidationLocalSettings, _ = types.ObjectValueFrom(ctx, redirectValidationLocalSettingsAttrTypes, testRedirectValidationLocalSettings)

	/*whiteListAttrTypes := map[string]attr.Type{
		"target_resource_sso":      basetypes.BoolType{},
		"target_resource_slo":      basetypes.BoolType{},
		"in_error_resource":        basetypes.BoolType{},
		"idp_discovery":            basetypes.BoolType{},
		"valid_domain":             basetypes.StringType{},
		"valid_path":               basetypes.StringType{},
		"allow_query_and_fragment": basetypes.BoolType{},
		"require_https":            basetypes.BoolType{},
	}

	whiteListAttrs := r.GetRedirectValidationLocalSettings().WhiteList
	var whiteListSliceAttrVal = []attr.Value{}
	whiteListSliceType := types.ObjectType{AttrTypes: whiteListAttrTypes}
	for i := 0; i < len(whiteListAttrs); i++ {
		whiteListAttrValues := map[string]attr.Value{
			"target_resource_sso":      types.BoolPointerValue(whiteListAttrs[i].TargetResourceSSO),
			"target_resource_slo":      types.BoolPointerValue(whiteListAttrs[i].TargetResourceSLO),
			"in_error_resource":        types.BoolPointerValue(whiteListAttrs[i].InErrorResource),
			"idp_discovery":            types.BoolPointerValue(whiteListAttrs[i].IdpDiscovery),
			"valid_domain":             types.StringValue(whiteListAttrs[i].ValidDomain),
			"valid_path":               types.StringPointerValue(whiteListAttrs[i].ValidPath),
			"allow_query_and_fragment": types.BoolPointerValue(whiteListAttrs[i].AllowQueryAndFragment),
			"require_https":            types.BoolPointerValue(whiteListAttrs[i].RequireHttps),
		}
		whiteListObj, _ := types.ObjectValue(whiteListAttrTypes, whiteListAttrValues)
		whiteListSliceAttrVal = append(whiteListSliceAttrVal, whiteListObj)
	}
	whiteListSlice, _ := types.SetValue(whiteListSliceType, whiteListSliceAttrVal)

	redirectValidationLocalSettingsAttrTypes := map[string]attr.Type{
		"enable_target_resource_validation_for_sso":           basetypes.BoolType{},
		"enable_target_resource_validation_for_slo":           basetypes.BoolType{},
		"enable_target_resource_validation_for_idp_discovery": basetypes.BoolType{},
		"enable_in_error_resource_validation":                 basetypes.BoolType{},
		"white_list":                                          basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: whiteListAttrTypes}},
	}

	redirectValidationLocalSettings := r.GetRedirectValidationLocalSettings()
	redirectValidationLocalSettingsAttrVals := map[string]attr.Value{
		"enable_target_resource_validation_for_sso":           types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForSSO()),
		"enable_target_resource_validation_for_slo":           types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForSLO()),
		"enable_target_resource_validation_for_idp_discovery": types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForIdpDiscovery()),
		"enable_in_error_resource_validation":                 types.BoolValue(redirectValidationLocalSettings.GetEnableInErrorResourceValidation()),
		"white_list":                                          whiteListSlice,
	}
	redirectValidationLocalSettingsObjVal := internaltypes.MaptoObjValue(redirectValidationLocalSettingsAttrTypes, redirectValidationLocalSettingsAttrVals, diag.Diagnostics{})
	*/
	redirectValidationPartnerSettingsAttrTypes := map[string]attr.Type{
		"enable_wreply_validation_slo": basetypes.BoolType{},
	}

	redirectValidationPartnerSettingsSlo := r.GetRedirectValidationPartnerSettings().EnableWreplyValidationSLO
	redirectValidationPartnerSettingsAttrVals := map[string]attr.Value{
		"enable_wreply_validation_slo": types.BoolPointerValue(redirectValidationPartnerSettingsSlo),
	}

	redirectValidationPartnerSettingsObjVal := internaltypes.MaptoObjValue(redirectValidationPartnerSettingsAttrTypes, redirectValidationPartnerSettingsAttrVals, diag.Diagnostics{})

	//state.RedirectValidationLocalSettings = redirectValidationLocalSettingsObjVal
	state.RedirectValidationPartnerSettings = redirectValidationPartnerSettingsObjVal
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

	apiCreateRedirectValidation := r.apiClient.RedirectValidationApi.UpdateRedirectValidationSettings(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateRedirectValidation = apiCreateRedirectValidation.Body(*createRedirectValidation)
	redirectValidationResponse, httpResp, err := r.apiClient.RedirectValidationApi.UpdateRedirectValidationSettingsExecute(apiCreateRedirectValidation)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the RedirectValidation", err, httpResp)
		return
	}
	responseJson, err := redirectValidationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state redirectValidationResourceModel

	readRedirectValidationResponse(ctx, redirectValidationResponse, &state)
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
	apiReadRedirectValidation, httpResp, err := apiClient.RedirectValidationApi.GetRedirectValidationSettings(ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a RedirectValidation", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadRedirectValidation.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readRedirectValidationResponse(ctx, apiReadRedirectValidation, &state)

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
	updateRedirectValidation := apiClient.RedirectValidationApi.UpdateRedirectValidationSettings(ProviderBasicAuthContext(ctx, providerConfig))
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
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating RedirectValidation", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateRedirectValidationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readRedirectValidationResponse(ctx, updateRedirectValidationResponse, &state)

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
