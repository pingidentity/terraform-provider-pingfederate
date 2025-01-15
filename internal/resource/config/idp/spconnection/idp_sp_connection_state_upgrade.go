package idpspconnection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/connectioncert"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func (r *idpSpConnectionResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	adapterOverrideSettingsAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"masked": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"pseudonym": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}

	spBrowserSSOAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"name_format": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	wsTrustAttribute := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"namespace": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}

	messageCustomizationsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"context_name": schema.StringAttribute{
				Optional: true,
			},
			"message_expression": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	channelsAttributeMappingNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"field_name": schema.StringAttribute{
				Required: true,
			},
			"saas_field_info": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"attribute_names": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
					},
					"character_case": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"create_only": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"default_value": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"expression": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"masked": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"parser": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"trim": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
				},
				Required: true,
			},
		},
	}

	outboundProvisionTargetSettingsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"value": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}

	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"connection_id": schema.StringAttribute{
						Required: true,
					},
					"active": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"additional_allowed_entities_configuration": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"additional_allowed_entities": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"entity_description": schema.StringAttribute{
											Optional: true,
										},
										"entity_id": schema.StringAttribute{
											Optional: true,
										},
									},
								},
								Optional: true,
								Computed: true,
							},
							"allow_additional_entities": schema.BoolAttribute{
								Optional: true,
							},
							"allow_all_entities": schema.BoolAttribute{
								Optional: true,
							},
						},
						Optional: true,
					},
					"application_icon_url": schema.StringAttribute{
						Optional: true,
					},
					"application_name": schema.StringAttribute{
						Optional: true,
					},
					"attribute_query": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
							"attribute_sources":              attributesources.ToSchema(1, false),
							"attributes": schema.SetAttribute{
								ElementType: types.StringType,
								Required:    true,
							},
							"issuance_criteria": issuancecriteria.ToSchema(),
							"policy": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"encrypt_assertion": schema.BoolAttribute{
										Optional: true,
									},
									"require_encrypted_name_id": schema.BoolAttribute{
										Optional: true,
									},
									"require_signed_attribute_query": schema.BoolAttribute{
										Optional: true,
									},
									"sign_assertion": schema.BoolAttribute{
										Optional: true,
									},
									"sign_response": schema.BoolAttribute{
										Optional: true,
									},
								},
								Optional: true,
							},
						},
						Optional: true,
					},
					"base_url": schema.StringAttribute{
						Optional: true,
					},
					"connection_target_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"contact_info": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"company": schema.StringAttribute{
								Optional: true,
							},
							"email": schema.StringAttribute{
								Optional: true,
							},
							"first_name": schema.StringAttribute{
								Optional: true,
							},
							"last_name": schema.StringAttribute{
								Optional: true,
							},
							"phone": schema.StringAttribute{
								Optional: true,
							},
						},
						Optional: true,
					},
					"creation_date": schema.StringAttribute{
						Optional: false,
						Computed: true,
					},
					"credentials": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"verification_issuer_dn": schema.StringAttribute{
								Optional: true,
							},
							"verification_subject_dn": schema.StringAttribute{
								Optional: true,
							},
							"certs": schema.ListNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"cert_view": schema.SingleNestedAttribute{
											Optional: false,
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"serial_number": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"subject_dn": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"subject_alternative_names": schema.SetAttribute{
													ElementType: types.StringType,
													Optional:    false,
													Computed:    true,
												},
												"issuer_dn": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"valid_from": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"expires": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"key_algorithm": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"key_size": schema.Int64Attribute{
													Optional: false,
													Computed: true,
												},
												"signature_algorithm": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"version": schema.Int64Attribute{
													Optional: false,
													Computed: true,
												},
												"sha1fingerprint": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"sha256fingerprint": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"status": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
												"crypto_provider": schema.StringAttribute{
													Optional: false,
													Computed: true,
												},
											},
										},
										"x509file": schema.SingleNestedAttribute{
											Required: true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Optional: true,
													Computed: true,
												},
												"file_data": schema.StringAttribute{
													Required: true,
												},
												"formatted_file_data": schema.StringAttribute{
													Computed: true,
												},
												"crypto_provider": schema.StringAttribute{
													Optional: true,
												},
											},
										},
										"active_verification_cert": schema.BoolAttribute{
											Optional: true,
											Computed: true,
										},
										"primary_verification_cert": schema.BoolAttribute{
											Optional: true,
											Computed: true,
										},
										"secondary_verification_cert": schema.BoolAttribute{
											Optional: true,
											Computed: true,
										},
										"encryption_cert": schema.BoolAttribute{
											Optional: true,
											Computed: true,
										},
									},
								},
							},
							"block_encryption_algorithm": schema.StringAttribute{
								Optional: true,
							},
							"key_transport_algorithm": schema.StringAttribute{
								Optional: true,
							},
							"signing_settings": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"signing_key_pair_ref": schema.SingleNestedAttribute{
										Attributes: resourcelink.ToSchema(),
										Optional:   true,
									},
									"alternative_signing_key_pair_refs": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: resourcelink.ToSchema(),
										},
										Optional: true,
									},
									"algorithm": schema.StringAttribute{
										Optional: true,
									},
									"include_cert_in_signature": schema.BoolAttribute{
										Optional: true,
										Computed: true,
									},
									"include_raw_key_in_signature": schema.BoolAttribute{
										Optional: true,
									},
								},
								Optional: true,
							},
							"decryption_key_pair_ref": schema.SingleNestedAttribute{
								Attributes: resourcelink.ToSchema(),
								Optional:   true,
							},
							"secondary_decryption_key_pair_ref": schema.SingleNestedAttribute{
								Attributes: resourcelink.ToSchema(),
								Optional:   true,
							},
							"outbound_back_channel_auth": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"http_basic_credentials": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"username": schema.StringAttribute{
												Optional: true,
											},
											"password": schema.StringAttribute{
												Optional:  true,
												Sensitive: true,
											},
										},
										Optional: true,
									},
									"digital_signature": schema.BoolAttribute{
										Optional: true,
									},
									"ssl_auth_key_pair_ref": schema.SingleNestedAttribute{
										Attributes: resourcelink.ToSchema(),
										Optional:   true,
									},
									"validate_partner_cert": schema.BoolAttribute{
										Optional: true,
										Computed: true,
									},
								},
								Optional: true,
							},
							"inbound_back_channel_auth": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"http_basic_credentials": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"username": schema.StringAttribute{
												Optional: true,
											},
											"password": schema.StringAttribute{
												Optional:  true,
												Sensitive: true,
											},
										},
										Optional: true,
									},
									"digital_signature": schema.BoolAttribute{
										Optional: true,
									},
									"verification_subject_dn": schema.StringAttribute{
										Optional: true,
									},
									"verification_issuer_dn": schema.StringAttribute{
										Optional: true,
									},
									"certs": schema.ListNestedAttribute{
										Optional: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"cert_view": schema.SingleNestedAttribute{
													Optional: false,
													Computed: true,
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"serial_number": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"subject_dn": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"subject_alternative_names": schema.SetAttribute{
															ElementType: types.StringType,
															Optional:    false,
															Computed:    true,
														},
														"issuer_dn": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"valid_from": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"expires": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"key_algorithm": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"key_size": schema.Int64Attribute{
															Optional: false,
															Computed: true,
														},
														"signature_algorithm": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"version": schema.Int64Attribute{
															Optional: false,
															Computed: true,
														},
														"sha1fingerprint": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"sha256fingerprint": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"status": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
														"crypto_provider": schema.StringAttribute{
															Optional: false,
															Computed: true,
														},
													},
												},
												"x509file": schema.SingleNestedAttribute{
													Required: true,
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional: true,
															Computed: true,
														},
														"file_data": schema.StringAttribute{
															Required: true,
														},
														"formatted_file_data": schema.StringAttribute{
															Computed: true,
														},
														"crypto_provider": schema.StringAttribute{
															Optional: true,
														},
													},
												},
												"active_verification_cert": schema.BoolAttribute{
													Optional: true,
													Computed: true,
												},
												"primary_verification_cert": schema.BoolAttribute{
													Optional: true,
													Computed: true,
												},
												"secondary_verification_cert": schema.BoolAttribute{
													Optional: true,
													Computed: true,
												},
												"encryption_cert": schema.BoolAttribute{
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"require_ssl": schema.BoolAttribute{
										Optional: true,
									},
								},
								Optional: true,
							},
						},
						Optional: true,
					},
					"default_virtual_entity_id": schema.StringAttribute{
						Optional: true,
					},
					"entity_id": schema.StringAttribute{
						Required: true,
					},
					"extended_properties": schema.MapNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"values": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
								},
							},
						},
						Optional: true,
					},
					"id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"license_connection_group": schema.StringAttribute{
						Optional: true,
					},
					"logging_mode": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"metadata_reload_settings": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"enable_auto_metadata_update": schema.BoolAttribute{
								Optional: true,
								Computed: true,
							},
							"metadata_url_ref": resourcelink.SingleNestedAttribute(),
						},
						Optional: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"outbound_provision": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"channels": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"active": schema.BoolAttribute{
											Required: true,
										},
										"attribute_mapping_all": schema.SetNestedAttribute{
											NestedObject: channelsAttributeMappingNestedObject,
											Optional:     false,
											Computed:     true,
										},
										"attribute_mapping": schema.SetNestedAttribute{
											NestedObject: channelsAttributeMappingNestedObject,
											Required:     true,
										},
										"channel_source": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"account_management_settings": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"account_status_algorithm": schema.StringAttribute{
															Required: true,
														},
														"account_status_attribute_name": schema.StringAttribute{
															Required: true,
														},
														"default_status": schema.BoolAttribute{
															Optional: true,
														},
														"flag_comparison_status": schema.BoolAttribute{
															Optional: true,
														},
														"flag_comparison_value": schema.StringAttribute{
															Optional: true,
														},
													},
													Required: true,
												},
												"base_dn": schema.StringAttribute{
													Required: true,
												},
												"change_detection_settings": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"changed_users_algorithm": schema.StringAttribute{
															Required: true,
														},
														"group_object_class": schema.StringAttribute{
															Required: true,
														},
														"time_stamp_attribute_name": schema.StringAttribute{
															Required: true,
														},
														"user_object_class": schema.StringAttribute{
															Required: true,
														},
														"usn_attribute_name": schema.StringAttribute{
															Optional: true,
														},
													},
													Required: true,
												},
												"data_source": resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to an LDAP datastore."),
												"group_membership_detection": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"group_member_attribute_name": schema.StringAttribute{
															Optional: true,
														},
														"member_of_group_attribute_name": schema.StringAttribute{
															Optional: true,
														},
													},
													Required: true,
												},
												"group_source_location": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"filter": schema.StringAttribute{
															Optional: true,
														},
														"group_dn": schema.StringAttribute{
															Optional: true,
														},
														"nested_search": schema.BoolAttribute{
															Optional: true,
															Computed: true,
														},
													},
													Optional: true,
													Computed: true,
												},
												"guid_attribute_name": schema.StringAttribute{
													Required: true,
												},
												"guid_binary": schema.BoolAttribute{
													Required: true,
												},
												"user_source_location": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"filter": schema.StringAttribute{
															Optional: true,
														},
														"group_dn": schema.StringAttribute{
															Optional: true,
														},
														"nested_search": schema.BoolAttribute{
															Optional: true,
															Computed: true,
														},
													},
													Required: true,
												},
											},
											Required: true,
										},
										"max_threads": schema.Int64Attribute{
											Optional: true,
											Computed: true,
										},
										"name": schema.StringAttribute{
											Required: true,
										},
										"timeout": schema.Int64Attribute{
											Optional: true,
											Computed: true,
										},
									},
								},
								Required: true,
							},
							"custom_schema": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attributes": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"multi_valued": schema.BoolAttribute{
													Optional: true,
												},
												"name": schema.StringAttribute{
													Optional: true,
												},
												"sub_attributes": schema.SetAttribute{
													ElementType: types.StringType,
													Optional:    true,
													Computed:    true,
												},
												"types": schema.SetAttribute{
													ElementType: types.StringType,
													Optional:    true,
													Computed:    true,
												},
											},
										},
										Optional: true,
										Computed: true,
									},
									"namespace": schema.StringAttribute{
										Optional: true,
									},
								},
								Optional: true,
							},
							"target_settings_all": schema.SetNestedAttribute{
								NestedObject: outboundProvisionTargetSettingsNestedObject,
								Optional:     false,
								Computed:     true,
							},
							"target_settings": schema.SetNestedAttribute{
								NestedObject: outboundProvisionTargetSettingsNestedObject,
								Required:     true,
							},
							"type": schema.StringAttribute{
								Required: true,
							},
						},
						Optional: true,
					},
					"sp_browser_sso": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"adapter_mappings": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
											Optional: true,
											Computed: true,
										},
										"adapter_override_settings": schema.SingleNestedAttribute{
											Attributes: map[string]schema.Attribute{
												"attribute_contract": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"core_attributes": schema.SetNestedAttribute{
															NestedObject: adapterOverrideSettingsAttribute,
															Required:     true,
														},
														"extended_attributes": schema.SetNestedAttribute{
															NestedObject: adapterOverrideSettingsAttribute,
															Optional:     true,
															Computed:     true,
														},
														"mask_ognl_values": schema.BoolAttribute{
															Optional: true,
															Computed: true,
														},
														"unique_user_key_attribute": schema.StringAttribute{
															Optional: true,
														},
													},
													Optional: true,
												},
												"attribute_mapping": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
														"attribute_sources":              attributesources.ToSchema(0, false),
														"issuance_criteria":              issuancecriteria.ToSchema(),
													},
													Optional: true,
												},
												"authn_ctx_class_ref": schema.StringAttribute{
													Optional: true,
												},
												"configuration": pluginconfiguration.ToSchema(),
												"id": schema.StringAttribute{
													Required: true,
												},
												"name": schema.StringAttribute{
													Required: true,
												},
												"parent_ref":            resourcelink.CompleteSingleNestedAttribute(true, false, false, ""),
												"plugin_descriptor_ref": resourcelink.CompleteSingleNestedAttribute(false, false, true, ""),
											},
											Optional: true,
										},
										"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
										"attribute_sources":              attributesources.ToSchema(0, false),
										"idp_adapter_ref":                resourcelink.CompleteSingleNestedAttribute(true, false, false, ""),
										"issuance_criteria":              issuancecriteria.ToSchema(),
										"restrict_virtual_entity_ids": schema.BoolAttribute{
											Optional: true,
										},
										"restricted_virtual_entity_ids": schema.SetAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Computed:    true,
										},
									},
								},
								Required: true,
							},
							"always_sign_artifact_response": schema.BoolAttribute{
								Optional: true,
							},
							"artifact": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"lifetime": schema.Int64Attribute{
										Required: true,
									},
									"resolver_locations": schema.SetNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"index": schema.Int64Attribute{
													Required: true,
												},
												"url": schema.StringAttribute{
													Required: true,
												},
											},
										},
										Required: true,
									},
									"source_id": schema.StringAttribute{
										Optional: true,
									},
								},
								Optional: true,
							},
							"assertion_lifetime": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"minutes_after": schema.Int64Attribute{
										Required: true,
									},
									"minutes_before": schema.Int64Attribute{
										Required: true,
									},
								},
								Required: true,
							},
							"attribute_contract": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"core_attributes": schema.SetNestedAttribute{
										NestedObject: spBrowserSSOAttribute,
										Optional:     true,
									},
									"extended_attributes": schema.SetNestedAttribute{
										NestedObject: spBrowserSSOAttribute,
										Optional:     true,
										Computed:     true,
									},
								},
								Required: true,
							},
							"authentication_policy_contract_assertion_mappings": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"abort_sso_transaction_as_fail_safe": schema.BoolAttribute{
											Optional: true,
											Computed: true,
										},
										"attribute_contract_fulfillment":     attributecontractfulfillment.ToSchema(true, false, false),
										"attribute_sources":                  attributesources.ToSchema(0, false),
										"authentication_policy_contract_ref": resourcelink.CompleteSingleNestedAttribute(false, false, true, ""),
										"issuance_criteria":                  issuancecriteria.ToSchema(),
										"restrict_virtual_entity_ids": schema.BoolAttribute{
											Optional: true,
										},
										"restricted_virtual_entity_ids": schema.SetAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Computed:    true,
										},
									},
								},
								Optional: true,
								Computed: true,
							},
							"default_target_url": schema.StringAttribute{
								Optional: true,
							},
							"enabled_profiles": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},
							"encryption_policy": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"encrypt_assertion": schema.BoolAttribute{
										Optional: true,
									},
									"encrypt_slo_subject_name_id": schema.BoolAttribute{
										Optional: true,
									},
									"encrypted_attributes": schema.SetAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Computed:    true,
									},
									"slo_subject_name_id_encrypted": schema.BoolAttribute{
										Optional: true,
									},
								},
								Optional: true,
							},
							"incoming_bindings": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},
							"message_customizations": schema.SetNestedAttribute{
								NestedObject: messageCustomizationsNestedObject,
								Optional:     true,
								Computed:     true,
							},
							"protocol": schema.StringAttribute{
								Required: true,
							},
							"require_signed_authn_requests": schema.BoolAttribute{
								Optional: true,
							},
							"sign_assertions": schema.BoolAttribute{
								Optional: true,
							},
							"sign_response_as_required": schema.BoolAttribute{
								Optional: true,
								Computed: true,
							},
							"slo_service_endpoints": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"binding": schema.StringAttribute{
											Optional: true,
										},
										"response_url": schema.StringAttribute{
											Optional: true,
										},
										"url": schema.StringAttribute{
											Required: true,
										},
									},
								},
								Optional: true,
								Computed: true,
							},
							"sp_saml_identity_mapping": schema.StringAttribute{
								Optional: true,
							},
							"sp_ws_fed_identity_mapping": schema.StringAttribute{
								Optional: true,
							},
							"sso_service_endpoints": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"binding": schema.StringAttribute{
											Optional: true,
										},
										"index": schema.Int64Attribute{
											Optional: true,
										},
										"is_default": schema.BoolAttribute{
											Optional: true,
											Computed: true,
										},
										"url": schema.StringAttribute{
											Required: true,
										},
									},
								},
								Required: true,
							},
							"url_whitelist_entries": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"allow_query_and_fragment": schema.BoolAttribute{
											Optional: true,
										},
										"require_https": schema.BoolAttribute{
											Optional: true,
										},
										"valid_domain": schema.StringAttribute{
											Optional: true,
										},
										"valid_path": schema.StringAttribute{
											Optional: true,
										},
									},
								},
								Optional: true,
							},
							"ws_fed_token_type": schema.StringAttribute{
								Optional: true,
							},
							"ws_trust_version": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
							"sso_application_endpoint": schema.StringAttribute{
								Optional: false,
								Computed: true,
							},
						},
						Optional: true,
					},
					"type": schema.StringAttribute{
						Optional: false,
						Computed: true,
					},
					"virtual_entity_ids": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
					},
					"ws_trust": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"abort_if_not_fulfilled_from_request": schema.BoolAttribute{
								Optional: true,
							},
							"attribute_contract": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"core_attributes": schema.SetNestedAttribute{
										NestedObject: wsTrustAttribute,
										Optional:     true,
									},
									"extended_attributes": schema.SetNestedAttribute{
										NestedObject: wsTrustAttribute,
										Optional:     true,
										Computed:     true,
									},
								},
								Required: true,
							},
							"default_token_type": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
							"encrypt_saml2_assertion": schema.BoolAttribute{
								Optional: true,
							},
							"generate_key": schema.BoolAttribute{
								Optional: true,
							},
							"message_customizations": schema.SetNestedAttribute{
								NestedObject: messageCustomizationsNestedObject,
								Optional:     true,
								Computed:     true,
							},
							"minutes_after": schema.Int64Attribute{
								Optional: true,
								Computed: true,
							},
							"minutes_before": schema.Int64Attribute{
								Optional: true,
								Computed: true,
							},
							"oauth_assertion_profiles": schema.BoolAttribute{
								Optional: true,
							},
							"partner_service_ids": schema.SetAttribute{
								ElementType: types.StringType,
								Required:    true,
							},
							"request_contract_ref": resourcelink.CompleteSingleNestedAttribute(true, false, false, ""),
							"token_processor_mappings": schema.SetNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"attribute_contract_fulfillment": attributecontractfulfillment.ToSchema(true, false, false),
										"attribute_sources":              attributesources.ToSchema(0, false),
										"idp_token_processor_ref":        resourcelink.CompleteSingleNestedAttribute(false, false, true, "Reference to the associated token processor."),
										"issuance_criteria":              issuancecriteria.ToSchema(),
										"restricted_virtual_entity_ids": schema.SetAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Computed:    true,
										},
									},
								},
								Required: true,
							},
						},
						Optional: true,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var diags diag.Diagnostics
				var priorStateData idpSpConnectionModel

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				priorStateData.Credentials, diags = priorStateData.schemaUpgradeCredentialsV0toV1(ctx)
				resp.Diagnostics.Append(diags...)

				resp.Diagnostics.Append(resp.State.Set(ctx, priorStateData)...)
			},
		},
	}
}

