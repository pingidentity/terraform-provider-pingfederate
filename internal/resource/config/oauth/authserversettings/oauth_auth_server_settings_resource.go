package oauthauthserversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAuthServerSettingsResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthServerSettingsResource{}
	_ resource.ResourceWithImportState = &oauthAuthServerSettingsResource{}
)

// OauthAuthServerSettingsResource is a helper function to simplify the provider implementation.
func OauthAuthServerSettingsResource() resource.Resource {
	return &oauthAuthServerSettingsResource{
		impl: oauthServerSettingsResource{},
	}
}

// oauthAuthServerSettingsResource is the resource implementation.
type oauthAuthServerSettingsResource struct {
	impl oauthServerSettingsResource
}

// GetSchema defines the schema for the resource.
func (r *oauthAuthServerSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_oauth_auth_server_settings` resource is deprecated. Use the `pingfederate_oauth_server_settings` resource instead."
}

// Metadata returns the resource type name.
func (r *oauthAuthServerSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings"
}

func (r *oauthAuthServerSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	r.impl.ValidateConfig(ctx, req, resp)
}

func (r *oauthAuthServerSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	r.impl.ModifyPlan(ctx, req, resp)
}

func (r *oauthAuthServerSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *oauthAuthServerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *oauthAuthServerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAuthServerSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthAuthServerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *oauthAuthServerSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
