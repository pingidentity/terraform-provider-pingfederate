SHELL := /bin/bash

.PHONY: install generate fmt vet test starttestcontainer removetestcontainer spincontainer clearstates kaboom testacc testacccomplete generateresource openlocalwebapi golangcilint tfproviderlint tflint terrafmtlint importfmtlint devcheck devchecknotest openapp testoneacc verifycontent

default: install

install:
	go mod tidy
	go install .

generate:
	go generate ./...
	go fmt ./...
	go vet ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

define productversiondir
 	PRODUCT_VERSION_DIR=$$(echo "$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-12.1.0}" | cut -b 1-4)
endef

#TODO: replace `edge` with `latest` here once a sprint release is available for PF 12.1
starttestcontainer:
	$(call productversiondir) && docker run --name pingfederate_terraform_provider_container \
		-d -p 9031:9031 \
		-p 9999:9999 \
		--env-file "${HOME}/.pingidentity/config" \
		-v $$(pwd)/server-profiles/shared-profile:/opt/in \
		-v $$(pwd)/server-profiles/$${PRODUCT_VERSION_DIR}/data.json.subst:/opt/in/instance/bulk-config/data.json.subst \
		pingidentity/pingfederate:$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-12.1.0}-edge
# Wait for the instance to become ready
	sleep 1
	duration=0
	while (( duration < 240 )) && ! docker logs pingfederate_terraform_provider_container 2>&1 | grep -q "Removing Imported Bulk File\|CONTAINER FAILURE"; \
	do \
	    duration=$$((duration+1)); \
		sleep 1; \
	done
# Fail if the container didn't become ready in time
	docker logs pingfederate_terraform_provider_container 2>&1 | grep -q "Removing Imported Bulk File" || \
		{ echo "PingFederate container did not become ready in time or contains errors. Logs:"; docker logs pingfederate_terraform_provider_container; exit 1; }

removetestcontainer:
	docker rm -f pingfederate_terraform_provider_container
	
spincontainer: removetestcontainer starttestcontainer

define test_acc_common_env_vars
	PINGFEDERATE_PROVIDER_HTTPS_HOST=https://localhost:9999 PINGFEDERATE_PROVIDER_ADMIN_API_PATH="/pf-admin-api/v1" PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS=true PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER=true PINGFEDERATE_PROVIDER_PRODUCT_VERSION=$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-12.1}
endef

define test_acc_basic_auth_env_vars
	PINGFEDERATE_PROVIDER_USERNAME=administrator PINGFEDERATE_PROVIDER_PASSWORD=2FederateM0re
endef

define test_acc_oauth_env_vars
	PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID=test PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET=2FederateM0re! PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL=https://localhost:9031/as/token.oauth2 PINGFEDERATE_PROVIDER_OAUTH_SCOPES=email
endef

# Set ACC_TEST_NAME to name of test in cli
testoneacc:
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test ./internal/acctest/... -timeout 10m -run ${ACC_TEST_NAME} -v -count=1

testaccfolder:
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test ./internal/acctest/config/${ACC_TEST_FOLDER}... -timeout 10m -v -count=1

testoneacccomplete: spincontainer testoneacc

# Some tests can step on each other's toes so run those tests in single threaded mode. Run the rest in parallel
testacc:
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test `go list ./internal/acctest/config... | grep -v -e authenticationapi -e oauth/authserversettings -e oauth/openidconnect/policy -e oauth/openidconnect/settings -e sp/targeturlmappings` -timeout 10m -v -p 4; \
	firstTestResult=$$?; \
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test `go list ./internal/acctest/config... | grep -e authenticationapi -e oauth/authserversettings -e oauth/openidconnect/policy -e oauth/openidconnect/settings -e sp/targeturlmappings` -timeout 10m -v -p 1; \
	secondTestResult=$$?; \
	if test "$$firstTestResult" != "0" || test "$$secondTestResult" != "0"; then \
		false; \
	fi

testauthacc:
	$(call test_acc_common_env_vars) $(call test_acc_oauth_env_vars) TF_ACC=1 go test ./internal/acctest/authentication/oauth_test.go -timeout 5m -v; \
	oauthResult=$$?; \
	$(call test_acc_common_env_vars) $(call test_acc_oauth_env_vars) TF_ACC=1 go test ./internal/acctest/authentication/access_token_test.go -timeout 5m -v; \
	atResult=$$?; \
	if test "$$oauthResult" != 0 || test "$$atResult" != 0; then \
		false; \
	fi

testacccomplete: spincontainer testauthacc testacc

clearstates:
	find . -name "*tfstate*" -delete
	
kaboom: clearstates spincontainer install

devchecknotest: verifycontent install golangcilint generate tfproviderlint tflint terrafmtlint importfmtlint

verifycontent:
	python3 ./scripts/verifyContent.py

devcheck: devchecknotest kaboom testauthacc testacc

generateresource:
	PINGFEDERATE_GENERATED_ENDPOINT=openIdConnectSettings \
	PINGFEDERATE_RESOURCE_DEFINITION_NAME=OpenIdConnectSettings \
	PINGFEDERATE_ALLOW_REQUIRED_BYPASS=False \
	OVERWRITE_EXISTING_RESOURCE_FILE=False \
	PINGFEDERATE_PUT_ONLY_RESOURCE=True \
	GENERATE_SCHEMA=True \
	python3 dev/generate_resource.py
	make fmt
	
openlocalwebapi:
	open "https://localhost:9999/pf-admin-api/api-docs/#/"

openapp:
	open "https://localhost:9999/pingfederate/app"

golangcilint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m ./internal/...

tfproviderlint: 
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlintx \
									-c 1 \
									-AT001.ignored-filename-suffixes=_test.go \
									-AT003=false \
									-XAT001=false \
									-XR004=false \
									-XS002=false ./internal/...

tflint:
	go run github.com/terraform-linters/tflint --recursive --disable-rule "terraform_unused_declarations" --disable-rule "terraform_required_providers" --disable-rule "terraform_required_version"

terrafmtlint:
	find ./internal/acctest -type f -name '*_test.go' \
		| sort -u \
		| xargs -I {} go run github.com/katbyte/terrafmt -f fmt {} -v

importfmtlint:
	go run github.com/pavius/impi/cmd/impi --local . --scheme stdThirdPartyLocal ./internal/...
