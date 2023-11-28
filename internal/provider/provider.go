package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/administrativeaccount"
	authenticationapisettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationapi/settings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicycontract"
	certificate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/certificate/ca"
	datastore "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/datastore"
	idpadapter "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/adapter"
	idpdefaulturls "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/defaulturls"
	keypairsigningimport "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypair/signing/import"
	keypairsslserverimport "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypair/sslserver/import"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/license"
	licenseagreement "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/license/agreement"
	localidentity "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/localidentity/identityprofile"
	oauthaccesstokenmanager "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/accesstokenmanager"
	oauthauthserversettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authserversettings"
	oauthauthserversettingsscopescommonscope "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authserversettings/scopes/commonscope"
	oauthauthserversettingsscopesexclusivescope "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authserversettings/scopes/exclusivescope"
	oauthclient "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/client"
	oauthissuer "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/issuer"
	oauthopenidconnectpolicy "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/openidconnect/policy"
	oauthtokenexchangetokengeneratormappings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/tokenexchange/tokengeneratormappings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/passwordcredentialvalidator"
	protocolmetadatalifetimesettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/protocolmetadata/lifetimesettings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/redirectvalidation"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings"
	serversettingsgeneralsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/generalsettings"
	serversettingslogsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/logsettings"
	serversettingssystemkeys "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/serversettings/systemkeys"
	sessionapplicationsessionpolicy "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session/applicationsessionpolicy"
	sessionauthenticationsessionpoliciesglobal "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session/authenticationsessionpolicies/global"
	sessionsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/session/settings"
	spauthenticationpolicycontractmapping "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/sp/authenticationpolicycontractmapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/tokenprocessortotokengeneratormapping"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/virtualhostnames"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfacesß
var (
	_ provider.Provider = &pingfederateProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func NewFactory(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pingfederateProvider{
			version: version,
		}
	}
}

// NewTestProvider is a helper function to simplify testing implementation.
func NewTestProvider() provider.Provider {
	return NewFactory("test")()
}

// PingFederate ProviderModel maps provider schema data to a Go type.
type pingfederateProviderModel struct {
	HttpsHost                       types.String `tfsdk:"https_host"`
	Username                        types.String `tfsdk:"username"`
	Password                        types.String `tfsdk:"password"`
	InsecureTrustAllTls             types.Bool   `tfsdk:"insecure_trust_all_tls"`
	CACertificatePEMFiles           types.Set    `tfsdk:"ca_certificate_pem_files"`
	XBypassExternalValidationHeader types.Bool   `tfsdk:"x_bypass_external_validation_header"`
}

// pingfederateProvider is the provider implementation.
type pingfederateProvider struct {
	version string
}

// Metadata returns the provider type name.
func (p *pingfederateProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingfederate"
}

