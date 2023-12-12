---
name: 🐛 Bug Report
about: Create a report to help us improve
title: ''
labels: bug
assignees: ''

---

<!--- Please keep this note for the community --->

### Community Note

* Please vote on this issue by adding a 👍 [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original issue to help the community and maintainers prioritize this request
* Please do not leave "+1" or other comments that do not add relevant new information or questions, they generate extra noise for issue followers and do not help prioritize the request
* If you are interested in working on this issue or have submitted a pull request, please leave a comment

<!--- Thank you for keeping this note for the community --->

Thank you for opening an issue. Please note that we try to keep the Terraform issue tracker reserved for bug reports and feature requests. For general usage questions, please see: https://www.terraform.io/community.html.

### PingFederate Terraform provider Version
Check the version you have configured in your .tf files. If you are not running the latest version of the provider, please upgrade because your issue may have already been fixed.

### PingFederate Version
Check the version of the PingFederate server being configured.

### Terraform Version
Run `terraform -v` to show the version.

### Affected Resource(s)
Please list the resources as a list, for example:
- pingfederate_administrative_accounts
- pingfederate_oauth_issuers

### Terraform Configuration Files
```hcl
# Copy-paste your Terraform configurations here - for large Terraform configs,
# please use a service like Dropbox and share a link to the ZIP file. For
# security, you can also encrypt the files using our GPG public key.

# Remember to replace any account/customer sensitive information in the configuration before submitting the issue
```

### Debug Output
Please provide your debug output with `TF_LOG=DEBUG` enabled on your `terraform plan` or `terraform apply`. Please provide a link to a GitHub Gist containing the complete debug output: https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

### Panic Output
If Terraform produced a panic, please provide your debug output from the GO panic

### Expected Behavior
What should have happened?

### Actual Behavior
What actually happened?

### Steps to Reproduce
Please list the steps required to reproduce the issue, for example:
1. `terraform apply`

### References
Are there any other GitHub issues (open or closed) or Pull Requests that should be linked here? For example:
- GH-1234
