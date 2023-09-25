# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_server_settings_system_keys" "serverSettingsSystemKeysExample" {
  current = {
    encrypted_key_data = "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMS4yLjUuMCIsInppcCI6IkRFRiJ9..J1yaOm2OdYCUDN402iIKPQ.LlpjecXwfHDiFJl_K6O57Mzp1RZxHN-TAbpKnypkRfeL1XgTHZrUkPgxO3ZcU7fb.q-X1zzd-de5svqDRbAE0lw"
  }
  pending = {
    encrypted_key_data = "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMS4yLjUuMCIsInppcCI6IkRFRiJ9..4Q-LeikGMQ-5dVVRMMDyfw.JLR4Yg1FfmaTdOpVHZ1V1BypiguCuKawnJsUD33weL3nYRvyEPFgMCuBV72GC-HG.2b2T22iR040xI4ro-Iemeg"
  }
}