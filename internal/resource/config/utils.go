package config

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Get BasicAuth context with a username and password
func BasicAuthContext(ctx context.Context, username, password string) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: username,
		Password: password,
	})
}

// Get a BasicAuth context from a ProviderConfiguration
func ProviderBasicAuthContext(ctx context.Context, providerConfig internaltypes.ProviderConfiguration) context.Context {
	return BasicAuthContext(ctx, *providerConfig.Username, *providerConfig.Password)
}

// Get an AccessToken context with an accessToken
func AccessTokenContext(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, client.ContextAccessToken, accessToken)
}

// Get an AccessToken context from a ProviderConfiguration
func ProviderAccessTokenContext(ctx context.Context, providerConfig internaltypes.ProviderConfiguration) context.Context {
	return AccessTokenContext(ctx, *providerConfig.AccessToken)
}

// Get an OAuth context with a tokenUrl, clientId, clientSecret, and scopes
func OAuthContext(ctx context.Context, transport *http.Transport, tokenUrl string, clientId string, clientSecret string, scopes []string) context.Context {
	return context.WithValue(ctx, client.ContextOAuth2, client.OAuthValues{
		Transport:    transport,
		TokenUrl:     tokenUrl,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scopes:       scopes,
	})
}

// Get an OAuth context from a ProviderConfiguration
func ProviderOAuthContext(ctx context.Context, providerConfig internaltypes.ProviderConfiguration) context.Context {
	return OAuthContext(ctx, providerConfig.Transport, *providerConfig.TokenUrl, *providerConfig.ClientId, *providerConfig.ClientSecret, providerConfig.Scopes)
}

func AuthContext(ctx context.Context, providerConfig internaltypes.ProviderConfiguration) context.Context {
	if providerConfig.Username != nil && providerConfig.Password != nil {
		return ProviderBasicAuthContext(ctx, providerConfig)
	} else if providerConfig.ClientId != nil && providerConfig.ClientSecret != nil && providerConfig.TokenUrl != nil {
		return ProviderOAuthContext(ctx, providerConfig)
	} else if providerConfig.AccessToken != nil {
		return ProviderAccessTokenContext(ctx, providerConfig)
	}
	return ctx
}

// Error from PF API
type pingFederateError struct {
	Schemas []string `json:"schemas"`
	Status  string   `json:"status"`
	Detail  string   `json:"detail"`
}

// Report an HTTP error as a warning
func ReportHttpErrorAsWarning(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response) {
	reportHttpResponse(ctx, diagnostics, errorSummary, err, httpResp, true)
}

func reportHttpResponse(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response, isWarning bool) {
	httpErrorPrinted := false
	var internalError error
	if httpResp != nil {
		body, internalError := io.ReadAll(httpResp.Body)
		if internalError == nil {
			tflog.Debug(ctx, "Error HTTP response body: "+string(body))
			var pdError pingFederateError
			internalError = json.Unmarshal(body, &pdError)
			if internalError == nil {
				if isWarning {
					diagnostics.AddWarning(errorSummary, err.Error()+" - Detail: "+pdError.Detail)
				} else {
					diagnostics.AddError(errorSummary, err.Error()+" - Detail: "+pdError.Detail)
				}
				httpErrorPrinted = true
			}
		}
	}
	if !httpErrorPrinted {
		if internalError != nil {
			tflog.Warn(ctx, "Failed to unmarshal HTTP response body: "+internalError.Error())
		}
		if isWarning {
			diagnostics.AddWarning(errorSummary, err.Error())
		} else {
			diagnostics.AddError(errorSummary, err.Error())
		}
	}
}

// Report an HTTP error
func ReportHttpError(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response) {
	httpErrorPrinted := false
	var internalError error
	if httpResp != nil {
		body, internalError := io.ReadAll(httpResp.Body)
		if internalError == nil {
			tflog.Debug(ctx, "Error HTTP response body: "+string(body))
			var paError pingFederateError
			internalError = json.Unmarshal(body, &paError)
			if internalError == nil {
				diagnostics.AddError(errorSummary, err.Error()+" - Detail: "+string(body))
				httpErrorPrinted = true
			}
		}
	}
	if !httpErrorPrinted {
		if internalError != nil {
			tflog.Warn(ctx, "Failed to unmarshal HTTP response body: "+internalError.Error())
		}
		diagnostics.AddError(errorSummary, err.Error())
	}
}
