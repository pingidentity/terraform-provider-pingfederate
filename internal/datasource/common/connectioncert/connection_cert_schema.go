// Copyright © 2025 Ping Identity Corporation

package connectioncert

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ToSchemaDataSource(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description:         description,
		MarkdownDescription: description,
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: toSchemaDataSourceAttributes(),
		},
	}
}

func toSchemaDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cert_view": schema.SingleNestedAttribute{
			Description: "Certificate details.",
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
				},
			},
		},
		"x509_file": schema.SingleNestedAttribute{
			Description: "Encoded certificate data.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed:    true,
					Description: "The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified.",
				},
				"file_data": schema.StringAttribute{
					Description: "The certificate data in PEM format. New line characters should be omitted or encoded in this value.",
					Computed:    true,
				},
				"crypto_provider": schema.StringAttribute{
					Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Optional values are `HSM` and `LOCAL`.",
					Computed:    true,
				},
			},
		},
		"active_verification_cert": schema.BoolAttribute{
			Description: "Indicates whether this is an active signature verification certificate. Default is `false`.",
			Computed:    true,
		},
		"primary_verification_cert": schema.BoolAttribute{
			Description: "Indicates whether this is the primary signature verification certificate. Only one certificate in the collection can have this flag set. Default is `false`.",
			Computed:    true,
		},
		"secondary_verification_cert": schema.BoolAttribute{
			Description: "Indicates whether this is the secondary signature verification certificate. Only one certificate in the collection can have this flag set. Default is `false`.",
			Computed:    true,
		},
		"encryption_cert": schema.BoolAttribute{
			Description: "Indicates whether to use this cert to encrypt outgoing assertions. Only one certificate in the collection can have this flag set. Default is `false`.",
			Computed:    true,
		},
	}
}
