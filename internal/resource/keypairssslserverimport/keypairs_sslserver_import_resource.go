package keyPairsSslServerImport

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
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
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
	FileData       types.String `tfsdk:"file_data"`
	Format         types.String `tfsdk:"format"`
	Password       types.String `tfsdk:"password"`
	CryptoProvider types.String `tfsdk:"crypto_provider"`
}

// GetSchema defines the schema for the resource.
func (r *keyPairsSslServerImportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	keyPairsSslServerImportResourceSchema(ctx, req, resp, false)
}

func keyPairsSslServerImportResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a KeyPairsSslServerImport.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
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

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"file_data", "format", "password"})
	}
	resp.Schema = schema
}
func addOptionalKeyPairsSslServerImportFields(ctx context.Context, addRequest *client.KeyPairFile, plan keyPairsSslServerImportResourceModel) error {

	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CryptoProvider) {
		addRequest.CryptoProvider = plan.CryptoProvider.ValueStringPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *keyPairsSslServerImportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_sslserver_import"
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for KeyPairsSslServerImport", err.Error())
		return
	}
	requestJson, err := createKeyPairsSslServerImport.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateKeyPairsSslServerImport := r.apiClient.KeyPairsSslServerApi.ImportSslServerKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateKeyPairsSslServerImport = apiCreateKeyPairsSslServerImport.Body(*createKeyPairsSslServerImport)
	keyPairsSslServerImportResponse, httpResp, err := r.apiClient.KeyPairsSslServerApi.ImportSslServerKeyPairExecute(apiCreateKeyPairsSslServerImport)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the KeyPairsSslServerImport", err, httpResp)
		return
	}
	responseJson, err := keyPairsSslServerImportResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state keyPairsSslServerImportResourceModel
	planFileData := plan.FileData.ValueString()
	planFormat := plan.Format.ValueString()
	planPassword := plan.Password.ValueString()

	readKeyPairsSslServerImportResponse(ctx, keyPairsSslServerImportResponse, &state, &plan, planFileData, planFormat, planPassword)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *keyPairsSslServerImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readKeyPairsSslServerImport(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readKeyPairsSslServerImport(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state keyPairsSslServerImportResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadKeyPairsSslServerImport, httpResp, err := apiClient.KeyPairsSslServerApi.GetSslServerKeyPair(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a KeyPairsSslServerImport", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadKeyPairsSslServerImport.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	stateFileData := state.FileData.ValueString()
	stateFormat := state.Format.ValueString()
	statePassword := state.Password.ValueString()
	readKeyPairsSslServerImportResponse(ctx, apiReadKeyPairsSslServerImport, &state, &state, stateFileData, stateFormat, statePassword)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *keyPairsSslServerImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Error(ctx, "Not sure how you got here, however this resource does not support update functionality.")
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *keyPairsSslServerImportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteKeyPairsSslServerImport(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteKeyPairsSslServerImport(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state keyPairsSslServerImportResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.KeyPairsSslServerApi.DeleteSslServerKeyPair(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a KeyPairsSslServerImport", err, httpResp)
		return
	}

}

func (r *keyPairsSslServerImportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
