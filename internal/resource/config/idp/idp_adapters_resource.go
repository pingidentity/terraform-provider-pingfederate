package idp

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpAdaptersResource{}
	_ resource.ResourceWithConfigure   = &idpAdaptersResource{}
	_ resource.ResourceWithImportState = &idpAdaptersResource{}
)

// IdpAdaptersResource is a helper function to simplify the provider implementation.
func IdpAdaptersResource() resource.Resource {
	return &idpAdaptersResource{}
}

// idpAdaptersResource is the resource implementation.
type idpAdaptersResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type idpAdaptersResourceModel struct {
	//TODO left out of schema
	AuthnCtxClassRef    types.String `tfsdk:"authn_ctx_class_ref"`
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
	AttributeMapping    types.Object `tfsdk:"attribute_mapping"`
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
}

// GetSchema defines the schema for the resource.
func (r *idpAdaptersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	idpAdaptersResourceSchema(ctx, req, resp, false)
}

func idpAdaptersResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages Idp Adapters",
		Attributes: map[string]schema.Attribute{
			"authn_ctx_class_ref": schema.StringAttribute{
				Description: "The fixed value that indicates how the user was authenticated.",
				Optional:    true,
			},
			//TODO don't add id in common schema
			"id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Optional:    true,
					},
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Optional:    true,
					},
				},
			},

			"configuration": schema.SingleNestedAttribute{
				Description: "Plugin instance configuration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"tables": schema.SetNestedAttribute{
						Description: "List of configuration tables.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the table.",
									Required:    true,
								},
								"rows": schema.SetNestedAttribute{
									Description: "List of table rows.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"fields": schema.SetNestedAttribute{
												Description: "The configuration fields in the row.",
												Required:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Description: "The name of the configuration field.",
															Required:    true,
														},
														"value": schema.StringAttribute{
															Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
															Optional:    true,
														},
														"encrypted_value": schema.StringAttribute{
															Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. If you do not want to update the stored value, this attribute should be passed back unchanged.",
															Optional:    true,
														},
														"inherited": schema.BoolAttribute{
															Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
															Optional:    true,
														},
													},
												},
											},
											"default_row": schema.BoolAttribute{
												Description: "Whether this row is the default.",
												Optional:    true,
											},
										},
									},
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this table is inherited from its parent instance. If true, the rows become read-only. The default value is false.",
									Optional:    true,
								},
							},
						},
					},
					"fields": schema.SetNestedAttribute{
						Description: "List of configuration fields.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the configuration field.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
									Optional:    true,
								},
								"encrypted_value": schema.StringAttribute{
									Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. If you do not want to update the stored value, this attribute should be passed back unchanged.",
									Optional:    true,
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
									Optional:    true,
								},
							},
						},
					},
				},
			},

			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the IdP adapter provides.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
								},
							},
						},
					},
					"unique_user_key_attribute": schema.StringAttribute{
						Description: "The attribute to use for uniquely identify a user's authentication sessions.",
						Optional:    true,
					},
					"mask_ognl_values": schema.BoolAttribute{
						Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to false.",
						Optional:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
						Optional:    true,
					},
				},
			},

			/*     "attribute_mapping": schema.SingleNestedAttribute{
			    Description: "The attributes mapping from attribute sources to attribute targets.",
			    Optional: true,
			    Attributes: map[string]schema.Attribute{
			  "attribute_sources": schema.SetNestedAttribute{
			    Description: "A list of configured data stores to look up attributes from.",
			    Optional: true,
			    NestedObject: schema.NestedAttributeObject{
			      Attributes: map[string]schema.Attribute{
			    "type": schema.StringAttribute{
			      Description: "The data store type of this attribute source.",
			    Required: true,
			Validators: []validator.String{
			  stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"}...),
			},
			    },
			  "data_store_ref": schema.SingleNestedAttribute{
			    Description: "Reference to the associated data store.",
			    Required: true,
			    Attributes: map[string]schema.Attribute{
			    "id": schema.StringAttribute{
			      Description: "The ID of the resource.",
			    Required: true,
			    },
			    "location": schema.StringAttribute{
			      Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			    Optional: true,
			    },
			          },
			        },
			    "id": schema.StringAttribute{
			      Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
			    Optional: true,
			    },
			    "description": schema.StringAttribute{
			      Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
			    Optional: true,
			    },
			  "attribute_contract_fulfillment": schema.MapNestedAttribute{
			    Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
			    Optional: true,
			    NestedObject: schema.NestedAttributeObject{
			      Attributes: map[string]schema.Attribute{
			      },
			    },
			  "attribute_contract_fulfillment": schema.MapNestedAttribute{
			    Description: "A list of mappings from attribute names to their fulfillment values.",
			    Required: true,
			    NestedObject: schema.NestedAttributeObject{
			      Attributes: map[string]schema.Attribute{
			    "type": schema.StringAttribute{
			      Description: "The data store type of this attribute source.",
			    Required: true,
			Validators: []validator.String{
			  stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"}...),
			},
			    },
			  "data_store_ref": schema.SingleNestedAttribute{
			    Description: "Reference to the associated data store.",
			    Required: true,
			    Attributes: map[string]schema.Attribute{
			    "id": schema.StringAttribute{
			      Description: "The ID of the resource.",
			    Required: true,
			    },
			    "location": schema.StringAttribute{
			      Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			    Optional: true,
			    },
			          },
			        },
			    "id": schema.StringAttribute{
			      Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
			    Optional: true,
			    },
			    "description": schema.StringAttribute{
			      Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
			    Optional: true,
			    },
			  "attribute_contract_fulfillment": schema.MapNestedAttribute{
			    Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
			    Optional: true,
			    NestedObject: schema.NestedAttributeObject{
			      Attributes: map[string]schema.Attribute{
			      },
			    },
			  "issuance_criteria": schema.SingleNestedAttribute{
			    Description: "The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.",
			    Optional: true,
			    Attributes: map[string]schema.Attribute{
			    "type": schema.StringAttribute{
			      Description: "The data store type of this attribute source.",
			    Required: true,
			Validators: []validator.String{
			  stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"}...),
			},
			    },
			  "data_store_ref": schema.SingleNestedAttribute{
			    Description: "Reference to the associated data store.",
			    Required: true,
			    Attributes: map[string]schema.Attribute{
			    "id": schema.StringAttribute{
			      Description: "The ID of the resource.",
			    Required: true,
			    },
			    "location": schema.StringAttribute{
			      Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			    Optional: true,
			    },
			          },
			        },
			    "id": schema.StringAttribute{
			      Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
			    Optional: true,
			    },
			    "description": schema.StringAttribute{
			      Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
			    Optional: true,
			    },
			  "attribute_contract_fulfillment": schema.MapNestedAttribute{
			    Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
			    Optional: true,
			    NestedObject: schema.NestedAttributeObject{
			      Attributes: map[string]schema.Attribute{
			      },
			    },
			    "inherited": schema.BoolAttribute{
			      Description: "Whether this attribute mapping is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
			    Optional: true,
			    },
			          },
			        },*/

		},
	}

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"FIX_ME"})
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

