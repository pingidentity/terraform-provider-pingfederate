package idp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpDefaultUrlsResource{}
	_ resource.ResourceWithConfigure   = &idpDefaultUrlsResource{}
	_ resource.ResourceWithImportState = &idpDefaultUrlsResource{}
)

// IdpDefaultUrlsResource is a helper function to simplify the provider implementation.
func IdpDefaultUrlsResource() resource.Resource {
	return &idpDefaultUrlsResource{}
}

// idpDefaultUrlsResource is the resource implementation.
type idpDefaultUrlsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type idpDefaultUrlsResourceModel struct {
	Id               types.String `tfsdk:"id"`
	ConfirmIdpSlo    types.Bool   `tfsdk:"confirm_idp_slo"`
	IdpSloSuccessUrl types.String `tfsdk:"idp_slo_success_url"`
	IdpErrorMsg      types.String `tfsdk:"idp_error_msg"`
}

// GetSchema defines the schema for the resource.
func (r *idpDefaultUrlsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	idpDefaultUrlsResourceSchema(ctx, req, resp, false)
}

func idpDefaultUrlsResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a IdpDefaultUrls.",
		Attributes: map[string]schema.Attribute{
			"confirm_idp_slo": schema.BoolAttribute{
				Description: "Prompt user to confirm Single Logout (SLO).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"idp_error_msg": schema.StringAttribute{
				Description: "Provide the error text displayed in a user's browser when an SSO operation fails.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"idp_slo_success_url": schema.StringAttribute{
				Description: "Provide the default URL you would like to send the user to when Single Logout has succeeded.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"idp_error_msg"})
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

func addOptionalIdpDefaultUrlsFields(ctx context.Context, addRequest *client.IdpDefaultUrl, plan idpDefaultUrlsResourceModel) error {
	if internaltypes.IsDefined(plan.ConfirmIdpSlo) {
		addRequest.ConfirmIdpSlo = plan.ConfirmIdpSlo.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IdpSloSuccessUrl) {
		addRequest.IdpSloSuccessUrl = plan.IdpSloSuccessUrl.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IdpErrorMsg) {
		addRequest.IdpErrorMsg = plan.IdpErrorMsg.ValueString()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *idpDefaultUrlsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_default_urls"
}

func (r *idpDefaultUrlsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readIdpDefaultUrlsResponse(ctx context.Context, r *client.IdpDefaultUrl, state *idpDefaultUrlsResourceModel, expectedValues *idpDefaultUrlsResourceModel) {
	state.Id = types.StringValue("id")
	state.ConfirmIdpSlo = types.BoolPointerValue(r.ConfirmIdpSlo)
	state.IdpSloSuccessUrl = internaltypes.StringTypeOrNil(r.IdpSloSuccessUrl, false)
	state.IdpErrorMsg = types.StringValue(r.IdpErrorMsg)
}

func (r *idpDefaultUrlsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpDefaultUrlsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIdpDefaultUrls := client.NewIdpDefaultUrl(plan.IdpErrorMsg.ValueString())
	err := addOptionalIdpDefaultUrlsFields(ctx, createIdpDefaultUrls, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpDefaultUrls", err.Error())
		return
	}
	requestJson, err := createIdpDefaultUrls.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateIdpDefaultUrls := r.apiClient.IdpDefaultUrlsApi.UpdateDefaultUrlSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateIdpDefaultUrls = apiCreateIdpDefaultUrls.Body(*createIdpDefaultUrls)
	idpDefaultUrlsResponse, httpResp, err := r.apiClient.IdpDefaultUrlsApi.UpdateDefaultUrlSettingsExecute(apiCreateIdpDefaultUrls)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdpDefaultUrls", err, httpResp)
		return
	}
	responseJson, err := idpDefaultUrlsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state idpDefaultUrlsResourceModel

	readIdpDefaultUrlsResponse(ctx, idpDefaultUrlsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *idpDefaultUrlsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdpDefaultUrls(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIdpDefaultUrls(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state idpDefaultUrlsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpDefaultUrls, httpResp, err := apiClient.IdpDefaultUrlsApi.GetDefaultUrl(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a IdpDefaultUrls", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadIdpDefaultUrls.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIdpDefaultUrlsResponse(ctx, apiReadIdpDefaultUrls, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpDefaultUrlsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdpDefaultUrls(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIdpDefaultUrls(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan idpDefaultUrlsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state idpDefaultUrlsResourceModel
	req.State.Get(ctx, &state)
	updateIdpDefaultUrls := apiClient.IdpDefaultUrlsApi.UpdateDefaultUrlSettings(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewIdpDefaultUrl(plan.IdpErrorMsg.ValueString())
	err := addOptionalIdpDefaultUrlsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpDefaultUrls", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateIdpDefaultUrls = updateIdpDefaultUrls.Body(*createUpdateRequest)
	updateIdpDefaultUrlsResponse, httpResp, err := apiClient.IdpDefaultUrlsApi.UpdateDefaultUrlSettingsExecute(updateIdpDefaultUrls)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdpDefaultUrls", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateIdpDefaultUrlsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readIdpDefaultUrlsResponse(ctx, updateIdpDefaultUrlsResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *idpDefaultUrlsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *idpDefaultUrlsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdpDefaultUrlsLocation(ctx, req, resp)
}
func importIdpDefaultUrlsLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
