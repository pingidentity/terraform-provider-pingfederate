package connectioncert

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ToSchema(description string, required, computed bool) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Description:         description,
		MarkdownDescription: description,
		Required:            required,
		Optional:            !required,
		Computed:            computed,
		NestedObject: schema.NestedAttributeObject{
			Attributes: ToSchemaAttributes(),
		},
	}
}

func ToSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cert_view": schema.SetNestedAttribute{
			Description: "A certificate used for signature verification or XML encryption.",
			Required:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The persistent, unique ID for the certificate.",
						Optional:    true,
					},
					"serial_number": schema.StringAttribute{
						Description: "The serial number assigned by the CA.",
						Optional:    true,
					},
					"subject_dn": schema.StringAttribute{
						Description: "The subject's distinguished name.",
						Optional:    true,
					},
					"subject_alternative_names": schema.SetAttribute{
						ElementType: types.StringType,
						Description: "The subject alternative names (SAN).",
						Optional:    true,
					},
					"issuer_dn": schema.StringAttribute{
						Description: "The issuer's distinguished name.",
						Optional:    true,
					},
					"valid_from": schema.StringAttribute{
						Description: "The start date from which the item is valid, in ISO 8601 format (UTC).",
						Optional:    true,
					},
					"expires": schema.StringAttribute{
						Description: "The end date up until which the item is valid, in ISO 8601 format (UTC).",
						Optional:    true,
					},
					"key_algorithm": schema.StringAttribute{
						Description: "The public key algorithm.",
						Optional:    true,
					},
					"key_size": schema.Int64Attribute{
						Description: "The public key size.",
						Optional:    true,
					},
					"signature_algorithm": schema.StringAttribute{
						Description: "The signature algorithm.",
						Optional:    true,
					},
					"version": schema.Int64Attribute{
						Description: "The X.509 version to which the item conforms.",
						Optional:    true,
					},
					"sha1_fingerprint": schema.StringAttribute{
						Description: "SHA-1 fingerprint in Hex encoding.",
						Optional:    true,
					},
					"sha256_fingerprint": schema.StringAttribute{
						Description: "SHA-256 fingerprint in Hex encoding.",
						Optional:    true,
					},
					"status": schema.StringAttribute{
						Description: "Status of the item.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("VALID", "EXPIRED", "NOT_YET_VALID", "REVOKED"),
						},
					},
					"crypto_provider": schema.StringAttribute{
						Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("HSM", "LOCAL"),
						},
					},
				},
			},
		},
		"x509_file": schema.SingleNestedAttribute{
			Description: "Encoded certificate data.",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description: "The persistent, unique ID for the certificate. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
					Optional:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"file_data": schema.StringAttribute{
					Description: "The certificate data in PEM format. New line characters should be omitted or encoded in this value.",
					Required:    true,
				},
				"crypto_provider": schema.StringAttribute{
					Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("HSM", "LOCAL"),
					},
				},
			},
		},
		"active_verification_cert": schema.BoolAttribute{
			Description: "Indicates if this is the active verification certificate.",
			Optional:    true,
		},
		"primary_verification_cert": schema.BoolAttribute{
			Description: "Indicates if this is the primary verification certificate.",
			Optional:    true,
		},
		"secondary_verification_cert": schema.BoolAttribute{
			Description: "Indicates if this is the secondary verification certificate.",
			Optional:    true,
		},
		"encryption_cert": schema.BoolAttribute{
			Description: "Indicates if this is the encryption certificate.",
			Optional:    true,
		},
	}
}
