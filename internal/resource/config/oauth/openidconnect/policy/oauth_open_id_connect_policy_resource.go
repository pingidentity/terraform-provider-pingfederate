package oauthopenidconnectpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthOpenIdConnectPolicyResource{}
	_ resource.ResourceWithConfigure   = &oauthOpenIdConnectPolicyResource{}
	_ resource.ResourceWithImportState = &oauthOpenIdConnectPolicyResource{}
)

// OauthOpenIdConnectPolicyResource is a helper function to simplify the provider implementation.
func OauthOpenIdConnectPolicyResource() resource.Resource {
	return &oauthOpenIdConnectPolicyResource{
		impl: openidConnectPolicyResource{},
	}
}

// oauthOpenIdConnectPolicyResource is the resource implementation.
type oauthOpenIdConnectPolicyResource struct {
	impl openidConnectPolicyResource
}

// GetSchema defines the schema for the resource.
func (r *oauthOpenIdConnectPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_oauth_open_id_connect_policy` resource is deprecated. Use the `pingfederate_openid_connect_policy` resource instead."
}

// Metadata returns the resource type name.
func (r *oauthOpenIdConnectPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_open_id_connect_policy"
}

func (r *oauthOpenIdConnectPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *oauthOpenIdConnectPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *oauthOpenIdConnectPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthOpenIdConnectPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthOpenIdConnectPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *oauthOpenIdConnectPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
