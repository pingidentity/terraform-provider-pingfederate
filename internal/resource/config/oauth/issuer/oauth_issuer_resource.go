package oauthissuer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthIssuerResource{}
	_ resource.ResourceWithConfigure   = &oauthIssuerResource{}
	_ resource.ResourceWithImportState = &oauthIssuerResource{}
)

// OauthIssuerResource is a helper function to simplify the provider implementation.
func OauthIssuerResource() resource.Resource {
	return &oauthIssuerResource{}
}

// oauthIssuerResource is the resource implementation.
type oauthIssuerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthIssuerResourceModel struct {
	Id          types.String `tfsdk:"id"`
	IssuerId    types.String `tfsdk:"issuer_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Host        types.String `tfsdk:"host"`
	Path        types.String `tfsdk:"path"`
}

// GetSchema defines the schema for the resource.
func (r *oauthIssuerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an OAuth Issuer.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"host": schema.StringAttribute{
				Required: true,
			},
			"path": schema.StringAttribute{
				Required: true,
			},
		},
	}

	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"issuer_id",
		true,
		"The persistent, unique ID for the virtual issuer. It can be any combination of [a-zA-Z0-9._-].")
	resp.Schema = schema
}
func addOptionalOauthIssuerFields(ctx context.Context, addRequest *client.Issuer, plan oauthIssuerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.IssuerId) {
		addRequest.Id = plan.IssuerId.ValueStringPointer()
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
func (r *oauthIssuerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_issuer"
}

func (r *oauthIssuerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthIssuerResponse(ctx context.Context, r *client.Issuer, state *oauthIssuerResourceModel, expectedValues *oauthIssuerResourceModel) {
	// Why doesn't PF guarantee the id will be returned on this resource?
	state.Id = types.StringPointerValue(r.Id)
	state.IssuerId = types.StringPointerValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringValue(*r.Description)
	state.Host = types.StringValue(r.Host)
	state.Path = types.StringValue(*r.Path)
}

func (r *oauthIssuerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthIssuerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	oauthIssuer := client.NewIssuer(plan.Name.ValueString(), plan.Host.ValueString())
	err := addOptionalOauthIssuerFields(ctx, oauthIssuer, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for an OAuth Issuer", err.Error())
		return
	}

	apiCreateOauthIssuer := r.apiClient.OauthIssuersAPI.AddOauthIssuer(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthIssuer = apiCreateOauthIssuer.Body(*oauthIssuer)
	oauthIssuerResponse, httpResp, err := r.apiClient.OauthIssuersAPI.AddOauthIssuerExecute(apiCreateOauthIssuer)

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating an OAuth Issuer", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthIssuerResourceModel

	readOauthIssuerResponse(ctx, oauthIssuerResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthIssuerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthIssuerResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthIssuer, httpResp, err := r.apiClient.OauthIssuersAPI.GetOauthIssuerById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.IssuerId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Issuer", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Issuer", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readOauthIssuerResponse(ctx, apiReadOauthIssuer, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthIssuerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan oauthIssuerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthIssuerResourceModel
	req.State.Get(ctx, &state)
	updateOauthIssuer := r.apiClient.OauthIssuersAPI.UpdateOauthIssuer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.IssuerId.ValueString())
	createUpdateRequest := client.NewIssuer(plan.Name.ValueString(), plan.Host.ValueString())
	err := addOptionalOauthIssuerFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to update request for an OAuth Issuer", err.Error())
		return
	}

	updateOauthIssuer = updateOauthIssuer.Body(*createUpdateRequest)
	updateOauthIssuerResponse, httpResp, err := r.apiClient.OauthIssuersAPI.UpdateOauthIssuerExecute(updateOauthIssuer)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating an OAuth Issuer", err, httpResp)
		return
	}

	// Read the response
	readOauthIssuerResponse(ctx, updateOauthIssuerResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *oauthIssuerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthIssuerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthIssuersAPI.DeleteOauthIssuer(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.IssuerId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Issuer", err, httpResp)
		return
	}
}

func (r *oauthIssuerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("issuer_id"), req, resp)
}
