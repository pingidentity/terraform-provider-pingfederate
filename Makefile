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

starttestcontainer:
	docker run --name pingfederate_terraform_provider_container \
		-d -p 9031:9031 \
		-d -p 9999:9999 \
		--env-file "${HOME}/.pingidentity/config" \
		-e SERVER_PROFILE_URL=https://github.com/pingidentity/pingidentity-server-profiles.git \
		-e SERVER_PROFILE_BRANCH=terraform-provider-pingfederate-1125 \
		-e SERVER_PROFILE_PATH=terraform-provider-pingfederate/pingfederate \
		pingidentity/pingfederate:$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-12.0.0}-latest
# Wait for the instance to become ready
	sleep 1
	duration=0
	while (( duration < 240 )) && ! docker logs pingfederate_terraform_provider_container 2>&1 | grep -q "PingFederate is up"; \
	do \
	    duration=$$((duration+1)); \
		sleep 1; \
	done
# Fail if the container didn't become ready in time
# docker logs pingfederate_terraform_provider_container 2>&1 | grep -q "PingFederate is up" || \
# 	{ echo "PingFederate container did not become ready in time. Logs:"; docker logs pingfederate_terraform_provider_container; exit 1; }
	docker logs pingfederate_terraform_provider_container -f
		
removetestcontainer:
	docker rm -f pingfederate_terraform_provider_container
	
spincontainer: removetestcontainer starttestcontainer

define test_acc_env_vars
	PINGFEDERATE_PROVIDER_HTTPS_HOST=https://localhost:9999 PINGFEDERATE_PROVIDER_USERNAME=administrator PINGFEDERATE_PROVIDER_PASSWORD=2FederateM0re PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS=true PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER=true PINGFEDERATE_PROVIDER_PRODUCT_VERSION=$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-12.0}
endef

# Set ACC_TEST_NAME to name of test in cli
testoneacc:
	$(call test_acc_env_vars) TF_ACC=1 go test ./internal/acctest/... -timeout 10m -run ${ACC_TEST_NAME} -v count=1

testoneacccomplete: spincontainer testoneacc

# Some tests can step on each other's toes so run those tests in single threaded mode. Run the rest in parallel
testacc:
	$(call test_acc_env_vars) TF_ACC=1 go test `go list ./internal/acctest/... | grep -v -e authenticationapiapplication -e authenticationapisettings -e oauthauthserversettings -e oauthopenidconnectpolicy` -timeout 10m -v -p 4; \
	firstTestResult=$$?; \
	$(call test_acc_env_vars) TF_ACC=1 go test `go list ./internal/acctest/... | grep -e authenticationapiapplication -e authenticationapisettings -e oauthauthserversettings -e oauthopenidconnectpolicy` -timeout 10m -v -p 1; \
	secondTestResult=$$?; \
	if test "$$firstTestResult" != "0" || test "$$secondTestResult" != "0" ; then \
		false; \
	fi

testacccomplete: spincontainer testacc

clearstates:
	find . -name "*tfstate*" -delete
	
kaboom: clearstates spincontainer install

devchecknotest: install golangcilint generate tfproviderlint tflint terrafmtlint importfmtlint

verifycontent:
	python3 ./devcheck/verifyContent.py

devcheck: verifycontent devchecknotest kaboom testacc

generateresource:
	PINGFEDERATE_GENERATED_ENDPOINT=serverSettings \
	PINGFEDERATE_RESOURCE_DEFINITION_NAME=ServerSettings \
	PINGFEDERATE_ALLOW_REQUIRED_BYPASS=False \
	OVERWRITE_EXISTING_RESOURCE_FILE=False \
	PINGFEDERATE_PUT_ONLY_RESOURCE=True \
	GENERATE_SCHEMA=True \
	python3 scripts/generate_resource.py
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
