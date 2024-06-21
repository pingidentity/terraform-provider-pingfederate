// Code generated by ping-terraform-plugin-framework-generator

package kerberosrealmssettings

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &kerberosRealmSettingsResource{}
	_ resource.ResourceWithConfigure   = &kerberosRealmSettingsResource{}
	_ resource.ResourceWithImportState = &kerberosRealmSettingsResource{}
)

func KerberosRealmSettingsResource() resource.Resource {
	return &kerberosRealmSettingsResource{}
}

type kerberosRealmSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *kerberosRealmSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kerberos_realm_settings"
}

func (r *kerberosRealmSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type kerberosRealmSettingsResourceModel struct {
	DebugLogOutput            types.Bool  `tfsdk:"debug_log_output"`
	ForceTcp                  types.Bool  `tfsdk:"force_tcp"`
	KdcRetries                types.Int64 `tfsdk:"kdc_retries"`
	KdcTimeout                types.Int64 `tfsdk:"kdc_timeout"`
	KeySetRetentionPeriodMins types.Int64 `tfsdk:"key_set_retention_period_mins"`
}

func (r *kerberosRealmSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"debug_log_output": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Reference to the default logging.",
			},
			"force_tcp": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Reference to the default security.",
			},
			"kdc_retries": schema.Int64Attribute{
				Required:    true,
				Description: "Reference to the default Key Distribution Center Retries.",
			},
			"kdc_timeout": schema.Int64Attribute{
				Required:    true,
				Description: "Reference to the default Key Distribution Center Timeout (in seconds).",
			},
			"key_set_retention_period_mins": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(610),
				Description: "The key set retention period in minutes. When 'retainPreviousKeysOnPasswordChange' is set to true for a realm, this setting determines how long keys will be retained after a password change occurs. If this field is omitted in a PUT request, the default of 610 minutes is applied.",
			},
		},
	}
}

func (model *kerberosRealmSettingsResourceModel) buildClientStruct() *client.KerberosRealmsSettings {
	result := &client.KerberosRealmsSettings{}
	// debug_log_output
	result.DebugLogOutput = model.DebugLogOutput.ValueBoolPointer()
	// force_tcp
	result.ForceTcp = model.ForceTcp.ValueBoolPointer()
	// kdc_retries
	result.KdcRetries = strconv.FormatInt(model.KdcRetries.ValueInt64(), 10)
	// kdc_timeout
	result.KdcTimeout = strconv.FormatInt(model.KdcTimeout.ValueInt64(), 10)
	// key_set_retention_period_mins
	result.KeySetRetentionPeriodMins = model.KeySetRetentionPeriodMins.ValueInt64Pointer()
	return result
}

func (r *kerberosRealmSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data kerberosRealmSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.KerberosRealmsAPI.UpdateKerberosRealmSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.KerberosRealmsAPI.UpdateKerberosRealmSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the kerberosRealmSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kerberosRealmSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data kerberosRealmSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.KerberosRealmsAPI.GetKerberosRealmSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the kerberosRealmSettings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the kerberosRealmSettings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kerberosRealmSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data kerberosRealmSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.KerberosRealmsAPI.UpdateKerberosRealmSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.KerberosRealmsAPI.UpdateKerberosRealmSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the kerberosRealmSettings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kerberosRealmSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
}

func (r *kerberosRealmSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState kerberosRealmSettingsResourceModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
