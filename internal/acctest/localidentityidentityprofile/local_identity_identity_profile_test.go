package acctest_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const localIdentityIdentityProfilesId = "testLocalIdProfile"

// Attributes to test with. Add optional properties to test here if desired.
type localIdentityIdentityProfilesResourceModel struct {
	id                      string
	name                    string
	registrationEnabled     bool
	profileEnabled          bool
	authSources             []string
	authSourceUpdatePolicy  *client.LocalIdentityAuthSourceUpdatePolicy
	registrationConfig      *client.RegistrationConfig
	profileConfig           *client.ProfileConfig
	fieldConfig             *client.FieldConfig
	emailVerificationConfig *client.EmailVerificationConfig
	dataStoreConfig         *client.LdapDataStoreConfig
}

func authSourcesHcl(sources []string) string {
	if len(sources) == 0 {
		// Just leave out empty auth sources for the sake of this test
		return ""
	}

	var hcl strings.Builder
	hcl.WriteString("auth_sources = [\n")
	for _, source := range sources {
		hcl.WriteString("    {\n")
		hcl.WriteString(fmt.Sprintf("        source = \"%s\"\n", source))
		hcl.WriteString("    },\n")
	}
	hcl.WriteString("]\n")

	return hcl.String()
}

func updatedAuthSourceUpdatePolicy() *client.LocalIdentityAuthSourceUpdatePolicy {
	return &client.LocalIdentityAuthSourceUpdatePolicy{
		StoreAttributes:  pointers.Bool(false),
		RetainAttributes: pointers.Bool(false),
		UpdateAttributes: pointers.Bool(false),
		UpdateInterval:   pointers.Float64(0),
	}
}

func authSourceUpdatePolicyHcl(policy *client.LocalIdentityAuthSourceUpdatePolicy) string {
	if policy == nil {
		return ""
	}

	return fmt.Sprintf(`
	auth_source_update_policy = {
		store_attributes = %[1]t
		retain_attributes = %[2]t
		update_attributes = %[3]t
		update_interval = %[4]f
	}
	`, *policy.StoreAttributes, *policy.RetainAttributes, *policy.UpdateAttributes, *policy.UpdateInterval)
}

func updatedRegistrationConfig() *client.RegistrationConfig {
	return &client.RegistrationConfig{
		CaptchaEnabled: pointers.Bool(true),
		CaptchaProviderRef: &client.ResourceLink{
			Id: "exampleCaptchaProvider",
		},
		TemplateName:                        "local.identity.registration.html",
		CreateAuthnSessionAfterRegistration: pointers.Bool(true),
		UsernameField:                       pointers.String("cn"),
		ThisIsMyDeviceEnabled:               pointers.Bool(false),
	}
}

func registrationConfigHcl(config *client.RegistrationConfig) string {
	if config == nil {
		return ""
	}

	return fmt.Sprintf(`
	registration_config = {
		captcha_enabled = %[1]t
		captcha_provider_ref = {
			id = "%[2]s"
		}
		template_name = "%[3]s"
		create_authn_session_after_registration = %[4]t
		username_field = "%[5]s"
		this_is_my_device_enabled = %[6]t
	}
	`, *config.CaptchaEnabled,
		config.CaptchaProviderRef.Id,
		config.TemplateName,
		*config.CreateAuthnSessionAfterRegistration,
		*config.UsernameField,
		*config.ThisIsMyDeviceEnabled)
}

func updatedProfileConfig() *client.ProfileConfig {
	return &client.ProfileConfig{
		DeleteIdentityEnabled: pointers.Bool(true),
		TemplateName:          "local.identity.profile.html",
	}
}

func profileConfigHcl(config *client.ProfileConfig) string {
	if config == nil {
		return ""
	}

	return fmt.Sprintf(`
	profile_config = {
		delete_identity_enabled = %[1]t
		template_name = "%[2]s"
	}
	`, *config.DeleteIdentityEnabled, config.TemplateName)
}

func updatedFieldConfig() *client.FieldConfig {
	emailAttrs := map[string]bool{
		"Read-Only":       false,
		"Required":        true,
		"Unique ID Field": true,
		"Mask Log Values": false,
	}
	textAttrs := map[string]bool{
		"Read-Only":       false,
		"Required":        true,
		"Unique ID Field": false,
		"Mask Log Values": false,
	}
	hiddenAttrs := map[string]bool{
		"Unique ID Field": false,
		"Mask Log Values": false,
	}

	return &client.FieldConfig{
		Fields: []client.LocalIdentityField{
			{
				Type:                  "EMAIL",
				Id:                    "mail",
				Label:                 "Email address",
				RegistrationPageField: pointers.Bool(true),
				ProfilePageField:      pointers.Bool(true),
				Attributes:            &emailAttrs,
			},
			{
				Type:                  "TEXT",
				Id:                    "cn",
				Label:                 "First Name",
				RegistrationPageField: pointers.Bool(true),
				ProfilePageField:      pointers.Bool(true),
				Attributes:            &textAttrs,
			},
			{
				Type:                  "HIDDEN",
				Id:                    "entryUUID",
				Label:                 "entryUUID",
				RegistrationPageField: pointers.Bool(true),
				ProfilePageField:      pointers.Bool(true),
				Attributes:            &hiddenAttrs,
			},
		},
		StripSpaceFromUniqueField: pointers.Bool(true),
	}
}

