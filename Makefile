SHELL := /bin/bash

.PHONY: install generate fmt vet test starttestcontainer removetestcontainer spincontainer clearstates kaboom testacc testacccomplete generateresource openlocalwebapi golangcilint tfproviderlint tflint terrafmtlint importfmtlint devcheck devchecknotest openapp testoneacc verifycontent testupgradeacc testupgradeoneacc testupgradecomplete spinupgradelanes testserverupgradeacc testserverupgradeoneacc

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
 	PRODUCT_VERSION_DIR=$$(echo "$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-13.1.0}" | cut -b 1-4)
endef

starttestcontainer:
	$(call productversiondir) && docker run --name pingfederate_terraform_provider_container \
		-d -p 9031:9031 \
		-p 9999:9999 \
		--env-file "${HOME}/.pingidentity/config" \
		-e "OPERATIONAL_MODE=${OPERATIONAL_MODE}" \
		-v $$(pwd)/server-profiles/shared-profile:/opt/in \
		-v $$(pwd)/server-profiles/$${PRODUCT_VERSION_DIR}/data.json$${DATA_JSON_SUFFIX}.subst:/opt/in/instance/bulk-config/data.json.subst \
		pingidentity/pingfederate:$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-13.1.0}-latest
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
	PINGFEDERATE_PROVIDER_HTTPS_HOST=https://localhost:9999 PINGFEDERATE_PROVIDER_ADMIN_API_PATH="/pf-admin-api/v1" PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS=true PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER=true PINGFEDERATE_PROVIDER_PRODUCT_VERSION=$${PINGFEDERATE_PROVIDER_PRODUCT_VERSION:-13.1}
endef

define test_acc_basic_auth_env_vars
	PINGFEDERATE_PROVIDER_USERNAME=administrator PINGFEDERATE_PROVIDER_PASSWORD=2FederateM0re
endef

define test_acc_oauth_env_vars
	PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID=test PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET=2FederateM0re! PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL=https://localhost:9031/as/token.oauth2 PINGFEDERATE_PROVIDER_OAUTH_SCOPES=email
endef

# Set ACC_TEST_NAME to name of test in cli
testoneacc:
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test ./internal/acctest/config/${ACC_TEST_FOLDER}... -timeout 20m -run ${ACC_TEST_NAME} -v -count=1

testaccfolder:
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test ./internal/acctest/config/${ACC_TEST_FOLDER}... -timeout 20m -v -count=1

testoneacccomplete: spincontainer testoneacc

# Some tests can step on each other's toes so run those tests in single threaded mode. Run the rest in parallel
testacc:
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test `go list ./internal/acctest/config... | grep -v -e authenticationapi -e oauth/authserversettings -e oauth/openidconnect/policy -e oauth/openidconnect/settings -e oauth/clientsettings -e oauth/tokenexchange/processor -e serversettings/wstruststssettings -e sp/targeturlmappings -e serversettings/systemkeys/rotate -e oauth/cibaserverpolicy/requestpolicies -e notificationpublishers -e oauth/accesstokenmanagers/settings -e oauth/accesstokenmapping -e captchaproviders/settings` -timeout 20m -v -p 4; \
	firstTestResult=$$?; \
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test `go list ./internal/acctest/config... | grep -e authenticationapi -e oauth/authserversettings -e oauth/openidconnect/policy -e oauth/openidconnect/settings -e oauth/clientsettings -e oauth/tokenexchange/processor -e serversettings/wstruststssettings -e sp/targeturlmappings -e serversettings/systemkeys/rotate -e oauth/cibaserverpolicy/requestpolicies -e notificationpublishers -e oauth/accesstokenmanagers/settings -e oauth/accesstokenmapping -e captchaproviders/settings` -timeout 20m -v -p 1; \
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

testaccclustered:
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 go test ./internal/acctest/config/cluster/... -timeout 5m -v

testacccomplete: spincontainer testacc