func (p *idpSpConnectionModel) schemaUpgradeCredentialsV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags, respDiags diag.Diagnostics
	planAttribute := p.Credentials

	if planAttribute.IsNull() {
		return types.ObjectNull(credentialsAttrTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(credentialsAttrTypes), diags
	} else {
		// Update the individual certs to use the new attribute names
		credentialsAttrs := planAttribute.Attributes()
		if credentialsAttrs["certs"].IsNull() {
			credentialsAttrs["certs"] = types.ListNull(connectioncert.ObjType())
		} else if credentialsAttrs["certs"].IsUnknown() {
			credentialsAttrs["certs"] = types.ListUnknown(connectioncert.ObjType())
		} else {
			finalCertValues := []attr.Value{}
			for _, cert := range credentialsAttrs["certs"].(types.List).Elements() {
				updatedCert, respDiags := p.schemaUpgradeCertV0toV1(ctx, cert.(types.Object))
				diags.Append(respDiags...)
				finalCertValues = append(finalCertValues, updatedCert)
			}
			credentialsAttrs["certs"], respDiags = types.ListValue(connectioncert.ObjType(), finalCertValues)
			diags.Append(respDiags...)
		}

		// Update the inbound_back_channel_auth certs as well
		if credentialsAttrs["inbound_back_channel_auth"].IsUnknown() {
			credentialsAttrs["inbound_back_channel_auth"] = types.ObjectUnknown(credentialsInboundBackChannelAuthAttrTypes)
		} else if credentialsAttrs["inbound_back_channel_auth"].IsNull() {
			credentialsAttrs["inbound_back_channel_auth"] = types.ObjectNull(credentialsInboundBackChannelAuthAttrTypes)
		} else {
			inboundBackChannelAuthAttrs := credentialsAttrs["inbound_back_channel_auth"].(types.Object).Attributes()
			if inboundBackChannelAuthAttrs["certs"].IsNull() {
				inboundBackChannelAuthAttrs["certs"] = types.ListNull(connectioncert.ObjType())
			} else if inboundBackChannelAuthAttrs["certs"].IsUnknown() {
				inboundBackChannelAuthAttrs["certs"] = types.ListUnknown(connectioncert.ObjType())
			} else {
				finalCertValues := []attr.Value{}
				for _, cert := range inboundBackChannelAuthAttrs["certs"].(types.List).Elements() {
					updatedCert, respDiags := p.schemaUpgradeCertV0toV1(ctx, cert.(types.Object))
					diags.Append(respDiags...)
					finalCertValues = append(finalCertValues, updatedCert)
				}
				inboundBackChannelAuthAttrs["certs"], respDiags = types.ListValue(connectioncert.ObjType(), finalCertValues)
				diags.Append(respDiags...)
			}
			credentialsAttrs["inbound_back_channel_auth"], respDiags = types.ObjectValue(credentialsInboundBackChannelAuthAttrTypes, inboundBackChannelAuthAttrs)
			diags.Append(respDiags...)
		}

		result, respDiags := types.ObjectValue(credentialsAttrTypes, credentialsAttrs)
		diags.Append(respDiags...)

		return result, diags
	}
}