func fieldHcl(field client.LocalIdentityField) string {
	var attributes strings.Builder
	if field.Attributes != nil {
		attributes.WriteString("attributes = {\n")
		for name, val := range *field.Attributes {
			attributes.WriteString(fmt.Sprintf("\"%s\" = %t,\n", name, val))
		}
		attributes.WriteString("}\n")
	}

	return fmt.Sprintf(`
	{
		type = "%s"
		id = "%s"
		label = "%s"
		registration_page_field = %t
		profile_page_field = %t
		%s
	},
	`, field.Type,
		field.Id,
		field.Label,
		*field.RegistrationPageField,
		*field.ProfilePageField,
		attributes.String())
}

func fieldConfigHcl(config *client.FieldConfig) string {
	if config == nil {
		return ""
	}

	var fields strings.Builder
	fields.WriteString("fields = [\n")
	for _, field := range config.Fields {
		fields.WriteString(fieldHcl(field))
	}
	fields.WriteString("]\n")

	return fmt.Sprintf(`
	field_config = {
		%s
		strip_space_from_unique_field = %t
	}
	`, fields.String(), *config.StripSpaceFromUniqueField)
}

func updatedEmailVerificationConfig() *client.EmailVerificationConfig {
	return &client.EmailVerificationConfig{
		EmailVerificationEnabled:             pointers.Bool(true),
		VerifyEmailTemplateName:              pointers.String("message-template-email-ownership-verification.html"),
		EmailVerificationSuccessTemplateName: pointers.String("local.identity.email.verification.success.html"),
		EmailVerificationErrorTemplateName:   pointers.String("local.identity.email.verification.error.html"),
		EmailVerificationType:                pointers.String("OTP"),
		AllowedOtpCharacterSet:               pointers.String("23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz"),
		EmailVerificationOtpTemplateName:     pointers.String("message-template-email-ownership-verification.html"),
		OtpLength:                            pointers.Int64(8),
		OtpRetryAttempts:                     pointers.Int64(3),
		OtpTimeToLive:                        pointers.Int64(1440),
		FieldForEmailToVerify:                "mail",
		FieldStoringVerificationStatus:       "entryUUID",
		NotificationPublisherRef: &client.ResourceLink{
			Id: "exampleSmtpPublisher",
		},
		RequireVerifiedEmail: pointers.Bool(true),
	}
}

func emailVerificationConfigHcl(config *client.EmailVerificationConfig) string {
	if config == nil {
		return ""
	}

	return fmt.Sprintf(`
	email_verification_config = {
		email_verification_enabled = %t
		verify_email_template_name = "%s"
		email_verification_success_template_name = "%s"
		email_verification_error_template_name   = "%s"
		email_verification_type              = "%s"
		allowed_otp_character_set            = "%s"
		email_verification_otp_template_name = "%s"
		otp_length                           = %d
		otp_retry_attempts                   = %d
		otp_time_to_live                     = %d
		field_for_email_to_verify         = "%s"
		field_storing_verification_status = "%s"
		notification_publisher_ref = {
		  id = "%s",
		}
		require_verified_email = %t
	}`, *config.EmailVerificationEnabled,
		*config.VerifyEmailTemplateName,
		*config.EmailVerificationSuccessTemplateName,
		*config.EmailVerificationErrorTemplateName,
		*config.EmailVerificationType,
		*config.AllowedOtpCharacterSet,
		*config.EmailVerificationOtpTemplateName,
		*config.OtpLength,
		*config.OtpRetryAttempts,
		*config.OtpTimeToLive,
		config.FieldForEmailToVerify,
		config.FieldStoringVerificationStatus,
		config.NotificationPublisherRef.Id,
		*config.RequireVerifiedEmail)
}

func updatedDataStoreConfig() *client.LdapDataStoreConfig {
	return &client.LdapDataStoreConfig{
		Type: "LDAP",
		DataStoreRef: client.ResourceLink{
			Id: "pingdirectory",
		},
		BaseDn:        "ou=people,dc=example,dc=com",
		CreatePattern: "uid=$${mail}",
		ObjectClass:   "inetOrgPerson",
		DataStoreMapping: map[string]client.DataStoreAttribute{
			"entryUUID": {
				Type: "LDAP",
				Name: "entryUUID",
			},
			"cn": {
				Type: "LDAP",
				Name: "cn",
			},
			"mail": {
				Type: "LDAP",
				Name: "mail",
			},
		},
	}
}

func dataStoreAttributeHcl(attribute client.DataStoreAttribute) string {
	return fmt.Sprintf(`
	"%[1]s" = {
		type = "%[2]s"
		name = "%[1]s"
		metadata = {}
	}
	`, attribute.Name, attribute.Type)
}

