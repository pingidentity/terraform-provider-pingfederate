// Copyright © 2026 Ping Identity Corporation

package authenticationapiapplication

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationApiApplicationResource{}
	_ resource.ResourceWithConfigure   = &authenticationApiApplicationResource{}
	_ resource.ResourceWithImportState = &authenticationApiApplicationResource{}

	emptyStringSet, _ = types.SetValue(types.StringType, nil)
	customId          = "application_id"
)

// AuthenticationApiApplicationResource is a helper function to simplify the provider implementation.
func AuthenticationApiApplicationResource() resource.Resource {
	return &authenticationApiApplicationResource{}
}

// authenticationApiApplicationResource is the resource implementation.
type authenticationApiApplicationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *authenticationApiApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an Authentication Api Application",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the Authentication API application. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					configvalidators.PingFederateId(),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "The Authentication API Application Name. Name must be unique.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"url": schema.StringAttribute{
				Description: "The Authentication API Application redirect URL.",
				Required:    true,
				Validators: []validator.String{
					configvalidators.ValidUrl(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The Authentication API Application description.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"additional_allowed_origins": schema.SetAttribute{
				Description: "The domain in the redirect URL is always whitelisted. This field contains a list of additional allowed origin URL's for cross-origin resource sharing.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(emptyStringSet),
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						configvalidators.ValidUrl(),
					),
				},
			},
			"client_for_redirectless_mode_ref": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Required:    false,
				Default:     objectdefault.StaticValue(types.ObjectNull(resourcelink.AttrType())),
				Description: "The client this application must use if it invokes the authentication API in redirectless mode. No client may be specified if `restrict_access_to_redirectless_mode` is `false` under `pingfederate_authentication_api_settings`.",
				Attributes:  resourcelink.ToSchema(),
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *authenticationApiApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_api_application"
}

func (r *authenticationApiApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func addOptionalAuthenticationApiApplicationFields(ctx context.Context, addRequest *client.AuthnApiApplication, plan authenticationApiApplicationModel) error {

	addRequest.Description = plan.Description.ValueStringPointer()

	var slice []string
	plan.AdditionalAllowedOrigins.ElementsAs(ctx, &slice, false)
	addRequest.AdditionalAllowedOrigins = slice

	var err error
	addRequest.ClientForRedirectlessModeRef, err = resourcelink.ClientStruct(plan.ClientForRedirectlessModeRef)
	if err != nil {
		return err
	}

	return nil

}

func (r *authenticationApiApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authenticationApiApplicationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAuthenticationApiApplication := client.NewAuthnApiApplication(plan.ApplicationId.ValueString(), plan.Name.ValueString(), plan.Url.ValueString())
	err := addOptionalAuthenticationApiApplicationFields(ctx, createAuthenticationApiApplication, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for an Authentication Api Application: "+err.Error())
		return
	}

	apiCreateAuthenticationApiApplication := r.apiClient.AuthenticationApiAPI.CreateApplication(config.AuthContext(ctx, r.providerConfig))
	apiCreateAuthenticationApiApplication = apiCreateAuthenticationApiApplication.Body(*createAuthenticationApiApplication)
	authenticationApiApplicationResponse, httpResp, err := r.apiClient.AuthenticationApiAPI.CreateApplicationExecute(apiCreateAuthenticationApiApplication)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating an Authentication Api Application", err, httpResp, &customId)
		return
	}

	// Read the response into the state
	var state authenticationApiApplicationModel
	// This is a workaround for the fact that the API does not return the location of oauth client ClientForRedirectlessModeRef. This specifically applies to creates/updates, but location is returned normally on reads, which is why we are running an additional get here.
	if internaltypes.IsDefined(plan.ClientForRedirectlessModeRef) {
		authenticationApiApplicationResponse, httpResp, err = r.apiClient.AuthenticationApiAPI.GetApplication(config.AuthContext(ctx, r.providerConfig), plan.ApplicationId.ValueString()).Execute()
		if err != nil {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Api Application", err, httpResp, &customId)
		}
	}

	diags = readAuthenticationApiApplicationResponse(ctx, authenticationApiApplicationResponse, &state)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationApiApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationApiApplicationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAuthenticationApiApplication, httpResp, err := r.apiClient.AuthenticationApiAPI.GetApplication(config.AuthContext(ctx, r.providerConfig), state.ApplicationId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Authentication API Application", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Api Application", err, httpResp, &customId)
		}
		return
	}

	diags = readAuthenticationApiApplicationResponse(ctx, apiReadAuthenticationApiApplication, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationApiApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan authenticationApiApplicationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateAuthenticationApiApplication := r.apiClient.AuthenticationApiAPI.UpdateApplication(config.AuthContext(ctx, r.providerConfig), plan.ApplicationId.ValueString())
	createUpdateRequest := client.NewAuthnApiApplication(plan.ApplicationId.ValueString(), plan.Name.ValueString(), plan.Url.ValueString())
	err := addOptionalAuthenticationApiApplicationFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for an Authentication Api Application: "+err.Error())
		return
	}

	updateAuthenticationApiApplication = updateAuthenticationApiApplication.Body(*createUpdateRequest)
	updateAuthenticationApiApplicationResponse, httpResp, err := r.apiClient.AuthenticationApiAPI.UpdateApplicationExecute(updateAuthenticationApiApplication)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating an Authentication Api Application", err, httpResp, &customId)
		return
	}

	// Read the response
	var state authenticationApiApplicationModel
	// This is a workaround for the fact that the API does not return the location of oauth client ClientForRedirectlessModeRef. This specifically applies to creates/updates, but location is returned normally on reads, which is why we are running an additional get here.
	if internaltypes.IsDefined(plan.ClientForRedirectlessModeRef) {
		updateAuthenticationApiApplicationResponse, httpResp, err = r.apiClient.AuthenticationApiAPI.GetApplication(config.AuthContext(ctx, r.providerConfig), plan.ApplicationId.ValueString()).Execute()
		if err != nil {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting an Authentication Api Application", err, httpResp, &customId)
		}
	}

	diags = readAuthenticationApiApplicationResponse(ctx, updateAuthenticationApiApplicationResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *authenticationApiApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state authenticationApiApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AuthenticationApiAPI.DeleteApplication(config.AuthContext(ctx, r.providerConfig), state.ApplicationId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting an Authentication Api Application", err, httpResp, &customId)
		return
	}

}

func (r *authenticationApiApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("application_id"), req, resp)
}
