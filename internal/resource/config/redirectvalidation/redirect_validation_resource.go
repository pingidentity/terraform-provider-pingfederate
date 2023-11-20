package redirectvalidation

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &redirectValidationResource{}
	_ resource.ResourceWithConfigure   = &redirectValidationResource{}
	_ resource.ResourceWithImportState = &redirectValidationResource{}
)

var (
	whiteListAttrTypes = map[string]attr.Type{
		"target_resource_sso":      basetypes.BoolType{},
		"target_resource_slo":      basetypes.BoolType{},
		"in_error_resource":        basetypes.BoolType{},
		"idp_discovery":            basetypes.BoolType{},
		"valid_domain":             basetypes.StringType{},
		"valid_path":               basetypes.StringType{},
		"allow_query_and_fragment": basetypes.BoolType{},
		"require_https":            basetypes.BoolType{},
	}

	redirectValidationLocalSettingsAttrTypes = map[string]attr.Type{
		"enable_target_resource_validation_for_sso":           basetypes.BoolType{},
		"enable_target_resource_validation_for_slo":           basetypes.BoolType{},
		"enable_target_resource_validation_for_idp_discovery": basetypes.BoolType{},
		"enable_in_error_resource_validation":                 basetypes.BoolType{},
		"white_list":                                          basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: whiteListAttrTypes}},
	}

	redirectValidationPartnerSettingsAttrTypes = map[string]attr.Type{
		"enable_wreply_validation_slo": basetypes.BoolType{},
	}

	whiteListDefault, _                       = types.ListValue(types.ObjectType{AttrTypes: whiteListAttrTypes}, nil)
	redirectValidationLocalSettingsDefault, _ = types.ObjectValue(redirectValidationLocalSettingsAttrTypes, map[string]attr.Value{
		"enable_target_resource_validation_for_sso":           types.BoolValue(false),
		"enable_target_resource_validation_for_slo":           types.BoolValue(false),
		"enable_target_resource_validation_for_idp_discovery": types.BoolValue(false),
		"enable_in_error_resource_validation":                 types.BoolValue(false),
		"white_list":                                          whiteListDefault,
	})

	redirectValidationPartnerSettingsDefault, _ = types.ObjectValue(redirectValidationPartnerSettingsAttrTypes, map[string]attr.Value{
		"enable_wreply_validation_slo": types.BoolValue(false),
	})
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
	schema := schema.Schema{
		Description: "Manages the settings for redirect validation.",
		Attributes: map[string]schema.Attribute{
			"redirect_validation_local_settings": schema.SingleNestedAttribute{
				Description: "Settings for local redirect validation.",
				Computed:    true,
				Optional:    true,
				Default:     objectdefault.StaticValue(redirectValidationLocalSettingsDefault),
				Attributes: map[string]schema.Attribute{
					"enable_target_resource_validation_for_sso": schema.BoolAttribute{
						Description: "Enable target resource validation for SSO.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"enable_target_resource_validation_for_slo": schema.BoolAttribute{
						Description: "Enable target resource validation for SLO.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"enable_target_resource_validation_for_idp_discovery": schema.BoolAttribute{
						Description: "Enable target resource validation for IdP discovery.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"enable_in_error_resource_validation": schema.BoolAttribute{
						Description: "Enable validation for error resource.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"white_list": schema.ListNestedAttribute{
						Description: "List of URLs that are designated as valid target resources.",
						Computed:    true,
						Optional:    true,
						Default:     listdefault.StaticValue(whiteListDefault),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"target_resource_sso": schema.BoolAttribute{
									Description: "Enable this target resource for SSO redirect validation.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"target_resource_slo": schema.BoolAttribute{
									Description: "Enable this target resource for SLO redirect validation.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"in_error_resource": schema.BoolAttribute{
									Description: "Enable this target resource for in error resource validation.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"idp_discovery": schema.BoolAttribute{
									Description: "Enable this target resource for IdP discovery validation.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"valid_domain": schema.StringAttribute{
									Description: "Domain of a valid resource.",
									Required:    true,
								},
								"valid_path": schema.StringAttribute{
									Description: "Path of a valid resource.",
									Computed:    true,
									Optional:    true,
									Default:     stringdefault.StaticString(""),
								},
								"allow_query_and_fragment": schema.BoolAttribute{
									Description: "Allow any query parameters and fragment in the resource.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"require_https": schema.BoolAttribute{
									Description: "Require HTTPS for accessing this resource.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
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
				Default:     objectdefault.StaticValue(redirectValidationPartnerSettingsDefault),
				Attributes: map[string]schema.Attribute{
					"enable_wreply_validation_slo": schema.BoolAttribute{
						Description: "Enable wreply validation for SLO.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalRedirectValidationFields(ctx context.Context, addRequest *client.RedirectValidationSettings, plan redirectValidationResourceModel) error {
	if internaltypes.IsDefined(plan.RedirectValidationLocalSettings) {
		addRequest.RedirectValidationLocalSettings = client.NewRedirectValidationLocalSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RedirectValidationLocalSettings, false)), addRequest.RedirectValidationLocalSettings)
		if err != nil {
			return err
		}
	}
	if internaltypes.IsDefined(plan.RedirectValidationPartnerSettings) {
		addRequest.RedirectValidationPartnerSettings = client.NewRedirectValidationPartnerSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RedirectValidationPartnerSettings, false)), addRequest.RedirectValidationPartnerSettings)
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

func readRedirectValidationResponse(ctx context.Context, r *client.RedirectValidationSettings, state *redirectValidationResourceModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = id.GenerateUUIDToState(existingId)
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
		whiteListObj, respDiags := types.ObjectValue(whiteListAttrTypes, whiteListAttrValues)
		diags.Append(respDiags...)
		whiteListSliceAttrVal = append(whiteListSliceAttrVal, whiteListObj)
	}
	whiteListSlice, respDiags := types.ListValue(whiteListSliceType, whiteListSliceAttrVal)
	diags.Append(respDiags...)
	redirectValidationLocalSettings := r.GetRedirectValidationLocalSettings()
	redirectValidationLocalSettingsAttrVals := map[string]attr.Value{
		"enable_target_resource_validation_for_sso":           types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForSSO()),
		"enable_target_resource_validation_for_slo":           types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForSLO()),
		"enable_target_resource_validation_for_idp_discovery": types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForIdpDiscovery()),
		"enable_in_error_resource_validation":                 types.BoolValue(redirectValidationLocalSettings.GetEnableInErrorResourceValidation()),
		"white_list":                                          whiteListSlice,
	}
	redirectValidationLocalSettingsObjVal := internaltypes.MaptoObjValue(redirectValidationLocalSettingsAttrTypes, redirectValidationLocalSettingsAttrVals, &diags)
	redirectValidationPartnerSettingsSlo := r.GetRedirectValidationPartnerSettings().EnableWreplyValidationSLO
	redirectValidationPartnerSettingsAttrVals := map[string]attr.Value{
		"enable_wreply_validation_slo": types.BoolPointerValue(redirectValidationPartnerSettingsSlo),
	}

	redirectValidationPartnerSettingsObjVal := internaltypes.MaptoObjValue(redirectValidationPartnerSettingsAttrTypes, redirectValidationPartnerSettingsAttrVals, &diags)
	state.RedirectValidationLocalSettings = redirectValidationLocalSettingsObjVal
	state.RedirectValidationPartnerSettings = redirectValidationPartnerSettingsObjVal
	return diags
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Redirect Validation", err.Error())
		return
	}

	apiCreateRedirectValidation := r.apiClient.RedirectValidationAPI.UpdateRedirectValidationSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateRedirectValidation = apiCreateRedirectValidation.Body(*createRedirectValidation)
	redirectValidationResponse, httpResp, err := r.apiClient.RedirectValidationAPI.UpdateRedirectValidationSettingsExecute(apiCreateRedirectValidation)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Redirect Validation", err, httpResp)
		return
	}

	// Read the response into the state
	var state redirectValidationResourceModel
	diags = readRedirectValidationResponse(ctx, redirectValidationResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *redirectValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state redirectValidationResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadRedirectValidation, httpResp, err := r.apiClient.RedirectValidationAPI.GetRedirectValidationSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Redirect Validation", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Redirect Validation", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readRedirectValidationResponse(ctx, apiReadRedirectValidation, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *redirectValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan redirectValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRedirectValidation := r.apiClient.RedirectValidationAPI.UpdateRedirectValidationSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewRedirectValidationSettings()
	err := addOptionalRedirectValidationFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Redirect Validation", err.Error())
		return
	}

	updateRedirectValidation = updateRedirectValidation.Body(*createUpdateRequest)
	updateRedirectValidationResponse, httpResp, err := r.apiClient.RedirectValidationAPI.UpdateRedirectValidationSettingsExecute(updateRedirectValidation)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Redirect Validation", err, httpResp)
		return
	}

	// Read the response
	var state redirectValidationResourceModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readRedirectValidationResponse(ctx, updateRedirectValidationResponse, &state, id)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *redirectValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *redirectValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
