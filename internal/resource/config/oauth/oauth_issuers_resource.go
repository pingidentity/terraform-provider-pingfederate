package oauth

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &oauthIssuersResource{}
	_ resource.ResourceWithConfigure   = &oauthIssuersResource{}
	_ resource.ResourceWithImportState = &oauthIssuersResource{}
)

// OauthIssuersResource is a helper function to simplify the provider implementation.
func OauthIssuersResource() resource.Resource {
	return &oauthIssuersResource{}
}

// oauthIssuersResource is the resource implementation.
type oauthIssuersResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthIssuersResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Host        types.String `tfsdk:"host"`
	Path        types.String `tfsdk:"path"`
}

// GetSchema defines the schema for the resource.
func (r *oauthIssuersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthIssuersResourceSchema(ctx, req, resp, false)
}

func oauthIssuersResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages an OAuth Issuer.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The persistent, unique ID for the virtual issuer. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "host"})
	}
	resp.Schema = schema
}
func addOptionalOauthIssuersFields(ctx context.Context, addRequest *client.Issuer, plan oauthIssuersResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Path) {
		addRequest.Path = plan.Path.ValueStringPointer()
	}

	return nil
}

// Metadata returns the resource type name.
func (r *oauthIssuersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_issuers"
}

func (r *oauthIssuersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthIssuersResponse(ctx context.Context, r *client.Issuer, state *oauthIssuersResourceModel, expectedValues *oauthIssuersResourceModel) {
	state.Id = types.StringValue(*r.Id)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringValue(*r.Description)
	state.Host = types.StringValue(r.Host)
	state.Path = types.StringValue(*r.Path)
}

func (r *oauthIssuersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthIssuersResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	oauthIssuer := client.NewIssuer(plan.Name.ValueString(), plan.Host.ValueString())
	err := addOptionalOauthIssuersFields(ctx, oauthIssuer, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Issuers", err.Error())
		return
	}
	requestJson, err := oauthIssuer.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateOauthIssuer := r.apiClient.OauthIssuersApi.AddOauthIssuer(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthIssuer = apiCreateOauthIssuer.Body(*oauthIssuer)
	oauthIssuerResponse, httpResp, err := r.apiClient.OauthIssuersApi.AddOauthIssuerExecute(apiCreateOauthIssuer)

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Issuers", err, httpResp)
		return
	}
	responseJson, err := oauthIssuerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthIssuersResourceModel

	readOauthIssuersResponse(ctx, oauthIssuerResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *oauthIssuersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOauthIssuers(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readOauthIssuers(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state oauthIssuersResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthIssuer, httpResp, err := apiClient.OauthIssuersApi.GetOauthIssuerById(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an OAuth Issuers", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadOauthIssuer.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthIssuersResponse(ctx, apiReadOauthIssuer, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthIssuersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOauthIssuers(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateOauthIssuers(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan oauthIssuersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthIssuersResourceModel
	req.State.Get(ctx, &state)
	updateOauthIssuer := apiClient.OauthIssuersApi.UpdateOauthIssuer(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewIssuer(plan.Name.ValueString(), plan.Host.ValueString())
	err := addOptionalOauthIssuersFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to update request for OAuth Issuers", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateOauthIssuer = updateOauthIssuer.Body(*createUpdateRequest)
	updateOauthIssuerResponse, httpResp, err := apiClient.OauthIssuersApi.UpdateOauthIssuerExecute(updateOauthIssuer)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Issuers", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateOauthIssuerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readOauthIssuersResponse(ctx, updateOauthIssuerResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *oauthIssuersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteOauthIssuers(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteOauthIssuers(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state oauthIssuersResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.OauthIssuersApi.DeleteOauthIssuer(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Issuers", err, httpResp)
		return
	}

}

func (r *oauthIssuersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOauthIssuersLocation(ctx, req, resp)
}
func importOauthIssuersLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
