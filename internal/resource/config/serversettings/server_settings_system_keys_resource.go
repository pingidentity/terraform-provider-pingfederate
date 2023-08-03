package serversettings

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverSettingsSystemKeysResource{}
	_ resource.ResourceWithConfigure   = &serverSettingsSystemKeysResource{}
	_ resource.ResourceWithImportState = &serverSettingsSystemKeysResource{}
)

// ServerSettingsSystemKeysResource is a helper function to simplify the provider implementation.
func ServerSettingsSystemKeysResource() resource.Resource {
	return &serverSettingsSystemKeysResource{}
}

// serverSettingsSystemKeysResource is the resource implementation.
type serverSettingsSystemKeysResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsSystemKeysResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Current  types.Object `tfsdk:"current"`
	Previous types.Object `tfsdk:"previous"`
	Pending  types.Object `tfsdk:"pending"`
}

// GetSchema defines the schema for the resource.
func (r *serverSettingsSystemKeysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	serverSettingsSystemKeysResourceSchema(ctx, req, resp, false)
}

func serverSettingsSystemKeysResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	resp.Schema = schema.Schema{
		Description: "Manages a Server Settings SystemKeys.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource",
				Computed:    true,
				Optional:    false,
				Required:    false,
			},
			"current": schema.SingleNestedAttribute{
				Description: "Current SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
				},
			},
			"previous": schema.SingleNestedAttribute{
				Description: "Previous SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
				},
			},
			"pending": schema.SingleNestedAttribute{
				Description: "Pending SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
						Required:    false,
					},
				},
			},
		},
	}
}

func addServerSettingsSystemKeysFields(ctx context.Context, addRequest *client.SystemKeys, plan serverSettingsSystemKeysResourceModel) {

	if internaltypes.IsDefined(plan.Current) {
		currentAttrs := plan.Current.Attributes()
		encryptedKeyDataAttrcurrent := currentAttrs["encrypted_key_data"].(types.String)
		if internaltypes.IsNonEmptyString(encryptedKeyDataAttrcurrent) {
			addRequest.Current = *client.NewSystemKey()
			currentEncryptedKeyData := encryptedKeyDataAttrcurrent.ValueString()
			addRequest.Current.EncryptedKeyData = &currentEncryptedKeyData
		}
	}
	if internaltypes.IsDefined(plan.Previous) {
		previousAttrs := plan.Previous.Attributes()
		encryptedKeyDataAttrPrevious := previousAttrs["encrypted_key_data"].(types.String)
		if internaltypes.IsNonEmptyString(encryptedKeyDataAttrPrevious) {
			addRequest.Previous = client.NewSystemKey()
			previousEncryptedKeyData := encryptedKeyDataAttrPrevious.ValueString()
			addRequest.Previous.EncryptedKeyData = &previousEncryptedKeyData
		}
	}
	if internaltypes.IsDefined(plan.Pending) {
		pendingAttrs := plan.Pending.Attributes()
		encryptedKeyDataAttrPending := pendingAttrs["encrypted_key_data"].(types.String)
		if internaltypes.IsNonEmptyString(encryptedKeyDataAttrPending) {
			addRequest.Pending = *client.NewSystemKey()
			pendingEncryptedKeyData := encryptedKeyDataAttrPending.ValueString()
			addRequest.Pending.EncryptedKeyData = &pendingEncryptedKeyData
		}
	}
}

// Metadata returns the resource type name.
func (r *serverSettingsSystemKeysResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_system_keys"
}

func (r *serverSettingsSystemKeysResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsSystemKeysResponse(ctx context.Context, r *client.SystemKeys, state *serverSettingsSystemKeysResourceModel) {
	state.Id = types.StringValue("id")
	currentAttrTypes := map[string]attr.Type{
		"creation_date":      basetypes.StringType{},
		"encrypted_key_data": basetypes.StringType{},
		"key_data":           basetypes.StringType{},
	}
	currentAttrs := r.GetCurrent()
	currentAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(currentAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(currentAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(currentAttrs.GetKeyData()),
	}
	currentAttrsObjVal := internaltypes.MaptoObjValue(currentAttrTypes, currentAttrVals, diag.Diagnostics{})

	previousAttrTypes := map[string]attr.Type{
		"creation_date":      basetypes.StringType{},
		"encrypted_key_data": basetypes.StringType{},
		"key_data":           basetypes.StringType{},
	}
	previousAttrs := r.GetPrevious()

	previousAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(previousAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(previousAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(previousAttrs.GetKeyData()),
	}
	previousAttrsObjVal := internaltypes.MaptoObjValue(previousAttrTypes, previousAttrVals, diag.Diagnostics{})
	pendingAttrTypes := map[string]attr.Type{
		"creation_date":      basetypes.StringType{},
		"encrypted_key_data": basetypes.StringType{},
		"key_data":           basetypes.StringType{},
	}
	pendingAttrs := r.GetPending()
	pendingAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(pendingAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(pendingAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(pendingAttrs.GetKeyData()),
	}
	pendingAttrsObjVal := internaltypes.MaptoObjValue(pendingAttrTypes, pendingAttrVals, diag.Diagnostics{})

	state.Current = currentAttrsObjVal
	state.Pending = pendingAttrsObjVal
	state.Previous = previousAttrsObjVal
}

func (r *serverSettingsSystemKeysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverSettingsSystemKeysResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createServerSettingsSystemKeys := client.NewSystemKeysWithDefaults()
	addServerSettingsSystemKeysFields(ctx, createServerSettingsSystemKeys, plan)
	requestJson, err := createServerSettingsSystemKeys.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateServerSettingsSystemKeys := r.apiClient.ServerSettingsApi.UpdateSystemKeys(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsSystemKeys = apiCreateServerSettingsSystemKeys.Body(*createServerSettingsSystemKeys)
	serverSettingsSystemKeysResponse, httpResp, err := r.apiClient.ServerSettingsApi.UpdateSystemKeysExecute(apiCreateServerSettingsSystemKeys)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the ServerSettingsSystemKeys", err, httpResp)
		return
	}
	responseJson, err := serverSettingsSystemKeysResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state serverSettingsSystemKeysResourceModel

	readServerSettingsSystemKeysResponse(ctx, serverSettingsSystemKeysResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serverSettingsSystemKeysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readServerSettingsSystemKeys(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readServerSettingsSystemKeys(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state serverSettingsSystemKeysResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsSystemKeys, httpResp, err := apiClient.ServerSettingsApi.GetSystemKeys(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a ServerSettingsSystemKeys", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadServerSettingsSystemKeys.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readServerSettingsSystemKeysResponse(ctx, apiReadServerSettingsSystemKeys, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsSystemKeysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateServerSettingsSystemKeys(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateServerSettingsSystemKeys(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan serverSettingsSystemKeysResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createServerSettingsSystemKeys := client.NewSystemKeysWithDefaults()
	addServerSettingsSystemKeysFields(ctx, createServerSettingsSystemKeys, plan)
	requestJson, err := createServerSettingsSystemKeys.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateServerSettingsSystemKeys := apiClient.ServerSettingsApi.UpdateSystemKeys(config.ProviderBasicAuthContext(ctx, providerConfig))
	apiCreateServerSettingsSystemKeys = apiCreateServerSettingsSystemKeys.Body(*createServerSettingsSystemKeys)
	serverSettingsSystemKeysResponse, httpResp, err := apiClient.ServerSettingsApi.UpdateSystemKeysExecute(apiCreateServerSettingsSystemKeys)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the ServerSettingsSystemKeys", err, httpResp)
		return
	}
	responseJson, err := serverSettingsSystemKeysResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state serverSettingsSystemKeysResourceModel

	readServerSettingsSystemKeysResponse(ctx, serverSettingsSystemKeysResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsSystemKeysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *serverSettingsSystemKeysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importServerSettingsSystemKeysLocation(ctx, req, resp)
}
func importServerSettingsSystemKeysLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//  id  doesn't matter because it is a singleton resource.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
