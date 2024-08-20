package connectioncert

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
)

func ObjType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: AttrTypes(),
	}
}

func CertAttrType() map[string]attr.Type {
	return map[string]attr.Type{
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
	}
}

func X509FileAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                  types.StringType,
		"file_data":           types.StringType,
		"formatted_file_data": types.StringType,
		"crypto_provider":     types.StringType,
	}
}

func AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cert_view": types.ObjectType{AttrTypes: CertAttrType()},
		"x509_file": types.ObjectType{
			AttrTypes: X509FileAttrType(),
		},
		"active_verification_cert":    types.BoolType,
		"primary_verification_cert":   types.BoolType,
		"secondary_verification_cert": types.BoolType,
		"encryption_cert":             types.BoolType,
	}
}

func ToState(ctx context.Context, planFileData types.String, clientConnectionCert client.ConnectionCert, diags *diag.Diagnostics) (types.Object, diag.Diagnostics) {
	var certViewValue types.Object
	if clientConnectionCert.CertView == nil {
		certViewValue = types.ObjectNull(CertAttrType())
	} else {
		certViewSubjectAlternativeNamesValue, objDiags := types.SetValueFrom(ctx, types.StringType, clientConnectionCert.CertView.SubjectAlternativeNames)
		diags.Append(objDiags...)

		certViewAttrValues := map[string]attr.Value{
			"crypto_provider":           types.StringPointerValue(clientConnectionCert.CertView.CryptoProvider),
			"expires":                   types.StringValue(clientConnectionCert.CertView.Expires.Format(time.RFC3339)),
			"id":                        types.StringPointerValue(clientConnectionCert.CertView.Id),
			"issuer_dn":                 types.StringPointerValue(clientConnectionCert.CertView.IssuerDN),
			"key_algorithm":             types.StringPointerValue(clientConnectionCert.CertView.KeyAlgorithm),
			"key_size":                  types.Int64PointerValue(clientConnectionCert.CertView.KeySize),
			"serial_number":             types.StringPointerValue(clientConnectionCert.CertView.SerialNumber),
			"sha1_fingerprint":          types.StringPointerValue(clientConnectionCert.CertView.Sha1Fingerprint),
			"sha256_fingerprint":        types.StringPointerValue(clientConnectionCert.CertView.Sha256Fingerprint),
			"signature_algorithm":       types.StringPointerValue(clientConnectionCert.CertView.SignatureAlgorithm),
			"status":                    types.StringPointerValue(clientConnectionCert.CertView.Status),
			"subject_alternative_names": certViewSubjectAlternativeNamesValue,
			"subject_dn":                types.StringPointerValue(clientConnectionCert.CertView.SubjectDN),
			"valid_from":                types.StringValue(clientConnectionCert.CertView.ValidFrom.Format(time.RFC3339)),
			"version":                   types.Int64PointerValue(clientConnectionCert.CertView.Version),
		}

		certViewValue, objDiags = types.ObjectValue(CertAttrType(), certViewAttrValues)
		diags.Append(objDiags...)
	}

	// Get the current file_data value
	fileDataAttr := types.StringNull()
	if planFileData.ValueString() != "" {
		fileDataAttr = planFileData
	}

	var objDiags diag.Diagnostics
	certsX509fileValue, objDiags := types.ObjectValue(X509FileAttrType(), map[string]attr.Value{
		"crypto_provider":     types.StringPointerValue(clientConnectionCert.X509File.CryptoProvider),
		"formatted_file_data": types.StringValue(clientConnectionCert.X509File.FileData),
		"file_data":           fileDataAttr,
		"id":                  types.StringPointerValue(clientConnectionCert.X509File.Id),
	})
	diags.Append(objDiags...)

	return types.ObjectValue(AttrTypes(), map[string]attr.Value{
		"active_verification_cert":    types.BoolPointerValue(clientConnectionCert.ActiveVerificationCert),
		"cert_view":                   certViewValue,
		"encryption_cert":             types.BoolPointerValue(clientConnectionCert.EncryptionCert),
		"primary_verification_cert":   types.BoolPointerValue(clientConnectionCert.PrimaryVerificationCert),
		"secondary_verification_cert": types.BoolPointerValue(clientConnectionCert.SecondaryVerificationCert),
		"x509_file":                   certsX509fileValue,
	})
}
