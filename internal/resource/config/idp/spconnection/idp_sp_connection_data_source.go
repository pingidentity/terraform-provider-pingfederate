package idpspconnection

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	datasourceattributecontractfulfillment "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributecontractfulfillment"
	datasourceattributesources "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/connectioncert"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	datasourceissuancecriteria "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/issuancecriteria"
	datasourceresourcelink "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
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
		"target_settings": types.SetType{ElemType: targetSettingsDataSourceElemAttrType},
		"custom_schema":   types.ObjectType{AttrTypes: customSchemaAttrTypes},
		"channels":        types.ListType{ElemType: channelsElemDataSourceAttrType},
	}

	credentialsInboundBackChannelAuthHttpBasicCredentialsDataSourceAttrTypes = map[string]attr.Type{
		"encrypted_password": types.StringType,
		"username":           types.StringType,
	}
	credentialsInboundBackChannelAuthDataSourceAttrTypes = map[string]attr.Type{
		"certs":                   types.ListType{ElemType: types.ObjectType{AttrTypes: connectioncert.AttrTypesDataSource()}},
		"digital_signature":       types.BoolType,
		"http_basic_credentials":  types.ObjectType{AttrTypes: credentialsInboundBackChannelAuthHttpBasicCredentialsDataSourceAttrTypes},
		"require_ssl":             types.BoolType,
		"type":                    types.StringType,
		"verification_issuer_dn":  types.StringType,
		"verification_subject_dn": types.StringType,
	}
	credentialsOutboundBackChannelAuthHttpBasicCredentialsDataSourceAttrTypes = map[string]attr.Type{
		"encrypted_password": types.StringType,
		"username":           types.StringType,
	}

	credentialsOutboundBackChannelAuthDataSourceAttrTypes = map[string]attr.Type{
		"digital_signature":      types.BoolType,
		"http_basic_credentials": types.ObjectType{AttrTypes: credentialsOutboundBackChannelAuthHttpBasicCredentialsDataSourceAttrTypes},
		"ssl_auth_key_pair_ref":  types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"type":                   types.StringType,
		"validate_partner_cert":  types.BoolType,
	}

	credentialsDataSourceAttrTypes = map[string]attr.Type{
		"block_encryption_algorithm":        types.StringType,
		"certs":                             types.ListType{ElemType: types.ObjectType{AttrTypes: connectioncert.AttrTypesDataSource()}},
		"decryption_key_pair_ref":           types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"inbound_back_channel_auth":         types.ObjectType{AttrTypes: credentialsInboundBackChannelAuthDataSourceAttrTypes},
		"key_transport_algorithm":           types.StringType,
		"outbound_back_channel_auth":        types.ObjectType{AttrTypes: credentialsOutboundBackChannelAuthDataSourceAttrTypes},
		"secondary_decryption_key_pair_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"signing_settings":                  types.ObjectType{AttrTypes: credentialsSigningSettingsAttrTypes},
		"verification_issuer_dn":            types.StringType,
		"verification_subject_dn":           types.StringType,
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
	httpBasicCredentialsSchema := schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"encrypted_password": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Description: "For GET requests, this field contains the encrypted password, if one exists.",
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
					"additional_allowed_entities": schema.SetNestedAttribute{
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
					"attributes": schema.SetAttribute{
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
					"certs":                   connectioncert.ToSchemaDataSource("The certificates used for signature verification and XML encryption."),
					"decryption_key_pair_ref": resourcelink.SingleNestedAttribute(),
					"inbound_back_channel_auth": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"certs": connectioncert.ToSchemaDataSource("The certificates used for signature verification and XML encryption."),
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
								Description: "If `verification_subject_dn` is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
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
							"alternative_signing_key_pair_refs": schema.SetNestedAttribute{
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
						Description: "If `verification_subject_dn` is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.",
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
						"values": schema.SetAttribute{
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
							"attributes": schema.SetNestedAttribute{
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
										"sub_attributes": schema.SetAttribute{
											ElementType: types.StringType,
											Computed:    true,
											Optional:    false,
											Description: "List of sub-attributes for an attribute.",
										},
										"types": schema.SetAttribute{
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
					"target_settings": schema.SetNestedAttribute{
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
					"adapter_mappings": schema.SetNestedAttribute{
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
												"core_attributes": schema.SetNestedAttribute{
													NestedObject: adapterOverrideSettingsAttribute,
													Computed:     true,
													Optional:     false,
													Description:  "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
												},
												"extended_attributes": schema.SetNestedAttribute{
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
											Description: "The ID of the plugin instance.",
										},
										"name": schema.StringAttribute{
											Computed:    true,
											Optional:    false,
											Description: "The plugin instance name.",
										},
										"parent_ref":            datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)"),
										"plugin_descriptor_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Reference to the plugin descriptor for this instance."),
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
								"restricted_virtual_entity_ids": schema.SetAttribute{
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
							"resolver_locations": schema.SetNestedAttribute{
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
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: spBrowserSSOAttribute,
								Computed:     true,
								Optional:     false,
								Description:  "A list of read-only assertion attributes (for example, SAML_SUBJECT) that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.SetNestedAttribute{
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
					"authentication_policy_contract_assertion_mappings": schema.SetNestedAttribute{
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
								"restricted_virtual_entity_ids": schema.SetAttribute{
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
					"enabled_profiles": schema.SetAttribute{
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
							"encrypted_attributes": schema.SetAttribute{
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
					"incoming_bindings": schema.SetAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Optional:    false,
						Description: "The SAML bindings that are enabled for browser-based SSO.",
					},
					"message_customizations": schema.SetNestedAttribute{
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
					"slo_service_endpoints": schema.SetNestedAttribute{
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
					"sso_service_endpoints": schema.SetNestedAttribute{
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
					"url_whitelist_entries": schema.SetNestedAttribute{
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
							"core_attributes": schema.SetNestedAttribute{
								NestedObject: wsTrustAttribute,
								Computed:     true,
								Optional:     false,
								Description:  "A list of read-only assertion attributes that are automatically populated by PingFederate.",
							},
							"extended_attributes": schema.SetNestedAttribute{
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
					"message_customizations": schema.SetNestedAttribute{
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
					"partner_service_ids": schema.SetAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Optional:    false,
						Description: "The partner service identifiers.",
					},
					"request_contract_ref": datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Request Contract to be used to map attribute values into the security token."),
					"token_processor_mappings": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"attribute_contract_fulfillment": datasourceattributecontractfulfillment.ToDataSourceSchema(),
								"attribute_sources":              datasourceattributesources.ToDataSourceSchema(),
								"idp_token_processor_ref":        datasourceresourcelink.ToDataSourceSchemaSingleNestedAttributeCustomDescription("Reference to the associated token processor."),
								"issuance_criteria":              datasourceissuancecriteria.ToDataSourceSchema(),
								"restricted_virtual_entity_ids": schema.SetAttribute{
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
	var respDiags, diags diag.Diagnostics
	// active
	state.Active = types.BoolPointerValue(r.Active)
	// additional_allowed_entities_configuration
	additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesAttrTypes := map[string]attr.Type{
		"entity_description": types.StringType,
		"entity_id":          types.StringType,
	}
	additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesElementType := types.ObjectType{AttrTypes: additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesAttrTypes}
	additionalAllowedEntitiesConfigurationAttrTypes := map[string]attr.Type{
		"additional_allowed_entities": types.SetType{ElemType: additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesElementType},
		"allow_additional_entities":   types.BoolType,
		"allow_all_entities":          types.BoolType,
	}
	var additionalAllowedEntitiesConfigurationValue types.Object
	if r.AdditionalAllowedEntitiesConfiguration == nil {
		additionalAllowedEntitiesConfigurationValue = types.ObjectNull(additionalAllowedEntitiesConfigurationAttrTypes)
	} else {
		var additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues []attr.Value
		for _, additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesResponseValue := range r.AdditionalAllowedEntitiesConfiguration.AdditionalAllowedEntities {
			additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue, diags := types.ObjectValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesAttrTypes, map[string]attr.Value{
				"entity_description": types.StringPointerValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesResponseValue.EntityDescription),
				"entity_id":          types.StringPointerValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesResponseValue.EntityId),
			})
			respDiags.Append(diags...)
			additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues = append(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues, additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue)
		}
		additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue, diags := types.SetValue(additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesElementType, additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValues)
		respDiags.Append(diags...)
		additionalAllowedEntitiesConfigurationValue, diags = types.ObjectValue(additionalAllowedEntitiesConfigurationAttrTypes, map[string]attr.Value{
			"additional_allowed_entities": additionalAllowedEntitiesConfigurationAdditionalAllowedEntitiesValue,
			"allow_additional_entities":   types.BoolPointerValue(r.AdditionalAllowedEntitiesConfiguration.AllowAdditionalEntities),
			"allow_all_entities":          types.BoolPointerValue(r.AdditionalAllowedEntitiesConfiguration.AllowAllEntities),
		})
		respDiags.Append(diags...)
	}

	state.AdditionalAllowedEntitiesConfiguration = additionalAllowedEntitiesConfigurationValue
	// application_icon_url
	state.ApplicationIconUrl = types.StringPointerValue(r.ApplicationIconUrl)
	// application_name
	state.ApplicationName = types.StringPointerValue(r.ApplicationName)
	// attribute_query
	attributeQueryAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	attributeQueryAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: attributeQueryAttributeContractFulfillmentAttrTypes}
	attributeQueryAttributeSourcesAttrTypes := attributesources.AttrTypes()
	attributeQueryAttributeSourcesElementType := types.ObjectType{AttrTypes: attributeQueryAttributeSourcesAttrTypes}
	attributeQueryIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	attributeQueryPolicyAttrTypes := map[string]attr.Type{
		"encrypt_assertion":              types.BoolType,
		"require_encrypted_name_id":      types.BoolType,
		"require_signed_attribute_query": types.BoolType,
		"sign_assertion":                 types.BoolType,
		"sign_response":                  types.BoolType,
	}
	attributeQueryAttrTypes := map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{ElemType: attributeQueryAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: attributeQueryAttributeSourcesElementType},
		"attributes":                     types.SetType{ElemType: types.StringType},
		"issuance_criteria":              types.ObjectType{AttrTypes: attributeQueryIssuanceCriteriaAttrTypes},
		"policy":                         types.ObjectType{AttrTypes: attributeQueryPolicyAttrTypes},
	}
	var attributeQueryValue types.Object
	if r.AttributeQuery == nil {
		attributeQueryValue = types.ObjectNull(attributeQueryAttrTypes)
	} else {
		attributeQueryAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &r.AttributeQuery.AttributeContractFulfillment)
		respDiags.Append(diags...)
		attributeQueryAttributeSourcesValue, diags := attributesources.ToState(context.Background(), r.AttributeQuery.AttributeSources)
		respDiags.Append(diags...)
		attributeQueryAttributesValue, diags := types.SetValueFrom(context.Background(), types.StringType, r.AttributeQuery.Attributes)
		respDiags.Append(diags...)
		attributeQueryIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), r.AttributeQuery.IssuanceCriteria)
		respDiags.Append(diags...)
		var attributeQueryPolicyValue types.Object
		if r.AttributeQuery.Policy == nil {
			attributeQueryPolicyValue = types.ObjectNull(attributeQueryPolicyAttrTypes)
		} else {
			attributeQueryPolicyValue, diags = types.ObjectValue(attributeQueryPolicyAttrTypes, map[string]attr.Value{
				"encrypt_assertion":              types.BoolPointerValue(r.AttributeQuery.Policy.EncryptAssertion),
				"require_encrypted_name_id":      types.BoolPointerValue(r.AttributeQuery.Policy.RequireEncryptedNameId),
				"require_signed_attribute_query": types.BoolPointerValue(r.AttributeQuery.Policy.RequireSignedAttributeQuery),
				"sign_assertion":                 types.BoolPointerValue(r.AttributeQuery.Policy.SignAssertion),
				"sign_response":                  types.BoolPointerValue(r.AttributeQuery.Policy.SignResponse),
			})
			respDiags.Append(diags...)
		}
		attributeQueryValue, diags = types.ObjectValue(attributeQueryAttrTypes, map[string]attr.Value{
			"attribute_contract_fulfillment": attributeQueryAttributeContractFulfillmentValue,
			"attribute_sources":              attributeQueryAttributeSourcesValue,
			"attributes":                     attributeQueryAttributesValue,
			"issuance_criteria":              attributeQueryIssuanceCriteriaValue,
			"policy":                         attributeQueryPolicyValue,
		})
		respDiags.Append(diags...)
	}

	state.AttributeQuery = attributeQueryValue
	// base_url
	state.BaseUrl = types.StringPointerValue(r.BaseUrl)
	// connection_id
	state.ConnectionId = types.StringPointerValue(r.Id)
	state.Id = types.StringPointerValue(r.Id)
	// connection_target_type
	state.ConnectionTargetType = types.StringPointerValue(r.ConnectionTargetType)
	// contact_info
	contactInfoAttrTypes := map[string]attr.Type{
		"company":    types.StringType,
		"email":      types.StringType,
		"first_name": types.StringType,
		"last_name":  types.StringType,
		"phone":      types.StringType,
	}
	var contactInfoValue types.Object
	if r.ContactInfo == nil {
		contactInfoValue = types.ObjectNull(contactInfoAttrTypes)
	} else {
		contactInfoValue, diags = types.ObjectValue(contactInfoAttrTypes, map[string]attr.Value{
			"company":    types.StringPointerValue(r.ContactInfo.Company),
			"email":      types.StringPointerValue(r.ContactInfo.Email),
			"first_name": types.StringPointerValue(r.ContactInfo.FirstName),
			"last_name":  types.StringPointerValue(r.ContactInfo.LastName),
			"phone":      types.StringPointerValue(r.ContactInfo.Phone),
		})
		respDiags.Append(diags...)
	}

	state.ContactInfo = contactInfoValue
	// creation_date
	state.CreationDate = types.StringValue(r.CreationDate.Format(time.RFC3339))
	// credentials
	var credentialsCertsValues []attr.Value
	for _, cert := range r.Credentials.Certs {
		credentailsCert, respDiags := connectioncert.ToStateDataSource(ctx, cert, &diags)
		diags.Append(respDiags...)
		credentialsCertsValues = append(credentialsCertsValues, credentailsCert)
	}
	credentialsCerts, respDiags := types.ListValue(types.ObjectType{AttrTypes: connectioncert.AttrTypesDataSource()}, credentialsCertsValues)
	diags.Append(respDiags...)
	var decryptionKeyPairRef types.Object
	if r.Credentials.DecryptionKeyPairRef == nil {
		decryptionKeyPairRef = types.ObjectNull(resourcelink.AttrType())
	} else {
		decryptionKeyPairRef, respDiags = types.ObjectValueFrom(ctx, resourcelink.AttrType(), r.Credentials.DecryptionKeyPairRef)
		diags.Append(respDiags...)
	}
	var inboundBackChannelAuth types.Object
	if r.Credentials.InboundBackChannelAuth == nil {
		inboundBackChannelAuth = types.ObjectNull(credentialsInboundBackChannelAuthDataSourceAttrTypes)
	} else {
		inboundBackChannelAuth, respDiags = types.ObjectValueFrom(ctx, credentialsInboundBackChannelAuthDataSourceAttrTypes, r.Credentials.InboundBackChannelAuth)
		diags.Append(respDiags...)
	}
	var outboundBackChannelAuth types.Object
	if r.Credentials.OutboundBackChannelAuth == nil {
		outboundBackChannelAuth = types.ObjectNull(credentialsOutboundBackChannelAuthDataSourceAttrTypes)
	} else {
		outboundBackChannelAuth, respDiags = types.ObjectValueFrom(ctx, credentialsOutboundBackChannelAuthDataSourceAttrTypes, r.Credentials.OutboundBackChannelAuth)
		diags.Append(respDiags...)
	}
	var secondaryDecryptionKeyPairRef types.Object
	if r.Credentials.SecondaryDecryptionKeyPairRef == nil {
		secondaryDecryptionKeyPairRef = types.ObjectNull(resourcelink.AttrType())
	} else {
		secondaryDecryptionKeyPairRef, respDiags = types.ObjectValueFrom(ctx, resourcelink.AttrType(), r.Credentials.SecondaryDecryptionKeyPairRef)
		diags.Append(respDiags...)
	}
	var signingSettings types.Object
	if r.Credentials.SigningSettings == nil {
		signingSettings = types.ObjectNull(credentialsSigningSettingsAttrTypes)
	} else {
		signingSettings, respDiags = types.ObjectValueFrom(ctx, credentialsSigningSettingsAttrTypes, r.Credentials.SigningSettings)
		diags.Append(respDiags...)
	}

	if r.Credentials != nil && r.Credentials.SigningSettings != nil && r.Credentials.SigningSettings.IncludeCertInSignature == nil {
		// PF returns false for include_cert_in_signature as nil. If nil is returned, just set it to false
		signingSettingsAttrs := signingSettings.Attributes()
		signingSettingsAttrs["include_cert_in_signature"] = types.BoolValue(false)
		signingSettings, respDiags = types.ObjectValue(credentialsSigningSettingsAttrTypes, signingSettingsAttrs)
		diags.Append(respDiags...)
	}

	state.Credentials, respDiags = types.ObjectValue(credentialsDataSourceAttrTypes, map[string]attr.Value{
		"block_encryption_algorithm":        types.StringPointerValue(r.Credentials.BlockEncryptionAlgorithm),
		"certs":                             credentialsCerts,
		"decryption_key_pair_ref":           decryptionKeyPairRef,
		"inbound_back_channel_auth":         inboundBackChannelAuth,
		"key_transport_algorithm":           types.StringPointerValue(r.Credentials.KeyTransportAlgorithm),
		"outbound_back_channel_auth":        outboundBackChannelAuth,
		"secondary_decryption_key_pair_ref": secondaryDecryptionKeyPairRef,
		"signing_settings":                  signingSettings,
		"verification_issuer_dn":            types.StringPointerValue(r.Credentials.VerificationIssuerDN),
		"verification_subject_dn":           types.StringPointerValue(r.Credentials.VerificationSubjectDN),
	})
	diags.Append(respDiags...)
	// default_virtual_entity_id
	state.DefaultVirtualEntityId = types.StringPointerValue(r.DefaultVirtualEntityId)
	// entity_id
	state.EntityId = types.StringValue(r.EntityId)
	// extended_properties
	extendedPropertiesAttrTypes := map[string]attr.Type{
		"values": types.SetType{ElemType: types.StringType},
	}
	extendedPropertiesElementType := types.ObjectType{AttrTypes: extendedPropertiesAttrTypes}
	var extendedPropertiesValue types.Map
	if r.ExtendedProperties == nil {
		extendedPropertiesValue = types.MapNull(extendedPropertiesElementType)
	} else {
		extendedPropertiesValues := make(map[string]attr.Value)
		for key, extendedPropertiesResponseValue := range *r.ExtendedProperties {
			extendedPropertiesValuesValue, diags := types.SetValueFrom(context.Background(), types.StringType, extendedPropertiesResponseValue.Values)
			respDiags.Append(diags...)
			extendedPropertiesValue, diags := types.ObjectValue(extendedPropertiesAttrTypes, map[string]attr.Value{
				"values": extendedPropertiesValuesValue,
			})
			respDiags.Append(diags...)
			extendedPropertiesValues[key] = extendedPropertiesValue
		}
		extendedPropertiesValue, diags = types.MapValue(extendedPropertiesElementType, extendedPropertiesValues)
		respDiags.Append(diags...)
	}

	state.ExtendedProperties = extendedPropertiesValue
	// license_connection_group
	state.LicenseConnectionGroup = types.StringPointerValue(r.LicenseConnectionGroup)
	// logging_mode
	// If the plan logging mode does not match the state logging mode, report that the error might be being controlled
	// by the `server_settings_general` resource
	if r.LoggingMode != nil && state.LoggingMode.ValueString() != *r.LoggingMode {
		diags.AddAttributeError(path.Root("logging_mode"), providererror.ConflictingValueReturnedError,
			"PingFederate returned a different value for `logging_mode` for this resource than was planned. "+
				"If `sp_connection_transaction_logging_override` is configured to anything other than `DONT_OVERRIDE` in the `server_settings_general` resource,"+
				" `logging_mode` should be configured to the same value in this resource.")
	}
	state.LoggingMode = types.StringPointerValue(r.LoggingMode)
	// metadata_reload_settings
	metadataReloadSettingsMetadataUrlRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	metadataReloadSettingsAttrTypes := map[string]attr.Type{
		"enable_auto_metadata_update": types.BoolType,
		"metadata_url_ref":            types.ObjectType{AttrTypes: metadataReloadSettingsMetadataUrlRefAttrTypes},
	}
	var metadataReloadSettingsValue types.Object
	if r.MetadataReloadSettings == nil {
		metadataReloadSettingsValue = types.ObjectNull(metadataReloadSettingsAttrTypes)
	} else {
		metadataReloadSettingsMetadataUrlRefValue, diags := types.ObjectValue(metadataReloadSettingsMetadataUrlRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(r.MetadataReloadSettings.MetadataUrlRef.Id),
		})
		respDiags.Append(diags...)
		metadataReloadSettingsValue, diags = types.ObjectValue(metadataReloadSettingsAttrTypes, map[string]attr.Value{
			"enable_auto_metadata_update": types.BoolPointerValue(r.MetadataReloadSettings.EnableAutoMetadataUpdate),
			"metadata_url_ref":            metadataReloadSettingsMetadataUrlRefValue,
		})
		respDiags.Append(diags...)
	}

	state.MetadataReloadSettings = metadataReloadSettingsValue
	// name
	state.Name = types.StringValue(r.Name)
	// outbound_provision
	outboundProvisionChannelsAttributeMappingSaasFieldInfoAttrTypes := map[string]attr.Type{
		"attribute_names": types.ListType{ElemType: types.StringType},
		"character_case":  types.StringType,
		"create_only":     types.BoolType,
		"default_value":   types.StringType,
		"expression":      types.StringType,
		"masked":          types.BoolType,
		"parser":          types.StringType,
		"trim":            types.BoolType,
	}
	outboundProvisionChannelsAttributeMappingAttrTypes := map[string]attr.Type{
		"field_name":      types.StringType,
		"saas_field_info": types.ObjectType{AttrTypes: outboundProvisionChannelsAttributeMappingSaasFieldInfoAttrTypes},
	}
	outboundProvisionChannelsAttributeMappingElementType := types.ObjectType{AttrTypes: outboundProvisionChannelsAttributeMappingAttrTypes}
	outboundProvisionChannelsChannelSourceAccountManagementSettingsAttrTypes := map[string]attr.Type{
		"account_status_algorithm":      types.StringType,
		"account_status_attribute_name": types.StringType,
		"default_status":                types.BoolType,
		"flag_comparison_status":        types.BoolType,
		"flag_comparison_value":         types.StringType,
	}
	outboundProvisionChannelsChannelSourceChangeDetectionSettingsAttrTypes := map[string]attr.Type{
		"changed_users_algorithm":   types.StringType,
		"group_object_class":        types.StringType,
		"time_stamp_attribute_name": types.StringType,
		"user_object_class":         types.StringType,
		"usn_attribute_name":        types.StringType,
	}
	outboundProvisionChannelsChannelSourceDataSourceAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	outboundProvisionChannelsChannelSourceGroupMembershipDetectionAttrTypes := map[string]attr.Type{
		"group_member_attribute_name":    types.StringType,
		"member_of_group_attribute_name": types.StringType,
	}
	outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes := map[string]attr.Type{
		"filter":        types.StringType,
		"group_dn":      types.StringType,
		"nested_search": types.BoolType,
	}
	outboundProvisionChannelsChannelSourceUserSourceLocationAttrTypes := map[string]attr.Type{
		"filter":        types.StringType,
		"group_dn":      types.StringType,
		"nested_search": types.BoolType,
	}
	outboundProvisionChannelsChannelSourceAttrTypes := map[string]attr.Type{
		"account_management_settings": types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceAccountManagementSettingsAttrTypes},
		"base_dn":                     types.StringType,
		"change_detection_settings":   types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceChangeDetectionSettingsAttrTypes},
		"data_source":                 types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceDataSourceAttrTypes},
		"group_membership_detection":  types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceGroupMembershipDetectionAttrTypes},
		"group_source_location":       types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes},
		"guid_attribute_name":         types.StringType,
		"guid_binary":                 types.BoolType,
		"user_source_location":        types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceUserSourceLocationAttrTypes},
	}
	outboundProvisionChannelsAttrTypes := map[string]attr.Type{
		"active":            types.BoolType,
		"attribute_mapping": types.SetType{ElemType: outboundProvisionChannelsAttributeMappingElementType},
		"channel_source":    types.ObjectType{AttrTypes: outboundProvisionChannelsChannelSourceAttrTypes},
		"max_threads":       types.Int64Type,
		"name":              types.StringType,
		"timeout":           types.Int64Type,
	}
	outboundProvisionChannelsElementType := types.ObjectType{AttrTypes: outboundProvisionChannelsAttrTypes}
	outboundProvisionCustomSchemaAttributesAttrTypes := map[string]attr.Type{
		"multi_valued":   types.BoolType,
		"name":           types.StringType,
		"sub_attributes": types.SetType{ElemType: types.StringType},
		"types":          types.SetType{ElemType: types.StringType},
	}
	outboundProvisionCustomSchemaAttributesElementType := types.ObjectType{AttrTypes: outboundProvisionCustomSchemaAttributesAttrTypes}
	outboundProvisionCustomSchemaAttrTypes := map[string]attr.Type{
		"attributes": types.SetType{ElemType: outboundProvisionCustomSchemaAttributesElementType},
		"namespace":  types.StringType,
	}
	outboundProvisionTargetSettingsAttrTypes := map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
	}
	outboundProvisionTargetSettingsElementType := types.ObjectType{AttrTypes: outboundProvisionTargetSettingsAttrTypes}
	outboundProvisionAttrTypes := map[string]attr.Type{
		"channels":        types.ListType{ElemType: outboundProvisionChannelsElementType},
		"custom_schema":   types.ObjectType{AttrTypes: outboundProvisionCustomSchemaAttrTypes},
		"target_settings": types.SetType{ElemType: outboundProvisionTargetSettingsElementType},
		"type":            types.StringType,
	}
	var outboundProvisionValue types.Object
	if r.OutboundProvision == nil {
		outboundProvisionValue = types.ObjectNull(outboundProvisionAttrTypes)
	} else {
		var outboundProvisionChannelsValues []attr.Value
		for _, outboundProvisionChannelsResponseValue := range r.OutboundProvision.Channels {
			var outboundProvisionChannelsAttributeMappingValues []attr.Value
			for _, outboundProvisionChannelsAttributeMappingResponseValue := range outboundProvisionChannelsResponseValue.AttributeMapping {
				outboundProvisionChannelsAttributeMappingSaasFieldInfoAttributeNamesValue, diags := types.ListValueFrom(context.Background(), types.StringType, outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.AttributeNames)
				respDiags.Append(diags...)
				outboundProvisionChannelsAttributeMappingSaasFieldInfoValue, diags := types.ObjectValue(outboundProvisionChannelsAttributeMappingSaasFieldInfoAttrTypes, map[string]attr.Value{
					"attribute_names": outboundProvisionChannelsAttributeMappingSaasFieldInfoAttributeNamesValue,
					"character_case":  types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.CharacterCase),
					"create_only":     types.BoolPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.CreateOnly),
					"default_value":   types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.DefaultValue),
					"expression":      types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Expression),
					"masked":          types.BoolPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Masked),
					"parser":          types.StringPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Parser),
					"trim":            types.BoolPointerValue(outboundProvisionChannelsAttributeMappingResponseValue.SaasFieldInfo.Trim),
				})
				respDiags.Append(diags...)
				outboundProvisionChannelsAttributeMappingValue, diags := types.ObjectValue(outboundProvisionChannelsAttributeMappingAttrTypes, map[string]attr.Value{
					"field_name":      types.StringValue(outboundProvisionChannelsAttributeMappingResponseValue.FieldName),
					"saas_field_info": outboundProvisionChannelsAttributeMappingSaasFieldInfoValue,
				})
				respDiags.Append(diags...)
				outboundProvisionChannelsAttributeMappingValues = append(outboundProvisionChannelsAttributeMappingValues, outboundProvisionChannelsAttributeMappingValue)
			}
			outboundProvisionChannelsAttributeMappingValue, diags := types.SetValue(outboundProvisionChannelsAttributeMappingElementType, outboundProvisionChannelsAttributeMappingValues)
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceAccountManagementSettingsValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceAccountManagementSettingsAttrTypes, map[string]attr.Value{
				"account_status_algorithm":      types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.AccountStatusAlgorithm),
				"account_status_attribute_name": types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.AccountStatusAttributeName),
				"default_status":                types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.DefaultStatus),
				"flag_comparison_status":        types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.FlagComparisonStatus),
				"flag_comparison_value":         types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.AccountManagementSettings.FlagComparisonValue),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceChangeDetectionSettingsValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceChangeDetectionSettingsAttrTypes, map[string]attr.Value{
				"changed_users_algorithm":   types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.ChangedUsersAlgorithm),
				"group_object_class":        types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.GroupObjectClass),
				"time_stamp_attribute_name": types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.TimeStampAttributeName),
				"user_object_class":         types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.UserObjectClass),
				"usn_attribute_name":        types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.ChangeDetectionSettings.UsnAttributeName),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceDataSourceValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceDataSourceAttrTypes, map[string]attr.Value{
				"id": types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.DataSource.Id),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceGroupMembershipDetectionValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceGroupMembershipDetectionAttrTypes, map[string]attr.Value{
				"group_member_attribute_name":    types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupMembershipDetection.GroupMemberAttributeName),
				"member_of_group_attribute_name": types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupMembershipDetection.MemberOfGroupAttributeName),
			})
			respDiags.Append(diags...)
			var outboundProvisionChannelsChannelSourceGroupSourceLocationValue types.Object
			if outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation == nil {
				outboundProvisionChannelsChannelSourceGroupSourceLocationValue = types.ObjectNull(outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes)
			} else {
				outboundProvisionChannelsChannelSourceGroupSourceLocationValue, diags = types.ObjectValue(outboundProvisionChannelsChannelSourceGroupSourceLocationAttrTypes, map[string]attr.Value{
					"filter":        types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation.Filter),
					"group_dn":      types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation.GroupDN),
					"nested_search": types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.GroupSourceLocation.NestedSearch),
				})
				respDiags.Append(diags...)
			}
			outboundProvisionChannelsChannelSourceUserSourceLocationValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceUserSourceLocationAttrTypes, map[string]attr.Value{
				"filter":        types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.UserSourceLocation.Filter),
				"group_dn":      types.StringPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.UserSourceLocation.GroupDN),
				"nested_search": types.BoolPointerValue(outboundProvisionChannelsResponseValue.ChannelSource.UserSourceLocation.NestedSearch),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsChannelSourceValue, diags := types.ObjectValue(outboundProvisionChannelsChannelSourceAttrTypes, map[string]attr.Value{
				"account_management_settings": outboundProvisionChannelsChannelSourceAccountManagementSettingsValue,
				"base_dn":                     types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.BaseDn),
				"change_detection_settings":   outboundProvisionChannelsChannelSourceChangeDetectionSettingsValue,
				"data_source":                 outboundProvisionChannelsChannelSourceDataSourceValue,
				"group_membership_detection":  outboundProvisionChannelsChannelSourceGroupMembershipDetectionValue,
				"group_source_location":       outboundProvisionChannelsChannelSourceGroupSourceLocationValue,
				"guid_attribute_name":         types.StringValue(outboundProvisionChannelsResponseValue.ChannelSource.GuidAttributeName),
				"guid_binary":                 types.BoolValue(outboundProvisionChannelsResponseValue.ChannelSource.GuidBinary),
				"user_source_location":        outboundProvisionChannelsChannelSourceUserSourceLocationValue,
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsValue, diags := types.ObjectValue(outboundProvisionChannelsAttrTypes, map[string]attr.Value{
				"active":            types.BoolValue(outboundProvisionChannelsResponseValue.Active),
				"attribute_mapping": outboundProvisionChannelsAttributeMappingValue,
				"channel_source":    outboundProvisionChannelsChannelSourceValue,
				"max_threads":       types.Int64Value(outboundProvisionChannelsResponseValue.MaxThreads),
				"name":              types.StringValue(outboundProvisionChannelsResponseValue.Name),
				"timeout":           types.Int64Value(outboundProvisionChannelsResponseValue.Timeout),
			})
			respDiags.Append(diags...)
			outboundProvisionChannelsValues = append(outboundProvisionChannelsValues, outboundProvisionChannelsValue)
		}
		outboundProvisionChannelsValue, diags := types.ListValue(outboundProvisionChannelsElementType, outboundProvisionChannelsValues)
		respDiags.Append(diags...)
		var outboundProvisionCustomSchemaValue types.Object
		if r.OutboundProvision.CustomSchema == nil {
			outboundProvisionCustomSchemaValue = types.ObjectNull(outboundProvisionCustomSchemaAttrTypes)
		} else {
			var outboundProvisionCustomSchemaAttributesValues []attr.Value
			for _, outboundProvisionCustomSchemaAttributesResponseValue := range r.OutboundProvision.CustomSchema.Attributes {
				outboundProvisionCustomSchemaAttributesSubAttributesValue, diags := types.SetValueFrom(context.Background(), types.StringType, outboundProvisionCustomSchemaAttributesResponseValue.SubAttributes)
				respDiags.Append(diags...)
				outboundProvisionCustomSchemaAttributesTypesValue, diags := types.SetValueFrom(context.Background(), types.StringType, outboundProvisionCustomSchemaAttributesResponseValue.Types)
				respDiags.Append(diags...)
				outboundProvisionCustomSchemaAttributesValue, diags := types.ObjectValue(outboundProvisionCustomSchemaAttributesAttrTypes, map[string]attr.Value{
					"multi_valued":   types.BoolPointerValue(outboundProvisionCustomSchemaAttributesResponseValue.MultiValued),
					"name":           types.StringPointerValue(outboundProvisionCustomSchemaAttributesResponseValue.Name),
					"sub_attributes": outboundProvisionCustomSchemaAttributesSubAttributesValue,
					"types":          outboundProvisionCustomSchemaAttributesTypesValue,
				})
				respDiags.Append(diags...)
				outboundProvisionCustomSchemaAttributesValues = append(outboundProvisionCustomSchemaAttributesValues, outboundProvisionCustomSchemaAttributesValue)
			}
			outboundProvisionCustomSchemaAttributesValue, diags := types.SetValue(outboundProvisionCustomSchemaAttributesElementType, outboundProvisionCustomSchemaAttributesValues)
			respDiags.Append(diags...)
			outboundProvisionCustomSchemaValue, diags = types.ObjectValue(outboundProvisionCustomSchemaAttrTypes, map[string]attr.Value{
				"attributes": outboundProvisionCustomSchemaAttributesValue,
				"namespace":  types.StringPointerValue(r.OutboundProvision.CustomSchema.Namespace),
			})
			respDiags.Append(diags...)
		}
		var outboundProvisionTargetSettingsValues []attr.Value
		for _, outboundProvisionTargetSettingsResponseValue := range r.OutboundProvision.TargetSettings {
			outboundProvisionTargetSettingsValue, diags := types.ObjectValue(outboundProvisionTargetSettingsAttrTypes, map[string]attr.Value{
				"name":            types.StringValue(outboundProvisionTargetSettingsResponseValue.Name),
				"value":           types.StringPointerValue(outboundProvisionTargetSettingsResponseValue.Value),
				"encrypted_value": types.StringPointerValue(outboundProvisionTargetSettingsResponseValue.EncryptedValue),
			})
			respDiags.Append(diags...)
			outboundProvisionTargetSettingsValues = append(outboundProvisionTargetSettingsValues, outboundProvisionTargetSettingsValue)
		}
		outboundProvisionTargetSettingsValue, diags := types.SetValue(outboundProvisionTargetSettingsElementType, outboundProvisionTargetSettingsValues)
		respDiags.Append(diags...)
		outboundProvisionValue, diags = types.ObjectValue(outboundProvisionAttrTypes, map[string]attr.Value{
			"channels":        outboundProvisionChannelsValue,
			"custom_schema":   outboundProvisionCustomSchemaValue,
			"target_settings": outboundProvisionTargetSettingsValue,
			"type":            types.StringValue(r.OutboundProvision.Type),
		})
		respDiags.Append(diags...)
	}

	state.OutboundProvision = outboundProvisionValue
	// sp_browser_sso
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"masked":    types.BoolType,
		"name":      types.StringType,
		"pseudonym": types.BoolType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"masked":    types.BoolType,
		"name":      types.StringType,
		"pseudonym": types.BoolType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":           types.SetType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType},
		"extended_attributes":       types.SetType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType},
		"mask_ognl_values":          types.BoolType,
		"unique_user_key_attribute": types.StringType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesAttrTypes := attributesources.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesAttrTypes}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes := map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesElementType},
		"issuance_criteria":              types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaAttrTypes},
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationAttrTypes := pluginconfiguration.AttrTypes()
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes := map[string]attr.Type{
		"attribute_contract":    types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes},
		"attribute_mapping":     types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes},
		"authn_ctx_class_ref":   types.StringType,
		"configuration":         types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationAttrTypes},
		"id":                    types.StringType,
		"name":                  types.StringType,
		"parent_ref":            types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes},
		"plugin_descriptor_ref": types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrTypes},
	}
	spBrowserSsoAdapterMappingsAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	spBrowserSsoAdapterMappingsAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAttributeContractFulfillmentAttrTypes}
	spBrowserSsoAdapterMappingsAttributeSourcesAttrTypes := attributesources.AttrTypes()
	spBrowserSsoAdapterMappingsAttributeSourcesElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAttributeSourcesAttrTypes}
	spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAdapterMappingsIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	spBrowserSsoAdapterMappingsAttrTypes := map[string]attr.Type{
		"abort_sso_transaction_as_fail_safe": types.BoolType,
		"adapter_override_settings":          types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes},
		"attribute_contract_fulfillment":     types.MapType{ElemType: spBrowserSsoAdapterMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":                  types.SetType{ElemType: spBrowserSsoAdapterMappingsAttributeSourcesElementType},
		"idp_adapter_ref":                    types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes},
		"issuance_criteria":                  types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsIssuanceCriteriaAttrTypes},
		"restrict_virtual_entity_ids":        types.BoolType,
		"restricted_virtual_entity_ids":      types.SetType{ElemType: types.StringType},
	}
	spBrowserSsoAdapterMappingsElementType := types.ObjectType{AttrTypes: spBrowserSsoAdapterMappingsAttrTypes}
	spBrowserSsoArtifactResolverLocationsAttrTypes := map[string]attr.Type{
		"index": types.Int64Type,
		"url":   types.StringType,
	}
	spBrowserSsoArtifactResolverLocationsElementType := types.ObjectType{AttrTypes: spBrowserSsoArtifactResolverLocationsAttrTypes}
	spBrowserSsoArtifactAttrTypes := map[string]attr.Type{
		"lifetime":           types.Int64Type,
		"resolver_locations": types.SetType{ElemType: spBrowserSsoArtifactResolverLocationsElementType},
		"source_id":          types.StringType,
	}
	spBrowserSsoAssertionLifetimeAttrTypes := map[string]attr.Type{
		"minutes_after":  types.Int64Type,
		"minutes_before": types.Int64Type,
	}
	spBrowserSsoAttributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"name":        types.StringType,
		"name_format": types.StringType,
	}
	spBrowserSsoAttributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAttributeContractCoreAttributesAttrTypes}
	spBrowserSsoAttributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"name":        types.StringType,
		"name_format": types.StringType,
	}
	spBrowserSsoAttributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: spBrowserSsoAttributeContractExtendedAttributesAttrTypes}
	spBrowserSsoAttributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: spBrowserSsoAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: spBrowserSsoAttributeContractExtendedAttributesElementType},
	}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentAttrTypes}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesAttrTypes := attributesources.AttrTypes()
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesElementType := types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesAttrTypes}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttrTypes := map[string]attr.Type{
		"abort_sso_transaction_as_fail_safe": types.BoolType,
		"attribute_contract_fulfillment":     types.MapType{ElemType: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":                  types.SetType{ElemType: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesElementType},
		"authentication_policy_contract_ref": types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrTypes},
		"issuance_criteria":                  types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaAttrTypes},
		"restrict_virtual_entity_ids":        types.BoolType,
		"restricted_virtual_entity_ids":      types.SetType{ElemType: types.StringType},
	}
	spBrowserSsoAuthenticationPolicyContractAssertionMappingsElementType := types.ObjectType{AttrTypes: spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttrTypes}
	spBrowserSsoEncryptionPolicyAttrTypes := map[string]attr.Type{
		"encrypt_assertion":             types.BoolType,
		"encrypt_slo_subject_name_id":   types.BoolType,
		"encrypted_attributes":          types.SetType{ElemType: types.StringType},
		"slo_subject_name_id_encrypted": types.BoolType,
	}
	spBrowserSsoMessageCustomizationsAttrTypes := map[string]attr.Type{
		"context_name":       types.StringType,
		"message_expression": types.StringType,
	}
	spBrowserSsoMessageCustomizationsElementType := types.ObjectType{AttrTypes: spBrowserSsoMessageCustomizationsAttrTypes}
	spBrowserSsoSloServiceEndpointsAttrTypes := map[string]attr.Type{
		"binding":      types.StringType,
		"response_url": types.StringType,
		"url":          types.StringType,
	}
	spBrowserSsoSloServiceEndpointsElementType := types.ObjectType{AttrTypes: spBrowserSsoSloServiceEndpointsAttrTypes}
	spBrowserSsoSsoServiceEndpointsAttrTypes := map[string]attr.Type{
		"binding":    types.StringType,
		"index":      types.Int64Type,
		"is_default": types.BoolType,
		"url":        types.StringType,
	}
	spBrowserSsoSsoServiceEndpointsElementType := types.ObjectType{AttrTypes: spBrowserSsoSsoServiceEndpointsAttrTypes}
	spBrowserSsoUrlWhitelistEntriesAttrTypes := map[string]attr.Type{
		"allow_query_and_fragment": types.BoolType,
		"require_https":            types.BoolType,
		"valid_domain":             types.StringType,
		"valid_path":               types.StringType,
	}
	spBrowserSsoUrlWhitelistEntriesElementType := types.ObjectType{AttrTypes: spBrowserSsoUrlWhitelistEntriesAttrTypes}
	spBrowserSsoAttrTypes := map[string]attr.Type{
		"adapter_mappings":              types.SetType{ElemType: spBrowserSsoAdapterMappingsElementType},
		"always_sign_artifact_response": types.BoolType,
		"artifact":                      types.ObjectType{AttrTypes: spBrowserSsoArtifactAttrTypes},
		"assertion_lifetime":            types.ObjectType{AttrTypes: spBrowserSsoAssertionLifetimeAttrTypes},
		"attribute_contract":            types.ObjectType{AttrTypes: spBrowserSsoAttributeContractAttrTypes},
		"authentication_policy_contract_assertion_mappings": types.SetType{ElemType: spBrowserSsoAuthenticationPolicyContractAssertionMappingsElementType},
		"default_target_url":            types.StringType,
		"enabled_profiles":              types.SetType{ElemType: types.StringType},
		"encryption_policy":             types.ObjectType{AttrTypes: spBrowserSsoEncryptionPolicyAttrTypes},
		"incoming_bindings":             types.SetType{ElemType: types.StringType},
		"message_customizations":        types.SetType{ElemType: spBrowserSsoMessageCustomizationsElementType},
		"protocol":                      types.StringType,
		"require_signed_authn_requests": types.BoolType,
		"sign_assertions":               types.BoolType,
		"sign_response_as_required":     types.BoolType,
		"slo_service_endpoints":         types.SetType{ElemType: spBrowserSsoSloServiceEndpointsElementType},
		"sp_saml_identity_mapping":      types.StringType,
		"sp_ws_fed_identity_mapping":    types.StringType,
		"sso_application_endpoint":      types.StringType,
		"sso_service_endpoints":         types.SetType{ElemType: spBrowserSsoSsoServiceEndpointsElementType},
		"url_whitelist_entries":         types.SetType{ElemType: spBrowserSsoUrlWhitelistEntriesElementType},
		"ws_fed_token_type":             types.StringType,
		"ws_trust_version":              types.StringType,
	}
	var spBrowserSsoValue types.Object
	if r.SpBrowserSso == nil {
		spBrowserSsoValue = types.ObjectNull(spBrowserSsoAttrTypes)
	} else {
		var spBrowserSsoAdapterMappingsValues []attr.Value
		for adapterMappingIndex, spBrowserSsoAdapterMappingsResponseValue := range r.SpBrowserSso.AdapterMappings {
			var spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue types.Object
			if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings == nil {
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes)
			} else {
				var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue types.Object
				if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract == nil {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes)
				} else {
					var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues []attr.Value
					for _, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue := range spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.CoreAttributes {
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
							"masked":    types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Masked),
							"name":      types.StringValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Name),
							"pseudonym": types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesResponseValue.Pseudonym),
						})
						respDiags.Append(diags...)
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues = append(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue)
					}
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue, diags := types.SetValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesElementType, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValues)
					respDiags.Append(diags...)
					var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues []attr.Value
					for _, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue := range spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.ExtendedAttributes {
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
							"masked":    types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Masked),
							"name":      types.StringValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Name),
							"pseudonym": types.BoolPointerValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesResponseValue.Pseudonym),
						})
						respDiags.Append(diags...)
						spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues = append(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue)
					}
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue, diags := types.SetValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesElementType, spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValues)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractAttrTypes, map[string]attr.Value{
						"core_attributes":           spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractCoreAttributesValue,
						"extended_attributes":       spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractExtendedAttributesValue,
						"mask_ognl_values":          types.BoolPointerValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.MaskOgnlValues),
						"unique_user_key_attribute": types.StringPointerValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeContract.UniqueUserKeyAttribute),
					})
					respDiags.Append(diags...)
				}
				var spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue types.Object
				if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping == nil {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes)
				} else {
					contractFulfillment := spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping.AttributeContractFulfillment
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesValue, diags := attributesources.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping.AttributeSources)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AttributeMapping.IssuanceCriteria)
					respDiags.Append(diags...)
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttrTypes, map[string]attr.Value{
						"attribute_contract_fulfillment": spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeContractFulfillmentValue,
						"attribute_sources":              spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingAttributeSourcesValue,
						"issuance_criteria":              spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingIssuanceCriteriaValue,
					})
					respDiags.Append(diags...)
				}
				overrideSettingsConfiguration := spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Configuration
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue, diags := pluginconfiguration.ToState(state.getSpBrowserSsoAdapterMappingsAdapterOverrideSettingsConfiguration(adapterMappingIndex), &overrideSettingsConfiguration, true)
				respDiags.Append(diags...)
				var spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue types.Object
				if spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef == nil {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue = types.ObjectNull(spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes)
				} else {
					spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefAttrTypes, map[string]attr.Value{
						"id": types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.ParentRef.Id),
					})
					respDiags.Append(diags...)
				}
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefAttrTypes, map[string]attr.Value{
					"id": types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.PluginDescriptorRef.Id),
				})
				respDiags.Append(diags...)
				spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttrTypes, map[string]attr.Value{
					"attribute_contract":    spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeContractValue,
					"attribute_mapping":     spBrowserSsoAdapterMappingsAdapterOverrideSettingsAttributeMappingValue,
					"authn_ctx_class_ref":   types.StringPointerValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.AuthnCtxClassRef),
					"configuration":         spBrowserSsoAdapterMappingsAdapterOverrideSettingsConfigurationValue,
					"id":                    types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Id),
					"name":                  types.StringValue(spBrowserSsoAdapterMappingsResponseValue.AdapterOverrideSettings.Name),
					"parent_ref":            spBrowserSsoAdapterMappingsAdapterOverrideSettingsParentRefValue,
					"plugin_descriptor_ref": spBrowserSsoAdapterMappingsAdapterOverrideSettingsPluginDescriptorRefValue,
				})
				respDiags.Append(diags...)
			}
			contractFulfillment := spBrowserSsoAdapterMappingsResponseValue.AttributeContractFulfillment
			spBrowserSsoAdapterMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			var spBrowserSsoAdapterMappingsIdpAdapterRefValue types.Object
			if spBrowserSsoAdapterMappingsResponseValue.IdpAdapterRef == nil {
				spBrowserSsoAdapterMappingsIdpAdapterRefValue = types.ObjectNull(spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes)
			} else {
				spBrowserSsoAdapterMappingsIdpAdapterRefValue, diags = types.ObjectValue(spBrowserSsoAdapterMappingsIdpAdapterRefAttrTypes, map[string]attr.Value{
					"id": types.StringValue(spBrowserSsoAdapterMappingsResponseValue.IdpAdapterRef.Id),
				})
				respDiags.Append(diags...)
			}
			spBrowserSsoAdapterMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), spBrowserSsoAdapterMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, spBrowserSsoAdapterMappingsResponseValue.RestrictedVirtualEntityIds)
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsValue, diags := types.ObjectValue(spBrowserSsoAdapterMappingsAttrTypes, map[string]attr.Value{
				"abort_sso_transaction_as_fail_safe": types.BoolPointerValue(spBrowserSsoAdapterMappingsResponseValue.AbortSsoTransactionAsFailSafe),
				"adapter_override_settings":          spBrowserSsoAdapterMappingsAdapterOverrideSettingsValue,
				"attribute_contract_fulfillment":     spBrowserSsoAdapterMappingsAttributeContractFulfillmentValue,
				"attribute_sources":                  spBrowserSsoAdapterMappingsAttributeSourcesValue,
				"idp_adapter_ref":                    spBrowserSsoAdapterMappingsIdpAdapterRefValue,
				"issuance_criteria":                  spBrowserSsoAdapterMappingsIssuanceCriteriaValue,
				"restrict_virtual_entity_ids":        types.BoolPointerValue(spBrowserSsoAdapterMappingsResponseValue.RestrictVirtualEntityIds),
				"restricted_virtual_entity_ids":      spBrowserSsoAdapterMappingsRestrictedVirtualEntityIdsValue,
			})
			respDiags.Append(diags...)
			spBrowserSsoAdapterMappingsValues = append(spBrowserSsoAdapterMappingsValues, spBrowserSsoAdapterMappingsValue)
		}
		spBrowserSsoAdapterMappingsValue, diags := types.SetValue(spBrowserSsoAdapterMappingsElementType, spBrowserSsoAdapterMappingsValues)
		respDiags.Append(diags...)
		var spBrowserSsoArtifactValue types.Object
		if r.SpBrowserSso.Artifact == nil {
			spBrowserSsoArtifactValue = types.ObjectNull(spBrowserSsoArtifactAttrTypes)
		} else {
			var spBrowserSsoArtifactResolverLocationsValues []attr.Value
			for _, spBrowserSsoArtifactResolverLocationsResponseValue := range r.SpBrowserSso.Artifact.ResolverLocations {
				spBrowserSsoArtifactResolverLocationsValue, diags := types.ObjectValue(spBrowserSsoArtifactResolverLocationsAttrTypes, map[string]attr.Value{
					"index": types.Int64Value(spBrowserSsoArtifactResolverLocationsResponseValue.Index),
					"url":   types.StringValue(spBrowserSsoArtifactResolverLocationsResponseValue.Url),
				})
				respDiags.Append(diags...)
				spBrowserSsoArtifactResolverLocationsValues = append(spBrowserSsoArtifactResolverLocationsValues, spBrowserSsoArtifactResolverLocationsValue)
			}
			spBrowserSsoArtifactResolverLocationsValue, diags := types.SetValue(spBrowserSsoArtifactResolverLocationsElementType, spBrowserSsoArtifactResolverLocationsValues)
			respDiags.Append(diags...)
			spBrowserSsoArtifactValue, diags = types.ObjectValue(spBrowserSsoArtifactAttrTypes, map[string]attr.Value{
				"lifetime":           types.Int64PointerValue(r.SpBrowserSso.Artifact.Lifetime),
				"resolver_locations": spBrowserSsoArtifactResolverLocationsValue,
				"source_id":          types.StringPointerValue(r.SpBrowserSso.Artifact.SourceId),
			})
			respDiags.Append(diags...)
		}
		spBrowserSsoAssertionLifetimeValue, diags := types.ObjectValue(spBrowserSsoAssertionLifetimeAttrTypes, map[string]attr.Value{
			"minutes_after":  types.Int64Value(r.SpBrowserSso.AssertionLifetime.MinutesAfter),
			"minutes_before": types.Int64Value(r.SpBrowserSso.AssertionLifetime.MinutesBefore),
		})
		respDiags.Append(diags...)
		var spBrowserSsoAttributeContractCoreAttributesValues []attr.Value
		for _, spBrowserSsoAttributeContractCoreAttributesResponseValue := range r.SpBrowserSso.AttributeContract.CoreAttributes {
			spBrowserSsoAttributeContractCoreAttributesValue, diags := types.ObjectValue(spBrowserSsoAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"name":        types.StringValue(spBrowserSsoAttributeContractCoreAttributesResponseValue.Name),
				"name_format": types.StringPointerValue(spBrowserSsoAttributeContractCoreAttributesResponseValue.NameFormat),
			})
			respDiags.Append(diags...)
			spBrowserSsoAttributeContractCoreAttributesValues = append(spBrowserSsoAttributeContractCoreAttributesValues, spBrowserSsoAttributeContractCoreAttributesValue)
		}
		spBrowserSsoAttributeContractCoreAttributesValue, diags := types.SetValue(spBrowserSsoAttributeContractCoreAttributesElementType, spBrowserSsoAttributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var spBrowserSsoAttributeContractExtendedAttributesValues []attr.Value
		for _, spBrowserSsoAttributeContractExtendedAttributesResponseValue := range r.SpBrowserSso.AttributeContract.ExtendedAttributes {
			spBrowserSsoAttributeContractExtendedAttributesValue, diags := types.ObjectValue(spBrowserSsoAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"name":        types.StringValue(spBrowserSsoAttributeContractExtendedAttributesResponseValue.Name),
				"name_format": types.StringPointerValue(spBrowserSsoAttributeContractExtendedAttributesResponseValue.NameFormat),
			})
			respDiags.Append(diags...)
			spBrowserSsoAttributeContractExtendedAttributesValues = append(spBrowserSsoAttributeContractExtendedAttributesValues, spBrowserSsoAttributeContractExtendedAttributesValue)
		}
		spBrowserSsoAttributeContractExtendedAttributesValue, diags := types.SetValue(spBrowserSsoAttributeContractExtendedAttributesElementType, spBrowserSsoAttributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		spBrowserSsoAttributeContractValue, diags := types.ObjectValue(spBrowserSsoAttributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     spBrowserSsoAttributeContractCoreAttributesValue,
			"extended_attributes": spBrowserSsoAttributeContractExtendedAttributesValue,
		})
		respDiags.Append(diags...)
		var spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues []attr.Value
		for _, spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue := range r.SpBrowserSso.AuthenticationPolicyContractAssertionMappings {
			contractFulfillment := spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AttributeContractFulfillment
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefValue, diags := types.ObjectValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AuthenticationPolicyContractRef.Id),
			})
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsRestrictedVirtualEntityIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.RestrictedVirtualEntityIds)
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue, diags := types.ObjectValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttrTypes, map[string]attr.Value{
				"abort_sso_transaction_as_fail_safe": types.BoolPointerValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.AbortSsoTransactionAsFailSafe),
				"attribute_contract_fulfillment":     spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeContractFulfillmentValue,
				"attribute_sources":                  spBrowserSsoAuthenticationPolicyContractAssertionMappingsAttributeSourcesValue,
				"authentication_policy_contract_ref": spBrowserSsoAuthenticationPolicyContractAssertionMappingsAuthenticationPolicyContractRefValue,
				"issuance_criteria":                  spBrowserSsoAuthenticationPolicyContractAssertionMappingsIssuanceCriteriaValue,
				"restrict_virtual_entity_ids":        types.BoolPointerValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsResponseValue.RestrictVirtualEntityIds),
				"restricted_virtual_entity_ids":      spBrowserSsoAuthenticationPolicyContractAssertionMappingsRestrictedVirtualEntityIdsValue,
			})
			respDiags.Append(diags...)
			spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues = append(spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues, spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue)
		}
		spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue, diags := types.SetValue(spBrowserSsoAuthenticationPolicyContractAssertionMappingsElementType, spBrowserSsoAuthenticationPolicyContractAssertionMappingsValues)
		respDiags.Append(diags...)
		spBrowserSsoEnabledProfilesValue, diags := types.SetValueFrom(context.Background(), types.StringType, r.SpBrowserSso.EnabledProfiles)
		respDiags.Append(diags...)
		var spBrowserSsoEncryptionPolicyValue types.Object
		if r.SpBrowserSso.EncryptionPolicy == nil {
			spBrowserSsoEncryptionPolicyValue = types.ObjectNull(spBrowserSsoEncryptionPolicyAttrTypes)
		} else {
			spBrowserSsoEncryptionPolicyEncryptedAttributesValue, diags := types.SetValueFrom(context.Background(), types.StringType, r.SpBrowserSso.EncryptionPolicy.EncryptedAttributes)
			respDiags.Append(diags...)
			spBrowserSsoEncryptionPolicyValue, diags = types.ObjectValue(spBrowserSsoEncryptionPolicyAttrTypes, map[string]attr.Value{
				"encrypt_assertion":             types.BoolPointerValue(r.SpBrowserSso.EncryptionPolicy.EncryptAssertion),
				"encrypt_slo_subject_name_id":   types.BoolPointerValue(r.SpBrowserSso.EncryptionPolicy.EncryptSloSubjectNameId),
				"encrypted_attributes":          spBrowserSsoEncryptionPolicyEncryptedAttributesValue,
				"slo_subject_name_id_encrypted": types.BoolPointerValue(r.SpBrowserSso.EncryptionPolicy.SloSubjectNameIDEncrypted),
			})
			respDiags.Append(diags...)
		}
		spBrowserSsoIncomingBindingsValue, diags := types.SetValueFrom(context.Background(), types.StringType, r.SpBrowserSso.IncomingBindings)
		respDiags.Append(diags...)
		var spBrowserSsoMessageCustomizationsValues []attr.Value
		for _, spBrowserSsoMessageCustomizationsResponseValue := range r.SpBrowserSso.MessageCustomizations {
			spBrowserSsoMessageCustomizationsValue, diags := types.ObjectValue(spBrowserSsoMessageCustomizationsAttrTypes, map[string]attr.Value{
				"context_name":       types.StringPointerValue(spBrowserSsoMessageCustomizationsResponseValue.ContextName),
				"message_expression": types.StringPointerValue(spBrowserSsoMessageCustomizationsResponseValue.MessageExpression),
			})
			respDiags.Append(diags...)
			spBrowserSsoMessageCustomizationsValues = append(spBrowserSsoMessageCustomizationsValues, spBrowserSsoMessageCustomizationsValue)
		}
		spBrowserSsoMessageCustomizationsValue, diags := types.SetValue(spBrowserSsoMessageCustomizationsElementType, spBrowserSsoMessageCustomizationsValues)
		respDiags.Append(diags...)
		var spBrowserSsoSloServiceEndpointsValues []attr.Value
		for _, spBrowserSsoSloServiceEndpointsResponseValue := range r.SpBrowserSso.SloServiceEndpoints {
			spBrowserSsoSloServiceEndpointsValue, diags := types.ObjectValue(spBrowserSsoSloServiceEndpointsAttrTypes, map[string]attr.Value{
				"binding":      types.StringPointerValue(spBrowserSsoSloServiceEndpointsResponseValue.Binding),
				"response_url": types.StringPointerValue(spBrowserSsoSloServiceEndpointsResponseValue.ResponseUrl),
				"url":          types.StringValue(spBrowserSsoSloServiceEndpointsResponseValue.Url),
			})
			respDiags.Append(diags...)
			spBrowserSsoSloServiceEndpointsValues = append(spBrowserSsoSloServiceEndpointsValues, spBrowserSsoSloServiceEndpointsValue)
		}
		spBrowserSsoSloServiceEndpointsValue, diags := types.SetValue(spBrowserSsoSloServiceEndpointsElementType, spBrowserSsoSloServiceEndpointsValues)
		respDiags.Append(diags...)
		var spBrowserSsoSsoServiceEndpointsValues []attr.Value
		for _, spBrowserSsoSsoServiceEndpointsResponseValue := range r.SpBrowserSso.SsoServiceEndpoints {
			// PF will return nil for false for the is_default boolean
			isDefault := types.BoolValue(false)
			if spBrowserSsoSsoServiceEndpointsResponseValue.IsDefault != nil {
				isDefault = types.BoolPointerValue(spBrowserSsoSsoServiceEndpointsResponseValue.IsDefault)
			}
			spBrowserSsoSsoServiceEndpointsValue, diags := types.ObjectValue(spBrowserSsoSsoServiceEndpointsAttrTypes, map[string]attr.Value{
				"binding":    types.StringPointerValue(spBrowserSsoSsoServiceEndpointsResponseValue.Binding),
				"index":      types.Int64PointerValue(spBrowserSsoSsoServiceEndpointsResponseValue.Index),
				"is_default": isDefault,
				"url":        types.StringValue(spBrowserSsoSsoServiceEndpointsResponseValue.Url),
			})
			respDiags.Append(diags...)
			spBrowserSsoSsoServiceEndpointsValues = append(spBrowserSsoSsoServiceEndpointsValues, spBrowserSsoSsoServiceEndpointsValue)
		}
		spBrowserSsoSsoServiceEndpointsValue, diags := types.SetValue(spBrowserSsoSsoServiceEndpointsElementType, spBrowserSsoSsoServiceEndpointsValues)
		respDiags.Append(diags...)
		var spBrowserSsoUrlWhitelistEntriesValue types.Set
		if r.SpBrowserSso.UrlWhitelistEntries == nil {
			spBrowserSsoUrlWhitelistEntriesValue = types.SetNull(spBrowserSsoUrlWhitelistEntriesElementType)
		} else {
			var spBrowserSsoUrlWhitelistEntriesValues []attr.Value
			for _, spBrowserSsoUrlWhitelistEntriesResponseValue := range r.SpBrowserSso.UrlWhitelistEntries {
				spBrowserSsoUrlWhitelistEntriesValue, diags := types.ObjectValue(spBrowserSsoUrlWhitelistEntriesAttrTypes, map[string]attr.Value{
					"allow_query_and_fragment": types.BoolPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.AllowQueryAndFragment),
					"require_https":            types.BoolPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.RequireHttps),
					"valid_domain":             types.StringPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.ValidDomain),
					"valid_path":               types.StringPointerValue(spBrowserSsoUrlWhitelistEntriesResponseValue.ValidPath),
				})
				respDiags.Append(diags...)
				spBrowserSsoUrlWhitelistEntriesValues = append(spBrowserSsoUrlWhitelistEntriesValues, spBrowserSsoUrlWhitelistEntriesValue)
			}
			spBrowserSsoUrlWhitelistEntriesValue, diags = types.SetValue(spBrowserSsoUrlWhitelistEntriesElementType, spBrowserSsoUrlWhitelistEntriesValues)
			respDiags.Append(diags...)
		}

		// always_sign_artifact_response, sign_assertions, require_signed_authn_requests can be returned as nil when set to false
		var alwaysSignArtifactResponse, signAssertions, requireSignedAuthnRequests bool
		if r.SpBrowserSso.AlwaysSignArtifactResponse != nil {
			alwaysSignArtifactResponse = *r.SpBrowserSso.AlwaysSignArtifactResponse
		}
		if r.SpBrowserSso.SignAssertions != nil {
			signAssertions = *r.SpBrowserSso.SignAssertions
		}
		if r.SpBrowserSso.RequireSignedAuthnRequests != nil {
			requireSignedAuthnRequests = *r.SpBrowserSso.RequireSignedAuthnRequests
		}

		spBrowserSsoValue, diags = types.ObjectValue(spBrowserSsoAttrTypes, map[string]attr.Value{
			"adapter_mappings":              spBrowserSsoAdapterMappingsValue,
			"always_sign_artifact_response": types.BoolValue(alwaysSignArtifactResponse),
			"artifact":                      spBrowserSsoArtifactValue,
			"assertion_lifetime":            spBrowserSsoAssertionLifetimeValue,
			"attribute_contract":            spBrowserSsoAttributeContractValue,
			"authentication_policy_contract_assertion_mappings": spBrowserSsoAuthenticationPolicyContractAssertionMappingsValue,
			"default_target_url":            types.StringPointerValue(r.SpBrowserSso.DefaultTargetUrl),
			"enabled_profiles":              spBrowserSsoEnabledProfilesValue,
			"encryption_policy":             spBrowserSsoEncryptionPolicyValue,
			"incoming_bindings":             spBrowserSsoIncomingBindingsValue,
			"message_customizations":        spBrowserSsoMessageCustomizationsValue,
			"protocol":                      types.StringValue(r.SpBrowserSso.Protocol),
			"require_signed_authn_requests": types.BoolValue(requireSignedAuthnRequests),
			"sign_assertions":               types.BoolValue(signAssertions),
			"sign_response_as_required":     types.BoolPointerValue(r.SpBrowserSso.SignResponseAsRequired),
			"slo_service_endpoints":         spBrowserSsoSloServiceEndpointsValue,
			"sp_saml_identity_mapping":      types.StringPointerValue(r.SpBrowserSso.SpSamlIdentityMapping),
			"sp_ws_fed_identity_mapping":    types.StringPointerValue(r.SpBrowserSso.SpWsFedIdentityMapping),
			"sso_application_endpoint":      types.StringPointerValue(r.SpBrowserSso.SsoApplicationEndpoint),
			"sso_service_endpoints":         spBrowserSsoSsoServiceEndpointsValue,
			"url_whitelist_entries":         spBrowserSsoUrlWhitelistEntriesValue,
			"ws_fed_token_type":             types.StringPointerValue(r.SpBrowserSso.WsFedTokenType),
			"ws_trust_version":              types.StringPointerValue(r.SpBrowserSso.WsTrustVersion),
		})
		respDiags.Append(diags...)
	}

	state.SpBrowserSso = spBrowserSsoValue
	// virtual_entity_ids
	state.VirtualEntityIds, diags = types.SetValueFrom(context.Background(), types.StringType, r.VirtualEntityIds)
	respDiags.Append(diags...)
	// ws_trust
	wsTrustAttributeContractCoreAttributesAttrTypes := map[string]attr.Type{
		"name":      types.StringType,
		"namespace": types.StringType,
	}
	wsTrustAttributeContractCoreAttributesElementType := types.ObjectType{AttrTypes: wsTrustAttributeContractCoreAttributesAttrTypes}
	wsTrustAttributeContractExtendedAttributesAttrTypes := map[string]attr.Type{
		"name":      types.StringType,
		"namespace": types.StringType,
	}
	wsTrustAttributeContractExtendedAttributesElementType := types.ObjectType{AttrTypes: wsTrustAttributeContractExtendedAttributesAttrTypes}
	wsTrustAttributeContractAttrTypes := map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: wsTrustAttributeContractCoreAttributesElementType},
		"extended_attributes": types.SetType{ElemType: wsTrustAttributeContractExtendedAttributesElementType},
	}
	wsTrustMessageCustomizationsAttrTypes := map[string]attr.Type{
		"context_name":       types.StringType,
		"message_expression": types.StringType,
	}
	wsTrustMessageCustomizationsElementType := types.ObjectType{AttrTypes: wsTrustMessageCustomizationsAttrTypes}
	wsTrustRequestContractRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	wsTrustTokenProcessorMappingsAttributeContractFulfillmentAttrTypes := attributecontractfulfillment.AttrTypes()
	wsTrustTokenProcessorMappingsAttributeContractFulfillmentElementType := types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsAttributeContractFulfillmentAttrTypes}
	wsTrustTokenProcessorMappingsAttributeSourcesAttrTypes := attributesources.AttrTypes()
	wsTrustTokenProcessorMappingsAttributeSourcesElementType := types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsAttributeSourcesAttrTypes}
	wsTrustTokenProcessorMappingsIdpTokenProcessorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	wsTrustTokenProcessorMappingsIssuanceCriteriaAttrTypes := issuancecriteria.AttrTypes()
	wsTrustTokenProcessorMappingsAttrTypes := map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{ElemType: wsTrustTokenProcessorMappingsAttributeContractFulfillmentElementType},
		"attribute_sources":              types.SetType{ElemType: wsTrustTokenProcessorMappingsAttributeSourcesElementType},
		"idp_token_processor_ref":        types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsIdpTokenProcessorRefAttrTypes},
		"issuance_criteria":              types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsIssuanceCriteriaAttrTypes},
		"restricted_virtual_entity_ids":  types.SetType{ElemType: types.StringType},
	}
	wsTrustTokenProcessorMappingsElementType := types.ObjectType{AttrTypes: wsTrustTokenProcessorMappingsAttrTypes}
	wsTrustAttrTypes := map[string]attr.Type{
		"abort_if_not_fulfilled_from_request": types.BoolType,
		"attribute_contract":                  types.ObjectType{AttrTypes: wsTrustAttributeContractAttrTypes},
		"default_token_type":                  types.StringType,
		"encrypt_saml2_assertion":             types.BoolType,
		"generate_key":                        types.BoolType,
		"message_customizations":              types.SetType{ElemType: wsTrustMessageCustomizationsElementType},
		"minutes_after":                       types.Int64Type,
		"minutes_before":                      types.Int64Type,
		"oauth_assertion_profiles":            types.BoolType,
		"partner_service_ids":                 types.SetType{ElemType: types.StringType},
		"request_contract_ref":                types.ObjectType{AttrTypes: wsTrustRequestContractRefAttrTypes},
		"token_processor_mappings":            types.SetType{ElemType: wsTrustTokenProcessorMappingsElementType},
	}
	var wsTrustValue types.Object
	if r.WsTrust == nil {
		wsTrustValue = types.ObjectNull(wsTrustAttrTypes)
	} else {
		var wsTrustAttributeContractCoreAttributesValues []attr.Value
		for _, wsTrustAttributeContractCoreAttributesResponseValue := range r.WsTrust.AttributeContract.CoreAttributes {
			wsTrustAttributeContractCoreAttributesValue, diags := types.ObjectValue(wsTrustAttributeContractCoreAttributesAttrTypes, map[string]attr.Value{
				"name":      types.StringValue(wsTrustAttributeContractCoreAttributesResponseValue.Name),
				"namespace": types.StringValue(wsTrustAttributeContractCoreAttributesResponseValue.Namespace),
			})
			respDiags.Append(diags...)
			wsTrustAttributeContractCoreAttributesValues = append(wsTrustAttributeContractCoreAttributesValues, wsTrustAttributeContractCoreAttributesValue)
		}
		wsTrustAttributeContractCoreAttributesValue, diags := types.SetValue(wsTrustAttributeContractCoreAttributesElementType, wsTrustAttributeContractCoreAttributesValues)
		respDiags.Append(diags...)
		var wsTrustAttributeContractExtendedAttributesValues []attr.Value
		for _, wsTrustAttributeContractExtendedAttributesResponseValue := range r.WsTrust.AttributeContract.ExtendedAttributes {
			wsTrustAttributeContractExtendedAttributesValue, diags := types.ObjectValue(wsTrustAttributeContractExtendedAttributesAttrTypes, map[string]attr.Value{
				"name":      types.StringValue(wsTrustAttributeContractExtendedAttributesResponseValue.Name),
				"namespace": types.StringValue(wsTrustAttributeContractExtendedAttributesResponseValue.Namespace),
			})
			respDiags.Append(diags...)
			wsTrustAttributeContractExtendedAttributesValues = append(wsTrustAttributeContractExtendedAttributesValues, wsTrustAttributeContractExtendedAttributesValue)
		}
		wsTrustAttributeContractExtendedAttributesValue, diags := types.SetValue(wsTrustAttributeContractExtendedAttributesElementType, wsTrustAttributeContractExtendedAttributesValues)
		respDiags.Append(diags...)
		wsTrustAttributeContractValue, diags := types.ObjectValue(wsTrustAttributeContractAttrTypes, map[string]attr.Value{
			"core_attributes":     wsTrustAttributeContractCoreAttributesValue,
			"extended_attributes": wsTrustAttributeContractExtendedAttributesValue,
		})
		respDiags.Append(diags...)
		var wsTrustMessageCustomizationsValues []attr.Value
		for _, wsTrustMessageCustomizationsResponseValue := range r.WsTrust.MessageCustomizations {
			wsTrustMessageCustomizationsValue, diags := types.ObjectValue(wsTrustMessageCustomizationsAttrTypes, map[string]attr.Value{
				"context_name":       types.StringPointerValue(wsTrustMessageCustomizationsResponseValue.ContextName),
				"message_expression": types.StringPointerValue(wsTrustMessageCustomizationsResponseValue.MessageExpression),
			})
			respDiags.Append(diags...)
			wsTrustMessageCustomizationsValues = append(wsTrustMessageCustomizationsValues, wsTrustMessageCustomizationsValue)
		}
		wsTrustMessageCustomizationsValue, diags := types.SetValue(wsTrustMessageCustomizationsElementType, wsTrustMessageCustomizationsValues)
		respDiags.Append(diags...)
		wsTrustPartnerServiceIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, r.WsTrust.PartnerServiceIds)
		respDiags.Append(diags...)
		var wsTrustRequestContractRefValue types.Object
		if r.WsTrust.RequestContractRef == nil {
			wsTrustRequestContractRefValue = types.ObjectNull(wsTrustRequestContractRefAttrTypes)
		} else {
			wsTrustRequestContractRefValue, diags = types.ObjectValue(wsTrustRequestContractRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(r.WsTrust.RequestContractRef.Id),
			})
			respDiags.Append(diags...)
		}
		var wsTrustTokenProcessorMappingsValues []attr.Value
		for _, wsTrustTokenProcessorMappingsResponseValue := range r.WsTrust.TokenProcessorMappings {
			contractFulfillment := wsTrustTokenProcessorMappingsResponseValue.AttributeContractFulfillment
			wsTrustTokenProcessorMappingsAttributeContractFulfillmentValue, diags := attributecontractfulfillment.ToState(context.Background(), &contractFulfillment)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsAttributeSourcesValue, diags := attributesources.ToState(context.Background(), wsTrustTokenProcessorMappingsResponseValue.AttributeSources)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsIdpTokenProcessorRefValue, diags := types.ObjectValue(wsTrustTokenProcessorMappingsIdpTokenProcessorRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(wsTrustTokenProcessorMappingsResponseValue.IdpTokenProcessorRef.Id),
			})
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsIssuanceCriteriaValue, diags := issuancecriteria.ToState(context.Background(), wsTrustTokenProcessorMappingsResponseValue.IssuanceCriteria)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsRestrictedVirtualEntityIdsValue, diags := types.SetValueFrom(context.Background(), types.StringType, wsTrustTokenProcessorMappingsResponseValue.RestrictedVirtualEntityIds)
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsValue, diags := types.ObjectValue(wsTrustTokenProcessorMappingsAttrTypes, map[string]attr.Value{
				"attribute_contract_fulfillment": wsTrustTokenProcessorMappingsAttributeContractFulfillmentValue,
				"attribute_sources":              wsTrustTokenProcessorMappingsAttributeSourcesValue,
				"idp_token_processor_ref":        wsTrustTokenProcessorMappingsIdpTokenProcessorRefValue,
				"issuance_criteria":              wsTrustTokenProcessorMappingsIssuanceCriteriaValue,
				"restricted_virtual_entity_ids":  wsTrustTokenProcessorMappingsRestrictedVirtualEntityIdsValue,
			})
			respDiags.Append(diags...)
			wsTrustTokenProcessorMappingsValues = append(wsTrustTokenProcessorMappingsValues, wsTrustTokenProcessorMappingsValue)
		}
		wsTrustTokenProcessorMappingsValue, diags := types.SetValue(wsTrustTokenProcessorMappingsElementType, wsTrustTokenProcessorMappingsValues)
		respDiags.Append(diags...)
		// Ensure that nil values are handled as false for encrypt_saml2_assertion and generate_key
		var encryptSaml2Assertion, generateKey bool
		if r.WsTrust.EncryptSaml2Assertion != nil {
			encryptSaml2Assertion = *r.WsTrust.EncryptSaml2Assertion
		}
		if r.WsTrust.GenerateKey != nil {
			generateKey = *r.WsTrust.GenerateKey
		}
		wsTrustValue, diags = types.ObjectValue(wsTrustAttrTypes, map[string]attr.Value{
			"abort_if_not_fulfilled_from_request": types.BoolPointerValue(r.WsTrust.AbortIfNotFulfilledFromRequest),
			"attribute_contract":                  wsTrustAttributeContractValue,
			"default_token_type":                  types.StringPointerValue(r.WsTrust.DefaultTokenType),
			"encrypt_saml2_assertion":             types.BoolValue(encryptSaml2Assertion),
			"generate_key":                        types.BoolValue(generateKey),
			"message_customizations":              wsTrustMessageCustomizationsValue,
			"minutes_after":                       types.Int64PointerValue(r.WsTrust.MinutesAfter),
			"minutes_before":                      types.Int64PointerValue(r.WsTrust.MinutesBefore),
			"oauth_assertion_profiles":            types.BoolPointerValue(r.WsTrust.OAuthAssertionProfiles),
			"partner_service_ids":                 wsTrustPartnerServiceIdsValue,
			"request_contract_ref":                wsTrustRequestContractRefValue,
			"token_processor_mappings":            wsTrustTokenProcessorMappingsValue,
		})
		respDiags.Append(diags...)
	}

	state.WsTrust = wsTrustValue
	return respDiags
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
