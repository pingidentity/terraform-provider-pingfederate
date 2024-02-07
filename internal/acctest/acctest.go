package acctest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/authentication"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Verify that any required environment variables are set before the test begins
func ConfigurationPreCheck(t *testing.T) {
	envVars := []string{
		"PINGFEDERATE_PROVIDER_HTTPS_HOST",
		"PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS",
		"PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER",
		"PINGFEDERATE_PROVIDER_PRODUCT_VERSION",
	}

	errorFound := false
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Errorf("The '%s' environment variable must be set to run acceptance tests", envVar)
			errorFound = true
		}
	}

	// Verify that the version supplied in the environment can be parsed
	versionVar := os.Getenv("PINGFEDERATE_PROVIDER_PRODUCT_VERSION")
	_, err := version.Parse(versionVar)
	if err != nil {
		t.Errorf("The '%s' value for the 'PINGFEDERATE_PROVIDER_PRODUCT_VERSION' environment variable is not a valid version: %s", versionVar, err.Error())
		errorFound = true
	}

	if errorFound {
		t.FailNow()
	}
}

func GetTransport() *http.Transport {
	// Trusting all for the acceptance tests, since they run on localhost
	// May want to incorporate actual trust here in the future.
	//#nosec G402
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func TestClient() *client.APIClient {
	httpsHost := os.Getenv("PINGFEDERATE_PROVIDER_HTTPS_HOST")
	adminApiPath := os.Getenv("PINGFEDERATE_PROVIDER_ADMIN_API_PATH")
	clientConfig := client.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	clientConfig.DefaultHeader["X-BypassExternalValidation"] = os.Getenv("PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER")
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + adminApiPath,
		},
	}

	httpClient := &http.Client{Transport: GetTransport()}
	clientConfig.HTTPClient = httpClient
	return client.NewAPIClient(clientConfig)
}

// lintignore:AT008
func TestAccessTokenContext(accessToken string) context.Context {
	ctx := context.Background()
	if accessToken == "" {
		fmt.Println("No access token found in environment")
		return nil
	}

	return config.AccessTokenContext(ctx, accessToken)
}

func TestBasicAuthContext() context.Context {
	ctx := context.Background()
	envVars, errors := authentication.TestEnvVarSlice([]string{"PINGFEDERATE_PROVIDER_USERNAME", "PINGFEDERATE_PROVIDER_PASSWORD"}, "acctest.go")
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		return nil
	}

	return config.BasicAuthContext(ctx, envVars["PINGFEDERATE_PROVIDER_USERNAME"], envVars["PINGFEDERATE_PROVIDER_PASSWORD"])
}

func TestOauth2Context() context.Context {
	ctx := context.Background()
	return config.OAuthContext(ctx, GetTransport(), os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL"), os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID"), os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET"), []string{os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_SCOPES")})
}

// Convert a string slice to the format used in Terraform files
func StringSliceToTerraformString(values []string) string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, str := range values {
		builder.WriteString(fmt.Sprintf("\"%s\"", str))
		if i < len(values)-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("]")
	return builder.String()
}

// Convert a string slice to the format used in Terraform files
func ObjectSliceOfKvStringsToTerraformString(keyValue string, values []string) string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, str := range values {
		builder.WriteString(fmt.Sprintf("{%s = \"%s\"}", keyValue, str))
		if i < len(values)-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("]")
	return builder.String()
}

// Convert a float64 slice to the format used in Terraform files
func FloatSliceToTerraformString(values []float64) string {
	var builder strings.Builder
	builder.WriteString("[")
	string := ""
	for _, v := range values {
		if len(string) > 0 {
			string += ","
		}
		string += fmt.Sprintf("%f", v)
	}
	builder.WriteString(string)
	builder.WriteString("]")
	return builder.String()
}

func FloatSliceToStringSlice(values []float64) []string {
	stringSlice := make([]string, 0, len(values))
	for _, v := range values {
		element := fmt.Sprintf("%f", v)
		stringSlice = append(stringSlice, element)
	}
	return stringSlice
}

func InterfaceSliceToStringSlice(values []interface{}) []string {
	stringSlice := make([]string, 0, len(values))
	for _, v := range values {
		element := fmt.Sprintf("%s", v)
		stringSlice = append(stringSlice, element)
	}
	return stringSlice
}