# Provider version-ladder tests: for each resource whose gen test file calls
# RunUpgradeLadder (inlined in internal/acctest/config/.../*_gen_test.go), apply
# its config with the oldest provider version that supports the current
# PINGFEDERATE_PROVIDER_PRODUCT_VERSION lane, then hop the provider one minor
# version at a time to the local build, asserting an empty plan after every
# hop. Compiled only under the 'upgradeladder' build tag (without it the tests
# skip instantly). Requires registry access (rungs download from the Terraform
# Registry); TF_PLUGIN_CACHE_DIR dedupes those downloads across rungs and
# resources.
#
# TF_CLI_CONFIG_FILE points at a minimal isolated CLI config so a developer's
# ~/.terraformrc dev_overrides cannot silently replace every rung's pinned
# registry version with the local dev binary.
# Ladder rungs are collision-prone singletons, so the tests run serially.
# Set ACC_TEST_NAME=<resource type substring> to run one resource's test.
# Set PINGFEDERATE_UPGRADE_LADDER=full|last2|<csv rungs> to override the ladder.
define upgrade_ladder_tfrc
	mkdir -p $(CURDIR)/.tfplugincache && \
	printf 'provider_installation {\n  direct {}\n}\n' > $(CURDIR)/.tfplugincache/tfrc
endef

testupgradeacc:
	$(call upgrade_ladder_tfrc)
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 TF_PLUGIN_CACHE_DIR=$(CURDIR)/.tfplugincache TF_CLI_CONFIG_FILE=$(CURDIR)/.tfplugincache/tfrc \
		go test -tags upgradeladder ./internal/acctest/config/... -run 'TestUpgradeLadder_.*' -timeout 60m -v -p 1 -count=1

testupgradeoneacc:
	$(call upgrade_ladder_tfrc)
	$(call test_acc_common_env_vars) $(call test_acc_basic_auth_env_vars) TF_ACC=1 TF_PLUGIN_CACHE_DIR=$(CURDIR)/.tfplugincache TF_CLI_CONFIG_FILE=$(CURDIR)/.tfplugincache/tfrc \
		go test -tags upgradeladder ./internal/acctest/config/... -run 'TestUpgradeLadder_.*$(ACC_TEST_NAME).*' -timeout 60m -v -p 1 -count=1

testupgradecomplete: spincontainer testupgradeacc

# PingFederate server-upgrade ladder: one container per supported lane on
# distinct ports; each ladder step re-points the provider (via its step
# provider block) at the next lane's server while the Terraform state carries
# forward, so the tests exercise resources applied on an older server and
# re-planned/re-applied after a server upgrade. LANES is a comma-separated
# major.minor list (default: every supported lane); port base for lane N is
# 9999 + N*100 (lane 1 -> 9999, lane 2 -> 10099, ...). The first lane's host
# feeds PINGFEDERATE_PROVIDER_HTTPS_HOST; later lanes feed
# PINGFEDERATE_UPGRADE_LANE_HOSTS.
# Set ACC_TEST_NAME to run one resource's test.
# The bulk config is a cleaned copy (CI-secret ops dropped) of that lane's own
# server-profiles/<verdir>/data.json.subst, built once per verdir into
# /tmp/pf-env/bulk-<verdir>/data.json.subst if absent.
define spin_server_lane
	docker rm -f pf-ladder-lane-$(1) 2>/dev/null; \
	docker run --name pf-ladder-lane-$(1) -d \
		-p $$(($(1)*100+31)):9031 -p $(2):9999 \
		--env-file /tmp/pf-env/config-noexport \
		-e "OPERATIONAL_MODE=STANDALONE" \
		-v $$(pwd)/server-profiles/shared-profile:/opt/in \
		-v /tmp/pf-env/bulk-$(3)/data.json.subst:/opt/in/instance/bulk-config/data.json.subst \
		pingidentity/pingfederate:$(4)-latest && \
	duration=0; \
	while (( duration < 240 )) && ! docker logs pf-ladder-lane-$(1) 2>&1 | grep -q "Removing Imported Bulk File\|CONTAINER FAILURE"; \
	do duration=$$((duration+1)); sleep 1; done; \
	docker logs pf-ladder-lane-$(1) 2>&1 | grep -q "Removing Imported Bulk File" || \
		{ echo "Lane container $(4) did not become ready in time"; docker logs pf-ladder-lane-$(1) | tail -20; exit 1; }
endef

