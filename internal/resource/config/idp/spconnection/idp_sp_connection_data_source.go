package idpspconnection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	datasourceattributecontractfulfillment "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	datasourceattributesources "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	datasourceissuancecriteria "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
	datasourceresourcelink "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &idpSpConnectionDataSource{}
	_ datasource.DataSourceWithConfigure = &idpSpConnectionDataSource{}
)

var (
	targetSettingsDataSourceElemAttrType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
	}}

	channelsElemDataSourceAttrType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"active":            types.BoolType,
		"channel_source":    types.ObjectType{AttrTypes: channelSourceAttrTypes},
		"attribute_mapping": types.SetType{ElemType: attributeMappingElemAttrTypes},
		"name":              types.StringType,
		"max_threads":       types.Int64Type,
		"timeout":           types.Int64Type,
	}}

	outboundProvisionDataSourceAttrTypes = map[string]attr.Type{
		"type":            types.StringType,
		"target_settings": types.ListType{ElemType: targetSettingsDataSourceElemAttrType},
		"custom_schema":   types.ObjectType{AttrTypes: customSchemaAttrTypes},
		"channels":        types.ListType{ElemType: channelsElemDataSourceAttrType},
	}
)

// IdpSpConnectionDataSource is a helper function to simplify the provider implementation.
func IdpSpConnectionDataSource() datasource.DataSource {
	return &idpSpConnectionDataSource{}
}