// GetSchema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *pingfederateProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"https_host": schema.StringAttribute{
				MarkdownDescription: "URI for PingFederate HTTPS port. Default value can be set with the `PINGFEDERATE_PROVIDER_HTTPS_HOST` environment variable.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for PingFederate Admin user. Default value can be set with the `PINGFEDERATE_PROVIDER_USERNAME` environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for PingFederate Admin user. Default value can be set with the `PINGFEDERATE_PROVIDER_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"insecure_trust_all_tls": schema.BoolAttribute{
				Description: "Set to true to trust any certificate when connecting to the PingFederate server. This is insecure and should not be enabled outside of testing. Default value can be set with the `PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS` environment variable.",
				Optional:    true,
			},
			"ca_certificate_pem_files": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Paths to files containing PEM-encoded certificates to be trusted as root CAs when connecting to the PingFederate server over HTTPS. If not set, the host's root CA set will be used. Default value can be set with the `PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES` environment variable, using commas to delimit multiple PEM files if necessary.",
				Optional:    true,
			},
			"x_bypass_external_validation_header": schema.BoolAttribute{
				Description: "Header value in request for PingFederate. The connection test will be bypassed when set to true. Default value can be set with the `PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER` environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *pingfederateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config pingfederateProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// User must provide a https host to the provider
	var httpsHost string
	if config.HttpsHost.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as https_host",
		)
	} else {
		if config.HttpsHost.IsNull() {
			httpsHost = os.Getenv("PINGFEDERATE_PROVIDER_HTTPS_HOST")
		} else {
			httpsHost = config.HttpsHost.ValueString()
		}
		if httpsHost == "" {
			resp.Diagnostics.AddError(
				"Unable to find https_host",
				"https_host cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_HTTPS_HOST environment variable.",
			)
		}
	}

	// User must provide a username to the provider
	var username string
	if config.Username.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as username",
		)
	} else {
		if config.Username.IsNull() {
			username = os.Getenv("PINGFEDERATE_PROVIDER_USERNAME")
		} else {
			username = config.Username.ValueString()
		}
		if username == "" {
			resp.Diagnostics.AddError(
				"Unable to find username",
				"username cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_USERNAME environment variable.",
			)
		}
	}

	// User must provide a username to the provider
	var password string
	if config.Password.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as password",
		)
	} else {
		if config.Password.IsNull() {
			password = os.Getenv("PINGFEDERATE_PROVIDER_PASSWORD")
		} else {
			password = config.Password.ValueString()
		}
		if password == "" {
			resp.Diagnostics.AddError(
				"Unable to find password",
				"password cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_PASSWORD environment variable.",
			)
		}
	}

	// Optional attributes
	var insecureTrustAllTls bool
	var err error
	if !config.InsecureTrustAllTls.IsUnknown() && !config.InsecureTrustAllTls.IsNull() {
		insecureTrustAllTls = config.InsecureTrustAllTls.ValueBool()
	} else {
		insecureTrustAllTls, err = strconv.ParseBool(os.Getenv("PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS"))
		if err != nil {
			insecureTrustAllTls = false
			tflog.Info(ctx, "Failed to parse boolean from 'PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS' environment variable, defaulting 'insecure_trust_all_tls' to false")
		}
	}

	var caCertPemFiles []string
	if !config.CACertificatePEMFiles.IsUnknown() && !config.CACertificatePEMFiles.IsNull() {
		config.CACertificatePEMFiles.ElementsAs(ctx, &caCertPemFiles, false)
	} else {
		pemFilesEnvVar := os.Getenv("PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES")
		if len(pemFilesEnvVar) == 0 {
			tflog.Info(ctx, "Did not find any certificate paths specified via the 'PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES' environment variable, using the host's root CA set")
		} else {
			caCertPemFiles = strings.Split(pemFilesEnvVar, ",")
		}
	}

	var caCertPool *x509.CertPool
	if len(caCertPemFiles) == 0 {
		tflog.Info(ctx, "No CA certs specified, using the host's root CA set")
		caCertPool = nil
	} else {
		caCertPool = x509.NewCertPool()
		for _, pemFilename := range caCertPemFiles {
			// Load CA cert
			pemFilename := filepath.Clean(pemFilename)
			caCert, err := os.ReadFile(pemFilename)
			if err != nil {
				resp.Diagnostics.AddError("Failed to read CA PEM certificate file: "+pemFilename, err.Error())
			}
			tflog.Info(ctx, "Adding CA cert from file: "+pemFilename)
			if !caCertPool.AppendCertsFromPEM(caCert) {
				resp.Diagnostics.AddWarning("Failed to parse certificate", "Failed to parse CA PEM certificate from file: "+pemFilename)
			}
		}
	}

	var xBypassExternalValidation bool
	var xBypassExternalValidationErr error
	if !config.XBypassExternalValidationHeader.IsUnknown() && !config.XBypassExternalValidationHeader.IsNull() {
		xBypassExternalValidation = config.XBypassExternalValidationHeader.ValueBool()
	} else {
		xBypassExternalValidation, xBypassExternalValidationErr = strconv.ParseBool(os.Getenv("PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER"))
		if xBypassExternalValidationErr != nil {
			xBypassExternalValidation = false
			tflog.Info(ctx, "Failed to parse boolean from 'PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER' environment variable, defaulting 'x_bypass_external_validation_header' to false")
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Make the PingFederate config and API client info available during DataSource and Resource
	// type Configure methods.
	var resourceConfig internaltypes.ResourceConfiguration
	providerConfig := internaltypes.ProviderConfiguration{
		HttpsHost: httpsHost,
		Username:  username,
		Password:  password,
	}
	resourceConfig.ProviderConfig = providerConfig
	clientConfig := client.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	clientConfig.DefaultHeader["X-BypassExternalValidation"] = strconv.FormatBool(xBypassExternalValidation)
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + "/pf-admin-api/v1",
		},
	}
	// #nosec G402
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureTrustAllTls,
			RootCAs:            caCertPool,
		},
	}
	httpClient := &http.Client{Transport: tr}
	clientConfig.HTTPClient = httpClient
	clientConfig.UserAgent = fmt.Sprintf("pingtools terraform-provider-pingfederate/%s go", p.version)
	resourceConfig.ApiClient = client.NewAPIClient(clientConfig)
	resp.ResourceData = resourceConfig
	resp.DataSourceData = resourceConfig

	tflog.Info(ctx, "Configured PingFederate client", map[string]interface{}{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pingfederateProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		administrativeaccount.AdministrativeAccountDataSource,
		authenticationapisettings.AuthenticationApiSettingsDataSource,
		authenticationpolicycontract.AuthenticationPolicyContractDataSource,
		certificate.CertificateDataSource,
		idpdefaulturls.IdpDefaultUrlsDataSource,
		keypairsigningimport.KeyPairsSigningImportDataSource,
		keypairsslserverimport.KeyPairsSslServerImportDataSource,
		license.LicenseDataSource,
		licenseagreement.LicenseAgreementDataSource,
		localidentity.LocalIdentityIdentityProfileDataSource,
		oauthaccesstokenmanager.OauthAccessTokenManagerDataSource,
		oauthauthserversettings.OauthAuthServerSettingsDataSource,
		oauthauthserversettingsscopescommonscope.OauthAuthServerSettingsScopesCommonScopeDataSource,
		oauthauthserversettingsscopesexclusivescope.OauthAuthServerSettingsScopesExclusiveScopeDataSource,
		oauthclient.OauthClientDataSource,
		oauthissuer.OauthIssuerDataSource,
		passwordcredentialvalidator.PasswordCredentialValidatorDataSource,
		protocolmetadatalifetimesettings.ProtocolMetadataLifetimeSettingsDataSource,
		redirectvalidation.RedirectValidationDataSource,
		serversettingsgeneralsettings.ServerSettingsGeneralSettingsDataSource,
		serversettingslogsettings.ServerSettingsLogSettingsDataSource,
		serversettingssystemkeys.ServerSettingsSystemKeysDataSource,
		sessionapplicationsessionpolicy.SessionApplicationSessionPolicyDataSource,
		sessionauthenticationsessionpoliciesglobal.SessionAuthenticationSessionPoliciesGlobalDataSource,
		sessionsettings.SessionSettingsDataSource,
		tokenprocessortotokengeneratormapping.TokenProcessorToTokenGeneratorMappingDataSource,
		virtualhostnames.VirtualHostNamesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *pingfederateProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		administrativeaccount.AdministrativeAccountResource,
		authenticationapisettings.AuthenticationApiSettingsResource,
		authenticationpolicycontract.AuthenticationPolicyContractResource,
		certificate.CertificateCAResource,
		idpadapter.IdpAdapterResource,
		idpdefaulturls.IdpDefaultUrlsResource,
		keypairsigningimport.KeyPairsSigningImportResource,
		keypairsslserverimport.KeyPairsSslServerImportResource,
		datastore.DataStoreResource,
		license.LicenseResource,
		licenseagreement.LicenseAgreementResource,
		localidentity.LocalIdentityIdentityProfileResource,
		oauthaccesstokenmanager.OauthAccessTokenManagerResource,
		oauthauthserversettings.OauthAuthServerSettingsResource,
		oauthauthserversettingsscopescommonscope.OauthAuthServerSettingsScopesCommonScopeResource,
		oauthauthserversettingsscopesexclusivescope.OauthAuthServerSettingsScopesExclusiveScopeResource,
		oauthclient.OauthClientResource,
		oauthissuer.OauthIssuerResource,
		oauthopenidconnectpolicy.OauthOpenIdConnectPolicyResource,
		oauthtokenexchangetokengeneratormappings.OauthTokenExchangeTokenGeneratorMappingResource,
		passwordcredentialvalidator.PasswordCredentialValidatorResource,
		protocolmetadatalifetimesettings.ProtocolMetadataLifetimeSettingsResource,
		redirectvalidation.RedirectValidationResource,
		serversettings.ServerSettingsResource,
		serversettingsgeneralsettings.ServerSettingsGeneralSettingsResource,
		serversettingslogsettings.ServerSettingsLogSettingsResource,
		serversettingssystemkeys.ServerSettingsSystemKeysResource,
		sessionapplicationsessionpolicy.SessionApplicationSessionPolicyResource,
		sessionauthenticationsessionpoliciesglobal.SessionAuthenticationSessionPoliciesGlobalResource,
		sessionsettings.SessionSettingsResource,
		spauthenticationpolicycontractmapping.SpAuthenticationPolicyContractMappingResource,
		tokenprocessortotokengeneratormapping.TokenProcessorToTokenGeneratorMappingResource,
		virtualhostnames.VirtualHostNamesResource,
	}
}
