package serversettingssystemkeys

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsSystemKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsSystemKeysDataSource{}
)

// ServerSettingsSystemKeysDataSource is a helper function to simplify the provider implementation.
func NewServerSettingsSystemKeysDataSource() datasource.DataSource {
	return &serverSettingsSystemKeysDataSource{}
}

// serverSettingsSystemKeysDataSource is the resource implementation.
type serverSettingsSystemKeysDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsSystemKeysDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	Current  types.Object `tfsdk:"current"`
	Previous types.Object `tfsdk:"previous"`
	Pending  types.Object `tfsdk:"pending"`
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsSystemKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes Server Settings SystemKeys.",
		Attributes: map[string]schema.Attribute{
			"current": schema.SingleNestedAttribute{
				Description: "Current SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    false,
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
						Optional:    false,
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"previous": schema.SingleNestedAttribute{
				Description: "Previous SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Optional:    true,
				Computed:    true,
				Default:     objectdefault.StaticValue(previousKeyDefault),
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
						Required:    false,
						Default:     stringdefault.StaticString(creationTimeDefault),
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("key_data")),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("encrypted_key_data")),
						},
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
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("key_data")),
						},
					},
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("encrypted_key_data")),
						},
					},
				},
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addServerSettingsSystemKeysFields(ctx context.Context, addRequest *client.SystemKeys, plan serverSettingsSystemKeysDataSourceModel) {

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
func (r *serverSettingsSystemKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_system_keys"
}

func (r *serverSettingsSystemKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.DataSourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsSystemKeysResponse(ctx context.Context, r *client.SystemKeys, state *serverSettingsSystemKeysDataSourceModel, existingId *string) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = id.GenerateUUIDToState(existingId)
	currentAttrs := r.GetCurrent()
	currentAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(currentAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(currentAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(currentAttrs.GetKeyData()),
	}
	currentAttrsObjVal := internaltypes.MaptoObjValue(systemKeyAttrTypes, currentAttrVals, &diags)

	previousAttrs := r.GetPrevious()
	previousAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(previousAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(previousAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(previousAttrs.GetKeyData()),
	}
	previousAttrsObjVal := internaltypes.MaptoObjValue(systemKeyAttrTypes, previousAttrVals, &diags)

	pendingAttrs := r.GetPending()
	pendingAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(pendingAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(pendingAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(pendingAttrs.GetKeyData()),
	}
	pendingAttrsObjVal := internaltypes.MaptoObjValue(systemKeyAttrTypes, pendingAttrVals, &diags)

	state.Current = currentAttrsObjVal
	state.Pending = pendingAttrsObjVal
	state.Previous = previousAttrsObjVal
	return diags
}

func (r *serverSettingsSystemKeysDataSource) Create(ctx context.Context, req datasource.CreateRequest, resp *datasource.CreateResponse) {
	var plan serverSettingsSystemKeysDataSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createServerSettingsSystemKeys := client.NewSystemKeysWithDefaults()
	addServerSettingsSystemKeysFields(ctx, createServerSettingsSystemKeys, plan)

	apiCreateServerSettingsSystemKeys := r.apiClient.ServerSettingsAPI.UpdateSystemKeys(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsSystemKeys = apiCreateServerSettingsSystemKeys.Body(*createServerSettingsSystemKeys)
	serverSettingsSystemKeysResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateSystemKeysExecute(apiCreateServerSettingsSystemKeys)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings System Keys", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsSystemKeysDataSourceModel
	diags = readServerSettingsSystemKeysResponse(ctx, serverSettingsSystemKeysResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverSettingsSystemKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsSystemKeysDataSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsSystemKeys, httpResp, err := r.apiClient.ServerSettingsAPI.GetSystemKeys(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings System Keys", err, httpResp)
			resp.State.RemoveDataSource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings System Keys", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readServerSettingsSystemKeysResponse(ctx, apiReadServerSettingsSystemKeys, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverSettingsSystemKeysDataSource) Update(ctx context.Context, req datasource.UpdateRequest, resp *datasource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverSettingsSystemKeysDataSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createServerSettingsSystemKeys := client.NewSystemKeysWithDefaults()
	addServerSettingsSystemKeysFields(ctx, createServerSettingsSystemKeys, plan)

	apiCreateServerSettingsSystemKeys := r.apiClient.ServerSettingsAPI.UpdateSystemKeys(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateServerSettingsSystemKeys = apiCreateServerSettingsSystemKeys.Body(*createServerSettingsSystemKeys)
	serverSettingsSystemKeysResponse, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateSystemKeysExecute(apiCreateServerSettingsSystemKeys)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Server Settings System Keys", err, httpResp)
		return
	}

	// Read the response into the state
	var state serverSettingsSystemKeysDataSourceModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = readServerSettingsSystemKeysResponse(ctx, serverSettingsSystemKeysResponse, &state, id)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *serverSettingsSystemKeysDataSource) Delete(ctx context.Context, req datasource.DeleteRequest, resp *datasource.DeleteResponse) {
}

func (r *serverSettingsSystemKeysDataSource) ImportState(ctx context.Context, req datasource.ImportStateRequest, resp *datasource.ImportStateResponse) {
	//  id  doesn't matter because it is a singleton datasource.
	datasource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
