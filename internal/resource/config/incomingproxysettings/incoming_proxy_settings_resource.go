package incomingproxysettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &incomingProxySettingsResource{}
	_ resource.ResourceWithConfigure   = &incomingProxySettingsResource{}
	_ resource.ResourceWithImportState = &incomingProxySettingsResource{}
)

// IncomingProxySettingsResource is a helper function to simplify the provider implementation.
func IncomingProxySettingsResource() resource.Resource {
	return &incomingProxySettingsResource{}
}

// incomingProxySettingsResource is the resource implementation.
type incomingProxySettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type incomingProxySettingsResourceModel struct {
	ForwardedIpAddressHeaderName  types.String `tfsdk:"forwarded_ip_address_header_name"`
	ForwardedIpAddressHeaderIndex types.String `tfsdk:"forwarded_ip_address_header_index"`
	ForwardedHostHeaderName       types.String `tfsdk:"forwarded_host_header_name"`
	ForwardedHostHeaderIndex      types.String `tfsdk:"forwarded_host_header_index"`
	ClientCertSSLHeaderName       types.String `tfsdk:"client_cert_ssl_header_name"`
	ClientCertChainSSLHeaderName  types.String `tfsdk:"client_cert_chain_ssl_header_name"`
	ProxyTerminatesHttpsConns     types.Bool   `tfsdk:"proxy_terminates_https_conns"`
}

// GetSchema defines the schema for the resource.
func (r *incomingProxySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages incoming proxy settings",
		Attributes: map[string]schema.Attribute{
			"forwarded_ip_address_header_name": schema.StringAttribute{
				Description: "Globally specify the header name (for example, X-Forwarded-For) where PingFederate should attempt to retrieve the client IP address in all HTTP requests.",
				Computed:    false,
				Optional:    true,
			},
			// Default value for the index is set in ModifyPlan method if:
			//    ForwardedIpAddressHeaderName is set in HCL AND
			//    ForwardedIpAddressHeaderIndex is not
			"forwarded_ip_address_header_index": schema.StringAttribute{
				Description: "PingFederate combines multiple comma-separated header values into the same order that they are received. Define which IP address you want to use. Default is to use the last address.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("FIRST", "LAST"),
					stringvalidator.AlsoRequires(path.MatchRoot("forwarded_ip_address_header_name")),
				},
			},
			"forwarded_host_header_name": schema.StringAttribute{
				Description: "Globally specify the header name (for example, X-Forwarded-Host) where PingFederate should attempt to retrieve the hostname and port in all HTTP requests.",
				Computed:    false,
				Optional:    true,
			},
			// Default value for the index is set in ModifyPlan method if:
			//    ForwardedHostHeaderName is set in HCL AND
			//    ForwardedHostHeaderIndex is not
			"forwarded_host_header_index": schema.StringAttribute{
				Description: "PingFederate combines multiple comma-separated header values into the same order that they are received. Define which hostname you want to use. Default is to use the last hostname.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("FIRST", "LAST"),
					stringvalidator.AlsoRequires(path.MatchRoot("forwarded_host_header_name")),
				},
			},
			"client_cert_ssl_header_name": schema.StringAttribute{
				Description: "While the proxy server is configured to pass client certificates as HTTP request headers, specify the header name here.",
				Computed:    false,
				Optional:    true,
			},
			"client_cert_chain_ssl_header_name": schema.StringAttribute{
				Description: "While the proxy server is configured to pass client certificates as HTTP request headers, specify the chain header name here.",
				Computed:    false,
				Optional:    true,
			},
			"proxy_terminates_https_conns": schema.BoolAttribute{
				Description: "Allows you to globally specify that connections to the reverse proxy are made over HTTPS even when HTTP is used between the reverse proxy and PingFederate. Default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	resp.Schema = schema
}

func addOptionalIncomingProxySettingsFields(ctx context.Context, addRequest *client.IncomingProxySettings, plan incomingProxySettingsResourceModel) {

	addRequest.ForwardedIpAddressHeaderName = plan.ForwardedIpAddressHeaderName.ValueStringPointer()
	addRequest.ForwardedIpAddressHeaderIndex = plan.ForwardedIpAddressHeaderIndex.ValueStringPointer()
	addRequest.ForwardedHostHeaderName = plan.ForwardedHostHeaderName.ValueStringPointer()
	addRequest.ForwardedHostHeaderIndex = plan.ForwardedHostHeaderIndex.ValueStringPointer()
	addRequest.ClientCertSSLHeaderName = plan.ClientCertSSLHeaderName.ValueStringPointer()
	addRequest.ClientCertChainSSLHeaderName = plan.ClientCertChainSSLHeaderName.ValueStringPointer()
	addRequest.ProxyTerminatesHttpsConns = plan.ProxyTerminatesHttpsConns.ValueBoolPointer()
}