// idpSpConnectionDataSource is the datasource implementation.
type idpSpConnectionDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *idpSpConnectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	certsSchema := schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"active_verification_cert": schema.BoolAttribute{
					Computed:    true,
					Optional:    false,
					Description: "Indicates whether this is an active signature verification certificate.",
				},
				"cert_view": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"crypto_provider": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
						},
						"expires": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The end date up until which the item is valid, in ISO 8601 format (UTC).",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The persistent, unique ID for the certificate.",
						},
						"issuer_dn": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The issuer's distinguished name.",
						},
						"key_algorithm": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The public key algorithm.",
						},
						"key_size": schema.Int64Attribute{
							Computed:    true,
							Optional:    false,
							Description: "The public key size.",
						},
						"serial_number": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The serial number assigned by the CA.",
						},
						"sha1fingerprint": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "SHA-1 fingerprint in Hex encoding.",
						},
						"sha256fingerprint": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "SHA-256 fingerprint in Hex encoding.",
						},
						"signature_algorithm": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The signature algorithm.",
						},
						"status": schema.StringAttribute{
							Computed: true,
							Optional: false,
						},
						"subject_alternative_names": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Optional:    false,
							Description: "The subject alternative names (SAN).",
						},
						"subject_dn": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The subject's distinguished name.",
						},
						"valid_from": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The start date from which the item is valid, in ISO 8601 format (UTC).",
						},
						"version": schema.Int64Attribute{
							Computed:    true,
							Optional:    false,
							Description: "The X.509 version to which the item conforms.",
						},
					},
					Computed:    true,
					Optional:    false,
					Description: "Certificate details.",
				},
				"encryption_cert": schema.BoolAttribute{
					Computed:    true,
					Optional:    false,
					Description: "Indicates whether to use this cert to encrypt outgoing assertions.",
				},
				"primary_verification_cert": schema.BoolAttribute{
					Computed:    true,
					Optional:    false,
					Description: "Indicates whether this is the primary signature verification certificate.",
				},
				"secondary_verification_cert": schema.BoolAttribute{
					Computed:    true,
					Optional:    false,
					Description: "Indicates whether this is the secondary signature verification certificate.",
				},
				"x509file": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"crypto_provider": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "Cryptographic Provider.",
						},
						"file_data": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The certificate data in PEM format.",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Optional:    false,
							Description: "The persistent, unique ID for the certificate.",
						},
					},
					Computed:    true,
					Optional:    false,
					Description: "Encoded certificate data.",
				},
			},
		},
		Computed:    true,
		Optional:    false,
		Description: "The certificates used for signature verification and XML encryption.",
	}

	httpBasicCredentialsSchema := schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"encrypted_password": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "For GET requests, this field contains the encrypted password, if one exists.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "User password.  To update the password, specify the plaintext value in this field.",
			},
			"username": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The username.",
			},
		},
		Computed:    true,
		Optional:    false,
		Description: "Username and password credentials.",
	}

	adapterOverrideSettingsAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"masked": schema.BoolAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Specifies whether this attribute is masked in PingFederate logs.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The name of this attribute.",
			},
			"pseudonym": schema.BoolAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Specifies whether this attribute is used to construct a pseudonym for the SP.",
			},
		},
	}

	spBrowserSSOAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The name of this attribute.",
			},
			"name_format": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The SAML Name Format for the attribute.",
			},
		},
	}

	wsTrustAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The name of this attribute.",
			},
			"namespace": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The attribute namespace.",
			},
		},
	}

	messageCustomizationsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"context_name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The context in which the customization will be applied.",
			},
			"message_expression": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The OGNL expression that will be executed.",
			},
		},
	}

	channelsAttributeMappingNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"field_name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The name of target field.",
			},
			"saas_field_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_names": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Optional:    false,
						Description: "The list of source attribute names used to generate or map to a target field",
					},
					"character_case": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The character case of the field value.",
					},
					"create_only": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Indicates whether this field is a create only field and cannot be updated.",
					},
					"default_value": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The default value for the target field",
					},
					"expression": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "An OGNL expression to obtain a value.",
					},
					"masked": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Indicates whether the attribute should be masked in server logs.",
					},
					"parser": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Indicates how the field shall be parsed.",
					},
					"trim": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Indicates whether field should be trimmed before provisioning.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "The settings that represent how attribute values from source data store will be mapped into Fields specified by the service provider.",
			},
		},
	}

	outboundProvisionTargetSettingsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The name of the configuration field.",
			},
			"value": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The value for the configuration field.",
			},
			"encrypted_value": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The encrypted value for the configuration field.",
			},
		},
	}

	schema := schema.Schema{
		Description: "Describes an IdP SP Connection",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Specifies whether the connection is active and ready to process incoming requests.",
			},
			"additional_allowed_entities_configuration": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"additional_allowed_entities": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"entity_description": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Entity description.",
								},
								"entity_id": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Unique entity identifier.",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "An array of additional allowed entities or issuers to be accepted during entity or issuer validation.",
					},
					"allow_additional_entities": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Set to true to configure additional entities or issuers to be accepted during entity or issuer validation.",
					},
					"allow_all_entities": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Set to true to accept any entity or issuer during entity or issuer validation.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "Additional allowed entities or issuers configuration. Currently only used in OIDC IdP (RP) connection.",
			},
			"application_icon_url": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The application icon url.",
			},
			"application_name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The application name.",
			},
			"attribute_query": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_contract_fulfillment": datasourceattributecontractfulfillment.ToDataSourceSchema(),
					"attribute_sources":              datasourceattributesources.ToDataSourceSchema(),
					"attributes": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Optional:    false,
						Description: "The list of attributes that may be returned to the SP in the response to an attribute request.",
					},
					"issuance_criteria": datasourceissuancecriteria.ToDataSourceSchema(),
					"policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt_assertion": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Encrypt the assertion.",
							},
							"require_encrypted_name_id": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Require an encrypted name identifier.",
							},
							"require_signed_attribute_query": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Require signed attribute query.",
							},
							"sign_assertion": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Sign the assertion.",
							},
							"sign_response": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Sign the response.",
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "The attribute query profile's security policy.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "The attribute query profile supports SPs in requesting user attributes.",
			},
			"base_url": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The fully-qualified hostname and port on which your partner's federation deployment runs.",
			},
			"connection_target_type": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The connection target type. This field is intended for bulk import/export usage.",
			},
			"contact_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"company": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Company name.",
					},
					"email": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Contact email address.",
					},
					"first_name": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Contact first name.",
					},
					"last_name": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Contact last name.",
					},
					"phone": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Contact phone number.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "Contact information.",
			},
			"creation_date": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The time at which the connection was created.",
			},
			"credentials": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"block_encryption_algorithm": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "The algorithm used to encrypt assertions sent to this partner.",
					},
					"certs":                   certsSchema,
					"decryption_key_pair_ref": resourcelink.SingleNestedAttribute(),
					"inbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"certs": certsSchema,
							"digital_signature": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "If incoming or outgoing messages must be signed.",
							},
							"http_basic_credentials": httpBasicCredentialsSchema,
							"require_ssl": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Incoming HTTP transmissions must use a secure channel.",
							},
							"type": schema.StringAttribute{
								Computed:    true,
								Optional:    false,
								Description: "The back channel authentication type.",
							},
							"verification_issuer_dn": schema.StringAttribute{
								Computed:    true,
								Optional:    false,
								Description: "If a verification Subject DN is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
							},
							"verification_subject_dn": schema.StringAttribute{
								Computed:    true,
								Optional:    false,
								Description: "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.",
							},
						},
						Computed: true,
						Optional: false,
					},
					"key_transport_algorithm": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "The algorithm used to transport keys to this partner.",
					},
					"outbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"digital_signature": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "If incoming or outgoing messages must be signed.",
							},
							"http_basic_credentials": httpBasicCredentialsSchema,
							"ssl_auth_key_pair_ref":  datasourceresourcelink.ToDataSourceSchemaSingleNestedAttribute(),
							"type": schema.StringAttribute{
								Computed:    true,
								Optional:    false,
								Description: "The back channel authentication type.",
							},
							"validate_partner_cert": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Validate the partner server certificate.",
							},
						},
						Optional: true,
					},
					"secondary_decryption_key_pair_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttribute(),
					"signing_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Computed:    true,
								Optional:    false,
								Description: "The algorithm used to sign messages sent to this partner. The default is SHA1withDSA for DSA certs, SHA256withRSA for RSA certs, and SHA256withECDSA for EC certs. For RSA certs, SHA1withRSA, SHA384withRSA, SHA512withRSA, SHA256withRSAandMGF1, SHA384withRSAandMGF1 and SHA512withRSAandMGF1 are also supported. For EC certs, SHA384withECDSA and SHA512withECDSA are also supported. If the connection is WS-Federation with JWT token type, then the possible values are RSA SHA256, RSA SHA384, RSA SHA512, RSASSA-PSS SHA256, RSASSA-PSS SHA384, RSASSA-PSS SHA512, ECDSA SHA256, ECDSA SHA384, ECDSA SHA512",
							},
							"alternative_signing_key_pair_refs": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: datasourceresourcelink.ToDataSourceSchema(),
								},
								Computed:    true,
								Optional:    false,
								Description: "The list of IDs of alternative key pairs used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console.",
							},
							"include_cert_in_signature": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Determines whether the signing certificate is included in the signature <KeyInfo> element.",
							},
							"include_raw_key_in_signature": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Determines whether the <KeyValue> element with the raw public key is included in the signature <KeyInfo> element.",
							},
							"signing_key_pair_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttribute(),
						},
						Computed:    true,
						Optional:    false,
						Description: "Settings related to signing messages sent to this partner.",
					},
					"verification_issuer_dn": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "If a verification Subject DN is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
					},
					"verification_subject_dn": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "The certificates and settings for encryption, signing, and signature verification.",
			},
			"default_virtual_entity_id": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The default alternate entity ID that identifies the local server to this partner.",
			},
			"entity_id": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The partner's entity ID (connection ID) or issuer value (for OIDC Connections).",
			},
			"extended_properties": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Optional:    false,
							Description: "A List of values",
						},
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "Extended Properties allows to store additional information for IdP/SP Connections. The names of these extended properties should be defined in /extendedProperties.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The persistent, unique ID for the connection.",
			},
			"license_connection_group": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The license connection group.",
			},
			"logging_mode": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The level of transaction logging applicable for this connection.",
			},
			"metadata_reload_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable_auto_metadata_update": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Specifies whether the metadata of the connection will be automatically reloaded.",
					},
					"metadata_url_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttribute(),
				},
				Computed:    true,
				Optional:    false,
				Description: "Configuration settings to enable automatic reload of partner's metadata.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "The connection name.",
			},
			"outbound_provision": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"channels": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"active": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Indicates whether the channel is the active channel for this connection.",
								},
								"attribute_mapping": schema.SetNestedAttribute{
									NestedObject: channelsAttributeMappingNestedObject,
									Computed:     true,
									Optional:     false,
									Description:  "The mapping of attributes from the local data store into Fields specified by the service provider.",
								},
								"channel_source": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"account_management_settings": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"account_status_algorithm": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The account status algorithm name.",
												},
												"account_status_attribute_name": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The account status attribute name.",
												},
												"default_status": schema.BoolAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The default status of the account.",
												},
												"flag_comparison_status": schema.BoolAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The flag that represents comparison status.",
												},
												"flag_comparison_value": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The flag that represents comparison value.",
												},
											},
											Computed:    true,
											Optional:    false,
											Description: "Account management settings.",
										},
										"base_dn": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "The base DN where the user records are located.",
										},
										"change_detection_settings": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"changed_users_algorithm": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The changed user algorithm. \nACTIVE_DIRECTORY_USN - For Active Directory only, this algorithm queries for update sequence numbers on user records that are larger than the last time records were checked. \nTIMESTAMP - Queries for timestamps on user records that are not older than the last time records were checked. This check is more efficient from the point of view of the PingFederate provisioner but can be more time consuming on the LDAP side, particularly with the Oracle Directory Server. \nTIMESTAMP_NO_NEGATION - Queries for timestamps on user records that are newer than the last time records were checked. This algorithm is recommended for the Oracle Directory Server.",
												},
												"group_object_class": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The group object class.",
												},
												"time_stamp_attribute_name": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The timestamp attribute name.",
												},
												"user_object_class": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The user object class.",
												},
												"usn_attribute_name": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The USN attribute name.",
												},
											},
											Computed:    true,
											Optional:    false,
											Description: "Setting to detect changes to a user or a group.",
										},
										"data_source": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Reference to an LDAP datastore."),
										"group_membership_detection": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"group_member_attribute_name": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The name of the attribute that represents group members in a group, also known as group member attribute.",
												},
												"member_of_group_attribute_name": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The name of the attribute that indicates the entity is a member of a group, also known as member of attribute.",
												},
											},
											Computed:    true,
											Optional:    false,
											Description: "Settings to detect group memberships.",
										},
										"group_source_location": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"filter": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "An LDAP filter.",
												},
												"group_dn": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The group DN for users or groups.",
												},
												"nested_search": schema.BoolAttribute{
													Computed:    true,
													Optional:    false,
													Description: "Indicates whether the search is nested.",
												},
											},
											Computed:    true,
											Optional:    false,
											Description: "The location settings that includes a DN and a LDAP filter.",
										},
										"guid_attribute_name": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "the GUID attribute name.",
										},
										"guid_binary": schema.BoolAttribute{
											Computed:    true,
											Optional:    false,
											Description: "Indicates whether the GUID is stored in binary format.",
										},
										"user_source_location": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"filter": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "An LDAP filter.",
												},
												"group_dn": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The group DN for users or groups.",
												},
												"nested_search": schema.BoolAttribute{
													Computed:    true,
													Optional:    false,
													Description: "Indicates whether the search is nested.",
												},
											},
											Computed:    true,
											Optional:    false,
											Description: "The location settings that includes a DN and a LDAP filter.",
										},
									},
									Computed:    true,
									Optional:    false,
									Description: "The source data source and LDAP settings.",
								},
								"max_threads": schema.Int64Attribute{
									Computed:    true,
									Optional:    false,
									Description: "The number of processing threads.",
								},
								"name": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "The name of the channel.",
								},
								"timeout": schema.Int64Attribute{
									Computed:    true,
									Optional:    false,
									Description: "Timeout, in seconds, for individual user and group provisioning operations on the target service provider.",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "Includes settings of a source data store, managing provisioning threads and mapping of attributes.",
					},
					"custom_schema": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attributes": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"multi_valued": schema.BoolAttribute{
											Computed:    true,
											Optional:    false,
											Description: "Indicates whether the attribute is multi-valued.",
										},
										"name": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "Name of the attribute.",
										},
										"sub_attributes": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
											Optional:    false,
											Description: "List of sub-attributes for an attribute.",
										},
										"types": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
											Optional:    false,
											Description: "Represents the name of each attribute type in case of multi-valued attribute.",
										},
									},
								},
								Computed: true,
								Optional: false,
							},
							"namespace": schema.StringAttribute{
								Computed: true,
								Optional: false,
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "Custom SCIM Attributes configuration.",
					},
					"target_settings": schema.ListNestedAttribute{
						NestedObject: outboundProvisionTargetSettingsNestedObject,
						Computed:     true,
						Optional:     false,
						Description:  "Configuration fields that includes credentials to target SaaS application.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "The SaaS plugin type.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "Outbound Provisioning allows an IdP to create and maintain user accounts at standards-based partner sites using SCIM as well as select-proprietary provisioning partner sites that are protocol-enabled.",
			},
			"sp_browser_sso": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"adapter_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used.",
								},
								"adapter_override_settings": schema.SingleNestedAttribute{
									Attributes: map[string]schema.Attribute{
										"attribute_contract": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"core_attributes": schema.ListNestedAttribute{
													NestedObject: adapterOverrideSettingsAttribute,
													Computed:     true,
													Optional:     false,
													Description:  "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
												},
												"extended_attributes": schema.ListNestedAttribute{
													NestedObject: adapterOverrideSettingsAttribute,
													Computed:     true,
													Optional:     false,
													Description:  "A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them.",
												},
												"mask_ognl_values": schema.BoolAttribute{
													Computed:    true,
													Optional:    false,
													Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs.",
												},
												"unique_user_key_attribute": schema.StringAttribute{
													Computed:    true,
													Optional:    false,
													Description: "The attribute to use for uniquely identify a user's authentication sessions.",
												},
											},
											Computed:    true,
											Optional:    false,
											Description: "A set of attributes exposed by an IdP adapter.",
										},
										"attribute_mapping": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"attribute_contract_fulfillment": datasourceattributecontractfulfillment.ToDataSourceSchema(),
												"attribute_sources":              datasourceattributesources.ToDataSourceSchema(),
												"issuance_criteria":              datasourceissuancecriteria.ToDataSourceSchema(),
											},
											Computed:    true,
											Optional:    false,
											Description: "An IdP Adapter Contract Mapping.",
										},
										"authn_ctx_class_ref": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "The fixed value that indicates how the user was authenticated.",
										},
										"configuration": pluginconfiguration.ToSchema(),
										"id": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
										},
										"name": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
										},
										"parent_ref":            datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)"),
										"plugin_descriptor_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override."),
									},
									Computed: true,
									Optional: false,
								},
								"attribute_contract_fulfillment": datasourceattributecontractfulfillment.ToDataSourceSchema(),
								"attribute_sources":              datasourceattributesources.ToDataSourceSchema(),
								"idp_adapter_ref":                datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Reference to the associated IdP adapter. Note: This is ignored if adapter overrides for this mapping exists. In this case, the override's parent adapter reference is used."),
								"issuance_criteria":              datasourceissuancecriteria.ToDataSourceSchema(),
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Restricts this mapping to specific virtual entity IDs.",
								},
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Optional:    false,
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "A list of adapters that map to outgoing assertions.",
					},
					"always_sign_artifact_response": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Specify to always sign the SAML ArtifactResponse.",
					},
					"artifact": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"lifetime": schema.Int64Attribute{
								Computed:    true,
								Optional:    false,
								Description: "The lifetime of the artifact in seconds.",
							},
							"resolver_locations": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"index": schema.Int64Attribute{
											Computed:    true,
											Optional:    false,
											Description: "The priority of the endpoint.",
										},
										"url": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
										},
									},
								},
								Computed:    true,
								Optional:    false,
								Description: "Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message",
							},
							"source_id": schema.StringAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Source ID for SAML1.x connections",
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "The settings for an Artifact binding.",
					},
					"assertion_lifetime": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"minutes_after": schema.Int64Attribute{
								Computed:    true,
								Optional:    false,
								Description: "Assertion validity in minutes after the assertion issuance.",
							},
							"minutes_before": schema.Int64Attribute{
								Computed:    true,
								Optional:    false,
								Description: "Assertion validity in minutes before the assertion issuance.",
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "The timeframe of validity before and after the issuance of the assertion.",
					},
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.ListNestedAttribute{
								NestedObject: spBrowserSSOAttribute,
								Computed:     true,
								Optional:     false,
								Description:  "A list of read-only assertion attributes (for example, SAML_SUBJECT) that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.ListNestedAttribute{
								NestedObject: spBrowserSSOAttribute,
								Computed:     true,
								Optional:     false,
								Description:  "A list of additional attributes that are added to the outgoing assertion.",
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "A set of user attributes that the IdP sends in the SAML assertion.",
					},
					"authentication_policy_contract_assertion_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used. By default, this value is false.",
								},
								"attribute_contract_fulfillment":     datasourceattributecontractfulfillment.ToDataSourceSchema(),
								"attribute_sources":                  datasourceattributesources.ToDataSourceSchema(),
								"authentication_policy_contract_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Reference to the associated Authentication Policy Contract."),
								"issuance_criteria":                  datasourceissuancecriteria.ToDataSourceSchema(),
								"restrict_virtual_entity_ids": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Restricts this mapping to specific virtual entity IDs.",
								},
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Optional:    false,
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "A list of authentication policy contracts that map to outgoing assertions.",
					},
					"default_target_url": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Default Target URL for SAML1.x connections.",
					},
					"enabled_profiles": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Optional:    false,
						Description: "The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO.",
					},
					"encryption_policy": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"encrypt_assertion": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Whether the outgoing SAML assertion will be encrypted.",
							},
							"encrypt_slo_subject_name_id": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Encrypt the name-identifier attribute in outbound SLO messages.",
							},
							"encrypted_attributes": schema.ListAttribute{
								ElementType: types.StringType,
								Computed:    true,
								Optional:    false,
								Description: "The list of outgoing SAML assertion attributes that will be encrypted.",
							},
							"slo_subject_name_id_encrypted": schema.BoolAttribute{
								Computed:    true,
								Optional:    false,
								Description: "Allow the encryption of the name-identifier attribute for inbound SLO messages.",
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "Defines what to encrypt in the browser-based SSO profile.",
					},
					"incoming_bindings": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Optional:    false,
						Description: "The SAML bindings that are enabled for browser-based SSO.",
					},
					"message_customizations": schema.ListNestedAttribute{
						NestedObject: messageCustomizationsNestedObject,
						Computed:     true,
						Optional:     false,
						Description:  "The message customizations for browser-based SSO.",
					},
					"protocol": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "The browser-based SSO protocol to use.",
					},
					"require_signed_authn_requests": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Require AuthN requests to be signed when received via the POST or Redirect bindings.",
					},
					"sign_assertions": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Always sign the SAML Assertion.",
					},
					"sign_response_as_required": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Sign SAML Response as required by the associated binding and encryption policy.",
					},
					"slo_service_endpoints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.",
								},
								"response_url": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "The absolute or relative URL to which logout responses are sent.",
								},
								"url": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "The absolute or relative URL of the endpoint.",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "A list of possible endpoints to send SLO requests and responses.",
					},
					"sp_saml_identity_mapping": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Process in which users authenticated by the IdP are associated with user accounts local to the SP.",
					},
					"sp_ws_fed_identity_mapping": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "Process in which users authenticated by the IdP are associated with user accounts local to the SP for WS-Federation connection types.",
					},
					"sso_service_endpoints": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"binding": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints.",
								},
								"index": schema.Int64Attribute{
									Computed:    true,
									Optional:    false,
									Description: "The priority of the endpoint.",
								},
								"is_default": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Whether or not this endpoint is the default endpoint.",
								},
								"url": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "The absolute or relative URL of the endpoint.",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "A list of possible endpoints to send assertions to.",
					},
					"url_whitelist_entries": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"allow_query_and_fragment": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Allow Any Query/Fragment",
								},
								"require_https": schema.BoolAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Require HTTPS",
								},
								"valid_domain": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Valid Domain Name (leading wildcard '*.' allowed)",
								},
								"valid_path": schema.StringAttribute{
									Computed:    true,
									Optional:    false,
									Description: "Valid Path (leave blank to allow any path)",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "For WS-Federation connections, a whitelist of additional allowed domains and paths used to validate wreply for SLO, if enabled.",
					},
					"ws_fed_token_type": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "The WS-Federation Token Type to use.",
					},
					"ws_trust_version": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "The WS-Trust version for a WS-Federation connection.",
					},
					"sso_application_endpoint": schema.StringAttribute{
						Optional:    false,
						Computed:    true,
						Description: "Application endpoint that can be used to invoke single sign-on (SSO) for the connection. This is a read-only parameter. Supported in PF version 11.3 or later.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "The SAML settings used to enable secure browser-based SSO to resources at your partner's site.",
			},
			"type": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The type of this connection.",
			},
			"virtual_entity_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    false,
				Description: "List of alternate entity IDs that identifies the local server to this partner.",
			},
			"ws_trust": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"abort_if_not_fulfilled_from_request": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "If the attribute contract cannot be fulfilled using data from the Request, abort the transaction.",
					},
					"attribute_contract": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"core_attributes": schema.ListNestedAttribute{
								NestedObject: wsTrustAttribute,
								Computed:     true,
								Optional:     false,
								Description:  "A list of read-only assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.ListNestedAttribute{
								NestedObject: wsTrustAttribute,
								Computed:     true,
								Optional:     false,
								Description:  "A list of additional attributes that are added to the outgoing assertion.",
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "A set of user attributes that this server will send in the token.",
					},
					"default_token_type": schema.StringAttribute{
						Computed:    true,
						Optional:    false,
						Description: "The default token type when a web service client (WSC) does not specify in the token request which token type the STS should issue. Defaults to SAML 2.0.",
					},
					"encrypt_saml2_assertion": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "When selected, the STS encrypts the SAML 2.0 assertion. Applicable only to SAML 2.0 security token.  This option does not apply to OAuth assertion profiles.",
					},
					"generate_key": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "When selected, the STS generates a symmetric key to be used in conjunction with the \"Holder of Key\" (HoK) designation for the assertion's Subject Confirmation Method.  This option does not apply to OAuth assertion profiles.",
					},
					"message_customizations": schema.ListNestedAttribute{
						NestedObject: messageCustomizationsNestedObject,
						Computed:     true,
						Optional:     false,
						Description:  "The message customizations for WS-Trust. Depending on server settings, connection type, and protocol this may or may not be supported.",
					},
					"minutes_after": schema.Int64Attribute{
						Computed:    true,
						Optional:    false,
						Description: "The amount of time after the SAML token was issued during which it is to be considered valid.",
					},
					"minutes_before": schema.Int64Attribute{
						Computed:    true,
						Optional:    false,
						Description: "The amount of time before the SAML token was issued during which it is to be considered valid.",
					},
					"oauth_assertion_profiles": schema.BoolAttribute{
						Computed:    true,
						Optional:    false,
						Description: "When selected, four additional token-type requests become available.",
					},
					"partner_service_ids": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Optional:    false,
						Description: "The partner service identifiers.",
					},
					"request_contract_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Request Contract to be used to map attribute values into the security token."),
					"token_processor_mappings": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_contract_fulfillment": datasourceattributecontractfulfillment.ToDataSourceSchema(),
								"attribute_sources":              datasourceattributesources.ToDataSourceSchema(),
								"idp_token_processor_ref":        datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Reference to the associated token processor."),
								"issuance_criteria":              datasourceissuancecriteria.ToDataSourceSchema(),
								"restricted_virtual_entity_ids": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Optional:    false,
									Description: "The list of virtual server IDs that this mapping is restricted to.",
								},
							},
						},
						Computed:    true,
						Optional:    false,
						Description: "A list of token processors to validate incoming tokens.",
					},
				},
				Computed:    true,
				Optional:    false,
				Description: "Ws-Trust STS provides security-token validation and creation to extend SSO access to identity-enabled Web Services",
			},
		},
	}

	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"connection_id",
		true,
		"The persistent, unique ID for the connection.")
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *idpSpConnectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_sp_connection"
}

func (r *idpSpConnectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func readIdpSpconnectionDataSourceResponse(ctx context.Context, r *client.SpConnection, state *idpSpConnectionModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	diags.Append(readIdpSpconnectionResponseCommon(ctx, r, state)...)

	state.AttributeQuery, respDiags = types.ObjectValueFrom(ctx, attributeQueryAttrTypes, r.AttributeQuery)
	diags.Append(respDiags...)

	state.OutboundProvision, respDiags = types.ObjectValueFrom(ctx, outboundProvisionDataSourceAttrTypes, r.OutboundProvision)
	diags.Append(respDiags...)

	return diags
}

func (r *idpSpConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state idpSpConnectionModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpSpconnection, httpResp, err := r.apiClient.IdpSpConnectionsAPI.GetSpConnection(config.AuthContext(ctx, r.providerConfig), state.ConnectionId.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the IdP SP Connection", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readIdpSpconnectionDataSourceResponse(ctx, apiReadIdpSpconnection, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
