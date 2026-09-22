# Development Environment Setup
## Prerequisites
The following applications must be installed locally to run the provider:
- [Terraform](https://www.terraform.io/downloads.html) 1.0+ (to run acceptance tests)
- [Go](https://golang.org/doc/install) 1.20+ (to build the provider plugin)

To run the example in this repository, you will also need
- [Docker](https://docs.docker.com/get-docker/)

## File system overview

Before starting, take a moment to familiarize yourself with the structure of this repository, [found here](filelayout.md)

## Preparing your Terraform environment to run locally-built providers
By default, Terraform attempts to pull providers from remote registries. This behavior can be overwritten by modifying the `~/.terraformrc` file to enable the local use of this provider. This configuration can be used to test in-development changes to the provider.

First, find the **GOBIN** path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.

```sh
$ go env GOBIN
/Users/<Username>/go/bin
```

If the GOBIN go environment variable is not set, use the default path, **/Users/\<Username\>/go/bin**. Create a `~/.terraformrc` file with the following contents, changing the \<PATH\> value to the value returned from `go env GOBIN`.

```text
provider_installation {
  dev_overrides {
    "pingidentity/pingfederate" = "<PATH>"
  }
  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this next line, Terraform will _only_ use
  # the dev_overrides block, meaning no other providers will be available.
  direct {}
}
```

### Running a locally-built PingFederate Go client
The PingFederate Terraform provider relies on the [PingFederate Go Client](https://github.com/pingidentity/pingfederate-go-client).

If changes are needed in the Go client, the `replace` command in the `go.mod` file can be used to point to a modified local Go client while testing.

```txt
replace github.com/pingidentity/pingfederate-go-client v0.0.1 => ./pinfederate-go-client
```

In this example above, the `replace` path points to `../pinfederate-go-client`, meaning you would need to clone the client repo and place it alongside this repo in your filesystem.

## Install the provider
Run `go mod tidy` to get any required dependencies.

Run `make install` (or just `make`) to install the provider locally.

## Running acceptance tests
Acceptance tests for the provider use a local PingFederate instance running in Docker. The following `make` targets will help with running acceptance tests:

- `make testacc`: Runs the acceptance tests, with the assumption that a local PingFederate instance is available
- `make starttestcontainer`: Starts a PingFederate Docker container and waits for it to become ready
- `make removetestcontainer`: Stops and removes the PingFederate Docker container used for testing
- `make testacccomplete`: Starts the PingFederate Docker container, waits for it to become ready, and runs the acceptance tests. This option is good for running the tests from scratch and for use in automation, but you will have to wait for the container startup each time.
  
**Tip**: If you plan on running tests multiple times and do not mind reusing the same server, then it is recommended to use the first three options above to perform each step individually.

## Provider version-ladder tests

"Provider produced inconsistent result after apply" errors (see CDI-532 and CDI-533) surface when the **provider version changes** against a live PingFederate server that injects defaults for unset optional properties. The version-ladder tests in `internal/acctest/upgradeladder` catch that bug class: for each participating resource, the same configuration is applied with the oldest provider version that supports the current PingFederate server lane, then the provider is stepped one minor version at a time up to the local build, asserting an empty plan after every hop.

Each PingFederate server version defines a **lane** (`PINGFEDERATE_PROVIDER_PRODUCT_VERSION` major.minor), and the ladder starts at the first provider release that supports the lane:

| PingFederate lane | Ladder rungs (provider versions) |
|---|---|
| 12.2 | 1.3.0 → 1.4.5 → 1.5.0 → 1.6.2 → 1.7.1 → 1.8.1 → 1.9.0 → 1.10.0 → local |
| 12.3 | 1.6.2 → 1.7.1 → 1.8.1 → 1.9.0 → 1.10.0 → local |
| 13.0 | 1.7.1 → 1.8.1 → 1.9.0 → 1.10.0 → local |
| 13.1 | 1.9.0 → 1.10.0 → local |

Rungs come from the public Terraform Registry (`pingidentity/pingfederate`, pinned with exact versions); the final rung is the locally built provider. Registry rungs require network access — downloads are cached in `.tfplugincache/` (gitignored).

### Running the ladder tests

```sh
# Start a container for the desired lane (12.2 exercises the full ladder)
PINGFEDERATE_PROVIDER_PRODUCT_VERSION=12.2 make spincontainer
# Run every participating resource through its ladder
PINGFEDERATE_PROVIDER_PRODUCT_VERSION=12.2 make testupgradeacc
# Run one resource's ladder test
make testupgradeoneacc ACC_TEST_NAME=oauth_server_settings
# Shortest smoke: final registry rung + local only
PINGFEDERATE_UPGRADE_LADDER=last2 make testupgradeacc
# From scratch (container + ladder)
make testupgradecomplete
```

The targets write an isolated Terraform CLI config (`.tfplugincache/tfrc`) so a developer's `~/.terraformrc` `dev_overrides` cannot silently substitute the local dev binary for the pinned registry rungs. `make testupgradeacc` runs with `-p 1` (participating resources include singletons that must not run concurrently with each other, or with `make testacc`, against the same container). `make testacc` never live-runs these tests: without the `upgradeladder` build tag they skip before touching the server.

In CI, the ladder runs only in the nightly [Scheduled Acceptance Tests workflow](../.github/workflows/scheduled-acctests.yaml) (and on manual dispatch) — never on pull requests, where it is too slow and registry-dependent. The scheduled job runs the full 12.2 lane plus the newest lane.

### Environment knobs

- `PINGFEDERATE_UPGRADE_LADDER`: `full` (default), `last2` (final registry rung + local), or explicit comma-separated rungs, e.g. `1.9.0,local`.
- `PINGFEDERATE_UPGRADE_RESOURCES`: comma-separated substrings selecting participating resource types, e.g. `oauth_server_settings,incoming_proxy_settings`. Non-matching resources skip.
- `ACC_TEST_NAME`: test-name filter for `make testupgradeoneacc`.

### Adding a resource to the ladder

Add a `TestUpgradeLadder_<ResourceName>` test to the resource's existing `*_gen_test.go` file, declaring one `upgradeladder.Spec` and calling `upgradeladder.RunUpgradeLadder`. There is no central registry and no extra test file: the test lives next to the resource's regular tests and reuses the package's generated `*_MinimalHCL()` (or an inline func), so the ladder always exercises the same HCL shape the package's regular tests use — no frozen HCL copy to keep in sync.

```go
func TestUpgradeLadder_SessionSettings(t *testing.T) {
	upgradeladder.RunUpgradeLadder(t, upgradeladder.Spec{
		ResourceType: "pingfederate_session_settings",
		HCL:          sessionSettings_MinimalHCL,
	})
}
```

Import the harness in the gen test file (as `upgradeladder "github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/upgradeladder"`). The gen test files carry no build tag: ladder tests compile in every build but call `RunUpgradeLadder`, which skips instantly unless the build sets the `upgradeladder` tag (i.e. `make testupgradeacc`) — so plain `go test`/`make testacc` runs see only a zero-cost skip, never a live run.

Rules:

- The Spec's HCL must apply cleanly with the *oldest* rung of every lane it runs on, so do not reference attributes added after that release. Resources shipped after v1.3.0 set `AvailableSince` to their first release (older rungs are dropped — e.g. `pingfederate_incoming_proxy_settings` is `"1.4.5"`).
- Top-level `data "pingfederate_<resourceType>" "example"` blocks (present in some generated HCL) are stripped automatically by the harness; other data sources are left in and fail the run loudly.
- The HCL must use the resource label `example` (enforced by `ValidateSpec`).
- `ClusterModeProbe` optionally gates the ladder on a live-server condition and skips when it fails — e.g. cluster settings on a standalone server, or a required CI-secret environment variable like `PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1` for the certificate resources.
- `Allowlist` is reserved for a future phase that tolerates attribute-level plan deltas (e.g. server-injected defaults); entries require a JIRA reference in `Reason`.

The table in `ladder.go` is maintained alongside releases: append a row when a new provider minor ships, and drop rows when the provider drops an EOL PingFederate lane. Everything else about PingFederate versions derives from `internal/version/version.go` — the supported-lane list (`version.SupportedMajorMinorVersions()`), lane validation (`version.Parse` in `ParseLane`, the same contract as the provider's own configure-time check), and the table's lane references (`MaxPFMinor` references `version.PingFederate*` constants). Supporting a new PingFederate version is a one-line change in `internal/version/version.go` plus the new `ladderTable` rows — `TestLadderTableCoversSupportedLanes` fails until the rows are added, and any stale row for a dropped lane fails compilation. Do not hand-write PingFederate version strings in ladder code or tests; reference the constants (or the helpers `version.Parse`/`version.MajorMinor`/`version.SupportedMajorMinorVersions`). Spec invariants are checked by `go test ./internal/acctest/upgradeladder/` without a server.

## Run an example
### Start the PingFederate server
Start a PingFederate server running locally with the provided **docker-compose.yaml** file. Change to the `docker-compose` directory and run `docker compose up`. (Alternatively, use the `make starttestcontainer` command from the previous section.) The server will take a couple of minutes to become ready. When you see the following output in the terminal, the server is ready to process requests:
```
pingfederate-1  | PingFederate is up
```

### Run Terraform
Change to the `examples/resources/<desired resource>` directory. The `resource.tf` file in this directory defines the Terraform configuration.

Run `terraform plan` to view what changes will be made by Terraform. Run `terraform apply` to apply them.

You can verify the location is created via administrator API(https://hostname/pf-admin-api/api-docs) or the UI:

You can make changes to the location and use `terraform apply` to apply them, and use the above commands to view those changes in PingFederate.

Run `terraform destroy` to destroy any objects managed by Terraform.

## Debugging with VSCode
You can attach a debugger to the provider with VSCode. The `.vscode/launch.json` file defines the debug configuration.

To debug the provider, navigate to **Run > Start Debugging**. Then, open the Debug Console and wait for a message like this:

```text
Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

	TF_REATTACH_PROVIDERS='{"registry.terraform.io/pingidentity/pingfederate":{"Protocol":"grpc","ProtocolVersion":6,"Pid":94877,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/h0/myp0srpd29d7mr9_0rgvjwtm0000gn/T/plugin2376654838"}}}'
```

You can then use this to attach the debugger to command-line terraform commands by pasting this line before each command.

```sh
$ TF_REATTACH_PROVIDERS='{"registry.terraform.io/pingidentity/pingfederate":{"Protocol":"grpc","ProtocolVersion":6,"Pid":94877,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/h0/myp0srpd29d7mr9_0rgvjwtm0000gn/T/plugin2376654838"}}}' terraform apply
```

**Note**: The `TF_REATTACH_PROVIDERS` variable changes each time you run the debugger. You will need to copy the output for use each time you start a new debugger.