// Metadata returns the resource type name.
func (r *incomingProxySettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incoming_proxy_settings"
}

func (r *incomingProxySettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *incomingProxySettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *incomingProxySettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}

	// PingFederate sets index to "LAST" if the header name is set and the index is not
	// Need these to match the behavior in state
	if internaltypes.IsDefined(plan.ForwardedIpAddressHeaderName) && !internaltypes.IsDefined(plan.ForwardedIpAddressHeaderIndex) {
		plan.ForwardedIpAddressHeaderIndex = types.StringValue("LAST")
	}

	if internaltypes.IsDefined(plan.ForwardedHostHeaderName) && !internaltypes.IsDefined(plan.ForwardedHostHeaderIndex) {
		plan.ForwardedHostHeaderIndex = types.StringValue("LAST")
	}

	// Plan checks against nil values, not empty strings
	if !internaltypes.IsDefined(plan.ForwardedIpAddressHeaderName) && !internaltypes.IsDefined(plan.ForwardedIpAddressHeaderIndex) {
		plan.ForwardedIpAddressHeaderIndex = types.StringNull()
	}

	// Plan checks against nil values, not empty strings
	if !internaltypes.IsDefined(plan.ForwardedHostHeaderName) && !internaltypes.IsDefined(plan.ForwardedHostHeaderIndex) {
		plan.ForwardedHostHeaderIndex = types.StringNull()
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)

}

func readIncomingProxySettingsResponse(ctx context.Context, r *client.IncomingProxySettings, state *incomingProxySettingsResourceModel) {
	state.ForwardedIpAddressHeaderName = types.StringPointerValue(r.ForwardedIpAddressHeaderName)
	state.ForwardedIpAddressHeaderIndex = types.StringPointerValue(r.ForwardedIpAddressHeaderIndex)
	state.ForwardedHostHeaderName = types.StringPointerValue(r.ForwardedHostHeaderName)
	state.ForwardedHostHeaderIndex = types.StringPointerValue(r.ForwardedHostHeaderIndex)
	state.ClientCertSSLHeaderName = types.StringPointerValue(r.ClientCertSSLHeaderName)
	state.ClientCertChainSSLHeaderName = types.StringPointerValue(r.ClientCertChainSSLHeaderName)
	state.ProxyTerminatesHttpsConns = types.BoolPointerValue(r.ProxyTerminatesHttpsConns)

}

func (r *incomingProxySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan incomingProxySettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIncomingProxySettings := client.NewIncomingProxySettings()
	addOptionalIncomingProxySettingsFields(ctx, createIncomingProxySettings, plan)

	apiCreateIncomingProxySettings := r.apiClient.IncomingProxySettingsAPI.UpdateIncomingProxySettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateIncomingProxySettings = apiCreateIncomingProxySettings.Body(*createIncomingProxySettings)
	incomingProxySettingsResponse, httpResp, err := r.apiClient.IncomingProxySettingsAPI.UpdateIncomingProxySettingsExecute(apiCreateIncomingProxySettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the incoming proxy settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state incomingProxySettingsResourceModel

	readIncomingProxySettingsResponse(ctx, incomingProxySettingsResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *incomingProxySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state incomingProxySettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIncomingProxySettings, httpResp, err := r.apiClient.IncomingProxySettingsAPI.GetIncomingProxySettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Incoming Proxy Settings", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the incoming proxy settings", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readIncomingProxySettingsResponse(ctx, apiReadIncomingProxySettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *incomingProxySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan incomingProxySettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateIncomingProxySettings := r.apiClient.IncomingProxySettingsAPI.UpdateIncomingProxySettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewIncomingProxySettings()
	addOptionalIncomingProxySettingsFields(ctx, createUpdateRequest, plan)

	updateIncomingProxySettings = updateIncomingProxySettings.Body(*createUpdateRequest)
	updateIncomingProxySettingsResponse, httpResp, err := r.apiClient.IncomingProxySettingsAPI.UpdateIncomingProxySettingsExecute(updateIncomingProxySettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating incoming proxy settings.", err, httpResp)
		return
	}

	// Read the response
	var state incomingProxySettingsResourceModel
	readIncomingProxySettingsResponse(ctx, updateIncomingProxySettingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *incomingProxySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *incomingProxySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState incomingProxySettingsResourceModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