func TfKeyValuePairToString(key string, value string, addDoubleQuotes bool) string {
	if len(value) > 0 && value != "0" {
		quoteVal := func() string {
			if addDoubleQuotes {
				return "\""
			}
			return ""
		}

		q := quoteVal()
		return fmt.Sprintf("%s = %s%s%s", key, q, value, q)
	} else {
		return ""
	}
}

// Utility methods for testing whether attributes match the expected values

// Test if string attributes match
func TestAttributesMatchString(resourceType string, resourceName *string, attributeName, expected, found string) error {
	if expected != found {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, expected, found)
	}
	return nil
}

// Test if expected string matches found string pointer
func TestAttributesMatchStringPointer(resourceType string, resourceName *string, attributeName, expected string, found *string) error {
	if found == nil && expected != "" {
		// Expect empty string to match nil pointer
		return missingAttributeError(resourceType, resourceName, attributeName, expected)
	}
	if found != nil {
		return TestAttributesMatchString(resourceType, resourceName, attributeName, expected, *found)
	}
	return nil
}

// Test if boolean attributes match
func TestAttributesMatchBool(resourceType string, resourceName *string, attributeName string, expected, found bool) error {
	if expected != found {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, strconv.FormatBool(expected), strconv.FormatBool(found))
	}
	return nil
}

// Test if float64 attributes match
func TestAttributesMatchFloat(resourceType string, resourceName *string, attributeName string, expected, found float64) error {
	if expected != found {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, fmt.Sprintf("%f", expected), fmt.Sprintf("%f", found))
	}
	return nil
}

// Test if int attributes match
func TestAttributesMatchInt(resourceType string, resourceName *string, attributeName string, expected, found int64) error {
	if expected != found {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, strconv.FormatInt(expected, 10), strconv.FormatInt(found, 10))
	}
	return nil
}

// Test if string slice attributes match
func TestAttributesMatchStringSlice(resourceType string, resourceName *string, attributeName string, expected, found []string) error {
	if !types.StringSlicesEqual(expected, found) {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, StringSliceToTerraformString(expected), StringSliceToTerraformString(found))
	}
	return nil
}

// Test if float slice attributes match
func TestAttributesMatchFloatSlice(resourceType string, resourceName *string, attributeName string, expected, found []float64) error {
	if !types.FloatSlicesEqual(expected, found) {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, FloatSliceToTerraformString(expected), FloatSliceToTerraformString(found))
	}
	return nil
}

func ExpectedDestroyError(resourceType, resourceName string) error {
	return fmt.Errorf("%s '%s' still exists after tests. Expected it to be destroyed", resourceType, resourceName)
}

func mismatchedAttributeError(resourceType string, resourceName *string, attributeName, expected, found string) error {
	if resourceName == nil {
		return mismatchedAttributeErrorSingletonResource(resourceType, attributeName, expected, found)
	}
	return fmt.Errorf("mismatched %s attribute for %s '%s'. expected '%s', found '%s'", attributeName, resourceType, *resourceName, expected, found)
}

func mismatchedAttributeErrorSingletonResource(resourceType, attributeName, expected, found string) error {
	return fmt.Errorf("mismatched %s attribute for %s. expected '%s', found '%s'", attributeName, resourceType, expected, found)
}

func missingAttributeError(resourceType string, resourceName *string, attributeName, expected string) error {
	if resourceName == nil {
		return missingAttributeErrorSingletonResource(resourceType, attributeName, expected)
	}
	return fmt.Errorf("missing %s attribute for %s '%s'. expected '%s'", attributeName, resourceType, *resourceName, expected)
}

func missingAttributeErrorSingletonResource(resourceType, attributeName, expected string) error {
	return fmt.Errorf("missing %s attribute for %s. expected '%s'", attributeName, resourceType, expected)
}

// Check that the version being tested is at least the given minimum version
var versionAtLeastResults = map[version.SupportedVersion]bool{}

func VersionAtLeast(minimumVersion version.SupportedVersion) bool {
	savedResult, ok := versionAtLeastResults[minimumVersion]
	if ok {
		return savedResult
	}

	// Just swallow the errors here since we already verify that the product version environment variable is valid in the pre-check.
	// Assume the version passed in is valid.
	supportedVersion, _ := version.Parse(os.Getenv("PINGFEDERATE_PROVIDER_PRODUCT_VERSION"))
	compare, _ := version.Compare(supportedVersion, minimumVersion)
	return compare >= 0
}
