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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/administrativeaccount"
	authenticationapiapplication "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationapi/application"
	authenticationapisettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationapi/settings"
	authenticationpoliciesfragments "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicies/fragments"
	authenticationpoliciessettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicies/settings"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/authenticationpolicycontract"
	certificate "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/certificate/ca"
	datastore "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/datastore"
	idpadapter "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/adapter"
	idpdefaulturls "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/defaulturls"
	idpspconnection "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/idp/spconnection"
	kerberosrealms "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/kerberos/realms"
	keypairsigningimport "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypair/signing/import"
	keypairsslserverimport "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/keypair/sslserver/import"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/license"
	licenseagreement "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/license/agreement"
	localidentity "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/localidentity/identityprofile"
	oauthaccesstokenmanager "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/accesstokenmanager"
	oauthauthserversettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authserversettings"
	oauthauthserversettingsscopescommonscope "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authserversettings/scopes/commonscope"
	oauthauthserversettingsscopesexclusivescope "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/authserversettings/scopes/exclusivescope"
	oauthcibaserverpolicysettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/cibaserverpolicy/settings"
	oauthclient "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/client"
	oauthissuer "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/issuer"
	oauthopenidconnectpolicy "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/openidconnect/policy"
	oauthtokenexchangegeneratorsettings "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/tokenexchange/generator/settings"
	oauthtokenexchangetokengeneratormapping "github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config/oauth/tokenexchange/tokengeneratormapping"
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
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
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
	AdminApiPath                    types.String `tfsdk:"admin_api_path"`
	Username                        types.String `tfsdk:"username"`
	Password                        types.String `tfsdk:"password"`
	AccessToken                     types.String `tfsdk:"access_token"`
	OAuth                           types.Object `tfsdk:"oauth"`
	InsecureTrustAllTls             types.Bool   `tfsdk:"insecure_trust_all_tls"`
	CACertificatePEMFiles           types.Set    `tfsdk:"ca_certificate_pem_files"`
	XBypassExternalValidationHeader types.Bool   `tfsdk:"x_bypass_external_validation_header"`
	ProductVersion                  types.String `tfsdk:"product_version"`
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
			"admin_api_path": schema.StringAttribute{
				MarkdownDescription: "Path for PingFederate Admin API. Default value can be set with the `PINGFEDERATE_PROVIDER_ADMIN_API_PATH` environment variable. If no value is supplied, the value used will be `/pf-admin-api/v1`.",
				Optional:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token for PingFederate Admin API. Default value can be set with the `PINGFEDERATE_PROVIDER_ACCESS_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("oauth")),
					stringvalidator.ConflictsWith(path.MatchRoot("username")),
					stringvalidator.ConflictsWith(path.MatchRoot("password")),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for PingFederate Admin user. Default value can be set with the `PINGFEDERATE_PROVIDER_USERNAME` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("oauth")),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for PingFederate Admin user. Default value can be set with the `PINGFEDERATE_PROVIDER_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("oauth")),
				},
			},
			"oauth": schema.SingleNestedAttribute{
				MarkdownDescription: "OAuth configuration for requesting access token. Default values can be set with the `PINGFEDERATE_PROVIDER_OAUTH_*` environment variables.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						MarkdownDescription: "OAuth client ID for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID` environment variable.",
						Required:            true,
					},
					"client_secret": schema.StringAttribute{
						MarkdownDescription: "OAuth client secret for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET` environment variable.",
						Required:            true,
					},
					"scopes": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "OAuth scopes for access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_SCOPES` environment variable.",
						Optional:            true,
					},
					"token_url": schema.StringAttribute{
						MarkdownDescription: "OAuth token URL for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL` environment variable.",
						Required:            true,
					},
				},
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
			"product_version": schema.StringAttribute{
				Description: "Version of the PingFederate server being configured. Default value can be set with the `PINGFEDERATE_PROVIDER_PRODUCT_VERSION` environment variable.",
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

	// User must provide a admin api base path to the provider
	var adminApiPath string
	if config.AdminApiPath.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as admin_api_path",
		)
	} else {
		if config.AdminApiPath.IsNull() {
			adminApiPath = os.Getenv("PINGFEDERATE_PROVIDER_ADMIN_API_PATH")
			if adminApiPath == "" {
				adminApiPath = "/pf-admin-api/v1"
			}
		} else {
			adminApiPath = config.AdminApiPath.ValueString()
		}
	}

	// Check if the user has provided a username to the provider
	var username string
	if config.Username.IsNull() {
		username = os.Getenv("PINGFEDERATE_PROVIDER_USERNAME")
	} else {
		username = config.Username.ValueString()
	}

	// Check if the user has provided a password to the provider
	var password string
	if config.Password.IsNull() {
		password = os.Getenv("PINGFEDERATE_PROVIDER_PASSWORD")
	} else {
		password = config.Password.ValueString()
	}

	// Check if the user has provided an access token to the provider
	var accessToken string
	if config.AccessToken.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the PingFederate Server",
			"Cannot use unknown value as access_token",
		)
	} else {
		if config.AccessToken.IsNull() {
			accessToken = os.Getenv("PINGFEDERATE_PROVIDER_ACCESS_TOKEN")
		} else {
			accessToken = config.AccessToken.ValueString()
		}
	}

	// Check if the user has provided an OAuth configuration to the provider
	var clientId string
	var hasOauthConfig bool
	if internaltypes.IsDefined(config.OAuth) {
		hasOauthConfig = true
		configClientId := config.OAuth.Attributes()["client_id"].(types.String)
		configClientSecret := config.OAuth.Attributes()["client_secret"].(types.String)
		configScopes := config.OAuth.Attributes()["scopes"].(types.List)
		configTokenUrl := config.OAuth.Attributes()["token_url"].(types.String)

		// Check if the user has provided a client_id to the provider
		if configClientId.IsNull() {
			clientId = os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID")
		} else {
			clientId = configClientId.ValueString()
		}
		if clientId == "" {
			resp.Diagnostics.AddError(
				"Unable to find client_id",
				"client_id cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID environment variable.",
			)
		}

		// Check if the user has provided a client_secret to the provider
		var clientSecret string
		if configClientSecret.IsNull() {
			clientSecret = os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET")
		} else {
			clientSecret = configClientSecret.ValueString()
		}
		if clientSecret == "" {
			resp.Diagnostics.AddError(
				"Unable to find client_secret",
				"client_secret cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET environment variable.",
			)
		}

		// Check if the user has provided scopes to the provider
		var scopes []string
		if len(configScopes.Elements()) > 0 {
			scopes = strings.Split(os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_SCOPES"), ",")
		} else {
			configScopes.ElementsAs(ctx, &scopes, false)
		}
		if len(scopes) == 0 {
			resp.Diagnostics.AddError(
				"Unable to find scopes",
				"scopes cannot be an empty list. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_OAUTH_SCOPES environment variable.",
			)
		}

		// Check if the user has provided a token_url to the provider
		var tokenUrl string
		if configTokenUrl.IsNull() {
			tokenUrl = os.Getenv("PINGFEDERATE_PROVIDER_TOKEN_URL")
		} else {
			tokenUrl = configTokenUrl.ValueString()
		}
		if tokenUrl == "" {
			resp.Diagnostics.AddError(
				"Unable to find token_url",
				"token_url cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL environment variable.",
			)
		}
	}

	// Validate the configuration
	if username == "" && password == "" && accessToken == "" && !hasOauthConfig {
		resp.Diagnostics.AddError(
			"Unable to find username and password, access_token, or OAuth configuration",
			"username and password, access_token, or oauth configuration cannot be empty. Either set them in the configuration or use the PINGFEDERATE_PROVIDER_* environment variables.",
		)
	} else {
		// If user has not provided an OAuth configuration or access token, they must provide username and password
		if !hasOauthConfig && accessToken == "" {
			if username == "" {
				resp.Diagnostics.AddError(
					"Unable to find username",
					"username cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_USERNAME environment variable.",
				)
			}

			if password == "" {
				resp.Diagnostics.AddError(
					"Unable to find password",
					"password cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_PASSWORD environment variable.",
				)
			}
		}

		// If user has not provided username and password or an OAuth configuration, they must provide an access token
		if username == "" && password == "" && !internaltypes.IsDefined(config.OAuth) {
			if accessToken == "" {
				resp.Diagnostics.AddError(
					"Unable to find access_token",
					"access_token cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_ACCESS_TOKEN environment variable.",
				)
			}
		}

		// If user has not provided username and password or an access token, they must provide an OAuth configuration
		if !internaltypes.IsDefined(config.OAuth) && (accessToken != "" && (username != "" || password != "")) {
			resp.Diagnostics.AddError(
				"Unable to find OAuth configuration",
				"Oauth configuration cannot be empty. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_OAUTH_* environment variables.",
			)
		}
	}

	// User must provide a product version to the provider
	var productVersion string
	var parsedProductVersion version.SupportedVersion
	var err error
	if !config.ProductVersion.IsUnknown() && !config.ProductVersion.IsNull() {
		productVersion = config.ProductVersion.ValueString()
	} else {
		productVersion = os.Getenv("PINGFEDERATE_PROVIDER_PRODUCT_VERSION")
	}

	if productVersion == "" {
		resp.Diagnostics.AddError(
			"Unable to find PingFederate version",
			"product_version cannot be an empty string. Either set it in the configuration or use the PINGFEDERATE_PROVIDER_PRODUCT_VERSION environment variable.",
		)
	} else {
		// Validate the PingFederate version
		parsedProductVersion, err = version.Parse(productVersion)
		if err != nil {
			resp.Diagnostics.AddError("Invalid PingFederate version", err.Error())
		}
	}

	// Optional attributes
	var insecureTrustAllTls bool
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

	// The extra suffix for the user-agent is optional and is not considered a provider parameter.
	// We just use it directly from the environment variable, if set.
	userAgentExtraSuffix := os.Getenv("PINGFEDERATE_TF_APPEND_USER_AGENT")

	// Make the PingFederate config and API client info available during DataSource and Resource
	// type Configure methods.
	var resourceConfig internaltypes.ResourceConfiguration
	providerConfig := internaltypes.ProviderConfiguration{
		HttpsHost:      httpsHost,
		ProductVersion: parsedProductVersion,
	}

	if username != "" {
		providerConfig.Username = &username
	}

	if password != "" {
		providerConfig.Password = &password
	}

	if accessToken != "" {
		providerConfig.AccessToken = &accessToken
	}

	if hasOauthConfig {
		providerConfig.ClientId = config.OAuth.Attributes()["client_id"].(types.String).ValueStringPointer()
		providerConfig.ClientSecret = config.OAuth.Attributes()["client_secret"].(types.String).ValueStringPointer()
		var scopes []string
		config.OAuth.Attributes()["scopes"].(types.List).ElementsAs(ctx, &scopes, false)
		providerConfig.Scopes = scopes
		providerConfig.TokenUrl = config.OAuth.Attributes()["token_url"].(types.String).ValueStringPointer()
	}

	resourceConfig.ProviderConfig = providerConfig
	clientConfig := client.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	clientConfig.DefaultHeader["X-BypassExternalValidation"] = strconv.FormatBool(xBypassExternalValidation)
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + adminApiPath,
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
	resourceConfig.ProviderConfig.Transport = tr
	clientConfig.HTTPClient = httpClient
	userAgentSuffix := fmt.Sprintf("terraform-provider-pingfederate/%s %s", p.version, productVersion)
	if userAgentExtraSuffix != "" {
		userAgentSuffix += fmt.Sprintf(" %s", userAgentExtraSuffix)
	}
	clientConfig.UserAgentSuffix = pointers.String(userAgentSuffix)
	resourceConfig.ApiClient = client.NewAPIClient(clientConfig)
	resp.ResourceData = resourceConfig
	resp.DataSourceData = resourceConfig

	tflog.Info(ctx, "Configured PingFederate client", map[string]interface{}{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pingfederateProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		administrativeaccount.AdministrativeAccountDataSource,
		authenticationapiapplication.AuthenticationApiApplicationDataSource,
		authenticationapisettings.AuthenticationApiSettingsDataSource,
		authenticationpoliciesfragments.AuthenticationPoliciesFragmentDataSource,
		authenticationpoliciessettings.AuthenticationPoliciesSettingsDataSource,
		authenticationpolicycontract.AuthenticationPolicyContractDataSource,
		certificate.CertificateDataSource,
		datastore.DataStoreDataSource,
		idpadapter.IdpAdapterDataSource,
		idpdefaulturls.IdpDefaultUrlsDataSource,
		idpspconnection.IdpSpConnectionDataSource,
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
		oauthtokenexchangetokengeneratormapping.OauthTokenExchangeTokenGeneratorMappingDataSource,
		oauthopenidconnectpolicy.OauthOpenIdConnectPolicyDataSource,
		passwordcredentialvalidator.PasswordCredentialValidatorDataSource,
		protocolmetadatalifetimesettings.ProtocolMetadataLifetimeSettingsDataSource,
		redirectvalidation.RedirectValidationDataSource,
		serversettings.ServerSettingsDataSource,
		serversettingsgeneralsettings.ServerSettingsGeneralSettingsDataSource,
		serversettingslogsettings.ServerSettingsLogSettingsDataSource,
		serversettingssystemkeys.ServerSettingsSystemKeysDataSource,
		sessionapplicationsessionpolicy.SessionApplicationSessionPolicyDataSource,
		sessionauthenticationsessionpoliciesglobal.SessionAuthenticationSessionPoliciesGlobalDataSource,
		sessionsettings.SessionSettingsDataSource,
		spauthenticationpolicycontractmapping.SpAuthenticationPolicyContractMappingDataSource,
		tokenprocessortotokengeneratormapping.TokenProcessorToTokenGeneratorMappingDataSource,
		virtualhostnames.VirtualHostNamesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *pingfederateProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		administrativeaccount.AdministrativeAccountResource,
		authenticationapiapplication.AuthenticationApiApplicationResource,
		authenticationapisettings.AuthenticationApiSettingsResource,
		authenticationpoliciesfragments.AuthenticationPoliciesFragmentResource,
		authenticationpoliciessettings.AuthenticationPoliciesSettingsResource,
		authenticationpolicycontract.AuthenticationPolicyContractResource,
		certificate.CertificateCAResource,
		idpadapter.IdpAdapterResource,
		idpdefaulturls.IdpDefaultUrlsResource,
		idpspconnection.IdpSpConnectionResource,
		kerberosrealms.KerberosRealmsResource,
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
		oauthcibaserverpolicysettings.OauthCibaServerPolicySettingsResource,
		oauthclient.OauthClientResource,
		oauthissuer.OauthIssuerResource,
		oauthopenidconnectpolicy.OauthOpenIdConnectPolicyResource,
		oauthtokenexchangegeneratorsettings.OauthTokenExchangeGeneratorSettingsResource,
		oauthtokenexchangetokengeneratormapping.OauthTokenExchangeTokenGeneratorMappingResource,
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
