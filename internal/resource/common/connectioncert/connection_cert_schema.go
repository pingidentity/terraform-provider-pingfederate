package connectioncert

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/planmodifiers"
)

func ToSchema(description string, required, computed bool) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description:         description,
		MarkdownDescription: description,
		Required:            required,
		Optional:            !required,
		Computed:            computed,
		NestedObject: schema.NestedAttributeObject{
			Attributes: ToSchemaAttributes(),
		},
		Validators: []validator.List{
			listvalidator.UniqueValues(),
		},
	}
}

func ToSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cert_view": schema.SingleNestedAttribute{
			Description: "A certificate used for signature verification or XML encryption.",
			Optional:    false,
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: "The persistent, unique ID for the certificate.",
					Optional:    false,
					Computed:    true,
				},
				"serial_number": schema.StringAttribute{
					Description: "The serial number assigned by the CA.",
					Optional:    false,
					Computed:    true,
				},
				"subject_dn": schema.StringAttribute{
					Description: "The subject's distinguished name.",
					Optional:    false,
					Computed:    true,
				},
				"subject_alternative_names": schema.SetAttribute{
					ElementType: types.StringType,
					Description: "The subject alternative names (SAN).",
					Optional:    false,
					Computed:    true,
				},
				"issuer_dn": schema.StringAttribute{
					Description: "The issuer's distinguished name.",
					Optional:    false,
					Computed:    true,
				},
				"valid_from": schema.StringAttribute{
					Description: "The start date from which the item is valid, in ISO 8601 format (UTC).",
					Optional:    false,
					Computed:    true,
				},
				"expires": schema.StringAttribute{
					Description: "The end date up until which the item is valid, in ISO 8601 format (UTC).",
					Optional:    false,
					Computed:    true,
				},
				"key_algorithm": schema.StringAttribute{
					Description: "The public key algorithm.",
					Optional:    false,
					Computed:    true,
				},
				"key_size": schema.Int64Attribute{
					Description: "The public key size.",
					Optional:    false,
					Computed:    true,
				},
				"signature_algorithm": schema.StringAttribute{
					Description: "The signature algorithm.",
					Optional:    false,
					Computed:    true,
				},
				"version": schema.Int64Attribute{
					Description: "The X.509 version to which the item conforms.",
					Optional:    false,
					Computed:    true,
				},
				"sha1_fingerprint": schema.StringAttribute{
					Description: "SHA-1 fingerprint in Hex encoding.",
					Optional:    false,
					Computed:    true,
				},
				"sha256_fingerprint": schema.StringAttribute{
					Description: "SHA-256 fingerprint in Hex encoding.",
					Optional:    false,
					Computed:    true,
				},
				"status": schema.StringAttribute{
					Description: "Status of the item.",
					Optional:    false,
					Computed:    true,
				},
				"crypto_provider": schema.StringAttribute{
					Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
					Optional:    false,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
		"x509_file": schema.SingleNestedAttribute{
			Description: "Encoded certificate data.",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						configvalidators.LowercaseId(),
						stringvalidator.LengthAtLeast(1),
					},
				},
				"file_data": schema.StringAttribute{
					Description: "The certificate data in PEM format. New line characters should be omitted or encoded in this value.",
					Required:    true,
				},
				"formatted_file_data": schema.StringAttribute{
					Description: "The certificate data in PEM format, formatted by PingFederate. This attribute is read-only.",
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"crypto_provider": schema.StringAttribute{
					Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Optional values are `HSM` and `LOCAL`.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("HSM", "LOCAL"),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				planmodifiers.ValidateX509FileData(),
			},
		},
		"active_verification_cert": schema.BoolAttribute{
			Description: "Indicates if this is the active verification certificate. Default is false.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"primary_verification_cert": schema.BoolAttribute{
			Description: "Indicates if this is the primary verification certificate. Default is false.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"secondary_verification_cert": schema.BoolAttribute{
			Description: "Indicates if this is the secondary verification certificate. Default is false.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"encryption_cert": schema.BoolAttribute{
			Description: "Indicates if this is the encryption certificate. Default is false.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
	}
}
