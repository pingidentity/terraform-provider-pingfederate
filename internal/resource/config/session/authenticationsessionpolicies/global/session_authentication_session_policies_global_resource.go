package sessionauthenticationsessionpoliciesglobal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sessionAuthenticationSessionPoliciesGlobalResource{}
	_ resource.ResourceWithConfigure   = &sessionAuthenticationSessionPoliciesGlobalResource{}
	_ resource.ResourceWithImportState = &sessionAuthenticationSessionPoliciesGlobalResource{}
)

// SessionAuthenticationSessionPoliciesGlobalResource is a helper function to simplify the provider implementation.
func SessionAuthenticationSessionPoliciesGlobalResource() resource.Resource {
	return &sessionAuthenticationSessionPoliciesGlobalResource{
		impl: sessionAuthenticationPoliciesGlobalResource{},
	}
}

// sessionAuthenticationSessionPoliciesGlobalResource is the resource implementation.
type sessionAuthenticationSessionPoliciesGlobalResource struct {
	impl sessionAuthenticationPoliciesGlobalResource
}

// GetSchema defines the schema for the resource.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.impl.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "The `pingfederate_session_authentication_session_policies_global` resource is deprecated. Use the `pingfederate_session_authentication_policies_global` resource instead."
}

// Metadata returns the resource type name.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_authentication_session_policies_global"
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.impl.Configure(ctx, req, resp)
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.impl.Create(ctx, req, resp)
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.impl.Read(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.impl.Update(ctx, req, resp)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.impl.Delete(ctx, req, resp)
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.impl.ImportState(ctx, req, resp)
}
