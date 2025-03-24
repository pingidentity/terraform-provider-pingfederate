#!/usr/bin/env bash

# This script is used to set the developer profile for the CLI
# It is sourced by the CLI to set the developer profile

profile_name=pingfederate-terraform-dev

set -e

list_profiles=$(pingcli config list-profiles --output-format json)
profile_exists=$(echo "[$list_profiles]" | jq --arg profile_name "$profile_name" 'any(.[]; .Message | contains($profile_name))')

# If the profile exists, prompt whether to overwrite it
if [ "$profile_exists" = "true" ]; then
  read -p "The \"$profile_name\" profile already exists. Do you want to overwrite it? (y/n) " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Removing existing profile..."
    pingcli config delete-profile -y $profile_name
  else
    echo "Exiting..."
    exit 0
  fi
fi

echo "Creating new Ping CLI profile \"$profile_name\"..."
pingcli config add-profile --name $profile_name --description "Developer profile for PingFederate Terraform"
pingcli config list-profiles

echo "Adding profile configuration..."
pingcli config set --profile $profile_name service.pingfederate.httpsHost=https://localhost:9999
pingcli config set --profile $profile_name service.pingfederate.insecureTrustAllTLS=true
pingcli config set --profile $profile_name service.pingfederate.xBypassExternalValidationHeader=true
pingcli config set --profile $profile_name service.pingfederate.authentication.type=basicAuth
pingcli config set --profile $profile_name service.pingfederate.authentication.basicAuth.username=administrator
pingcli config set --profile $profile_name service.pingfederate.authentication.basicAuth.password=2FederateM0re
pingcli config set --profile $profile_name export.services=pingfederate
pingcli config set --profile $profile_name export.outputDirectory=`pwd`/tf-export

echo "PingFederate service configuration..."
pingcli config get --profile $profile_name service.pingfederate

echo "Export command configuration..."
pingcli config get --profile $profile_name export