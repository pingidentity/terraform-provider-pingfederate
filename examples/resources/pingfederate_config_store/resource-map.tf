resource "pingfederate_config_store" "signature_algorithms" {
  bundle     = "com.pingidentity.crypto.SignatureAlgorithms"
  setting_id = "signature-algorithms"
  map_value = {
    "DSA_SHA1" : "http://www.w3.org/2000/09/xmldsig#dsa-sha1"
    "RSA_SHA1" : "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
    "RSA_SHA256" : "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
    "RSA_SHA384" : "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384"
    "RSA_SHA512" : "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512"
    "ECDSA_SHA256" : "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"
    "ECDSA_SHA384" : "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha384"
    "ECDSA_SHA512" : "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha512"
  }
}