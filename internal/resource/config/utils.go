// Copyright © 2026 Ping Identity Corporation

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	// lowerOrDigitToUpper will match a sequence of a lowercase letter or digit followed by an uppercase letter
	lowerOrDigitToUpper = regexp.MustCompile(`([a-z0-9])[A-Z]`)
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
type pingFederateValidationError struct {
	Message          string `json:"message"`
	DeveloperMessage string `json:"developerMessage,omitempty"`
	FieldPath        string `json:"fieldPath"`
	ErrorId          string `json:"errorId"`
}

type pingFederateErrorResponse struct {
	ResultId         string                        `json:"resultId"`
	Message          string                        `json:"message"`
	ValidationErrors []pingFederateValidationError `json:"validationErrors"`
}

// Report a 404 as a warning for resources
func AddResourceNotFoundWarning(ctx context.Context, diagnostics *diag.Diagnostics, resourceType string, httpResp *http.Response) {
	diagnostics.AddWarning("Resource not found", fmt.Sprintf("The requested %s resource configuration cannot be found in the PingFederate service.  If the requested resource is managed in Terraform's state, it may have been removed outside of Terraform.", resourceType))
	if httpResp != nil {
		body, err := io.ReadAll(httpResp.Body)
		if err == nil {
			tflog.Debug(ctx, "Error HTTP response body: "+string(body))
		} else {
			tflog.Warn(ctx, "Failed to read HTTP response body: "+err.Error())
		}
	}
}

func toTerraformIdentifier(pfIdentifier string) string {
	// Insert an underscore between lowercase letter followed by uppercase letter
	insertedUnderscores := lowerOrDigitToUpper.ReplaceAllStringFunc(pfIdentifier, func(s string) string {
		firstRune, size := utf8.DecodeRuneInString(s)
		if firstRune == utf8.RuneError && size <= 1 {
			// The string is empty, return it
			return s
		}

		return fmt.Sprintf("%s_%s", string(firstRune), strings.ToLower(s[size:]))
	})

	// Lowercase the final string
	return strings.ToLower(insertedUnderscores)
}

// Report an HTTP error
func ReportHttpError(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response) {
	ReportHttpErrorCustomId(ctx, diagnostics, errorSummary, err, httpResp, nil)
}

// Report an HTTP error
func ReportHttpErrorCustomId(ctx context.Context, diagnostics *diag.Diagnostics, errorSummary string, err error, httpResp *http.Response, customId *string) {
	httpErrorPrinted := false
	var internalError error
	var body []byte
	if httpResp != nil {
		body, internalError = io.ReadAll(httpResp.Body)
		if internalError == nil {
			tflog.Debug(ctx, "Error HTTP response body: "+string(body))
			var pfError pingFederateErrorResponse
			internalError = json.Unmarshal(body, &pfError)
			if internalError == nil {
				if len(pfError.ValidationErrors) == 0 {
					var errorDetail strings.Builder
					errorDetail.WriteString("Error summary: ")
					errorDetail.WriteString(errorSummary)
					errorDetail.WriteString("\nMessage: ")
					errorDetail.WriteString(pfError.Message)
					errorDetail.WriteString("\nHTTP status: ")
					errorDetail.WriteString(httpResp.Status)
					errorDetail.WriteString("\nResult ID: ")
					errorDetail.WriteString(pfError.ResultId)
					diagnostics.AddError(providererror.PingFederateAPIError, errorDetail.String())
				}
				for _, validationError := range pfError.ValidationErrors {
					var errorDetail strings.Builder
					errorDetail.WriteString("Error summary: ")
					errorDetail.WriteString(errorSummary)
					errorDetail.WriteString("\nMessage: ")
					errorDetail.WriteString(validationError.Message)
					errorDetail.WriteString("\nHTTP status: ")
					errorDetail.WriteString(httpResp.Status)
					if validationError.FieldPath != "" {
						errorDetail.WriteString("\nPingFederate field path: ")
						errorDetail.WriteString(validationError.FieldPath)
					}
					if validationError.ErrorId != "" {
						errorDetail.WriteString("\nError ID: ")
						errorDetail.WriteString(validationError.ErrorId)
					}
					if validationError.DeveloperMessage != "" {
						errorDetail.WriteString("\nDeveloper message: ")
						errorDetail.WriteString(validationError.DeveloperMessage)
					}
					if validationError.FieldPath != "" {
						tfFieldPath := validationError.FieldPath
						if customId != nil && tfFieldPath == "id" {
							tfFieldPath = *customId
						}
						tfAttrName := toTerraformIdentifier(tfFieldPath)
						// Attempt to build the terraform field path from the attribute name
						fieldSteps := strings.Split(tfAttrName, ".")
						var fieldPath path.Path
						for i, step := range fieldSteps {
							if i == 0 {
								fieldPath = path.Root(step)
							} else {
								fieldPath = fieldPath.AtMapKey(step)
							}
						}
						diagnostics.AddAttributeError(fieldPath, providererror.PingFederateValidationError, errorDetail.String())
					} else {
						diagnostics.AddError(providererror.PingFederateValidationError, errorDetail.String())
					}
				}
			} else {
				diagnostics.AddError(providererror.PingFederateAPIError, errorSummary+"\n"+err.Error()+" - Detail:\n"+string(body))
			}
			httpErrorPrinted = true
		}
	}
	if !httpErrorPrinted {
		if internalError != nil {
			tflog.Warn(ctx, "Failed to read HTTP response body: "+internalError.Error())
		}
		diagnostics.AddError(providererror.PingFederateAPIError, errorSummary+"\n"+err.Error())
	}
}
