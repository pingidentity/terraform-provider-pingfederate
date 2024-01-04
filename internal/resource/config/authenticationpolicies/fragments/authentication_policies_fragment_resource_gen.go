package authenticationpoliciesfragments

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type AuthenticationPoliciesFragmentModel struct {
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Inputs      types.Object `tfsdk:"inputs"`
	Name        types.String `tfsdk:"name"`
	Outputs     types.Object `tfsdk:"outputs"`
	RootNode    types.Object `tfsdk:"root_node"`
}

func AuthenticationPoliciesFragmentResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the authentication policy fragment.",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "The authentication policy fragment ID. ID is unique.",
			},
			"inputs": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
					"location": schema.StringAttribute{
						Optional:    true,
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
					},
				},
				Optional:    true,
				Description: "A reference to a resource.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The authentication policy fragment name. Name is unique.",
			},
			"outputs": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
					"location": schema.StringAttribute{
						Optional:    true,
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
					},
				},
				Optional:    true,
				Description: "A reference to a resource.",
			},
			"root_node": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"apc_mapping_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_mapping": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract_fulfillment": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:    true,
														Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:    true,
														Description: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:    true,
												Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "The value for this attribute.",
											},
										},
										Required:    true,
										Description: "Defines how an attribute in an attribute contract should be populated.",
									},
									"attribute_sources": schema.StringAttribute{
										Optional:    true,
										Description: "A list of configured data stores to look up attributes from.",
									},
									"issuance_criteria": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"conditional_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"attribute_name": schema.StringAttribute{
															Required:    true,
															Description: "The name of the attribute to use in this issuance criterion.",
														},
														"condition": schema.StringAttribute{
															Required:    true,
															Description: "The condition that will be applied to the source attribute's value and the expected value.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"EQUALS",
																	"EQUALS_CASE_INSENSITIVE",
																	"EQUALS_DN",
																	"NOT_EQUAL",
																	"NOT_EQUAL_CASE_INSENSITIVE",
																	"NOT_EQUAL_DN",
																	"MULTIVALUE_CONTAINS",
																	"MULTIVALUE_CONTAINS_CASE_INSENSITIVE",
																	"MULTIVALUE_CONTAINS_DN",
																	"MULTIVALUE_DOES_NOT_CONTAIN",
																	"MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE",
																	"MULTIVALUE_DOES_NOT_CONTAIN_DN",
																),
															},
														},
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"source": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"id": schema.StringAttribute{
																	Optional:    true,
																	Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																},
																"type": schema.StringAttribute{
																	Required:    true,
																	Description: "The source type of this key.",
																	Validators: []validator.String{
																		stringvalidator.OneOf(
																			"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																			"ACCOUNT_LINK",
																			"ADAPTER",
																			"ASSERTION",
																			"CONTEXT",
																			"CUSTOM_DATA_STORE",
																			"EXPRESSION",
																			"JDBC_DATA_STORE",
																			"LDAP_DATA_STORE",
																			"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																			"MAPPED_ATTRIBUTES",
																			"NO_MAPPING",
																			"TEXT",
																			"TOKEN",
																			"REQUEST",
																			"OAUTH_PERSISTENT_GRANT",
																			"SUBJECT_TOKEN",
																			"ACTOR_TOKEN",
																			"PASSWORD_CREDENTIAL_VALIDATOR",
																			"IDP_CONNECTION",
																			"AUTHENTICATION_POLICY_CONTRACT",
																			"CLAIMS",
																			"LOCAL_IDENTITY_PROFILE",
																			"EXTENDED_CLIENT_METADATA",
																			"EXTENDED_PROPERTIES",
																			"TRACKED_HTTP_PARAMS",
																			"FRAGMENT",
																			"INPUTS",
																			"ATTRIBUTE_QUERY",
																			"IDENTITY_STORE_USER",
																			"IDENTITY_STORE_GROUP",
																			"SCIM_USER",
																			"SCIM_GROUP",
																		),
																	},
																},
															},
															Required:    true,
															Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														},
														"value": schema.StringAttribute{
															Required:    true,
															Description: "The expected value of this issuance criterion.",
														},
													},
												},
												Optional:    true,
												Description: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											},
											"expression_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"expression": schema.StringAttribute{
															Required:    true,
															Description: "The OGNL expression to evaluate.",
														},
													},
												},
												Optional:    true,
												Description: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											},
										},
										Optional:    true,
										Description: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									},
								},
								Required:    true,
								Description: "A list of mappings from attribute sources to attribute targets.",
							},
							"authentication_policy_contract_ref": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:    true,
										Description: "The ID of the resource.",
									},
									"location": schema.StringAttribute{
										Optional:    true,
										Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
									},
								},
								Required:    true,
								Description: "A reference to a resource.",
							},
							"context": schema.StringAttribute{
								Optional:    true,
								Description: "The result context.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The authentication selection type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"APC_MAPPING",
										"LOCAL_IDENTITY_MAPPING",
										"AUTHN_SELECTOR",
										"AUTHN_SOURCE",
										"DONE",
										"CONTINUE",
										"RESTART",
										"FRAGMENT",
									),
								},
							},
						},
						Optional:    true,
						Description: "An authentication policy contract selection action.",
					},
					"authn_selector_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"authentication_selector_ref": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:    true,
										Description: "The ID of the resource.",
									},
									"location": schema.StringAttribute{
										Optional:    true,
										Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
									},
								},
								Required:    true,
								Description: "A reference to a resource.",
							},
							"context": schema.StringAttribute{
								Optional:    true,
								Description: "The result context.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The authentication selection type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"APC_MAPPING",
										"LOCAL_IDENTITY_MAPPING",
										"AUTHN_SELECTOR",
										"AUTHN_SOURCE",
										"DONE",
										"CONTINUE",
										"RESTART",
										"FRAGMENT",
									),
								},
							},
						},
						Optional:    true,
						Description: "An authentication selector selection action.",
					},
					"authn_source_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_rules": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"fallback_to_success": schema.BoolAttribute{
										Optional:    true,
										Description: "When all the rules fail, you may choose to default to the general success action or fail. Default to success.",
									},
									"items": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"attribute_name": schema.StringAttribute{
													Optional:    true,
													Description: "The name of the attribute to use in this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
												},
												"attribute_source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:    true,
															Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:    true,
															Description: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Optional:    true,
													Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"condition": schema.StringAttribute{
													Optional:    true,
													Description: "The condition that will be applied to the attribute's expected value. This field is required if the Attribute Source type is not 'EXPRESSION'.",
													Validators: []validator.String{
														stringvalidator.OneOf(
															"EQUALS",
															"EQUALS_CASE_INSENSITIVE",
															"EQUALS_DN",
															"NOT_EQUAL",
															"NOT_EQUAL_CASE_INSENSITIVE",
															"NOT_EQUAL_DN",
															"MULTIVALUE_CONTAINS",
															"MULTIVALUE_CONTAINS_CASE_INSENSITIVE",
															"MULTIVALUE_CONTAINS_DN",
															"MULTIVALUE_DOES_NOT_CONTAIN",
															"MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE",
															"MULTIVALUE_DOES_NOT_CONTAIN_DN",
														),
													},
												},
												"expected_value": schema.StringAttribute{
													Optional:    true,
													Description: "The expected value of this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
												},
												"expression": schema.StringAttribute{
													Optional:    true,
													Description: "The expression of this attribute rule. This field is required if the Attribute Source type is 'EXPRESSION'.",
												},
												"result": schema.StringAttribute{
													Required:    true,
													Description: "The result of this attribute rule.",
												},
											},
										},
										Optional:    true,
										Description: "The actual list of attribute rules.",
									},
								},
								Optional:    true,
								Description: "A collection of attribute rules",
							},
							"authentication_source": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"source_ref": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required:    true,
												Description: "The ID of the resource.",
											},
											"location": schema.StringAttribute{
												Optional:    true,
												Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
											},
										},
										Required:    true,
										Description: "A reference to a resource.",
									},
									"type": schema.StringAttribute{
										Required:    true,
										Description: "The type of this authentication source.",
										Validators: []validator.String{
											stringvalidator.OneOf(
												"IDP_ADAPTER",
												"IDP_CONNECTION",
											),
										},
									},
								},
								Required:    true,
								Description: "An authentication source (IdP adapter or IdP connection).",
							},
							"input_user_id_mapping": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"source": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Optional:    true,
												Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
											},
											"type": schema.StringAttribute{
												Required:    true,
												Description: "The source type of this key.",
												Validators: []validator.String{
													stringvalidator.OneOf(
														"TOKEN_EXCHANGE_PROCESSOR_POLICY",
														"ACCOUNT_LINK",
														"ADAPTER",
														"ASSERTION",
														"CONTEXT",
														"CUSTOM_DATA_STORE",
														"EXPRESSION",
														"JDBC_DATA_STORE",
														"LDAP_DATA_STORE",
														"PING_ONE_LDAP_GATEWAY_DATA_STORE",
														"MAPPED_ATTRIBUTES",
														"NO_MAPPING",
														"TEXT",
														"TOKEN",
														"REQUEST",
														"OAUTH_PERSISTENT_GRANT",
														"SUBJECT_TOKEN",
														"ACTOR_TOKEN",
														"PASSWORD_CREDENTIAL_VALIDATOR",
														"IDP_CONNECTION",
														"AUTHENTICATION_POLICY_CONTRACT",
														"CLAIMS",
														"LOCAL_IDENTITY_PROFILE",
														"EXTENDED_CLIENT_METADATA",
														"EXTENDED_PROPERTIES",
														"TRACKED_HTTP_PARAMS",
														"FRAGMENT",
														"INPUTS",
														"ATTRIBUTE_QUERY",
														"IDENTITY_STORE_USER",
														"IDENTITY_STORE_GROUP",
														"SCIM_USER",
														"SCIM_GROUP",
													),
												},
											},
										},
										Required:    true,
										Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
									},
									"value": schema.StringAttribute{
										Required:    true,
										Description: "The value for this attribute.",
									},
								},
								Optional:    true,
								Description: "Defines how an attribute in an attribute contract should be populated.",
							},
							"user_id_authenticated": schema.BoolAttribute{
								Optional:    true,
								Description: "Indicates whether the user ID obtained by the user ID mapping is authenticated.",
							},
						},
						Optional:    true,
						Description: "An authentication source selection action.",
					},
					"continue_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context": schema.StringAttribute{
								Optional:    true,
								Description: "The result context.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The authentication selection type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"APC_MAPPING",
										"LOCAL_IDENTITY_MAPPING",
										"AUTHN_SELECTOR",
										"AUTHN_SOURCE",
										"DONE",
										"CONTINUE",
										"RESTART",
										"FRAGMENT",
									),
								},
							},
						},
						Optional:    true,
						Description: "The continue selection action.",
					},
					"done_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context": schema.StringAttribute{
								Optional:    true,
								Description: "The result context.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The authentication selection type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"APC_MAPPING",
										"LOCAL_IDENTITY_MAPPING",
										"AUTHN_SELECTOR",
										"AUTHN_SOURCE",
										"DONE",
										"CONTINUE",
										"RESTART",
										"FRAGMENT",
									),
								},
							},
						},
						Optional:    true,
						Description: "The done selection action.",
					},
					"fragment_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"attribute_rules": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"fallback_to_success": schema.BoolAttribute{
										Optional:    true,
										Description: "When all the rules fail, you may choose to default to the general success action or fail. Default to success.",
									},
									"items": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"attribute_name": schema.StringAttribute{
													Optional:    true,
													Description: "The name of the attribute to use in this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
												},
												"attribute_source": schema.SingleNestedAttribute{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Optional:    true,
															Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
														},
														"type": schema.StringAttribute{
															Required:    true,
															Description: "The source type of this key.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																	"ACCOUNT_LINK",
																	"ADAPTER",
																	"ASSERTION",
																	"CONTEXT",
																	"CUSTOM_DATA_STORE",
																	"EXPRESSION",
																	"JDBC_DATA_STORE",
																	"LDAP_DATA_STORE",
																	"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																	"MAPPED_ATTRIBUTES",
																	"NO_MAPPING",
																	"TEXT",
																	"TOKEN",
																	"REQUEST",
																	"OAUTH_PERSISTENT_GRANT",
																	"SUBJECT_TOKEN",
																	"ACTOR_TOKEN",
																	"PASSWORD_CREDENTIAL_VALIDATOR",
																	"IDP_CONNECTION",
																	"AUTHENTICATION_POLICY_CONTRACT",
																	"CLAIMS",
																	"LOCAL_IDENTITY_PROFILE",
																	"EXTENDED_CLIENT_METADATA",
																	"EXTENDED_PROPERTIES",
																	"TRACKED_HTTP_PARAMS",
																	"FRAGMENT",
																	"INPUTS",
																	"ATTRIBUTE_QUERY",
																	"IDENTITY_STORE_USER",
																	"IDENTITY_STORE_GROUP",
																	"SCIM_USER",
																	"SCIM_GROUP",
																),
															},
														},
													},
													Optional:    true,
													Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
												},
												"condition": schema.StringAttribute{
													Optional:    true,
													Description: "The condition that will be applied to the attribute's expected value. This field is required if the Attribute Source type is not 'EXPRESSION'.",
													Validators: []validator.String{
														stringvalidator.OneOf(
															"EQUALS",
															"EQUALS_CASE_INSENSITIVE",
															"EQUALS_DN",
															"NOT_EQUAL",
															"NOT_EQUAL_CASE_INSENSITIVE",
															"NOT_EQUAL_DN",
															"MULTIVALUE_CONTAINS",
															"MULTIVALUE_CONTAINS_CASE_INSENSITIVE",
															"MULTIVALUE_CONTAINS_DN",
															"MULTIVALUE_DOES_NOT_CONTAIN",
															"MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE",
															"MULTIVALUE_DOES_NOT_CONTAIN_DN",
														),
													},
												},
												"expected_value": schema.StringAttribute{
													Optional:    true,
													Description: "The expected value of this attribute rule. This field is required if the Attribute Source type is not 'EXPRESSION'.",
												},
												"expression": schema.StringAttribute{
													Optional:    true,
													Description: "The expression of this attribute rule. This field is required if the Attribute Source type is 'EXPRESSION'.",
												},
												"result": schema.StringAttribute{
													Required:    true,
													Description: "The result of this attribute rule.",
												},
											},
										},
										Optional:    true,
										Description: "The actual list of attribute rules.",
									},
								},
								Optional:    true,
								Description: "A collection of attribute rules",
							},
							"context": schema.StringAttribute{
								Optional:    true,
								Description: "The result context.",
							},
							"fragment": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:    true,
										Description: "The ID of the resource.",
									},
									"location": schema.StringAttribute{
										Optional:    true,
										Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
									},
								},
								Required:    true,
								Description: "A reference to a resource.",
							},
							"fragment_mapping": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract_fulfillment": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:    true,
														Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:    true,
														Description: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:    true,
												Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "The value for this attribute.",
											},
										},
										Required:    true,
										Description: "Defines how an attribute in an attribute contract should be populated.",
									},
									"attribute_sources": schema.StringAttribute{
										Optional:    true,
										Description: "A list of configured data stores to look up attributes from.",
									},
									"issuance_criteria": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"conditional_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"attribute_name": schema.StringAttribute{
															Required:    true,
															Description: "The name of the attribute to use in this issuance criterion.",
														},
														"condition": schema.StringAttribute{
															Required:    true,
															Description: "The condition that will be applied to the source attribute's value and the expected value.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"EQUALS",
																	"EQUALS_CASE_INSENSITIVE",
																	"EQUALS_DN",
																	"NOT_EQUAL",
																	"NOT_EQUAL_CASE_INSENSITIVE",
																	"NOT_EQUAL_DN",
																	"MULTIVALUE_CONTAINS",
																	"MULTIVALUE_CONTAINS_CASE_INSENSITIVE",
																	"MULTIVALUE_CONTAINS_DN",
																	"MULTIVALUE_DOES_NOT_CONTAIN",
																	"MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE",
																	"MULTIVALUE_DOES_NOT_CONTAIN_DN",
																),
															},
														},
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"source": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"id": schema.StringAttribute{
																	Optional:    true,
																	Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																},
																"type": schema.StringAttribute{
																	Required:    true,
																	Description: "The source type of this key.",
																	Validators: []validator.String{
																		stringvalidator.OneOf(
																			"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																			"ACCOUNT_LINK",
																			"ADAPTER",
																			"ASSERTION",
																			"CONTEXT",
																			"CUSTOM_DATA_STORE",
																			"EXPRESSION",
																			"JDBC_DATA_STORE",
																			"LDAP_DATA_STORE",
																			"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																			"MAPPED_ATTRIBUTES",
																			"NO_MAPPING",
																			"TEXT",
																			"TOKEN",
																			"REQUEST",
																			"OAUTH_PERSISTENT_GRANT",
																			"SUBJECT_TOKEN",
																			"ACTOR_TOKEN",
																			"PASSWORD_CREDENTIAL_VALIDATOR",
																			"IDP_CONNECTION",
																			"AUTHENTICATION_POLICY_CONTRACT",
																			"CLAIMS",
																			"LOCAL_IDENTITY_PROFILE",
																			"EXTENDED_CLIENT_METADATA",
																			"EXTENDED_PROPERTIES",
																			"TRACKED_HTTP_PARAMS",
																			"FRAGMENT",
																			"INPUTS",
																			"ATTRIBUTE_QUERY",
																			"IDENTITY_STORE_USER",
																			"IDENTITY_STORE_GROUP",
																			"SCIM_USER",
																			"SCIM_GROUP",
																		),
																	},
																},
															},
															Required:    true,
															Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														},
														"value": schema.StringAttribute{
															Required:    true,
															Description: "The expected value of this issuance criterion.",
														},
													},
												},
												Optional:    true,
												Description: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											},
											"expression_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"expression": schema.StringAttribute{
															Required:    true,
															Description: "The OGNL expression to evaluate.",
														},
													},
												},
												Optional:    true,
												Description: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											},
										},
										Optional:    true,
										Description: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									},
								},
								Optional:    true,
								Description: "A list of mappings from attribute sources to attribute targets.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The authentication selection type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"APC_MAPPING",
										"LOCAL_IDENTITY_MAPPING",
										"AUTHN_SELECTOR",
										"AUTHN_SOURCE",
										"DONE",
										"CONTINUE",
										"RESTART",
										"FRAGMENT",
									),
								},
							},
						},
						Optional:    true,
						Description: "A authentication policy fragment selection action.",
					},
					"local_identity_mapping_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context": schema.StringAttribute{
								Optional:    true,
								Description: "The result context.",
							},
							"inbound_mapping": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract_fulfillment": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:    true,
														Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:    true,
														Description: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:    true,
												Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "The value for this attribute.",
											},
										},
										Required:    true,
										Description: "Defines how an attribute in an attribute contract should be populated.",
									},
									"attribute_sources": schema.StringAttribute{
										Optional:    true,
										Description: "A list of configured data stores to look up attributes from.",
									},
									"issuance_criteria": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"conditional_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"attribute_name": schema.StringAttribute{
															Required:    true,
															Description: "The name of the attribute to use in this issuance criterion.",
														},
														"condition": schema.StringAttribute{
															Required:    true,
															Description: "The condition that will be applied to the source attribute's value and the expected value.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"EQUALS",
																	"EQUALS_CASE_INSENSITIVE",
																	"EQUALS_DN",
																	"NOT_EQUAL",
																	"NOT_EQUAL_CASE_INSENSITIVE",
																	"NOT_EQUAL_DN",
																	"MULTIVALUE_CONTAINS",
																	"MULTIVALUE_CONTAINS_CASE_INSENSITIVE",
																	"MULTIVALUE_CONTAINS_DN",
																	"MULTIVALUE_DOES_NOT_CONTAIN",
																	"MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE",
																	"MULTIVALUE_DOES_NOT_CONTAIN_DN",
																),
															},
														},
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"source": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"id": schema.StringAttribute{
																	Optional:    true,
																	Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																},
																"type": schema.StringAttribute{
																	Required:    true,
																	Description: "The source type of this key.",
																	Validators: []validator.String{
																		stringvalidator.OneOf(
																			"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																			"ACCOUNT_LINK",
																			"ADAPTER",
																			"ASSERTION",
																			"CONTEXT",
																			"CUSTOM_DATA_STORE",
																			"EXPRESSION",
																			"JDBC_DATA_STORE",
																			"LDAP_DATA_STORE",
																			"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																			"MAPPED_ATTRIBUTES",
																			"NO_MAPPING",
																			"TEXT",
																			"TOKEN",
																			"REQUEST",
																			"OAUTH_PERSISTENT_GRANT",
																			"SUBJECT_TOKEN",
																			"ACTOR_TOKEN",
																			"PASSWORD_CREDENTIAL_VALIDATOR",
																			"IDP_CONNECTION",
																			"AUTHENTICATION_POLICY_CONTRACT",
																			"CLAIMS",
																			"LOCAL_IDENTITY_PROFILE",
																			"EXTENDED_CLIENT_METADATA",
																			"EXTENDED_PROPERTIES",
																			"TRACKED_HTTP_PARAMS",
																			"FRAGMENT",
																			"INPUTS",
																			"ATTRIBUTE_QUERY",
																			"IDENTITY_STORE_USER",
																			"IDENTITY_STORE_GROUP",
																			"SCIM_USER",
																			"SCIM_GROUP",
																		),
																	},
																},
															},
															Required:    true,
															Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														},
														"value": schema.StringAttribute{
															Required:    true,
															Description: "The expected value of this issuance criterion.",
														},
													},
												},
												Optional:    true,
												Description: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											},
											"expression_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"expression": schema.StringAttribute{
															Required:    true,
															Description: "The OGNL expression to evaluate.",
														},
													},
												},
												Optional:    true,
												Description: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											},
										},
										Optional:    true,
										Description: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									},
								},
								Optional:    true,
								Description: "A list of mappings from attribute sources to attribute targets.",
							},
							"local_identity_ref": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:    true,
										Description: "The ID of the resource.",
									},
									"location": schema.StringAttribute{
										Optional:    true,
										Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
									},
								},
								Required:    true,
								Description: "A reference to a resource.",
							},
							"outbound_attribute_mapping": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"attribute_contract_fulfillment": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"source": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Optional:    true,
														Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													},
													"type": schema.StringAttribute{
														Required:    true,
														Description: "The source type of this key.",
														Validators: []validator.String{
															stringvalidator.OneOf(
																"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																"ACCOUNT_LINK",
																"ADAPTER",
																"ASSERTION",
																"CONTEXT",
																"CUSTOM_DATA_STORE",
																"EXPRESSION",
																"JDBC_DATA_STORE",
																"LDAP_DATA_STORE",
																"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																"MAPPED_ATTRIBUTES",
																"NO_MAPPING",
																"TEXT",
																"TOKEN",
																"REQUEST",
																"OAUTH_PERSISTENT_GRANT",
																"SUBJECT_TOKEN",
																"ACTOR_TOKEN",
																"PASSWORD_CREDENTIAL_VALIDATOR",
																"IDP_CONNECTION",
																"AUTHENTICATION_POLICY_CONTRACT",
																"CLAIMS",
																"LOCAL_IDENTITY_PROFILE",
																"EXTENDED_CLIENT_METADATA",
																"EXTENDED_PROPERTIES",
																"TRACKED_HTTP_PARAMS",
																"FRAGMENT",
																"INPUTS",
																"ATTRIBUTE_QUERY",
																"IDENTITY_STORE_USER",
																"IDENTITY_STORE_GROUP",
																"SCIM_USER",
																"SCIM_GROUP",
															),
														},
													},
												},
												Required:    true,
												Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "The value for this attribute.",
											},
										},
										Required:    true,
										Description: "Defines how an attribute in an attribute contract should be populated.",
									},
									"attribute_sources": schema.StringAttribute{
										Optional:    true,
										Description: "A list of configured data stores to look up attributes from.",
									},
									"issuance_criteria": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"conditional_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"attribute_name": schema.StringAttribute{
															Required:    true,
															Description: "The name of the attribute to use in this issuance criterion.",
														},
														"condition": schema.StringAttribute{
															Required:    true,
															Description: "The condition that will be applied to the source attribute's value and the expected value.",
															Validators: []validator.String{
																stringvalidator.OneOf(
																	"EQUALS",
																	"EQUALS_CASE_INSENSITIVE",
																	"EQUALS_DN",
																	"NOT_EQUAL",
																	"NOT_EQUAL_CASE_INSENSITIVE",
																	"NOT_EQUAL_DN",
																	"MULTIVALUE_CONTAINS",
																	"MULTIVALUE_CONTAINS_CASE_INSENSITIVE",
																	"MULTIVALUE_CONTAINS_DN",
																	"MULTIVALUE_DOES_NOT_CONTAIN",
																	"MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE",
																	"MULTIVALUE_DOES_NOT_CONTAIN_DN",
																),
															},
														},
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"source": schema.SingleNestedAttribute{
															Attributes: map[string]schema.Attribute{
																"id": schema.StringAttribute{
																	Optional:    true,
																	Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																},
																"type": schema.StringAttribute{
																	Required:    true,
																	Description: "The source type of this key.",
																	Validators: []validator.String{
																		stringvalidator.OneOf(
																			"TOKEN_EXCHANGE_PROCESSOR_POLICY",
																			"ACCOUNT_LINK",
																			"ADAPTER",
																			"ASSERTION",
																			"CONTEXT",
																			"CUSTOM_DATA_STORE",
																			"EXPRESSION",
																			"JDBC_DATA_STORE",
																			"LDAP_DATA_STORE",
																			"PING_ONE_LDAP_GATEWAY_DATA_STORE",
																			"MAPPED_ATTRIBUTES",
																			"NO_MAPPING",
																			"TEXT",
																			"TOKEN",
																			"REQUEST",
																			"OAUTH_PERSISTENT_GRANT",
																			"SUBJECT_TOKEN",
																			"ACTOR_TOKEN",
																			"PASSWORD_CREDENTIAL_VALIDATOR",
																			"IDP_CONNECTION",
																			"AUTHENTICATION_POLICY_CONTRACT",
																			"CLAIMS",
																			"LOCAL_IDENTITY_PROFILE",
																			"EXTENDED_CLIENT_METADATA",
																			"EXTENDED_PROPERTIES",
																			"TRACKED_HTTP_PARAMS",
																			"FRAGMENT",
																			"INPUTS",
																			"ATTRIBUTE_QUERY",
																			"IDENTITY_STORE_USER",
																			"IDENTITY_STORE_GROUP",
																			"SCIM_USER",
																			"SCIM_GROUP",
																		),
																	},
																},
															},
															Required:    true,
															Description: "A key that is meant to reference a source from which an attribute can be retrieved. This model is usually paired with a value which, depending on the SourceType, can be a hardcoded value or a reference to an attribute name specific to that SourceType. Not all values are applicable - a validation error will be returned for incorrect values.<br>For each SourceType, the value should be:<br>ACCOUNT_LINK - If account linking was enabled for the browser SSO, the value must be 'Local User ID', unless it has been overridden in PingFederate's server configuration.<br>ADAPTER - The value is one of the attributes of the IdP Adapter.<br>ASSERTION - The value is one of the attributes coming from the SAML assertion.<br>AUTHENTICATION_POLICY_CONTRACT - The value is one of the attributes coming from an authentication policy contract.<br>LOCAL_IDENTITY_PROFILE - The value is one of the fields coming from a local identity profile.<br>CONTEXT - The value must be one of the following ['TargetResource' or 'OAuthScopes' or 'ClientId' or 'AuthenticationCtx' or 'ClientIp' or 'Locale' or 'StsBasicAuthUsername' or 'StsSSLClientCertSubjectDN' or 'StsSSLClientCertChain' or 'VirtualServerId' or 'AuthenticatingAuthority' or 'DefaultPersistentGrantLifetime'.]<br>CLAIMS - Attributes provided by the OIDC Provider.<br>CUSTOM_DATA_STORE - The value is one of the attributes returned by this custom data store.<br>EXPRESSION - The value is an OGNL expression.<br>EXTENDED_CLIENT_METADATA - The value is from an OAuth extended client metadata parameter. This source type is deprecated and has been replaced by EXTENDED_PROPERTIES.<br>EXTENDED_PROPERTIES - The value is from an OAuth Client's extended property.<br>IDP_CONNECTION - The value is one of the attributes passed in by the IdP connection.<br>JDBC_DATA_STORE - The value is one of the column names returned from the JDBC attribute source.<br>LDAP_DATA_STORE - The value is one of the LDAP attributes supported by your LDAP data store.<br>MAPPED_ATTRIBUTES - The value is the name of one of the mapped attributes that is defined in the associated attribute mapping.<br>OAUTH_PERSISTENT_GRANT - The value is one of the attributes from the persistent grant.<br>PASSWORD_CREDENTIAL_VALIDATOR - The value is one of the attributes of the PCV.<br>NO_MAPPING - A placeholder value to indicate that an attribute currently has no mapped source.TEXT - A hardcoded value that is used to populate the corresponding attribute.<br>TOKEN - The value is one of the token attributes.<br>REQUEST - The value is from the request context such as the CIBA identity hint contract or the request contract for Ws-Trust.<br>TRACKED_HTTP_PARAMS - The value is from the original request parameters.<br>SUBJECT_TOKEN - The value is one of the OAuth 2.0 Token exchange subject_token attributes.<br>ACTOR_TOKEN - The value is one of the OAuth 2.0 Token exchange actor_token attributes.<br>TOKEN_EXCHANGE_PROCESSOR_POLICY - The value is one of the attributes coming from a Token Exchange Processor policy.<br>FRAGMENT - The value is one of the attributes coming from an authentication policy fragment.<br>INPUTS - The value is one of the attributes coming from an attribute defined in the input authentication policy contract for an authentication policy fragment.<br>ATTRIBUTE_QUERY - The value is one of the user attributes queried from an Attribute Authority.<br>IDENTITY_STORE_USER - The value is one of the attributes from a user identity store provisioner for SCIM processing.<br>IDENTITY_STORE_GROUP - The value is one of the attributes from a group identity store provisioner for SCIM processing.<br>SCIM_USER - The value is one of the attributes passed in from the SCIM user request.<br>SCIM_GROUP - The value is one of the attributes passed in from the SCIM group request.<br>",
														},
														"value": schema.StringAttribute{
															Required:    true,
															Description: "The expected value of this issuance criterion.",
														},
													},
												},
												Optional:    true,
												Description: "A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.",
											},
											"expression_criteria": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"error_result": schema.StringAttribute{
															Optional:    true,
															Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
														},
														"expression": schema.StringAttribute{
															Required:    true,
															Description: "The OGNL expression to evaluate.",
														},
													},
												},
												Optional:    true,
												Description: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
											},
										},
										Optional:    true,
										Description: "A list of criteria that determines whether a transaction (usually a SSO transaction) is continued. All criteria must pass in order for the transaction to continue.",
									},
								},
								Required:    true,
								Description: "A list of mappings from attribute sources to attribute targets.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The authentication selection type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"APC_MAPPING",
										"LOCAL_IDENTITY_MAPPING",
										"AUTHN_SELECTOR",
										"AUTHN_SOURCE",
										"DONE",
										"CONTINUE",
										"RESTART",
										"FRAGMENT",
									),
								},
							},
						},
						Optional:    true,
						Description: "A local identity profile selection action.",
					},
					"restart_policy_action": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"context": schema.StringAttribute{
								Optional:    true,
								Description: "The result context.",
							},
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The authentication selection type.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"APC_MAPPING",
										"LOCAL_IDENTITY_MAPPING",
										"AUTHN_SELECTOR",
										"AUTHN_SOURCE",
										"DONE",
										"CONTINUE",
										"RESTART",
										"FRAGMENT",
									),
								},
							},
						},
						Optional:    true,
						Description: "The restart selection action.",
					},
				},
				Optional:    true,
				Description: "An authentication policy tree node.",
			},
		},
	}
}