func dataStoreConfigHcl(config *client.LdapDataStoreConfig) string {
	if config == nil {
		return ""
	}

	var dataStoreMapping strings.Builder
	dataStoreMapping.WriteString("data_store_mapping = {\n")
	for _, attr := range config.DataStoreMapping {
		dataStoreMapping.WriteString(dataStoreAttributeHcl(attr))
	}
	dataStoreMapping.WriteString("}\n")

	return fmt.Sprintf(`
	data_store_config = {
		type = "%s"
		data_store_ref = {
			id = "%s"
		}
		base_dn = "%s"
		create_pattern = "%s"
		object_class = "%s"
		%s
	}
	`, config.Type,
		config.DataStoreRef.Id,
		config.BaseDn,
		config.CreatePattern,
		config.ObjectClass,
		dataStoreMapping.String())
}

func TestAccLocalIdentityIdentityProfiles(t *testing.T) {
	resourceName := "myLocalIdentityIdentityProfiles"
	initialResourceModel := localIdentityIdentityProfilesResourceModel{
		// Test is only run on attributes that do not require a PD dataStore.
		id:                  localIdentityIdentityProfilesId,
		name:                "example",
		registrationEnabled: false,
		profileEnabled:      false,
	}
	updatedResourceModel := localIdentityIdentityProfilesResourceModel{
		id:                      localIdentityIdentityProfilesId,
		name:                    "example1",
		registrationEnabled:     true,
		profileEnabled:          true,
		authSources:             []string{"example", "test"},
		authSourceUpdatePolicy:  updatedAuthSourceUpdatePolicy(),
		registrationConfig:      updatedRegistrationConfig(),
		profileConfig:           updatedProfileConfig(),
		fieldConfig:             updatedFieldConfig(),
		emailVerificationConfig: updatedEmailVerificationConfig(),
		dataStoreConfig:         updatedDataStoreConfig(),
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckLocalIdentityIdentityProfilesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalIdentityIdentityProfiles(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalIdentityIdentityProfilesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccLocalIdentityIdentityProfiles(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLocalIdentityIdentityProfilesAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLocalIdentityIdentityProfiles(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_local_identity_identity_profile." + resourceName,
				ImportStateId:     localIdentityIdentityProfilesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccLocalIdentityIdentityProfiles(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalIdentityIdentityProfilesAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.LocalIdentityIdentityProfilesAPI.DeleteIdentityProfile(ctx, updatedResourceModel.id).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccLocalIdentityIdentityProfiles(resourceName string, resourceModel localIdentityIdentityProfilesResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractsExample" {
  contract_id         = "%[2]s"
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "%[2]s"
}

resource "pingfederate_local_identity_identity_profile" "%[1]s" {
  profile_id = "%[2]s"
  name       = "%[3]s"
  apc_id = {
    id = pingfederate_authentication_policy_contract.authenticationPolicyContractsExample.id
  }
  registration_enabled = %[4]t
  profile_enabled      = %[5]t
  %[6]s
  %[7]s
  %[8]s
  %[9]s
  %[10]s
  %[11]s
  %[12]s
}

data "pingfederate_local_identity_identity_profile" "%[1]s" {
  profile_id = pingfederate_local_identity_identity_profile.%[1]s.id
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.registrationEnabled,
		resourceModel.profileEnabled,
		authSourcesHcl(resourceModel.authSources),
		authSourceUpdatePolicyHcl(resourceModel.authSourceUpdatePolicy),
		registrationConfigHcl(resourceModel.registrationConfig),
		profileConfigHcl(resourceModel.profileConfig),
		fieldConfigHcl(resourceModel.fieldConfig),
		emailVerificationConfigHcl(resourceModel.emailVerificationConfig),
		dataStoreConfigHcl(resourceModel.dataStoreConfig),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedLocalIdentityIdentityProfilesAttributes(config localIdentityIdentityProfilesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "LocalIdentityIdentityProfiles"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LocalIdentityIdentityProfilesAPI.GetIdentityProfile(ctx, localIdentityIdentityProfilesId).Execute()

		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, &config.id, "id",
			config.id, *response.Id)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "registration_enabled",
			config.registrationEnabled, *response.RegistrationEnabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "profile_enabled",
			config.profileEnabled, *response.ProfileEnabled)
		if err != nil {
			return err
		}

		if config.registrationConfig != nil {
			err = acctest.TestAttributesMatchBool(resourceType, &config.id, "registration_config.captcha_enabled",
				*config.registrationConfig.CaptchaEnabled, *response.RegistrationConfig.CaptchaEnabled)
			if err != nil {
				return err
			}
		}

		if config.dataStoreConfig != nil {
			err = acctest.TestAttributesMatchString(resourceType, &config.id, "data_store_config.base_dn",
				config.dataStoreConfig.BaseDn, response.DataStoreConfig.BaseDn)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLocalIdentityIdentityProfilesDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.LocalIdentityIdentityProfilesAPI.DeleteIdentityProfile(ctx, localIdentityIdentityProfilesId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("LocalIdentityIdentityProfiles", localIdentityIdentityProfilesId)
	}
	return nil
}
