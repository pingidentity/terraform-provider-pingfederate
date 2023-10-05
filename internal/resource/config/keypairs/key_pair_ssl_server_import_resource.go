package keypairs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &keyPairsSslServerImportResource{}
	_ resource.ResourceWithConfigure   = &keyPairsSslServerImportResource{}
	_ resource.ResourceWithImportState = &keyPairsSslServerImportResource{}
)

// KeyPairsSslServerImportResource is a helper function to simplify the provider implementation.
func KeyPairsSslServerImportResource() resource.Resource {
	return &keyPairsSslServerImportResource{}
}

// keyPairsSslServerImportResource is the resource implementation.
type keyPairsSslServerImportResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keyPairsSslServerImportResourceModel struct {
	Id             types.String `tfsdk:"id"`
	CustomId       types.String `tfsdk:"custom_id"`
	FileData       types.String `tfsdk:"file_data"`
	Format         types.String `tfsdk:"format"`
	Password       types.String `tfsdk:"password"`
	CryptoProvider types.String `tfsdk:"crypto_provider"`
}

// GetSchema defines the schema for the resource.
func (r *keyPairsSslServerImportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a KeyPairsSslServerImport.",
		Attributes: map[string]schema.Attribute{
			"custom_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the certificate. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_data": schema.StringAttribute{
				Description: "Base-64 encoded PKCS12 or PEM file data. In the case of PEM, the raw (non-base-64) data is also accepted. In BCFIPS mode, only PEM with PBES2 and AES or Triple DES encryption is accepted and 128-bit salt is required.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"format": schema.StringAttribute{
				Description: "Key pair file format. If specified, this field will control what file format is expected, otherwise the format will be auto-detected. In BCFIPS mode, only PEM is supported. (PKCS12, PEM)",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password for the file. In BCFIPS mode, the password must be at least 14 characters.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. (LOCAL, HSM)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	config.AddCommonSchema(&schema)
	resp.Schema = schema
}

func addOptionalKeyPairsSslServerImportFields(ctx context.Context, addRequest *client.KeyPairFile, plan keyPairsSslServerImportResourceModel) error {

	if internaltypes.IsDefined(plan.CustomId) {
		addRequest.Id = plan.CustomId.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CryptoProvider) {
		addRequest.CryptoProvider = plan.CryptoProvider.ValueStringPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *keyPairsSslServerImportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_pair_ssl_server_import"
}

func (r *keyPairsSslServerImportResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readKeyPairsSslServerImportResponse(ctx context.Context, r *client.KeyPairView, state *keyPairsSslServerImportResourceModel, expectedValues *keyPairsSslServerImportResourceModel, planFileData string, planFormat string, planPassword string) {
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.CustomId = internaltypes.StringTypeOrNil(r.Id, false)
	state.FileData = internaltypes.StringTypeOrNil(&planFileData, false)
	state.Format = internaltypes.StringTypeOrNil(&planFormat, false)
	state.Password = types.StringValue(planPassword)
	state.CryptoProvider = internaltypes.StringTypeOrNil(r.CryptoProvider, false)
}

func (r *keyPairsSslServerImportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan keyPairsSslServerImportResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createKeyPairsSslServerImport := client.NewKeyPairFile(plan.FileData.ValueString(), plan.Password.ValueString())
	err := addOptionalKeyPairsSslServerImportFields(ctx, createKeyPairsSslServerImport, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for KeyPair SSL Server Import", err.Error())
		return
	}
	_, requestErr := createKeyPairsSslServerImport.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of the KeyPair SSL Server Import: %s", requestErr.Error())
	}

	apiCreateKeyPairsSslServerImport := r.apiClient.KeyPairsSslServerAPI.ImportSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateKeyPairsSslServerImport = apiCreateKeyPairsSslServerImport.Body(*createKeyPairsSslServerImport)
	keyPairsSslServerImportResponse, httpResp, err := r.apiClient.KeyPairsSslServerAPI.ImportSslServerKeyPairExecute(apiCreateKeyPairsSslServerImport)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the KeyPair SSL Server Import", err, httpResp)
		return
	}
	_, responseErr := keyPairsSslServerImportResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of the KeyPair SSL Server Import: %s", responseErr.Error())
	}

	// Read the response into the state
	var state keyPairsSslServerImportResourceModel
	planFileData := plan.FileData.ValueString()
	planFormat := plan.Format.ValueString()
	planPassword := plan.Password.ValueString()

	readKeyPairsSslServerImportResponse(ctx, keyPairsSslServerImportResponse, &state, &plan, planFileData, planFormat, planPassword)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *keyPairsSslServerImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state keyPairsSslServerImportResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadKeyPairsSslServerImport, httpResp, err := r.apiClient.KeyPairsSslServerAPI.GetSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the KeyPair SSL Server Import", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the KeyPair SSL Server Import", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadKeyPairsSslServerImport.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of the KeyPair SSL Server Import: %s", responseErr.Error())
	}

	// Read the response into the state
	stateFileData := state.FileData.ValueString()
	stateFormat := state.Format.ValueString()
	statePassword := state.Password.ValueString()
	readKeyPairsSslServerImportResponse(ctx, apiReadKeyPairsSslServerImport, &state, &state, stateFileData, stateFormat, statePassword)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *keyPairsSslServerImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *keyPairsSslServerImportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state keyPairsSslServerImportResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.KeyPairsSslServerAPI.DeleteSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CustomId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a KeyPair SSL Server Import", err, httpResp)
		return
	}
}

func (r *keyPairsSslServerImportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("custom_id"), req, resp)
}
