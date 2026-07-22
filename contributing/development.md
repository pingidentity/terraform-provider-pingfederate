# Development Environment Setup

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.4+ (to run acceptance tests)
- [Go](https://golang.org/doc/install) 1.23.3+ (to build and test the provider plugin)
- [Docker](https://docs.docker.com/get-docker/) (to run acceptance tests and examples)

## Quick Start

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](#requirements) before proceeding).

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/pingidentity/`).

Clone repository to: `$HOME/development/pingidentity/`

```sh
mkdir -p $HOME/development/pingidentity/; cd $HOME/development/pingidentity/
git clone git@github.com:pingidentity/terraform-provider-pingfederate.git
```

To compile the provider, run `make generate`. This will generate, format, and vet the code.

```sh
make generate
```

To install the provider for local use, run `make install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make install
```

## File system overview

Before starting, take a moment to familiarize yourself with the structure of this repository, [found here](filelayout.md)

## Preparing your Terraform environment to run locally-built providers

By default, Terraform attempts to pull providers from remote registries. This behavior can be overwritten by modifying the `~/.terraformrc` file to enable the local use of this provider. This configuration can be used to test in-development changes to the provider.

With Terraform v0.14 and later, [development overrides for provider developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers) can be leveraged in order to use the provider built from source.

First, find the **GOBIN** path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.

```sh
$ go env GOBIN
/Users/<Username>/go/bin
```

If the GOBIN go environment variable is not set, use the default path, **/Users/\<Username\>/go/bin**. 

To do this, populate a Terraform CLI configuration file (`~/.terraformrc` for all platforms other than Windows; `terraform.rc` in the `%APPDATA%` directory when using Windows) with at least the following options, changing the \<PATH\> value to the value returned from `go env GOBIN`:

```hcl
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

## Local SDK Changes

### Running a locally-built PingFederate Go client

The PingFederate Terraform provider relies on the [PingFederate Go Client](https://github.com/pingidentity/pingfederate-go-client).

Occasionally, development may include changes to the PingFederate Go SDK. If you'd like to develop this provider locally using a local, modified version of the SDK, this can be achieved by adding a `replace` directive in the `go.mod` file. 

If changes are needed in the Go client, the `replace` command in the `go.mod` file can be used to point to a modified local Go client while testing.

For example, the start of the `go.mod` file may look like the following example, where the local cloned SDK is in the `../pingfederate-go-client` relative path:

```
module github.com/pingidentity/terraform-provider-pingfederate

go 1.23.3

replace github.com/pingidentity/pingfederate-go-client v0.0.1 => ../pingfederate-go-client

require (
	github.com/pingidentity/pingfederate-go-client v0.x.x
  
  ...
)

...
```

In this example above, the `replace` path points to `../pingfederate-go-client`, meaning you would need to clone the client repo and place it alongside this repo in your filesystem.

Once updated, run the following to build the project:

```shell
make generate
```

## Install the provider

Run `go mod tidy` to get any required dependencies.

Run `make install` (or just `make`) to install the provider locally.

## Development Commands

The following make targets are available for development:

**Code Generation and Formatting:**
- `make generate`: Generates code, formats, and vets the code
- `make fmt`: Formats Go code
- `make vet`: Runs Go vet

**Development Validation:**
- `make devcheck`: Full development check (includes linting, unit tests, and acceptance tests)
- `make devchecknotest`: Development check without running tests
- `make test`: Run unit tests only (fast, no external dependencies)
- `make verifycontent`: Verifies content using Python scripts

**Linting:**
- `make golangcilint`: Runs golangci-lint
- `make tfproviderlint`: Runs Terraform provider linting
- `make tflint`: Runs Terraform linting
- `make terrafmtlint`: Runs Terraform format linting
- `make importfmtlint`: Runs import format linting
- `make lint`: Runs all linting checks above

**Utility:**
- `make generateresource`: Generates a new resource
- `make openlocalwebapi`: Opens the local web API documentation
- `make openapp`: Opens the PingFederate application

## Testing the Provider

In order to test the provider locally, you can run different types of tests:

**Unit Tests**: Run fast unit tests without external dependencies:
```sh
make test
```

**Code Validation**: Run basic validation and vetting:
```sh
make vet
```

## Running acceptance tests

In order to run the full suite of Acceptance tests against a live PingFederate tenant, run `make testacc`.

*Note:* Acceptance tests create real configuration in PingFederate. Please ensure you have a trial PingFederate account or licensed subscription to proceed with these tests.

Acceptance tests for the provider use a local PingFederate instance running in Docker. The following `make` targets will help with running acceptance tests:

- `make testacc`: Runs the acceptance tests, with the assumption that a local PingFederate instance is available
- `make starttestcontainer`: Starts a PingFederate Docker container and waits for it to become ready
- `make removetestcontainer`: Stops and removes the PingFederate Docker container used for testing
- `make testacccomplete`: Starts the PingFederate Docker container, waits for it to become ready, and runs the acceptance tests. This option is good for running the tests from scratch and for use in automation, but you will have to wait for the container startup each time
- `make testoneacc`: Run a single acceptance test (requires setting `ACC_TEST_NAME` and `ACC_TEST_FOLDER` environment variables)
- `make spincontainer`: Shortcut for `removetestcontainer` then `starttestcontainer`

**Additional development targets:**
- `make devcheck`: Runs the full development check including all linting and tests
- `make devchecknotest`: Runs development checks without running tests
- `make golangcilint`: Runs Go linting
- `make verifycontent`: Verifies content using Python scripts
- `make clearstates`: Clears all Terraform state files
- `make kaboom`: Full reset - clears states, spins container, and installs
  
**Tip**: If you plan on running tests multiple times and do not mind reusing the same server, then it is recommended to use the first three options above to perform each step individually.

```sh
make testacc
```

## Using the Provider

The development overrides configuration described in the [Preparing your Terraform environment](#preparing-your-terraform-environment-to-run-locally-built-providers) section allows you to use the provider built from source in your Terraform configurations.

## Run an example

### Start the PingFederate server

Start a PingFederate server running locally with the provided **docker-compose.yaml** file. Change to the `docker-compose` directory and run `docker compose up`. (Alternatively, use the `make starttestcontainer` command from the previous section.) The server will take a couple of minutes to become ready. When you see the following output in the terminal, the server is ready to process requests:

```
pingfederate-1  | PingFederate is up
```

### Run Terraform

Change to the `examples/resources/<desired resource>` directory. The `resource.tf` file in this directory defines the Terraform configuration.

Run `terraform plan` to view what changes will be made by Terraform. Run `terraform apply` to apply them.

You can verify the configuration is created via administrator API (https://hostname/pf-admin-api/api-docs) or the UI.

You can make changes to the configuration and use `terraform apply` to apply them, and use the above commands to view those changes in PingFederate.

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