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

	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Verify that any required environment variables are set before the test begins
func ConfigurationPreCheck(t *testing.T) {
	envVars := []string{
		"PINGFEDERATE_PROVIDER_HTTPS_HOST",
		"PINGFEDERATE_PROVIDER_USERNAME",
		"PINGFEDERATE_PROVIDER_PASSWORD",
		"PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS",
		"X_BYPASS_EXTERNAL_VALIDATION",
	}

	errorFound := false
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Errorf("The '%s' environment variable must be set to run acceptance tests", envVar)
			errorFound = true
		}
	}

	if errorFound {
		t.FailNow()
	}
}

func TestClient() *client.APIClient {
	httpsHost := os.Getenv("PINGFEDERATE_PROVIDER_HTTPS_HOST")
	clientConfig := client.NewConfiguration()
	clientConfig.DefaultHeader["X-Xsrf-Header"] = "PingFederate"
	clientConfig.DefaultHeader["X-BypassExternalValidation"] = os.Getenv("X_BYPASS_EXTERNAL_VALIDATION")
	clientConfig.Servers = client.ServerConfigurations{
		{
			URL: httpsHost + "/pf-admin-api/v1",
		},
	}
	// Trusting all for the acceptance tests, since they run on localhost
	// May want to incorporate actual trust here in the future.
	//#nosec G402
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	clientConfig.HTTPClient = httpClient
	return client.NewAPIClient(clientConfig)
}

func TestBasicAuthContext() context.Context {
	ctx := context.Background()
	return config.BasicAuthContext(ctx, os.Getenv("PINGFEDERATE_PROVIDER_USERNAME"), os.Getenv("PINGFEDERATE_PROVIDER_PASSWORD"))
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
	if !types.SetsEqual(expected, found) {
		return mismatchedAttributeError(resourceType, resourceName, attributeName, StringSliceToTerraformString(expected), StringSliceToTerraformString(found))
	}
	return nil
}

// Test if float slice attributes match
func TestAttributesMatchFloatSlice(resourceType string, resourceName *string, attributeName string, expected, found []float64) error {
	if !types.FloatSetsEqual(expected, found) {
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
