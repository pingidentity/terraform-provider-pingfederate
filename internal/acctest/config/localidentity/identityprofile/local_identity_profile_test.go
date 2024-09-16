package localidentityprofile_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const localIdentityProfilesId = "testLocalIdProfile"

// Attributes to test with. Add optional properties to test here if desired.
type localIdentityProfilesResourceModel struct {
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
		UsernameField:                       pointers.String("mail"),
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
	textAttrs := map[string]bool{
		"Read-Only":       false,
		"Required":        true,
		"Unique ID Field": true,
		"Mask Log Values": false,
	}
	checkboxGroupAttrs := map[string]bool{
		"Read-Only":       false,
		"Must Pick One":   false,
		"Mask Log Values": false,
	}
	hiddenAttrs := map[string]bool{
		"Unique ID Field": false,
		"Mask Log Values": false,
	}

	return &client.FieldConfig{
		Fields: []client.LocalIdentityField{
			{
				Type:                  "TEXT",
				Id:                    "mail",
				Label:                 "First Name",
				RegistrationPageField: pointers.Bool(true),
				ProfilePageField:      pointers.Bool(true),
				Attributes:            &textAttrs,
				DefaultValue:          pointers.String("defaultValue"),
			},
			{
				Type:                  "HIDDEN",
				Id:                    "cn",
				Label:                 "cn",
				RegistrationPageField: pointers.Bool(false),
				ProfilePageField:      pointers.Bool(true),
				Attributes:            &hiddenAttrs,
			},
			{
				Type:                  "CHECKBOX_GROUP",
				Id:                    "entryUUID",
				Label:                 "entryUUID",
				RegistrationPageField: pointers.Bool(false),
				ProfilePageField:      pointers.Bool(true),
				Attributes:            &checkboxGroupAttrs,
				Options:               []string{"option1, option2"},
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

	defaultValue := ""
	if field.DefaultValue != nil {
		defaultValue = fmt.Sprintf("default_value = \"%s\"", *field.DefaultValue)
	}

	options := ""
	if len(field.Options) > 0 {
		options = fmt.Sprintf("options = %s", acctest.StringSliceToTerraformString(field.Options))
	}

	return fmt.Sprintf(`
	{
		type = "%s"
		id = "%s"
		label = "%s"
		registration_page_field = %t
		profile_page_field = %t
		%s
		%s
		%s
	},
	`, field.Type,
		field.Id,
		field.Label,
		*field.RegistrationPageField,
		*field.ProfilePageField,
		attributes.String(),
		defaultValue,
		options)
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
		FieldStoringVerificationStatus:       "cn",
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

func TestAccLocalIdentityProfiles(t *testing.T) {
	resourceName := "myLocalIdentityProfiles"
	initialResourceModel := localIdentityProfilesResourceModel{
		// Test is only run on attributes that do not require a PD dataStore.
		id:                  localIdentityProfilesId,
		name:                "example",
		registrationEnabled: false,
		profileEnabled:      false,
	}

	updatedResourceModel := localIdentityProfilesResourceModel{
		id:                      localIdentityProfilesId,
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
		CheckDestroy: testAccCheckLocalIdentityProfilesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalIdentityProfiles(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalIdentityProfilesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccLocalIdentityProfiles(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLocalIdentityProfilesAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr("pingfederate_local_identity_profile.myLocalIdentityProfiles", "auth_source_update_policy.store_attributes", fmt.Sprintf("%t", *updatedResourceModel.authSourceUpdatePolicy.StoreAttributes)),
					resource.TestCheckResourceAttr("pingfederate_local_identity_profile.myLocalIdentityProfiles", "auth_source_update_policy.retain_attributes", fmt.Sprintf("%t", *updatedResourceModel.authSourceUpdatePolicy.RetainAttributes)),
					resource.TestCheckResourceAttr("pingfederate_local_identity_profile.myLocalIdentityProfiles", "auth_source_update_policy.update_attributes", fmt.Sprintf("%t", *updatedResourceModel.authSourceUpdatePolicy.UpdateAttributes)),
					resource.TestCheckResourceAttr("pingfederate_local_identity_profile.myLocalIdentityProfiles", "auth_source_update_policy.update_interval", "0"),
				),
			},
			{
				// Test importing the resource
				Config:            testAccLocalIdentityProfiles(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_local_identity_profile." + resourceName,
				ImportStateId:     localIdentityProfilesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccLocalIdentityProfiles(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalIdentityProfilesAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.LocalIdentityIdentityProfilesAPI.DeleteIdentityProfile(ctx, localIdentityProfilesId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccLocalIdentityProfiles(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalIdentityProfilesAttributes(initialResourceModel),
			},
		},
	})
}

func testAccLocalIdentityProfiles(resourceName string, resourceModel localIdentityProfilesResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractsExample" {
  contract_id         = "%[2]s"
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "%[2]s"
}

resource "pingfederate_local_identity_profile" "%[1]s" {
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

data "pingfederate_local_identity_profile" "%[1]s" {
  profile_id = pingfederate_local_identity_profile.%[1]s.id
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
func testAccCheckExpectedLocalIdentityProfilesAttributes(config localIdentityProfilesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "LocalIdentityProfiles"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LocalIdentityIdentityProfilesAPI.GetIdentityProfile(ctx, localIdentityProfilesId).Execute()
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

		if config.authSourceUpdatePolicy != nil {
			err = acctest.TestAttributesMatchBool(resourceType, &config.id, "store_attributes",
				*config.authSourceUpdatePolicy.StoreAttributes, *response.AuthSourceUpdatePolicy.StoreAttributes)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, &config.id, "retain_attributes",
				*config.authSourceUpdatePolicy.RetainAttributes, *response.AuthSourceUpdatePolicy.RetainAttributes)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, &config.id, "update_attributes",
				*config.authSourceUpdatePolicy.UpdateAttributes, *response.AuthSourceUpdatePolicy.UpdateAttributes)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchInt(resourceType, &config.id, "update_interval",
				int64(*config.authSourceUpdatePolicy.UpdateInterval), int64(*response.AuthSourceUpdatePolicy.UpdateInterval))
			if err != nil {
				return err
			}
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
func testAccCheckLocalIdentityProfilesDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.LocalIdentityIdentityProfilesAPI.DeleteIdentityProfile(ctx, localIdentityProfilesId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("LocalIdentityProfiles", localIdentityProfilesId)
	}
	return nil
}
