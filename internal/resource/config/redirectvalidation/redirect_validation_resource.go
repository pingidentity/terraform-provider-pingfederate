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
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &redirectValidationResource{}
	_ resource.ResourceWithConfigure   = &redirectValidationResource{}
	_ resource.ResourceWithImportState = &redirectValidationResource{}
)

var (
	whiteListAttrTypes = map[string]attr.Type{
		"target_resource_sso":      types.BoolType,
		"target_resource_slo":      types.BoolType,
		"in_error_resource":        types.BoolType,
		"idp_discovery":            types.BoolType,
		"valid_domain":             types.StringType,
		"valid_path":               types.StringType,
		"allow_query_and_fragment": types.BoolType,
		"require_https":            types.BoolType,
	}
	uriAllowListAttrTypes = map[string]attr.Type{
		"target_resource_sso":      types.BoolType,
		"target_resource_slo":      types.BoolType,
		"in_error_resource":        types.BoolType,
		"idp_discovery":            types.BoolType,
		"allow_query_and_fragment": types.BoolType,
		"valid_uri":                types.StringType,
	}

	redirectValidationLocalSettingsAttrTypes = map[string]attr.Type{
		"enable_target_resource_validation_for_sso":           types.BoolType,
		"enable_target_resource_validation_for_slo":           types.BoolType,
		"enable_target_resource_validation_for_idp_discovery": types.BoolType,
		"enable_in_error_resource_validation":                 types.BoolType,
		"white_list":                                          types.ListType{ElemType: types.ObjectType{AttrTypes: whiteListAttrTypes}},
		"uri_allow_list":                                      types.ListType{ElemType: types.ObjectType{AttrTypes: uriAllowListAttrTypes}},
	}

	redirectValidationPartnerSettingsAttrTypes = map[string]attr.Type{
		"enable_wreply_validation_slo": types.BoolType,
	}

	whiteListDefault, _                       = types.ListValue(types.ObjectType{AttrTypes: whiteListAttrTypes}, nil)
	redirectValidationLocalSettingsDefault, _ = types.ObjectValue(redirectValidationLocalSettingsAttrTypes, map[string]attr.Value{
		"enable_target_resource_validation_for_sso":           types.BoolValue(false),
		"enable_target_resource_validation_for_slo":           types.BoolValue(false),
		"enable_target_resource_validation_for_idp_discovery": types.BoolValue(false),
		"enable_in_error_resource_validation":                 types.BoolValue(false),
		"white_list":                                          whiteListDefault,
		"uri_allow_list":                                      types.ListNull(types.ObjectType{AttrTypes: uriAllowListAttrTypes}),
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
					"uri_allow_list": schema.ListNestedAttribute{
						Description: "List of URIs that are designated as valid target resources.",
						Optional:    true,
						Computed:    true,
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
								"allow_query_and_fragment": schema.BoolAttribute{
									Description: "Allow any query parameters and fragment in the resource.",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
								"valid_uri": schema.StringAttribute{
									Description: "URI of a valid resource.",
									Required:    true,
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

func (r *redirectValidationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Compare to version 12.1 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1210)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast121 := compare >= 0
	var plan redirectValidationModel
	req.Plan.Get(ctx, &plan)
	// If redirect_validation_local_settings.uri_allow_list is set prior to PF version 11.3, throw an error
	var diags diag.Diagnostics
	var localSettingsAttrs map[string]attr.Value
	localSettingsModified := false
	if !pfVersionAtLeast121 {
		if internaltypes.IsDefined(plan.RedirectValidationLocalSettings) {
			localSettingsAttrs = plan.RedirectValidationLocalSettings.Attributes()
			uriAllowList := localSettingsAttrs["uri_allow_list"]
			if internaltypes.IsDefined(uriAllowList) {
				version.AddUnsupportedAttributeError("redirect_validation_local_settings.uri_allow_list",
					r.providerConfig.ProductVersion, version.PingFederate1210, &resp.Diagnostics)
			} else if uriAllowList.IsUnknown() {
				localSettingsAttrs["uri_allow_list"] = types.ListNull(types.ObjectType{AttrTypes: uriAllowListAttrTypes})
				localSettingsModified = true
			}
		}
	} else if internaltypes.IsDefined(plan.RedirectValidationLocalSettings) {
		localSettingsAttrs = plan.RedirectValidationLocalSettings.Attributes()
		if localSettingsAttrs["uri_allow_list"].IsUnknown() {
			// Default to empty list
			localSettingsAttrs["uri_allow_list"], diags = types.ListValue(types.ObjectType{AttrTypes: uriAllowListAttrTypes}, nil)
			resp.Diagnostics.Append(diags...)
			localSettingsModified = true
		}
	}

	if localSettingsModified {
		plan.RedirectValidationLocalSettings, diags = types.ObjectValue(redirectValidationLocalSettingsAttrTypes, localSettingsAttrs)
		resp.Diagnostics.Append(diags...)
		resp.Plan.Set(ctx, &plan)
	}
}

func addOptionalRedirectValidationFields(ctx context.Context, addRequest *client.RedirectValidationSettings, plan redirectValidationModel) error {
	if internaltypes.IsDefined(plan.RedirectValidationLocalSettings) {
		addRequest.RedirectValidationLocalSettings = client.NewRedirectValidationLocalSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.RedirectValidationLocalSettings, true)), addRequest.RedirectValidationLocalSettings)
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

func (r *redirectValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan redirectValidationModel

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

	apiCreateRedirectValidation := r.apiClient.RedirectValidationAPI.UpdateRedirectValidationSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateRedirectValidation = apiCreateRedirectValidation.Body(*createRedirectValidation)
	redirectValidationResponse, httpResp, err := r.apiClient.RedirectValidationAPI.UpdateRedirectValidationSettingsExecute(apiCreateRedirectValidation)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Redirect Validation", err, httpResp)
		return
	}

	// Read the response into the state
	var state redirectValidationModel
	diags = readRedirectValidationResponse(ctx, redirectValidationResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *redirectValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state redirectValidationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadRedirectValidation, httpResp, err := r.apiClient.RedirectValidationAPI.GetRedirectValidationSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Redirect Validation", httpResp)
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
	var plan redirectValidationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRedirectValidation := r.apiClient.RedirectValidationAPI.UpdateRedirectValidationSettings(config.AuthContext(ctx, r.providerConfig))
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
	var state redirectValidationModel
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
