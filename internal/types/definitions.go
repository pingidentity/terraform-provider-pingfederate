// Copyright © 2026 Ping Identity Corporation

package types

import (
	"net/http"

	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Configuration used by the provider and resources
type ProviderConfiguration struct {
	HttpsHost      string
	Transport      *http.Transport
	Username       *string
	Password       *string
	AccessToken    *string
	TokenUrl       *string
	ClientId       *string
	ClientSecret   *string
	Scopes         []string
	ProductVersion version.SupportedVersion
}

// Configuration passed to resources
type ResourceConfiguration struct {
	ProviderConfig ProviderConfiguration
	ApiClient      *client.APIClient
}
