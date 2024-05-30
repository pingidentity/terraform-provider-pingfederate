# Instructions on how to using custom schema in docs using templates

### Steps
1. Run `make generate` to generate the full document using the `tfplugindocs` package.
2. Copy the lines under the `###Schema` section until above the `###Import` from newly generated file and paste into a new markdown document in the `./schema-docs` directory - be sure to match the resource name as the file name. Make any changes you wish to add key notes about recursion, etc. at the top of this file. These will get injected into the templates found in `./templates/resources/*`
3. Remove the newly generated documentation from step #1 (leave all other existing documentation) in `./docs/resources` (or `data-sources` directory).
4. The templates are found in the `./templates/*` directory. You will need to create a new one to match the resource you desire to generate - use the others here as reference
5. Add the `<resource_type>.md` to the list of documents on line #21 in the `main.go` file in the root of this repo
6. Run `make generate`. **Keep in mind any file you put in the `./templates` directory (that is not a template) will be added to the `./docs/` directory after `make generate` is run**
7. To verify your generated doc, Hashicorp provides a nice registry documentation preview tool, found [here](https://registry.terraform.io/tools/doc-preview). This is extremely helpful for verifying formatting.

### Notes
* The `authentication_fragments` and `authentication_policies` resources use initial generated schema content from `make generate`. If you wish to modify the schema here, it's easiest to change the `const MaxPolicyNodeRecursiveDepth = 10` to have a value of `0` in `authentication_policy_tree_node_schema.go`

* The `./scripts/markdownDocFormatting.py` script is needed to remove the hard-coded backticks that are in the `tfplugindocs` package for using the `codefile` argument. As this point in time, this is the best way to accomplish desired schema formatting for the registry.