spinupgradelanes:
	docker rm -f pingfederate_terraform_provider_container 2>/dev/null; \
	mkdir -p /tmp/pf-env && grep -vE "^export " "${HOME}/.pingidentity/config" > /tmp/pf-env/config-noexport; \
	lane_i=1; lanes=""; hosts=""; \
	for ver in $(shell echo $(LANES) | tr ',' ' '); do \
		port=$$((9999 + lane_i * 100)); \
		verdir=$$(echo $$ver | cut -b 1-4); \
		mkdir -p /tmp/pf-env/bulk-$${verdir}; \
		if [ ! -s /tmp/pf-env/bulk-$${verdir}/data.json.subst ]; then \
			python3 -c "import json; d=json.load(open('server-profiles/$${verdir}/data.json.subst')); ops=[op for op in d['operations'] if op.get('resourceType') not in ('/pingOneConnections','/oauth/outOfBandAuthPlugins')]; json.dump({'metadata': d.get('metadata',{}), 'operations': ops}, open('/tmp/pf-env/bulk-$${verdir}/data.json.subst','w'))"; \
		fi; \
		$(call spin_server_lane,$${lane_i},$${port},$${verdir},$${ver}); \
		lanes="$$lanes,$${ver%.*}"; hosts="$$hosts,https://localhost:$${port}"; \
		lane_i=$$((lane_i+1)); \
	done; \
	lanes=$${lanes#,}; hosts=$${hosts#,}; rest=$${hosts#*,}; \
	echo "Export these for testserverupgradeacc:"; \
	echo "  export PINGFEDERATE_UPGRADE_SERVER_LANES='$$lanes'"; \
	echo "  export PINGFEDERATE_UPGRADE_LANE_HOSTS='$$rest'"; \
	echo "  export PINGFEDERATE_PROVIDER_HTTPS_HOST='$${hosts%%,*}'"; \
	echo "  export PINGFEDERATE_PROVIDER_PRODUCT_VERSION=$${lanes%%,*}"; \
	echo "  export PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1=... PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_2=...  # optional: enables the certificate resources"

# The shared test functions run both ladders: this target sets the server-ladder
# opt-in (EnvServerLadderEnabled) and the same registry isolation as the
# provider-ladder targets, since RunUpgradeLadder runs first in every test.
testserverupgradeacc:
	$(call upgrade_ladder_tfrc)
	$(call test_acc_basic_auth_env_vars) TF_ACC=1 \
		PINGFEDERATE_UPGRADE_SERVER_LADDER=1 \
		PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS=true PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER=true PINGFEDERATE_PROVIDER_ADMIN_API_PATH="/pf-admin-api/v1" \
		TF_PLUGIN_CACHE_DIR=$(CURDIR)/.tfplugincache TF_CLI_CONFIG_FILE=$(CURDIR)/.tfplugincache/tfrc \
		go test -tags upgradeladder ./internal/acctest/config/... -run 'TestUpgradeLadder_.*' -timeout 60m -v -p 1 -count=1

testserverupgradeoneacc:
	$(call upgrade_ladder_tfrc)
	$(call test_acc_basic_auth_env_vars) TF_ACC=1 \
		PINGFEDERATE_UPGRADE_SERVER_LADDER=1 \
		PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS=true PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER=true PINGFEDERATE_PROVIDER_ADMIN_API_PATH="/pf-admin-api/v1" \
		TF_PLUGIN_CACHE_DIR=$(CURDIR)/.tfplugincache TF_CLI_CONFIG_FILE=$(CURDIR)/.tfplugincache/tfrc \
		go test -tags upgradeladder ./internal/acctest/config/... -run 'TestUpgradeLadder_.*$(ACC_TEST_NAME).*' -timeout 60m -v -p 1 -count=1

clearstates:
	find . -name "*tfstate*" -delete
	
kaboom: clearstates spincontainer install

devchecknotest: verifycontent install golangcilint generate tfproviderlint tflint terrafmtlint importfmtlint

verifycontent:
	python3 ./scripts/verifyContent.py

devcheck: devchecknotest kaboom testacc

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
	go tool golangci-lint run --timeout 5m ./internal/...

tfproviderlint:
	go tool tfproviderlintx \
						-c 1 \
						-AT001.ignored-filename-suffixes=_test.go \
						-AT003=false \
						-R018=false \
						-XAT001=false \
						-XR004=false \
						-XS002=false ./internal/...

tflint:
	go tool tflint --recursive --disable-rule "terraform_unused_declarations" --disable-rule "terraform_required_providers" --disable-rule "terraform_required_version"

terrafmtlint:
	find ./internal/acctest -type f -name '*_test.go' \
		| sort -u \
		| xargs -I {} go tool terrafmt -f fmt {} -v

importfmtlint:
	go tool impi --local . --scheme stdThirdPartyLocal ./internal/...
