---
page_title: "PingFederate Server Support"
description: |-
  Details of the PingFederate server version support in the provider.
---

# PingFederate Server Support

The support of PingFederate product versions aligns with the [Ping Identity End of Life Policy](https://www.pingidentity.com/en/legal/end-of-life-policy.html).  Once a PingFederate server version becomes end of life, future releases of the provider will no longer be tested against that version and it's support in the provider will be removed.

If issues are encountered when using the Terraform provider, customers are encouraged to first:
1. Ensure that the PingFederate server version is supported according to the [Ping Identity End of Life Policy](https://www.pingidentity.com/en/legal/end-of-life-policy.html).
2. Ensure that the PingFederate server is compatible with the Terraform provider version according to the table below.

The following table lists the minimum and maximum versions of PingFederate server that are compatible with previously released Terraform provider versions.

| Provider Version | Minimum PingFederate Version | Maximum PingFederate Version |
|------------------|------------------------------|------------------------------|
| `1.0`            | `11.2`                       | `12.1`                       |