func addOptionalIdpAdaptersFields(ctx context.Context, addRequest *client.IdpAdapter, plan idpAdaptersResourceModel) error {

	if internaltypes.IsDefined(plan.AuthnCtxClassRef) {
		addRequest.AuthnCtxClassRef = plan.AuthnCtxClassRef.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueString()
	}

	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueString()
	}

	if internaltypes.IsDefined(plan.PluginDescriptorRef) {
		addRequest.PluginDescriptorRef = client.NewPluginDescriptorRef()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef)), addRequest.PluginDescriptorRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ParentRef) {
		addRequest.ParentRef = client.NewParentRef()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ParentRef)), addRequest.ParentRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Configuration) {
		addRequest.Configuration = client.NewConfiguration()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration)), addRequest.Configuration)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeMapping) {
		addRequest.AttributeMapping = client.NewAttributeMapping()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeMapping)), addRequest.AttributeMapping)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = client.NewAttributeContract()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
	}

	return nil

}

// Metadata returns the resource type name.
func (r *idpAdaptersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_adapters"
}

func (r *idpAdaptersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readIdpAdaptersResponse(ctx context.Context, r *client.IdpAdapter, state *idpAdaptersResourceModel) {
	state.AuthnCtxClassRef = internaltypes.StringTypeOrNil(r.AuthnCtxClassRef)
	state.Id = internaltypes.StringTypeOrNil(r.Id)
	state.Name = internaltypes.StringTypeOrNil(r.Name)
	state.PluginDescriptorRef = (r.PluginDescriptorRef)
	state.ParentRef = (r.ParentRef)
	state.Configuration = (r.Configuration)
	state.AttributeMapping = (r.AttributeMapping)
	state.AttributeContract = (r.AttributeContract)
}

func (r *idpAdaptersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpAdaptersResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createIdpAdapters := client.NewIdpAdapter()
	err := addOptionalIdpAdaptersFields(ctx, createIdpAdapters, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapters", err.Error())
		return
	}
	requestJson, err := createIdpAdapters.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateIdpAdapters := r.apiClient.IdpAdaptersApi.CreateIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateIdpAdapters = apiCreateIdpAdapters.Body(*createIdpAdapters)
	idpAdaptersResponse, httpResp, err := r.apiClient.IdpAdaptersApi.CreateIdpAdapterExecute(apiCreateIdpAdapters)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdpAdapters", err, httpResp)
		return
	}
	responseJson, err := idpAdaptersResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state idpAdaptersResourceModel

	readIdpAdaptersResponse(ctx, idpAdaptersResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *idpAdaptersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdpAdapters(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIdpAdapters(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state idpAdaptersResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpAdapters, httpResp, err := apiClient.IdpAdaptersApi.GetIdpAdapter(config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a IdpAdapters", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadIdpAdapters.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIdpAdaptersResponse(ctx, apiReadIdpAdapters, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpAdaptersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdpAdapters(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIdpAdapters(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan idpAdaptersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state idpAdaptersResourceModel
	req.State.Get(ctx, &state)
	updateIdpAdapters := apiClient.IdpAdaptersApi.UpdateIdpAdapter(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewIdpAdapter() //TODO
	err := addOptionalIdpAdaptersFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapters", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateIdpAdapters = updateIdpAdapters.Body(*createUpdateRequest)
	updateIdpAdaptersResponse, httpResp, err := apiClient.IdpAdaptersApi.UpdateIdpAdapterExecute(updateIdpAdapters)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdpAdapters", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateIdpAdaptersResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readIdpAdaptersResponse(ctx, updateIdpAdaptersResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *idpAdaptersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *idpAdaptersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdpAdaptersLocation(ctx, req, resp)
}
func importIdpAdaptersLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
