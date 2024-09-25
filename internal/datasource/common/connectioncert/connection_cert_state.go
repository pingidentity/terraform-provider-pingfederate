package connectioncert

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AttrTypesDataSource() map[string]attr.Type {
	return map[string]attr.Type{
		"cert_view": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"crypto_provider":           types.StringType,
				"expires":                   types.StringType,
				"id":                        types.StringType,
				"issuer_dn":                 types.StringType,
				"key_algorithm":             types.StringType,
				"key_size":                  types.Int64Type,
				"serial_number":             types.StringType,
				"sha1_fingerprint":          types.StringType,
				"sha256_fingerprint":        types.StringType,
				"signature_algorithm":       types.StringType,
				"status":                    types.StringType,
				"subject_alternative_names": types.SetType{ElemType: types.StringType},
				"subject_dn":                types.StringType,
				"valid_from":                types.StringType,
				"version":                   types.Int64Type,
			},
		},
		"x509_file": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":              types.StringType,
				"file_data":       types.StringType,
				"crypto_provider": types.StringType,
			},
		},
		"active_verification_cert":    types.BoolType,
		"primary_verification_cert":   types.BoolType,
		"secondary_verification_cert": types.BoolType,
		"encryption_cert":             types.BoolType,
	}
}