func (p *idpSpConnectionModel) schemaUpgradeCertV0toV1(_ context.Context, certv1 types.Object) (types.Object, diag.Diagnostics) {
	var diags, respDiags diag.Diagnostics
	finalAttrs := map[string]attr.Value{}
	certAttrs := certv1.Attributes()

	finalAttrs["x509_file"] = certAttrs["x509file"]
	finalAttrs["active_verification_cert"] = certAttrs["active_verification_cert"]
	finalAttrs["primary_verification_cert"] = certAttrs["primary_verification_cert"]
	finalAttrs["secondary_verification_cert"] = certAttrs["secondary_verification_cert"]
	finalAttrs["encryption_cert"] = certAttrs["encryption_cert"]

	if certAttrs["cert_view"].IsNull() {
		finalAttrs["cert_view"] = types.ObjectNull(connectioncert.CertViewAttrType())
	} else if certAttrs["cert_view"].IsUnknown() {
		finalAttrs["cert_view"] = types.ObjectUnknown(connectioncert.CertViewAttrType())
	} else {
		finalCertViewAttrs := map[string]attr.Value{}
		certViewAttrs := certAttrs["cert_view"].(types.Object).Attributes()
		finalCertViewAttrs["crypto_provider"] = certViewAttrs["crypto_provider"]
		finalCertViewAttrs["expires"] = certViewAttrs["expires"]
		finalCertViewAttrs["id"] = certViewAttrs["id"]
		finalCertViewAttrs["issuer_dn"] = certViewAttrs["issuer_dn"]
		finalCertViewAttrs["key_algorithm"] = certViewAttrs["key_algorithm"]
		finalCertViewAttrs["key_size"] = certViewAttrs["key_size"]
		finalCertViewAttrs["serial_number"] = certViewAttrs["serial_number"]
		// Update the sha attribute names
		finalCertViewAttrs["sha1_fingerprint"] = certViewAttrs["sha1fingerprint"]
		finalCertViewAttrs["sha256_fingerprint"] = certViewAttrs["sha256fingerprint"]
		finalCertViewAttrs["signature_algorithm"] = certViewAttrs["signature_algorithm"]
		finalCertViewAttrs["status"] = certViewAttrs["status"]
		finalCertViewAttrs["subject_alternative_names"] = certViewAttrs["subject_alternative_names"]
		finalCertViewAttrs["subject_dn"] = certViewAttrs["subject_dn"]
		finalCertViewAttrs["valid_from"] = certViewAttrs["valid_from"]
		finalCertViewAttrs["version"] = certViewAttrs["version"]

		finalAttrs["cert_view"], respDiags = types.ObjectValue(connectioncert.CertViewAttrType(), finalCertViewAttrs)
		diags.Append(respDiags...)
	}

	result, respDiags := types.ObjectValue(connectioncert.AttrTypes(), finalAttrs)
	diags.Append(respDiags...)
	return result, diags
}
